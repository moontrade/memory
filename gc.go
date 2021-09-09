//go:build !gc.conservative && !gc.extalloc

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
	gc_OBJECT_OVERHEAD = (unsafe.Sizeof(gcObject{}) - tlsf_BLOCK_OVERHEAD + tlsf_AL_MASK) & ^tlsf_AL_MASK
	//gc_OBJECT_OVERHEAD = unsafe.Sizeof(gcObject{}) - tlsf_BLOCK_OVERHEAD + unsafe.Sizeof(uintptr(0))

	// Maximum size of a garbage collector object's payload.
	gc_OBJECT_MAXSIZE = tlsf_BLOCK_MAXSIZE - gc_OBJECT_OVERHEAD

	// Overhead of a garbage collector object. Excludes memory manager block overhead.
	gc_TOTAL_OVERHEAD = tlsf_BLOCK_OVERHEAD + gc_OBJECT_OVERHEAD
)

var (
	collector *GC
)

// GC is a Two-Color Mark & Sweep collector on top of a Two-Level Segmented Fit (TLSF) allocator.
// Similar features to the internal extalloc GC in TinyGo except GC uses a robinhood
// hashset instead of a treap structure and without the need for a linked list.
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
	allocs      gcSet
	first, last uintptr
	allocator   *Allocator
	stats       GCStats
}

// GCStats provides all the monitoring metrics needed to see how the GC
// is operating and performing.
type GCStats struct {
	Started           int64   // Epoch in nanos when GC was first started
	Cycles            int64   // Number of times GC collect has ran
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
	println("Moon GC cycle")
	println("\tlive:			", uint(s.Live))
	println("\tlive bytes:		", uint(s.LiveBytes))
	println("\tfrees:			", uint(s.Frees))
	println("\tallocs:			", uint(s.TotalAllocs))
	println("\tfreed bytes:		", uint(s.FreedBytes))
	println("\tsweep bytes:		", uint(s.SweepBytes))
	println("\ttotal bytes:		", uint(s.TotalBytes))
	println("\tlast sweep:		", uint(s.LastSweep))
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
	PrintDebugInfo()
}

//goland:noinspection ALL
func newGC(allocator *Allocator, initialCapacity uintptr) *GC {
	t := (*GC)(allocator.Alloc(unsafe.Sizeof(GC{})))
	t.allocator = allocator
	t.allocs = newGCSet(allocator, initialCapacity)
	t.first = ^uintptr(0)
	t.last = 0
	t.stats.Started = time.Now().UnixNano()
	return t
}

// Object Represents a managed object in memory, consisting of a header followed by the object's data.
type gcObject struct {
	tlsfBLOCK
	color  uint32 // Pointer to the next object with color flags stored in the alignment bits.
	rtSize uint32 // Runtime size.
}

// Gets the size of this object in memory.
func (o *gcObject) size() uintptr {
	return tlsf_BLOCK_OVERHEAD + (o.mmInfo & ^uintptr(3))
}

//goland:noinspection ALL
func (tc *GC) scan(start uintptr, end uintptr) {
	if start < tc.first || end > tc.last {
		return
	}

	// Align start and end pointers.
	start = (start + unsafe.Alignof(unsafe.Pointer(nil)) - 1) &^ (unsafe.Alignof(unsafe.Pointer(nil)) - 1)
	end &^= unsafe.Alignof(unsafe.Pointer(nil)) - 1

	// Mark all pointers.
	for ptr := start; ptr < end; ptr += unsafe.Alignof(unsafe.Pointer(nil)) {
		p := *(*uintptr)(unsafe.Pointer(ptr))
		tc.mark(p)
	}
}

//goland:noinspection ALL
func (tc *GC) mark(root uintptr) {
	root -= gc_TOTAL_OVERHEAD
	if root < tc.first || root > tc.last {
		return
	}
	if tc.allocs.has(root) {
		*(*uint32)(unsafe.Pointer(root + tlsf_BLOCK_OVERHEAD)) = gc_BLACK
	}
}

//goland:noinspection ALL
func (tc *GC) markRecursive(root uintptr, depth int) {
	root -= gc_TOTAL_OVERHEAD
	if !tc.allocs.has(root) {
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
		if uintptr(obj.rtSize)%sizeofPointer != 0 {
			return
		}
		start := root + gc_TOTAL_OVERHEAD
		end := start + uintptr(obj.rtSize)
		//start = (start + unsafe.Alignof(unsafe.Pointer(nil)) - 1) &^ (unsafe.Alignof(unsafe.Pointer(nil)) - 1)
		//end &^= unsafe.Alignof(unsafe.Pointer(nil)) - 1

		for ptr := start; ptr < end; ptr += unsafe.Alignof(unsafe.Pointer(nil)) {
			p := *(*uintptr)(unsafe.Pointer(ptr))
			if (p < tc.first || p > tc.last) || (p >= start && p < end) {
				continue
			}
			if !tc.allocs.has(p - gc_TOTAL_OVERHEAD) {
				continue
			}
			tc.markRecursive(p, depth+1)
		}
	}
}

func (tc *GC) markRoots(start, end uintptr) {
	if gc_TRACE {
		println("markRoots", uint(start), uint(end))
	}
	tc.scan(start, end)
}

//goland:noinspection ALL
func (tc *GC) markGraph(root uintptr) {
	var (
		obj   = (*gcObject)(unsafe.Pointer(root))
		start = root + gc_TOTAL_OVERHEAD
		end   = start + uintptr(obj.rtSize)
	)

	// unaligned allocation must be some sort of string or data buffer. skip it.
	if uintptr(obj.rtSize)%sizeofPointer != 0 {
		return
	}

	pointersToCount := (uint(end) - uint(start)) / uint(unsafe.Sizeof(unsafe.Pointer(nil)))
	if pointersToCount > 128 {
		//return
	}

	// Mark all pointers.
	for ptr := start; ptr < end; ptr += unsafe.Alignof(unsafe.Pointer(nil)) {
		p := *(*uintptr)(unsafe.Pointer(ptr))
		if (p < tc.first || p > tc.last) || (p >= start && p < end) {
			continue
		}
		if !tc.allocs.has(p - gc_TOTAL_OVERHEAD) {
			continue
		}
		tc.markRecursive(p, 0)
	}
}

//goland:noinspection ALL
func (tc *GC) new(size uintptr) uintptr {
	// Is the size too large?
	if size > gc_OBJECT_MAXSIZE {
		panic("allocation too large")
	}

	// Allocate memory
	obj := (*gcObject)(unsafe.Pointer(uintptr(tc.allocator.Alloc(gc_OBJECT_OVERHEAD+size)) - tlsf_BLOCK_OVERHEAD))

	// set the runtime size and set to WHITE
	obj.rtSize = uint32(size)
	obj.color = gc_WHITE
	tc.stats.LiveBytes += obj.size()
	tc.stats.TotalBytes += int64(obj.size())
	tc.stats.Live++
	tc.stats.TotalAllocs++

	// Convert to uint pointer
	ptr := uintptr(unsafe.Pointer(obj))

	// Zero out the allocation
	memzero(unsafe.Pointer(ptr+gc_TOTAL_OVERHEAD), size)

	// Add to allocations map
	tc.allocs.set(ptr)

	// Update first pointer if necessary
	if ptr < tc.first {
		tc.first = ptr
	}
	// Update last pointer if necessary
	if ptr > tc.last {
		tc.last = ptr
	}

	// Return pointer to data
	return ptr + gc_TOTAL_OVERHEAD
}

// free will immediately remove the GC'd object from the collector
// and free up the memory in the underlying tlsf allocator.
//goland:noinspection ALL
func (tc *GC) free(ptr unsafe.Pointer) bool {
	p := uintptr(ptr) - gc_TOTAL_OVERHEAD
	if !tc.allocs.has(p) {
		return false
	}

	if gc_TRACE {
		println("GC free", uint(uintptr(ptr)))
	}

	obj := (*gcObject)(unsafe.Pointer(p))
	size := obj.size()
	tc.stats.LiveBytes -= size
	tc.stats.FreedBytes += int64(size)
	tc.stats.Live--
	tc.stats.Frees++
	tc.allocs.del(p)

	println("GC free", uint(uintptr(ptr)), "size", uint(size), "rtSize", obj.rtSize)
	tc.allocator.Free(unsafe.Pointer(p + tlsf_BLOCK_OVERHEAD))

	return true
}

//goland:noinspection ALL
func (tc *GC) collect() {
	if gc_TRACE {
		println("moon GC collect started...")
	}
	//println("tcmsCollect")
	var (
		start = time.Now().UnixNano()
		stats = &tc.stats
		k     uintptr
		obj   *gcObject
		min   = ^uintptr(0)
		max   = uintptr(0)
	)
	stats.Cycles++

	////////////////////////////////////////////////////////////////////////
	// Mark Roots Phase
	////////////////////////////////////////////////////////////////////////
	markStack()
	markGlobals()
	markScheduler()
	// End of mark roots
	markTime := time.Now().UnixNano() - start

	////////////////////////////////////////////////////////////////////////
	// Mark Graph Phase
	////////////////////////////////////////////////////////////////////////
	start = markTime + start
	stats.LastSweep = 0
	stats.LastSweepBytes = 0
	var (
		items     = tc.allocs.items
		itemsSize = tc.allocs.size
		itemsEnd  = items + (itemsSize * unsafe.Sizeof(gcSetItem{}))
	)
	for ; items < itemsEnd; items += unsafe.Sizeof(gcSetItem{}) {
		k = *(*uintptr)(unsafe.Pointer(items))
		if k == 0 {
			continue
		}
		tc.markGraph(k)
	}

	// End of mark graph
	markGraphTime := time.Now().UnixNano() - start

	////////////////////////////////////////////////////////////////////////
	// Sweep Phase
	////////////////////////////////////////////////////////////////////////
	start = markGraphTime + start

	// reset items iterator
	items = tc.allocs.items
	itemsSize = tc.allocs.size
	itemsEnd = items + (itemsSize * unsafe.Sizeof(gcSetItem{}))
	for ; items < itemsEnd; items += unsafe.Sizeof(gcSetItem{}) {
		k = *(*uintptr)(unsafe.Pointer(items))
		// Empty item?
		if k == 0 {
			continue
		}
		obj = (*gcObject)(unsafe.Pointer(k))
		if obj.color == gc_WHITE {
			stats.LiveBytes -= obj.size()
			stats.LastSweepBytes += int64(obj.size())
			tc.stats.Live--
			stats.LastSweep++

			if gc_TRACE {
				println("GC sweep", uint(k), "size", uint(obj.size()))
			}

			println("GC sweep", uint(k+gc_TOTAL_OVERHEAD), "size", uint(obj.size()))

			// Free memory
			tc.allocator.Free(unsafe.Pointer(k + tlsf_BLOCK_OVERHEAD))

			// Remove from alloc map
			tc.allocs.del(k)
			//items -= unsafe.Sizeof(gcSetItem{})
		} else {
			if k < min {
				min = k
			}
			if k > max {
				max = k
			}
			if gc_TRACE {
				println("GC retained", uint(k), "size", uint(obj.size()))
			}
			obj.color = gc_WHITE
		}
	}

	tc.first = min
	tc.last = max
	sweepTime := time.Now().UnixNano() - start
	stats.LastMarkRootsTime = markTime
	stats.LastMarkGraphTime = markGraphTime
	stats.SweepTime += sweepTime
	stats.SweepBytes += stats.LastSweepBytes

	stats.Sweeps += stats.LastSweep
	stats.LastGCTime = markTime + markGraphTime + sweepTime
	stats.TotalTime += stats.LastGCTime

	if gc_TRACE {
		println("moon GC collect finished")
	}
	//stats.Print()
}

// gcSet is a hashset that uses the robinhood algorithm. This
// implementation is not concurrent safe.
type gcSet struct {
	// items are the slots of the hashmap for items.
	items uintptr
	end   uintptr
	size  uintptr

	// Number of keys in the gcSet.
	count     uintptr
	allocator *Allocator
	// When any item's distance gets too large, grow the gcSet.
	// Defaults to 10.
	maxDistance int32
	growBy      int32
	// Number of hash slots to grow by
	growthFactor float32
}

// Item represents an entry in the gcSet.
type gcSetItem struct {
	key      uintptr
	distance int32 // How far item is from its best position.
}

const (
	sizeOfGCSetItem = unsafe.Sizeof(gcSetItem{})
)

// gcSetHash uses the fnv hashing algorithm for 32bit integers.
// fnv was chosen for its consistent low collision rate even with tight monotonic numbers (WASM)
// and for its performance. Other potential candidates are wyhash, metro, and adler32. Each of these
// have optimized 32bit version in hash.go in this package.
//go:inline
var gcSetHash = fnv32

//func gcSetHash(v uint32) uint32 {
//	return fnv32(v)
//}

// newGCSet returns a new robinhood hashmap.
//goland:noinspection ALL
func newGCSet(allocator *Allocator, size uintptr) gcSet {
	items := allocator.Alloc(unsafe.Sizeof(gcSetItem{}) * size)
	memzero(items, unsafe.Sizeof(gcSetItem{})*size)
	return gcSet{
		items:        uintptr(items),
		size:         size,
		end:          uintptr(items) + (size * sizeOfGCSetItem),
		maxDistance:  10,
		growBy:       64,
		growthFactor: 2.0,
	}
}

// reset clears gcSet, where already allocated memory will be reused.
//goland:noinspection ALL
func (m *gcSet) reset() {
	memzero(unsafe.Pointer(m.items), unsafe.Sizeof(gcSetItem{})*m.size)
	m.count = 0
}

// has returns whether the key exists in the set.
//goland:noinspection ALL
func (m *gcSet) has(k uintptr) bool {
	var (
		idx      = m.items + (uintptr(gcSetHash(uint32(k))%uint32(m.size)) * sizeOfGCSetItem)
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
		idx += sizeOfGCSetItem
		if idx >= m.end {
			idx = m.items
		}
		// Went all the way around?
		if idx == idxStart {
			return false
		}
	}
}

// set inserts or updates a key into the gcSet. The returned
// wasNew will be true if the mutation was on a newly seen, inserted
// key, and wasNew will be false if the mutation was an update to an
// existing key.
//goland:noinspection ALL
func (m *gcSet) set(k uintptr) bool {
	var (
		idx      = m.items + (uintptr(gcSetHash(uint32(k))%uint32(m.size)) * sizeOfGCSetItem)
		idxStart = idx
		incoming = gcSetItem{k, 0}
	)
	for {
		entry := (*gcSetItem)(unsafe.Pointer(idx))
		if entry.key == 0 {
			*(*gcSetItem)(unsafe.Pointer(idx)) = incoming
			m.count++
			return true
		}

		if entry.key == incoming.key {
			entry.distance = incoming.distance
			return false
		}

		// Swap if the incoming item is further from its best idx.
		if entry.distance < incoming.distance {
			incoming, *(*gcSetItem)(unsafe.Pointer(idx)) = *(*gcSetItem)(unsafe.Pointer(idx)), incoming
		}

		// One step further away from best idx.
		incoming.distance++

		idx += sizeOfGCSetItem
		if idx >= m.end {
			idx = m.items
		}

		// Grow if distances become big or we went all the way around.
		if incoming.distance > m.maxDistance || idx == idxStart {
			m.grow()
			return m.set(incoming.key)
		}
	}
}

// del removes a key from the gcSet.
//goland:noinspection ALL
func (m *gcSet) del(k uintptr) (uintptr, bool) {
	if k == 0 {
		return 0, false
	}

	var (
		idx      = m.items + (uintptr(gcSetHash(uint32(k))%uint32(m.size)) * sizeOfGCSetItem)
		idxStart = idx
	)
	for {
		entry := (*gcSetItem)(unsafe.Pointer(idx))
		if entry.key == 0 {
			return 0, false
		}
		if entry.key == k {
			break // Found the item.
		}
		idx += sizeOfGCSetItem
		if idx >= m.end {
			idx = m.items
		}
		if idx == idxStart {
			return 0, false
		}
	}
	// Left-shift succeeding items in the linear chain.
	for {
		next := idx + sizeOfGCSetItem
		if next >= m.end {
			next = m.items
		}
		// Went all the way around?
		if next == idx {
			break
		}
		f := (*gcSetItem)(unsafe.Pointer(next))
		if f.key == 0 || f.distance <= 0 {
			break
		}
		f.distance--
		*(*gcSetItem)(unsafe.Pointer(idx)) = *f
		idx = next
	}
	// Clear entry
	*(*gcSetItem)(unsafe.Pointer(idx)) = gcSetItem{}
	m.count--
	return idxStart, true
}

//goland:noinspection ALL
func (m *gcSet) grow() {
	// Calculate new size
	// m.size + 128
	if m.growthFactor <= 1.0 {
		m.growthFactor = 2.0
	}
	//newSize := m.size + 32 // uintptr(float32(m.size) * m.growthFactor)
	newSize := uintptr(float32(m.size) * m.growthFactor)

	if gc_TRACE {
		println("gcSet.grow", "newSize", uint(newSize), "oldSize", uint(m.size))
	}

	// Allocate new items table
	items := uintptr(m.allocator.Alloc(newSize * sizeOfGCSetItem))
	// Calculate end
	itemsEnd := items + (newSize * sizeOfGCSetItem)
	// Zero the allocation out
	memzero(unsafe.Pointer(items), newSize*sizeOfGCSetItem)
	// Init next structure
	next := gcSet{
		items:        items,
		size:         newSize,
		end:          itemsEnd,
		count:        0,
		growthFactor: m.growthFactor,
		maxDistance:  m.maxDistance,
	}

	// Add all entries from old to next
	for i := m.items; i < m.end; i += sizeOfGCSetItem {
		_ = next.set(*(*uintptr)(unsafe.Pointer(i)))
	}

	// Free old items
	m.allocator.Free(unsafe.Pointer(m.items))

	// Update to next
	*m = next
}
