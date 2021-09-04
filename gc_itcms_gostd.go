package runtime

var tlfsPool = newTLSF(1)
var _ = initITCMS()

func visitGlobals() {
	markRoots(globalsStart, globalsEnd)
}

// Visits all objects on the stack.
func visitStack() {}
