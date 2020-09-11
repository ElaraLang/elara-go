package parser

import "elara/lexer"

type Type interface {
	typeOf()
}

type ElementaryTypeContract struct {
	Identifier string
}
type BinaryTypeContract struct {
	lhs    Type
	typeOp TokenType
	rhs    Type
}

func (p *Parser) typeContract() (contract Type, err error) {
	return p.contractualOr()
}

func (p *Parser) contractualOr() (contract Type, err error) {
	contract, err = p.contractualAnd()
	if err != nil {
		return
	}
	for p.match(lexer.Or) {
		op := p.previous()
		rhs, error := p.contractualAnd()
		if error != nil {
			err = error
			return
		}
		contract = BinaryTypeContract{
			lhs:    contract,
			typeOp: op.TokenType,
			rhs:    rhs,
		}
	}
	return
}

func (p *Parser) contractualAnd() (contract Type, err error) {
	contract, err = p.primaryContract()
	if err != nil {
		return
	}
	for p.match(lexer.And) {
		op := p.previous()
		rhs, error := p.primaryContract()
		if error != nil {
			err = error
			return
		}
		contract = BinaryTypeContract{
			lhs:    contract,
			typeOp: op.TokenType,
			rhs:    rhs,
		}
	}
	return
}

func (p *Parser) primaryContract() (contract Type, err error) {
	if p.peek().TokenType == lexer.Identifier {
		return ElementaryTypeContract{Identifier: p.advance().Text}, nil
	}
	return p.contractualGroup()
}

func (p *Parser) contractualGroup() (contract Type, err error) {
	if p.match(lexer.LParen) {
		contract, err = p.contractualOr()
		if err != nil {
			return nil, err
		}
		_, err = p.consume(lexer.RBrace, "contract group not closed. Expected '}'")
		if err != nil {
			return nil, err
		}
		return
	}
	err = ParseError{
		token:   p.previous(),
		message: "Invalid type contract",
	}
	return nil, err
}

func (t ElementaryTypeContract) typeOf() {}
func (t BinaryTypeContract) typeOf()     {}
