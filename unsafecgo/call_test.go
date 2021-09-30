package unsafecgo

import (
	"github.com/moontrade/memory/unsafecgo/cgo"
	"runtime"
	"testing"
	"time"
)

func BenchmarkCall(b *testing.B) {
	b.Run("asm", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			NonBlocking((*byte)(cgo.Stub), 0, 0)
		}
	})
	b.Run("libcCall", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			cgo.NonBlocking((*byte)(cgo.Stub), 0, 0)
		}
	})
	b.Run("cgo", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			cgo.CGO()
		}
	})
	b.Run("cgo trampoline", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			cgo.Blocking((*byte)(cgo.Stub), 0, 0)
		}
	})
}

func TestSleep(t *testing.T) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	for i := 0; i < 10000; i++ {
		//NonBlocking((*byte)(cgo.Usleep), uintptr(time.Second), 0)
		cgo.NonBlocking((*byte)(cgo.Usleep), uintptr(time.Second/1000), 0)
		//cgo.DoUsleep(int64(time.Second/1000))
		println(time.Now().UnixNano())
	}
}
