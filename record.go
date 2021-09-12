package mem

import (
	"reflect"
	"unsafe"
)

type Record struct {
	Bytes
}

func (p *Record) SetPointer(offset int, value Bytes) {
	if p.CheckBounds(offset + value.Len()) {
		return
	}
	dst := *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: p.addr + uintptr(offset),
		Len:  value.Len(),
		Cap:  value.Len(),
	}))
	copy(dst, value.String())
}
func (p *Record) SetPointerUnsafe(offset int, value Bytes) {
	if p.CheckBounds(offset + value.Len()) {
		return
	}
	dst := *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: p.addr + uintptr(offset),
		Len:  value.Len(),
		Cap:  value.Len(),
	}))
	copy(dst, value.String())
}

func (p *Record) BytesFixed(offset, size int) []byte {
	if p.CheckBounds(offset + size) {
		return nil
	}
	return *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: p.addr + uintptr(offset),
		Len:  size,
		Cap:  size,
	}))
}

func (p *Record) BytesFixedUnsafe(offset, size int) []byte {
	return *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: p.addr + uintptr(offset),
		Len:  size,
		Cap:  size,
	}))
}

//goland:noinspection GoVetUnsafePointer
func (p *Record) StringFixed(offset, max int) string {
	if p.CheckBounds(offset + max) {
		return EmptyString
	}
	sizeBytes := fixedStringLenSize
	// Read length
	var length int
	switch sizeBytes {
	case 1:
		length = int(*(*byte)(unsafe.Pointer(uintptr(p.addr) + uintptr(offset+max-sizeBytes))))
	case 2:
		length = int(*(*uint16)(unsafe.Pointer(uintptr(p.addr) + uintptr(offset+max-sizeBytes))))
		//length = int(*(*byte)(unsafe.Bytes(uintptr(p.ptr) + end))) |
		//	int(*(*byte)(unsafe.Bytes(uintptr(p.ptr) + end+1)))
	case 4:
		length = int(*(*uint32)(unsafe.Pointer(uintptr(p.addr) + uintptr(offset+max-sizeBytes))))
	default:
		return EmptyString
	}

	return *(*string)(unsafe.Pointer(&reflect.StringHeader{
		Data: p.addr + uintptr(offset),
		Len:  length,
	}))
}

//goland:noinspection GoVetUnsafePointer
func (p *Record) StringFixedUnsafe(offset, max int) string {
	// Read length
	length := int(*(*uint16)(unsafe.Pointer(p.addr + uintptr(offset+max-fixedStringLenSize))))

	return *(*string)(unsafe.Pointer(&reflect.StringHeader{
		Data: p.addr + uintptr(offset),
		Len:  length,
	}))
}

//goland:noinspection GoVetUnsafePointer
func (p *Record) StringFixedBytes(offset, max int) []byte {
	if p.CheckBounds(offset + max) {
		return nil
	}
	length := int(*(*uint16)(unsafe.Pointer(p.addr + uintptr(offset+max-fixedStringLenSize))))

	return *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: p.addr + uintptr(offset),
		Len:  length,
		Cap:  length,
	}))
}

func (p *Record) StringFixedBytesUnsafe(offset, max int) []byte {
	length := int(*(*uint16)(unsafe.Pointer(p.addr + uintptr(offset+max-fixedStringLenSize))))
	return *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: p.addr + uintptr(offset),
		Len:  length,
		Cap:  length,
	}))
}

//goland:noinspection GoVetUnsafePointer
func (p *Record) SetStringFixed(offset, max int, value string) {
	if p.CheckBounds(offset + max) {
		return
	}
	length := len(value)
	if length > max-fixedStringLenSize {
		length = max - fixedStringLenSize
		value = value[:length]
	}
	dst := *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: p.addr + uintptr(offset),
		Len:  len(value),
		Cap:  len(value),
	}))
	copy(dst, value)

	end := uintptr(offset + max - fixedStringLenSize)
	*(*byte)(unsafe.Pointer(p.addr + end)) = byte(length)
	*(*byte)(unsafe.Pointer(p.addr + end + 1)) = byte(length >> 8)
}

//goland:noinspection GoVetUnsafePointer
func (p *Record) SetStringFixedUnsafe(offset, max int, value string) {
	length := len(value)
	if length > max-fixedStringLenSize {
		length = max - fixedStringLenSize
		value = value[:length]
	}
	dst := *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: p.addr + uintptr(offset),
		Len:  len(value),
		Cap:  len(value),
	}))
	copy(dst, value)

	end := uintptr(offset + max - fixedStringLenSize)
	*(*byte)(unsafe.Pointer(p.addr + end)) = byte(length)
	*(*byte)(unsafe.Pointer(p.addr + end + 1)) = byte(length >> 8)
}

//func (p *Record) SetSize(offset int, size int32) int {
//	if size < 0 {
//		return -1
//	}
//	if size < 255 {
//		p.SetByte(offset, byte(size))
//		return 1
//	}
//	if size < 65535 {
//		p.SetUInt16(offset, uint16(size))
//		return 2
//	}
//	p.SetInt32(offset, size)
//	return 4
//}
