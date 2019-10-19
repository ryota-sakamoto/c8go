package main

import (
	"fmt"
	"os"

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

	fmt.Println(".intel_syntax noprefix")
	fmt.Println(".global main")
	fmt.Println("main:")

	gen(node)
	fmt.Println("    pop rax")
	fmt.Println("    ret")
}

func gen(n *node.Node) {
	if n.IsNum() {
		fmt.Println(fmt.Sprintf("    push %d", n.Val))
		return
	}

	gen(n.Left)
	gen(n.Right)

	fmt.Println("    pop rdi")
	fmt.Println("    pop rax")

	switch n.Kind {
	case node.ND_ADD:
		fmt.Println("    add rax, rdi")
	case node.ND_SUB:
		fmt.Println("    sub rax, rdi")
	case node.ND_MUL:
		fmt.Println("    imul rax, rdi")
	case node.ND_DIV:
		fmt.Println("    cqo")
		fmt.Println("    idiv rdi")
	}

	fmt.Println("    push rax")
}
