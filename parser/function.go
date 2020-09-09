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
		param, error := p.expression()
		if error != nil {
			err = error
			return
		}
		params = append(params, param)
		if separator != nil {
			_, err = p.consume(*separator, "Expected separator in function parameters")
			if error != nil {
				err = error
				return
			}
		}
	}
	expr = params
	return
}

func (p *Parser) functionArguments() (args []FunctionArgument, err error) {
	args = make([]FunctionArgument, 0)
	_, error := p.consume(lexer.LParen, "Expected left paren before starting function definition")
	if error != nil {
		err = error
		return
	}
	for !p.match(lexer.RParen) {
		arg, error := p.functionArgument()
		if error != nil {
			err = error
			return
		}
		args = append(args, arg)
	}
	return
}

func (p *Parser) functionArgument() (arg FunctionArgument, err error) {
	i1, error := p.consume(lexer.Identifier, "Invalid argument in function def")
	if error != nil {
		err = error
		return
	}
	if p.match(lexer.Equal) {
		expr, error := p.expression()
		if error != nil {
			err = error
			return
		}
		arg = FunctionArgument{
			Type:     nil,
			Variable: VariableExpr{Identifier: i1.Text},
			Default:  &expr,
		}
		return
	}
	id, error := p.consume(lexer.Identifier, "Invalid argument in function def")
	if error != nil {
		err = error
		return
	}

	var def *Expr

	if p.match(lexer.Equal) {
		expr, error := p.expression()
		if error != nil {
			err = error
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
	closing, error := p.findParenClosingPoint()
	if error != nil {
		err = error
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
				return
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
