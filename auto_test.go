package memory

import "testing"

func TestAuto_Add(t *testing.T) {
	//sizeBefore := int64(0)
	auto := NewAuto(20)

	for i := 0; i < 100; i++ {
		auto.Alloc(64)
	}

	//if sizeBefore >= a.Stats().AllocSize {
	//	t.Fatal("allocator should have some allocs")
	//}
	//if a.Stats().Allocs != 105 {
	//	t.Fatalf("expected 105 allocs not %d", a.Stats().Allocs)
	//}

	auto.Print()
	auto.Free()

	//if a.Stats().AllocSize != 0 {
	//	t.Fatal("allocator should have zero allocs")
	//}
}

func TestAuto_Scope(t *testing.T) {
	Scope(func(a Auto) {
		for i := 0; i < 100; i++ {
			a.Alloc(64)
			a.Str(128)
		}
	})
	//if a.Stats().Allocs != 0 {
	//	t.Fatalf("allocator should have 0 allocs not: %d", a.Stats().Allocs)
	//}
	//if a.Stats().FreeSize == 0 {
	//	t.Fatal("allocator should have a FreeSize greater than zero")
	//}
	//if a.Stats().AllocSize != 0 {
	//	t.Fatal("allocator should have zero allocs")
	//}
}
