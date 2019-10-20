package code

import (
	"fmt"

	"github.com/ryota-sakamoto/c8go/node"
)

type Generator struct {
	node *node.Node
}

func NewGenerator(n *node.Node) *Generator {
	return &Generator{
		node: n,
	}
}

func (g *Generator) Run() {
	fmt.Println(".intel_syntax noprefix")
	fmt.Println(".global main")
	fmt.Println("main:")

	gen(g.node)

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
	case node.ND_EQ:
		fmt.Println("    cmp rax, rdi")
		fmt.Println("    sete al")
		fmt.Println("    movzb rax, al")
	case node.ND_NE:
		fmt.Println("    cmp rax, rdi")
		fmt.Println("    setne al")
		fmt.Println("    movzb rax, al")
	case node.ND_LT:
		fmt.Println("    cmp rax, rdi")
		fmt.Println("    setl al")
		fmt.Println("    movzb rax, al")
	case node.ND_LE:
		fmt.Println("    cmp rax, rdi")
		fmt.Println("    setle al")
		fmt.Println("    movzb rax, al")
	}

	fmt.Println("    push rax")
}
