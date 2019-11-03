package token

import (
	"fmt"
	"strings"
)

var (
	NotReserverdError = tokenError{errorType: "NotReserverdError"}
	NotExpectedError  = tokenError{errorType: "NotExpectedError"}
	NotVariableError  = tokenError{errorType: "NotVariableError"}
	EmptyVarName      = tokenError{errorType: "EmptyVarName"}
	NotNumberError    = tokenError{errorType: "NotNumberError"}
)

type tokenError struct {
	errorType string
	input     string
	message   string
	pos       int
}

func (t tokenError) Error() string {
	s := `compile error: %s
----------
%s
%s
%s
----------`

	return fmt.Sprintf(s, t.errorType, t.input, strings.Repeat(" ", t.pos-1)+"^", t.message)
}

func (t *tokenError) New(input string, message string, pos int) error {
	t.pos = pos
	t.input = input
	t.message = message
	return t
}
