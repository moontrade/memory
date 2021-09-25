package queue

import (
	"github.com/moontrade/memory"
	"sync"
	"sync/atomic"
	"testing"
)

func TestLockFreeQueue(t *testing.T) {
	const taskNum = 10000
	memory.AllocZeroed(24)
	q := AllocLockFreeQueue()

	b := memory.AllocBytes(24)
	b.Free()

	var wg sync.WaitGroup
	wg.Add(4)
	go func() {
		for i := 0; i < taskNum; i++ {
			task := memory.AllocBytes(24)
			q.Enqueue(task)
		}
		wg.Done()
	}()
	go func() {
		for i := 0; i < taskNum; i++ {
			task := memory.AllocBytes(32)
			q.Enqueue(task)
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
