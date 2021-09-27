//go:build !cgo_safe && !libfuzzer
// +build !cgo_safe,!libfuzzer

package rax

/*
#include "rax.h"
#include <stdlib.h>

typedef struct {
	size_t ptr;
} rax_new_t;

void do_rax_new(rax_new_t* args) {
	args->ptr = (size_t)raxNew();
}

void do_rax_free(rax *rax) {
	raxFree(rax);
}

typedef struct {
	size_t ptr;
	size_t size;
} rax_size_t;

void do_rax_size(rax_size_t* args) {
	args->size = (size_t)raxSize((rax*)(void*)args->ptr);
}

void do_rax_show(rax* args) {
	raxShow(args);
}

typedef struct {
	size_t rax;
	size_t key;
	size_t len;
	size_t data;
	size_t old;
	size_t result;
} rax_insert_t;

void do_rax_insert(rax_insert_t* args) {
	void* old = NULL;
	args->result = (size_t)raxInsert((rax*)(void*)args->rax, (unsigned char*)args->key, args->len, (void*)args->data, &old);
	args->old = (size_t)old;
}

*/
import "C"
import (
	"github.com/moontrade/memory"
	"unsafe"
)

//go:linkname libcCall runtime.libcCall
func libcCall(fn, arg unsafe.Pointer)

type rax_new_t struct {
	ptr uintptr
}

func New() *Rax {
	args := rax_new_t{}
	ptr := uintptr(unsafe.Pointer(&args))
	libcCall(unsafe.Pointer(C.do_rax_new), unsafe.Pointer(ptr))
	return (*Rax)(unsafe.Pointer(args.ptr))
}

func (r *Rax) Free() {
	ptr := uintptr(unsafe.Pointer(r))
	libcCall(unsafe.Pointer(C.do_rax_free), unsafe.Pointer(ptr))
}

type rax_size_t struct {
	ptr  uintptr
	size uintptr
}

func (r *Rax) Size() int {
	args := rax_size_t{ptr: uintptr(unsafe.Pointer(r))}
	ptr := uintptr(unsafe.Pointer(&args))
	libcCall(unsafe.Pointer(C.do_rax_size), unsafe.Pointer(ptr))
	return int(args.size)
}

func (r *Rax) Print() {
	libcCall(unsafe.Pointer(C.do_rax_show), unsafe.Pointer(r))
}

type rax_insert_t struct {
	rax    uintptr
	s      uintptr
	len    uintptr
	data   uintptr
	old    uintptr
	result uintptr
}

func (r *Rax) Insert(key memory.Pointer, size int, data memory.Pointer) (int, memory.Pointer) {
	args := rax_insert_t{
		rax:  uintptr(unsafe.Pointer(r)),
		s:    uintptr(key),
		len:  uintptr(size),
		data: uintptr(data),
	}
	ptr := uintptr(unsafe.Pointer(&args))
	libcCall(unsafe.Pointer(C.do_rax_insert), unsafe.Pointer(ptr))
	return int(args.result), memory.Pointer(args.old)
}

func (r *Rax) InsertBytes(key memory.Bytes, data memory.Pointer) (int, memory.Pointer) {
	args := rax_insert_t{
		rax:  uintptr(unsafe.Pointer(r)),
		s:    uintptr(key.Pointer),
		len:  uintptr(key.Len()),
		data: uintptr(data),
	}
	ptr := uintptr(unsafe.Pointer(&args))
	libcCall(unsafe.Pointer(C.do_rax_insert), unsafe.Pointer(ptr))
	return int(args.result), memory.Pointer(args.old)
}
