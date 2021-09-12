package mem

import (
	"unsafe"
)

// Auto is a singly linked list (Stack) of nodes which contain pointers
// that will all free when Free is called. uintptr is used to ensure the
// compiler doesn't confuse it for a Go GC managed pointer.
type Auto uintptr

//goland:noinspection GoVetUnsafePointer
func (a Auto) Allocator() *Allocator {
	return (*Allocator)(unsafe.Pointer(((*autoHead)(unsafe.Pointer(a))).allocator))
}

//goland:noinspection GoVetUnsafePointer
func (a Auto) Count() uintptr {
	return ((*autoHead)(unsafe.Pointer(a))).count
}

//goland:noinspection GoVetUnsafePointer
func (a Auto) Size() uintptr {
	return ((*autoHead)(unsafe.Pointer(a))).bytes
}

//goland:noinspection GoVetUnsafePointer
func (a Auto) Max() uintptr {
	return ((*autoHead)(unsafe.Pointer(a))).max
}

type autoHead struct {
	head      uintptr // pointer to the head node
	allocator uintptr // allocator instance used
	max       uintptr // max number of entries per node
	count     uintptr
	bytes     uintptr
}

type autoNode struct {
	len   uintptr
	next  uintptr
	first struct{}
}

//goland:noinspection GoVetUnsafePointer
func NewAuto(a *Allocator, nodeSize uintptr) Auto {
	if nodeSize == 0 {
		nodeSize = 32
	}
	p := a.Alloc(unsafe.Sizeof(autoHead{}) + unsafe.Sizeof(autoNode{}) + (nodeSize * unsafe.Sizeof(uintptr(0))))
	h := (*autoHead)(p.Unsafe())
	h.head = uintptr(p) + unsafe.Sizeof(autoHead{})
	h.allocator = uintptr(unsafe.Pointer(a))
	h.max = nodeSize
	h.count = 0
	h.bytes = allocationSize(p)
	n := (*autoNode)(unsafe.Pointer(h.head))
	n.len = 0
	n.next = 0
	return Auto(p)
}

func (au *Auto) Scope(fn func(Auto)) {
	if fn != nil {
		fn(*au)
	}
	au.Free()
}

//goland:noinspection GoVetUnsafePointer
func (au *Auto) HasNext() bool {
	return *(*uintptr)(unsafe.Pointer(*au)) != 0
}

//goland:noinspection GoVetUnsafePointer
func (au *Auto) Next() Auto {
	if au == nil {
		return 0
	}
	p := uintptr(*au)
	return Auto(*(*uintptr)(unsafe.Pointer(p)))
}

//goland:noinspection GoVetUnsafePointer
func (au Auto) Alloc(size uintptr) Pointer {
	if au == 0 {
		return Pointer(0)
	}
	h := (*autoHead)(unsafe.Pointer(au))
	p := ((*Allocator)(unsafe.Pointer(h.allocator))).Alloc(size)
	au.add(p)
	return p
}

//goland:noinspection GoVetUnsafePointer
func (au Auto) Bytes(length, capacity uintptr) Bytes {
	if au == 0 {
		return Bytes{}
	}
	h := (*autoHead)(unsafe.Pointer(au))
	p := ((*Allocator)(unsafe.Pointer(h.allocator))).Bytes(length, capacity)
	au.add(p.Pointer)
	return p
}

//goland:noinspection GoVetUnsafePointer
func (au *Auto) add(ptr Pointer) {
	if au == nil {
		return
	}
	h := (*autoHead)(unsafe.Pointer(*au))
	n := (*autoNode)(unsafe.Pointer(h.head))
	if n == nil {
		return
	}
	if n.len == h.max {
		nextPtr := ((*Allocator)(unsafe.Pointer(h.allocator))).Alloc(unsafe.Sizeof(autoNode{}) + (h.max * unsafe.Sizeof(uintptr(0))))
		h.bytes += allocationSize(nextPtr) + allocationSize(ptr)
		next := (*autoNode)(nextPtr.Unsafe())
		// Add length to 1
		next.len = 1
		// Link to current n
		next.next = uintptr(unsafe.Pointer(n))
		// Update reference to new n
		h.head = uintptr(nextPtr)
		// Add first item
		*(*uintptr)(unsafe.Pointer(&next.first)) = uintptr(ptr)
	} else {
		h.bytes += allocationSize(ptr)
		// Add item
		*(*uintptr)(unsafe.Pointer(uintptr(unsafe.Pointer(&n.first)) + (n.len * unsafe.Sizeof(uintptr(0))))) = uintptr(ptr)
		// Increment length
		n.len++
	}
	h.count++
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
//goland:noinspection GoVetUnsafePointer
func (au *Auto) Free() {
	if au == nil {
		return
	}
	head := (*autoHead)(unsafe.Pointer(*au))
	n := (*autoNode)(unsafe.Pointer(head.head))
	if n == nil {
		return
	}
	a := (*Allocator)(unsafe.Pointer(head.allocator))
	for n != nil {
		var (
			start = uintptr(unsafe.Pointer(&n.first))
			end   = start + (n.len * unsafe.Sizeof(uintptr(0)))
			item  uintptr
		)
		for i := start; i < end; i += unsafe.Sizeof(uintptr(0)) {
			item = *(*uintptr)(unsafe.Pointer(i))
			if item == 0 {
				break
			}
			a.Free(Pointer(item))
		}

		if n.next == 0 {
			// Free header node
			a.Free(Pointer(uintptr(unsafe.Pointer(n)) - unsafe.Sizeof(autoHead{})))
			break
		}

		a.Free(Pointer(unsafe.Pointer(n)))
		n = (*autoNode)(unsafe.Pointer(n.next))
	}
	*au = 0
}

//goland:noinspection GoVetUnsafePointer
func (au *Auto) Print() {
	head := (*autoHead)(unsafe.Pointer(*au))
	n := (*autoNode)(unsafe.Pointer(head.head))
	if n == nil {
		return
	}

	println("Auto =>", " count:", head.count, " bytes:", head.bytes, " addr:", uint(uintptr(unsafe.Pointer(head))))
	count := -1
	for n != nil {
		count++
		var (
			start = uintptr(unsafe.Pointer(&n.first))
			end   = start + (n.len * unsafe.Sizeof(uintptr(0)))
			item  uintptr
			index uintptr
		)
		println("\t[", count, "] ->", n.len, "items")
		for i := start; i < end; i += unsafe.Sizeof(uintptr(0)) {
			item = *(*uintptr)(unsafe.Pointer(i))
			if item == 0 {
				break
			}
			index = (i - start) / unsafe.Sizeof(uintptr(0))
			space := ""
			if index < 10 {
				space = "   "
			} else if index < 100 {
				space = "  "
			} else if index < 1000 {
				space = " "
			}
			println("\t\t", space, uint(index), "->", uint(item))
		}

		if n.next == 0 {
			break
		}
		n = (*autoNode)(unsafe.Pointer(n.next))
	}
}
