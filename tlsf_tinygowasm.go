//go:build (tinygo.wasm || wasi) && !gc.conservative
// +build tinygo.wasm wasi
// +build !gc.conservative

package mem

import "unsafe"

const wasmPageSize = 64 * 1024

//export llvm.wasm.memory.size.i32
func wasm_memory_size(index int32) int32

//export llvm.wasm.memory.grow.i32
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
	if heapStart == 0 || Pointer(heapStart+unsafe.Sizeof(TLSF{}))+_TLSFRootSize+_TLSFALMask+16 >
		Pointer(uintptr(pagesBefore)*uintptr(wasmPageSize)) {
		wasm_memory_grow(0, 1)
		pagesAfter = wasm_memory_size(0)
	}

	// Bootstrap allocator which will take over all allocations from now on.
	allocator = bootstrap(
		Pointer(heapStart),
		Pointer(uintptr(pagesAfter)*uintptr(wasmPageSize)),
		1,
		GrowMin,
	)
}

// Alloc calls Alloc on the system allocator
func Alloc(size Pointer) Pointer {
	return allocator.Alloc(size)
}

// Realloc calls Realloc on the system allocator
func Realloc(p Pointer, size Pointer) Pointer {
	return allocator.Realloc(p, size)
}

// Free calls Free on the system allocator
func Free(p Pointer) {
	allocator.Free(p)
}

func Scope(fn func(a Auto)) {
	a := NewAuto(allocator.AsAllocator(), 32)
	fn(a)
	a.Free()
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
func GrowByDouble(pagesBefore, pagesNeeded int32, minSize Pointer) (pagesAdded int32, start, end Pointer) {
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
	return func(pagesBefore, pagesNeeded int32, minSize Pointer) (pagesAdded int32, start, end Pointer) {
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
func GrowMin(pagesBefore, pagesNeeded int32, minSize Pointer) (int32, Pointer, Pointer) {
	start, end := growBy(pagesNeeded)
	if start == 0 {
		return 0, 0, 0
	}
	return pagesNeeded, start, end
}

func growBy(pages int32) (Pointer, Pointer) {
	before := wasm_memory_size(0)
	wasm_memory_grow(0, pages)
	after := wasm_memory_size(0)
	if before == after {
		return 0, 0
	}
	return Pointer(before * wasmPageSize), Pointer(after * wasmPageSize)
}

//export heapAlloc
func heapAlloc(size uintptr) unsafe.Pointer {
	//println("heapAlloc")
	return unsafe.Pointer(allocator.Alloc(Pointer(size)))
}
