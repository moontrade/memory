package alloc

import (
	"github.com/moontrade/memory/hash"
	. "github.com/moontrade/memory/mem"
	"github.com/moontrade/memory/tlsf"
	"unsafe"
)

// PointerSet is a hashset that uses the robinhood algorithm. This
// implementation is not concurrent safe.
type PointerSet struct {
	// items are the slots of the hashmap for items.
	items uintptr
	end   uintptr
	size  uintptr

	// Number of keys in the PointerSet.
	count     uintptr
	allocator *tlsf.Heap
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
var pointerSetHash = hash.FNV32a

//func pointerSetHash(v uint32) uint32 {
//	return fnv32(v)
//}

// NewPointerSet returns a new robinhood hashmap.
//goland:noinspection ALL
func NewPointerSet(allocator *tlsf.Heap, size uintptr) PointerSet {
	items := allocator.Alloc(unsafe.Sizeof(pointerSetItem{}) * size)
	Zero(unsafe.Pointer(items), unsafe.Sizeof(pointerSetItem{})*uintptr(size))
	return PointerSet{
		items:        uintptr(items),
		size:         uintptr(size),
		end:          uintptr(items) + (uintptr(size) * unsafe.Sizeof(pointerSetItem{})),
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
	ps.allocator.Free(ps.items)
	ps.items = 0
	return nil
}

// Reset clears PointerSet, where already allocated memory will be reused.
//goland:noinspection ALL
func (ps *PointerSet) Reset() {
	Zero(unsafe.Pointer(ps.items), unsafe.Sizeof(pointerSetItem{})*ps.size)
	ps.count = 0
}

//goland:noinspection GoVetUnsafePointer
func (ps *PointerSet) isCollision(key uintptr) bool {
	return *(*uintptr)(unsafe.Pointer(
		ps.items + (uintptr(pointerSetHash(uint32(key))%uint32(ps.size)) * unsafe.Sizeof(pointerSetItem{})))) != 0
}

// Has returns whether the key exists in the Add.
//goland:noinspection ALL
func (ps *PointerSet) Has(key uintptr) bool {
	var (
		idx      = ps.items + (uintptr(pointerSetHash(uint32(key))%uint32(ps.size)) * unsafe.Sizeof(pointerSetItem{}))
		idxStart = idx
	)
	for {
		entry := *(*uintptr)(unsafe.Pointer(idx))
		if entry == 0 {
			return false
		}
		if entry == key {
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

func (ps *PointerSet) Set(key uintptr) (bool, bool) {
	return ps.Add(key, 0)
}

// Add inserts or updates a key into the PointerSet. The returned
// wasNew will be true if the mutation was on a newly seen, inserted
// key, and wasNew will be false if the mutation was an update to an
// existing key.
//goland:noinspection GoVetUnsafePointer
func (ps *PointerSet) Add(key uintptr, depth int) (bool, bool) {
	var (
		idx      = ps.items + (uintptr(pointerSetHash(uint32(key))%uint32(ps.size)) * unsafe.Sizeof(pointerSetItem{}))
		idxStart = idx
		incoming = pointerSetItem{key, 0}
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
func (ps *PointerSet) Delete(key uintptr) (uintptr, bool) {
	if key == 0 {
		return 0, false
	}

	var (
		idx      = ps.items + (uintptr(pointerSetHash(uint32(key))%uint32(ps.size)) * unsafe.Sizeof(pointerSetItem{}))
		idxStart = idx
	)
	for {
		entry := (*pointerSetItem)(unsafe.Pointer(idx))
		if entry.key == 0 {
			return 0, false
		}
		if entry.key == key {
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

	//if gc_TRACE {
	//	println("PointerSet.Grow", "newSize", uint(newSize), "oldSize", uint(ps.size))
	//}

	// Allocate new items table
	items := uintptr(ps.allocator.AllocZeroed(newSize * unsafe.Sizeof(pointerSetItem{})))
	// Calculate end
	itemsEnd := items + (newSize * unsafe.Sizeof(pointerSetItem{}))
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
	ps.allocator.Free(ps.items)
	ps.items = 0

	// Update to next
	*ps = next
	return true
}
