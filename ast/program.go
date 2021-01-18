package ast

import "github.com/ElaraLang/elara/lexer"

type Program struct {
	Statements []Statement
}

type Node interface {
	toString() string
	TokenValue() string
}

type Statement interface {
	toString() string
	statementNode()
}

type Expression interface {
	toString() string
	expressionNode()
}

type Identifier struct {
	token lexer.Token
	name  string
}
