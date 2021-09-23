package hash

import (
	"encoding/binary"
	"hash/crc32"
	"hash/crc64"
	"math/rand"
	"testing"
)

func TestHashCollisions_WASM(t *testing.T) {
	var (
		c32   int
		a32   int
		c16   int
		fn    int
		mur   int
		wy    int
		wy64  int
		met   int
		met64 int
		total int
	)

	type config struct {
		low            int
		high           int
		adder          int
		factor         float64
		addressStart   int
		multiplierLow  int
		multiplierHigh int
		multiplierAdd  int
	}
	for _, cfg := range []config{
		//{2, 10, 4, 8, 65580, 64, 4096, 128},
		// WASM like pointer values
		{40, 1024, 88, 3, 65580, 512, 256000, 96},
		{1024, 4096, 512, 3, 65580, 56, 52050, 512},
	} {
		for i := cfg.low; i < cfg.high; i += cfg.adder {
			var (
				entries = i
				slots   = int(float64(i) * cfg.factor)
			)
			for multiplier := cfg.multiplierLow; multiplier < cfg.multiplierHigh; multiplier += cfg.multiplierAdd {
				total += entries
				c32 += testCollisions(entries, multiplier, slots, crc32h)
				c16 += testCollisions(entries, multiplier, slots, crc16a)
				a32 += testCollisions(entries, multiplier, slots, Adler32)
				fn += testCollisions(entries, multiplier, slots, FNV32)
				mur += testCollisions(entries, multiplier, slots, Murmur32)
				wy += testCollisions(entries, multiplier, slots, WyHash32)
				wy64 += testCollisions64(entries, multiplier, slots, WyHash64)
				met += testCollisions(entries, multiplier, slots, Metro32)
				met64 += testCollisions64(entries, multiplier, slots, Metro64)
			}
		}

		//println("\tcrc32		", c32)
		//println("\tcrc16		", c16)
		//println("\tadler32		", a32)
		//println("\tfnv			", fn)
		//println("\tmurmur32		", mur)
		//println("\twyhash32	", wy)
		//println("\twyhash64	", wy64)
		//println("\tmetro32		", met)
		//println("\tmetro64		", met64)
	}

	println("")
	println("total		", total)
	println("\tcrc32		", c32)
	println("\tcrc16		", c16)
	println("\tAdler32		", a32)
	println("\tfnv			", fn)
	println("\tMurmur32	", mur)
	println("\tWyHash32	", wy)
	println("\tWyHash64	", wy64)
	println("\tMetro32		", met)
	println("\tMetro64		", met64)
}

func rangeRandom(min, max uint32) uint32 {
	return uint32(rand.Int31n(int32(max-min)) + int32(min))
}

func testCollisions(entries, multiplier, slots int, hasher func(uint32) uint32) int {
	m := make(map[uint32]struct{})
	count := 0

	ptr := 65680

	for i := 0; i < entries; i++ {
		v := hasher(uint32(ptr))
		ptr += multiplier
		index := v % uint32(slots)

		_, ok := m[index]
		if ok {
			count++
		} else {
			m[index] = struct{}{}
		}
	}
	return count
}

func testCollisionsRandom(entries, slots int, hasher func(uint32) uint32) int {
	m := make(map[uint32]struct{})
	count := 0

	for i := 0; i < entries; i++ {
		v := hasher(uint32(rand.Uint32()))
		index := v % uint32(slots)

		_, ok := m[index]
		if ok {
			count++
		} else {
			m[index] = struct{}{}
		}
	}
	return count
}

func testCollisions64(entries, multiplier, slots int, hasher func(uint64) uint64) int {
	m := make(map[uint64]struct{})
	count := 0

	ptr := 65536

	for i := uint64(0); i < uint64(entries); i++ {
		v := hasher(uint64(ptr))
		ptr += multiplier
		index := v % uint64(slots)

		_, ok := m[index]
		if ok {
			count++
		} else {
			m[index] = struct{}{}
		}
	}
	return count
}

func crc64h(v uint32) uint32 {
	h := crc64.New(crc64.MakeTable(crc64.ECMA))
	h.Reset()
	b := [4]byte{}
	binary.LittleEndian.PutUint32(b[0:], v)
	h.Write(b[0:4])
	return uint32(h.Sum64())
}

func crc32h(v uint32) uint32 {
	h := crc32.NewIEEE()
	h.Reset()
	b := [4]byte{}
	binary.LittleEndian.PutUint32(b[0:], v)
	h.Write(b[0:4])
	return uint32(h.Sum32())
}

var crc16tab = [256]uint16{
	0x0000, 0x1021, 0x2042, 0x3063, 0x4084, 0x50a5, 0x60c6, 0x70e7,
	0x8108, 0x9129, 0xa14a, 0xb16b, 0xc18c, 0xd1ad, 0xe1ce, 0xf1ef,
	0x1231, 0x0210, 0x3273, 0x2252, 0x52b5, 0x4294, 0x72f7, 0x62d6,
	0x9339, 0x8318, 0xb37b, 0xa35a, 0xd3bd, 0xc39c, 0xf3ff, 0xe3de,
	0x2462, 0x3443, 0x0420, 0x1401, 0x64e6, 0x74c7, 0x44a4, 0x5485,
	0xa56a, 0xb54b, 0x8528, 0x9509, 0xe5ee, 0xf5cf, 0xc5ac, 0xd58d,
	0x3653, 0x2672, 0x1611, 0x0630, 0x76d7, 0x66f6, 0x5695, 0x46b4,
	0xb75b, 0xa77a, 0x9719, 0x8738, 0xf7df, 0xe7fe, 0xd79d, 0xc7bc,
	0x48c4, 0x58e5, 0x6886, 0x78a7, 0x0840, 0x1861, 0x2802, 0x3823,
	0xc9cc, 0xd9ed, 0xe98e, 0xf9af, 0x8948, 0x9969, 0xa90a, 0xb92b,
	0x5af5, 0x4ad4, 0x7ab7, 0x6a96, 0x1a71, 0x0a50, 0x3a33, 0x2a12,
	0xdbfd, 0xcbdc, 0xfbbf, 0xeb9e, 0x9b79, 0x8b58, 0xbb3b, 0xab1a,
	0x6ca6, 0x7c87, 0x4ce4, 0x5cc5, 0x2c22, 0x3c03, 0x0c60, 0x1c41,
	0xedae, 0xfd8f, 0xcdec, 0xddcd, 0xad2a, 0xbd0b, 0x8d68, 0x9d49,
	0x7e97, 0x6eb6, 0x5ed5, 0x4ef4, 0x3e13, 0x2e32, 0x1e51, 0x0e70,
	0xff9f, 0xefbe, 0xdfdd, 0xcffc, 0xbf1b, 0xaf3a, 0x9f59, 0x8f78,
	0x9188, 0x81a9, 0xb1ca, 0xa1eb, 0xd10c, 0xc12d, 0xf14e, 0xe16f,
	0x1080, 0x00a1, 0x30c2, 0x20e3, 0x5004, 0x4025, 0x7046, 0x6067,
	0x83b9, 0x9398, 0xa3fb, 0xb3da, 0xc33d, 0xd31c, 0xe37f, 0xf35e,
	0x02b1, 0x1290, 0x22f3, 0x32d2, 0x4235, 0x5214, 0x6277, 0x7256,
	0xb5ea, 0xa5cb, 0x95a8, 0x8589, 0xf56e, 0xe54f, 0xd52c, 0xc50d,
	0x34e2, 0x24c3, 0x14a0, 0x0481, 0x7466, 0x6447, 0x5424, 0x4405,
	0xa7db, 0xb7fa, 0x8799, 0x97b8, 0xe75f, 0xf77e, 0xc71d, 0xd73c,
	0x26d3, 0x36f2, 0x0691, 0x16b0, 0x6657, 0x7676, 0x4615, 0x5634,
	0xd94c, 0xc96d, 0xf90e, 0xe92f, 0x99c8, 0x89e9, 0xb98a, 0xa9ab,
	0x5844, 0x4865, 0x7806, 0x6827, 0x18c0, 0x08e1, 0x3882, 0x28a3,
	0xcb7d, 0xdb5c, 0xeb3f, 0xfb1e, 0x8bf9, 0x9bd8, 0xabbb, 0xbb9a,
	0x4a75, 0x5a54, 0x6a37, 0x7a16, 0x0af1, 0x1ad0, 0x2ab3, 0x3a92,
	0xfd2e, 0xed0f, 0xdd6c, 0xcd4d, 0xbdaa, 0xad8b, 0x9de8, 0x8dc9,
	0x7c26, 0x6c07, 0x5c64, 0x4c45, 0x3ca2, 0x2c83, 0x1ce0, 0x0cc1,
	0xef1f, 0xff3e, 0xcf5d, 0xdf7c, 0xaf9b, 0xbfba, 0x8fd9, 0x9ff8,
	0x6e17, 0x7e36, 0x4e55, 0x5e74, 0x2e93, 0x3eb2, 0x0ed1, 0x1ef0}

func crc16a(v uint32) uint32 {
	crc := uint16(0)
	crc = ((crc << 8) & 0xff00) ^ crc16tab[((crc>>8)&0xff)^uint16(byte(v))]
	crc = ((crc << 8) & 0xff00) ^ crc16tab[((crc>>8)&0xff)^uint16(byte(v<<8))]
	crc = ((crc << 8) & 0xff00) ^ crc16tab[((crc>>8)&0xff)^uint16(byte(v<<16))]
	crc = ((crc << 8) & 0xff00) ^ crc16tab[((crc>>8)&0xff)^uint16(byte(v<<24))]
	return uint32(crc)
}

func BenchmarkHash(b *testing.B) {
	//b.Run("crc32", func(b *testing.B) {
	//	for i := 0; i < b.N; i++ {
	//		crc32h(23)
	//	}
	//})
	//b.Run("crc16", func(b *testing.B) {
	//	for i := 0; i < b.N; i++ {
	//		crc16a(23)
	//	}
	//})
	//b.Run("crc16", func(b *testing.B) {
	//	for i := 0; i < b.N; i++ {
	//		crc16a(23)
	//	}
	//})
	b.Run("Adler32", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			Adler32(Adler32(uint32(i)))
		}
	})
	b.Run("FNV32", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			FNV32(FNV32(uint32(i)))
		}
	})
	b.Run("Murmur32", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			Murmur32(Murmur32(uint32(i)))
		}
	})
	b.Run("WyHash32", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			WyHash32(WyHash32(uint32(i)))
		}
	})
	b.Run("WyHash64", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			WyHash64(WyHash64(uint64(i)))
		}
	})
	b.Run("Metro32", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			Metro32(Metro32(uint32(i)))
		}
	})
	b.Run("Metro64", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			Metro64(Metro64(uint64(i)))
		}
	})
	//b.Run("WyHash", func(b *testing.B) {
	//	h := fnv.New32a()
	//	for i := 0; i < b.N; i++ {
	//		h.Reset()
	//
	//		WyHash32(23)
	//	}
	//})
}
