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

wasm_trap_t* wasm_func_call_no_copy(wasm_store_t*,
										   const wasm_func_t*,
										   const wasm_val_t args[],
										   wasm_val_t results[]);

wasm_trap_t* wasm_func_call_no_trap(wasm_store_t*,
										   const wasm_func_t*,
										   const wasm_val_t args[],
										   wasm_val_t results[]);

wasm_trap_t* wasm_func_call_no_copy_no_trap(wasm_store_t*,
										   const wasm_func_t*,
										   const wasm_val_t args[],
										   wasm_val_t results[]);

wasm_module_t* wasm_module_precompiled_new(wasm_engine_t*,
											  const char* binary,
											  size_t num_binary_bytes);

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

int32_t wasm_val_get_i32(wasm_val_t* val) {
	return val->i32;
}


*/
import "C"
import (
	"math"
	"reflect"
	"unsafe"
)

type (
	WASMConfigT       C.wasm_config_t
	WASMEngineT       C.wasm_engine_t
	WASMCompartmentT  C.wasm_compartment_t
	WASMStoreT        C.wasm_store_t
	WASMValTypeT      C.wasm_valtype_t
	WASMFuncTypeT     C.wasm_functype_t
	WASMTableTypeT    C.wasm_tabletype_t
	WASMMemoryTypeT   C.wasm_memorytype_t
	WASMGlobalTypeT   C.wasm_globaltype_t
	WASMExternTypeT   C.wasm_externtype_t
	WASMRefT          C.wasm_ref_t
	WASMTrapT         C.wasm_trap_t
	WASMForeignT      C.wasm_foreign_t
	WASMModuleT       C.wasm_module_t
	WASMFuncT         C.wasm_func_t
	WASMFuncCallbackT C.wasm_func_callback_t
	WASMTableT        C.wasm_table_t
	WASMMemoryT       C.wasm_memory_t
	WASMGlobalT       C.wasm_global_t
	WASMExternT       C.wasm_extern_t
	WASMInstanceT     C.wasm_instance_t
	WASMValT          C.wasm_val_t

	/*
		typedef struct wasm_frame_t
		{
			wasm_func_t* function;
			size_t instr_index;
		} wasm_frame_t;
	*/
	WASMFrameT struct {
		Function   *WASMFuncT
		InstrIndex C.size_t
	}

	/*
		typedef struct wasm_limits_t
		{
			uint32_t min;
			uint32_t max;
		} wasm_limits_t;
	*/
	WASMLimitsT struct {
		Min uint32
		Max uint32
	}

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

type WASMMutabilityT uint8

const (
	CONST = WASMMutabilityT(C.WASM_CONST)
	VAR   = WASMMutabilityT(C.WASM_VAR)
)

type WASMSharedT uint8

const (
	NOTSHARED = WASMSharedT(C.WASM_NOTSHARED)
	SHARED    = WASMSharedT(C.WASM_SHARED)
)

type WASMIndexT uint8

const (
	INDEX_I32 = WASMIndexT(C.WASM_INDEX_I32)
	INDEX_I64 = WASMIndexT(C.WASM_INDEX_I64)
)

type WASMValKindT uint8

const (
	I32     = WASMValKindT(C.WASM_I32)
	I64     = WASMValKindT(C.WASM_I64)
	F32     = WASMValKindT(C.WASM_F32)
	F64     = WASMValKindT(C.WASM_F64)
	V128    = WASMValKindT(C.WASM_V128)
	ANYREF  = WASMValKindT(C.WASM_ANYREF)
	FUNCREF = WASMValKindT(C.WASM_FUNCREF)
)

func (k WASMValKindT) IsNum() bool {
	return k < ANYREF
}
func (k WASMValKindT) IsRef() bool {
	return k >= ANYREF
}

const (
	WASM_MEMORY_PAGE_SIZE        = 0x10000
	WASM_MEMORY_PAGES_MAX        = math.MaxUint32
	WASM_TABLE_SIZE_MAX          = math.MaxUint32
	LIMITS_MAX_DEFAULT    uint32 = 0xffffffff
)

var _EMPTY = C.CString("")

func (v *WASMValT) I32() int32 {
	return *(*int32)(unsafe.Pointer(v))
}
func WASMValGetI32(val *WASMValT) int32 {
	return int32(C.wasm_val_get_i32((*C.wasm_val_t)(val)))
}
func (v *WASMValT) SetI32(value int32) {
	*(*int32)(unsafe.Pointer(v)) = value
}

func (v *WASMValT) I64() int64 {
	return *(*int64)(unsafe.Pointer(v))
}
func (v *WASMValT) SetI64(value int64) {
	*(*int64)(unsafe.Pointer(v)) = value
}

func (v *WASMValT) F32() float32 {
	return *(*float32)(unsafe.Pointer(v))
}
func (v *WASMValT) SetF32(value float32) {
	*(*float32)(unsafe.Pointer(v)) = value
}

func (v *WASMValT) F64() float64 {
	return *(*float64)(unsafe.Pointer(v))
}
func (v *WASMValT) SetF64(value float64) {
	*(*float64)(unsafe.Pointer(v)) = value
}

func (v *WASMValT) V128() complex128 {
	return *(*complex128)(unsafe.Pointer(v))
}
func (v *WASMValT) SetV128(value complex128) {
	*(*complex128)(unsafe.Pointer(v)) = value
}

func (v *WASMValT) Ref32() uint32 {
	return *(*uint32)(unsafe.Pointer(v))
}
func (v *WASMValT) SetRef32(value uint32) {
	*(*uint32)(unsafe.Pointer(v)) = value
}

func (v *WASMValT) Ref64() uint64 {
	return *(*uint64)(unsafe.Pointer(v))
}
func (v *WASMValT) SetRef64(value uint64) {
	*(*uint64)(unsafe.Pointer(v)) = value
}

//func (v *WASMValTypeT) IsNum() bool {
//	return
//}

/////////////////////////////////////////////////////////////////////////////
// wasm_config_t
/////////////////////////////////////////////////////////////////////////////

func WASMConfigNew() *WASMConfigT {
	return (*WASMConfigT)(C.wasm_config_new())
}

func WASMConfigDelete(config *WASMConfigT) {
	C.wasm_config_delete((*C.wasm_config_t)(config))
}

func (c *WASMConfigT) Close() error {
	WASMConfigDelete(c)
	return nil
}

func WASMConfigFeatureSetImportExportMutableGlobals(config *WASMConfigT, enable bool) {
	C.wasm_config_feature_set_import_export_mutable_globals((*C.wasm_config_t)(config), (C.bool)(enable))
}

func (c *WASMConfigT) SetImportExportMutableGlobals(enable bool) *WASMConfigT {
	WASMConfigFeatureSetImportExportMutableGlobals(c, enable)
	return c
}

func WASMConfigFeatureSetNonTrappingFloatToInt(config *WASMConfigT, enable bool) {
	C.wasm_config_feature_set_nontrapping_float_to_int((*C.wasm_config_t)(config), (C.bool)(enable))
}

func (c *WASMConfigT) SetNonTrappingFloatToInt(enable bool) *WASMConfigT {
	WASMConfigFeatureSetNonTrappingFloatToInt(c, enable)
	return c
}

func WASMConfigFeatureSetSignExtension(config *WASMConfigT, enable bool) {
	C.wasm_config_feature_set_sign_extension((*C.wasm_config_t)(config), (C.bool)(enable))
}

func (c *WASMConfigT) SetSignExtension(enable bool) *WASMConfigT {
	WASMConfigFeatureSetSignExtension(c, enable)
	return c
}

func WASMConfigFeatureSetBulkMemoryOps(config *WASMConfigT, enable bool) {
	C.wasm_config_feature_set_bulk_memory_ops((*C.wasm_config_t)(config), (C.bool)(enable))
}

func (c *WASMConfigT) SetBulkMemoryOps(enable bool) *WASMConfigT {
	WASMConfigFeatureSetBulkMemoryOps(c, enable)
	return c
}

func WASMConfigFeatureSetSIMD(config *WASMConfigT, enable bool) {
	C.wasm_config_feature_set_simd((*C.wasm_config_t)(config), (C.bool)(enable))
}

func (c *WASMConfigT) SetSIMD(enable bool) *WASMConfigT {
	WASMConfigFeatureSetSIMD(c, enable)
	return c
}

func WASMConfigFeatureSetAtomics(config *WASMConfigT, enable bool) {
	C.wasm_config_feature_set_atomics((*C.wasm_config_t)(config), (C.bool)(enable))
}

func (c *WASMConfigT) SetAtomics(enable bool) *WASMConfigT {
	WASMConfigFeatureSetAtomics(c, enable)
	return c
}

func WASMConfigFeatureSetExceptionHandling(config *WASMConfigT, enable bool) {
	C.wasm_config_feature_set_exception_handling((*C.wasm_config_t)(config), (C.bool)(enable))
}

func (c *WASMConfigT) SetExceptionHandling(enable bool) *WASMConfigT {
	WASMConfigFeatureSetExceptionHandling(c, enable)
	return c
}

func WASMConfigFeatureSetMultiValue(config *WASMConfigT, enable bool) {
	C.wasm_config_feature_set_multivalue((*C.wasm_config_t)(config), (C.bool)(enable))
}

func (c *WASMConfigT) SetMultiValue(enable bool) *WASMConfigT {
	WASMConfigFeatureSetMultiValue(c, enable)
	return c
}

func WASMConfigFeatureSetReferenceTypes(config *WASMConfigT, enable bool) {
	C.wasm_config_feature_set_reference_types((*C.wasm_config_t)(config), (C.bool)(enable))
}

func (c *WASMConfigT) SetReferenceTypes(enable bool) *WASMConfigT {
	WASMConfigFeatureSetReferenceTypes(c, enable)
	return c
}

func WASMConfigFeatureSetExtendedNameSection(config *WASMConfigT, enable bool) {
	C.wasm_config_feature_set_extended_name_section((*C.wasm_config_t)(config), (C.bool)(enable))
}

func (c *WASMConfigT) SetExtendedNameSection(enable bool) *WASMConfigT {
	WASMConfigFeatureSetExtendedNameSection(c, enable)
	return c
}

func WASMConfigFeatureSetMultiMemory(config *WASMConfigT, enable bool) {
	C.wasm_config_feature_set_multimemory((*C.wasm_config_t)(config), (C.bool)(enable))
}

func (c *WASMConfigT) SetMultiMemory(enable bool) *WASMConfigT {
	WASMConfigFeatureSetMultiMemory(c, enable)
	return c
}

/////////////////////////////////////////////////////////////////////////////
// wasm_engine_t
/////////////////////////////////////////////////////////////////////////////

func WASMEngineNew() *WASMEngineT {
	return (*WASMEngineT)(C.wasm_engine_new())
}

func WASMEngineDelete(engine *WASMEngineT) {
	C.wasm_engine_delete((*C.wasm_engine_t)(engine))
}

func (e *WASMEngineT) Close() error {
	WASMEngineDelete(e)
	return nil
}

func (e *WASMEngineT) Delete() {
	WASMEngineDelete(e)
}

func WASMEngineNewWithConfig(config *WASMConfigT) *WASMEngineT {
	return (*WASMEngineT)(C.wasm_engine_new_with_config((*C.wasm_config_t)(config)))
}

func (e *WASMEngineT) NewCompartment(debugName string) *WASMCompartmentT {
	return WASMCompartmentNew(e, debugName)
}

func (e *WASMEngineT) NewModule(binary []byte) *WASMModuleT {
	return WASMModuleNew(e, binary)
}

func (e *WASMEngineT) NewPrecompiledModule(binary []byte) *WASMModuleT {
	return WASMModulePrecompiledNew(e, binary)
}

/////////////////////////////////////////////////////////////////////////////
// wasm_compartment_t
/////////////////////////////////////////////////////////////////////////////

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

func WASMCompartmentDelete(compartment *WASMCompartmentT) {
	C.wasm_compartment_delete((*C.wasm_compartment_t)(compartment))
}

func (c *WASMCompartmentT) Close() error {
	WASMCompartmentDelete(c)
	return nil
}

func (c *WASMCompartmentT) Delete() {
	WASMCompartmentDelete(c)
}

func WASMCompartmentClone(compartment *WASMCompartmentT) *WASMCompartmentT {
	return (*WASMCompartmentT)(C.wasm_compartment_clone((*C.wasm_compartment_t)(compartment)))
}

func (c *WASMCompartmentT) Clone() *WASMCompartmentT {
	return WASMCompartmentClone(c)
}

func WASMCompartmentContains(compartment *WASMCompartmentT, ref *WASMRefT) bool {
	return bool(C.wasm_compartment_contains((*C.wasm_compartment_t)(compartment), (*C.wasm_ref_t)(ref)))
}

func (c *WASMCompartmentT) Contains(ref *WASMRefT) bool {
	return WASMCompartmentContains(c, ref)
}

func (c *WASMCompartmentT) NewStore(debugName string) *WASMStoreT {
	return WASMStoreNew(c, debugName)
}

/////////////////////////////////////////////////////////////////////////////
// wasm_store_t
/////////////////////////////////////////////////////////////////////////////

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

func (s *WASMStoreT) Close() error {
	WASMStoreDelete(s)
	return nil
}

func (s *WASMStoreT) Delete() {
	WASMStoreDelete(s)
}

/////////////////////////////////////////////////////////////////////////////
// wasm_instance_t
/////////////////////////////////////////////////////////////////////////////

func WASMInstanceNew(
	store *WASMStoreT,
	module *WASMModuleT,
	imports []*WASMExternT,
	outTrap **WASMTrapT,
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
		(**C.wasm_trap_t)(unsafe.Pointer(outTrap)),
		cstr,
	))
}

func WASMInstanceDelete(instance *WASMInstanceT) {
	C.wasm_instance_delete((*C.wasm_instance_t)(instance))
}

func (inst *WASMInstanceT) Close() error {
	WASMInstanceDelete(inst)
	return nil
}
func (inst *WASMInstanceT) Delete() {
	WASMInstanceDelete(inst)
}

func WASMInstanceNumExports(inst *WASMInstanceT) int {
	return int(C.wasm_instance_num_exports((*C.wasm_instance_t)(inst)))
}

func (inst *WASMInstanceT) NumExports() int {
	return WASMInstanceNumExports(inst)
}

func WASMInstanceExport(inst *WASMInstanceT, index int) *WASMExternT {
	return (*WASMExternT)(C.wasm_instance_export((*C.wasm_instance_t)(inst), C.size_t(index)))
}

func (inst *WASMInstanceT) Export(index int) *WASMExternT {
	return WASMInstanceExport(inst, index)
}

/////////////////////////////////////////////////////////////////////////////
// wasm_functype_t
/////////////////////////////////////////////////////////////////////////////

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

func WASMFuncTypeDelete(t *WASMFuncTypeT) {
	C.wasm_functype_delete((*C.wasm_functype_t)(t))
}

func (f *WASMFuncTypeT) Close() error {
	WASMFuncTypeDelete(f)
	return nil
}

func WASMFuncTypeNumParams(funcType *WASMFuncTypeT) int {
	return int(C.wasm_functype_num_params((*C.wasm_functype_t)(funcType)))
}

func (f *WASMFuncTypeT) NumParams() int {
	return WASMFuncTypeNumParams(f)
}

func WASMFuncTypeParam(funcType *WASMFuncTypeT, index int) *WASMValTypeT {
	return (*WASMValTypeT)(C.wasm_functype_param((*C.wasm_functype_t)(funcType), (C.size_t)(index)))
}

func (f *WASMFuncTypeT) Param(index int) *WASMValTypeT {
	return WASMFuncTypeParam(f, index)
}

func WASMFuncTypeNumResults(funcType *WASMFuncTypeT) int {
	return int(C.wasm_functype_num_results((*C.wasm_functype_t)(funcType)))
}

func (f *WASMFuncTypeT) NumResults() int {
	return WASMFuncTypeNumResults(f)
}

func WASMFuncTypeResult(funcType *WASMFuncTypeT, index int) *WASMValTypeT {
	return (*WASMValTypeT)(C.wasm_functype_result((*C.wasm_functype_t)(funcType), (C.size_t)(index)))
}

func (f *WASMFuncTypeT) Result(index int) *WASMValTypeT {
	return WASMFuncTypeResult(f, index)
}

/////////////////////////////////////////////////////////////////////////////
// wasm_globaltype_t
/////////////////////////////////////////////////////////////////////////////

func WASMGlobalTypeNew(valType *WASMValTypeT, mutability WASMMutabilityT) *WASMGlobalTypeT {
	return (*WASMGlobalTypeT)(C.wasm_globaltype_new(
		(*C.wasm_valtype_t)(valType),
		(C.wasm_mutability_t)(mutability),
	))
}

func WASMGlobalTypeDelete(globalType *WASMGlobalTypeT) {
	C.wasm_globaltype_delete((*C.wasm_globaltype_t)(globalType))
}

func (g *WASMGlobalTypeT) Close() error {
	if g == nil {
		return nil
	}
	WASMGlobalTypeDelete(g)
	return nil
}

func WASMGlobalTypeContent(globalType *WASMGlobalTypeT) *WASMValTypeT {
	return (*WASMValTypeT)(C.wasm_globaltype_content((*C.wasm_globaltype_t)(globalType)))
}

func (g *WASMGlobalTypeT) Content() *WASMValTypeT {
	return WASMGlobalTypeContent(g)
}

func WASMGlobalTypeMutability(globalType *WASMGlobalTypeT) WASMMutabilityT {
	return (WASMMutabilityT)(C.wasm_globaltype_mutability((*C.wasm_globaltype_t)(globalType)))
}

func (g *WASMGlobalTypeT) Mutability() WASMMutabilityT {
	return WASMGlobalTypeMutability(g)
}

/////////////////////////////////////////////////////////////////////////////
// wasm_tabletype_t
/////////////////////////////////////////////////////////////////////////////

func WASMTableTypeNew(t *WASMValTypeT, limits *WASMLimitsT, shared WASMSharedT, index int) *WASMTableTypeT {
	return (*WASMTableTypeT)(C.wasm_tabletype_new(
		(*C.wasm_valtype_t)(t),
		(*C.wasm_limits_t)(unsafe.Pointer(limits)),
		(C.wasm_shared_t)(shared),
		(C.wasm_index_t)(index),
	))
}

func WASMTableTypeDelete(t *WASMTableTypeT) {
	C.wasm_tabletype_delete((*C.wasm_tabletype_t)(t))
}

func (t *WASMTableTypeT) Close() error {
	WASMTableTypeDelete(t)
	return nil
}

func WASMTableTypeElement(tableType *WASMTableTypeT) *WASMValTypeT {
	return (*WASMValTypeT)(C.wasm_tabletype_element((*C.wasm_tabletype_t)(tableType)))
}

func (t *WASMTableTypeT) Element() *WASMValTypeT {
	return WASMTableTypeElement(t)
}

func WASMTableTypeLimits(tableType *WASMTableTypeT) *WASMLimitsT {
	return (*WASMLimitsT)(unsafe.Pointer(C.wasm_tabletype_limits((*C.wasm_tabletype_t)(tableType))))
}

func (t *WASMTableTypeT) Limits() *WASMLimitsT {
	return WASMTableTypeLimits(t)
}

func WASMTableTypeShared(tableType *WASMTableTypeT) WASMSharedT {
	return (WASMSharedT)(C.wasm_tabletype_shared((*C.wasm_tabletype_t)(tableType)))
}

func (t *WASMTableTypeT) Shared() WASMSharedT {
	return WASMTableTypeShared(t)
}

// Memory Types

func WASMMemoryTypeNew(limits *WASMLimitsT, shared WASMSharedT, index int) *WASMMemoryTypeT {
	return (*WASMMemoryTypeT)(C.wasm_memorytype_new((*C.wasm_limits_t)(unsafe.Pointer(limits)), (C.wasm_shared_t)(shared), (C.wasm_index_t)(index)))
}

func WASMMemoryTypeLimits(memoryType *WASMMemoryTypeT) *WASMLimitsT {
	return (*WASMLimitsT)(unsafe.Pointer(C.wasm_memorytype_limits((*C.wasm_memorytype_t)(memoryType))))
}

func (m *WASMMemoryTypeT) Limits() *WASMLimitsT {
	return WASMMemoryTypeLimits(m)
}

func WASMMemoryTypeShared(memoryType *WASMMemoryTypeT) WASMSharedT {
	return (WASMSharedT)(C.wasm_memorytype_shared((*C.wasm_memorytype_t)(memoryType)))
}

func (m *WASMMemoryTypeT) Shared() WASMSharedT {
	return WASMMemoryTypeShared(m)
}

/////////////////////////////////////////////////////////////////////////////
// wasm_import_t
/////////////////////////////////////////////////////////////////////////////

func (w *WASMImportT) Module() string {
	return string(*(*string)(unsafe.Pointer(&reflect.StringHeader{
		Data: uintptr(unsafe.Pointer(w.module)),
		Len:  int(w.num_module_bytes),
	})))
}

func (w *WASMImportT) ModuleUnsafe() string {
	return *(*string)(unsafe.Pointer(&reflect.StringHeader{
		Data: uintptr(unsafe.Pointer(w.module)),
		Len:  int(w.num_module_bytes),
	}))
}

func (w *WASMImportT) Name() string {
	return string(w.NameUnsafe())
}

func (w *WASMImportT) NameUnsafe() string {
	return *(*string)(unsafe.Pointer(&reflect.StringHeader{
		Data: uintptr(unsafe.Pointer(w.name)),
		Len:  int(w.num_name_bytes),
	}))
}

/////////////////////////////////////////////////////////////////////////////
// wasm_export_t
/////////////////////////////////////////////////////////////////////////////

func (w *WASMExportT) NameUnsafe() string {
	return *(*string)(unsafe.Pointer(&reflect.StringHeader{
		Data: uintptr(unsafe.Pointer(w.name)),
		Len:  int(w.num_name_bytes),
	}))
}

/////////////////////////////////////////////////////////////////////////////
// wasm_val_t
/////////////////////////////////////////////////////////////////////////////

func WASMValDelete(kind C.wasm_valkind_t, v *C.wasm_val_t) {
	C.wasm_val_delete(kind, v)
}

//func wasm_val_copy(kind C.wasm_valkind_t, out *C.wasm_val_t, v *C.wasm_val_t)

/////////////////////////////////////////////////////////////////////////////
// wasm_trap_t
/////////////////////////////////////////////////////////////////////////////

func WASMTrapNew(compartment *WASMCompartmentT, message string) *WASMTrapT {
	msgbytes := []byte(message)
	return (*WASMTrapT)(C.wasm_trap_new(
		(*C.wasm_compartment_t)(compartment),
		(*C.char)(unsafe.Pointer(&msgbytes[0])),
		(C.size_t)(len(msgbytes))))
}

func WASMTrapDelete(trap *WASMTrapT) {
	C.wasm_trap_delete((*C.wasm_trap_t)(trap))
}

func (t *WASMTrapT) Close() error {
	WASMTrapDelete(t)
	return nil
}

func WASMTrapMessageString(trap *WASMTrapT) string {
	return string(WASMTrapMessage(trap, make([]byte, 1024)))
}

func (t *WASMTrapT) String() string {
	return WASMTrapMessageString(t)
}

func (t *WASMTrapT) Error() string {
	return t.String()
}

func WASMTrapMessage(trap *WASMTrapT, buffer []byte) []byte {
	if buffer == nil {
		buffer = make([]byte, 1024)
	}
	size := C.size_t(len(buffer))
	C.wasm_trap_message((*C.wasm_trap_t)(trap), (*C.char)(unsafe.Pointer(&buffer[0])), &size)
	if int(size) > len(buffer) {
		buffer = make([]byte, int(size))
		C.wasm_trap_message((*C.wasm_trap_t)(trap), (*C.char)(unsafe.Pointer(&buffer[0])), &size)
	}
	return buffer[0:int(size)]
}

func (t *WASMTrapT) Message(b []byte) []byte {
	return WASMTrapMessage(t, b)
}

func WASMTrapStackNumFrames(trap *WASMTrapT) int {
	return int(C.wasm_trap_stack_num_frames((*C.wasm_trap_t)(trap)))
}

func (t *WASMTrapT) StackNumFrames() int {
	return WASMTrapStackNumFrames(t)
}

func WASMTrapStackFrame(trap *WASMTrapT, index int, outFrame *WASMFrameT) {
	C.wasm_trap_stack_frame((*C.wasm_trap_t)(trap), (C.size_t)(index), (*C.wasm_frame_t)(unsafe.Pointer(outFrame)))
}

func (t *WASMTrapT) StackFrame(index int, outFrame *WASMFrameT) {
	WASMTrapStackFrame(t, index, outFrame)
}

/////////////////////////////////////////////////////////////////////////////
// wasm_foreign_t
/////////////////////////////////////////////////////////////////////////////

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

func WASMForeignDelete(foreign *WASMForeignT) {
	C.wasm_foreign_delete((*C.wasm_foreign_t)(foreign))
}

func (f *WASMForeignT) Close() error {
	WASMForeignDelete(f)
	return nil
}

/////////////////////////////////////////////////////////////////////////////
// wasm_module_t
/////////////////////////////////////////////////////////////////////////////

func WASMModuleNew(engine *WASMEngineT, binary []byte) *WASMModuleT {
	slice := (*reflect.StringHeader)(unsafe.Pointer(&binary))
	return (*WASMModuleT)(C.wasm_module_new((*C.wasm_engine_t)(engine), (*C.char)(unsafe.Pointer(slice.Data)), (C.size_t)(slice.Len)))
}

func WASMModulePrecompiledNew(engine *WASMEngineT, binary []byte) *WASMModuleT {
	slice := (*reflect.StringHeader)(unsafe.Pointer(&binary))
	return (*WASMModuleT)(C.wasm_module_precompiled_new((*C.wasm_engine_t)(engine), (*C.char)(unsafe.Pointer(slice.Data)), (C.size_t)(slice.Len)))
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

func WASMModuleDelete(module *WASMModuleT) {
	C.wasm_module_delete((*C.wasm_module_t)(module))
}

func (m *WASMModuleT) Close() error {
	WASMModuleDelete(m)
	return nil
}
func (m *WASMModuleT) Delete() {
	WASMModuleDelete(m)
}

func WASMModulePrint(module *WASMModuleT) string {
	var out C.size_t
	ptr := C.wasm_module_print((*C.wasm_module_t)(module), &out)
	return C.GoStringN(ptr, (C.int)(out))
}

func (m *WASMModuleT) Print() string {
	return WASMModulePrint(m)
}

func (m *WASMModuleT) PrintTo(b []byte) []byte {
	var out C.size_t
	ptr := C.wasm_module_print((*C.wasm_module_t)(m), &out)
	return append(b, *(*string)(unsafe.Pointer(&reflect.StringHeader{
		Data: uintptr(unsafe.Pointer(ptr)),
		Len:  int(out),
	}))...)
}

func WASMModuleValidate(binary []byte) bool {
	if len(binary) == 0 {
		return false
	}
	return bool(C.wasm_module_validate((*C.char)(unsafe.Pointer(&binary[0])), (C.size_t)(len(binary))))
}

func WASMModuleNumImports(module *WASMModuleT) int {
	return int(C.wasm_module_num_imports((*C.wasm_module_t)(module)))
}

func (m *WASMModuleT) NumImports() int {
	return WASMModuleNumImports(m)
}

func WASMModuleImport(module *WASMModuleT, index int) WASMImportT {
	var out C.wasm_import_t
	C.wasm_module_import((*C.wasm_module_t)(module), (C.size_t)(index), &out)
	return *(*WASMImportT)(unsafe.Pointer(&out))
}

func (m *WASMModuleT) Import(index int) WASMImportT {
	return WASMModuleImport(m, index)
}

func (m *WASMModuleT) Imports(imports []WASMImportT) []WASMImportT {
	count := m.NumImports()
	if len(imports) > 0 {
		imports = imports[:0]
	}
	for i := 0; i < count; i++ {
		imports = append(imports, m.Import(i))
	}
	return imports
}

func WASMModuleNumExports(module *WASMModuleT) int {
	return int(C.wasm_module_num_exports((*C.wasm_module_t)(module)))
}

func (m *WASMModuleT) NumExports() int {
	return WASMModuleNumExports(m)
}

func WASMModuleExport(module *WASMModuleT, index int) WASMExportT {
	var out C.wasm_export_t
	C.wasm_module_export((*C.wasm_module_t)(module), (C.size_t)(index), &out)
	return *(*WASMExportT)(unsafe.Pointer(&out))
}

func (m *WASMModuleT) Export(index int) WASMExportT {
	return WASMModuleExport(m, index)
}

func (m *WASMModuleT) Exports(exports []WASMExportT) []WASMExportT {
	count := m.NumExports()
	if len(exports) > 0 {
		exports = exports[:0]
	}
	for i := 0; i < count; i++ {
		exports = append(exports, m.Export(i))
	}
	return exports
}

/////////////////////////////////////////////////////////////////////////////
// wasm_func_t
/////////////////////////////////////////////////////////////////////////////

func WASMFuncNew(
	compartment *WASMCompartmentT,
	funcType *WASMFuncTypeT,
	callback WASMFuncCallbackT,
	debugName string) *WASMFuncT {
	var cstr *C.char
	if debugName == "" {
		cstr = _EMPTY
	} else {
		cstr = C.CString(debugName)
		defer C.free(unsafe.Pointer(cstr))
	}
	return (*WASMFuncT)(C.wasm_func_new(
		(*C.wasm_compartment_t)(compartment),
		(*C.wasm_functype_t)(funcType),
		(C.wasm_func_callback_t)(callback),
		cstr,
	))
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

func WASMFuncDelete(fn *WASMFuncT) {
	C.wasm_func_delete((*C.wasm_func_t)(fn))
}

func (f *WASMFuncT) Close() error {
	WASMFuncDelete(f)
	return nil
}

func (f *WASMFuncT) Delete() {
	WASMFuncDelete(f)
}

func WASMFuncType(fn *WASMFuncT) *WASMFuncTypeT {
	return (*WASMFuncTypeT)(C.wasm_func_type((*C.wasm_func_t)(fn)))
}

func (f *WASMFuncT) Type() *WASMFuncTypeT {
	return WASMFuncType(f)
}

func WASMFuncParamArity(fn *WASMFuncT) int {
	return int(C.wasm_func_param_arity((*C.wasm_func_t)(fn)))
}

func (f *WASMFuncT) ParamArity() int {
	return WASMFuncParamArity(f)
}

func WASMFuncResultArity(fn *WASMFuncT) int {
	return int(C.wasm_func_result_arity((*C.wasm_func_t)(fn)))
}

func (f *WASMFuncT) ResultArity() int {
	return WASMFuncResultArity(f)
}

func WASMFuncCall(store *WASMStoreT, fn *WASMFuncT, args *WASMValT, results *WASMValT) *WASMTrapT {
	return (*WASMTrapT)(C.wasm_func_call(
		(*C.wasm_store_t)(store),
		(*C.wasm_func_t)(fn),
		(*C.wasm_val_t)(args),
		(*C.wasm_val_t)(results),
	))
}

func (f *WASMFuncT) Call(store *WASMStoreT, args *WASMValT, results *WASMValT) *WASMTrapT {
	return (*WASMTrapT)(C.wasm_func_call(
		(*C.wasm_store_t)(store),
		(*C.wasm_func_t)(f),
		(*C.wasm_val_t)(args),
		(*C.wasm_val_t)(results),
	))
}

func WASMFuncCallNoTrap(store *WASMStoreT, fn *WASMFuncT, args *WASMValT, results *WASMValT) *WASMTrapT {
	return (*WASMTrapT)(C.wasm_func_call_no_trap(
		(*C.wasm_store_t)(store),
		(*C.wasm_func_t)(fn),
		(*C.wasm_val_t)(args),
		(*C.wasm_val_t)(results),
	))
}

func (f *WASMFuncT) CallNoTrap(store *WASMStoreT, args *WASMValT, results *WASMValT) *WASMTrapT {
	return (*WASMTrapT)(C.wasm_func_call_no_trap(
		(*C.wasm_store_t)(store),
		(*C.wasm_func_t)(f),
		(*C.wasm_val_t)(args),
		(*C.wasm_val_t)(results),
	))
}

func WASMFuncCallNoCopy(store *WASMStoreT, fn *WASMFuncT, args *WASMValT, results *WASMValT) *WASMTrapT {
	return (*WASMTrapT)(C.wasm_func_call_no_copy(
		(*C.wasm_store_t)(store),
		(*C.wasm_func_t)(fn),
		(*C.wasm_val_t)(args),
		(*C.wasm_val_t)(results),
	))
}

func (f *WASMFuncT) CallNoCopy(store *WASMStoreT, args *WASMValT, results *WASMValT) *WASMTrapT {
	return (*WASMTrapT)(C.wasm_func_call_no_copy(
		(*C.wasm_store_t)(store),
		(*C.wasm_func_t)(f),
		(*C.wasm_val_t)(args),
		(*C.wasm_val_t)(results),
	))
}

func WASMFuncCallNoCopyNoTrap(store *WASMStoreT, fn *WASMFuncT, args *WASMValT, results *WASMValT) *WASMTrapT {
	return (*WASMTrapT)(C.wasm_func_call_no_copy_no_trap(
		(*C.wasm_store_t)(store),
		(*C.wasm_func_t)(fn),
		(*C.wasm_val_t)(args),
		(*C.wasm_val_t)(results),
	))
}

func (f *WASMFuncT) CallUnsafe(store *WASMStoreT, args *WASMValT, results *WASMValT) *WASMTrapT {
	return (*WASMTrapT)(C.wasm_func_call_no_copy_no_trap(
		(*C.wasm_store_t)(store),
		(*C.wasm_func_t)(f),
		(*C.wasm_val_t)(args),
		(*C.wasm_val_t)(results),
	))
}

/////////////////////////////////////////////////////////////////////////////
// wasm_global_t
/////////////////////////////////////////////////////////////////////////////

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

func WASMGlobalDelete(global *WASMGlobalT) {
	C.wasm_global_delete((*C.wasm_global_t)(global))
}

func (g *WASMGlobalT) Close() error {
	WASMGlobalDelete(g)
	return nil
}

func WASMGlobalType(global *WASMGlobalT) *WASMGlobalTypeT {
	return (*WASMGlobalTypeT)(C.wasm_global_type((*C.wasm_global_t)(global)))
}

func (g *WASMGlobalT) Type() *WASMGlobalTypeT {
	return WASMGlobalType(g)
}

func WASMGlobalGet(store *WASMStoreT, global *WASMGlobalT) *WASMValT {
	var out C.wasm_val_t //
	C.wasm_global_get((*C.wasm_store_t)(store), (*C.wasm_global_t)(global), &out)
	return (*WASMValT)(&out)
}

func (g *WASMGlobalT) Get(store *WASMStoreT) *WASMValT {
	return WASMGlobalGet(store, g)
}

//func WASMGlobalSet(global *WASMGlobalT, val *WASMValT) {
//	C.wasm_global_set((*C.wasm_global_t)(global), (*C.wasm_val_t)(val))
//}
//
//func (g *WASMGlobalT) Set(val *WASMValT) {
//	WASMGlobalSet(g, val)
//}

/////////////////////////////////////////////////////////////////////////////
// wasm_table_t
/////////////////////////////////////////////////////////////////////////////

func WASMTableNew(
	compartment *WASMCompartmentT,
	tableType *WASMTableTypeT,
	init *WASMRefT,
	debugName string) *WASMTableT {
	cstr := C.CString(debugName)
	defer C.free(unsafe.Pointer(cstr))
	return (*WASMTableT)(C.wasm_table_new((*C.wasm_compartment_t)(compartment), (*C.wasm_tabletype_t)(tableType), (*C.wasm_ref_t)(init), cstr))
}

func WASMTableDelete(table *WASMTableT) {
	C.wasm_table_delete((*C.wasm_table_t)(table))
}

func (t *WASMTableT) Close() error {
	WASMTableDelete(t)
	return nil
}

func WASMTableType(table *WASMTableT) *WASMTableTypeT {
	return (*WASMTableTypeT)(C.wasm_table_type((*C.wasm_table_t)(table)))
}

func (t *WASMTableT) Type() *WASMTableTypeT {
	return WASMTableType(t)
}

func WASMTableGet(table *WASMTableT, index int) *WASMRefT {
	return (*WASMRefT)(C.wasm_table_get((*C.wasm_table_t)(table), (C.wasm_table_size_t)(index)))
}

func (t *WASMTableT) Get(index int) *WASMRefT {
	return WASMTableGet(t, index)
}

func WASMTableSet(table *WASMTableT, index int, value *WASMRefT) bool {
	return bool(C.wasm_table_set((*C.wasm_table_t)(table), (C.wasm_table_size_t)(index), (*C.wasm_ref_t)(value)))
}

func (t *WASMTableT) Set(index int, value *WASMRefT) bool {
	return WASMTableSet(t, index, value)
}

func WASMTableSize(table *WASMTableT) int {
	return int(C.wasm_table_size((*C.wasm_table_t)(table)))
}

func (t *WASMTableT) Size() int {
	return WASMTableSize(t)
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

func (t *WASMTableT) Grow(delta int, init *WASMRefT) (bool, int) {
	return WASMTableGrow(t, delta, init)
}

/////////////////////////////////////////////////////////////////////////////
// wasm_memory_t
/////////////////////////////////////////////////////////////////////////////

func WASMMemoryNew(compartment *C.wasm_compartment_t, memoryType *C.wasm_memorytype_t, debugName string) *C.wasm_memory_t {
	var cstr *C.char
	if debugName == "" {
		cstr = _EMPTY
	} else {
		cstr = C.CString(debugName)
		defer C.free(unsafe.Pointer(cstr))
	}
	return C.wasm_memory_new(compartment, memoryType, cstr)
}

func WASMMemoryDelete(memory *WASMMemoryT) {
	C.wasm_memory_delete((*C.wasm_memory_t)(memory))
}

func (m *WASMMemoryT) Close() error {
	WASMMemoryDelete(m)
	return nil
}

func WASMMemoryData(memory *WASMMemoryT) *C.char {
	return C.wasm_memory_data((*C.wasm_memory_t)(memory))
}

func (m *WASMMemoryT) Data() *C.char {
	return WASMMemoryData(m)
}

func WASMMemoryDataSize(memory *WASMMemoryT) uintptr {
	return uintptr(C.wasm_memory_data_size((*C.wasm_memory_t)(memory)))
}

func (m *WASMMemoryT) Size() int {
	return int(WASMMemoryDataSize(m))
}

func WASMMemorySize(memory *WASMMemoryT) int {
	return int(C.wasm_memory_size((*C.wasm_memory_t)(memory)))
}

func (m *WASMMemoryT) Pages() int {
	return WASMMemorySize(m)
}

func WASMMemoryGrow(
	memory *WASMMemoryT,
	delta int,
) (bool, int) {
	var previousSize C.wasm_memory_pages_t
	ok := bool(C.wasm_memory_grow((*C.wasm_memory_t)(memory), (C.wasm_memory_pages_t)(delta), &previousSize))
	return ok, int(previousSize)
}

func (m *WASMMemoryT) Grow(delta int) (bool, int) {
	return WASMMemoryGrow(m, delta)
}

/////////////////////////////////////////////////////////////////////////////
// wasm_externkind_t
/////////////////////////////////////////////////////////////////////////////

type WASMExternKindT uint8

const (
	EXTERN_FUNC   = WASMExternKindT(C.WASM_EXTERN_FUNC)
	EXTERN_TABLE  = WASMExternKindT(C.WASM_EXTERN_TABLE)
	EXTERN_MEMORY = WASMExternKindT(C.WASM_EXTERN_MEMORY)
	EXTERN_GLOBAL = WASMExternKindT(C.WASM_EXTERN_GLOBAL)
)

func WASMExternKind(extern *WASMExternT) WASMExternKindT {
	return WASMExternKindT(C.wasm_extern_kind((*C.wasm_extern_t)(extern)))
}

func (e *WASMExternT) AsKind() WASMExternKindT {
	return WASMExternKind(e)
}

/////////////////////////////////////////////////////////////////////////////
// wasm_externtype_t
/////////////////////////////////////////////////////////////////////////////

func WASMExternType(extern *WASMExternT) *WASMExternTypeT {
	return (*WASMExternTypeT)(C.wasm_extern_type((*C.wasm_extern_t)(extern)))
}

func WASMFuncTypeAsExternType(t *WASMFuncTypeT) *WASMExternTypeT {
	return (*WASMExternTypeT)(C.wasm_functype_as_externtype((*C.wasm_functype_t)(t)))
}

func WASMGlobalTypeAsExternType(t *WASMGlobalTypeT) *WASMExternTypeT {
	return (*WASMExternTypeT)(C.wasm_globaltype_as_externtype((*C.wasm_globaltype_t)(t)))
}

func WASMTableTypeAsExternType(t *WASMTableTypeT) *WASMExternTypeT {
	return (*WASMExternTypeT)(C.wasm_tabletype_as_externtype((*C.wasm_tabletype_t)(t)))
}

func WASMMemoryTypeAsExternType(t *WASMMemoryTypeT) *WASMExternTypeT {
	return (*WASMExternTypeT)(C.wasm_memorytype_as_externtype((*C.wasm_memorytype_t)(t)))
}

func WASMExternTypeAsFuncType(t *WASMExternTypeT) *WASMFuncTypeT {
	return (*WASMFuncTypeT)(C.wasm_externtype_as_functype((*C.wasm_externtype_t)(t)))
}

func WASMExternTypeAsGlobalType(t *WASMExternTypeT) *WASMGlobalTypeT {
	return (*WASMGlobalTypeT)(C.wasm_externtype_as_globaltype((*C.wasm_externtype_t)(t)))
}

func WASMExternTypeAsTableType(t *WASMExternTypeT) *WASMTableTypeT {
	return (*WASMTableTypeT)(C.wasm_externtype_as_tabletype((*C.wasm_externtype_t)(t)))
}

func WASMExternTypeAsMemoryType(t *WASMExternTypeT) *WASMMemoryTypeT {
	return (*WASMMemoryTypeT)(C.wasm_externtype_as_memorytype((*C.wasm_externtype_t)(t)))
}

//func WASMExternType(extern *WASMExternT) *WASMExternTypeT {
//	return (*WASMExternTypeT)(C.wasm_extern_type_const((*C.wasm_extern_t)(extern)))
//}

func WASMFuncTypeAsExternTypeConst(t *WASMFuncTypeT) *WASMExternTypeT {
	return (*WASMExternTypeT)(C.wasm_functype_as_externtype_const((*C.wasm_functype_t)(t)))
}

func WASMGlobalTypeAsExternTypeConst(t *WASMGlobalTypeT) *WASMExternTypeT {
	return (*WASMExternTypeT)(C.wasm_globaltype_as_externtype_const((*C.wasm_globaltype_t)(t)))
}

func WASMTableTypeAsExternTypeConst(t *WASMTableTypeT) *WASMExternTypeT {
	return (*WASMExternTypeT)(C.wasm_tabletype_as_externtype_const((*C.wasm_tabletype_t)(t)))
}

func WASMMemoryTypeAsExternTypeConst(t *WASMMemoryTypeT) *WASMExternTypeT {
	return (*WASMExternTypeT)(C.wasm_memorytype_as_externtype_const((*C.wasm_memorytype_t)(t)))
}

func WASMExternTypeAsFuncTypeConst(t *WASMExternTypeT) *WASMFuncTypeT {
	return (*WASMFuncTypeT)(C.wasm_externtype_as_functype_const((*C.wasm_externtype_t)(t)))
}

func WASMExternTypeAsGlobalTypeConst(t *WASMExternTypeT) *WASMGlobalTypeT {
	return (*WASMGlobalTypeT)(C.wasm_externtype_as_globaltype_const((*C.wasm_externtype_t)(t)))
}

func WASMExternTypeAsTableTypeConst(t *WASMExternTypeT) *WASMTableTypeT {
	return (*WASMTableTypeT)(C.wasm_externtype_as_tabletype_const((*C.wasm_externtype_t)(t)))
}

func WASMExternTypeAsMemoryTypeConst(t *WASMExternTypeT) *WASMMemoryTypeT {
	return (*WASMMemoryTypeT)(C.wasm_externtype_as_memorytype_const((*C.wasm_externtype_t)(t)))
}

/////////////////////////////////////////////////////////////////////////////
// wasm_extern_t
/////////////////////////////////////////////////////////////////////////////

func WASMExternDelete(extern *WASMExternT) {
	C.wasm_extern_delete((*C.wasm_extern_t)(extern))
}

func (e *WASMExternT) Close() error {
	WASMExternDelete(e)
	return nil
}

func WASMFuncAsExtern(fn *WASMFuncT) *WASMExternT {
	return (*WASMExternT)(C.wasm_func_as_extern((*C.wasm_func_t)(fn)))
}
func (f *WASMFuncT) AsExtern() *WASMExternT {
	return WASMFuncAsExtern(f)
}
func WASMGlobalAsExtern(global *WASMGlobalT) *WASMExternT {
	return (*WASMExternT)(C.wasm_global_as_extern((*C.wasm_global_t)(global)))
}
func (f *WASMGlobalT) AsExtern() *WASMExternT {
	return WASMGlobalAsExtern(f)
}
func WASMTableAsExtern(table *WASMTableT) *WASMExternT {
	return (*WASMExternT)(C.wasm_table_as_extern((*C.wasm_table_t)(table)))
}
func (f *WASMTableT) AsExtern() *WASMExternT {
	return WASMTableAsExtern(f)
}
func WASMMemoryAsExtern(memory *WASMMemoryT) *WASMExternT {
	return (*WASMExternT)(C.wasm_memory_as_extern((*C.wasm_memory_t)(memory)))
}
func (f *WASMMemoryT) AsExtern() *WASMExternT {
	return WASMMemoryAsExtern(f)
}
func WASMExternAsFunc(extern *WASMExternT) *WASMFuncT {
	return (*WASMFuncT)(C.wasm_extern_as_func((*C.wasm_extern_t)(extern)))
}
func (e *WASMExternT) AsFunc() *WASMFuncT {
	return WASMExternAsFunc(e)
}
func WASMExternAsGlobal(extern *WASMExternT) *WASMGlobalT {
	return (*WASMGlobalT)(C.wasm_extern_as_global((*C.wasm_extern_t)(extern)))
}
func (e *WASMExternT) AsGlobal() *WASMGlobalT {
	return WASMExternAsGlobal(e)
}
func WASMExternAsTable(extern *WASMExternT) *WASMTableT {
	return (*WASMTableT)(C.wasm_extern_as_table((*C.wasm_extern_t)(extern)))
}
func (e *WASMExternT) AsTable() *WASMTableT {
	return WASMExternAsTable(e)
}
func WASMExternAsMemory(extern *WASMExternT) *WASMMemoryT {
	return (*WASMMemoryT)(C.wasm_extern_as_memory((*C.wasm_extern_t)(extern)))
}
func (e *WASMExternT) AsMemory() *WASMMemoryT {
	return WASMExternAsMemory(e)
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

/////////////////////////////////////////////////////////////////////////////
// wasm_valtype_t
/////////////////////////////////////////////////////////////////////////////

func WASMValTypeNew(kind WASMValKindT) *WASMValTypeT {
	return (*WASMValTypeT)(C.wasm_valtype_new((C.wasm_valkind_t)(kind)))
}

func WASMValTypeDelete(t *WASMValTypeT) {
	C.wasm_valtype_delete((*C.wasm_valtype_t)(t))
}

func WASMValTypeKind(t *WASMValTypeT) WASMValKindT {
	return WASMValKindT(C.wasm_valtype_kind((*C.wasm_valtype_t)(t)))
}

func (t *WASMValTypeT) Kind() WASMValKindT {
	return WASMValTypeKind(t)
}

func WASMValTypeIsNum(t *WASMValTypeT) bool {
	return WASMValTypeKind(t).IsNum()
}

func (t *WASMValTypeT) IsNum() bool {
	return WASMValTypeKind(t).IsNum()
}

func WASMValTypeIsRef(t *WASMValTypeT) bool {
	return WASMValTypeKind(t).IsRef()
}

func (t *WASMValTypeT) IsRef() bool {
	return WASMValTypeKind(t).IsRef()
}

func WASMValTypeNewI32() *WASMValTypeT {
	return WASMValTypeNew(I32)
}

func WASMValTypeNewI64() *WASMValTypeT {
	return WASMValTypeNew(I64)
}

func WASMValTypeNewF32() *WASMValTypeT {
	return WASMValTypeNew(F32)
}

func WASMValTypeNewF64() *WASMValTypeT {
	return WASMValTypeNew(F64)
}

func WASMValTypeNewV128() *WASMValTypeT {
	return WASMValTypeNew(V128)
}

func WASMValTypeNewAnyref() *WASMValTypeT {
	return WASMValTypeNew(ANYREF)
}

func WASMValTypeNewFuncref() *WASMValTypeT {
	return WASMValTypeNew(FUNCREF)
}

/////////////////////////////////////////////////////////////////////////////
// wasm_functype_t
/////////////////////////////////////////////////////////////////////////////

func WASMFuncTypeNew_0_0() *WASMFuncTypeT {
	return WASMFuncTypeNew(nil, nil)
}

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
