package runtime

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
	"unsafe"
)

func Test_Counts(t *testing.T) {
	p := NewPool(1)
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
	stress(NewPool(1), 1000000, 100, 7000, 10000)
}

func stress(pool *Pool, iterations, allocsPerIteration, minAllocs, maxAllocs int) {
	type allocation struct {
		ptr  unsafe.Pointer
		size int
	}

	sz := make([]int, 0, allocsPerIteration)
	for i := 0; i < int(float64(allocsPerIteration)*0.99); i++ {
		switch i % 3 {
		case 0:
			sz = append(sz, 128)
		case 1:
			sz = append(sz, 32)
		case 2:
			sz = append(sz, 64)
		}

		//sz = append(sz, randomSize(8, 1024))
	}
	//for i := 0; i < int(float64(allocsPerIteration)*0.15); i++ {
	//	sz = append(sz, randomSize(256, 2048))
	//}
	//for i := 0; i < int(float64(allocsPerIteration)*0.04); i++ {
	//	sz = append(sz, randomSize(2049, 16384))
	//}
	//for i := 0; i < int(float64(allocsPerIteration)*0.01); i++ {
	//	sz = append(sz, randomSize(16385, 65536))
	//}

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
				ptr:  pool.Alloc(uintptr(size)), //tlsfalloc(uintptr(size)),
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
		max := randomSize(minAllocs, maxAllocs)
		totalFrees += len(allocs) - max
		for x := max; x < len(allocs); x++ {
			alloc := allocs[x]
			pool.Free(alloc.ptr)
			allocSize -= alloc.size
		}
		allocs = allocs[:max]
	}

	println("alloc bytes			", allocSize)
	println("alloc count			", len(allocs))
	println("total allocs		", totalAllocs)
	println("total frees			", totalFrees)
	println("total frees			", totalFrees)
	println("memory pages		", pool.pages)
	println("heap size			", pool.heapSize)
	println("free size			", pool.freeSize)
	println("alloc size			", pool.allocSize)
	//println("alloc size			", allocSize)
	println("max allocs			", maxAllocCount)
	println("max alloc size		", pool.maxAllocSize)
	println("fragmentation		", fmt.Sprintf("%.2f", float64(pool.heapSize-pool.maxAllocSize)/float64(pool.heapSize)))
}

func Test_Allocator(t *testing.T) {
	pool := NewPool(1)
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
	pool := NewPool(1)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if i%2 == 0 {
			pool.Free(pool.Alloc(38))
		} else {
			pool.Free(pool.Alloc(16))
		}
	}
}

func randomSize(min, max int) int {
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
//	return tlsf_HL_START + (((fl << tlsf_SL_BITS) + uintptr(sl)) << tlsf_ALIGNOF_USIZE)
//}
