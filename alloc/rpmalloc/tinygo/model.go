package rpmalloc

// #include "rpmalloc.h"
import "C"

type Config struct {
	//! Map memory pages for the given number of bytes. The returned address MUST be
	//  aligned to the rpmalloc span size, which will always be a power of two.
	//  Optionally the function can store an alignment offset in the offset variable
	//  in case it performs alignment and the returned pointer is offset from the
	//  actual start of the memory region due to this alignment. The alignment offset
	//  will be passed to the memory unmap function. The alignment offset MUST NOT be
	//  larger than 65535 (storable in an uint16_t), if it is you must use natural
	//  alignment to shift it into 16 bits. If you set a memory_map function, you
	//  must also set a memory_unmap function or else the default implementation will
	//  be used for both.
	MemoryMap uintptr
	//! Unmap the memory pages starting at address and spanning the given number of bytes.
	//  If release is set to non-zero, the unmap is for an entire span range as returned by
	//  a previous libfuzzerCall to memory_map and that the entire range should be released. The
	//  release argument holds the size of the entire span range. If release is set to 0,
	//  the unmap is a partial decommit of a subset of the mapped memory range.
	//  If you set a memory_unmap function, you must also set a memory_map function or
	//  else the default implementation will be used for both.
	MemoryUnmap uintptr
	//! Called when an assert fails, if asserts are enabled. Will use the standard assert()
	//  if this is not set.
	ErrorCallback uintptr
	//! Called when a libfuzzerCall to map memory pages fails (out of memory). If this callback is
	//  not set or returns zero the library will return a null pointer in the allocation
	//  libfuzzerCall. If this callback returns non-zero the map libfuzzerCall will be retried. The argument
	//  passed is the number of bytes that was requested in the map libfuzzerCall. Only used if
	//  the default system memory map function is used (memory_map callback is not set).
	MapFailCallback uintptr
	//! Size of memory pages. The page size MUST be a power of two. All memory mapping
	//  requests to memory_map will be made with size set to a multiple of the page size.
	//  Used if RPMALLOC_CONFIGURABLE is defined to 1, otherwise system page size is used.
	PageSize uintptr
	//! Size of a span of memory blocks. MUST be a power of two, and in [4096,262144]
	//  range (unless 0 - set to 0 to use the default span size). Used if RPMALLOC_CONFIGURABLE
	//  is defined to 1.
	SpanSize uintptr
	//! Number of spans to map at each request to map new virtual memory blocks. This can
	//  be used to minimize the system libfuzzerCall overhead at the cost of virtual memory address
	//  space. The extra mapped pages will not be written until actually used, so physical
	//  committed memory should not be affected in the default implementation. Will be
	//  aligned to a multiple of spans that match memory page size in case of huge pages.
	SpanMapCount uintptr
	//! Enable use of large/huge pages. If this flag is set to non-zero and page size is
	//  zero, the allocator will try to enable huge pages and auto detect the configuration.
	//  If this is set to non-zero and page_size is also non-zero, the allocator will
	//  assume huge pages have been configured and enabled prior to initializing the
	//  allocator.
	//  For Windows, see https://docs.microsoft.com/en-us/windows/desktop/memory/large-page-support
	//  For Linux, see https://www.kernel.org/doc/Documentation/vm/hugetlbpage.txt
	EnableHugePages int32
	Unused          int32
}

type GlobalStats struct {
	//! Current amount of virtual memory mapped, all of which might not have been committed (only if ENABLE_STATISTICS=1)
	Mapped uintptr
	//! Peak amount of virtual memory mapped, all of which might not have been committed (only if ENABLE_STATISTICS=1)
	MappedPeak uintptr
	//! Current amount of memory in global caches for small and medium sizes (<32KiB)
	Cached uintptr
	//! Current amount of memory allocated in huge allocations, i.e larger than LARGE_SIZE_LIMIT which is 2MiB by default (only if ENABLE_STATISTICS=1)
	HugeAlloc uintptr
	//! Peak amount of memory allocated in huge allocations, i.e larger than LARGE_SIZE_LIMIT which is 2MiB by default (only if ENABLE_STATISTICS=1)
	HugeAllocPeak uintptr
	//! Total amount of memory mapped since initialization (only if ENABLE_STATISTICS=1)
	MappedTotal uintptr
	//! Total amount of memory unmapped since initialization  (only if ENABLE_STATISTICS=1)
	UnmappedTotal uintptr
}

type ThreadStats struct {
	//! Current number of bytes available in thread size class caches for small and medium sizes (<32KiB)
	SizeCache uintptr
	//! Current number of bytes available in thread span caches for small and medium sizes (<32KiB)
	SpanCache uintptr
	//! Total number of bytes transitioned from thread cache to global cache (only if ENABLE_STATISTICS=1)
	ThreadToGlobal uintptr
	//! Total number of bytes transitioned from global cache to thread cache (only if ENABLE_STATISTICS=1)
	GlobalToThread uintptr
	//! Per span count statistics (only if ENABLE_STATISTICS=1)
	SpanUse [64]SpanStats
	//! Per size class statistics (only if ENABLE_STATISTICS=1)
	SizeUse [128]SizeUse
}

type SpanStats struct {
	//! Currently used number of spans
	Current uintptr
	//! High water mark of spans used
	Peak uintptr
	//! Number of spans transitioned to global cache
	ToGlobal uintptr
	//! Number of spans transitioned from global cache
	FromGlobal uintptr
	//! Number of spans transitioned to thread cache
	ToCache uintptr
	//! Number of spans transitioned from thread cache
	FromCache uintptr
	//! Number of spans transitioned to reserved state
	ToReserved uintptr
	//! Number of spans transitioned from reserved state
	FromReserved uintptr
	//! Number of raw memory map calls (not hitting the reserve spans but resulting in actual OS mmap calls)
	MapCalls uintptr
}

type SizeUse struct {
	//! Current number of allocations
	AllocCurrent uintptr
	//! Peak number of allocations
	AllocPeak uintptr
	//! Total number of allocations
	AllocTotal uintptr
	//! Total number of frees
	FreeTotal uintptr
	//! Number of spans transitioned to cache
	SpansToCache uintptr
	//! Number of spans transitioned from cache
	SpansFromCache uintptr
	//! Number of spans transitioned from reserved state
	SpansFromReserved uintptr
	//! Number of raw memory map calls (not hitting the reserve spans but resulting in actual OS mmap calls)
	MapCalls uintptr
}
