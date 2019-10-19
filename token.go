package main

import (
	"errors"
	"fmt"

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
	return fmt.Sprintf("kind: %s, val: %d, s: %s", t.kind, t.val, t.s)
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

	return nil
}

func (t *Token) ConsumeNumber() (int, error) {
	if !t.isNumber() {
		return 0, fmt.Errorf("current is not number: %+v", t)
	}
	v := t.val
	if err := t.Consume(); err != nil {
		return 0, err
	}

	return v, nil
}

func (t *Token) Expect(c byte) bool {
	if !t.isReserved() || t.s[0] != c {
		return false
	}
	if err := t.Consume(); err != nil {
		return false
	}
	return true
}

func Tokenize(s string) (*Token, error) {
	token := Token{}
	current := &token
	for len(s) > 0 {
		switch s[0] {
		case ' ':
			s = s[1:]
		case '+', '-', '*', '/', '(', ')':
			current = newToken(TK_RESERVED, current, s)
			s = s[1:]
		default:
			tmp := s
			num, err := util.ParseInt(&s)
			if err != nil {
				return nil, err
			}
			current = newToken(TK_NUM, current, tmp)
			current.val = num
		}
	}
	current = newToken(TK_EOF, current, s)

	return token.next, nil
}

func newToken(kind TokenKind, current *Token, s string) *Token {
	next := Token{
		kind: kind,
		next: nil,
		s:    s,
	}
	current.next = &next

	return &next
}
