package main

import (
	"fmt"
	mem "github.com/moontrade/memory"
	"github.com/moontrade/memory/art"
)

var alloc = mem.NewTLSF(1)

func keyOf(v string) art.Key {
	b := alloc.Bytes(mem.Pointer(len(v)))
	b.AppendString(v)
	return art.Key{b}
}

func keyOfBytes(v []byte) art.Key {
	b := alloc.Bytes(mem.Pointer(len(v)))
	b.AppendBytes(v)
	return art.Key{b}
}

func DumpTree() {
	tree := art.New()
	terms := []string{"A", "a", "aa"}
	for _, term := range terms {
		tree.Insert(keyOf(term), term)
	}
	fmt.Println(tree)
}

func SimpleTree() {
	tree := art.New()

	tree.Insert(keyOf("Hi, I'm Key"), "Nice to meet you, I'm Value")
	value, found := tree.Search(keyOf("Hi, I'm Key"))
	if found {
		fmt.Printf("Search value=%v\n", value)
	}

	tree.ForEach(func(node art.Node) bool {
		fmt.Printf("Callback value=%v\n", node.Value())
		return true
	})

	for it := tree.Iterator(); it.HasNext(); {
		value, _ := it.Next()
		fmt.Printf("Iterator value=%v\n", value.Value())
	}
}

func main() {
	DumpTree()
	SimpleTree()
}
