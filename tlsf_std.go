//go:build !tinygo && !wasm && !wasi && !tinygo.wasm && (darwin || linux)
// +build !tinygo
// +build !wasm
// +build !wasi
// +build !tinygo.wasm
// +build darwin linux

package mem

import (
	"sync"
	"unsafe"
)

func NewTLSF(pages int32) *TLSF {
	return NewTLSFArena(pages, NewSliceArena(), GrowMin)
}

func NewTLSFArena(pages int32, arena Arena, grow GrowFactory) *TLSF {
	if pages <= 0 {
		pages = 1
	}
	g := grow(arena)
	if g == nil {
		g = GrowMin(arena)
	}
	pagesAdded, start, end := g(0, pages, 0)
	a := bootstrap(start, end, pagesAdded, g)
	a.arena = Pointer(unsafe.Pointer(&arena))
	return a
}

type TLSFSync struct {
	a     *TLSF
	stats TLSFStats
	sync.Mutex
}

func (a *TLSF) ToSync() *TLSFSync {
	return &TLSFSync{a: a}
}

func (a *TLSFSync) AsAllocator() Allocator {
	return Allocator(unsafe.Pointer(a)) | _TLSFSync
}

// Alloc allocates a block of memory that fits the size provided.
// The allocation IS cleared / zeroed out.
func (a *TLSFSync) Alloc(size Pointer) Pointer {
	a.Lock()
	defer a.Unlock()
	return a.a.Alloc(size)
}

// AllocNotCleared allocates a block of memory that fits the size provided.
// The allocation is NOT cleared / zeroed out.
func (a *TLSFSync) AllocNotCleared(size Pointer) Pointer {
	a.Lock()
	defer a.Unlock()
	return a.a.AllocNotCleared(size)
}

// Realloc determines the best way to resize an allocation.
// Any extra size added is NOT cleared / zeroed out.
func (a *TLSFSync) Realloc(ptr Pointer, size Pointer) Pointer {
	a.Lock()
	defer a.Unlock()
	return Pointer(unsafe.Pointer(a.a.moveBlock(tlsfCheckUsedBlock(ptr), size))) + _TLSFBlockOverhead
}

// Free release the allocation back into the free list.
func (a *TLSFSync) Free(ptr Pointer) {
	a.Lock()
	defer a.Unlock()
	a.a.freeBlock(tlsfCheckUsedBlock(ptr))
}

//goland:noinspection GoVetUnsafePointer
func (a *TLSFSync) Bytes(length Pointer) Bytes {
	a.Lock()
	defer a.Unlock()
	return a.a.BytesCap(length, length)
}

//goland:noinspection GoVetUnsafePointer
func (a *TLSFSync) BytesCap(length, capacity Pointer) Bytes {
	a.Lock()
	defer a.Unlock()
	return a.a.BytesCap(length, capacity)
}

func (a *TLSFSync) BytesCapNotCleared(length, capacity Pointer) Bytes {
	a.Lock()
	defer a.Unlock()
	return a.a.BytesCapNotCleared(length, capacity)
}
