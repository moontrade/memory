//go:build !tinygo.wasm
// +build !tinygo.wasm

package tlsf

func assert(truth bool, message string) {
	if !truth {
		panic(message)
	}
}
