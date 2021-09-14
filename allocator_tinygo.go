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
