package parser

import "elara/lexer"

type FunctionArgument struct {
	Lazy    bool
	Type    Type
	Name    string
	Default Expr
}

func (p *Parser) invocationParameters(separator *TokenType) (expr []Expr) {
	params := make([]Expr, 0)
	for !p.match(lexer.RParen) {
		param := p.expression()
		params = append(params, param)
		if p.peek().TokenType == lexer.RParen {
			p.advance()
			break
		}
		if separator != nil {
			p.consume(*separator, "Expected separator "+separator.String()+" in function parameters")
		}
	}
	expr = params
	return
}

func (p *Parser) functionArguments() (args []FunctionArgument) {
	args = make([]FunctionArgument, 0)
	p.consume(lexer.LParen, "Expected left paren before starting function definition")

	for !p.match(lexer.RParen) {
		arg := p.functionArgument()
		args = append(args, arg)
	}
	return
}

func (p *Parser) functionArgument() FunctionArgument {
	lazy := p.parseProperties(lexer.Lazy)[0]
	i1 := p.consume(lexer.Identifier, "Invalid argument in function def")
	if p.match(lexer.Equal) {
		return FunctionArgument{
			Lazy:    lazy,
			Type:    nil,
			Name:    string(i1.Text),
			Default: p.expression(),
		}
	}
	id := p.consume(lexer.Identifier, "Invalid argument in function def")
	var def Expr
	if p.match(lexer.Equal) {

		def = p.expression()
	}
	typ := ElementaryTypeContract{Identifier: string(i1.Text)}
	return FunctionArgument{
		Lazy:    lazy,
		Type:    typ,
		Name:    string(id.Text),
		Default: def,
	}
}

func (p *Parser) isFuncDef() (result bool) {
	closing := p.findParenClosingPoint(p.current)
	return p.tokens[closing+1].TokenType == lexer.Arrow ||
		(p.tokens[closing+1].TokenType == lexer.Identifier && p.tokens[closing+2].TokenType == lexer.Arrow)
}

func (p *Parser) findParenClosingPoint(start int) (index int) {
	if p.tokens[start].TokenType != lexer.LParen {
		return -1
	}
	cur := start + 1
	for p.tokens[cur].TokenType != lexer.RParen {
		if p.tokens[cur].TokenType == lexer.LParen {
			cur = p.findParenClosingPoint(cur)
		}
		cur++
		if cur > len(p.tokens) {
			panic(ParseError{
				token:   p.previous(),
				message: "Unexpected end before closing parenthesis",
			})
		}
	}
	return cur
}
