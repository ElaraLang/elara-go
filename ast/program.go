package ast

import "github.com/ElaraLang/elara/lexer"

type Program struct {
	Statements []Statement
}

type Node interface {
	ToString() string
	TokenValue() string
}

type Statement interface {
	ToString() string
	statementNode()
}

type Expression interface {
	ToString() string
	expressionNode()
}

type Type interface {
	ToString() string
	typeNode()
}

type Identifier struct {
	token lexer.Token
	name  string
}
