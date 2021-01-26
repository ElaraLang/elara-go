package parserlegacy

import (
	"fmt"
	"github.com/ElaraLang/elara/lexer"
)

type Scanner = lexer.TokenReader
type Token = lexer.Token
type TokenType = lexer.TokenType

type ParseError struct {
	token   Token
	message string
}

func (pe ParseError) Error() string {
	return fmt.Sprintf("Parse ErrorChannel: %s at %s", pe.message, pe.token.String())
}

type Parser struct {
	tokens  []Token
	current int
}

func NewEmptyParser() *Parser {
	return &Parser{}
}

func NewParser(tokens []Token) *Parser {
	return &Parser{
		tokens: tokens,
	}
}

func (p *Parser) Reset(tokens []Token) {
	p.tokens = tokens
	p.current = 0
}

func (p *Parser) Parse() (result []Stmt, error []ParseError) {
	p.current = 0
	result = make([]Stmt, 0)
	error = make([]ParseError, 0)

	for !p.isAtEnd() {
		p.parseLine(&result, &error)
	}
	return
}

func (p *Parser) parseLine(result *[]Stmt, error *[]ParseError) {
	defer p.handleError(error)
	if p.peek().TokenType == lexer.NEWLINE {
		p.advance()
		return
	}

	if p.current == 0 && p.check(lexer.Namespace) {
		ns, importStmt := p.parseFileMeta()
		*result = append(*result, ns, importStmt)
		return
	}

	stmt := p.declaration()
	*result = append(*result, stmt)
	if !(p.match(lexer.NEWLINE) || p.isAtEnd()) {
		panic(ParseError{
			token:   p.peek(),
			message: "Expected new line",
		})
	}
}

func (p *Parser) handleError(errors *[]ParseError) {
	if r := recover(); r != nil {
		switch err := r.(type) {
		case ParseError:
			*errors = append(*errors, err)
			break
		case []ParseError:
			*errors = append(*errors, err...)
		case error:
			*errors = append(*errors, ParseError{
				token:   p.previous(),
				message: err.Error(),
			})
		default:
			*errors = append(*errors, ParseError{
				token:   p.previous(),
				message: "Invalid errors thrown by Parser: ",
			})
			break
		}
		p.syncError()
	}
}

func (p *Parser) peek() Token {
	if p.current >= len(p.tokens) {
		return Token{
			TokenType: lexer.EOF,
		}
	}
	return p.tokens[p.current]
}

func (p *Parser) previous() Token {
	return p.tokens[p.current-1]
}

func (p *Parser) isAtEnd() bool {
	if p.current == len(p.tokens) {
		return true
	}
	return p.peek().TokenType == lexer.EOF
}

func (p *Parser) check(tokenType TokenType) bool {
	return !p.isAtEnd() && p.peek().TokenType == tokenType
}

func (p *Parser) advance() Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}
func (p *Parser) reverse() {
	if p.current == 0 {
		return
	}
	p.current--
}

func (p *Parser) match(types ...TokenType) bool {
	for _, t := range types {
		if p.check(t) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) consume(tokenType TokenType, msg string) (token Token) {
	if p.check(tokenType) {
		token = p.advance()
		return
	}
	panic(ParseError{token: p.peek(), message: msg})
}

func (p *Parser) consumeValidIdentifier(msg string) (token Token) {
	next := p.peek()
	nextType := next.TokenType
	if nextType != lexer.Identifier && nextType != lexer.Add && nextType != lexer.Subtract && nextType != lexer.Slash && nextType != lexer.Multiply {
		panic(ParseError{token: p.peek(), message: msg})
	}
	p.advance()
	return next
}

func (p *Parser) cleanNewLines() {
	for p.match(lexer.NEWLINE) {
	}
}
func (p *Parser) insert(index int, value ...Token) {
	if len(p.tokens) == index {
		p.tokens = append(p.tokens, value...)
	}
	p.tokens = append(p.tokens[:index+len(value)], p.tokens[index:]...)
	for i := 0; i < len(value); i++ {
		p.tokens[index+i] = value[i]
	}
}

func (p *Parser) insertBlankType(index int, value ...TokenType) {
	blankTokens := make([]Token, len(value))
	for i := range value {
		blankTokens[i] = Token{
			TokenType: value[i],
			Text:      nil,
			Position:  lexer.CreatePosition(-1, 1),
		}
	}
	p.insert(index, blankTokens...)
}

func (p *Parser) syncError() {
	for !p.isAtEnd() && !p.check(lexer.NEWLINE) && !p.check(lexer.EOF) {
		p.advance()
	}
	p.cleanNewLines()
}

func (p *Parser) parseProperties(propTypes ...lexer.TokenType) []bool {
	result := make([]bool, len(propTypes))
	for contains(propTypes, p.peek().TokenType) {
		tokTyp := p.advance().TokenType
		for i := 0; i < len(propTypes); i++ {
			if propTypes[i] == tokTyp {
				if result[i] {
					panic(ParseError{
						token:   p.previous(),
						message: "Multiple variable properties of same type defined",
					})
				}
				result[i] = true
				break
			}
		}
	}
	return result
}

func contains(s []lexer.TokenType, e lexer.TokenType) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
