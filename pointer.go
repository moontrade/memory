package memory

import (
	"github.com/moontrade/memory/hash"
	"unsafe"
)

// Pointer is a wrapper around a raw pointer that is not unsafe.Pointer
// so Go won't confuse it for a potential GC managed pointer.
type Pointer uintptr

func (p Pointer) ToFat(length int) FatPointer {
	return FatPointer{Pointer: p, len: uintptr(length)}
}

//goland:noinspection GoVetUnsafePointer
func (p Pointer) Unsafe() unsafe.Pointer {
	return unsafe.Pointer(p)
}

// Add is Pointer arithmetic.
func (p Pointer) Add(offset int) Pointer {
	return Pointer(uintptr(int(p) + offset))
}

// Free deallocates memory pointed by Pointer
func (p *Pointer) Free() {
	if p == nil || *p == 0 {
		return
	}
	Free(*p)
	*p = 0
}

// SizeOf returns the size of the allocation provided by the platform allocator.
func (p Pointer) SizeOf() uintptr {
	return SizeOf(p)
}

// Clone the memory starting at offset for size number of bytes and return the new Pointer.
func (p Pointer) Clone(offset, size int) Pointer {
	clone := Alloc(uintptr(size))
	p.Copy(offset, size, clone)
	return clone
}

// Zero zeroes out the entire allocation.
func (p Pointer) Zero(size uintptr) {
	Zero(p.Unsafe(), size)
}

// Move does a memmove
//goland:noinspection GoVetUnsafePointer
func (p Pointer) Move(offset, size int, to Pointer) {
	Move(unsafe.Pointer(to), unsafe.Pointer(uintptr(int(p)+offset)), uintptr(size))
}

// Copy does a memcpy
//goland:noinspection GoVetUnsafePointer
func (p Pointer) Copy(offset, size int, to Pointer) {
	Copy(unsafe.Pointer(to), unsafe.Pointer(uintptr(int(p)+offset)), uintptr(size))
}

// Equals does a memequal
//goland:noinspection GoVetUnsafePointer
func (p Pointer) Equals(offset, size int, to Pointer) bool {
	return Equals(unsafe.Pointer(to), unsafe.Pointer(uintptr(int(p)+offset)), uintptr(size))
}

// Compare does a memcmp
//goland:noinspection GoVetUnsafePointer
func (p Pointer) Compare(offset, size int, to Pointer) int {
	return Compare(unsafe.Pointer(to), unsafe.Pointer(uintptr(int(p)+offset)), uintptr(size))
}

///////////////////////////////////////////////////////////////////////////////////////////////
// Byte
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) Int8(offset int) int8 {
	return *(*int8)(unsafe.Pointer(uintptr(int(p) + offset)))
}

//goland:noinspection GoVetUnsafePointer
func (p Pointer) UInt8(offset int) uint8 {
	return *(*uint8)(unsafe.Pointer(uintptr(int(p) + offset)))
}

//goland:noinspection GoVetUnsafePointer
func (p Pointer) Byte(offset int) byte {
	return *(*byte)(unsafe.Pointer(uintptr(int(p) + offset)))
}

///////////////////////////////////////////////////////////////////////////////////////////////
// Set Byte
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) SetInt8(offset int, v int8) {
	*(*int8)(unsafe.Pointer(uintptr(int(p) + offset))) = v
}

//goland:noinspection GoVetUnsafePointer
func (p Pointer) SetUInt8(offset int, v uint8) {
	*(*uint8)(unsafe.Pointer(uintptr(int(p) + offset))) = v
}

//goland:noinspection GoVetUnsafePointer
func (p Pointer) SetByte(offset int, v byte) {
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset))) = v
}

///////////////////////////////////////////////////////////////////////////////////////////////
// Int16 Native Endian
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) Int16(offset int) int16 {
	return *(*int16)(unsafe.Pointer(uintptr(int(p) + offset)))
}

//goland:noinspection GoVetUnsafePointer
func (p Pointer) SetInt16(offset int, v int16) {
	*(*int16)(unsafe.Pointer(uintptr(int(p) + offset))) = v
}

///////////////////////////////////////////////////////////////////////////////////////////////
// UInt16 Native Endian
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) UInt16(offset int) uint16 {
	return *(*uint16)(unsafe.Pointer(uintptr(int(p) + offset)))
}

//goland:noinspection GoVetUnsafePointer
func (p Pointer) SetUInt16(offset int, v uint16) {
	*(*uint16)(unsafe.Pointer(uintptr(int(p) + offset))) = v
}

///////////////////////////////////////////////////////////////////////////////////////////////
// Int32 Native Endian
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) Int32(offset int) int32 {
	return *(*int32)(unsafe.Pointer(uintptr(int(p) + offset)))
}

//goland:noinspection GoVetUnsafePointer
func (p Pointer) SetInt32(offset int, v int32) {
	*(*int32)(unsafe.Pointer(uintptr(int(p) + offset))) = v
}

///////////////////////////////////////////////////////////////////////////////////////////////
// UInt32 Native Endian
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) UInt32(offset int) uint32 {
	return *(*uint32)(unsafe.Pointer(uintptr(int(p) + offset)))
}

//goland:noinspection GoVetUnsafePointer
func (p Pointer) SetUInt32(offset int, v uint32) {
	*(*uint32)(unsafe.Pointer(uintptr(int(p) + offset))) = v
}

///////////////////////////////////////////////////////////////////////////////////////////////
// Int64 Native Endian
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) Int64(offset int) int64 {
	return *(*int64)(unsafe.Pointer(uintptr(int(p) + offset)))
}

//goland:noinspection GoVetUnsafePointer
func (p Pointer) SetInt64(offset int, v int64) {
	*(*int64)(unsafe.Pointer(uintptr(int(p) + offset))) = v
}

///////////////////////////////////////////////////////////////////////////////////////////////
// UInt64 Native Endian
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) UInt64(offset int) uint64 {
	return *(*uint64)(unsafe.Pointer(uintptr(int(p) + offset)))
}

//goland:noinspection GoVetUnsafePointer
func (p Pointer) SetUInt64(offset int, v uint64) {
	*(*uint64)(unsafe.Pointer(uintptr(int(p) + offset))) = v
}

///////////////////////////////////////////////////////////////////////////////////////////////
// Float32 Native Endian
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) Float32(offset int) float32 {
	return *(*float32)(unsafe.Pointer(uintptr(int(p) + offset)))
}

//goland:noinspection GoVetUnsafePointer
func (p Pointer) SetFloat32(offset int, v float32) {
	*(*float32)(unsafe.Pointer(uintptr(int(p) + offset))) = v
}

///////////////////////////////////////////////////////////////////////////////////////////////
// Float64 Native Endian
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) Float64(offset int) float64 {
	return *(*float64)(unsafe.Pointer(uintptr(int(p) + offset)))
}

//goland:noinspection GoVetUnsafePointer
func (p Pointer) SetFloat64(offset int, v float64) {
	*(*float64)(unsafe.Pointer(uintptr(int(p) + offset))) = v
}

///////////////////////////////////////////////////////////////////////////////////////////////
// int
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) Int(offset int) int {
	return *(*int)(unsafe.Pointer(uintptr(int(p) + offset)))
}

//goland:noinspection GoVetUnsafePointer
func (p Pointer) SetInt(offset int, v int) {
	*(*int)(unsafe.Pointer(uintptr(int(p) + offset))) = v
}

///////////////////////////////////////////////////////////////////////////////////////////////
// uint
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) UInt(offset int) uint {
	return *(*uint)(unsafe.Pointer(uintptr(int(p) + offset)))
}

//goland:noinspection GoVetUnsafePointer
func (p Pointer) SetUInt(offset int, v uint) {
	*(*uint)(unsafe.Pointer(uintptr(int(p) + offset))) = v
}

///////////////////////////////////////////////////////////////////////////////////////////////
// uintptr
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) Uintptr(offset int) uintptr {
	return *(*uintptr)(unsafe.Pointer(uintptr(int(p) + offset)))
}

//goland:noinspection GoVetUnsafePointer
func (p Pointer) SetUintptr(offset int, v uintptr) {
	*(*uintptr)(unsafe.Pointer(uintptr(int(p) + offset))) = v
}

///////////////////////////////////////////////////////////////////////////////////////////////
// Pointer
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) Pointer(offset int) Pointer {
	return *(*Pointer)(unsafe.Pointer(uintptr(int(p) + offset)))
}

//goland:noinspection GoVetUnsafePointer
func (p Pointer) SetPointer(offset int, v Pointer) {
	*(*Pointer)(unsafe.Pointer(uintptr(int(p) + offset))) = v
}

///////////////////////////////////////////////////////////////////////////////////////////////
// String
///////////////////////////////////////////////////////////////////////////////////////////////

type _string struct {
	ptr uintptr
	len int
}

//goland:noinspection GoVetUnsafePointer
func (p Pointer) String(offset, size int) string {
	return *(*string)(unsafe.Pointer(&_string{
		ptr: uintptr(int(p) + offset),
		len: size,
	}))
}

func (p Pointer) SetString(offset int, value string) {
	dst := *(*[]byte)(unsafe.Pointer(&_bytes{
		Data: uintptr(int(p) + offset),
		Len:  len(value),
		Cap:  len(value),
	}))
	copy(dst, value)
}

func (p Pointer) SetBytes(offset int, value []byte) {
	dst := *(*[]byte)(unsafe.Pointer(&_bytes{
		Data: uintptr(int(p) + offset),
		Len:  len(value),
		Cap:  len(value),
	}))
	copy(dst, value)
}

///////////////////////////////////////////////////////////////////////////////////////////////
// Byte Slice
///////////////////////////////////////////////////////////////////////////////////////////////

type _bytes struct {
	Data uintptr
	Len  int
	Cap  int
}

//goland:noinspection GoVetUnsafePointer
func (p Pointer) Bytes(offset, length, capacity int) []byte {
	return *(*[]byte)(unsafe.Pointer(&_bytes{
		Data: uintptr(int(p) + offset),
		Len:  length,
		Cap:  capacity,
	}))
}

///////////////////////////////////////////////////////////////////////////////////////////////
// Int24 Little Endian
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) Int24LE(offset int) int32 {
	return int32(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset)))) |
		int32(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 1))))<<8 |
		int32(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 2))))<<16
}

//goland:noinspection GoVetUnsafePointer
func (p Pointer) SetInt24LE(offset int, v int32) {
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset))) = byte(v)
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 1))) = byte(v >> 8)
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 2))) = byte(v >> 16)
}

///////////////////////////////////////////////////////////////////////////////////////////////
// Int24 Big Endian
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) Int24BE(offset int) int32 {
	return int32(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 2)))) |
		int32(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 1))))<<8 |
		int32(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset))))<<16
}

//goland:noinspection GoVetUnsafePointer
func (p Pointer) SetInt24BE(offset int, v int32) {
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset))) = byte(v >> 16)
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 1))) = byte(v >> 8)
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 2))) = byte(v)
}

///////////////////////////////////////////////////////////////////////////////////////////////
// UInt24 Little Endian
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) UInt24LE(offset int) uint32 {
	return uint32(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset)))) |
		uint32(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 1))))<<8 |
		uint32(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 2))))<<16
}

//goland:noinspection GoVetUnsafePointer
func (p Pointer) SetUInt24LE(offset int, v uint32) {
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset))) = byte(v)
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 1))) = byte(v >> 8)
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 2))) = byte(v >> 16)
}

///////////////////////////////////////////////////////////////////////////////////////////////
// UInt24 Big Endian
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) UInt24BE(offset int) uint32 {
	return uint32(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 2)))) |
		uint32(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 1))))<<8 |
		uint32(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset))))<<16
}

//goland:noinspection GoVetUnsafePointer
func (p Pointer) SetUInt24BE(offset int, v uint32) {
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset))) = byte(v >> 16)
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 1))) = byte(v >> 8)
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 2))) = byte(v)
}

///////////////////////////////////////////////////////////////////////////////////////////////
// Int40 Little Endian
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) Int40LE(offset int) int64 {
	return int64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset)))) |
		int64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 1))))<<8 |
		int64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 2))))<<16 |
		int64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 3))))<<24 |
		int64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 4))))<<32
}

//goland:noinspection GoVetUnsafePointer
func (p Pointer) SetInt40LE(offset int, v int64) {
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset))) = byte(v)
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 1))) = byte(v >> 8)
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 2))) = byte(v >> 16)
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 3))) = byte(v >> 24)
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 4))) = byte(v >> 32)
}

///////////////////////////////////////////////////////////////////////////////////////////////
// Int40 Big Endian
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) Int40BE(offset int) int64 {
	return int64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 4)))) |
		int64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 3))))<<8 |
		int64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 2))))<<16 |
		int64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 1))))<<24 |
		int64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset))))<<32
}

//goland:noinspection GoVetUnsafePointer
func (p Pointer) SetInt40BE(offset int, v int64) {
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset))) = byte(v >> 32)
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 1))) = byte(v >> 24)
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 2))) = byte(v >> 16)
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 3))) = byte(v >> 8)
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 4))) = byte(v)
}

///////////////////////////////////////////////////////////////////////////////////////////////
// UInt40 Little Endian
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) UInt40LE(offset int) uint64 {
	return uint64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset)))) |
		uint64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 1))))<<8 |
		uint64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 2))))<<16 |
		uint64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 3))))<<24 |
		uint64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 4))))<<32
}

//goland:noinspection GoVetUnsafePointer
func (p Pointer) SetUInt40LE(offset int, v uint64) {
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset))) = byte(v)
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 1))) = byte(v >> 8)
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 2))) = byte(v >> 16)
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 3))) = byte(v >> 24)
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 4))) = byte(v >> 32)
}

///////////////////////////////////////////////////////////////////////////////////////////////
// UInt40 Big Endian
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) UInt40BE(offset int) uint64 {
	return uint64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 4)))) |
		uint64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 3))))<<8 |
		uint64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 2))))<<16 |
		uint64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 1))))<<24 |
		uint64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset))))<<32
}

//goland:noinspection GoVetUnsafePointer
func (p Pointer) SetUInt40BE(offset int, v uint64) {
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset))) = byte(v >> 32)
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 1))) = byte(v >> 24)
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 2))) = byte(v >> 16)
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 3))) = byte(v >> 8)
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 4))) = byte(v)
}

///////////////////////////////////////////////////////////////////////////////////////////////
// Int48 Little Endian
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) Int48LE(offset int) int64 {
	return int64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset)))) |
		int64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 1))))<<8 |
		int64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 2))))<<16 |
		int64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 3))))<<24 |
		int64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 4))))<<32 |
		int64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 5))))<<40
}

//goland:noinspection GoVetUnsafePointer
func (p Pointer) SetInt48LE(offset int, v int64) {
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset))) = byte(v)
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 1))) = byte(v >> 8)
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 2))) = byte(v >> 16)
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 3))) = byte(v >> 24)
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 4))) = byte(v >> 32)
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 5))) = byte(v >> 40)
}

///////////////////////////////////////////////////////////////////////////////////////////////
// Int48 Big Endian
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) Int48BE(offset int) int64 {
	return int64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 5)))) |
		int64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 4))))<<8 |
		int64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 3))))<<16 |
		int64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 2))))<<24 |
		int64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 1))))<<32 |
		int64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset))))<<40
}

//goland:noinspection GoVetUnsafePointer
func (p Pointer) SetInt48BE(offset int, v int64) {
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset))) = byte(v >> 40)
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 1))) = byte(v >> 32)
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 2))) = byte(v >> 24)
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 3))) = byte(v >> 16)
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 4))) = byte(v >> 8)
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 5))) = byte(v)
}

///////////////////////////////////////////////////////////////////////////////////////////////
// UInt48 Little Endian
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) UInt48LE(offset int) uint64 {
	return uint64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset)))) |
		uint64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 1))))<<8 |
		uint64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 2))))<<16 |
		uint64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 3))))<<24 |
		uint64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 4))))<<32 |
		uint64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 5))))<<40
}

//goland:noinspection GoVetUnsafePointer
func (p Pointer) SetUInt48LE(offset int, v uint64) {
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset))) = byte(v)
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 1))) = byte(v >> 8)
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 2))) = byte(v >> 16)
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 3))) = byte(v >> 24)
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 4))) = byte(v >> 32)
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 5))) = byte(v >> 40)
}

///////////////////////////////////////////////////////////////////////////////////////////////
// UInt48 Big Endian
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) UInt48BE(offset int) uint64 {
	return uint64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 5)))) |
		uint64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 4))))<<8 |
		uint64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 3))))<<16 |
		uint64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 2))))<<24 |
		uint64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 1))))<<32 |
		uint64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset))))<<40
}

//goland:noinspection GoVetUnsafePointer
func (p Pointer) SetUInt48BE(offset int, v uint64) {
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset))) = byte(v >> 40)
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 1))) = byte(v >> 32)
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 2))) = byte(v >> 24)
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 3))) = byte(v >> 16)
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 4))) = byte(v >> 8)
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 5))) = byte(v)
}

///////////////////////////////////////////////////////////////////////////////////////////////
// Int56 Little Endian
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) Int56LE(offset int) int64 {
	return int64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset)))) |
		int64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 1))))<<8 |
		int64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 2))))<<16 |
		int64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 3))))<<24 |
		int64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 4))))<<32 |
		int64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 5))))<<40 |
		int64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 6))))<<48
}

//goland:noinspection GoVetUnsafePointer
func (p Pointer) SetInt56LE(offset int, v int64) {
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset))) = byte(v)
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 1))) = byte(v >> 8)
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 2))) = byte(v >> 16)
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 3))) = byte(v >> 24)
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 4))) = byte(v >> 32)
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 5))) = byte(v >> 40)
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 6))) = byte(v >> 48)
}

///////////////////////////////////////////////////////////////////////////////////////////////
// Int56 Big Endian
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) Int56BE(offset int) int64 {
	return int64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 6)))) |
		int64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 5))))<<8 |
		int64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 4))))<<16 |
		int64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 3))))<<24 |
		int64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 2))))<<32 |
		int64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 1))))<<40 |
		int64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset))))<<48
}

//goland:noinspection GoVetUnsafePointer
func (p Pointer) SetInt56BE(offset int, v int64) {
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset))) = byte(v >> 48)
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 1))) = byte(v >> 40)
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 2))) = byte(v >> 32)
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 3))) = byte(v >> 24)
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 4))) = byte(v >> 16)
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 5))) = byte(v >> 8)
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 6))) = byte(v)
}

///////////////////////////////////////////////////////////////////////////////////////////////
// UInt56 Little Endian
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) UInt56LE(offset int) uint64 {
	return uint64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset)))) |
		uint64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 1))))<<8 |
		uint64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 2))))<<16 |
		uint64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 3))))<<24 |
		uint64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 4))))<<32 |
		uint64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 5))))<<40 |
		uint64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 6))))<<48
}

//goland:noinspection GoVetUnsafePointer
func (p Pointer) SetUInt56LE(offset int, v uint64) {
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset))) = byte(v)
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 1))) = byte(v >> 8)
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 2))) = byte(v >> 16)
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 3))) = byte(v >> 24)
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 4))) = byte(v >> 32)
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 5))) = byte(v >> 40)
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 6))) = byte(v >> 48)
}

///////////////////////////////////////////////////////////////////////////////////////////////
// UInt56 Big Endian
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) UInt56BE(offset int) uint64 {
	return uint64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 6)))) |
		uint64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 5))))<<8 |
		uint64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 4))))<<16 |
		uint64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 3))))<<24 |
		uint64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 2))))<<32 |
		uint64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 1))))<<40 |
		uint64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset))))<<48
}

//goland:noinspection GoVetUnsafePointer
func (p Pointer) SetUInt56BE(offset int, v uint64) {
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset))) = byte(v >> 48)
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 1))) = byte(v >> 40)
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 2))) = byte(v >> 32)
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 3))) = byte(v >> 24)
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 4))) = byte(v >> 16)
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 5))) = byte(v >> 8)
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 6))) = byte(v)
}

func (p Pointer) Hash32(length int) uint32 {
	return p.Hash32At(0, length)
}

func (p Pointer) Hash32At(offset, length int) uint32 {
	const (
		offset32 = uint32(2166136261)
		prime32  = uint32(16777619)
	)
	hash := offset32

	start := uintptr(int(p) + offset)
	end := start + uintptr(length)
	for ; start < end; start++ {
		hash ^= uint32(*(*byte)(unsafe.Pointer(start)))
		hash *= prime32
	}
	return hash
}

func (p Pointer) Hash64(length int) uint64 {
	return p.Hash64At(0, length)
}

func (p Pointer) Hash64At(offset, length int) uint64 {
	const (
		offset64 = uint64(14695981039346656037)
		prime64  = uint64(1099511628211)
	)
	hash := offset64
	start := uintptr(int(p) + offset)
	end := start + uintptr(length)
	for ; start < end; start++ {
		hash ^= uint64(*(*byte)(unsafe.Pointer(start)))
		hash *= prime64
	}
	return hash
}

func (p Pointer) WyHash64(seed uint64, offset, length int) uint64 {
	return hash.WyHash(p.Bytes(offset, length, length), seed)
}

func (p Pointer) Metro64(seed uint64, offset, length int) uint64 {
	return hash.Metro(p.Bytes(offset, length, length), seed)
}
