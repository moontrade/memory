# TLSF (Two-Level Segregate Fit) memory allocator for TinyGo/Go


```text
=== TLSF (Two-Level Segregate Fit) memory allocator ===

TLSF is a general purpose dynamic memory allocator specifically designed to meet real-time requirements:

    Bounded Response Time - The worst-case execution time (WCET) of memory allocation
                            and deallocation Has got to be known in advance and be
                            independent of application data. Allocator Has a constant
                            cost O(1).

                     Fast - Additionally to a bounded cost, the allocator Has to be
                            efficient and fast enough. Allocator executes a maximum
                            of 168 processor instructions in a x86 architecture.
                            Depending on the compiler version and optimisation flags,
                            it can be slightly lower or higher.

    Efficient Memory Use - 	Traditionally, real-time systems run for long periods of
                            time and some (embedded applications), have strong constraints
                            of memory size. Fragmentation can have a significant impact on
                            such systems. It can increase  dramatically, and degrade the
                            system performance. A way to measure this efficiency is the
                            memory fragmentation incurred by the allocator. Allocator has
                            been tested in hundreds of different loads (real-time tasks,
                            general purpose applications, etc.) obtaining an average
                            fragmentation lower than 15 %. The maximum fragmentation
                            measured is lower than 25%.

Memory can be added on demand and is a multiple of 64kb pages. Grow is used to allocate new
memory to be added to the allocator. Each Grow must provide a contiguous chunk of memory.
However, the allocator may be comprised of many contiguous chunks which are not contiguous
of each other. There is not a mechanism for shrinking the memory. Supplied Grow function
can effectively limit how big the allocator can get. If a zero pointer is returned it will
cause an out-of-memory situation which is propagated as a nil pointer being returned from
Alloc. It's up to the application to decide how to handle such scenarios.

see: http:www.gii.upv.es/tlsf/
see: https://github.com/AssemblyScript/assemblyscript/blob/main/std/assembly/rt/tlsf.ts
     The allocator is largely ported from AssemblyScript

- `ffs(x)` is equivalent to `ctz(x)` with x != 0
- `fls(x)` is equivalent to `sizeof(x) * 8 - clz(x) - 1`

╒══════════════ Block size interpretation (32-bit) ═════════════╕
   3                   2                   1
 1 0 9 8 7 6 5 4 3 2 1 0 9 8 7 6 5 4 3 2 1 0 9 8 7 6 5 4 3 2 1 0  bits
├─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┼─┴─┴─┴─╫─┴─┴─┴─┤
│ |                    FL                       │ SB = SL + AL  │ ◄─ usize
└───────────────────────────────────────────────┴───────╨───────┘
FL: first level, SL: second level, AL: alignment, SB: small block

 ╒══════════════ Block size interpretation (32-bit) ═════════════╕
    3                   2                   1
  1 0 9 8 7 6 5 4 3 2 1 0 9 8 7 6 5 4 3 2 1 0 9 8 7 6 5 4 3 2 1 0  bits
 ├─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┼─┴─┴─┴─╫─┴─┴─┴─┤
 │ |                    FL                       │ SB = SL + AL  │ ◄─ usize
 └───────────────────────────────────────────────┴───────╨───────┘
 FL: first level, SL: second level, AL: alignment, SB: small block


 Memory manager

 ╒════════════ Memory manager block layout (32-bit) ═════════════╕
    3                   2                   1
  1 0 9 8 7 6 5 4 3 2 1 0 9 8 7 6 5 4 3 2 1 0 9 8 7 6 5 4 3 2 1 0  bits
 ├─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┤
 │                           MM info                             │ -4
 ╞>ptr═══════════════════════════════════════════════════════════╡
 │                              ...                              │


 ╒════════════════════ Block layout (32-bit) ════════════════════╕
    3                   2                   1
  1 0 9 8 7 6 5 4 3 2 1 0 9 8 7 6 5 4 3 2 1 0 9 8 7 6 5 4 3 2 1 0  bits
 ├─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┼─┼─┤            ┐
 │                          size                             │L│F│ ◄─┐ info   overhead
 ╞>ptr═══════════════════════════════════════════════════════╧═╧═╡   │        ┘
 │                        if free: ◄ prev                        │ ◄─┤ usize
 ├───────────────────────────────────────────────────────────────┤   │
 │                        if free: next ►                        │ ◄─┤
 ├───────────────────────────────────────────────────────────────┤   │
 │                             ...                               │   │ >= 0
 ├───────────────────────────────────────────────────────────────┤   │
 │                        if free: back ▲                        │ ◄─┘
 └───────────────────────────────────────────────────────────────┘ >= MIN SIZE
 F: FREE, L: LEFTFREE
 
 
 ╒═════════════════════ Root layout (32-bit) ════════════════════╕
    3                   2                   1
  1 0 9 8 7 6 5 4 3 2 1 0 9 8 7 6 5 4 3 2 1 0 9 8 7 6 5 4 3 2 1 0  bits
 ├─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┴─┤          ┐
 │        0        |           flMap                            S│ ◄────┐
 ╞═══════════════════════════════════════════════════════════════╡      │
 │                           slMap[0] S                          │ ◄─┐  │
 ├───────────────────────────────────────────────────────────────┤   │  │
 │                           slMap[1]                            │ ◄─┤  │
 ├───────────────────────────────────────────────────────────────┤  uint32 │
 │                           slMap[22]                           │ ◄─┘  │
 ╞═══════════════════════════════════════════════════════════════╡    usize
 │                            head[0]                            │ ◄────┤
 ├───────────────────────────────────────────────────────────────┤      │
 │                              ...                              │ ◄────┤
 ├───────────────────────────────────────────────────────────────┤      │
 │                           head[367]                           │ ◄────┤
 ╞═══════════════════════════════════════════════════════════════╡      │
 │                             tail                              │ ◄────┘
 └───────────────────────────────────────────────────────────────┘   SIZE   ┘
 S: Small blocks map
 
                    [00]: < 256B (SB)  [12]: < 1M
                    [01]: < 512B       [13]: < 2M
                    [02]: < 1K         [14]: < 4M
                    [03]: < 2K         [15]: < 8M
                    [04]: < 4K         [16]: < 16M
                    [05]: < 8K         [17]: < 32M
                    [06]: < 16K        [18]: < 64M
                    [07]: < 32K        [19]: < 128M
                    [08]: < 64K        [20]: < 256M
                    [09]: < 128K       [21]: < 512M
                    [10]: < 256K       [22]: <= 1G - OVERHEAD
                    [11]: < 512K

WASM VMs limit to 2GB total (currently), making one 1G block max (or three 512M etc.) due to block overhead
```