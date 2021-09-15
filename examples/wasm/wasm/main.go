package main

import (
	mem "github.com/moontrade/memory"
	"time"
)

//export llvm.wasm.memory.size.i32
func wasm_memory_size(index int32) int32

//export llvm.wasm.memory.grow.i32
func wasm_memory_grow(index, pages int32) int32

var done = make(chan bool, 1)
var b []byte

func main() {
	println("hi moontrade!")

	go func() {
		for {
			mem.Scope(func(a mem.Auto) {
				a.Alloc(512)
			})
			b = make([]byte, 4096)
			b[0] = 10
			println(time.Now().UnixNano())
			time.Sleep(time.Second)
		}
	}()

	<-done
}
