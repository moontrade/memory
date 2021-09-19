package main

import (
	"github.com/wasmerio/wasmer-go/wasmer"
	"os"
	"testing"
)

func BenchmarkLoad(b *testing.B) {
	dylib, err := os.ReadFile("../wasm-local/main.dylib")
	if err != nil {
		b.Fatal(err)
	}

	b.Run("Dylib", func(b *testing.B) {
		engine := wasmer.NewDylibEngine()
		store := wasmer.NewStore(engine)
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			module, err := wasmer.DeserializeModule(store, dylib)
			if err != nil {
				b.Fatal(err)
			}
			module.Close()
		}
	})

	wasm, err := os.ReadFile("../wasm-local/main.wasm")
	if err != nil {
		b.Fatal(err)
	}

	b.Run("JIT", func(b *testing.B) {
		engine := wasmer.NewEngine()
		store := wasmer.NewStore(engine)
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			module, err := wasmer.NewModule(store, wasm)
			if err != nil {
				b.Fatal(err)
			}
			module.Close()
		}
	})

	b.Run("Dylib", func(b *testing.B) {
		engine := wasmer.NewDylibEngine()
		store := wasmer.NewStore(engine)
		module, err := wasmer.DeserializeModule(store, dylib)
		if err != nil {
			b.Fatal(err)
		}
		b.ResetTimer()

		for i := 0; i < b.N; i++ {

			if err != nil {
				b.Fatal(err)
			}
			module.Close()
		}
	})
}
