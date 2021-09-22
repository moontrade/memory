package art

import (
	mem "github.com/moontrade/memory"
	"unsafe"
)

type nodeFactory interface {
	newNode4() *artNode
	newNode16() *artNode
	newNode48() *artNode
	newNode256() *artNode
	newLeaf(key Key, value interface{}) *artNode
}

// make sure that objFactory implements all methods of nodeFactory interface
var _ nodeFactory = &objFactory{}

var allocator = mem.NextAllocator()
var factory = newObjFactory(allocator)

func newTree() *tree {
	return (*tree)(unsafe.Pointer(allocator.AllocZeroed(unsafe.Sizeof(tree{}))))
}

//type objFactory struct{}
//
//func newObjFactory() nodeFactory {
//	return &objFactory{}
//}
//
//// Simple obj factory implementation
//func (f *objFactory) newNode4() *artNode {
//	return &artNode{kind: Node4, ref: unsafe.Pointer(new(node4))}
//}
//
//func (f *objFactory) newNode16() *artNode {
//	return &artNode{kind: Node16, ref: unsafe.Pointer(&node16{})}
//}
//
//func (f *objFactory) newNode48() *artNode {
//	return &artNode{kind: Node48, ref: unsafe.Pointer(&node48{})}
//}
//
//func (f *objFactory) newNode256() *artNode {
//	return &artNode{kind: Node256, ref: unsafe.Pointer(&node256{})}
//}
//
//func (f *objFactory) newLeaf(key Key, value interface{}) *artNode {
//	clonedKey := Key{key.Clone()}
//	//clonedKey := make(Key, len(key))
//	//copy(clonedKey, key)
//	return &artNode{
//		kind: Leaf,
//		ref:  unsafe.Pointer(&leaf{key: clonedKey, value: value}),
//	}
//}

type objFactory struct {
	a mem.Allocator
}

func newObjFactory(a mem.Allocator) nodeFactory {
	f := (*objFactory)(unsafe.Pointer(a.AllocZeroed(unsafe.Sizeof(objFactory{}))))
	f.a = a
	return f
}

// Simple obj factory implementation
func (f *objFactory) newNode4() *artNode {
	n := (*artNode)(unsafe.Pointer(f.a.AllocZeroed(unsafe.Sizeof(artNode{}))))
	n4 := (*node48)(unsafe.Pointer(f.a.AllocZeroed(unsafe.Sizeof(node4{}))))
	n.kind = Node4
	n.ref = unsafe.Pointer(n4)
	return n
}

func (f *objFactory) newNode16() *artNode {
	n := (*artNode)(unsafe.Pointer(f.a.AllocZeroed(unsafe.Sizeof(artNode{}))))
	n16 := (*node48)(unsafe.Pointer(f.a.AllocZeroed(unsafe.Sizeof(node16{}))))
	n.kind = Node16
	n.ref = unsafe.Pointer(n16)
	return n
}

func (f *objFactory) newNode48() *artNode {
	n := (*artNode)(unsafe.Pointer(f.a.AllocZeroed(unsafe.Sizeof(artNode{}))))
	n48 := (*node48)(unsafe.Pointer(f.a.AllocZeroed(unsafe.Sizeof(node48{}))))
	n.kind = Node48
	n.ref = unsafe.Pointer(n48)
	return n
}

func (f *objFactory) newNode256() *artNode {
	n := (*artNode)(unsafe.Pointer(f.a.AllocZeroed(unsafe.Sizeof(artNode{}))))
	n256 := (*node256)(unsafe.Pointer(f.a.AllocZeroed(unsafe.Sizeof(node256{}))))
	n.kind = Node256
	n.ref = unsafe.Pointer(n256)
	return n
}

func (f *objFactory) newLeaf(key Key, value interface{}) *artNode {
	clonedKey := key.Clone()
	n := (*artNode)(unsafe.Pointer(f.a.AllocZeroed(unsafe.Sizeof(artNode{}))))
	l := (*leaf)(unsafe.Pointer(f.a.AllocZeroed(unsafe.Sizeof(leaf{}))))
	n.kind = Leaf
	l.key = clonedKey
	l.value = value
	n.ref = unsafe.Pointer(l)
	//clonedKey := make(Key, len(key))
	//copy(clonedKey, key)
	return &artNode{
		kind: Leaf,
		ref:  unsafe.Pointer(&leaf{key: clonedKey, value: value}),
	}
}
