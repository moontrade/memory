package alloc

import (
	"github.com/moontrade/memory/tlsf"
	"testing"
)

func TestGC(t *testing.T) {
	// Backing Allocator for GC
	a := tlsf.NewHeap(1)

	// Create a simple roots marking system
	// This will provided by the runtime / compiler in TinyGo.
	roots := make(map[uintptr]struct{})
	var gc *gc
	markGlobals := func() {
		for k := range roots {
			gc.markRoot(k)
		}
	}
	// Create GC
	gc = newGC(a, 16, markGlobals, nil)

	// Allocate root
	root := func(size uintptr) uintptr {
		p := gc.New(size)
		roots[p] = struct{}{}
		return p
	}
	// Allocate and leak
	leak := func(size uintptr) uintptr {
		return gc.New(size)
	}

	// Do some allocations
	root(64)
	root(72)
	gc.Free(leak(128))
	leak(512)

	// Print before
	gc.Print()
	println()

	// Run full GC collect
	gc.Collect()

	// Print after
	gc.Print()
	println()
}

//func BenchmarkHashSet(b *testing.B) {
//	m := make(map[uintptr]struct{}, 16)
//	m[1000] = struct{}{}
//	a := NewTLSFArena(1, NewSliceArena(), GrowMin)
//	set := NewPointerSet(a.AsAllocator(), 16)
//	set.Set(1000)
//	s := &set
//
//	b.Run("map get exists", func(b *testing.B) {
//		for i := 0; i < b.N; i++ {
//			_, _ = m[1000]
//		}
//	})
//	b.Run("moon get exists", func(b *testing.B) {
//		for i := 0; i < b.N; i++ {
//			_ = s.Has(1000)
//		}
//	})
//
//	b.Run("map get not exists", func(b *testing.B) {
//		for i := 0; i < b.N; i++ {
//			_, _ = m[1001]
//		}
//	})
//	b.Run("moon get not exists", func(b *testing.B) {
//		for i := 0; i < b.N; i++ {
//			_ = s.Has(1001)
//		}
//	})
//
//	b.Run("map Add", func(b *testing.B) {
//		for i := 0; i < b.N; i++ {
//			m[1001] = struct{}{}
//		}
//	})
//	b.Run("moon Add", func(b *testing.B) {
//		for i := 0; i < b.N; i++ {
//			_, _ = s.Set(1001)
//		}
//	})
//
//	b.Run("map Delete", func(b *testing.B) {
//		for i := 0; i < b.N; i++ {
//			delete(m, 1001)
//		}
//	})
//	b.Run("moon Delete", func(b *testing.B) {
//		for i := 0; i < b.N; i++ {
//			_, _ = s.Delete(1001)
//		}
//	})
//}
//
//func BenchmarkHashSetHashAlgos(b *testing.B) {
//	a := NewTLSFArena(1, NewSliceArena(), GrowMin)
//	set := NewPointerSet(a.AsAllocator(), 16)
//	set.Set(1000)
//	s := &set
//
//	get := func(name string, algo func(uint32) uint32) {
//		b.Run(name+" get exists", func(b *testing.B) {
//			pointerSetHash = algo
//			for i := 0; i < b.N; i++ {
//				_ = s.Has(1000)
//			}
//		})
//	}
//	getNot := func(name string, algo func(uint32) uint32) {
//		b.Run(name+"get not exists", func(b *testing.B) {
//			pointerSetHash = algo
//			for i := 0; i < b.N; i++ {
//				_ = s.Has(1001)
//			}
//		})
//	}
//
//	doSet := func(name string, algo func(uint32) uint32) {
//		b.Run(name+" Add", func(b *testing.B) {
//			pointerSetHash = algo
//			for i := 0; i < b.N; i++ {
//				_, _ = s.Set(1001)
//			}
//		})
//	}
//	del := func(name string, algo func(uint32) uint32) {
//		b.Run(name+" Delete", func(b *testing.B) {
//			pointerSetHash = algo
//			for i := 0; i < b.N; i++ {
//				_, _ = s.Delete(1001)
//			}
//		})
//	}
//
//	get("fnv", fnv32)
//	get("adler32", adler32)
//	get("wyhash", wyhash32)
//	get("metro", metro32)
//
//	getNot("fnv", fnv32)
//	getNot("adler32", adler32)
//	getNot("wyhash", wyhash32)
//	getNot("metro", metro32)
//
//	doSet("fnv", fnv32)
//	doSet("adler32", adler32)
//	doSet("wyhash", wyhash32)
//	doSet("metro", metro32)
//
//	del("fnv", fnv32)
//	del("adler32", adler32)
//	del("wyhash", wyhash32)
//	del("metro", metro32)
//}
//
//func Test_ThrashPointerSet(t *testing.T) {
//	allocator := NewTLSF(25)
//
//	var (
//		iterations                 = 1000
//		heapBase           uintptr = 67056
//		allocsPerIteration         = 100
//		minAllocs                  = 32
//		maxAllocs                  = 512
//	)
//
//	run := func(name string, fn func(uint32) uint32) {
//		_set := NewPointerSet(allocator.AsAllocator(), uintptr(maxAllocs*128))
//		set := &_set
//		defer set.Close()
//		pointerSetHash = fn
//		println(name+" collisions", thrashPointerSet(set, heapBase, true,
//			iterations, allocsPerIteration, minAllocs, maxAllocs,
//			//randomSize(0.95, 16, 48),
//			randomSize(1, 24, 384),
//			//randomSize(0.55, 64, 512),
//			//randomSize(0.70, 128, 512),
//			//randomSize(0.15, 128, 512),
//			//randomSize(0.30, 128, 1024),
//		))
//	}
//
//	run("fnv", fnv32)
//	run("wy", wyhash32)
//	run("metro", metro32)
//	run("adler", adler32)
//}
//
//func thrashPointerSet(
//	set *PointerSet,
//	heapBase uintptr,
//	shuffle bool,
//	iterations, allocsPerIteration, minAllocs, maxAllocs int,
//	sizeClasses ...*sizeClass,
//) int64 {
//	type allocation struct {
//		ptr  Pointer
//		size uintptr
//	}
//
//	sz := make([]int, 0, allocsPerIteration)
//	for _, sc := range sizeClasses {
//		for i := 0; i < int(float64(allocsPerIteration)*sc.pct); i++ {
//			sz = append(sz, sc.next())
//		}
//	}
//
//	allocs := make([]allocation, 0, maxAllocs)
//	var (
//		collisions    int64   = 0
//		allocSize     uintptr = 0
//		totalAllocs           = 0
//		totalFrees            = 0
//		maxAllocCount         = 0
//		maxAllocSize  uintptr = 0
//	)
//
//	rand.Seed(time.Now().UnixNano())
//
//	start := time.Now()
//	for i := 0; i < iterations; i++ {
//		rand.Shuffle(len(sz), func(i, j int) { sz[i], sz[j] = sz[j], sz[i] })
//
//		nextPtr := Pointer(heapBase)
//		for _, size := range sz {
//			allocs = append(allocs, allocation{
//				ptr:  nextPtr,
//				size: uintptr(size),
//			})
//			allocSize += uintptr(size)
//			if !set.Has(nextPtr) {
//				if set.isCollision(nextPtr) {
//					collisions++
//				}
//				set.Set(nextPtr)
//			}
//			nextPtr += Pointer(size)
//		}
//		totalAllocs += len(sz)
//
//		if maxAllocCount < len(allocs) {
//			maxAllocCount = len(allocs)
//		}
//		if allocSize > maxAllocSize {
//			maxAllocSize = allocSize
//		}
//
//		if len(allocs) < minAllocs || len(allocs) < maxAllocs {
//			continue
//		}
//
//		rand.Shuffle(len(allocs), func(i, j int) { allocs[i], allocs[j] = allocs[j], allocs[i] })
//		max := randomRange(minAllocs, maxAllocs)
//		//max := maxAllocs
//		totalFrees += len(allocs) - max
//		for x := max; x < len(allocs); x++ {
//			alloc := allocs[x]
//			set.Delete(alloc.ptr)
//			allocSize -= alloc.size
//		}
//		allocs = allocs[:max]
//	}
//
//	_ = start
//
//	return collisions
//}
