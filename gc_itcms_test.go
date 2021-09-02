package runtime

import (
	"fmt"
	"testing"
)

func Test_New(t *testing.T) {
	obj := itcmsNew(24)
	fmt.Println(obj)

	mark(obj)
	itcmsCollect()
	mark(obj)
	itcmsCollect()
	itcmsCollect()
	itcmsCollect()
}
