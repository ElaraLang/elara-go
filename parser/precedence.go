package parser

import (
	"github.com/ElaraLang/elara/lexer"
)

const (
	_ int = iota
	// Expression precedences
	LOWEST
	EQUALS
	COMPARISON
	SUM
	PRODUCT
	PREFIX
	CALL

	// Type precedences
	TYPE_LOWEST
	TYPE_UNION
	TYPE_INTERSECTION
)

var precedences = map[lexer.TokenType]int{
	lexer.Equals:       EQUALS,
	lexer.NotEquals:    EQUALS,
	lexer.GreaterEqual: COMPARISON,
	lexer.LesserEqual:  COMPARISON,
	lexer.LAngle:       COMPARISON,
	lexer.RAngle:       COMPARISON,
	lexer.Add:          SUM,
	lexer.Subtract:     SUM,
	lexer.Multiply:     PRODUCT,
	lexer.Slash:        PRODUCT,
	lexer.LParen:       CALL,
	lexer.LSquare:      CALL,
}

var typePrecedences = map[lexer.TokenType]int{
	lexer.TypeAnd: TYPE_UNION,
	lexer.TypeOr:  TYPE_INTERSECTION,
}

func precedenceOf(tok lexer.TokenType) int {
	if value, contains := precedences[tok]; contains {
		return value
	}
	return LOWEST
}

func typePrecedenceOf(tok lexer.TokenType) int {
	if value, contains := typePrecedences[tok]; contains {
		return value
	}
	return TYPE_LOWEST
}
