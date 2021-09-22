//go:build tinygo && tinygo.wasm
// +build tinygo,tinygo.wasm

package mem

import (
	"unsafe"
)

////////////////////////////////////////////////////////////////////////////////////
// Global Allocator
////////////////////////////////////////////////////////////////////////////////////

var allocator *TLSF
var allocator_ Allocator

func TLSFInstance() *TLSF {
	return allocator
}

func NextAllocator() Allocator {
	return allocator.AsAllocator()
}

////////////////////////////////////////////////////////////////////////////////////
// Global allocator convenience
////////////////////////////////////////////////////////////////////////////////////

// Alloc calls Alloc on the system allocator
func Alloc(size uintptr) Pointer {
	return allocator.Alloc(size)
}

// Alloc calls Alloc on the system allocator
//export alloc
func AllocZeroed(size uintptr) Pointer {
	return allocator.AllocZeroed(size)
}

// Realloc calls Realloc on the system allocator
//export realloc
func Realloc(p Pointer, size uintptr) Pointer {
	return allocator.Realloc(p, size)
}

// Free calls Free on the system allocator
//export free
func Free(p Pointer) {
	allocator.Free(p)
}

//// Alloc calls Alloc on the system allocator
//func AllocBytes(length uintptr) Bytes {
//	return allocator.Bytes(length)
//}
//
//// Alloc calls Alloc on the system allocator
//func AllocBytesCap(length, capacity uintptr) Bytes {
//	return allocator.BytesCapNotCleared(length, capacity)
//}
//
//// Alloc calls Alloc on the system allocator
//func AllocBytesCapNotCleared(length, capacity uintptr) Bytes {
//	return allocator.BytesCapNotCleared(length, capacity)
//}

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

func extalloc(size uintptr) unsafe.Pointer {
	ptr := unsafe.Pointer(allocator.Alloc(size))
	//println("extalloc", uint(uintptr(ptr)))
	return ptr
}

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

func (a Allocator) Alloc(size uintptr) Pointer {
	return (*TLSF)(unsafe.Pointer(a)).Alloc(size)
}

func (a Allocator) AllocZeroed(size uintptr) Pointer {
	return (*TLSF)(unsafe.Pointer(a)).AllocZeroed(size)
}

func (a Allocator) Realloc(ptr Pointer, size uintptr) Pointer {
	return (*TLSF)(unsafe.Pointer(a)).Realloc(ptr, size)
}

func (a Allocator) Free(ptr Pointer) {
	(*TLSF)(unsafe.Pointer(a)).Free(ptr)
}

func (a Allocator) Str(size uintptr) Str {
	return AllocString(size)
}

////////////////////////////////////////////////////////////////////////////////////
// tinygo hooks
////////////////////////////////////////////////////////////////////////////////////

const wasmPageSize = 64 * 1024

func WASMMemorySize(index int32) int32 {
	return wasm_memory_size(index)
}

func WASMMemorySizeBytes(index int32) int32 {
	return wasm_memory_size(index) * wasmPageSize
}

//export llvm.wasm.memory.size.i32
func wasm_memory_size(index int32) int32

//export llvm.wasm.memory.grow.i32
func wasm_memory_grow(index, pages int32) int32

func growBy(pages int32) (Pointer, Pointer) {
	before := wasm_memory_size(0)
	wasm_memory_grow(0, pages)
	after := wasm_memory_size(0)
	if before == after {
		return 0, 0
	}
	return Pointer(before * wasmPageSize), Pointer(after * wasmPageSize)
}

////go:linkname gcMemcpy runtime.gcMemcpy
//func gcMemcpy(dst, src unsafe.Pointer, n uintptr)

//go:linkname memcpy runtime.gcMemcpy
func memcpy(dst, src unsafe.Pointer, n uintptr)

//go:linkname memmove runtime.gcMemmove
func memmove(dst, src unsafe.Pointer, size uintptr)

//func memmove(dst, src unsafe.Pointer, size uintptr) {
//	gcMemmove(dst, src, size)
//}

//go:linkname memzero runtime.gcZero
func memzero(ptr unsafe.Pointer, size uintptr)

//func memzero(ptr unsafe.Pointer, size uintptr) {
//	gcZero(ptr, size)
//}

func memzeroSlow(ptr unsafe.Pointer, size uintptr) {
	b := *(*[]byte)(unsafe.Pointer(&_bytes{
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

//go:linkname memequal runtime.gcMemequal
func memequal(x, y unsafe.Pointer, size uintptr) bool

//func memequal(x, y unsafe.Pointer, size uintptr) bool {
//	return gcMemequal(x, y, size)
//	//if a == nil {
//	//	return b == nil
//	//}
//	//return *(*string)(unsafe.Pointer(&_string{
//	//	ptr: uintptr(a),
//	//	len: int(n),
//	//})) == *(*string)(unsafe.Pointer(&_string{
//	//	ptr: uintptr(b),
//	//	len: int(n),
//	//}))
//}

//go:linkname initAllocator runtime.initAllocator
func initAllocator(heapStart, heapEnd uintptr) {
	//println("initAllocator", uint(heapStart))
	//println("globals", uint(globalsStart), uint(globalsEnd))

	// Did we get called twice?
	if allocator != nil {
		return
	}

	var (
		pagesBefore = wasm_memory_size(0)
		pagesAfter  = pagesBefore
	)

	// Need to add a page initially to fit minimum size required by allocator?
	if heapStart == 0 || Pointer(heapStart+unsafe.Sizeof(TLSF{}))+_TLSFRootSize+_TLSFALMask+16 >
		Pointer(uintptr(pagesBefore)*uintptr(wasmPageSize)) {
		// Just need a single page. Root size is much smaller than a single WASM page.
		wasm_memory_grow(0, 1)
		pagesAfter = wasm_memory_size(0)
	}

	// Bootstrap allocator which will take over all allocations from now on.
	allocator = bootstrap(
		Pointer(heapStart),
		Pointer(uintptr(pagesAfter)*uintptr(wasmPageSize)),
		1,
		GrowMin,
	)
}

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

// GrowByDouble will double the heap on each Grow
func GrowByDouble(pagesBefore, pagesNeeded int32, minSize Pointer) (pagesAdded int32, start, end Pointer) {
	if pagesBefore > pagesNeeded {
		pagesAdded = pagesBefore
	} else {
		pagesAdded = pagesNeeded
	}
	start, end = growBy(pagesAdded)
	if start == 0 {
		pagesAdded = pagesNeeded
		start, end = growBy(pagesAdded)
		if start == 0 {
			return 0, 0, 0
		}
	}
	return
}

var (
	growByPages int32 = 1
)

// GrowBy will Grow by the number of pages specified or by the minimum needed, whichever is greater.
func GrowBy(pages int32) Grow {
	growByPages = pages
	return doGrowByPages
}

func doGrowByPages(pagesBefore, pagesNeeded int32, minSize Pointer) (pagesAdded int32, start, end Pointer) {
	if growByPages > pagesNeeded {
		pagesAdded = growByPages
	} else {
		pagesAdded = pagesNeeded
	}
	start, end = growBy(pagesAdded)
	if start == 0 {
		pagesAdded = pagesNeeded
		start, end = growBy(pagesAdded)
		if start == 0 {
			return 0, 0, 0
		}
	}
	return
}

// GrowByMin will Grow by a single page or by the minimum needed, whichever is greater.
func GrowMin(pagesBefore, pagesNeeded int32, minSize Pointer) (int32, Pointer, Pointer) {
	start, end := growBy(pagesNeeded)
	if start == 0 {
		return 0, 0, 0
	}
	return pagesNeeded, start, end
}
