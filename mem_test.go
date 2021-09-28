package memory

import (
	"reflect"
	"testing"
	"unsafe"
)

func TestCmp(t *testing.T) {
	println(Cmp("hello", "hello0"))
	println(Cmp("hello0", "hello"))
	println(Cmp("hello", "hello"))
}

func BenchmarkCompare(b *testing.B) {
	b.Run("Slow", func(b *testing.B) {
		_a := "hello"
		_b := "hello8"
		__a := unsafe.Pointer(*(**uintptr)(unsafe.Pointer(&_a)))
		__b := unsafe.Pointer(*(**uintptr)(unsafe.Pointer(&_b)))

		for i := 0; i < b.N; i++ {
			compareSlow(__a, __b, uintptr(len(_a)))
		}
	})
	b.Run("Fast", func(b *testing.B) {
		_a := "hello"
		_b := "hello8"
		__a := unsafe.Pointer(*(**uintptr)(unsafe.Pointer(&_a)))
		__b := unsafe.Pointer(*(**uintptr)(unsafe.Pointer(&_b)))

		for i := 0; i < b.N; i++ {
			Compare(__a, __b, uintptr(len(_a)))
		}
	})
}

func compareSlow(a, b unsafe.Pointer, n uintptr) int {
	if a == nil {
		if b == nil {
			return 0
		}
		return -1
	}
	ab := *(*string)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(a),
		Len:  int(n),
	}))
	bb := *(*string)(unsafe.Pointer(&reflect.SliceHeader{
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
