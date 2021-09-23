//go:build tinygo.wasm || 386 || amd64 || amd64p32 || arm || arm64 || loong64 || mips64le || mips64p32 || mips64p32le || mipsle || ppc64le || riscv || riscv64 || wasm
// +build tinygo.wasm 386 amd64 amd64p32 arm arm64 loong64 mips64le mips64p32 mips64p32le mipsle ppc64le riscv riscv64 wasm

package alloc

import (
	"math/bits"
	"unsafe"
)

///////////////////////////////////////////////////////////////////////////////////////////////
// Int16 Little Endian
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) Int16LE(offset int) int16 {
	return *(*int16)(unsafe.Pointer(uintptr(int(p) + offset)))
}

//goland:noinspection GoVetUnsafePointer
func (p Pointer) SetInt16LE(offset int, v int16) {
	*(*int16)(unsafe.Pointer(uintptr(int(p) + offset))) = v
}

///////////////////////////////////////////////////////////////////////////////////////////////
// Int16 Big Endian
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) Int16BE(offset int) int16 {
	return int16(bits.ReverseBytes16(*(*uint16)(unsafe.Pointer(uintptr(int(p) + offset)))))
}

//goland:noinspection GoVetUnsafePointer
func (p Pointer) SetInt16BE(offset int, v int16) {
	*(*int16)(unsafe.Pointer(uintptr(int(p) + offset))) = int16(bits.ReverseBytes16(uint16(v)))
}

///////////////////////////////////////////////////////////////////////////////////////////////
// UInt16 Little Endian
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) UInt16LE(offset int) uint16 {
	return *(*uint16)(unsafe.Pointer(uintptr(int(p) + offset)))
}

//goland:noinspection GoVetUnsafePointer
func (p Pointer) SetUInt16LE(offset int, v uint16) {
	*(*uint16)(unsafe.Pointer(uintptr(int(p) + offset))) = v
}

///////////////////////////////////////////////////////////////////////////////////////////////
// UInt16 Big Endian
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) UInt16BE(offset int) uint16 {
	return bits.ReverseBytes16(*(*uint16)(unsafe.Pointer(uintptr(int(p) + offset))))
}

//goland:noinspection GoVetUnsafePointer
func (p Pointer) SetUInt16BE(offset int, v uint16) {
	*(*uint16)(unsafe.Pointer(uintptr(int(p) + offset))) = bits.ReverseBytes16(v)
}

///////////////////////////////////////////////////////////////////////////////////////////////
// Int32 Little Endian
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) Int32LE(offset int) int32 {
	return *(*int32)(unsafe.Pointer(uintptr(int(p) + offset)))
}

//goland:noinspection GoVetUnsafePointer
func (p Pointer) SetInt32LE(offset int, v int32) {
	*(*int32)(unsafe.Pointer(uintptr(int(p) + offset))) = v
}

///////////////////////////////////////////////////////////////////////////////////////////////
// Int32 Big Endian
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) Int32BE(offset int) int32 {
	return int32(bits.ReverseBytes32(*(*uint32)(unsafe.Pointer(p + Pointer(offset)))))
}

//goland:noinspection GoVetUnsafePointer
func (p Pointer) SetInt32BE(offset int, v int32) {
	*(*int32)(unsafe.Pointer(uintptr(int(p) + offset))) = int32(bits.ReverseBytes32(uint32(v)))
}

///////////////////////////////////////////////////////////////////////////////////////////////
// UInt32 Little Endian
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) UInt32LE(offset int) uint32 {
	return *(*uint32)(unsafe.Pointer(uintptr(int(p) + offset)))
}

//goland:noinspection GoVetUnsafePointer
func (p Pointer) SetUInt32LE(offset int, v uint32) {
	*(*uint32)(unsafe.Pointer(uintptr(int(p) + offset))) = v
}

///////////////////////////////////////////////////////////////////////////////////////////////
// UInt32 Big Endian
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) UInt32BESlow() uint32 {
	return uint32(*(*byte)(unsafe.Pointer(p + 3))) |
		uint32(*(*byte)(unsafe.Pointer(p + 2)))<<8 |
		uint32(*(*byte)(unsafe.Pointer(p + 1)))<<16 |
		uint32(*(*byte)(unsafe.Pointer(p)))<<24
}

//goland:noinspection GoVetUnsafePointer
func (p Pointer) UInt32BE(offset int) uint32 {
	return bits.ReverseBytes32(*(*uint32)(unsafe.Pointer(uintptr(int(p) + offset))))
}

//goland:noinspection GoVetUnsafePointer
func (p Pointer) SetUInt32BE(offset int, v uint32) {
	*(*uint32)(unsafe.Pointer(uintptr(int(p) + offset))) = bits.ReverseBytes32(v)
}

///////////////////////////////////////////////////////////////////////////////////////////////
// Int64 Little Endian
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) Int64LE(offset int) int64 {
	return *(*int64)(unsafe.Pointer(uintptr(int(p) + offset)))
}

//goland:noinspection GoVetUnsafePointer
func (p Pointer) SetInt64LE(offset int, v int64) {
	*(*int64)(unsafe.Pointer(uintptr(int(p) + offset))) = v
}

///////////////////////////////////////////////////////////////////////////////////////////////
// Int64 Big Endian
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) Int64BE(offset int) int64 {
	return int64(bits.ReverseBytes64(*(*uint64)(unsafe.Pointer(uintptr(int(p) + offset)))))
}

//goland:noinspection GoVetUnsafePointer
func (p Pointer) SetInt64BE(offset int, v int64) {
	*(*int64)(unsafe.Pointer(uintptr(int(p) + offset))) = int64(bits.ReverseBytes64(uint64(v)))
}

///////////////////////////////////////////////////////////////////////////////////////////////
// UInt64 Little Endian
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) UInt64LE(offset int) uint64 {
	return *(*uint64)(unsafe.Pointer(uintptr(int(p) + offset)))
}

//goland:noinspection GoVetUnsafePointer
func (p Pointer) SetUInt64LE(offset int, v uint64) {
	*(*uint64)(unsafe.Pointer(uintptr(int(p) + offset))) = v
}

///////////////////////////////////////////////////////////////////////////////////////////////
// UInt64 Big Endian
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) UInt64BE(offset int) uint64 {
	return bits.ReverseBytes64(*(*uint64)(unsafe.Pointer(uintptr(int(p) + offset))))
}

//goland:noinspection GoVetUnsafePointer
func (p Pointer) SetUInt64BE(offset int, v uint64) {
	*(*uint64)(unsafe.Pointer(uintptr(int(p) + offset))) = bits.ReverseBytes64(v)
}

///////////////////////////////////////////////////////////////////////////////////////////////
// Float32 Little Endian
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) Float32LE(offset int) float32 {
	return *(*float32)(unsafe.Pointer(uintptr(int(p) + offset)))
}

//goland:noinspection GoVetUnsafePointer
func (p Pointer) SetFloat32LE(offset int, v float32) {
	*(*float32)(unsafe.Pointer(uintptr(int(p) + offset))) = v
}

///////////////////////////////////////////////////////////////////////////////////////////////
// Float32 Big Endian
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) Float32BE(offset int) float32 {
	return float32(bits.ReverseBytes32(*(*uint32)(unsafe.Pointer(uintptr(int(p) + offset)))))
}

//goland:noinspection GoVetUnsafePointer
func (p Pointer) SetFloat32BE(offset int, v float32) {
	*(*float32)(unsafe.Pointer(uintptr(int(p) + offset))) = float32(bits.ReverseBytes32(uint32(v)))
}

///////////////////////////////////////////////////////////////////////////////////////////////
// Float64 Little Endian
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) Float64LE(offset int) float64 {
	return *(*float64)(unsafe.Pointer(uintptr(int(p) + offset)))
}

//goland:noinspection GoVetUnsafePointer
func (p Pointer) SetFloat64LE(offset int, v float64) {
	*(*float64)(unsafe.Pointer(uintptr(int(p) + offset))) = v
}

///////////////////////////////////////////////////////////////////////////////////////////////
// Float64 Big Endian
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) Float64BE(offset int) float64 {
	return float64(bits.ReverseBytes64(*(*uint64)(unsafe.Pointer(uintptr(int(p) + offset)))))
}

//goland:noinspection GoVetUnsafePointer
func (p Pointer) SetFloat64BE(offset int, v float64) {
	*(*float64)(unsafe.Pointer(uintptr(int(p) + offset))) = float64(bits.ReverseBytes64(uint64(v)))
}

///////////////////////////////////////////////////////////////////////////////////////////////
// Int24 Native Endian
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) Int24(offset int) int32 {
	return p.Int24LE(offset)
}

//goland:noinspection GoVetUnsafePointer
func (p Pointer) SetInt24(offset int, v int32) {
	p.SetInt24LE(offset, v)
}

///////////////////////////////////////////////////////////////////////////////////////////////
// UInt24 Native Endian
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) UInt24(offset int) uint32 {
	return p.UInt24LE(offset)
}

//goland:noinspection GoVetUnsafePointer
func (p Pointer) SetUInt24(offset int, v uint32) {
	p.SetUInt24LE(offset, v)
}

///////////////////////////////////////////////////////////////////////////////////////////////
// Int40 Native Endian
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) Int40(offset int) int64 {
	return p.Int40LE(offset)
}

//goland:noinspection GoVetUnsafePointer
func (p Pointer) SetInt40(offset int, v int64) {
	p.SetInt40LE(offset, v)
}

///////////////////////////////////////////////////////////////////////////////////////////////
// UInt40 Native Endian
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) UInt40(offset int) uint64 {
	return p.UInt40LE(offset)
}

//goland:noinspection GoVetUnsafePointer
func (p Pointer) SetUInt40(offset int, v uint64) {
	p.SetUInt40LE(offset, v)
}

///////////////////////////////////////////////////////////////////////////////////////////////
// Int48 Native Endian
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) Int48(offset int) int64 {
	return p.Int48LE(offset)
}

//goland:noinspection GoVetUnsafePointer
func (p Pointer) SetInt48(offset int, v int64) {
	p.SetInt48LE(offset, v)
}

///////////////////////////////////////////////////////////////////////////////////////////////
// UInt48 Native Endian
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) UInt48(offset int) uint64 {
	return p.UInt48LE(offset)
}

//goland:noinspection GoVetUnsafePointer
func (p Pointer) SetUInt48(offset int, v uint64) {
	p.SetUInt48LE(offset, v)
}

///////////////////////////////////////////////////////////////////////////////////////////////
// Int56 Native Endian
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) Int56(offset int) int64 {
	return p.Int56LE(offset)
}

//goland:noinspection GoVetUnsafePointer
func (p Pointer) SetInt56(offset int, v int64) {
	p.SetInt56LE(offset, v)
}

///////////////////////////////////////////////////////////////////////////////////////////////
// UInt56 Native Endian
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) UInt56(offset int) uint64 {
	return p.UInt56LE(offset)
}

//goland:noinspection GoVetUnsafePointer
func (p Pointer) SetUInt56(offset int, v uint64) {
	p.SetUInt56LE(offset, v)
}
