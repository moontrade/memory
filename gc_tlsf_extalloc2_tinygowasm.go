//go:build tinygo.wasm && gc.leaking && gc.extalloc
// +build tinygo.wasm,gc.leaking,gc.extalloc

package runtime

// This GC implementation is the simplest useful memory allocator possible: it
// only allocates memory and never frees it. For some constrained systems, it
// may be the only memory allocator possible.

import (
	"unsafe"
)

func GC() {
	// No-op.
}

func KeepAlive(x interface{}) {
	// Unimplemented. Only required with SetFinalizer().
}

func SetFinalizer(obj interface{}, finalizer interface{}) {
	// Unimplemented.
}

func initHeap() {
	// Nothing to initialize.
	initTSLF()
}

// setHeapEnd sets a new (larger) heapEnd pointer.
func setHeapEnd(newHeapEnd uintptr) {
	// This "heap" is so simple that simply assigning a new value is good
	// enough.
	heapEnd = newHeapEnd
}

func initTSLF() uintptr {
	if tlsfRoot != nil {
		return heapStart
	}

	var heapBase uintptr
	heapBase = heapStart
	rootOffset := (heapBase + tlsf_AL_MASK) & ^tlsf_AL_MASK

	var pagesBefore = wasm_memory_size(0)
	var pagesNeeded = int32((((int(rootOffset) + int(tlsf_ROOT_SIZE)) + 0xffff) & ^0xffff) >> 16)

	if pagesNeeded > pagesBefore && growHeapBy(pagesNeeded-pagesBefore) < 0 {
	}

	heapEnd = uintptr(wasm_memory_size(0) * wasmPageSize)

	tlsfRoot = (*Root)(unsafe.Pointer(rootOffset))
	tlsfRoot.flMap = 0
	setTail(tlsfRoot, nil)
	for fl := uintptr(0); fl < uintptr(tlsf_FL_BITS); fl++ {
		setSL(tlsfRoot, fl, 0)
		for sl := uint32(0); sl < tlsf_SL_SIZE; sl++ {
			setHead(tlsfRoot, fl, sl, nil)
		}
	}
	var memStart = rootOffset + tlsf_ROOT_SIZE
	memSize := uintptr(wasm_memory_size(0) * wasmPageSize)
	addMemory(tlsfRoot, memStart, memSize)
	return heapBase
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
