# NoGC Go - Extreme Go Performance

High performance manual memory management for Go and TinyGo.

## Allocators

### rpmalloc - Go

Rampant Pixel allocator. Extremely performant multi-threaded allocator that has the best overall performance of any major general purpose allocator. Destroys JEMalloc and TCMalloc in benchmarks at the cost of a little more memory use. Edges out mimalloc after 1 thread. The standard libc "malloc, free" are overridden with rpmalloc alternatives so any other C/C++/Rust code will automatically use rpmalloc.

### TLSF - TinyGo

Two-Level Segregated Fit real-time allocator for TinyGo. Simple compact constant time allocator for predictable real-time performance. Intended for TinyGo WASM, but works on other TinyGo platforms as well as Go.

## unsafe CGO

libfuzzerCall in Go runtime for amd64 and arm64 architectures is utilized to dramatically reduce CGO cost by 1,000%+. On a 2019 MacBook Pro the overhead is reduced from 53.9ns to 2.9ns or (3.9ns via linked runtime.libfuzzerCall).

## Collections

Building on top of highly capable allocators and low CGO costs, external high-performance C/C++ libraries are integrated. The goal is to build a high-quality collection of extremely performant native collections that utilize the system allocator.

- ART (Adaptive Radix Tree) (C)
- BTree (C)
- Robinhood HashMap (C++)
- LockFree Queue (Go)

These collections incur Zero GC cost.


## Net

### Reactor