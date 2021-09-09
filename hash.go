package mem

import "math/bits"

func toUint16LE(b []byte) uint16 {
	return uint16(b[0]) | uint16(b[1])<<8
}

func toUint32LE(b []byte) uint32 {
	return uint32(b[0]) | uint32(b[1])<<8 | uint32(b[2])<<16 | uint32(b[3])<<24
}

func toUint64LE(b []byte) uint64 {
	return uint64(b[0]) | uint64(b[1])<<8 | uint64(b[2])<<16 | uint64(b[3])<<24 |
		uint64(b[4])<<32 | uint64(b[5])<<40 | uint64(b[6])<<48 | uint64(b[7])<<56
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////
// wyhash
/////////////////////////////////////////////////////////////////////////////////////////////////////////
const (
	wyp0 uint64 = 0xa0761d6478bd642f
	wyp1 uint64 = 0xe7037ed1a0b428db
	wyp2 uint64 = 0x8ebc6af09c88c6e3
	wyp3 uint64 = 0x589965cc75374cc3
	wyp4 uint64 = 0x1d8e4e27c47d124f
	wyp5 uint64 = 0xeb44accab455d165

	wyseed uint64 = 0xeb44accab455d165
	//wyseed uint64 = 0xeb44accab455d16f
	//wyseed uint64 = 0x8ebc6af09c88c6e3
)

func wymum(A, B uint64) uint64 {
	hi, lo := bits.Mul64(A, B)
	return hi ^ lo
}

////go:inline
//func wymum(x, y uint64) uint64 {
//	const mask32 = 1<<32 - 1
//	x0 := x & mask32
//	x1 := x >> 32
//	y0 := y & mask32
//	y1 := y >> 32
//	w0 := x0 * y0
//	t := x1*y0 + w0>>32
//	w1 := t & mask32
//	w2 := t >> 32
//	w1 += x0 * y1
//	return (x1*y1 + w2 + w1>>32) ^ (x * y)
//}

func wyhash32(v uint32) uint32 {
	return uint32(wymum(wymum(uint64(v)^wyseed^wyp0, uint64(v)^wyseed^wyp1)^wyseed, uint64(4)^wyp4))
}

func wyhash64(v uint64) uint64 {
	return wymum(wymum(v^wyseed^wyp0, v^wyseed^wyp1)^wyseed, 8^wyp4)
}

func wyr3(p []byte) uint64 {
	k := len(p)
	return (uint64(p[0]) << 16) | (uint64(p[k>>1]) << 8) | uint64(p[k-1])
}

func wyr8mix(b []byte) uint64 {
	return uint64(uint32(b[0])|uint32(b[1])<<8|uint32(b[2])<<16|uint32(b[3])<<24)<<32 |
		uint64(uint32(b[4])|uint32(b[5])<<8|uint32(b[6])<<16|uint32(b[7])<<24)
}

func wyhash(key []byte, wyseed uint64) uint64 {
	p := key
	if len(p) == 0 {
		return wyseed
	}
	switch {
	case len(p) < 4:
		return wymum(wymum(wyr3(p)^wyseed^wyp0, wyseed^wyp1)^wyseed, uint64(len(p))^wyp4)
	case len(p) <= 8:
		return wymum(wymum(uint64(toUint32LE(p[:4]))^wyseed^wyp0, uint64(toUint32LE(p[len(p)-4:len(p)-4+4]))^wyseed^wyp1)^wyseed, uint64(len(p))^wyp4)
	case len(p) <= 16:
		return wymum(wymum(wyr8mix(p)^wyseed^wyp0, wyr8mix(p[len(p)-8:])^wyseed^wyp1)^wyseed, uint64(len(p))^wyp4)
	case len(p) <= 24:
		return wymum(wymum(wyr8mix(p)^wyseed^wyp0, wyr8mix(p[8:])^wyseed^wyp1)^wymum(wyr8mix(p[len(key)-8:])^wyseed^wyp2, wyseed^wyp3), uint64(len(p))^wyp4)
	case len(p) <= 32:
		return wymum(wymum(wyr8mix(p)^wyseed^wyp0, wyr8mix(p[8:])^wyseed^wyp1)^wymum(wyr8mix(p[16:])^wyseed^wyp2, wyr8mix(p[len(key)-8:])^wyseed^wyp3), uint64(len(p))^wyp4)
	}

	see1 := wyseed
	for len(p) > 256 {
		wyseed = wymum(toUint64LE(p[:8])^wyseed^wyp0, toUint64LE(p[8:8+8])^wyseed^wyp1) ^ wymum(toUint64LE(p[16:16+8])^wyseed^wyp2, toUint64LE(p[24:24+8])^wyseed^wyp3)
		see1 = wymum(toUint64LE(p[32:32+8])^see1^wyp1, toUint64LE(p[40:40+8])^see1^wyp2) ^ wymum(toUint64LE(p[48:48+8])^see1^wyp3, toUint64LE(p[56:56+8])^see1^wyp0)
		wyseed = wymum(toUint64LE(p[64:64+8])^wyseed^wyp0, toUint64LE(p[72:72+8])^wyseed^wyp1) ^ wymum(toUint64LE(p[80:80+8])^wyseed^wyp2, toUint64LE(p[88:88+8])^wyseed^wyp3)
		see1 = wymum(toUint64LE(p[96:96+8])^see1^wyp1, toUint64LE(p[104:104+8])^see1^wyp2) ^ wymum(toUint64LE(p[112:112+8])^see1^wyp3, toUint64LE(p[120:120+8])^see1^wyp0)
		wyseed = wymum(toUint64LE(p[128:128+8])^wyseed^wyp0, toUint64LE(p[136:136+8])^wyseed^wyp1) ^ wymum(toUint64LE(p[144:144+8])^wyseed^wyp2, toUint64LE(p[152:152+8])^wyseed^wyp3)
		see1 = wymum(toUint64LE(p[160:160+8])^see1^wyp1, toUint64LE(p[168:168+8])^see1^wyp2) ^ wymum(toUint64LE(p[176:176+8])^see1^wyp3, toUint64LE(p[184:184+8])^see1^wyp0)
		wyseed = wymum(toUint64LE(p[192:192+8])^wyseed^wyp0, toUint64LE(p[200:200+8])^wyseed^wyp1) ^ wymum(toUint64LE(p[208:208+8])^wyseed^wyp2, toUint64LE(p[216:216+8])^wyseed^wyp3)
		see1 = wymum(toUint64LE(p[224:224+8])^see1^wyp1, toUint64LE(p[232:232+8])^see1^wyp2) ^ wymum(toUint64LE(p[240:240+8])^see1^wyp3, toUint64LE(p[248:248+8])^see1^wyp0)
		p = p[256:]
	}

	for len(p) > 32 {
		wyseed = wymum(toUint64LE(p[:8])^wyseed^wyp0, toUint64LE(p[8:8+8])^wyseed^wyp1)
		see1 = wymum(toUint64LE(p[16:16+8])^see1^wyp2, toUint64LE(p[24:24+8])^see1^wyp3)
		p = p[32:]
	}

	switch {
	case len(p) < 4:
		wyseed = wymum(wyr3(p)^wyseed^wyp0, wyseed^wyp1)
	case len(p) <= 8:
		wyseed = wymum(uint64(toUint32LE(p[:4]))^wyseed^wyp0, uint64(toUint32LE(p[len(p)-4:len(p)-4+4]))^wyseed^wyp1)
	case len(p) <= 16:
		wyseed = wymum(wyr8mix(p)^wyseed^wyp0, wyr8mix(p[len(p)-8:])^wyseed^wyp1)
	case len(p) <= 24:
		wyseed = wymum(wyr8mix(p)^wyseed^wyp0, wyr8mix(p[8:])^wyseed^wyp1)
		see1 = wymum(wyr8mix(p[len(p)-8:])^see1^wyp2, see1^wyp3)
	default:
		wyseed = wymum(wyr8mix(p)^wyseed^wyp0, wyr8mix(p[8:])^wyseed^wyp1)
		see1 = wymum(wyr8mix(p[16:])^see1^wyp2, wyr8mix(p[len(p)-8:])^see1^wyp3)

	}
	return wymum(wyseed^see1, uint64(len(key))^wyp4)
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////
// adler32
/////////////////////////////////////////////////////////////////////////////////////////////////////////

func adler32(v uint32) uint32 {
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

/////////////////////////////////////////////////////////////////////////////////////////////////////////
// fnv
/////////////////////////////////////////////////////////////////////////////////////////////////////////

func fnv32(v uint32) uint32 {
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

/////////////////////////////////////////////////////////////////////////////////////////////////////////
// metro
/////////////////////////////////////////////////////////////////////////////////////////////////////////

func metro32(v uint32) uint32 {
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

func metro64(v uint64) uint64 {
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

func metro(buffer []byte, seed uint64) uint64 {
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
