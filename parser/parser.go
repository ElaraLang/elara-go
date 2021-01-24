package parser

import (
	"github.com/ElaraLang/elara/ast"
	"github.com/ElaraLang/elara/lexer"
)

type Parser struct {
	Tape              *TokenTape
	statementParslets map[lexer.TokenType]statementParslet
	prefixParslets    map[lexer.TokenType]prefixParslet
	infixParslets     map[lexer.TokenType]infixParslet
}

func NewParser(tokens []lexer.Token, channel chan lexer.Token) Parser {
	tape := NewTokenTape(tokens, channel)
	p := Parser{Tape: &tape}
	p.initPrefixParselets()
	p.initInfixParselets()
	p.initStatementParselets()
	return p
}

func NewReplParser(channel chan lexer.Token) Parser {
	tape := NewReplTokenTape(channel)
	p := Parser{Tape: &tape}
	p.initPrefixParselets()
	p.initInfixParselets()
	p.initStatementParselets()
	return p
}

type (
	statementParslet func() ast.Statement
	prefixParslet    func() ast.Expression
	infixParslet     func(ast.Expression) ast.Expression
)

func (p *Parser) registerPrefix(tokenType lexer.TokenType, function prefixParslet) {
	p.prefixParslets[tokenType] = function
}
func (p *Parser) registerInfix(tokenType lexer.TokenType, function infixParslet) {
	p.infixParslets[tokenType] = function
}
func (p *Parser) registerStatement(tokenType lexer.TokenType, function statementParslet) {
	p.statementParslets[tokenType] = function
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
