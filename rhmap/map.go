package rhmap

import "unsafe"

const (
	map_TRACE = false
)

// Map is a hashset that uses the robinhood algorithm. This
// implementation is not concurrent safe.
type Map struct {
	// items are the slots of the hashmap for items.
	items uintptr
	end   uintptr
	size  uintptr

	// Number of keys in the Map.
	count     uintptr
	allocator *TLSF

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
	key      Bytes
	value    Bytes
	distance int32 // How far item is from its best position.
}

// NewMap returns a new robinhood hashmap.
//goland:noinspection ALL
func NewMap(allocator *TLSF, size uintptr) Map {
	items := allocator.Alloc(Pointer(unsafe.Sizeof(mapItem{}) * size))
	memzero(items.Unsafe(), unsafe.Sizeof(mapItem{})*size)
	return Map{
		items:        uintptr(items),
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

	ps.allocator.Free(Pointer(ps.items))
	ps.items = 0
	return nil
}

func (ps *Map) freeAll() {
	for i := ps.items; i < ps.end; i += unsafe.Sizeof(mapItem{}) {
		item := (*mapItem)(unsafe.Pointer(i))
		if item.key.Pointer == 0 {
			continue
		}
		item.key.Free()
		if item.value.Pointer != 0 {
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
	memzero(unsafe.Pointer(ps.items), unsafe.Sizeof(mapItem{})*ps.size)
	ps.count = 0
}

//goland:noinspection GoVetUnsafePointer
func (ps *Map) isCollision(key Bytes) bool {
	return *(*uintptr)(unsafe.Pointer(
		ps.items + (uintptr(key.Hash32()%uint32(ps.size)) * unsafe.Sizeof(pointerSetItem{})))) != 0
}

// Has returns whether the key exists in the Add.
//goland:noinspection ALL
func (ps *Map) Has(key Bytes) bool {
	var (
		idx      = ps.items + (uintptr(key.Hash32()%uint32(ps.size)) * unsafe.Sizeof(pointerSetItem{}))
		idxStart = idx
	)
	for {
		entry := (*mapItem)(unsafe.Pointer(idx))
		if entry == nil {
			return false
		}
		if entry.key.Equals(&key) {
			return true
		}
		idx += unsafe.Sizeof(mapItem{})
		if idx >= ps.end {
			idx = ps.items
		}
		// Went all the way around?
		if idx == idxStart {
			return false
		}
	}
}

// Get returns whether the key exists in the Add.
//goland:noinspection ALL
func (ps *Map) Get(key Bytes) (Bytes, bool) {
	var (
		idx      = ps.items + (uintptr(key.Hash32()%uint32(ps.size)) * unsafe.Sizeof(pointerSetItem{}))
		idxStart = idx
	)
	for {
		entry := (*mapItem)(unsafe.Pointer(idx))
		if entry.key.Pointer == 0 {
			return Bytes{}, false
		}
		if entry.key.Equals(&key) {
			return entry.value, true
		}
		idx += unsafe.Sizeof(mapItem{})
		if idx >= ps.end {
			idx = ps.items
		}
		// Went all the way around?
		if idx == idxStart {
			return Bytes{}, false
		}
	}
}

func (ps *Map) Set(key Bytes, value Bytes) (Bytes, bool, bool) {
	return ps.set(key, value, 0)
}

// Set inserts or updates a key into the Map. The returned
// wasNew will be true if the mutation was on a newly seen, inserted
// key, and wasNew will be false if the mutation was an update to an
// existing key.
//goland:noinspection GoVetUnsafePointer
func (ps *Map) set(key Bytes, value Bytes, depth int) (Bytes, bool, bool) {
	var (
		idx      = ps.items + (uintptr(key.Hash32()%uint32(ps.size)) * unsafe.Sizeof(pointerSetItem{}))
		idxStart = idx
		incoming = mapItem{key, value, 0}
	)
	for {
		entry := (*mapItem)(unsafe.Pointer(idx))
		if entry.key.Pointer == 0 {
			*(*mapItem)(unsafe.Pointer(idx)) = incoming
			ps.count++
			return Bytes{}, true, true
		}

		if entry.key.Equals(&incoming.key) {
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
			idx = ps.items
		}

		// Grow if distances become big or we went all the way around.
		if incoming.distance > ps.maxDistance || idx == idxStart {
			if depth > 5 {
				return Bytes{}, false, false
			}
			if !ps.Grow() {
				return Bytes{}, false, false
			}
			return ps.set(incoming.key, incoming.value, depth+1)
		}
	}
}

// Delete removes a key from the Map.
//goland:noinspection GoVetUnsafePointer
func (ps *Map) Delete(key Bytes) (Bytes, bool) {
	if key.Pointer == 0 {
		return Bytes{}, false
	}

	var (
		idx      = ps.items + (uintptr(key.Hash32()%uint32(ps.size)) * unsafe.Sizeof(pointerSetItem{}))
		idxStart = idx
		prev     Bytes
	)
	for {
		entry := (*mapItem)(unsafe.Pointer(idx))
		if entry.key.Pointer == 0 {
			return Bytes{}, false
		}

		if entry.key.Equals(&key) {
			// Found the item.
			prev = entry.value
			break
		}
		idx += unsafe.Sizeof(mapItem{})
		if idx >= ps.end {
			idx = ps.items
		}
		if idx == idxStart {
			return Bytes{}, false
		}
	}
	// Left-shift succeeding items in the linear chain.
	for {
		next := idx + unsafe.Sizeof(mapItem{})
		if next >= ps.end {
			next = ps.items
		}
		// Went all the way around?
		if next == idx {
			break
		}
		f := (*mapItem)(unsafe.Pointer(next))
		if f.key.Pointer == 0 || f.distance <= 0 {
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
	items := uintptr(ps.allocator.Alloc(Pointer(newSize * unsafe.Sizeof(mapItem{}))))
	// Calculate end
	itemsEnd := items + (newSize * unsafe.Sizeof(mapItem{}))
	// Zero the allocation out
	memzero(unsafe.Pointer(items), newSize*unsafe.Sizeof(mapItem{}))
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
	for i := ps.items; i < ps.end; i += unsafe.Sizeof(mapItem{}) {
		item = (*mapItem)(unsafe.Pointer(i))
		if _, _, success = next.set(item.key, item.value, 0); !success {
			return false
		}
	}

	// Free old items
	ps.allocator.Free(Pointer(ps.items))
	ps.items = 0

	// Update to next
	*ps = next
	return true
}
