package mem

import (
	"unsafe"
)

// Auto is a singly linked list (Stack) of nodes which contain pointers
// that will all free when Free is called. uintptr is used to ensure the
// compiler doesn't confuse it for a Go GC managed pointer.
type Auto uintptr

//goland:noinspection GoVetUnsafePointer
func (a Auto) Allocator() Allocator {
	return ((*autoHead)(unsafe.Pointer(a))).alloc
}

//goland:noinspection GoVetUnsafePointer
func (a Auto) Count() Pointer {
	return ((*autoHead)(unsafe.Pointer(a))).count
}

//goland:noinspection GoVetUnsafePointer
func (a Auto) Size() Pointer {
	return ((*autoHead)(unsafe.Pointer(a))).bytes
}

//goland:noinspection GoVetUnsafePointer
func (a Auto) Max() Pointer {
	return ((*autoHead)(unsafe.Pointer(a))).max
}

type autoHead struct {
	head  Pointer   // pointer to the head node
	alloc Allocator // allocator instance used
	max   Pointer   // max number of entries per node
	count Pointer
	bytes Pointer
}

type autoNode struct {
	len   Pointer
	next  Pointer
	first struct{}
}

//goland:noinspection GoVetUnsafePointer
func NewAuto(a Allocator, nodeSize Pointer) Auto {
	if nodeSize == 0 {
		nodeSize = 32
	}
	p := a.Alloc(Pointer(unsafe.Sizeof(autoHead{}) + unsafe.Sizeof(autoNode{}) + (uintptr(nodeSize) * unsafe.Sizeof(Pointer(0)))))
	h := (*autoHead)(p.Unsafe())
	h.head = Pointer(uintptr(p) + unsafe.Sizeof(autoHead{}))
	h.alloc = a
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
func (au Auto) Alloc(size Pointer) Pointer {
	if au == 0 {
		return Pointer(0)
	}
	h := (*autoHead)(unsafe.Pointer(au))
	p := h.alloc.Alloc(size)
	au.add(p)
	return p
}

//goland:noinspection GoVetUnsafePointer
func (au Auto) Bytes(length, capacity Pointer) Bytes {
	if au == 0 {
		return Bytes{}
	}
	h := (*autoHead)(unsafe.Pointer(au))
	p := h.alloc.Bytes(length)
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
		nextPtr := h.alloc.Alloc(Pointer(unsafe.Sizeof(autoNode{}) + (uintptr(h.max) * unsafe.Sizeof(Pointer(0)))))
		h.bytes += allocationSize(nextPtr) + allocationSize(ptr)
		next := (*autoNode)(nextPtr.Unsafe())
		// Add length to 1
		next.len = 1
		// Link to current n
		next.next = Pointer(unsafe.Pointer(n))
		// Update reference to new n
		h.head = nextPtr
		// Add first item
		*(*uintptr)(unsafe.Pointer(&next.first)) = uintptr(ptr)
	} else {
		h.bytes += allocationSize(ptr)
		// Add item
		*(*Pointer)(unsafe.Pointer(uintptr(unsafe.Pointer(&n.first)) + (uintptr(n.len) * unsafe.Sizeof(uintptr(0))))) = ptr
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
	a := head.alloc
	for n != nil {
		var (
			start = Pointer(unsafe.Pointer(&n.first))
			end   = start + (n.len * Pointer(unsafe.Sizeof(Pointer(0))))
			item  Pointer
		)
		for i := start; i < end; i += Pointer(unsafe.Sizeof(Pointer(0))) {
			item = *(*Pointer)(unsafe.Pointer(i))
			if item == 0 {
				break
			}
			a.Free(item)
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

	println("Auto =>", " count:", head.count, " bytes:", head.bytes, " addr:", uint(Pointer(unsafe.Pointer(head))))
	count := -1
	for n != nil {
		count++
		var (
			start = Pointer(unsafe.Pointer(&n.first))
			end   = start + (n.len * Pointer(unsafe.Sizeof(Pointer(0))))
			item  Pointer
			index Pointer
		)
		println("\t[", count, "] ->", n.len, "items")
		for i := start; i < end; i += Pointer(unsafe.Sizeof(Pointer(0))) {
			item = *(*Pointer)(unsafe.Pointer(i))
			if item == 0 {
				break
			}
			index = (i - start) / Pointer(unsafe.Sizeof(Pointer(0)))
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
