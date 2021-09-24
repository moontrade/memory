//go:build !cgo_safe && !libfuzzer
// +build !cgo_safe,!libfuzzer

package rpmalloc

/*
#include "rpmalloc.h"

void do_rpmalloc_thread_statistics(rpmalloc_thread_statistics_t* stats) {
	rpmalloc_thread_statistics(stats);
}

void do_rpmalloc_global_statistics(rpmalloc_global_statistics_t* stats) {
	rpmalloc_global_statistics(stats);
}

typedef struct {
	size_t size;
	size_t ptr;
} malloc_t;

void do_rpmalloc(malloc_t* args) {
	args->ptr = (size_t)rpmalloc((size_t)args->size);
}

typedef struct {
	size_t size;
	size_t ptr;
	size_t cap;
} malloc_cap_t;

void do_rpmalloc_cap(malloc_cap_t* args) {
	args->ptr = (size_t)rpmalloc((size_t)args->size);
	args->cap = (size_t)rpmalloc_usable_size((void*)args->ptr);
}

typedef struct {
	size_t num;
	size_t size;
	size_t ptr;
} calloc_t;

void do_rpcalloc(calloc_t* args) {
	args->ptr = (size_t)rpcalloc(args->num, args->size);
}

typedef struct {
	size_t num;
	size_t size;
	size_t ptr;
	size_t cap;
} calloc_cap_t;

void do_rpcalloc_cap(calloc_cap_t* args) {
	args->ptr = (size_t)rpcalloc(args->num, args->size);
	args->cap = (size_t)rpmalloc_usable_size((void*)args->ptr);
}

typedef struct {
	size_t ptr;
	size_t size;
	size_t newptr;
} realloc_t;

void do_rprealloc(realloc_t* args) {
	args->newptr = (size_t)rprealloc((void*)args->ptr, args->size);
}

typedef struct {
	size_t ptr;
	size_t size;
	size_t newptr;
	size_t cap;
} realloc_cap_t;

void do_rprealloc_cap(realloc_cap_t* args) {
	args->newptr = (size_t)rprealloc((void*)args->ptr, args->size);
	args->cap = (size_t)rpmalloc_usable_size((void*)args->newptr);
}

void do_rpfree(void* ptr) {
	rpfree(ptr);
}

typedef struct {
	size_t ptr;
	size_t size;
} usable_size_t;

void do_rpmalloc_usable_size(usable_size_t* args) {
	args->size = (size_t)rpmalloc_usable_size((void*)args->ptr);
}
*/
import "C"
import "unsafe"

//go:linkname asmcgocall runtime.asmcgocall
func asmcgocall(fn, arg unsafe.Pointer) int32

// ReadThreadStats get thread statistics
func ReadThreadStats(stats *ThreadStats) {
	asmcgocall(unsafe.Pointer(C.do_rpmalloc_thread_statistics), unsafe.Pointer(stats))
}

// ReadGlobalStats get global statistics
func ReadGlobalStats(stats *GlobalStats) {
	asmcgocall(unsafe.Pointer(C.do_rpmalloc_global_statistics), unsafe.Pointer(stats))
}

type malloc_t struct {
	size   uintptr
	retval uintptr
}

// Malloc allocate a memory block of at least the given size
func Malloc(size uintptr) uintptr {
	args := malloc_t{size: size}
	ptr := uintptr(unsafe.Pointer(&args))
	asmcgocall(unsafe.Pointer(C.do_rpmalloc), unsafe.Pointer(ptr))
	return args.retval
}

type malloc_cap_t struct {
	size uintptr
	ptr  uintptr
	cap  uintptr
}

// MallocCap allocate a memory block of at least the given size and return capacity
func MallocCap(size uintptr) (uintptr, uintptr) {
	args := malloc_cap_t{size: size}
	ptr := uintptr(unsafe.Pointer(&args))
	asmcgocall(unsafe.Pointer(C.do_rpmalloc_cap), unsafe.Pointer(ptr))
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
	asmcgocall(unsafe.Pointer(C.do_rpcalloc), unsafe.Pointer(ptr))
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
	asmcgocall(unsafe.Pointer(C.do_rpcalloc_cap), unsafe.Pointer(ptr))
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
	asmcgocall(unsafe.Pointer(C.do_rprealloc), unsafe.Pointer(p))
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
	asmcgocall(unsafe.Pointer(C.do_rprealloc_cap), unsafe.Pointer(p))
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
	asmcgocall(unsafe.Pointer(C.do_rpmalloc_usable_size), unsafe.Pointer(p))
	return args.ret
}

// Free the given memory block
func Free(ptr uintptr) {
	asmcgocall(unsafe.Pointer(C.do_rpfree), unsafe.Pointer(ptr))
}
