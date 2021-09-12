package mem

import "unsafe"

var arenas = make(map[unsafe.Pointer]Arena)

// SliceArena allocates memory by creating Go byte slices
type SliceArena struct {
	allocs map[uintptr][]byte
}

func NewSliceArena() *SliceArena {
	a := &SliceArena{
		allocs: make(map[uintptr][]byte, 16),
	}
	arenasMu.Lock()
	arenas[unsafe.Pointer(a)] = a
	arenasMu.Unlock()
	return a
}

func (a *SliceArena) Alloc(size uintptr) (uintptr, uintptr) {
	b := make([]byte, size)
	p := unsafe.Pointer(&b[0])
	a.allocs[uintptr(p)] = b
	return uintptr(p), uintptr(p) + uintptr(cap(b))
}

func (a *SliceArena) Free() {
	arenasMu.Lock()
	defer arenasMu.Unlock()
	// Clear all allocs to help GC marking
	for k, _ := range a.allocs {
		delete(a.allocs, k)
	}
	// Remove from global arenas map and the GC will free when needed
	delete(arenas, unsafe.Pointer(a))
}
