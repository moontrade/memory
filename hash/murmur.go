package hash

import "math/bits"

func Murmur32(v uint32) uint32 {
	const (
		c1   uint32 = 0xcc9e2d51
		c2   uint32 = 0x1b873593
		seed uint32 = 0x1f576b93
	)
	h1 := seed
	k1 := v

	k1 *= c1
	k1 = bits.RotateLeft32(k1, 15)
	k1 *= c2

	h1 ^= k1
	h1 = bits.RotateLeft32(h1, 13)
	h1 = h1*5 + 0xe6546b64

	h1 ^= 4

	h1 ^= h1 >> 16
	h1 *= 0x85ebca6b
	h1 ^= h1 >> 13
	h1 *= 0xc2b2ae35
	h1 ^= h1 >> 16

	return h1
}
