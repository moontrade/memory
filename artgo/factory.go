package artgo

import (
	"github.com/moontrade/memory"
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

var factory = newObjFactory()

func newTree() *tree {
	return (*tree)(unsafe.Pointer(memory.AllocZeroed(unsafe.Sizeof(tree{}))))
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

type objFactory struct{}

func newObjFactory() nodeFactory {
	f := (*objFactory)(unsafe.Pointer(memory.AllocZeroed(unsafe.Sizeof(objFactory{}))))
	return f
}

// Simple obj factory implementation
func (f *objFactory) newNode4() *artNode {
	n := (*artNode)(unsafe.Pointer(memory.AllocZeroed(unsafe.Sizeof(artNode{}) + unsafe.Sizeof(node4{}))))
	n4 := (*node4)(unsafe.Pointer(uintptr(unsafe.Pointer(n)) + unsafe.Sizeof(artNode{})))
	n.kind = Node4
	n.ref = uintptr(unsafe.Pointer(n4))
	return n
}

func (f *objFactory) newNode16() *artNode {
	n := (*artNode)(unsafe.Pointer(memory.AllocZeroed(unsafe.Sizeof(artNode{}) + unsafe.Sizeof(node16{}))))
	n16 := (*node16)(unsafe.Pointer(uintptr(unsafe.Pointer(n)) + unsafe.Sizeof(artNode{})))
	//n16 := (*node48)(unsafe.Pointer(memory.AllocZeroed(unsafe.Sizeof(node16{}))))
	n.kind = Node16
	n.ref = uintptr(unsafe.Pointer(n16))
	return n
}

func (f *objFactory) newNode48() *artNode {
	n := (*artNode)(unsafe.Pointer(memory.AllocZeroed(unsafe.Sizeof(artNode{}) + unsafe.Sizeof(node48{}))))
	n48 := (*node48)(unsafe.Pointer(uintptr(unsafe.Pointer(n)) + unsafe.Sizeof(artNode{})))
	//n48 := (*node48)(unsafe.Pointer(memory.AllocZeroed(unsafe.Sizeof(node48{}))))
	n.kind = Node48
	n.ref = uintptr(unsafe.Pointer(n48))
	return n
}

func (f *objFactory) newNode256() *artNode {
	n := (*artNode)(unsafe.Pointer(memory.AllocZeroed(unsafe.Sizeof(artNode{}) + unsafe.Sizeof(node256{}))))
	n256 := (*node256)(unsafe.Pointer(uintptr(unsafe.Pointer(n)) + unsafe.Sizeof(artNode{})))
	//n256 := (*node256)(unsafe.Pointer(memory.AllocZeroed(unsafe.Sizeof(node256{}))))
	n.kind = Node256
	n.ref = uintptr(unsafe.Pointer(n256))
	return n
}

func (f *objFactory) newLeaf(key Key, value interface{}) *artNode {
	clonedKey := key.Clone()
	n := (*artNode)(unsafe.Pointer(memory.AllocZeroed(unsafe.Sizeof(artNode{}) + unsafe.Sizeof(leaf{}))))
	l := (*leaf)(unsafe.Pointer(uintptr(unsafe.Pointer(n)) + unsafe.Sizeof(artNode{})))
	//l := (*leaf)(unsafe.Pointer(memory.AllocZeroed(unsafe.Sizeof(leaf{}))))
	n.kind = Leaf
	l.key = clonedKey
	l.value = value
	n.ref = uintptr(unsafe.Pointer(l))
	//clonedKey := make(Key, len(key))
	//copy(clonedKey, key)
	return n
	//return &artNode{
	//	kind: Leaf,
	//	ref:  unsafe.Pointer(&leaf{key: clonedKey, value: value}),
	//}
}
