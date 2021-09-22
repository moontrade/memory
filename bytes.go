package mem

import (
	"unsafe"
)

type BytesSlice struct {
	Bytes
	p Bytes
}

func (p *BytesSlice) Free() {
	// Noop
}

type bytes uintptr

type bytesLayout struct {
	p     Pointer
	len   uintptr
	cap   uintptr
	alloc Allocator
	data  struct{}
}

//goland:noinspection GoVetUnsafePointer
func (b bytes) Data() Pointer {
	return *(*Pointer)(unsafe.Pointer(b))
}

//goland:noinspection GoVetUnsafePointer
func (b bytes) layout() *bytesLayout {
	return (*bytesLayout)(unsafe.Pointer(b))
}

//goland:noinspection GoVetUnsafePointer
func (b bytes) Len() uintptr {
	return b.layout().len
}

//goland:noinspection GoVetUnsafePointer
func (b bytes) Cap() uintptr {
	return b.layout().cap
}

// Bytes is a fat pointer to a heap allocation from an Allocator
type Bytes struct {
	Pointer Pointer
	len     int
	cap     int
	alloc   Allocator
}

//goland:noinspection GoVetUnsafePointer
func (b *Bytes) Free() {
	if b == nil || b.Pointer == 0 {
		return
	}
	b.alloc.Free(b.Pointer)
	*b = Bytes{}
}

func (b *Bytes) Len() int {
	return int(b.len)
}

func (b *Bytes) Cap() int {
	return int(b.cap)
}

func (b *Bytes) Allocator() Allocator {
	return b.alloc
}

func (b *Bytes) Unsafe() unsafe.Pointer {
	return b.Pointer.Unsafe()
}

// Reset zeroes out the entire allocation and sets the length back to 0
func (b *Bytes) Reset() {
	memzero(b.Unsafe(), uintptr(b.cap))
	b.len = 0
}

// Zero zeroes out the entire allocation.
func (b *Bytes) Zero() {
	memzero(b.Unsafe(), uintptr(b.cap))
}

//goland:noinspection GoVetUnsafePointer
func (b *Bytes) Equals(o *Bytes) bool {
	return b.len == o.len && (b.Pointer == o.Pointer || memequal(
		unsafe.Pointer(b.Pointer),
		unsafe.Pointer(o.Pointer), uintptr(o.len)))
}

//goland:noinspection GoVetUnsafePointer
func (b *Bytes) EqualsSlice(s *BytesSlice) bool {
	return b.len == s.len && (b.Pointer == s.Pointer || memequal(
		unsafe.Pointer(b.Pointer),
		unsafe.Pointer(s.Pointer), uintptr(s.len)))
}

func (b *Bytes) String() string {
	if b.IsEmpty() {
		return ""
	}
	return b.Pointer.String(0, int(b.len))
}

func (b *Bytes) Substring(offset, length int) string {
	if b.IsEmpty() {
		return ""
	}
	if offset < 0 {
		offset = 0
	}
	if length+offset > int(b.cap) {
		length = int(b.cap) - offset
	}
	if length <= 0 {
		return ""
	}
	return *(*string)(unsafe.Pointer(&_string{
		ptr: uintptr(b.Pointer) + uintptr(offset),
		len: length,
	}))
}

func (b *Bytes) Bytes() []byte {
	if b.IsNil() {
		return nil
	}
	return b.Pointer.Bytes(0, int(b.len), int(b.len))
}

func (b *Bytes) IsNil() bool {
	return b == nil || b.Pointer == 0
}

func (b *Bytes) IsEmpty() bool {
	return b.Pointer == 0 || b.len == 0
}

func (b *Bytes) CheckBounds(offset int) bool {
	return uintptr(b.Pointer) == 0 || int(b.len) < offset
}

func (b *Bytes) Slice(offset, length int) BytesSlice {
	if b.IsNil() {
		return BytesSlice{}
	}
	if offset+length > b.len {
		return BytesSlice{}
	}
	return BytesSlice{
		Bytes: Bytes{
			Pointer: Pointer(int(b.Pointer) + offset),
			len:     length,
			cap:     b.cap - offset,
		},
		p: *b,
	}
}

///////////////////////////////////////////////////////////////////////////////////////////
// Int
///////////////////////////////////////////////////////////////////////////////////////////

func (b *Bytes) Int(offset int) int {
	b.EnsureLen(offset + int(unsafe.Sizeof(int(0))))
	return b.Pointer.Int(offset)
}

func (b *Bytes) SetInt(offset int, value int) {
	b.EnsureLen(offset + int(unsafe.Sizeof(int(0))))
	b.Pointer.SetInt(offset, value)
}

//goland:noinspection GoVetUnsafePointer
func (b *Bytes) AppendInt(value int) {
	b.Pointer.SetInt(b.ensureAppend(int(unsafe.Sizeof(int(0)))), value)
}

///////////////////////////////////////////////////////////////////////////////////////////
// UInt
///////////////////////////////////////////////////////////////////////////////////////////

func (b *Bytes) UInt(offset int) int {
	b.EnsureLen(offset + int(unsafe.Sizeof(int(0))))
	return b.Pointer.Int(offset)
}

func (b *Bytes) SetUInt(offset int, value uint) {
	b.EnsureLen(offset + int(unsafe.Sizeof(uint(0))))
	b.Pointer.SetUInt(offset, value)
}

//goland:noinspection GoVetUnsafePointer
func (b *Bytes) AppendUInt(value uint) {
	b.Pointer.SetUInt(b.ensureAppend(int(unsafe.Sizeof(uint(0)))), value)
}

///////////////////////////////////////////////////////////////////////////////////////////
// Pointer
///////////////////////////////////////////////////////////////////////////////////////////

func (b *Bytes) PointerAt(offset int) Pointer {
	b.EnsureLen(offset + 8)
	return b.Pointer.Pointer(offset)
}

func (b *Bytes) SetPointer(offset int, value Pointer) {
	b.EnsureLen(offset + int(unsafe.Sizeof(uintptr(0))))
	b.Pointer.SetPointer(offset, value)
}

func (b *Bytes) AppendPointer(value Pointer) {
	b.Pointer.SetPointer(b.ensureAppend(int(unsafe.Sizeof(Pointer(0)))), value)
}

///////////////////////////////////////////////////////////////////////////////////////////
// uintptr
///////////////////////////////////////////////////////////////////////////////////////////

func (b *Bytes) Uintptr(offset int) uintptr {
	b.EnsureLen(offset + int(unsafe.Sizeof(uintptr(0))))
	return b.Pointer.Uintptr(offset)
}

func (b *Bytes) SetUintptr(offset int, value uintptr) {
	b.EnsureLen(offset + int(unsafe.Sizeof(uintptr(0))))
	b.Pointer.SetUintptr(offset, value)
}

func (b *Bytes) AppendUintptr(value uintptr) {
	b.Pointer.SetUintptr(b.ensureAppend(int(unsafe.Sizeof(uintptr(0)))), value)
}

///////////////////////////////////////////////////////////////////////////////////////////
// Int8
///////////////////////////////////////////////////////////////////////////////////////////

func (b *Bytes) Int8(offset int) int8 {
	b.EnsureLen(offset + 1)
	return b.Pointer.Int8(offset)
}

func (b *Bytes) SetInt8(offset int, value int8) {
	b.EnsureLen(offset + 1)
	b.Pointer.SetInt8(offset, value)
}

func (b *Bytes) AppendInt8(value int8) {
	b.Pointer.SetInt8(b.ensureAppend(1), value)
}

///////////////////////////////////////////////////////////////////////////////////////////
// UInt8
///////////////////////////////////////////////////////////////////////////////////////////

func (b *Bytes) UInt8(offset int) uint8 {
	b.EnsureLen(offset + 1)
	return b.Pointer.UInt8(offset)
}

// SetUInt8 is safe version. Will grow allocation if needed.
func (b *Bytes) SetUInt8(offset int, value uint8) {
	b.EnsureLen(offset + 1)
	b.Pointer.SetUInt8(offset, value)
}

func (b *Bytes) AppendUInt8(value uint8) {
	b.Pointer.SetUInt8(b.ensureAppend(1), value)
}

///////////////////////////////////////////////////////////////////////////////////////////
// Byte
///////////////////////////////////////////////////////////////////////////////////////////

func (b *Bytes) Byte(offset int) byte {
	b.EnsureLen(offset + 1)
	return b.Pointer.Byte(offset)
}

func (b *Bytes) SetByte(offset int, value byte) {
	b.EnsureLen(offset + 1)
	b.Pointer.SetByte(offset, value)
}

func (b *Bytes) AppendByte(value byte) {
	b.Pointer.SetByte(b.ensureAppend(1), value)
}

///////////////////////////////////////////////////////////////////////////////////////////
// Int16 Native Endian
///////////////////////////////////////////////////////////////////////////////////////////

func (b *Bytes) Int16(offset int) int16 {
	b.EnsureLen(offset + 2)
	return b.Pointer.Int16(offset)
}

func (b *Bytes) SetInt16(offset int, value int16) {
	b.EnsureLen(offset + 2)
	b.Pointer.SetInt16(offset, value)
}

func (b *Bytes) AppendInt16(value int16) {
	b.Pointer.SetInt16(b.ensureAppend(2), value)
}

///////////////////////////////////////////////////////////////////////////////////////////
// Int16 Little Endian
///////////////////////////////////////////////////////////////////////////////////////////

func (b *Bytes) Int16LE(offset int) int16 {
	b.EnsureLen(offset + 2)
	return b.Pointer.Int16LE(offset)
}

func (b *Bytes) SetInt16LE(offset int, value int16) {
	b.EnsureLen(offset + 2)
	b.Pointer.SetInt16LE(offset, value)
}

func (b *Bytes) AppendInt16LE(value int16) {
	b.Pointer.SetInt16LE(b.ensureAppend(2), value)
}

///////////////////////////////////////////////////////////////////////////////////////////
// Int16 Little Endian
///////////////////////////////////////////////////////////////////////////////////////////

func (b *Bytes) Int16BE(offset int) int16 {
	b.EnsureLen(offset + 2)
	return b.Pointer.Int16BE(offset)
}

func (b *Bytes) SetInt16BE(offset int, value int16) {
	b.EnsureLen(offset + 2)
	b.Pointer.SetInt16BE(offset, value)
}

func (b *Bytes) AppendInt16BE(value int16) {
	b.Pointer.SetInt16BE(b.ensureAppend(2), value)
}

///////////////////////////////////////////////////////////////////////////////////////////
// UInt16 Native Endian
///////////////////////////////////////////////////////////////////////////////////////////

func (b *Bytes) UInt16(offset int) uint16 {
	b.EnsureLen(offset + 2)
	return b.Pointer.UInt16(offset)
}

func (b *Bytes) SetUInt16(offset int, value uint16) {
	b.EnsureLen(offset + 2)
	b.Pointer.SetUInt16(offset, value)
}

func (b *Bytes) AppendUInt16(value uint16) {
	b.Pointer.SetUInt16(b.ensureAppend(2), value)
}

///////////////////////////////////////////////////////////////////////////////////////////
// UInt16 Little Endian
///////////////////////////////////////////////////////////////////////////////////////////

func (b *Bytes) UInt16LE(offset int) uint16 {
	b.EnsureLen(offset + 2)
	return b.Pointer.UInt16LE(offset)
}

func (b *Bytes) SetUInt16LE(offset int, value uint16) {
	b.EnsureLen(offset + 2)
	b.Pointer.SetUInt16LE(offset, value)
}

func (b *Bytes) AppendUInt16LE(value uint16) {
	b.Pointer.SetUInt16LE(b.ensureAppend(2), value)
}

///////////////////////////////////////////////////////////////////////////////////////////
// UInt16 Big Endian
///////////////////////////////////////////////////////////////////////////////////////////

func (b *Bytes) UInt16BE(offset int) uint16 {
	b.EnsureLen(offset + 2)
	return b.Pointer.UInt16BE(offset)
}

func (b *Bytes) SetUInt16BE(offset int, value uint16) {
	b.EnsureLen(offset + 2)
	b.Pointer.SetUInt16BE(offset, value)
}

func (b *Bytes) AppendUInt16BE(value uint16) {
	b.Pointer.SetUInt16BE(b.ensureAppend(2), value)
}

///////////////////////////////////////////////////////////////////////////////////////////
// Int32 Native Endian
///////////////////////////////////////////////////////////////////////////////////////////

func (b *Bytes) Int32(offset int) int32 {
	b.EnsureLen(offset + 4)
	return b.Pointer.Int32(offset)
}

func (b *Bytes) SetInt32(offset int, value int32) {
	b.EnsureLen(offset + 4)
	b.Pointer.SetInt32(offset, value)
}

func (b *Bytes) AppendInt32(value int32) {
	b.Pointer.SetInt32(b.ensureAppend(4), value)
}

///////////////////////////////////////////////////////////////////////////////////////////
// Int32 Little Endian
///////////////////////////////////////////////////////////////////////////////////////////

func (b *Bytes) Int32LE(offset int) int32 {
	b.EnsureLen(offset + 4)
	return b.Pointer.Int32LE(offset)
}

func (b *Bytes) SetInt32LE(offset int, value int32) {
	b.EnsureLen(offset + 4)
	b.Pointer.SetInt32LE(offset, value)
}

func (b *Bytes) AppendInt32LE(value int32) {
	b.Pointer.SetInt32LE(b.ensureAppend(4), value)
}

///////////////////////////////////////////////////////////////////////////////////////////
// Int32 Little Endian
///////////////////////////////////////////////////////////////////////////////////////////

func (b *Bytes) Int32BE(offset int) int32 {
	b.EnsureLen(offset + 4)
	return b.Pointer.Int32BE(offset)
}

func (b *Bytes) SetInt32BE(offset int, value int32) {
	b.EnsureLen(offset + 4)
	b.Pointer.SetInt32BE(offset, value)
}

func (b *Bytes) AppendInt32BE(value int32) {
	b.Pointer.SetInt32BE(b.ensureAppend(4), value)
}

///////////////////////////////////////////////////////////////////////////////////////////
// UInt32 Native Endian
///////////////////////////////////////////////////////////////////////////////////////////

func (b *Bytes) UInt32(offset int) uint32 {
	b.EnsureLen(offset + 4)
	return b.Pointer.UInt32(offset)
}

func (b *Bytes) SetUInt32(offset int, value uint32) {
	b.EnsureLen(offset + 4)
	b.Pointer.SetUInt32(offset, value)
}

func (b *Bytes) AppendUInt32(value uint32) {
	b.Pointer.SetUInt32(b.ensureAppend(4), value)
}

///////////////////////////////////////////////////////////////////////////////////////////
// UInt32 Little Endian
///////////////////////////////////////////////////////////////////////////////////////////

func (b *Bytes) UInt32LE(offset int) uint32 {
	b.EnsureLen(offset + 4)
	return b.Pointer.UInt32LE(offset)
}

func (b *Bytes) SetUInt32LE(offset int, value uint32) {
	b.EnsureLen(offset + 4)
	b.Pointer.SetUInt32LE(offset, value)
}

func (b *Bytes) AppendUInt32LE(value uint32) {
	b.Pointer.SetUInt32LE(b.ensureAppend(4), value)
}

///////////////////////////////////////////////////////////////////////////////////////////
// UInt32 Big Endian
///////////////////////////////////////////////////////////////////////////////////////////

func (b *Bytes) UInt32BE(offset int) uint32 {
	b.EnsureLen(offset + 4)
	return b.Pointer.UInt32BE(offset)
}

func (b *Bytes) SetUInt32BE(offset int, value uint32) {
	b.EnsureLen(offset + 4)
	b.Pointer.SetUInt32BE(offset, value)
}

func (b *Bytes) AppendUInt32BE(value uint32) {
	b.Pointer.SetUInt32BE(b.ensureAppend(4), value)
}

///////////////////////////////////////////////////////////////////////////////////////////
// Int64 Native Endian
///////////////////////////////////////////////////////////////////////////////////////////

func (b *Bytes) Int64(offset int) int64 {
	b.EnsureLen(offset + 8)
	return b.Pointer.Int64(offset)
}

func (b *Bytes) SetInt64(offset int, value int64) {
	b.EnsureLen(offset + 8)
	b.Pointer.SetInt64(offset, value)
}

func (b *Bytes) AppendInt64(value int64) {
	b.Pointer.SetInt64(b.ensureAppend(8), value)
}

///////////////////////////////////////////////////////////////////////////////////////////
// Int64 Little Endian
///////////////////////////////////////////////////////////////////////////////////////////

func (b *Bytes) Int64LE(offset int) int64 {
	b.EnsureLen(offset + 8)
	return b.Pointer.Int64LE(offset)
}

func (b *Bytes) SetInt64LE(offset int, value int64) {
	b.EnsureLen(offset + 8)
	b.Pointer.SetInt64LE(offset, value)
}

func (b *Bytes) AppendInt64LE(value int64) {
	b.Pointer.SetInt64LE(b.ensureAppend(8), value)
}

///////////////////////////////////////////////////////////////////////////////////////////
// Int64 Big Endian
///////////////////////////////////////////////////////////////////////////////////////////

func (b *Bytes) Int64BE(offset int) int64 {
	b.EnsureLen(offset + 8)
	return b.Pointer.Int64BE(offset)
}

func (b *Bytes) SetInt64BE(offset int, value int64) {
	b.EnsureLen(offset + 8)
	b.Pointer.SetInt64BE(offset, value)
}

func (b *Bytes) AppendInt64BE(value int64) {
	b.Pointer.SetInt64BE(b.ensureAppend(8), value)
}

///////////////////////////////////////////////////////////////////////////////////////////
// UInt64 Native Endian
///////////////////////////////////////////////////////////////////////////////////////////

func (b *Bytes) UInt64(offset int) uint64 {
	b.EnsureLen(offset + 8)
	return b.Pointer.UInt64(offset)
}

func (b *Bytes) SetUInt64(offset int, value uint64) {
	b.EnsureLen(offset + 8)
	b.Pointer.SetUInt64(offset, value)
}

func (b *Bytes) AppendUInt64(value uint64) {
	b.Pointer.SetUInt64(b.ensureAppend(8), value)
}

///////////////////////////////////////////////////////////////////////////////////////////
// UInt64 Little Endian
///////////////////////////////////////////////////////////////////////////////////////////

func (b *Bytes) UInt64LE(offset int) uint64 {
	b.EnsureLen(offset + 8)
	return b.Pointer.UInt64LE(offset)
}

func (b *Bytes) SetUInt64LE(offset int, value uint64) {
	b.EnsureLen(offset + 8)
	b.Pointer.SetUInt64LE(offset, value)
}

func (b *Bytes) AppendUInt64LE(value uint64) {
	b.Pointer.SetUInt64LE(b.ensureAppend(8), value)
}

///////////////////////////////////////////////////////////////////////////////////////////
// UInt64 Big Endian
///////////////////////////////////////////////////////////////////////////////////////////

func (b *Bytes) UInt64BE(offset int) uint64 {
	b.EnsureLen(offset + 8)
	return b.Pointer.UInt64BE(offset)
}

func (b *Bytes) SetUInt64BE(offset int, value uint64) {
	b.EnsureLen(offset + 8)
	b.Pointer.SetUInt64BE(offset, value)
}

func (b *Bytes) AppendUInt64BE(value uint64) {
	b.Pointer.SetUInt64BE(b.ensureAppend(8), value)
}

///////////////////////////////////////////////////////////////////////////////////////////
// Float32 Native Endian
///////////////////////////////////////////////////////////////////////////////////////////

func (b *Bytes) Float32(offset int) float32 {
	b.EnsureLen(offset + 4)
	return b.Pointer.Float32(offset)
}

func (b *Bytes) SetFloat32(offset int, value float32) {
	b.EnsureLen(offset + 4)
	b.Pointer.SetFloat32(offset, value)
}

func (b *Bytes) AppendFloat32(value float32) {
	b.Pointer.SetFloat32(b.ensureAppend(4), value)
}

///////////////////////////////////////////////////////////////////////////////////////////
// Float32 Little Endian
///////////////////////////////////////////////////////////////////////////////////////////

func (b *Bytes) Float32LE(offset int) float32 {
	b.EnsureLen(offset + 4)
	return b.Pointer.Float32LE(offset)
}

func (b *Bytes) SetFloat32LE(offset int, value float32) {
	b.EnsureLen(offset + 4)
	b.Pointer.SetFloat32LE(offset, value)
}

func (b *Bytes) AppendFloat32LE(value float32) {
	b.Pointer.SetFloat32LE(b.ensureAppend(4), value)
}

///////////////////////////////////////////////////////////////////////////////////////////
// Float32 Big Endian
///////////////////////////////////////////////////////////////////////////////////////////

func (b *Bytes) Float32BE(offset int) float32 {
	b.EnsureLen(offset + 4)
	return b.Pointer.Float32BE(offset)
}

func (b *Bytes) SetFloat32BE(offset int, value float32) {
	b.EnsureLen(offset + 4)
	b.Pointer.SetFloat32BE(offset, value)
}

func (b *Bytes) AppendFloat32BE(value float32) {
	b.Pointer.SetFloat32BE(b.ensureAppend(4), value)
}

///////////////////////////////////////////////////////////////////////////////////////////
// Float64 Native Endian
///////////////////////////////////////////////////////////////////////////////////////////

func (b *Bytes) Float64(offset int) float64 {
	b.EnsureLen(offset + 8)
	return b.Pointer.Float64(offset)
}

func (b *Bytes) SetFloat64(offset int, value float64) {
	b.EnsureLen(offset + 8)
	b.Pointer.SetFloat64(offset, value)
}

func (b *Bytes) AppendFloat64(value float64) {
	b.Pointer.SetFloat64(b.ensureAppend(8), value)
}

///////////////////////////////////////////////////////////////////////////////////////////
// Float64 Little Endian
///////////////////////////////////////////////////////////////////////////////////////////

func (b *Bytes) Float64LE(offset int) float64 {
	b.EnsureLen(offset + 8)
	return b.Pointer.Float64LE(offset)
}

func (b *Bytes) SetFloat64LE(offset int, value float64) {
	b.EnsureLen(offset + 8)
	b.Pointer.SetFloat64(offset, value)
}

func (b *Bytes) AppendFloat64LE(value float64) {
	b.Pointer.SetFloat64LE(b.ensureAppend(8), value)
}

///////////////////////////////////////////////////////////////////////////////////////////
// Float64 Big Endian
///////////////////////////////////////////////////////////////////////////////////////////

func (b *Bytes) Float64BE(offset int) float64 {
	b.EnsureLen(offset + 8)
	return b.Pointer.Float64BE(offset)
}

func (b *Bytes) SetFloat64BE(offset int, value float64) {
	b.EnsureLen(offset + 8)
	b.Pointer.SetFloat64(offset, value)
}

func (b *Bytes) AppendFloat64BE(value float64) {
	b.Pointer.SetFloat64BE(b.ensureAppend(8), value)
}

///////////////////////////////////////////////////////////////////////////////////////////
// String
///////////////////////////////////////////////////////////////////////////////////////////

func (b *Bytes) SetString(offset int, value string) {
	b.EnsureCap(offset + len(value))
	length := offset + len(value)
	if int(b.len) < length {
		b.len = length
	}
	b.Pointer.SetString(offset, value)
}

///////////////////////////////////////////////////////////////////////////////////////////
// Bytes
///////////////////////////////////////////////////////////////////////////////////////////

func (b *Bytes) Set(offset int, value Bytes) {
	if value.IsNil() || value.len == 0 {
		return
	}
	length := offset + int(value.len)
	b.EnsureCap(offset + length)
	if int(b.len) < length {
		b.len = int(length)
	}
	memcpy(b.Unsafe(), value.Unsafe(), uintptr(value.len))
}

func (b *Bytes) SetBytes(offset int, value []byte) {
	b.EnsureCap(offset + len(value))
	length := offset + len(value)
	if int(b.len) < length {
		b.len = int(length)
	}
	b.Pointer.SetBytes(offset, value)
}

///////////////////////////////////////////////////////////////////////////////////////////
// Int24 Native Endian
///////////////////////////////////////////////////////////////////////////////////////////

func (b *Bytes) Int24(offset int) int32 {
	b.EnsureLen(offset + 3)
	return b.Pointer.Int24(offset)
}

func (b *Bytes) SetInt24(offset int, value int32) {
	b.EnsureLen(offset + 3)
	b.Pointer.SetInt24(offset, value)
}

func (b *Bytes) AppendInt24(value int32) {
	b.Pointer.SetInt24(b.ensureAppend(3), value)
}

///////////////////////////////////////////////////////////////////////////////////////////
// Int24 Little Endian
///////////////////////////////////////////////////////////////////////////////////////////

func (b *Bytes) Int24LE(offset int) int32 {
	b.EnsureLen(offset + 3)
	return b.Pointer.Int24LE(offset)
}

func (b *Bytes) SetInt24LE(offset int, value int32) {
	b.EnsureLen(offset + 3)
	b.Pointer.SetInt24LE(offset, value)
}

func (b *Bytes) AppendInt24LE(value int32) {
	b.Pointer.SetInt24LE(b.ensureAppend(3), value)
}

///////////////////////////////////////////////////////////////////////////////////////////
// Int24 Big Endian
///////////////////////////////////////////////////////////////////////////////////////////

func (b *Bytes) Int24BE(offset int) int32 {
	b.EnsureLen(offset + 3)
	return b.Pointer.Int24BE(offset)
}

func (b *Bytes) SetInt24BE(offset int, value int32) {
	b.EnsureLen(offset + 3)
	b.Pointer.SetInt24BE(offset, value)
}

func (b *Bytes) AppendInt24BE(value int32) {
	b.Pointer.SetInt24BE(b.ensureAppend(3), value)
}

///////////////////////////////////////////////////////////////////////////////////////////
// UInt24 Native Endian
///////////////////////////////////////////////////////////////////////////////////////////

func (b *Bytes) UInt24(offset int) uint32 {
	b.EnsureLen(offset + 3)
	return b.Pointer.UInt24(offset)
}

func (b *Bytes) SetUInt24(offset int, value uint32) {
	b.EnsureLen(offset + 3)
	b.Pointer.SetUInt24(offset, value)
}

func (b *Bytes) AppendUInt24(value uint32) {
	b.Pointer.SetUInt24(b.ensureAppend(3), value)
}

///////////////////////////////////////////////////////////////////////////////////////////
// UInt24 Little Endian
///////////////////////////////////////////////////////////////////////////////////////////

func (b *Bytes) UInt24LE(offset int) uint32 {
	b.EnsureLen(offset + 3)
	return b.Pointer.UInt24LE(offset)
}

func (b *Bytes) SetUInt24LE(offset int, value uint32) {
	b.EnsureLen(offset + 3)
	b.Pointer.SetUInt24LE(offset, value)
}

func (b *Bytes) AppendUInt24LE(value uint32) {
	b.Pointer.SetUInt24LE(b.ensureAppend(3), value)
}

///////////////////////////////////////////////////////////////////////////////////////////
// UInt24 Big Endian
///////////////////////////////////////////////////////////////////////////////////////////

func (b *Bytes) UInt24BE(offset int) uint32 {
	b.EnsureLen(offset + 3)
	return b.Pointer.UInt24BE(offset)
}

func (b *Bytes) SetUInt24BE(offset int, value uint32) {
	b.EnsureLen(offset + 3)
	b.Pointer.SetUInt24BE(offset, value)
}

func (b *Bytes) AppendUInt24BE(value uint32) {
	b.Pointer.SetUInt24BE(b.ensureAppend(3), value)
}

///////////////////////////////////////////////////////////////////////////////////////////
// Int40 Native Endian
///////////////////////////////////////////////////////////////////////////////////////////

func (b *Bytes) Int40(offset int) int64 {
	b.EnsureLen(offset + 5)
	return b.Pointer.Int40(offset)
}

func (b *Bytes) SetInt40(offset int, value int64) {
	b.EnsureLen(offset + 5)
	b.Pointer.SetInt40(offset, value)
}

func (b *Bytes) AppendInt40(value int64) {
	b.Pointer.SetInt40(b.ensureAppend(5), value)
}

///////////////////////////////////////////////////////////////////////////////////////////
// Int40 Little Endian
///////////////////////////////////////////////////////////////////////////////////////////

func (b *Bytes) Int40LE(offset int) int64 {
	b.EnsureLen(offset + 5)
	return b.Pointer.Int40LE(offset)
}

func (b *Bytes) SetInt40LE(offset int, value int64) {
	b.EnsureLen(offset + 5)
	b.Pointer.SetInt40LE(offset, value)
}

func (b *Bytes) AppendInt40LE(value int64) {
	b.Pointer.SetInt40LE(b.ensureAppend(5), value)
}

///////////////////////////////////////////////////////////////////////////////////////////
// Int40 Big Endian
///////////////////////////////////////////////////////////////////////////////////////////

func (b *Bytes) Int40BE(offset int) int64 {
	b.EnsureLen(offset + 5)
	return b.Pointer.Int40BE(offset)
}

func (b *Bytes) SetInt40BE(offset int, value int64) {
	b.EnsureLen(offset + 5)
	b.Pointer.SetInt40BE(offset, value)
}

func (b *Bytes) AppendInt40BE(value int64) {
	b.Pointer.SetInt40BE(b.ensureAppend(5), value)
}

///////////////////////////////////////////////////////////////////////////////////////////
// UInt40 Native Endian
///////////////////////////////////////////////////////////////////////////////////////////

func (b *Bytes) UInt40(offset int) uint64 {
	b.EnsureLen(offset + 5)
	return b.Pointer.UInt40(offset)
}

func (b *Bytes) SetUInt40(offset int, value uint64) {
	b.EnsureLen(offset + 5)
	b.Pointer.SetUInt40(offset, value)
}

func (b *Bytes) AppendUInt40(value uint64) {
	b.Pointer.SetUInt40(b.ensureAppend(5), value)
}

///////////////////////////////////////////////////////////////////////////////////////////
// UInt40 Little Endian
///////////////////////////////////////////////////////////////////////////////////////////

func (b *Bytes) UInt40LE(offset int) uint64 {
	b.EnsureLen(offset + 5)
	return b.Pointer.UInt40LE(offset)
}

func (b *Bytes) SetUInt40LE(offset int, value uint64) {
	b.EnsureLen(offset + 5)
	b.Pointer.SetUInt40LE(offset, value)
}

func (b *Bytes) AppendUInt40LE(value uint64) {
	b.Pointer.SetUInt40LE(b.ensureAppend(5), value)
}

///////////////////////////////////////////////////////////////////////////////////////////
// UInt40 Big Endian
///////////////////////////////////////////////////////////////////////////////////////////

func (b *Bytes) UInt40BE(offset int) uint64 {
	b.EnsureLen(offset + 5)
	return b.Pointer.UInt40BE(offset)
}

func (b *Bytes) SetUInt40BE(offset int, value uint64) {
	b.EnsureLen(offset + 5)
	b.Pointer.SetUInt40BE(offset, value)
}

func (b *Bytes) AppendUInt40BE(value uint64) {
	b.Pointer.SetUInt40BE(b.ensureAppend(5), value)
}

///////////////////////////////////////////////////////////////////////////////////////////
// Int48 Native Endian
///////////////////////////////////////////////////////////////////////////////////////////

func (b *Bytes) Int48(offset int) int64 {
	b.EnsureLen(offset + 6)
	return b.Pointer.Int48(offset)
}

func (b *Bytes) SetInt48(offset int, value int64) {
	b.EnsureLen(offset + 6)
	b.Pointer.SetInt48(offset, value)
}

func (b *Bytes) AppendInt48(value int64) {
	b.Pointer.SetInt48(b.ensureAppend(6), value)
}

///////////////////////////////////////////////////////////////////////////////////////////
// Int48 Little Endian
///////////////////////////////////////////////////////////////////////////////////////////

func (b *Bytes) Int48LE(offset int) int64 {
	b.EnsureLen(offset + 6)
	return b.Pointer.Int48LE(offset)
}

func (b *Bytes) SetInt48LE(offset int, value int64) {
	b.EnsureLen(offset + 6)
	b.Pointer.SetInt48LE(offset, value)
}

func (b *Bytes) AppendInt48LE(value int64) {
	b.Pointer.SetInt48LE(b.ensureAppend(6), value)
}

///////////////////////////////////////////////////////////////////////////////////////////
// Int48 Big Endian
///////////////////////////////////////////////////////////////////////////////////////////

func (b *Bytes) Int48BE(offset int) int64 {
	b.EnsureLen(offset + 6)
	return b.Pointer.Int48BE(offset)
}

func (b *Bytes) SetInt48BE(offset int, value int64) {
	b.EnsureLen(offset + 6)
	b.Pointer.SetInt48BE(offset, value)
}

func (b *Bytes) AppendInt48BE(value int64) {
	b.Pointer.SetInt48BE(b.ensureAppend(6), value)
}

///////////////////////////////////////////////////////////////////////////////////////////
// UInt48 Native Endian
///////////////////////////////////////////////////////////////////////////////////////////

func (b *Bytes) UInt48(offset int) uint64 {
	b.EnsureLen(offset + 6)
	return b.Pointer.UInt48(offset)
}

func (b *Bytes) SetUInt48(offset int, value uint64) {
	b.EnsureLen(offset + 6)
	b.Pointer.SetUInt48(offset, value)
}

func (b *Bytes) AppendUInt48(value uint64) {
	b.Pointer.SetUInt48(b.ensureAppend(6), value)
}

///////////////////////////////////////////////////////////////////////////////////////////
// UInt48 Little Endian
///////////////////////////////////////////////////////////////////////////////////////////

func (b *Bytes) UInt48LE(offset int) uint64 {
	b.EnsureLen(offset + 6)
	return b.Pointer.UInt48LE(offset)
}

func (b *Bytes) SetUInt48LE(offset int, value uint64) {
	b.EnsureLen(offset + 6)
	b.Pointer.SetUInt48LE(offset, value)
}

func (b *Bytes) AppendUInt48LE(value uint64) {
	b.Pointer.SetUInt48LE(b.ensureAppend(6), value)
}

///////////////////////////////////////////////////////////////////////////////////////////
// UInt48 Big Endian
///////////////////////////////////////////////////////////////////////////////////////////

func (b *Bytes) UInt48BE(offset int) uint64 {
	b.EnsureLen(offset + 6)
	return b.Pointer.UInt48BE(offset)
}

func (b *Bytes) SetUInt48BE(offset int, value uint64) {
	b.EnsureLen(offset + 6)
	b.Pointer.SetUInt48BE(offset, value)
}

func (b *Bytes) AppendUInt48BE(value uint64) {
	b.Pointer.SetUInt48BE(b.ensureAppend(6), value)
}

///////////////////////////////////////////////////////////////////////////////////////////
// Int56 Native Endian
///////////////////////////////////////////////////////////////////////////////////////////

func (b *Bytes) Int56(offset int) int64 {
	b.EnsureLen(offset + 7)
	return b.Pointer.Int56(offset)
}

func (b *Bytes) SetInt56(offset int, value int64) {
	b.EnsureLen(offset + 7)
	b.Pointer.SetInt56(offset, value)
}

func (b *Bytes) AppendInt56(value int64) {
	b.Pointer.SetInt56(b.ensureAppend(7), value)
}

///////////////////////////////////////////////////////////////////////////////////////////
// Int56 Little Endian
///////////////////////////////////////////////////////////////////////////////////////////

func (b *Bytes) Int56LE(offset int) int64 {
	b.EnsureLen(offset + 7)
	return b.Pointer.Int56LE(offset)
}

func (b *Bytes) SetInt56LE(offset int, value int64) {
	b.EnsureLen(offset + 7)
	b.Pointer.SetInt56LE(offset, value)
}

func (b *Bytes) AppendInt56LE(value int64) {
	b.Pointer.SetInt56LE(b.ensureAppend(7), value)
}

///////////////////////////////////////////////////////////////////////////////////////////
// Int56 Big Endian
///////////////////////////////////////////////////////////////////////////////////////////

func (b *Bytes) Int56BE(offset int) int64 {
	b.EnsureLen(offset + 7)
	return b.Pointer.Int56BE(offset)
}

func (b *Bytes) SetInt56BE(offset int, value int64) {
	b.EnsureLen(offset + 7)
	b.Pointer.SetInt56BE(offset, value)
}

func (b *Bytes) AppendInt56BE(value int64) {
	b.Pointer.SetInt56BE(b.ensureAppend(7), value)
}

///////////////////////////////////////////////////////////////////////////////////////////
// UInt56 Native Endian
///////////////////////////////////////////////////////////////////////////////////////////

func (b *Bytes) UInt56(offset int) uint64 {
	b.EnsureLen(offset + 7)
	return b.Pointer.UInt56(offset)
}

func (b *Bytes) SetUInt56(offset int, value uint64) {
	b.EnsureLen(offset + 7)
	b.Pointer.SetUInt56(offset, value)
}

func (b *Bytes) AppendUInt56(value uint64) {
	b.Pointer.SetUInt56(b.ensureAppend(7), value)
}

///////////////////////////////////////////////////////////////////////////////////////////
// UInt56 Little Endian
///////////////////////////////////////////////////////////////////////////////////////////

func (b *Bytes) UInt56LE(offset int) uint64 {
	b.EnsureLen(offset + 7)
	return b.Pointer.UInt56LE(offset)
}

func (b *Bytes) SetUInt56LE(offset int, value uint64) {
	b.EnsureLen(offset + 7)
	b.Pointer.SetUInt56LE(offset, value)
}

func (b *Bytes) AppendUInt56LE(value uint64) {
	b.Pointer.SetUInt56LE(b.ensureAppend(7), value)
}

///////////////////////////////////////////////////////////////////////////////////////////
// UInt56 Big Endian
///////////////////////////////////////////////////////////////////////////////////////////

func (b *Bytes) UInt56BE(offset int) uint64 {
	b.EnsureLen(offset + 7)
	return b.Pointer.UInt56BE(offset)
}

func (b *Bytes) SetUInt56BE(offset int, value uint64) {
	b.EnsureLen(offset + 7)
	b.Pointer.SetUInt56BE(offset, value)
}

func (b *Bytes) AppendUInt56BE(value uint64) {
	b.Pointer.SetUInt56BE(b.ensureAppend(7), value)
}

func (b *Bytes) ensureAppend(extra int) int {
	offset := b.len
	b.EnsureCap(offset + extra)
	b.len += extra
	return offset
}

// EnsureLen ensures the length is at least neededLen in size
// If not, EnsureCap(neededLen) is called and the length set to neededLen.
func (b *Bytes) EnsureLen(neededLen int) {
	if b.len > neededLen {
		return
	}
	b.EnsureCap(neededLen)
	b.len = neededLen
}

// EnsureCap ensures the capacity is at least neededCap in size
//goland:noinspection GoVetUnsafePointer
func (b *Bytes) EnsureCap(neededCap int) bool {
	if b.cap >= neededCap {
		return true
	}
	newCap := neededCap - b.cap
	addr := b.alloc.Realloc(b.Pointer, uintptr(newCap))
	//addr := ((*Allocator)(unsafe.Pointer(p.alloc))).Realloc(p.Pointer, Pointer(newCap))
	if addr == 0 {
		return false
	}
	b.Pointer = addr
	b.cap = newCap
	return true
}

//// Clone creates a copy of this instance of Bytes
//func (b *Bytes) Clone() Bytes {
//	c := b.alloc.Bytes(Pointer(b.len))
//	memcpy(c.Unsafe(), b.Unsafe(), uintptr(b.len))
//	return c
//}

func (b *Bytes) Append(value Bytes) int {
	b.EnsureCap(b.len + value.len)
	i := b.len
	b.len += value.len
	return int(i)
}

//goland:noinspection GoVetUnsafePointer
func (b *Bytes) AppendBytes(value []byte) int {
	b.EnsureCap(b.len + len(value))
	i := b.len
	memcpy(unsafe.Pointer(b.Pointer), unsafe.Pointer(&value[0]), uintptr(len(value)))
	b.len += len(value)
	return i
}

func (b *Bytes) AppendString(value string) int {
	b.EnsureCap(b.len + len(value))
	i := b.len
	b.len += len(value)
	return i
}

func (b *Bytes) SetLength(length int) {
	b.EnsureCap(length)
	b.len = length
}
