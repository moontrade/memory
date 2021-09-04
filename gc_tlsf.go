package runtime

import (
	"math/bits"
	"unsafe"
)

// Allocator
// === The TLSF (Two-Level Segregate Fit) memory allocator ===
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

const (
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

	tlsfDebug   = false
	tlsfTrace   = false
	tlsfRTrace  = false
	tlsfProfile = false

	// Overhead of a memory manager block.
	tlsf_BLOCK_OVERHEAD = unsafe.Sizeof(tlsfBLOCK{})
	// Block constants. A block must have a minimum size of three pointers so it can hold `prev`,
	// `next` and `back` if free.
	tlsf_BLOCK_MINSIZE = ((3*tlsf_sizeofPointer + tlsf_BLOCK_OVERHEAD + tlsf_AL_MASK) & ^tlsf_AL_MASK) - tlsf_BLOCK_OVERHEAD
	// Maximum size of a memory manager block's payload.
	tlsf_BLOCK_MAXSIZE = (1 << 30) - tlsf_BLOCK_OVERHEAD
	//tlsf_BLOCK_MAXSIZE = (1 << ((tlsf_ALIGN_SIZE_LOG2 + 1)*10)) - tlsf_BLOCK_OVERHEAD

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

func tlsfPrintInfo() {
	println("ALIGNOF_U32		", int64(tlsf_ALIGN_U32))
	println("ALIGN_SIZE_LOG2)	", int64(tlsf_ALIGN_SIZE_LOG2))
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

var allocator *tlsf

//export tlsfallocator
func tlsfGetAllocator() *tlsf {
	if allocator != nil {
		return allocator
	}
	allocator = newTLSF(1)
	return allocator
}

// tlsf
type tlsf struct {
	root         *tlsfRoot
	heapStart    uintptr
	heapEnd      uintptr
	heapSize     int64
	allocSize    int64
	maxAllocSize int64
	freeSize     int64
	allocs       int32
	pages        int32
}

//goland:noinspection GoVetUnsafePointer
func initTLSF(start, end uintptr, pages int32) *tlsf {
	// init allocator
	a := (*tlsf)(unsafe.Pointer(start))
	*a = tlsf{
		pages:     pages,
		heapStart: start,
		heapEnd:   end,
	}

	// init root
	rootOffset := (start + unsafe.Sizeof(tlsf{}) + tlsf_AL_MASK) & ^tlsf_AL_MASK
	root := (*tlsfRoot)(unsafe.Pointer(rootOffset))
	a.root = tlsfRootInit(root)

	// add initial memory
	tlsfAddMemory(a, rootOffset+tlsf_ROOT_SIZE, a.heapEnd)
	return a
}

//goland:noinspection GoVetUnsafePointer
func (p *tlsf) Alloc(size uintptr) unsafe.Pointer {
	return unsafe.Pointer(tlsfalloc(p, size))
}

//goland:noinspection GoVetUnsafePointer
func (p *tlsf) Realloc(ptr unsafe.Pointer, size uintptr) unsafe.Pointer {
	return unsafe.Pointer(tlsfrealloc(p, uintptr(ptr), size))
}

//goland:noinspection GoVetUnsafePointer
func (p *tlsf) Free(ptr unsafe.Pointer) {
	tlsffree(p, uintptr(ptr))
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
func tlsfGetFreeLeft(block *tlsfBlock) *tlsfBlock {
	return *(**tlsfBlock)(unsafe.Pointer(uintptr(unsafe.Pointer(block)) - tlsf_sizeofPointer))
}

// Gets the right block of a block by advancing to the right by its size.
func tlsfGetRight(block *tlsfBlock) *tlsfBlock {
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
type tlsfRoot struct {
	flMap uintptr
}

func tlsfRootInit(root *tlsfRoot) *tlsfRoot {
	root.flMap = 0
	tlsfSetTail(root, nil)
	for fl := uintptr(0); fl < uintptr(tlsf_FL_BITS); fl++ {
		tlsfSetSL(root, fl, 0)
		for sl := uint32(0); sl < tlsf_SL_SIZE; sl++ {
			tlsfSetHead(root, fl, sl, nil)
		}
	}
	return root
}

const (
	tlsf_SL_START  = tlsf_sizeofPointer
	tlsf_SL_END    = tlsf_SL_START + (uintptr(tlsf_FL_BITS) << tlsf_ALIGN_U32)
	tlsf_HL_START  = (tlsf_SL_END + tlsf_AL_MASK) &^ tlsf_AL_MASK
	tlsf_HL_END    = tlsf_HL_START + uintptr(tlsf_FL_BITS)*uintptr(tlsf_SL_SIZE)*tlsf_sizeofPointer
	tlsf_ROOT_SIZE = tlsf_HL_END + tlsf_sizeofPointer
)

// Gets the second level map of the specified first level.
func tlsfGetSL(root *tlsfRoot, fl uintptr) uint32 {
	return *(*uint32)(unsafe.Pointer(uintptr(unsafe.Pointer(root)) + (fl << tlsf_ALIGN_U32) + tlsf_SL_START))
}

// Sets the second level map of the specified first level.
func tlsfSetSL(root *tlsfRoot, fl uintptr, slMap uint32) {
	*(*uint32)(unsafe.Pointer(uintptr(unsafe.Pointer(root)) + (fl << tlsf_ALIGN_U32) + tlsf_SL_START)) = slMap
}

// Gets the head of the free list for the specified combination of first and second level.
func tlsfGetHead(root *tlsfRoot, fl uintptr, sl uint32) *tlsfBlock {
	return *(**tlsfBlock)(unsafe.Pointer(
		uintptr(unsafe.Pointer(root)) + (((fl << tlsf_SL_BITS) + uintptr(sl)) << tlsf_ALIGN_SIZE_LOG2) + tlsf_HL_START))
}

// Sets the head of the free list for the specified combination of first and second level.
func tlsfSetHead(root *tlsfRoot, fl uintptr, sl uint32, head *tlsfBlock) {
	*(*uintptr)(unsafe.Pointer(uintptr(
		unsafe.Pointer(root)) + tlsf_HL_START + (((fl << tlsf_SL_BITS) + uintptr(sl)) << tlsf_ALIGN_SIZE_LOG2))) = uintptr(unsafe.Pointer(head))
}

// Gets the tail block.
func tlsfGetTail(root *tlsfRoot) *tlsfBlock {
	return *(**tlsfBlock)(unsafe.Pointer(uintptr(unsafe.Pointer(root)) + tlsf_HL_END))
}

// Sets the tail block.
func tlsfSetTail(root *tlsfRoot, tail *tlsfBlock) {
	*(*uintptr)(unsafe.Pointer(uintptr(unsafe.Pointer(root)) + tlsf_HL_END)) = uintptr(unsafe.Pointer(tail))
}

// Inserts a previously used block back into the free list.
func tlsfInsertBlock(p *tlsf, block *tlsfBlock) {
	var (
		root      = p.root
		blockInfo = block.mmInfo
		right     = tlsfGetRight(block)
		rightInfo = right.mmInfo
	)
	//(blockInfo & FREE)

	// merge with right block if also free
	if rightInfo&tlsf_FREE != 0 {
		tlsfRemoveBlock(root, right)
		blockInfo = blockInfo + tlsf_BLOCK_OVERHEAD + (rightInfo & ^tlsf_TAGS_MASK) // keep block tags
		block.mmInfo = blockInfo
		right = tlsfGetRight(block)
		rightInfo = right.mmInfo
		// 'back' is set below
	}

	// merge with left block if also free
	if blockInfo&tlsf_LEFTFREE != 0 {
		left := tlsfGetFreeLeft(block)
		leftInfo := left.mmInfo
		if tlsfDebug {
			tlsfAssert(leftInfo&tlsf_FREE != 0, "must be free according to right tags")
		}
		tlsfRemoveBlock(root, left)
		block = left
		blockInfo = leftInfo + tlsf_BLOCK_OVERHEAD + (blockInfo & ^tlsf_TAGS_MASK) // keep left tags
		block.mmInfo = blockInfo
		// 'back' is set below
	}

	right.mmInfo = rightInfo | tlsf_LEFTFREE
	// reference to right is no longer used now, hence rightInfo is not synced

	// we now know the size of the block
	size := blockInfo & ^tlsf_TAGS_MASK

	// set 'back' to itself at the end of block
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
		boundedSize := tlsfMin(size, tlsf_BLOCK_MAXSIZE)
		fl = inv - tlsfClz(boundedSize)
		sl = uint32((boundedSize >> (fl - uintptr(tlsf_SL_BITS))) ^ (1 << tlsf_SL_BITS))
		fl -= uintptr(tlsf_SB_BITS) - 1
	}

	// perform insertion
	head := tlsfGetHead(root, fl, sl)
	block.prev = 0
	block.next = uintptr(unsafe.Pointer(head))
	if head != nil {
		head.prev = uintptr(unsafe.Pointer(block))
	}
	tlsfSetHead(root, fl, sl, block)

	// update first and second level maps
	root.flMap |= 1 << fl
	tlsfSetSL(root, fl, tlsfGetSL(root, fl)|(1<<sl))
}

//goland:noinspection GoVetUnsafePointer
func tlsfRemoveBlock(root *tlsfRoot, block *tlsfBlock) {
	blockInfo := block.mmInfo
	if tlsfDebug {
		tlsfAssert(blockInfo&tlsf_FREE != 0, "must be free")
	}
	size := blockInfo & ^tlsf_TAGS_MASK
	if tlsfDebug {
		tlsfAssert(size >= tlsf_BLOCK_MINSIZE, "must be valid")
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
		boundedSize := tlsfMin(size, tlsf_BLOCK_MAXSIZE)
		fl = inv - tlsfClz(boundedSize)
		sl = uint32((boundedSize >> (fl - uintptr(tlsf_SL_BITS))) ^ (1 << uintptr(tlsf_SL_BITS)))
		fl -= uintptr(tlsf_SB_BITS) - 1
	}
	if tlsfDebug {
		tlsfAssert(fl < uintptr(tlsf_FL_BITS) && sl < tlsf_SL_SIZE, "fl/sl out of range")
	}

	// link previous and next free block
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
	if block == tlsfGetHead(root, fl, sl) {
		tlsfSetHead(root, fl, sl, (*tlsfBlock)(unsafe.Pointer(next)))

		// clear second level map if head is empty now
		if next == 0 {
			slMap := tlsfGetSL(root, fl)
			slMap &= ^(1 << sl)
			tlsfSetSL(root, fl, slMap)

			// clear first level map if second level is empty now
			if slMap == 0 {
				root.flMap &= ^(1 << fl)
			}
		}
	}
	// note: does not alter left/back because it is likely that splitting
	// is performed afterwards, invalidating those changes. so, the caller
	// must perform those updates.
}

// Searches for a free block of at least the specified size.
func tlsfSearchBlock(root *tlsfRoot, size uintptr) *tlsfBlock {
	// mapping_search
	var (
		fl uintptr
		sl uint32
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
			requestSize = size + (1 << (invRound - tlsfClz(size))) - 1
		} else {
			requestSize = size
		}

		fl = inv - tlsfClz(requestSize)
		sl = uint32((requestSize >> (fl - uintptr(tlsf_SL_BITS))) ^ (1 << tlsf_SL_BITS))
		fl -= uintptr(tlsf_SB_BITS) - 1
	}
	if tlsfDebug {
		tlsfAssert(fl < uintptr(tlsf_FL_BITS) && sl < tlsf_SL_SIZE, "fl/sl out of range")
	}

	// search second level
	var slMap = tlsfGetSL(root, fl) & (^uint32(0) << sl)
	var head *tlsfBlock
	if slMap == 0 {
		// search next larger first level
		flMap := root.flMap & (^uintptr(0) << (fl + 1))
		if flMap == 0 {
			head = nil
		} else {
			fl = tlsfCtz(flMap)
			slMap = tlsfGetSL(root, fl)
			if tlsfDebug {
				tlsfAssert(slMap != 0, "can't be zero if fl points here")
			}
			head = tlsfGetHead(root, fl, tlsfCtz32(slMap))
		}
	} else {
		head = tlsfGetHead(root, fl, tlsfCtz32(slMap))
	}

	return head
}

func tlsfPrepareBlock(p *tlsf, block *tlsfBlock, size uintptr) {
	blockInfo := block.mmInfo
	if tlsfDebug {
		tlsfAssert(((size+tlsf_BLOCK_OVERHEAD)&tlsf_AL_MASK) == 0,
			"size must be aligned so the new block is")
	}
	// split if the block can hold another MINSIZE block incl. overhead
	remaining := (blockInfo & ^tlsf_TAGS_MASK) - size
	if remaining >= tlsf_BLOCK_OVERHEAD+tlsf_BLOCK_MINSIZE {
		block.mmInfo = size | (blockInfo & tlsf_LEFTFREE) // also discards FREE

		spare := (*tlsfBlock)(unsafe.Pointer(uintptr(unsafe.Pointer(block)) + tlsf_BLOCK_OVERHEAD + size))
		spare.mmInfo = (remaining - tlsf_BLOCK_OVERHEAD) | tlsf_FREE // not LEFTFREE
		tlsfInsertBlock(p, spare)                                    // also sets 'back'

		// otherwise tag block as no longer FREE and right as no longer LEFTFREE
	} else {
		block.mmInfo = blockInfo & ^tlsf_FREE
		tlsfGetRight(block).mmInfo &= ^tlsf_LEFTFREE
	}
}

// tlsfGrowMemory grows the pool by a number of 64kb pages to fit the required size
func tlsfGrowMemory(p *tlsf, size uintptr) {
	// Here, both rounding performed in searchBlock ...
	const halfMaxSize = tlsf_BLOCK_MAXSIZE >> 1
	if size < halfMaxSize { // don't round last fl
		const invRound = (tlsf_sizeofPointer*8 - 1) - uintptr(tlsf_SL_BITS)
		size += (1 << (invRound - tlsfClz(size))) - 1
	}
	// and additional BLOCK_OVERHEAD must be taken into account. If we are going
	// to merge with the tail block, that's one time, otherwise it's two times.
	var pagesBefore = int32(p.pages)
	var offset uintptr = 0
	if tlsf_BLOCK_OVERHEAD != uintptr(unsafe.Pointer(tlsfGetTail(p.root))) {
		offset = 1
	}
	size += tlsf_BLOCK_OVERHEAD << ((uintptr(pagesBefore) << 16) - offset)
	var pagesNeeded = ((int32(size) + 0xffff) & ^0xffff) >> 16
	var pagesWanted = tlsfMax32(pagesBefore, pagesNeeded) // double memory

	start, end := p.Grow(pagesWanted)
	if start == 0 {
		start, end = p.Grow(pagesNeeded)
		if start == 0 {
			panic("out of memory")
		}
	}
	tlsfAddMemory(p, start, end)
}

// tlsfAddMemory adds the newly allocated memory to the TLSF bitmaps
//goland:noinspection GoVetUnsafePointer
func tlsfAddMemory(p *tlsf, start, end uintptr) bool {
	if tlsfDebug {
		tlsfAssert(start <= end, "start must be <= end")
	}
	start = ((start + tlsf_BLOCK_OVERHEAD + tlsf_AL_MASK) & ^tlsf_AL_MASK) - tlsf_BLOCK_OVERHEAD
	end &= ^tlsf_AL_MASK

	var tail = tlsfGetTail(p.root)
	var tailInfo uintptr = 0
	if tail != nil { // more memory
		if tlsfDebug {
			tlsfAssert(start >= uintptr(unsafe.Pointer(tail))+tlsf_BLOCK_OVERHEAD, "out of bounds")
		}

		// merge with current tail if adjacent
		const offsetToTail = tlsf_AL_SIZE
		if start-offsetToTail == uintptr(unsafe.Pointer(tail)) {
			start -= offsetToTail
			tailInfo = tail.mmInfo
		} else {
			// We don't do this, but a user might `memory.grow` manually
			// leading to non-adjacent pages managed by TLSF.
		}

	} else if tlsfDebug { // first memory
		tlsfAssert(start >= uintptr(unsafe.Pointer(p.root))+tlsf_ROOT_SIZE, "starts after root")
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
	tlsfSetTail(p.root, tail)

	p.freeSize += int64(leftSize)
	tlsfInsertBlock(p, left) // also merges with free left before tail / sets 'back'

	return true
}

// Computes the size (excl. header) of a block.
func tlsfComputeSize(size uintptr) uintptr {
	// Size must be large enough and aligned minus preceeding overhead
	if size <= tlsf_BLOCK_MINSIZE {
		return tlsf_BLOCK_MINSIZE
	} else {
		return ((size + tlsf_BLOCK_OVERHEAD + tlsf_AL_MASK) & ^tlsf_AL_MASK) - tlsf_BLOCK_OVERHEAD
	}
}

// Prepares and checks an allocation size.
func tlsfPrepareSize(size uintptr) uintptr {
	if size > tlsf_BLOCK_MAXSIZE {
		panic("allocation too large")
	}
	return tlsfComputeSize(size)
}

// Allocates a block of the specified size.
func tlsfAllocateBlock(p *tlsf, size uintptr) *tlsfBlock {
	var payloadSize = tlsfPrepareSize(size)
	var block = tlsfSearchBlock(p.root, payloadSize)
	if block == nil {
		tlsfGrowMemory(p, payloadSize)
		block = tlsfSearchBlock(p.root, payloadSize)
		if tlsfDebug {
			tlsfAssert(block != nil, "block must be found now")
		}
	}
	if tlsfDebug {
		tlsfAssert((block.mmInfo & ^tlsf_TAGS_MASK) >= payloadSize, "must fit")
	}

	tlsfRemoveBlock(p.root, block)
	tlsfPrepareBlock(p, block, payloadSize)

	// update stats
	payloadSize = block.mmInfo & ^tlsf_TAGS_MASK
	p.allocSize += int64(payloadSize)
	if p.allocSize > p.maxAllocSize {
		p.maxAllocSize = p.allocSize
	}
	p.freeSize -= int64(payloadSize)
	p.allocs++

	// return block
	return block
}

func tlsfReallocateBlock(p *tlsf, block *tlsfBlock, size uintptr) *tlsfBlock {
	var payloadSize = tlsfPrepareSize(size)
	var blockInfo = block.mmInfo
	var blockSize = blockInfo & ^tlsf_TAGS_MASK

	// possibly split and update runtime size if it still fits
	if payloadSize <= blockSize {
		tlsfPrepareBlock(p, block, payloadSize)
		//if (isDefined(ASC_RTRACE)) {
		//	if (payloadSize != blockSize) onresize(block, BLOCK_OVERHEAD + blockSize);
		//}
		return block
	}

	// merge with right free block if merger is large enough
	var (
		right     = tlsfGetRight(block)
		rightInfo = right.mmInfo
	)
	if rightInfo&tlsf_FREE != 0 {
		mergeSize := blockSize + tlsf_BLOCK_OVERHEAD + (rightInfo & ^tlsf_TAGS_MASK)
		if mergeSize >= payloadSize {
			tlsfRemoveBlock(p.root, right)
			block.mmInfo = (blockInfo & tlsf_TAGS_MASK) | mergeSize
			tlsfPrepareBlock(p, block, payloadSize)
			//if (isDefined(ASC_RTRACE)) onresize(block, BLOCK_OVERHEAD + blockSize);
			return block
		}
	}

	// otherwise, move the block
	return tlsfMoveBlock(p, block, size)
}

func tlsfMoveBlock(p *tlsf, block *tlsfBlock, newSize uintptr) *tlsfBlock {
	var newBlock = tlsfAllocateBlock(p, newSize)

	memcpy(unsafe.Pointer(uintptr(unsafe.Pointer(newBlock))+tlsf_BLOCK_OVERHEAD),
		unsafe.Pointer(uintptr(unsafe.Pointer(block))+tlsf_BLOCK_OVERHEAD),
		block.mmInfo & ^tlsf_TAGS_MASK)

	tlsfFreeBlock(p, block)
	//maybeFreeBlock(p, block)

	return newBlock
}

func tlsfFreeBlock(p *tlsf, block *tlsfBlock) {
	size := block.mmInfo & ^tlsf_TAGS_MASK
	p.freeSize += int64(size)
	p.allocSize -= int64(size)
	p.allocs--

	block.mmInfo = block.mmInfo | tlsf_FREE
	tlsfInsertBlock(p, block)
}

func tlsfMin(l, r uintptr) uintptr {
	if l < r {
		return l
	}
	return r
}

func tlsfMax32(l, r int32) int32 {
	if l > r {
		return l
	}
	return r
}

func tlsfClz(value uintptr) uintptr {
	return uintptr(bits.LeadingZeros(uint(value)))
}

func tlsfCtz(value uintptr) uintptr {
	return uintptr(bits.TrailingZeros(uint(value)))
}

func tlsfCtz32(value uint32) uint32 {
	return uint32(bits.TrailingZeros32(value))
}

func tlsfCheckUsedBlock(ptr uintptr) *tlsfBlock {
	block := (*tlsfBlock)(unsafe.Pointer(ptr - tlsf_BLOCK_OVERHEAD))
	if !(ptr != 0 && ((ptr & tlsf_AL_MASK) == 0) && ((block.mmInfo & tlsf_FREE) == 0)) {
		panic("used block is not valid to be freed or reallocated")
	}
	return block
}

func tlsfAssert(truthy bool, message string) {
	if !truthy {
		panic(message)
	}
}

func tlsfalloc(pool *tlsf, size uintptr) uintptr {
	p := uintptr(unsafe.Pointer(tlsfAllocateBlock(pool, size)))
	if p == 0 {
		return 0
	}
	return p + tlsf_BLOCK_OVERHEAD
}

func tlsfrealloc(p *tlsf, ptr, size uintptr) uintptr {
	//if ptr < uintptr(unsafe.Pointer(p.root)) {
	//	//if ptr < heapStart {
	//	return uintptr(unsafe.Pointer(tlsfMoveBlock(p, tlsfCheckUsedBlock(ptr), size))) + tlsf_BLOCK_OVERHEAD
	//} else {
	//	return uintptr(unsafe.Pointer(tlsfReallocateBlock(p, tlsfCheckUsedBlock(ptr), size))) + tlsf_BLOCK_OVERHEAD
	//}

	//return uintptr(unsafe.Pointer(tlsfReallocateBlock(p, tlsfCheckUsedBlock(ptr), size))) + tlsf_BLOCK_OVERHEAD
	return uintptr(unsafe.Pointer(tlsfMoveBlock(p, tlsfCheckUsedBlock(ptr), size))) + tlsf_BLOCK_OVERHEAD
}

func tlsffree(p *tlsf, ptr uintptr) {
	tlsfFreeBlock(p, tlsfCheckUsedBlock(ptr))
}

//export memalloc
func memalloc(size uintptr) unsafe.Pointer {
	p := allocator.Alloc(size)
	memzero(p, size)
	return p
}

//export memrealloc
func memrealloc(ptr unsafe.Pointer, size uintptr) unsafe.Pointer {
	return allocator.Realloc(ptr, size)
}

//export memfree
func memfree(ptr unsafe.Pointer) {
	allocator.Free(ptr)
}