package node

import (
	"github.com/pkg/errors"

	"github.com/ryota-sakamoto/c8go/token"
)

type NodeKind int

const (
	ND_UNKNOWN NodeKind = iota + 1
	ND_ADD
	ND_SUB
	ND_MUL
	ND_DIV
	ND_LVAR
	ND_NUM

	ND_ASSIGN

	ND_GT // > , but not use
	ND_GE // >=, but not use
	ND_LT // <
	ND_LE // <=
	ND_EQ // ==
	ND_NE // !=

	ND_RETURN // return
)

func (nk NodeKind) String() string {
	switch nk {
	case ND_ADD:
		return "ND_ADD"
	case ND_SUB:
		return "ND_SUB"
	case ND_MUL:
		return "ND_MUL"
	case ND_DIV:
		return "ND_DIV"
	case ND_NUM:
		return "ND_NUM"
	case ND_RETURN:
		return "ND_RETURN"
	default:
		return "Unknown"
	}
}

type Node struct {
	Kind   NodeKind
	Left   *Node
	Right  *Node
	Val    int
	Offset int
}

func (n Node) IsNum() bool {
	return n.Kind == ND_NUM
}

func NewNode(kind NodeKind, left *Node, right *Node) *Node {
	node := Node{
		Kind:  kind,
		Left:  left,
		Right: right,
	}

	return &node
}

func NewNodeNum(n int) *Node {
	node := Node{
		Kind: ND_NUM,
		Val:  n,
	}

	return &node
}

func NewNodeLVar(offset int) *Node {
	node := Node{
		Kind:   ND_LVAR,
		Offset: offset,
	}

	return &node
}

type NodeParser struct {
	token *token.Token
}

func NewNodeParser(token *token.Token) *NodeParser {
	np := NodeParser{
		token: token,
	}

	return &np
}

func (np *NodeParser) Program() ([]*Node, error) {
	result := []*Node{}
	for !np.token.IsEOF() {
		node, err := np.Stmt()
		if err != nil {
			return nil, err
		}
		result = append(result, node)
	}
	return result, nil
}

func (np *NodeParser) Stmt() (*Node, error) {
	re := np.token.ConsumeReserved("return")

	node, err := np.Expr()
	if err != nil {
		return nil, err
	}

	if re == nil {
		node = NewNode(ND_RETURN, nil, node)
	}

	return node, np.token.ConsumeReserved(";")
}

func (np *NodeParser) Expr() (*Node, error) {
	return np.Assign()
}

func (np *NodeParser) Assign() (*Node, error) {
	node, err := np.Equality()
	if err != nil {
		return nil, err
	}

	err = np.token.ConsumeReserved("=")
	if err == nil {
		right, err := np.Assign()
		if err != nil {
			return nil, err
		}
		node = NewNode(ND_ASSIGN, node, right)
	}

	return node, nil
}

func (np *NodeParser) Equality() (*Node, error) {
	node, err := np.Relational()
	if err != nil {
		return nil, err
	}

	for {
		err := np.token.ConsumeReserved("==")
		if err == nil {
			right, err := np.Relational()
			if err != nil {
				return nil, err
			}
			node = NewNode(ND_EQ, node, right)
			continue
		}

		err = np.token.ConsumeReserved("!=")
		if err == nil {
			right, err := np.Relational()
			if err != nil {
				return nil, err
			}
			node = NewNode(ND_NE, node, right)
			continue
		}

		return node, nil
	}
}

func (np *NodeParser) Relational() (*Node, error) {
	node, err := np.Add()
	if err != nil {
		return nil, err
	}

	for {
		err := np.token.ConsumeReserved("<")
		if err == nil {
			right, err := np.Add()
			if err != nil {
				return nil, err
			}
			node = NewNode(ND_LT, node, right)
			continue
		}

		err = np.token.ConsumeReserved("<=")
		if err == nil {
			right, err := np.Add()
			if err != nil {
				return nil, err
			}
			node = NewNode(ND_LE, node, right)
			continue
		}

		err = np.token.ConsumeReserved(">")
		if err == nil {
			right, err := np.Add()
			if err != nil {
				return nil, err
			}
			node = NewNode(ND_LT, right, node)
			continue
		}

		err = np.token.ConsumeReserved(">=")
		if err == nil {
			right, err := np.Add()
			if err != nil {
				return nil, err
			}
			node = NewNode(ND_LE, right, node)
			continue
		}

		return node, nil
	}
}

func (np *NodeParser) Add() (*Node, error) {
	node, err := np.Mul()
	if err != nil {
		return nil, err
	}

	for {
		err := np.token.ConsumeReserved("+")
		if err == nil {
			right, err := np.Mul()
			if err != nil {
				return nil, err
			}
			node = NewNode(ND_ADD, node, right)
			continue
		}

		err = np.token.ConsumeReserved("-")
		if err == nil {
			right, err := np.Mul()
			if err != nil {
				return nil, err
			}
			node = NewNode(ND_SUB, node, right)
			continue
		}

		return node, nil
	}
}

func (np *NodeParser) Mul() (*Node, error) {
	node, err := np.Unary()
	if err != nil {
		return nil, err
	}

	for {
		err := np.token.ConsumeReserved("*")
		if err == nil {
			right, err := np.Unary()
			if err != nil {
				return nil, err
			}
			node = NewNode(ND_MUL, node, right)
			continue
		}

		err = np.token.ConsumeReserved("/")
		if err == nil {
			right, err := np.Unary()
			if err != nil {
				return nil, err
			}
			node = NewNode(ND_DIV, node, right)
		}

		return node, nil
	}
}

func (np *NodeParser) Unary() (*Node, error) {
	err := np.token.ConsumeReserved("+")
	if err == nil {
		return np.Primary()
	}

	err = np.token.ConsumeReserved("-")
	if err == nil {
		right, err := np.Primary()
		if err != nil {
			return nil, err
		}
		return NewNode(ND_SUB, NewNodeNum(0), right), nil
	}

	return np.Primary()
}

func (np *NodeParser) Primary() (*Node, error) {
	err := np.token.ConsumeReserved("(")
	if err == nil {
		node, err := np.Expr()
		if err != nil {
			return nil, err
		}

		err = np.token.ConsumeReserved(")")
		if err != nil {
			return nil, errors.WithStack(err)
		}

		return node, nil
	}

	n, err := np.token.ConsumeNumber()
	if err == nil {
		return NewNodeNum(n), errors.WithStack(err)
	}

	name, err := np.token.GetVariableName()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if _, ok := locals.get(name); !ok {
		locals.set(name)
	}
	offset, _ := locals.get(name)
	return NewNodeLVar(offset), np.token.Consume()
}

var locals = localVariale{
	vars:      map[string]int{},
	maxOffset: 0,
}

type localVariale struct {
	vars      map[string]int
	maxOffset int
}

func (l localVariale) get(name string) (int, bool) {
	v, ok := l.vars[name]
	return v, ok
}

func (l *localVariale) set(name string) {
	l.maxOffset = l.maxOffset + 8
	l.vars[name] = l.maxOffset
}
