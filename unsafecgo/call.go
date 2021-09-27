//go:build !libfuzzer && (amd64 || arm64) && (linux || darwin)
// +build !libfuzzer
// +build amd64 arm64
// +build linux darwin

package unsafecgo

// Call C function fn without going all the way through cgo.
// Example: Call((*byte)(C.my_c_func), 0, 0)
// 			void my_c_func(size_t arg0, size_t arg1) {
//			}
//go:nosplit
//go:noescape
func Call(fn *byte, arg0, arg1 uintptr)
