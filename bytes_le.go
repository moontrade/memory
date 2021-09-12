package mem

import (
	"reflect"
	"unsafe"
)

const (
	EmptyString        = ""
	fixedStringLenSize = 2
)

type BytesSlice struct {
	Bytes
	p Bytes
}

func (p *BytesSlice) Drop() {
	// Noop
}

// Bytes is a fat pointer to a heap allocation from an Allocator
type Bytes struct {
	addr  uintptr
	len   uint32
	cap   uint32
	alloc uintptr
}

//goland:noinspection GoVetUnsafePointer
func (p *Bytes) Drop() {
	if p == nil || p.addr == 0 {
		return
	}
	((*Allocator)(unsafe.Pointer(p.alloc))).Free(unsafe.Pointer(p.addr))
	*p = Bytes{}
}

//goland:noinspection GoVetUnsafePointer
func (p *Bytes) Unsafe() unsafe.Pointer {
	return unsafe.Pointer(p.addr)
}

func (p *Bytes) Len() int {
	return int(p.len)
}

//func Wrap(b []byte) Bytes {
//	return Bytes{
//		ptr: unsafe.Bytes(&b[0]),
//		len: uint32(len(b)),
//	}
//}
//
//func WrapMut(b []byte) BytesMut {
//	return BytesMut{Bytes{
//		ptr: unsafe.Bytes(&b[0]),
//		len: uint32(len(b)),
//	}}
//}
//
//func WrapString(s string) Bytes {
//	h := (*reflect.StringHeader)(unsafe.Bytes(&s))
//	return Bytes{
//		ptr: unsafe.Bytes(h.Data),
//		len: uint32(len(s)),
//	}
//}
//
//func WrapStringMut(s string) BytesMut {
//	h := (*reflect.StringHeader)(unsafe.Bytes(&s))
//	return BytesMut{Bytes{
//		ptr: unsafe.Bytes(h.Data),
//		len: uint32(len(s)),
//	}}
//}

func (p *Bytes) String() string {
	if p.IsEmpty() {
		return EmptyString
	}
	return *(*string)(unsafe.Pointer(&reflect.StringHeader{
		Data: p.addr,
		Len:  int(p.len),
	}))
}

func (p *Bytes) Substring(offset, length int) string {
	if p.IsEmpty() {
		return EmptyString
	}
	if offset < 0 {
		offset = 0
	}
	if length+offset > int(p.cap) {
		length = int(p.cap) - offset
	}
	if length <= 0 {
		return EmptyString
	}
	return *(*string)(unsafe.Pointer(&reflect.StringHeader{
		Data: p.addr + uintptr(offset),
		Len:  length,
	}))
}

func (p *Bytes) Bytes() []byte {
	if p.IsNil() {
		return nil
	}
	return *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: p.addr,
		Len:  int(p.len),
		Cap:  int(p.len),
	}))
}

func (p *Bytes) IsNil() bool {
	return p == nil || p.addr == 0
}

func (p *Bytes) IsEmpty() bool {
	return uintptr(p.addr) == 0 || p.len == 0
}

func (p *Bytes) CheckBounds(offset int) bool {
	return uintptr(p.addr) == 0 || int(p.len) < offset
}

//goland:noinspection GoVetUnsafePointer
func (p *Bytes) Int8Unsafe(offset int) int8 {
	return *(*int8)(unsafe.Pointer(p.addr + uintptr(offset)))
}

//goland:noinspection GoVetUnsafePointer
func (p *Bytes) Int8(offset int) int8 {
	if p.CheckBounds(offset + 1) {
		return 0
	}
	return *(*int8)(unsafe.Pointer(p.addr + uintptr(offset)))
}

//goland:noinspection GoVetUnsafePointer
func (p *Bytes) Byte(offset int) byte {
	if p.CheckBounds(offset + 1) {
		return 0
	}
	return *(*byte)(unsafe.Pointer(p.addr + uintptr(offset)))
}

//goland:noinspection GoVetUnsafePointer
func (p *Bytes) ByteUnsafe(offset int) byte {
	return *(*byte)(unsafe.Pointer(p.addr + uintptr(offset)))
}

//goland:noinspection GoVetUnsafePointer
func (p *Bytes) UInt8(offset int) uint8 {
	if p.CheckBounds(offset + 1) {
		return 0
	}
	return *(*uint8)(unsafe.Pointer(p.addr + uintptr(offset)))
}

//goland:noinspection GoVetUnsafePointer
func (p *Bytes) UInt8Unsafe(offset int) uint8 {
	return *(*uint8)(unsafe.Pointer(p.addr + uintptr(offset)))
}

//goland:noinspection GoVetUnsafePointer
func (p *Bytes) Int16(offset int) int16 {
	if p.CheckBounds(offset + 2) {
		return 0
	}
	return *(*int16)(unsafe.Pointer(p.addr + uintptr(offset)))
}

//goland:noinspection GoVetUnsafePointer
func (p *Bytes) Int16Unsafe(offset int) int16 {
	return *(*int16)(unsafe.Pointer(p.addr + uintptr(offset)))
}

//goland:noinspection GoVetUnsafePointer
func (p *Bytes) UInt16(offset int) uint16 {
	if p.CheckBounds(offset + 2) {
		return 0
	}
	return *(*uint16)(unsafe.Pointer(p.addr + uintptr(offset)))
}

//goland:noinspection GoVetUnsafePointer
func (p *Bytes) UInt16Unsafe(offset int) uint16 {
	return *(*uint16)(unsafe.Pointer(p.addr + uintptr(offset)))
}

//goland:noinspection GoVetUnsafePointer
func (p *Bytes) UInt16BE(offset int) uint16 {
	if p.CheckBounds(offset + 2) {
		return 0
	}
	return uint16(*(*byte)(unsafe.Pointer(p.addr + 1 + uintptr(offset)))) |
		uint16(*(*byte)(unsafe.Pointer(p.addr + uintptr(offset))))
}

//goland:noinspection GoVetUnsafePointer
func (p *Bytes) UInt16BEUnsafe(offset int) uint16 {
	return uint16(*(*byte)(unsafe.Pointer(p.addr + 1 + uintptr(offset)))) |
		uint16(*(*byte)(unsafe.Pointer(p.addr + uintptr(offset))))
}

//goland:noinspection GoVetUnsafePointer
func (p *Bytes) Int32(offset int) int32 {
	if p.CheckBounds(offset + 4) {
		return 0
	}
	return *(*int32)(unsafe.Pointer(p.addr + uintptr(offset)))
}

//goland:noinspection GoVetUnsafePointer
func (p *Bytes) Int32Unsafe(offset int) int32 {
	return *(*int32)(unsafe.Pointer(p.addr + uintptr(offset)))
}

//goland:noinspection GoVetUnsafePointer
func (p *Bytes) UInt32(offset int) uint32 {
	if p.CheckBounds(offset + 4) {
		return 0
	}
	return *(*uint32)(unsafe.Pointer(p.addr + uintptr(offset)))
}

//goland:noinspection GoVetUnsafePointer
func (p *Bytes) UInt32Unsafe(offset int) uint32 {
	return *(*uint32)(unsafe.Pointer(p.addr + uintptr(offset)))
}

//goland:noinspection GoVetUnsafePointer
func (p *Bytes) Int64(offset int) int64 {
	if p.CheckBounds(offset + 8) {
		return 0
	}
	return *(*int64)(unsafe.Pointer(p.addr + uintptr(offset)))
}

//goland:noinspection GoVetUnsafePointer
func (p *Bytes) Int64Unsafe(offset int) int64 {
	return *(*int64)(unsafe.Pointer(p.addr + uintptr(offset)))
}

//goland:noinspection GoVetUnsafePointer
func (p *Bytes) UInt64(offset int) uint64 {
	if p.CheckBounds(offset + 8) {
		return 0
	}
	return *(*uint64)(unsafe.Pointer(p.addr + uintptr(offset)))
}

//goland:noinspection GoVetUnsafePointer
func (p *Bytes) UInt64Unsafe(offset int) uint64 {
	return *(*uint64)(unsafe.Pointer(p.addr + uintptr(offset)))
}

//goland:noinspection GoVetUnsafePointer
func (p *Bytes) Float32(offset int) float32 {
	if p.CheckBounds(offset + 4) {
		return 0
	}
	return *(*float32)(unsafe.Pointer(p.addr + uintptr(offset)))
}

//goland:noinspection GoVetUnsafePointer
func (p *Bytes) Float32Unsafe(offset int) float32 {
	return *(*float32)(unsafe.Pointer(p.addr + uintptr(offset)))
}

//goland:noinspection GoVetUnsafePointer
func (p *Bytes) Float64(offset int) float64 {
	if p.CheckBounds(offset + 8) {
		return 0
	}
	return *(*float64)(unsafe.Pointer(p.addr + uintptr(offset)))
}

//goland:noinspection GoVetUnsafePointer
func (p *Bytes) Float64Unsafe(offset int) float64 {
	return *(*float64)(unsafe.Pointer(p.addr + uintptr(offset)))
}

//func (p *Bytes) Substr(offset, length int) string {
//	return p.Slice(offset, length).String()
//}
//
//func (p *Bytes) SliceBytes(offset, length int) []byte {
//	return p.Slice(offset, length).Bytes()
//}
//
func (p *Bytes) Slice(offset, length int) BytesSlice {
	if p.IsNil() {
		return BytesSlice{}
	}
	if offset+length > int(p.len) {
		return BytesSlice{}
	}
	return BytesSlice{
		Bytes: Bytes{
			addr: p.addr + uintptr(offset),
			len:  uint32(length),
			cap:  p.cap - uint32(offset),
		},
		p: *p,
	}
}

func (p *Bytes) Mut() BytesMut {
	if p == nil {
		return BytesMut{}
	}
	return *(*BytesMut)(unsafe.Pointer(p))
}

type BytesMut struct {
	Bytes
}

//goland:noinspection GoVetUnsafePointer
func (p *BytesMut) SetInt8(offset int, value int8) {
	if p.CheckBounds(offset + 1) {
		return
	}
	*(*int8)(unsafe.Pointer(p.addr + uintptr(offset))) = value
}

//goland:noinspection GoVetUnsafePointer
func (p *BytesMut) SetInt8Unsafe(offset int, value int8) {
	*(*int8)(unsafe.Pointer(p.addr + uintptr(offset))) = value
}

//goland:noinspection GoVetUnsafePointer
func (p *BytesMut) SetUInt8(offset int, value uint8) {
	if p.CheckBounds(offset + 1) {
		return
	}
	*(*uint8)(unsafe.Pointer(p.addr + uintptr(offset))) = value
}

//goland:noinspection GoVetUnsafePointer
func (p *BytesMut) SetUInt8Unsafe(offset int, value uint8) {
	*(*uint8)(unsafe.Pointer(p.addr + uintptr(offset))) = value
}

//goland:noinspection GoVetUnsafePointer
func (p *BytesMut) SetByte(offset int, value byte) {
	if p.CheckBounds(offset + 1) {
		return
	}
	*(*byte)(unsafe.Pointer(p.addr + uintptr(offset))) = value
}

//goland:noinspection GoVetUnsafePointer
func (p *BytesMut) SetByteUnsafe(offset int, value byte) {
	*(*byte)(unsafe.Pointer(p.addr + uintptr(offset))) = value
}

//goland:noinspection GoVetUnsafePointer
func (p *BytesMut) SetInt16(offset int, value int16) {
	if p.CheckBounds(offset + 2) {
		return
	}
	*(*int16)(unsafe.Pointer(p.addr + uintptr(offset))) = value
}

//goland:noinspection GoVetUnsafePointer
func (p *BytesMut) SetInt16Unsafe(offset int, value int16) {
	*(*int16)(unsafe.Pointer(p.addr + uintptr(offset))) = value
}
func (p *BytesMut) SetUInt16(offset int, value uint16) {
	if p.CheckBounds(offset + 2) {
		return
	}
	*(*uint16)(unsafe.Pointer(p.addr + uintptr(offset))) = value
}

//goland:noinspection GoVetUnsafePointer
func (p *BytesMut) SetUInt16Unsafe(offset int, value uint16) {
	*(*uint16)(unsafe.Pointer(p.addr + uintptr(offset))) = value
}

//goland:noinspection GoVetUnsafePointer
func (p *BytesMut) SetInt32(offset int, value int32) {
	if p.CheckBounds(offset + 4) {
		return
	}
	*(*int32)(unsafe.Pointer(p.addr + uintptr(offset))) = value
}

//goland:noinspection GoVetUnsafePointer
func (p *BytesMut) SetInt32Unsafe(offset int, value int32) {
	*(*int32)(unsafe.Pointer(p.addr + uintptr(offset))) = value
}

//goland:noinspection GoVetUnsafePointer
func (p *BytesMut) SetUInt32(offset int, value uint32) {
	if p.CheckBounds(offset + 4) {
		return
	}
	*(*uint32)(unsafe.Pointer(p.addr + uintptr(offset))) = value
}

//goland:noinspection GoVetUnsafePointer
func (p *BytesMut) SetUInt32Unsafe(offset int, value uint32) {
	*(*uint32)(unsafe.Pointer(p.addr + uintptr(offset))) = value
}

//goland:noinspection GoVetUnsafePointer
func (p *BytesMut) SetInt64(offset int, value int64) {
	if p.CheckBounds(offset + 8) {
		return
	}
	*(*int64)(unsafe.Pointer(p.addr + uintptr(offset))) = value
}

//goland:noinspection GoVetUnsafePointer
func (p *BytesMut) SetInt64Unsafe(offset int, value int64) {
	*(*int64)(unsafe.Pointer(p.addr + uintptr(offset))) = value
}

//goland:noinspection GoVetUnsafePointer
func (p *BytesMut) SetUInt64(offset int, value uint64) {
	if p.CheckBounds(offset + 8) {
		return
	}
	*(*uint64)(unsafe.Pointer(p.addr + uintptr(offset))) = value
}

//goland:noinspection GoVetUnsafePointer
func (p *BytesMut) SetUInt64Unsafe(offset int, value uint64) {
	*(*uint64)(unsafe.Pointer(p.addr + uintptr(offset))) = value
}

//goland:noinspection GoVetUnsafePointer
func (p *BytesMut) SetFloat32(offset int, value float32) {
	if p.CheckBounds(offset + 4) {
		return
	}
	*(*float32)(unsafe.Pointer(p.addr + uintptr(offset))) = value
}

//goland:noinspection GoVetUnsafePointer
func (p *BytesMut) SetFloat32Unsafe(offset int, value float32) {
	*(*float32)(unsafe.Pointer(p.addr + uintptr(offset))) = value
}

//goland:noinspection GoVetUnsafePointer
func (p *BytesMut) SetFloat64(offset int, value float64) {
	if p.CheckBounds(offset + 8) {
		return
	}
	*(*float64)(unsafe.Pointer(p.addr + uintptr(offset))) = value
}

//goland:noinspection GoVetUnsafePointer
func (p *BytesMut) SetFloat64Unsafe(offset int, value float64) {
	*(*float64)(unsafe.Pointer(p.addr + uintptr(offset))) = value
}
func (p *BytesMut) SetString(offset int, value string) {
	if p.CheckBounds(offset + len(value)) {
		return
	}
	dst := *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: p.addr + uintptr(offset),
		Len:  len(value),
		Cap:  len(value),
	}))
	copy(dst, value)
}
func (p *BytesMut) SetStringUnsafe(offset int, value string) {
	dst := *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: p.addr + uintptr(offset),
		Len:  len(value),
		Cap:  len(value),
	}))
	copy(dst, value)
}
func (p *BytesMut) SetBytes(offset int, value []byte) {
	if p.CheckBounds(offset + len(value)) {
		return
	}
	dst := *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: p.addr + uintptr(offset),
		Len:  len(value),
		Cap:  len(value),
	}))
	copy(dst, value)
}
func (p *BytesMut) SetBytesUnsafe(offset int, value []byte) {
	dst := *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: p.addr + uintptr(offset),
		Len:  len(value),
		Cap:  len(value),
	}))
	copy(dst, value)
}

//goland:noinspection GoVetUnsafePointer
func (p *BytesMut) grow(extra int) int {
	minNeeded := (int(p.len) + extra) - int(p.cap)
	if minNeeded <= 0 {
		return 0
	}
	newCap := p.len + uint32(extra)
	addr := uintptr(((*Allocator)(unsafe.Pointer(p.alloc))).Realloc(unsafe.Pointer(p.addr), uintptr(newCap)))
	if addr == 0 {
		return -1
	}
	p.addr = addr
	p.cap = newCap
	return int(newCap)
}

func (p *BytesMut) Append(value Bytes) int {
	l := int(value.len)
	if int(p.cap-p.len) < l {
		p.grow(l)
	}
	i := p.len
	p.len += value.len
	return int(i)
}

func (p *BytesMut) AppendString(value string) int {
	if int(p.cap-p.len) < len(value) {
		p.grow(len(value))
	}
	i := p.len
	p.len += uint32(len(value))
	return int(i)
}

func (p *BytesMut) Ensure(length int) int {
	if int(p.cap) < length {
		return p.grow(length)
	}
	return int(p.len)
}

func (p *BytesMut) Reserve(length int) int {
	if int(p.cap-p.len) < length {
		p.grow(length)
	}
	i := p.len
	p.len += uint32(length)
	return int(i)
}
