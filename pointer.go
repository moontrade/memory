package mem

import (
	"unsafe"
)

// Pointer is a wrapper around a raw pointer that is not unsafe.Pointer
// so Go won't confuse it for a potential GC managed pointer.
type Pointer uintptr

//goland:noinspection GoVetUnsafePointer
func (p Pointer) Unsafe() unsafe.Pointer {
	return unsafe.Pointer(p)
}

// Add is Pointer arithmetic.
func (p Pointer) Add(offset int) Pointer {
	return Pointer(uintptr(int(p) + offset))
}

///////////////////////////////////////////////////////////////////////////////////////////////
// Byte
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) LoadInt8() int8 {
	return *(*int8)(unsafe.Pointer(p))
}

//goland:noinspection GoVetUnsafePointer
func (p Pointer) Int8(offset int) int8 {
	return *(*int8)(unsafe.Pointer(uintptr(int(p) + offset)))
}

//goland:noinspection GoVetUnsafePointer
func (p Pointer) LoadUInt8() uint8 {
	return *(*uint8)(unsafe.Pointer(p))
}

//goland:noinspection GoVetUnsafePointer
func (p Pointer) UInt8(offset int) uint8 {
	return *(*uint8)(unsafe.Pointer(uintptr(int(p) + offset)))
}

//goland:noinspection GoVetUnsafePointer
func (p Pointer) LoadByte() byte {
	return *(*byte)(unsafe.Pointer(p))
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
// String
///////////////////////////////////////////////////////////////////////////////////////////////

type _string struct {
	ptr uintptr
	len int
}

//goland:noinspection GoVetUnsafePointer
func (p Pointer) LoadString(size int) string {
	return *(*string)(unsafe.Pointer(&_string{
		ptr: uintptr(p),
		len: size,
	}))
}

//goland:noinspection GoVetUnsafePointer
func (p Pointer) String(offset, size int) string {
	return *(*string)(unsafe.Pointer(&_string{
		ptr: uintptr(int(p) + offset),
		len: size,
	}))
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
func (p Pointer) LoadBytes(length, capacity int) []byte {
	return *(*[]byte)(unsafe.Pointer(&_bytes{
		Data: uintptr(p),
		Len:  length,
		Cap:  capacity,
	}))
}

//goland:noinspection GoVetUnsafePointer
func (p Pointer) Bytes(offset, length, capacity int) []byte {
	return *(*[]byte)(unsafe.Pointer(&_bytes{
		Data: uintptr(int(p) + offset),
		Len:  length,
		Cap:  capacity,
	}))
}
