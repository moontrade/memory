//go:build !tinygo && !wasm && !wasi && !tinygo.wasm
// +build !tinygo,!wasm,!wasi,!tinygo.wasm

package mem

import (
	"sync"
	"unsafe"
)

var (
	arenas   map[unsafe.Pointer]Arena
	arenasMu sync.Mutex
)

// ArenaSlice allocates memory by creating Go byte slices
type ArenaSlice struct {
	allocs map[Pointer][]byte
}

func NewSliceArena() *ArenaSlice {
	// Allocate on Go GC
	a := &ArenaSlice{
		allocs: make(map[Pointer][]byte, 16),
	}
	arenasMu.Lock()
	if arenas == nil {
		arenas = make(map[unsafe.Pointer]Arena)
	}
	arenas[unsafe.Pointer(a)] = a
	arenasMu.Unlock()
	return a
}

func (a *ArenaSlice) Alloc(size Pointer) (Pointer, Pointer) {
	b := make([]byte, size)
	p := unsafe.Pointer(&b[0])
	a.allocs[Pointer(p)] = b
	return Pointer(p), Pointer(p) + Pointer(cap(b))
}

func (a *ArenaSlice) Free() {
	arenasMu.Lock()
	defer arenasMu.Unlock()
	// Clear all allocs to help GC marking
	for k := range a.allocs {
		delete(a.allocs, k)
	}
	// Remove from global arenas map and the GC will free when needed
	delete(arenas, unsafe.Pointer(a))
}