package parser

import (
	lexer "elara/lexer"
	"fmt"
)

type Scanner = lexer.Scanner
type Token = lexer.Token
type TokenType = lexer.TokenType

type ParseError struct {
	token   Token
	message string
}

func (pe ParseError) Error() string {
	return fmt.Sprintf("Parse Error: %s", pe.message)
}

type Parser struct {
	tokens  []Token
	current int
}

type Type string

func (p *Parser) New(tokens []Token) *Parser {
	return &Parser{
		tokens: tokens,
	}
}

func (p *Parser) Parse() ([]Stmt, []error) {
	p.current = 0
	result := make([]Stmt, 1)
	for !p.isAtEnd() {
		stmt, error := p.declaration()
		if error != nil {
			_ = ParseError{
				token:   p.peek(),
				message: error.Error(),
			}.Error()
		} else {
			result = append(result, stmt)
		}
	}
	return result, nil
}

func (p *Parser) peek() Token {
	return p.tokens[p.current]
}

func (p *Parser) previous() Token {
	return p.tokens[p.current-1]
}

func (p *Parser) isAtEnd() bool {
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

func (p *Parser) match(types ...TokenType) bool {
	for _, t := range types {
		if p.check(t) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) consume(tokenType TokenType, msg string) (token Token, err error) {
	if p.check(tokenType) {
		token = p.advance()
	} else {
		err = ParseError{
			token:   p.peek(),
			message: msg,
		}
	}
	return
}

/*
	Declaration -> let [mut] id = value || Statement
	Statement ->
*/

// ----- Statements -----

func (p *Parser) declaration() (stmt Stmt, err error) {
	if p.match(lexer.Let) {
		mut := p.match(lexer.Mut)

		id, err := p.consume(lexer.Identifier, "Expected identifier for variable declaration")
		if err != nil {
			return
		}

		var typ *Type
		if p.match(lexer.Colon) {
			typStr, err := p.consume(lexer.Identifier, "Expected type after colon in variable declaration")
			if err != nil {
				return
			}
			typI := Type(typStr)
			typ = &typI
		}

		expr, err := p.expression()

		if err != nil {
			return nil, err
		}

		stmt = VarDefStmt{
			Mutable:    mut,
			Identifier: id.Text,
			Type:       typ,
			Value:      expr,
		}
		return
	}

	return p.statement()
}

func (p *Parser) statement() (Stmt, error) {
	switch p.peek().TokenType {
	case lexer.While:
		return p.while()
	case lexer.If:
		return p.ifStmt()
	case lexer.LBrace:
		return p.blockStmt()
	default:
		return p.exprStatement()
	}
}

func (p *Parser) while() (stmt Stmt, err error) {
	_, err = p.consume(lexer.While, "Expected while at beginning of while loop")
	if err != nil {
		return
	}

	expr, err := p.expression()
	if err != nil {
		return
	}

	body, err := p.statement()
	if err != nil {
		return
	}

	stmt = WhileStmt{
		Condition: expr,
		Body:      body,
	}
	return
}

func (p *Parser) ifStmt() (stmt Stmt, err error) {
	_, err = p.consume(lexer.If, "Expected while at beginning of while loop")
	if err != nil {
		return
	}

	condition, err := p.expression()
	if err != nil {
		return
	}
	_, err = p.consume(lexer.Arrow, "Expected arrow after condition for if statement")
	if err != nil {
		return
	}

	mainBranch, err := p.statement()
	if err != nil {
		return
	}
	var elseBranch *Stmt = nil
	if p.match(lexer.Else) {
		if p.check(lexer.If) {
			ebr, err := p.ifStmt()
			if err != nil {
				return
			}
			elseBranch = &ebr
		} else {
			_, err = p.consume(lexer.Arrow, "Expected arrow after condition for else statement")
			if err != nil {
				return
			}
			ebr, err := p.statement()
			if err != nil {
				return
			}
			elseBranch = &ebr
		}
	}
	stmt = IfElseStmt{
		Condition:  condition,
		MainBranch: mainBranch,
		ElseBranch: elseBranch,
	}
	return
}

func (p *Parser) blockStmt() (stmt Stmt, err error) {
	result := make([]Stmt, 1)
	_, err = p.consume(lexer.LBrace, "Expected { at beginning of block")
	if err != nil {
		return
	}
	for !p.check(lexer.RBrace) {
		decl, err := p.declaration()
		if err != nil {
			return
		}
		result = append(result, decl)
	}
	_, err = p.consume(lexer.RBrace, "Expected } at beginning of block")
	if err != nil {
		return
	}
	stmt = BlockStmt{Stmts: result}
	return
}
func (p *Parser) exprStatement() (stmt Stmt, err error) {
	expr, err := p.expression()
	if err != nil {
		return
	}
	stmt = ExpressionStmt{Expr: expr}
	return
}

// ----- Expressions -----

func (p *Parser) expression() (Expr, error) {

}
