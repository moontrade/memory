package mem

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
