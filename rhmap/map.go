package rhmap

import (
	mem "github.com/moontrade/memory"
	"unsafe"
)

const (
	map_TRACE = false
)

// Map is a hashset that uses the robinhood algorithm. This
// implementation is not concurrent safe.
type Map struct {
	// items are the slots of the hashmap for items.
	items mem.Pointer
	end   uintptr
	size  uintptr

	// Number of keys in the Map.
	count     uintptr
	allocator mem.Allocator

	// When any item's distance gets too large, Grow the Map.
	// Defaults to 10.
	maxDistance int32
	growBy      int32
	// Number of hash slots to Grow by
	growthFactor float32

	freeOnClose bool
}

// Item represents an entry in the Map.
type mapItem struct {
	key      mem.Str
	value    mem.Str
	distance int32 // How far item is from its best position.
}

// NewMap returns a new robinhood hashmap.
//goland:noinspection ALL
func NewMap(allocator mem.Allocator, size uintptr) Map {
	items := allocator.AllocZeroed(unsafe.Sizeof(mapItem{}) * size)
	return Map{
		items:        items,
		size:         size,
		end:          uintptr(items) + (size * unsafe.Sizeof(mapItem{})),
		allocator:    allocator,
		maxDistance:  10,
		growBy:       64,
		growthFactor: 2.0,
	}
}

//goland:noinspection GoVetUnsafePointer
func (ps *Map) Close() error {
	if ps == nil {
		return nil
	}
	if ps.items == 0 {
		return nil
	}

	if ps.freeOnClose {
		ps.freeAll()
	}

	ps.allocator.Free(ps.items)
	ps.items = 0
	return nil
}

func (ps *Map) freeAll() {
	for i := uintptr(ps.items); i < ps.end; i += unsafe.Sizeof(mapItem{}) {
		item := (*mapItem)(unsafe.Pointer(i))
		if item.key == 0 {
			continue
		}
		item.key.Free()
		if item.value != 0 {
			item.value.Free()
		}
	}
}

// Reset clears Map, where already allocated memory will be reused.
//goland:noinspection ALL
func (ps *Map) Reset() {
	if ps.freeOnClose {
		ps.freeAll()
	}
	ps.items.Zero(unsafe.Sizeof(mapItem{}) * ps.size)
	ps.count = 0
}

//goland:noinspection GoVetUnsafePointer
func (ps *Map) isCollision(key mem.Str) bool {
	return *(*uintptr)(unsafe.Pointer(
		uintptr(ps.items) + (uintptr(key.Hash32()%uint32(ps.size)) * unsafe.Sizeof(mapItem{})))) != 0
}

// Has returns whether the key exists in the Add.
//goland:noinspection ALL
func (ps *Map) Has(key mem.Str) bool {
	var (
		idx      = uintptr(ps.items) + (uintptr(key.Hash32()%uint32(ps.size)) * unsafe.Sizeof(mapItem{}))
		idxStart = idx
	)
	for {
		entry := (*mapItem)(unsafe.Pointer(idx))
		if entry == nil {
			return false
		}
		if entry.key.Equals(key) {
			return true
		}
		idx += unsafe.Sizeof(mapItem{})
		if idx >= ps.end {
			idx = uintptr(ps.items)
		}
		// Went all the way around?
		if idx == idxStart {
			return false
		}
	}
}

// Get returns whether the key exists in the Add.
//goland:noinspection ALL
func (ps *Map) Get(key mem.Str) (mem.Str, bool) {
	var (
		idx      = uintptr(ps.items) + (uintptr(key.Hash32()%uint32(ps.size)) * unsafe.Sizeof(mapItem{}))
		idxStart = idx
	)
	for {
		entry := (*mapItem)(unsafe.Pointer(idx))
		if entry.key == 0 {
			return 0, false
		}
		if entry.key.Equals(key) {
			return entry.value, true
		}
		idx += unsafe.Sizeof(mapItem{})
		if idx >= ps.end {
			idx = uintptr(ps.items)
		}
		// Went all the way around?
		if idx == idxStart {
			return 0, false
		}
	}
}

func (ps *Map) Set(key mem.Str, value mem.Str) (mem.Str, bool, bool) {
	return ps.set(key, value, 0)
}

// Set inserts or updates a key into the Map. The returned
// wasNew will be true if the mutation was on a newly seen, inserted
// key, and wasNew will be false if the mutation was an update to an
// existing key.
//goland:noinspection GoVetUnsafePointer
func (ps *Map) set(key mem.Str, value mem.Str, depth int) (mem.Str, bool, bool) {
	var (
		idx      = uintptr(ps.items) + (uintptr(key.Hash32()%uint32(ps.size)) * unsafe.Sizeof(mapItem{}))
		idxStart = idx
		incoming = mapItem{key, value, 0}
	)
	for {
		entry := (*mapItem)(unsafe.Pointer(idx))
		if entry.key == 0 {
			*(*mapItem)(unsafe.Pointer(idx)) = incoming
			ps.count++
			return 0, true, true
		}

		if entry.key.Equals(incoming.key) {
			old := entry.value
			entry.value = incoming.value
			entry.distance = incoming.distance
			return old, false, true
		}

		// Swap if the incoming item is further from its best idx.
		if entry.distance < incoming.distance {
			incoming, *(*mapItem)(unsafe.Pointer(idx)) = *(*mapItem)(unsafe.Pointer(idx)), incoming
		}

		// One step further away from best idx.
		incoming.distance++

		idx += unsafe.Sizeof(mapItem{})
		if idx >= ps.end {
			idx = uintptr(ps.items)
		}

		// Grow if distances become big or we went all the way around.
		if incoming.distance > ps.maxDistance || idx == idxStart {
			if depth > 5 {
				return 0, false, false
			}
			if !ps.Grow() {
				return 0, false, false
			}
			return ps.set(incoming.key, incoming.value, depth+1)
		}
	}
}

// Delete removes a key from the Map.
//goland:noinspection GoVetUnsafePointer
func (ps *Map) Delete(key mem.Str) (mem.Str, bool) {
	if key == 0 {
		return 0, false
	}

	var (
		idx      = uintptr(ps.items) + (uintptr(key.Hash32()%uint32(ps.size)) * unsafe.Sizeof(mapItem{}))
		idxStart = idx
		prev     mem.Str
	)
	for {
		entry := (*mapItem)(unsafe.Pointer(idx))
		if entry.key == 0 {
			return 0, false
		}

		if entry.key.Equals(key) {
			// Found the item.
			prev = entry.value
			break
		}
		idx += unsafe.Sizeof(mapItem{})
		if idx >= ps.end {
			idx = uintptr(ps.items)
		}
		if idx == idxStart {
			return 0, false
		}
	}
	// Left-shift succeeding items in the linear chain.
	for {
		next := idx + unsafe.Sizeof(mapItem{})
		if next >= ps.end {
			next = uintptr(ps.items)
		}
		// Went all the way around?
		if next == idx {
			break
		}
		f := (*mapItem)(unsafe.Pointer(next))
		if f.key == 0 || f.distance <= 0 {
			break
		}
		f.distance--
		*(*mapItem)(unsafe.Pointer(idx)) = *f
		idx = next
	}
	// Clear entry
	*(*mapItem)(unsafe.Pointer(idx)) = mapItem{}
	ps.count--
	return prev, true
}

//goland:noinspection GoVetUnsafePointer
func (ps *Map) Grow() bool {
	// Calculate new size
	// ps.size + 128
	if ps.growthFactor <= 1.0 {
		ps.growthFactor = 2.0
	}
	newSize := uintptr(float32(ps.size) * ps.growthFactor)

	if map_TRACE {
		println("Map.Grow", "newSize", uint(newSize), "oldSize", uint(ps.size))
	}

	// Allocate new items table
	items := ps.allocator.AllocZeroed(newSize * unsafe.Sizeof(mapItem{}))
	// Calculate end
	itemsEnd := uintptr(items) + (newSize * unsafe.Sizeof(mapItem{}))
	// Init next structure
	next := Map{
		items:        items,
		size:         newSize,
		end:          itemsEnd,
		allocator:    ps.allocator,
		count:        0,
		growthFactor: ps.growthFactor,
		maxDistance:  ps.maxDistance,
	}

	// Add all entries from old to next
	var (
		success bool
		item    *mapItem
	)
	for i := uintptr(ps.items); i < ps.end; i += unsafe.Sizeof(mapItem{}) {
		item = (*mapItem)(unsafe.Pointer(i))
		if _, _, success = next.set(item.key, item.value, 0); !success {
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
