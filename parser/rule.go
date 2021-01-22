package parser

import (
	"github.com/ElaraLang/elara/ast"
	"github.com/ElaraLang/elara/lexer"
	"strconv"
)

func (p *Parser) initializePrefixParselets() {
	p.prefixParseFunctions = make(map[lexer.TokenType]parsePrefix)
	p.registerPrefix(lexer.Int, p.parseInteger)
	p.registerPrefix(lexer.Float, p.parseFloat)
	p.registerPrefix(lexer.Char, p.parseFloat)
	p.registerPrefix(lexer.String, p.parseFloat)
	// p.registerPrefix(lexer.LParen, p.parseFunction) // TODO Ambiguity resolver required
	p.registerPrefix(lexer.BooleanTrue, p.parseBoolean)
	p.registerPrefix(lexer.BooleanFalse, p.parseBoolean)
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
