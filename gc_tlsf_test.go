package runtime

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
	"unsafe"
)

func Test_Counts(t *testing.T) {
	p := newTLSF(1)
	p1 := p.Alloc(38)
	println("alloc size", p.allocSize)
	p2 := p.Alloc(81)
	println("alloc size", p.allocSize)
	p.Free(p1)
	println("alloc size", p.allocSize)
	p.Free(p2)
	println("alloc size", p.allocSize)
}

func Test_Stress(t *testing.T) {
	stress(newTLSF(12), 1000000, 100, 7000, 10000,
		randomSize(0.80, 24, 128),
		randomSize(0.70, 128, 512),
		//randomSize(0.15, 128, 512),
		//randomSize(0.25, 512, 1024),
	)
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
	return rand.Intn(s.max-s.min) + s.min
}

func stress(allocator *tlsf, iterations, allocsPerIteration, minAllocs, maxAllocs int, sizeClasses ...*sizeClass) {
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

	for i := 0; i < iterations; i++ {
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(sz), func(i, j int) { sz[i], sz[j] = sz[j], sz[i] })

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

		rand.Shuffle(len(allocs), func(i, j int) { allocs[i], allocs[j] = allocs[j], allocs[i] })
		max := randomRange(minAllocs, maxAllocs)
		totalFrees += len(allocs) - max
		for x := max; x < len(allocs); x++ {
			alloc := allocs[x]
			allocator.Free(alloc.ptr)
			allocSize -= alloc.size
		}
		allocs = allocs[:max]
	}

	println("alloc bytes			", allocSize)
	println("alloc count			", len(allocs))
	println("total allocs		", totalAllocs)
	println("total frees			", totalFrees)
	println("total frees			", totalFrees)
	println("memory pages		", allocator.pages)
	println("heap size			", allocator.heapSize)
	println("free size			", allocator.freeSize)
	println("alloc size			", allocator.allocSize)
	//println("alloc size			", allocSize)
	println("max allocs			", maxAllocCount)
	println("max alloc size		", allocator.maxAllocSize)
	println("fragmentation		", fmt.Sprintf("%.2f", float64(allocator.heapSize-allocator.maxAllocSize)/float64(allocator.heapSize)))
}

func Test_Allocator(t *testing.T) {
	println("ALIGN_SIZE", 10<<3)
	pool := newTLSF(1)
	tlsfPrintInfo()
	ptr := pool.Alloc(16)
	ptr2 := pool.Alloc(49)
	ptr4 := pool.Alloc(8224)
	println("ptr", uint(uintptr(ptr)))
	println("ptr2", uint(uintptr(ptr2)))
	ptr3 := pool.Alloc(517)
	pool.Free(ptr)
	pool.Free(ptr4)
	pool.Free(ptr2)
	pool.Free(ptr3)
}

func BenchmarkAllocator_Alloc(b *testing.B) {
	pool := newTLSF(1)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if i%2 == 0 {
			b.SetBytes(1024)
			pool.Free(pool.Alloc(1024))
		} else {
			b.SetBytes(512)
			pool.Free(pool.Alloc(512))
		}
	}
}

func randomRange(min, max int) int {
	return rand.Intn(max-min) + min
}

//func address(block *tlsfBlock) uint {
//	if block == nil {
//		return uint(0)
//	}
//	return uint(uintptr(unsafe.Pointer(block)) - heapStart)
//}
//
//func getHeadOffset(fl uintptr, sl uint32) uintptr {
//	return tlsf_HL_START + (((fl << tlsf_SL_BITS) + uintptr(sl)) << tlsf_ALIGN_SIZE_LOG2)
//}
