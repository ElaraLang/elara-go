package parser

import "elara/lexer"

type Type interface {
	typeOf()
}

type DefinedTypeContract struct {
	DefType []DefinedType
}
type ElementaryTypeContract struct {
	Identifier string
}
type InvocableTypeContract struct {
	Args       []Type
	ReturnType Type
}

type BinaryTypeContract struct {
	Lhs    Type
	TypeOp TokenType
	Rhs    Type
}

func (p *Parser) typeContract() (contract Type, err error) {
	return p.contractualOr(false)
}

func (p *Parser) typeContractDefinable() (contract Type, err error) {
	return p.contractualOr(true)
}

func (p *Parser) contractualOr(allowDef bool) (contract Type, err error) {
	contract, err = p.contractualAnd(allowDef)
	if err != nil {
		return
	}
	for p.match(lexer.Or) {
		op := p.previous()
		rhs, error := p.contractualAnd(allowDef)
		if error != nil {
			err = error
			return
		}
		contract = BinaryTypeContract{
			Lhs:    contract,
			TypeOp: op.TokenType,
			Rhs:    rhs,
		}
	}
	return
}

func (p *Parser) contractualAnd(allowDef bool) (contract Type, err error) {
	contract, err = p.primaryContract(allowDef)
	if err != nil {
		return
	}
	for p.match(lexer.And) {
		op := p.previous()
		rhs, error := p.primaryContract(allowDef)
		if error != nil {
			err = error
			return
		}
		contract = BinaryTypeContract{
			Lhs:    contract,
			TypeOp: op.TokenType,
			Rhs:    rhs,
		}
	}
	return
}

func (p *Parser) primaryContract(allowDef bool) (contract Type, err error) {

	if p.peek().TokenType == lexer.Identifier {
		return ElementaryTypeContract{Identifier: p.advance().Text}, nil
	} else if p.match(lexer.LParen) {
		isfunc, error := p.isFuncDef()
		if error != nil {
			return nil, error
		}
		if isfunc {
			args := make([]Type, 0)
			for !p.check(lexer.RParen) {
				argTyp, error := p.typeContract()
				if error != nil {
					err = error
					return
				}
				args = append(args, argTyp)
				if !p.match(lexer.Comma) {
					break
				}
			}
			_, err = p.consume(lexer.RParen, "Function type args not ended properly with ')'")
			if err != nil {
				return
			}
			_, err = p.consume(lexer.Arrow, "Expected arrow after function type args")
			if err != nil {
				return
			}

			ret, error := p.typeContract()
			if error != nil {
				err = error
				return
			}
			return InvocableTypeContract{
				Args:       args,
				ReturnType: ret,
			}, nil
		} else {
			contract, err = p.contractualOr(allowDef)
			if err != nil {
				return nil, err
			}
			_, err = p.consume(lexer.RBrace, "contract group not closed. Expected '}'")
			if err != nil {
				return nil, err
			}
			return
		}
	}
	return p.definedContract(allowDef)
}

func (p *Parser) definedContract(allowDef bool) (contract Type, err error) {
	if allowDef && p.check(lexer.LBrace) {

		defTyp, error := p.definedTypes()

		if error != nil {
			return nil, error
		}

		return DefinedTypeContract{DefType: defTyp}, nil
	}
	err = ParseError{
		token:   p.previous(),
		message: "Invalid type contract",
	}
	return nil, err
}

func (t ElementaryTypeContract) typeOf() {}
func (t BinaryTypeContract) typeOf()     {}
func (t InvocableTypeContract) typeOf()  {}
func (t DefinedTypeContract) typeOf()    {}
