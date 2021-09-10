package mem

import "testing"

func TestAuto_Add(t *testing.T) {
	a := NewAllocator(1)

	sizeBefore := a.AllocSize
	auto := NewAuto(a, 20)

	for i := 0; i < 99; i++ {
		auto.Alloc(64)
	}

	if sizeBefore >= a.AllocSize {
		t.Fatal("allocator should have some allocs")
	}
	if a.Allocs != 104 {
		t.Fatalf("expected 105 allocs not %d", a.Allocs)
	}

	auto.Print()
	auto.Free()

	if a.AllocSize != 0 {
		t.Fatal("allocator should have zero allocs")
	}
}
