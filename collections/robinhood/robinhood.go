//go:build !tinygo
// +build !tinygo

package robinhood

/*

#cgo CXXFLAGS: -std=c++14
#cgo darwin,amd64 LDFLAGS: -Wl,-rpath,${SRCDIR}/lib/darwin_amd64 -L${SRCDIR}/lib/darwin_amd64
#cgo darwin,amd64 LDFLAGS: -lrobinhood
#include <stdlib.h>
#include "library.h"

void do_robinhood_hello(size_t arg0, size_t arg1) {
	robinhood_hello(0, 0);
}
*/
import "C"
import (
	"github.com/moontrade/memory/unsafecgo"
	"unsafe"
)

type Tree uintptr

type rhNewT struct {
	ptr  uintptr
	code uintptr
}

func New() *Tree {
	args := rhNewT{}
	ptr := uintptr(unsafe.Pointer(&args))
	unsafecgo.Call((*byte)(C.do_robinhood_hello), ptr, 0)
	return (*Tree)(unsafe.Pointer(args.ptr))
}
