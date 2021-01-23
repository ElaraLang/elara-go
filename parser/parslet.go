package parser

import (
	"github.com/ElaraLang/elara/ast"
	"github.com/ElaraLang/elara/lexer"
	"strconv"
)

func (p *Parser) initPrefixParselets() {
	p.prefixParslets = make(map[lexer.TokenType]prefixParslet, 0)
	p.registerPrefix(lexer.Int, p.parseInteger)
	p.registerPrefix(lexer.Float, p.parseFloat)
	p.registerPrefix(lexer.Char, p.parseChar)
	p.registerPrefix(lexer.String, p.parseString)
	p.registerPrefix(lexer.LParen, p.resolvingParslet(p.functionGroupResolver()))
	p.registerPrefix(lexer.BooleanTrue, p.parseBoolean)
	p.registerPrefix(lexer.BooleanFalse, p.parseBoolean)
	p.registerPrefix(lexer.Subtract, p.parseUnaryExpression)
	p.registerPrefix(lexer.Not, p.parseUnaryExpression)
}

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

func (p *Parser) initStatementParselets() {
	p.statementParslets = make(map[lexer.TokenType]statementParslet, 0)
}

func (p *Parser) parseExpressionStatement() ast.Statement {
	return &ast.ExpressionStatement{
		Token:      p.Tape.Current(),
		Expression: p.parseExpression(LOWEST),
	}
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

func (p *Parser) parseUnaryExpression() ast.Expression {
	operator := p.Tape.Consume(lexer.Dot)
	expr := p.parseExpression(PREFIX)
	return &ast.UnaryExpression{
		Token:    operator,
		Operator: operator,
		Right:    expr,
	}
}

func (p *Parser) parseIdentifier() ast.Identifier {
	token := p.Tape.Consume(lexer.Int)
	return ast.Identifier{Token: token, Name: string(token.Text)}
}

func (p *Parser) parseInteger() ast.Expression {
	token := p.Tape.Consume(lexer.Int)
	value, err := strconv.ParseInt(string(token.Text), 10, 64)
	if err != nil {
		// panic
	}
	return &ast.IntegerLiteral{Token: token, Value: value}
}

func (p *Parser) parseFloat() ast.Expression {
	token := p.Tape.Consume(lexer.Float)
	value, err := strconv.ParseFloat(string(token.Text), 10)
	if err != nil {
		// panic
	}
	return &ast.FloatLiteral{Token: token, Value: value}
}

func (p *Parser) parseBoolean() ast.Expression {
	token := p.Tape.Consume(lexer.BooleanTrue, lexer.BooleanFalse)
	value := token.TokenType == lexer.BooleanTrue
	return &ast.BooleanLiteral{Token: token, Value: value}
}

func (p *Parser) parseChar() ast.Expression {
	token := p.Tape.Consume(lexer.Char)
	value := token.Text[0]
	return &ast.CharLiteral{Token: token, Value: value}
}

func (p *Parser) parseString() ast.Expression {
	token := p.Tape.Consume(lexer.String)
	value := string(token.Text)
	return &ast.StringLiteral{Token: token, Value: value}
}

func (p *Parser) parseFunction() ast.Expression {
	token := p.Tape.Consume(lexer.LParen)
	params := p.parseFunctionParameters()
	p.Tape.Expect(lexer.RParen)
	p.Tape.Expect(lexer.Arrow)

	var typ ast.Type

	if !p.Tape.ValidationPeek(0, lexer.LBrace) {
		typ = p.parseType()
	}

	body := p.parseStatement()
	return &ast.FunctionLiteral{
		Token:      token,
		ReturnType: typ,
		Parameters: params,
		Body:       body,
	}
}

func (p *Parser) parseGroupExpression() ast.Expression {
	p.Tape.Expect(lexer.LParen)
	expr := p.parseExpression(LOWEST)
	p.Tape.Expect(lexer.RParen)
	return expr
}

func (p *Parser) parseType() ast.Type {
	return nil // TODO
}

func (p *Parser) parseCollection() ast.Expression {
	tok := p.Tape.Consume(lexer.LSquare)
	elements := p.parseCollectionElements()
	p.Tape.Expect(lexer.RSquare)
	return &ast.CollectionLiteral{Token: tok, Elements: elements}
}

func (p *Parser) parseMap() ast.Expression {
	tok := p.Tape.Consume(lexer.LBrace)
	elements := p.parseMapEntries()
	p.Tape.Consume(lexer.RBrace)
	return &ast.MapLiteral{Token: tok, Entries: elements}
}
