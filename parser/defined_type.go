package parser

import "elara/lexer"

type DefinedType struct {
	Identifier string
	DefType    Type
}

func (p *Parser) definedTypes() (types []DefinedType, err error) {
	types = make([]DefinedType, 0)
	_, err = p.consume(lexer.LBrace, "Expected '{' where defined type starts")
	if err != nil {
		return nil, err
	}

	for !p.check(lexer.RBrace) {
		id, error := p.consume(lexer.Identifier, "Expected identifier for type in defined type contract")
		if error != nil {
			return nil, error
		}
		typ, error := p.primaryContract(true)
		dTyp := DefinedType{
			Identifier: id.Text,
			DefType:    typ,
		}
		types = append(types, dTyp)
		if !p.match(lexer.Comma) {
			break
		}
	}

	_, err = p.consume(lexer.RBrace, "Expected '}' where defined type ends")
	if err != nil {
		return nil, err
	}
	return types, nil
}
