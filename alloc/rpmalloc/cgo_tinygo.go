//go:build tinygo
// +build tinygo

package rpmalloc

import (
	"github.com/moontrade/memory/alloc/rpmalloc/tinygo"
)

//// ReadThreadStats get thread statistics
//func ReadThreadStats(stats *ThreadStats) {
//	rpmalloc.ReadThreadStats(stats)
//}
//
//// ReadGlobalStats get global statistics
////go:nosplit
////go:noescape
//func ReadGlobalStats(stats *GlobalStats) {
//	rpmalloc.ReadGlobalStats(stats)
//}

// Malloc allocate a memory block of at least the given size
//go:nosplit
//go:noescape
func Malloc(size uintptr) uintptr {
	return rpmalloc.Malloc(size)
}

// MallocCap allocate a memory block of at least the given size
//go:nosplit
//go:noescape
func MallocCap(size uintptr) (uintptr, uintptr) {
	return rpmalloc.MallocCap(size)
}

// MallocZero allocate a memory block of at least the given size
//go:nosplit
//go:noescape
func MallocZeroed(size uintptr) uintptr {
	return rpmalloc.MallocZeroed(size)
}

// MallocZeroCap allocate a memory block of at least the given size
//go:nosplit
//go:noescape
func MallocZeroedCap(size uintptr) (uintptr, uintptr) {
	return rpmalloc.MallocZeroedCap(size)
}

// Calloc Allocates a memory block of at least the given size and zero initialize it
//go:nosplit
//go:noescape
func Calloc(num, size uintptr) uintptr {
	return rpmalloc.Calloc(num, size)
}

// CallocCap Allocates a memory block of at least the given size and zero initialize it
//go:nosplit
//go:noescape
func CallocCap(num, size uintptr) (uintptr, uintptr) {
	return rpmalloc.CallocCap(num, size)
}

// Realloc the given block to at least the given size
//go:nosplit
//go:noescape
func Realloc(ptr, size uintptr) uintptr {
	return rpmalloc.Realloc(ptr, size)
}

// ReallocCap the given block to at least the given size
func ReallocCap(ptr, size uintptr) (uintptr, uintptr) {
	return rpmalloc.ReallocCap(ptr, size)
}

// UsableSize Query the usable size of the given memory block (from given pointer to the end of block)
func UsableSize(ptr uintptr) uintptr {
	return rpmalloc.UsableSize(ptr)
}

// Free the given memory block
func Free(ptr uintptr) {
	rpmalloc.Free(ptr)
}
