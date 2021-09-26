package main

import (
	"github.com/moontrade/memory"
	_ "github.com/moontrade/memory"
	"runtime"
	"time"
)

func main() {
	memory.Init()
	memory.Free(memory.Alloc(128))

	for i := 0; i < 100; i++ {
		go func() {
			runtime.LockOSThread()
			//C.free(C.malloc(128))
			time.Sleep(time.Second * 5)
		}()
	}

	time.Sleep(time.Hour)
}
