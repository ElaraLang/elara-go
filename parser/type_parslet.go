package parser

import (
	"github.com/ElaraLang/elara/ast"
	"github.com/ElaraLang/elara/lexer"
)

func (p *Parser) initTypeParselets() {
}

func (p *Parser) parseFunctionType() ast.Type {
	tok := p.Tape.Consume(lexer.LParen)
	params := make([]ast.Type, 0)

	for !p.Tape.ValidateHead(lexer.RParen) {
		param := p.parseType()
		params = append(params, param)
		if !(p.Tape.Match(lexer.Comma) || p.Tape.ValidateHead(lexer.RParen)) {
			// panic
		}
	}
	p.Tape.Expect(lexer.Arrow)

	retType := p.parseType()

	return &ast.FunctionType{
		Token:      tok,
		ParamTypes: params,
		ReturnType: retType,
	}
}

func (p *Parser) parseMapType() ast.Type {
	tok := p.Tape.Consume(lexer.LBrace)
	keyType := p.parseType()
	p.Tape.Expect(lexer.Comma)
	valueType := p.parseType()
	p.Tape.Expect(lexer.RBrace)
	return &ast.MapType{
		Token:     tok,
		KeyType:   keyType,
		ValueType: valueType,
	}
}

func (p *Parser) parseCollectionType() ast.Type {
	tok := p.Tape.Consume(lexer.LSquare)
	typ := p.parseType()
	p.Tape.Expect(lexer.RSquare)
	return &ast.CollectionType{
		Token: tok,
		Type:  typ,
	}
}

func (p *Parser) parseContractualType() ast.Type {
	tok := p.Tape.Consume(lexer.Type)
	p.Tape.Expect(lexer.LBrace)
	contracts := make([]ast.Contract, 0)
	for !p.Tape.ValidateHead(lexer.RBrace) {
		contract := p.parseContract()
		contracts = append(contracts, contract)
		if !(p.Tape.Match(lexer.Comma) || p.Tape.ValidateHead(lexer.RBrace)) {
			// panic
		}
	}
	p.Tape.Expect(lexer.RBrace)
	return &ast.ContractualType{
		Token:     tok,
		Contracts: contracts,
	}
}

func (p *Parser) parseContract() ast.Contract {
	id := p.Tape.Consume(lexer.Identifier)
	typ := p.parseType()
	return ast.Contract{
		Token: id,
		Identifier: ast.Identifier{
			Token: id,
			Name:  string(id.Text),
		},
		Type: typ,
	}
}

func (p *Parser) parsePrimaryType() ast.Type {
	idTok := p.Tape.Consume(lexer.Identifier)
	id := ast.Identifier{
		Token: idTok,
		Name:  string(idTok.Text),
	}
	return &ast.PrimaryType{
		Token:      idTok,
		Identifier: id,
	}
}