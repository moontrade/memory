package art

/*
#include "art.h"
#include <stdlib.h>

typedef struct {
	size_t ptr;
	size_t code;
} art_new_t;

void do_art_new(size_t arg0, size_t arg1) {
	art_new_t* args = (art_new_t*)(void*)arg0;
	art_tree* tree = (art_tree*)calloc(1, sizeof(art_tree));
	args->code = (size_t)art_tree_init(tree);
	args->ptr = (size_t)(void*)tree;
}

void do_art_destroy(size_t arg0, size_t arg1) {
	art_tree_destroy((art_tree*)(void*)arg0);
	free((art_tree*)(void*)arg0);
}

typedef struct {
	size_t tree;
	size_t size;
} art_size_t;

void do_art_size(size_t arg0, size_t arg1) {
	art_size_t* args = (art_size_t*)(void*)arg0;
	args->size = (size_t)art_size((art_tree*)(void*)args->tree);
}

typedef struct {
	size_t tree;
	size_t key;
	size_t len;
	size_t data;
	size_t old;
} art_insert_t;

void do_art_insert(size_t arg0, size_t arg1) {
	art_insert_t* args = (art_insert_t*)(void*)arg0;
	args->old = (size_t)art_insert((art_tree*)(void*)args->tree, (unsigned char*)args->key, (int)args->len, (void*)args->data);
}

void do_art_insert_no_replace(size_t arg0, size_t arg1) {
	art_insert_t* args = (art_insert_t*)(void*)arg0;
	args->old = (size_t)art_insert_no_replace((art_tree*)(void*)args->tree, (unsigned char*)args->key, (int)args->len, (void*)args->data);
}

typedef struct {
	size_t tree;
	size_t key;
	size_t len;
	size_t value;
} art_delete_t;

void do_art_delete(size_t arg0, size_t arg1) {
	art_delete_t* args = (art_delete_t*)(void*)arg0;
	args->value = (size_t)art_delete((art_tree*)(void*)args->tree, (unsigned char*)args->key, (int)args->len);
}

typedef struct {
	size_t tree;
	size_t key;
	size_t len;
	size_t result;
} art_search_t;

void do_art_search(size_t arg0, size_t arg1) {
	art_search_t* args = (art_search_t*)(void*)arg0;
	args->result = (size_t)art_search((art_tree*)(void*)args->tree, (unsigned char*)args->key, (int)args->len);
}

typedef struct {
	size_t tree;
	size_t result;
} art_minmax_t;

void do_art_minimum(size_t arg0, size_t arg1) {
	art_minmax_t* args = (art_minmax_t*)(void*)arg0;
	args->result = (size_t)art_minimum((art_tree*)(void*)args->tree);
}

void do_art_maximum(size_t arg0, size_t arg1) {
	art_minmax_t* args = (art_minmax_t*)(void*)arg0;
	args->result = (size_t)art_maximum((art_tree*)(void*)args->tree);
}

*/
import "C"
import (
	"github.com/moontrade/memory"
	"github.com/moontrade/memory/unsafecgo"
	"reflect"
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

type artNewT struct {
	ptr  uintptr
	code uintptr
}

func New() *Tree {
	args := artNewT{}
	ptr := uintptr(unsafe.Pointer(&args))
	unsafecgo.NonBlocking((*byte)(C.do_art_new), ptr, 0)
	return (*Tree)(unsafe.Pointer(args.ptr))
}

func (r *Tree) Free() {
	ptr := uintptr(unsafe.Pointer(r))
	unsafecgo.NonBlocking((*byte)(C.do_art_destroy), ptr, 0)
}

type artSizeT struct {
	ptr  uintptr
	size uintptr
}

func (r *Tree) Size() int {
	args := artSizeT{ptr: uintptr(unsafe.Pointer(r))}
	ptr := uintptr(unsafe.Pointer(&args))
	unsafecgo.NonBlocking((*byte)(C.do_art_size), ptr, 0)
	return int(args.size)
}

//func (r *Art) Print() {
//	unsafecgo.NonBlocking((*byte)(C.do_rax_show), uintptr(unsafe.Pointer(r)), 0)
//}

type artInsertT struct {
	tree  uintptr
	key   uintptr
	len   uintptr
	value uintptr
	old   uintptr
}

func (r *Tree) Insert(key memory.Pointer, size int, value memory.Pointer) memory.Pointer {
	args := artInsertT{
		tree:  uintptr(unsafe.Pointer(r)),
		key:   uintptr(key),
		len:   uintptr(size),
		value: uintptr(value),
	}
	ptr := uintptr(unsafe.Pointer(&args))
	unsafecgo.NonBlocking((*byte)(C.do_art_insert), ptr, 0)
	return memory.Pointer(args.old)
}

func (r *Tree) InsertBytes(key memory.Bytes, value memory.Pointer) memory.Pointer {
	return r.Insert(key.Pointer, key.Len(), value)
}

func (r *Tree) InsertString(key string, value memory.Pointer) memory.Pointer {
	k := (*reflect.StringHeader)(unsafe.Pointer(&key))
	return r.Insert(memory.Pointer(k.Data), k.Len, value)
}

func (r *Tree) InsertSlice(key []byte, value memory.Pointer) memory.Pointer {
	k := (*reflect.SliceHeader)(unsafe.Pointer(&key))
	return r.Insert(memory.Pointer(k.Data), k.Len, value)
}

func (r *Tree) InsertNoReplace(key memory.Pointer, size int, value memory.Pointer) memory.Pointer {
	args := artInsertT{
		tree:  uintptr(unsafe.Pointer(r)),
		key:   uintptr(key),
		len:   uintptr(size),
		value: uintptr(value),
	}
	ptr := uintptr(unsafe.Pointer(&args))
	unsafecgo.NonBlocking((*byte)(C.do_art_insert_no_replace), ptr, 0)
	return memory.Pointer(args.old)
}

func (r *Tree) InsertNoReplaceBytes(key memory.Bytes, value memory.Pointer) memory.Pointer {
	return r.InsertNoReplace(key.Pointer, key.Len(), value)
}

func (r *Tree) InsertNoReplaceString(key string, value memory.Pointer) memory.Pointer {
	k := (*reflect.StringHeader)(unsafe.Pointer(&key))
	return r.InsertNoReplace(memory.Pointer(k.Data), k.Len, value)
}

func (r *Tree) InsertNoReplaceSlice(key []byte, value memory.Pointer) memory.Pointer {
	k := (*reflect.SliceHeader)(unsafe.Pointer(&key))
	return r.InsertNoReplace(memory.Pointer(k.Data), k.Len, value)
}

type artDeleteT struct {
	tree uintptr
	key  uintptr
	len  uintptr
	item uintptr
}

func (r *Tree) Delete(key memory.Pointer, size int) memory.Pointer {
	args := artDeleteT{
		tree: uintptr(unsafe.Pointer(r)),
		key:  uintptr(key),
		len:  uintptr(size),
	}
	ptr := uintptr(unsafe.Pointer(&args))
	unsafecgo.NonBlocking((*byte)(C.do_art_delete), ptr, 0)
	return memory.Pointer(args.item)
}

func (r *Tree) DeleteBytes(key memory.Bytes) memory.Pointer {
	return r.Delete(key.Pointer, key.Len())
}

type artSearchT struct {
	tree   uintptr
	s      uintptr
	len    uintptr
	result uintptr
}

func (r *Tree) Find(key memory.Pointer, size int) memory.Pointer {
	args := artSearchT{
		tree: uintptr(unsafe.Pointer(r)),
		s:    uintptr(key),
		len:  uintptr(size),
	}
	ptr := uintptr(unsafe.Pointer(&args))
	unsafecgo.NonBlocking((*byte)(C.do_art_search), ptr, 0)
	return memory.Pointer(args.result)
}

func (r *Tree) FindBytes(key memory.Bytes) memory.Pointer {
	return r.Find(key.Pointer, key.Len())
}

type artTreeMinmaxT struct {
	tree   uintptr
	result uintptr
}

func (r *Tree) Minimum() *Leaf {
	args := artTreeMinmaxT{
		tree: uintptr(unsafe.Pointer(r)),
	}
	ptr := uintptr(unsafe.Pointer(&args))
	unsafecgo.NonBlocking((*byte)(C.do_art_minimum), ptr, 0)
	return (*Leaf)(unsafe.Pointer(args.result))
}

// Maximum Returns the maximum valued leaf
// @return The maximum leaf or NULL
func (r *Tree) Maximum() *Leaf {
	args := artTreeMinmaxT{
		tree: uintptr(unsafe.Pointer(r)),
	}
	ptr := uintptr(unsafe.Pointer(&args))
	unsafecgo.NonBlocking((*byte)(C.do_art_maximum), ptr, 0)
	return (*Leaf)(unsafe.Pointer(args.result))
}
