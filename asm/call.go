package asm

/*
#include <stdlib.h>
#include <stdio.h>

void do_call() {
	fprintf(stderr, "hello\n");
}
*/
import "C"
import "github.com/moontrade/memory/call"

func Call() {
	C.do_call()
}

func Call0() {
	call.InvokeC((*byte)(C.do_call), 0, 0)
}
