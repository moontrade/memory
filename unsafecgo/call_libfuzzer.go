//go:build libfuzzer && (amd64 || arm64)
// +build libfuzzer
// +build amd64 arm64

package unsafecgo

import _ "unsafe"

// NonBlocking C function fn without going all the way through cgo.
// Example: NonBlocking((*byte)(C.my_c_func), 0, 0)
// 			void my_c_func(size_t arg0, size_t arg1) {
//			}
//go:noescape
//go:nosplit
//go:linkname NonBlocking runtime.libfuzzerCall
func NonBlocking(fn *byte, arg0, arg1 uintptr)

func Blocking(fn *byte, arg0, arg1 uintptr) {
	cgo.Call(fn, arg0, arg1)
}
