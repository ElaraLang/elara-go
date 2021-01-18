package parserlegacy

import (
	"github.com/ElaraLang/elara/lexer"
)

type GenericContract struct {
	Identifier string
	Contract   Type
}

func (p *Parser) generic() (contracts []GenericContract) {
	p.consume(lexer.LAngle, "Expected generic declaration to start with `<`")
	contracts = make([]GenericContract, 0)
	for {
		contract := p.genericContract()
		contracts = append(contracts, contract)
		if !p.match(lexer.Comma) {
			break
		}
	}
	p.consume(lexer.RAngle, "Expected generic declaration to end with `>`")
	return
}

func (p *Parser) genericContract() (typContract GenericContract) {
	typID := p.consume(lexer.Identifier, "Expected identifier for generic type")
	p.consume(lexer.Colon, "Expected colon after generic type id")
	contract := p.typeContractDefinable()
	typContract = GenericContract{
		Identifier: string(typID.Text),
		Contract:   contract,
	}
	return
}

func (p *Parser) typeStatement() (typStmt Stmt) {
	p.consume(lexer.Type, "Expected 'type' at the start of type declaration")
	id := p.consume(lexer.Identifier, "Expected identifier for type")
	p.consume(lexer.Equal, "Expected equals after type identifier")
	contract := p.typeContractDefinable()
	typStmt = TypeStmt{
		Identifier: string(id.Text),
		Contract:   contract,
	}
	return
}

func (p *Parser) genericStatement() (genericStmt Stmt) {
	generic := p.generic()

	p.cleanNewLines()

	stmt := p.declaration()
	return GenerifiedStmt{
		Contracts: generic,
		Statement: stmt,
	}
}
