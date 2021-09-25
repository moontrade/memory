//go:build tinygo && !tinygo.wasm && (darwin || linux || windows)
// +build tinygo
// +build !tinygo.wasm
// +build darwin linux windows

package tlsf

import "unsafe"

//go:linkname Malloc runtime.malloc
func Malloc(size uintptr) unsafe.Pointer

//go:linkname Copy runtime.Memcpy
func Copy(dst, src unsafe.Pointer, n uintptr)

//go:linkname Zero runtime.gcZero
func Zero(ptr unsafe.Pointer, size uintptr)

//go:linkname Equals runtime.gcMemequal
func Equals(x, y unsafe.Pointer, size uintptr) bool

//// TODO: Faster LLVM way?
//func memequal(a, b unsafe.Pointer, n uintptr) bool {
//	return gcMemequal(a, b, n)
//	//if a == nil {
//	//	return b == nil
//	//}
//	//return *(*string)(unsafe.Pointer(&_string{
//	//	ptr: uintptr(a),
//	//	len: int(n),
//	//})) == *(*string)(unsafe.Pointer(&_string{
//	//	ptr: uintptr(b),
//	//	len: int(n),
//	//}))
//}
