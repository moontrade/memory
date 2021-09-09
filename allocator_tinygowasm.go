//go:build (tinygo.wasm || wasi) && !gc.conservative
// +build tinygo.wasm wasi
// +build !gc.conservative

package mem

const wasmPageSize = 64 * 1024

//export llvm.wasm.memory.size.i32
func wasm_memory_size(index int32) int32

//export llvm.wasm.memory.grow.i32
func wasm_memory_grow(index, pages int32) int32

// GrowByDouble will double the heap on each grow
func GrowByDouble() Grow {
	return func(pagesBefore, pagesNeeded int32, minSize uintptr) (pagesAdded int32, start, end uintptr) {
		if pagesBefore > pagesNeeded {
			pagesAdded = pagesBefore
		} else {
			pagesAdded = pagesNeeded
		}
		start, end = growBy(pagesAdded)
		if start == 0 {
			pagesAdded = pagesNeeded
			start, end = growBy(pagesAdded)
			if start == 0 {
				return 0, 0, 0
			}
		}
		return
	}
}

// GrowBy will grow by the number of pages specified or by the minimum needed, whichever is greater.
func GrowBy(pages int32) Grow {
	return func(pagesBefore, pagesNeeded int32, minSize uintptr) (pagesAdded int32, start, end uintptr) {
		if pages > pagesNeeded {
			pagesAdded = pages
		} else {
			pagesAdded = pagesNeeded
		}
		start, end = growBy(pagesAdded)
		if start == 0 {
			pagesAdded = pagesNeeded
			start, end = growBy(pagesAdded)
			if start == 0 {
				return 0, 0, 0
			}
		}
		return
	}
}

// GrowByMin will grow by a single page or by the minimum needed, whichever is greater.
func GrowMin() Grow {
	return func(pagesBefore, pagesNeeded int32, minSize uintptr) (int32, uintptr, uintptr) {
		start, end := growBy(pagesNeeded)
		if start == 0 {
			return 0, 0, 0
		}
		return pagesNeeded, start, end
	}
}

func growBy(pages int32) (uintptr, uintptr) {
	before := wasm_memory_size(0)
	wasm_memory_grow(0, pages)
	after := wasm_memory_size(0)
	if before == after {
		return 0, 0
	}
	return uintptr(before * wasmPageSize), uintptr(after * wasmPageSize)
}

func NewTLSF(pages int32, grow Grow) *Allocator {
	if pages <= 0 {
		pages = 1
	}
	pagesAdded, start, end := grow(0, pages, 0)
	return InitTLSF(start, end, pagesAdded, grow)
}

func newTLSF(pages int32) *Allocator {
	return NewTLSF(pages, GrowMin())
}
