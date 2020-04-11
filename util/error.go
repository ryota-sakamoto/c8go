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
}

func (t CompileError) Error() string {
	s := `compile error: %s
----------
%s
%s
%s
----------`

	return fmt.Sprintf(s, t.errorType, t.Input, strings.Repeat("~", t.Pos-1)+"^", t.Message)
}

func (t *CompileError) New(input string, message string, pos int) error {
	t.Pos = pos
	t.Input = input
	t.Message = message
	return t
}
