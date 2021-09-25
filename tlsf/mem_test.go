package tlsf

import (
	"testing"
	"unsafe"
)

func BenchmarkMemzero(b *testing.B) {
	buf := make([]byte, 16)
	buf[0] = 77
	buf[15] = 88

	Zero(unsafe.Pointer(&buf[0]), uintptr(len(buf)))

	ptr := unsafe.Pointer(&buf[0])

	b.Run("memclr", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			Zero(ptr, uintptr(len(buf)))
		}
	})

	b.Run("slow", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			zeroSlow(ptr, uintptr(len(buf)))
		}
	})
}
