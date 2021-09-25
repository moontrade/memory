//go:build tinygo && !tinygo.wasm && (darwin || linux || windows)
// +build tinygo
// +build !tinygo.wasm
// +build darwin linux windows

package memory

import (
	"unsafe"
)

////////////////////////////////////////////////////////////////////////////////////
// Global Allocator
////////////////////////////////////////////////////////////////////////////////////

var allocator *TLSF

func TLSFInstance() *TLSF {
	return allocator
}

////////////////////////////////////////////////////////////////////////////////////
// Global allocator convenience
////////////////////////////////////////////////////////////////////////////////////

// Alloc calls Alloc on the system allocator
func Alloc(size uintptr) Pointer {
	return Pointer(allocator.Alloc(size))
}
func AllocCap(size uintptr) (Pointer, uintptr) {
	ptr := Pointer(allocator.Alloc(size))
	return ptr, tlsf.SizeOf(ptr)
}
func AllocZeroed(size uintptr) Pointer {
	return Pointer(allocator.AllocZeroed(size))
}
func AllocZeroedCap(size uintptr) (Pointer, uintptr) {
	ptr := Pointer(allocator.AllocZeroed(size))
	return ptr, tlsf.SizeOf(ptr)
}

// Alloc calls Alloc on the system allocator
//export alloc
func Calloc(num, size uintptr) Pointer {
	return Pointer(allocator.AllocZeroed(num * size))
}
func CallocCap(num, size uintptr) (Pointer, uintptr) {
	ptr := Pointer(allocator.AllocZeroed(num * size))
	return ptr, tlsf.SizeOf(ptr)
}

// Realloc calls Realloc on the system allocator
//export realloc
func Realloc(p Pointer, size uintptr) Pointer {
	return Pointer(allocator.Realloc(uintptr(p), size))
}
func ReallocCap(p Pointer, size uintptr) (Pointer, uintptr) {
	newPtr := Pointer(allocator.Realloc(uintptr(p), size))
	return newPtr, tlsf.SizeOf(newPtr)
}

// Free calls Free on the system allocator
//export free
func Free(p Pointer) {
	allocator.Free(uintptr(p))
}

func SizeOf(p Pointer) uintptr {
	return uintptr(tlsf.SizeOf(uintptr(p)))
}

func Scope(fn func(a Auto)) {
	a := NewAuto(32)
	fn(a)
	a.Free()
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
// tinygo hooks
////////////////////////////////////////////////////////////////////////////////////

const wasmPageSize = 64 * 1024

func newTLSF(pages int32, grow Grow) *Heap {
	if pages <= 0 {
		pages = 1
	}
	size := uintptr(pages * wasmPageSize)
	segment := uintptr(malloc(uintptr(size)))
	return Bootstrap(segment, segment+size, pages, grow)
}

//go:linkname initAllocator runtime.initAllocator
func initAllocator(heapStart, heapEnd uintptr) {
	if allocator != nil {
		return
	}
	allocator = newTLSF(1, GrowMin)
	allocator_ = Allocator(unsafe.Pointer(allocator))
}

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
