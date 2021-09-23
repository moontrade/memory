package hash

import "math/bits"

func Metro32(v uint32) uint32 {
	const (
		k0   uint64 = 0xD6D018F5
		k1   uint64 = 0xA2AA033B
		k2   uint64 = 0x62992FC1
		k3   uint64 = 0x30BC5B29
		seed uint64 = 979532
	)
	hash := (seed + k2) * k0
	hash += uint64(v) * k3
	hash ^= bits.RotateLeft64(hash, -26) * k1
	hash ^= bits.RotateLeft64(hash, -28)
	hash *= k0
	hash ^= bits.RotateLeft64(hash, -29)
	return uint32(hash)
}

func Metro64(v uint64) uint64 {
	const (
		k0   uint64 = 0xD6D018F5
		k1   uint64 = 0xA2AA033B
		k2   uint64 = 0x62992FC1
		k3   uint64 = 0x30BC5B29
		seed uint64 = 979532
	)
	hash := (seed + k2) * k0
	hash += v * k3
	hash ^= bits.RotateLeft64(hash, -26) * k1
	hash ^= bits.RotateLeft64(hash, -28)
	hash *= k0
	hash ^= bits.RotateLeft64(hash, -29)
	return hash
}

func Metro(buffer []byte, seed uint64) uint64 {
	const (
		k0 = 0xD6D018F5
		k1 = 0xA2AA033B
		k2 = 0x62992FC1
		k3 = 0x30BC5B29
	)

	var (
		ptr  = buffer
		hash = (seed + k2) * k0
	)
	if len(ptr) >= 32 {
		v0, v1, v2, v3 := hash, hash, hash, hash

		for len(ptr) >= 32 {
			v0 += toUint64LE(ptr[:8]) * k0
			v0 = bits.RotateLeft64(v0, -29) + v2
			v1 += toUint64LE(ptr[8:16]) * k1
			v1 = bits.RotateLeft64(v1, -29) + v3
			v2 += toUint64LE(ptr[16:24]) * k2
			v2 = bits.RotateLeft64(v2, -29) + v0
			v3 += toUint64LE(ptr[24:32]) * k3
			v3 = bits.RotateLeft64(v3, -29) + v1
			ptr = ptr[32:]
		}

		v2 ^= bits.RotateLeft64(((v0+v3)*k0)+v1, -37) * k1
		v3 ^= bits.RotateLeft64(((v1+v2)*k1)+v0, -37) * k0
		v0 ^= bits.RotateLeft64(((v0+v2)*k0)+v3, -37) * k1
		v1 ^= bits.RotateLeft64(((v1+v3)*k1)+v2, -37) * k0
		hash += v0 ^ v1
	}
	if len(ptr) >= 16 {
		v0 := hash + (toUint64LE(ptr[:8]) * k2)
		v0 = bits.RotateLeft64(v0, -29) * k3
		v1 := hash + (toUint64LE(ptr[8:16]) * k2)
		v1 = bits.RotateLeft64(v1, -29) * k3
		v0 ^= bits.RotateLeft64(v0*k0, -21) + v1
		v1 ^= bits.RotateLeft64(v1*k3, -21) + v0
		hash += v1
		ptr = ptr[16:]
	}
	if len(ptr) >= 8 {
		hash += toUint64LE(ptr[:8]) * k3
		ptr = ptr[8:]
		hash ^= bits.RotateLeft64(hash, -55) * k1
	}
	if len(ptr) >= 4 {
		hash += uint64(toUint32LE(ptr[:4])) * k3
		hash ^= bits.RotateLeft64(hash, -26) * k1
		ptr = ptr[4:]
	}
	if len(ptr) >= 2 {
		hash += uint64(toUint16LE(ptr[:2])) * k3
		ptr = ptr[2:]
		hash ^= bits.RotateLeft64(hash, -48) * k1
	}
	if len(ptr) >= 1 {
		hash += uint64(ptr[0]) * k3
		hash ^= bits.RotateLeft64(hash, -37) * k1
	}

	hash ^= bits.RotateLeft64(hash, -28)
	hash *= k0
	hash ^= bits.RotateLeft64(hash, -29)

	return hash
}
