package parser

import (
	"github.com/ElaraLang/elara/ast"
)

func (p *Parser) initTypeInfixParselets() {
}

func (p *Parser) parseAlgebraicType(left ast.Type) ast.Type {
	operator := p.Tape.ConsumeAny()
	precedence := typePrecedenceOf(operator.TokenType)
	right := p.parseType(precedence)
	return &ast.AlgebraicType{
		Token:     operator,
		Left:      left,
		Operation: operator,
		Right:     right,
	}
}
