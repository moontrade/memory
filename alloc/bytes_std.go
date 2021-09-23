//go:build !tinygo && !wasm && !wasi && !tinygo.wasm
// +build !tinygo,!wasm,!wasi,!tinygo.wasm

package alloc

import (
	"github.com/moontrade/memory/tlsf"
	"math"
	"unsafe"
)

const (
	//_SDSType5 = 0
	_Type8    = 1
	_Type16   = 2
	_Type32   = 3
	_Type64   = 4
	_TypeMask = 7
	_TypeBits = 3

	_8HeaderSize       = 3
	_8HeaderTotalSize  = Pointer(3 + int(tlsf.BlockOverhead))
	_16HeaderSize      = 4
	_16HeaderTotalSize = Pointer(4 + int(tlsf.BlockOverhead))
	_32HeaderSize      = 6
	_32HeaderTotalSize = Pointer(6 + int(tlsf.BlockOverhead))
	_64HeaderSize      = 10
	_64HeaderTotalSize = Pointer(10 + int(tlsf.BlockOverhead))
)

func NewString(size uintptr) Bytes {
	return newString(NextAllocator(), size)
}

func WrapString(s string) Bytes {
	str := newString(NextAllocator(), uintptr(len(s)))
	str.Pointer.SetString(0, s)
	str.setLen(len(s))
	return str
}

func WrapBytes(b []byte) Bytes {
	str := newString(NextAllocator(), uintptr(len(b)))
	str.Pointer.SetBytes(0, b)
	str.setLen(len(b))
	return str
}

func NewStringWith(a Allocator, size uintptr) Bytes {
	return newString(a, size)
}

func newString(a Allocator, size uintptr) Bytes {
	switch {
	case size <= math.MaxUint8:
		return Bytes{a.Alloc(size+_8HeaderSize) + Pointer(_8HeaderSize)}.init(a.Slot(), _Type8)
	case size <= math.MaxUint16:
		return Bytes{a.Alloc(size+_16HeaderSize) + Pointer(_16HeaderSize)}.init(a.Slot(), _Type16)
	case size <= math.MaxUint32:
		return Bytes{a.Alloc(size+_32HeaderSize) + Pointer(_32HeaderSize)}.init(a.Slot(), _Type32)
	default:
		return Bytes{a.Alloc(size+_64HeaderSize) + Pointer(_64HeaderSize)}.init(a.Slot(), _Type64)
	}
}

func (s Bytes) init(allocator, flags uint8) Bytes {
	*(*uint8)(unsafe.Pointer(s.Pointer - 2)) = allocator
	*(*uint8)(unsafe.Pointer(s.Pointer - 1)) = flags
	return s
}

func (s Bytes) Allocator() Allocator {
	return AllocatorBySlot(*(*uint8)(unsafe.Pointer(s.Pointer - 2)))
	//return Allocators[*(*uint8)(unsafe.Pointer(s.Pointer - 2))]
}

// Go doesn't support packed structs so below are templates of the packed memory layout.
/*
type sdsHeader5 struct {
	tlsf.BLOCK
	alloc byte
	flags byte
	data struct{}
}

type sdsHeader8 struct {
	tlsf.BLOCK
	len uint8
	alloc byte
	flags byte
	data struct{}
}

type sdsHeader16 struct {
	tlsf.BLOCK
	len uint16
	alloc byte
	flags byte
	data struct{}
}

type sdsHeader32 struct {
	tlsf.BLOCK
	len uint32
	alloc byte
	flags byte
	data struct{}
}

type sdsHeader64 struct {
	tlsf.BLOCK
	len uint64
	alloc byte
	flags byte
	data struct{}
}
*/

func (s Bytes) Len() int {
	flags := *(*uint8)(unsafe.Pointer(s.Pointer - 1))
	switch flags & _TypeMask {
	case _Type8:
		return int(*(*uint8)(unsafe.Pointer(s.Pointer - _8HeaderSize)))
	case _Type16:
		return int(*(*uint16)(unsafe.Pointer(s.Pointer - _16HeaderSize)))
	case _Type32:
		return int(*(*uint32)(unsafe.Pointer(s.Pointer - _32HeaderSize)))
	case _Type64:
		return int(*(*uint64)(unsafe.Pointer(s.Pointer - _64HeaderSize)))
	}
	return 0
}

func (s Bytes) setLen(l int) {
	flags := *(*uint8)(unsafe.Pointer(s.Pointer - 1))
	switch flags & _TypeMask {
	case _Type8:
		*(*uint8)(unsafe.Pointer(s.Pointer - _8HeaderSize)) = uint8(l)
	case _Type16:
		*(*uint16)(unsafe.Pointer(s.Pointer - _16HeaderSize)) = uint16(l)
	case _Type32:
		*(*uint32)(unsafe.Pointer(s.Pointer - _32HeaderSize)) = uint32(l)
	case _Type64:
		*(*uint64)(unsafe.Pointer(s.Pointer - _64HeaderSize)) = uint64(l)
	}
}

func (s Bytes) Cap() int {
	flags := *(*uint8)(unsafe.Pointer(s.Pointer - 1))
	switch flags & _TypeMask {
	case _Type8:
		return int((*(*uintptr)(unsafe.Pointer(s.Pointer - _8HeaderTotalSize))) & ^tlsf.TagsMask) - _8HeaderSize
	case _Type16:
		return int((*(*uintptr)(unsafe.Pointer(s.Pointer - _16HeaderTotalSize))) & ^tlsf.TagsMask) - _16HeaderSize
	case _Type32:
		return int((*(*uintptr)(unsafe.Pointer(s.Pointer - _32HeaderTotalSize))) & ^tlsf.TagsMask) - _32HeaderSize
	case _Type64:
		return int((*(*uintptr)(unsafe.Pointer(s.Pointer - _64HeaderTotalSize))) & ^tlsf.TagsMask) - _64HeaderSize
	}
	return 0
}

func (s Bytes) allocationPointer() Pointer {
	flags := *(*uint8)(unsafe.Pointer(s.Pointer - 1))
	switch flags & _TypeMask {
	case _Type8:
		return s.Pointer - _8HeaderSize
	case _Type16:
		return s.Pointer - _16HeaderSize
	case _Type32:
		return s.Pointer - _32HeaderSize
	case _Type64:
		return s.Pointer - _64HeaderSize
	}
	return 0
}
