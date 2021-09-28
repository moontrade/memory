package tlsf

import (
	"math"
	"math/bits"
	"unsafe"
)

// Heap === Heap (Two-Level Segregate Fit) memory allocator ===
//
// Heap is a general purpose dynamic memory allocator specifically designed to meet
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
// see: https://github.com/AssemblyScript/assemblyscript
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
type Heap struct {
	root      *root
	HeapStart uintptr
	HeapEnd   uintptr
	arena     uintptr
	Grow      Grow
	Slot      uint8
	Stats
}

// Stats provides the metrics of an Allocator
type Stats struct {
	HeapSize        int64
	AllocSize       int64
	PeakAllocSize   int64
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
	if s.HeapSize == 0 || s.PeakAllocSize == 0 {
		return 0
	}
	pct := float64(s.HeapSize-s.PeakAllocSize) / float64(s.HeapSize)
	s.fragmentation = float32(math.Floor(pct*100) / 100)
	return s.fragmentation
}

// Grow provides the ability to Grow the heap and allocate a contiguous
// chunk of system memory to add to the allocator.
type Grow func(pagesBefore, pagesNeeded int32, minSize uintptr) (pagesAdded int32, start, end uintptr)

const (
	PageSize = uintptr(64 * 1024)

	_TLSFAlignU32 = 2
	// All allocation sizes and addresses are aligned to 4 or 8 bytes.
	// 32bit = 2
	// 64bit = 3
	// <expr> = bits.UintSize / 8 / 4 + 1
	_TLSFAlignSizeLog2 uintptr = ((32 << (^uint(0) >> 63)) / 8 / 4) + 1
	_TLSFSizeofPointer         = uintptr(unsafe.Sizeof(uintptr(0)))

	ALBits uint32  = 4 // 16 bytes to fit up to v128
	ALSize uintptr = 1 << uintptr(ALBits)
	ALMask         = ALSize - 1

	// Overhead of a memory manager block.
	BlockOverhead = unsafe.Sizeof(BLOCK{})
	// Block constants. A block must have a minimum size of three pointers so it can hold `prev`,
	// `prev` and `back` if free.
	BlockMinSize = ((3*_TLSFSizeofPointer + BlockOverhead + ALMask) & ^ALMask) - BlockOverhead
	// Maximum size of a memory manager block's payload.
	BlockMaxSize = (1 << 30) - BlockOverhead
	//BlockMaxSize = (1 << ((_TLSFAlignSizeLog2 + 1)*10)) - BlockOverhead

	_TLSFDebug = false

	_TLSFSLBits uint32 = 4
	_TLSFSLSize uint32 = 1 << _TLSFSLBits
	_TLSFSBBits        = _TLSFSLBits + ALBits
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
	TagsMask              = _TLSFFREE | _TLSFLEFTFREE
)

// Alloc allocates a block of memory that fits the size provided
//goland:noinspection GoVetUnsafePointer
func (a *Heap) Alloc(size uintptr) uintptr {
	if a == nil {
		panic("nil")
	}
	p := uintptr(unsafe.Pointer(a.allocateBlock(uintptr(size))))
	if p == 0 {
		return 0
	}
	p = p + BlockOverhead
	return p
}

// AllocZeroed allocates a block of memory that fits the size provided
//goland:noinspection GoVetUnsafePointer
func (a *Heap) AllocZeroed(size uintptr) uintptr {
	p := uintptr(unsafe.Pointer(a.allocateBlock(uintptr(size))))
	if p == 0 {
		return 0
	}
	p = p + BlockOverhead
	Zero(unsafe.Pointer(p), size)
	return p
}

// Realloc determines the best way to resize an allocation.
func (a *Heap) Realloc(ptr uintptr, size uintptr) uintptr {
	p := uintptr(unsafe.Pointer(a.moveBlock(checkUsedBlock(ptr), uintptr(size))))
	if p == 0 {
		return 0
	}
	return p + BlockOverhead
}

// Free release the allocation back into the free list.
//goland:noinspection GoVetUnsafePointer
func (a *Heap) Free(ptr uintptr) {
	//println("Free", uint(ptr))
	//a.freeBlock((*tlsfBlock)(unsafe.Pointer(ptr - BlockOverhead)))
	a.freeBlock(checkUsedBlock(ptr))
}

// Bootstrap bootstraps the Allocator with the initial block of contiguous memory
// that at least fits the minimum required to fit the bitmap.
//goland:noinspection GoVetUnsafePointer
func Bootstrap(start, end uintptr, pages int32, grow Grow) *Heap {
	start = (start + unsafe.Alignof(unsafe.Pointer(nil)) - 1) &^ (unsafe.Alignof(unsafe.Pointer(nil)) - 1)

	//if a.T {
	//	println("Bootstrap", "pages", pages, uint(start), uint(end), uint(end-start))
	//}
	// init allocator
	a := (*Heap)(unsafe.Pointer(start))
	*a = Heap{
		HeapStart: start,
		HeapEnd:   end,
		Stats: Stats{
			InitialPages: pages,
			Pages:        pages,
		},
		Grow: grow,
	}

	// init root
	rootOffset := uintptr(unsafe.Sizeof(Heap{})) + ((start + ALMask) & ^ALMask)
	a.root = (*root)(unsafe.Pointer(rootOffset))
	a.root.init()

	// add initial memory
	a.addMemory(rootOffset+RootSize, end)
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
type BLOCK struct {
	MMInfo uintptr
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
	BLOCK
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
	return (*tlsfBlock)(unsafe.Pointer(uintptr(unsafe.Pointer(block)) + BlockOverhead + (block.MMInfo & ^TagsMask)))
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
	SLStart  = _TLSFSizeofPointer
	SLEnd    = SLStart + (uintptr(_TLSFFLBits) << _TLSFAlignU32)
	HLStart  = (SLEnd + ALMask) &^ ALMask
	HLEnd    = HLStart + uintptr(_TLSFFLBits)*uintptr(_TLSFSLSize)*_TLSFSizeofPointer
	RootSize = HLEnd + _TLSFSizeofPointer
)

// Gets the second level map of the specified first level.
func (r *root) getSL(fl uintptr) uint32 {
	return *(*uint32)(unsafe.Pointer(uintptr(unsafe.Pointer(r)) + (fl << _TLSFAlignU32) + SLStart))
}

// Sets the second level map of the specified first level.
func (r *root) setSL(fl uintptr, slMap uint32) {
	*(*uint32)(unsafe.Pointer(uintptr(unsafe.Pointer(r)) + (fl << _TLSFAlignU32) + SLStart)) = slMap
}

// Gets the head of the free list for the specified combination of first and second level.
func (r *root) getHead(fl uintptr, sl uint32) *tlsfBlock {
	return *(**tlsfBlock)(unsafe.Pointer(uintptr(unsafe.Pointer(r)) + HLStart +
		(((fl << _TLSFSLBits) + uintptr(sl)) << _TLSFAlignSizeLog2)))
}

// Sets the head of the free list for the specified combination of first and second level.
func (r *root) setHead(fl uintptr, sl uint32, head *tlsfBlock) {
	*(*uintptr)(unsafe.Pointer(uintptr(unsafe.Pointer(r)) + HLStart +
		(((fl << _TLSFSLBits) + uintptr(sl)) << _TLSFAlignSizeLog2))) = uintptr(unsafe.Pointer(head))
}

// Gets the tail block.
func (r *root) getTail() *tlsfBlock {
	return *(**tlsfBlock)(unsafe.Pointer(uintptr(unsafe.Pointer(r)) + HLEnd))
}

// Sets the tail block.
func (r *root) setTail(tail *tlsfBlock) {
	*(*uintptr)(unsafe.Pointer(uintptr(unsafe.Pointer(r)) + HLEnd)) = uintptr(unsafe.Pointer(tail))
}

// Inserts a previously used block back into the free list.
func (a *Heap) insertBlock(block *tlsfBlock) {
	var (
		r         = a.root
		blockInfo = block.MMInfo
		right     = block.getRight()
		rightInfo = right.MMInfo
	)
	//(blockInfo & FREE)

	// merge with right block if also free
	if rightInfo&_TLSFFREE != 0 {
		a.removeBlock(right)
		blockInfo = blockInfo + BlockOverhead + (rightInfo & ^TagsMask) // keep block tags
		block.MMInfo = blockInfo
		right = block.getRight()
		rightInfo = right.MMInfo
		// 'back' is Add below
	}

	// merge with left block if also free
	if blockInfo&_TLSFLEFTFREE != 0 {
		left := block.getFreeLeft()
		leftInfo := left.MMInfo
		if _TLSFDebug {
			assert(leftInfo&_TLSFFREE != 0, "must be free according to right tags")
		}
		a.removeBlock(left)
		block = left
		blockInfo = leftInfo + BlockOverhead + (blockInfo & ^TagsMask) // keep left tags
		block.MMInfo = blockInfo
		// 'back' is Add below
	}

	right.MMInfo = rightInfo | _TLSFLEFTFREE
	// reference to right is no longer used now, hence rightInfo is not synced

	// we now know the size of the block
	size := blockInfo & ^TagsMask

	// Add 'back' to itself at the end of block
	*(*uintptr)(unsafe.Pointer(uintptr(unsafe.Pointer(right)) - _TLSFSizeofPointer)) = uintptr(unsafe.Pointer(block))

	// mapping_insert
	var (
		fl uintptr
		sl uint32
	)
	if size < uintptr(_TLSFSBSize) {
		fl = 0
		sl = uint32(size >> ALBits)
	} else {
		const inv = _TLSFSizeofPointer*8 - 1
		boundedSize := min(size, BlockMaxSize)
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
func (a *Heap) removeBlock(block *tlsfBlock) {
	r := a.root
	blockInfo := block.MMInfo
	if _TLSFDebug {
		assert(blockInfo&_TLSFFREE != 0, "must be free")
	}
	size := blockInfo & ^TagsMask
	if _TLSFDebug {
		assert(size >= BlockMinSize, "must be valid")
	}

	// mapping_insert
	var (
		fl uintptr
		sl uint32
	)
	if size < uintptr(_TLSFSBSize) {
		fl = 0
		sl = uint32(size >> ALBits)
	} else {
		const inv = _TLSFSizeofPointer*8 - 1
		boundedSize := min(size, BlockMaxSize)
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
func (a *Heap) searchBlock(size uintptr) *tlsfBlock {
	// mapping_search
	var (
		fl uintptr
		sl uint32
		r  = a.root
	)
	if size < uintptr(_TLSFSBSize) {
		fl = 0
		sl = uint32(size >> ALBits)
	} else {
		const (
			halfMaxSize = BlockMaxSize >> 1 // don't round last fl
			inv         = _TLSFSizeofPointer*8 - 1
			invRound    = inv - uintptr(_TLSFSLBits)
		)

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

func (a *Heap) prepareBlock(block *tlsfBlock, size uintptr) {
	blockInfo := block.MMInfo
	if _TLSFDebug {
		assert(((size+BlockOverhead)&ALMask) == 0,
			"size must be aligned so the New block is")
	}
	// split if the block can hold another MINSIZE block incl. overhead
	remaining := (blockInfo & ^TagsMask) - size
	if remaining >= BlockOverhead+BlockMinSize {
		block.MMInfo = size | (blockInfo & _TLSFLEFTFREE) // also discards FREE

		spare := (*tlsfBlock)(unsafe.Pointer(uintptr(unsafe.Pointer(block)) + BlockOverhead + size))
		spare.MMInfo = (remaining - BlockOverhead) | _TLSFFREE // not LEFTFREE
		a.insertBlock(spare)                                   // also sets 'back'

		// otherwise tag block as no longer FREE and right as no longer LEFTFREE
	} else {
		block.MMInfo = blockInfo & ^_TLSFFREE
		block.getRight().MMInfo &= ^_TLSFLEFTFREE
	}
}

// growMemory grows the pool by a number of 64kb pages to fit the required size
func (a *Heap) growMemory(size uintptr) bool {
	if a.Grow == nil {
		return false
	}
	// Here, both rounding performed in searchBlock ...
	const halfMaxSize = BlockMaxSize >> 1
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
	if BlockOverhead != uintptr(unsafe.Pointer(a.root.getTail())) {
		offset = 1
	}
	size += BlockOverhead << ((uintptr(pagesBefore) << 16) - offset)
	pagesNeeded := ((int32(size) + 0xffff) & ^0xffff) >> 16

	addedPages, start, end := a.Grow(pagesBefore, pagesNeeded, size)
	if start == 0 || end == 0 {
		return false
	}
	if addedPages == 0 {
		addedPages = int32((end - start) / PageSize)
		if (end-start)%PageSize > 0 {
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
func (a *Heap) addMemory(start, end uintptr) bool {
	if _TLSFDebug {
		assert(start <= end, "start must be <= end")
	}
	start = ((start + BlockOverhead + ALMask) & ^ALMask) - BlockOverhead
	end &= ^ALMask

	var tail = a.root.getTail()
	var tailInfo uintptr = 0
	if tail != nil { // more memory
		if _TLSFDebug {
			assert(start >= uintptr(unsafe.Pointer(tail))+BlockOverhead, "out of bounds")
		}

		// merge with current tail if adjacent
		const offsetToTail = ALSize
		if start-offsetToTail == uintptr(unsafe.Pointer(tail)) {
			start -= offsetToTail
			tailInfo = tail.MMInfo
		} else {
			// We don't do this, but a user might `memory.Grow` manually
			// leading to non-adjacent pages managed by Allocator.
		}
	} else if _TLSFDebug { // first memory
		assert(start >= uintptr(unsafe.Pointer(a.root))+RootSize, "starts after root")
	}

	// check if size is large enough for a free block and the tail block
	var size = end - start
	if size < BlockOverhead+BlockMinSize+BlockOverhead {
		return false
	}

	// left size is total minus its own and the zero-length tail's header
	var (
		leftSize = size - 2*BlockOverhead
		left     = (*tlsfBlock)(unsafe.Pointer(start))
	)
	left.MMInfo = leftSize | _TLSFFREE | (tailInfo & _TLSFLEFTFREE)
	left.prev = 0
	left.next = 0

	// tail is a zero-length used block
	tail = (*tlsfBlock)(unsafe.Pointer(start + BlockOverhead + leftSize))
	tail.MMInfo = 0 | _TLSFLEFTFREE
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
	if size <= BlockMinSize {
		return BlockMinSize
	} else {
		return ((size + BlockOverhead + ALMask) & ^ALMask) - BlockOverhead
	}
}

// Prepares and checks an allocation size.
func prepareSize(size uintptr) uintptr {
	if size > BlockMaxSize {
		panic("allocation too large")
	}
	return computeSize(size)
}

// Allocates a block of the specified size.
func (a *Heap) allocateBlock(size uintptr) *tlsfBlock {
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
		assert((block.MMInfo & ^TagsMask) >= payloadSize, "must fit")
	}

	a.removeBlock(block)
	a.prepareBlock(block, payloadSize)

	// update stats
	payloadSize = block.MMInfo & ^TagsMask
	a.AllocSize += int64(payloadSize)
	if a.AllocSize > a.PeakAllocSize {
		a.PeakAllocSize = a.AllocSize
	}
	a.FreeSize -= int64(payloadSize)
	a.Allocs++

	// return block
	return block
}

func (a *Heap) reallocateBlock(block *tlsfBlock, size uintptr) *tlsfBlock {
	var payloadSize = prepareSize(size)
	var blockInfo = block.MMInfo
	var blockSize = blockInfo & ^TagsMask

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
		rightInfo = right.MMInfo
	)
	if rightInfo&_TLSFFREE != 0 {
		mergeSize := blockSize + BlockOverhead + (rightInfo & ^TagsMask)
		if mergeSize >= payloadSize {
			a.removeBlock(right)
			block.MMInfo = (blockInfo & TagsMask) | mergeSize
			a.prepareBlock(block, payloadSize)
			//if (isDefined(ASC_RTRACE)) onresize(block, BLOCK_OVERHEAD + blockSize);
			return block
		}
	}

	// otherwise, move the block
	return a.moveBlock(block, size)
}

func (a *Heap) moveBlock(block *tlsfBlock, newSize uintptr) *tlsfBlock {
	newBlock := a.allocateBlock(newSize)
	if newBlock == nil {
		return nil
	}

	Copy(unsafe.Pointer(uintptr(unsafe.Pointer(newBlock))+BlockOverhead),
		unsafe.Pointer(uintptr(unsafe.Pointer(block))+BlockOverhead),
		uintptr(block.MMInfo & ^TagsMask))

	a.freeBlock(block)
	//maybeFreeBlock(a, block)

	return newBlock
}

func (a *Heap) freeBlock(block *tlsfBlock) {
	size := block.MMInfo & ^TagsMask
	a.FreeSize += int64(size)
	a.AllocSize -= int64(size)
	a.Allocs--

	block.MMInfo = block.MMInfo | _TLSFFREE
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

//goland:noinspection GoVetUnsafePointer
func checkUsedBlock(ptr uintptr) *tlsfBlock {
	block := (*tlsfBlock)(unsafe.Pointer(ptr - BlockOverhead))
	if !(ptr != 0 && ((ptr & ALMask) == 0) && ((block.MMInfo & _TLSFFREE) == 0)) {
		panic("used block is not valid to be freed or reallocated")
	}
	return block
}

//goland:noinspection GoVetUnsafePointer
func SizeOf(ptr uintptr) uintptr {
	return ((*tlsfBlock)(unsafe.Pointer(ptr - BlockOverhead))).MMInfo & ^TagsMask
}

func PrintDebugInfo() {
	println("ALIGNOF_U32		", int64(_TLSFAlignU32))
	println("ALIGN_SIZE_LOG2	", int64(_TLSFAlignSizeLog2))
	println("U32_MAX			", ^uint32(0))
	println("PTR_MAX			", ^uintptr(0))
	println("AL_BITS			", int64(ALBits))
	println("AL_SIZE			", int64(ALSize))
	println("AL_MASK			", int64(ALMask))
	println("BLOCK_OVERHEAD	", int64(BlockOverhead))
	println("BLOCK_MAXSIZE	", int64(BlockMaxSize))
	println("SL_BITS			", int64(_TLSFSLBits))
	println("SL_SIZE			", int64(_TLSFSLSize))
	println("SB_BITS			", int64(_TLSFSBBits))
	println("SB_SIZE			", int64(_TLSFSBSize))
	println("FL_BITS			", int64(_TLSFFLBits))
	println("FREE			", int64(_TLSFFREE))
	println("LEFTFREE		", int64(_TLSFLEFTFREE))
	println("TAGS_MASK		", int64(TagsMask))
	println("BLOCK_MINSIZE	", int64(BlockMinSize))
	println("SL_START		", int64(SLStart))
	println("SL_END			", int64(SLEnd))
	println("HL_START		", int64(HLStart))
	println("HL_END			", int64(HLEnd))
	println("ROOT_SIZE		", int64(RootSize))
}
