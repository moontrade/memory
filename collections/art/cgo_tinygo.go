//go:build tinygo && (darwin || linux || unix) && !tinygo.wasm
// +build tinygo
// +build darwin linux unix
// +build !tinygo.wasm

package art

/*
#include "art.h"
#include <stdlib.h>
*/
import "C"
import (
	"github.com/moontrade/memory"
	"unsafe"
)

type Tree C.art_tree

type Leaf C.art_leaf

func (l *Leaf) Data() memory.Pointer {
	return *(*memory.Pointer)(unsafe.Pointer(l))
}
func (l *Leaf) Key() memory.FatPointer {
	return memory.FatPointerOf(
		memory.Pointer(uintptr(unsafe.Pointer(l))+unsafe.Sizeof(uintptr(0))+4),
		uintptr(*(*uint32)(unsafe.Pointer(uintptr(unsafe.Pointer(l)) + unsafe.Sizeof(uintptr(0))))))
}

//go:nosplit
//go:noescape
func New() (*Tree, int) {
	tree := C.calloc(1, C.ulong(unsafe.Sizeof(C.art_tree{})))
	return (*Tree)(tree), int(C.art_tree_init((*C.art_tree)(tree)))
}

//go:nosplit
//go:noescape
func (t *Tree) Free() {
	C.art_tree_destroy((*C.art_tree)(t))
}

func (t *Tree) Size() int {
	return int(*(*uintptr)(unsafe.Pointer(uintptr(unsafe.Pointer(t)) + unsafe.Sizeof(uintptr(0)))))
}

type _string struct {
	Data, Len uintptr
}

type _bytes struct {
	Data, Len, Cap uintptr
}

//go:nosplit
//go:noescape
func (t *Tree) Insert(key memory.Pointer, size int, value memory.Pointer) memory.Pointer {
	return memory.Pointer(C.art_insert((*C.art_tree)(t), (*C.uchar)(unsafe.Pointer(key)), C.int(size), unsafe.Pointer(value)))
}

func (t *Tree) InsertBytes(key memory.Bytes, value memory.Pointer) memory.Pointer {
	return t.Insert(key.Pointer, key.Len(), value)
}

func (t *Tree) InsertString(key string, value memory.Pointer) memory.Pointer {
	k := (*_string)(unsafe.Pointer(&key))
	return t.Insert(memory.Pointer(k.Data), int(k.Len), value)
}

func (t *Tree) InsertSlice(key []byte, value memory.Pointer) memory.Pointer {
	k := (*_bytes)(unsafe.Pointer(&key))
	return t.Insert(memory.Pointer(k.Data), int(k.Len), value)
}

//go:nosplit
//go:noescape
func (t *Tree) InsertNoReplace(key memory.Pointer, size int, value memory.Pointer) memory.Pointer {
	return memory.Pointer(C.art_insert_no_replace((*C.art_tree)(t), (*C.uchar)(unsafe.Pointer(key)), C.int(size), unsafe.Pointer(value)))
}

func (t *Tree) InsertNoReplaceBytes(key memory.Bytes, value memory.Pointer) memory.Pointer {
	return t.InsertNoReplace(key.Pointer, key.Len(), value)
}

func (t *Tree) InsertNoReplaceString(key string, value memory.Pointer) memory.Pointer {
	k := (*_string)(unsafe.Pointer(&key))
	return t.InsertNoReplace(memory.Pointer(k.Data), int(k.Len), value)
}

func (t *Tree) InsertNoReplaceSlice(key []byte, value memory.Pointer) memory.Pointer {
	k := (*_bytes)(unsafe.Pointer(&key))
	return t.InsertNoReplace(memory.Pointer(k.Data), int(k.Len), value)
}

//go:nosplit
//go:noescape
func (t *Tree) Delete(key memory.Pointer, size int) memory.Pointer {
	return memory.Pointer(C.art_delete((*C.art_tree)(t), (*C.uchar)(unsafe.Pointer(key)), C.int(size)))
}

func (t *Tree) DeleteBytes(key memory.Bytes) memory.Pointer {
	return t.Delete(key.Pointer, key.Len())
}

//go:nosplit
//go:noescape
func (t *Tree) Find(key memory.Pointer, size int) memory.Pointer {
	return memory.Pointer(C.art_search((*C.art_tree)(t), (*C.uchar)(unsafe.Pointer(key)), C.int(size)))
}

func (t *Tree) FindBytes(key memory.Bytes) memory.Pointer {
	return t.Find(key.Pointer, key.Len())
}

//go:nosplit
//go:noescape
func (t *Tree) Minimum() memory.Pointer {
	return memory.Pointer(unsafe.Pointer(C.art_minimum((*C.art_tree)(t))))
}

//go:nosplit
//go:noescape
func (t *Tree) Maximum() memory.Pointer {
	return memory.Pointer(unsafe.Pointer(C.art_maximum((*C.art_tree)(t))))
}
