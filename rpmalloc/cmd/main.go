package main

// #include <stdlib.h>
import "C"
import (
	"github.com/moontrade/memory"
	"runtime"
	"time"
)

func main() {
	memory.Free(memory.Alloc(128))

	for i := 0; i < 100; i++ {
		go func() {
			runtime.LockOSThread()
			C.free(C.malloc(128))
			time.Sleep(time.Minute)
		}()
	}

	time.Sleep(time.Hour)
}
