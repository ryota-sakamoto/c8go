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
	case node.ND_FUNC:
		fmt.Println(n.Name + ":")

		fmt.Println("    push rbp")
		fmt.Println("    mov rbp, rsp")
		fmt.Println("    sub rsp, 208")

		for i, offset := range n.DefineArgsOffset {
			fmt.Println("    mov rax, rbp")
			fmt.Println(fmt.Sprintf("    sub rax, %d", offset))

			switch i {
			case 0:
				fmt.Println("    mov [rax], rdi")
			case 1:
				fmt.Println("    mov [rax], rsi")
			case 2:
				fmt.Println("    mov [rax], rdx")
			case 3:
				fmt.Println("    mov [rax], rcx")
			case 4:
				fmt.Println("    mov [rax], r8")
			case 5:
				fmt.Println("    mov [rax], r9")
			default:
				panic(fmt.Sprintf("not support args len: %d", len(n.DefineArgsOffset)))
			}
		}

		for _, n := range n.Block {
			gen(n)
		}
		return
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
		end := getLabelCount()

		gen(n.Left)

		fmt.Println("    pop rax")
		fmt.Println("    cmp rax, 0")
		fmt.Println(fmt.Sprintf("    je .Lend%d", end))

		gen(n.Right)

		fmt.Println(fmt.Sprintf(".Lend%d:", end))
		return
	case node.ND_IF_ELSE:
		gen(n.Left)
		gen(n.Right)
		return
	case node.ND_ELSE:
		ec := getLabelCount()
		end := getLabelCount()

		fmt.Println("    pop rax")
		fmt.Println("    cmp rax, 0")
		fmt.Println(fmt.Sprintf("    je .Lelse%d", ec))

		gen(n.Left)
		fmt.Println(fmt.Sprintf("    jmp .Lend%d", end))
		fmt.Println(fmt.Sprintf(".Lelse%d:", ec))
		gen(n.Right)
		fmt.Println(fmt.Sprintf(".Lend%d:", end))

		return
	case node.ND_WHILE:
		begin := getLabelCount()
		end := getLabelCount()

		fmt.Println(fmt.Sprintf(".Lbegin%d:", begin))
		gen(n.Left)

		fmt.Println("    pop rax")
		fmt.Println("    cmp rax, 0")
		fmt.Println(fmt.Sprintf("    je .Lend%d", end))

		gen(n.Right)

		fmt.Println(fmt.Sprintf("    jmp .Lbegin%d", begin))
		fmt.Println(fmt.Sprintf(".Lend%d:", end))
		return
	case node.ND_BLOCK:
		for _, n := range n.Block {
			gen(n)
		}
		return
	case node.ND_CALL_FUNC:
		if len(n.Args) > 6 {
			panic(fmt.Sprintf("not support args len: %d", len(n.Args)))
		}

		for i, argsNode := range n.Args {
			gen(argsNode)

			fmt.Println("    pop rax")

			switch i {
			case 0:
				fmt.Println("    mov rdi, rax")
			case 1:
				fmt.Println("    mov rsi, rax")
			case 2:
				fmt.Println("    mov rdx, rax")
			case 3:
				fmt.Println("    mov rcx, rax")
			case 4:
				fmt.Println("    mov r8, rax")
			case 5:
				fmt.Println("    mov r9, rax")
			}
		}

		fmt.Println("    push rbp")
		fmt.Println("    mov rbp, rsp")
		fmt.Println(fmt.Sprintf("    call %s", n.Name))
		fmt.Println("    pop rbp")

		fmt.Println("    push rax")
		return
	case node.ND_ADDR:
		genLabel(n.Left)
		return
	case node.ND_DEREF:
		gen(n.Right)
		fmt.Println("    pop rax")
		fmt.Println("    mov rax, [rax]")
		fmt.Println("    push rax")
		return
	case node.ND_DEFINE_VAR:
		gen(n.Left)
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
	switch n.Kind {
	case node.ND_LVAR:
		fmt.Println("    mov rax, rbp")
		fmt.Println(fmt.Sprintf("    sub rax, %d", n.Variable.Offset+n.ArrayIndex*8))
		fmt.Println("    push rax")
	case node.ND_DEREF:
		current := n.Right

		fmt.Println("    mov rax, rbp")
		fmt.Println(fmt.Sprintf("    sub rax, %d", current.Variable.Offset))
		fmt.Println("    push rax")

		for current.Variable.IsPointerType() {
			if err := current.Variable.Next(); err != nil {
				panic(err)
			}
			fmt.Println("    pop rax")
			fmt.Println("    mov rax, [rax]")
			fmt.Println("    push rax")
		}
		// log.Println(current.Variable)
	default:
		panic(fmt.Sprintf("%s is not supported type", n.Kind))
	}
}

var counter = 0

func getLabelCount() int {
	counter++
	return counter
}
