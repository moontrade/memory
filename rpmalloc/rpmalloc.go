package rpmalloc

/*
#cgo LDFLAGS: -lc++
#include "rpmalloc.h"
#include "malloc.h"
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

func init() {
	C.rpmalloc_initialize()
	C.malloc(12)
	//C.do_hook_malloc_a()
}

type Heap C.rpmalloc_heap_t
