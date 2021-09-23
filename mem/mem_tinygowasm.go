//go:build tinygo && tinygo.wasm
// +build tinygo,tinygo.wasm

package mem

import (
	"unsafe"
)

////go:linkname gcMemcpy runtime.gcMemcpy
//func gcMemcpy(dst, src unsafe.Pointer, n uintptr)

//go:linkname Copy runtime.gcMemcpy
func Copy(dst, src unsafe.Pointer, n uintptr)

//go:linkname Move runtime.gcMemmove
func Move(dst, src unsafe.Pointer, size uintptr)

//func Memmove(dst, src unsafe.Pointer, size uintptr) {
//	gcMemmove(dst, src, size)
//}

//go:linkname Zero runtime.gcZero
func Zero(ptr unsafe.Pointer, size uintptr)

//func Memzero(ptr unsafe.Pointer, size uintptr) {
//	gcZero(ptr, size)
//}

type _bytes struct {
	Data uintptr
	Len  int
	Cap  int
}

func zeroSlow(ptr unsafe.Pointer, size uintptr) {
	b := *(*[]byte)(unsafe.Pointer(&_bytes{
		Data: uintptr(ptr),
		Len:  int(size),
		Cap:  int(size),
	}))
	switch {
	case size < 8:
		for i := 0; i < len(b); i++ {
			b[i] = 0
		}
	case size == 8:
		*(*uint64)(unsafe.Pointer(&b[0])) = 0
	default:
		var i = 0
		for ; i <= len(b)-8; i += 8 {
			*(*uint64)(unsafe.Pointer(&b[i])) = 0
		}
		for ; i < len(b); i++ {
			b[i] = 0
		}
	}
}

//go:linkname Equals runtime.gcMemequal
func Equals(x, y unsafe.Pointer, size uintptr) bool
