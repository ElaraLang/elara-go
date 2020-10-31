package parser

import "elara/lexer"

type StructField struct {
	Mutable    bool
	Identifier string
	FieldType  *Type
	Default    *Expr
}

func (p *Parser) structFields() (fields []StructField) {
	p.consume(lexer.LBrace, "Expected '{' at struct field start")

	fields = make([]StructField, 0)
	p.cleanNewLines()
	for !p.check(lexer.RBrace) {
		field := p.structField()
		fields = append(fields, *field)
		if !p.match(lexer.NEWLINE) && !p.check(lexer.RBrace) {
			panic(ParseError{
				token:   p.previous(),
				message: "Expected newline after struct field",
			})
		}
	}
	p.consume(lexer.RBrace, "Expected '}' at struct def end")
	return
}

func (p *Parser) structField() (field *StructField) {
	mutable := p.match(lexer.Mut)
	t1 := p.advance()
	t2 := p.advance()
	var typ Type
	var identifier string
	var def Expr
	if t1.TokenType == lexer.Identifier {
		switch t2.TokenType {
		case lexer.Identifier:
			typ = ElementaryTypeContract{Identifier: string(t1.Text)}
			identifier = string(t2.Text)
			if p.match(lexer.Equal) {
				def = p.logicalOr()
			}
			break
		case lexer.Equal:
			identifier = string(t2.Text)
			def = p.logicalOr()
			break
		default:
			panic(ParseError{
				token:   t1,
				message: "Invalid struct field",
			})
		}
	} else {
		panic(ParseError{
			token:   t1,
			message: "Invalid struct field",
		})
	}
	return &StructField{
		Mutable:    mutable,
		Identifier: identifier,
		FieldType:  &typ,
		Default:    &def,
	}
}
