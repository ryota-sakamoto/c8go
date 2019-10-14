package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		panic("Invalud args len")
	}

	s := os.Args[1]
	token, err := Tokenize(s)
	if err != nil {
		panic(err)
	}

	fmt.Println(".intel_syntax noprefix")
	fmt.Println(".global main")
	fmt.Println("main:")

	fn, err := token.ConsumeNumber()
	if err != nil {
		panic(err)
	}

	fmt.Println(fmt.Sprintf("    mov rax, %d", fn))
	for !token.isEOF() {
		if token.Expect('+') {
			n, e := token.ConsumeNumber()
			if e != nil {
				panic(e)
			}
			fmt.Println(fmt.Sprintf("    add rax, %d", n))
		}

		if token.Expect('-') {
			n, e := token.ConsumeNumber()
			if e != nil {
				panic(e)
			}
			fmt.Println(fmt.Sprintf("    sub rax, %d", n))
		}
	}

	fmt.Println("    ret")
}
