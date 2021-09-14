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
	stats Stats
	sync.Mutex
}

func (a *TLSF) ToSync() *TLSFSync {
	return &TLSFSync{a: a}
}

func (a *TLSFSync) AsAllocator() Allocator {
	return Allocator(unsafe.Pointer(a)) | _TLSFSync
}

// Alloc allocates a block of memory that fits the size provided
//goland:noinspection GoVetUnsafePointer
func (a *TLSFSync) Alloc(size Pointer) Pointer {
	a.Lock()
	defer a.Unlock()
	p := Pointer(unsafe.Pointer(a.a.allocateBlock(size)))
	if p == 0 {
		return p
	}
	return p + _TLSFBlockOverhead
}

// Alloc allocates a block of memory that fits the size provided
//goland:noinspection GoVetUnsafePointer
func (a *TLSFSync) AllocZeroed(size Pointer) Pointer {
	a.Lock()
	defer a.Unlock()
	p := Pointer(unsafe.Pointer(a.a.allocateBlock(size)))
	if p == 0 {
		return p
	}
	p = p + _TLSFBlockOverhead
	memzero(unsafe.Pointer(p), uintptr(size))
	return p
}

//goland:noinspection GoVetUnsafePointer
func (a *TLSFSync) Bytes(length Pointer) Bytes {
	return a.BytesCapacity(length, length)
}

//goland:noinspection GoVetUnsafePointer
func (a *TLSFSync) BytesCapacity(length, capacity Pointer) Bytes {
	a.Lock()
	defer a.Unlock()
	if capacity < length {
		capacity = length
	}
	p := Pointer(unsafe.Pointer(a.a.allocateBlock(capacity)))
	if p == 0 {
		return Bytes{}
	}
	return Bytes{
		Pointer: p + _TLSFBlockOverhead,
		len:     uint32(length),
		cap:     uint32(*(*Pointer)(unsafe.Pointer(p)) & ^_TLSFTagsMask),
		alloc:   Allocator(unsafe.Pointer(a)) | _TLSFSync,
	}
}

// Realloc determines the best way to resize an allocation.
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
