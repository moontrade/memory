package alloc

import "math"

// HeapStats provides the metrics of an Allocator
type HeapStats struct {
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

func (s *HeapStats) Fragmentation() float32 {
	if s.fragmentation != 0 {
		return s.fragmentation
	}
	if s.HeapSize == 0 || s.PeakAllocSize == 0 {
		return 0
	}
	pct := float64(s.HeapSize-s.PeakAllocSize) / float64(s.HeapSize)
	s.fragmentation = float32(math.Floor(pct*100) / 100)
	return s.fragmentation
}
