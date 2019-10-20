package token

import (
	"errors"
	"fmt"
	"strings"

	"github.com/ryota-sakamoto/c8go/util"
)

type TokenKind int

const (
	Unknown TokenKind = iota + 1
	TK_RESERVED
	TK_NUM
	TK_EOF
)

func (tk TokenKind) String() string {
	switch tk {
	case TK_RESERVED:
		return "TK_RESERVED"
	case TK_NUM:
		return "TK_NUM"
	case TK_EOF:
		return "TK_EOF"
	default:
		return "Unknown"
	}
}

type Token struct {
	kind TokenKind
	next *Token
	val  int
	s    string
	len  int

	input string
	pos   int
}

func (t *Token) isNumber() bool {
	return t.kind == TK_NUM
}

func (t *Token) isReserved() bool {
	return t.kind == TK_RESERVED
}

func (t *Token) isEOF() bool {
	return t.kind == TK_EOF
}

func (t Token) String() string {
	return fmt.Sprintf("s: %q, pos: %d, kind: %s, val: %d, tl: %d", t.s, t.pos, t.kind, t.val, t.len)
}

func (t *Token) Consume() error {
	if t.next == nil {
		return errors.New("next is nil")
	}
	next := t.next
	t.kind = next.kind
	t.val = next.val
	t.s = next.s
	t.next = next.next
	t.pos = next.pos
	t.len = next.len

	return nil
}

func (t *Token) ConsumeNumber() (int, error) {
	if !t.isNumber() {
		return 0, t.NewTokenError("current is not number: %+v", t)
	}
	v := t.val
	if err := t.Consume(); err != nil {
		return 0, err
	}

	return v, nil
}

func (t *Token) Expect(c string) error {
	if !t.isReserved() || t.s[0:t.len] != c {
		return t.NewTokenError("current is not reversed: %+v", t)
	}
	if err := t.Consume(); err != nil {
		return err
	}
	return nil
}

func Tokenize(s string) (*Token, error) {
	token := Token{
		input: s,
		pos:   0,
	}
	current := &token
	tl := 1
	back := false
	for len(s) > 0 {
		switch s[:tl] {
		case " ":
			s = s[1:]
		case "+", "-", "*", "/", "(", ")":
			current = newToken(TK_RESERVED, current, s, 1)
			s = s[1:]
		case "<", ">", "=", "!":
			if !back {
				tl++
			} else {
				current = newToken(TK_RESERVED, current, s, 1)
				s = s[1:]

				back = false
			}
		case "<=", ">=", "==", "!=":
			current = newToken(TK_RESERVED, current, s, 2)
			s = s[2:]

			tl = 1
		default:
			if tl != 1 {
				tl = 1
				back = true
				continue
			}

			tmp := s
			num, err := util.ParseInt(&s)
			if err != nil {
				return nil, tokenError(token.input, current.pos+1, err.Error())
			}
			current = newToken(TK_NUM, current, tmp, 1)
			current.val = num
		}
		current.pos++
	}
	current = newToken(TK_EOF, current, s, 1)
	current.pos++

	return token.next, nil
}

func newToken(kind TokenKind, current *Token, s string, len int) *Token {
	next := Token{
		kind:  kind,
		next:  nil,
		input: current.input,
		s:     s,
		len:   len,
		pos:   current.pos,
	}
	current.next = &next

	return &next
}

func (t *Token) NewTokenError(format string, a ...interface{}) error {
	return tokenError(t.input, t.pos, format, a...)
}

func tokenError(input string, pos int, format string, a ...interface{}) error {
	s := `%s
%s
%s`
	return fmt.Errorf(s, input, strings.Repeat(" ", pos-1)+"^", fmt.Sprintf(format, a...))
}
