package parser

import "C"
import (
	lexer "elara/lexer"
	"fmt"
	"strconv"
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

func (p *Parser) Parse() (result []Stmt, error []ParseError) {
	p.current = 0
	result = make([]Stmt, 0)
	error = make([]ParseError, 0)
	for !p.isAtEnd() {
		stmt, err := p.declaration()
		if err != nil {
			switch pErr := err.(type) {
			case ParseError:
				error = append(error, pErr)
				break
			default:
				break
			}
			p.syncError()
		} else {
			result = append(result, stmt)
			if !(p.match(lexer.NEWLINE) || p.match(lexer.EOF)) {
				pErr := ParseError{
					token:   p.peek(),
					message: "Expected new line",
				}
				error = append(error, pErr)
			}
		}
	}
	return
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
		return p.whileStatement()
	case lexer.If:
		return p.ifStatement()
	case lexer.LBrace:
		return p.blockStatement()
	default:
		return p.exprStatement()
	}
}

func (p *Parser) whileStatement() (stmt Stmt, err error) {
	_, err = p.consume(lexer.While, "Expected whileStatement at beginning of whileStatement loop")
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

func (p *Parser) ifStatement() (stmt Stmt, err error) {
	_, err = p.consume(lexer.If, "Expected whileStatement at beginning of whileStatement loop")
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
			ebr, err := p.ifStatement()
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

func (p *Parser) blockStatement() (stmt Stmt, err error) {
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
	return p.assignment()
}

func (p *Parser) assignment() (expr Expr, err error) {
	expr, err = p.logicalOr()
	if err != nil {
		return
	}

	if p.check(lexer.Equal) {
		eqlTok := p.advance()
		rhs, err := p.logicalOr()
		if err != nil {
			return
		}
		switch v := expr.(type) {
		default:
			err = ParseError{
				token:   eqlTok,
				message: "Invalid type found behind assignment",
			}
		case VariableExpr:
			expr = AssignmentExpr{
				Identifier: v.Identifier,
				Value:      rhs,
			}
		case ContextExpr:
			expr = AssignmentExpr{
				Context:    &v.Context,
				Identifier: v.Variable.Identifier,
				Value:      rhs,
			}
		}
	}
	return
}

func (p *Parser) logicalOr() (expr Expr, err error) {
	expr, err = p.logicalAnd()
	if err != nil {
		return
	}

	for p.match(lexer.Or) {
		op := p.previous()
		rhs, err := p.logicalAnd()
		if err != nil {
			return
		}
		expr = BinaryExpr{
			Lhs: expr,
			Op:  op.TokenType,
			Rhs: rhs,
		}
	}
	return
}

func (p *Parser) logicalAnd() (expr Expr, err error) {
	expr, err = p.referenceEquality()
	if err != nil {
		return
	}

	for p.match(lexer.And) {
		op := p.previous()
		rhs, err := p.referenceEquality()
		if err != nil {
			return
		}
		expr = BinaryExpr{
			Lhs: expr,
			Op:  op.TokenType,
			Rhs: rhs,
		}
	}
	return
}

func (p *Parser) referenceEquality() (expr Expr, err error) {
	expr, err = p.comparison()
	if err != nil {
		return
	}

	for p.match(lexer.Equals, lexer.NotEquals) {
		op := p.previous()
		rhs, err := p.comparison()
		if err != nil {
			return
		}
		expr = BinaryExpr{
			Lhs: expr,
			Op:  op.TokenType,
			Rhs: rhs,
		}
	}
	return
}

func (p *Parser) comparison() (expr Expr, err error) {
	expr, err = p.addition()
	if err != nil {
		return
	}

	for p.match(lexer.GreaterEqual, lexer.RAngle, lexer.LesserEqual, lexer.LAngle) {
		op := p.previous()
		rhs, err := p.addition()
		if err != nil {
			return
		}
		expr = BinaryExpr{
			Lhs: expr,
			Op:  op.TokenType,
			Rhs: rhs,
		}
	}
	return
}

func (p *Parser) addition() (expr Expr, err error) {
	expr, err = p.multiplication()
	if err != nil {
		return
	}

	for p.match(lexer.Add, lexer.Subtract) {
		op := p.previous()
		rhs, err := p.multiplication()
		if err != nil {
			return
		}
		expr = BinaryExpr{
			Lhs: expr,
			Op:  op.TokenType,
			Rhs: rhs,
		}
	}
	return
}

func (p *Parser) multiplication() (expr Expr, err error) {
	expr, err = p.unary()
	if err != nil {
		return
	}

	for p.match(lexer.Multiply, lexer.Slash, lexer.Mod) {
		op := p.previous()
		rhs, err := p.unary()
		if err != nil {
			return
		}
		expr = BinaryExpr{
			Lhs: expr,
			Op:  op.TokenType,
			Rhs: rhs,
		}
	}
	return
}

func (p *Parser) unary() (expr Expr, err error) {
	if p.match(lexer.Subtract, lexer.Not, lexer.Add) {
		op := p.previous()
		rhs, err := p.unary()
		if err != nil {
			return
		}
		expr = UnaryExpr{
			Op:  op.TokenType,
			Rhs: rhs,
		}
		return
	}
	expr, err = p.invoke()
	return
}

func (p *Parser) invoke() (expr Expr, err error) {
	expr, err = p.funDef()
	if err != nil {
		return
	}
	for p.match(lexer.LParen, lexer.Dot) {
		switch p.previous().TokenType {
		case lexer.LParen:
			separator := lexer.Comma
			args, err := p.invocationParameters(&separator)
			if err != nil {
				return
			}
			expr = InvocationExpr{
				Invoker: expr,
				Args:    args,
			}
			break
		case lexer.Dot:
			id, err := p.consume(lexer.Identifier, "Expected identifier inside context getter/setter")
			if err != nil {
				return
			}
			expr = ContextExpr{
				Context:  expr,
				Variable: VariableExpr{Identifier: id.Text},
			}
			break
		}
	}
	return
}

func (p *Parser) funDef() (expr Expr, err error) {
	if p.match(lexer.LParen) {
		isFunc, err := p.isFuncDef()
		if err != nil {
			return
		}
		if isFunc {
			args, err := p.functionArguments()
			if err != nil {
				return
			}
			var typ *Type

			if p.match(lexer.Identifier) {
				typ1 := Type(p.previous().Text)
				typ = &typ1
			}

			_, err = p.consume(lexer.Arrow, "Expected arrow at function definition")
			if err != nil {
				return
			}

			stmt, err := p.statement()
			if err != nil {
				return
			}
			expr = FuncDefExpr{
				Arguments:  args,
				ReturnType: typ,
				Statement:  stmt,
			}
			return
		} else {
			expr, err = p.expression()
			if err != nil {
				return
			}
			expr = GroupExpr{Group: expr}
			return
		}
	}
	return p.primary()
}

func (p *Parser) primary() (expr Expr, err error) {
	switch p.peek().TokenType {
	case lexer.Boolean:
		truth, err := p.consume(lexer.Boolean, "Expected boolean")
		if err != nil {
			return
		}
		boolVal, err := strconv.ParseBool(truth.Text)
		if err != nil {
			return
		}
		expr = BooleanLiteralExpr{Value: boolVal}
		break
	case lexer.String:
		str, err := p.consume(lexer.String, "Expected string")
		if err != nil {
			return
		}
		expr = StringLiteralExpr{Value: str.Text}
		break
	case lexer.Int:
		integr, err := p.consume(lexer.Int, "Expected integer")
		if err != nil {
			return
		}
		intVal, err := strconv.ParseInt(integr.Text, 10, 64)
		if err != nil {
			return
		}
		expr = IntegerLiteralExpr{Value: intVal}
		break
	case lexer.Float:
		flt, err := p.consume(lexer.Float, "Expected float")
		if err != nil {
			return
		}
		fltVal, err := strconv.ParseFloat(flt.Text, 64)
		if err != nil {
			return
		}
		expr = FloatLiteralExpr{Value: fltVal}
		break
	case lexer.Identifier:
		str, err := p.consume(lexer.Identifier, "Expected identifier")
		if err != nil {
			return
		}
		expr = StringLiteralExpr{Value: str.Text}
		break
	}

	if expr != nil {
		err = ParseError{
			token:   p.peek(),
			message: "Invalid expression",
		}
		return
	}
	return
}

func (p *Parser) syncError() {
	for !p.isAtEnd() || p.peek().TokenType != lexer.NEWLINE || p.peek().TokenType != lexer.EOF {
		p.advance()
	}
	for p.match(lexer.NEWLINE) {
	}
}
