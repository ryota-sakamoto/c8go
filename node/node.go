package node

import (
	"github.com/ryota-sakamoto/c8go/token"
)

type NodeKind int

const (
	ND_UNKNOWN NodeKind = iota + 1
	ND_ADD
	ND_SUB
	ND_MUL
	ND_DIV
	ND_NUM
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
	default:
		return "Unknown"
	}
}

type Node struct {
	Kind  NodeKind
	Left  *Node
	Right *Node
	Val   int
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

type NodeParser struct {
	token *token.Token
}

func NewNodeParser(token *token.Token) *NodeParser {
	np := NodeParser{
		token: token,
	}

	return &np
}

func (np *NodeParser) Expr() (*Node, error) {
	node, err := np.Mul()
	if err != nil {
		return nil, err
	}

	for {
		err := np.token.Expect('+')
		if err == nil {
			right, err := np.Mul()
			if err != nil {
				return nil, err
			}
			node = NewNode(ND_ADD, node, right)
			continue
		}

		err = np.token.Expect('-')
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
	node, err := np.Primary()
	if err != nil {
		return nil, err
	}

	for {
		err := np.token.Expect('*')
		if err == nil {
			right, err := np.Primary()
			if err != nil {
				return nil, err
			}
			node = NewNode(ND_MUL, node, right)
			continue
		}

		err = np.token.Expect('/')
		if err == nil {
			right, err := np.Primary()
			if err != nil {
				return nil, err
			}
			node = NewNode(ND_DIV, node, right)
		}

		return node, nil
	}
}

func (np *NodeParser) Primary() (*Node, error) {
	err := np.token.Expect('(')
	if err == nil {
		node, err := np.Expr()
		if err != nil {
			return nil, err
		}

		err = np.token.Expect(')')
		if err != nil {
			return nil, err
		}

		return node, nil
	}
	n, err := np.token.ConsumeNumber()
	return NewNodeNum(n), err
}
