package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) < 2 {
		panic("Invalud args len")
	}

	n := os.Args[1]
	num, err := strconv.Atoi(n)
	if err != nil {
		panic(err)
	}

	fmt.Println(".intel_syntax noprefix")
	fmt.Println(".global main")
	fmt.Println("main:")
	fmt.Println(fmt.Sprintf("    mov rax, %d", num))
	fmt.Println("    ret")
}
