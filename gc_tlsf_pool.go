//go:build !tinygo.wasm
// +build !tinygo.wasm

package runtime

import (
	"sync"
	"unsafe"
)

type Pool struct {
	root         *Root
	heapStart    uintptr
	heapEnd      uintptr
	segments     [][]byte
	pages        int
	heapSize     int64
	allocSize    int64
	maxAllocSize int64
	freeSize     int64
	allocs       int32
	mu           sync.Mutex
}

func NewPool(pages int) *Pool {
	if pages <= 0 {
		pages = 1
	}
	page := make([]byte, pages*wasmPageSize)
	p := &Pool{
		heapStart: uintptr(unsafe.Pointer(&page[0])),
		heapSize:  int64(len(page)),
		segments:  [][]byte{page},
		pages:     pages,
	}

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

	var (
		memStart = rootOffset + tlsf_ROOT_SIZE
	)
	addMemory(p, memStart, p.heapStart+uintptr(len(page)))
	return p
}

func (p *Pool) Grow(pages int) (uintptr, uintptr) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if pages <= 0 {
		pages = 1
	}
	buf := make([]byte, pages*wasmPageSize)
	p.pages += pages
	p.segments = append(p.segments, buf)

	start, end := uintptr(unsafe.Pointer(&buf[0])), uintptr(unsafe.Pointer(&buf[0]))+uintptr(len(buf))
	p.heapSize += int64(len(buf))
	p.heapEnd = end
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
