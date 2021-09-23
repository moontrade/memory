package hash

func Adler32(v uint32) uint32 {
	const mod = 65521
	var d uint32 = 1
	s1, s2 := d&0xffff, d>>16

	s1 += uint32(byte(v))
	s2 += s1
	s1 += uint32(byte(v >> 8))
	s2 += s1
	s1 += uint32(byte(v >> 16))
	s2 += s1
	s1 += uint32(byte(v >> 24))
	s2 += s1
	s1 %= mod
	s2 %= mod
	return s2<<16 | s1
}
