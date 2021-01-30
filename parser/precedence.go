package parser

import (
	"github.com/ElaraLang/elara/lexer"
)

const (
	_ int = iota
	// Expression precedences
	Lowest
	InfixCall
	Assign
	CastAs
	CheckIs
	LogicOr
	LogicAnd
	Equals
	Comparison
	Sum
	Product
	Prefix
	Invoke

	// Type precedences
	TypeLowest
	TypeUnion
	TypeIntersection
	TypeMap
	TypeProperty
)

var precedences = map[lexer.TokenType]int{
	lexer.Identifier:   InfixCall,
	lexer.Equal:        Assign,
	lexer.As:           CastAs,
	lexer.Is:           CheckIs,
	lexer.Or:           LogicOr,
	lexer.And:          LogicAnd,
	lexer.Equals:       Equals,
	lexer.NotEquals:    Equals,
	lexer.GreaterEqual: Comparison,
	lexer.LesserEqual:  Comparison,
	lexer.LAngle:       Comparison,
	lexer.RAngle:       Comparison,
	lexer.Add:          Sum,
	lexer.Subtract:     Sum,
	lexer.Multiply:     Product,
	lexer.Slash:        Product,
	lexer.Dot:          Invoke,
	lexer.LParen:       Invoke,
	lexer.LSquare:      Invoke,
}

var typePrecedences = map[lexer.TokenType]int{
	lexer.TypeAnd: TypeUnion,
	lexer.TypeOr:  TypeIntersection,
	lexer.Colon:   TypeMap,
	lexer.Dot:     TypeProperty,
}

func precedenceOf(tok lexer.TokenType) int {
	if value, contains := precedences[tok]; contains {
		return value
	}
	return Lowest
}

func typePrecedenceOf(tok lexer.TokenType) int {
	if value, contains := typePrecedences[tok]; contains {
		return value
	}
	return TypeLowest
}
