//go:build !tinygo.wasm

package runtime

import (
	"reflect"
	"unsafe"
)

var (
	globalsStart uintptr = 0
	globalsEnd   uintptr = 0
)

const (
	wasmPageSize = 64 * 1024
)

func maybeFreeBlock(p *Pool, block *tlsfBlock) {
	freeBlock(p, block)
}

func tlsfPrintInfo() {
	//println("heapStart		", int64(heapStart))
	//println("heapEnd			", int64(heapEnd))
	println("ALIGNOF_U32		", int64(tlsf_ALIGN_U32))
	println("ALIGNOF_USIZE	", int64(tlsf_ALIGN_SIZE_LOG2))
	println("U32_MAX			", ^uint32(0))
	println("PTR_MAX			", ^uintptr(0))
	println("AL_BITS			", int64(tlsf_AL_BITS))
	println("AL_SIZE			", int64(tlsf_AL_SIZE))
	println("AL_MASK			", int64(tlsf_AL_MASK))
	println("BLOCK_OVERHEAD	", int64(tlsf_BLOCK_OVERHEAD))
	println("BLOCK_MAXSIZE	", int64(tlsf_BLOCK_MAXSIZE))
	println("SL_BITS			", int64(tlsf_SL_BITS))
	println("SL_SIZE			", int64(tlsf_SL_SIZE))
	println("SB_BITS			", int64(tlsf_SB_BITS))
	println("SB_SIZE			", int64(tlsf_SB_SIZE))
	println("FL_BITS			", int64(tlsf_FL_BITS))
	println("FREE			", int64(tlsf_FREE))
	println("LEFTFREE		", int64(tlsf_LEFTFREE))
	println("TAGS_MASK		", int64(tlsf_TAGS_MASK))
	println("BLOCK_MINSIZE	", int64(tlsf_BLOCK_MINSIZE))
	println("SL_START		", int64(tlsf_SL_START))
	println("SL_END			", int64(tlsf_SL_END))
	println("HL_START		", int64(tlsf_HL_START))
	println("HL_END			", int64(tlsf_HL_END))
	println("ROOT_SIZE		", int64(tlsf_ROOT_SIZE))
}

func memcpy(dst, src unsafe.Pointer, n uintptr) {
	dstB := *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(dst),
		Len:  int(n),
		Cap:  int(n),
	}))
	srcB := *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(src),
		Len:  int(n),
		Cap:  int(n),
	}))
	copy(dstB, srcB)
}

func memzero(ptr unsafe.Pointer, size uintptr) {
	b := *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(ptr),
		Len:  int(size),
		Cap:  int(size),
	}))
	switch {
	case size < 8:
		for i := 0; i < len(b); i++ {
			b[i] = 0
		}
	case size == 8:
		*(*uint64)(unsafe.Pointer(&b[0])) = 0
	default:
		var i = 0
		for ; i <= len(b)-8; i += 8 {
			*(*uint64)(unsafe.Pointer(&b[i])) = 0
		}
		for ; i < len(b); i++ {
			b[i] = 0
		}
	}
}
