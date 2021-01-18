package ast

import "github.com/ElaraLang/elara/lexer"

type ExpressionStatement struct {
	Expression Expression
}

type DeclarationStatement struct {
	Token      lexer.Token
	Mutable    bool
	Lazy       bool
	Open       bool
	Identifier string
	Type       Type
	Value      Expression
}

type StructDefStatement struct {
	Token  lexer.Token
	Id     Identifier
	Fields []StructField
}

type WhileStatement struct {
	Condition Expression
	Body      Statement
}

type ExtendStatement struct {
	Identifier Identifier
	Alias      Identifier
	Body       BlockStatement
}

type BlockStatement struct {
	Block []Statement
}

type TypeStatement struct {
	Identifier Identifier
	Contract   Type
}

type GenerifiedStatement struct {
	Contracts []GenericContract
	Statement Statement
}

type ReturnStatement struct {
	Value Expression
}
