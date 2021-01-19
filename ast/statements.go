package ast

import "github.com/ElaraLang/elara/lexer"

type ExpressionStatement struct {
	Token      lexer.Token
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
	Token     lexer.Token
	Condition Expression
	Body      Statement
}

type ExtendStatement struct {
	Token      lexer.Token
	Identifier Identifier
	Alias      Identifier
	Body       BlockStatement
}

type BlockStatement struct {
	Token lexer.Token
	Block []Statement
}

type TypeStatement struct {
	Token      lexer.Token
	Identifier Identifier
	Contract   Type
}

type GenerifiedStatement struct {
	Token     lexer.Token
	Contracts []GenericContract
	Statement Statement
}

type ReturnStatement struct {
	Token lexer.Token
	Value Expression
}
