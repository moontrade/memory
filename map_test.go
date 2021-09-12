package mem

import "testing"

func TestMap(t *testing.T) {
	a := NewAllocator(1)
	m := NewMap(a, 16)

	key := a.Bytes(0, 16)
	value := a.Bytes(0, 16)

	key.SetString(0, "MYID")
	value.SetString(0, "MYVALUE")

	m.Set(key, value)

	v, ok := m.Get(key)
	if !ok {

	}

	println(key.String(), v.String())
}
