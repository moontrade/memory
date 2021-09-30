package main

import (
	"github.com/moontrade/memory"
	//_ "github.com/moontrade/memory/alloc"
	//"github.com/moontrade/memory/alloc"
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
				//p := mem.Alloc(128000)
				//b = *(*[]byte)(unsafe.Pointer(&_bytes{uintptr(p), 128000,128000}))
				b = make([]byte, 128)
				b[0] = 10
			}
			b = make([]byte, 65536)
			println(time.Now().UnixNano())
			//start := time.Now().UnixNano()
			runtime.GC()
			//println("full GC", time.Now().UnixNano()-start)
			time.Sleep(time.Second)
		}
	}()

	<-done
}
