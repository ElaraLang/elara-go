package parser

import "github.com/ElaraLang/elara/lexer"

type DefinedType struct {
	Identifier string
	DefType    Type
}

func (p *Parser) definedTypes() (types []DefinedType) {
	types = make([]DefinedType, 0)
	p.consume(lexer.LBrace, "Expected '{' where defined type starts")
	p.cleanNewLines()
	for !p.check(lexer.RBrace) {
		id := p.consume(lexer.Identifier, "Expected identifier for type in defined type contract")
		typ := p.primaryContract(true)
		dTyp := DefinedType{
			Identifier: string(id.Text),
			DefType:    typ,
		}
		types = append(types, dTyp)
		if !p.match(lexer.Comma) {
			break
		}
		p.cleanNewLines()
	}
	p.cleanNewLines()
	p.consume(lexer.RBrace, "Expected '}' where defined type ends")
	return types
}
