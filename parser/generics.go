package parser

import (
	"elara/lexer"
)

type GenericContract struct {
	Identifier string
	Contract   Type
}

func (p *Parser) generic() (contracts []GenericContract, err error) {
	_, err = p.consume(lexer.LAngle, "Expected generic declaration to start with `<`")
	if err != nil {
		return
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
	_, err = p.consume(lexer.RAngle, "Expected generic declaration to end with `>`")
	if err != nil {
		return
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

func (p *Parser) typeStatement() (typStmt Stmt, err error) {
	_, err = p.consume(lexer.Type, "Expected 'type' at the start of type declaration")
	if err != nil {
		return
	}

	id, error := p.consume(lexer.Identifier, "Expected identifier for type")

	if error != nil {
		err = error
		return
	}

	_, err = p.consume(lexer.Arrow, "Expected arrow after type identifier")

	if err != nil {
		return
	}

	contract, error := p.typeContractDefinable()
	if error != nil {
		err = error
		return
	}
	typStmt = TypeStmt{
		Identifier: id.Text,
		Contract:   contract,
	}
	return
}

func (p *Parser) genericStatement() (genericStmt Stmt, err error) {
	generic, error := p.generic()
	if error != nil {
		return nil, error
	}
	p.cleanNewLines()

	stmt, error := p.declaration()
	return GenerifiedStmt{
		Contracts: generic,
		Statement: stmt,
	}, nil
}
