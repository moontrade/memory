package runtime

var tlfsPool = NewPool(1)
var _ = initITCMS()

func visitGlobals() {
	markRoots(globalsStart, globalsEnd)
}

// Visits all objects on the stack.
func visitStack() {}
