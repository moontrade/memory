//go:build !tinygo && !wasm && !wasi && !tinygo.wasm
// +build !tinygo,!wasm,!wasi,!tinygo.wasm

package mem

import (
	"reflect"
	"unsafe"
)

var allocator *TLSFSync

func init() {
	allocator = NewTLSF(1).ToSync()
}

func Scope(fn func(a Auto)) {
	a := NewAuto(allocator.AsAllocator(), 32)
	defer a.Free()
	fn(a)
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

func (a Allocator) Alloc(size Pointer) Pointer {
	if a&_TLSFNoSync != 0 {
		return (*TLSF)(unsafe.Pointer(a & ^_AllocatorMask)).Alloc(size)
	} else {
		return (*TLSFSync)(unsafe.Pointer(a & ^_AllocatorMask)).Alloc(size)
	}
}

func (a Allocator) AllocNotCleared(size Pointer) Pointer {
	if a&_TLSFNoSync != 0 {
		return (*TLSF)(unsafe.Pointer(a & ^_AllocatorMask)).AllocNotCleared(size)
	} else {
		return (*TLSFSync)(unsafe.Pointer(a & ^_AllocatorMask)).AllocNotCleared(size)
	}
}

func (a Allocator) Realloc(ptr, size Pointer) Pointer {
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

func (a Allocator) Bytes(length Pointer) Bytes {
	if a&_TLSFNoSync != 0 {
		return (*TLSF)(unsafe.Pointer(a & ^_AllocatorMask)).Bytes(length)
	}
	return (*TLSFSync)(unsafe.Pointer(a & ^_AllocatorMask)).Bytes(length)
}

func (a Allocator) BytesCap(length, capacity Pointer) Bytes {
	if a&_TLSFNoSync != 0 {
		return (*TLSF)(unsafe.Pointer(a & ^_AllocatorMask)).BytesCap(length, capacity)
	}
	return (*TLSFSync)(unsafe.Pointer(a & ^_AllocatorMask)).BytesCap(length, capacity)
}

func (a Allocator) BytesCapNotCleared(length, capacity Pointer) Bytes {
	if a&_TLSFNoSync != 0 {
		return (*TLSF)(unsafe.Pointer(a & ^_AllocatorMask)).BytesCapNotCleared(length, capacity)
	}
	return (*TLSFSync)(unsafe.Pointer(a & ^_AllocatorMask)).BytesCapNotCleared(length, capacity)
}

func (a *TLSF) AsAllocator() Allocator {
	return Allocator(unsafe.Pointer(a)) | _TLSFNoSync
}

func memequal(a, b unsafe.Pointer, n uintptr) bool {
	if a == nil {
		return b == nil
	}
	return *(*string)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(a),
		Len:  int(n),
	})) == *(*string)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(b),
		Len:  int(n),
	}))
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

func memcpy(dst, src unsafe.Pointer, n uintptr) {
	dstB := *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(dst),
		Len:  int(n),
		Cap:  int(n),
	}))
	srcB := *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(src),
		Len:  int(n),
		Cap:  int(n),
	}))
	copy(dstB, srcB)
}

////go:linkname memclrNoHeapPointers runtime.memclrNoHeapPointers
//func memclrNoHeapPointers(ptr unsafe.Pointer, n uintptr)

func memzero(ptr unsafe.Pointer, n uintptr) {
	//memclrNoHeapPointers(ptr, n)
	memzeroSlow(ptr, n)
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
