package memory

import (
	"unsafe"
)

// Auto is a singly linked list (Stack) of nodes which contain pointers
// that will all free when Free is called. uintptr is used to ensure the
// compiler doesn't confuse it for a Go GC managed pointer.
type Auto uintptr

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
	head  Pointer // pointer to the head node
	max   uintptr // max number of entries per node
	count uintptr
	bytes uintptr
}

type autoNode struct {
	len   uintptr
	next  Pointer
	first struct{}
}

//goland:noinspection GoVetUnsafePointer
func NewAuto(nodeSize uintptr) Auto {
	if nodeSize == 0 {
		nodeSize = 32
	}
	p, s := AllocCap(unsafe.Sizeof(autoHead{}) + unsafe.Sizeof(autoNode{}) + (nodeSize * unsafe.Sizeof(Pointer(0))))
	h := (*autoHead)(p.Unsafe())
	h.head = Pointer(uintptr(p) + unsafe.Sizeof(autoHead{}))
	h.max = nodeSize
	h.count = 0
	h.bytes = s
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
	p := Alloc(size)
	au.add(p)
	return p
}

//goland:noinspection GoVetUnsafePointer
func (au Auto) Str(size uintptr) Bytes {
	if au == 0 {
		return Bytes{}
	}
	//p := h.alloc.Str(size)
	//au.add(p.allocationPointer())
	//return p
	return Bytes{}
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
		nextPtr, sz := AllocCap(unsafe.Sizeof(autoNode{}) + (h.max * unsafe.Sizeof(Pointer(0))))
		h.bytes += sz + SizeOf(ptr)
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
		h.bytes += SizeOf(ptr)
		// Add item
		*(*Pointer)(unsafe.Pointer(uintptr(unsafe.Pointer(&n.first)) + (n.len * unsafe.Sizeof(uintptr(0))))) = ptr
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
	for n != nil {
		var (
			start = uintptr(unsafe.Pointer(&n.first))
			end   = start + (n.len * unsafe.Sizeof(Pointer(0)))
			item  Pointer
		)
		for i := start; i < end; i += unsafe.Sizeof(Pointer(0)) {
			item = *(*Pointer)(unsafe.Pointer(i))
			if item == 0 {
				break
			}
			Free(item)
		}

		if n.next == 0 {
			// Free header node
			Free(Pointer(uintptr(unsafe.Pointer(n)) - unsafe.Sizeof(autoHead{})))
			break
		}

		Free(Pointer(unsafe.Pointer(n)))
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
			start = uintptr(unsafe.Pointer(&n.first))
			end   = start + (n.len * unsafe.Sizeof(Pointer(0)))
			item  Pointer
			index uintptr
		)
		println("\t[", count, "] ->", n.len, "items")
		for i := start; i < end; i += unsafe.Sizeof(Pointer(0)) {
			item = *(*Pointer)(unsafe.Pointer(i))
			if item == 0 {
				break
			}
			index = (i - start) / unsafe.Sizeof(Pointer(0))
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
