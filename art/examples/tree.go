package main

import (
	"fmt"
	"github.com/moontrade/memory/art"
)

func DumpTree() {
	tree := art.New()
	terms := []string{"A", "a", "aa"}
	for _, term := range terms {
		tree.Insert(art.StringKey(term), term)
	}
	fmt.Println(tree)
}

func SimpleTree() {
	tree := art.New()

	tree.Insert(art.StringKey("Hi, I'm Key"), "Nice to meet you, I'm Value")
	value, found := tree.Search(art.StringKey("Hi, I'm Key"))
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
