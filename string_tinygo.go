//go:build tinygo && tinygo.wasm
// +build tinygo,tinygo.wasm

package mem

import (
	"math"
	"unsafe"
)

const (
	//_SDSType5 = 0
	_SDSType8    = 1
	_SDSType16   = 2
	_SDSType32   = 3
	_SDSType64   = 4
	_SDSTypeMask = 7
	_SDSTypeBits = 3

	_SDS8HeaderSize       = 2
	_SDS8HeaderTotalSize  = Str(2 + int(_TLSFBlockOverhead))
	_SDS16HeaderSize      = 3
	_SDS16HeaderTotalSize = Str(3 + int(_TLSFBlockOverhead))
	_SDS32HeaderSize      = 5
	_SDS32HeaderTotalSize = Str(5 + int(_TLSFBlockOverhead))
	_SDS64HeaderSize      = 9
	_SDS64HeaderTotalSize = Str(9 + int(_TLSFBlockOverhead))
)

func AllocString(size int) Str {
	return newString(NextAllocator(), size)
}

func newString(a Allocator, size int) Str {
	switch {
	case size <= math.MaxUint8:
		return Str(a.Alloc(Pointer(size+_SDS8HeaderSize)) + Pointer(_SDS8HeaderSize)).init(_SDSType8)
	case size <= math.MaxUint16:
		return Str(a.Alloc(Pointer(size+_SDS16HeaderSize)) + Pointer(_SDS16HeaderSize)).init(_SDSType16)
	case size <= math.MaxInt32:
		return Str(a.Alloc(Pointer(size+_SDS32HeaderSize)) + Pointer(_SDS32HeaderSize)).init(_SDSType32)
	default:
		return Str(a.Alloc(Pointer(size+_SDS64HeaderSize)) + Pointer(_SDS64HeaderSize)).init(_SDSType64)
	}
}

func (s Str) init(flags uint8) Str {
	*(*uint8)(unsafe.Pointer(s - 1)) = flags
	return s
}

func (s Str) Allocator() *TLSF {
	return allocator
}
