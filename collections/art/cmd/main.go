package main

/*
#include "api.h"
//#include "robinhood.h"
//#include "simpleclass.h"
#include <stdlib.h>

 */
import "C"
import (
	"fmt"
	"github.com/moontrade/memory"
	"github.com/moontrade/memory/alloc/rpmalloc"
	"github.com/moontrade/memory/collections/art"
	"github.com/moontrade/memory/collections/rhmap"
	"strconv"
	"time"
)

func main() {
	rpmalloc.Malloc(128)
	println("malloc")
	C.malloc(128)

	for i := 0; i < 5; i++ {
		go func() {
			//runtime.LockOSThread()
			C.malloc(128)
		}()
	}

	println("ART")
	tree, _ := art.New()
	////println("tree", uint(uintptr(unsafe.Pointer(tree))))
	key := memory.WrapString("hello")
	value := memory.WrapString("world!")
	existing := tree.Insert(key.Pointer, key.Len(), value.Pointer)
	println("existing", uint(uintptr(existing.Unsafe())))
	existing = tree.InsertBytes(key, memory.WrapString("world 2!").Pointer)
	println("existing", memory.BytesRef(existing).String())
	existing = tree.FindBytes(key)
	println("found", memory.BytesRef(existing).String())
	tree.Free()
	println("ART size", tree.Size())
	println("ART done")
	//
	for i := 0; i < 10; i++ {
		//benchmarkCGO(100000000)
		println("running...")
		benchmarkARTInsert(1000000)
		//benchmarkGoMapInsert(1000000)
		benchmarkRHMapInsert(10000000)
		benchmarkRHMapInt64Insert(10000000)
		//time.Sleep(time.Second)
	}
	time.Sleep(time.Hour)
}

func benchmarkCGO(iterations int) {
	start := time.Now().UnixNano()
	for i := 0; i < iterations; i++ {
		C.do_stub()
	}
	end := time.Now().UnixNano()-start
	//println(float64(end)/float64(iterations))
	fmt.Printf("CGO overhead %d -> %.1fns\n", iterations, float64(end)/float64(iterations))
	//fmt.Printf("%.1fns\n", float64(end)/float64(iterations))
}

func benchmarkARTInsert(iterations int) {
	tree, _ := art.New()
	key := memory.AllocBytes(8)

	start := time.Now().UnixNano()
	for i := 0; i < iterations; i++ {
		key.SetInt64BE(0, int64(i))
		tree.InsertBytes(key, 0)
	}
	end := time.Now().UnixNano()-start
	//println(float64(end)/float64(iterations))
	fmt.Printf("ART insert %d -> %.1fns\n", iterations, float64(end)/float64(iterations))

	key.SetInt64BE(0, 500)//int64(iterations/2))
	start = time.Now().UnixNano()
	for i := 0; i < iterations; i++ {
		tree.FindBytes(key)
	}
	end = time.Now().UnixNano()-start
	//println(float64(end)/float64(iterations))
	fmt.Printf("ART get %d -> %.1fns\n", iterations, float64(end)/float64(iterations))

	tree.Free()
	key.Free()
	//fmt.Printf("%.1fns\n", float64(end)/float64(iterations))
}


func benchmarkGoMapInsert(iterations int) {
	mp := make(map[string]struct{})
	key := memory.AllocBytes(8)

	start := time.Now().UnixNano()
	for i := 0; i < iterations; i++ {
		//key.SetInt64BE(0, int64(i))
		mp[strconv.Itoa(i)] = struct{}{}
	}
	end := time.Now().UnixNano()-start
	//println(float64(end)/float64(iterations))
	fmt.Printf("GoMap insert %d -> %.1fns\n", iterations, float64(end)/float64(iterations))

	key.SetInt64BE(0, 500)//int64(iterations/2))
	searchKey := strconv.Itoa(iterations/2)
	start = time.Now().UnixNano()
	for i := 0; i < iterations; i++ {
		_, _ = mp[searchKey]
	}
	end = time.Now().UnixNano()-start
	//println(float64(end)/float64(iterations))
	fmt.Printf("GoMap get %d -> %.1fns\n", iterations, float64(end)/float64(iterations))

	key.Free()
	//fmt.Printf("%.1fns\n", float64(end)/float64(iterations))
}

func benchmarkRHMapInsert(iterations int) {
	mp := rhmap.NewMap(uintptr(iterations*2))
	key := memory.AllocBytes(8)

	start := time.Now().UnixNano()
	for i := 0; i < iterations; i++ {
		key.SetInt64BE(0, int64(i))
		mp.Set(key, memory.Bytes{0})
	}
	end := time.Now().UnixNano()-start
	//println(float64(end)/float64(iterations))
	fmt.Printf("RHMap insert %d -> %.1fns\n", iterations, float64(end)/float64(iterations))

	key.SetInt64BE(0, int64(iterations/2))//int64(iterations/2))
	start = time.Now().UnixNano()
	for i := 0; i < iterations; i++ {
		mp.Get(key)
	}
	end = time.Now().UnixNano()-start
	//println(float64(end)/float64(iterations))
	fmt.Printf("RHMap get %d -> %.1fns\n", iterations, float64(end)/float64(iterations))

	_ = mp.Close()
	key.Free()
	//fmt.Printf("%.1fns\n", float64(end)/float64(iterations))
}

func benchmarkRHMapInt64Insert(iterations int64) {
	mp := rhmap.NewMapInt64(uintptr(iterations*2))

	start := time.Now().UnixNano()
	for i := int64(1); i < iterations; i++ {
		mp.Set(i, 1)
	}
	end := time.Now().UnixNano()-start
	//println(float64(end)/float64(iterations))
	fmt.Printf("RHMapInt64 insert %d -> %.1fns\n", iterations, float64(end)/float64(iterations))

	key := iterations/2
	start = time.Now().UnixNano()
	for i := int64(1); i < iterations; i++ {
		mp.Get(key)
	}
	end = time.Now().UnixNano()-start
	//println(float64(end)/float64(iterations))
	fmt.Printf("RHMapInt64 get %d -> %.1fns\n", iterations, float64(end)/float64(iterations))

	_ = mp.Close()
	//fmt.Printf("%.1fns\n", float64(end)/float64(iterations))
}