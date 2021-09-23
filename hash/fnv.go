package hash

/////////////////////////////////////////////////////////////////////////////////////////////////////////
// fnv
/////////////////////////////////////////////////////////////////////////////////////////////////////////

func FNV32a(v uint32) uint32 {
	const (
		offset32 = uint32(2166136261)
		prime32  = uint32(16777619)
	)
	hash := offset32
	hash ^= uint32(byte(v))
	hash *= prime32
	hash ^= uint32(byte(v >> 8))
	hash *= prime32
	hash ^= uint32(byte(v >> 16))
	hash *= prime32
	hash ^= uint32(byte(v >> 24))
	hash *= prime32
	return hash
}

func FNV32(v uint32) uint32 {
	const (
		offset32 = uint32(2166136261)
		prime32  = uint32(16777619)
	)
	hash := offset32
	hash *= prime32
	hash ^= uint32(byte(v))
	hash *= prime32
	hash ^= uint32(byte(v >> 8))
	hash *= prime32
	hash ^= uint32(byte(v >> 16))
	hash *= prime32
	hash ^= uint32(byte(v >> 24))
	return hash
}

func FNV64a(v uint64) uint64 {
	const (
		offset64 = uint64(14695981039346656037)
		prime64  = uint64(1099511628211)
	)
	hash := offset64
	hash ^= uint64(byte(v))
	hash *= prime64
	hash ^= uint64(byte(v >> 8))
	hash *= prime64
	hash ^= uint64(byte(v >> 16))
	hash *= prime64
	hash ^= uint64(byte(v >> 24))
	hash *= prime64
	hash ^= uint64(byte(v >> 32))
	hash *= prime64
	hash ^= uint64(byte(v >> 40))
	hash *= prime64
	hash ^= uint64(byte(v >> 48))
	hash *= prime64
	hash ^= uint64(byte(v >> 56))
	hash *= prime64
	return hash
}

func FNV64(v uint64) uint64 {
	const (
		offset64 = uint64(14695981039346656037)
		prime64  = uint64(1099511628211)
	)
	hash := offset64
	hash *= prime64
	hash ^= uint64(byte(v))
	hash *= prime64
	hash ^= uint64(byte(v >> 8))
	hash *= prime64
	hash ^= uint64(byte(v >> 16))
	hash *= prime64
	hash ^= uint64(byte(v >> 24))
	hash *= prime64
	hash ^= uint64(byte(v >> 32))
	hash *= prime64
	hash ^= uint64(byte(v >> 40))
	hash *= prime64
	hash ^= uint64(byte(v >> 48))
	hash *= prime64
	hash ^= uint64(byte(v >> 56))
	return hash
}
