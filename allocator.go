package mem

import (
	"math"
	"math/bits"
	"unsafe"
)

// Allocator === TLSF (Two-Level Segregate Fit) memory allocator ===
//
// TLSF is a general purpose dynamic memory allocator specifically designed to meet
// real-time requirements:
//
// 		Bounded Response Time - The worst-case execution time (WCET) of memory allocation
//								and deallocation Has got to be known in advance and be
//								independent of application data. Allocator Has a constant
//								cost O(1).
//
//						 Fast - Additionally to a bounded cost, the allocator Has to be
//								efficient and fast enough. Allocator executes a maximum
//								of 168 processor instructions in a x86 architecture.
//								Depending on the compiler version and optimisation flags,
//								it can be slightly lower or higher.
//
// 		Efficient Memory Use - 	Traditionally, real-time systems run for long periods of
//								time and some (embedded applications), have strong constraints
//								of memory size. Fragmentation can have a significant impact on
//								such systems. It can increase  dramatically, and degrade the
//								system performance. A way to measure this efficiency is the
//								memory fragmentation incurred by the allocator. Allocator has
//								been tested in hundreds of different loads (real-time tasks,
//								general purpose applications, etc.) obtaining an average
//								fragmentation lower than 15 %. The maximum fragmentation
//								measured is lower than 25%.
//
// Memory can be added on demand and is a multiple of 64kb pages. Grow is used to allocate new
// memory to be added to the allocator. Each Grow must provide a contiguous chunk of memory.
// However, the allocator may be comprised of many contiguous chunks which are not contiguous
// of each other. There is not a mechanism for shrinking the memory. Supplied Grow function
// can effectively limit how big the allocator can get. If a zero pointer is returned it will
// cause an out-of-memory situation which is propagated as a nil pointer being returned from
// Alloc. It's up to the application to decide how to handle such scenarios.
//
// see: http://www.gii.upv.es/tlsf/
//
// - `ffs(x)` is equivalent to `ctz(x)` with x != 0
// - `fls(x)` is equivalent to `sizeof(x) * 8 - clz(x) - 1`
//
// ╒══════════════ Block size interpretation (32-bit) ═════════════╕
//    3                   2                   1
//  1 0 9 8 7 6 5 4 3 2 1 0 9 8 7 6 5 4 3 2 1 0 9 8 7 6 5 4 3 2 1 0  bits
// ├─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┼─┴─┴─┴─╫─┴─┴─┴─┤
// │ |                    FL                       │ SB = SL + AL  │ ◄─ usize
// └───────────────────────────────────────────────┴───────╨───────┘
// FL: first level, SL: second level, AL: alignment, SB: small block
type Allocator struct {
	root      *root
	HeapStart uintptr
	HeapEnd   uintptr
	arena     uintptr
	Grow      Grow
	Stats
}

// Stats provides the metrics of an Allocator
type Stats struct {
	HeapSize        int64
	AllocSize       int64
	MaxUsedSize     int64
	FreeSize        int64
	Allocs          int32
	InitialPages    int32
	ConsecutiveLow  int32
	ConsecutiveHigh int32
	Pages           int32
	Grows           int32
	fragmentation   float32
}

func (s *Stats) Fragmentation() float32 {
	if s.HeapSize == 0 || s.MaxUsedSize == 0 {
		return 0
	}
	pct := float64(s.HeapSize-s.MaxUsedSize) / float64(s.HeapSize)
	s.fragmentation = float32(math.Floor(pct*100) / 100)
	return s.fragmentation
}

// Grow provides the ability to Grow the heap and allocate a contiguous
// chunk of system memory to add to the allocator.
type Grow func(pagesBefore, pagesNeeded int32, minSize uintptr) (pagesAdded int32, start, end uintptr)

const (
	_TLSFPageSize = uintptr(64 * 1024)

	_TLSFAlignU32 = 2
	// All allocation sizes and addresses are aligned to 4 or 8 bytes.
	// 32bit = 2
	// 64bit = 3
	// <expr> = bits.UintSize / 8 / 4 + 1
	_TLSFAlignSizeLog2 = ((32 << (^uint(0) >> 63)) / 8 / 4) + 1
	_TLSFSizeofPointer = unsafe.Sizeof(uintptr(0))

	_TLSFALBits uint32  = 4 // 16 bytes to fit up to v128
	_TLSFALSize uintptr = 1 << uintptr(_TLSFALBits)
	_TLSFALMask         = _TLSFALSize - 1

	// Overhead of a memory manager block.
	_TLSFBlockOverhead = unsafe.Sizeof(tlsfBLOCK{})
	// Block constants. A block must have a minimum size of three pointers so it can hold `prev`,
	// `prev` and `back` if free.
	_TLSFBlockMinSize = ((3*_TLSFSizeofPointer + _TLSFBlockOverhead + _TLSFALMask) & ^_TLSFALMask) - _TLSFBlockOverhead
	// Maximum size of a memory manager block's payload.
	_TLSFBlockMaxSize = (1 << 30) - _TLSFBlockOverhead
	//_TLSFBlockMaxSize = (1 << ((_TLSFAlignSizeLog2 + 1)*10)) - _TLSFBlockOverhead

	_TLSFDebug = false

	_TLSFSLBits uint32 = 4
	_TLSFSLSize uint32 = 1 << _TLSFSLBits
	_TLSFSBBits        = _TLSFSLBits + _TLSFALBits
	_TLSFSBSize uint32 = 1 << _TLSFSBBits
	_TLSFFLBits        = 31 - _TLSFSBBits

	// [00]: < 256B (SB)  [12]: < 1M
	// [01]: < 512B       [13]: < 2M
	// [02]: < 1K         [14]: < 4M
	// [03]: < 2K         [15]: < 8M
	// [04]: < 4K         [16]: < 16M
	// [05]: < 8K         [17]: < 32M
	// [06]: < 16K        [18]: < 64M
	// [07]: < 32K        [19]: < 128M
	// [08]: < 64K        [20]: < 256M
	// [09]: < 128K       [21]: < 512M
	// [10]: < 256K       [22]: <= 1G - OVERHEAD
	// [11]: < 512K
	// WASM VMs limit to 2GB total (currently), making one 1G block max
	// (or three 512M etc.) due to block overhead

	// Tags stored in otherwise unused alignment bits
	_TLSFFREE     uintptr = 1 << 0
	_TLSFLEFTFREE uintptr = 1 << 1
	_TLSFTagsMask         = _TLSFFREE | _TLSFLEFTFREE
)

func (a *Allocator) Bytes(length, capacity uintptr) Bytes {
	if capacity < length {
		capacity = length
	}
	p := uintptr(unsafe.Pointer(a.allocateBlock(capacity)))
	if p == 0 {
		return Bytes{}
	}
	return Bytes{
		addr:  p + _TLSFBlockOverhead,
		len:   uint32(length),
		cap:   uint32(*(*uintptr)(unsafe.Pointer(p)) & ^_TLSFTagsMask),
		alloc: uintptr(unsafe.Pointer(a)),
	}
}

func (a *Allocator) BytesMut(length, capacity uintptr) BytesMut {
	b := a.Bytes(length, capacity)
	if b.addr == 0 {
		return BytesMut{}
	}
	return BytesMut{
		Bytes: b,
	}
}

// Alloc allocates a block of memory that fits the size provided
//goland:noinspection GoVetUnsafePointer
func (a *Allocator) Alloc(size uintptr) unsafe.Pointer {
	p := uintptr(unsafe.Pointer(a.allocateBlock(size)))
	if p == 0 {
		return nil
	}
	return unsafe.Pointer(p + _TLSFBlockOverhead)
}

// Realloc determines the best way to resize an allocation.
func (a *Allocator) Realloc(ptr unsafe.Pointer, size uintptr) unsafe.Pointer {
	p := uintptr(unsafe.Pointer(a.moveBlock(checkUsedBlock(uintptr(ptr)), size)))
	if p == 0 {
		return nil
	}
	return unsafe.Pointer(p + _TLSFBlockOverhead)
}

// Free release the allocation back into the free list.
func (a *Allocator) Free(ptr unsafe.Pointer) {
	a.freeBlock((*tlsfBlock)(unsafe.Pointer(uintptr(ptr) - _TLSFBlockOverhead)))
	//a.freeBlock(checkUsedBlock(uintptr(ptr)))
}

// Scope creates an Auto free list that automatically reclaims memory
// after callback finishes.
func (a *Allocator) Scope(fn func(a Auto)) {
	if fn == nil {
		return
	}
	auto := NewAuto(a, 32)
	defer auto.Free()
	fn(auto)
}

// bootstrap bootstraps the Allocator with the initial block of contiguous memory
// that at least fits the minimum required to fit the bitmap.
//goland:noinspection GoVetUnsafePointer
func bootstrap(start, end uintptr, pages int32, grow Grow) *Allocator {
	start = (start + unsafe.Alignof(unsafe.Pointer(nil)) - 1) &^ (unsafe.Alignof(unsafe.Pointer(nil)) - 1)

	//if a.T {
	//	println("bootstrap", "pages", pages, uint(start), uint(end), uint(end-start))
	//}
	// init allocator
	a := (*Allocator)(unsafe.Pointer(start))
	*a = Allocator{
		HeapStart: start,
		HeapEnd:   end,
		Stats: Stats{
			InitialPages: pages,
			Pages:        pages,
		},
		Grow: grow,
	}

	// init root
	rootOffset := unsafe.Sizeof(Allocator{}) + ((start + _TLSFALMask) & ^_TLSFALMask)
	a.root = (*root)(unsafe.Pointer(rootOffset))
	a.root.init()

	// add initial memory
	a.addMemory(rootOffset+_TLSFRootSize, end)
	return a
}

// Memory manager

// ╒════════════ Memory manager block layout (32-bit) ═════════════╕
//    3                   2                   1
//  1 0 9 8 7 6 5 4 3 2 1 0 9 8 7 6 5 4 3 2 1 0 9 8 7 6 5 4 3 2 1 0  bits
// ├─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┤
// │                           MM info                             │ -4
// ╞>ptr═══════════════════════════════════════════════════════════╡
// │                              ...                              │
type tlsfBLOCK struct {
	mmInfo uintptr
}

// ╒════════════════════ Block layout (32-bit) ════════════════════╕
//    3                   2                   1
//  1 0 9 8 7 6 5 4 3 2 1 0 9 8 7 6 5 4 3 2 1 0 9 8 7 6 5 4 3 2 1 0  bits
// ├─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┼─┼─┤            ┐
// │                          size                             │L│F│ ◄─┐ info   overhead
// ╞>ptr═══════════════════════════════════════════════════════╧═╧═╡   │        ┘
// │                        if free: ◄ prev                        │ ◄─┤ usize
// ├───────────────────────────────────────────────────────────────┤   │
// │                        if free: next ►                        │ ◄─┤
// ├───────────────────────────────────────────────────────────────┤   │
// │                             ...                               │   │ >= 0
// ├───────────────────────────────────────────────────────────────┤   │
// │                        if free: back ▲                        │ ◄─┘
// └───────────────────────────────────────────────────────────────┘ >= MIN SIZE
// F: FREE, L: LEFTFREE
type tlsfBlock struct {
	tlsfBLOCK
	// Previous free block, if any. Only valid if free, otherwise part of payload.
	//prev *Block
	prev uintptr
	// Next free block, if any. Only valid if free, otherwise part of payload.
	//next *Block
	next uintptr

	// If the block is free, there is a 'back'reference at its end pointing at its start.
}

// Gets the left block of a block. Only valid if the left block is free.
func (block *tlsfBlock) getFreeLeft() *tlsfBlock {
	return *(**tlsfBlock)(unsafe.Pointer(uintptr(unsafe.Pointer(block)) - _TLSFSizeofPointer))
}

// Gets the right block of a block by advancing to the right by its size.
func (block *tlsfBlock) getRight() *tlsfBlock {
	return (*tlsfBlock)(unsafe.Pointer(uintptr(unsafe.Pointer(block)) + _TLSFBlockOverhead + (block.mmInfo & ^_TLSFTagsMask)))
}

// ╒═════════════════════ Root layout (32-bit) ════════════════════╕
//    3                   2                   1
//  1 0 9 8 7 6 5 4 3 2 1 0 9 8 7 6 5 4 3 2 1 0 9 8 7 6 5 4 3 2 1 0  bits
// ├─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┤          ┐
// │        0        |           flMap                            S│ ◄────┐
// ╞═══════════════════════════════════════════════════════════════╡      │
// │                           slMap[0] S                          │ ◄─┐  │
// ├───────────────────────────────────────────────────────────────┤   │  │
// │                           slMap[1]                            │ ◄─┤  │
// ├───────────────────────────────────────────────────────────────┤  uint32 │
// │                           slMap[22]                           │ ◄─┘  │
// ╞═══════════════════════════════════════════════════════════════╡    usize
// │                            head[0]                            │ ◄────┤
// ├───────────────────────────────────────────────────────────────┤      │
// │                              ...                              │ ◄────┤
// ├───────────────────────────────────────────────────────────────┤      │
// │                           head[367]                           │ ◄────┤
// ╞═══════════════════════════════════════════════════════════════╡      │
// │                             tail                              │ ◄────┘
// └───────────────────────────────────────────────────────────────┘   SIZE   ┘
// S: Small blocks map
type root struct {
	flMap uintptr
}

func (r *root) init() {
	r.flMap = 0
	r.setTail(nil)
	for fl := uintptr(0); fl < uintptr(_TLSFFLBits); fl++ {
		r.setSL(fl, 0)
		for sl := uint32(0); sl < _TLSFSLSize; sl++ {
			r.setHead(fl, sl, nil)
		}
	}
}

const (
	_TLSFSLStart  = _TLSFSizeofPointer
	_TLSFSLEnd    = _TLSFSLStart + (uintptr(_TLSFFLBits) << _TLSFAlignU32)
	_TLSFHLStart  = (_TLSFSLEnd + _TLSFALMask) &^ _TLSFALMask
	_TLSFHLEnd    = _TLSFHLStart + uintptr(_TLSFFLBits)*uintptr(_TLSFSLSize)*_TLSFSizeofPointer
	_TLSFRootSize = _TLSFHLEnd + _TLSFSizeofPointer
)

// Gets the second level map of the specified first level.
func (r *root) getSL(fl uintptr) uint32 {
	return *(*uint32)(unsafe.Pointer(uintptr(unsafe.Pointer(r)) + (fl << _TLSFAlignU32) + _TLSFSLStart))
}

// Sets the second level map of the specified first level.
func (r *root) setSL(fl uintptr, slMap uint32) {
	*(*uint32)(unsafe.Pointer(uintptr(unsafe.Pointer(r)) + (fl << _TLSFAlignU32) + _TLSFSLStart)) = slMap
}

// Gets the head of the free list for the specified combination of first and second level.
func (r *root) getHead(fl uintptr, sl uint32) *tlsfBlock {
	return *(**tlsfBlock)(unsafe.Pointer(uintptr(unsafe.Pointer(r)) + _TLSFHLStart +
		(((fl << _TLSFSLBits) + uintptr(sl)) << _TLSFAlignSizeLog2)))
}

// Sets the head of the free list for the specified combination of first and second level.
func (r *root) setHead(fl uintptr, sl uint32, head *tlsfBlock) {
	*(*uintptr)(unsafe.Pointer(uintptr(unsafe.Pointer(r)) + _TLSFHLStart +
		(((fl << _TLSFSLBits) + uintptr(sl)) << _TLSFAlignSizeLog2))) = uintptr(unsafe.Pointer(head))
}

// Gets the tail block.
func (r *root) getTail() *tlsfBlock {
	return *(**tlsfBlock)(unsafe.Pointer(uintptr(unsafe.Pointer(r)) + _TLSFHLEnd))
}

// Sets the tail block.
func (r *root) setTail(tail *tlsfBlock) {
	*(*uintptr)(unsafe.Pointer(uintptr(unsafe.Pointer(r)) + _TLSFHLEnd)) = uintptr(unsafe.Pointer(tail))
}

// Inserts a previously used block back into the free list.
func (a *Allocator) insertBlock(block *tlsfBlock) {
	var (
		r         = a.root
		blockInfo = block.mmInfo
		right     = block.getRight()
		rightInfo = right.mmInfo
	)
	//(blockInfo & FREE)

	// merge with right block if also free
	if rightInfo&_TLSFFREE != 0 {
		a.removeBlock(right)
		blockInfo = blockInfo + _TLSFBlockOverhead + (rightInfo & ^_TLSFTagsMask) // keep block tags
		block.mmInfo = blockInfo
		right = block.getRight()
		rightInfo = right.mmInfo
		// 'back' is Add below
	}

	// merge with left block if also free
	if blockInfo&_TLSFLEFTFREE != 0 {
		left := block.getFreeLeft()
		leftInfo := left.mmInfo
		if _TLSFDebug {
			assert(leftInfo&_TLSFFREE != 0, "must be free according to right tags")
		}
		a.removeBlock(left)
		block = left
		blockInfo = leftInfo + _TLSFBlockOverhead + (blockInfo & ^_TLSFTagsMask) // keep left tags
		block.mmInfo = blockInfo
		// 'back' is Add below
	}

	right.mmInfo = rightInfo | _TLSFLEFTFREE
	// reference to right is no longer used now, hence rightInfo is not synced

	// we now know the size of the block
	size := blockInfo & ^_TLSFTagsMask

	// Add 'back' to itself at the end of block
	*(*uintptr)(unsafe.Pointer(uintptr(unsafe.Pointer(right)) - _TLSFSizeofPointer)) = uintptr(unsafe.Pointer(block))

	// mapping_insert
	var (
		fl uintptr
		sl uint32
	)
	if size < uintptr(_TLSFSBSize) {
		fl = 0
		sl = uint32(size >> _TLSFALBits)
	} else {
		const inv = _TLSFSizeofPointer*8 - 1
		boundedSize := min(size, _TLSFBlockMaxSize)
		fl = inv - clz(boundedSize)
		sl = uint32((boundedSize >> (fl - uintptr(_TLSFSLBits))) ^ (1 << _TLSFSLBits))
		fl -= uintptr(_TLSFSBBits) - 1
	}

	// perform insertion
	head := r.getHead(fl, sl)
	block.prev = 0
	block.next = uintptr(unsafe.Pointer(head))
	if head != nil {
		head.prev = uintptr(unsafe.Pointer(block))
	}
	r.setHead(fl, sl, block)

	// update first and second level maps
	r.flMap |= 1 << fl
	r.setSL(fl, r.getSL(fl)|(1<<sl))
}

//goland:noinspection GoVetUnsafePointer
func (a *Allocator) removeBlock(block *tlsfBlock) {
	r := a.root
	blockInfo := block.mmInfo
	if _TLSFDebug {
		assert(blockInfo&_TLSFFREE != 0, "must be free")
	}
	size := blockInfo & ^_TLSFTagsMask
	if _TLSFDebug {
		assert(size >= _TLSFBlockMinSize, "must be valid")
	}

	// mapping_insert
	var (
		fl uintptr
		sl uint32
	)
	if size < uintptr(_TLSFSBSize) {
		fl = 0
		sl = uint32(size >> _TLSFALBits)
	} else {
		const inv = _TLSFSizeofPointer*8 - 1
		boundedSize := min(size, _TLSFBlockMaxSize)
		fl = inv - clz(boundedSize)
		sl = uint32((boundedSize >> (fl - uintptr(_TLSFSLBits))) ^ (1 << uintptr(_TLSFSLBits)))
		fl -= uintptr(_TLSFSBBits) - 1
	}
	if _TLSFDebug {
		assert(fl < uintptr(_TLSFFLBits) && sl < _TLSFSLSize, "fl/sl out of range")
	}

	// link previous and prev free block
	var (
		prev = block.prev
		next = block.next
	)
	if prev != 0 {
		(*tlsfBlock)(unsafe.Pointer(prev)).next = next
	}
	if next != 0 {
		(*tlsfBlock)(unsafe.Pointer(next)).prev = prev
	}

	// update head if we are removing it
	if block == r.getHead(fl, sl) {
		r.setHead(fl, sl, (*tlsfBlock)(unsafe.Pointer(next)))

		// clear second level map if head is empty now
		if next == 0 {
			slMap := r.getSL(fl)
			slMap &= ^(1 << sl)
			r.setSL(fl, slMap)

			// clear first level map if second level is empty now
			if slMap == 0 {
				r.flMap &= ^(1 << fl)
			}
		}
	}
	// note: does not alter left/back because it is likely that splitting
	// is performed afterwards, invalidating those changes. so, the caller
	// must perform those updates.
}

// Searches for a free block of at least the specified size.
func (a *Allocator) searchBlock(size uintptr) *tlsfBlock {
	// mapping_search
	var (
		fl uintptr
		sl uint32
		r  = a.root
	)
	if size < uintptr(_TLSFSBSize) {
		fl = 0
		sl = uint32(size >> _TLSFALBits)
	} else {
		const halfMaxSize = _TLSFBlockMaxSize >> 1 // don't round last fl
		const inv = _TLSFSizeofPointer*8 - 1
		const invRound = inv - uintptr(_TLSFSLBits)

		var requestSize uintptr
		if size < halfMaxSize {
			requestSize = size + (1 << (invRound - clz(size))) - 1
		} else {
			requestSize = size
		}

		fl = inv - clz(requestSize)
		sl = uint32((requestSize >> (fl - uintptr(_TLSFSLBits))) ^ (1 << _TLSFSLBits))
		fl -= uintptr(_TLSFSBBits) - 1
	}
	if _TLSFDebug {
		assert(fl < uintptr(_TLSFFLBits) && sl < _TLSFSLSize, "fl/sl out of range")
	}

	// search second level
	var (
		slMap = r.getSL(fl) & (^uint32(0) << sl)
		head  *tlsfBlock
	)
	if slMap == 0 {
		// search prev larger first level
		flMap := r.flMap & (^uintptr(0) << (fl + 1))
		if flMap == 0 {
			head = nil
		} else {
			fl = ctz(flMap)
			slMap = r.getSL(fl)
			if _TLSFDebug {
				assert(slMap != 0, "can't be zero if fl points here")
			}
			head = r.getHead(fl, ctz32(slMap))
		}
	} else {
		head = r.getHead(fl, ctz32(slMap))
	}

	return head
}

func (a *Allocator) prepareBlock(block *tlsfBlock, size uintptr) {
	blockInfo := block.mmInfo
	if _TLSFDebug {
		assert(((size+_TLSFBlockOverhead)&_TLSFALMask) == 0,
			"size must be aligned so the New block is")
	}
	// split if the block can hold another MINSIZE block incl. overhead
	remaining := (blockInfo & ^_TLSFTagsMask) - size
	if remaining >= _TLSFBlockOverhead+_TLSFBlockMinSize {
		block.mmInfo = size | (blockInfo & _TLSFLEFTFREE) // also discards FREE

		spare := (*tlsfBlock)(unsafe.Pointer(uintptr(unsafe.Pointer(block)) + _TLSFBlockOverhead + size))
		spare.mmInfo = (remaining - _TLSFBlockOverhead) | _TLSFFREE // not LEFTFREE
		a.insertBlock(spare)                                        // also sets 'back'

		// otherwise tag block as no longer FREE and right as no longer LEFTFREE
	} else {
		block.mmInfo = blockInfo & ^_TLSFFREE
		block.getRight().mmInfo &= ^_TLSFLEFTFREE
	}
}

// growMemory grows the pool by a number of 64kb pages to fit the required size
func (a *Allocator) growMemory(size uintptr) bool {
	if a.Grow == nil {
		return false
	}
	// Here, both rounding performed in searchBlock ...
	const halfMaxSize = _TLSFBlockMaxSize >> 1
	if size < halfMaxSize { // don't round last fl
		const invRound = (_TLSFSizeofPointer*8 - 1) - uintptr(_TLSFSLBits)
		size += (1 << (invRound - clz(size))) - 1
	}
	// and additional BLOCK_OVERHEAD must be taken into account. If we are going
	// to merge with the tail block, that's one time, otherwise it's two times.
	var (
		pagesBefore         = a.Pages
		offset      uintptr = 0
	)
	if _TLSFBlockOverhead != uintptr(unsafe.Pointer(a.root.getTail())) {
		offset = 1
	}
	size += _TLSFBlockOverhead << ((uintptr(pagesBefore) << 16) - offset)
	pagesNeeded := ((int32(size) + 0xffff) & ^0xffff) >> 16

	addedPages, start, end := a.Grow(pagesBefore, pagesNeeded, size)
	if start == 0 || end == 0 {
		return false
	}
	if addedPages == 0 {
		addedPages = int32((end - start) / _TLSFPageSize)
		if (end-start)%_TLSFPageSize > 0 {
			addedPages++
		}
	}
	a.Pages += addedPages
	a.HeapEnd = end
	a.addMemory(start, end)
	return true
}

// addMemory adds the newly allocated memory to the Allocator bitmaps
//goland:noinspection GoVetUnsafePointer
func (a *Allocator) addMemory(start, end uintptr) bool {
	if _TLSFDebug {
		assert(start <= end, "start must be <= end")
	}
	start = ((start + _TLSFBlockOverhead + _TLSFALMask) & ^_TLSFALMask) - _TLSFBlockOverhead
	end &= ^_TLSFALMask

	var tail = a.root.getTail()
	var tailInfo uintptr = 0
	if tail != nil { // more memory
		if _TLSFDebug {
			assert(start >= uintptr(unsafe.Pointer(tail))+_TLSFBlockOverhead, "out of bounds")
		}

		// merge with current tail if adjacent
		const offsetToTail = _TLSFALSize
		if start-offsetToTail == uintptr(unsafe.Pointer(tail)) {
			start -= offsetToTail
			tailInfo = tail.mmInfo
		} else {
			// We don't do this, but a user might `memory.Grow` manually
			// leading to non-adjacent pages managed by Allocator.
		}

	} else if _TLSFDebug { // first memory
		assert(start >= uintptr(unsafe.Pointer(a.root))+_TLSFRootSize, "starts after root")
	}

	// check if size is large enough for a free block and the tail block
	var size = end - start
	if size < _TLSFBlockOverhead+_TLSFBlockMinSize+_TLSFBlockOverhead {
		return false
	}

	// left size is total minus its own and the zero-length tail's header
	var leftSize = size - 2*_TLSFBlockOverhead
	var left = (*tlsfBlock)(unsafe.Pointer(start))
	left.mmInfo = leftSize | _TLSFFREE | (tailInfo & _TLSFLEFTFREE)
	left.prev = 0
	left.next = 0

	// tail is a zero-length used block
	tail = (*tlsfBlock)(unsafe.Pointer(start + _TLSFBlockOverhead + leftSize))
	tail.mmInfo = 0 | _TLSFLEFTFREE
	a.root.setTail(tail)

	a.FreeSize += int64(leftSize)
	a.HeapSize += int64(end - start)

	// also merges with free left before tail / sets 'back'
	a.insertBlock(left)

	return true
}

// Computes the size (excl. header) of a block.
func computeSize(size uintptr) uintptr {
	// Size must be large enough and aligned minus preceeding overhead
	if size <= _TLSFBlockMinSize {
		return _TLSFBlockMinSize
	} else {
		return ((size + _TLSFBlockOverhead + _TLSFALMask) & ^_TLSFALMask) - _TLSFBlockOverhead
	}
}

// Prepares and checks an allocation size.
func prepareSize(size uintptr) uintptr {
	if size > _TLSFBlockMaxSize {
		panic("allocation too large")
	}
	return computeSize(size)
}

// Allocates a block of the specified size.
func (a *Allocator) allocateBlock(size uintptr) *tlsfBlock {
	var payloadSize = prepareSize(size)
	var block = a.searchBlock(payloadSize)
	if block == nil {
		if !a.growMemory(payloadSize) {
			return nil
		}
		block = a.searchBlock(payloadSize)
		if _TLSFDebug {
			assert(block != nil, "block must be found now")
		}
		if block == nil {
			return nil
		}
	}
	if _TLSFDebug {
		assert((block.mmInfo & ^_TLSFTagsMask) >= payloadSize, "must fit")
	}

	a.removeBlock(block)
	a.prepareBlock(block, payloadSize)

	// update stats
	payloadSize = block.mmInfo & ^_TLSFTagsMask
	a.AllocSize += int64(payloadSize)
	if a.AllocSize > a.MaxUsedSize {
		a.MaxUsedSize = a.AllocSize
	}
	a.FreeSize -= int64(payloadSize)
	a.Allocs++

	// return block
	return block
}

func (a *Allocator) reallocateBlock(block *tlsfBlock, size uintptr) *tlsfBlock {
	var payloadSize = prepareSize(size)
	var blockInfo = block.mmInfo
	var blockSize = blockInfo & ^_TLSFTagsMask

	// possibly split and update runtime size if it still fits
	if payloadSize <= blockSize {
		a.prepareBlock(block, payloadSize)
		//if (isDefined(ASC_RTRACE)) {
		//	if (payloadSize != blockSize) onresize(block, BLOCK_OVERHEAD + blockSize);
		//}
		return block
	}

	// merge with right free block if merger is large enough
	var (
		right     = block.getRight()
		rightInfo = right.mmInfo
	)
	if rightInfo&_TLSFFREE != 0 {
		mergeSize := blockSize + _TLSFBlockOverhead + (rightInfo & ^_TLSFTagsMask)
		if mergeSize >= payloadSize {
			a.removeBlock(right)
			block.mmInfo = (blockInfo & _TLSFTagsMask) | mergeSize
			a.prepareBlock(block, payloadSize)
			//if (isDefined(ASC_RTRACE)) onresize(block, BLOCK_OVERHEAD + blockSize);
			return block
		}
	}

	// otherwise, move the block
	return a.moveBlock(block, size)
}

func (a *Allocator) moveBlock(block *tlsfBlock, newSize uintptr) *tlsfBlock {
	newBlock := a.allocateBlock(newSize)
	if newBlock == nil {
		return nil
	}

	memcpy(unsafe.Pointer(uintptr(unsafe.Pointer(newBlock))+_TLSFBlockOverhead),
		unsafe.Pointer(uintptr(unsafe.Pointer(block))+_TLSFBlockOverhead),
		block.mmInfo & ^_TLSFTagsMask)

	a.freeBlock(block)
	//maybeFreeBlock(a, block)

	return newBlock
}

func (a *Allocator) freeBlock(block *tlsfBlock) {
	size := block.mmInfo & ^_TLSFTagsMask
	a.FreeSize += int64(size)
	a.AllocSize -= int64(size)
	a.Allocs--

	block.mmInfo = block.mmInfo | _TLSFFREE
	a.insertBlock(block)
}

func min(l, r uintptr) uintptr {
	if l < r {
		return l
	}
	return r
}

func clz(value uintptr) uintptr {
	return uintptr(bits.LeadingZeros(uint(value)))
}

func ctz(value uintptr) uintptr {
	return uintptr(bits.TrailingZeros(uint(value)))
}

func ctz32(value uint32) uint32 {
	return uint32(bits.TrailingZeros32(value))
}

func checkUsedBlock(ptr uintptr) *tlsfBlock {
	block := (*tlsfBlock)(unsafe.Pointer(ptr - _TLSFBlockOverhead))
	if !(ptr != 0 && ((ptr & _TLSFALMask) == 0) && ((block.mmInfo & _TLSFFREE) == 0)) {
		panic("used block is not valid to be freed or reallocated")
	}
	return block
}

func allocationSize0(ptr uintptr) uintptr {
	return ((*tlsfBlock)(unsafe.Pointer(ptr - _TLSFBlockOverhead))).mmInfo & ^_TLSFTagsMask
}

func allocationSize(ptr unsafe.Pointer) uintptr {
	return ((*tlsfBlock)(unsafe.Pointer(uintptr(ptr) - _TLSFBlockOverhead))).mmInfo & ^_TLSFTagsMask
}

func AllocatorPrintDebugInfo() {
	println("ALIGNOF_U32		", int64(_TLSFAlignU32))
	println("ALIGN_SIZE_LOG2	", int64(_TLSFAlignSizeLog2))
	println("U32_MAX			", ^uint32(0))
	println("PTR_MAX			", ^uintptr(0))
	println("AL_BITS			", int64(_TLSFALBits))
	println("AL_SIZE			", int64(_TLSFALSize))
	println("AL_MASK			", int64(_TLSFALMask))
	println("BLOCK_OVERHEAD	", int64(_TLSFBlockOverhead))
	println("BLOCK_MAXSIZE	", int64(_TLSFBlockMaxSize))
	println("SL_BITS			", int64(_TLSFSLBits))
	println("SL_SIZE			", int64(_TLSFSLSize))
	println("SB_BITS			", int64(_TLSFSBBits))
	println("SB_SIZE			", int64(_TLSFSBSize))
	println("FL_BITS			", int64(_TLSFFLBits))
	println("FREE			", int64(_TLSFFREE))
	println("LEFTFREE		", int64(_TLSFLEFTFREE))
	println("TAGS_MASK		", int64(_TLSFTagsMask))
	println("BLOCK_MINSIZE	", int64(_TLSFBlockMinSize))
	println("SL_START		", int64(_TLSFSLStart))
	println("SL_END			", int64(_TLSFSLEnd))
	println("HL_START		", int64(_TLSFHLStart))
	println("HL_END			", int64(_TLSFHLEnd))
	println("ROOT_SIZE		", int64(_TLSFRootSize))
}
