package rhmap

import (
	mem "github.com/moontrade/memory"
	"testing"
)

func TestMap(t *testing.T) {
	m := NewMap(mem.NextAllocator(), 16)

	key := mem.AllocString(16)
	value := mem.AllocString(16)

	key.AppendString("MYID")
	value.AppendString("MYVALUE")

	m.Set(key, value)

	v, ok := m.Get(key)
	if !ok {

	}

	println(key.String(), v.String())
}
