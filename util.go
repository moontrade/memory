package memory

// divOrZero safely divides a by b or returns zero if either are zero
func divOrZero(a, b int64) int64 {
	if a == 0 || b == 0 {
		return 0
	}
	return a / b
}

func toMicros(v int64) int64 {
	return divOrZero(v, 1000)
}

func toMillis(v int64) int64 {
	return divOrZero(v, 1000000)
}

const (
	nanosSuffix   = "ns"
	microsSuffix  = "µs"
	millisSuffix  = "ms"
	secondsSuffix = "secs"
)

func toTimeName(v int64) string {
	switch {
	case v < 1000:
		return nanosSuffix
	case v < 1000000:
		return microsSuffix
	case v < 1000000000:
		return millisSuffix
	default:
		return secondsSuffix
	}
}

const (
	bitsize       = 32 << (^uint(0) >> 63)
	maxint        = int(1<<(bitsize-1) - 1)
	maxintHeadBit = 1 << (bitsize - 2)
)

// LogarithmicRange iterates from ceiled to power of two min to max,
// calling cb on each iteration.
func LogarithmicRange(min, max int, cb func(int)) {
	if min == 0 {
		min = 1
	}
	for n := ceilToPowerOfTwo(min); n <= max; n <<= 1 {
		cb(n)
	}
}

// IsPowerOfTwo reports whether given integer is a power of two.
func IsPowerOfTwo(n int) bool {
	return n&(n-1) == 0
}

// ceilToPowerOfTwo returns the least power of two integer value greater than
// or equal to n.
func ceilToPowerOfTwo(n int) int {
	if n&maxintHeadBit != 0 && n > maxintHeadBit {
		panic("argument is too large")
	}
	if n <= 2 {
		return n
	}
	n--
	n = fillBits(n)
	n++
	return n
}

// FloorToPowerOfTwo returns the greatest power of two integer value less than
// or equal to n.
func FloorToPowerOfTwo(n int) int {
	if n <= 2 {
		return n
	}
	n = fillBits(n)
	n >>= 1
	n++
	return n
}

func fillBits(n int) int {
	n |= n >> 1
	n |= n >> 2
	n |= n >> 4
	n |= n >> 8
	n |= n >> 16
	n |= n >> 32
	return n
}
