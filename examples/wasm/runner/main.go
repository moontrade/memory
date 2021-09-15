// A Wasm module can sometimes not work, it will produce errors in some
// cases.
//
// In this example we'll see how to handle such errors in the most
// basic way. To do that we'll use a Wasm module that we know will
// produce an error.
//
// You can run the example directly by executing in Wasmer root:
//
// ```shell
// go test examples/example_early_exit_test.go
// ```
//
// Ready?

package main

import (
	"encoding/binary"
	"fmt"
	"github.com/wasmerio/wasmer-go/wasmer"
	"os"
	"runtime"
	"strings"
	"time"
	"unsafe"
)

type exitCode struct {
	code int32
}

func (self *exitCode) Error() string {
	return fmt.Sprintf("exit code: %d", self.code)
}

func earlyExit(args []wasmer.Value) ([]wasmer.Value, error) {
	return nil, &exitCode{1}
}

//func main2() {
//	dylib, err := os.ReadFile("./rust/myproject/pkg/main.dylib")
//	if err != nil {
//
//		panic(err)
//	}
//
//	// Create an Engine
//	engine := wasmer.NewEngineWithConfig(wasmer.NewConfig().
//		UseDylibEngine(),
//	)
//
//	// Create a Store
//	store := wasmer.NewStore(engine)
//
//	fmt.Println("Compiling module...")
//	module, err := wasmer.DeserializeModule(store, dylib)
//	//module, err := wasmer.NewModule(store, dylib)
//
//	if err != nil {
//		fmt.Println("Failed to compile module:", err)
//	}
//
//	importObject := wasmer.NewImportObject()
//	fmt.Println("Instantiating module...")
//	// Let's instantiate the Wasm module.
//	instance, err := wasmer.NewInstance(module, importObject)
//
//	instance.Close()
//}

func main() {
	dylib, err := os.ReadFile("../wasm-local/main.dylib")
	if err != nil {
		dylib, err = os.ReadFile("../wasm/main.dylib")
		if err != nil {
			panic(err)
		}
	}

	// Create an Engine
	engine := wasmer.NewEngineWithConfig(wasmer.NewConfig().
		UseDylibEngine(),
	)

	// Create a Store
	store := wasmer.NewStore(engine)

	fmt.Println("Compiling module...")
	module, err := wasmer.DeserializeModule(store, dylib)
	//module, err := wasmer.NewModule(store, dylib)

	if err != nil {
		fmt.Println("Failed to compile module:", err)
	}

	env, err := wasmer.NewWasiStateBuilder("main").
		Argument("-conf hi").
		//Argument("-port 80").
		CaptureStderr().
		CaptureStdout().
		Finalize()
	if err != nil {
		panic(err)
	}

	importObject, err := env.GenerateImportObject(store, module)
	if err != nil {
		panic(err)
	}

	imports := make(map[string]wasmer.IntoExtern)
	environment := &Env{
		env:   env,
		store: store,
	}
	environment.register(imports)
	// Create an import object with the expected function.
	//importObject := wasmer.NewImportObject()
	importObject.Register(
		"env",
		imports,
	)

	fmt.Println("Instantiating module...")
	// Let's instantiate the Wasm module.
	instance, err := wasmer.NewInstance(module, importObject)

	if err != nil {
		panic(fmt.Sprintln("Failed to instantiate the module:", err))
	}

	mem, err := instance.Exports.GetMemory("memory")
	environment.mem = mem

	resume, err := instance.Exports.GetFunction("resume")
	environment.resume = resume
	gc, err := instance.Exports.GetFunction("gc")
	environment.gc = gc
	start, err := instance.Exports.GetWasiStartFunction()
	if err != nil {
		panic(err)
	}
	environment.start = start

	go environment.run()

	time.Sleep(time.Hour)
}

type Env struct {
	env         *wasmer.WasiEnvironment
	store       *wasmer.Store
	goScheduler wasmer.NativeFunction
	mem         *wasmer.Memory
	nowGet      []wasmer.Value

	start  wasmer.NativeFunction
	resume wasmer.NativeFunction
	gc     wasmer.NativeFunction

	nextSleep time.Duration
}

func btos(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func (env *Env) run() {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	var err error
	_, err = env.start()
	if err != nil {
		return
	}
	fmt.Println(string(env.env.ReadStdout()))

	for {
		//println("before resume")
		_, _ = env.resume()

		output := strings.TrimSpace(btos(env.env.ReadStdout()))
		if len(output) > 0 {
			fmt.Println(output)
		}

		//println("after resume")
		//println("sleeping for...", env.nextSleep.String())
		//println("memory pages", env.mem.Size())
		time.Sleep(env.nextSleep)

		//start := time.Now()
		//_, _ = env.gc()
		//println("GC", time.Now().Sub(start).String())
	}
}

func (env *Env) setTimeout(values []wasmer.Value) ([]wasmer.Value, error) {
	sleepFor := time.Duration(values[0].I64())
	env.nextSleep = sleepFor

	if env.nextSleep == 0 {
		env.nextSleep = sleepFor
	} else {
		if env.nextSleep > sleepFor {
			env.nextSleep = sleepFor
		}
	}
	return nil, nil
}

func (env *Env) register(m map[string]wasmer.IntoExtern) {
	env.nowGet = make([]wasmer.Value, 1)
	m["main.now"] = wasmer.NewFunction(
		env.store,
		wasmer.NewFunctionType(
			wasmer.NewValueTypes(wasmer.I32, wasmer.I32),
			wasmer.NewValueTypes(wasmer.I64)),
		env.Now,
	)
	m["setTimeout"] = wasmer.NewFunction(
		env.store,
		wasmer.NewFunctionType(
			wasmer.NewValueTypes(wasmer.I64),
			wasmer.NewValueTypes()),
		env.setTimeout,
	)
	m["time.startTimer"] = wasmer.NewFunction(
		env.store,
		wasmer.NewFunctionType(
			wasmer.NewValueTypes(wasmer.I32, wasmer.I32, wasmer.I32),
			wasmer.NewValueTypes()),
		env.StartTimer,
	)
	m["time.resetTimer"] = wasmer.NewFunction(
		env.store,
		wasmer.NewFunctionType(
			wasmer.NewValueTypes(wasmer.I32, wasmer.I64, wasmer.I32, wasmer.I32),
			wasmer.NewValueTypes(wasmer.I32)),
		env.ResetTimer,
	)
	//m["runtime.ticks"] = wasmer.NewFunction(
	//	env.store,
	//	wasmer.NewFunctionType(
	//		wasmer.NewValueTypes(),
	//		wasmer.NewValueTypes(wasmer.F64)),
	//	env.RuntimeTicks,
	//)
	//m["runtime.sleepTicks"] = wasmer.NewFunction(
	//	env.store,
	//	wasmer.NewFunctionType(
	//		wasmer.NewValueTypes(wasmer.F64),
	//		wasmer.NewValueTypes()),
	//	env.RuntimeSleepTicks,
	//)
	//m["syscall/js.valueGet"] = wasmer.NewFunction(
	//	env.store,
	//	wasmer.NewFunctionType(
	//		wasmer.NewValueTypes(wasmer.I32, wasmer.I32, wasmer.I32, wasmer.I32, wasmer.I32, wasmer.I32),
	//		wasmer.NewValueTypes()),
	//	env.syscallJsValueGet,
	//)
	//m["syscall/js.valueSet"] = wasmer.NewFunction(
	//	env.store,
	//	wasmer.NewFunctionType(
	//		wasmer.NewValueTypes(wasmer.I32, wasmer.I32, wasmer.I32, wasmer.I32, wasmer.I32, wasmer.I32),
	//		wasmer.NewValueTypes()),
	//	env.syscallJsValueSet,
	//)
	//m["syscall/js.valueDelete"] = wasmer.NewFunction(
	//	env.store,
	//	wasmer.NewFunctionType(
	//		wasmer.NewValueTypes(wasmer.I32, wasmer.I32, wasmer.I32),
	//		wasmer.NewValueTypes()),
	//	env.syscallJsValueDelete,
	//)
	//m["syscall/js.valueIndex"] = wasmer.NewFunction(
	//	env.store,
	//	wasmer.NewFunctionType(
	//		wasmer.NewValueTypes(wasmer.I32, wasmer.I32, wasmer.I32, wasmer.I32, wasmer.I32),
	//		wasmer.NewValueTypes()),
	//	env.syscallJsValueIndex,
	//)
	//m["syscall/js.valueSetIndex"] = wasmer.NewFunction(
	//	env.store,
	//	wasmer.NewFunctionType(
	//		wasmer.NewValueTypes(wasmer.I32, wasmer.I32, wasmer.I32),
	//		wasmer.NewValueTypes()),
	//	env.syscallJsValueSetIndex,
	//)
	//m["syscall/js.valueLength"] = wasmer.NewFunction(
	//	env.store,
	//	wasmer.NewFunctionType(
	//		wasmer.NewValueTypes(wasmer.I32, wasmer.I32, wasmer.I32),
	//		wasmer.NewValueTypes(wasmer.I32)),
	//	env.syscallJsValueLength,
	//)
	//m["syscall/js.finalizeRef"] = wasmer.NewFunction(
	//	env.store,
	//	wasmer.NewFunctionType(
	//		wasmer.NewValueTypes(wasmer.I32, wasmer.I32, wasmer.I32),
	//		wasmer.NewValueTypes()),
	//	env.syscallJsFinalizeRef,
	//)
	//
	//m["syscall/js.stringVal"] = wasmer.NewFunction(
	//	env.store,
	//	wasmer.NewFunctionType(
	//		wasmer.NewValueTypes(wasmer.I32, wasmer.I32, wasmer.I32, wasmer.I32, wasmer.I32),
	//		wasmer.NewValueTypes()),
	//	env.syscallJsStringVal,
	//)
	//
	//m["syscall/js.valuePrepareString"] = wasmer.NewFunction(
	//	env.store,
	//	wasmer.NewFunctionType(
	//		wasmer.NewValueTypes(wasmer.I32, wasmer.I32, wasmer.I32, wasmer.I32),
	//		wasmer.NewValueTypes()),
	//	env.syscallJsValuePrepareString,
	//)
	//
	//m["syscall/js.valueLoadString"] = wasmer.NewFunction(
	//	env.store,
	//	wasmer.NewFunctionType(
	//		wasmer.NewValueTypes(wasmer.I32, wasmer.I32, wasmer.I32, wasmer.I32, wasmer.I32, wasmer.I32),
	//		wasmer.NewValueTypes()),
	//	env.syscallJsValueLoadString,
	//)
	//
	//m["syscall/js.valueCall"] = wasmer.NewFunction(
	//	env.store,
	//	wasmer.NewFunctionType(
	//		wasmer.NewValueTypes(wasmer.I32, wasmer.I32, wasmer.I32, wasmer.I32, wasmer.I32, wasmer.I32, wasmer.I32, wasmer.I32, wasmer.I32),
	//		wasmer.NewValueTypes()),
	//	env.syscallJsValueLoadString,
	//)
}

func (env *Env) StartTimer(values []wasmer.Value) ([]wasmer.Value, error) {
	//env.nowGet[0] = wasmer.NewValue(time.Now().UnixNano(), wasmer.I64)
	addr := values[0].I32()
	ptr := values[1].I32()
	l := values[2].I32()

	//
	_ = addr
	_ = ptr
	_ = l
	//
	//binary.LittleEndian.PutUint64(env.mem.Data()[addr:], uint64(time.Now().UnixNano()))
	return nil, nil
}

func (env *Env) ResetTimer(values []wasmer.Value) ([]wasmer.Value, error) {
	//env.nowGet[0] = wasmer.NewValue(time.Now().UnixNano(), wasmer.I64)
	addr := values[0].I32()
	tm := values[1].I64()
	ptr := values[2].I32()
	l := values[3].I32()
	//
	_ = addr
	_ = tm
	_ = ptr
	_ = l
	//
	//binary.LittleEndian.PutUint64(env.mem.Data()[addr:], uint64(time.Now().UnixNano()))
	return nil, nil
}

func (env *Env) Now(values []wasmer.Value) ([]wasmer.Value, error) {
	env.nowGet[0] = wasmer.NewValue(time.Now().UnixNano(), wasmer.I64)
	//addr := values[0].I32()
	//ptr := values[1].I32()
	//l := values[2].I32()
	//
	//_ = addr
	//_ = ptr
	//_ = l
	//
	//binary.LittleEndian.PutUint64(env.mem.Data()[addr:], uint64(time.Now().UnixNano()))
	return env.nowGet, nil
}

func (env *Env) RuntimeTicks(values []wasmer.Value) ([]wasmer.Value, error) {
	r := make([]wasmer.Value, 1)
	r[0] = wasmer.NewValue(float64(time.Now().UnixNano()), wasmer.F64)
	return r, nil
}

func (env *Env) RuntimeSleepTicks(values []wasmer.Value) ([]wasmer.Value, error) {
	env.goScheduler()
	return nil, nil
}

func (env *Env) syscallJsFinalizeRef(values []wasmer.Value) ([]wasmer.Value, error) {
	sp := values[0].I32()

	_ = sp

	return nil, nil
}

func (env *Env) syscallJsStringVal(values []wasmer.Value) ([]wasmer.Value, error) {
	ret_ptr := values[0].I32()
	value_ptr := values[1].I32()
	value_len := values[1].I32()

	_ = ret_ptr
	_ = value_ptr
	_ = value_len

	return nil, nil
}

func (env *Env) syscallJsValueGet(values []wasmer.Value) ([]wasmer.Value, error) {
	retval := values[0].I32()
	p_addr := values[1].I32()
	p_ptr := values[2].I32()
	p_len := values[3].I32()

	data := env.mem.Data()
	prop := string(data[p_ptr : p_ptr+p_len])
	value := data[p_addr]

	_ = prop
	_ = value

	_ = retval
	_ = p_addr
	_ = p_ptr
	_ = p_len

	binary.LittleEndian.PutUint32(data[p_addr:], uint32(value))

	return nil, nil
}

func (env *Env) syscallJsValueSet(values []wasmer.Value) ([]wasmer.Value, error) {
	retval := values[0].I32()
	p_addr := values[1].I32()
	p_ptr := values[2].I32()
	p_len := values[3].I32()

	_ = retval
	_ = p_addr
	_ = p_ptr
	_ = p_len

	return nil, nil
}

func (env *Env) syscallJsValueDelete(values []wasmer.Value) ([]wasmer.Value, error) {
	v_addr := values[0].I32()
	p_ptr := values[1].I32()
	p_len := values[2].I32()

	_ = v_addr
	_ = p_ptr
	_ = p_len

	return nil, nil
}

func (env *Env) syscallJsValueIndex(values []wasmer.Value) ([]wasmer.Value, error) {
	v_addr := values[0].I32()
	p_ptr := values[1].I32()
	p_len := values[2].I32()

	_ = v_addr
	_ = p_ptr
	_ = p_len

	return nil, nil
}

func (env *Env) syscallJsValueSetIndex(values []wasmer.Value) ([]wasmer.Value, error) {
	v_addr := values[0].I32()
	p_ptr := values[1].I32()
	p_len := values[2].I32()

	_ = v_addr
	_ = p_ptr
	_ = p_len

	return nil, nil
}

func (env *Env) syscallJsValueLength(values []wasmer.Value) ([]wasmer.Value, error) {
	v_addr := values[0].I32()

	_ = v_addr
	return nil, nil
}

func (env *Env) syscallJsValuePrepareString(values []wasmer.Value) ([]wasmer.Value, error) {
	ret_addr := values[0].I32()
	v_addr := values[1].I32()

	_ = ret_addr
	_ = v_addr

	return nil, nil
}

func (env *Env) syscallJsValueLoadString(values []wasmer.Value) ([]wasmer.Value, error) {
	v_addr := values[0].I32()
	slice_ptr := values[1].I32()
	slice_len := values[2].I32()
	slice_cap := values[3].I32()

	_ = v_addr
	_ = slice_ptr
	_ = slice_len
	_ = slice_cap

	return nil, nil
}

func (env *Env) syscallJsValueCall(values []wasmer.Value) ([]wasmer.Value, error) {
	ret_addr := values[0].I32()
	_ = ret_addr
	//v_addr := values[0].I32()
	//slice_ptr := values[1].I32()
	//slice_len := values[2].I32()
	//slice_cap := values[3].I32()
	//
	//_ = v_addr
	//_ = slice_ptr
	//_ = slice_len
	//_ = slice_cap

	return nil, nil
}
