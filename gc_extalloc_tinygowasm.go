//go:build (tinygo.wasm || wasi) && gc.extalloc && gc.tlsf
// +build tinygo.wasm wasi
// +build gc.extalloc
// +build gc.tlsf

package runtime

import "unsafe"

func initHeap() {
	var (
		heapBase    = heapStart
		pagesBefore = wasm_memory_size(0)
		pagesAfter  = pagesBefore
	)

	// Need to add a page?
	if heapBase+unsafe.Sizeof(tlsf{})+tlsf_ROOT_SIZE+tlsf_AL_MASK+16 > uintptr(pagesBefore)*uintptr(wasmPageSize) {
		wasm_memory_grow(0, 1)
		pagesAfter = wasm_memory_size(0)
	}

	allocator = initTLSF(
		heapStart,
		uintptr(pagesAfter)*uintptr(wasmPageSize),
		1,
	)
	heapEnd = allocator.heapEnd
}

func setHeapEnd(end uintptr) {
	heapEnd = end
}

func extalloc(size uintptr) unsafe.Pointer {
	return allocator.Alloc(size)
}

func extfree(ptr unsafe.Pointer) {
	allocator.Free(ptr)
}
