package parser

import (
	"elara/lexer"
)

type GenericContract struct {
	Identifier string
}


func (p *Parser) generic()  {
	_, error = p.consume(lexer.LAngle, "Expected generic declaration to start with `<`")

}

func (p *Parser) genericType() (string, err error)  {
	typID, error := p.consume(lexer.Identifier, "Expected identifier for generic type")
	if (error != nil) {
		return nil, error
	}
	_, error = p.consume(lexer.Colon, "Expected colon after generic type id")
	if (error != nil) {
		return nil, error
	}

	var typeContracts []Type
	if (p.check(lexer.Identifier))
}
