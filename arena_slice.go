package mem

import (
	"sync"
	"unsafe"
)

var (
	arenas   = make(map[unsafe.Pointer]Arena)
	arenasMu sync.Mutex
)

// SliceArena allocates memory by creating Go byte slices
type SliceArena struct {
	allocs map[Pointer][]byte
}

func NewSliceArena() *SliceArena {
	a := &SliceArena{
		allocs: make(map[Pointer][]byte, 16),
	}
	arenasMu.Lock()
	arenas[unsafe.Pointer(a)] = a
	arenasMu.Unlock()
	return a
}

func (a *SliceArena) Alloc(size Pointer) (Pointer, Pointer) {
	b := make([]byte, size)
	p := unsafe.Pointer(&b[0])
	a.allocs[Pointer(p)] = b
	return Pointer(p), Pointer(p) + Pointer(cap(b))
}

func (a *SliceArena) Free() {
	arenasMu.Lock()
	defer arenasMu.Unlock()
	// Clear all allocs to help GC marking
	for k := range a.allocs {
		delete(a.allocs, k)
	}
	// Remove from global arenas map and the GC will free when needed
	delete(arenas, unsafe.Pointer(a))
}
