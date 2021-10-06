//go:build tinygo && !tiny.wasm && (darwin || (linux && !baremetal && !wasi) || (freebsd && !baremetal)) && !nintendoswitch
// +build tinygo
// +build !tiny.wasm
// +build darwin linux,!baremetal,!wasi freebsd,!baremetal
// +build !nintendoswitch

package memory

import (
	"github.com/moontrade/memory/alloc/rpmalloc"
	"unsafe"
)

////////////////////////////////////////////////////////////////////////////////////
// Global allocator convenience
////////////////////////////////////////////////////////////////////////////////////

func Init() {}

func Scope(fn func(a AutoFree)) {
	scope(fn)
}

func scope(fn func(a AutoFree)) {
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
	return Pointer(rpmalloc.MallocZeroed(size))
}

func AllocZeroedCap(size uintptr) (Pointer, uintptr) {
	p, c := rpmalloc.MallocZeroedCap(size)
	return Pointer(p), c
}

// Alloc calls Alloc on the system allocator
func Calloc(num, size uintptr) Pointer {
	return Pointer(rpmalloc.Calloc(num, size))
}

func CallocCap(num, size uintptr) (Pointer, uintptr) {
	p, c := rpmalloc.CallocCap(num, size)
	return Pointer(p), c
}

// Realloc calls Realloc on the system allocator
func Realloc(p Pointer, size uintptr) Pointer {
	return Pointer(rpmalloc.Realloc(uintptr(p), size))
}

func ReallocCap(p Pointer, size uintptr) (Pointer, uintptr) {
	newPtr, c := rpmalloc.ReallocCap(uintptr(p), size)
	return Pointer(newPtr), c
}

// Free calls Free on the system allocator
func Free(p Pointer) {
	rpmalloc.Free(uintptr(p))
}

func SizeOf(ptr Pointer) uintptr {
	return rpmalloc.UsableSize(uintptr(ptr))
}

//// Scope creates an AutoFree free list that automatically reclaims memory
//// after callback finishes.
//func (a *Heap) Scope(fn func(a AutoFree)) {
//	auto := NewAuto(a.AsAllocator(), 32)
//	fn(auto)
//	auto.Free()
//}

func initAllocator() {}

func extalloc(size uintptr) unsafe.Pointer {
	ptr := unsafe.Pointer(Alloc(size))
	//println("extalloc", uint(uintptr(ptr)))
	return ptr
}

func extfree(ptr unsafe.Pointer) {
	//println("extfree", uint(uintptr(ptr)))
	Free(Pointer(ptr))
}
