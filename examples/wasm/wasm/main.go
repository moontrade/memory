package main

import (
	"runtime"
	"time"
)

var done = make(chan bool, 1)
var b []byte

func main() {
	println("hi moontrade!")

	go func() {
		for {
			//mem.Scope(func(a mem.Auto) {
			//	a.Alloc(512)
			//})
			if b == nil {
				b = make([]byte, 4096)
				b[0] = 10
			}
			println(time.Now().UnixNano())
			runtime.GC()
			time.Sleep(time.Second)
		}
	}()

	<-done
}
