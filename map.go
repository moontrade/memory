package mem

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
	allocator *Allocator

	hash func(uint32) uint32

	// When any item's distance gets too large, Grow the Map.
	// Defaults to 10.
	maxDistance int32
	growBy      int32
	// Number of hash slots to Grow by
	growthFactor float32
}

// Item represents an entry in the Map.
type mapItem struct {
	key      uintptr
	value    uintptr
	distance int32 // How far item is from its best position.
}

// NewMap returns a new robinhood hashmap.
//goland:noinspection ALL
func NewMap(allocator *Allocator, size uintptr) Map {
	items := allocator.Alloc(unsafe.Sizeof(mapItem{}) * size)
	memzero(items, unsafe.Sizeof(mapItem{})*size)
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
	ps.allocator.Free(unsafe.Pointer(ps.items))
	ps.items = 0
	return nil
}

// Reset clears Map, where already allocated memory will be reused.
//goland:noinspection ALL
func (ps *Map) Reset() {
	memzero(unsafe.Pointer(ps.items), unsafe.Sizeof(mapItem{})*ps.size)
	ps.count = 0
}

//goland:noinspection GoVetUnsafePointer
func (ps *Map) isCollision(k unsafe.Pointer) bool {
	return *(*uintptr)(unsafe.Pointer(
		ps.items + (uintptr(ps.hash(uint32(uintptr(k)))%uint32(ps.size)) * unsafe.Sizeof(mapItem{})))) != 0
}

// Has returns whether the key exists in the Add.
//goland:noinspection ALL
func (ps *Map) Has(k uintptr) bool {
	var (
		idx      = ps.items + (uintptr(ps.hash(uint32(k))%uint32(ps.size)) * unsafe.Sizeof(mapItem{}))
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
func (ps *Map) Get(key unsafe.Pointer) (unsafe.Pointer, bool) {
	var (
		k        = uintptr(key)
		idx      = ps.items + (uintptr(ps.hash(uint32(k))%uint32(ps.size)) * unsafe.Sizeof(mapItem{}))
		idxStart = idx
	)
	for {
		entry := *(*mapItem)(unsafe.Pointer(idx))
		if entry.key == 0 {
			return nil, false
		}
		if entry.key == k {
			return unsafe.Pointer(entry.value), true
		}
		idx += unsafe.Sizeof(mapItem{})
		if idx >= ps.end {
			idx = ps.items
		}
		// Went all the way around?
		if idx == idxStart {
			return nil, false
		}
	}
}

// Set inserts or updates a key into the Map. The returned
// wasNew will be true if the mutation was on a newly seen, inserted
// key, and wasNew will be false if the mutation was an update to an
// existing key.
//goland:noinspection GoVetUnsafePointer
func (ps *Map) Set(key unsafe.Pointer, value unsafe.Pointer, depth int) (bool, bool) {
	var (
		k        = uintptr(key)
		v        = uintptr(value)
		idx      = ps.items + (uintptr(ps.hash(uint32(k))%uint32(ps.size)) * unsafe.Sizeof(mapItem{}))
		idxStart = idx
		incoming = mapItem{k, v, 0}
	)
	for {
		entry := (*mapItem)(unsafe.Pointer(idx))
		if entry.key == 0 {
			*(*mapItem)(unsafe.Pointer(idx)) = incoming
			ps.count++
			return true, true
		}

		if entry.key == incoming.key {
			entry.distance = incoming.distance
			return false, true
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
				return false, false
			}
			if !ps.Grow() {
				return false, false
			}
			return ps.Set(unsafe.Pointer(incoming.key), unsafe.Pointer(incoming.value), depth+1)
		}
	}
}

// Delete removes a key from the Map.
//goland:noinspection GoVetUnsafePointer
func (ps *Map) Delete(key unsafe.Pointer) (unsafe.Pointer, bool) {
	if key == nil {
		return nil, false
	}

	var (
		k        = uintptr(key)
		idx      = ps.items + (uintptr(ps.hash(uint32(k))%uint32(ps.size)) * unsafe.Sizeof(mapItem{}))
		idxStart = idx
		prev     unsafe.Pointer
	)
	for {
		entry := (*mapItem)(unsafe.Pointer(idx))
		if entry.key == 0 {
			return nil, false
		}

		if entry.key == k {
			// Found the item.
			prev = unsafe.Pointer(entry.value)
			break
		}
		idx += unsafe.Sizeof(mapItem{})
		if idx >= ps.end {
			idx = ps.items
		}
		if idx == idxStart {
			return nil, false
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
	//newSize := ps.size + 32 // uintptr(float32(ps.size) * ps.growthFactor)
	newSize := uintptr(float32(ps.size) * ps.growthFactor)

	if map_TRACE {
		println("Map.Grow", "newSize", uint(newSize), "oldSize", uint(ps.size))
	}

	// Allocate new items table
	items := uintptr(ps.allocator.Alloc(newSize * unsafe.Sizeof(mapItem{})))
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
		if _, success = next.Set(unsafe.Pointer(item.key), unsafe.Pointer(item.value), 0); !success {
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
