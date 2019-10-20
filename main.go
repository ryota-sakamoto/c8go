package main

import (
	"fmt"
	"os"

	"github.com/ryota-sakamoto/c8go/code"
	"github.com/ryota-sakamoto/c8go/node"
	"github.com/ryota-sakamoto/c8go/token"
)

func main() {
	if len(os.Args) < 2 {
		panic("Invalud args len")
	}

	s := os.Args[1]
	token, err := token.Tokenize(s)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	parser := node.NewNodeParser(token)
	node, err := parser.Expr()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	generator := code.NewGenerator(node)
	generator.Run()
}
