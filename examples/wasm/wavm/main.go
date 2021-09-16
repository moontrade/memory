package main

import (
	"fmt"
	"io/ioutil"
	"log"
)

func main() {
	log.Print("Initializing WAVM")
	engine := NewEngine(&Config{
		WASIEnabled: true,
	})
	rawModule, _ := ioutil.ReadFile("quickjs.wasm.precompiled")
	module := engine.LoadModule(rawModule, true)
	numExports := module.NumExports()
	for i := 0; i < numExports; i++ {
		export := module.GetExport(i)
		fmt.Println(i, export)
	}
	instance := engine.WASIRun(module)
	log.Print("Instance loaded ", instance)
}
