package ast

import (
	"github.com/ElaraLang/elara/lexer"
)

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

func JoinToString(slice interface{}, separator string) string {
	res := ""

	switch slice := slice.(type) {
	case []string:
		l := len(slice)
		for i, v := range slice {
			res += v
			if i < l {
				res += separator
			}
		}
	case []Entry:
		l := len(slice)
		for i, v := range slice {
			res += "(" + v.Key.ToString() + ") :" + "(" + v.Value.ToString() + ")"
			if i < l {
				res += separator
			}
		}

	case []Parameter:
		l := len(slice)
		for i, v := range slice {
			res += v.ToString()
			if i < l {
				res += separator
			}
		}
	case []StructField:
		l := len(slice)
		for i, v := range slice {
			res += v.ToString()
			if i < l {
				res += separator
			}
		}
	case []Node:
		l := len(slice)
		for i, v := range slice {
			res += v.ToString()
			if i < l {
				res += separator
			}
		}
	default:
		res = "Unknown"
	}
	return res
}
