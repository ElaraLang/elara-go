package parser

import "github.com/ElaraLang/elara/lexer"

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

func (p *Parser) typeContract() (contract Type) {
	return p.contractualOr(false)
}

func (p *Parser) typeContractDefinable() (contract Type) {
	return p.contractualOr(true)
}

func (p *Parser) contractualOr(allowDef bool) (contract Type) {
	contract = p.contractualAnd(allowDef)
	for p.match(lexer.Or) {
		op := p.previous()
		rhs := p.contractualAnd(allowDef)
		contract = BinaryTypeContract{
			Lhs:    contract,
			TypeOp: op.TokenType,
			Rhs:    rhs,
		}
	}
	return
}

func (p *Parser) contractualAnd(allowDef bool) (contract Type) {
	contract = p.primaryContract(allowDef)

	for p.match(lexer.And) {
		op := p.previous()
		rhs := p.primaryContract(allowDef)

		contract = BinaryTypeContract{
			Lhs:    contract,
			TypeOp: op.TokenType,
			Rhs:    rhs,
		}
	}
	return
}

func (p *Parser) primaryContract(allowDef bool) (contract Type) {

	if p.peek().TokenType == lexer.Identifier {
		return ElementaryTypeContract{Identifier: string(p.advance().Text)}
	} else if p.match(lexer.LParen) {
		isFunc := p.isFuncDef()
		if isFunc {
			args := make([]Type, 0)
			for !p.check(lexer.RParen) {
				argTyp := p.typeContract()
				args = append(args, argTyp)
				if !p.match(lexer.Comma) {
					break
				}
			}
			p.consume(lexer.RParen, "Function type args not ended properly with ')'")
			p.consume(lexer.Arrow, "Expected arrow after function type args")

			ret := p.typeContract()
			return InvocableTypeContract{
				Args:       args,
				ReturnType: ret,
			}
		} else {
			contract = p.contractualOr(allowDef)

			p.consume(lexer.RBrace, "contract group not closed. Expected '}'")
			return
		}
	}
	return p.definedContract(allowDef)
}

func (p *Parser) definedContract(allowDef bool) (contract Type) {
	if allowDef && p.check(lexer.LBrace) {
		defTyp := p.definedTypes()
		return DefinedTypeContract{DefType: defTyp}
	}
	panic(ParseError{
		token:   p.previous(),
		message: "Invalid type contract",
	})
}

func (t ElementaryTypeContract) typeOf() {}
func (t BinaryTypeContract) typeOf()     {}
func (t InvocableTypeContract) typeOf()  {}
func (t DefinedTypeContract) typeOf()    {}
