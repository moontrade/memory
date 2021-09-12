//go:build (tinygo.wasm || wasi) && !gc.conservative
// +build tinygo.wasm wasi
// +build !gc.conservative

package mem

import "unsafe"

const wasmPageSize = 64 * 1024

// WASM environments only have a single system allocator
var allocator *Allocator

//export llvm.wasm.memory.size.i32
func wasm_memory_size(index int32) int32

//export llvm.wasm.memory.Grow.i32
func wasm_memory_grow(index, pages int32) int32

//go:export initAllocator
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
	if heapStart == 0 || heapStart+unsafe.Sizeof(Allocator{})+_TLSFRootSize+_TLSFALMask+16 > uintptr(pagesBefore)*uintptr(wasmPageSize) {
		wasm_memory_grow(0, 1)
		pagesAfter = wasm_memory_size(0)
	}

	// Bootstrap allocator which will take over all allocations from now on.
	allocator = bootstrap(
		heapStart,
		uintptr(pagesAfter)*uintptr(wasmPageSize),
		1,
	)
}

// Alloc calls Alloc on the system allocator
func Alloc(size uintptr) unsafe.Pointer {
	return allocator.Alloc(size)
}

// Realloc calls Realloc on the system allocator
func Realloc(p unsafe.Pointer, size uintptr) unsafe.Pointer {
	return allocator.Realloc(p, size)
}

// Free calls Free on the system allocator
func Free(p unsafe.Pointer) {
	allocator.Free(p)
}

func SetGrow(grow Grow) {
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

// GrowBy will Grow by the number of pages specified or by the minimum needed, whichever is greater.
func GrowBy(pages int32) Grow {
	return func(pagesBefore, pagesNeeded int32, minSize uintptr) (pagesAdded int32, start, end uintptr) {
		if pages > pagesNeeded {
			pagesAdded = pages
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
}

// GrowByMin will Grow by a single page or by the minimum needed, whichever is greater.
func GrowMin(pagesBefore, pagesNeeded int32, minSize uintptr) (int32, uintptr, uintptr) {
	start, end := growBy(pagesNeeded)
	if start == 0 {
		return 0, 0, 0
	}
	return pagesNeeded, start, end
}

func growBy(pages int32) (uintptr, uintptr) {
	before := wasm_memory_size(0)
	wasm_memory_grow(0, pages)
	after := wasm_memory_size(0)
	if before == after {
		return 0, 0
	}
	return uintptr(before * wasmPageSize), uintptr(after * wasmPageSize)
}
