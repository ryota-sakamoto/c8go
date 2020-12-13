package util

import (
	"fmt"
	"strings"
)

var (
	NotReserverdError = CompileError{errorType: "NotReserverdError"}
	NotExpectedError  = CompileError{errorType: "NotExpectedError"}
	NotVariableError  = CompileError{errorType: "NotVariableError"}
	EmptyVarName      = CompileError{errorType: "EmptyVarName"}
	NotNumberError    = CompileError{errorType: "NotNumberError"}
)

type CompileError struct {
	errorType string
	Input     string
	Message   string
	Pos       int
	Line      int
}

func (t CompileError) Error() string {
	s := fmt.Sprintf(`compile error: %s
----------
`, t.errorType)
	for i, v := range strings.Split(t.Input, "\n") {
		s += fmt.Sprintf("%s\n", v)
		if t.Line == i {
			s += fmt.Sprintf("%s\n%s\n", strings.Repeat("~", t.Pos-1)+"^", t.Message)
		}
	}

	return s
}

func (t *CompileError) New(input string, message string, pos int, line int) error {
	t.Pos = pos
	t.Line = line
	t.Input = input
	t.Message = message
	return t
}
