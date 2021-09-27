//go:build libfuzzer
// +build libfuzzer

package rax

/*
#include "rax.h"
#include <stdlib.h>

typedef struct {
	size_t ptr;
} rax_new_t;

void do_rax_new(size_t arg0, size_t arg1) {
	rax_new_t* args = (rax_new_t*)(void*)arg0;
	args->ptr = (size_t)raxNew();
}

void do_rax_free(rax *rax) {
	raxFree(rax);
}

typedef struct {
	size_t ptr;
	size_t size;
} rax_size_t;

void do_rax_size(size_t arg0, size_t arg1) {
	rax_size_t* args = (rax_size_t*)(void*)arg0;
	args->size = (size_t)raxSize((rax*)(void*)args->ptr);
}

void do_rax_show(size_t arg0, size_t arg1) {
	raxShow((rax*)(void*)arg0);
}

typedef struct {
	size_t rax;
	size_t key;
	size_t len;
	size_t data;
	size_t old;
	size_t result;
} rax_insert_t;

void do_rax_insert(size_t arg0, size_t arg1) {
	rax_insert_t* args = (rax_insert_t*)(void*)arg0;
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

//go:linkname libfuzzerCall runtime.libfuzzerCall
func libfuzzerCall(fn *byte, arg0, arg1 uintptr)

type rax_new_t struct {
	ptr uintptr
}

func New() *Rax {
	args := rax_new_t{}
	ptr := uintptr(unsafe.Pointer(&args))
	libfuzzerCall((*byte)(C.do_rax_new), ptr, 0)
	return (*Rax)(unsafe.Pointer(args.ptr))
}

func (r *Rax) Free() {
	ptr := uintptr(unsafe.Pointer(r))
	libfuzzerCall((*byte)(C.do_rax_free), ptr, 0)
}

type rax_size_t struct {
	ptr  uintptr
	size uintptr
}

func (r *Rax) Size() int {
	args := rax_size_t{ptr: uintptr(unsafe.Pointer(r))}
	ptr := uintptr(unsafe.Pointer(&args))
	libfuzzerCall((*byte)(C.do_rax_size), ptr, 0)
	return int(args.size)
}

func (r *Rax) Print() {
	libfuzzerCall((*byte)(C.do_rax_show), uintptr(unsafe.Pointer(r)), 0)
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
	libfuzzerCall((*byte)(C.do_rax_insert), ptr, 0)
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
	libfuzzerCall((*byte)(C.do_rax_insert), ptr, 0)
	return int(args.result), memory.Pointer(args.old)
}
