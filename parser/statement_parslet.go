package parser

import (
	"github.com/ElaraLang/elara/ast"
	"github.com/ElaraLang/elara/lexer"
)

func (p *Parser) initStatementParselets() {
	p.statementParslets = make(map[lexer.TokenType]statementParslet, 0)
	p.registerStatement(lexer.Let, p.parseLetStatement)

}

func (p *Parser) parseLetStatement() ast.Statement {
	token := p.Tape.Consume(lexer.Let)
	prop := p.Tape.MatchInorderedSequence(lexer.Mut, lexer.Lazy, lexer.Restricted)
	id := p.parseIdentifier()
	var varType ast.Type
	var value ast.Expression
	if p.Tape.ValidationPeek(0, lexer.LParen) {
		value = p.parseExpression(LOWEST)
	} else {
		if p.Tape.ValidationPeek(0, lexer.Colon) {
			varType = p.parseType()
		}
		p.Tape.Expect(lexer.Equal)
		value = p.parseExpression(LOWEST)
	}
	return &ast.DeclarationStatement{
		Token:      token,
		Mutable:    prop[lexer.Mut],
		Lazy:       prop[lexer.Lazy],
		Open:       prop[lexer.Restricted], // TODO:: Introduce OPEN token to lexer
		Identifier: id,
		Type:       varType,
		Value:      value,
	}
}

func (p *Parser) parseWhileStatement() ast.Statement {
	token := p.Tape.Consume(lexer.While)
	condition := p.parseExpression(LOWEST)
	var body ast.Statement
	if p.Tape.Match(lexer.Arrow) {
		body = p.parseExpressionStatement()
	} else {
		body = p.parseBlockStatement()
	}
	return &ast.WhileStatement{
		Token:     token,
		Condition: condition,
		Body:      body,
	}
}

func (p *Parser) parseExtendStatement() ast.Statement {
	token := p.Tape.Consume(lexer.Extend)
	id := p.parseIdentifier()
	var alias ast.Identifier
	if p.Tape.Match(lexer.As) {
		alias = p.parseIdentifier()
	} else {
		alias = ast.Identifier{
			Token: token,
			Name:  "this",
		}
	}
	body := p.parseBlockStatement()
	return &ast.ExtendStatement{
		Token:      token,
		Identifier: id,
		Alias:      alias,
		Body:       body,
	}
}

func (p *Parser) parseBlockStatement() ast.Statement {
	token := p.Tape.Consume(lexer.LBrace)
	block := make([]ast.Statement, 0)
	for p.Tape.Match(lexer.RBrace) {
		block = append(block, p.parseStatement())
	}
	return &ast.BlockStatement{
		Token: token,
		Block: block,
	}
}

func (p *Parser) parseExpressionStatement() ast.Statement {
	return &ast.ExpressionStatement{
		Token:      p.Tape.Current(),
		Expression: p.parseExpression(LOWEST),
	}
}
