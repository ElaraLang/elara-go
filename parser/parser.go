package parser

import (
	"github.com/ElaraLang/elara/ast"
	"github.com/ElaraLang/elara/lexer"
)

type Parser struct {
	Tape                 TokenTape
	prefixParseFunctions map[lexer.TokenType]prefixParslet
	infixParseFunctions  map[lexer.TokenType]parseInfix
}

func NewParser(tokens []lexer.Token, channel chan lexer.Token) Parser {
	return Parser{Tape: NewTokenTape(tokens, channel)}
}

func NewReplParser(channel chan lexer.Token) Parser {
	return Parser{Tape: NewReplTokenTape(channel)}
}

type (
	prefixParslet func() ast.Expression
	parseInfix    func(ast.Expression) ast.Expression
)

func (p *Parser) registerPrefix(tokenType lexer.TokenType, function prefixParslet) {
	p.prefixParseFunctions[tokenType] = function
}
func (p *Parser) registerInfix(tokenType lexer.TokenType, function parseInfix) {
	p.infixParseFunctions[tokenType] = function
}
