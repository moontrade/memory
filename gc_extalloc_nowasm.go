//go:build !tinygo.wasm && !wasm && !wasi && gc.extalloc && !gc.tlsf
// +build !tinygo.wasm,!wasm,!wasi,gc.extalloc,!gc.tlsf

package runtime

import "unsafe"

func initHeap() {
}

//export memalloc
func memalloc(size uintptr) unsafe.Pointer {
	p := malloc(size)
	memzero(p, size)
	return p
}

//export memrealloc
func memrealloc(ptr unsafe.Pointer, size uintptr) unsafe.Pointer {
	return nil
}

//export memfree
func memfree(ptr unsafe.Pointer) {
	extfree(ptr)
}
