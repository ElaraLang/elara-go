package parser

import "elara/lexer"

type FunctionArgument struct {
	Type     *Type
	Variable VariableExpr
	Default  Expr
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

func (p *Parser) functionArgument() (args []FunctionArgument, err error) {
	args := make([]FunctionArgument, 0)
	for !p.match(lexer.RParen) {

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
