package mem

import (
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
	Pointer
	//addr  uintptr
	len   uint32
	cap   uint32
	alloc uintptr
}

//goland:noinspection GoVetUnsafePointer
func (p *Bytes) Equals(b *Bytes) bool {
	if p.len == b.len {
		return p.Pointer == b.Pointer || memequal(unsafe.Pointer(p.Pointer), unsafe.Pointer(b.Pointer), uintptr(p.len))
	}
	return false
}

//goland:noinspection GoVetUnsafePointer
func (p *Bytes) Drop() {
	if p == nil || p.Pointer == 0 {
		return
	}
	((*Allocator)(unsafe.Pointer(p.alloc))).Free(p.Pointer)
	*p = Bytes{}
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
//func WrapMut(b []byte) Bytes {
//	return Bytes{Bytes{
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
//func WrapStringMut(s string) Bytes {
//	h := (*reflect.StringHeader)(unsafe.Bytes(&s))
//	return Bytes{Bytes{
//		ptr: unsafe.Bytes(h.Data),
//		len: uint32(len(s)),
//	}}
//}

func (p *Bytes) String() string {
	if p.IsEmpty() {
		return EmptyString
	}
	return p.Pointer.String(0, int(p.len))
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
	return *(*string)(unsafe.Pointer(&_string{
		ptr: uintptr(p.Pointer) + uintptr(offset),
		len: length,
	}))
}

func (p *Bytes) Bytes() []byte {
	if p.IsNil() {
		return nil
	}
	return p.Pointer.Bytes(0, int(p.len), int(p.len))
}

func (p *Bytes) IsNil() bool {
	return p == nil || p.Pointer == 0
}

func (p *Bytes) IsEmpty() bool {
	return p.Pointer == 0 || p.len == 0
}

func (p *Bytes) CheckBounds(offset int) bool {
	return uintptr(p.Pointer) == 0 || int(p.len) < offset
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
			Pointer: p.Add(offset),
			len:     uint32(length),
			cap:     p.cap - uint32(offset),
		},
		p: *p,
	}
}

//goland:noinspection GoVetUnsafePointer
func (p *Bytes) SetInt8(offset int, value int8) {
	if p.CheckBounds(offset + 1) {
		return
	}
	p.Pointer.SetInt8(offset, value)
}

//goland:noinspection GoVetUnsafePointer
func (p *Bytes) SetInt8Unsafe(offset int, value int8) {
	p.Pointer.SetInt8(offset, value)
}

//goland:noinspection GoVetUnsafePointer
func (p *Bytes) SetUInt8(offset int, value uint8) {
	if p.CheckBounds(offset + 1) {
		return
	}
	p.Pointer.SetUInt8(offset, value)
}

//goland:noinspection GoVetUnsafePointer
func (p *Bytes) SetByte(offset int, value byte) {
	if p.CheckBounds(offset + 1) {
		return
	}
	p.Pointer.SetByte(offset, value)
}

//goland:noinspection GoVetUnsafePointer
func (p *Bytes) SetInt16(offset int, value int16) {
	if p.CheckBounds(offset + 2) {
		return
	}
	p.Pointer.SetInt16(offset, value)
}

//goland:noinspection GoVetUnsafePointer
func (p *Bytes) SetUInt16(offset int, value uint16) {
	if p.CheckBounds(offset + 2) {
		return
	}
	p.Pointer.SetUInt16(offset, value)
}

//goland:noinspection GoVetUnsafePointer
func (p *Bytes) SetInt32(offset int, value int32) {
	if p.CheckBounds(offset + 4) {
		return
	}
	p.Pointer.SetInt32(offset, value)
}

//goland:noinspection GoVetUnsafePointer
func (p *Bytes) SetUInt32(offset int, value uint32) {
	if p.CheckBounds(offset + 4) {
		return
	}
	p.Pointer.SetUInt32(offset, value)
}

//goland:noinspection GoVetUnsafePointer
func (p *Bytes) SetInt64(offset int, value int64) {
	if p.CheckBounds(offset + 8) {
		return
	}
	p.Pointer.SetInt64(offset, value)
}

//goland:noinspection GoVetUnsafePointer
func (p *Bytes) SetUInt64(offset int, value uint64) {
	if p.CheckBounds(offset + 8) {
		return
	}
	p.Pointer.SetUInt64(offset, value)
}

//goland:noinspection GoVetUnsafePointer
func (p *Bytes) SetFloat32(offset int, value float32) {
	if p.CheckBounds(offset + 4) {
		return
	}
	p.Pointer.SetFloat32(offset, value)
}

//goland:noinspection GoVetUnsafePointer
func (p *Bytes) SetFloat64(offset int, value float64) {
	p.ensureCap(offset + 8)
	p.Pointer.SetFloat64(offset, value)
}

func (p *Bytes) SetString(offset int, value string) {
	length := offset + len(value)
	p.ensureCap(length)
	if int(p.len) < length {
		p.len = uint32(length)
	}
	dst := *(*[]byte)(unsafe.Pointer(&_bytes{
		Data: uintptr(p.Pointer.Add(offset)),
		Len:  len(value),
		Cap:  len(value),
	}))
	copy(dst, value)
}
func (p *Bytes) SetBytes(offset int, value []byte) {
	if p.CheckBounds(offset + len(value)) {
		return
	}
	dst := *(*[]byte)(unsafe.Pointer(&_bytes{
		Data: uintptr(p.Pointer.Add(offset)),
		Len:  len(value),
		Cap:  len(value),
	}))
	copy(dst, value)
}
func (p *Bytes) SetBytesUnsafe(offset int, value []byte) {
	dst := *(*[]byte)(unsafe.Pointer(&_bytes{
		Data: uintptr(p.Pointer.Add(offset)),
		Len:  len(value),
		Cap:  len(value),
	}))
	copy(dst, value)
}

//goland:noinspection GoVetUnsafePointer
func (p *Bytes) ensureCap(neededCap int) bool {
	if int(p.cap) >= neededCap {
		return true
	}
	newCap := uint32(neededCap - int(p.cap))
	addr := ((*Allocator)(unsafe.Pointer(p.alloc))).Realloc(p.Pointer, uintptr(newCap))
	if addr == 0 {
		return false
	}
	p.Pointer = addr
	p.cap = newCap
	return true
}

//goland:noinspection GoVetUnsafePointer
func (p *Bytes) ensureCapU32(neededCap uint32) bool {
	if p.cap >= neededCap {
		return true
	}
	newCap := neededCap - p.cap
	addr := ((*Allocator)(unsafe.Pointer(p.alloc))).Realloc(p.Pointer, uintptr(newCap))
	if addr == 0 {
		return false
	}
	p.Pointer = addr
	p.cap = newCap
	return true
}

func (p *Bytes) Clone() Bytes {
	b := ((*Allocator)(p.Unsafe())).Bytes(uintptr(p.len), uintptr(p.len))
	memcpy(b.Unsafe(), p.Unsafe(), uintptr(p.len))
	return b
}

// Reset zeroes out the entire allocation and sets the length back to 0
func (p *Bytes) Reset() {
	memzero(p.Unsafe(), uintptr(p.cap))
	p.len = 0
}

// Zero zeroes out the entire allocation.
func (p *Bytes) Zero() {
	memzero(p.Unsafe(), uintptr(p.cap))
}

func (p *Bytes) Append(value Bytes) int {
	p.ensureCapU32(p.len + value.len)
	i := p.len
	p.len += value.len
	return int(i)
}

func (p *Bytes) AppendString(value string) int {
	p.ensureCapU32(p.len + uint32(len(value)))
	i := p.len
	p.len += uint32(len(value))
	return int(i)
}

func (p *Bytes) SetLength(length int) {
	p.ensureCapU32(uint32(length))
	p.len = uint32(length)
}

func (p *Bytes) Extend(length int) {
	p.ensureCapU32(p.len + uint32(length))
}

func (p *Bytes) Reserve(length int) {
	p.ensureCapU32(uint32(length))
}
