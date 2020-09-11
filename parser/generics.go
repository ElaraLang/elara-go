package parser

import (
	"elara/lexer"
)

type GenericContract struct {
	Identifier string
	Contract   Type
}

func (p *Parser) generic() (contracts []GenericContract, err error) {
	_, error := p.consume(lexer.LAngle, "Expected generic declaration to start with `<`")
	if error != nil {
		return nil, error
	}
	contracts = make([]GenericContract, 0)
	for {
		contract, error := p.genericContract()
		if error != nil {
			return nil, error
		}
		contracts = append(contracts, contract)
		if !p.match(lexer.Comma) {
			break
		}
	}
	return
}

func (p *Parser) genericContract() (typContract GenericContract, err error) {
	typID, error := p.consume(lexer.Identifier, "Expected identifier for generic type")
	if error != nil {
		err = error
		return
	}

	_, error = p.consume(lexer.Colon, "Expected colon after generic type id")

	if error != nil {
		err = error
		return
	}

	contract, error := p.typeContractDefinable()
	if error != nil {
		err = error
		return
	}
	typContract = GenericContract{
		Identifier: typID.Text,
		Contract:   contract,
	}
	return
}
