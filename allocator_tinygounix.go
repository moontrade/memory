//go:build tinygo && !tinygo.wasm && (darwin || linux || windows)
// +build tinygo
// +build !tinygo.wasm
// +build darwin linux windows

package mem

import (
	"unsafe"
)

////////////////////////////////////////////////////////////////////////////////////
// Global Allocator
////////////////////////////////////////////////////////////////////////////////////

var allocator *TLSF

func TLSFInstance() *TLSF {
	return allocator
}

////////////////////////////////////////////////////////////////////////////////////
// Global allocator convenience
////////////////////////////////////////////////////////////////////////////////////

// Alloc calls Alloc on the system allocator
func Alloc(size Pointer) Pointer {
	return allocator.Alloc(size)
}

// Alloc calls Alloc on the system allocator
func AllocNotCleared(size Pointer) Pointer {
	return allocator.AllocZeroed(size)
}

// Realloc calls Realloc on the system allocator
func Realloc(p Pointer, size Pointer) Pointer {
	return allocator.Realloc(p, size)
}

// Free calls Free on the system allocator
func Free(p Pointer) {
	allocator.Free(p)
}

// Alloc calls Alloc on the system allocator
func AllocBytes(length Pointer) Bytes {
	return allocator.Bytes(length)
}

// Alloc calls Alloc on the system allocator
func AllocBytesCap(length, capacity Pointer) Bytes {
	return allocator.BytesCapNotCleared(length, capacity)
}

// Alloc calls Alloc on the system allocator
func AllocBytesCapNotCleared(length, capacity Pointer) Bytes {
	return allocator.BytesCapNotCleared(length, capacity)
}

func Scope(fn func(a Auto)) {
	a := NewAuto(allocator.AsAllocator(), 32)
	fn(a)
	a.Free()
}

// Scope creates an Auto free list that automatically reclaims memory
// after callback finishes.
func (a *TLSF) Scope(fn func(a Auto)) {
	auto := NewAuto(a.AsAllocator(), 32)
	fn(auto)
	auto.Free()
}

//go:export extalloc
func extalloc(size uintptr) unsafe.Pointer {
	ptr := unsafe.Pointer(allocator.Alloc(Pointer(size)))
	//println("extalloc", uint(uintptr(ptr)))
	return ptr
}

//go:export extfree
func extfree(ptr unsafe.Pointer) {
	//println("extfree", uint(uintptr(ptr)))
	allocator.Free(Pointer(ptr))
}

////////////////////////////////////////////////////////////////////////////////////
// Allocator facade
////////////////////////////////////////////////////////////////////////////////////

func (a *TLSF) AsAllocator() Allocator {
	return Allocator(unsafe.Pointer(a))
}

func (a Allocator) Alloc(size Pointer) Pointer {
	return (*TLSF)(unsafe.Pointer(a)).Alloc(size)
}

func (a Allocator) AllocNotCleared(size Pointer) Pointer {
	return (*TLSF)(unsafe.Pointer(a)).AllocZeroed(size)
}

func (a Allocator) Realloc(ptr, size Pointer) Pointer {
	return (*TLSF)(unsafe.Pointer(a)).Realloc(ptr, size)
}

func (a Allocator) Free(ptr Pointer) {
	(*TLSF)(unsafe.Pointer(a)).Free(ptr)
}

func (a Allocator) Bytes(length Pointer) Bytes {
	return (*TLSF)(unsafe.Pointer(a)).Bytes(length)
}

func (a Allocator) BytesCap(length, capacity Pointer) Bytes {
	return (*TLSF)(unsafe.Pointer(a)).BytesCap(length, capacity)
}

func (a Allocator) BytesCapNotCleared(length, capacity Pointer) Bytes {
	return (*TLSF)(unsafe.Pointer(a)).BytesCapNotCleared(length, capacity)
}

////////////////////////////////////////////////////////////////////////////////////
// tinygo hooks
////////////////////////////////////////////////////////////////////////////////////

const wasmPageSize = 64 * 1024

//go:linkname malloc runtime.malloc
func malloc(size uintptr) unsafe.Pointer

//go:linkname memcpy runtime.memcpy
func memcpy(dst, src unsafe.Pointer, n uintptr)

////go:linkname memzero runtime.memzero
//func memzero(ptr unsafe.Pointer, size uintptr)

//go:linkname gcZero runtime.gcZero
func gcZero(ptr unsafe.Pointer, size uintptr)

func memzero(ptr unsafe.Pointer, size uintptr) {
	gcZero(ptr, size)
}

//go:linkname gcZero runtime.gcZero
func gcMemequal(x, y unsafe.Pointer, size uintptr) bool

// TODO: Faster LLVM way?
func memequal(a, b unsafe.Pointer, n uintptr) bool {
	return gcMemequal(a, b, n)
	//if a == nil {
	//	return b == nil
	//}
	//return *(*string)(unsafe.Pointer(&_string{
	//	ptr: uintptr(a),
	//	len: int(n),
	//})) == *(*string)(unsafe.Pointer(&_string{
	//	ptr: uintptr(b),
	//	len: int(n),
	//}))
}

//go:linkname initAllocator runtime.initAllocator
func initAllocator(heapStart, heapEnd uintptr) {
	if allocator != nil {
		return
	}
	allocator = newTLSF(1, GrowMin)
}

//func NewTLSF(pages int32) *TLSF {
//	if pages <= 0 {
//		pages = 1
//	}
//
//	size := uintptr(pages * wasmPageSize)
//	start := uintptr(wasm_memory_size(0) * wasmPageSize)
//	wasm_memory_grow(0, pages)
//	end := uintptr(wasm_memory_size(0) * wasmPageSize)
//	if start == end {
//		panic("out of memory")
//	}
//	return initTLSF(start, start+size, pages)
//}
//
//func (a *TLSF) Grow(pages int32) (uintptr, uintptr) {
//	if pages <= 0 {
//		pages = 1
//	}
//
//	// wasm memory grow
//	before := wasm_memory_size(0)
//	wasm_memory_grow(0, pages)
//	after := wasm_memory_size(0)
//	if before == after {
//		return 0, 0
//	}
//
//	a.Pages += int32(pages)
//	start := uintptr(before * wasmPageSize)
//	end := uintptr(after * wasmPageSize)
//	a.HeapEnd = end
//	return start, end
//}

////////////////////////////////////////////////////////////////////////////////////
// Grow Strategy
////////////////////////////////////////////////////////////////////////////////////

func SetGrow(grow Grow) {
	if allocator == nil {
		initAllocator(0, 0)
	}
	if allocator != nil {
		allocator.Grow = grow
	}
}

func GrowByDouble() Grow {
	return func(pagesBefore, pagesNeeded int32, minSize Pointer) (pagesAdded int32, start, end Pointer) {
		if pagesBefore > pagesNeeded {
			pagesAdded = pagesBefore
		} else {
			pagesAdded = pagesNeeded
		}
		ptr := Pointer(malloc(uintptr(pagesAdded) * uintptr(_TLSFPageSize)))
		if ptr == 0 {
			pagesAdded = pagesNeeded
			ptr = Pointer(malloc(uintptr(pagesAdded) * uintptr(_TLSFPageSize)))
			if ptr == 0 {
				return 0, 0, 0
			}
		}
		start = ptr
		end = start + (Pointer(pagesAdded) * _TLSFPageSize)
		return
	}
}

func GrowBy(pages int32) Grow {
	return func(pagesBefore, pagesNeeded int32, minSize Pointer) (pagesAdded int32, start, end Pointer) {
		if pages > pagesNeeded {
			pagesAdded = pages
		} else {
			pagesAdded = pagesNeeded
		}
		ptr := Pointer(malloc(uintptr(pagesAdded) * uintptr(_TLSFPageSize)))
		if ptr == 0 {
			pagesAdded = pagesNeeded
			ptr = Pointer(malloc(uintptr(pagesAdded) * uintptr(_TLSFPageSize)))
			if ptr == 0 {
				return 0, 0, 0
			}
		}
		start = ptr
		end = start + (Pointer(pagesAdded) * _TLSFPageSize)
		return
	}
}

func GrowMin(pagesBefore, pagesNeeded int32, minSize Pointer) (int32, Pointer, Pointer) {
	ptr := Pointer(malloc(uintptr(pagesNeeded) * uintptr(_TLSFPageSize)))
	if ptr == 0 {
		return 0, 0, 0
	}
	return pagesNeeded, ptr, Pointer(ptr) + (Pointer(pagesNeeded) * _TLSFPageSize)
}

func newTLSF(pages int32, grow Grow) *TLSF {
	if pages <= 0 {
		pages = 1
	}
	size := Pointer(pages * wasmPageSize)
	segment := Pointer(malloc(uintptr(size)))
	return bootstrap(segment, segment+size, pages, grow)
}
