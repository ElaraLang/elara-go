package ast

import "github.com/ElaraLang/elara/lexer"

type Contract struct {
	Token      lexer.Token
	Identifier Identifier
	Type       Type
}

type ContractualType struct {
	Token     lexer.Token
	Contracts []Contract
}

type AlgebraicType struct {
	Token     lexer.Token
	Left      Type
	Operation lexer.Token
	Right     Type
}

type FunctionType struct {
	Token      lexer.Token
	ParamTypes []Type
	ReturnType Type
}

type CollectionType struct {
	Token lexer.Token
	Type  Type
}

type MapType struct {
	Token     lexer.Token
	KeyType   Type
	ValueType Type
}

type PrimaryType struct {
	Token      lexer.Token
	Identifier Identifier
}
