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
	p.registerPrefix(lexer.If, p.parseUnaryExpression)
}

func (p *Parser) parseIfExpression() ast.Expression {
	operator := p.Tape.Consume(lexer.If)
	condition := p.parseExpression(Lowest)
	var mainBranch ast.Statement
	var elseBranch ast.Statement
	if p.Tape.Match(lexer.Arrow) {
		mainBranch = p.parseExpressionStatement()
	} else {
		mainBranch = p.parseBlockStatement()
	}
	if p.Tape.Match(lexer.Else) {
		switch p.Tape.Current().TokenType {
		case lexer.Arrow:
			p.Tape.Match(lexer.Arrow)
			fallthrough
		case lexer.If:
			elseBranch = p.parseExpressionStatement()
		default:
			elseBranch = p.parseBlockStatement()

		}
	}

	return &ast.IfExpression{
		Token:      operator,
		Condition:  condition,
		MainBranch: mainBranch,
		ElseBranch: elseBranch,
	}
}

func (p *Parser) parseUnaryExpression() ast.Expression {
	operator := p.Tape.Consume(lexer.Dot)
	expr := p.parseExpression(Prefix)
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
		typ = p.parseType(TypeLowest)
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
	expr := p.parseExpression(Lowest)
	p.Tape.Expect(lexer.RParen)
	return expr
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
