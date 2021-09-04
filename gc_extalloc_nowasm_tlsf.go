//go:build !tinygo.wasm && !wasm && !wasi && gc.extalloc && gc.tlsf
// +build !tinygo.wasm,!wasm,!wasi,gc.extalloc,gc.tlsf

package runtime

func initHeap() {
}
