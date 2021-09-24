package alloc

import "testing"

func BenchmarkNextAllocator(b *testing.B) {
	b.Run("fastrand", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = NextAllocatorRandom()
		}
	})
	b.Run("counter", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			NextAllocator()
		}
	})
}
