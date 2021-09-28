package rax

/*
#cgo darwin,amd64 CFLAGS: -I${SRCDIR}/src
#cgo darwin,amd64 LDFLAGS: -Wl,-rpath,${SRCDIR}/lib/darwin_amd64 -L${SRCDIR}/lib/darwin_amd64
#cgo darwin,amd64 LDFLAGS: -lrax -ldl -lc -lm
#cgo linux,amd64 CFLAGS: -I${SRCDIR}/src
#cgo linux,amd64 LDFLAGS: -Wl,-rpath,${SRCDIR}/lib/linux_amd64 -L${SRCDIR}/lib/linux_amd64
#cgo linux,amd64 LDFLAGS: -lrax -ldl -lc -lm
#cgo linux,amd64 CFLAGS: -D_GNU_SOURCE
#include "rax.h"
*/
import "C"

type (
	Rax      C.rax
	Node     C.raxNode
	Stack    C.raxStack
	Iterator C.raxIterator
)
