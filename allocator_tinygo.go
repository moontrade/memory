//go:build tinygo && tinygo.wasm
// +build tinygo,tinygo.wasm

package mem

import (
	"unsafe"
)

func (a *TLSF) AsAllocator() Allocator {
	return Allocator(unsafe.Pointer(a))
}

func (a Allocator) Alloc(size Pointer) Pointer {
	return (*TLSF)(unsafe.Pointer(a)).Alloc(size)
}

func (a Allocator) AllocZeroed(size Pointer) Pointer {
	return (*TLSF)(unsafe.Pointer(a)).AllocZeroed(size)
}

func (a Allocator) Bytes(length Pointer) Bytes {
	return (*TLSF)(unsafe.Pointer(a)).Bytes(length)
}

func (a Allocator) BytesCapacity(length, capacity Pointer) Bytes {
	return (*TLSF)(unsafe.Pointer(a)).BytesCapacity(length, capacity)
}

func (a Allocator) Realloc(ptr, size Pointer) Pointer {
	return (*TLSF)(unsafe.Pointer(a)).Realloc(ptr, size)
}

func (a Allocator) Free(ptr Pointer) {
	(*TLSF)(unsafe.Pointer(a)).Free(ptr)
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
