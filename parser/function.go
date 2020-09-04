package parser

import "elara/lexer"

type FunctionArgument struct {
	Type     *Type
	Variable VariableExpr
	Default  *Expr
}

func (p *Parser) invocationParameters(separator *TokenType) (expr []Expr, err error) {
	params := make([]Expr, 0)
	for !p.match(lexer.RParen) {
		param, err := p.expression()
		if err != nil {
			return
		}
		params = append(params, param)
		if separator != nil {
			_, err = p.consume(*separator, "Expected separator in function parameters")
			if err != nil {
				return
			}
		}
	}
	expr = params
	return
}

func (p *Parser) functionArguments() (args []FunctionArgument, err error) {
	args = make([]FunctionArgument, 0)
	_, err = p.consume(lexer.LParen, "Expected left paren before starting function definition")
	if err != nil {
		return
	}
	for !p.match(lexer.RParen) {
		arg, err := p.functionArgument()
		if err != nil {
			return
		}
		args = append(args, arg)
	}
	return
}

func (p *Parser) functionArgument() (arg FunctionArgument, err error) {
	i1, err := p.consume(lexer.Identifier, "Invalid argument in function def")
	if err != nil {
		return
	}
	if p.match(lexer.Equal) {
		expr, err := p.expression()
		if err != nil {
			return
		}
		arg = FunctionArgument{
			Type:     nil,
			Variable: VariableExpr{Identifier: i1.Text},
			Default:  &expr,
		}
		return
	}
	id, err := p.consume(lexer.Identifier, "Invalid argument in function def")
	if err != nil {
		return
	}

	var def *Expr

	if p.match(lexer.Equal) {
		expr, err := p.expression()
		if err != nil {
			return
		}
		def = &expr
	}

	typ := Type(i1.Text)

	arg = FunctionArgument{
		Type:     &typ,
		Variable: VariableExpr{Identifier: id.Text},
		Default:  def,
	}

	return
}

func (p *Parser) isFuncDef() (result bool, err error) {
	closing, err := p.findParenClosingPoint()
	if err != nil {
		return
	}

	return p.tokens[closing+1].TokenType == lexer.Arrow ||
		(p.tokens[closing+1].TokenType == lexer.Identifier && p.tokens[closing+2].TokenType == lexer.Arrow), nil
}

func (p *Parser) findParenClosingPoint() (index int, err error) {
	cur := p.current
	for p.tokens[cur].TokenType != lexer.RBrace {
		if p.match(lexer.LBrace) {
			cur, err = p.findParenClosingPoint()
			if err != nil {
				return -1, err
			}
		}
		cur++
		if cur > len(p.tokens) {
			return -1, ParseError{
				token:   p.previous(),
				message: "Unexpected end before closing parenthesis",
			}
		}
	}
	return cur, nil
}
