package mem

import (
	"testing"
)

func BenchmarkHashSet(b *testing.B) {
	m := make(map[uintptr]struct{}, 16)
	m[1000] = struct{}{}
	a := NewAllocator(1, GrowMin(DefaultMalloc))
	set := newGCSet(a, 16)
	set.set(1000)
	s := &set

	b.Run("map get exists", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = m[1000]
		}
	})
	b.Run("moon get exists", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = s.has(1000)
		}
	})

	b.Run("map get not exists", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = m[1001]
		}
	})
	b.Run("moon get not exists", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = s.has(1001)
		}
	})

	b.Run("map set", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			m[1001] = struct{}{}
		}
	})
	b.Run("moon set", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = s.set(1001)
		}
	})

	b.Run("map del", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			delete(m, 1001)
		}
	})
	b.Run("moon del", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = s.del(1001)
		}
	})
}

func BenchmarkHashSetHashAlgos(b *testing.B) {
	a := NewAllocator(1, GrowMin(DefaultMalloc))
	set := newGCSet(a, 16)
	set.set(1000)
	s := &set

	get := func(name string, algo func(uint32) uint32) {
		b.Run(name+" get exists", func(b *testing.B) {
			gcSetHash = algo
			for i := 0; i < b.N; i++ {
				_ = s.has(1000)
			}
		})
	}
	getNot := func(name string, algo func(uint32) uint32) {
		b.Run(name+"get not exists", func(b *testing.B) {
			gcSetHash = algo
			for i := 0; i < b.N; i++ {
				_ = s.has(1001)
			}
		})
	}

	doSet := func(name string, algo func(uint32) uint32) {
		b.Run(name+" set", func(b *testing.B) {
			gcSetHash = algo
			for i := 0; i < b.N; i++ {
				_ = s.set(1001)
			}
		})
	}
	del := func(name string, algo func(uint32) uint32) {
		b.Run(name+" del", func(b *testing.B) {
			gcSetHash = algo
			for i := 0; i < b.N; i++ {
				_, _ = s.del(1001)
			}
		})
	}

	get("fnv", fnv32)
	get("adler32", adler32)
	get("wyhash", wyhash32)
	get("metro", metro32)

	getNot("fnv", fnv32)
	getNot("adler32", adler32)
	getNot("wyhash", wyhash32)
	getNot("metro", metro32)

	doSet("fnv", fnv32)
	doSet("adler32", adler32)
	doSet("wyhash", wyhash32)
	doSet("metro", metro32)

	del("fnv", fnv32)
	del("adler32", adler32)
	del("wyhash", wyhash32)
	del("metro", metro32)
}
