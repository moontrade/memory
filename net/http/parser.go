package http

/*
#include "picohttpparser.h"
#include <stdlib.h>

typedef struct {
	uintptr_t buf;
	uintptr_t len;
	uintptr_t method;
	size_t method_len;
	uintptr_t path;
	size_t path_len;
	uintptr_t headers;
	size_t num_headers;
	size_t max_num_headers;
	uintptr_t last_len;
	int minor_version;
	int result;
} phr_parse_request_t;

void do_phr_parse_request(size_t arg0, size_t arg1) {
	phr_parse_request_t* args = (phr_parse_request_t*)(void*)arg0;

	const char *method;
    const char *path;
	args->result = phr_parse_request(
		(const char*)args->buf,
		args->len,
		&method,
		&args->method_len,
		&path,
		&args->path_len,
		&args->minor_version,
		(struct phr_header*)args->headers,
		&args->num_headers,
		args->last_len
	);
	args->method = (uintptr_t)method;
	args->path = (uintptr_t)path;
}


*/
import "C"
import (
	"github.com/moontrade/memory"
	"github.com/moontrade/memory/unsafecgo"
	"reflect"
	"unsafe"
)

type Request struct {
	input        uintptr
	inputLen     uintptr
	method       uintptr
	methodLen    uintptr
	path         uintptr
	pathLen      uintptr
	headers      uintptr
	numHeaders   uintptr
	maxHeaders   uintptr
	lastLen      uintptr
	minorVersion int32
	result       int32
}

func AllocRequest(headers int) *Request {
	if headers < 10 {
		headers = 10
	}
	if headers > 4096 {
		headers = 4096
	}
	p := memory.AllocZeroed(unsafe.Sizeof(Request{}) + (unsafe.Sizeof(Header{}) * uintptr(headers)))
	r := (*Request)(unsafe.Pointer(p))
	r.headers = uintptr(p.Add(int(unsafe.Sizeof(Request{}))))
	r.numHeaders = uintptr(headers)
	r.maxHeaders = uintptr(headers)
	return r
}

func (p *Request) Free() {
	memory.Free(memory.Pointer(unsafe.Pointer(p)))
}

func (p *Request) Method() string {
	return *(*string)(unsafe.Pointer(&reflect.StringHeader{
		Data: p.method,
		Len:  int(p.methodLen),
	}))
}
func (p *Request) Path() string {
	return *(*string)(unsafe.Pointer(&reflect.StringHeader{
		Data: p.path,
		Len:  int(p.pathLen),
	}))
}

func (p *Request) NumHeaders() int {
	return int(p.numHeaders)
}

func (p *Request) Header(index int) *Header {
	return (*Header)(unsafe.Pointer(uintptr(unsafe.Pointer(p)) + unsafe.Sizeof(Request{}) + (unsafe.Sizeof(Header{}) * uintptr(index))))
}

type Header struct {
	Name  memory.FatPointer
	Value memory.FatPointer
}

//func printSizeOfs() {
//	type SizeT C.size_t
//	type UintptrT C.uintptr_t
//	type Int C.int
//	println("sizeof(size_t)", uint(unsafe.Sizeof(SizeT(0))),
//		"sizeof(uintptr_t)", uint(unsafe.Sizeof(UintptrT(0))),
//		"sizeof(struct phr_header)", uint(unsafe.Sizeof(Header{})),
//		"sizeof(int)", uint(unsafe.Sizeof(Int(0))))
//}

func ParseRequest(args *Request) int {
	//args := Request{}
	args.numHeaders = args.maxHeaders
	ptr := uintptr(unsafe.Pointer(args))
	unsafecgo.Call((*byte)(C.do_phr_parse_request), ptr, 0)
	return int(args.result)
}
