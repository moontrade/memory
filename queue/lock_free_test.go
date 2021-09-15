package queue

import (
	mem "github.com/moontrade/memory"
	"sync"
	"sync/atomic"
	"testing"
	"unsafe"
)

func TestLockFreeQueue(t *testing.T) {
	const taskNum = 10000
	a := mem.NewTLSF(1000).ToSync()
	a2 := mem.NewTLSF(1000).ToSync()
	_ = a2
	a.AllocNotCleared(24)
	q := AllocLockFreeQueue(mem.Allocator(unsafe.Pointer(a)))

	b := a.Bytes(24)
	b.Free()

	var wg sync.WaitGroup
	wg.Add(4)
	go func() {
		for i := 0; i < taskNum; i++ {
			task := a.Bytes(24)
			q.Enqueue(a, task)
		}
		wg.Done()
	}()
	go func() {
		for i := 0; i < taskNum; i++ {
			task := a.Bytes(32)
			q.Enqueue(a, task)
		}
		wg.Done()
	}()

	var counter int32
	go func() {
		for {
			task := q.Dequeue()
			if !task.IsNil() {
				atomic.AddInt32(&counter, 1)
				task.Free()
			} else if task.IsNil() && atomic.LoadInt32(&counter) == 2*taskNum {
				break
			}

			//count := atomic.LoadInt32(&counter)
			//if count % 100 == 0 {
			//	println("counter", count)
			//}
		}
		wg.Done()
	}()
	go func() {
		for {
			task := q.Dequeue()
			if !task.IsNil() {
				atomic.AddInt32(&counter, 1)
				task.Free()
			} else if task.IsNil() && atomic.LoadInt32(&counter) == 2*taskNum {
				break
			}

			//count := atomic.LoadInt32(&counter)
			//if count % 100 == 0 {
			//	println("counter", count)
			//}
		}
		wg.Done()
	}()
	wg.Wait()

	t.Logf("sent and received all %d tasks", 2*taskNum)
}
