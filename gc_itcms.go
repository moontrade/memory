package runtime

import "unsafe"

var (
	itcms_white int32 = 0

	itcmsFromSpace *itcmsObject
	itcmsToSpace   *itcmsObject
	itcmsPinSpace  *itcmsObject
	itcmsIter      *itcmsObject
)

const (
	itcms_gray        = 2
	itcms_transparent = 3
	itcms_COLOR_MASK  = 3

	itcms_GC_INCR_DEBUG = true
)

const (
	// Overhead of a garbage collector object. Excludes memory manager block overhead.
	itcms_OBJECT_OVERHEAD = (unsafe.Sizeof(itcms_OBJECT{}) - tlsf_BLOCK_OVERHEAD + tlsf_AL_MASK) & ^tlsf_AL_MASK

	// Maximum size of a garbage collector object's payload.
	itcms_OBJECT_MAXSIZE = tlsf_BLOCK_MAXSIZE - itcms_OBJECT_OVERHEAD

	// Overhead of a garbage collector object. Excludes memory manager block overhead.
	itcms_TOTAL_OVERHEAD = tlsf_BLOCK_OVERHEAD + itcms_OBJECT_OVERHEAD
)

// Size in memory of all objects currently managed by the GC.
var itcmsTotal uintptr = 0

const (
	itcms_STATE_IDLE  = 0 // Currently transitioning from SWEEP to MARK state.
	itcms_STATE_MARK  = 1 // Currently marking reachable objects.
	itcms_STATE_SWEEP = 2 // Currently sweeping unreachable objects.
)

// Current collector state
var itcmsState = itcms_STATE_IDLE

func initLazy(space *itcmsObject) *itcmsObject {
	*space = itcmsObject{} // zero out
	space.nextWithColor = uintptr(unsafe.Pointer(space))
	space.prev = uintptr(unsafe.Pointer(space))
	return space
}

var defaultPool = NewPool(1)
var _ = initIncrementalGC()

func initIncrementalGC() int {
	itcmsFromSpace = initLazy((*itcmsObject)(defaultPool.Alloc(unsafe.Sizeof(itcmsObject{}))))
	itcmsToSpace = initLazy((*itcmsObject)(defaultPool.Alloc(unsafe.Sizeof(itcmsObject{}))))
	itcmsPinSpace = initLazy((*itcmsObject)(defaultPool.Alloc(unsafe.Sizeof(itcmsObject{}))))
	itcmsIter = nil
	return 0
}

// Visit cookie indicating scanning of an object.
const itcms_VISIT_SCAN = 0

// Garbage collector

// ╒══════════ Garbage collector object layout (32-bit) ═══════════╕
//    3                   2                   1
//  1 0 9 8 7 6 5 4 3 2 1 0 9 8 7 6 5 4 3 2 1 0 9 8 7 6 5 4 3 2 1 0  bits
// ├─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┤
// │                     Memory manager block                      │ -20
// ╞═══════════════════════════════════════════════════════════════╡
// │                            GC info                            │ -16
// ├───────────────────────────────────────────────────────────────┤
// │                            GC info                            │ -12
// ├───────────────────────────────────────────────────────────────┤
// │                            RT id                              │ -8
// ├───────────────────────────────────────────────────────────────┤
// │                            RT size                            │ -4
// ╞>ptr═══════════════════════════════════════════════════════════╡
// │                              ...                              │
type itcms_OBJECT struct {
	tlsfBLOCK
	gcInfo  uint32 // Garbage collector info.
	gcInfo2 uint32 // Garbage collector info.
	rtId    uint32 // Runtime class id
	rtSize  uint32 // Runtime object size
}

// ╒═══════════════ Managed object layout (32-bit) ════════════════╕
//    3                   2                   1
//  1 0 9 8 7 6 5 4 3 2 1 0 9 8 7 6 5 4 3 2 1 0 9 8 7 6 5 4 3 2 1 0  bits
// ├─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┤
// │                      Memory manager block                     │
// ╞═══════════════════════════════════════════════════════════╤═══╡
// │                              next                         │ C │ = nextWithColor
// ├───────────────────────────────────────────────────────────┴───┤
// │                              prev                             │
// ├───────────────────────────────────────────────────────────────┤
// │                              rtId                             │
// ├───────────────────────────────────────────────────────────────┤
// │                              rtSize                           │
// ╞>ptr═══════════════════════════════════════════════════════════╡
// │                               ...                             │
// C: color

// Object Represents a managed object in memory, consisting of a header followed by the object's data.
type itcmsObject struct {
	tlsfBLOCK
	nextWithColor uintptr // Pointer to the next object with color flags stored in the alignment bits.
	prev          uintptr // Pointer to the previous object.
	//rtId          uint32  // Runtime id.
	rtSize uint32 // Runtime size.
}

// Gets the pointer to the next object.
func (o *itcmsObject) next() *itcmsObject {
	return (*itcmsObject)(unsafe.Pointer(o.nextWithColor & ^uintptr(itcms_COLOR_MASK)))
}

// Sets the pointer to the next object.
func (o *itcmsObject) setNext(obj *itcmsObject) {
	o.nextWithColor = uintptr(unsafe.Pointer(obj)) | (o.nextWithColor & itcms_COLOR_MASK)
}

// Gets this object's color.
func (o *itcmsObject) color() int32 {
	return int32(o.nextWithColor & itcms_COLOR_MASK)
}

// Sets this object's color.
func (o *itcmsObject) setColor(color int32) {
	o.nextWithColor = (o.nextWithColor & ^uintptr(itcms_COLOR_MASK)) | uintptr(color)
}

// Gets the size of this object in memory.
func (o *itcmsObject) size() uintptr {
	return tlsf_BLOCK_OVERHEAD + (o.mmInfo & ^uintptr(3))
}

//func (o *Object) isPointerfree() bool {
//	rtId := o.rtId
//	if rtId
//}

func (o *itcmsObject) getPrev() *itcmsObject {
	return (*itcmsObject)(unsafe.Pointer(o.prev))
}

// Unlinks this object from its list.
func (o *itcmsObject) unlink() {
	next := o.next()
	if next == nil {
		if itcms_GC_INCR_DEBUG {
			//tlsfAssert(o.getPrev() == nil && uintptr(unsafe.Pointer(o)) < heapStart, "")
		}
		return
	}
	prev := o.getPrev()
	if itcms_GC_INCR_DEBUG {
		tlsfAssert(prev != nil, "prev is nil")
	}
	next.prev = o.prev
	prev.setNext(next)
}

// Links this object to the specified list, with the given color.
func (o *itcmsObject) linkTo(list *itcmsObject, withColor int32) {
	prev := list.getPrev()
	o.nextWithColor = uintptr(unsafe.Pointer(list)) | uintptr(withColor)
	o.prev = uintptr(unsafe.Pointer(prev))
	prev.setNext(o)
	list.prev = uintptr(unsafe.Pointer(o))
}

func (o *itcmsObject) makeGray() {
	if o == itcmsIter {
		itcmsIter = (*itcmsObject)(unsafe.Pointer(o.prev))
	}
	o.unlink()
	o.linkTo(itcmsToSpace, itcms_gray)
}

func markRoot(addr uintptr, root uintptr) {
	before := itcmsVisitCount
	mark(root)
	after := itcmsVisitCount
	if itcms_GC_INCR_DEBUG {
		if before < after {
			println("marked root:", root, "at", addr)
		} else if addr != 0 {
			println("did not mark root:", root, "at", addr)
		}
	}
}

func markRoots(start uintptr, end uintptr) {
	scan(start, end)
}

// scan loads all pointer-aligned words and marks any pointers that it finds.
func scan(start uintptr, end uintptr) {
	// Align start and end pointers.
	start = (start + unsafe.Alignof(unsafe.Pointer(nil)) - 1) &^ (unsafe.Alignof(unsafe.Pointer(nil)) - 1)
	end &^= unsafe.Alignof(unsafe.Pointer(nil)) - 1

	// Mark all pointers.
	for ptr := start; ptr < end; ptr += unsafe.Alignof(unsafe.Pointer(nil)) {
		mark(*(*uintptr)(unsafe.Pointer(ptr)))
	}
}

// Visits all objects considered to be program roots.
func visitRoots() {
	visitGlobals()
	var (
		pn   = itcmsPinSpace
		iter = pn.next()
	)
	for iter != pn {
		itcmsVisitMembers(uintptr(unsafe.Pointer(iter)) + itcms_TOTAL_OVERHEAD)
		iter = iter.next()
	}
}

// Visits all objects on the stack.
func visitStack() {
	//var ptr = heapStart // TODO __stack_pointer
	//for ptr < heapStart {
	//	mark(*(*uintptr)(unsafe.Pointer(ptr)))
	//	ptr += _PTR_SIZE
	//}
}

func itcmsVisitMembers(ptr uintptr) {
	mark(ptr)
}

func visitGlobals() {
	markRoots(globalsStart, globalsEnd)
}

var itcmsVisitCount = uintptr(0)

func mark(ptr uintptr) {
	if ptr == 0 {
		return
	}
	obj := (*itcmsObject)(unsafe.Pointer(ptr - itcms_TOTAL_OVERHEAD))
	if obj.color() == itcms_white {
		obj.makeGray()
		itcmsVisitCount++
	}
}

func itcmsStep() uintptr {

	var obj *itcmsObject
	switch itcmsState {
	case itcms_STATE_IDLE:
		itcmsState = itcms_STATE_MARK
		itcmsVisitCount = 0
		visitRoots()
		itcmsIter = itcmsToSpace
		return itcmsVisitCount * MARKCOST

	case itcms_STATE_MARK:
		var black int32
		if itcms_white == 1 {
			black = 0
		} else {
			black = 1
		}
		obj = itcmsIter.next()
		for obj != itcmsToSpace {
			itcmsIter = obj
			if obj.color() != black {
				obj.setColor(black)
				itcmsVisitCount = 0
				itcmsVisitMembers(uintptr(unsafe.Pointer(obj)) + itcms_TOTAL_OVERHEAD)
				return itcmsVisitCount * MARKCOST
			}
			obj = obj.next()
		}
		itcmsVisitCount = 0
		visitRoots()
		obj = itcmsIter.next()
		if obj == itcmsToSpace {
			visitStack()
			obj = itcmsIter.next()
			for obj != itcmsToSpace {
				if obj.color() != black {
					obj.setColor(black)
					itcmsVisitMembers(uintptr(unsafe.Pointer(obj)) + itcms_TOTAL_OVERHEAD)
				}
				obj = obj.next()
			}
			from := itcmsFromSpace
			itcmsFromSpace = itcmsToSpace
			itcmsToSpace = from
			itcms_white = black
			itcmsIter = from.next()
			itcmsState = itcms_STATE_SWEEP
		}
		return itcmsVisitCount * MARKCOST

	case itcms_STATE_SWEEP:
		obj = itcmsIter
		if obj != itcmsToSpace {
			itcmsIter = obj.next()
			itcmsFree(obj)
			return SWEEPCOST
		}
		itcmsToSpace.nextWithColor = uintptr(unsafe.Pointer(itcmsToSpace))
		itcmsToSpace.prev = uintptr(unsafe.Pointer(itcmsToSpace))
		itcmsState = itcms_STATE_IDLE
	}
	return 0
}

func itcmsFree(obj *itcmsObject) {
	if uintptr(unsafe.Pointer(obj)) < defaultPool.heapStart {
		obj.nextWithColor = 0 // may become linked again
		obj.prev = 0
	} else {
		itcmsTotal -= obj.size()
		defaultPool.Free(unsafe.Pointer(uintptr(unsafe.Pointer(obj)) + tlsf_BLOCK_OVERHEAD))
	}
}

func itcmsNew(size uintptr) uintptr {
	if size >= itcms_OBJECT_MAXSIZE {
		panic("allocation too large")
	}
	if itcmsTotal >= itcmsThreshold {
		itcmsInterrupt()
	}
	obj := (*itcmsObject)(unsafe.Pointer(uintptr(defaultPool.Alloc(itcms_OBJECT_OVERHEAD+size)) - tlsf_BLOCK_OVERHEAD))
	obj.rtSize = uint32(size)
	obj.linkTo(itcmsFromSpace, itcms_white) // inits next/prev
	itcmsTotal += obj.size()
	ptr := uintptr(unsafe.Pointer(obj)) + itcms_TOTAL_OVERHEAD
	// may be visited before being fully initialized, so must fill
	memzero(unsafe.Pointer(ptr), size)
	return ptr
}

func itcmsLink(parentPtr, childPtr uintptr, expectMultiple bool) {
	if childPtr == 0 {
		return
	}
	child := (*itcmsObject)(unsafe.Pointer(childPtr - itcms_TOTAL_OVERHEAD))
	if child.color() == itcms_white {
		parent := (*itcmsObject)(unsafe.Pointer(parentPtr - itcms_TOTAL_OVERHEAD))
		parentColor := parent.color()
		var notWhite int32
		if itcms_white == 0 {
			notWhite = 1
		} else {
			notWhite = 0
		}
		if parentColor == notWhite {
			// Maintain the invariant that no black object may point to a white object.
			if expectMultiple {
				// Move the barrier "backward". Suitable for containers receiving multiple stores.
				// Avoids a barrier for subsequent objects stored into the same container.
				parent.makeGray()
			} else {
				// Move the barrier "forward". Suitable for objects receiving isolated stores.
				child.makeGray()
			}
		} else if parentColor == itcms_transparent && itcmsState == itcms_STATE_MARK {
			// Pinned objects are considered 'black' during the mark phase.
			child.makeGray()
		}
	}
}

func itcmsPin(ptr uintptr) uintptr {
	if ptr != 0 {
		obj := (*itcmsObject)(unsafe.Pointer(ptr - itcms_TOTAL_OVERHEAD))
		if obj.color() == itcms_transparent {
			panic("already pinned")
		}
		obj.unlink() // from fromSpace
		obj.linkTo(itcmsPinSpace, itcms_transparent)
	}
	return ptr
}

func itcmsUnpin(ptr uintptr) {
	if ptr == 0 {
		return
	}
	obj := (*itcmsObject)(unsafe.Pointer(ptr - itcms_TOTAL_OVERHEAD))
	if obj.color() != itcms_transparent {
		panic("not pinned")
	}
	if itcmsState == itcms_STATE_MARK {
		// We may be right at the point after marking roots for the second time and
		// entering the sweep phase, in which case the object would be missed if it
		// is not only pinned but also a root. Make sure it isn't missed.
		obj.makeGray()
	} else {
		obj.unlink()
		obj.linkTo(itcmsFromSpace, itcms_white)
	}
}

func itcmsCollect() {
	if itcmsState > itcms_STATE_IDLE {
		for itcmsState != itcms_STATE_IDLE {
			itcmsStep()
		}
	}
	itcmsStep()
	for itcmsState != itcms_STATE_IDLE {
		itcmsStep()
	}
	itcmsThreshold = uintptr((uint64(itcmsTotal) * uint64(itcms_IDLEFACTOR) / 100) + uint64(itcms_GRANULARITY))
}

// Magic constants responsible for pause times. Obtained experimentally
// using the compiler compiling itself. 2048 budget pro run by default.
const MARKCOST = uintptr(1)
const SWEEPCOST = uintptr(10)

// How often to interrupt. The default of 1024 means “interrupt each 1024 bytes allocated”.
const itcms_GRANULARITY uintptr = 1024

// How long to interrupt. The default of 200% means “run at double the speed of allocations”.
const itcms_STEPFACTOR uintptr = 200

// How long to idle. The default of 200% means “wait for memory to double before kicking in again”.
const itcms_IDLEFACTOR uintptr = 200

// Threshold of memory used by objects to exceed before interrupting again.
var itcmsThreshold uintptr = 1024 //((uintptr(memoryPageCount()) << 16) - heapStart) >> 1

// Performs a reasonable amount of incremental GC steps.
func itcmsInterrupt() {
	budget := itcms_GRANULARITY * itcms_STEPFACTOR / 100
LOOP:
	for {
		budget -= itcmsStep()
		if itcmsState == itcms_STATE_IDLE {
			itcmsThreshold = uintptr((uint64(itcmsTotal) * uint64(itcms_IDLEFACTOR) / 100) + uint64(itcms_GRANULARITY))
		}
		//if budget <= 0 {
		break LOOP
		//}
	}
	if itcmsThreshold < itcms_GRANULARITY {
		itcmsThreshold = itcmsTotal + itcms_GRANULARITY*(itcmsTotal-1)
	} else {
		itcmsThreshold = itcmsTotal + itcms_GRANULARITY*itcmsTotal
	}
}
