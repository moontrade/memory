package memory

import "unsafe"

type FatPointer struct {
	Pointer
	len uintptr
}

func (fp *FatPointer) Len() uintptr {
	return fp.len
}

func FatPointerOf(p Pointer, length uintptr) FatPointer {
	return FatPointer{Pointer: p, len: length}
}

func (fp FatPointer) String() string {
	return fp.Pointer.String(0, int(fp.len))
}

func (fp FatPointer) Bytes() []byte {
	return fp.Pointer.Bytes(0, int(fp.len), int(fp.len))
}
func (fp FatPointer) Clone() FatPointer {
	b := Alloc(uintptr(fp.len))
	Copy(unsafe.Pointer(b), unsafe.Pointer(fp.Pointer), uintptr(fp.len))
	return FatPointer{Pointer: b, len: fp.len}
}

func (fp FatPointer) CloneAsBytes() Bytes {
	b := AllocBytes(uintptr(fp.len))
	Copy(unsafe.Pointer(b.Pointer), unsafe.Pointer(fp.Pointer), uintptr(fp.len))
	b.setLen(int(fp.len))
	return b
}
