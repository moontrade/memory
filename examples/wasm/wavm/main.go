package main

/*
#cgo CFLAGS: -I/usr/local/include/WAVM
#cgo LDFLAGS: -lWAVM
#include <stdio.h>
#include <stdlib.h>
#include <wavm-c/wavm-c.h>
typedef struct export_t
{
	const char* name;
	size_t num_name_bytes;
	wasm_externtype_t* typ;
} export_t;
*/
import "C"
import (
	"log"
)

func main() {
	log.Print("Initializing WAVM")
	//engine := NewEngine(&Config{
	//	WASIEnabled: true,
	//})
	//rawModule, _ := ioutil.ReadFile("quickjs.wasm.precompiled")

	engine := NewEngine(&Config{})

	_ = engine
	//module := engine.LoadModule(rawModule, true)
	//numExports := module.NumExports()
	//for i := 0; i < numExports; i++ {
	//	export := module.GetExport(i)
	//	fmt.Println(i, export)
	//}
	//instance := engine.WASIRun(module)
	//log.Print("Instance loaded ", instance)
}
