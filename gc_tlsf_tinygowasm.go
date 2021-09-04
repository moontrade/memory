//go:build gc.tlsf && (tinygo.wasm || wasi)
// +build gc.tlsf
// +build tinygo.wasm wasi

package runtime

func newTLSF(pages int32) *tlsf {
	if pages <= 0 {
		pages = 1
	}
	size := uintptr(pages * wasmPageSize)
	start := uintptr(wasm_memory_size(0) * wasmPageSize)
	wasm_memory_grow(0, pages)
	end := uintptr(wasm_memory_size(0) * wasmPageSize)
	if start == end {
		panic("out of memory")
	}
	return initTLSF(start, start+size, pages)
}

func (p *tlsf) Grow(pages int32) (uintptr, uintptr) {
	if pages <= 0 {
		pages = 1
	}

	// wasm memory grow
	before := wasm_memory_size(0)
	wasm_memory_grow(0, pages)
	after := wasm_memory_size(0)
	if before == after {
		return 0, 0
	}

	p.pages += int32(pages)
	start := uintptr(before * wasmPageSize)
	end := uintptr(after * wasmPageSize)
	p.heapSize += int64(end - start)
	p.heapEnd = end
	heapEnd = end
	return start, end
}
