package main

import (
	"fmt"
	"github.com/moontrade/memory/artgo"
)

func DumpTree() {
	tree := artgo.New()
	terms := []string{"A", "a", "aa"}
	for _, term := range terms {
		tree.Insert(artgo.StringKey(term), term)
	}
	fmt.Println(tree)
}

func SimpleTree() {
	tree := artgo.New()

	tree.Insert(artgo.StringKey("Hi, I'm Key"), "Nice to meet you, I'm Value")
	value, found := tree.Search(artgo.StringKey("Hi, I'm Key"))
	if found {
		fmt.Printf("Search value=%v\n", value)
	}

	tree.ForEach(func(node artgo.Node) bool {
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
