//go:build tinygo && tinygo.wasm
// +build tinygo,tinygo.wasm

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

func Scope(fn func(a Auto)) {
	a := NewAuto(32)
	fn(a)
	a.Free()
}

//// Scope creates an Auto free list that automatically reclaims memory
//// after callback finishes.
//func (a *Heap) Scope(fn func(a Auto)) {
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

const wasmPageSize = 64 * 1024

func WASMMemorySize(index int32) int32 {
	return wasm_memory_size(index)
}

func WASMMemorySizeBytes(index int32) int32 {
	return wasm_memory_size(index) * wasmPageSize
}

//export llvm.wasm.memory.size.i32
func wasm_memory_size(index int32) int32

//export llvm.wasm.memory.grow.i32
func wasm_memory_grow(index, pages int32) int32

func growBy(pages int32) (uintptr, uintptr) {
	before := wasm_memory_size(0)
	wasm_memory_grow(0, pages)
	after := wasm_memory_size(0)
	if before == after {
		return 0, 0
	}
	return uintptr(before * wasmPageSize), uintptr(after * wasmPageSize)
}

//go:linkname initAllocator runtime.initAllocator
func initAllocator(heapStart, heapEnd uintptr) {
	//println("initAllocator", uint(heapStart))
	//println("globals", uint(globalsStart), uint(globalsEnd))

	// Did we get called twice?
	if allocator != nil {
		return
	}

	var (
		pagesBefore = wasm_memory_size(0)
		pagesAfter  = pagesBefore
	)

	// Need to add a page initially to fit minimum size required by allocator?
	if heapStart == 0 || heapStart+unsafe.Sizeof(tlsf.Heap{})+tlsf.RootSize+tlsf.ALMask+16 >
		uintptr(pagesBefore)*uintptr(wasmPageSize) {
		// Just need a single page. Root size is much smaller than a single WASM page.
		wasm_memory_grow(0, 1)
		pagesAfter = wasm_memory_size(0)
	}

	// Bootstrap allocator which will take over all allocations from now on.
	allocator = tlsf.Bootstrap(
		heapStart,
		uintptr(pagesAfter)*uintptr(wasmPageSize),
		1,
		GrowMin,
	)
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
	start, end = growBy(pagesAdded)
	if start == 0 {
		pagesAdded = pagesNeeded
		start, end = growBy(pagesAdded)
		if start == 0 {
			return 0, 0, 0
		}
	}
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
	if growByPages > pagesNeeded {
		pagesAdded = growByPages
	} else {
		pagesAdded = pagesNeeded
	}
	start, end = growBy(pagesAdded)
	if start == 0 {
		pagesAdded = pagesNeeded
		start, end = growBy(pagesAdded)
		if start == 0 {
			return 0, 0, 0
		}
	}
	return
}

// GrowByMin will Grow by a single page or by the minimum needed, whichever is greater.
func GrowMin(pagesBefore, pagesNeeded int32, minSize uintptr) (int32, uintptr, uintptr) {
	start, end := growBy(pagesNeeded)
	if start == 0 {
		return 0, 0, 0
	}
	return pagesNeeded, start, end
}
