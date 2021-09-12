//go:build !tinygo
// +build !tinygo

package mem

// Arena allocates memory from the underlying platform. It is used to add
// new memory to an Allocator.
type Arena interface {
	Alloc(size uintptr) (uintptr, uintptr)

	Free()
}
