package main

import (
	"runtime"
	"testing"
	"unsafe"
)

func TestStub(t *testing.T) {
	worker := Load(true, true)

	var (
		args    = make([]WASMValT, 2)
		results = make([]WASMValT, 2)
		store   = worker.store
		stub    = worker.stub
		trap    *WASMTrapT
	)

	trap = stub.Call(
		store,
		(*WASMValT)(unsafe.Pointer(&args[0])),
		(*WASMValT)(unsafe.Pointer(&results[0])),
	)
	if trap != nil {
		t.Fatal(trap.String())
	}

	_ = worker.Close()
}

func BenchmarkLoad(b *testing.B) {
	engine := WASMEngineNew()
	b.Run("WASM", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			WASMModuleNew(engine, fileWASM)
			//Load(false, false).Close()
		}
	})

	b.Run("Compiled", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			WASMModulePrecompiledNew(engine, file)
			//Load(true, false).Close()
		}
	})
}

func BenchmarkStub(b *testing.B) {
	worker := Load(true, false)

	b.Run("Safe", func(b *testing.B) {
		var (
			//args    = make([]WASMValT, 2)
			//results = make([]WASMValT, 2)
			store = worker.store
			stub  = worker.stub
			trap  *WASMTrapT
		)
		runtime.LockOSThread()
		defer runtime.UnlockOSThread()
		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			trap = WASMFuncCall(
				store,
				stub,
				nil, nil,
				//(*WASMValT)(unsafe.Pointer(&args[0])),
				//(*WASMValT)(unsafe.Pointer(&results[0])),
			)
			if trap != nil {
				b.Fatal(WASMTrapMessageString(trap))
			}
		}
	})

	b.Run("No Copy", func(b *testing.B) {
		var (
			//args    = make([]WASMValT, 2)
			//results = make([]WASMValT, 2)
			store = worker.store
			stub  = worker.stub
			trap  *WASMTrapT
		)
		runtime.LockOSThread()
		defer runtime.UnlockOSThread()
		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			trap = WASMFuncCallNoCopy(
				store,
				stub,
				nil, nil,
				//(*WASMValT)(unsafe.Pointer(&args[0])),
				//(*WASMValT)(unsafe.Pointer(&results[0])),
			)
			if trap != nil {
				b.Fatal(WASMTrapMessageString(trap))
			}
		}
	})

	b.Run("No Trap", func(b *testing.B) {
		var (
			//args    = make([]WASMValT, 2)
			//results = make([]WASMValT, 2)
			store = worker.store
			stub  = worker.stub
			trap  *WASMTrapT
		)
		runtime.LockOSThread()
		defer runtime.UnlockOSThread()
		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			trap = WASMFuncCallNoTrap(
				store,
				stub,
				nil, nil,
				//(*WASMValT)(unsafe.Pointer(&args[0])),
				//(*WASMValT)(unsafe.Pointer(&results[0])),
			)
			if trap != nil {
				b.Fatal(WASMTrapMessageString(trap))
			}
		}
	})

	b.Run("No Copy and No Trap", func(b *testing.B) {
		var (
			//args    = make([]WASMValT, 2)
			//results = make([]WASMValT, 2)
			store = worker.store
			stub  = worker.stub
			trap  *WASMTrapT
		)
		runtime.LockOSThread()
		defer runtime.UnlockOSThread()
		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			trap = WASMFuncCallNoCopyNoTrap(
				store,
				stub,
				nil, nil,
				//(*WASMValT)(unsafe.Pointer(&args[0])),
				//(*WASMValT)(unsafe.Pointer(&results[0])),
			)
			if trap != nil {
				b.Fatal(WASMTrapMessageString(trap))
			}
		}
	})

	//worker.Close()
}
