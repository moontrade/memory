//go:build !libfuzzer && !tinygo && (amd64 || arm64) && (linux || darwin)
// +build !libfuzzer
// +build !tinygo
// +build amd64 arm64
// +build linux darwin

package unsafecgo

import "github.com/moontrade/memory/unsafecgo/cgo"

// NonBlocking C function fn without going all the way through cgo
// Be very careful using it. If the C code blocks it can/will
// lock up your app.
// Example: NonBlocking((*byte)(C.my_c_func), 0, 0)
// 			void my_c_func(size_t arg0, size_t arg1) {
//			}
//go:nosplit
//go:noescape
func NonBlocking(fn *byte, arg0, arg1 uintptr)

func Blocking(fn *byte, arg0, arg1 uintptr) {
	cgo.Blocking(fn, arg0, arg1)
}
