package mem

import (
	"testing"
)

func TestPointer(t *testing.T) {
	a := NewAllocator(1)
	b := a.BytesMut(128, 252)
	b.SetUInt32(0, 5)
	println(b.UInt32(0))

	b.SetUInt32(7, 13)
	println(b.UInt32(7))

	str := "hello there"
	b.SetString(24, str)
	println(b.Substring(24, len(str)))
}
