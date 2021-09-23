package alloc

//func TestPointerBE(t *testing.T) {
//	a := NewTLSF(1)
//	p := Pointer(a.Alloc(32))
//	p.SetInt32BE(0, 100)
//	println(p.Int32BE(0), p.Int32BESlow(), p.Int32(0))
//}
//
//func BenchmarkPointerBE(b *testing.B) {
//	a := NewTLSF(1)
//	p := Pointer(a.Alloc(32))
//
//	b.Run("Int32", func(b *testing.B) {
//		for i := 0; i < b.N; i++ {
//			p.Int32(0)
//		}
//	})
//	b.Run("Int32BE", func(b *testing.B) {
//		for i := 0; i < b.N; i++ {
//			p.Int32BE(0)
//		}
//	})
//	b.Run("Int32BE Slow", func(b *testing.B) {
//		for i := 0; i < b.N; i++ {
//			p.Int32BESlow()
//		}
//	})
//	b.Run("UInt32BE", func(b *testing.B) {
//		for i := 0; i < b.N; i++ {
//			p.UInt32BE(0)
//		}
//	})
//	b.Run("UInt32BE Slow", func(b *testing.B) {
//		for i := 0; i < b.N; i++ {
//			p.UInt32BESlow()
//		}
//	})
//	b.Run("Int64", func(b *testing.B) {
//		for i := 0; i < b.N; i++ {
//			p.Int64(0)
//		}
//	})
//	b.Run("Int64BE", func(b *testing.B) {
//		for i := 0; i < b.N; i++ {
//			p.Int64BE(0)
//		}
//	})
//	b.Run("UInt64BE", func(b *testing.B) {
//		for i := 0; i < b.N; i++ {
//			p.UInt64BE(0)
//		}
//	})
//	b.Run("UInt64BE Slow", func(b *testing.B) {
//		for i := 0; i < b.N; i++ {
//			p.UInt64BESlow()
//		}
//	})
//}
//
////goland:noinspection GoVetUnsafePointer
//func (p Pointer) Int16BESlow() int16 {
//	return int16(*(*byte)(unsafe.Pointer(p + 1))) |
//		int16(*(*byte)(unsafe.Pointer(p)))<<8
//}
//
////goland:noinspection GoVetUnsafePointer
//func (p Pointer) UInt24BESlow() uint32 {
//	return uint32(*(*byte)(unsafe.Pointer(p + 2))) |
//		uint32(*(*byte)(unsafe.Pointer(p + 1)))<<8 |
//		uint32(*(*byte)(unsafe.Pointer(p)))<<16
//}
//
////goland:noinspection GoVetUnsafePointer
//func (p Pointer) UInt32LESlow() uint32 {
//	return uint32(*(*byte)(unsafe.Pointer(p + 3))) |
//		uint32(*(*byte)(unsafe.Pointer(p + 2)))<<8 |
//		uint32(*(*byte)(unsafe.Pointer(p + 1)))<<16 |
//		uint32(*(*byte)(unsafe.Pointer(p)))<<24
//}
//
////goland:noinspection GoVetUnsafePointer
//func (p Pointer) UInt64BESlow() uint64 {
//	return uint64(*(*byte)(unsafe.Pointer(p + 7))) |
//		uint64(*(*byte)(unsafe.Pointer(p + 6)))<<8 |
//		uint64(*(*byte)(unsafe.Pointer(p + 5)))<<16 |
//		uint64(*(*byte)(unsafe.Pointer(p + 4)))<<24 |
//		uint64(*(*byte)(unsafe.Pointer(p + 3)))<<32 |
//		uint64(*(*byte)(unsafe.Pointer(p + 2)))<<40 |
//		uint64(*(*byte)(unsafe.Pointer(p + 1)))<<48 |
//		uint64(*(*byte)(unsafe.Pointer(p)))<<56
//}
//
////goland:noinspection GoVetUnsafePointer
//func (p Pointer) Int32BESlow() int32 {
//	return int32(*(*byte)(unsafe.Pointer(p + 3))) |
//		int32(*(*byte)(unsafe.Pointer(p + 2)))<<8 |
//		int32(*(*byte)(unsafe.Pointer(p + 1)))<<16 |
//		int32(*(*byte)(unsafe.Pointer(p)))<<24
//}
//
//func BenchmarkPointerBESet(b *testing.B) {
//	a := NewTLSF(1)
//	p := Pointer(a.Alloc(32))
//
//	b.Run("SetInt32At", func(b *testing.B) {
//		for i := 0; i < b.N; i++ {
//			p.SetInt32BE(8, 1)
//		}
//	})
//	b.Run("SetInt32BEAtSlow", func(b *testing.B) {
//		for i := 0; i < b.N; i++ {
//			//p.SetInt32BEAtSlow(8, 1)
//		}
//	})
//
//}
