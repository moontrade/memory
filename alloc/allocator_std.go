//go:build !tinygo && !wasm && !wasi && !tinygo.wasm
// +build !tinygo,!wasm,!wasi,!tinygo.wasm

package alloc

import (
	"github.com/moontrade/memory/mem"
	"github.com/moontrade/memory/tlsf"
	"runtime"
	"unsafe"
)

var allocator Allocator

func Scope(fn func(a Auto)) {
	scope(NextAllocator(), fn)
}

func scope(alloc Allocator, fn func(a Auto)) {
	a := NewAuto(alloc, 32)
	defer a.Free()
	fn(a)
	a.Print()
}

var (
	allocatorSlots = runtime.NumCPU() * 2
)

const (
	_TLSFNoSync    Allocator = 1 << 0
	_TLSFSync      Allocator = 1 << 1
	_RPMalloc      Allocator = 1 << 2
	_AllocatorMask           = _TLSFNoSync | _TLSFSync
)

func init() {
	if allocatorSlots > cap(_Allocators) {
		allocatorSlots = cap(_Allocators)
	}
	for i := 0; i < allocatorSlots; i++ {
		a := tlsf.NewHeap(1)
		s := a.ToSync()
		a.Slot = uint8(i)
		s.Slot = uint8(i)
		_SyncAllocators[i] = s
		_Allocators[i] = toTLSFSyncAllocator(s)
	}
	if allocatorSlots < cap(_Allocators) {
		// Distribute evenly amongst the remaining slots
		for i := allocatorSlots; i < cap(_Allocators); i++ {
			_SyncAllocators[i] = _SyncAllocators[i%allocatorSlots]
			_Allocators[i] = _Allocators[i%allocatorSlots]
		}
	}
}

var (
	_AllocatorCount uint64
	_SyncAllocators [255]*tlsf.Sync
	_Allocators     [255]Allocator
)

func AllocatorBySlot(slot uint8) Allocator {
	return _Allocators[slot]
}
func NextAllocator() Allocator {
	_AllocatorCount++
	return _Allocators[_AllocatorCount%255]
}
func NextAllocatorRandom() Allocator {
	return _Allocators[mem.Fastrand()%255]
}

func toTLSFAllocator(a *tlsf.Heap) Allocator {
	return Allocator(unsafe.Pointer(a)) | _TLSFNoSync
}

func toTLSFSyncAllocator(a *tlsf.Sync) Allocator {
	return Allocator(unsafe.Pointer(a)) | _TLSFSync
}

func (a Allocator) Scope(fn func(a Auto)) {
	scope(a, fn)
}

func (a Allocator) Slot() uint8 {
	if a&_TLSFNoSync != 0 {
		return (*tlsf.Heap)(unsafe.Pointer(a & ^_AllocatorMask)).Slot
	} else {
		return (*tlsf.Sync)(unsafe.Pointer(a & ^_AllocatorMask)).Slot
	}
}

func (a Allocator) Alloc(size uintptr) Pointer {
	if a&_TLSFNoSync != 0 {
		return Pointer((*tlsf.Heap)(unsafe.Pointer(a & ^_AllocatorMask)).Alloc(size))
	} else {
		return Pointer((*tlsf.Sync)(unsafe.Pointer(a & ^_AllocatorMask)).Alloc(size))
	}
}

func (a Allocator) AllocZeroed(size uintptr) Pointer {
	if a&_TLSFNoSync != 0 {
		return Pointer((*tlsf.Heap)(unsafe.Pointer(a & ^_AllocatorMask)).AllocZeroed(size))
	} else {
		return Pointer((*tlsf.Sync)(unsafe.Pointer(a & ^_AllocatorMask)).AllocZeroed(size))
	}
}

func (a Allocator) Realloc(ptr Pointer, size uintptr) Pointer {
	if a&_TLSFNoSync != 0 {
		return Pointer((*tlsf.Heap)(unsafe.Pointer(a & ^_AllocatorMask)).Realloc(uintptr(ptr), size))
	}
	return Pointer((*tlsf.Sync)(unsafe.Pointer(a & ^_AllocatorMask)).Realloc(uintptr(ptr), size))
}

func (a Allocator) Free(ptr Pointer) {
	if a&_TLSFNoSync != 0 {
		(*tlsf.Heap)(unsafe.Pointer(a & ^_AllocatorMask)).Free(uintptr(ptr))
		return
	}
	(*tlsf.Sync)(unsafe.Pointer(a & ^_AllocatorMask)).Free(uintptr(ptr))
}

func (a Allocator) Str(size uintptr) Bytes {
	return NewString(size)
}

func (a Allocator) SizeOf(ptr Pointer) uintptr {
	return tlsf.SizeOf(uintptr(ptr))
}

func (a Allocator) Stats() tlsf.Stats {
	if a&_TLSFNoSync != 0 {
		s := (*tlsf.Heap)(unsafe.Pointer(a & ^_AllocatorMask)).Stats
		return s
		//return *(*GCStats)(unsafe.Pointer(&s))
	} else {
		s := (*tlsf.Sync)(unsafe.Pointer(a & ^_AllocatorMask)).Stats()
		return s
		//return *(*GCStats)(unsafe.Pointer(&s))
	}
}
