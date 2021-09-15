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
	return allocator.AllocNotCleared(size)
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
	return (*TLSF)(unsafe.Pointer(a)).AllocNotCleared(size)
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

//go:export memcpy
func memcpy(dst, src unsafe.Pointer, n uintptr)

//export memzero
func memzero(ptr unsafe.Pointer, size uintptr)

// TODO: Faster LLVM way?
func memequal(a, b unsafe.Pointer, n uintptr) bool {
	if a == nil {
		return b == nil
	}
	return *(*string)(unsafe.Pointer(&_string{
		ptr: uintptr(a),
		len: int(n),
	})) == *(*string)(unsafe.Pointer(&_string{
		ptr: uintptr(b),
		len: int(n),
	}))
}

//export heapAlloc
func heapAlloc(size uintptr) unsafe.Pointer {
	//println("heapAlloc")
	return unsafe.Pointer(allocator.Alloc(Pointer(size)))
}

//go:export initAllocator
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

// GrowBy will Grow by the number of pages specified or by the minimum needed, whichever is greater.
func GrowBy(pages int32) Grow {
	return func(pagesBefore, pagesNeeded int32, minSize Pointer) (pagesAdded int32, start, end Pointer) {
		if pages > pagesNeeded {
			pagesAdded = pages
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
}

// GrowByMin will Grow by a single page or by the minimum needed, whichever is greater.
func GrowMin(pagesBefore, pagesNeeded int32, minSize Pointer) (int32, Pointer, Pointer) {
	start, end := growBy(pagesNeeded)
	if start == 0 {
		return 0, 0, 0
	}
	return pagesNeeded, start, end
}
