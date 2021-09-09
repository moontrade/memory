package mem

import (
	"unsafe"
)

const (
	sizeofPointer       = unsafe.Sizeof(uintptr(0))
	autoNodeMax         = uintptr(24)
	autoSizeofNode      = sizeofPointer + (autoNodeMax * sizeofPointer) + sizeofPointer
	autoFirstOffset     = sizeofPointer + sizeofPointer
	autoNextOffset      = sizeofPointer + (autoNodeMax * sizeofPointer)
	autoFirstNextOffset = autoFirstOffset + (autoNodeMax * sizeofPointer) + sizeofPointer

	audoSizeOfNode      = unsafe.Sizeof(node{}) + (sizeofPointer * autoNodeMax)
	autoSizeOfFirstNode = unsafe.Sizeof(header{}) + audoSizeOfNode
)

// Auto is a linked list of pointers that will all free when Free is called.
type Auto uintptr

type header struct {
	head      uintptr
	allocator uintptr
}

type node struct {
	len   uintptr
	prev  uintptr
	first struct{}
}

func (a *Allocator) Scope(fn func(a Auto)) {
	au := NewAuto(a)
	defer au.Free()
	fn(au)
}

func NewAuto(a *Allocator) Auto {
	p := a.Alloc(autoSizeOfFirstNode)
	h := (*header)(p)
	h.head = uintptr(p) + unsafe.Sizeof(header{})
	h.allocator = uintptr(unsafe.Pointer(a))
	n := (*node)(unsafe.Pointer(h.head))
	n.len = 0
	n.prev = 0
	return Auto(p)
}

func (au Auto) newAuto() Auto {
	h := (*header)(unsafe.Pointer(au))
	p := ((*Allocator)(unsafe.Pointer(h.allocator))).Alloc(audoSizeOfNode)
	next := (*node)(p)
	next.len = 0
	next.prev = 0
	h.head = uintptr(p)
	return Auto(p)
}

func (au *Auto) Scope(fn func(Auto)) {
	if fn != nil {
		fn(*au)
	}
	au.Free()
}

//goland:noinspection ALL
func (au *Auto) HasNext() bool {
	return *(*uintptr)(unsafe.Pointer(*au)) != 0
}

//goland:noinspection ALL
func (au *Auto) Next() Auto {
	if au == nil {
		return 0
	}
	p := uintptr(*au)
	return Auto(*(*uintptr)(unsafe.Pointer(p)))
}

func (au Auto) Alloc(size uintptr) unsafe.Pointer {
	if au == 0 {
		return nil
	}
	h := (*header)(unsafe.Pointer(au))
	p := ((*Allocator)(unsafe.Pointer(h.allocator))).Alloc(size)
	au.Add(p)
	return p
}

//goland:noinspection ALL
func (au *Auto) Add(ptr unsafe.Pointer) {
	if au == nil {
		return
	}
	head := (*header)(unsafe.Pointer(uintptr(*au)))
	n := (*node)(unsafe.Pointer(head.head))
	if n == nil {
		return
	}
	if n.len == autoNodeMax {
		nextPtr := ((*Allocator)(unsafe.Pointer(head.allocator))).Alloc(audoSizeOfNode)
		next := (*node)(nextPtr)
		// Set length to 1
		next.len = 1
		// Link to previous n
		next.prev = uintptr(unsafe.Pointer(n))
		// Update reference to new n
		head.head = uintptr(nextPtr)
		// Set first item
		*(*uintptr)(unsafe.Pointer(uintptr(nextPtr) + autoFirstOffset)) = uintptr(ptr)
	} else {
		// Add item
		*(*uintptr)(unsafe.Pointer(head.head + autoFirstOffset + (n.len * sizeofPointer))) = uintptr(ptr)
		// Increment length
		n.len++
	}
}

// Close releases / frees every allocation
func (au *Auto) Close() error {
	if au == nil {
		return nil
	}
	au.Free()
	return nil
}

// Free releases every allocation
//goland:noinspection ALL
func (au *Auto) Free() {
	if au == nil {
		return
	}
	//au.Print()
	head := (*header)(unsafe.Pointer(uintptr(*au)))
	n := (*node)(unsafe.Pointer(head.head))
	if n == nil {
		return
	}
	//for n != nil {
	//	length := *(*uintptr)(unsafe.Pointer(n))
	//	end := n + autoFirstOffset + (length * sizeofPointer)
	//	var item uintptr
	//	for i := n + autoFirstOffset; i < end; i += sizeofPointer {
	//		item = *(*uintptr)(unsafe.Pointer(i))
	//		if item != 0 {
	//			Free(unsafe.Pointer(item))
	//		}
	//	}
	//
	//	next := *(*uintptr)(unsafe.Pointer(n + autoNextOffset))
	//	if next == 0 {
	//		Free(unsafe.Pointer(n - sizeofPointer))
	//		break
	//	} else {
	//		Free(unsafe.Pointer(n))
	//		n = next
	//	}
	//}
	*au = 0
}

//goland:noinspection ALL
func (au *Auto) Print() {
	node := uintptr(*au)
	if node == 0 {
		return
	}
	node = *(*uintptr)(unsafe.Pointer(node - sizeofPointer))
	for node != 0 {
		length := *(*uintptr)(unsafe.Pointer(node))
		println("")
		println("AUTOFREE", "len", uint(length))
		end := node + autoFirstOffset + (length * sizeofPointer)
		var item uintptr
		for i := node + autoFirstOffset; i < end; i += sizeofPointer {
			item = *(*uintptr)(unsafe.Pointer(i))
			index := ((i - node) / sizeofPointer) - 1
			if item != 0 {
				println("\t", uint(index), "->", uint(item))
			}
		}
		node = *(*uintptr)(unsafe.Pointer(node + autoNextOffset))
	}
}
