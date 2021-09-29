//go:build tinygo && (darwin || (linux && !baremetal && !wasi) || (freebsd && !baremetal)) && !nintendoswitch
// +build tinygo
// +build darwin linux,!baremetal,!wasi freebsd,!baremetal
// +build !nintendoswitch

package memory

import "unsafe"

//go:linkname malloc malloc
func malloc(size uintptr) unsafe.Pointer

func Compare(a, b unsafe.Pointer, n uintptr) int {
	if a == nil {
		if b == nil {
			return 0
		}
		return -1
	}
	ab := *(*string)(unsafe.Pointer(&_bytes{
		Data: uintptr(a),
		Len:  int(n),
	}))
	bb := *(*string)(unsafe.Pointer(&_bytes{
		Data: uintptr(b),
		Len:  int(n),
	}))
	if ab < bb {
		return -1
	}
	if ab == bb {
		return 0
	}
	return 1
}

//go:linkname Copy runtime.gcMemcpy
func Copy(dst, src unsafe.Pointer, n uintptr)

//go:linkname Move runtime.gcMemmove
func Move(dst, src unsafe.Pointer, size uintptr)

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
