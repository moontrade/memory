package main

import (
	_ "github.com/moontrade/memory"
	mem "github.com/moontrade/memory"
	"runtime"
	"time"
)

var done = make(chan bool, 1)
var b []byte

type _bytes struct {
	data uintptr
	len  uintptr
	cap  uintptr
}

func main() {
	//mem.IsPowerOfTwo(0)
	println("hi moontrade!")

	go func() {
		for {
			mem.Scope(func(a mem.Auto) {
				a.Alloc(512)
			})
			if b == nil {
				//p := mem.Alloc(128000)
				//b = *(*[]byte)(unsafe.Pointer(&_bytes{uintptr(p), 128000,128000}))
				b = make([]byte, 128)
				b[0] = 10
			}
			println(time.Now().UnixNano())
			//start := time.Now().UnixNano()
			runtime.GC()
			//println("full GC", time.Now().UnixNano()-start)
			time.Sleep(time.Second)
		}
	}()

	<-done
}
