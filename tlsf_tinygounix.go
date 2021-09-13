//go:build tinygo && (darwin || linux || freebsd) && !baremetal && !wasi && !tinygo.wasm
// +build tinygo
// +build darwin linux freebsd
// +build !baremetal
// +build !wasi
// +build !tinygo.wasm

package mem

import "unsafe"

const wasmPageSize = 64 * 1024

func init() {
	initAllocator(0)
}

func Alloc(size uintptr) unsafe.Pointer {
	return allocator.Alloc(size)
}

func Realloc(ptr unsafe.Pointer, size uintptr) unsafe.Pointer {
	return allocator.Realloc(ptr, size)
}

func Free(ptr unsafe.Pointer) {
	allocator.Free(ptr)
}

//go:export initAllocator
func initAllocator(heapStart uintptr) {
	if allocator != nil {
		return
	}
	allocator = newTLSF(3)
}

//go:export extalloc
func extalloc(size uintptr) unsafe.Pointer {
	return allocator.Alloc(size)
}

//go:export extfree
func extfree(ptr unsafe.Pointer) {
	allocator.Free(ptr)
}

func getAllocator() *TLSF {
	return allocator
}

//go:export malloc
func malloc(size uintptr) unsafe.Pointer

////go:export memcpy
////export memcpy
//func Copy(dst, src unsafe.Pointer, size uintptr)
//
////go:export memzero
////export memzero
//func Zero(ptr unsafe.Pointer, size uintptr)

func Copy(dst, src unsafe.Pointer, n uintptr) {
	memcpy(dst, src, n)
}

//export memcpy
func memcpy(dst, src unsafe.Pointer, n uintptr)

//func memcpy0(dst, src unsafe.Pointer, n uintptr) {
//	dstB := *(*[]byte)(unsafe.Pointer(&struct{ Data, Len, Cap uintptr }{
//		Data: uintptr(dst),
//		Len:  n,
//		Cap:  n,
//	}))
//	srcB := *(*[]byte)(unsafe.Pointer(&struct{ Data, Len, Cap uintptr }{
//		Data: uintptr(src),
//		Len:  n,
//		Cap:  n,
//	}))
//	copy(dstB, srcB)
//}

func Zero(ptr unsafe.Pointer, size uintptr) {
	memzero(ptr, size)
}

//export memzero
func memzero(ptr unsafe.Pointer, size uintptr)

//func memzero0(ptr unsafe.Pointer, size uintptr) {
//	b := *(*[]byte)(unsafe.Pointer(&struct{ Data, Len, Cap uintptr }{
//		Data: uintptr(ptr),
//		Len:  size,
//		Cap:  size,
//	}))
//	//println("memzero", uint(uintptr(unsafe.Pointer(&b[0]))), len(b))
//	switch {
//	case size < 8:
//		for i := 0; i < len(b); i++ {
//			b[i] = 0
//		}
//	case size == 8:
//		*(*uint64)(unsafe.Pointer(&b[0])) = 0
//	default:
//		var i = 0
//		for ; i <= len(b)-8; i += 8 {
//			*(*uint64)(unsafe.Pointer(&b[i])) = 0
//		}
//		for ; i < len(b); i++ {
//			b[i] = 0
//		}
//	}
//}

func GrowByDouble() Grow {
	return func(pagesBefore, pagesNeeded int32, minSize uintptr) (pagesAdded int32, start, end uintptr) {
		if pagesBefore > pagesNeeded {
			pagesAdded = pagesBefore
		} else {
			pagesAdded = pagesNeeded
		}
		ptr := malloc(uintptr(pagesAdded) * _TLSFPageSize)
		if ptr == nil {
			pagesAdded = pagesNeeded
			ptr = malloc(uintptr(pagesAdded) * _TLSFPageSize)
			if ptr == nil {
				return 0, 0, 0
			}
		}
		start = uintptr(ptr)
		end = start + (uintptr(pagesAdded) * _TLSFPageSize)
		return
	}
}

func GrowBy(pages int32) Grow {
	return func(pagesBefore, pagesNeeded int32, minSize uintptr) (pagesAdded int32, start, end uintptr) {
		if pages > pagesNeeded {
			pagesAdded = pages
		} else {
			pagesAdded = pagesNeeded
		}
		ptr := malloc(uintptr(pagesAdded) * _TLSFPageSize)
		if ptr == nil {
			pagesAdded = pagesNeeded
			ptr = malloc(uintptr(pagesAdded) * _TLSFPageSize)
			if ptr == nil {
				return 0, 0, 0
			}
		}
		start = uintptr(ptr)
		end = start + (uintptr(pagesAdded) * _TLSFPageSize)
		return
	}
}

func GrowMin() Grow {
	return func(pagesBefore, pagesNeeded int32, minSize uintptr) (int32, uintptr, uintptr) {
		ptr := malloc(uintptr(pagesNeeded) * _TLSFPageSize)
		if ptr == nil {
			return 0, 0, 0
		}
		return pagesNeeded, uintptr(ptr), uintptr(ptr) + (uintptr(pagesNeeded) * _TLSFPageSize)
	}
}

func newTLSF(pages int32) *TLSF {
	if pages <= 0 {
		pages = 1
	}
	size := uintptr(pages * wasmPageSize)
	segment := uintptr(malloc(size))
	return bootstrap(segment, segment+size, pages)
}

func (a *TLSF) Grow(pages int32) (uintptr, uintptr) {
	if pages <= 0 {
		pages = 1
	}
	// Allocate new segment
	segment := uintptr(malloc(uintptr(pages * wasmPageSize)))
	if segment == 0 {
		return 0, 0
	}
	a.Pages += pages
	start, end := segment, segment+uintptr(pages*wasmPageSize)
	a.HeapEnd = end
	return start, end
}

//func markStack() {}
//
//func markGlobals() {}
//
//func markScheduler() {}
