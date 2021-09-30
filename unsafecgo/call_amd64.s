// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !libfuzzer && (amd64 || arm64) && (linux || darwin)
// +build !libfuzzer
// +build amd64 arm64
// +build linux darwin

//#include "headers.h"
#include "go_asm.h"
#include "go_tls.h"
#include "textflag.h"

// Based on race_amd64.s; see commentary there.

#ifdef GOOS_windows
#define RARG0 CX
#define RARG1 DX
#else
#define RARG0 DI
#define RARG1 SI
#endif

// WARNING!!!
// Go doesn't allow packages outside of runtime to include "go_asm.h" so the below
// defines required were pulled from generated "go_asm.h" by running make in the
// cmd sub-directory of this package. Navigate into the build (WORK) directory and
// look for a go_asm.h file that's big (>10kb). The below defines will be in there.
// The below defines have been observed to be the same across both linux and darwin
// given it appears to be CPU arch based (amd64) only. The below defines are also
// the same for arm64.
//
// Safer to add build tag "libfuzzer" to hook into the auto-generated "go_asm.h".
// However, it's about ~1ns slower per call because of linking overhead somehow.

#define g_m 48
#define g_sched 56
#define gobuf_sp 0
#define m_g0 0

// void runtime·libfuzzerCall(fn, arg0, arg1 uintptr)
// Calls C function fn from libFuzzer and passes 2 arguments to it.
TEXT ·NonBlocking(SB), NOSPLIT, $0-24
	MOVQ	fn+0(FP), AX
	MOVQ	arg0+8(FP), RARG0
	MOVQ	arg1+16(FP), RARG1

	get_tls(R12)
	MOVQ	g(R12), R14
	MOVQ	g_m(R14), R13

	// Switch to g0 stack.
	MOVQ	SP, R12		// callee-saved, preserved across the CALL
	MOVQ	m_g0(R13), R10
	CMPQ	R10, R14
	JE	call	// already on g0
	MOVQ	(g_sched+gobuf_sp)(R10), SP
call:
	ANDQ	$~15, SP	// alignment for gcc ABI
	CALL	AX
	MOVQ	R12, SP
	// Back to Go world, set special registers.
    // The g register (R14) is preserved in C.
    //XORPS	X15, X15
	RET
