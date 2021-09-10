package mem

import "testing"

func TestAuto_Add(t *testing.T) {
	a := NewAllocator(1)

	sizeBefore := a.AllocSize
	auto := NewAuto(a, 20)

	for i := 0; i < 100; i++ {
		auto.Alloc(64)
	}

	if sizeBefore >= a.AllocSize {
		t.Fatal("allocator should have some allocs")
	}
	if a.Allocs != 105 {
		t.Fatalf("expected 105 allocs not %d", a.Allocs)
	}

	auto.Print()
	auto.Free()

	if a.AllocSize != 0 {
		t.Fatal("allocator should have zero allocs")
	}
}

func TestAuto_Scope(t *testing.T) {
	a := NewAllocator(1)
	a.Scope(func(a Auto) {
		for i := 0; i < 100; i++ {
			a.Alloc(64)
		}
	})
	if a.Allocs != 0 {
		t.Fatalf("allocator should have 0 allocs not: %d", a.Allocs)
	}
	if a.FreeSize == 0 {
		t.Fatal("allocator should have a FreeSize greater than zero")
	}
	if a.AllocSize != 0 {
		t.Fatal("allocator should have zero allocs")
	}
}
