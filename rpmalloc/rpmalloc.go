package rpmalloc

/*
#include "rpmalloc.h"
#include <stdlib.h>
#include <string.h>
#include <pthread.h>
#include <stdio.h>
#include <inttypes.h>

//void pthread_cleanup_push(void (*routine)(void *), void *arg);
//void pthread_cleanup_pop(int execute);

void on_thread_cleanup(void* data) {
}

void hook() {
	//pthread_cleanup_push(on_thread_cleanup, NULL);
	printf("Hi I'm C!\n");
}

void stub() {
	//pthread_cleanup_push(on_thread_cleanup, NULL);
	//printf("Hi I'm C!\n");
}

typedef struct {
	uintptr_t* args;
	uintptr_t retval;
} argset_t;

typedef struct {
	uintptr_t size;
	uintptr_t retval;
} malloc_t;

void alloc(malloc_t* args) {
	//size_t sz = (size_t)*args->args;
	//printf("sz: %i\n", (int)sz);
	args->retval = (uintptr_t)rpmalloc((size_t)args->size);
	//rpfree((void*)args->retval);
}

void alloc_32() {
	//size_t sz = (size_t)*args->args;
	//printf("sz: %i\n", (int)sz);
	rpfree(rpmalloc(32));
	//args->retval = (uintptr_t)rpmalloc((size_t)*args->args);
}

void do_free(void* ptr) {
	rpfree(ptr);
}

void run_rpallocs(size_t size, int n) {
	for (int i = 0; i < n; i++) {
		rpfree(rpmalloc(size));
	}
}

void run_rpalloczero(size_t size, int n) {
	for (int i = 0; i < n; i++) {
		void* ptr = rpmalloc(size);
		memset(ptr, 0, size);
		rpfree(ptr);
	}
}

void* ml(size_t size) {
	return malloc(size);
}

void run_mallocs(size_t size, int n) {
	for (int i = 0; i < n; i++) {
		void* ptr = ml(size);
		free(ptr);
	}
}
*/
import "C"
import (
	"github.com/moontrade/memory/alloc"
	"unsafe"
)

func init() {
	C.rpmalloc_initialize()
}

//func Hook() {
//	C.hook()
//}
//
//func HookDirect() {
//	asmcgocall(unsafe.Pointer(C.hook), nil)
//}
//
//func Stub() {
//	C.stub()
//}
//
//func StubDirect() {
//	asmcgocall(unsafe.Pointer(C.stub), nil)
//}
//
//func AllocDirect(size int) uintptr {
//	a := &args[0]
//	a.size = uintptr(size)
//	a.retval = 0
//	asmcgocall(unsafe.Pointer(C.alloc), unsafe.Pointer(a))
//	return a.retval
//}
//
//var (
//	args [255]malloc_t
//)
//
//func AllocDirect32() uintptr {
//	sz := C.size_t(32)
//	as := argset{args: unsafe.Pointer(&sz)}
//	asmcgocall(unsafe.Pointer(C.alloc_32), nil)
//	_ = as
//	return 0
//	//return as.retval
//}
//
//func FreeDirect(ptr uintptr) {
//	asmcgocall(unsafe.Pointer(C.do_free), unsafe.Pointer(ptr))
//}

func InitThread() {
	C.rpmalloc_thread_initialize()
}

func runAllocs(size uintptr, n int32) {
	C.run_rpallocs((C.size_t)(size), (C.int)(n))
}

func runAllocZeroed(size uintptr, n int32) {
	C.run_rpalloczero((C.size_t)(size), (C.int)(n))
}

func runMallocs(size uintptr, n int32) {
	C.run_mallocs((C.size_t)(size), (C.int)(n))
}

func Alloc(size uintptr) alloc.Pointer {
	return alloc.Pointer(unsafe.Pointer(C.rpmalloc((C.size_t)(size))))
}

// argset matches runtime/cgo/linux_syscall.c:argset_t
type argset struct {
	args   unsafe.Pointer
	retval uintptr
}

////go:linkname asmcgocall runtime.asmcgocall
//func asmcgocall(fn, arg unsafe.Pointer) int32
