//go:build gc.provided
// +build gc.provided

package mem

import (
	"time"
	"unsafe"
)

const (
	gcAllocAny        = 0
	gcAllocChan       = 1
	gcAllocChanBuffer = 2
	gcAllocMap        = 3
	gcAllocSlice      = 4
	gcTask            = 5
	gcTaskStack       = 6
)

var collector *GC

//export gcAlloc
func gcAlloc(size uintptr, kind int32) unsafe.Pointer {
	if gc_TRACE {
		println("gcAlloc", uint(size), "kind", kind)
	}

	//switch kind {
	//case gcTaskStack:
	//	//ptr := allocator.Alloc(size)
	//	//return ptr
	//}
	//println("gcAlloc", uint(size), "kind", kind)

	//if size == 65536 || size == 40 {
	//	return allocator.Alloc(size)
	//}
	//return allocator.Alloc(size)
	//switch kind {
	//case gcAllocChan:
	//	return allocator.Alloc(size)
	//case gcAllocChanBuffer:
	//	return allocator.Alloc(size)
	//}
	ptr := collector.New(Pointer(size))
	//println("alloc ptr", uint(ptr))
	return unsafe.Pointer(ptr)
}

//go:export gcFree
func gcFree(ptr unsafe.Pointer) {
	if gc_TRACE {
		println("gcFree", uint(uintptr(ptr)))
	}
	if !collector.Free(Pointer(ptr)) {
		allocator.Free(Pointer(ptr))
	}
}

//go:export gcRun
func gcRun() {
	start := time.Now().UnixNano()

	collector.Collect()

	println("full GC", time.Now().UnixNano()-start)
}

//export KeepAlive
func gcKeepAlive(x interface{}) {
	println("gcKeepAlive")
}

//export gcSetFinalizer
func gcSetFinalizer(obj interface{}, finalizer interface{}) {
	println("gcSetFinalizer")
}

//export gcInitHeap
func gcInitHeap(heapStart, heapEnd uintptr) {
	println("gcInitHeap", uint(heapStart), uint(heapEnd))
	if allocator == nil {
		initAllocator(heapStart, heapEnd)
	}
	collector = NewGC(allocator, 64, doMarkGlobals, doMarkStack)
}

//
//export gcSetHeapEnd
func gcSetHeapEnd(newHeapEnd uintptr) {
	println("gcSetHeapEnd", uint(newHeapEnd))
}

func doMarkGlobals() {
	markGlobals()
	markScheduler()
}

//export markGlobals
func markGlobals()

//export gcMarkGlobals
func gcMarkGlobals(addr, root uintptr) {
	//println("gcMarkGlobals", uint(addr), uint(root))
	collector.markRoot(Pointer(root))
}

func doMarkStack() {
	markStack()
}

//export markStack
func markStack()

//export gcMarkRoots
func gcMarkRoots(start, end uintptr) {
	//println("gcMarkRoots", uint(start), uint(end))
	collector.markRoots(Pointer(start), Pointer(end))
}

//export gcMarkRoot
func gcMarkRoot(addr, root uintptr) {
	//println("gcMarkRoot", uint(addr), uint(root))
	collector.markRoot(Pointer(root))
}

//export markScheduler
func markScheduler()

//export gcMarkTask
func gcMarkTask(runQueuePtr, taskPtr uintptr) {
	//println("gcMarkTask", uint(allocator.HeapStart), uint(runQueuePtr), uint(taskPtr))
	collector.markRoot(Pointer(taskPtr))
}
