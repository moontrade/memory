//go:build !tinygo
// +build !tinygo

package queue

import (
	mem "github.com/moontrade/memory"
	"sync/atomic"
	"unsafe"
)

// LockFree is a simple non-blocking and concurrent queue.
type LockFree struct {
	head   uintptr
	tail   uintptr
	alloc  mem.Allocator
	length int32
}

type lockFreeNode struct {
	value mem.Bytes
	next  uintptr
	alloc mem.Allocator
}

//goland:noinspection GoVetUnsafePointer
func AllocLockFreeQueue(a mem.Allocator) *LockFree {
	q := (*LockFree)(unsafe.Pointer(a.AllocZeroed(mem.Pointer(unsafe.Sizeof(LockFree{})))))
	n := a.AllocZeroed(mem.Pointer(unsafe.Sizeof(lockFreeNode{})))
	q.head = uintptr(n)
	q.tail = uintptr(n)
	q.alloc = a
	(*lockFreeNode)(unsafe.Pointer(n)).alloc = a
	return q
}

// Enqueue puts the given value v at the tail of the queue.
//goland:noinspection GoVetUnsafePointer
func (q *LockFree) Enqueue(task mem.Bytes) {
	n := uintptr(task.Allocator().AllocZeroed(mem.Pointer(unsafe.Sizeof(lockFreeNode{}))))
	node := (*lockFreeNode)(unsafe.Pointer(n))
	node.value = task
	node.alloc = task.Allocator()
retry:
	last := atomic.LoadUintptr(&q.tail)
	lastV := (*lockFreeNode)(unsafe.Pointer(last))
	next := atomic.LoadUintptr(&lastV.next)
	// Are tail and next consistent?
	if last == atomic.LoadUintptr(&q.tail) {
		if next == 0 {
			// Try to link node at the end of the linked list.
			if atomic.CompareAndSwapUintptr(&lastV.next, next, n) { // enqueue is done.
				// Try to swing tail to the inserted node.
				atomic.CompareAndSwapUintptr(&q.tail, last, n)
				atomic.AddInt32(&q.length, 1)
				return
			}
		} else { // tail was not pointing to the last node
			// Try to swing tail to the next node.
			atomic.CompareAndSwapUintptr(&q.tail, last, next)
		}
	}
	goto retry
}

// Dequeue removes and returns the value at the head of the queue.
// It returns nil if the queue is empty.
//goland:noinspection GoVetUnsafePointer
func (q *LockFree) Dequeue() mem.Bytes {
retry:
	first := atomic.LoadUintptr(&q.head)
	firstV := (*lockFreeNode)(unsafe.Pointer(first))
	last := atomic.LoadUintptr(&q.tail)
	next := atomic.LoadUintptr(&firstV.next)
	// Are first, tail, and next consistent?
	if first == atomic.LoadUintptr(&q.head) {
		// Is queue empty or tail falling behind?
		if first == last {
			// Is queue empty?
			if next == 0 {
				//println("empty")
				return mem.Bytes{}
			}
			//println("first == tail")
			atomic.CompareAndSwapUintptr(&q.tail, last, next) // tail is falling behind, try to advance it.
		} else {
			// Read value before CAS, otherwise another dequeue might free the next node.
			task := (*lockFreeNode)(unsafe.Pointer(next)).value
			if atomic.CompareAndSwapUintptr(&q.head, first, next) { // dequeue is done, return value.
				atomic.AddInt32(&q.length, -1)
				firstV.alloc.Free(mem.Pointer(first))
				return task
			}
		}
	}
	goto retry
}

// Empty indicates whether this queue is empty or not.
func (q *LockFree) Empty() bool {
	return atomic.LoadInt32(&q.length) == 0
}

//func load(p *uintptr) uintptr {
//	return (*node)(unsafe.Pointer(atomic.LoadUintptr(p)))
//}
//
//func cas(p *uintptr, old, new uintptr) bool {
//	return atomic.CompareAndSwapUintptr(p, old, new)
//}
