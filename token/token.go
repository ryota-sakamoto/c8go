package token

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/ryota-sakamoto/c8go/util"
)

type TokenKind int

const (
	Unknown TokenKind = iota + 1
	TK_RESERVED
	TK_RETURN
	TK_IF
	TK_ELSE
	TK_WHILE
	TK_IDENT
	TK_NUM
	TK_EOF
)

func (tk TokenKind) String() string {
	switch tk {
	case TK_RESERVED:
		return "TK_RESERVED"
	case TK_RETURN:
		return "TK_RETURN"
	case TK_IDENT:
		return "TK_IDENT"
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
	return t.kind == TK_RESERVED ||
		t.kind == TK_RETURN ||
		t.kind == TK_IF ||
		t.kind == TK_ELSE ||
		t.kind == TK_WHILE
}

func (t *Token) IsEOF() bool {
	return t.kind == TK_EOF
}

func (t *Token) GetVariableName() (string, error) {
	if t.kind != TK_IDENT {
		return "", t.NewTokenError(NotVariableError, "current is not variable: %+v", t)
	}
	return t.s[:t.len], nil
}

func (t Token) String() string {
	return fmt.Sprintf("s: %q, pos: %d, kind: %s, val: %d, tl: %d", t.s, t.pos, t.kind, t.val, t.len)
}

func (t *Token) Expect(c string) bool {
	return t.len == len(c) && t.s[0:t.len] == c
}

func (t *Token) Consume() error {
	if t.next == nil {
		return errors.New("next is nil")
	}
	next := t.next
	*t = *next

	return nil
}

func (t *Token) ConsumeNumber() (int, error) {
	if !t.isNumber() {
		return 0, t.NewTokenError(NotNumberError, "current is not number: %+v", t)
	}
	v := t.val
	if err := t.Consume(); err != nil {
		return 0, err
	}

	return v, nil
}

func (t *Token) ConsumeReserved(c string) error {
	if !t.isReserved() {
		return t.NewTokenError(NotReserverdError, "current is not reversed: %+v, want: %+v", t, c)
	}
	if !t.Expect(c) {
		return t.NewTokenError(NotExpectedError, "current is not expected reversed: %+v, want: %+v", t, c)
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
	for len(s) > 0 {
		if s[:1] == " " {
			s = s[1:]
			current.pos++
			continue
		}

		isReserved := false
		for _, v := range []string{"+", "-", "*", "/", "(", ")", ";", "{", "}"} {
			if s[:1] == v {
				isReserved = true
				break
			}
		}
		if isReserved {
			current = newToken(TK_RESERVED, current, s, 1)
			s = s[1:]
			current.pos++
			continue
		}

		isComparisonReserved := false
		for _, v := range []string{"<", ">", "=", "!"} {
			if s[:1] == v {
				isComparisonReserved = true
				break
			}
		}
		if isComparisonReserved {
			f := false
			for _, v := range []string{"<=", ">=", "==", "!="} {
				if s[:2] == v {
					f = true
					break
				}
			}
			if f {
				current = newToken(TK_RESERVED, current, s, 2)
				s = s[2:]
				current.pos += 2
			} else {
				current = newToken(TK_RESERVED, current, s, 1)
				s = s[1:]
				current.pos++
			}
			continue
		}

		if len(s) >= 2 && s[:2] == "if" && !util.IsAlnum(s[2]) {
			current = newToken(TK_IF, current, s, 2)
			s = s[2:]
			current.pos += 2
			continue
		}

		if len(s) >= 4 && s[:4] == "else" && !util.IsAlnum(s[4]) {
			current = newToken(TK_IF, current, s, 4)
			s = s[4:]
			current.pos += 4
			continue
		}

		if len(s) >= 6 && s[:6] == "return" && !util.IsAlnum(s[6]) {
			current = newToken(TK_RETURN, current, s, 6)
			s = s[6:]
			current.pos += 6
			continue
		}

		if len(s) >= 5 && s[:5] == "while" && !util.IsAlnum(s[5]) {
			current = newToken(TK_WHILE, current, s, 5)
			s = s[5:]
			current.pos += 5
			continue
		}

		if _, err := strconv.Atoi(s[:1]); err == nil {
			tmp := s
			num, err := util.ParseInt(&s)
			if err != nil {
				return nil, tokenError{
					input:   token.input,
					message: err.Error(),
					pos:     current.pos + 1,
				}
			}
			current = newToken(TK_NUM, current, tmp, 1)
			current.val = num
			current.pos++
			continue
		}

		tmp := s
		varName := ""
		for len(s) > 0 {
			c := s[:1]
			if util.IsAlnum(c[0]) {
				s = s[1:]
				varName += c
				continue
			}
			break
		}
		if len(varName) == 0 {
			return nil, tokenError{
				input:   token.input,
				message: "varName is empty",
				pos:     current.pos + 1,
			}
		}

		current = newToken(TK_IDENT, current, tmp, len(varName))
		current.pos++
		continue
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

func (t *Token) NewTokenError(e tokenError, format string, a ...interface{}) error {
	return e.New(t.input, fmt.Sprintf(format, a...), t.pos)
}
