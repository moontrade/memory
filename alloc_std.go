//go:build !tinygo && !wasm && !wasi && !tinygo.wasm
// +build !tinygo,!wasm,!wasi,!tinygo.wasm

package memory

import (
	"github.com/moontrade/memory/rpmalloc"
)

func Init() {}

func Scope(fn func(a Auto)) {
	scope(fn)
}

func scope(fn func(a Auto)) {
	a := NewAuto(32)
	defer a.Free()
	fn(a)
	a.Print()
}

// Alloc calls Alloc on the system allocator
func Alloc(size uintptr) Pointer {
	return Pointer(rpmalloc.Malloc(size))
}
func AllocCap(size uintptr) (Pointer, uintptr) {
	p, c := rpmalloc.MallocCap(size)
	return Pointer(p), c
}
func AllocZeroed(size uintptr) Pointer {
	return Calloc(1, size)
}
func AllocZeroedCap(size uintptr) (Pointer, uintptr) {
	return CallocCap(1, size)
}

// Alloc calls Alloc on the system allocator
//export alloc
func Calloc(num, size uintptr) Pointer {
	return Pointer(rpmalloc.Calloc(num, size))
}
func CallocCap(num, size uintptr) (Pointer, uintptr) {
	p, c := rpmalloc.CallocCap(num, size)
	return Pointer(p), c
}

// Realloc calls Realloc on the system allocator
//export realloc
func Realloc(p Pointer, size uintptr) Pointer {
	return Pointer(rpmalloc.Realloc(uintptr(p), size))
}
func ReallocCap(p Pointer, size uintptr) (Pointer, uintptr) {
	newPtr, c := rpmalloc.ReallocCap(uintptr(p), size)
	return Pointer(newPtr), c
}

// Free calls Free on the system allocator
//export free
func Free(p Pointer) {
	rpmalloc.Free(uintptr(p))
}

func SizeOf(ptr Pointer) uintptr {
	return rpmalloc.UsableSize(uintptr(ptr))
}
