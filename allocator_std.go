//go:build !tinygo && !wasm && !wasi && !tinygo.wasm
// +build !tinygo,!wasm,!wasi,!tinygo.wasm

package mem

import (
	"reflect"
	"sync"
	"unsafe"
)

var (
	allocs   = make(map[uintptr][]byte)
	allocsMu sync.Mutex
)

func malloc(size uintptr) unsafe.Pointer {
	b := make([]byte, size)
	p := unsafe.Pointer(&b[0])
	allocsMu.Lock()
	defer allocsMu.Unlock()
	allocs[uintptr(p)] = b
	return p
}

func realloc(ptr unsafe.Pointer, size uintptr) unsafe.Pointer {
	allocsMu.Lock()
	old := allocs[uintptr(ptr)]
	allocsMu.Unlock()
	if old == nil {
		return nil
	}

	b := make([]byte, size)
	copy(b, old)
	np := unsafe.Pointer(&b[0])
	allocsMu.Lock()
	allocs[uintptr(np)] = b
	allocsMu.Unlock()
	return np
}

func free(ptr unsafe.Pointer) {
	allocsMu.Lock()
	defer allocsMu.Unlock()
	b := allocs[uintptr(ptr)]
	if b != nil {
		delete(allocs, uintptr(ptr))
	}
}

func memcpy(dst, src unsafe.Pointer, n uintptr) {
	dstB := *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(dst),
		Len:  int(n),
		Cap:  int(n),
	}))
	srcB := *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(src),
		Len:  int(n),
		Cap:  int(n),
	}))
	copy(dstB, srcB)
}

func memzero(ptr unsafe.Pointer, size uintptr) {
	b := *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(ptr),
		Len:  int(size),
		Cap:  int(size),
	}))
	switch {
	case size < 8:
		for i := 0; i < len(b); i++ {
			b[i] = 0
		}
	case size == 8:
		*(*uint64)(unsafe.Pointer(&b[0])) = 0
	default:
		var i = 0
		for ; i <= len(b)-8; i += 8 {
			*(*uint64)(unsafe.Pointer(&b[i])) = 0
		}
		for ; i < len(b); i++ {
			b[i] = 0
		}
	}
}

func DefaultMalloc(size uintptr) unsafe.Pointer {
	return malloc(size)
}

func GrowByDouble(malloc Malloc) Grow {
	return func(pagesBefore, pagesNeeded int32, minSize uintptr) (pagesAdded int32, start, end uintptr) {
		if pagesBefore > pagesNeeded {
			pagesAdded = pagesBefore
		} else {
			pagesAdded = pagesNeeded
		}
		ptr := malloc(uintptr(pagesAdded) * tlsf_PAGE_SIZE)
		if ptr == nil {
			pagesAdded = pagesNeeded
			ptr = malloc(uintptr(pagesAdded) * tlsf_PAGE_SIZE)
			if ptr == nil {
				return 0, 0, 0
			}
		}
		start = uintptr(ptr)
		end = start + (uintptr(pagesAdded) * tlsf_PAGE_SIZE)
		return
	}
}

func GrowBy(pages int32, malloc Malloc) Grow {
	return func(pagesBefore, pagesNeeded int32, minSize uintptr) (pagesAdded int32, start, end uintptr) {
		if pages > pagesNeeded {
			pagesAdded = pages
		} else {
			pagesAdded = pagesNeeded
		}
		ptr := malloc(uintptr(pagesAdded) * tlsf_PAGE_SIZE)
		if ptr == nil {
			pagesAdded = pagesNeeded
			ptr = malloc(uintptr(pagesAdded) * tlsf_PAGE_SIZE)
			if ptr == nil {
				return 0, 0, 0
			}
		}
		start = uintptr(ptr)
		end = start + (uintptr(pagesAdded) * tlsf_PAGE_SIZE)
		return
	}
}

func GrowMin(malloc Malloc) Grow {
	return func(pagesBefore, pagesNeeded int32, minSize uintptr) (int32, uintptr, uintptr) {
		ptr := malloc(uintptr(pagesNeeded) * tlsf_PAGE_SIZE)
		if ptr == nil {
			return 0, 0, 0
		}
		return pagesNeeded, uintptr(ptr), uintptr(ptr) + (uintptr(pagesNeeded) * tlsf_PAGE_SIZE)
	}
}

func NewAllocator(pages int32) *Allocator {
	return NewAllocatorWithGrow(pages, GrowMin(DefaultMalloc))
}

func NewAllocatorWithGrow(pages int32, grow Grow) *Allocator {
	if pages <= 0 {
		pages = 1
	}
	pagesAdded, start, end := grow(0, pages, 0)
	return Bootstrap(start, end, pagesAdded, grow)
}

type SyncAllocator struct {
	*Allocator
	mu sync.Mutex
}

func (a *Allocator) ToSync() *SyncAllocator {
	return &SyncAllocator{Allocator: a}
}

// Alloc allocates a block of memory that fits the size provided
//goland:noinspection GoVetUnsafePointer
func (a *SyncAllocator) Alloc(size uintptr) unsafe.Pointer {
	a.mu.Lock()
	defer a.mu.Unlock()
	p := uintptr(unsafe.Pointer(a.allocateBlock(size)))
	if p == 0 {
		return nil
	}
	return unsafe.Pointer(p + tlsf_BLOCK_OVERHEAD)
}

// Realloc determines the best way to resize an allocation.
func (a *SyncAllocator) Realloc(ptr unsafe.Pointer, size uintptr) unsafe.Pointer {
	a.mu.Lock()
	defer a.mu.Unlock()
	return unsafe.Pointer(uintptr(unsafe.Pointer(a.moveBlock(checkUsedBlock(uintptr(ptr)), size))) + tlsf_BLOCK_OVERHEAD)
}

// Free release the allocation back into the free list.
func (a *SyncAllocator) Free(ptr unsafe.Pointer) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.freeBlock(checkUsedBlock(uintptr(ptr)))
}
