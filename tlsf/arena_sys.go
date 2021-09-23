//go:build !tinygo
// +build !tinygo

package tlsf

import (
	"sync"
	"unsafe"
)

type span struct {
	p uintptr
	n uintptr
}

type sysArena struct {
	allocs map[uintptr]span
	stat   sysMemStat
	mu     sync.Mutex
}

func NewSysArena() *sysArena {
	allocatorsMu.Lock()
	a := &sysArena{}
	if allocators == nil {
		allocators = make(map[unsafe.Pointer]Arena)
	}
	allocators[unsafe.Pointer(a)] = a
	allocatorsMu.Unlock()
	return a
}

func (s *sysArena) Size() uint64 {
	return uint64(s.stat)
}

func (a *sysArena) Alloc(size uintptr) (uintptr, uintptr) {
	a.mu.Lock()
	defer a.mu.Unlock()
	ptr := sysAlloc(size, &a.stat)
	if ptr == nil {
		return 0, 0
	}
	if a.allocs == nil {
		a.allocs = make(map[uintptr]span, 16)
	}
	a.allocs[uintptr(ptr)] = span{uintptr(ptr), size}
	return uintptr(ptr), uintptr(ptr) + size
}

func (a *sysArena) Free() {
	a.mu.Lock()
	defer a.mu.Unlock()

	if len(a.allocs) == 0 {
		return
	}

	// Free all sysAllocs by calling sysFree
	for k, v := range a.allocs {
		sysFree(unsafe.Pointer(v.p), v.n, &a.stat)
		delete(a.allocs, k)
	}

	allocatorsMu.Lock()
	delete(allocators, unsafe.Pointer(a))
	allocatorsMu.Unlock()
}

//go:linkname sysMemStat runtime.sysMemStat
type sysMemStat uint64

////go:linkname persistentalloc runtime.persistentalloc
//func persistentalloc(size, align uintptr, sysStat *sysMemStat) unsafe.Pointer

//go:linkname sysAlloc runtime.sysAlloc
func sysAlloc(size uintptr, sysStat *sysMemStat) unsafe.Pointer

//go:linkname sysFree runtime.sysFree
func sysFree(ptr unsafe.Pointer, n uintptr, sysStat *sysMemStat)
