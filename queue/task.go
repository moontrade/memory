package queue

import "github.com/moontrade/memory"

type Task struct {
	Data memory.Bytes
	next uintptr
}

//func AllocTask(a mem.Allocator, size uintptr) *Task {
//	a.AllocZeroed(mem.Pointer(size + unsafe.Sizeof(Task{})))
//}
