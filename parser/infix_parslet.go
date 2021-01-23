package parser

import (
	"github.com/ElaraLang/elara/ast"
	"github.com/ElaraLang/elara/lexer"
)

func (p *Parser) initInfixParselets() {
	p.infixParslets = make(map[lexer.TokenType]infixParselet, 0)
	p.registerInfix(lexer.Add, p.parseBinaryExpression)
	p.registerInfix(lexer.Subtract, p.parseBinaryExpression)
	p.registerInfix(lexer.Multiply, p.parseBinaryExpression)
	p.registerInfix(lexer.Slash, p.parseBinaryExpression)
	p.registerInfix(lexer.Equals, p.parseBinaryExpression)
	p.registerInfix(lexer.NotEquals, p.parseBinaryExpression)
	p.registerInfix(lexer.GreaterEqual, p.parseBinaryExpression)
	p.registerInfix(lexer.LesserEqual, p.parseBinaryExpression)
	p.registerInfix(lexer.LAngle, p.parseBinaryExpression)
	p.registerInfix(lexer.RAngle, p.parseBinaryExpression)
	p.registerInfix(lexer.Dot, p.parsePropertyExpression)
	p.registerInfix(lexer.LParen, p.parseFunctionCall)
	p.registerInfix(lexer.LSquare, p.parseAccessOperator)
}

func (p *Parser) parseFunctionCall(left ast.Expression) ast.Expression {
	opening := p.Tape.Consume(lexer.LParen)
	args := p.parseFunctionCallArguments()
	p.Tape.Expect(lexer.RParen)
	return &ast.CallExpression{
		Token:      opening,
		Expression: left,
		Arguments:  args,
	}
}

func (p *Parser) parseAccessOperator(left ast.Expression) ast.Expression {
	opening := p.Tape.Consume(lexer.LSquare)
	index := p.parseExpression(LOWEST)
	p.Tape.Expect(lexer.RSquare)
	return &ast.AccessExpression{
		Token:      opening,
		Expression: left,
		Index:      index,
	}
}

func (p *Parser) parseBinaryExpression(left ast.Expression) ast.Expression {
	operator := p.Tape.ConsumeAny()
	precedence := precedenceOf(operator.TokenType)
	right := p.parseExpression(precedence)
	return &ast.BinaryExpression{
		Token:    operator,
		Left:     left,
		Operator: operator,
		Right:    right,
	}
}

func (p *Parser) parsePropertyExpression(left ast.Expression) ast.Expression {
	token := p.Tape.Consume(lexer.Dot)
	right := p.parseIdentifier()
	return &ast.PropertyExpression{
		Token:    token,
		Context:  left,
		Variable: right,
	}
}
