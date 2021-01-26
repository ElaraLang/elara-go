package parser

import (
	"github.com/ElaraLang/elara/ast"
	"github.com/ElaraLang/elara/lexer"
)

func (p *Parser) parseFunctionParameters() []ast.Parameter {
	params := make([]ast.Parameter, 0)
	if p.Tape.ValidationPeek(0, lexer.RParen) {
		return params
	}
outer:
	for {
		if p.Tape.ValidationPeek(0, lexer.Identifier) &&
			p.Tape.ValidationPeek(1, lexer.Comma) {
			id := p.Tape.Consume(lexer.Identifier)
			param := ast.Parameter{
				Type:       nil,
				Identifier: ast.Identifier{Token: id, Name: id.String()},
			}
			params = append(params, param)
			p.Tape.Expect(lexer.Comma)
			p.Tape.skipLineBreaks()
			continue
		}
		typ := p.parseType(TypeLowest)
		id := p.Tape.Consume(lexer.Identifier)
		param := ast.Parameter{
			Type:       typ,
			Identifier: ast.Identifier{Token: id, Name: id.String()},
		}
		params = append(params, param)
		switch p.Tape.Current().TokenType {
		case lexer.RParen:
			break outer
		case lexer.Comma:
			p.Tape.advance()
			p.Tape.skipLineBreaks()
		default:
			p.error(id, "Parameter separator missing!")
		}
	}
	return params
}

func (p *Parser) parseFunctionCallArguments() []ast.Expression {
	args := make([]ast.Expression, 0)
	if p.Tape.ValidationPeek(0, lexer.RParen) {
		return args
	}
	args = append(args, p.parseExpression(Lowest))
	for p.Tape.Match(lexer.Comma) {
		p.Tape.skipLineBreaks()
		args = append(args, p.parseExpression(Lowest))
	}
	return args
}
