//go:build tinygo.wasm || gc.extalloc
// +build tinygo.wasm gc.extalloc

package runtime

import "unsafe"

var tlsfPool *Pool

// Attach TLSF to heapStart
func initDefaultTLSF() {
	var (
		heapBase    = heapStart
		pagesBefore = wasm_memory_size(0)
	)

	// Need to add a page?
	if heapBase+unsafe.Sizeof(Pool{})+tlsf_ROOT_SIZE+tlsf_AL_MASK > uintptr(pagesBefore)*uintptr(wasmPageSize) {
		wasm_memory_grow(0, 1)
	}

	tlsfPool = (*Pool)(unsafe.Pointer(heapBase))
	tlsfPool.pages = int(wasm_memory_size(0))
	tlsfPool.heapStart = heapStart
	tlsfPool.heapEnd = uintptr(tlsfPool.pages) * uintptr(wasmPageSize)

	heapBase += unsafe.Sizeof(Pool{})
	rootOffset := (heapBase + tlsf_AL_MASK) & ^tlsf_AL_MASK

	heapEnd = uintptr(wasm_memory_size(0) * wasmPageSize)
	tlsfPool.heapEnd = heapEnd
	tlsfPool.root = (*Root)(unsafe.Pointer(rootOffset))
	tlsfPool.root.flMap = 0
	setTail(tlsfPool.root, nil)
	for fl := uintptr(0); fl < uintptr(tlsf_FL_BITS); fl++ {
		setSL(tlsfPool.root, fl, 0)
		for sl := uint32(0); sl < tlsf_SL_SIZE; sl++ {
			setHead(tlsfPool.root, fl, sl, nil)
		}
	}
	var memStart = rootOffset + tlsf_ROOT_SIZE
	memSize := uintptr(wasm_memory_size(0) * wasmPageSize)
	addMemory(tlsfPool.root, memStart, memSize)
}

func extalloc(size uintptr) unsafe.Pointer {
	p := unsafe.Pointer(tlsfalloc(size))
	//println("extalloc", uint(size), uint(uintptr(p)))
	return p
}

func extfree(ptr unsafe.Pointer) {
	//println("extfree", uint(uintptr(ptr)))
	tlsffree(uintptr(ptr))
}

func maybeFreeBlock(root *Root, block *tlsfBlock) {
	if uintptr(unsafe.Pointer(block)) >= heapStart {
		//if uintptr(unsafe.Pointer(block)) >= heapStart {
		freeBlock(root, block)
	}
}
