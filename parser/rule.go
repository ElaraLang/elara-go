package parser

import (
	"github.com/ElaraLang/elara/ast"
	"github.com/ElaraLang/elara/lexer"
	"strconv"
)

func (p *Parser) initializePrefixParselets() {
	p.prefixParseFunctions = make(map[lexer.TokenType]parsePrefix)
	p.registerPrefix(lexer.Int, p.parseInteger)
	p.registerPrefix(lexer.Float, p.parseFloat)
}

func (p *Parser) parseInteger() ast.Expression {
	integer := p.Tape.Consume(lexer.Int)
	value, err := strconv.ParseInt(string(integer.Text), 10, 64)
	if err != nil {
		// panic
	}
	return &ast.IntegerLiteral{Token: integer, Value: value}
}

func (p *Parser) parseFloat() ast.Expression {
	integer := p.Tape.Consume(lexer.Float)
	value, err := strconv.ParseFloat(string(integer.Text), 10)
	if err != nil {
		// panic
	}
	return &ast.FloatLiteral{Token: integer, Value: value}
}
