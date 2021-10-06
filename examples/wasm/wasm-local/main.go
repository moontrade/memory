package main

import (
	"runtime"
	"time"

	"github.com/moontrade/memory"
)

var done = make(chan bool, 1)
var b []byte

//export stub
func stub() {}

func main() {
	//mem.IsPowerOfTwo(0)
	println("hi moontrade!")

	go func() {
		for {
			memory.Scope(func(a memory.AutoFree) {
				a.Alloc(512)
			})
			if b == nil {
				b = make([]byte, 128)
				b[0] = 10
			}
			//b = make([]byte, 65536)
			println(time.Now().UnixNano())
			//start := time.Now().UnixNano()
			runtime.GC()
			//println("full GC", time.Now().UnixNano()-start)
			time.Sleep(time.Second)
		}
	}()

	<-done
}
