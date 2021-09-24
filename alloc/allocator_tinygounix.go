//go:build tinygo && !tinygo.wasm && (darwin || linux || windows)
// +build tinygo
// +build !tinygo.wasm
// +build darwin linux windows

package alloc

import (
	"unsafe"
)

////////////////////////////////////////////////////////////////////////////////////
// Global Allocator
////////////////////////////////////////////////////////////////////////////////////

var allocator *TLSF
var allocator_ Allocator

func TLSFInstance() *TLSF {
	return allocator
}

////////////////////////////////////////////////////////////////////////////////////
// Global allocator convenience
////////////////////////////////////////////////////////////////////////////////////

// Alloc calls Alloc on the system allocator
func Alloc(size Pointer) Pointer {
	return allocator.Alloc(size)
}

// Alloc calls Alloc on the system allocator
func AllocNotCleared(size Pointer) Pointer {
	return allocator.AllocZeroed(size)
}

// Realloc calls Realloc on the system allocator
func Realloc(p Pointer, size Pointer) Pointer {
	return allocator.Realloc(p, size)
}

// Free calls Free on the system allocator
func Free(p Pointer) {
	allocator.Free(p)
}

// Alloc calls Alloc on the system allocator
func AllocBytes(length Pointer) Bytes {
	return allocator.Bytes(length)
}

// Alloc calls Alloc on the system allocator
func AllocBytesCap(length, capacity Pointer) Bytes {
	return allocator.BytesCapNotCleared(length, capacity)
}

// Alloc calls Alloc on the system allocator
func AllocBytesCapNotCleared(length, capacity Pointer) Bytes {
	return allocator.BytesCapNotCleared(length, capacity)
}

func Scope(fn func(a Auto)) {
	a := NewAuto(allocator.AsAllocator(), 32)
	fn(a)
	a.Free()
}

// Scope creates an Auto free list that automatically reclaims memory
// after callback finishes.
func (a *TLSF) Scope(fn func(a Auto)) {
	auto := NewAuto(a.AsAllocator(), 32)
	fn(auto)
	auto.Free()
}

//go:export extalloc
func extalloc(size uintptr) unsafe.Pointer {
	ptr := unsafe.Pointer(allocator.Alloc(Pointer(size)))
	//println("extalloc", uint(uintptr(ptr)))
	return ptr
}

//go:export extfree
func extfree(ptr unsafe.Pointer) {
	//println("extfree", uint(uintptr(ptr)))
	allocator.Free(Pointer(ptr))
}

////////////////////////////////////////////////////////////////////////////////////
// Allocator facade
////////////////////////////////////////////////////////////////////////////////////

// Alloc calls Alloc on the system allocator
func Alloc(size uintptr) Pointer {
	return allocator.Alloc(size)
}

// Alloc calls Alloc on the system allocator
//export alloc
func AllocZeroed(size uintptr) Pointer {
	return allocator.AllocZeroed(size)
}

// Realloc calls Realloc on the system allocator
//export realloc
func Realloc(p Pointer, size uintptr) Pointer {
	return allocator.Realloc(p, size)
}

// Free calls Free on the system allocator
//export free
func Free(p Pointer) {
	allocator.Free(p)
}

func SizeOf(p Pointer) uintptr {
	return uintptr(tlsf.AllocationSize(p))
}

func ReadStats() HeapStats {
	return *(*HeapStats)(unsafe.Pointer(&allocator.Stats))
}

func (a Allocator) Stats() HeapStats {
	return *(*HeapStats)(unsafe.Pointer(&allocator.Stats))
}

////////////////////////////////////////////////////////////////////////////////////
// tinygo hooks
////////////////////////////////////////////////////////////////////////////////////

const wasmPageSize = 64 * 1024

//go:linkname initAllocator runtime.initAllocator
func initAllocator(heapStart, heapEnd uintptr) {
	if allocator != nil {
		return
	}
	allocator = newTLSF(1, GrowMin)
	allocator_ = Allocator(unsafe.Pointer(allocator))
}

//func NewTLSF(pages int32) *TLSF {
//	if pages <= 0 {
//		pages = 1
//	}
//
//	size := uintptr(pages * wasmPageSize)
//	start := uintptr(wasm_memory_size(0) * wasmPageSize)
//	wasm_memory_grow(0, pages)
//	end := uintptr(wasm_memory_size(0) * wasmPageSize)
//	if start == end {
//		panic("out of memory")
//	}
//	return initTLSF(start, start+size, pages)
//}
//
//func (a *TLSF) Grow(pages int32) (uintptr, uintptr) {
//	if pages <= 0 {
//		pages = 1
//	}
//
//	// wasm memory grow
//	before := wasm_memory_size(0)
//	wasm_memory_grow(0, pages)
//	after := wasm_memory_size(0)
//	if before == after {
//		return 0, 0
//	}
//
//	a.Pages += int32(pages)
//	start := uintptr(before * wasmPageSize)
//	end := uintptr(after * wasmPageSize)
//	a.HeapEnd = end
//	return start, end
//}

func SetGrow(grow Grow) {
	if allocator == nil {
		initAllocator(0, 0)
	}
	if allocator != nil {
		allocator.Grow = grow
	}
}

func GrowByDouble() Grow {
	return func(pagesBefore, pagesNeeded int32, minSize uintptr) (pagesAdded int32, start, end uintptr) {
		if pagesBefore > pagesNeeded {
			pagesAdded = pagesBefore
		} else {
			pagesAdded = pagesNeeded
		}
		ptr := uintptr(malloc(uintptr(pagesAdded) * uintptr(PageSize)))
		if ptr == 0 {
			pagesAdded = pagesNeeded
			ptr = uintptr(malloc(uintptr(pagesAdded) * uintptr(PageSize)))
			if ptr == 0 {
				return 0, 0, 0
			}
		}
		start = ptr
		end = start + (uintptr(pagesAdded) * PageSize)
		return
	}
}

func GrowBy(pages int32) Grow {
	return func(pagesBefore, pagesNeeded int32, minSize uintptr) (pagesAdded int32, start, end uintptr) {
		if pages > pagesNeeded {
			pagesAdded = pages
		} else {
			pagesAdded = pagesNeeded
		}
		ptr := uintptr(malloc(uintptr(pagesAdded) * uintptr(PageSize)))
		if ptr == 0 {
			pagesAdded = pagesNeeded
			ptr = uintptr(malloc(uintptr(pagesAdded) * uintptr(PageSize)))
			if ptr == 0 {
				return 0, 0, 0
			}
		}
		start = ptr
		end = start + (uintptr(pagesAdded) * PageSize)
		return
	}
}

func GrowMin(pagesBefore, pagesNeeded int32, minSize uintptr) (int32, uintptr, uintptr) {
	ptr := uintptr(malloc(uintptr(pagesNeeded) * uintptr(PageSize)))
	if ptr == 0 {
		return 0, 0, 0
	}
	return pagesNeeded, ptr, uintptr(ptr) + (uintptr(pagesNeeded) * PageSize)
}

func newTLSF(pages int32, grow Grow) *Heap {
	if pages <= 0 {
		pages = 1
	}
	size := uintptr(pages * wasmPageSize)
	segment := uintptr(malloc(uintptr(size)))
	return Bootstrap(segment, segment+size, pages, grow)
}
