package mem

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
	"unsafe"
)

func Test_TLSFCounts(t *testing.T) {
	p := NewTLSF(1, GrowMin(DefaultMalloc))
	p1 := p.Alloc(38)
	println("alloc size", p.AllocSize)
	p2 := p.Alloc(81)
	println("alloc size", p.AllocSize)
	p.Free(p1)
	println("alloc size", p.AllocSize)
	p.Free(p2)
	println("alloc size", p.AllocSize)
}

func Test_TLSFThrash(t *testing.T) {
	thrashTLSF(NewTLSF(50, GrowMin(DefaultMalloc)), false,
		1000000, 100, 7000, 10000,
		randomSize(0.80, 24, 256),
		randomSize(0.70, 128, 512),
		//randomSize(0.15, 128, 512),
		randomSize(0.30, 128, 1024),
	)

	//thrashTLSF(newTLSF(2), 100000, 100, 12000, 17000,
	//	randomSize(0.80, 24, 96),
	//	//randomSize(0.70, 128, 512),
	//	//randomSize(0.15, 128, 512),
	//	//randomSize(0.30, 128, 1024),
	//)
}

type sizeClass struct {
	pct      float64
	min, max int
	next     func() int
}

func randomSize(pct float64, min, max int) *sizeClass {
	sz := &sizeClass{pct, min, max, nil}
	sz.next = sz.nextRandom
	return sz
}

func (s *sizeClass) nextRandom() int {
	if s.max == s.min {
		return s.max
	}
	return rand.Intn(s.max-s.min) + s.min
}

func thrashTLSF(
	allocator *Allocator, shuffle bool,
	iterations, allocsPerIteration, minAllocs, maxAllocs int,
	sizeClasses ...*sizeClass,
) {
	type allocation struct {
		ptr  unsafe.Pointer
		size int
	}

	sz := make([]int, 0, allocsPerIteration)
	for _, sc := range sizeClasses {
		for i := 0; i < int(float64(allocsPerIteration)*sc.pct); i++ {
			sz = append(sz, sc.next())
		}
	}

	allocs := make([]allocation, 0, maxAllocs)
	allocSize := 0
	totalAllocs := 0
	totalFrees := 0
	maxAllocCount := 0
	maxAllocSize := 0

	if shuffle {
		rand.Seed(time.Now().UnixNano())
	}

	start := time.Now()
	for i := 0; i < iterations; i++ {
		if shuffle {
			rand.Shuffle(len(sz), func(i, j int) { sz[i], sz[j] = sz[j], sz[i] })
		}

		for _, size := range sz {
			allocs = append(allocs, allocation{
				ptr:  allocator.Alloc(uintptr(size)), //tlsfalloc(uintptr(size)),
				size: size,
			})
			allocSize += size
		}
		totalAllocs += len(sz)

		if maxAllocCount < len(allocs) {
			maxAllocCount = len(allocs)
		}
		if allocSize > maxAllocSize {
			maxAllocSize = allocSize
		}

		if len(allocs) < minAllocs || len(allocs) < maxAllocs {
			continue
		}

		//rand.Shuffle(len(allocs), func(i, j int) { allocs[i], allocs[j] = allocs[j], allocs[i] })
		max := randomRange(minAllocs, maxAllocs)
		//max := maxAllocs
		totalFrees += len(allocs) - max
		for x := max; x < len(allocs); x++ {
			alloc := allocs[x]
			allocator.Free(alloc.ptr)
			allocSize -= alloc.size
		}
		allocs = allocs[:max]
	}

	elapsed := time.Now().Sub(start)
	seconds := float64(elapsed) / float64(time.Second)
	println("total time			", elapsed.String())
	fmt.Printf("allocs per sec		 %.1f million / sec\n", float64(float64(totalAllocs)/seconds/1000000))
	//println("allocs per sec		", float64(totalAllocs) / seconds)
	println("alloc bytes			", allocSize)
	println("alloc count			", len(allocs))
	println("total allocs		", totalAllocs)
	println("total frees			", totalFrees)
	println("total frees			", totalFrees)
	println("memory pages		", allocator.Pages)
	println("heap size			", allocator.HeapSize)
	println("free size			", allocator.FreeSize)
	println("alloc size			", allocator.AllocSize)
	//println("alloc size			", AllocSize)
	println("max allocs			", maxAllocCount)
	println("max alloc size		", allocator.MaxUsedSize)
	println("fragmentation		", fmt.Sprintf("%.2f", float64(allocator.HeapSize-allocator.MaxUsedSize)/float64(allocator.HeapSize)))
}

func Test_Allocator(t *testing.T) {
	println("ALIGN_SIZE", 10<<3)
	a := NewTLSF(1, GrowMin(DefaultMalloc))
	PrintDebugInfo()
	ptr := a.Alloc(16)
	ptr2 := a.Alloc(49)
	ptr4 := a.Alloc(8224)
	println("ptr", uint(uintptr(ptr)))
	println("ptr2", uint(uintptr(ptr2)))
	ptr3 := a.Alloc(517)
	a.Free(ptr)
	a.Free(ptr4)
	a.Free(ptr2)
	a.Free(ptr3)
}

func BenchmarkAllocator_Alloc(b *testing.B) {
	b.Run("Allocator alloc", func(b *testing.B) {
		pool := NewTLSF(1, GrowMin(DefaultMalloc))
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			size := randomRange(24, 4096)
			b.SetBytes(int64(size))
			pool.Free(pool.Alloc(512))
		}
	})

	b.Run("Go GC alloc", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			size := randomRange(24, 4096)
			b.SetBytes(int64(size))
			_ = make([]byte, size)
		}
	})
}

func randomRange(min, max int) int {
	return rand.Intn(max-min) + min
}
