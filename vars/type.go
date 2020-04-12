package vars

import (
	"github.com/pkg/errors"
)

type LocalVariales struct {
	vars      map[string]Variable
	maxOffset int
}

func (l LocalVariales) Get(name string) (Variable, bool) {
	v, ok := l.vars[name]
	return v, ok
}

func (l *LocalVariales) Set(v Variable) {
	l.maxOffset = l.maxOffset + 8
	v.Offset = l.maxOffset
	l.vars[v.Name] = v
}

func NewLocalVariales() LocalVariales {
	return LocalVariales{
		vars:      map[string]Variable{},
		maxOffset: 0,
	}
}

type Variable struct {
	Name    string
	Type    Type
	Pointer *Variable
	Offset  int
}

func NewVariable(name string, t Type) Variable {
	return Variable{
		Name: name,
		Type: t,
	}
}

func (v Variable) IsPointerType() bool {
	return v.Pointer != nil && v.Type == PointerType
}

func (v *Variable) Next() error {
	if v.Type != PointerType {
		return errors.Errorf("%s is not PointerType", v.Type)
	}

	v.Type = v.Pointer.Type
	v.Pointer = v.Pointer.Pointer

	return nil
}

type Type int

const (
	_ Type = iota
	IntType
	PointerType
)

var s = []string{"Unknown", "IntType", "PointerType"}

func (t Type) String() string {
	return s[t]
}
