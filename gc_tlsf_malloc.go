//go:build !tinygo.wasm && !wasm

package runtime

func newTLSF(pages int32) *tlsf {
	if pages <= 0 {
		pages = 1
	}
	size := uintptr(pages * wasmPageSize)
	segment := uintptr(malloc(size))
	return initTLSF(segment, segment+size, pages)
}

func (p *tlsf) Grow(pages int32) (uintptr, uintptr) {
	if pages <= 0 {
		pages = 1
	}
	// Allocate new segment
	segment := uintptr(malloc(uintptr(pages * wasmPageSize)))
	if segment == 0 {
		return 0, 0
	}
	p.pages += pages
	start, end := segment, segment+uintptr(pages*wasmPageSize)
	p.heapSize += int64(pages * wasmPageSize)
	p.heapEnd = end
	return start, end
}
