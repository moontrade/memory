//go:build cgo_safe && !libfuzzer
// +build cgo_safe,!libfuzzer

package rpmalloc

// #include "rpmalloc.h"
import "C"
import (
	"unsafe"
)

// ReadThreadStats get thread statistics
func ReadThreadStats(stats *ThreadStats) {
	C.rpmalloc_thread_statistics((*C.rpmalloc_thread_statistics_t)(unsafe.Pointer(stats)))
}

// ReadGlobalStats get global statistics
func ReadGlobalStats(stats *GlobalStats) {
	C.rpmalloc_global_statistics((*C.rpmalloc_global_statistics_t)(unsafe.Pointer(stats)))
}

// Malloc allocate a memory block of at least the given size
func Malloc(size uintptr) uintptr {
	return uintptr(unsafe.Pointer(C.rpmalloc((C.size_t)(size))))
}

// MallocCap allocate a memory block of at least the given size
func MallocCap(size uintptr) (uintptr, uintptr) {
	ptr := uintptr(unsafe.Pointer(C.rpmalloc((C.size_t)(size))))
	return ptr, uintptr(C.rpmalloc_usable_size((*C.void)(ptr)))
}

// Calloc Allocates a memory block of at least the given size and zero initialize it
func Calloc(num, size uintptr) uintptr {
	return uintptr(C.rpcalloc((C.size_t)(num), (C.size_t)(size)))
}

// CallocCap Allocates a memory block of at least the given size and zero initialize it
func CallocCap(num, size uintptr) uintptr {
	ptr := uintptr(C.rpcalloc((C.size_t)(num), (C.size_t)(size)))
	return ptr, uintptr(C.rpmalloc_usable_size((*C.void)(ptr)))
}

// Realloc the given block to at least the given size
func Realloc(ptr, size uintptr) uintptr {
	return uintptr(C.rprealloc((*C.void)(ptr), (C.size_t)(size)))
}

// ReallocCap the given block to at least the given size
func ReallocCap(ptr, size uintptr) uintptr {
	newptr := uintptr(C.rprealloc((*C.void)(ptr), (C.size_t)(size)))
	return newptr, uintptr(C.rpmalloc_usable_size((*C.void)(newptr)))
}

// UsableSize Query the usable size of the given memory block (from given pointer to the end of block)
func UsableSize(ptr uintptr) uintptr {
	return uintptr(C.rpmalloc_usable_size((*C.void)(ptr)))
}

// Free the given memory block
func Free(ptr uintptr) {
	C.rpfree(unsafe.Pointer(ptr))
}
