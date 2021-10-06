//go:build !tinygo
// +build !tinygo

package rpmalloc

/*
#cgo darwin,amd64 LDFLAGS: -Wl,-rpath,${SRCDIR}/lib/darwin_amd64 -L${SRCDIR}/lib/darwin_amd64
#cgo darwin,amd64 LDFLAGS: -lrpmalloc -ldl -lc -lm
#cgo linux,amd64 LDFLAGS: -Wl,-rpath,${SRCDIR}/lib/linux_amd64 -L${SRCDIR}/lib/linux_amd64
#cgo linux,amd64 LDFLAGS: -lrpmalloc -ldl -lc -lm
#cgo linux,amd64 CFLAGS: -D_GNU_SOURCE
#include "rpmalloc.h"
*/
import "C"

func init() {
	C.rpmalloc_initialize()
}
func Init() {}

type Heap C.rpmalloc_heap_t
