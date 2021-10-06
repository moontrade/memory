//go:build tinygo
// +build tinygo

package rpmalloc

/*
#include "rpmalloc.h"
#include <stdlib.h>
*/
import "C"
import (
	"unsafe"
)

func init() {
	C.rpmalloc_initialize()
}
func Init() {}

// ReadThreadStats get thread statistics
func ReadThreadStats(stats *ThreadStats) {
	C.rpmalloc_thread_statistics((*C.rpmalloc_thread_statistics_t)(unsafe.Pointer(stats)))
}

// ReadGlobalStats get global statistics
//go:nosplit
//go:noescape
func ReadGlobalStats(stats *GlobalStats) {
	C.rpmalloc_global_statistics((*C.rpmalloc_global_statistics_t)(unsafe.Pointer(stats)))
}

// Malloc allocate a memory block of at least the given size
//go:nosplit
//go:noescape
func Malloc(size uintptr) uintptr {
	return uintptr(unsafe.Pointer(C.rpmalloc((C.size_t)(size))))
}

// MallocCap allocate a memory block of at least the given size
//go:nosplit
//go:noescape
func MallocCap(size uintptr) (uintptr, uintptr) {
	ptr := uintptr(unsafe.Pointer(C.rpmalloc((C.size_t)(size))))
	return ptr, uintptr(C.rpmalloc_usable_size(unsafe.Pointer(ptr)))
}

// MallocZero allocate a memory block of at least the given size
//go:nosplit
//go:noescape
func MallocZeroed(size uintptr) uintptr {
	return uintptr(unsafe.Pointer(C.rpcalloc((C.size_t)(1), (C.size_t)(size))))
}

// MallocZeroCap allocate a memory block of at least the given size
//go:nosplit
//go:noescape
func MallocZeroedCap(size uintptr) (uintptr, uintptr) {
	ptr := uintptr(unsafe.Pointer(C.rpcalloc((C.size_t)(1), (C.size_t)(size))))
	return ptr, uintptr(C.rpmalloc_usable_size(unsafe.Pointer(ptr)))
}

// Calloc Allocates a memory block of at least the given size and zero initialize it
//go:nosplit
//go:noescape
func Calloc(num, size uintptr) uintptr {
	return uintptr(C.rpcalloc((C.size_t)(num), (C.size_t)(size)))
}

// CallocCap Allocates a memory block of at least the given size and zero initialize it
//go:nosplit
//go:noescape
func CallocCap(num, size uintptr) (uintptr, uintptr) {
	ptr := uintptr(C.rpcalloc((C.size_t)(num), (C.size_t)(size)))
	return ptr, uintptr(C.rpmalloc_usable_size(unsafe.Pointer(ptr)))
}

// Realloc the given block to at least the given size
//go:nosplit
//go:noescape
func Realloc(ptr, size uintptr) uintptr {
	return uintptr(C.rprealloc(unsafe.Pointer(ptr), (C.size_t)(size)))
}

// ReallocCap the given block to at least the given size
func ReallocCap(ptr, size uintptr) (uintptr, uintptr) {
	newptr := uintptr(C.rprealloc(unsafe.Pointer(ptr), (C.size_t)(size)))
	return newptr, uintptr(C.rpmalloc_usable_size(unsafe.Pointer(newptr)))
}

// UsableSize Query the usable size of the given memory block (from given pointer to the end of block)
func UsableSize(ptr uintptr) uintptr {
	return uintptr(C.rpmalloc_usable_size(unsafe.Pointer(ptr)))
}

// Free the given memory block
func Free(ptr uintptr) {
	C.rpfree(unsafe.Pointer(ptr))
}
