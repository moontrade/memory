package main

import (
	"github.com/moontrade/memory/unsafecgo"
)

func main() {
	//cgo.CGO()
	unsafecgo.NonBlocking((*byte)(nil), 0, 0)
}
