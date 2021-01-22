package parser

import (
	"github.com/ElaraLang/elara/ast"
	"github.com/ElaraLang/elara/lexer"
)

func (p *Parser) parseCollectionElements() []ast.Expression {
	elements := make([]ast.Expression, 0)
	for !p.Tape.ValidationPeek(0, lexer.RSquare) {
		if len(elements) > 0 {
			p.Tape.Expect(lexer.Comma)
		}
		expr := p.parseExpression()
		elements = append(elements, expr)
	}
	return elements
}
