package main

/*
#cgo CFLAGS: -I/usr/local/include/WAVM
#cgo LDFLAGS: -lWAVM
#include <stdio.h>
#include <stdlib.h>
#include <wavm-c/wavm-c.h>
#include <time.h>
typedef struct export_t
{
	const char* name;
	size_t num_name_bytes;
	wasm_externtype_t* typ;
} export_t;

typedef struct moontrade_worker_t
{
	wasm_store_t* store;
	wasm_memory_t* memory;
	wasm_func_t* start;
	wasm_func_t* resume;
	wasm_func_t* alloc;
	wasm_func_t* realloc;
	wasm_func_t* free;
	const char* log;
	size_t log_len;
	size_t log_cap;
} moontrade_worker_t;

moontrade_worker_t MOON_001 = {
	.store = NULL,
	.memory = NULL,
	.start = NULL,
	.resume = NULL,
	.alloc = NULL,
	.realloc = NULL,
	.free = NULL,
	.log = NULL,
	.log_len = 0,
	.log_cap = 0,
};

wasm_memory_t* MEM = NULL;

void moontrade_memory_set(wasm_memory_t* mem) {
	MEM = mem;
}

void moontrade_resume(wasm_instance_t* instance, wasm_memory_t* mem) {
	MEM = mem;
}

// fd_write
wasm_trap_t* moontrade_fd_write(const wasm_val_t args[], wasm_val_t results[])
{
	// Append to worker log buffer
	return NULL;
}

// clock_time_get
wasm_trap_t* moontrade_clock_time_get(const wasm_val_t args[], wasm_val_t results[])
{
	struct timespec now;
	clock_gettime(CLOCK_REALTIME, &now);
	//system_clock::time_point begin = system_clock::now();
	int64_t epoch = ((int64_t)now.tv_sec * (int64_t)1000000000) + (int64_t)now.tv_nsec;
	return NULL;
}

// args_sizes_get
wasm_trap_t* moontrade_args_sizes_get(const wasm_val_t args[], wasm_val_t results[])
{
	return NULL;
}

// args_get
wasm_trap_t* moontrade_args_get(const wasm_val_t args[], wasm_val_t results[])
{
	return NULL;
}

// set_timeout
wasm_trap_t* moontrade_set_timeout(const wasm_val_t args[], wasm_val_t results[])
{
	// Add timeout to VM
	return NULL;
}
*/
import "C"
import (
	"fmt"
	"log"
	"os"
	"unsafe"
)

type Engine struct {
	ptr *C.wasm_engine_t
}

type Env struct {
	fd_write *C.wasm_func_t
	memory   *C.wasm_memory_t
}

type WorkerT struct {
	engine      *WASMEngineT
	compartment *WASMCompartmentT
	store       *WASMStoreT
	memory      *WASMMemoryT
	start       *WASMFuncT
	resume      *WASMFuncT
	alloc       *WASMFuncT
	realloc     *WASMFuncT
	free        *WASMFuncT
	log         *C.char
	len         C.size_t
	cap         C.size_t
}

func main() {
	log.Print("Initializing WAVM")
	file, _ := os.ReadFile("examples/wasm/wasm-local/main.wasm")

	engine := WASMEngineNew()
	compartment := WASMCompartmentNew(engine, "")
	store := WASMStoreNew(compartment, "")
	module := WASMModuleNew(engine, file)

	worker := &WorkerT{
		engine:      engine,
		compartment: compartment,
		store:       store,
	}

	importCount := int(WASMModuleNumImports(module))
	moduleImports := make([]WASMImportT, importCount)
	for i := 0; i < importCount; i++ {
		moduleImports[i] = WASMModuleImport(module, i)
		println("import", moduleImports[i].ModuleUnsafe(), "name", moduleImports[i].NameUnsafe())
	}

	exportCount := int(WASMModuleNumExports(module))
	moduleExports := make([]WASMExportT, exportCount)
	for i := 0; i < exportCount; i++ {
		//var out_import *C.WASMImportT
		moduleExports[i] = WASMModuleExport(module, i)
		println("export", moduleExports[i].NameUnsafe())
	}

	println("num imports", int(WASMModuleNumImports(module)))
	println("num exports", WASMModuleNumExports(module))

	// Create import func types
	fd_write_type := WASMFuncTypeNew_4_1(
		WASMValTypeNew(I32), WASMValTypeNew(I32), WASMValTypeNew(I32), WASMValTypeNew(I32),
		WASMValTypeNew(I32),
	)
	//println("fd_write", "params", wasm_functype_num_params(fd_write_type), "results", wasm_functype_num_results(fd_write_type))
	clock_time_get_type := WASMFuncTypeNew_3_1(
		WASMValTypeNew(I32), WASMValTypeNew(I64), WASMValTypeNew(I32),
		WASMValTypeNew(I32),
	)
	args_sizes_get_type := WASMFuncTypeNew_2_1(
		WASMValTypeNew(I32), WASMValTypeNew(I32),
		WASMValTypeNew(I32),
	)
	args_get_type := WASMFuncTypeNew_2_1(
		WASMValTypeNew(I32), WASMValTypeNew(I32),
		WASMValTypeNew(I32))
	setTimeout_type := WASMFuncTypeNew_1_0(WASMValTypeNew(I64))

	// Create import funcs
	fd_write := WASMFuncNew(
		compartment,
		fd_write_type,
		(C.wasm_func_callback_t)(C.moontrade_fd_write),
		"fd_write",
	)
	clock_time_get := WASMFuncNew(
		compartment,
		clock_time_get_type,
		(C.wasm_func_callback_t)(C.moontrade_clock_time_get),
		"clock_time_get",
	)
	args_sizes_get := WASMFuncNew(
		compartment,
		args_sizes_get_type,
		(C.wasm_func_callback_t)(C.moontrade_args_sizes_get),
		"args_sizes_get",
	)
	args_get := WASMFuncNew(
		compartment,
		args_get_type,
		(C.wasm_func_callback_t)(C.moontrade_args_get),
		"args_get",
	)
	setTimeout := WASMFuncNew(
		compartment,
		setTimeout_type,
		(C.wasm_func_callback_t)(C.moontrade_set_timeout),
		"setTimeout",
	)

	// Delete func types
	WASMFuncTypeDelete(fd_write_type)
	WASMFuncTypeDelete(clock_time_get_type)
	WASMFuncTypeDelete(args_sizes_get_type)
	WASMFuncTypeDelete(args_get_type)
	WASMFuncTypeDelete(setTimeout_type)

	imports := []*WASMExternT{
		WASMFuncAsExtern(fd_write),
		WASMFuncAsExtern(clock_time_get),
		WASMFuncAsExtern(args_sizes_get),
		WASMFuncAsExtern(args_get),
		WASMFuncAsExtern(setTimeout),
		WASMFuncAsExtern(setTimeout),
	}

	var trap *C.wasm_trap_t
	instance := WASMInstanceNew(store, module, imports, &trap, "")

	WASMFuncDelete(fd_write)
	WASMFuncDelete(clock_time_get)
	WASMFuncDelete(args_sizes_get)
	WASMFuncDelete(args_get)
	WASMFuncDelete(setTimeout)

	// Extract exports
	instNumExports := int(WASMInstanceNumExports(instance))
	if instNumExports != len(moduleExports) {
		panic("instance exports and module exports don't match")
	}

	for i := 0; i < instNumExports; i++ {
		export := moduleExports[i]
		extern := WASMInstanceExport(instance, i)
		switch WASMExternKind(extern) {
		case EXTERN_FUNC:
			println("func extern", export.NameUnsafe())
			fn := WASMExternAsFunc(extern)
			switch export.NameUnsafe() {
			case "_start":
				worker.start = fn
			case "resume":
				worker.resume = fn
			case "alloc":
				worker.alloc = fn
			case "realloc":
				worker.realloc = fn
			case "free":
				worker.free = fn
			}
		case EXTERN_TABLE:
		case EXTERN_MEMORY:
			worker.memory = WASMExternAsMemory(extern)
		case EXTERN_GLOBAL:
		}
	}

	// Delete module and instance
	WASMModuleDelete(module)
	WASMInstanceDelete(instance)

	// Init memory
	data := WASMMemoryData(worker.memory)
	pages := int(WASMMemorySize(worker.memory))
	size := uint(WASMMemoryDataSize(worker.memory))
	//C.moontrade_memory_set(memory)
	fmt.Println("data", uintptr(unsafe.Pointer(data)), "pages", uint(pages), "size", size)

	//out := make([]byte, 256000)
	//o := C.CBytes(out)

	//size := C.size_t(0)
	//out := C.WASMModulePrint(module, &size)
	//os.WriteFile("main.wast", C.GoBytes(unsafe.Pointer(out), C.int(size)), 0755)
	//println(string(C.GoBytes(unsafe.Pointer(out), C.int(size))))
	//wast := []byte(helloWast)
	//wast = append(wast, 0)
	//helloWast = string(wast)
	//module := wasm_module_new_text(engine, helloWast)

	//imports := make([]*C.wasm_extern_t, 0, 16)
	//
	//helloType := WASMFuncTypeNew_2_1(wasm_valty

	// Shutdown
	WASMStoreDelete(store)
	WASMCompartmentDelete(compartment)
	WASMEngineDelete(engine)
}

// Worker represents a single running instance / worker
type Worker struct {
}

// Reactor executes and schedules Workers. Reactors are single threaded.
type Reactor struct {
}
