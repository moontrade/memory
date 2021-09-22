package mem

import (
	"testing"
)

func TestSDS(t *testing.T) {
	a := NewTLSF(1)
	str := newString(a.AsAllocator(), 32)
	println("size 32", "cap", str.Cap())
	str = newString(a.AsAllocator(), 8)
	println("size 8", "cap", str.Cap())
	str = newString(a.AsAllocator(), 37)
	println("size 37", "cap", str.Cap())
	str = newString(a.AsAllocator(), 39)
	println("size 39", "cap", str.Cap())
	str = newString(a.AsAllocator(), 41)
	println("size 41", "cap", str.Cap())
}

func BenchmarkSDS(b *testing.B) {
	buf := make([]byte, 128, 256)
	b.Run("len", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = len(buf)
		}
	})

	a := NewTLSF(1)
	str := newString(a.AsAllocator(), 32)
	str.Cap()
	str = newString(a.AsAllocator(), 39)
	str = newString(a.AsAllocator(), 41)
	b.Run("Str.len", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = str.Len()
		}
	})
	b.Run("Str.cap", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = str.Cap()
		}
	})
}
