package ast

import "github.com/ElaraLang/elara/lexer"

// Program represents the parsed program
// Information about the namespace/module of the  and imports required are present
type Program struct {
	Statements []Statement
}

// Node represents a node in the syntax tree
// Helps in printing nodes out for debugging
type Node interface {
	ToString() string
	TokenValue() string
}

// Statement is a node that syntactically does not have a value
// but may be expressed an "Unit" or similar
type Statement interface {
	ToString() string
	statementNode()
}

// Expression is a node that composes values with nodes
type Expression interface {
	ToString() string
	expressionNode()
}

// Type is a syntax tree node that represents a Type or a contract
type Type interface {
	ToString() string
	typeNode()
}

// Identifier represents an identifier leaf used by the syntax tree to represent a "Name"
type Identifier struct {
	Token lexer.Token
	Name  string
}
