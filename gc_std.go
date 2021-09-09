//go:build !tinygo && !tinygo.wasm && !wasi

package mem

var heapStart uintptr = 0

func markStack() {

}

func markGlobals() {

}

func markScheduler() {

}
