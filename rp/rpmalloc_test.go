package rp

import (
	"runtime"
	"sync"
	"testing"
)

func TestAlloc(t *testing.T) {
	//InitThread()
	a := Alloc(24)
	Free(a)

	for i := 0; i < 100; i++ {
		go func() {
			InitThread()
			Free(Alloc(32))
		}()
	}
}

func BenchmarkCAlloc(b *testing.B) {
	var (
		goroutines = runtime.NumCPU()
	)

	b.Run("rp", func(b *testing.B) {
		init := sync.RWMutex{}
		init.Lock()
		start := &sync.WaitGroup{}
		start.Add(1)
		wg := &sync.WaitGroup{}
		wg.Add(goroutines)

		for g := 0; g < goroutines; g++ {
			go func() {
				defer wg.Done()
				init.RLock()
				start.Wait()
				runAllocs(32, int32(b.N))
				init.RUnlock()
			}()
		}

		init.Unlock()
		b.ResetTimer()
		b.ReportAllocs()
		start.Done()
		//InitThread()
		wg.Wait()
	})
	b.Run("rpZeroed", func(b *testing.B) {
		//InitThread()
		b.ResetTimer()
		b.ReportAllocs()
		runAllocZeroed(32, int32(b.N))
	})
}
