//go:build !gc.conservative && !gc.extalloc
// +build !gc.conservative,!gc.extalloc

package alloc

import (
	"github.com/moontrade/memory/mem"
	"github.com/moontrade/memory/tlsf"
	"time"
	"unsafe"
)

const (
	//gc_WHITE       uintptr = 1 << (bits.UintSize-1)
	//gc_BLACK       uintptr = 1 << (bits.UintSize-2)
	//gc_COLOR_MASK = gc_WHITE | gc_BLACK
	gc_WHITE uint32 = 0
	gc_BLACK uint32 = 1
	gc_DEBUG        = false
	gc_TRACE        = false

	// Overhead of a garbage collector object. Excludes memory manager block overhead.
	gc_ObjectOverhead = ((unsafe.Sizeof(gcObject{})) - tlsf.BlockOverhead + tlsf.ALMask) & ^tlsf.ALMask
	//gc_ObjectOverhead = unsafe.Sizeof(gcObject{}) - _TLSFBlockOverhead + unsafe.Sizeof(uintptr(0))

	// Maximum size of a garbage collector object's payload.
	gc_ObjectMaxSize = tlsf.BlockMaxSize - gc_ObjectOverhead

	// Overhead of a garbage collector object. Excludes memory manager block overhead.
	gc_TotalOverhead = tlsf.BlockOverhead + gc_ObjectOverhead
)

type GCObject uintptr

func (o GCObject) Ptr() Pointer {
	return Pointer(o)
}

// GC is a Two-Color Mark & Sweep collector on top of a Two-Level Segmented Fit (Heap)
// allocator built for TinyGo. Similar features to the internal extalloc GC in TinyGo
// except GC uses a robinhood hashset instead of a treap structure and without the need
// for a linked list. Instead, a single linear allocation is used for the hashset. Both
// colors reside in the same hashset. It should provide faster scan performance which
// becomes more noticeable as the scan size increases. For large objects, it's more ideal
// to allocate directly in the Heap allocator manually.
//
// Given the constraints of TinyGo, this is a conservative collector. However, GC
// is tuned for more manual use of the underlying Heap memory allocator. Heap is an O(1)
// time allocator and a great fit for real-time embedded systems. GC compliments
// it with a simple design and extremely quick operation for small object graphs.
//
// Large object graphs should be manually allocated and use the various tools available
// like Auto and Ref containers. GC supports a manual free as well as provided by
// the TinyGo compiler. TinyGo LLVM coroutines utilize this feature for internal coroutine
// lifecycle objects. It's quite simple to write Go programs with goroutines and channels
// that never require a GC cycle / sweep.
//
// Goal pause times are less than 10 microseconds. GC aims to complete as quickly
// as possible, but it is largely dependent on the application minimizing root scanning
// by placing manually allocated globals where possible. This effectively removes that
// graph from the marking phase.
//
// Relatively large TinyGo object graphs should still complete under 50 microseconds.
type gc struct {
	allocs      PointerSet
	first, last uintptr
	allocator   *tlsf.Heap
	markGlobals markFn
	markStack   markFn
	GCStats
}

type markFn func()

// GCStats provides all the monitoring metrics needed to see how the GC
// is operating and performing.
type GCStats struct {
	Started           int64   // Epoch in nanos when GC was first started
	Cycles            int64   // Number of times GC collect Has ran
	Live              int64   // Number of live objects
	TotalAllocs       int64   // Count of all allocations created
	TotalBytes        int64   // Sum of all allocation's size in bytes
	Frees             int64   // Count of times an allocation was freed instead of swept
	FreedBytes        int64   // Sum of all freed allocation's size in bytes
	Sweeps            int64   // Count of times an allocation was swept instead of freed
	SweepBytes        int64   // Sum of all swept allocation's size in bytes
	SweepTime         int64   // Sum of all time in nanos spent during the Sweep phase
	SweepTimeMin      int64   // Minimum time in nanos spent during a single Sweep phase
	SweepTimeMax      int64   // Maximum time in nanos spent during a single Sweep phase
	SweepTimeAvg      int64   // Average time in nanos spent during a single Sweep phase
	Roots             int64   //
	RootsMin          int64   //
	RootsMax          int64   //
	RootsTimeMin      int64   //
	RootsTimeMax      int64   //
	RootsTimeAvg      int64   //
	GraphDepth        int64   //
	GraphMinDepth     int64   //
	GraphMaxDepth     int64   //
	GraphAvgDepth     int64   //
	GraphTimeMin      int64   //
	GraphTimeMax      int64   //
	GraphTimeAvg      int64   //
	TotalTime         int64   // Sum of all time in nanos spent doing GC collect
	MinTime           int64   // Minimum time in nanos spent during a single GC collect
	MaxTime           int64   // Maximum time in nanos spent during a single GC collect
	AvgTime           int64   // Average time in nanos spent during a single GC collect
	LastMarkRootsTime int64   // Time in nanos spent during the most recent GC collect "Mark Roots" phase
	LastMarkGraphTime int64   // Time in nanos spent during the most recent GC collect "Mark Graph" phase
	LastSweepTime     int64   // Time in nanos spent during the most recent GC collect "Sweep" phase
	LastGCTime        int64   // Time in nanos spent during the most recent GC collect
	LastSweep         int64   // Number of allocations that were swept during the most recent GC collect "Sweep" phase
	LastSweepBytes    int64   // Number of bytes reclaimed during the most recent GC collect "Sweep" phase
	LiveBytes         uintptr // Sum of all live allocation's size in bytes
}

func (s *GCStats) Print() {
	println("GC cycle")
	println("\tlive:				", uint(s.Live))
	println("\tlive bytes:			", uint(s.LiveBytes))
	println("\tfrees:				", uint(s.Frees))
	println("\tallocs:				", uint(s.TotalAllocs))
	println("\tfreed bytes:		", uint(s.FreedBytes))
	println("\tsweep bytes:		", uint(s.SweepBytes))
	println("\ttotal bytes:		", uint(s.TotalBytes))
	println("\tlast sweep:			", uint(s.LastSweep))
	println("\tlast sweep bytes:	", uint(s.LastSweepBytes))
	println("\tlast mark time:		", toMicros(s.LastMarkRootsTime), microsSuffix)
	println("\tlast graph time:	", toMicros(s.LastMarkGraphTime), microsSuffix)
	println("\tlast sweep time:	", toMicros(s.LastSweepTime), microsSuffix)
	println("\tlast GC time:		", toMicros(s.LastGCTime), microsSuffix)
}

func PrintDebugInfo() {
	println("gc_ObjectOverhead	", uint(gc_ObjectOverhead))
	println("gc_ObjectMaxSize		", uint(gc_ObjectMaxSize))
	println("gc_TotalOverhead		", uint(gc_TotalOverhead))
	tlsf.PrintDebugInfo()
}

//goland:noinspection ALL
func newGC(
	allocator *tlsf.Heap,
	initialCapacity uintptr,
	markGlobals, markStack markFn,
) *gc {
	gc := (*gc)(unsafe.Pointer(allocator.Alloc(unsafe.Sizeof(gc{}))))
	gc.allocator = allocator
	gc.allocs = NewPointerSet(allocator, initialCapacity)
	gc.first = ^uintptr(0)
	gc.last = 0
	gc.markGlobals = markGlobals
	gc.markStack = markStack
	gc.Started = time.Now().UnixNano()
	return gc
}

// GCObject Represents a managed object in memory, consisting of a header followed by the object's data.
type gcObject struct {
	tlsf.BLOCK
	color  uint32 // alloc.Pointer to the next object with color flags stored in the alignment bits.
	rtSize uint32 // Runtime size.
}

// Gets the size of this object in memory.
func (o *gcObject) size() uintptr {
	return tlsf.BlockOverhead + (o.MMInfo & ^uintptr(3))
}

func (gc *gc) Allocator() *tlsf.Heap {
	return gc.allocator
}

// MarkRoot marks a single pointer as a root
//goland:noinspection ALL
func (gc *gc) markRoot(root uintptr) {
	//root = root.Add(-int(gc_TotalOverhead))
	if root < gc.first || root > gc.last {
		return
	}
	root -= gc_TotalOverhead
	if gc.allocs.Has(root) {
		// Mark as gc_BLACK
		(*(*gcObject)(unsafe.Pointer(root))).color = gc_BLACK
	}
}

// MarkRoots scans a block of contiguous memory for root pointers.
//goland:noinspection ALL
func (gc *gc) markRoots(start, end uintptr) {
	if gc_TRACE {
		println("MarkRoots", uint(start), uint(end))
	}

	if end <= start {
		return
	}
	if start == 0 || end == 0 {
		return
	}

	// Adjust to keep within range GC range
	//println("MarkRoots", uint(start), uint(end))

	// Align start and end pointers.
	start = (uintptr(start) + unsafe.Alignof(unsafe.Pointer(nil)) - 1) &^ (unsafe.Alignof(unsafe.Pointer(nil)) - 1)
	end &^= unsafe.Alignof(unsafe.Pointer(nil)) - 1

	// Mark all pointers.
	for ptr := start; ptr < end; ptr += unsafe.Alignof(unsafe.Pointer(nil)) {
		p := *(*uintptr)(unsafe.Pointer(ptr))
		gc.markRoot(p)
	}
}

//goland:noinspection ALL
func (gc *gc) markRecursive(root uintptr, depth int) {
	// Are we too deep?
	if depth > 64 {
		return
	}

	if gc_TRACE {
		println("markRecursive", uint(root), "depth", depth)
	}
	obj := (*gcObject)(unsafe.Pointer(root))
	if obj.color == gc_WHITE {
		obj.color = gc_BLACK

		if gc_TRACE {
			println(uint(root), "color", obj.color, "rtSize", obj.rtSize, "size", uint(obj.size()))
		}
		if uintptr(obj.rtSize)%unsafe.Sizeof(uintptr(0)) != 0 {
			return
		}
		start := root + gc_TotalOverhead
		end := start + uintptr(obj.rtSize)
		//start = (start + unsafe.Alignof(unsafe.Pointer(nil)) - 1) &^ (unsafe.Alignof(unsafe.Pointer(nil)) - 1)
		//end &^= unsafe.Alignof(unsafe.Pointer(nil)) - 1

		for ptr := start; ptr < end; ptr += unsafe.Alignof(unsafe.Pointer(nil)) {
			p := *(*uintptr)(unsafe.Pointer(ptr)) - gc_TotalOverhead
			if p < gc.first || p > gc.last {
				continue
			}
			if !gc.allocs.Has(p) {
				continue
			}
			gc.markRecursive(p, depth+1)
		}
	}
}

//goland:noinspection ALL
func (gc *gc) markGraph(root uintptr) {
	var (
		obj   = (*gcObject)(unsafe.Pointer(root))
		start = root + gc_TotalOverhead
		end   = start + uintptr(obj.rtSize)
	)

	// unaligned allocation must be some sort of string or data buffer. skip it.
	if uintptr(obj.rtSize)%unsafe.Sizeof(uintptr(0)) != 0 {
		return
	}

	//pointersToCount := (uint(end) - uint(start)) / uint(unsafe.Sizeof(unsafe.Pointer(nil)))
	//if pointersToCount > 4096 {
	//	//println("huge object found", uint(pointersToCount), "pointers to scan")
	//	return
	//}

	// Mark all pointers.
	for ptr := start; ptr < end; ptr += unsafe.Alignof(unsafe.Pointer(nil)) {
		p := *(*uintptr)(unsafe.Pointer(ptr)) - gc_TotalOverhead
		if p < gc.first || p > gc.last {
			continue
		}
		if !gc.allocs.Has(p) {
			continue
		}
		gc.markRecursive(p, 0)
	}
}

// New allocates a new GC GCObject
//goland:noinspection ALL
func (gc *gc) New(size uintptr) uintptr {
	// Is the size too large?
	if size > uintptr(gc_ObjectMaxSize) {
		panic("allocation too large")
	}

	// Allocate memory
	obj := (*gcObject)(unsafe.Pointer((gc.allocator.AllocZeroed(uintptr(gc_ObjectOverhead) + size)) - tlsf.BlockOverhead))
	if obj == nil {
		return 0
	}

	// Add the runtime size and Add to gc_WHITE
	obj.rtSize = uint32(size)
	obj.color = gc_WHITE
	gc.LiveBytes += obj.size()
	gc.TotalBytes += int64(obj.size())
	gc.Live++
	gc.TotalAllocs++

	// Convert to uint pointer
	ptr := uintptr(unsafe.Pointer(obj))

	// Zero out the allocation
	mem.Zero(unsafe.Pointer(ptr+gc_TotalOverhead), uintptr(size))

	// Add to allocations map
	gc.allocs.Add(ptr, 0)

	// Update first pointer if necessary
	if ptr < gc.first {
		gc.first = ptr
	}
	// Update last pointer if necessary
	if ptr > gc.last {
		gc.last = ptr
	}

	ptr += gc_TotalOverhead
	//println("New", uint(ptr))

	// Return pointer to data
	return ptr
}

// Free will immediately remove the GC GCObject and free up the memory in the allocator.
//goland:noinspection ALL
func (gc *gc) Free(ptr uintptr) bool {
	p := ptr - gc_TotalOverhead
	//if !gc.allocs.Has(p) {
	if _, ok := gc.allocs.Delete(p); !ok {
		return false
	}

	if gc_TRACE {
		println("GC free", uint(ptr))
	}

	obj := (*gcObject)(unsafe.Pointer(p))
	size := obj.size()
	gc.LiveBytes -= size
	gc.FreedBytes += int64(size)
	gc.Live--
	gc.Frees++
	//gc.allocs.Delete(p)

	//println("GC free", uint(uintptr(ptr)), "size", uint(size), "rtSize", obj.rtSize)
	gc.allocator.Free(p + tlsf.BlockOverhead)

	return true
}

//goland:noinspection ALL
func (gc *gc) Collect() {
	if gc_TRACE {
		println("GC collect started...")
	}
	//println("tcmsCollect")
	var (
		start = time.Now().UnixNano()
		k     uintptr
		obj   *gcObject
		first = ^uintptr(0)
		last  = uintptr(0)
	)
	gc.Cycles++

	////////////////////////////////////////////////////////////////////////
	// Mark Roots Phase
	////////////////////////////////////////////////////////////////////////
	if gc.markStack != nil {
		gc.markStack()
	}
	if gc.markGlobals != nil {
		gc.markGlobals()
	}
	// End of mark roots
	end := time.Now().UnixNano()
	markTime := end - start

	////////////////////////////////////////////////////////////////////////
	// Mark Graph Phase
	////////////////////////////////////////////////////////////////////////
	start = end
	gc.LastSweep = 0
	gc.LastSweepBytes = 0
	var (
		items     = gc.allocs.items
		itemsSize = gc.allocs.size
		itemsEnd  = items + (itemsSize * unsafe.Sizeof(pointerSetItem{}))
	)
	for ; items < itemsEnd; items += unsafe.Sizeof(pointerSetItem{}) {
		k = *(*uintptr)(unsafe.Pointer(items))
		if k == 0 {
			continue
		}
		gc.markGraph(k)
	}

	// End of mark graph
	end = time.Now().UnixNano()
	markGraphTime := end - start

	////////////////////////////////////////////////////////////////////////
	// Sweep Phase
	////////////////////////////////////////////////////////////////////////
	start = markGraphTime + start

	// Reset items iterator
	items = gc.allocs.items
	itemsSize = gc.allocs.size
	itemsEnd = items + (itemsSize * unsafe.Sizeof(pointerSetItem{}))
	for ; items < itemsEnd; items += unsafe.Sizeof(pointerSetItem{}) {
		// dereference pointer
		k = *(*uintptr)(unsafe.Pointer(items))
		// Empty item?
		if k == 0 {
			continue
		}
		// cast to object
		obj = (*gcObject)(unsafe.Pointer(k))
		// free all gc_WHITE objects
		if obj.color == gc_WHITE {
			gc.LiveBytes -= obj.size()
			gc.LastSweepBytes += int64(obj.size())
			gc.Live--
			gc.LastSweep++

			if gc_TRACE {
				println("GC sweep", uint(k), "size", uint(obj.size()))
			}

			println("GC sweep", uint(uintptr(k)+gc_TotalOverhead), "size", uint(obj.size()), "rtSize", obj.rtSize)

			// Free memory
			gc.allocator.Free(k + tlsf.BlockOverhead)

			// Remove from alloc map
			gc.allocs.Delete(k)
			//items -= unsafe.Sizeof(pointerSetItem{})
		} else { // turn all gc_BLACK objects into gc_WHITE objects
			//k += gc_TotalOverhead
			if k < first {
				first = k
			}
			if k > last {
				last = k
			}
			if gc_TRACE {
				//println("GC retained", uint(k), "size", uint(obj.size()))
			}
			obj.color = gc_WHITE
		}
	}

	gc.first = first
	gc.last = last
	end = time.Now().UnixNano()
	sweepTime := end - start
	gc.LastMarkRootsTime = markTime
	gc.LastMarkGraphTime = markGraphTime
	gc.LastSweepTime = sweepTime
	gc.SweepTime += sweepTime
	gc.SweepBytes += gc.LastSweepBytes
	gc.Sweeps += gc.LastSweep
	gc.LastGCTime = markTime + markGraphTime + sweepTime
	gc.TotalTime += gc.LastGCTime

	if gc_TRACE {
		println("GC collect finished")
	}
	//stats.Print()
}
