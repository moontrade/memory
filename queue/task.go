package queue

import (
	mem "github.com/moontrade/memory"
)

type Task struct {
	Data mem.Bytes
	next uintptr
}

//func AllocTask(a mem.Allocator, size uintptr) *Task {
//	a.AllocNotCleared(mem.Pointer(size + unsafe.Sizeof(Task{})))
//}
