package mem

import (
	"math/bits"
	"unsafe"
)

// Allocator === TLSF (Two-Level Segregate Fit) memory allocator ===
//
// TLSF is a general purpose dynamic memory allocator specifically designed to
// meet real-time requirements:
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
	Grow      Grow
	Stats
}

// Stats provides the metrics of an Allocator
type Stats struct {
	HeapSize     int64
	AllocSize    int64
	MaxUsedSize  int64
	FreeSize     int64
	Allocs       int32
	InitialPages int32
	Pages        int32
	Grows        int32
}

// Grow provides the ability to Grow the heap and allocate a contiguous
// chunk of system memory to add to the allocator.
type Grow func(pagesBefore, pagesNeeded int32, minSize uintptr) (pagesAdded int32, start, end uintptr)

// Malloc provides the actual system allocation
type Malloc func(size uintptr) unsafe.Pointer

const (
	tlsf_PAGE_SIZE = uintptr(64 * 1024)

	tlsf_ALIGN_U32 = 2
	// All allocation sizes and addresses are aligned to 4 or 8 bytes.
	// 32bit = 2
	// 64bit = 3
	// <expr> = bits.UintSize / 8 / 4 + 1
	tlsf_ALIGN_SIZE_LOG2 = ((32 << (^uint(0) >> 63)) / 8 / 4) + 1
	tlsf_sizeofPointer   = unsafe.Sizeof(uintptr(0))

	tlsf_AL_BITS uint32  = 4 // 16 bytes to fit up to v128
	tlsf_AL_SIZE uintptr = 1 << uintptr(tlsf_AL_BITS)
	tlsf_AL_MASK         = tlsf_AL_SIZE - 1

	// Overhead of a memory manager block.
	tlsf_BLOCK_OVERHEAD = unsafe.Sizeof(tlsfBLOCK{})
	// Block constants. A block must have a minimum size of three pointers so it can hold `prev`,
	// `prev` and `back` if free.
	tlsf_BLOCK_MINSIZE = ((3*tlsf_sizeofPointer + tlsf_BLOCK_OVERHEAD + tlsf_AL_MASK) & ^tlsf_AL_MASK) - tlsf_BLOCK_OVERHEAD
	// Maximum size of a memory manager block's payload.
	tlsf_BLOCK_MAXSIZE = (1 << 30) - tlsf_BLOCK_OVERHEAD
	//tlsf_BLOCK_MAXSIZE = (1 << ((tlsf_ALIGN_SIZE_LOG2 + 1)*10)) - tlsf_BLOCK_OVERHEAD

	tlsf_DEBUG = false

	tlsf_SL_BITS uint32 = 4
	tlsf_SL_SIZE uint32 = 1 << tlsf_SL_BITS
	tlsf_SB_BITS        = tlsf_SL_BITS + tlsf_AL_BITS
	tlsf_SB_SIZE uint32 = 1 << tlsf_SB_BITS
	tlsf_FL_BITS        = 31 - tlsf_SB_BITS

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
	// VMs limit to 2GB total (currently), making one 1G block max (or three 512M etc.) due to block overhead

	// Tags stored in otherwise unused alignment bits
	tlsf_FREE      uintptr = 1 << 0
	tlsf_LEFTFREE  uintptr = 1 << 1
	tlsf_TAGS_MASK         = tlsf_FREE | tlsf_LEFTFREE
)

// Alloc allocates a block of memory that fits the size provided
//goland:noinspection GoVetUnsafePointer
func (a *Allocator) Alloc(size uintptr) unsafe.Pointer {
	p := uintptr(unsafe.Pointer(a.allocateBlock(size)))
	if p == 0 {
		return nil
	}
	return unsafe.Pointer(p + tlsf_BLOCK_OVERHEAD)
}

// Realloc determines the best way to resize an allocation.
func (a *Allocator) Realloc(ptr unsafe.Pointer, size uintptr) unsafe.Pointer {
	p := uintptr(unsafe.Pointer(a.moveBlock(checkUsedBlock(uintptr(ptr)), size)))
	if p == 0 {
		return nil
	}
	return unsafe.Pointer(p + tlsf_BLOCK_OVERHEAD)
}

// Free release the allocation back into the free list.
func (a *Allocator) Free(ptr unsafe.Pointer) {
	a.freeBlock(checkUsedBlock(uintptr(ptr)))
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

// Bootstrap bootstraps the Allocator with the initial block of contiguous memory
// that at least fits the minimum required to fit the bitmap.
//goland:noinspection GoVetUnsafePointer
func Bootstrap(start, end uintptr, pages int32, grow Grow) *Allocator {
	start = (start + unsafe.Alignof(unsafe.Pointer(nil)) - 1) &^ (unsafe.Alignof(unsafe.Pointer(nil)) - 1)

	//if a.T {
	//	println("Bootstrap", "pages", pages, uint(start), uint(end), uint(end-start))
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
	rootOffset := unsafe.Sizeof(Allocator{}) + ((start + tlsf_AL_MASK) & ^tlsf_AL_MASK)
	a.root = (*root)(unsafe.Pointer(rootOffset))
	a.root.init()

	// add initial memory
	a.addMemory(rootOffset+tlsf_ROOT_SIZE, end)
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
	return *(**tlsfBlock)(unsafe.Pointer(uintptr(unsafe.Pointer(block)) - tlsf_sizeofPointer))
}

// Gets the right block of a block by advancing to the right by its size.
func (block *tlsfBlock) getRight() *tlsfBlock {
	return (*tlsfBlock)(unsafe.Pointer(uintptr(unsafe.Pointer(block)) + tlsf_BLOCK_OVERHEAD + (block.mmInfo & ^tlsf_TAGS_MASK)))
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
	for fl := uintptr(0); fl < uintptr(tlsf_FL_BITS); fl++ {
		r.setSL(fl, 0)
		for sl := uint32(0); sl < tlsf_SL_SIZE; sl++ {
			r.setHead(fl, sl, nil)
		}
	}
}

const (
	tlsf_SL_START  = tlsf_sizeofPointer
	tlsf_SL_END    = tlsf_SL_START + (uintptr(tlsf_FL_BITS) << tlsf_ALIGN_U32)
	tlsf_HL_START  = (tlsf_SL_END + tlsf_AL_MASK) &^ tlsf_AL_MASK
	tlsf_HL_END    = tlsf_HL_START + uintptr(tlsf_FL_BITS)*uintptr(tlsf_SL_SIZE)*tlsf_sizeofPointer
	tlsf_ROOT_SIZE = tlsf_HL_END + tlsf_sizeofPointer
)

// Gets the second level map of the specified first level.
func (r *root) getSL(fl uintptr) uint32 {
	return *(*uint32)(unsafe.Pointer(uintptr(unsafe.Pointer(r)) + (fl << tlsf_ALIGN_U32) + tlsf_SL_START))
}

// Sets the second level map of the specified first level.
func (r *root) setSL(fl uintptr, slMap uint32) {
	*(*uint32)(unsafe.Pointer(uintptr(unsafe.Pointer(r)) + (fl << tlsf_ALIGN_U32) + tlsf_SL_START)) = slMap
}

// Gets the head of the free list for the specified combination of first and second level.
func (r *root) getHead(fl uintptr, sl uint32) *tlsfBlock {
	return *(**tlsfBlock)(unsafe.Pointer(uintptr(unsafe.Pointer(r)) + tlsf_HL_START +
		(((fl << tlsf_SL_BITS) + uintptr(sl)) << tlsf_ALIGN_SIZE_LOG2)))
}

// Sets the head of the free list for the specified combination of first and second level.
func (r *root) setHead(fl uintptr, sl uint32, head *tlsfBlock) {
	*(*uintptr)(unsafe.Pointer(uintptr(unsafe.Pointer(r)) + tlsf_HL_START +
		(((fl << tlsf_SL_BITS) + uintptr(sl)) << tlsf_ALIGN_SIZE_LOG2))) = uintptr(unsafe.Pointer(head))
}

// Gets the tail block.
func (r *root) getTail() *tlsfBlock {
	return *(**tlsfBlock)(unsafe.Pointer(uintptr(unsafe.Pointer(r)) + tlsf_HL_END))
}

// Sets the tail block.
func (r *root) setTail(tail *tlsfBlock) {
	*(*uintptr)(unsafe.Pointer(uintptr(unsafe.Pointer(r)) + tlsf_HL_END)) = uintptr(unsafe.Pointer(tail))
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
	if rightInfo&tlsf_FREE != 0 {
		a.removeBlock(right)
		blockInfo = blockInfo + tlsf_BLOCK_OVERHEAD + (rightInfo & ^tlsf_TAGS_MASK) // keep block tags
		block.mmInfo = blockInfo
		right = block.getRight()
		rightInfo = right.mmInfo
		// 'back' is Add below
	}

	// merge with left block if also free
	if blockInfo&tlsf_LEFTFREE != 0 {
		left := block.getFreeLeft()
		leftInfo := left.mmInfo
		if tlsf_DEBUG {
			assert(leftInfo&tlsf_FREE != 0, "must be free according to right tags")
		}
		a.removeBlock(left)
		block = left
		blockInfo = leftInfo + tlsf_BLOCK_OVERHEAD + (blockInfo & ^tlsf_TAGS_MASK) // keep left tags
		block.mmInfo = blockInfo
		// 'back' is Add below
	}

	right.mmInfo = rightInfo | tlsf_LEFTFREE
	// reference to right is no longer used now, hence rightInfo is not synced

	// we now know the size of the block
	size := blockInfo & ^tlsf_TAGS_MASK

	// Add 'back' to itself at the end of block
	*(*uintptr)(unsafe.Pointer(uintptr(unsafe.Pointer(right)) - tlsf_sizeofPointer)) = uintptr(unsafe.Pointer(block))

	// mapping_insert
	var (
		fl uintptr
		sl uint32
	)
	if size < uintptr(tlsf_SB_SIZE) {
		fl = 0
		sl = uint32(size >> tlsf_AL_BITS)
	} else {
		const inv = tlsf_sizeofPointer*8 - 1
		boundedSize := min(size, tlsf_BLOCK_MAXSIZE)
		fl = inv - clz(boundedSize)
		sl = uint32((boundedSize >> (fl - uintptr(tlsf_SL_BITS))) ^ (1 << tlsf_SL_BITS))
		fl -= uintptr(tlsf_SB_BITS) - 1
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
	if tlsf_DEBUG {
		assert(blockInfo&tlsf_FREE != 0, "must be free")
	}
	size := blockInfo & ^tlsf_TAGS_MASK
	if tlsf_DEBUG {
		assert(size >= tlsf_BLOCK_MINSIZE, "must be valid")
	}

	// mapping_insert
	var (
		fl uintptr
		sl uint32
	)
	if size < uintptr(tlsf_SB_SIZE) {
		fl = 0
		sl = uint32(size >> tlsf_AL_BITS)
	} else {
		const inv = tlsf_sizeofPointer*8 - 1
		boundedSize := min(size, tlsf_BLOCK_MAXSIZE)
		fl = inv - clz(boundedSize)
		sl = uint32((boundedSize >> (fl - uintptr(tlsf_SL_BITS))) ^ (1 << uintptr(tlsf_SL_BITS)))
		fl -= uintptr(tlsf_SB_BITS) - 1
	}
	if tlsf_DEBUG {
		assert(fl < uintptr(tlsf_FL_BITS) && sl < tlsf_SL_SIZE, "fl/sl out of range")
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
	if size < uintptr(tlsf_SB_SIZE) {
		fl = 0
		sl = uint32(size >> tlsf_AL_BITS)
	} else {
		const halfMaxSize = tlsf_BLOCK_MAXSIZE >> 1 // don't round last fl
		const inv = tlsf_sizeofPointer*8 - 1
		const invRound = inv - uintptr(tlsf_SL_BITS)

		var requestSize uintptr
		if size < halfMaxSize {
			requestSize = size + (1 << (invRound - clz(size))) - 1
		} else {
			requestSize = size
		}

		fl = inv - clz(requestSize)
		sl = uint32((requestSize >> (fl - uintptr(tlsf_SL_BITS))) ^ (1 << tlsf_SL_BITS))
		fl -= uintptr(tlsf_SB_BITS) - 1
	}
	if tlsf_DEBUG {
		assert(fl < uintptr(tlsf_FL_BITS) && sl < tlsf_SL_SIZE, "fl/sl out of range")
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
			if tlsf_DEBUG {
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
	if tlsf_DEBUG {
		assert(((size+tlsf_BLOCK_OVERHEAD)&tlsf_AL_MASK) == 0,
			"size must be aligned so the New block is")
	}
	// split if the block can hold another MINSIZE block incl. overhead
	remaining := (blockInfo & ^tlsf_TAGS_MASK) - size
	if remaining >= tlsf_BLOCK_OVERHEAD+tlsf_BLOCK_MINSIZE {
		block.mmInfo = size | (blockInfo & tlsf_LEFTFREE) // also discards FREE

		spare := (*tlsfBlock)(unsafe.Pointer(uintptr(unsafe.Pointer(block)) + tlsf_BLOCK_OVERHEAD + size))
		spare.mmInfo = (remaining - tlsf_BLOCK_OVERHEAD) | tlsf_FREE // not LEFTFREE
		a.insertBlock(spare)                                         // also sets 'back'

		// otherwise tag block as no longer FREE and right as no longer LEFTFREE
	} else {
		block.mmInfo = blockInfo & ^tlsf_FREE
		block.getRight().mmInfo &= ^tlsf_LEFTFREE
	}
}

// growMemory grows the pool by a number of 64kb pages to fit the required size
func (a *Allocator) growMemory(size uintptr) bool {
	if a.Grow == nil {
		return false
	}
	// Here, both rounding performed in searchBlock ...
	const halfMaxSize = tlsf_BLOCK_MAXSIZE >> 1
	if size < halfMaxSize { // don't round last fl
		const invRound = (tlsf_sizeofPointer*8 - 1) - uintptr(tlsf_SL_BITS)
		size += (1 << (invRound - clz(size))) - 1
	}
	// and additional BLOCK_OVERHEAD must be taken into account. If we are going
	// to merge with the tail block, that's one time, otherwise it's two times.
	var (
		pagesBefore         = a.Pages
		offset      uintptr = 0
	)
	if tlsf_BLOCK_OVERHEAD != uintptr(unsafe.Pointer(a.root.getTail())) {
		offset = 1
	}
	size += tlsf_BLOCK_OVERHEAD << ((uintptr(pagesBefore) << 16) - offset)
	pagesNeeded := ((int32(size) + 0xffff) & ^0xffff) >> 16

	addedPages, start, end := a.Grow(pagesBefore, pagesNeeded, size)
	if start == 0 || end == 0 {
		return false
	}
	if addedPages == 0 {
		addedPages = int32((end - start) / tlsf_PAGE_SIZE)
		if (end-start)%tlsf_PAGE_SIZE > 0 {
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
	if tlsf_DEBUG {
		assert(start <= end, "start must be <= end")
	}
	start = ((start + tlsf_BLOCK_OVERHEAD + tlsf_AL_MASK) & ^tlsf_AL_MASK) - tlsf_BLOCK_OVERHEAD
	end &= ^tlsf_AL_MASK

	var tail = a.root.getTail()
	var tailInfo uintptr = 0
	if tail != nil { // more memory
		if tlsf_DEBUG {
			assert(start >= uintptr(unsafe.Pointer(tail))+tlsf_BLOCK_OVERHEAD, "out of bounds")
		}

		// merge with current tail if adjacent
		const offsetToTail = tlsf_AL_SIZE
		if start-offsetToTail == uintptr(unsafe.Pointer(tail)) {
			start -= offsetToTail
			tailInfo = tail.mmInfo
		} else {
			// We don't do this, but a user might `memory.Grow` manually
			// leading to non-adjacent pages managed by Allocator.
		}

	} else if tlsf_DEBUG { // first memory
		assert(start >= uintptr(unsafe.Pointer(a.root))+tlsf_ROOT_SIZE, "starts after root")
	}

	// check if size is large enough for a free block and the tail block
	var size = end - start
	if size < tlsf_BLOCK_OVERHEAD+tlsf_BLOCK_MINSIZE+tlsf_BLOCK_OVERHEAD {
		return false
	}

	// left size is total minus its own and the zero-length tail's header
	var leftSize = size - 2*tlsf_BLOCK_OVERHEAD
	var left = (*tlsfBlock)(unsafe.Pointer(start))
	left.mmInfo = leftSize | tlsf_FREE | (tailInfo & tlsf_LEFTFREE)
	left.prev = 0
	left.next = 0

	// tail is a zero-length used block
	tail = (*tlsfBlock)(unsafe.Pointer(start + tlsf_BLOCK_OVERHEAD + leftSize))
	tail.mmInfo = 0 | tlsf_LEFTFREE
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
	if size <= tlsf_BLOCK_MINSIZE {
		return tlsf_BLOCK_MINSIZE
	} else {
		return ((size + tlsf_BLOCK_OVERHEAD + tlsf_AL_MASK) & ^tlsf_AL_MASK) - tlsf_BLOCK_OVERHEAD
	}
}

// Prepares and checks an allocation size.
func prepareSize(size uintptr) uintptr {
	if size > tlsf_BLOCK_MAXSIZE {
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
		if tlsf_DEBUG {
			assert(block != nil, "block must be found now")
		}
		if block == nil {
			return nil
		}
	}
	if tlsf_DEBUG {
		assert((block.mmInfo & ^tlsf_TAGS_MASK) >= payloadSize, "must fit")
	}

	a.removeBlock(block)
	a.prepareBlock(block, payloadSize)

	// update stats
	payloadSize = block.mmInfo & ^tlsf_TAGS_MASK
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
	var blockSize = blockInfo & ^tlsf_TAGS_MASK

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
	if rightInfo&tlsf_FREE != 0 {
		mergeSize := blockSize + tlsf_BLOCK_OVERHEAD + (rightInfo & ^tlsf_TAGS_MASK)
		if mergeSize >= payloadSize {
			a.removeBlock(right)
			block.mmInfo = (blockInfo & tlsf_TAGS_MASK) | mergeSize
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

	memcpy(unsafe.Pointer(uintptr(unsafe.Pointer(newBlock))+tlsf_BLOCK_OVERHEAD),
		unsafe.Pointer(uintptr(unsafe.Pointer(block))+tlsf_BLOCK_OVERHEAD),
		block.mmInfo & ^tlsf_TAGS_MASK)

	a.freeBlock(block)
	//maybeFreeBlock(a, block)

	return newBlock
}

func (a *Allocator) freeBlock(block *tlsfBlock) {
	size := block.mmInfo & ^tlsf_TAGS_MASK
	a.FreeSize += int64(size)
	a.AllocSize -= int64(size)
	a.Allocs--

	block.mmInfo = block.mmInfo | tlsf_FREE
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
	block := (*tlsfBlock)(unsafe.Pointer(ptr - tlsf_BLOCK_OVERHEAD))
	if !(ptr != 0 && ((ptr & tlsf_AL_MASK) == 0) && ((block.mmInfo & tlsf_FREE) == 0)) {
		panic("used block is not valid to be freed or reallocated")
	}
	return block
}

func assert(truthy bool, message string) {
	if !truthy {
		panic(message)
	}
}

func allocationSize(ptr unsafe.Pointer) uintptr {
	return ((*tlsfBlock)(unsafe.Pointer(uintptr(ptr) - tlsf_BLOCK_OVERHEAD))).mmInfo & ^tlsf_TAGS_MASK
}

func PrintDebugInfo() {
	println("ALIGNOF_U32		", int64(tlsf_ALIGN_U32))
	println("ALIGN_SIZE_LOG2	", int64(tlsf_ALIGN_SIZE_LOG2))
	println("U32_MAX			", ^uint32(0))
	println("PTR_MAX			", ^uintptr(0))
	println("AL_BITS			", int64(tlsf_AL_BITS))
	println("AL_SIZE			", int64(tlsf_AL_SIZE))
	println("AL_MASK			", int64(tlsf_AL_MASK))
	println("BLOCK_OVERHEAD	", int64(tlsf_BLOCK_OVERHEAD))
	println("BLOCK_MAXSIZE	", int64(tlsf_BLOCK_MAXSIZE))
	println("SL_BITS			", int64(tlsf_SL_BITS))
	println("SL_SIZE			", int64(tlsf_SL_SIZE))
	println("SB_BITS			", int64(tlsf_SB_BITS))
	println("SB_SIZE			", int64(tlsf_SB_SIZE))
	println("FL_BITS			", int64(tlsf_FL_BITS))
	println("FREE			", int64(tlsf_FREE))
	println("LEFTFREE		", int64(tlsf_LEFTFREE))
	println("TAGS_MASK		", int64(tlsf_TAGS_MASK))
	println("BLOCK_MINSIZE	", int64(tlsf_BLOCK_MINSIZE))
	println("SL_START		", int64(tlsf_SL_START))
	println("SL_END			", int64(tlsf_SL_END))
	println("HL_START		", int64(tlsf_HL_START))
	println("HL_END			", int64(tlsf_HL_END))
	println("ROOT_SIZE		", int64(tlsf_ROOT_SIZE))
}
