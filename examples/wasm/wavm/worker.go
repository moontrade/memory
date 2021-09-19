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
	printf("WASM fd_write\n");
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
	"os"
	"unsafe"
)

type WorkerT struct {
	engine      *WASMEngineT
	compartment *WASMCompartmentT
	store       *WASMStoreT
	memory      *WASMMemoryT
	stub        *WASMFuncT
	start       *WASMFuncT
	resume      *WASMFuncT
	alloc       *WASMFuncT
	realloc     *WASMFuncT
	free        *WASMFuncT
	log         *C.char
	len         C.size_t
	cap         C.size_t
}

func (w *WorkerT) Close() error {
	w.store.Delete()
	w.compartment.Delete()
	w.engine.Delete()
	return nil
}

var (
	fd_write_type = WASMFuncTypeNew_4_1(
		WASMValTypeNew(I32), WASMValTypeNew(I32), WASMValTypeNew(I32), WASMValTypeNew(I32),
		WASMValTypeNew(I32),
	)

	//println("fd_write", "params", wasm_functype_num_params(fd_write_type), "results", wasm_functype_num_results(fd_write_type))
	clock_time_get_type = WASMFuncTypeNew_3_1(
		WASMValTypeNew(I32), WASMValTypeNew(I64), WASMValTypeNew(I32),
		WASMValTypeNew(I32),
	)

	args_sizes_get_type = WASMFuncTypeNew_2_1(
		WASMValTypeNew(I32), WASMValTypeNew(I32),
		WASMValTypeNew(I32),
	)

	args_get_type = WASMFuncTypeNew_2_1(
		WASMValTypeNew(I32), WASMValTypeNew(I32),
		WASMValTypeNew(I32))

	setTimeout_type = WASMFuncTypeNew_1_0(WASMValTypeNew(I64))
)

var file []byte
var fileWASM []byte

func init() {
	file, _ = os.ReadFile("examples/wasm/wasm-local/main.compiled")
	if len(file) == 0 {
		file, _ = os.ReadFile("../wasm-local/main.compiled")
		fileWASM, _ = os.ReadFile("../wasm-local/main.wasm")
	} else {
		fileWASM, _ = os.ReadFile("examples/wasm/wasm-local/main.wasm")
	}
}

func Load(precompiled, trace bool) *WorkerT {
	engine := WASMEngineNewWithConfig(WASMConfigNew().
		SetBulkMemoryOps(true).SetMultiMemory(true))
	compartment := engine.NewCompartment("")
	store := compartment.NewStore("")

	var module *WASMModuleT

	if precompiled {
		module = engine.NewPrecompiledModule(file)
	} else {
		module = engine.NewModule(fileWASM)
	}

	worker := &WorkerT{
		engine:      engine,
		compartment: compartment,
		store:       store,
	}

	modImports := module.Imports(nil)
	for _, imp := range modImports {
		if trace {
			println("import", imp.ModuleUnsafe(), "name", imp.NameUnsafe())
		}
	}

	modExports := module.Exports(nil)
	for _, export := range modExports {
		if trace {
			println("export", export.NameUnsafe())
		}
	}

	// Create import funcs
	fd_write := WASMFuncNew(
		compartment,
		fd_write_type,
		(WASMFuncCallbackT)(C.moontrade_fd_write),
		"fd_write",
	)
	clock_time_get := WASMFuncNew(
		compartment,
		clock_time_get_type,
		(WASMFuncCallbackT)(C.moontrade_clock_time_get),
		"clock_time_get",
	)
	args_sizes_get := WASMFuncNew(
		compartment,
		args_sizes_get_type,
		(WASMFuncCallbackT)(C.moontrade_args_sizes_get),
		"args_sizes_get",
	)
	args_get := WASMFuncNew(
		compartment,
		args_get_type,
		(WASMFuncCallbackT)(C.moontrade_args_get),
		"args_get",
	)
	setTimeout := WASMFuncNew(
		compartment,
		setTimeout_type,
		(WASMFuncCallbackT)(C.moontrade_set_timeout),
		"setTimeout",
	)

	// Delete func types
	//WASMFuncTypeDelete(fd_write_type)
	//WASMFuncTypeDelete(clock_time_get_type)
	//WASMFuncTypeDelete(args_sizes_get_type)
	//WASMFuncTypeDelete(args_get_type)
	//WASMFuncTypeDelete(setTimeout_type)

	//imports := []*WASMExternT{
	//	WASMFuncAsExtern(fd_write),
	//	WASMFuncAsExtern(clock_time_get),
	//	WASMFuncAsExtern(args_sizes_get),
	//	WASMFuncAsExtern(args_get),
	//	WASMFuncAsExtern(setTimeout),
	//}
	imports := []*WASMExternT{
		fd_write.AsExtern(),
		clock_time_get.AsExtern(),
		args_sizes_get.AsExtern(),
		args_get.AsExtern(),
		setTimeout.AsExtern(),
	}

	var trap *WASMTrapT
	instance := WASMInstanceNew(store, module, imports, &trap, "")

	fd_write.Delete()
	clock_time_get.Delete()
	args_sizes_get.Delete()
	args_get.Delete()
	setTimeout.Delete()
	//WASMFuncDelete(fd_write)
	//WASMFuncDelete(clock_time_get)
	//WASMFuncDelete(args_sizes_get)
	//WASMFuncDelete(args_get)
	//WASMFuncDelete(setTimeout)

	// Extract exports
	instNumExports := instance.NumExports()
	if instNumExports != len(modExports) {
		panic("instance exports and module exports don't match")
	}

	for i := 0; i < instNumExports; i++ {
		export := modExports[i]
		extern := instance.Export(i)

		switch extern.AsKind() {
		case EXTERN_FUNC:
			if trace {
				println("func extern", export.NameUnsafe())
			}

			fn := extern.AsFunc()
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
			case "stub":
				worker.stub = fn
			}
		case EXTERN_TABLE:
		case EXTERN_MEMORY:
			worker.memory = extern.AsMemory()
		case EXTERN_GLOBAL:
		}
	}

	// Delete module and instance
	module.Delete()
	instance.Delete()

	// Init memory
	data := worker.memory.Data()
	pages := worker.memory.Pages()
	size := worker.memory.Size()
	//C.moontrade_memory_set(memory)
	if trace {
		fmt.Println("data", uintptr(unsafe.Pointer(data)), "pages", uint(pages), "size", size)
	}

	var (
		args    = make([]WASMValT, 2)
		results = make([]WASMValT, 2)
	)

	trap = worker.stub.Call(
		store,
		(*WASMValT)(unsafe.Pointer(&args[0])),
		(*WASMValT)(unsafe.Pointer(&results[0])),
	)

	if trap != nil {
		println(trap.String())
	}

	trap = worker.start.Call(
		store,
		(*WASMValT)(unsafe.Pointer(&args[0])),
		(*WASMValT)(unsafe.Pointer(&results[0])),
	)

	if trap != nil {
		println(trap.String())
	}

	args[0].SetI32(64)

	if trace {
		println("sizeof(WASMValT{})", unsafe.Sizeof(WASMValT{}))
	}

	trap = worker.alloc.Call(
		store,
		(*WASMValT)(unsafe.Pointer(&args[0])),
		(*WASMValT)(unsafe.Pointer(&results[0])),
	)

	if trap != nil {
		println(trap.String())
	}

	if trace {
		result := WASMValGetI32(&results[0])
		println(result, results[0].I32())
	}

	return worker
}

// Worker represents a single running instance / worker
type Worker struct {
}

// Reactor executes and schedules Workers. Reactors are single threaded.
type Reactor struct {
}
