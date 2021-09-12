//go:build !tinygo.wasm
// +build !tinygo.wasm

package mem

func assert(truthy bool, message string) {
	if !truthy {
		panic(message)
	}
}
