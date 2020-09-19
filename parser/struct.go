package parser

import "elara/lexer"

type StructField struct {
	Mutable    bool
	Identifier string
	FieldType  *Type
	Default    *Expr
}

func (p *Parser) structFields() (fields []StructField, err error) {
	tok, error := p.consume(lexer.LBrace, "Expected '{' at struct field start")
	if error != nil {
		return nil, error
	}
	fields = make([]StructField, 0)
	p.cleanNewLines()
	for !p.check(lexer.RBrace) {
		field, error := p.structField()
		if error != nil {
			return nil, error
		}
		fields = append(fields, *field)
		if !p.match(lexer.NEWLINE) && !p.check(lexer.RBrace) {
			return nil, ParseError{
				token:   tok,
				message: "Expected newline after struct field",
			}
		}
	}
	_, error = p.consume(lexer.RBrace, "Expected '}' at struct def end")
	if error != nil {
		return nil, error
	}
	return
}

func (p *Parser) structField() (field *StructField, err error) {
	mutable := p.match(lexer.Mut)
	t1 := p.advance()
	t2 := p.advance()
	var typ Type
	var identifier string
	var def Expr
	if t1.TokenType == lexer.Identifier {
		switch t2.TokenType {
		case lexer.Identifier:
			typ = ElementaryTypeContract{Identifier: identifier}
			identifier = t2.Text
			if p.match(lexer.Equal) {
				def, err = p.logicalOr()
			}
			break
		case lexer.Equal:
			identifier = t2.Text
			def, err = p.logicalOr()
			break
		default:
			err = ParseError{
				token:   t1,
				message: "Invalid struct field",
			}
			break
		}
	} else {
		err = ParseError{
			token:   t1,
			message: "Invalid struct field",
		}
	}
	if err != nil {
		return nil, err
	}
	return &StructField{
		Mutable:    mutable,
		Identifier: identifier,
		FieldType:  &typ,
		Default:    &def,
	}, nil
}
