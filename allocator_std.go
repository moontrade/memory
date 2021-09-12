//go:build !tinygo && !wasm && !wasi && !tinygo.wasm && (darwin || linux)
// +build !tinygo
// +build !wasm
// +build !wasi
// +build !tinygo.wasm
// +build darwin linux

package mem

import (
	"reflect"
	"sync"
	"unsafe"
)

var (
	allocs   = make(map[uintptr][]byte)
	arenasMu sync.Mutex
)

//func malloc(size uintptr) unsafe.Pointer {
//	b := make([]byte, size)
//	p := unsafe.Pointer(&b[0])
//	arenasMu.Lock()
//	defer arenasMu.Unlock()
//	allocs[uintptr(p)] = b
//	return p
//}

func memequal(a, b unsafe.Pointer, n uintptr) bool {
	if a == nil {
		return b == nil
	}
	return *(*string)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(a),
		Len:  int(n),
	})) == *(*string)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(b),
		Len:  int(n),
	}))
}

func memcmp(a, b unsafe.Pointer, n uintptr) int {
	if a == nil {
		if b == nil {
			return 0
		}
		return -1
	}
	ab := *(*string)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(a),
		Len:  int(n),
	}))
	bb := *(*string)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(b),
		Len:  int(n),
	}))
	if ab < bb {
		return -1
	}
	if ab == bb {
		return 0
	}
	return 1
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

type GrowFactory func(arena Arena) Grow

func GrowByDouble(arena Arena) Grow {
	return func(pagesBefore, pagesNeeded int32, minSize uintptr) (pagesAdded int32, start, end uintptr) {
		if pagesBefore > pagesNeeded {
			pagesAdded = pagesBefore
		} else {
			pagesAdded = pagesNeeded
		}
		start, end = arena.Alloc(uintptr(pagesAdded) * _TLSFPageSize)
		if start == 0 {
			pagesAdded = pagesNeeded
			start, end = arena.Alloc(uintptr(pagesAdded) * _TLSFPageSize)
			if start == 0 {
				return 0, 0, 0
			}
		}
		return
	}
}

func GrowBy(pages int32, arena Arena) Grow {
	return func(pagesBefore, pagesNeeded int32, minSize uintptr) (pagesAdded int32, start, end uintptr) {
		if pages > pagesNeeded {
			pagesAdded = pages
		} else {
			pagesAdded = pagesNeeded
		}
		start, end = arena.Alloc(uintptr(pagesAdded) * _TLSFPageSize)
		if start == 0 {
			pagesAdded = pagesNeeded
			start, end = arena.Alloc(uintptr(pagesAdded) * _TLSFPageSize)
			if start == 0 {
				return 0, 0, 0
			}
		}
		return
	}
}

func GrowMin(arena Arena) Grow {
	return func(pagesBefore, pagesNeeded int32, minSize uintptr) (int32, uintptr, uintptr) {
		start, end := arena.Alloc(uintptr(pagesNeeded) * _TLSFPageSize)
		if start == 0 {
			return 0, 0, 0
		}
		return pagesNeeded, start, end
	}
}

func NewAllocator(pages int32) *Allocator {
	return NewAllocatorWithGrow(pages, NewSliceArena(), GrowMin)
}

func NewAllocatorWithGrow(pages int32, arena Arena, grow GrowFactory) *Allocator {
	if pages <= 0 {
		pages = 1
	}
	g := grow(arena)
	if g == nil {
		g = GrowMin(arena)
	}
	pagesAdded, start, end := g(0, pages, 0)
	a := bootstrap(start, end, pagesAdded, g)
	a.arena = uintptr(unsafe.Pointer(&arena))
	return a
}

type SyncAllocator struct {
	a     *Allocator
	stats Stats
	sync.Mutex
}

func (a *Allocator) ToSync() *SyncAllocator {
	return &SyncAllocator{a: a}
}

// Alloc allocates a block of memory that fits the size provided
//goland:noinspection GoVetUnsafePointer
func (a *SyncAllocator) Alloc(size uintptr) unsafe.Pointer {
	a.Lock()
	defer a.Unlock()
	p := uintptr(unsafe.Pointer(a.a.allocateBlock(size)))
	if p == 0 {
		return nil
	}
	return unsafe.Pointer(p + _TLSFBlockOverhead)
}

// Realloc determines the best way to resize an allocation.
func (a *SyncAllocator) Realloc(ptr unsafe.Pointer, size uintptr) unsafe.Pointer {
	a.Lock()
	defer a.Unlock()
	return unsafe.Pointer(uintptr(unsafe.Pointer(a.a.moveBlock(checkUsedBlock(uintptr(ptr)), size))) + _TLSFBlockOverhead)
}

// Free release the allocation back into the free list.
func (a *SyncAllocator) Free(ptr unsafe.Pointer) {
	a.Lock()
	defer a.Unlock()
	a.a.freeBlock(checkUsedBlock(uintptr(ptr)))
}
