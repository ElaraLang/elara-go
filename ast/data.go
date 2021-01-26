package ast

import "github.com/ElaraLang/elara/lexer"

// Module represents the namespace of a program.
// Helps segregate different components and avoid identifier conflicts
// Represented as Root / Sub or just Root if Sub is not provided
type Module struct {
	Pkg    string
	PkgIds []Identifier
}

func (m *Module) ToString() string {
	return m.Pkg
}

type Parameter struct {
	Type       Type
	Identifier Identifier
}

func (p *Parameter) ToString() string {
	res := ""
	if p.Type != nil {
		res += p.Type.ToString() + " "
	}
	res += p.Identifier.Name
	return res
}

type Entry struct {
	Key   Expression
	Value Expression
}

func (p *Entry) ToString() string {
	return "(" + p.Key.ToString() + " : " + p.Value.ToString() + ")"
}

type StructField struct {
	Mutable    bool
	Lazy       bool
	Open       bool
	Type       Type
	Identifier Identifier
	Default    Expression
}

func (p *StructField) ToString() string {
	return p.Type.ToString() + " " + p.Identifier.Name
}

type NamedContract struct {
	Token      lexer.Token
	Identifier Identifier
	Type       Type
}
