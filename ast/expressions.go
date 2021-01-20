package ast

import (
	"github.com/ElaraLang/elara/lexer"
)

type BinaryExpression struct {
	Token    lexer.Token
	Left     Expression
	Operator lexer.Token
	Right    Expression
}

type UnaryExpression struct {
	Token    lexer.Token
	Operator lexer.Token
	Right    Expression
}

type PropertyExpression struct {
	Token    lexer.Token
	Context  Expression
	Variable Identifier
}

type IfExpression struct {
	Token      lexer.Token
	Condition  Expression
	MainBranch Statement
	ElseBranch Statement
}

type AccessExpression struct {
	Token      lexer.Token
	Expression Expression
	Index      Expression
}

type CallExpression struct {
	Token      lexer.Token
	Expression Expression
	Arguments  []Expression
}

type TypeCastExpression struct {
	Token      lexer.Token
	Expression Expression
	Type       Type
}

type TypeCheckExpression struct {
	Token      lexer.Token
	Expression Expression
	Type       Type
}

// Literals

type BooleanLiteral struct {
	Token lexer.Token
	Value bool
}

type IntegerLiteral struct {
	Token lexer.Token
	Value int64
}

type FloatLiteral struct {
	Token lexer.Token
	Value float64
}

type CharLiteral struct {
	Token lexer.Token
	Value rune
}

type StringLiteral struct {
	Token lexer.Token
	Value string
}

type FunctionLiteral struct {
	Token      lexer.Token
	ReturnType Type
	Parameters []Parameter
	Body       Statement
}

type MapLiteral struct {
	Token   lexer.Token
	Entries []Entry
}

type CollectionLiteral struct {
	Token    lexer.Token
	Elements []Expression
}
