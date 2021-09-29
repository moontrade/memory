package art

import (
	"testing"
)

func TestNew(t *testing.T) {
	loop := New(1000)
	if err := loop.BindAcceptor(20000); err != nil {
		t.Fatal(err)
	}
	loop.Destroy()
}
