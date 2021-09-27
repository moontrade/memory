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
//export alloc
func Alloc(size uintptr) Pointer {
	return Pointer(rpmalloc.Malloc(size))
}

//export allocCap
func AllocCap(size uintptr) (Pointer, uintptr) {
	p, c := rpmalloc.MallocCap(size)
	return Pointer(p), c
}

//export allocZeroed
func AllocZeroed(size uintptr) Pointer {
	return Pointer(rpmalloc.MallocZeroed(size))
}

//export allocZeroedCap
func AllocZeroedCap(size uintptr) (Pointer, uintptr) {
	p, c := rpmalloc.MallocZeroedCap(size)
	return Pointer(p), c
}

// Alloc calls Alloc on the system allocator
//export calloc
func Calloc(num, size uintptr) Pointer {
	return Pointer(rpmalloc.Calloc(num, size))
}

//export callocCap
func CallocCap(num, size uintptr) (Pointer, uintptr) {
	p, c := rpmalloc.CallocCap(num, size)
	return Pointer(p), c
}

// Realloc calls Realloc on the system allocator
//export realloc
func Realloc(p Pointer, size uintptr) Pointer {
	return Pointer(rpmalloc.Realloc(uintptr(p), size))
}

//export reallocCap
func ReallocCap(p Pointer, size uintptr) (Pointer, uintptr) {
	newPtr, c := rpmalloc.ReallocCap(uintptr(p), size)
	return Pointer(newPtr), c
}

// Free calls Free on the system allocator
//export free
func Free(p Pointer) {
	rpmalloc.Free(uintptr(p))
}

//export sizeof
func SizeOf(ptr Pointer) uintptr {
	return rpmalloc.UsableSize(uintptr(ptr))
}
