//go:build !tinygo && !wasm && !wasi && !tinygo.wasm
// +build !tinygo,!wasm,!wasi,!tinygo.wasm

package mem

import (
	"reflect"
	"runtime"
	"unsafe"
)

var (
	allocator      *TLSFSync
	allocatorSlots = runtime.NumCPU() * 2
)

func init() {
	if allocatorSlots > cap(_Allocators) {
		allocatorSlots = cap(_Allocators)
	}
	for i := 0; i < allocatorSlots; i++ {
		allocator = NewTLSF(1).ToSync()
		allocator.slot = uint8(i)
		_Allocators[i] = allocator.AsAllocator()
	}
	if allocatorSlots < cap(_Allocators) {
		// Distribute evenly amongst the remaining slots
		for i := allocatorSlots; i < cap(_Allocators); i++ {
			_Allocators[i] = _Allocators[i%allocatorSlots]
		}
	}
}

var (
	_AllocatorCount uint64
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
	return _Allocators[fastrand()%255]
}

func Scope(fn func(a Auto)) {
	a := NewAuto(allocator.AsAllocator(), 32)
	defer a.Free()
	fn(a)
}

// Scope creates an Auto free list that automatically reclaims memory
// after callback finishes.
func (a *TLSF) Scope(fn func(a Auto)) {
	if fn == nil {
		return
	}
	auto := NewAuto(a.AsAllocator(), 32)
	defer auto.Free()
	fn(auto)
}

// Scope creates an Auto free list that automatically reclaims memory
// after callback finishes.
func (a *TLSFSync) Scope(fn func(a Auto)) {
	if fn == nil {
		return
	}
	auto := NewAuto(a.AsAllocator(), 32)
	defer auto.Free()
	fn(auto)
}

const (
	_TLSFNoSync    Allocator = 1 << 0
	_TLSFSync      Allocator = 1 << 1
	_AllocatorMask           = _TLSFNoSync | _TLSFSync
)

func toTLSFAllocator(a *TLSF) Allocator {
	return Allocator(unsafe.Pointer(a)) | _TLSFNoSync
}

func toTLSFSyncAllocator(a *TLSFSync) Allocator {
	return Allocator(unsafe.Pointer(a)) | _TLSFSync
}

func (a Allocator) Slot() uint8 {
	if a&_TLSFNoSync != 0 {
		return 0
	} else {
		return (*TLSFSync)(unsafe.Pointer(a & ^_AllocatorMask)).slot
	}
}

func (a Allocator) Alloc(size uintptr) Pointer {
	if a&_TLSFNoSync != 0 {
		return (*TLSF)(unsafe.Pointer(a & ^_AllocatorMask)).Alloc(size)
	} else {
		return (*TLSFSync)(unsafe.Pointer(a & ^_AllocatorMask)).Alloc(size)
	}
}

func (a Allocator) AllocZeroed(size uintptr) Pointer {
	if a&_TLSFNoSync != 0 {
		return (*TLSF)(unsafe.Pointer(a & ^_AllocatorMask)).AllocZeroed(size)
	} else {
		return (*TLSFSync)(unsafe.Pointer(a & ^_AllocatorMask)).AllocZeroed(size)
	}
}

func (a Allocator) Realloc(ptr Pointer, size uintptr) Pointer {
	if a&_TLSFNoSync != 0 {
		return (*TLSF)(unsafe.Pointer(a & ^_AllocatorMask)).Realloc(ptr, size)
	}
	return (*TLSFSync)(unsafe.Pointer(a & ^_AllocatorMask)).Realloc(ptr, size)
}

func (a Allocator) Free(ptr Pointer) {
	if a&_TLSFNoSync != 0 {
		(*TLSF)(unsafe.Pointer(a & ^_AllocatorMask)).Free(ptr)
		return
	}
	(*TLSFSync)(unsafe.Pointer(a & ^_AllocatorMask)).Free(ptr)
}

func (a Allocator) Str(size uintptr) Str {
	return newString(a, size)
}

//func (a Allocator) Bytes(length Pointer) Bytes {
//	if a&_TLSFNoSync != 0 {
//		return (*TLSF)(unsafe.Pointer(a & ^_AllocatorMask)).Bytes(length)
//	}
//	return (*TLSFSync)(unsafe.Pointer(a & ^_AllocatorMask)).Bytes(length)
//}
//
//func (a Allocator) BytesCap(length, capacity Pointer) Bytes {
//	if a&_TLSFNoSync != 0 {
//		return (*TLSF)(unsafe.Pointer(a & ^_AllocatorMask)).BytesCap(length, capacity)
//	}
//	return (*TLSFSync)(unsafe.Pointer(a & ^_AllocatorMask)).BytesCap(length, capacity)
//}
//
//func (a Allocator) BytesCapNotCleared(length, capacity Pointer) Bytes {
//	if a&_TLSFNoSync != 0 {
//		return (*TLSF)(unsafe.Pointer(a & ^_AllocatorMask)).BytesCapNotCleared(length, capacity)
//	}
//	return (*TLSFSync)(unsafe.Pointer(a & ^_AllocatorMask)).BytesCapNotCleared(length, capacity)
//}

func (a *TLSF) AsAllocator() Allocator {
	return Allocator(unsafe.Pointer(a)) | _TLSFNoSync
}

func memcmp(a, b unsafe.Pointer, n uintptr) int {
	if a == nil {
		if b == nil {
			return 0
		}
		return -1
	}
	ab := *(*string)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(a),
		Len:  int(n),
	}))
	bb := *(*string)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(b),
		Len:  int(n),
	}))
	if ab < bb {
		return -1
	}
	if ab == bb {
		return 0
	}
	return 1
}

type GrowFactory func(arena Arena) Grow

func GrowByDouble(arena Arena) Grow {
	return func(pagesBefore, pagesNeeded int32, minSize Pointer) (pagesAdded int32, start, end Pointer) {
		if pagesBefore > pagesNeeded {
			pagesAdded = pagesBefore
		} else {
			pagesAdded = pagesNeeded
		}
		start, end = arena.Alloc(Pointer(pagesAdded) * _TLSFPageSize)
		if start == 0 {
			pagesAdded = pagesNeeded
			start, end = arena.Alloc(Pointer(pagesAdded) * _TLSFPageSize)
			if start == 0 {
				return 0, 0, 0
			}
		}
		return
	}
}

func GrowBy(pages int32, arena Arena) Grow {
	return func(pagesBefore, pagesNeeded int32, minSize Pointer) (pagesAdded int32, start, end Pointer) {
		if pages > pagesNeeded {
			pagesAdded = pages
		} else {
			pagesAdded = pagesNeeded
		}
		start, end = arena.Alloc(Pointer(pagesAdded) * _TLSFPageSize)
		if start == 0 {
			pagesAdded = pagesNeeded
			start, end = arena.Alloc(Pointer(pagesAdded) * _TLSFPageSize)
			if start == 0 {
				return 0, 0, 0
			}
		}
		return
	}
}

func GrowMin(arena Arena) Grow {
	return func(pagesBefore, pagesNeeded int32, minSize Pointer) (int32, Pointer, Pointer) {
		start, end := arena.Alloc(Pointer(pagesNeeded) * _TLSFPageSize)
		if start == 0 {
			return 0, 0, 0
		}
		return pagesNeeded, start, end
	}
}

func memcpy(dst, src unsafe.Pointer, n uintptr) {
	memmove(dst, src, n)
	//memcpySlow(dst, src, n)
}

func memcpySlow(dst, src unsafe.Pointer, n uintptr) {
	//dstB := *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
	//	Data: uintptr(dst),
	//	Len:  int(n),
	//	Cap:  int(n),
	//}))
	//srcB := *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
	//	Data: uintptr(src),
	//	Len:  int(n),
	//	Cap:  int(n),
	//}))
	//copy(dstB, srcB)
}

// memmove copies n bytes from "from" to "to".
//
// memmove ensures that any pointer in "from" is written to "to" with
// an indivisible write, so that racy reads cannot observe a
// half-written pointer. This is necessary to prevent the garbage
// collector from observing invalid pointers, and differs from memmove
// in unmanaged languages. However, memmove is only required to do
// this if "from" and "to" may contain pointers, which can only be the
// case if "from", "to", and "n" are all be word-aligned.
//
// Implementations are in memmove_*.s.
//
//go:noescape
//go:linkname memmove runtime.memmove
func memmove(to, from unsafe.Pointer, n uintptr)

func memzero(ptr unsafe.Pointer, n uintptr) {
	memclrNoHeapPointers(ptr, n)
	//memzeroSlow(ptr, n)
}

func memzeroSlow(ptr unsafe.Pointer, size uintptr) {
	b := *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(ptr),
		Len:  int(size),
		Cap:  int(size),
	}))
	switch {
	case size < 8:
		for i := 0; i < len(b); i++ {
			b[i] = 0
		}
	case size == 8:
		*(*uint64)(unsafe.Pointer(&b[0])) = 0
	default:
		var i = 0
		for ; i <= len(b)-8; i += 8 {
			*(*uint64)(unsafe.Pointer(&b[i])) = 0
		}
		for ; i < len(b); i++ {
			b[i] = 0
		}
	}
}

// memclrNoHeapPointers clears n bytes starting at ptr.
//
// Usually you should use typedmemclr. memclrNoHeapPointers should be
// used only when the caller knows that *ptr contains no heap pointers
// because either:
//
// *ptr is initialized memory and its type is pointer-free, or
//
// *ptr is uninitialized memory (e.g., memory that's being reused
// for a new allocation) and hence contains only "junk".
//
// memclrNoHeapPointers ensures that if ptr is pointer-aligned, and n
// is a multiple of the pointer size, then any pointer-aligned,
// pointer-sized portion is cleared atomically. Despite the function
// name, this is necessary because this function is the underlying
// implementation of typedmemclr and memclrHasPointers. See the doc of
// memmove for more details.
//
// The (CPU-specific) implementations of this function are in memclr_*.s.
//
//go:noescape
//go:linkname memclrNoHeapPointers runtime.memclrNoHeapPointers
func memclrNoHeapPointers(ptr unsafe.Pointer, n uintptr)

//func memequal(a, b unsafe.Pointer, n uintptr) bool {
//	if a == nil {
//		return b == nil
//	}
//	return *(*string)(unsafe.Pointer(&reflect.SliceHeader{
//		Data: uintptr(a),
//		Len:  int(n),
//	})) == *(*string)(unsafe.Pointer(&reflect.SliceHeader{
//		Data: uintptr(b),
//		Len:  int(n),
//	}))
//}

//go:linkname memequal runtime.memequal
func memequal(a, b unsafe.Pointer, size uintptr) bool

//go:linkname fastrand runtime.fastrand
func fastrand() uint32
