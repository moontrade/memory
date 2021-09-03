package runtime

import (
	"fmt"
	"testing"
)

func Test_New(t *testing.T) {
	fmt.Println(((32 << (^uint(0) >> 63)) / 8 / 4) + 1)
	fmt.Println((32 / 8 / 4) + 1)
	obj := itcmsNew(24)
	fmt.Println(obj)

	mark(obj)
	itcmsCollect()
	mark(obj)
	itcmsCollect()
	itcmsCollect()
	itcmsCollect()
}
