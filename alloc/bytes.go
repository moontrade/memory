package alloc

import (
	"github.com/moontrade/memory/hash"
	. "github.com/moontrade/memory/mem"
	"unsafe"
)

//type StrSlice struct {
//	Pointer
//	str Bytes
//}

// Bytes is a compact single dynamic allocation to be used as an unsafe replacement for string.
type Bytes struct {
	Pointer // Use for unchecked unsafe access
}

func (s *Bytes) Free() {
	if s == nil || s.Pointer == 0 {
		return
	}
	s.Allocator().Free(s.allocationPointer())
	s.Pointer = 0
}

func (s *Bytes) CString() unsafe.Pointer {
	l := s.Len()
	if l == 0 {
		return nil
	}
	s.EnsureCap(l + 1)
	// Ensure it's NULL terminated
	*(*byte)(unsafe.Pointer(uintptr(s.Pointer) + uintptr(l))) = 0
	return s.Unsafe()
}

func (s Bytes) flags() uint8 {
	return *(*uint8)(unsafe.Pointer(s.Pointer - 1))
}

// Reset zeroes out the entire allocation and sets the length back to 0
func (s *Bytes) Reset() {
	Zero(s.Unsafe(), uintptr(s.Cap()))
	s.setLen(0)
}

// Zero zeroes out the entire allocation.
func (s Bytes) Zero() {
	Zero(s.Unsafe(), uintptr(s.Cap()))
}

//goland:noinspection GoVetUnsafePointer
func (s *Bytes) Equals(o Bytes) bool {
	l := s.Len()
	return l == o.Len() && (*s == o || Equals(
		s.Unsafe(),
		o.Unsafe(), uintptr(l)))
}

func (s *Bytes) IsNil() bool {
	return s == nil || s.Pointer == 0
}

func (s Bytes) IsEmpty() bool {
	return s.Pointer == 0 || s.Len() == 0
}

func (s *Bytes) CheckBounds(offset int) bool {
	return uintptr(s.Pointer) == 0 || s.Len() < offset
}

func (s *Bytes) String() string {
	l := s.Len()
	if l == 0 {
		return ""
	}
	return s.Pointer.String(0, l)
}

func (s *Bytes) Bytes() []byte {
	l := s.Len()
	if l == 0 {
		return nil
	}
	return s.Pointer.Bytes(0, l, l)
}

func (s *Bytes) ensureAppend(extra int) int {
	offset := s.Len()
	s.EnsureCap(offset + extra)
	s.setLen(offset + extra)
	return offset
}

// EnsureLen ensures the length is at least neededLen in size
// If not, EnsureCap(neededLen) is called and the length set to neededLen.
func (s *Bytes) EnsureLen(neededLen int) {
	l := s.Len()
	if l > neededLen {
		return
	}
	s.EnsureCap(neededLen)
	s.setLen(neededLen)
}

// EnsureCap ensures the capacity is at least neededCap in size
//goland:noinspection GoVetUnsafePointer
func (s *Bytes) EnsureCap(neededCap int) bool {
	cp := s.Cap()
	if cp >= neededCap {
		return true
	}
	newCap := neededCap - cp
	addr := s.Allocator().Realloc(s.allocationPointer(), uintptr(newCap))
	//addr := ((*Allocator)(unsafe.Pointer(p.alloc))).Realloc(p.Pointer, Pointer(newCap))
	if addr == 0 {
		return false
	}
	*s = Bytes{addr}
	return true
}

// Clone creates a copy of this instance of Bytes
func (s *Bytes) Clone() Bytes {
	l := s.Len()
	c := NewString(uintptr(l))
	Copy(c.Unsafe(), s.Pointer.Unsafe(), uintptr(l))
	c.setLen(l)
	return c
}

//goland:noinspection GoVetUnsafePointer
func (s *Bytes) Append(value Bytes) {
	if value.Pointer == 0 {
		return
	}
	l := s.Len()
	vl := value.Len()
	s.EnsureCap(l + vl)
	Move(unsafe.Pointer(uintptr(int(s.Pointer)+l)), value.Unsafe(), uintptr(vl))
	s.setLen(l + vl)
}

//goland:noinspection GoVetUnsafePointer
func (s *Bytes) AppendBytes(value []byte) {
	if len(value) == 0 {
		return
	}
	l := s.Len()
	s.EnsureCap(l + len(value))
	Move(unsafe.Pointer(uintptr(int(s.Pointer)+l)), unsafe.Pointer(&value[0]), uintptr(len(value)))
	s.setLen(l + len(value))
}

//goland:noinspection GoVetUnsafePointer
func (s *Bytes) AppendString(value string) {
	if len(value) == 0 {
		return
	}
	l := s.Len()
	s.EnsureCap(l + len(value))
	Move(
		unsafe.Pointer(uintptr(int(s.Pointer)+l)),
		unsafe.Pointer(*(*uintptr)(unsafe.Pointer(&value))),
		uintptr(len(value)),
	)
	s.setLen(l + len(value))
}

func (s *Bytes) SetLength(length int) {
	s.EnsureCap(length)
	s.setLen(length)
}

///////////////////////////////////////////////////////////////////////////////////////////
// Int
///////////////////////////////////////////////////////////////////////////////////////////

func (s *Bytes) Int(offset int) int {
	s.EnsureLen(offset + int(unsafe.Sizeof(int(0))))
	return s.Pointer.Int(offset)
}

func (s *Bytes) SetInt(offset int, value int) {
	s.EnsureLen(offset + int(unsafe.Sizeof(int(0))))
	s.Pointer.SetInt(offset, value)
}

//goland:noinspection GoVetUnsafePointer
func (s *Bytes) AppendInt(value int) {
	s.Pointer.SetInt(s.ensureAppend(int(unsafe.Sizeof(int(0)))), value)
}

///////////////////////////////////////////////////////////////////////////////////////////
// UInt
///////////////////////////////////////////////////////////////////////////////////////////

func (s *Bytes) UInt(offset int) int {
	s.EnsureLen(offset + int(unsafe.Sizeof(uint(0))))
	return s.Pointer.Int(offset)
}

func (s *Bytes) SetUInt(offset int, value uint) {
	s.EnsureLen(offset + int(unsafe.Sizeof(uint(0))))
	s.Pointer.SetUInt(offset, value)
}

//goland:noinspection GoVetUnsafePointer
func (s *Bytes) AppendUInt(value uint) {
	s.Pointer.SetUInt(s.ensureAppend(int(unsafe.Sizeof(uint(0)))), value)
}

///////////////////////////////////////////////////////////////////////////////////////////
// Pointer
///////////////////////////////////////////////////////////////////////////////////////////

func (s *Bytes) PointerAt(offset int) Pointer {
	s.EnsureLen(offset + 8)
	return s.Pointer.Pointer(offset)
}

func (s *Bytes) SetPointer(offset int, value Pointer) {
	s.EnsureLen(offset + int(unsafe.Sizeof(uintptr(0))))
	s.Pointer.SetPointer(offset, value)
}

func (s *Bytes) AppendPointer(value Pointer) {
	s.Pointer.SetPointer(s.ensureAppend(int(unsafe.Sizeof(Pointer(0)))), value)
}

///////////////////////////////////////////////////////////////////////////////////////////
// uintptr
///////////////////////////////////////////////////////////////////////////////////////////

func (s *Bytes) Uintptr(offset int) uintptr {
	s.EnsureLen(offset + int(unsafe.Sizeof(uintptr(0))))
	return s.Pointer.Uintptr(offset)
}

func (s *Bytes) SetUintptr(offset int, value uintptr) {
	s.EnsureLen(offset + int(unsafe.Sizeof(uintptr(0))))
	s.Pointer.SetUintptr(offset, value)
}

func (s *Bytes) AppendUintptr(value uintptr) {
	s.Pointer.SetUintptr(s.ensureAppend(int(unsafe.Sizeof(uintptr(0)))), value)
}

///////////////////////////////////////////////////////////////////////////////////////////
// Int8
///////////////////////////////////////////////////////////////////////////////////////////

func (s *Bytes) Int8(offset int) int8 {
	s.EnsureLen(offset + 1)
	return s.Pointer.Int8(offset)
}

func (s *Bytes) SetInt8(offset int, value int8) {
	s.EnsureLen(offset + 1)
	s.Pointer.SetInt8(offset, value)
}

func (s *Bytes) AppendInt8(value int8) {
	s.Pointer.SetInt8(s.ensureAppend(1), value)
}

///////////////////////////////////////////////////////////////////////////////////////////
// UInt8
///////////////////////////////////////////////////////////////////////////////////////////

func (s *Bytes) UInt8(offset int) uint8 {
	s.EnsureLen(offset + 1)
	return s.Pointer.UInt8(offset)
}

// SetUInt8 is safe version. Will grow allocation if needed.
func (s *Bytes) SetUInt8(offset int, value uint8) {
	s.EnsureLen(offset + 1)
	s.Pointer.SetUInt8(offset, value)
}

func (s *Bytes) AppendUInt8(value uint8) {
	s.Pointer.SetUInt8(s.ensureAppend(1), value)
}

///////////////////////////////////////////////////////////////////////////////////////////
// Byte
///////////////////////////////////////////////////////////////////////////////////////////

func (s *Bytes) Byte(offset int) byte {
	s.EnsureLen(offset + 1)
	return s.Pointer.Byte(offset)
}

func (s *Bytes) SetByte(offset int, value byte) {
	s.EnsureLen(offset + 1)
	s.Pointer.SetByte(offset, value)
}

func (s *Bytes) AppendByte(value byte) {
	s.Pointer.SetByte(s.ensureAppend(1), value)
}

///////////////////////////////////////////////////////////////////////////////////////////
// Int16 Native Endian
///////////////////////////////////////////////////////////////////////////////////////////

func (s *Bytes) Int16(offset int) int16 {
	s.EnsureLen(offset + 2)
	return s.Pointer.Int16(offset)
}

func (s *Bytes) SetInt16(offset int, value int16) {
	s.EnsureLen(offset + 2)
	s.Pointer.SetInt16(offset, value)
}

func (s *Bytes) AppendInt16(value int16) {
	s.Pointer.SetInt16(s.ensureAppend(2), value)
}

///////////////////////////////////////////////////////////////////////////////////////////
// Int16 Little Endian
///////////////////////////////////////////////////////////////////////////////////////////

func (s *Bytes) Int16LE(offset int) int16 {
	s.EnsureLen(offset + 2)
	return s.Pointer.Int16LE(offset)
}

func (s *Bytes) SetInt16LE(offset int, value int16) {
	s.EnsureLen(offset + 2)
	s.Pointer.SetInt16LE(offset, value)
}

func (s *Bytes) AppendInt16LE(value int16) {
	s.Pointer.SetInt16LE(s.ensureAppend(2), value)
}

///////////////////////////////////////////////////////////////////////////////////////////
// Int16 Little Endian
///////////////////////////////////////////////////////////////////////////////////////////

func (s *Bytes) Int16BE(offset int) int16 {
	s.EnsureLen(offset + 2)
	return s.Pointer.Int16BE(offset)
}

func (s *Bytes) SetInt16BE(offset int, value int16) {
	s.EnsureLen(offset + 2)
	s.Pointer.SetInt16BE(offset, value)
}

func (s *Bytes) AppendInt16BE(value int16) {
	s.Pointer.SetInt16BE(s.ensureAppend(2), value)
}

///////////////////////////////////////////////////////////////////////////////////////////
// UInt16 Native Endian
///////////////////////////////////////////////////////////////////////////////////////////

func (s *Bytes) UInt16(offset int) uint16 {
	s.EnsureLen(offset + 2)
	return s.Pointer.UInt16(offset)
}

func (s *Bytes) SetUInt16(offset int, value uint16) {
	s.EnsureLen(offset + 2)
	s.Pointer.SetUInt16(offset, value)
}

func (s *Bytes) AppendUInt16(value uint16) {
	s.Pointer.SetUInt16(s.ensureAppend(2), value)
}

///////////////////////////////////////////////////////////////////////////////////////////
// UInt16 Little Endian
///////////////////////////////////////////////////////////////////////////////////////////

func (s *Bytes) UInt16LE(offset int) uint16 {
	s.EnsureLen(offset + 2)
	return s.Pointer.UInt16LE(offset)
}

func (s *Bytes) SetUInt16LE(offset int, value uint16) {
	s.EnsureLen(offset + 2)
	s.Pointer.SetUInt16LE(offset, value)
}

func (s *Bytes) AppendUInt16LE(value uint16) {
	s.Pointer.SetUInt16LE(s.ensureAppend(2), value)
}

///////////////////////////////////////////////////////////////////////////////////////////
// UInt16 Big Endian
///////////////////////////////////////////////////////////////////////////////////////////

func (s *Bytes) UInt16BE(offset int) uint16 {
	s.EnsureLen(offset + 2)
	return s.Pointer.UInt16BE(offset)
}

func (s *Bytes) SetUInt16BE(offset int, value uint16) {
	s.EnsureLen(offset + 2)
	s.Pointer.SetUInt16BE(offset, value)
}

func (s *Bytes) AppendUInt16BE(value uint16) {
	s.Pointer.SetUInt16BE(s.ensureAppend(2), value)
}

///////////////////////////////////////////////////////////////////////////////////////////
// Int32 Native Endian
///////////////////////////////////////////////////////////////////////////////////////////

func (s *Bytes) Int32(offset int) int32 {
	s.EnsureLen(offset + 4)
	return s.Pointer.Int32(offset)
}

func (s *Bytes) SetInt32(offset int, value int32) {
	s.EnsureLen(offset + 4)
	s.Pointer.SetInt32(offset, value)
}

func (s *Bytes) AppendInt32(value int32) {
	s.Pointer.SetInt32(s.ensureAppend(4), value)
}

///////////////////////////////////////////////////////////////////////////////////////////
// Int32 Little Endian
///////////////////////////////////////////////////////////////////////////////////////////

func (s *Bytes) Int32LE(offset int) int32 {
	s.EnsureLen(offset + 4)
	return s.Pointer.Int32LE(offset)
}

func (s *Bytes) SetInt32LE(offset int, value int32) {
	s.EnsureLen(offset + 4)
	s.Pointer.SetInt32LE(offset, value)
}

func (s *Bytes) AppendInt32LE(value int32) {
	s.Pointer.SetInt32LE(s.ensureAppend(4), value)
}

///////////////////////////////////////////////////////////////////////////////////////////
// Int32 Little Endian
///////////////////////////////////////////////////////////////////////////////////////////

func (s *Bytes) Int32BE(offset int) int32 {
	s.EnsureLen(offset + 4)
	return s.Pointer.Int32BE(offset)
}

func (s *Bytes) SetInt32BE(offset int, value int32) {
	s.EnsureLen(offset + 4)
	s.Pointer.SetInt32BE(offset, value)
}

func (s *Bytes) AppendInt32BE(value int32) {
	s.Pointer.SetInt32BE(s.ensureAppend(4), value)
}

///////////////////////////////////////////////////////////////////////////////////////////
// UInt32 Native Endian
///////////////////////////////////////////////////////////////////////////////////////////

func (s *Bytes) UInt32(offset int) uint32 {
	s.EnsureLen(offset + 4)
	return s.Pointer.UInt32(offset)
}

func (s *Bytes) SetUInt32(offset int, value uint32) {
	s.EnsureLen(offset + 4)
	s.Pointer.SetUInt32(offset, value)
}

func (s *Bytes) AppendUInt32(value uint32) {
	s.Pointer.SetUInt32(s.ensureAppend(4), value)
}

///////////////////////////////////////////////////////////////////////////////////////////
// UInt32 Little Endian
///////////////////////////////////////////////////////////////////////////////////////////

func (s *Bytes) UInt32LE(offset int) uint32 {
	s.EnsureLen(offset + 4)
	return s.Pointer.UInt32LE(offset)
}

func (s *Bytes) SetUInt32LE(offset int, value uint32) {
	s.EnsureLen(offset + 4)
	s.Pointer.SetUInt32LE(offset, value)
}

func (s *Bytes) AppendUInt32LE(value uint32) {
	s.Pointer.SetUInt32LE(s.ensureAppend(4), value)
}

///////////////////////////////////////////////////////////////////////////////////////////
// UInt32 Big Endian
///////////////////////////////////////////////////////////////////////////////////////////

func (s *Bytes) UInt32BE(offset int) uint32 {
	s.EnsureLen(offset + 4)
	return s.Pointer.UInt32BE(offset)
}

func (s *Bytes) SetUInt32BE(offset int, value uint32) {
	s.EnsureLen(offset + 4)
	s.Pointer.SetUInt32BE(offset, value)
}

func (s *Bytes) AppendUInt32BE(value uint32) {
	s.Pointer.SetUInt32BE(s.ensureAppend(4), value)
}

///////////////////////////////////////////////////////////////////////////////////////////
// Int64 Native Endian
///////////////////////////////////////////////////////////////////////////////////////////

func (s *Bytes) Int64(offset int) int64 {
	s.EnsureLen(offset + 8)
	return s.Pointer.Int64(offset)
}

func (s *Bytes) SetInt64(offset int, value int64) {
	s.EnsureLen(offset + 8)
	s.Pointer.SetInt64(offset, value)
}

func (s *Bytes) AppendInt64(value int64) {
	s.Pointer.SetInt64(s.ensureAppend(8), value)
}

///////////////////////////////////////////////////////////////////////////////////////////
// Int64 Little Endian
///////////////////////////////////////////////////////////////////////////////////////////

func (s *Bytes) Int64LE(offset int) int64 {
	s.EnsureLen(offset + 8)
	return s.Pointer.Int64LE(offset)
}

func (s *Bytes) SetInt64LE(offset int, value int64) {
	s.EnsureLen(offset + 8)
	s.Pointer.SetInt64LE(offset, value)
}

func (s *Bytes) AppendInt64LE(value int64) {
	s.Pointer.SetInt64LE(s.ensureAppend(8), value)
}

///////////////////////////////////////////////////////////////////////////////////////////
// Int64 Big Endian
///////////////////////////////////////////////////////////////////////////////////////////

func (s *Bytes) Int64BE(offset int) int64 {
	s.EnsureLen(offset + 8)
	return s.Pointer.Int64BE(offset)
}

func (s *Bytes) SetInt64BE(offset int, value int64) {
	s.EnsureLen(offset + 8)
	s.Pointer.SetInt64BE(offset, value)
}

func (s *Bytes) AppendInt64BE(value int64) {
	s.Pointer.SetInt64BE(s.ensureAppend(8), value)
}

///////////////////////////////////////////////////////////////////////////////////////////
// UInt64 Native Endian
///////////////////////////////////////////////////////////////////////////////////////////

func (s *Bytes) UInt64(offset int) uint64 {
	s.EnsureLen(offset + 8)
	return s.Pointer.UInt64(offset)
}

func (s *Bytes) SetUInt64(offset int, value uint64) {
	s.EnsureLen(offset + 8)
	s.Pointer.SetUInt64(offset, value)
}

func (s *Bytes) AppendUInt64(value uint64) {
	s.Pointer.SetUInt64(s.ensureAppend(8), value)
}

///////////////////////////////////////////////////////////////////////////////////////////
// UInt64 Little Endian
///////////////////////////////////////////////////////////////////////////////////////////

func (s *Bytes) UInt64LE(offset int) uint64 {
	s.EnsureLen(offset + 8)
	return s.Pointer.UInt64LE(offset)
}

func (s *Bytes) SetUInt64LE(offset int, value uint64) {
	s.EnsureLen(offset + 8)
	s.Pointer.SetUInt64LE(offset, value)
}

func (s *Bytes) AppendUInt64LE(value uint64) {
	s.Pointer.SetUInt64LE(s.ensureAppend(8), value)
}

///////////////////////////////////////////////////////////////////////////////////////////
// UInt64 Big Endian
///////////////////////////////////////////////////////////////////////////////////////////

func (s *Bytes) UInt64BE(offset int) uint64 {
	s.EnsureLen(offset + 8)
	return s.Pointer.UInt64BE(offset)
}

func (s *Bytes) SetUInt64BE(offset int, value uint64) {
	s.EnsureLen(offset + 8)
	s.Pointer.SetUInt64BE(offset, value)
}

func (s *Bytes) AppendUInt64BE(value uint64) {
	s.Pointer.SetUInt64BE(s.ensureAppend(8), value)
}

///////////////////////////////////////////////////////////////////////////////////////////
// Float32 Native Endian
///////////////////////////////////////////////////////////////////////////////////////////

func (s *Bytes) Float32(offset int) float32 {
	s.EnsureLen(offset + 4)
	return s.Pointer.Float32(offset)
}

func (s *Bytes) SetFloat32(offset int, value float32) {
	s.EnsureLen(offset + 4)
	s.Pointer.SetFloat32(offset, value)
}

func (s *Bytes) AppendFloat32(value float32) {
	s.Pointer.SetFloat32(s.ensureAppend(4), value)
}

///////////////////////////////////////////////////////////////////////////////////////////
// Float32 Little Endian
///////////////////////////////////////////////////////////////////////////////////////////

func (s *Bytes) Float32LE(offset int) float32 {
	s.EnsureLen(offset + 4)
	return s.Pointer.Float32LE(offset)
}

func (s *Bytes) SetFloat32LE(offset int, value float32) {
	s.EnsureLen(offset + 4)
	s.Pointer.SetFloat32LE(offset, value)
}

func (s *Bytes) AppendFloat32LE(value float32) {
	s.Pointer.SetFloat32LE(s.ensureAppend(4), value)
}

///////////////////////////////////////////////////////////////////////////////////////////
// Float32 Big Endian
///////////////////////////////////////////////////////////////////////////////////////////

func (s *Bytes) Float32BE(offset int) float32 {
	s.EnsureLen(offset + 4)
	return s.Pointer.Float32BE(offset)
}

func (s *Bytes) SetFloat32BE(offset int, value float32) {
	s.EnsureLen(offset + 4)
	s.Pointer.SetFloat32BE(offset, value)
}

func (s *Bytes) AppendFloat32BE(value float32) {
	s.Pointer.SetFloat32BE(s.ensureAppend(4), value)
}

///////////////////////////////////////////////////////////////////////////////////////////
// Float64 Native Endian
///////////////////////////////////////////////////////////////////////////////////////////

func (s *Bytes) Float64(offset int) float64 {
	s.EnsureLen(offset + 8)
	return s.Pointer.Float64(offset)
}

func (s *Bytes) SetFloat64(offset int, value float64) {
	s.EnsureLen(offset + 8)
	s.Pointer.SetFloat64(offset, value)
}

func (s *Bytes) AppendFloat64(value float64) {
	s.Pointer.SetFloat64(s.ensureAppend(8), value)
}

///////////////////////////////////////////////////////////////////////////////////////////
// Float64 Little Endian
///////////////////////////////////////////////////////////////////////////////////////////

func (s *Bytes) Float64LE(offset int) float64 {
	s.EnsureLen(offset + 8)
	return s.Pointer.Float64LE(offset)
}

func (s *Bytes) SetFloat64LE(offset int, value float64) {
	s.EnsureLen(offset + 8)
	s.Pointer.SetFloat64(offset, value)
}

func (s *Bytes) AppendFloat64LE(value float64) {
	s.Pointer.SetFloat64LE(s.ensureAppend(8), value)
}

///////////////////////////////////////////////////////////////////////////////////////////
// Float64 Big Endian
///////////////////////////////////////////////////////////////////////////////////////////

func (s *Bytes) Float64BE(offset int) float64 {
	s.EnsureLen(offset + 8)
	return s.Pointer.Float64BE(offset)
}

func (s *Bytes) SetFloat64BE(offset int, value float64) {
	s.EnsureLen(offset + 8)
	s.Pointer.SetFloat64(offset, value)
}

func (s *Bytes) AppendFloat64BE(value float64) {
	s.Pointer.SetFloat64BE(s.ensureAppend(8), value)
}

///////////////////////////////////////////////////////////////////////////////////////////
// String
///////////////////////////////////////////////////////////////////////////////////////////

func (s *Bytes) SetString(offset int, value string) {
	s.EnsureCap(offset + len(value))
	length := offset + len(value)
	if s.Len() < length {
		s.setLen(length)
	}
	s.Pointer.SetString(offset, value)
}

///////////////////////////////////////////////////////////////////////////////////////////
// Bytes
///////////////////////////////////////////////////////////////////////////////////////////

func (s *Bytes) Set(offset int, value Bytes) {
	l := value.Len()
	if value.IsNil() || l == 0 {
		return
	}
	length := offset + int(l)
	s.EnsureCap(offset + length)
	if s.Len() < length {
		s.setLen(length)
	}
	Copy(s.Unsafe(), value.Unsafe(), uintptr(l))
}

func (s *Bytes) SetBytes(offset int, value []byte) {
	s.EnsureCap(offset + len(value))
	length := offset + len(value)
	if s.Len() < length {
		s.setLen(length)
	}
	s.Pointer.SetBytes(offset, value)
}

///////////////////////////////////////////////////////////////////////////////////////////
// Int24 Native Endian
///////////////////////////////////////////////////////////////////////////////////////////

func (s *Bytes) Int24(offset int) int32 {
	s.EnsureLen(offset + 3)
	return s.Pointer.Int24(offset)
}

func (s *Bytes) SetInt24(offset int, value int32) {
	s.EnsureLen(offset + 3)
	s.Pointer.SetInt24(offset, value)
}

func (s *Bytes) AppendInt24(value int32) {
	s.Pointer.SetInt24(s.ensureAppend(3), value)
}

///////////////////////////////////////////////////////////////////////////////////////////
// Int24 Little Endian
///////////////////////////////////////////////////////////////////////////////////////////

func (s *Bytes) Int24LE(offset int) int32 {
	s.EnsureLen(offset + 3)
	return s.Pointer.Int24LE(offset)
}

func (s *Bytes) SetInt24LE(offset int, value int32) {
	s.EnsureLen(offset + 3)
	s.Pointer.SetInt24LE(offset, value)
}

func (s *Bytes) AppendInt24LE(value int32) {
	s.Pointer.SetInt24LE(s.ensureAppend(3), value)
}

///////////////////////////////////////////////////////////////////////////////////////////
// Int24 Big Endian
///////////////////////////////////////////////////////////////////////////////////////////

func (s *Bytes) Int24BE(offset int) int32 {
	s.EnsureLen(offset + 3)
	return s.Pointer.Int24BE(offset)
}

func (s *Bytes) SetInt24BE(offset int, value int32) {
	s.EnsureLen(offset + 3)
	s.Pointer.SetInt24BE(offset, value)
}

func (s *Bytes) AppendInt24BE(value int32) {
	s.Pointer.SetInt24BE(s.ensureAppend(3), value)
}

///////////////////////////////////////////////////////////////////////////////////////////
// UInt24 Native Endian
///////////////////////////////////////////////////////////////////////////////////////////

func (s *Bytes) UInt24(offset int) uint32 {
	s.EnsureLen(offset + 3)
	return s.Pointer.UInt24(offset)
}

func (s *Bytes) SetUInt24(offset int, value uint32) {
	s.EnsureLen(offset + 3)
	s.Pointer.SetUInt24(offset, value)
}

func (s *Bytes) AppendUInt24(value uint32) {
	s.Pointer.SetUInt24(s.ensureAppend(3), value)
}

///////////////////////////////////////////////////////////////////////////////////////////
// UInt24 Little Endian
///////////////////////////////////////////////////////////////////////////////////////////

func (s *Bytes) UInt24LE(offset int) uint32 {
	s.EnsureLen(offset + 3)
	return s.Pointer.UInt24LE(offset)
}

func (s *Bytes) SetUInt24LE(offset int, value uint32) {
	s.EnsureLen(offset + 3)
	s.Pointer.SetUInt24LE(offset, value)
}

func (s *Bytes) AppendUInt24LE(value uint32) {
	s.Pointer.SetUInt24LE(s.ensureAppend(3), value)
}

///////////////////////////////////////////////////////////////////////////////////////////
// UInt24 Big Endian
///////////////////////////////////////////////////////////////////////////////////////////

func (s *Bytes) UInt24BE(offset int) uint32 {
	s.EnsureLen(offset + 3)
	return s.Pointer.UInt24BE(offset)
}

func (s *Bytes) SetUInt24BE(offset int, value uint32) {
	s.EnsureLen(offset + 3)
	s.Pointer.SetUInt24BE(offset, value)
}

func (s *Bytes) AppendUInt24BE(value uint32) {
	s.Pointer.SetUInt24BE(s.ensureAppend(3), value)
}

///////////////////////////////////////////////////////////////////////////////////////////
// Int40 Native Endian
///////////////////////////////////////////////////////////////////////////////////////////

func (s *Bytes) Int40(offset int) int64 {
	s.EnsureLen(offset + 5)
	return s.Pointer.Int40(offset)
}

func (s *Bytes) SetInt40(offset int, value int64) {
	s.EnsureLen(offset + 5)
	s.Pointer.SetInt40(offset, value)
}

func (s *Bytes) AppendInt40(value int64) {
	s.Pointer.SetInt40(s.ensureAppend(5), value)
}

///////////////////////////////////////////////////////////////////////////////////////////
// Int40 Little Endian
///////////////////////////////////////////////////////////////////////////////////////////

func (s *Bytes) Int40LE(offset int) int64 {
	s.EnsureLen(offset + 5)
	return s.Pointer.Int40LE(offset)
}

func (s *Bytes) SetInt40LE(offset int, value int64) {
	s.EnsureLen(offset + 5)
	s.Pointer.SetInt40LE(offset, value)
}

func (s *Bytes) AppendInt40LE(value int64) {
	s.Pointer.SetInt40LE(s.ensureAppend(5), value)
}

///////////////////////////////////////////////////////////////////////////////////////////
// Int40 Big Endian
///////////////////////////////////////////////////////////////////////////////////////////

func (s *Bytes) Int40BE(offset int) int64 {
	s.EnsureLen(offset + 5)
	return s.Pointer.Int40BE(offset)
}

func (s *Bytes) SetInt40BE(offset int, value int64) {
	s.EnsureLen(offset + 5)
	s.Pointer.SetInt40BE(offset, value)
}

func (s *Bytes) AppendInt40BE(value int64) {
	s.Pointer.SetInt40BE(s.ensureAppend(5), value)
}

///////////////////////////////////////////////////////////////////////////////////////////
// UInt40 Native Endian
///////////////////////////////////////////////////////////////////////////////////////////

func (s *Bytes) UInt40(offset int) uint64 {
	s.EnsureLen(offset + 5)
	return s.Pointer.UInt40(offset)
}

func (s *Bytes) SetUInt40(offset int, value uint64) {
	s.EnsureLen(offset + 5)
	s.Pointer.SetUInt40(offset, value)
}

func (s *Bytes) AppendUInt40(value uint64) {
	s.Pointer.SetUInt40(s.ensureAppend(5), value)
}

///////////////////////////////////////////////////////////////////////////////////////////
// UInt40 Little Endian
///////////////////////////////////////////////////////////////////////////////////////////

func (s *Bytes) UInt40LE(offset int) uint64 {
	s.EnsureLen(offset + 5)
	return s.Pointer.UInt40LE(offset)
}

func (s *Bytes) SetUInt40LE(offset int, value uint64) {
	s.EnsureLen(offset + 5)
	s.Pointer.SetUInt40LE(offset, value)
}

func (s *Bytes) AppendUInt40LE(value uint64) {
	s.Pointer.SetUInt40LE(s.ensureAppend(5), value)
}

///////////////////////////////////////////////////////////////////////////////////////////
// UInt40 Big Endian
///////////////////////////////////////////////////////////////////////////////////////////

func (s *Bytes) UInt40BE(offset int) uint64 {
	s.EnsureLen(offset + 5)
	return s.Pointer.UInt40BE(offset)
}

func (s *Bytes) SetUInt40BE(offset int, value uint64) {
	s.EnsureLen(offset + 5)
	s.Pointer.SetUInt40BE(offset, value)
}

func (s *Bytes) AppendUInt40BE(value uint64) {
	s.Pointer.SetUInt40BE(s.ensureAppend(5), value)
}

///////////////////////////////////////////////////////////////////////////////////////////
// Int48 Native Endian
///////////////////////////////////////////////////////////////////////////////////////////

func (s *Bytes) Int48(offset int) int64 {
	s.EnsureLen(offset + 6)
	return s.Pointer.Int48(offset)
}

func (s *Bytes) SetInt48(offset int, value int64) {
	s.EnsureLen(offset + 6)
	s.Pointer.SetInt48(offset, value)
}

func (s *Bytes) AppendInt48(value int64) {
	s.Pointer.SetInt48(s.ensureAppend(6), value)
}

///////////////////////////////////////////////////////////////////////////////////////////
// Int48 Little Endian
///////////////////////////////////////////////////////////////////////////////////////////

func (s *Bytes) Int48LE(offset int) int64 {
	s.EnsureLen(offset + 6)
	return s.Pointer.Int48LE(offset)
}

func (s *Bytes) SetInt48LE(offset int, value int64) {
	s.EnsureLen(offset + 6)
	s.Pointer.SetInt48LE(offset, value)
}

func (s *Bytes) AppendInt48LE(value int64) {
	s.Pointer.SetInt48LE(s.ensureAppend(6), value)
}

///////////////////////////////////////////////////////////////////////////////////////////
// Int48 Big Endian
///////////////////////////////////////////////////////////////////////////////////////////

func (s *Bytes) Int48BE(offset int) int64 {
	s.EnsureLen(offset + 6)
	return s.Pointer.Int48BE(offset)
}

func (s *Bytes) SetInt48BE(offset int, value int64) {
	s.EnsureLen(offset + 6)
	s.Pointer.SetInt48BE(offset, value)
}

func (s *Bytes) AppendInt48BE(value int64) {
	s.Pointer.SetInt48BE(s.ensureAppend(6), value)
}

///////////////////////////////////////////////////////////////////////////////////////////
// UInt48 Native Endian
///////////////////////////////////////////////////////////////////////////////////////////

func (s *Bytes) UInt48(offset int) uint64 {
	s.EnsureLen(offset + 6)
	return s.Pointer.UInt48(offset)
}

func (s *Bytes) SetUInt48(offset int, value uint64) {
	s.EnsureLen(offset + 6)
	s.Pointer.SetUInt48(offset, value)
}

func (s *Bytes) AppendUInt48(value uint64) {
	s.Pointer.SetUInt48(s.ensureAppend(6), value)
}

///////////////////////////////////////////////////////////////////////////////////////////
// UInt48 Little Endian
///////////////////////////////////////////////////////////////////////////////////////////

func (s *Bytes) UInt48LE(offset int) uint64 {
	s.EnsureLen(offset + 6)
	return s.Pointer.UInt48LE(offset)
}

func (s *Bytes) SetUInt48LE(offset int, value uint64) {
	s.EnsureLen(offset + 6)
	s.Pointer.SetUInt48LE(offset, value)
}

func (s *Bytes) AppendUInt48LE(value uint64) {
	s.Pointer.SetUInt48LE(s.ensureAppend(6), value)
}

///////////////////////////////////////////////////////////////////////////////////////////
// UInt48 Big Endian
///////////////////////////////////////////////////////////////////////////////////////////

func (s *Bytes) UInt48BE(offset int) uint64 {
	s.EnsureLen(offset + 6)
	return s.Pointer.UInt48BE(offset)
}

func (s *Bytes) SetUInt48BE(offset int, value uint64) {
	s.EnsureLen(offset + 6)
	s.Pointer.SetUInt48BE(offset, value)
}

func (s *Bytes) AppendUInt48BE(value uint64) {
	s.Pointer.SetUInt48BE(s.ensureAppend(6), value)
}

///////////////////////////////////////////////////////////////////////////////////////////
// Int56 Native Endian
///////////////////////////////////////////////////////////////////////////////////////////

func (s *Bytes) Int56(offset int) int64 {
	s.EnsureLen(offset + 7)
	return s.Pointer.Int56(offset)
}

func (s *Bytes) SetInt56(offset int, value int64) {
	s.EnsureLen(offset + 7)
	s.Pointer.SetInt56(offset, value)
}

func (s *Bytes) AppendInt56(value int64) {
	s.Pointer.SetInt56(s.ensureAppend(7), value)
}

///////////////////////////////////////////////////////////////////////////////////////////
// Int56 Little Endian
///////////////////////////////////////////////////////////////////////////////////////////

func (s *Bytes) Int56LE(offset int) int64 {
	s.EnsureLen(offset + 7)
	return s.Pointer.Int56LE(offset)
}

func (s *Bytes) SetInt56LE(offset int, value int64) {
	s.EnsureLen(offset + 7)
	s.Pointer.SetInt56LE(offset, value)
}

func (s *Bytes) AppendInt56LE(value int64) {
	s.Pointer.SetInt56LE(s.ensureAppend(7), value)
}

///////////////////////////////////////////////////////////////////////////////////////////
// Int56 Big Endian
///////////////////////////////////////////////////////////////////////////////////////////

func (s *Bytes) Int56BE(offset int) int64 {
	s.EnsureLen(offset + 7)
	return s.Pointer.Int56BE(offset)
}

func (s *Bytes) SetInt56BE(offset int, value int64) {
	s.EnsureLen(offset + 7)
	s.Pointer.SetInt56BE(offset, value)
}

func (s *Bytes) AppendInt56BE(value int64) {
	s.Pointer.SetInt56BE(s.ensureAppend(7), value)
}

///////////////////////////////////////////////////////////////////////////////////////////
// UInt56 Native Endian
///////////////////////////////////////////////////////////////////////////////////////////

func (s *Bytes) UInt56(offset int) uint64 {
	s.EnsureLen(offset + 7)
	return s.Pointer.UInt56(offset)
}

func (s *Bytes) SetUInt56(offset int, value uint64) {
	s.EnsureLen(offset + 7)
	s.Pointer.SetUInt56(offset, value)
}

func (s *Bytes) AppendUInt56(value uint64) {
	s.Pointer.SetUInt56(s.ensureAppend(7), value)
}

///////////////////////////////////////////////////////////////////////////////////////////
// UInt56 Little Endian
///////////////////////////////////////////////////////////////////////////////////////////

func (s *Bytes) UInt56LE(offset int) uint64 {
	s.EnsureLen(offset + 7)
	return s.Pointer.UInt56LE(offset)
}

func (s *Bytes) SetUInt56LE(offset int, value uint64) {
	s.EnsureLen(offset + 7)
	s.Pointer.SetUInt56LE(offset, value)
}

func (s *Bytes) AppendUInt56LE(value uint64) {
	s.Pointer.SetUInt56LE(s.ensureAppend(7), value)
}

///////////////////////////////////////////////////////////////////////////////////////////
// UInt56 Big Endian
///////////////////////////////////////////////////////////////////////////////////////////

func (s *Bytes) UInt56BE(offset int) uint64 {
	s.EnsureLen(offset + 7)
	return s.Pointer.UInt56BE(offset)
}

func (s *Bytes) SetUInt56BE(offset int, value uint64) {
	s.EnsureLen(offset + 7)
	s.Pointer.SetUInt56BE(offset, value)
}

func (s *Bytes) AppendUInt56BE(value uint64) {
	s.Pointer.SetUInt56BE(s.ensureAppend(7), value)
}

func (s *Bytes) Hash32() uint32 {
	const (
		offset32 = uint32(2166136261)
		prime32  = uint32(16777619)
	)
	hash := offset32
	end := s.Pointer + Pointer(s.Len())
	for i := s.Pointer; i < end; i++ {
		hash ^= uint32(*(*byte)(unsafe.Pointer(i)))
		hash *= prime32
	}
	return hash
}

func (s *Bytes) Hash64() uint64 {
	const (
		offset64 = uint64(14695981039346656037)
		prime64  = uint64(1099511628211)
	)
	hash := offset64
	end := s.Pointer + Pointer(s.Len())
	for i := s.Pointer; i < end; i++ {
		hash ^= uint64(*(*byte)(unsafe.Pointer(i)))
		hash *= prime64
	}
	return hash
}

func (s *Bytes) WyHash64(seed uint64, offset, length int) uint64 {
	return hash.WyHash(s.Bytes()[offset:length], seed)
}

func (s *Bytes) Metro64(seed uint64, offset, length int) uint64 {
	return hash.Metro(s.Bytes()[offset:length], seed)
}
