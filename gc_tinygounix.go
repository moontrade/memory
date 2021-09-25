//go:build tinygo && gc.provided && !tinygo.wasm && (darwin || linux || windows)
// +build tinygo
// +build gc.provided
// +build !tinygo.wasm
// +build darwin linux windows

package memory

import (
	"unsafe"
)

////////////////////////////////////////////////////////////////////////////////////
// GC Instance
////////////////////////////////////////////////////////////////////////////////////

var collector *gc

//export gcInitHeap
func gcInitHeap(heapStart, heapEnd uintptr) {
	println("gcInitHeap", uint(heapStart), uint(heapEnd))
	if allocator == nil {
		initAllocator(heapStart, heapEnd)
	}
	collector = newGC(allocator, 64, doMarkGlobals, doMarkStack)
}

////////////////////////////////////////////////////////////////////////////////////
// gcAlloc hook
////////////////////////////////////////////////////////////////////////////////////

//export gcAlloc
func gcAlloc(size uintptr) unsafe.Pointer {
	if gc_TRACE {
		println("gcAlloc", uint(size))
	}

	ptr := collector.New(Pointer(size))
	//println("alloc ptr", uint(ptr))
	return unsafe.Pointer(ptr)
}

////////////////////////////////////////////////////////////////////////////////////
// gcFree hook
////////////////////////////////////////////////////////////////////////////////////

//go:export gcFree
func gcFree(ptr unsafe.Pointer) {
	if gc_TRACE {
		println("gcFree", uint(uintptr(ptr)))
	}
	println("gcFree", uint(uintptr(ptr)))
	if !collector.Free(Pointer(ptr)) {
		allocator.Free(Pointer(ptr))
	}
}

////////////////////////////////////////////////////////////////////////////////////
// gcRun hook
////////////////////////////////////////////////////////////////////////////////////

//go:export gcRun
func gcRun() {
	//start := time.Now().UnixNano()

	collector.Collect()

	//println("full GC", time.Now().UnixNano()-start)
	collector.Print()
}

//export KeepAlive
func gcKeepAlive(x interface{}) {
	//println("gcKeepAlive")
}

//export gcSetFinalizer
func gcSetFinalizer(obj interface{}, finalizer interface{}) {
	//println("gcSetFinalizer")
}

////////////////////////////////////////////////////////////////////////////////////
// gcSetHeapEnd hook
////////////////////////////////////////////////////////////////////////////////////

//export gcSetHeapEnd
func gcSetHeapEnd(newHeapEnd uintptr) {
	//println("gcSetHeapEnd", uint(newHeapEnd))
}

////////////////////////////////////////////////////////////////////////////////////
// markGlobals hook
////////////////////////////////////////////////////////////////////////////////////

func doMarkGlobals() {
	markGlobals()
	markScheduler()
}

//export markGlobals
func markGlobals()

//export gcMarkGlobals
func gcMarkGlobals(start, end uintptr) {
	println("gcMarkGlobals", uint(start), uint(end))
	//collector.markRoots(Pointer(start), Pointer(end))
	collector.markRoot(Pointer(end))
}

////////////////////////////////////////////////////////////////////////////////////
// markStack hook
////////////////////////////////////////////////////////////////////////////////////

func doMarkStack() {
	markStack()
}

//export markStack
func markStack()

////////////////////////////////////////////////////////////////////////////////////
// gcMarkRoots hook
////////////////////////////////////////////////////////////////////////////////////

//export gcMarkRoots
func gcMarkRoots(start, end uintptr) {
	//println("gcMarkRoots", uint(start), uint(end))
	if start == 0 {
		collector.markRoot(Pointer(end))
	} else {
		if end-start < 1000000 {
			collector.markRoots(Pointer(start), Pointer(end))
		} else {
			collector.markRoot(Pointer(end))
		}
	}
}

////////////////////////////////////////////////////////////////////////////////////
// gcMarkRoot hook
////////////////////////////////////////////////////////////////////////////////////

//export gcMarkRoot
func gcMarkRoot(addr, root uintptr) {
	//println("gcMarkRoot", uint(addr), uint(root))
	collector.markRoot(Pointer(root))
}

////////////////////////////////////////////////////////////////////////////////////
// markScheduler hook
////////////////////////////////////////////////////////////////////////////////////

//export markScheduler
func markScheduler()

////////////////////////////////////////////////////////////////////////////////////
// gcMarkTask hook
////////////////////////////////////////////////////////////////////////////////////

//export gcMarkTask
func gcMarkTask(runQueuePtr, taskPtr uintptr) {
	println("gcMarkTask", uint(allocator.HeapStart), uint(runQueuePtr), uint(taskPtr))
	collector.markRoot(Pointer(taskPtr))
}