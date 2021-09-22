package mem

import (
	"unsafe"
)

// Str is a compact single dynamic allocation to be used as an unsafe replacement for string.
type Str Pointer

func (s Str) Pointer() Pointer {
	return Pointer(s)
}

func (s *Str) Free() {
	if *s == 0 {
		return
	}
	s.Allocator().Free(s.allocationPointer())
	*s = 0
}

func (s Str) flags() uint8 {
	return *(*uint8)(unsafe.Pointer(s - 1))
}

// Go doesn't support packed structs so below are templates of the packed memory layout.
/*
type sdsHeader5 struct {
	tlsfBLOCK
	alloc byte
	flags byte
	data struct{}
}

type sdsHeader8 struct {
	//tlsfBLOCK
	len uint8
	alloc byte
	flags byte
	data struct{}
}

type sdsHeader16 struct {
	tlsfBLOCK
	len uint16
	alloc byte
	flags byte
	data struct{}
}

type sdsHeader32 struct {
	tlsfBLOCK
	len uint32
	alloc byte
	flags byte
	data struct{}
}

type sdsHeader64 struct {
	tlsfBLOCK
	len uint64
	alloc byte
	flags byte
	data struct{}
}
*/

func (s Str) Len() int {
	flags := *(*uint8)(unsafe.Pointer(s - 1))
	switch flags & _SDSTypeMask {
	case _SDSType8:
		return int(*(*uint8)(unsafe.Pointer(s - _SDS8HeaderSize)))
	case _SDSType16:
		return int(*(*uint16)(unsafe.Pointer(s - _SDS16HeaderSize)))
	case _SDSType32:
		return int(*(*uint32)(unsafe.Pointer(s - _SDS32HeaderSize)))
	case _SDSType64:
		return int(*(*uint64)(unsafe.Pointer(s - _SDS64HeaderSize)))
	}
	return 0
}

func (s Str) setLen(l int) {
	flags := *(*uint8)(unsafe.Pointer(s - 1))
	switch flags & _SDSTypeMask {
	case _SDSType8:
		*(*uint8)(unsafe.Pointer(s - _SDS8HeaderSize)) = uint8(l)
	case _SDSType16:
		*(*uint16)(unsafe.Pointer(s - _SDS16HeaderSize)) = uint16(l)
	case _SDSType32:
		*(*uint32)(unsafe.Pointer(s - _SDS32HeaderSize)) = uint32(l)
	case _SDSType64:
		*(*uint64)(unsafe.Pointer(s - _SDS64HeaderSize)) = uint64(l)
	}
}

func (s Str) Cap() int {
	flags := *(*uint8)(unsafe.Pointer(s - 1))
	switch flags & _SDSTypeMask {
	case _SDSType8:
		return int((*(*tlsfBLOCK)(unsafe.Pointer(s - _SDS8HeaderTotalSize))).mmInfo & ^_TLSFTagsMask) - _SDS8HeaderSize
	case _SDSType16:
		return int((*(*tlsfBLOCK)(unsafe.Pointer(s - _SDS16HeaderTotalSize))).mmInfo & ^_TLSFTagsMask) - _SDS16HeaderSize
	case _SDSType32:
		return int((*(*tlsfBLOCK)(unsafe.Pointer(s - _SDS32HeaderTotalSize))).mmInfo & ^_TLSFTagsMask) - _SDS32HeaderSize
	case _SDSType64:
		return int((*(*tlsfBLOCK)(unsafe.Pointer(s - _SDS64HeaderTotalSize))).mmInfo & ^_TLSFTagsMask) - _SDS64HeaderSize
	}
	return 0
}

func (s Str) Byte(offset int) byte {
	return *(*byte)(unsafe.Pointer(uintptr(int(s) + offset)))
}

func (s Str) allocationPointer() Pointer {
	flags := *(*uint8)(unsafe.Pointer(s - 1))
	switch flags & _SDSTypeMask {
	case _SDSType8:
		return Pointer(s - _SDS8HeaderSize)
	case _SDSType16:
		return Pointer(s - _SDS16HeaderSize)
	case _SDSType32:
		return Pointer(s - _SDS32HeaderSize)
	case _SDSType64:
		return Pointer(s - _SDS64HeaderSize)
	}
	return 0
}

func (s *Str) Unsafe() unsafe.Pointer {
	return s.Pointer().Unsafe()
}

// Reset zeroes out the entire allocation and sets the length back to 0
func (s *Str) Reset() {
	memzero(s.Unsafe(), uintptr(s.Cap()))
	s.setLen(0)
}

// Zero zeroes out the entire allocation.
func (s Str) Zero() {
	flags := *(*uint8)(unsafe.Pointer(s - 1))
	switch flags & _SDSTypeMask {
	case _SDSType8:
		memzero(s.Unsafe(), uintptr(((*(*tlsfBLOCK)(unsafe.Pointer(s - _SDS8HeaderTotalSize))).mmInfo & ^_TLSFTagsMask)-Pointer(_SDS8HeaderSize)))
		*(*uint8)(unsafe.Pointer(s - _SDS8HeaderSize)) = 0
	case _SDSType16:
		memzero(s.Unsafe(), uintptr(((*(*tlsfBLOCK)(unsafe.Pointer(s - _SDS16HeaderTotalSize))).mmInfo & ^_TLSFTagsMask)-Pointer(_SDS16HeaderSize)))
		*(*uint16)(unsafe.Pointer(s - _SDS16HeaderSize)) = 0
	case _SDSType32:
		memzero(s.Unsafe(), uintptr(((*(*tlsfBLOCK)(unsafe.Pointer(s - _SDS32HeaderTotalSize))).mmInfo & ^_TLSFTagsMask)-Pointer(_SDS32HeaderSize)))
		*(*uint32)(unsafe.Pointer(s - _SDS32HeaderSize)) = 0
	case _SDSType64:
		memzero(s.Unsafe(), uintptr(((*(*tlsfBLOCK)(unsafe.Pointer(s - _SDS64HeaderTotalSize))).mmInfo & ^_TLSFTagsMask)-Pointer(_SDS64HeaderSize)))
		*(*uint64)(unsafe.Pointer(s - _SDS64HeaderSize)) = 0
	}
}

//goland:noinspection GoVetUnsafePointer
func (s *Str) Equals(o Str) bool {
	l := s.Len()
	return l == o.Len() && (*s == o || memequal(
		s.Pointer().Unsafe(),
		o.Pointer().Unsafe(), uintptr(l)))
}

func (s *Str) IsNil() bool {
	return s == nil || *s == 0
}

func (s Str) IsEmpty() bool {
	return s == 0 || s.Len() == 0
}

func (s *Str) CheckBounds(offset int) bool {
	return uintptr(*s) == 0 || s.Len() < offset
}

func (s *Str) String() string {
	l := s.Len()
	if l == 0 {
		return ""
	}
	return s.Pointer().String(0, l)
}

func (s *Str) Bytes() []byte {
	l := s.Len()
	if l == 0 {
		return nil
	}
	return s.Pointer().Bytes(0, l, l)
}

func (s *Str) ensureAppend(extra int) int {
	offset := s.Len()
	s.EnsureCap(offset + extra)
	s.setLen(offset + extra)
	return offset
}

// EnsureLen ensures the length is at least neededLen in size
// If not, EnsureCap(neededLen) is called and the length set to neededLen.
func (s *Str) EnsureLen(neededLen int) {
	l := s.Len()
	if l > neededLen {
		return
	}
	s.EnsureCap(neededLen)
	s.setLen(neededLen)
}

// EnsureCap ensures the capacity is at least neededCap in size
//goland:noinspection GoVetUnsafePointer
func (s *Str) EnsureCap(neededCap int) bool {
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
	*s = Str(addr)
	return true
}

// Clone creates a copy of this instance of Bytes
func (s *Str) Clone() Str {
	l := s.Len()
	c := AllocString(uintptr(l))
	memcpy(c.Unsafe(), s.Pointer().Unsafe(), uintptr(l))
	c.setLen(l)
	return c
}

//goland:noinspection GoVetUnsafePointer
func (s *Str) Append(value Str) {
	if value == 0 {
		return
	}
	l := s.Len()
	vl := value.Len()
	s.EnsureCap(l + vl)
	memmove(unsafe.Pointer(uintptr(int(*s)+l)), value.Unsafe(), uintptr(vl))
	s.setLen(l + vl)
}

//goland:noinspection GoVetUnsafePointer
func (s *Str) AppendBytes(value []byte) {
	if len(value) == 0 {
		return
	}
	l := s.Len()
	s.EnsureCap(l + len(value))
	memmove(unsafe.Pointer(uintptr(int(*s)+l)), unsafe.Pointer(&value[0]), uintptr(len(value)))
	s.setLen(l + len(value))
}

//goland:noinspection GoVetUnsafePointer
func (s *Str) AppendString(value string) {
	if len(value) == 0 {
		return
	}
	l := s.Len()
	s.EnsureCap(l + len(value))
	memmove(
		unsafe.Pointer(uintptr(int(*s)+l)),
		unsafe.Pointer(*(*uintptr)(unsafe.Pointer(&value))),
		uintptr(len(value)),
	)
	s.setLen(l + len(value))
}

func (s *Str) SetLength(length int) {
	s.EnsureCap(length)
	s.setLen(length)
}
