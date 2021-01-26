package parser

import (
	"github.com/ElaraLang/elara/ast"
	"github.com/ElaraLang/elara/lexer"
)

func (p *Parser) initStatementParselets() {
	p.statementParslets = make(map[lexer.TokenType]statementParslet, 0)
	p.registerStatement(lexer.Let, p.parseLetStatement)
	p.registerStatement(lexer.While, p.parseWhileStatement)
	p.registerStatement(lexer.Return, p.parseReturnStatement)
	p.registerStatement(lexer.Extend, p.parseExtendStatement)
	p.registerStatement(lexer.Type, p.parseTypeStatement)
	p.registerStatement(lexer.Hash, p.parseGenerifiedStatement)
	p.registerStatement(lexer.Namespace, p.parseNamespace)
	p.registerStatement(lexer.Import, p.parseImport)
}

func (p *Parser) parseLetStatement() ast.Statement {
	token := p.Tape.Consume(lexer.Let)
	prop := p.Tape.MatchInorderedSequence(lexer.Mut, lexer.Lazy, lexer.Open)
	id := p.parseIdentifier()
	var varType ast.Type
	var value ast.Expression
	if p.Tape.ValidationPeek(0, lexer.LParen) {
		value = p.parseExpression(Lowest)
	} else {
		if p.Tape.ValidationPeek(0, lexer.Colon) {
			varType = p.parseType(TypeLowest)
		}
		p.Tape.Expect(lexer.Equal)
		p.Tape.skipLineBreaks()
		value = p.parseExpression(Lowest)
	}
	return &ast.DeclarationStatement{
		Token:      token,
		Mutable:    prop[lexer.Mut],
		Lazy:       prop[lexer.Lazy],
		Open:       prop[lexer.Open], // TODO:: Introduce OPEN token to lexer
		Identifier: id,
		Type:       varType,
		Value:      value,
	}
}

func (p *Parser) parseWhileStatement() ast.Statement {
	token := p.Tape.Consume(lexer.While)
	condition := p.parseExpression(Lowest)
	var body ast.Statement
	if p.Tape.Match(lexer.Arrow) {
		body = p.parseExpressionStatement()
	} else {
		p.Tape.skipLineBreaks()
		body = p.parseBlockStatement()
	}
	return &ast.WhileStatement{
		Token:     token,
		Condition: condition,
		Body:      body,
	}
}

func (p *Parser) parseReturnStatement() ast.Statement {
	token := p.Tape.Consume(lexer.Return)
	p.Tape.skipLineBreaks()
	value := p.parseExpression(Lowest)
	return &ast.ReturnStatement{
		Token: token,
		Value: value,
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
	p.Tape.skipLineBreaks()
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
	p.Tape.skipLineBreaks()
	block := make([]ast.Statement, 0)
	for !p.Tape.Match(lexer.RBrace) {
		block = append(block, p.parseStatement())
		p.Tape.skipLineBreaks()
	}
	return &ast.BlockStatement{
		Token: token,
		Block: block,
	}
}

func (p *Parser) parseExpressionStatement() ast.Statement {
	return &ast.ExpressionStatement{
		Token:      p.Tape.Current(),
		Expression: p.parseExpression(Lowest),
	}
}

func (p *Parser) parseTypeStatement() ast.Statement {
	token := p.Tape.Consume(lexer.Type)
	id := p.Tape.Consume(lexer.Identifier)
	p.Tape.Expect(lexer.Equal)
	var internalAlias string
	if p.Tape.ValidationPeek(1, lexer.Where) {
		internalAlias = string(p.Tape.Consume(lexer.Identifier).Data)
		p.Tape.Expect(lexer.Where)
	} else {
		internalAlias = string(id.Data)
	}
	p.Tape.skipLineBreaks()
	typeContract := p.parseType(TypeLowest)
	return &ast.TypeStatement{
		Token: token,
		Identifier: ast.Identifier{
			Token: id,
			Name:  string(id.Data),
		},
		InternalId: ast.Identifier{
			Token: id,
			Name:  internalAlias,
		},
		Contract: typeContract,
	}
}

func (p *Parser) parseGenerifiedStatement() ast.Statement {
	token := p.Tape.Consume(lexer.Hash)
	id := p.Tape.Consume(lexer.Identifier)
	p.Tape.Expect(lexer.Where)
	typeContract := p.parseType(TypeLowest)
	p.Tape.skipLineBreaks()
	stmt := p.parseStatement()
	return &ast.GenerifiedStatement{
		Token: token,
		Contract: ast.NamedContract{
			Token: token,
			Identifier: ast.Identifier{
				Token: id,
				Name:  string(id.Data),
			},
			Type: typeContract,
		},
		Statement: stmt,
	}
}

func (p *Parser) parseNamespace() ast.Statement {
	tok := p.Tape.Consume(lexer.Namespace)
	mod := p.parseModule()
	return &ast.NamespaceStatement{
		Token:  tok,
		Module: *mod,
	}
}

func (p *Parser) parseImport() ast.Statement {
	tok := p.Tape.Consume(lexer.Import)
	mod := p.parseModule()
	return &ast.ImportStatement{
		Token:  tok,
		Module: *mod,
	}
}

func (p *Parser) parseModule() *ast.Module {
	baseTok := p.parseIdentifier()
	idSlice := make([]ast.Identifier, 1)
	idSlice[0] = baseTok
	mod := baseTok.Name
	for p.Tape.Match(lexer.Slash) {
		subId := p.parseIdentifier()
		mod += "/" + subId.Name
		idSlice = append(idSlice, subId)
	}
	return &ast.Module{
		Pkg:    mod,
		PkgIds: idSlice,
	}
}
