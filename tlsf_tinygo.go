//go:build tinygo
// +build tinygo

package mem

// tinygo has a single global allocator
var allocator *TLSF
