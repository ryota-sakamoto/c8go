package node

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/ryota-sakamoto/c8go/token"
	"github.com/ryota-sakamoto/c8go/util"
	"github.com/ryota-sakamoto/c8go/vars"
)

type NodeKind int

const (
	_ NodeKind = iota
	ND_ADD
	ND_SUB
	ND_DEREF
	ND_ADDR
	ND_MUL
	ND_DIV
	ND_LVAR
	ND_NUM
	ND_FUNC      // func()
	ND_CALL_FUNC // call func()

	ND_DEFINE_VAR

	ND_ASSIGN

	ND_GT // > , but not use
	ND_GE // >=, but not use
	ND_LT // <
	ND_LE // <=
	ND_EQ // ==
	ND_NE // !=

	ND_RETURN  // return
	ND_IF      // if
	ND_ELSE    // else
	ND_IF_ELSE // if & else
	ND_WHILE   // while

	ND_BLOCK // {}
)

type Node struct {
	Kind             NodeKind
	Left             *Node
	Right            *Node
	Block            []*Node
	Val              int
	Variable         vars.Variable
	Name             string
	Args             []*Node
	DefineArgsOffset []int
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

func NewNodeFunc(name string, block []*Node, args []int) *Node {
	node := Node{
		Kind:             ND_FUNC,
		Name:             name,
		Block:            block,
		DefineArgsOffset: args,
	}

	return &node
}

func NewNodeBlock(block []*Node) *Node {
	node := Node{
		Kind:  ND_BLOCK,
		Block: block,
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

func NewNodeLVar(v vars.Variable) *Node {
	node := Node{
		Kind:     ND_LVAR,
		Variable: v,
	}

	return &node
}

func NewNodeCallFunc(name string, args []*Node) *Node {
	node := Node{
		Kind: ND_CALL_FUNC,
		Name: name,
		Args: args,
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
		if err := np.token.ConsumeReserved("int"); err != nil {
			return nil, errors.WithStack(err)
		}

		name, err := np.token.ConsumeIndent()
		if err != nil {
			return nil, errors.WithStack(err)
		}
		if err := np.token.ConsumeReserved("("); err != nil {
			return nil, errors.WithStack(err)
		}

		args := []int{}
		first := true
		for !np.token.Expect(")") {
			if first {
				first = false
			} else {
				if err := np.token.ConsumeReserved(","); err != nil {
					return nil, errors.WithStack(err)
				}
			}

			if err := np.token.ConsumeReserved("int"); err != nil {
				return nil, errors.WithStack(err)
			}

			name, err := np.token.ConsumeIndent()
			if err != nil {
				return nil, errors.WithStack(err)
			}

			if _, ok := locals.Get(name); ok {
				return nil, util.CompileError{
					Input:   np.token.GetInput(),
					Message: fmt.Sprintf("%s is already defined.", name),
					Pos:     np.token.GetPos(),
				}
			}
			locals.Set(vars.NewVariable(name, vars.IntType))
			variable, _ := locals.Get(name)

			args = append(args, variable.Offset)
		}

		if err := np.token.ConsumeReserved(")"); err != nil {
			return nil, errors.WithStack(err)
		}
		if err := np.token.ConsumeReserved("{"); err != nil {
			return nil, errors.WithStack(err)
		}

		block := []*Node{}
		for !np.token.Expect("}") {
			node, err := np.Stmt()
			if err != nil {
				return nil, err
			}
			block = append(block, node)
		}

		if err := np.token.ConsumeReserved("}"); err != nil {
			return nil, errors.WithStack(err)
		}

		funcNode := NewNodeFunc(name, block, args)
		result = append(result, funcNode)
	}
	return result, nil
}

func (np *NodeParser) Stmt() (*Node, error) {
	if np.token.Expect("{") {
		if err := np.token.ConsumeReserved("{"); err != nil {
			return nil, errors.WithStack(err)
		}

		block := []*Node{}
		for !np.token.Expect("}") {
			node, err := np.Stmt()
			if err != nil {
				return nil, err
			}

			block = append(block, node)
		}

		if err := np.token.ConsumeReserved("}"); err != nil {
			return nil, errors.WithStack(err)
		}

		return NewNodeBlock(block), nil
	}

	if np.token.Expect("return") {
		if err := np.token.ConsumeReserved("return"); err != nil {
			return nil, errors.WithStack(err)
		}

		node, err := np.Expr()
		if err != nil {
			return nil, err
		}

		if err := np.token.ConsumeReserved(";"); err != nil {
			return nil, errors.WithStack(err)
		}

		return NewNode(ND_RETURN, nil, node), nil
	}

	if np.token.Expect("if") {
		if err := np.token.ConsumeReserved("if"); err != nil {
			return nil, errors.WithStack(err)
		}

		if err := np.token.ConsumeReserved("("); err != nil {
			return nil, errors.WithStack(err)
		}

		node1, err := np.Expr()
		if err != nil {
			return nil, err
		}

		if err := np.token.ConsumeReserved(")"); err != nil {
			return nil, errors.WithStack(err)
		}

		ifNode, err := np.Stmt()
		if err != nil {
			return nil, err
		}

		if np.token.Expect("else") {
			if err := np.token.ConsumeReserved("else"); err != nil {
				return nil, errors.WithStack(err)
			}

			elseNode, err := np.Stmt()
			if err != nil {
				return nil, err
			}

			ifNode = NewNode(ND_ELSE, ifNode, elseNode)
			return NewNode(ND_IF_ELSE, node1, ifNode), nil
		}

		return NewNode(ND_IF, node1, ifNode), nil
	}

	if np.token.Expect("while") {
		if err := np.token.ConsumeReserved("while"); err != nil {
			return nil, errors.WithStack(err)
		}

		if err := np.token.ConsumeReserved("("); err != nil {
			return nil, errors.WithStack(err)
		}

		node, err := np.Expr()
		if err != nil {
			return nil, err
		}

		if err := np.token.ConsumeReserved(")"); err != nil {
			return nil, errors.WithStack(err)
		}

		s, err := np.Stmt()
		if err != nil {
			return nil, err
		}

		return NewNode(ND_WHILE, node, s), nil
	}

	if np.token.Expect("int") {
		if err := np.token.ConsumeReserved("int"); err != nil {
			return nil, errors.WithStack(err)
		}

		head := vars.Variable{}
		current := &head
		isPointerType := false
		for np.token.Expect("*") {
			if err := np.token.Consume(); err != nil {
				return nil, errors.WithStack(err)
			}
			isPointerType = true

			next := vars.Variable{}
			current.Type = vars.PointerType
			current.Pointer = &next
			current = &next
		}

		name, err := np.token.ConsumeIndent()
		if err != nil {
			return nil, errors.WithStack(err)
		}

		if _, ok := locals.Get(name); ok {
			return nil, util.CompileError{
				Input:   np.token.GetInput(),
				Message: fmt.Sprintf("%s is already defined.", name),
				Pos:     np.token.GetPos(),
			}
		}

		if isPointerType {
			head.Name = name
			current.Type = vars.IntType
		} else {
			head = vars.NewVariable(name, vars.IntType)
		}
		locals.Set(head)

		if err := np.token.ConsumeReserved(";"); err != nil {
			return nil, errors.WithStack(err)
		}

		node, err := np.Stmt()
		if err != nil {
			return nil, err
		}

		return NewNode(ND_DEFINE_VAR, node, nil), nil
	}

	node, err := np.Expr()
	if err != nil {
		return nil, err
	}

	if err := np.token.ConsumeReserved(";"); err != nil {
		return nil, errors.WithStack(err)
	}

	return node, nil
}

func (np *NodeParser) Expr() (*Node, error) {
	return np.Assign()
}

func (np *NodeParser) Assign() (*Node, error) {
	node, err := np.Equality()
	if err != nil {
		return nil, err
	}

	if np.token.Expect("=") {
		err = np.token.ConsumeReserved("=")
		if err != nil {
			return nil, errors.WithStack(err)
		}

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
		if np.token.Expect("==") {
			err := np.token.ConsumeReserved("==")
			if err != nil {
				return nil, errors.WithStack(err)
			}

			right, err := np.Relational()
			if err != nil {
				return nil, err
			}
			node = NewNode(ND_EQ, node, right)
			continue
		}

		if np.token.Expect("!=") {
			err := np.token.ConsumeReserved("!=")
			if err != nil {
				return nil, errors.WithStack(err)
			}

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
		if np.token.Expect("<") {
			err := np.token.ConsumeReserved("<")
			if err != nil {
				return nil, errors.WithStack(err)
			}

			right, err := np.Add()
			if err != nil {
				return nil, err
			}
			node = NewNode(ND_LT, node, right)
			continue
		}

		if np.token.Expect("<=") {
			err := np.token.ConsumeReserved("<=")
			if err != nil {
				return nil, errors.WithStack(err)
			}

			right, err := np.Add()
			if err != nil {
				return nil, err
			}
			node = NewNode(ND_LE, node, right)
			continue
		}

		if np.token.Expect(">") {
			err := np.token.ConsumeReserved(">")
			if err != nil {
				return nil, errors.WithStack(err)
			}

			right, err := np.Add()
			if err != nil {
				return nil, err
			}
			node = NewNode(ND_LT, right, node)
			continue
		}

		if np.token.Expect(">=") {
			err := np.token.ConsumeReserved(">=")
			if err != nil {
				return nil, errors.WithStack(err)
			}

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
		if np.token.Expect("+") {
			err := np.token.ConsumeReserved("+")
			if err != nil {
				return nil, errors.WithStack(err)
			}

			right, err := np.Mul()
			if err != nil {
				return nil, err
			}
			node = NewNode(ND_ADD, node, right)
			continue
		}

		if np.token.Expect("-") {
			err := np.token.ConsumeReserved("-")
			if err != nil {
				return nil, errors.WithStack(err)
			}

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
		if np.token.Expect("*") {
			err := np.token.ConsumeReserved("*")
			if err != nil {
				return nil, errors.WithStack(err)
			}

			right, err := np.Unary()
			if err != nil {
				return nil, err
			}
			node = NewNode(ND_MUL, node, right)
			continue
		}

		if np.token.Expect("/") {
			err := np.token.ConsumeReserved("/")
			if err != nil {
				return nil, errors.WithStack(err)
			}

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

	if np.token.Expect("-") {
		err = np.token.ConsumeReserved("-")
		if err != nil {
			return nil, errors.WithStack(err)
		}

		right, err := np.Primary()
		if err != nil {
			return nil, err
		}
		return NewNode(ND_SUB, NewNodeNum(0), right), nil
	}

	if np.token.Expect("*") {
		err = np.token.ConsumeReserved("*")
		if err != nil {
			return nil, errors.WithStack(err)
		}

		right, err := np.Unary()
		if err != nil {
			return nil, err
		}

		return NewNode(ND_DEREF, nil, right), nil
	}

	if np.token.Expect("&") {
		err = np.token.ConsumeReserved("&")
		if err != nil {
			return nil, errors.WithStack(err)
		}

		left, err := np.Unary()
		if err != nil {
			return nil, err
		}

		return NewNode(ND_ADDR, left, nil), nil
	}

	return np.Primary()
}

func (np *NodeParser) Primary() (*Node, error) {
	if np.token.Expect("(") {
		err := np.token.ConsumeReserved("(")
		if err != nil {
			return nil, errors.WithStack(err)
		}

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

	name, err := np.token.ConsumeIndent()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if np.token.Expect("(") {
		err = np.token.ConsumeReserved("(")
		if err != nil {
			return nil, errors.WithStack(err)
		}

		args := []*Node{}
		first := true
		for !np.token.Expect(")") {
			if first {
				first = false
			} else {
				if err := np.token.ConsumeReserved(","); err != nil {
					return nil, errors.WithStack(err)
				}
			}

			argsNode, err := np.Expr()
			if err != nil {
				return nil, errors.WithStack(err)
			}

			args = append(args, argsNode)
		}

		err = np.token.ConsumeReserved(")")
		if err != nil {
			return nil, errors.WithStack(err)
		}

		return NewNodeCallFunc(name, args), nil
	}

	if variable, ok := locals.Get(name); !ok {
		return nil, util.CompileError{
			Input:   np.token.GetInput(),
			Message: fmt.Sprintf("%s is not defined.", name),
			Pos:     np.token.GetPos(),
		}
	} else {
		return NewNodeLVar(variable), nil
	}
}

var locals = vars.NewLocalVariales()
