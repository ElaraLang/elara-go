package parser

import (
	"github.com/ElaraLang/elara/ast"
	"github.com/ElaraLang/elara/lexer"
)

type Parser struct {
	Tape              TokenTape
	statementParslets map[lexer.TokenType]statementParslet
	prefixParslets    map[lexer.TokenType]prefixParslet
	infixParslets     map[lexer.TokenType]infixParselet
}

func NewParser(tokens []lexer.Token, channel chan lexer.Token) Parser {
	return Parser{Tape: NewTokenTape(tokens, channel)}
}

func NewReplParser(channel chan lexer.Token) Parser {
	return Parser{Tape: NewReplTokenTape(channel)}
}

type (
	statementParslet func() ast.Statement
	prefixParslet    func() ast.Expression
	infixParselet    func(ast.Expression) ast.Expression
)

func (p *Parser) registerPrefix(tokenType lexer.TokenType, function prefixParslet) {
	p.prefixParslets[tokenType] = function
}
func (p *Parser) registerInfix(tokenType lexer.TokenType, function infixParselet) {
	p.infixParslets[tokenType] = function
}

func (p *Parser) parseStatement() ast.Statement {
	parseStmt := p.statementParslets[p.Tape.Current().TokenType]
	if parseStmt == nil {
		return p.parseExpressionStatement()
	}
	return parseStmt()
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	parsePrefix := p.prefixParslets[p.Tape.Current().TokenType]
	if parsePrefix == nil {
		// panic
		return nil
	}
	expr := parsePrefix()
	for !p.Tape.ValidationPeek(0, lexer.NEWLINE) && precedence < precedenceOf(p.Tape.Current().TokenType) {
		infix := p.infixParslets[p.Tape.Current().TokenType]
		if infix == nil {
			return expr
		}
		expr = infix(expr)
	}
	return expr
}
