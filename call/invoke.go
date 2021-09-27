package call

//go:noescape
func InvokeC(fn *byte, arg0, arg1 uintptr)
