package rp

/*
#include "rpmalloc.h"
#include <stdlib.h>
#include <string.h>
#include <pthread.h>

//void pthread_cleanup_push(void (*routine)(void *), void *arg);
//void pthread_cleanup_pop(int execute);

void on_thread_cleanup(void* data) {
}

void hook() {
	//pthread_cleanup_push(on_thread_cleanup, NULL);
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
	mem "github.com/moontrade/memory"
	"unsafe"
)

func init() {
	C.rpmalloc_initialize()
}

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

func Alloc(size uintptr) mem.Pointer {
	return mem.Pointer(unsafe.Pointer(C.rpmalloc((C.size_t)(size))))
}

func Free(ptr mem.Pointer) {
	C.rpfree(unsafe.Pointer(ptr))
}
