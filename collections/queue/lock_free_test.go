package queue

import (
	"github.com/moontrade/memory"
	"sync"
	"sync/atomic"
	"testing"
)

func TestLockFreeQueue(t *testing.T) {
	const taskNum = 50000000
	memory.AllocZeroed(24)
	q := AllocLockFreeQueue()

	b := memory.AllocBytes(24)
	b.Free()

	var wg sync.WaitGroup
	goroutines := 100
	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			for i := 0; i < taskNum/goroutines; i++ {
				task := memory.AllocBytes(24)
				q.Enqueue(task)
			}
		}()
	}

	var counter int32
	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			for {
				task := q.Dequeue()
				if !task.IsNil() {
					atomic.AddInt32(&counter, 1)
					task.Free()
				} else {
					break
				}
			}
		}()
	}

	wg.Wait()

	t.Logf("sent and received all %d tasks", taskNum)
}
