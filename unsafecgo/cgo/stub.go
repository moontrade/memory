package cgo

/*
#include <stdio.h>
#include <time.h>
#include <unistd.h>

void do_usleep(size_t arg0, size_t arg1) {
	nanosleep((const struct timespec[]){{0, 500000000L}}, NULL);
}
void unsafecgo_stub() {}

typedef void unsafecgo_trampoline_handler(size_t arg0, size_t arg1);

void unsafecgo_cgo_call(size_t fn, size_t arg0, size_t arg1) {
	((unsafecgo_trampoline_handler*)fn)(arg0, arg1);
}
*/
import "C"
import "unsafe"

var (
	Stub   = C.unsafecgo_stub
	Usleep = C.do_usleep
)

func CGO() {
	C.unsafecgo_stub()
}

func NonBlocking(fn *byte, arg0, arg1 uintptr) {
	Blocking(fn, arg0, arg1)
	//libcCall(unsafe.Pointer(fn), unsafe.Pointer(arg0))
}

func Blocking(fn *byte, arg0, arg1 uintptr) {
	C.unsafecgo_cgo_call((C.size_t)(uintptr(unsafe.Pointer(fn))), (C.size_t)(arg0), (C.size_t)(arg1))
}

////go:linkname libcCall runtime.libcCall
//func libcCall(fn, arg unsafe.Pointer) int32

func DoUsleep(useconds int64) {
	C.do_usleep((C.size_t)(useconds), (C.size_t)(0))
}
