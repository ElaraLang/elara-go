package parser

import (
	"github.com/ElaraLang/elara/ast"
	"github.com/ElaraLang/elara/lexer"
)

func (p *Parser) initTypeInfixParselets() {
	p.infixTypeParslets = make(map[lexer.TokenType]infixTypeParslet, 0)
	p.registerTypeInfix(lexer.TypeAnd, p.parseAlgebraicType)
	p.registerTypeInfix(lexer.TypeOr, p.parseAlgebraicType)
	p.registerTypeInfix(lexer.Colon, p.parseMapType)
}

func (p *Parser) parseAlgebraicType(left ast.Type) ast.Type {
	operator := p.Tape.ConsumeAny()
	precedence := typePrecedenceOf(operator.TokenType)
	p.Tape.skipLineBreaks()
	right := p.parseType(precedence)
	return &ast.AlgebraicType{
		Token:     operator,
		Left:      left,
		Operation: operator,
		Right:     right,
	}
}

func (p *Parser) parseMapType(keyType ast.Type) ast.Type {
	tok := p.Tape.Consume(lexer.Colon)
	precedence := typePrecedenceOf(tok.TokenType)
	valueType := p.parseType(precedence)
	return &ast.MapType{
		Token:     tok,
		KeyType:   keyType,
		ValueType: valueType,
	}
}
