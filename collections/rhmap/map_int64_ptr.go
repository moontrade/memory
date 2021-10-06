package rhmap

import (
	. "github.com/moontrade/memory"
	"github.com/moontrade/memory/hash"
	"unsafe"
)

// MapInt64 is a hashset that uses the robinhood algorithm. This
// implementation is not concurrent safe.
type MapInt64 struct {
	// items are the slots of the hashmap for items.
	items Pointer
	end   uintptr
	size  uintptr

	// Number of keys in the Map.
	count uintptr

	// When any item's distance gets too large, Grow the Map.
	// Defaults to 10.
	maxDistance int32
	growBy      int32
	// Number of hash slots to Grow by
	growthFactor float32

	freeOnClose bool
}

// Item represents an entry in the Map.
type mapInt64Item struct {
	key      int64
	value    uintptr
	distance int32 // How far item is from its best position.
}

// NewMapInt64 returns a new robinhood hashmap.
//goland:noinspection ALL
func NewMapInt64(size uintptr) MapInt64 {
	//items := AllocZeroed(unsafe.Sizeof(mapInt64Item{}) * size)
	items := Calloc(uintptr(size), unsafe.Sizeof(mapInt64Item{}))
	return MapInt64{
		items:        items,
		size:         size,
		end:          uintptr(items) + (size * unsafe.Sizeof(mapInt64Item{})),
		maxDistance:  10,
		growBy:       64,
		growthFactor: 2.0,
	}
}

//goland:noinspection GoVetUnsafePointer
func (ps *MapInt64) Close() error {
	if ps == nil {
		return nil
	}
	if ps.items == 0 {
		return nil
	}

	Free(ps.items)
	ps.items = 0
	return nil
}

// Reset clears Map, where already allocated memory will be reused.
//goland:noinspection ALL
func (ps *MapInt64) Reset() {
	ps.items.Zero(unsafe.Sizeof(mapInt64Item{}) * ps.size)
	ps.count = 0
}

//goland:noinspection GoVetUnsafePointer
func (ps *MapInt64) isCollision(key int64) bool {
	return *(*uintptr)(unsafe.Pointer(
		uintptr(ps.items) + (uintptr(hash.FNV64a(uint64(key))%uint64(ps.size)) * unsafe.Sizeof(mapInt64Item{})))) != 0
}

// Has returns whether the key exists in the Add.
//goland:noinspection ALL
func (ps *MapInt64) Has(key int64) bool {
	var (
		idx      = uintptr(ps.items) + (uintptr(hash.FNV64a(uint64(key))%uint64(ps.size)) * unsafe.Sizeof(mapInt64Item{}))
		idxStart = idx
	)
	for {
		entry := (*mapInt64Item)(unsafe.Pointer(idx))
		if entry == nil {
			return false
		}
		if entry.key == key {
			return true
		}
		idx += unsafe.Sizeof(mapInt64Item{})
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
func (ps *MapInt64) Get(key int64) (uintptr, bool) {
	var (
		idx      = uintptr(ps.items) + (uintptr(hash.FNV64a(uint64(key))%uint64(ps.size)) * unsafe.Sizeof(mapInt64Item{}))
		idxStart = idx
	)
	for {
		entry := (*mapInt64Item)(unsafe.Pointer(idx))
		if entry.key == key {
			return entry.value, true
		}
		idx += unsafe.Sizeof(mapInt64Item{})
		if idx >= ps.end {
			idx = uintptr(ps.items)
		}
		// Went all the way around?
		if idx == idxStart {
			return 0, false
		}
	}
}

func (ps *MapInt64) Set(key int64, value uintptr) (uintptr, bool, bool) {
	return ps.set(key, value, 0)
}

// Set inserts or updates a key into the Map. The returned
// wasNew will be true if the mutation was on a newly seen, inserted
// key, and wasNew will be false if the mutation was an update to an
// existing key.
//goland:noinspection GoVetUnsafePointer
func (ps *MapInt64) set(key int64, value uintptr, depth int) (uintptr, bool, bool) {
	var (
		idx      = uintptr(ps.items) + (uintptr(hash.FNV64a(uint64(key))%uint64(ps.size)) * unsafe.Sizeof(mapInt64Item{}))
		idxStart = idx
		incoming = mapInt64Item{key, value, 0}
	)
	for {
		entry := (*mapInt64Item)(unsafe.Pointer(idx))
		//if entry.key == 0 {
		//	*(*mapInt64Item)(unsafe.Pointer(idx)) = incoming
		//	ps.count++
		//	return 0, true, true
		//}

		if entry.key == incoming.key {
			old := entry.value
			entry.value = incoming.value
			entry.distance = incoming.distance
			return old, false, true
		}

		// Swap if the incoming item is further from its best idx.
		if entry.distance < incoming.distance {
			incoming, *(*mapInt64Item)(unsafe.Pointer(idx)) = *(*mapInt64Item)(unsafe.Pointer(idx)), incoming
		}

		// One step further away from best idx.
		incoming.distance++

		idx += unsafe.Sizeof(mapInt64Item{})
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
func (ps *MapInt64) Delete(key int64) (uintptr, bool) {
	var (
		idx      = uintptr(ps.items) + (uintptr(hash.FNV64a(uint64(key))%uint64(ps.size)) * unsafe.Sizeof(mapInt64Item{}))
		idxStart = idx
		prev     uintptr
	)
	for {
		entry := (*mapInt64Item)(unsafe.Pointer(idx))

		if entry.key == key {
			// Found the item.
			prev = entry.value
			break
		}
		idx += unsafe.Sizeof(mapInt64Item{})
		if idx >= ps.end {
			idx = uintptr(ps.items)
		}
		if idx == idxStart {
			return 0, false
		}
	}
	// Left-shift succeeding items in the linear chain.
	for {
		next := idx + unsafe.Sizeof(mapInt64Item{})
		if next >= ps.end {
			next = uintptr(ps.items)
		}
		// Went all the way around?
		if next == idx {
			break
		}
		f := (*mapInt64Item)(unsafe.Pointer(next))
		if f.key == 0 || f.distance <= 0 {
			break
		}
		f.distance--
		*(*mapInt64Item)(unsafe.Pointer(idx)) = *f
		idx = next
	}
	// Clear entry
	*(*mapInt64Item)(unsafe.Pointer(idx)) = mapInt64Item{}
	ps.count--
	return prev, true
}

//goland:noinspection GoVetUnsafePointer
func (ps *MapInt64) Grow() bool {
	// Calculate new size
	// ps.size + 128
	if ps.growthFactor <= 1.0 {
		ps.growthFactor = 2.0
	}
	//newSize := uintptr(float32(ps.size) * ps.growthFactor)
	newSize := ps.size * 2

	if _TRACE {
		println("Map.Grow", "newSize", uint(newSize), "oldSize", uint(ps.size))
	}

	// Allocate new items table
	items := Calloc(newSize, unsafe.Sizeof(mapInt64Item{}))
	//items := AllocZeroed(newSize * unsafe.Sizeof(mapInt64Item{}))
	// Calculate end
	itemsEnd := uintptr(items) + (newSize * unsafe.Sizeof(mapInt64Item{}))
	// Init next structure
	next := MapInt64{
		items:        items,
		size:         newSize,
		end:          itemsEnd,
		count:        0,
		growthFactor: ps.growthFactor,
		maxDistance:  ps.maxDistance,
	}

	// Add all entries from old to next
	var (
		success bool
		item    *mapInt64Item
	)
	for i := uintptr(ps.items); i < ps.end; i += unsafe.Sizeof(mapInt64Item{}) {
		item = (*mapInt64Item)(unsafe.Pointer(i))
		if _, _, success = next.set(item.key, item.value, 0); !success {
			return false
		}
	}

	// Free old items
	Free(ps.items)
	ps.items = 0

	// Update to next
	*ps = next
	return true
}
