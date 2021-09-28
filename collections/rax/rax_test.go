package rax

import (
	"github.com/moontrade/memory"
	"testing"
)

func TestNew(t *testing.T) {
	rax := New()
	rax.Size()
	rax.Free()
}

func TestRax_Insert(t *testing.T) {
	rax := New()
	key := memory.WrapString("hello")
	code, old := rax.Insert(key.Pointer, key.Len(), memory.WrapString("world").Pointer)
	println("code", code, "old", uint(old))
	code, old = rax.InsertBytes(memory.WrapString("hello9"), memory.WrapString("world 9").Pointer)
	println("code", code, "old", uint(old))
	code, old = rax.InsertBytes(memory.WrapString("hello7"), memory.WrapString("world 7").Pointer)
	println("code", code, "old", uint(old))
	code, old = rax.InsertBytes(memory.WrapString("hello"), memory.WrapString("world!").Pointer)
	println("code", code, "old", uint(old))
	if old != 0 {
		memory.Free(old)
	}
	found := rax.FindBytes(memory.WrapString("hello"))
	println("found", uint(found))
	rax.Print()
	rax.Free()
}

func BenchmarkRax_Insert(b *testing.B) {
	b.Run("insert", func(b *testing.B) {
		rax := New()
		defer rax.Free()

		value := memory.WrapString("world")
		key := memory.AllocBytes(8)

		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			key.SetInt64BE(0, int64(i))
			code, old := rax.Insert(key.Pointer, key.Len(), value.Pointer)
			_, _ = code, old
		}
	})
}
