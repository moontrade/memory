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
	a.arena = Pointer(unsafe.Pointer(&arena))
	return a
}

const (
	_AllocatorNoSync IAllocator = 1 << 0
	_AllocatorSync   IAllocator = 1 << 1
	_AllocatorMask              = _AllocatorNoSync | _AllocatorSync
)

func toIAllocator(a *Allocator) IAllocator {
	return IAllocator(unsafe.Pointer(a)) | _AllocatorNoSync
}

func toIAllocatorSync(a *SyncAllocator) IAllocator {
	return IAllocator(unsafe.Pointer(a)) | _AllocatorSync
}

func (a IAllocator) Alloc(size Pointer) Pointer {
	if a&_AllocatorNoSync != 0 {
		return (*Allocator)(unsafe.Pointer(a & ^_AllocatorMask)).Alloc(size)
	} else {
		return (*SyncAllocator)(unsafe.Pointer(a & ^_AllocatorMask)).Alloc(size)
	}
}

func (a IAllocator) AllocZeroed(size Pointer) Pointer {
	if a&_AllocatorNoSync != 0 {
		return (*Allocator)(unsafe.Pointer(a & ^_AllocatorMask)).AllocZeroed(size)
	} else {
		return (*SyncAllocator)(unsafe.Pointer(a & ^_AllocatorMask)).AllocZeroed(size)
	}
}

func (a IAllocator) Bytes(length Pointer) Bytes {
	if a&_AllocatorNoSync != 0 {
		return (*Allocator)(unsafe.Pointer(a & ^_AllocatorMask)).Bytes(length)
	}
	return (*SyncAllocator)(unsafe.Pointer(a & ^_AllocatorMask)).Bytes(length)
}

func (a IAllocator) BytesCapacity(length, capacity Pointer) Bytes {
	if a&_AllocatorNoSync != 0 {
		return (*Allocator)(unsafe.Pointer(a & ^_AllocatorMask)).BytesCapacity(length, capacity)
	}
	return (*SyncAllocator)(unsafe.Pointer(a & ^_AllocatorMask)).BytesCapacity(length, capacity)
}

func (a IAllocator) Realloc(ptr, size Pointer) Pointer {
	if a&_AllocatorNoSync != 0 {
		return (*Allocator)(unsafe.Pointer(a & ^_AllocatorMask)).Realloc(ptr, size)
	}
	return (*SyncAllocator)(unsafe.Pointer(a & ^_AllocatorMask)).Realloc(ptr, size)
}

func (a IAllocator) Free(ptr Pointer) {
	if a&_AllocatorNoSync != 0 {
		(*Allocator)(unsafe.Pointer(a & ^_AllocatorMask)).Free(ptr)
	}
	(*SyncAllocator)(unsafe.Pointer(a & ^_AllocatorMask)).Free(ptr)
}

func (a *Allocator) AsIAllocator() IAllocator {
	return IAllocator(unsafe.Pointer(a)) | _AllocatorNoSync
}

type SyncAllocator struct {
	a     *Allocator
	stats Stats
	sync.Mutex
}

func (a *Allocator) ToSync() *SyncAllocator {
	return &SyncAllocator{a: a}
}

func (a *SyncAllocator) AsIAllocator() IAllocator {
	return IAllocator(unsafe.Pointer(a)) | _AllocatorSync
}

// Alloc allocates a block of memory that fits the size provided
//goland:noinspection GoVetUnsafePointer
func (a *SyncAllocator) Alloc(size Pointer) Pointer {
	a.Lock()
	defer a.Unlock()
	p := Pointer(unsafe.Pointer(a.a.allocateBlock(size)))
	if p == 0 {
		return p
	}
	return p + _TLSFBlockOverhead
}

// Alloc allocates a block of memory that fits the size provided
//goland:noinspection GoVetUnsafePointer
func (a *SyncAllocator) AllocZeroed(size Pointer) Pointer {
	a.Lock()
	defer a.Unlock()
	p := Pointer(unsafe.Pointer(a.a.allocateBlock(size)))
	if p == 0 {
		return p
	}
	p = p + _TLSFBlockOverhead
	memzero(unsafe.Pointer(p), uintptr(size))
	return p
}

//goland:noinspection GoVetUnsafePointer
func (a *SyncAllocator) Bytes(length Pointer) Bytes {
	return a.BytesCapacity(length, length)
}

//goland:noinspection GoVetUnsafePointer
func (a *SyncAllocator) BytesCapacity(length, capacity Pointer) Bytes {
	a.Lock()
	defer a.Unlock()
	if capacity < length {
		capacity = length
	}
	p := Pointer(unsafe.Pointer(a.a.allocateBlock(capacity)))
	if p == 0 {
		return Bytes{}
	}
	return Bytes{
		Pointer: p + _TLSFBlockOverhead,
		len:     uint32(length),
		cap:     uint32(*(*Pointer)(unsafe.Pointer(p)) & ^_TLSFTagsMask),
		alloc:   IAllocator(unsafe.Pointer(a)) | _AllocatorSync,
	}
}

// Realloc determines the best way to resize an allocation.
func (a *SyncAllocator) Realloc(ptr Pointer, size Pointer) Pointer {
	a.Lock()
	defer a.Unlock()
	return Pointer(unsafe.Pointer(a.a.moveBlock(checkUsedBlock(ptr), size))) + _TLSFBlockOverhead
}

// Free release the allocation back into the free list.
func (a *SyncAllocator) Free(ptr Pointer) {
	a.Lock()
	defer a.Unlock()
	a.a.freeBlock(checkUsedBlock(ptr))
}

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

////go:linkname memclrNoHeapPointers runtime.memclrNoHeapPointers
//func memclrNoHeapPointers(ptr unsafe.Pointer, n uintptr)

func memzero(ptr unsafe.Pointer, n uintptr) {
	//memclrNoHeapPointers(ptr, n)
	memzeroSlow(ptr, n)
}

func memzeroSlow(ptr unsafe.Pointer, size uintptr) {
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
	return func(pagesBefore, pagesNeeded int32, minSize Pointer) (pagesAdded int32, start, end Pointer) {
		if pagesBefore > pagesNeeded {
			pagesAdded = pagesBefore
		} else {
			pagesAdded = pagesNeeded
		}
		start, end = arena.Alloc(Pointer(pagesAdded) * _TLSFPageSize)
		if start == 0 {
			pagesAdded = pagesNeeded
			start, end = arena.Alloc(Pointer(pagesAdded) * _TLSFPageSize)
			if start == 0 {
				return 0, 0, 0
			}
		}
		return
	}
}

func GrowBy(pages int32, arena Arena) Grow {
	return func(pagesBefore, pagesNeeded int32, minSize Pointer) (pagesAdded int32, start, end Pointer) {
		if pages > pagesNeeded {
			pagesAdded = pages
		} else {
			pagesAdded = pagesNeeded
		}
		start, end = arena.Alloc(Pointer(pagesAdded) * _TLSFPageSize)
		if start == 0 {
			pagesAdded = pagesNeeded
			start, end = arena.Alloc(Pointer(pagesAdded) * _TLSFPageSize)
			if start == 0 {
				return 0, 0, 0
			}
		}
		return
	}
}

func GrowMin(arena Arena) Grow {
	return func(pagesBefore, pagesNeeded int32, minSize Pointer) (int32, Pointer, Pointer) {
		start, end := arena.Alloc(Pointer(pagesNeeded) * _TLSFPageSize)
		if start == 0 {
			return 0, 0, 0
		}
		return pagesNeeded, start, end
	}
}
