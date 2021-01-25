package parser

import (
	"github.com/ElaraLang/elara/lexer"
)

const (
	_ int = iota
	LOWEST
	EQUALS
	COMPARISON
	SUM
	PRODUCT
	PREFIX
	CALL
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
}

func precedenceOf(tok lexer.TokenType) int {
	if value, contains := precedences[tok]; contains {
		return value
	}
	return LOWEST
}
