//go:build tinygo.wasm || gc.leaking
// +build tinygo.wasm gc.leaking

package runtime

import "unsafe"

var tlsfPool *Pool

func initHeap() {
	initTLSF()
	initITCMS()
}

// Attach TLSF to heapStart
func initTLSF() {
	var (
		heapBase    = heapStart
		pagesBefore = wasm_memory_size(0)
	)

	// Need to add a page?
	if heapBase+unsafe.Sizeof(Pool{})+tlsf_ROOT_SIZE+tlsf_AL_MASK+16 > uintptr(pagesBefore)*uintptr(wasmPageSize) {
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

func GC() {
	println("itcms collect()")
	itcmsCollect()
}

func KeepAlive(x interface{}) {
	// Unimplemented. Only required with SetFinalizer().
}

func SetFinalizer(obj interface{}, finalizer interface{}) {
	// Unimplemented.
}

func alloc(size uintptr) unsafe.Pointer {
	addr := itcmsNew(size)
	println("itcmsAlloc", uint(size), uint(addr))
	return unsafe.Pointer(addr)
}

func free(ptr unsafe.Pointer) {
	println("itcmsFree -> noop", uint(uintptr(ptr)))
}
