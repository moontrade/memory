//go:build !amd64 && !arm64
// +build !amd64,!arm64

package unsafecgo

import (
	"github.com/moontrade/memory/unsafecgo/cgo"
	"unsafe"
)

func NonBlocking(fn *byte, arg0, arg1 uintptr) {
	cgo.Blocking(uintptr(unsafe.Pointer(f)))
}

func Blocking(fn *byte, arg0, arg1 uintptr) {
	cgo.Blocking(fn, arg0, arg1)
}
