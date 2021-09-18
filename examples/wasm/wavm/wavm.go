package main

/*
#cgo CFLAGS: -I/usr/local/include/WAVM
#cgo LDFLAGS: -lWAVM
#include <inttypes.h>
#include <stdio.h>
#include <stdlib.h>
#include <stdint.h>
#include <string.h>
#include <time.h>
//#include <chrono>
//#include <wavm-c/wavm-c.h>

#include "wavm-c.h"
typedef struct export_t
{
	const char* name;
	size_t num_name_bytes;
	wasm_externtype_t* typ;
} export_t;


static void handle_trap(wasm_trap_t* trap)
{
	char* message_buffer = malloc(1024);
	size_t num_message_bytes = 1024;
	if(!wasm_trap_message(trap, message_buffer, &num_message_bytes))
	{
		message_buffer = malloc(num_message_bytes);
		assert(wasm_trap_message(trap, message_buffer, &num_message_bytes));
	}

	fprintf(stderr, "Runtime error: %.*s\n", (int)num_message_bytes, message_buffer);

	wasm_trap_delete(trap);
	free(message_buffer);
}

// A function to be called from Wasm code.
static wasm_trap_t* hello_callback(const wasm_val_t args[], wasm_val_t results[])
{
	// Make a copy of the string passed as an argument, and ensure that it is null terminated.
	const uint32_t address = (uint32_t)args[0].i32;
	const uint32_t num_chars = (uint32_t)args[1].i32;
	char buffer[1025];
	if(num_chars > 1024)
	{
		fprintf(stderr, "Callback called with too many characters: num_chars=%u.\n", num_chars);
		exit(1);
	}

	//const size_t num_memory_bytes = wasm_memory_data_size(memory);
	//if(((uint64_t)address + (uint64_t)num_chars) > num_memory_bytes)
	//{
	//	fprintf(stderr,
	//			"Callback called with out-of-bounds string address:\n"
	//			"  address=%u\n"
	//			"  num_chars=%u\n"
	//			"  wasm_memory_data_size(memory)=%zu\n",
	//			address,
	//			num_chars,
	//			wasm_memory_data_size(memory));
	//	exit(1);
	//}
	//
	//memcpy(buffer, wasm_memory_data(memory) + address, num_chars);
	//buffer[num_chars] = 0;
	//
	//printf("Hello, %s!\n", buffer);
	results[0].i32 = num_chars;
	return NULL;
}

//wasm_trap_t* wasi_fd_write(const wasm_val_t args[], wasm_val_t results[]);
// fd_write
//wasm_trap_t* wasi_fd_write(const wasm_val_t args[], wasm_val_t results[])
//{
//	return NULL;
//}
//
//// clock_time_get
//wasm_trap_t* wasi_clock_time_get(const wasm_val_t args[], wasm_val_t results[])
//{
//	struct timespec now;
//	clock_gettime(CLOCK_REALTIME, &now);
//	//system_clock::time_point begin = system_clock::now();
//	int64_t epoch = ((int64_t)now.tv_sec * (int64_t)1000000000) + (int64_t)now.tv_nsec;
//	return NULL;
//}
//
//// args_sizes_get
//wasm_trap_t* wasi_args_sizes_get(const wasm_val_t args[], wasm_val_t results[])
//{
//	return NULL;
//}
//
//// args_get
//wasm_trap_t* wasi_args_get(const wasm_val_t args[], wasm_val_t results[])
//{
//	return NULL;
//}
//
//// set_timeout
//wasm_trap_t* set_timeout(const wasm_val_t args[], wasm_val_t results[])
//{
//	return NULL;
//}


*/
import "C"
import (
	"math"
	"reflect"
	"unsafe"
)

type (
	WASMConfigT      C.wasm_config_t
	WASMEngineT      C.wasm_engine_t
	WASMCompartmentT C.wasm_compartment_t
	WASMStoreT       C.wasm_store_t
	WASMValTypeT     C.wasm_valtype_t
	WASMFuncTypeT    C.wasm_functype_t
	WASMTableTypeT   C.wasm_tabletype_t
	WASMMemoryTypeT  C.wasm_memorytype_t
	WASMGlobalTypeT  C.wasm_globaltype_t
	WASMExternTypeT  C.wasm_externtype_t
	WASMRefT         C.wasm_ref_t
	WASMTrapT        C.wasm_trap_t
	WASMFrameT       C.wasm_frame_t
	WASMForeignT     C.wasm_foreign_t
	WASMModuleT      C.wasm_module_t
	WASMFuncT        C.wasm_func_t
	WASMTableT       C.wasm_table_t
	WASMMemoryT      C.wasm_memory_t
	WASMGlobalT      C.wasm_global_t
	WASMExternT      C.wasm_extern_t
	WASMInstanceT    C.wasm_instance_t
	WASMValT         C.wasm_val_t

	WASMLimitsT     C.wasm_limits_t
	WASMSharedT     C.wasm_shared_t
	WASMMutabilityT C.wasm_mutability_t

	/*
		typedef struct wasm_import_t
		{
			const char* module;
			size_t num_module_bytes;
			const char* name;
			size_t num_name_bytes;
			wasm_externtype_t* type;
		} wasm_import_t;
	*/
	WASMImportT struct {
		module           *C.char
		num_module_bytes C.size_t
		name             *C.char
		num_name_bytes   C.size_t
		_type            *C.wasm_externtype_t
	}

	/*
		typedef struct wasm_export_t
		{
			const char* name;
			size_t num_name_bytes;
			wasm_externtype_t* type;
		} wasm_export_t;
	*/
	WASMExportT struct {
		name           *C.char
		num_name_bytes C.size_t
		_type          *C.wasm_externtype_t
	}
)

const (
	WASM_MEMORY_PAGE_SIZE = 0x10000
	WASM_MEMORY_PAGES_MAX = math.MaxUint32
	WASM_TABLE_SIZE_MAX   = math.MaxUint32

	I32     = C.WASM_I32
	I64     = C.WASM_I64
	F32     = C.WASM_F32
	F64     = C.WASM_F64
	V128    = C.WASM_V128
	ANYREF  = C.WASM_ANYREF
	FUNCREF = C.WASM_FUNCREF

	EXTERN_FUNC   = C.WASM_EXTERN_FUNC
	EXTERN_TABLE  = C.WASM_EXTERN_TABLE
	EXTERN_MEMORY = C.WASM_EXTERN_MEMORY
	EXTERN_GLOBAL = C.WASM_EXTERN_GLOBAL
)

var _EMPTY = C.CString("")

// Configuration

func WASMConfigDelete(config *WASMConfigT) {
	C.wasm_config_delete((*C.wasm_config_t)(config))
}

func WASMConfigNew() *WASMConfigT {
	return (*WASMConfigT)(C.wasm_config_new())
}

func WASMConfigFeatureSetImportExportMutableGlobals(config *WASMConfigT, enable bool) {
	C.wasm_config_feature_set_import_export_mutable_globals((*C.wasm_config_t)(config), (C.bool)(enable))
}

func WASMConfigFeatureSetSIMD(config *WASMConfigT, enable bool) {
	C.wasm_config_feature_set_simd((*C.wasm_config_t)(config), (C.bool)(enable))
}

func WASMConfigFeatureSetAtomics(config *WASMConfigT, enable bool) {
	C.wasm_config_feature_set_atomics((*C.wasm_config_t)(config), (C.bool)(enable))
}

// Engine

func WASMEngineDelete(engine *WASMEngineT) {
	C.wasm_engine_delete((*C.wasm_engine_t)(engine))
}

func WASMEngineNew() *WASMEngineT {
	return (*WASMEngineT)(C.wasm_engine_new())
}

func WASMEngineNewWithConfig(config *WASMConfigT) *WASMEngineT {
	return (*WASMEngineT)(C.wasm_engine_new_with_config((*C.wasm_config_t)(config)))
}

// Compartments

func WASMCompartmentDelete(compartment *WASMCompartmentT) {
	C.wasm_compartment_delete((*C.wasm_compartment_t)(compartment))
}

func WASMCompartmentNew(engine *WASMEngineT, debugName string) *WASMCompartmentT {
	var cstr *C.char
	if debugName == "" {
		cstr = _EMPTY
	} else {
		cstr = C.CString(debugName)
		defer C.free(unsafe.Pointer(cstr))
	}
	return (*WASMCompartmentT)(C.wasm_compartment_new((*C.wasm_engine_t)(engine), cstr))
}

func WASMCompartmentClone(compartment *WASMCompartmentT) *WASMCompartmentT {
	return (*WASMCompartmentT)(C.wasm_compartment_clone((*C.wasm_compartment_t)(compartment)))
}

func WASMCompartmentContains(compartment *WASMCompartmentT, ref *C.wasm_ref_t) bool {
	return bool(C.wasm_compartment_contains((*C.wasm_compartment_t)(compartment), ref))
}

// Store

func WASMStoreNew(compartment *WASMCompartmentT, debugName string) *WASMStoreT {
	var cstr *C.char
	if debugName == "" {
		cstr = _EMPTY
	} else {
		cstr = C.CString(debugName)
		defer C.free(unsafe.Pointer(cstr))
	}
	return (*WASMStoreT)(C.wasm_store_new((*C.wasm_compartment_t)(compartment), cstr))
}

func WASMStoreDelete(store *WASMStoreT) {
	C.wasm_store_delete((*C.wasm_store_t)(store))
}

// Instance

func WASMInstanceDelete(instance *WASMInstanceT) {
	C.wasm_instance_delete((*C.wasm_instance_t)(instance))
}

func WASMInstanceNew(
	store *WASMStoreT,
	module *WASMModuleT,
	imports []*WASMExternT,
	outTrap **C.wasm_trap_t,
	debugName string,
) *WASMInstanceT {
	var cstr *C.char
	if debugName == "" {
		cstr = _EMPTY
	} else {
		cstr = C.CString(debugName)
		defer C.free(unsafe.Pointer(cstr))
	}
	return (*WASMInstanceT)(C.wasm_instance_new(
		(*C.wasm_store_t)(store),
		(*C.wasm_module_t)(module),
		(**C.wasm_extern_t)(unsafe.Pointer(&imports[0])),
		outTrap,
		cstr,
	))
}

func WASMInstanceNumExports(inst *WASMInstanceT) int {
	return int(C.wasm_instance_num_exports((*C.wasm_instance_t)(inst)))
}

func WASMInstanceExport(inst *WASMInstanceT, index int) *WASMExternT {
	return (*WASMExternT)(C.wasm_instance_export((*C.wasm_instance_t)(inst), C.size_t(index)))
}

// Function Types

func WASMFuncTypeDelete(t *WASMFuncTypeT) {
	C.wasm_functype_delete((*C.wasm_functype_t)(t))
}

func WASMFuncTypeNew(params []*WASMValTypeT, results []*WASMValTypeT) *WASMFuncTypeT {
	var _params **C.wasm_valtype_t
	var _results **C.wasm_valtype_t

	if len(params) > 0 {
		_params = (**C.wasm_valtype_t)(unsafe.Pointer(&params[0]))
	} else {
		_params = nil
	}
	if len(results) > 0 {
		_results = (**C.wasm_valtype_t)(unsafe.Pointer(&results[0]))
	} else {
		_results = nil
	}

	return (*WASMFuncTypeT)(C.wasm_functype_new(
		_params,
		(C.size_t)(len(params)),
		_results,
		(C.size_t)(len(results)),
	))
}

func WASMFuncTypeNumParams(funcType *WASMFuncTypeT) int {
	return int(C.wasm_functype_num_params((*C.wasm_functype_t)(funcType)))
}

func WASMFuncTypeParam(funcType *WASMFuncTypeT, index int) *WASMValTypeT {
	return (*WASMValTypeT)(C.wasm_functype_param((*C.wasm_functype_t)(funcType), (C.size_t)(index)))
}

func WASMFuncTypeNumResults(funcType *WASMFuncTypeT) int {
	return int(C.wasm_functype_num_results((*C.wasm_functype_t)(funcType)))
}

func WASMFuncTypeResult(funcType *WASMFuncTypeT, index int) *WASMValTypeT {
	return (*WASMValTypeT)(C.wasm_functype_result((*C.wasm_functype_t)(funcType), (C.size_t)(index)))
}

// Global Types

func WASMGlobalTypeNew(valType *WASMValTypeT, mutability WASMMutabilityT) *C.wasm_globaltype_t {
	return C.wasm_globaltype_new((*C.wasm_valtype_t)(valType), (C.wasm_mutability_t)(mutability))
}

func WASMGlobalTypeContent(globalType *WASMGlobalTypeT) *WASMValTypeT {
	return (*WASMValTypeT)(C.wasm_globaltype_content((*C.wasm_globaltype_t)(globalType)))
}

func WASMGlobalTypeMutability(globalType *WASMGlobalTypeT) WASMMutabilityT {
	return (WASMMutabilityT)(C.wasm_globaltype_mutability((*C.wasm_globaltype_t)(globalType)))
}

// Table Types

//func wasm_tabletype_new(
//	valType *C.wasm_valtype_t,
//	limits *C.wasm_limits_t,
//	shared C.wasm_shared_t,
//) *C.wasm_tabletype_t {
//	return C.wasm_tabletype_new(valType, limits, shared)
//}

func WASMTableTypeElement(tableType *WASMTableTypeT) *WASMValTypeT {
	return (*WASMValTypeT)(C.wasm_tabletype_element((*C.wasm_tabletype_t)(tableType)))
}

func WASMTableTypeLimits(tableType *WASMTableTypeT) *WASMLimitsT {
	return (*WASMLimitsT)(C.wasm_tabletype_limits((*C.wasm_tabletype_t)(tableType)))
}

func WASMTableTypeShared(tableType *WASMTableTypeT) WASMSharedT {
	return (WASMSharedT)(C.wasm_tabletype_shared((*C.wasm_tabletype_t)(tableType)))
}

// Memory Types

func WASMMemoryTypeNew(limits *WASMLimitsT, shared WASMSharedT, index int) *WASMMemoryTypeT {
	return (*WASMMemoryTypeT)(C.wasm_memorytype_new((*C.wasm_limits_t)(limits), (C.wasm_shared_t)(shared), (C.wasm_index_t)(index)))
}

func WASMMemoryTypeLimits(memoryType *WASMMemoryTypeT) *WASMLimitsT {
	return (*WASMLimitsT)(C.wasm_memorytype_limits((*C.wasm_memorytype_t)(memoryType)))
}

func WASMMemoryTypeShared(memoryType *WASMMemoryTypeT) WASMSharedT {
	return (WASMSharedT)(C.wasm_memorytype_shared((*C.wasm_memorytype_t)(memoryType)))
}

// Imports

func (w *WASMImportT) ModuleUnsafe() string {
	return *(*string)(unsafe.Pointer(&reflect.StringHeader{
		Data: uintptr(unsafe.Pointer(w.module)),
		Len:  int(w.num_module_bytes),
	}))
}

func (w *WASMImportT) NameUnsafe() string {
	return *(*string)(unsafe.Pointer(&reflect.StringHeader{
		Data: uintptr(unsafe.Pointer(w.name)),
		Len:  int(w.num_name_bytes),
	}))
}

// Exports

func (w *WASMExportT) NameUnsafe() string {
	return *(*string)(unsafe.Pointer(&reflect.StringHeader{
		Data: uintptr(unsafe.Pointer(w.name)),
		Len:  int(w.num_name_bytes),
	}))
}

// Values
func WASMValDelete(kind C.wasm_valkind_t, v *C.wasm_val_t) {
	C.wasm_val_delete(kind, v)
}

//func wasm_val_copy(kind C.wasm_valkind_t, out *C.wasm_val_t, v *C.wasm_val_t)

// Traps

func WASMTrapDelete(trap *WASMTrapT) {
	C.wasm_trap_delete((*C.wasm_trap_t)(trap))
}

func WASMTrapNew(compartment *WASMCompartmentT, message string) *WASMTrapT {
	msgbytes := []byte(message)
	return (*WASMTrapT)(C.wasm_trap_new(
		(*C.wasm_compartment_t)(compartment),
		(*C.char)(unsafe.Pointer(&msgbytes[0])),
		(C.size_t)(len(msgbytes))))
}

func WASMTrapMessage(trap *WASMTrapT) string {
	var (
		msg  *C.char
		size C.size_t
	)
	C.wasm_trap_message((*C.wasm_trap_t)(trap), msg, &size)
	return *(*string)(unsafe.Pointer(&reflect.StringHeader{Data: uintptr(unsafe.Pointer(msg)), Len: int(size)}))
}

func WASMTrapStackNumFrames(trap *WASMTrapT) int {
	return int(C.wasm_trap_stack_num_frames((*C.wasm_trap_t)(trap)))
}

func WASMTrapStackFrame(trap *WASMTrapT, index int, outFrame *WASMFrameT) {
	C.wasm_trap_stack_frame((*C.wasm_trap_t)(trap), (C.size_t)(index), (*C.wasm_frame_t)(outFrame))
}

// Foreign Objects

func WASMForeignDelete(foreign *WASMForeignT) {
	C.wasm_foreign_delete((*C.wasm_foreign_t)(foreign))
}

func WASMForeignNew(compartment *WASMCompartmentT, debugName string) *WASMForeignT {
	var cstr *C.char
	if debugName == "" {
		cstr = _EMPTY
	} else {
		cstr = C.CString(debugName)
		defer C.free(unsafe.Pointer(cstr))
	}
	return (*WASMForeignT)(C.wasm_foreign_new((*C.wasm_compartment_t)(compartment), cstr))
}

// Modules

func WASMModuleDelete(module *WASMModuleT) {
	C.wasm_module_delete((*C.wasm_module_t)(module))
}

func WASMModuleNew(engine *WASMEngineT, binary []byte) *WASMModuleT {
	slice := (*reflect.StringHeader)(unsafe.Pointer(&binary))
	return (*WASMModuleT)(C.wasm_module_new((*C.wasm_engine_t)(engine), (*C.char)(unsafe.Pointer(slice.Data)), (C.size_t)(slice.Len)))
}

//func wasm_module_new_text(engine *C.wasm_engine_t, wast string) *C.wasm_module_t {
//	return wasm_module_new_wast(engine, []byte(wast))
//}

func WASMModuleNewText(engine *WASMEngineT, wast string) *WASMModuleT {
	if len(wast) == 0 {
		return nil
	}
	ptr := C.CString(wast)
	defer C.free(unsafe.Pointer(ptr))

	println(C.strlen(ptr))
	return (*WASMModuleT)(C.wasm_module_new_text((*C.wasm_engine_t)(engine), ptr, C.size_t(C.strlen(ptr))))
}

func WASMModulePrint(module *WASMModuleT) string {
	var out C.size_t
	ptr := C.wasm_module_print((*C.wasm_module_t)(module), &out)
	return C.GoStringN(ptr, (C.int)(out))
}

func WASMModuleValidate(binary []byte) bool {
	if len(binary) == 0 {
		return false
	}
	slice := (*reflect.SliceHeader)(unsafe.Pointer(&binary))
	return bool(C.wasm_module_validate((*C.char)(unsafe.Pointer(slice.Data)), (C.size_t)(slice.Len)))
}

func WASMModuleNumImports(module *WASMModuleT) int {
	return int(C.wasm_module_num_imports((*C.wasm_module_t)(module)))
}

func WASMModuleImport(module *WASMModuleT, index int) WASMImportT {
	var out C.wasm_import_t
	C.wasm_module_import((*C.wasm_module_t)(module), (C.size_t)(index), &out)
	return *(*WASMImportT)(unsafe.Pointer(&out))
}

func WASMModuleNumExports(module *WASMModuleT) int {
	return int(C.wasm_module_num_exports((*C.wasm_module_t)(module)))
}

func WASMModuleExport(module *WASMModuleT, index int) WASMExportT {
	var out C.wasm_export_t
	C.wasm_module_export((*C.wasm_module_t)(module), (C.size_t)(index), &out)
	return *(*WASMExportT)(unsafe.Pointer(&out))
}

// Function Instances

func WASMFuncDelete(fn *WASMFuncT) {
	C.wasm_func_delete((*C.wasm_func_t)(fn))
}

func WASMFuncNew(
	compartment *WASMCompartmentT,
	funcType *WASMFuncTypeT,
	callback C.wasm_func_callback_t,
	debugName string) *WASMFuncT {
	cstr := C.CString(debugName)
	defer C.free(unsafe.Pointer(cstr))
	return (*WASMFuncT)(C.wasm_func_new((*C.wasm_compartment_t)(compartment), (*C.wasm_functype_t)(funcType), callback, cstr))
}

func WASMFuncNewWithEnv(
	compartment *WASMCompartmentT,
	funcType *WASMFuncTypeT,
	callback C.wasm_func_callback_with_env_t,
	env *C.void,
	finalizer func(env *C.void),
	debugName string,
) *WASMFuncT {
	cstr := C.CString(debugName)
	defer C.free(unsafe.Pointer(cstr))
	return (*WASMFuncT)(C.wasm_func_new((*C.wasm_compartment_t)(compartment), (*C.wasm_functype_t)(funcType), callback, cstr))
}

func WASMFuncType(fn *WASMFuncT) *WASMFuncTypeT {
	return (*WASMFuncTypeT)(C.wasm_func_type((*C.wasm_func_t)(fn)))
}

func WASMFuncParamArity(fn *WASMFuncT) int {
	return int(C.wasm_func_param_arity((*C.wasm_func_t)(fn)))
}

func WASMFuncResultArity(fn *WASMFuncT) int {
	return int(C.wasm_func_result_arity((*C.wasm_func_t)(fn)))
}

// Global Instances

func WASMGlobalDelete(global *C.wasm_global_t) {
	C.wasm_global_delete(global)
}

func WASMGlobalNew(
	compartment *WASMCompartmentT,
	globalType *WASMGlobalTypeT,
	val *WASMValT,
	debugName string,
) *C.wasm_global_t {
	var cstr *C.char
	if debugName == "" {
		cstr = _EMPTY
	} else {
		cstr = C.CString(debugName)
		defer C.free(unsafe.Pointer(cstr))
	}
	return C.wasm_global_new((*C.wasm_compartment_t)(compartment), (*C.wasm_globaltype_t)(globalType), (*C.wasm_val_t)(val), cstr)
}

func WASMGlobalType(global *WASMGlobalT) *WASMGlobalTypeT {
	return (*WASMGlobalTypeT)(C.wasm_global_type((*C.wasm_global_t)(global)))
}

func WASMGlobalGet(store *WASMStoreT, global *WASMGlobalT) *WASMValT {
	var out C.wasm_val_t //
	C.wasm_global_get((*C.wasm_store_t)(store), (*C.wasm_global_t)(global), &out)
	return (*WASMValT)(&out)
}

func WASMGlobalSet(global *C.wasm_global_t, val *C.wasm_val_t) {
	//C.wasm_global_set(global, val)
}

// Table Instances

func WASMTableDelete(table *WASMTableT) {
	C.wasm_table_delete((*C.wasm_table_t)(table))
}

func WASMTableNew(
	compartment *WASMCompartmentT,
	tableType *WASMTableTypeT,
	init *WASMRefT,
	debugName string) *WASMTableT {
	cstr := C.CString(debugName)
	defer C.free(unsafe.Pointer(cstr))
	return (*WASMTableT)(C.wasm_table_new((*C.wasm_compartment_t)(compartment), (*C.wasm_tabletype_t)(tableType), (*C.wasm_ref_t)(init), cstr))
}

func WASMTableType(table *WASMTableT) *WASMTableTypeT {
	return (*WASMTableTypeT)(C.wasm_table_type((*C.wasm_table_t)(table)))
}

func WASMTableGet(table *WASMTableT, index int) *WASMRefT {
	return (*WASMRefT)(C.wasm_table_get((*C.wasm_table_t)(table), (C.wasm_table_size_t)(index)))
}

func WASMTableSet(table *WASMTableT, index int, value *WASMRefT) bool {
	return bool(C.wasm_table_set((*C.wasm_table_t)(table), (C.wasm_table_size_t)(index), (*C.wasm_ref_t)(value)))
}

func WASMTableSize(table *WASMTableT) int {
	return int(C.wasm_table_size((*C.wasm_table_t)(table)))
}

func WASMTableGrow(
	table *WASMTableT,
	delta int,
	init *WASMRefT,
) (bool, int) {
	var previousSize C.wasm_table_size_t
	ok := bool(C.wasm_table_grow((*C.wasm_table_t)(table), (C.wasm_table_size_t)(delta), (*C.wasm_ref_t)(init), &previousSize))
	return ok, int(previousSize)
}

// Memory Instances

func WASMMemoryDelete(memory *WASMMemoryT) {
	C.wasm_memory_delete((*C.wasm_memory_t)(memory))
}

func WASMMemoryNew(compartment *C.wasm_compartment_t, memoryType *C.wasm_memorytype_t, debugName string) *C.wasm_memory_t {
	cstr := C.CString(debugName)
	defer C.free(unsafe.Pointer(cstr))
	return C.wasm_memory_new(compartment, memoryType, cstr)
}

func WASMMemoryData(memory *WASMMemoryT) *C.char {
	return C.wasm_memory_data((*C.wasm_memory_t)(memory))
}

func WASMMemoryDataSize(memory *WASMMemoryT) uintptr {
	return uintptr(C.wasm_memory_data_size((*C.wasm_memory_t)(memory)))
}

func WASMMemorySize(memory *WASMMemoryT) int {
	return int(C.wasm_memory_size((*C.wasm_memory_t)(memory)))
}

func WASMMemoryGrow(
	memory *WASMMemoryT,
	delta int,
) (bool, int) {
	var previousSize C.wasm_memory_pages_t
	ok := bool(C.wasm_memory_grow((*C.wasm_memory_t)(memory), (C.wasm_memory_pages_t)(delta), &previousSize))
	return ok, int(previousSize)
}

// Externals

func WASMExternDelete(extern *WASMExternT) {
	C.wasm_extern_delete((*C.wasm_extern_t)(extern))
}

func WASMExternKind(extern *WASMExternT) C.wasm_externkind_t {
	return C.wasm_extern_kind((*C.wasm_extern_t)(extern))
}

func WASMExternType(extern *WASMExternT) *C.wasm_externtype_t {
	return C.wasm_extern_type((*C.wasm_extern_t)(extern))
}

func WASMFuncAsExtern(fn *WASMFuncT) *WASMExternT {
	return (*WASMExternT)(C.wasm_func_as_extern((*C.wasm_func_t)(fn)))
}
func WASMGlobalAsExtern(global *WASMGlobalT) *WASMExternT {
	return (*WASMExternT)(C.wasm_global_as_extern((*C.wasm_global_t)(global)))
}
func WASMTableAsExtern(table *WASMTableT) *WASMExternT {
	return (*WASMExternT)(C.wasm_table_as_extern((*C.wasm_table_t)(table)))
}
func WASMMemoryAsExtern(memory *WASMMemoryT) *WASMExternT {
	return (*WASMExternT)(C.wasm_memory_as_extern((*C.wasm_memory_t)(memory)))
}
func WASMExternAsFunc(extern *WASMExternT) *WASMFuncT {
	return (*WASMFuncT)(C.wasm_extern_as_func((*C.wasm_extern_t)(extern)))
}
func WASMExternAsGlobal(extern *WASMExternT) *WASMGlobalT {
	return (*WASMGlobalT)(C.wasm_extern_as_global((*C.wasm_extern_t)(extern)))
}
func WASMExternAsTable(extern *WASMExternT) *WASMTableT {
	return (*WASMTableT)(C.wasm_extern_as_table((*C.wasm_extern_t)(extern)))
}
func WASMExternAsMemory(extern *WASMExternT) *WASMMemoryT {
	return (*WASMMemoryT)(C.wasm_extern_as_memory((*C.wasm_extern_t)(extern)))
}

func WASMFuncAsExternConst(fn *WASMFuncT) *WASMExternT {
	return (*WASMExternT)(C.wasm_func_as_extern_const((*C.wasm_func_t)(fn)))
}
func WASMGlobalAsExternConst(global *WASMGlobalT) *WASMExternT {
	return (*WASMExternT)(C.wasm_global_as_extern_const((*C.wasm_global_t)(global)))
}
func WASMTableAsExternConst(table *WASMTableT) *WASMExternT {
	return (*WASMExternT)(C.wasm_table_as_extern_const((*C.wasm_table_t)(table)))
}
func WASMMemoryAsExternConst(memory *WASMMemoryT) *WASMExternT {
	return (*WASMExternT)(C.wasm_memory_as_extern_const((*C.wasm_memory_t)(memory)))
}
func WASMExternAsFuncConst(extern *WASMExternT) *WASMFuncT {
	return (*WASMFuncT)(C.wasm_extern_as_func_const((*C.wasm_extern_t)(extern)))
}
func WASMExternAsGlobalConst(extern *WASMExternT) *WASMGlobalT {
	return (*WASMGlobalT)(C.wasm_extern_as_global_const((*C.wasm_extern_t)(extern)))
}
func WASMExternAsTableConst(extern *WASMExternT) *WASMTableT {
	return (*WASMTableT)(C.wasm_extern_as_table_const((*C.wasm_extern_t)(extern)))
}
func WASMExternAsMemoryConst(extern *WASMExternT) *WASMMemoryT {
	return (*WASMMemoryT)(C.wasm_extern_as_memory_const((*C.wasm_extern_t)(extern)))
}

// Convenience

func WASMValTypeNew(kind C.wasm_valkind_t) *WASMValTypeT {
	return (*WASMValTypeT)(C.wasm_valtype_new(kind))
}

func WASMValTypeDelete(t *WASMValTypeT) {
	C.wasm_valtype_delete((*C.wasm_valtype_t)(t))
}

func WASMValTypeNewI32() *WASMValTypeT {
	return (*WASMValTypeT)(C.wasm_valtype_new(I32))
}

func WASMValTypeNewI64() *WASMValTypeT {
	return (*WASMValTypeT)(C.wasm_valtype_new(I64))
}

func WASMValTypeNewF32() *WASMValTypeT {
	return (*WASMValTypeT)(C.wasm_valtype_new(F32))
}

func WASMValTypeNewF64() *WASMValTypeT {
	return (*WASMValTypeT)(C.wasm_valtype_new(F64))
}

func WASMValTypeNewV128() *WASMValTypeT {
	return (*WASMValTypeT)(C.wasm_valtype_new(V128))
}

func WASMValTypeNewAnyref() *WASMValTypeT {
	return (*WASMValTypeT)(C.wasm_valtype_new(ANYREF))
}

func WASMValTypeNewFuncref() *WASMValTypeT {
	return (*WASMValTypeT)(C.wasm_valtype_new(FUNCREF))
}

//
//func wasm_functype_new_0_0() *C.wasm_functype_t {
//	return WASMFuncTypeNew(nil, nil)
//}
//
//func WASMFuncTypeNew_1_0(p *C.wasm_valtype_t) *C.wasm_functype_t {
//	return WASMFuncTypeNew([]*C.wasm_valtype_t{p}, nil)
//}

func WASMFuncTypeNew_2_0(
	p1 *WASMValTypeT,
	p2 *WASMValTypeT,
) *WASMFuncTypeT {
	return WASMFuncTypeNew([]*WASMValTypeT{p1, p2}, nil)
}

func WASMFuncTypeNew_3_0(
	p1 *WASMValTypeT,
	p2 *WASMValTypeT,
	p3 *WASMValTypeT,
) *WASMFuncTypeT {
	return WASMFuncTypeNew([]*WASMValTypeT{p1, p2, p3}, nil)
}

func WASMFuncTypeNew_0_1(r *WASMValTypeT) *WASMFuncTypeT {
	return WASMFuncTypeNew(nil, []*WASMValTypeT{r})
}

func WASMFuncTypeNew_1_0(
	p *WASMValTypeT,
) *WASMFuncTypeT {
	return WASMFuncTypeNew([]*WASMValTypeT{p}, []*WASMValTypeT{})
}

func WASMFuncTypeNew_1_1(
	p *WASMValTypeT,
	r *WASMValTypeT,
) *WASMFuncTypeT {
	return WASMFuncTypeNew([]*WASMValTypeT{p}, []*WASMValTypeT{r})
}

func WASMFuncTypeNew_2_1(
	p1 *WASMValTypeT,
	p2 *WASMValTypeT,
	r *WASMValTypeT,
) *WASMFuncTypeT {
	return WASMFuncTypeNew([]*WASMValTypeT{p1, p2}, []*WASMValTypeT{r})
}

func WASMFuncTypeNew_3_1(
	p1 *WASMValTypeT,
	p2 *WASMValTypeT,
	p3 *WASMValTypeT,
	r *WASMValTypeT,
) *WASMFuncTypeT {
	return WASMFuncTypeNew([]*WASMValTypeT{p1, p2, p3}, []*WASMValTypeT{r})
}

func WASMFuncTypeNew_0_2(
	r1 *WASMValTypeT,
	r2 *WASMValTypeT,
) *WASMFuncTypeT {
	return WASMFuncTypeNew(nil, []*WASMValTypeT{r1, r2})
}

func WASMFuncTypeNew_1_2(
	p *WASMValTypeT,
	r1 *WASMValTypeT,
	r2 *WASMValTypeT,
) *WASMFuncTypeT {
	return WASMFuncTypeNew([]*WASMValTypeT{p}, []*WASMValTypeT{r1, r2})
}

func WASMFuncTypeNew_2_2(
	p1 *WASMValTypeT,
	p2 *WASMValTypeT,
	r1 *WASMValTypeT,
	r2 *WASMValTypeT,
) *WASMFuncTypeT {
	return WASMFuncTypeNew([]*WASMValTypeT{p1, p2}, []*WASMValTypeT{r1, r2})
}

func WASMFuncTypeNew_3_2(
	p1 *WASMValTypeT,
	p2 *WASMValTypeT,
	p3 *WASMValTypeT,
	r1 *WASMValTypeT,
	r2 *WASMValTypeT,
) *WASMFuncTypeT {
	return WASMFuncTypeNew([]*WASMValTypeT{p1, p2, p3}, []*WASMValTypeT{r1, r2})
}

func WASMFuncTypeNew_4_1(
	p1 *WASMValTypeT,
	p2 *WASMValTypeT,
	p3 *WASMValTypeT,
	p4 *WASMValTypeT,
	r1 *WASMValTypeT,
) *WASMFuncTypeT {
	return WASMFuncTypeNew([]*WASMValTypeT{p1, p2, p3, p4}, []*WASMValTypeT{r1})
}
