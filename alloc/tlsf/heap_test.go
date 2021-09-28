package tlsf

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"testing"
	"time"
)

func Test_AllocatorCounts(t *testing.T) {
	p := NewHeapWithConfig(1, NewSysArena(), GrowMin)
	p1 := p.Alloc(38)
	println("alloc size", p.AllocSize)
	p2 := p.Alloc(81)
	println("alloc size", p.AllocSize)
	p.Free(p1)
	println("alloc size", p.AllocSize)
	p.Free(p2)
	println("alloc size", p.AllocSize)
}

//func Benchmark_Ctz32(b *testing.B) {
//	b.Run("bits.TrailingZeroes32", func(b *testing.B) {
//		for i := 0; i < b.N; i++ {
//			bits.TrailingZeros32(100)
//			bits.TrailingZeros32(101)
//			bits.TrailingZeros32(102)
//		}
//	})
//
//	b.Run("Ctz32", func(b *testing.B) {
//		for i := 0; i < b.N; i++ {
//			Ctz32(100)
//			Ctz32(101)
//			Ctz32(102)
//		}
//	})
//}

func Test_AllocatorThrash(t *testing.T) {
	a := NewSysArena()
	statsBefore := runtime.MemStats{}
	runtime.ReadMemStats(&statsBefore)
	thrashAllocator(NewHeapWithConfig(1, a, GrowMin), false,
		1000000, 100, 15000, 21000,
		randomSize(0.95, 16, 48),
		randomSize(0.95, 48, 192),
		randomSize(0.55, 64, 512),
		//randomSize(0.70, 128, 512),
		//randomSize(0.15, 128, 512),
		//randomSize(0.30, 128, 1024),
	)

	var statsAfter runtime.MemStats
	runtime.ReadMemStats(&statsAfter)
	fmt.Println("SysAllocator Size", a.Size())
	fmt.Println("GCStats Before", statsBefore)
	fmt.Println("GCStats After", statsAfter)

	//thrashAllocator(newAllocator(2), 100000, 100, 12000, 17000,
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

func thrashAllocator(
	allocator *Heap, shuffle bool,
	iterations, allocsPerIteration, minAllocs, maxAllocs int,
	sizeClasses ...*sizeClass,
) {
	type allocation struct {
		ptr  uintptr
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
	println("max alloc size		", allocator.PeakAllocSize)
	println("fragmentation		", fmt.Sprintf("%.2f", allocator.Stats.Fragmentation()))
}

//func TestAllocatorKind(t *testing.T) {
//	a := NewHeapWithConfig(1, NewSysArena(), GrowMin)
//	s := a.ToSync()
//
//	p1 := toTLSFAllocator(a)
//	p2 := toTLSFSyncAllocator(s)
//
//	p1.Free(p1.Alloc(64))
//	p2.Free(p2.Alloc(128))
//
//	println(p1, p1|_TLSFNoSync)
//	println(p2, p2|_TLSFSync, (p2|_TLSFSync) & ^_AllocatorMask)
//	println((p2 | _TLSFSync) & ^_AllocatorMask)
//	println((p1 | _TLSFNoSync) & _AllocatorMask)
//	println((p2 | _TLSFSync) & _AllocatorMask)
//}

func Test_Allocator(t *testing.T) {
	a := NewHeapWithConfig(1, NewSysArena(), GrowMin)
	PrintDebugInfo()
	ptr := a.Alloc(16)
	ptr2 := a.Alloc(49)
	ptr4 := a.Alloc(8224)
	println("ptr", uint(ptr))
	println("ptr2", uint(ptr2))
	ptr3 := a.Alloc(517)
	a.Free(ptr)
	a.Free(ptr4)
	a.Free(ptr2)
	a.Free(ptr3)
}

func BenchmarkAllocator_Alloc(b *testing.B) {
	var (
		min, max    = 24, 2048
		showGCStats = true
	)
	doAfter := func(before, after runtime.MemStats) {
		if showGCStats {
			fmt.Println("Before", "GC CPU", before.GCCPUFraction, "TotalAllocs", before.TotalAlloc, "Frees", before.Frees, "PauseNs Total", before.PauseTotalNs)
			fmt.Println("After ", "GC CPU", after.GCCPUFraction, "TotalAllocs", after.TotalAlloc, "Frees", after.Frees, "PauseNs Total", after.PauseTotalNs)
			println()
		}
	}
	//b.Run("Random Range time", func(b *testing.B) {
	//	b.ReportAllocs()
	//	b.ResetTimer()
	//	for i := 0; i < b.N; i++ {
	//		randomRange(min, max)
	//	}
	//	after()
	//})

	randomRangeSizes := make([]uintptr, 0, 256)
	for i := 0; i < 1000; i++ {
		randomRangeSizes = append(randomRangeSizes, uintptr(randomRange(min, max)))
	}

	//b.Run("Heap alloc", func(b *testing.B) {
	//	a := NewHeapWithConfig(1, NewSysArena(), GrowMin)
	//	var before runtime.MemStats
	//	runtime.ReadMemStats(&before)
	//	b.ReportAllocs()
	//	b.ResetTimer()
	//	for i := 0; i < b.N; i++ {
	//		size := randomRangeSizes[i%len(randomRangeSizes)]
	//		b.SetBytes(int64(size))
	//		a.Free(a.Alloc(uintptr(size)))
	//	}
	//	after(before)
	//})
	//b.Run("Heap allocZeroed", func(b *testing.B) {
	//	a := NewHeapWithConfig(1, NewSysArena(), GrowMin)
	//	var before runtime.MemStats
	//	runtime.ReadMemStats(&before)
	//	b.ReportAllocs()
	//	b.ResetTimer()
	//	for i := 0; i < b.N; i++ {
	//		size := randomRangeSizes[i%len(randomRangeSizes)]
	//		b.SetBytes(int64(size))
	//		a.Free(a.AllocZeroed(uintptr(size)))
	//	}
	//	after(before)
	//})
	//b.Run("Allocator alloc", func(b *testing.B) {
	//	al := NewHeapWithConfig(1, NewSysArena(), GrowMin)
	//	a := toTLSFAllocator(al)
	//	b.ReportAllocs()
	//	b.ResetTimer()
	//	for i := 0; i < b.N; i++ {
	//		size := randomRangeSizes[i%len(randomRangeSizes)]
	//		b.SetBytes(int64(size))
	//		a.Free(a.Alloc(uintptr(size)))
	//	}
	//	after()
	//})
	//b.Run("Allocator allocZeroed", func(b *testing.B) {
	//	al := NewHeapWithConfig(1, NewSysArena(), GrowMin)
	//	a := toTLSFAllocator(al)
	//	b.ReportAllocs()
	//	b.ResetTimer()
	//	for i := 0; i < b.N; i++ {
	//		size := randomRangeSizes[i%len(randomRangeSizes)]
	//		b.SetBytes(int64(size))
	//		a.Free(a.AllocZeroed(uintptr(size)))
	//	}
	//	after()
	//})
	b.Run("Sync Heap alloc", func(b *testing.B) {
		a := NewHeapWithConfig(1, NewSysArena(), GrowMin).ToSync()
		runtime.GC()
		runtime.GC()
		b.ReportAllocs()
		b.ResetTimer()
		var before runtime.MemStats
		runtime.ReadMemStats(&before)
		for i := 0; i < b.N; i++ {
			size := randomRangeSizes[i%len(randomRangeSizes)]
			b.SetBytes(int64(size))
			a.Free(a.Alloc(uintptr(size)))
		}
		b.StopTimer()
		var after runtime.MemStats
		runtime.ReadMemStats(&after)
		doAfter(before, after)
	})
	//b.Run("Sync Heap allocZeroed", func(b *testing.B) {
	//	a := NewHeapWithConfig(1, NewSysArena(), GrowMin).ToSync()
	//	var before runtime.MemStats
	//	runtime.ReadMemStats(&before)
	//	b.ReportAllocs()
	//	b.ResetTimer()
	//	for i := 0; i < b.N; i++ {
	//		size := randomRangeSizes[i%len(randomRangeSizes)]
	//		b.SetBytes(int64(size))
	//		a.Free(a.AllocZeroed(uintptr(size)))
	//	}
	//	b.StopTimer()
	//	var after runtime.MemStats
	//	runtime.ReadMemStats(&after)
	//	doAfter(before, after)
	//})
	//b.Run("Sync Allocator alloc", func(b *testing.B) {
	//	al := NewHeapWithConfig(1, NewSysArena(), GrowMin).ToSync()
	//	a := toTLSFSyncAllocator(al)
	//	b.ReportAllocs()
	//	b.ResetTimer()
	//	for i := 0; i < b.N; i++ {
	//		size := randomRangeSizes[i%len(randomRangeSizes)]
	//		b.SetBytes(int64(size))
	//		a.Free(a.Alloc(uintptr(size)))
	//	}
	//	after()
	//})
	//b.Run("Sync Allocator allocZeroed", func(b *testing.B) {
	//	a := NewHeapWithConfig(1, NewSysArena(), GrowMin).ToSync()
	//	b.ReportAllocs()
	//	b.ResetTimer()
	//	for i := 0; i < b.N; i++ {
	//		size := randomRangeSizes[i%len(randomRangeSizes)]
	//		b.SetBytes(int64(size))
	//		a.Free(a.AllocZeroed(uintptr(size)))
	//	}
	//	after()
	//})

	b.Run("Go GC pool", func(b *testing.B) {
		runtime.GC()
		runtime.GC()
		b.ReportAllocs()
		b.ResetTimer()
		var before runtime.MemStats
		runtime.ReadMemStats(&before)
		for i := 0; i < b.N; i++ {
			size := randomRangeSizes[i%len(randomRangeSizes)]
			b.SetBytes(int64(size))
			PutBytes(GetBytes(int(size)))
		}
		b.StopTimer()
		var after runtime.MemStats
		runtime.ReadMemStats(&after)
		doAfter(before, after)
	})

	b.Run("Go GC alloc", func(b *testing.B) {
		runtime.GC()
		runtime.GC()
		b.ReportAllocs()
		b.ResetTimer()
		var before runtime.MemStats
		runtime.ReadMemStats(&before)
		for i := 0; i < b.N; i++ {
			size := randomRangeSizes[i%len(randomRangeSizes)]
			b.SetBytes(int64(size))
			_ = make([]byte, 0, size)
		}
		b.StopTimer()
		var after runtime.MemStats
		runtime.ReadMemStats(&after)
		doAfter(before, after)
	})
}

func randomRange(min, max int) int {
	return rand.Intn(max-min) + min
}

var (
	pool1 = &sync.Pool{New: func() interface{} {
		return make([]byte, 1)
	}}
	pool2 = &sync.Pool{New: func() interface{} {
		return make([]byte, 2)
	}}
	pool4 = &sync.Pool{New: func() interface{} {
		return make([]byte, 4)
	}}
	pool8 = &sync.Pool{New: func() interface{} {
		return make([]byte, 8)
	}}
	pool12 = &sync.Pool{New: func() interface{} {
		return make([]byte, 12)
	}}
	pool16 = &sync.Pool{New: func() interface{} {
		return make([]byte, 16)
	}}
	pool24 = &sync.Pool{New: func() interface{} {
		return make([]byte, 24)
	}}
	pool32 = &sync.Pool{New: func() interface{} {
		return make([]byte, 32)
	}}
	pool40 = &sync.Pool{New: func() interface{} {
		return make([]byte, 40)
	}}
	pool48 = &sync.Pool{New: func() interface{} {
		return make([]byte, 48)
	}}
	pool56 = &sync.Pool{New: func() interface{} {
		return make([]byte, 56)
	}}
	pool64 = &sync.Pool{New: func() interface{} {
		return make([]byte, 64)
	}}
	pool72 = &sync.Pool{New: func() interface{} {
		return make([]byte, 72)
	}}
	pool96 = &sync.Pool{New: func() interface{} {
		return make([]byte, 96)
	}}
	pool128 = &sync.Pool{New: func() interface{} {
		return make([]byte, 128)
	}}
	pool192 = &sync.Pool{New: func() interface{} {
		return make([]byte, 192)
	}}
	pool256 = &sync.Pool{New: func() interface{} {
		return make([]byte, 256)
	}}
	pool384 = &sync.Pool{New: func() interface{} {
		return make([]byte, 384)
	}}
	pool512 = &sync.Pool{New: func() interface{} {
		return make([]byte, 512)
	}}
	pool768 = &sync.Pool{New: func() interface{} {
		return make([]byte, 768)
	}}
	pool1024 = &sync.Pool{New: func() interface{} {
		return make([]byte, 1024)
	}}
	pool2048 = &sync.Pool{New: func() interface{} {
		return make([]byte, 2048)
	}}
	pool4096 = &sync.Pool{New: func() interface{} {
		return make([]byte, 4096)
	}}
	pool8192 = &sync.Pool{New: func() interface{} {
		return make([]byte, 8192)
	}}
	pool16384 = &sync.Pool{New: func() interface{} {
		return make([]byte, 16384)
	}}
	pool32768 = &sync.Pool{New: func() interface{} {
		return make([]byte, 32768)
	}}
	pool65536 = &sync.Pool{New: func() interface{} {
		return make([]byte, 65536)
	}}
)

func GetBytes(n int) []byte {
	v := ceilToPowerOfTwo(n)
	switch v {
	case 0, 1:
		return pool1.Get().([]byte)[:n]
	case 2:
		return pool2.Get().([]byte)[:n]
	case 4:
		return pool4.Get().([]byte)[:n]
	case 8:
		return pool8.Get().([]byte)[:n]
	case 16:
		return pool16.Get().([]byte)[:n]
	case 24:
		return pool24.Get().([]byte)[:n]
	case 32:
		return pool32.Get().([]byte)[:n]
	case 64:
		switch {
		case n < 41:
			return pool40.Get().([]byte)[:n]
		case n < 49:
			return pool48.Get().([]byte)[:n]
		case n < 57:
			return pool56.Get().([]byte)[:n]
		}
		return pool64.Get().([]byte)[:n]
	case 128:
		switch {
		case n < 73:
			return pool72.Get().([]byte)[:n]
		case n < 97:
			return pool96.Get().([]byte)[:n]
		}
		return pool128.Get().([]byte)[:n]
	case 256:
		switch {
		case n < 193:
			return pool192.Get().([]byte)[:n]
		}
		return pool256.Get().([]byte)[:n]
	case 512:
		if n <= 384 {
			return pool384.Get().([]byte)
		}
		return pool512.Get().([]byte)[:n]
	case 1024:
		if n <= 768 {
			return pool768.Get().([]byte)[:n]
		}
		return pool1024.Get().([]byte)[:n]
	case 2048:
		return pool2048.Get().([]byte)[:n]
	case 4096:
		return pool4096.Get().([]byte)[:n]
	case 8192:
		return pool8192.Get().([]byte)[:n]
	case 16384:
		return pool16384.Get().([]byte)[:n]
	case 32768:
		return pool32768.Get().([]byte)[:n]
	case 65536:
		return pool65536.Get().([]byte)[:n]
	}

	return make([]byte, n)
}

func PutBytes(b []byte) {
	switch cap(b) {
	case 1:
		pool1.Put(b)
	case 2:
		pool2.Put(b)
	case 4:
		pool4.Put(b)
	case 8:
		pool8.Put(b)
	case 12:
		pool12.Put(b)
	case 16:
		pool16.Put(b)
	case 24:
		pool24.Put(b)
	case 32:
		pool32.Put(b)
	case 40:
		pool40.Put(b)
	case 48:
		pool48.Put(b)
	case 56:
		pool56.Put(b)
	case 64:
		pool64.Put(b)
	case 72:
		pool72.Put(b)
	case 96:
		pool96.Put(b)
	case 128:
		pool128.Put(b)
	case 192:
		pool192.Put(b)
	case 256:
		pool256.Put(b)
	case 384:
		pool384.Put(b)
	case 512:
		pool512.Put(b)
	case 768:
		pool768.Put(b)
	case 1024:
		pool1024.Put(b)
	case 2048:
		pool2048.Put(b)
	case 4096:
		pool4096.Put(b)
	case 8192:
		pool8192.Put(b)
	case 16384:
		pool16384.Put(b)
	case 32768:
		pool32768.Put(b)
	case 65536:
		pool65536.Put(b)
	}
}

const (
	bitsize       = 32 << (^uint(0) >> 63)
	maxint        = int(1<<(bitsize-1) - 1)
	maxintHeadBit = 1 << (bitsize - 2)
)

// LogarithmicRange iterates from ceiled to power of two min to max,
// calling cb on each iteration.
func LogarithmicRange(min, max int, cb func(int)) {
	if min == 0 {
		min = 1
	}
	for n := ceilToPowerOfTwo(min); n <= max; n <<= 1 {
		cb(n)
	}
}

// IsPowerOfTwo reports whether given integer is a power of two.
func IsPowerOfTwo(n int) bool {
	return n&(n-1) == 0
}

// Identity is identity.
func Identity(n int) int {
	return n
}

// ceilToPowerOfTwo returns the least power of two integer value greater than
// or equal to n.
func ceilToPowerOfTwo(n int) int {
	if n&maxintHeadBit != 0 && n > maxintHeadBit {
		panic("argument is too large")
	}
	if n <= 2 {
		return n
	}
	n--
	n = fillBits(n)
	n++
	return n
}

// FloorToPowerOfTwo returns the greatest power of two integer value less than
// or equal to n.
func FloorToPowerOfTwo(n int) int {
	if n <= 2 {
		return n
	}
	n = fillBits(n)
	n >>= 1
	n++
	return n
}

func fillBits(n int) int {
	n |= n >> 1
	n |= n >> 2
	n |= n >> 4
	n |= n >> 8
	n |= n >> 16
	n |= n >> 32
	return n
}
