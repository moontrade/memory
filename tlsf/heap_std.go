//go:build !tinygo && !wasm && !wasi && !tinygo.wasm && (darwin || linux)
// +build !tinygo
// +build !wasm
// +build !wasi
// +build !tinygo.wasm
// +build darwin linux

package tlsf

import (
	"sync"
	"unsafe"
)

func NewHeap(pages int32) *Heap {
	return NewHeapWithConfig(pages, NewSliceArena(), GrowMin)
}

func NewHeapWithConfig(pages int32, pageAlloc Arena, grow GrowFactory) *Heap {
	if pages <= 0 {
		pages = 1
	}
	g := grow(pageAlloc)
	if g == nil {
		g = GrowMin(pageAlloc)
	}
	pagesAdded, start, end := g(0, pages, 0)
	a := Bootstrap(start, end, pagesAdded, g)
	a.arena = uintptr(unsafe.Pointer(&pageAlloc))
	return a
}

//// Scope creates an Auto free list that automatically reclaims memory
//// after callback finishes.
//func (a *Heap) Scope(fn func(a Auto)) {
//	if fn == nil {
//		return
//	}
//	auto := NewAuto(a.As, 32)
//	defer auto.Free()
//	fn(auto)
//}
//
//// Scope creates an Auto free list that automatically reclaims memory
//// after callback finishes.
//func (a *Sync) Scope(fn func(a Auto)) {
//	if fn == nil {
//		return
//	}
//	auto := NewAuto(a.AsAllocator(), 32)
//	defer auto.Free()
//	fn(auto)
//}

type SystemAllocator struct {
}

//// Alloc allocates a block of memory that fits the size provided.
//// The allocation IS cleared / zeroed out.
//func (SystemAllocator) AllocString(size uintptr) uintptr {
//	a.Lock()
//	defer a.Unlock()
//	return a.a.Alloc(size)
//}
//
//// AllocZeroed allocates a block of memory that fits the size provided.
//// The allocation is NOT cleared / zeroed out.
//func (SystemAllocator) AllocZeroed(size uintptr) uintptr {
//	a.Lock()
//	defer a.Unlock()
//	return a.a.AllocZeroed(size)
//}

type Sync struct {
	a    *Heap
	Slot uint8
	sync.Mutex
}

func (a *Heap) ToSync() *Sync {
	return &Sync{a: a}
}

func (a *Sync) Stats() Stats {
	return a.a.Stats
}

// Alloc allocates a block of memory that fits the size provided.
// The allocation IS cleared / zeroed out.
func (a *Sync) Alloc(size uintptr) uintptr {
	a.Lock()
	defer a.Unlock()
	return a.a.Alloc(size)
}

// AllocZeroed allocates a block of memory that fits the size provided.
// The allocation is NOT cleared / zeroed out.
func (a *Sync) AllocZeroed(size uintptr) uintptr {
	a.Lock()
	defer a.Unlock()
	return a.a.AllocZeroed(size)
}

// Realloc determines the best way to resize an allocation.
// Any extra size added is NOT cleared / zeroed out.
func (a *Sync) Realloc(ptr uintptr, size uintptr) uintptr {
	a.Lock()
	defer a.Unlock()
	return uintptr(unsafe.Pointer(a.a.moveBlock(checkUsedBlock(ptr), uintptr(size)))) + BlockOverhead
}

// Free release the allocation back into the free list.
func (a *Sync) Free(ptr uintptr) {
	a.Lock()
	defer a.Unlock()
	a.a.freeBlock(checkUsedBlock(ptr))
}

//// Alloc allocates a block of memory that fits the size provided.
//// The allocation IS cleared / zeroed out.
//func (a *Sync) AllocString(size uintptr) Bytes {
//	return newString(a.AsAllocator(), size)
//}

type GrowFactory func(arena Arena) Grow

func GrowByDouble(arena Arena) Grow {
	return func(pagesBefore, pagesNeeded int32, minSize uintptr) (pagesAdded int32, start, end uintptr) {
		if pagesBefore > pagesNeeded {
			pagesAdded = pagesBefore
		} else {
			pagesAdded = pagesNeeded
		}
		start, end = arena.Alloc(uintptr(pagesAdded) * PageSize)
		if start == 0 {
			pagesAdded = pagesNeeded
			start, end = arena.Alloc(uintptr(pagesAdded) * PageSize)
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
		start, end = arena.Alloc(uintptr(pagesAdded) * PageSize)
		if start == 0 {
			pagesAdded = pagesNeeded
			start, end = arena.Alloc(uintptr(pagesAdded) * PageSize)
			if start == 0 {
				return 0, 0, 0
			}
		}
		return
	}
}

func GrowMin(arena Arena) Grow {
	return func(pagesBefore, pagesNeeded int32, minSize uintptr) (int32, uintptr, uintptr) {
		start, end := arena.Alloc(uintptr(pagesNeeded) * PageSize)
		if start == 0 {
			return 0, 0, 0
		}
		return pagesNeeded, start, end
	}
}
