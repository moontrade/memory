// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !libfuzzer && (amd64 || arm64) && (linux || darwin)
// +build !libfuzzer
// +build amd64 arm64
// +build linux darwin

#include "headers.h"
#include "textflag.h"

// Based on race_arm64.s; see commentary there.

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

// func runtime·libfuzzerCall(fn, arg0, arg1 uintptr)
// Calls C function fn from libFuzzer and passes 2 arguments to it.
TEXT ·NonBlocking(SB), NOSPLIT, $0-24
	MOVD	fn+0(FP), R9
	MOVD	arg0+8(FP), R0
	MOVD	arg1+16(FP), R1

	MOVD	g_m(g), R10

	// Switch to g0 stack.
	MOVD	RSP, R19	// callee-saved, preserved across the CALL
	MOVD	m_g0(R10), R11
	CMP	R11, g
	BEQ	call	// already on g0
	MOVD	(g_sched+gobuf_sp)(R11), R12
	MOVD	R12, RSP
call:
	BL	R9
	MOVD	R19, RSP
	RET
