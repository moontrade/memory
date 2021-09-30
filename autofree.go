package memory

import (
	"unsafe"
)

// AutoFree is a singly linked list (Stack) of nodes which contain pointers
// that will all free when Free is called. uintptr is used to ensure the
// compiler doesn't confuse it for a Go GC managed pointer.
type AutoFree uintptr

//goland:noinspection GoVetUnsafePointer
func (af AutoFree) Count() uintptr {
	return ((*autoHead)(unsafe.Pointer(af))).count
}

//goland:noinspection GoVetUnsafePointer
func (af AutoFree) Size() uintptr {
	return ((*autoHead)(unsafe.Pointer(af))).bytes
}

//goland:noinspection GoVetUnsafePointer
func (af AutoFree) Max() uintptr {
	return ((*autoHead)(unsafe.Pointer(af))).max
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
func NewAuto(nodeSize uintptr) AutoFree {
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
	return AutoFree(p)
}

func (af *AutoFree) Scope(fn func(AutoFree)) {
	if fn != nil {
		fn(*af)
	}
	af.Free()
}

//goland:noinspection GoVetUnsafePointer
func (af *AutoFree) HasNext() bool {
	return *(*uintptr)(unsafe.Pointer(*af)) != 0
}

//goland:noinspection GoVetUnsafePointer
func (af *AutoFree) Next() AutoFree {
	if af == nil {
		return 0
	}
	p := uintptr(*af)
	return AutoFree(*(*uintptr)(unsafe.Pointer(p)))
}

//goland:noinspection GoVetUnsafePointer
func (af AutoFree) Alloc(size uintptr) Pointer {
	if af == 0 {
		return Pointer(0)
	}
	p := Alloc(size)
	af.add(p)
	return p
}

func (af AutoFree) AllocCap(size uintptr) FatPointer {
	p, c := AllocCap(size)
	if p == 0 {
		return FatPointer{}
	}
	af.add(p)
	return FatPointer{
		Pointer: p,
		len:     c,
	}
}

//goland:noinspection GoVetUnsafePointer
func (af AutoFree) Bytes(size uintptr) Bytes {
	if af == 0 {
		return Bytes{}
	}
	b := AllocBytes(size)
	af.add(b.allocationPointer())
	return b
}

//goland:noinspection GoVetUnsafePointer
func (af *AutoFree) add(ptr Pointer) {
	if af == nil {
		return
	}
	h := (*autoHead)(unsafe.Pointer(*af))
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
func (af *AutoFree) Close() error {
	if af == nil {
		return nil
	}
	af.Free()
	return nil
}

// Free releases every allocation
//goland:noinspection GoVetUnsafePointer
func (af *AutoFree) Free() {
	if af == nil {
		return
	}
	head := (*autoHead)(unsafe.Pointer(*af))
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
	*af = 0
}

//goland:noinspection GoVetUnsafePointer
func (af *AutoFree) Print() {
	head := (*autoHead)(unsafe.Pointer(*af))
	n := (*autoNode)(unsafe.Pointer(head.head))
	if n == nil {
		return
	}

	println("AutoFree =>", " count:", head.count, " bytes:", head.bytes, " addr:", uint(Pointer(unsafe.Pointer(head))))
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
