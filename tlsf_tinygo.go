//go:build tinygo
// +build tinygo

package mem

import (
	"unsafe"
)

//export memcpy
func memcpy(dst, src unsafe.Pointer, n uintptr)

//export memzero
func memzero(ptr unsafe.Pointer, size uintptr)
