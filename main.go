package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"unicode"
)

func main() {
	if len(os.Args) < 2 {
		panic("Invalud args len")
	}

	n := os.Args[1]
	num, err := parseInt(&n)
	if err != nil {
		panic(err)
	}

	fmt.Println(".intel_syntax noprefix")
	fmt.Println(".global main")
	fmt.Println("main:")

	fmt.Println(fmt.Sprintf("    mov rax, %d", num))
	for len(n) != 0 {
		c := n[0]
		n = n[1:]
		num, _ := parseInt(&n)

		switch c {
		case '+':
			fmt.Println(fmt.Sprintf("    add rax, %d", num))
		case '-':
			fmt.Println(fmt.Sprintf("    sub rax, %d", num))
		default:
			panic(fmt.Sprintf("Invalid char: %c", n[0]))
		}
	}

	fmt.Println("    ret")
}

func parseInt(s *string) (int, error) {
	t := []rune{}
	index := 0
	for _, c := range *s {
		if unicode.IsDigit(c) {
			t = append(t, c)
			index++
		} else {
			break
		}
	}

	if index == 0 {
		return 0, errors.New("")
	}

	*s = (*s)[index:]
	return strconv.Atoi(string(t))
}
