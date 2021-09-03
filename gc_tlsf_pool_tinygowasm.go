//go:build (tinygo.wasm && gc.extalloc) || wasm

package runtime

import (
	"unsafe"
)

type Pool struct {
	root         *Root
	heapStart    uintptr
	heapEnd      uintptr
	heapSize     int64
	allocSize    int64
	maxAllocSize int64
	freeSize     int64
	allocs       int32
	pages        int
}

func NewPool(pages int) *Pool {
	if pages <= 0 {
		pages = 1
	}
	p := &Pool{}

	before := wasm_memory_size(0)
	wasm_memory_grow(0, int32(pages))
	after := wasm_memory_size(0)

	p.pages = pages

	p.heapStart, p.heapEnd = uintptr(before)*uintptr(wasmPageSize), uintptr(after)*uintptr(wasmPageSize)
	heapEnd = p.heapEnd
	p.heapSize = int64(p.heapEnd - p.heapStart)

	rootOffset := (p.heapStart + tlsf_AL_MASK) & ^tlsf_AL_MASK
	p.root = (*Root)(unsafe.Pointer(rootOffset))
	p.root.flMap = 0
	setTail(p.root, nil)
	for fl := uintptr(0); fl < uintptr(tlsf_FL_BITS); fl++ {
		setSL(p.root, fl, 0)
		for sl := uint32(0); sl < tlsf_SL_SIZE; sl++ {
			setHead(p.root, fl, sl, nil)
		}
	}

	addMemory(p, rootOffset+tlsf_ROOT_SIZE, uintptr(after)*uintptr(wasmPageSize))
	return p
}

func (p *Pool) Grow(pages int) (uintptr, uintptr) {
	if pages <= 0 {
		pages = 1
	}
	before := wasm_memory_size(0)
	p.pages += pages
	wasm_memory_grow(0, int32(pages))
	after := wasm_memory_size(0)

	start := uintptr(before * wasmPageSize)
	end := uintptr(after * wasmPageSize)

	p.heapSize += int64(end - start)
	p.heapEnd = end
	heapEnd = end
	return start, end
}

func (p *Pool) Alloc(size uintptr) unsafe.Pointer {
	return unsafe.Pointer(tlsfalloc(p, size))
}

func (p *Pool) Realloc(ptr unsafe.Pointer, size uintptr) unsafe.Pointer {
	return unsafe.Pointer(tlsfrealloc(p, uintptr(ptr), size))
}

func (p *Pool) Free(ptr unsafe.Pointer) {
	tlsffree(p, uintptr(ptr))
}

func maybeFreeBlock(p *Pool, block *tlsfBlock) {
	if uintptr(unsafe.Pointer(block)) >= p.heapStart {
		freeBlock(p, block)
	}
}
