//go:build tinygo && tinygo.wasm
// +build tinygo,tinygo.wasm

package alloc

import (
	"github.com/moontrade/memory/tlsf"
	"math"
	"unsafe"
)

const (
	//_Type5 = 0
	_Type8    = 1
	_Type16   = 2
	_Type32   = 3
	_Type64   = 4
	_TypeMask = 7
	_TypeBits = 3

	_8HeaderSize       = 2
	_8HeaderTotalSize  = Pointer(2 + int(tlsf.BlockOverhead))
	_16HeaderSize      = 3
	_16HeaderTotalSize = Pointer(3 + int(tlsf.BlockOverhead))
	_32HeaderSize      = 5
	_32HeaderTotalSize = Pointer(5 + int(tlsf.BlockOverhead))
	_64HeaderSize      = 9
	_64HeaderTotalSize = Pointer(9 + int(tlsf.BlockOverhead))
)

func NewString(size uintptr) Bytes {
	switch {
	case size <= math.MaxUint8:
		return Bytes{allocator_.Alloc(size+_8HeaderSize) + Pointer(_8HeaderSize)}.init(_Type8)
	case size <= math.MaxUint16:
		return Bytes{allocator_.Alloc(size+_16HeaderSize) + Pointer(_16HeaderSize)}.init(_Type16)
	case size <= math.MaxInt32:
		return Bytes{allocator_.Alloc(size+_32HeaderSize) + Pointer(_32HeaderSize)}.init(_Type32)
	default:
		return Bytes{allocator_.Alloc(size+_64HeaderSize) + Pointer(_64HeaderSize)}.init(_Type64)
	}
}

func WrapString(s string) Bytes {
	str := NewString(uintptr(len(s)))
	str.Pointer.SetString(0, s)
	str.setLen(len(s))
	return str
}

func WrapBytes(b []byte) Bytes {
	str := NewString(uintptr(len(b)))
	str.Pointer.SetBytes(0, b)
	str.setLen(len(b))
	return str
}

func (s Bytes) init(flags uint8) Bytes {
	*(*uint8)(unsafe.Pointer(s.Pointer - 1)) = flags
	return s
}

func (s Bytes) Allocator() Allocator {
	return allocator_
}

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
