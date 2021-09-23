package hash

import "math/bits"

/////////////////////////////////////////////////////////////////////////////////////////////////////////
// WyHash
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

func WyHash32(v uint32) uint32 {
	return uint32(wymum(wymum(uint64(v)^wyseed^wyp0, uint64(v)^wyseed^wyp1)^wyseed, uint64(4)^wyp4))
}

func WyHash64(v uint64) uint64 {
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

func WyHash(key []byte, wyseed uint64) uint64 {
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
