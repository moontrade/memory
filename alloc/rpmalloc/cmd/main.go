package main

// #include <stdlib.h>
import "C"
import (
	"github.com/moontrade/memory"
	"runtime"
	"time"
)

func main() {
	memory.Init()

	for i := 0; i < 2; i++ {
		go func() {
			runtime.LockOSThread()
			memory.Free(memory.AllocZeroed(128))
			C.free(C.malloc(128))
			time.Sleep(time.Second * 5)
		}()
	}

	time.Sleep(time.Hour)
}
