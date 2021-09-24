//go:build libfuzzer
// +build libfuzzer

package rpmalloc

/*
#include "rpmalloc.h"
#include <stdlib.h>
#include <memory.h>

void do_rpmalloc_thread_statistics(size_t arg0, size_t arg1) {
	rpmalloc_thread_statistics((rpmalloc_thread_statistics_t*)(void*)arg0);
}

void do_rpmalloc_global_statistics(size_t arg0, size_t arg1) {
	rpmalloc_global_statistics((rpmalloc_global_statistics_t*)(void*)arg0);
}

typedef struct {
	size_t size;
	size_t ptr;
} malloc_t;

void do_rpmalloc(size_t arg0, size_t arg1) {
	malloc_t* args = (malloc_t*)arg0;
	args->ptr = (size_t)rpmalloc((size_t)args->size);
}

typedef struct {
	size_t size;
	size_t ptr;
	size_t cap;
} malloc_cap_t;

void do_rpmalloc_cap(size_t arg0, size_t arg1) {
	malloc_cap_t* args = (malloc_cap_t*)arg0;
	args->ptr = (size_t)rpmalloc((size_t)args->size);
	args->cap = (size_t)rpmalloc_usable_size((void*)args->ptr);
}

typedef struct {
	size_t num;
	size_t size;
	size_t ptr;
} calloc_t;

void do_rpcalloc(size_t arg0, size_t arg1) {
	calloc_t* args = (calloc_t*)(void*)arg0;
	args->ptr = (size_t)rpcalloc(args->num, args->size);
}

typedef struct {
	size_t num;
	size_t size;
	size_t ptr;
	size_t cap;
} calloc_cap_t;

void do_rpcalloc_cap(size_t arg0, size_t arg1) {
	calloc_cap_t* args = (calloc_cap_t*)(void*)arg0;
	args->ptr = (size_t)rpcalloc(args->num, args->size);
	args->cap = (size_t)rpmalloc_usable_size((void*)args->ptr);
}

typedef struct {
	size_t ptr;
	size_t size;
	size_t newptr;
} realloc_t;

void do_rprealloc(size_t arg0, size_t arg1) {
	realloc_t* args = (realloc_t*)(void*)arg0;
	args->newptr = (size_t)rprealloc((void*)args->ptr, args->size);
}

typedef struct {
	size_t ptr;
	size_t size;
	size_t newptr;
	size_t cap;
} realloc_cap_t;

void do_rprealloc_cap(size_t arg0, size_t arg1) {
	realloc_cap_t* args = (realloc_cap_t*)(void*)arg0;
	args->newptr = (size_t)rprealloc((void*)args->ptr, args->size);
	args->cap = (size_t)rpmalloc_usable_size((void*)args->newptr);
}

void do_rpfree(size_t ptr, size_t arg2) {
	rpfree((void*)ptr);
}

typedef struct {
	size_t ptr;
	size_t size;
} usable_size_t;

void do_rpmalloc_usable_size(size_t arg0, size_t arg1) {
	usable_size_t* args = (usable_size_t*)arg0;
	args->size = (size_t)rpmalloc_usable_size((void*)args->ptr);
}
*/
import "C"
import "unsafe"

//go:linkname libfuzzerCall runtime.libfuzzerCall
func libfuzzerCall(fn *byte, arg0, arg1 uintptr)

// ReadThreadStats get thread statistics
func ReadThreadStats(stats *ThreadStats) {
	libfuzzerCall((*byte)(C.do_rpmalloc_thread_statistics), uintptr(unsafe.Pointer(stats)), 0)
}

// ReadGlobalStats get global statistics
func ReadGlobalStats(stats *GlobalStats) {
	libfuzzerCall((*byte)(C.do_rpmalloc_global_statistics), uintptr(unsafe.Pointer(stats)), 0)
}

type malloc_t struct {
	size uintptr
	ptr  uintptr
}

// Malloc allocate a memory block of at least the given size
func Malloc(size uintptr) uintptr {
	args := malloc_t{size: size}
	ptr := uintptr(unsafe.Pointer(&args))
	libfuzzerCall((*byte)(C.do_rpmalloc), ptr, 0)
	return args.ptr
}

type malloc_cap_t struct {
	size uintptr
	ptr  uintptr
	cap  uintptr
}

// Malloc allocate a memory block of at least the given size
func MallocCap(size uintptr) (uintptr, uintptr) {
	args := malloc_cap_t{size: size}
	ptr := uintptr(unsafe.Pointer(&args))
	libfuzzerCall((*byte)(C.do_rpmalloc), ptr, 0)
	return args.ptr, args.cap
}

type calloc_t struct {
	num  uintptr
	size uintptr
	ptr  uintptr
}

// Calloc Allocates a memory block of at least the given size and zero initialize it.
func Calloc(num, size uintptr) uintptr {
	args := calloc_t{
		num:  num,
		size: size,
	}
	ptr := uintptr(unsafe.Pointer(&args))
	libfuzzerCall((*byte)(C.do_rpcalloc), ptr, 0)
	return args.ptr
}

type calloc_cap_t struct {
	num  uintptr
	size uintptr
	ptr  uintptr
	cap  uintptr
}

// Calloc Allocates a memory block of at least the given size and zero initialize it.
func CallocCap(num, size uintptr) (uintptr, uintptr) {
	args := calloc_cap_t{
		num:  num,
		size: size,
	}
	ptr := uintptr(unsafe.Pointer(&args))
	libfuzzerCall((*byte)(C.do_rpcalloc_cap), ptr, 0)
	return args.ptr, args.cap
}

type realloc_t struct {
	ptr    uintptr
	size   uintptr
	newptr uintptr
}

// Realloc the given block to at least the given size
func Realloc(ptr, size uintptr) uintptr {
	args := realloc_t{
		ptr:  ptr,
		size: size,
	}
	p := uintptr(unsafe.Pointer(&args))
	libfuzzerCall((*byte)(C.do_rprealloc), p, 0)
	return args.newptr
}

type realloc_cap_t struct {
	ptr    uintptr
	size   uintptr
	newptr uintptr
	cap    uintptr
}

// Realloc the given block to at least the given size
func ReallocCap(ptr, size uintptr) (uintptr, uintptr) {
	args := realloc_cap_t{
		ptr:  ptr,
		size: size,
	}
	p := uintptr(unsafe.Pointer(&args))
	libfuzzerCall((*byte)(C.do_rprealloc_cap), p, 0)
	return args.newptr, args.cap
}

type usable_size_t struct {
	ptr uintptr
	ret uintptr
}

// UsableSize Query the usable size of the given memory block (from given pointer to the end of block)
func UsableSize(ptr uintptr) uintptr {
	args := usable_size_t{ptr: ptr}
	p := uintptr(unsafe.Pointer(&args))
	libfuzzerCall((*byte)(C.do_rpmalloc_usable_size), p, 0)
	return args.ret
}

// Free the given memory block
func Free(ptr uintptr) {
	libfuzzerCall((*byte)(C.do_rpfree), ptr, 0)
}
