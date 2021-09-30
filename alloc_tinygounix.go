//go:build tinygo && !tiny.wasm && (darwin || (linux && !baremetal && !wasi) || (freebsd && !baremetal)) && !nintendoswitch
// +build tinygo
// +build !tiny.wasm
// +build darwin linux,!baremetal,!wasi freebsd,!baremetal
// +build !nintendoswitch

package memory

import (
	"github.com/moontrade/memory/alloc/tlsf"
	"unsafe"
)

////////////////////////////////////////////////////////////////////////////////////
// Global Allocator
////////////////////////////////////////////////////////////////////////////////////

var allocator *tlsf.Heap

func HeapInstance() *tlsf.Heap {
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
	ptr := allocator.Alloc(size)
	return Pointer(ptr), tlsf.SizeOf(ptr)
}
func AllocZeroed(size uintptr) Pointer {
	return Pointer(allocator.AllocZeroed(size))
}
func AllocZeroedCap(size uintptr) (Pointer, uintptr) {
	ptr := allocator.AllocZeroed(size)
	return Pointer(ptr), tlsf.SizeOf(ptr)
}

// Alloc calls Alloc on the system allocator
//export alloc
func Calloc(num, size uintptr) Pointer {
	return Pointer(allocator.AllocZeroed(num * size))
}
func CallocCap(num, size uintptr) (Pointer, uintptr) {
	ptr := allocator.AllocZeroed(num * size)
	return Pointer(ptr), tlsf.SizeOf(ptr)
}

// Realloc calls Realloc on the system allocator
//export realloc
func Realloc(p Pointer, size uintptr) Pointer {
	return Pointer(allocator.Realloc(uintptr(p), size))
}
func ReallocCap(p Pointer, size uintptr) (Pointer, uintptr) {
	newPtr := allocator.Realloc(uintptr(p), size)
	return Pointer(newPtr), tlsf.SizeOf(newPtr)
}

// Free calls Free on the system allocator
//export free
func Free(p Pointer) {
	allocator.Free(uintptr(p))
}

func SizeOf(p Pointer) uintptr {
	return uintptr(tlsf.SizeOf(uintptr(p)))
}

func Scope(fn func(a AutoFree)) {
	a := NewAuto(32)
	fn(a)
	a.Free()
}

//// Scope creates an AutoFree free list that automatically reclaims memory
//// after callback finishes.
//func (a *Heap) Scope(fn func(a AutoFree)) {
//	auto := NewAuto(a.AsAllocator(), 32)
//	fn(auto)
//	auto.Free()
//}

func extalloc(size uintptr) unsafe.Pointer {
	ptr := unsafe.Pointer(allocator.Alloc(size))
	//println("extalloc", uint(uintptr(ptr)))
	return ptr
}

func extfree(ptr unsafe.Pointer) {
	//println("extfree", uint(uintptr(ptr)))
	allocator.Free(uintptr(ptr))
}

////////////////////////////////////////////////////////////////////////////////////
// tinygo hooks
////////////////////////////////////////////////////////////////////////////////////

const PageSize = 64 * 1024

func newTLSF(pages int32, grow tlsf.Grow) *tlsf.Heap {
	if pages <= 0 {
		pages = 1
	}
	size := uintptr(pages * PageSize)
	segment := uintptr(malloc(uintptr(size)))
	return tlsf.Bootstrap(segment, segment+size, pages, grow)
}

//go:linkname initAllocator runtime.initAllocator
func initAllocator(heapStart, heapEnd uintptr) {
	if allocator != nil {
		return
	}
	allocator = newTLSF(1, GrowMin)
}

func SetGrow(grow tlsf.Grow) {
	if allocator == nil {
		initAllocator(0, 0)
	}
	if allocator != nil {
		allocator.Grow = grow
	}
}

// GrowByDouble will double the heap on each Grow
func GrowByDouble(pagesBefore, pagesNeeded int32, minSize uintptr) (pagesAdded int32, start, end uintptr) {
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

var (
	growByPages int32 = 1
)

// GrowBy will Grow by the number of pages specified or by the minimum needed, whichever is greater.
func GrowBy(pages int32) tlsf.Grow {
	growByPages = pages
	return doGrowByPages
}

func doGrowByPages(pagesBefore, pagesNeeded int32, minSize uintptr) (pagesAdded int32, start, end uintptr) {
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

// GrowByMin will Grow by a single page or by the minimum needed, whichever is greater.
func GrowMin(pagesBefore, pagesNeeded int32, minSize uintptr) (int32, uintptr, uintptr) {
	ptr := uintptr(malloc(uintptr(pagesNeeded) * uintptr(PageSize)))
	if ptr == 0 {
		return 0, 0, 0
	}
	return pagesNeeded, ptr, uintptr(ptr) + (uintptr(pagesNeeded) * PageSize)
}
