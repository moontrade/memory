//go:build !tinygo.wasm

package runtime

import (
	"reflect"
	"sync"
	"unsafe"
)

var (
	globalsStart uintptr = 0
	globalsEnd   uintptr = 0
)

const (
	wasmPageSize = 64 * 1024
)

var _allocs = make(map[uintptr][]byte)
var _allocsMu sync.Mutex

func malloc(size uintptr) unsafe.Pointer {
	b := make([]byte, size)
	p := unsafe.Pointer(&b[0])
	_allocsMu.Lock()
	defer _allocsMu.Unlock()
	_allocs[uintptr(p)] = b
	return p
}

func realloc(ptr unsafe.Pointer, size uintptr) unsafe.Pointer {
	_allocsMu.Lock()
	old := _allocs[uintptr(ptr)]
	_allocsMu.Unlock()
	if old == nil {
		return nil
	}

	b := make([]byte, size)
	copy(b, old)
	np := unsafe.Pointer(&b[0])
	_allocsMu.Lock()
	_allocs[uintptr(np)] = b
	_allocsMu.Unlock()
	return np
}

func free(ptr unsafe.Pointer) {
	_allocsMu.Lock()
	defer _allocsMu.Unlock()
	b := _allocs[uintptr(ptr)]
	if b != nil {
		delete(_allocs, uintptr(ptr))
	}
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
