package parser

import "github.com/ElaraLang/elara/lexer"

type Type interface {
	typeOf()
}

type DefinedTypeContract struct {
	DefType []DefinedType
	Name    string
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

type CollectionTypeContract struct {
	ElemType Type
}

type MapTypeContract struct {
	KeyType   Type
	ValueType Type
}

func (p *Parser) typeContract() (contract Type) {
	return p.contractualOr(false)
}

func (p *Parser) typeContractDefinable() (contract Type) {
	return p.contractualOr(true)
}

func (p *Parser) contractualOr(allowDef bool) (contract Type) {
	contract = p.contractualAnd(allowDef)
	for p.match(lexer.TypeOr) {
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

	for p.match(lexer.TypeAnd) {
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
	if p.peek().TokenType == lexer.LSquare {
		p.advance()
		collType := p.consume(lexer.Identifier, "Expected identifier after [ for collection type")
		p.consume(lexer.RSquare, "Expected ] after [ for collection type")
		return CollectionTypeContract{ElemType: ElementaryTypeContract{Identifier: string(collType.Text)}}
	}
	if p.peek().TokenType == lexer.Identifier {
		name := string(p.advance().Text)
		return ElementaryTypeContract{Identifier: name}
	} else if p.check(lexer.LParen) {
		isFunc := p.isFuncDef()
		if isFunc {
			args := make([]Type, 0)
			p.consume(lexer.LParen, "Function type args not started properly with '('")

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

			p.consume(lexer.LParen, "contract group not closed. Expected '}'")
			return
		}
	}
	if p.peek().TokenType == lexer.LBrace {
		p.advance()
		//Peek until reaching a closing brace
		count := 0
		seenColon := false
		for {
			count++
			next := p.advance().TokenType
			if next == lexer.Colon {
				seenColon = true
			}
			if next == lexer.RBrace {
				break
			}
		}
		for i := 0; i < count; i++ {
			p.reverse()
		}
		if !seenColon { //it's not a map type
			p.reverse() // Undo the lbrace read
			return p.definedContract(allowDef)
		}

		keyType := p.typeContract()
		p.consume(lexer.Colon, "Expected colon in map type")
		valueType := p.typeContract()

		p.consume(lexer.RBrace, "Expected closing brace for map type contract")

		return MapTypeContract{
			KeyType:   keyType,
			ValueType: valueType,
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
func (t CollectionTypeContract) typeOf() {}
func (t MapTypeContract) typeOf()        {}
