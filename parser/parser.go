package parser

import (
	"github.com/ElaraLang/elara/ast"
	"github.com/ElaraLang/elara/lexer"
)

type Parser struct {
	OutputChannel      chan ast.Statement
	ErrorChannel       chan ParseError
	Tape               *TokenTape
	statementParslets  map[lexer.TokenType]statementParslet
	prefixParslets     map[lexer.TokenType]prefixParslet
	infixParslets      map[lexer.TokenType]infixParslet
	prefixTypeParslets map[lexer.TokenType]prefixTypeParslet
	infixTypeParslets  map[lexer.TokenType]infixTypeParslet

	fileName string
}

func NewParser(inputChannel chan lexer.Token, outputChannel chan ast.Statement, errorChannel chan ParseError) Parser {
	tape := NewTokenTape(inputChannel)
	p := Parser{OutputChannel: outputChannel, ErrorChannel: errorChannel, Tape: &tape}
	p.initPrefixParselets()
	p.initInfixParselets()
	p.initStatementParselets()
	p.initTypePrefixParselets()
	p.initTypeInfixParselets()
	return p
}

func (p *Parser) Parse(fileName string) {
	p.fileName = fileName
	if p.Tape.IsClosed() {
		p.Tape.Unwind()
	}
	for !p.Tape.ValidateHead(lexer.EOF) {
		p.parseSafely()
	}
}

func (p *Parser) parseSafely() {
	defer p.handleParseError()
	p.Tape.skipLineBreaks()

	p.OutputChannel <- p.parseStatement()

	if !p.Tape.Match(lexer.NEWLINE, lexer.EOF) &&
		len(p.Tape.tokens) > 0 &&
		!p.Tape.ValidationPeek(-1, lexer.NEWLINE) {
		p.error(p.Tape.Current(), expect("NEWLINE", "end of statement"))
	}
}

type (
	statementParslet  func() ast.Statement
	prefixParslet     func() ast.Expression
	infixParslet      func(ast.Expression) ast.Expression
	prefixTypeParslet func() ast.Type
	infixTypeParslet  func(ast.Type) ast.Type
)

func (p *Parser) registerPrefix(tokenType lexer.TokenType, function prefixParslet) {
	p.prefixParslets[tokenType] = function
}
func (p *Parser) registerInfix(tokenType lexer.TokenType, function infixParslet) {
	p.infixParslets[tokenType] = function
}

func (p *Parser) registerTypePrefix(tokenType lexer.TokenType, function prefixTypeParslet) {
	p.prefixTypeParslets[tokenType] = function
}
func (p *Parser) registerTypeInfix(tokenType lexer.TokenType, function infixTypeParslet) {
	p.infixTypeParslets[tokenType] = function
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
		p.error(p.Tape.Current(), "Invalid Token found at expression")
		return nil
	}
	expr := parsePrefix()
	for !p.Tape.ValidationPeek(0, lexer.NEWLINE, lexer.EOF) && precedence < precedenceOf(p.Tape.Current().TokenType) {
		p.Tape.skipLineBreaks()
		infix := p.infixParslets[p.Tape.Current().TokenType]
		if infix == nil {
			return expr
		}
		expr = infix(expr)
	}
	return expr
}

func (p *Parser) parseType(precedence int) ast.Type {
	parsePrefixType := p.prefixTypeParslets[p.Tape.Current().TokenType]
	if parsePrefixType == nil {
		p.error(p.Tape.Current(), "Invalid Token found at type")
		return nil
	}
	typ := parsePrefixType()
	for !p.Tape.ValidationPeek(0, lexer.NEWLINE) && precedence < typePrecedenceOf(p.Tape.Current().TokenType) {
		infix := p.infixTypeParslets[p.Tape.Current().TokenType]
		if infix == nil {
			return typ
		}
		typ = infix(typ)
	}
	return typ
}
