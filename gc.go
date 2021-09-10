//go:build !gc.conservative && !gc.extalloc
// +build !gc.conservative,!gc.extalloc

package mem

import (
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
	gc_OBJECT_OVERHEAD = (unsafe.Sizeof(gcObject{}) - _TLSFBlockOverhead + _TLSFALMask) & ^_TLSFALMask
	//gc_OBJECT_OVERHEAD = unsafe.Sizeof(gcObject{}) - _TLSFBlockOverhead + unsafe.Sizeof(uintptr(0))

	// Maximum size of a garbage collector object's payload.
	gc_OBJECT_MAXSIZE = _TLSFBlockMaxSize - gc_OBJECT_OVERHEAD

	// Overhead of a garbage collector object. Excludes memory manager block overhead.
	gc_TOTAL_OVERHEAD = _TLSFBlockOverhead + gc_OBJECT_OVERHEAD
)

// GC is a Two-Color Mark & Sweep collector on top of a Two-Level Segmented Fit (TLSF)
// allocator. Similar features to the internal extalloc GC in TinyGo except GC uses a
// robinhood hashset instead of a treap structure and without the need for a linked list.
// Instead, a single linear allocation is used for the hashset. Both colors reside in
// the same hashset.
//
// Given the constraints of TinyGo, this is a conservative collector. However, GC
// is tuned for more manual use of the underlying TLSF memory allocator. TLSF is an O(1)
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
type GC struct {
	allocs      PointerSet
	first, last uintptr
	allocator   *Allocator
	markGlobals MarkFn
	markStack   MarkFn
	GCStats
}

type MarkFn func(mark func(root uintptr), markRoots func(start, end uintptr))

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

func GCPrintDebug() {
	println("gc_OBJECT_OVERHEAD	", uint(gc_OBJECT_OVERHEAD))
	println("gc_OBJECT_MAXSIZE		", uint(gc_OBJECT_MAXSIZE))
	println("gc_TOTAL_OVERHEAD		", uint(gc_TOTAL_OVERHEAD))
	AllocatorPrintDebugInfo()
}

//goland:noinspection ALL
func NewGC(
	allocator *Allocator,
	initialCapacity uintptr,
	markGlobals, markStack MarkFn,
) *GC {
	gc := (*GC)(allocator.Alloc(unsafe.Sizeof(GC{})))
	gc.allocator = allocator
	gc.allocs = NewPointerSet(allocator, initialCapacity)
	gc.first = ^uintptr(0)
	gc.last = 0
	gc.markGlobals = markGlobals
	gc.markStack = markStack
	gc.Started = time.Now().UnixNano()
	return gc
}

// Object Represents a managed object in memory, consisting of a header followed by the object's data.
type gcObject struct {
	tlsfBLOCK
	color  uint32 // Pointer to the next object with color flags stored in the alignment bits.
	rtSize uint32 // Runtime size.
}

// Gets the size of this object in memory.
func (o *gcObject) size() uintptr {
	return _TLSFBlockOverhead + (o.mmInfo & ^uintptr(3))
}

func (gc *GC) Allocator() *Allocator {
	return gc.allocator
}

// MarkRoot marks a single pointer as a root
//goland:noinspection ALL
func (gc *GC) markRoot(root uintptr) {
	root -= gc_TOTAL_OVERHEAD
	if root < gc.first || root > gc.last {
		return
	}
	if gc.allocs.Has(root) {
		// Mark as BLACK
		(*(*gcObject)(unsafe.Pointer(root))).color = gc_BLACK
	}
}

// MarkRoots scans a block of contiguous memory for root pointers.
//goland:noinspection ALL
func (gc *GC) markRoots(start, end uintptr) {
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
	if start < gc.first {
		if gc.first >= end {
			return
		}
		start = gc.first
	}
	if end > gc.last {
		end = gc.last
	}

	// Align start and end pointers.
	start = (start + unsafe.Alignof(unsafe.Pointer(nil)) - 1) &^ (unsafe.Alignof(unsafe.Pointer(nil)) - 1)
	end &^= unsafe.Alignof(unsafe.Pointer(nil)) - 1

	// Mark all pointers.
	for ptr := start; ptr < end; ptr += unsafe.Alignof(unsafe.Pointer(nil)) {
		p := *(*uintptr)(unsafe.Pointer(ptr))
		gc.markRoot(p)
	}
}

//goland:noinspection ALL
func (gc *GC) markRecursive(root uintptr, depth int) {
	root -= gc_TOTAL_OVERHEAD
	if !gc.allocs.Has(root) {
		return
	}
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
		start := root + gc_TOTAL_OVERHEAD
		end := start + uintptr(obj.rtSize)
		//start = (start + unsafe.Alignof(unsafe.Pointer(nil)) - 1) &^ (unsafe.Alignof(unsafe.Pointer(nil)) - 1)
		//end &^= unsafe.Alignof(unsafe.Pointer(nil)) - 1

		for ptr := start; ptr < end; ptr += unsafe.Alignof(unsafe.Pointer(nil)) {
			p := *(*uintptr)(unsafe.Pointer(ptr))
			if (p < gc.first || p > gc.last) || (p >= start && p < end) {
				continue
			}
			if !gc.allocs.Has(p - gc_TOTAL_OVERHEAD) {
				continue
			}
			gc.markRecursive(p, depth+1)
		}
	}
}

//goland:noinspection ALL
func (gc *GC) markGraph(root uintptr) {
	var (
		obj   = (*gcObject)(unsafe.Pointer(root))
		start = root + gc_TOTAL_OVERHEAD
		end   = start + uintptr(obj.rtSize)
	)

	// unaligned allocation must be some sort of string or data buffer. skip it.
	if uintptr(obj.rtSize)%unsafe.Sizeof(uintptr(0)) != 0 {
		return
	}

	pointersToCount := (uint(end) - uint(start)) / uint(unsafe.Sizeof(unsafe.Pointer(nil)))
	if pointersToCount > 128 {
		//return
	}

	// Mark all pointers.
	for ptr := start; ptr < end; ptr += unsafe.Alignof(unsafe.Pointer(nil)) {
		p := *(*uintptr)(unsafe.Pointer(ptr))
		if (p < gc.first || p > gc.last) || (p >= start && p < end) {
			continue
		}
		if !gc.allocs.Has(p - gc_TOTAL_OVERHEAD) {
			continue
		}
		gc.markRecursive(p, 0)
	}
}

// New allocates a new GC Object
//goland:noinspection ALL
func (gc *GC) New(size uintptr) unsafe.Pointer {
	// Is the size too large?
	if size > gc_OBJECT_MAXSIZE {
		panic("allocation too large")
	}

	// Allocate memory
	obj := (*gcObject)(unsafe.Pointer(uintptr(gc.allocator.Alloc(gc_OBJECT_OVERHEAD+size)) - _TLSFBlockOverhead))
	if obj == nil {
		return nil
	}

	// Add the runtime size and Add to WHITE
	obj.rtSize = uint32(size)
	obj.color = gc_WHITE
	gc.LiveBytes += obj.size()
	gc.TotalBytes += int64(obj.size())
	gc.Live++
	gc.TotalAllocs++

	// Convert to uint pointer
	ptr := uintptr(unsafe.Pointer(obj))

	// Zero out the allocation
	memzero(unsafe.Pointer(ptr+gc_TOTAL_OVERHEAD), size)

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

	// Return pointer to data
	return unsafe.Pointer(ptr + gc_TOTAL_OVERHEAD)
}

// Free will immediately remove the GC Object and free up the memory in the allocator.
//goland:noinspection ALL
func (gc *GC) Free(ptr unsafe.Pointer) bool {
	p := uintptr(ptr) - gc_TOTAL_OVERHEAD
	if !gc.allocs.Has(p) {
		return false
	}

	if gc_TRACE {
		println("GC free", uint(uintptr(ptr)))
	}

	obj := (*gcObject)(unsafe.Pointer(p))
	size := obj.size()
	gc.LiveBytes -= size
	gc.FreedBytes += int64(size)
	gc.Live--
	gc.Frees++
	gc.allocs.Delete(p)

	println("GC free", uint(uintptr(ptr)), "size", uint(size), "rtSize", obj.rtSize)
	gc.allocator.Free(unsafe.Pointer(p + _TLSFBlockOverhead))

	return true
}

//goland:noinspection ALL
func (gc *GC) Collect() {
	if gc_TRACE {
		println("moon GC collect started...")
	}
	//println("tcmsCollect")
	var (
		start = time.Now().UnixNano()
		k     uintptr
		obj   *gcObject
		min   = ^uintptr(0)
		max   = uintptr(0)
	)
	gc.Cycles++

	////////////////////////////////////////////////////////////////////////
	// Mark Roots Phase
	////////////////////////////////////////////////////////////////////////
	if gc.markStack != nil {
		gc.markStack(gc.markRoot, gc.markRoots)
	}
	if gc.markGlobals != nil {
		gc.markGlobals(gc.markRoot, gc.markRoots)
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
		k = *(*uintptr)(unsafe.Pointer(items))
		// Empty item?
		if k == 0 {
			continue
		}
		obj = (*gcObject)(unsafe.Pointer(k))
		if obj.color == gc_WHITE {
			gc.LiveBytes -= obj.size()
			gc.LastSweepBytes += int64(obj.size())
			gc.Live--
			gc.LastSweep++

			if gc_TRACE {
				println("GC sweep", uint(k), "size", uint(obj.size()))
			}

			//println("GC sweep", uint(k+gc_TOTAL_OVERHEAD), "size", uint(obj.size()), "rtSize", obj.rtSize)

			// Free memory
			gc.allocator.Free(unsafe.Pointer(k + _TLSFBlockOverhead))

			// Remove from alloc map
			gc.allocs.Delete(k)
			//items -= unsafe.Sizeof(pointerSetItem{})
		} else {
			if k < min {
				min = k
			}
			if k > max {
				max = k
			}
			if gc_TRACE {
				//println("GC retained", uint(k), "size", uint(obj.size()))
			}
			obj.color = gc_WHITE
		}
	}

	gc.first = min
	gc.last = max
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
		println("moon GC collect finished")
	}
	//stats.Print()
}

// PointerSet is a hashset that uses the robinhood algorithm. This
// implementation is not concurrent safe.
type PointerSet struct {
	// items are the slots of the hashmap for items.
	items uintptr
	end   uintptr
	size  uintptr

	// Number of keys in the PointerSet.
	count     uintptr
	allocator *Allocator
	// When any item's distance gets too large, Grow the PointerSet.
	// Defaults to 10.
	maxDistance int32
	growBy      int32
	// Number of hash slots to Grow by
	growthFactor float32
}

// Item represents an entry in the PointerSet.
type pointerSetItem struct {
	key      uintptr
	distance int32 // How far item is from its best position.
}

// pointerSetHash uses the fnv hashing algorithm for 32bit integers.
// fnv was chosen for its consistent low collision rate even with tight monotonic numbers (WASM)
// and for its performance. Other potential candidates are wyhash, metro, and adler32. Each of these
// have optimized 32bit version in hash.go in this package.
//go:inline
var pointerSetHash = fnv32

//func pointerSetHash(v uint32) uint32 {
//	return fnv32(v)
//}

// NewPointerSet returns a new robinhood hashmap.
//goland:noinspection ALL
func NewPointerSet(allocator *Allocator, size uintptr) PointerSet {
	items := allocator.Alloc(unsafe.Sizeof(pointerSetItem{}) * size)
	memzero(items, unsafe.Sizeof(pointerSetItem{})*size)
	return PointerSet{
		items:        uintptr(items),
		size:         size,
		end:          uintptr(items) + (size * unsafe.Sizeof(pointerSetItem{})),
		allocator:    allocator,
		maxDistance:  10,
		growBy:       64,
		growthFactor: 2.0,
	}
}

//goland:noinspection GoVetUnsafePointer
func (ps *PointerSet) Close() error {
	if ps == nil {
		return nil
	}
	if ps.items == 0 {
		return nil
	}
	ps.allocator.Free(unsafe.Pointer(ps.items))
	ps.items = 0
	return nil
}

// Reset clears PointerSet, where already allocated memory will be reused.
//goland:noinspection ALL
func (ps *PointerSet) Reset() {
	memzero(unsafe.Pointer(ps.items), unsafe.Sizeof(pointerSetItem{})*ps.size)
	ps.count = 0
}

//goland:noinspection GoVetUnsafePointer
func (ps *PointerSet) isCollision(k uintptr) bool {
	return *(*uintptr)(unsafe.Pointer(
		ps.items + (uintptr(pointerSetHash(uint32(k))%uint32(ps.size)) * unsafe.Sizeof(pointerSetItem{})))) != 0
}

// Has returns whether the key exists in the Add.
//goland:noinspection ALL
func (ps *PointerSet) Has(k uintptr) bool {
	var (
		idx      = ps.items + (uintptr(pointerSetHash(uint32(k))%uint32(ps.size)) * unsafe.Sizeof(pointerSetItem{}))
		idxStart = idx
	)
	for {
		entry := *(*uintptr)(unsafe.Pointer(idx))
		if entry == 0 {
			return false
		}
		if entry == k {
			return true
		}
		idx += unsafe.Sizeof(pointerSetItem{})
		if idx >= ps.end {
			idx = ps.items
		}
		// Went all the way around?
		if idx == idxStart {
			return false
		}
	}
}

func (ps *PointerSet) Set(k uintptr) (bool, bool) {
	return ps.Add(k, 0)
}

// Add inserts or updates a key into the PointerSet. The returned
// wasNew will be true if the mutation was on a newly seen, inserted
// key, and wasNew will be false if the mutation was an update to an
// existing key.
//goland:noinspection GoVetUnsafePointer
func (ps *PointerSet) Add(k uintptr, depth int) (bool, bool) {
	var (
		idx      = ps.items + (uintptr(pointerSetHash(uint32(k))%uint32(ps.size)) * unsafe.Sizeof(pointerSetItem{}))
		idxStart = idx
		incoming = pointerSetItem{k, 0}
	)
	for {
		entry := (*pointerSetItem)(unsafe.Pointer(idx))
		if entry.key == 0 {
			*(*pointerSetItem)(unsafe.Pointer(idx)) = incoming
			ps.count++
			return true, true
		}

		if entry.key == incoming.key {
			entry.distance = incoming.distance
			return false, true
		}

		// Swap if the incoming item is further from its best idx.
		if entry.distance < incoming.distance {
			incoming, *(*pointerSetItem)(unsafe.Pointer(idx)) = *(*pointerSetItem)(unsafe.Pointer(idx)), incoming
		}

		// One step further away from best idx.
		incoming.distance++

		idx += unsafe.Sizeof(pointerSetItem{})
		if idx >= ps.end {
			idx = ps.items
		}

		// Grow if distances become big or we went all the way around.
		if incoming.distance > ps.maxDistance || idx == idxStart {
			if depth > 5 {
				return false, false
			}
			if !ps.Grow() {
				return false, false
			}
			return ps.Add(incoming.key, depth+1)
		}
	}
}

// Delete removes a key from the PointerSet.
//goland:noinspection GoVetUnsafePointer
func (ps *PointerSet) Delete(k uintptr) (uintptr, bool) {
	if k == 0 {
		return 0, false
	}

	var (
		idx      = ps.items + (uintptr(pointerSetHash(uint32(k))%uint32(ps.size)) * unsafe.Sizeof(pointerSetItem{}))
		idxStart = idx
	)
	for {
		entry := (*pointerSetItem)(unsafe.Pointer(idx))
		if entry.key == 0 {
			return 0, false
		}
		if entry.key == k {
			break // Found the item.
		}
		idx += unsafe.Sizeof(pointerSetItem{})
		if idx >= ps.end {
			idx = ps.items
		}
		if idx == idxStart {
			return 0, false
		}
	}
	// Left-shift succeeding items in the linear chain.
	for {
		next := idx + unsafe.Sizeof(pointerSetItem{})
		if next >= ps.end {
			next = ps.items
		}
		// Went all the way around?
		if next == idx {
			break
		}
		f := (*pointerSetItem)(unsafe.Pointer(next))
		if f.key == 0 || f.distance <= 0 {
			break
		}
		f.distance--
		*(*pointerSetItem)(unsafe.Pointer(idx)) = *f
		idx = next
	}
	// Clear entry
	*(*pointerSetItem)(unsafe.Pointer(idx)) = pointerSetItem{}
	ps.count--
	return idxStart, true
}

//goland:noinspection GoVetUnsafePointer
func (ps *PointerSet) Grow() bool {
	// Calculate new size
	// ps.size + 128
	if ps.growthFactor <= 1.0 {
		ps.growthFactor = 2.0
	}
	//newSize := ps.size + 32 // uintptr(float32(ps.size) * ps.growthFactor)
	newSize := uintptr(float32(ps.size) * ps.growthFactor)

	if gc_TRACE {
		println("PointerSet.Grow", "newSize", uint(newSize), "oldSize", uint(ps.size))
	}

	// Allocate new items table
	items := uintptr(ps.allocator.Alloc(newSize * unsafe.Sizeof(pointerSetItem{})))
	// Calculate end
	itemsEnd := items + (newSize * unsafe.Sizeof(pointerSetItem{}))
	// Zero the allocation out
	memzero(unsafe.Pointer(items), newSize*unsafe.Sizeof(pointerSetItem{}))
	// Init next structure
	next := PointerSet{
		items:        items,
		size:         newSize,
		end:          itemsEnd,
		allocator:    ps.allocator,
		count:        0,
		growthFactor: ps.growthFactor,
		maxDistance:  ps.maxDistance,
	}

	// Add all entries from old to next
	var success bool
	for i := ps.items; i < ps.end; i += unsafe.Sizeof(pointerSetItem{}) {
		if _, success = next.Add(*(*uintptr)(unsafe.Pointer(i)), 0); !success {
			return false
		}
	}

	// Free old items
	ps.allocator.Free(unsafe.Pointer(ps.items))
	ps.items = 0

	// Update to next
	*ps = next
	return true
}
