//go:build libfuzzer
// +build libfuzzer

package rpmalloc

/*
#include "rpmalloc.h"
#include <stdlib.h>

typedef struct {
	size_t size;
	size_t ptr;
} malloc_t;

void do_malloc(size_t arg0, size_t arg1) {
	malloc_t* args = (malloc_t*)arg0;
	args->ptr = (size_t)malloc((size_t)args->size);
}
void do_free(size_t ptr, size_t arg1) {
	free((void*)ptr);
}

void do_rpmalloc_thread_statistics(size_t arg0, size_t arg1) {
	rpmalloc_thread_statistics((rpmalloc_thread_statistics_t*)(void*)arg0);
}

void do_rpmalloc_global_statistics(size_t arg0, size_t arg1) {
	rpmalloc_global_statistics((rpmalloc_global_statistics_t*)(void*)arg0);
}

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

typedef struct {
	size_t ptr;
} heap_acquire_t;

void do_rpmalloc_heap_acquire(size_t arg0, size_t arg1) {
	heap_acquire_t* args = (heap_acquire_t*)(void*)arg0;
	args->ptr = (size_t)rpmalloc_heap_acquire();
}

typedef struct {
	size_t ptr;
} heap_release_t;

void do_rpmalloc_heap_release(size_t arg0, size_t arg1) {
	heap_release_t* args = (heap_release_t*)(void*)arg0;
	rpmalloc_heap_release((rpmalloc_heap_t*)(void*)args->ptr);
}

typedef struct {
	size_t heap;
	size_t size;
	size_t ptr;
} heap_alloc_t;

void do_rpmalloc_heap_alloc(size_t arg0, size_t arg1) {
	heap_alloc_t* args = (heap_alloc_t*)(void*)arg0;
	args->ptr = (size_t)rpmalloc_heap_alloc((rpmalloc_heap_t*)(void*)args->heap, args->size);
}

typedef struct {
	size_t heap;
	size_t size;
	size_t ptr;
	size_t cap;
} heap_alloc_cap_t;

void do_rpmalloc_heap_alloc_cap(size_t arg0, size_t arg1) {
	heap_alloc_cap_t* args = (heap_alloc_cap_t*)(void*)arg0;
	args->ptr = (size_t)rpmalloc_heap_alloc((rpmalloc_heap_t*)(void*)args->heap, args->size);
	args->cap = (size_t)rpmalloc_usable_size((void*)args->ptr);
}

typedef struct {
	size_t heap;
	size_t num;
	size_t size;
	size_t ptr;
} heap_calloc_t;

void do_rpmalloc_heap_calloc(size_t arg0, size_t arg1) {
	heap_calloc_t* args = (heap_calloc_t*)(void*)arg0;
	args->ptr = (size_t)rpmalloc_heap_calloc((rpmalloc_heap_t*)(void*)args->heap, args->num, args->size);
}

typedef struct {
	size_t heap;
	size_t num;
	size_t size;
	size_t ptr;
	size_t cap;
} heap_calloc_cap_t;

void do_rpmalloc_heap_calloc_cap(size_t arg0, size_t arg1) {
	heap_calloc_cap_t* args = (heap_calloc_cap_t*)(void*)arg0;
	args->ptr = (size_t)rpmalloc_heap_calloc((rpmalloc_heap_t*)(void*)args->heap, args->num, args->size);
	args->cap = (size_t)rpmalloc_usable_size((void*)args->ptr);
}

typedef struct {
	size_t heap;
	size_t ptr;
	size_t size;
	size_t newptr;
	int flags;
} heap_realloc_t;

void do_rpmalloc_heap_realloc(size_t arg0, size_t arg1) {
	heap_realloc_t* args = (heap_realloc_t*)(void*)arg0;
	args->newptr = (size_t)rpmalloc_heap_realloc((rpmalloc_heap_t*)(void*)args->heap, (void*)args->ptr, args->size, args->flags);
}

typedef struct {
	size_t heap;
	size_t ptr;
	size_t size;
	size_t newptr;
	size_t cap;
	int flags;
} heap_realloc_cap_t;

void do_rpmalloc_heap_realloc_cap(size_t arg0, size_t arg1) {
	heap_realloc_cap_t* args = (heap_realloc_cap_t*)(void*)arg0;
	args->newptr = (size_t)rpmalloc_heap_realloc((rpmalloc_heap_t*)(void*)args->heap, (void*)args->ptr, args->size, args->flags);
	args->cap = (size_t)rpmalloc_usable_size((void*)args->newptr);
}

typedef struct {
	size_t heap;
	size_t ptr;
} heap_free_t;

void do_rpmalloc_heap_free(size_t arg0, size_t arg1) {
	heap_free_t* args = (heap_free_t*)(void*)arg0;
	rpmalloc_heap_free((rpmalloc_heap_t*)(void*)args->heap, (void*)args->ptr);
}

typedef struct {
	size_t heap;
} heap_free_all_t;

void do_rpmalloc_heap_free_all(size_t arg0, size_t arg1) {
	heap_free_all_t* args = (heap_free_all_t*)(void*)arg0;
	rpmalloc_heap_free_all((rpmalloc_heap_t*)(void*)args->heap);
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

// Malloc allocate a memory block of at least the given size
func StdMalloc(size uintptr) uintptr {
	args := malloc_t{size: size}
	ptr := uintptr(unsafe.Pointer(&args))
	libfuzzerCall((*byte)(C.do_malloc), ptr, 0)
	return args.ptr
}

// Free the given memory block
func StdFree(ptr uintptr) {
	libfuzzerCall((*byte)(C.do_free), ptr, 0)
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

func InitThread() {
	C.rpmalloc_thread_initialize()
}

type acquire_heap_t struct {
	ptr uintptr
}

func AcquireHeap() *Heap {
	args := acquire_heap_t{}
	ptr := uintptr(unsafe.Pointer(&args))
	libfuzzerCall((*byte)(C.do_rpmalloc_heap_acquire), ptr, 0)
	return (*Heap)(unsafe.Pointer(args.ptr))
}

type release_heap_t struct {
	heap uintptr
}

func (h *Heap) Release() {
	args := release_heap_t{heap: uintptr(unsafe.Pointer(h))}
	ptr := uintptr(unsafe.Pointer(&args))
	libfuzzerCall((*byte)(C.do_rpmalloc_heap_release), ptr, 0)
}

type heap_alloc_t struct {
	heap uintptr
	size uintptr
	ptr  uintptr
}

// Alloc Allocate a memory block of at least the given size using the given heap.
func (h *Heap) Alloc(size uintptr) uintptr {
	args := heap_alloc_t{heap: uintptr(unsafe.Pointer(h)), size: size}
	ptr := uintptr(unsafe.Pointer(&args))
	libfuzzerCall((*byte)(C.do_rpmalloc_heap_alloc), ptr, 0)
	return args.ptr
}

type heap_alloc_cap_t struct {
	heap uintptr
	size uintptr
	ptr  uintptr
	cap  uintptr
}

// AllocCap Allocate a memory block of at least the given size using the given heap.
func (h *Heap) AllocCap(size uintptr) (uintptr, uintptr) {
	args := heap_alloc_cap_t{heap: uintptr(unsafe.Pointer(h)), size: size}
	ptr := uintptr(unsafe.Pointer(&args))
	libfuzzerCall((*byte)(C.do_rpmalloc_heap_alloc), ptr, 0)
	return args.ptr, args.cap
}

type heap_calloc_t struct {
	heap uintptr
	num  uintptr
	size uintptr
	ptr  uintptr
}

// Calloc Allocate a memory block of at least the given size using the given heap and zero initialize it.
func (h *Heap) Calloc(num, size uintptr) uintptr {
	args := heap_calloc_t{heap: uintptr(unsafe.Pointer(h)), num: num, size: size}
	ptr := uintptr(unsafe.Pointer(&args))
	libfuzzerCall((*byte)(C.do_rpmalloc_heap_calloc), ptr, 0)
	return args.ptr
}

type heap_calloc_cap_t struct {
	heap uintptr
	num  uintptr
	size uintptr
	ptr  uintptr
	cap  uintptr
}

// Calloc Allocate a memory block of at least the given size using the given heap and zero initialize it.
func (h *Heap) CallocCap(num, size uintptr) (uintptr, uintptr) {
	args := heap_calloc_cap_t{heap: uintptr(unsafe.Pointer(h)), num: num, size: size}
	ptr := uintptr(unsafe.Pointer(&args))
	libfuzzerCall((*byte)(C.do_rpmalloc_heap_calloc_cap), ptr, 0)
	return args.ptr, args.cap
}

type heap_realloc_t struct {
	heap   uintptr
	ptr    uintptr
	size   uintptr
	newptr uintptr
	flags  int32
}

// Realloc Reallocate the given block to at least the given size. The memory block MUST be allocated
// by the same heap given to this function.
func (h *Heap) Realloc(ptr, size uintptr, flags int32) uintptr {
	args := heap_realloc_t{heap: uintptr(unsafe.Pointer(h)), ptr: ptr, size: size, flags: flags}
	p := uintptr(unsafe.Pointer(&args))
	libfuzzerCall((*byte)(C.do_rpmalloc_heap_realloc), p, 0)
	return args.newptr
}

type heap_realloc_cap_t struct {
	heap   uintptr
	ptr    uintptr
	size   uintptr
	newptr uintptr
	cap    uintptr
	flags  int32
}

// ReallocCap Reallocate the given block to at least the given size. The memory block MUST be allocated
// by the same heap given to this function.
func (h *Heap) ReallocCap(ptr, size uintptr, flags int32) (uintptr, uintptr) {
	args := heap_realloc_cap_t{heap: uintptr(unsafe.Pointer(h)), ptr: ptr, size: size, flags: flags}
	p := uintptr(unsafe.Pointer(&args))
	libfuzzerCall((*byte)(C.do_rpmalloc_heap_realloc_cap), p, 0)
	return args.newptr, args.cap
}

type heap_free_t struct {
	heap uintptr
	ptr  uintptr
}

// Free the given memory block from the given heap. The memory block MUST be allocated
// by the same heap given to this function.
func (h *Heap) Free(ptr uintptr) {
	args := heap_free_t{heap: uintptr(unsafe.Pointer(h)), ptr: ptr}
	p := uintptr(unsafe.Pointer(&args))
	libfuzzerCall((*byte)(C.do_rpmalloc_heap_free), p, 0)
}

type heap_free_all_t struct {
	heap uintptr
}

// FreeAll memory allocated by the heap
func (h *Heap) FreeAll() {
	args := heap_free_all_t{heap: uintptr(unsafe.Pointer(h))}
	p := uintptr(unsafe.Pointer(&args))
	libfuzzerCall((*byte)(C.do_rpmalloc_heap_free_all), p, 0)
}
