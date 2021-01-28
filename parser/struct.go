package parser

import (
	"github.com/ElaraLang/elara/ast"
	"github.com/ElaraLang/elara/lexer"
)

func (p *Parser) parseStructFields() *[]ast.StructField {
	fields := make([]ast.StructField, 0)
	p.Tape.skipLineBreaks()
	for !p.Tape.ValidateHead(lexer.RBrace) {
		var typ ast.Type
		var id ast.IdentifierLiteral
		var def ast.Expression
		prop := p.Tape.MatchUnorderedSequence(lexer.Mut, lexer.Open)
		if p.Tape.ValidationPeek(1, lexer.Equal) {
			id = *p.parseIdentifier().(*ast.IdentifierLiteral)
		} else {
			typ = p.parseType(TypeLowest)
			id = *p.parseIdentifier().(*ast.IdentifierLiteral)
		}

		if p.Tape.Match(lexer.Equal) {
			def = p.parseExpression(Lowest)
		}

		if typ == nil && def == nil {
			p.error(id.Token, "Neither Type nor Default provided in struct field!")
		}

		fields = append(fields, ast.StructField{
			Mutable:    prop[lexer.Mut],
			Open:       prop[lexer.Open],
			Type:       typ,
			Identifier: id,
			Default:    def,
		})

		p.Tape.skipLineBreaks()
	}
	return &fields
}
