//go:build !tinygo && !wasm && !wasi && !tinygo.wasm
// +build !tinygo,!wasm,!wasi,!tinygo.wasm

package tlsf

import (
	"sync"
	"unsafe"
)

var (
	allocators   map[unsafe.Pointer]Arena
	allocatorsMu sync.Mutex
)

// SliceArena allocates memory by creating Go byte slices
type SliceArena struct {
	allocs map[uintptr][]byte
}

func NewSliceArena() *SliceArena {
	// Allocate on Go GC
	a := &SliceArena{
		allocs: make(map[uintptr][]byte, 16),
	}
	allocatorsMu.Lock()
	if allocators == nil {
		allocators = make(map[unsafe.Pointer]Arena)
	}
	allocators[unsafe.Pointer(a)] = a
	allocatorsMu.Unlock()
	return a
}

func (a *SliceArena) Alloc(size uintptr) (uintptr, uintptr) {
	b := make([]byte, size)
	p := unsafe.Pointer(&b[0])
	a.allocs[uintptr(p)] = b
	return uintptr(p), uintptr(p) + uintptr(cap(b))
}

func (a *SliceArena) Free() {
	allocatorsMu.Lock()
	defer allocatorsMu.Unlock()
	// Clear all allocs to help GC marking
	for k := range a.allocs {
		delete(a.allocs, k)
	}
	// Remove from global allocators map and the GC will free when needed
	delete(allocators, unsafe.Pointer(a))
}
