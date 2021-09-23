package tlsf

import (
	"fmt"
	"testing"
	"unsafe"
)

func TestSysAlloc(t *testing.T) {
	var stat sysMemStat
	fmt.Println(unsafe.Sizeof(sysMemStat(0)))
	//ptr := persistentalloc(65536, unsafe.Sizeof(uintptr(0)), &stat)
	//ptr2 := persistentalloc(65536, unsafe.Sizeof(uintptr(0)), &stat)
	ptr3 := sysAlloc(65536, &stat)
	_ = ptr3
	//fmt.Println("ptr", ptr, "ptr2", ptr2, "ptr3", ptr3)
}
