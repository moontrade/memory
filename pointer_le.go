//go:build tinygo.wasm || 386 || amd64 || amd64p32 || arm || arm64 || loong64 || mips64le || mips64p32 || mips64p32le || mipsle || ppc64le || riscv || riscv64 || wasm
// +build tinygo.wasm 386 amd64 amd64p32 arm arm64 loong64 mips64le mips64p32 mips64p32le mipsle ppc64le riscv riscv64 wasm

package mem

import (
	"math/bits"
	"unsafe"
)

///////////////////////////////////////////////////////////////////////////////////////////////
// Int16 Native Endian
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) LoadInt16() int16 {
	return *(*int16)(unsafe.Pointer(p))
}

//goland:noinspection GoVetUnsafePointer
func (p Pointer) Int16(offset int) int16 {
	return *(*int16)(unsafe.Pointer(uintptr(int(p) + offset)))
}

///////////////////////////////////////////////////////////////////////////////////////////////
// Int16 Little Endian
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) LoadInt16LE() int16 {
	return *(*int16)(unsafe.Pointer(p))
}

//goland:noinspection GoVetUnsafePointer
func (p Pointer) Int16LE(offset int) int16 {
	return *(*int16)(unsafe.Pointer(uintptr(int(p) + offset)))
}

///////////////////////////////////////////////////////////////////////////////////////////////
// Int16 Big Endian
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) LoadInt16BE() int16 {
	return int16(bits.ReverseBytes16(*(*uint16)(unsafe.Pointer(p))))
}

//goland:noinspection GoVetUnsafePointer
func (p Pointer) Int16BE(offset int) int16 {
	return int16(bits.ReverseBytes16(*(*uint16)(unsafe.Pointer(uintptr(int(p) + offset)))))
}

///////////////////////////////////////////////////////////////////////////////////////////////
// UInt16 Native Endian
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) LoadUInt16() uint16 {
	return *(*uint16)(unsafe.Pointer(p))
}

//goland:noinspection GoVetUnsafePointer
func (p Pointer) UInt16(offset int) uint16 {
	return *(*uint16)(unsafe.Pointer(uintptr(int(p) + offset)))
}

///////////////////////////////////////////////////////////////////////////////////////////////
// UInt16 Little Endian
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) LoadUInt16LE() uint16 {
	return *(*uint16)(unsafe.Pointer(p))
}

//goland:noinspection GoVetUnsafePointer
func (p Pointer) UInt16LE(offset int) uint16 {
	return *(*uint16)(unsafe.Pointer(uintptr(int(p) + offset)))
}

///////////////////////////////////////////////////////////////////////////////////////////////
// UInt16 Big Endian
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) LoadUInt16BE() uint16 {
	return bits.ReverseBytes16(*(*uint16)(unsafe.Pointer(p)))
}

//goland:noinspection GoVetUnsafePointer
func (p Pointer) UInt16BE(offset int) uint16 {
	return bits.ReverseBytes16(*(*uint16)(unsafe.Pointer(uintptr(int(p) + offset))))
}

///////////////////////////////////////////////////////////////////////////////////////////////
// Int24 Native Endian
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) LoadInt24() int32 {
	return int32(*(*byte)(unsafe.Pointer(p))) |
		int32(*(*byte)(unsafe.Pointer(p + 1)))<<8 |
		int32(*(*byte)(unsafe.Pointer(p + 2)))<<16
}

//goland:noinspection GoVetUnsafePointer
func (p Pointer) Int24(offset int) int32 {
	return int32(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset)))) |
		int32(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 1))))<<8 |
		int32(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 2))))<<16
}

///////////////////////////////////////////////////////////////////////////////////////////////
// Int24 Little Endian
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) LoadInt24LE() int32 {
	return int32(*(*byte)(unsafe.Pointer(p))) |
		int32(*(*byte)(unsafe.Pointer(p + 1)))<<8 |
		int32(*(*byte)(unsafe.Pointer(p + 2)))<<16
}

//goland:noinspection GoVetUnsafePointer
func (p Pointer) Int24LE(offset int) int32 {
	return int32(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset)))) |
		int32(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 1))))<<8 |
		int32(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 2))))<<16
}

///////////////////////////////////////////////////////////////////////////////////////////////
// Int24 Big Endian
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) LoadInt24BE() int32 {
	return int32(*(*byte)(unsafe.Pointer(p + 2))) |
		int32(*(*byte)(unsafe.Pointer(p + 1)))<<8 |
		int32(*(*byte)(unsafe.Pointer(p)))<<16
}

//goland:noinspection GoVetUnsafePointer
func (p Pointer) Int24BE(offset int) int32 {
	return int32(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 2)))) |
		int32(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 1))))<<8 |
		int32(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset))))<<16
}

///////////////////////////////////////////////////////////////////////////////////////////////
// UInt24 Native Endian
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) LoadUInt24() uint32 {
	return uint32(*(*byte)(unsafe.Pointer(p))) |
		uint32(*(*byte)(unsafe.Pointer(p + 1)))<<8 |
		uint32(*(*byte)(unsafe.Pointer(p + 2)))<<16
}

//goland:noinspection GoVetUnsafePointer
func (p Pointer) UInt24(offset int) uint32 {
	return uint32(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset)))) |
		uint32(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 1))))<<8 |
		uint32(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 2))))<<16
}

///////////////////////////////////////////////////////////////////////////////////////////////
// UInt24 Little Endian
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) LoadUInt24LE() uint32 {
	return uint32(*(*byte)(unsafe.Pointer(p))) |
		uint32(*(*byte)(unsafe.Pointer(p + 1)))<<8 |
		uint32(*(*byte)(unsafe.Pointer(p + 2)))<<16
}

//goland:noinspection GoVetUnsafePointer
func (p Pointer) UInt24LE(offset int) uint32 {
	return uint32(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset)))) |
		uint32(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 1))))<<8 |
		uint32(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 2))))<<16
}

///////////////////////////////////////////////////////////////////////////////////////////////
// UInt24 Big Endian
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) LoadUInt24BE() uint32 {
	return uint32(*(*byte)(unsafe.Pointer(p + 2))) |
		uint32(*(*byte)(unsafe.Pointer(p + 1)))<<8 |
		uint32(*(*byte)(unsafe.Pointer(p)))<<16
}

//goland:noinspection GoVetUnsafePointer
func (p Pointer) UInt24BE(offset int) uint32 {
	return uint32(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 2)))) |
		uint32(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 1))))<<8 |
		uint32(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset))))<<16
}

///////////////////////////////////////////////////////////////////////////////////////////////
// Int32 Native Endian
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) LoadInt32() int32 {
	return *(*int32)(unsafe.Pointer(p))
}

//goland:noinspection GoVetUnsafePointer
func (p Pointer) Int32(offset int) int32 {
	return *(*int32)(unsafe.Pointer(uintptr(int(p) + offset)))
}

///////////////////////////////////////////////////////////////////////////////////////////////
// Int32 Little Endian
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) LoadInt32LE() int32 {
	return *(*int32)(unsafe.Pointer(p))
}

//goland:noinspection GoVetUnsafePointer
func (p Pointer) Int32LE(offset int) int32 {
	return *(*int32)(unsafe.Pointer(uintptr(int(p) + offset)))
}

///////////////////////////////////////////////////////////////////////////////////////////////
// Int32 Big Endian
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) LoadInt32BE() int32 {
	return int32(bits.ReverseBytes32(*(*uint32)(unsafe.Pointer(p))))
}

//goland:noinspection GoVetUnsafePointer
func (p Pointer) Int32BE(offset int) int32 {
	return int32(bits.ReverseBytes32(*(*uint32)(unsafe.Pointer(p + Pointer(offset)))))
}

///////////////////////////////////////////////////////////////////////////////////////////////
// UInt32 Native Endian
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) LoadUInt32() uint32 {
	return *(*uint32)(unsafe.Pointer(p))
}

//goland:noinspection GoVetUnsafePointer
func (p Pointer) UInt32(offset int) uint32 {
	return *(*uint32)(unsafe.Pointer(uintptr(int(p) + offset)))
}

///////////////////////////////////////////////////////////////////////////////////////////////
// UInt32 Little Endian
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) LoadUInt32LE() uint32 {
	return *(*uint32)(unsafe.Pointer(p))
}

//goland:noinspection GoVetUnsafePointer
func (p Pointer) UInt32LE(offset int) uint32 {
	return *(*uint32)(unsafe.Pointer(uintptr(int(p) + offset)))
}

///////////////////////////////////////////////////////////////////////////////////////////////
// UInt32 Big Endian
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) LoadUInt32BE() uint32 {
	return bits.ReverseBytes32(*(*uint32)(unsafe.Pointer(p)))
}

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

///////////////////////////////////////////////////////////////////////////////////////////////
// Int40 Native Endian
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) LoadInt40() int64 {
	return int64(*(*byte)(unsafe.Pointer(p))) |
		int64(*(*byte)(unsafe.Pointer(p + 1)))<<8 |
		int64(*(*byte)(unsafe.Pointer(p + 2)))<<16 |
		int64(*(*byte)(unsafe.Pointer(p + 3)))<<24 |
		int64(*(*byte)(unsafe.Pointer(p + 4)))<<32
}

//goland:noinspection GoVetUnsafePointer
func (p Pointer) Int40(offset int) int64 {
	return int64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset)))) |
		int64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 1))))<<8 |
		int64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 2))))<<16 |
		int64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 3))))<<24 |
		int64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 4))))<<32
}

///////////////////////////////////////////////////////////////////////////////////////////////
// Int40 Little Endian
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) LoadInt40LE() int64 {
	return int64(*(*byte)(unsafe.Pointer(p))) |
		int64(*(*byte)(unsafe.Pointer(p + 1)))<<8 |
		int64(*(*byte)(unsafe.Pointer(p + 2)))<<16 |
		int64(*(*byte)(unsafe.Pointer(p + 3)))<<24 |
		int64(*(*byte)(unsafe.Pointer(p + 4)))<<32
}

//goland:noinspection GoVetUnsafePointer
func (p Pointer) Int40LE(offset int) int64 {
	return int64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset)))) |
		int64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 1))))<<8 |
		int64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 2))))<<16 |
		int64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 3))))<<24 |
		int64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 4))))<<32
}

///////////////////////////////////////////////////////////////////////////////////////////////
// Int40 Big Endian
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) LoadInt40BE() int64 {
	return int64(*(*byte)(unsafe.Pointer(p + 4))) |
		int64(*(*byte)(unsafe.Pointer(p + 3)))<<8 |
		int64(*(*byte)(unsafe.Pointer(p + 2)))<<16 |
		int64(*(*byte)(unsafe.Pointer(p + 1)))<<24 |
		int64(*(*byte)(unsafe.Pointer(p)))<<32
}

//goland:noinspection GoVetUnsafePointer
func (p Pointer) Int40BE(offset int) int64 {
	return int64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 4)))) |
		int64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 3))))<<8 |
		int64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 2))))<<16 |
		int64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 1))))<<24 |
		int64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset))))<<32
}

///////////////////////////////////////////////////////////////////////////////////////////////
// UInt40 Native Endian
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) LoadUInt40() uint64 {
	return uint64(*(*byte)(unsafe.Pointer(p))) |
		uint64(*(*byte)(unsafe.Pointer(p + 1)))<<8 |
		uint64(*(*byte)(unsafe.Pointer(p + 2)))<<16 |
		uint64(*(*byte)(unsafe.Pointer(p + 3)))<<24 |
		uint64(*(*byte)(unsafe.Pointer(p + 4)))<<32
}

//goland:noinspection GoVetUnsafePointer
func (p Pointer) UInt40(offset int) uint64 {
	return uint64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset)))) |
		uint64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 1))))<<8 |
		uint64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 2))))<<16 |
		uint64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 3))))<<24 |
		uint64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 4))))<<32
}

///////////////////////////////////////////////////////////////////////////////////////////////
// UInt40 Little Endian
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) LoadUInt40LE() uint64 {
	return uint64(*(*byte)(unsafe.Pointer(p))) |
		uint64(*(*byte)(unsafe.Pointer(p + 1)))<<8 |
		uint64(*(*byte)(unsafe.Pointer(p + 2)))<<16 |
		uint64(*(*byte)(unsafe.Pointer(p + 3)))<<24 |
		uint64(*(*byte)(unsafe.Pointer(p + 4)))<<32
}

//goland:noinspection GoVetUnsafePointer
func (p Pointer) UInt40LE(offset int) uint64 {
	return uint64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset)))) |
		uint64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 1))))<<8 |
		uint64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 2))))<<16 |
		uint64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 3))))<<24 |
		uint64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 4))))<<32
}

///////////////////////////////////////////////////////////////////////////////////////////////
// UInt40 Big Endian
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) LoadUInt40BE() uint64 {
	return uint64(*(*byte)(unsafe.Pointer(uintptr(int(p) + 4)))) |
		uint64(*(*byte)(unsafe.Pointer(uintptr(int(p) + 3))))<<8 |
		uint64(*(*byte)(unsafe.Pointer(uintptr(int(p) + 2))))<<16 |
		uint64(*(*byte)(unsafe.Pointer(uintptr(int(p) + 1))))<<24 |
		uint64(*(*byte)(unsafe.Pointer(uintptr(int(p)))))<<32
}

//goland:noinspection GoVetUnsafePointer
func (p Pointer) UInt40BE(offset int) uint64 {
	return uint64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 4)))) |
		uint64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 3))))<<8 |
		uint64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 2))))<<16 |
		uint64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 1))))<<24 |
		uint64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset))))<<32
}

///////////////////////////////////////////////////////////////////////////////////////////////
// Int48 Native Endian
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) LoadInt48() int64 {
	return int64(*(*byte)(unsafe.Pointer(p))) |
		int64(*(*byte)(unsafe.Pointer(p + 1)))<<8 |
		int64(*(*byte)(unsafe.Pointer(p + 2)))<<16 |
		int64(*(*byte)(unsafe.Pointer(p + 3)))<<24 |
		int64(*(*byte)(unsafe.Pointer(p + 4)))<<32 |
		int64(*(*byte)(unsafe.Pointer(p + 5)))<<40
}

//goland:noinspection GoVetUnsafePointer
func (p Pointer) Int48(offset int) int64 {
	return int64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset)))) |
		int64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 1))))<<8 |
		int64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 2))))<<16 |
		int64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 3))))<<24 |
		int64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 4))))<<32 |
		int64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 5))))<<40
}

///////////////////////////////////////////////////////////////////////////////////////////////
// Int48 Little Endian
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) LoadInt48LE() int64 {
	return int64(*(*byte)(unsafe.Pointer(p))) |
		int64(*(*byte)(unsafe.Pointer(p + 1)))<<8 |
		int64(*(*byte)(unsafe.Pointer(p + 2)))<<16 |
		int64(*(*byte)(unsafe.Pointer(p + 3)))<<24 |
		int64(*(*byte)(unsafe.Pointer(p + 4)))<<32 |
		int64(*(*byte)(unsafe.Pointer(p + 5)))<<40
}

//goland:noinspection GoVetUnsafePointer
func (p Pointer) Int48LE(offset int) int64 {
	return int64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset)))) |
		int64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 1))))<<8 |
		int64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 2))))<<16 |
		int64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 3))))<<24 |
		int64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 4))))<<32 |
		int64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 5))))<<40
}

///////////////////////////////////////////////////////////////////////////////////////////////
// Int48 Big Endian
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) LoadInt48BE() int64 {
	return int64(*(*byte)(unsafe.Pointer(p + 5))) |
		int64(*(*byte)(unsafe.Pointer(p + 4)))<<8 |
		int64(*(*byte)(unsafe.Pointer(p + 3)))<<16 |
		int64(*(*byte)(unsafe.Pointer(p + 2)))<<24 |
		int64(*(*byte)(unsafe.Pointer(p + 1)))<<32 |
		int64(*(*byte)(unsafe.Pointer(p)))<<40
}

//goland:noinspection GoVetUnsafePointer
func (p Pointer) Int48BE(offset int) int64 {
	return int64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 5)))) |
		int64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 4))))<<8 |
		int64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 3))))<<16 |
		int64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 2))))<<24 |
		int64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 1))))<<32 |
		int64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset))))<<40
}

///////////////////////////////////////////////////////////////////////////////////////////////
// UInt48 Native Endian
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) LoadUInt48() uint64 {
	return uint64(*(*byte)(unsafe.Pointer(p))) |
		uint64(*(*byte)(unsafe.Pointer(p + 1)))<<8 |
		uint64(*(*byte)(unsafe.Pointer(p + 2)))<<16 |
		uint64(*(*byte)(unsafe.Pointer(p + 3)))<<24 |
		uint64(*(*byte)(unsafe.Pointer(p + 4)))<<32 |
		uint64(*(*byte)(unsafe.Pointer(p + 5)))<<40
}

//goland:noinspection GoVetUnsafePointer
func (p Pointer) UInt48(offset int) uint64 {
	return uint64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset)))) |
		uint64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 1))))<<8 |
		uint64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 2))))<<16 |
		uint64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 3))))<<24 |
		uint64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 4))))<<32 |
		uint64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 5))))<<40
}

///////////////////////////////////////////////////////////////////////////////////////////////
// UInt48 Little Endian
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) LoadUInt48LE() uint64 {
	return uint64(*(*byte)(unsafe.Pointer(p))) |
		uint64(*(*byte)(unsafe.Pointer(p + 1)))<<8 |
		uint64(*(*byte)(unsafe.Pointer(p + 2)))<<16 |
		uint64(*(*byte)(unsafe.Pointer(p + 3)))<<24 |
		uint64(*(*byte)(unsafe.Pointer(p + 4)))<<32 |
		uint64(*(*byte)(unsafe.Pointer(p + 5)))<<40
}

//goland:noinspection GoVetUnsafePointer
func (p Pointer) UInt48LE(offset int) uint64 {
	return uint64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset)))) |
		uint64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 1))))<<8 |
		uint64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 2))))<<16 |
		uint64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 3))))<<24 |
		uint64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 4))))<<32 |
		uint64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 5))))<<40
}

///////////////////////////////////////////////////////////////////////////////////////////////
// UInt48 Big Endian
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) LoadUInt48BE() uint64 {
	return uint64(*(*byte)(unsafe.Pointer(p + 6))) |
		uint64(*(*byte)(unsafe.Pointer(p + 5)))<<8 |
		uint64(*(*byte)(unsafe.Pointer(p + 4)))<<16 |
		uint64(*(*byte)(unsafe.Pointer(p + 3)))<<24 |
		uint64(*(*byte)(unsafe.Pointer(p + 2)))<<32 |
		uint64(*(*byte)(unsafe.Pointer(p + 1)))<<40
}

//goland:noinspection GoVetUnsafePointer
func (p Pointer) UInt48BE(offset int) uint64 {
	return uint64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 5)))) |
		uint64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 4))))<<8 |
		uint64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 3))))<<16 |
		uint64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 2))))<<24 |
		uint64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 1))))<<32 |
		uint64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset))))<<40
}

///////////////////////////////////////////////////////////////////////////////////////////////
// Int56 Native Endian
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) LoadInt56() int64 {
	return int64(*(*byte)(unsafe.Pointer(p))) |
		int64(*(*byte)(unsafe.Pointer(p + 1)))<<8 |
		int64(*(*byte)(unsafe.Pointer(p + 2)))<<16 |
		int64(*(*byte)(unsafe.Pointer(p + 3)))<<24 |
		int64(*(*byte)(unsafe.Pointer(p + 4)))<<32 |
		int64(*(*byte)(unsafe.Pointer(p + 5)))<<40 |
		int64(*(*byte)(unsafe.Pointer(p + 6)))<<48
}

//goland:noinspection GoVetUnsafePointer
func (p Pointer) Int56(offset int) int64 {
	return int64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset)))) |
		int64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 1))))<<8 |
		int64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 2))))<<16 |
		int64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 3))))<<24 |
		int64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 4))))<<32 |
		int64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 5))))<<40 |
		int64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 6))))<<48
}

///////////////////////////////////////////////////////////////////////////////////////////////
// Int56 Little Endian
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) LoadInt56LE() int64 {
	return int64(*(*byte)(unsafe.Pointer(p))) |
		int64(*(*byte)(unsafe.Pointer(p + 1)))<<8 |
		int64(*(*byte)(unsafe.Pointer(p + 2)))<<16 |
		int64(*(*byte)(unsafe.Pointer(p + 3)))<<24 |
		int64(*(*byte)(unsafe.Pointer(p + 4)))<<32 |
		int64(*(*byte)(unsafe.Pointer(p + 5)))<<40 |
		int64(*(*byte)(unsafe.Pointer(p + 6)))<<48
}

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

///////////////////////////////////////////////////////////////////////////////////////////////
// Int56 Big Endian
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) LoadInt56BE() int64 {
	return int64(*(*byte)(unsafe.Pointer(p + 6))) |
		int64(*(*byte)(unsafe.Pointer(p + 5)))<<8 |
		int64(*(*byte)(unsafe.Pointer(p + 4)))<<16 |
		int64(*(*byte)(unsafe.Pointer(p + 3)))<<24 |
		int64(*(*byte)(unsafe.Pointer(p + 2)))<<32 |
		int64(*(*byte)(unsafe.Pointer(p + 1)))<<40 |
		int64(*(*byte)(unsafe.Pointer(p)))<<48
}

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

///////////////////////////////////////////////////////////////////////////////////////////////
// UInt56 Native Endian
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) LoadUInt56() uint64 {
	return uint64(*(*byte)(unsafe.Pointer(p))) |
		uint64(*(*byte)(unsafe.Pointer(p + 1)))<<8 |
		uint64(*(*byte)(unsafe.Pointer(p + 2)))<<16 |
		uint64(*(*byte)(unsafe.Pointer(p + 3)))<<24 |
		uint64(*(*byte)(unsafe.Pointer(p + 4)))<<32 |
		uint64(*(*byte)(unsafe.Pointer(p + 5)))<<40 |
		uint64(*(*byte)(unsafe.Pointer(p + 6)))<<48
}

//goland:noinspection GoVetUnsafePointer
func (p Pointer) UInt56(offset int) uint64 {
	return uint64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset)))) |
		uint64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 1))))<<8 |
		uint64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 2))))<<16 |
		uint64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 3))))<<24 |
		uint64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 4))))<<32 |
		uint64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 5))))<<40 |
		uint64(*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 6))))<<48
}

///////////////////////////////////////////////////////////////////////////////////////////////
// UInt56 Little Endian
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) LoadUInt56LE() uint64 {
	return uint64(*(*byte)(unsafe.Pointer(p))) |
		uint64(*(*byte)(unsafe.Pointer(p + 1)))<<8 |
		uint64(*(*byte)(unsafe.Pointer(p + 2)))<<16 |
		uint64(*(*byte)(unsafe.Pointer(p + 3)))<<24 |
		uint64(*(*byte)(unsafe.Pointer(p + 4)))<<32 |
		uint64(*(*byte)(unsafe.Pointer(p + 5)))<<40 |
		uint64(*(*byte)(unsafe.Pointer(p + 6)))<<48
}

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

///////////////////////////////////////////////////////////////////////////////////////////////
// UInt56 Big Endian
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) LoadUInt56BE() uint64 {
	return uint64(*(*byte)(unsafe.Pointer(p + 6))) |
		uint64(*(*byte)(unsafe.Pointer(p + 5)))<<8 |
		uint64(*(*byte)(unsafe.Pointer(p + 4)))<<16 |
		uint64(*(*byte)(unsafe.Pointer(p + 3)))<<24 |
		uint64(*(*byte)(unsafe.Pointer(p + 2)))<<32 |
		uint64(*(*byte)(unsafe.Pointer(p + 1)))<<40 |
		uint64(*(*byte)(unsafe.Pointer(p)))<<48
}

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

///////////////////////////////////////////////////////////////////////////////////////////////
// Int64 Native Endian
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) LoadInt64() int64 {
	return *(*int64)(unsafe.Pointer(p))
}

//goland:noinspection GoVetUnsafePointer
func (p Pointer) Int64(offset int) int64 {
	return *(*int64)(unsafe.Pointer(uintptr(int(p) + offset)))
}

///////////////////////////////////////////////////////////////////////////////////////////////
// Int64 Little Endian
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) LoadInt64LE() int64 {
	return *(*int64)(unsafe.Pointer(p))
}

//goland:noinspection GoVetUnsafePointer
func (p Pointer) Int64LE(offset int) int64 {
	return *(*int64)(unsafe.Pointer(uintptr(int(p) + offset)))
}

///////////////////////////////////////////////////////////////////////////////////////////////
// Int64 Big Endian
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) LoadInt64BE() int64 {
	return int64(bits.ReverseBytes64(*(*uint64)(unsafe.Pointer(p))))
}

//goland:noinspection GoVetUnsafePointer
func (p Pointer) Int64BE(offset int) int64 {
	return int64(bits.ReverseBytes64(*(*uint64)(unsafe.Pointer(uintptr(int(p) + offset)))))
}

///////////////////////////////////////////////////////////////////////////////////////////////
// UInt64 Native Endian
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) LoadUInt64() uint64 {
	return *(*uint64)(unsafe.Pointer(p))
}

//goland:noinspection GoVetUnsafePointer
func (p Pointer) UInt64(offset int) uint64 {
	return *(*uint64)(unsafe.Pointer(uintptr(int(p) + offset)))
}

///////////////////////////////////////////////////////////////////////////////////////////////
// UInt64 Little Endian
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) LoadUInt64LE() uint64 {
	return *(*uint64)(unsafe.Pointer(p))
}

//goland:noinspection GoVetUnsafePointer
func (p Pointer) UInt64LE(offset int) uint64 {
	return *(*uint64)(unsafe.Pointer(uintptr(int(p) + offset)))
}

///////////////////////////////////////////////////////////////////////////////////////////////
// UInt64 Big Endian
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) LoadUInt64BE() uint64 {
	return bits.ReverseBytes64(*(*uint64)(unsafe.Pointer(p)))
}

//goland:noinspection GoVetUnsafePointer
func (p Pointer) UInt64BE(offset int) uint64 {
	return bits.ReverseBytes64(*(*uint64)(unsafe.Pointer(uintptr(int(p) + offset))))
}

///////////////////////////////////////////////////////////////////////////////////////////////
// uintptr
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) LoadUintptr() int64 {
	return *(*int64)(unsafe.Pointer(p))
}

//goland:noinspection GoVetUnsafePointer
func (p Pointer) Uintptr(offset int) uintptr {
	return *(*uintptr)(unsafe.Pointer(uintptr(int(p) + offset)))
}

///////////////////////////////////////////////////////////////////////////////////////////////
// Pointer
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) LoadPointer() Pointer {
	return *(*Pointer)(unsafe.Pointer(p))
}

//goland:noinspection GoVetUnsafePointer
func (p Pointer) Pointer(offset int) Pointer {
	return *(*Pointer)(unsafe.Pointer(uintptr(int(p) + offset)))
}

///////////////////////////////////////////////////////////////////////////////////////////////
// Float32 Native Endian
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) LoadFloat32() float32 {
	return *(*float32)(unsafe.Pointer(p))
}

//goland:noinspection GoVetUnsafePointer
func (p Pointer) Float32(offset int) float32 {
	return *(*float32)(unsafe.Pointer(uintptr(int(p) + offset)))
}

///////////////////////////////////////////////////////////////////////////////////////////////
// Float32 Little Endian
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) LoadFloat32LE() float32 {
	return *(*float32)(unsafe.Pointer(p))
}

//goland:noinspection GoVetUnsafePointer
func (p Pointer) Float32LE(offset int) float32 {
	return *(*float32)(unsafe.Pointer(uintptr(int(p) + offset)))
}

///////////////////////////////////////////////////////////////////////////////////////////////
// Float32 Big Endian
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) LoadFloat32BE() float32 {
	return float32(bits.ReverseBytes32(*(*uint32)(unsafe.Pointer(p))))
}

//goland:noinspection GoVetUnsafePointer
func (p Pointer) Float32BE(offset int) float32 {
	return float32(bits.ReverseBytes32(*(*uint32)(unsafe.Pointer(uintptr(int(p) + offset)))))
}

///////////////////////////////////////////////////////////////////////////////////////////////
// Float64 Native Endian
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) LoadFloat64() float64 {
	return *(*float64)(unsafe.Pointer(p))
}

//goland:noinspection GoVetUnsafePointer
func (p Pointer) Float64(offset int) float64 {
	return *(*float64)(unsafe.Pointer(uintptr(int(p) + offset)))
}

///////////////////////////////////////////////////////////////////////////////////////////////
// Float64 Little Endian
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) LoadFloat64LE() float64 {
	return *(*float64)(unsafe.Pointer(p))
}

//goland:noinspection GoVetUnsafePointer
func (p Pointer) Float64LE(offset int) float64 {
	return *(*float64)(unsafe.Pointer(uintptr(int(p) + offset)))
}

///////////////////////////////////////////////////////////////////////////////////////////////
// Float64 Big Endian
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) LoadFloat64BE() float64 {
	return float64(bits.ReverseBytes64(*(*uint64)(unsafe.Pointer(p))))
}

//goland:noinspection GoVetUnsafePointer
func (p Pointer) Float64BE(offset int) float64 {
	return float64(bits.ReverseBytes64(*(*uint64)(unsafe.Pointer(uintptr(int(p) + offset)))))
}

///////////////////////////////////////////////////////////////////////////////////////////////
// Int16 Native Endian
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) SetInt16(offset int, v int16) {
	*(*int16)(unsafe.Pointer(uintptr(int(p) + offset))) = v
}

///////////////////////////////////////////////////////////////////////////////////////////////
// Int16 Little Endian
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) SetInt16LE(offset int, v int16) {
	*(*int16)(unsafe.Pointer(uintptr(int(p) + offset))) = v
}

///////////////////////////////////////////////////////////////////////////////////////////////
// Int16 Big Endian
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) SetInt16BE(offset int, v int16) {
	*(*int16)(unsafe.Pointer(uintptr(int(p) + offset))) = int16(bits.ReverseBytes16(uint16(v)))
}

///////////////////////////////////////////////////////////////////////////////////////////////
// UInt16 Native Endian
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) SetUInt16(offset int, v uint16) {
	*(*uint16)(unsafe.Pointer(uintptr(int(p) + offset))) = v
}

///////////////////////////////////////////////////////////////////////////////////////////////
// UInt16 Little Endian
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) SetUInt16LE(offset int, v uint16) {
	*(*uint16)(unsafe.Pointer(uintptr(int(p) + offset))) = v
}

///////////////////////////////////////////////////////////////////////////////////////////////
// UInt16 Big Endian
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) SetUInt16BE(offset int, v uint16) {
	*(*uint16)(unsafe.Pointer(uintptr(int(p) + offset))) = bits.ReverseBytes16(v)
}

///////////////////////////////////////////////////////////////////////////////////////////////
// Int32 Native Endian
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) SetInt32(offset int, v int32) {
	*(*int32)(unsafe.Pointer(uintptr(int(p) + offset))) = v
}

///////////////////////////////////////////////////////////////////////////////////////////////
// Int32 Little Endian
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) SetInt32LE(offset int, v int32) {
	*(*int32)(unsafe.Pointer(uintptr(int(p) + offset))) = v
}

///////////////////////////////////////////////////////////////////////////////////////////////
// Int32 Big Endian
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) SetInt32BE(offset int, v int32) {
	*(*int32)(unsafe.Pointer(uintptr(int(p) + offset))) = int32(bits.ReverseBytes32(uint32(v)))
}

///////////////////////////////////////////////////////////////////////////////////////////////
// UInt32 Native Endian
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) SetUInt32(offset int, v uint32) {
	*(*uint32)(unsafe.Pointer(uintptr(int(p) + offset))) = v
}

///////////////////////////////////////////////////////////////////////////////////////////////
// UInt32 Little Endian
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) SetUInt32LE(offset int, v uint32) {
	*(*uint32)(unsafe.Pointer(uintptr(int(p) + offset))) = v
}

///////////////////////////////////////////////////////////////////////////////////////////////
// UInt32 Big Endian
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
//func (p Pointer) SetUInt32BE(v uint32) {
//	*(*uint32)(unsafe.Pointer(p)) = bits.ReverseBytes32(v)
//}

//goland:noinspection GoVetUnsafePointer
func (p Pointer) SetUInt32BE(offset int, v uint32) {
	*(*uint32)(unsafe.Pointer(uintptr(int(p) + offset))) = bits.ReverseBytes32(v)
}

///////////////////////////////////////////////////////////////////////////////////////////////
// Int64 Native Endian
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) SetInt64(offset int, v int64) {
	*(*int64)(unsafe.Pointer(uintptr(int(p) + offset))) = v
}

///////////////////////////////////////////////////////////////////////////////////////////////
// Int64 Little Endian
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) SetInt64LE(offset int, v int64) {
	*(*int64)(unsafe.Pointer(uintptr(int(p) + offset))) = v
}

///////////////////////////////////////////////////////////////////////////////////////////////
// Int64 Big Endian
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) SetInt64BE(offset int, v int64) {
	*(*int64)(unsafe.Pointer(uintptr(int(p) + offset))) = int64(bits.ReverseBytes64(uint64(v)))
}

///////////////////////////////////////////////////////////////////////////////////////////////
// UInt64 Native Endian
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) SetUInt64(offset int, v uint64) {
	*(*uint64)(unsafe.Pointer(uintptr(int(p) + offset))) = v
}

///////////////////////////////////////////////////////////////////////////////////////////////
// UInt64 Little Endian
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) SetUInt64LE(offset int, v uint64) {
	*(*uint64)(unsafe.Pointer(uintptr(int(p) + offset))) = v
}

///////////////////////////////////////////////////////////////////////////////////////////////
// UInt64 Big Endian
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) SetUInt64BE(offset int, v uint64) {
	*(*uint64)(unsafe.Pointer(uintptr(int(p) + offset))) = bits.ReverseBytes64(v)
}

///////////////////////////////////////////////////////////////////////////////////////////////
// uintptr
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) SetUintptr(offset int, v uintptr) {
	*(*uintptr)(unsafe.Pointer(uintptr(int(p) + offset))) = v
}

///////////////////////////////////////////////////////////////////////////////////////////////
// uintptr
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) SetPointer(offset int, v Pointer) {
	*(*Pointer)(unsafe.Pointer(uintptr(int(p) + offset))) = v
}

///////////////////////////////////////////////////////////////////////////////////////////////
// Float32 Native Endian
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) SetFloat32(offset int, v float32) {
	*(*float32)(unsafe.Pointer(uintptr(int(p) + offset))) = v
}

///////////////////////////////////////////////////////////////////////////////////////////////
// Float32 Little Endian
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) SetFloat32LE(offset int, v float32) {
	*(*float32)(unsafe.Pointer(uintptr(int(p) + offset))) = v
}

///////////////////////////////////////////////////////////////////////////////////////////////
// Float32 Big Endian
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) SetFloat32BE(offset int, v float32) {
	*(*float32)(unsafe.Pointer(uintptr(int(p) + offset))) = float32(bits.ReverseBytes32(uint32(v)))
}

///////////////////////////////////////////////////////////////////////////////////////////////
// Float64 Native Endian
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) SetFloat64(offset int, v float64) {
	*(*float64)(unsafe.Pointer(uintptr(int(p) + offset))) = v
}

///////////////////////////////////////////////////////////////////////////////////////////////
// Float64 Little Endian
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) SetFloat64LE(offset int, v float64) {
	*(*float64)(unsafe.Pointer(uintptr(int(p) + offset))) = v
}

///////////////////////////////////////////////////////////////////////////////////////////////
// Float64 Big Endian
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) SetFloat64BE(offset int, v float64) {
	*(*float64)(unsafe.Pointer(uintptr(int(p) + offset))) = float64(bits.ReverseBytes64(uint64(v)))
}

///////////////////////////////////////////////////////////////////////////////////////////////
// Int24 Native Endian
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) SetInt24(offset int, v int32) {
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset))) = byte(v)
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 1))) = byte(v >> 8)
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 2))) = byte(v >> 16)
}

///////////////////////////////////////////////////////////////////////////////////////////////
// Int24 Little Endian
///////////////////////////////////////////////////////////////////////////////////////////////

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
func (p Pointer) SetInt24BE(offset int, v int32) {
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset))) = byte(v >> 16)
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 1))) = byte(v >> 8)
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 2))) = byte(v)
}

///////////////////////////////////////////////////////////////////////////////////////////////
// Int40 Native Endian
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) SetInt40(offset int, v int64) {
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset))) = byte(v)
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 1))) = byte(v >> 8)
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 2))) = byte(v >> 16)
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 3))) = byte(v >> 24)
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 4))) = byte(v >> 32)
}

///////////////////////////////////////////////////////////////////////////////////////////////
// Int40 Little Endian
///////////////////////////////////////////////////////////////////////////////////////////////

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
func (p Pointer) SetInt40BE(offset int, v int64) {
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset))) = byte(v >> 32)
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 1))) = byte(v >> 24)
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 2))) = byte(v >> 16)
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 3))) = byte(v >> 8)
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 4))) = byte(v)
}

///////////////////////////////////////////////////////////////////////////////////////////////
// UInt40 Native Endian
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) SetUInt40(offset int, v uint64) {
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset))) = byte(v)
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 1))) = byte(v >> 8)
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 2))) = byte(v >> 16)
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 3))) = byte(v >> 24)
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 4))) = byte(v >> 32)
}

///////////////////////////////////////////////////////////////////////////////////////////////
// UInt40 Little Endian
///////////////////////////////////////////////////////////////////////////////////////////////

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
func (p Pointer) SetUInt40BE(offset int, v uint64) {
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset))) = byte(v >> 32)
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 1))) = byte(v >> 24)
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 2))) = byte(v >> 16)
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 3))) = byte(v >> 8)
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 4))) = byte(v)
}

///////////////////////////////////////////////////////////////////////////////////////////////
// Int48 Native Endian
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) SetInt48(offset int, v int64) {
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset))) = byte(v)
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 1))) = byte(v >> 8)
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 2))) = byte(v >> 16)
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 3))) = byte(v >> 24)
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 4))) = byte(v >> 32)
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 5))) = byte(v >> 40)
}

///////////////////////////////////////////////////////////////////////////////////////////////
// Int48 Little Endian
///////////////////////////////////////////////////////////////////////////////////////////////

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
func (p Pointer) SetInt48BE(offset int, v int64) {
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset))) = byte(v >> 40)
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 1))) = byte(v >> 32)
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 2))) = byte(v >> 24)
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 3))) = byte(v >> 16)
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 4))) = byte(v >> 8)
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 5))) = byte(v)
}

///////////////////////////////////////////////////////////////////////////////////////////////
// UInt48 Native Endian
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) SetUInt48(offset int, v uint64) {
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset))) = byte(v)
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 1))) = byte(v >> 8)
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 2))) = byte(v >> 16)
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 3))) = byte(v >> 24)
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 4))) = byte(v >> 32)
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 5))) = byte(v >> 40)
}

///////////////////////////////////////////////////////////////////////////////////////////////
// UInt48 Little Endian
///////////////////////////////////////////////////////////////////////////////////////////////

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
func (p Pointer) SetUInt48BE(offset int, v uint64) {
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset))) = byte(v >> 40)
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 1))) = byte(v >> 32)
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 2))) = byte(v >> 24)
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 3))) = byte(v >> 16)
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 4))) = byte(v >> 8)
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 5))) = byte(v)
}

///////////////////////////////////////////////////////////////////////////////////////////////
// Int56 Native Endian
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) SetInt56(offset int, v int64) {
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset))) = byte(v)
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 1))) = byte(v >> 8)
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 2))) = byte(v >> 16)
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 3))) = byte(v >> 24)
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 4))) = byte(v >> 32)
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 5))) = byte(v >> 40)
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 6))) = byte(v >> 48)
}

///////////////////////////////////////////////////////////////////////////////////////////////
// Int56 Little Endian
///////////////////////////////////////////////////////////////////////////////////////////////

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
// UInt56 Native Endian
///////////////////////////////////////////////////////////////////////////////////////////////

//goland:noinspection GoVetUnsafePointer
func (p Pointer) SetUInt56(offset int, v uint64) {
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset))) = byte(v)
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 1))) = byte(v >> 8)
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 2))) = byte(v >> 16)
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 3))) = byte(v >> 24)
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 4))) = byte(v >> 32)
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 5))) = byte(v >> 40)
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 6))) = byte(v >> 48)
}

///////////////////////////////////////////////////////////////////////////////////////////////
// UInt56 Little Endian
///////////////////////////////////////////////////////////////////////////////////////////////

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
func (p Pointer) SetUInt56BE(offset int, v uint64) {
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset))) = byte(v >> 48)
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 1))) = byte(v >> 40)
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 2))) = byte(v >> 32)
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 3))) = byte(v >> 24)
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 4))) = byte(v >> 16)
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 5))) = byte(v >> 8)
	*(*byte)(unsafe.Pointer(uintptr(int(p) + offset + 6))) = byte(v)
}
