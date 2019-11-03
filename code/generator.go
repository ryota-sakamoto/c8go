package code

import (
	"fmt"

	"github.com/ryota-sakamoto/c8go/node"
)

type Generator struct {
	node *node.Node
}

func NewGenerator() *Generator {
	return &Generator{}
}

func (g *Generator) Before() {
	fmt.Println(".intel_syntax noprefix")
	fmt.Println(".global main")
	fmt.Println("main:")

	fmt.Println("    push rbp")
	fmt.Println("    mov rbp, rsp")
	fmt.Println("    sub rsp, 208")
}

func (g *Generator) After() {
	fmt.Println("    mov rsp, rbp")
	fmt.Println("    pop rbp")
	fmt.Println("    ret")
}

func (g *Generator) Run(n *node.Node) {
	gen(n)
	fmt.Println("    pop rax")
}

func gen(n *node.Node) {
	switch n.Kind {
	case node.ND_NUM:
		fmt.Println(fmt.Sprintf("    push %d", n.Val))
		return
	case node.ND_LVAR:
		genLabel(n)
		fmt.Println("    pop rax")
		fmt.Println("    mov rax, [rax]")
		fmt.Println("    push rax")
		return
	case node.ND_ASSIGN:
		genLabel(n.Left)
		gen(n.Right)

		fmt.Println("    pop rdi")
		fmt.Println("    pop rax")
		fmt.Println("    mov [rax], rdi")
		fmt.Println("    push rdi")
		return
	case node.ND_RETURN:
		gen(n.Right)

		fmt.Println("    pop rax")
		fmt.Println("    mov rsp, rbp")
		fmt.Println("    pop rbp")
		fmt.Println("    ret")
		return
	case node.ND_IF:
		gen(n.Left)

		fmt.Println("    pop rax")
		fmt.Println("    cmp rax, 0")
		fmt.Println("    je .LendXXX")

		gen(n.Right)

		fmt.Println(".LendXXX:")
		return
	case node.ND_IF_ELSE:
		gen(n.Left)
		gen(n.Right)
		return
	case node.ND_ELSE:
		fmt.Println("    pop rax")
		fmt.Println("    cmp rax, 0")
		fmt.Println("    je .LelseXXX")

		gen(n.Left)
		fmt.Println("    jmp .LendXXX")
		fmt.Println(".LelseXXX:")
		gen(n.Right)
		fmt.Println(".LendXXX:")

		return
	case node.ND_WHILE:
		fmt.Println(".LbeginXXX:")
		gen(n.Left)

		fmt.Println("    pop rax")
		fmt.Println("    cmp rax, 0")
		fmt.Println("    je .LendXXX")

		gen(n.Right)

		fmt.Println("    jmp .LbeginXXX")
		fmt.Println(".LendXXX:")
		return
	case node.ND_BLOCK:
		for _, n := range n.Block {
			gen(n)
		}
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

func genLabel(n *node.Node) {
	if n.Kind != node.ND_LVAR {
		panic("")
	}
	fmt.Println("    mov rax, rbp")
	fmt.Println(fmt.Sprintf("    sub rax, %d", n.Offset))
	fmt.Println("    push rax")
}
