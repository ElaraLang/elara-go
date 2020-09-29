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
	return fmt.Sprintf("Parse Error: %s at %s", pe.message, pe.token.Text)
}

type Parser struct {
	tokens  []Token
	current int
}

func NewParser(tokens *[]Token) *Parser {
	return &Parser{
		tokens: *tokens,
	}
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
	stmt := p.declaration()
	*result = append(*result, stmt)
	if !(p.match(lexer.NEWLINE) || p.isAtEnd()) {
		panic(ParseError{
			token:   p.peek(),
			message: "Expected new line",
		})
	}
}

func (p *Parser) handleError(error *[]ParseError) {
	if r := recover(); r != nil {
		switch err := r.(type) {
		case ParseError:
			*error = append(*error, err)
			break
		case []ParseError:
			*error = append(*error, err...)
		default:
			*error = append(*error, ParseError{
				token:   p.previous(),
				message: "Invalid error thrown by Parser",
			})
			break
		}
		p.syncError()
	}
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

func (p *Parser) consume(tokenType TokenType, msg string) (token Token) {
	if p.check(tokenType) {
		token = p.advance()
		return
	}
	panic(ParseError{token: p.peek(), message: msg})
}

func (p *Parser) cleanNewLines() {
	for p.check(lexer.NEWLINE) {
		p.advance()
	}
}

// ----- Statements -----

func (p *Parser) declaration() (stmt Stmt) {
	if p.match(lexer.Let) {
		mut := p.match(lexer.Mut)

		id := p.consume(lexer.Identifier, "Expected identifier for variable declaration")
		if p.match(lexer.Arrow) {
			execStmt := p.statement()

			expr := FuncDefExpr{
				Arguments:  make([]FunctionArgument, 0),
				ReturnType: nil,
				Statement:  execStmt,
			}
			return VarDefStmt{
				Mutable:    mut,
				Identifier: id.Text,
				Type:       nil,
				Value:      expr,
			}
		}

		var typ Type
		if p.match(lexer.Colon) {
			typI := p.typeContract()
			typ = typI
		}

		p.consume(lexer.Equal, "Expected Equal on variable declaration")
		expr := p.expression()

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

func (p *Parser) statement() Stmt {
	switch p.peek().TokenType {
	case lexer.While:
		return p.whileStatement()
	case lexer.If:
		return p.ifStatement()
	case lexer.LBrace:
		return p.blockStatement()
	case lexer.Struct:
		return p.structStatement()
	case lexer.Type:
		return p.typeStatement()
	case lexer.LAngle:
		return p.genericStatement()
	default:
		return p.exprStatement()
	}
}

func (p *Parser) whileStatement() (stmt Stmt) {
	p.consume(lexer.While, "Expected whileStatement at beginning of whileStatement loop")
	expr := p.expression()

	p.consume(lexer.Arrow, "Expected arrow after condition for if statement")

	body := p.statement()

	stmt = WhileStmt{
		Condition: expr,
		Body:      body,
	}
	return
}

func (p *Parser) ifStatement() (stmt Stmt) {
	p.consume(lexer.If, "Expected whileStatement at beginning of whileStatement loop")

	condition := p.logicalOr()

	p.consume(lexer.Arrow, "Expected arrow after condition for if statement")

	mainBranch := p.statement()

	var elseBranch *Stmt = nil
	if p.match(lexer.Else) {
		if p.check(lexer.If) {
			ebr := p.ifStatement()
			elseBranch = &ebr
		} else {
			p.consume(lexer.Arrow, "Expected arrow after condition for else statement")

			ebr := p.statement()
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

func (p *Parser) blockStatement() (stmt Stmt) {
	result := make([]Stmt, 0)
	p.consume(lexer.LBrace, "Expected { at beginning of block")

	for !p.check(lexer.RBrace) {
		decl := p.declaration()
		result = append(result, decl)
	}
	p.consume(lexer.RBrace, "Expected } at beginning of block")

	stmt = BlockStmt{Stmts: result}
	return
}

func (p *Parser) structStatement() (stmt Stmt) {
	p.consume(lexer.Struct, "Expected struct start to begin with `struct` keyword")

	identifier := p.consume(lexer.Identifier, "Expected identifier after `struct` keyword")

	fields := p.structFields()
	return StructDefStmt{
		Identifier:   identifier.Text,
		StructFields: fields,
	}
}

func (p *Parser) exprStatement() (stmt Stmt) {
	expr := p.expression()
	stmt = ExpressionStmt{Expr: expr}
	return
}

// ----- Expressions -----

func (p *Parser) expression() Expr {
	return p.assignment()
}

func (p *Parser) assignment() (expr Expr) {
	expr = p.logicalOr()

	if p.check(lexer.Equal) {
		eqlTok := p.advance()
		rhs := p.logicalOr()

		switch v := expr.(type) {
		case VariableExpr:
			expr = AssignmentExpr{
				Identifier: v.Identifier,
				Value:      rhs,
			}
			break
		case ContextExpr:
			expr = AssignmentExpr{
				Context:    &v.Context,
				Identifier: v.Variable.Identifier,
				Value:      rhs,
			}
			break
		default:
			panic(ParseError{
				token:   eqlTok,
				message: "Invalid type found behind assignment",
			})
		}
	}
	return
}

func (p *Parser) logicalOr() (expr Expr) {
	expr = p.logicalAnd()

	for p.match(lexer.Or) {
		op := p.previous()
		rhs := p.logicalAnd()
		expr = BinaryExpr{
			Lhs: expr,
			Op:  op.TokenType,
			Rhs: rhs,
		}
	}
	return
}

func (p *Parser) logicalAnd() (expr Expr) {
	expr = p.referenceEquality()

	for p.match(lexer.And) {
		op := p.previous()
		rhs := p.referenceEquality()

		expr = BinaryExpr{
			Lhs: expr,
			Op:  op.TokenType,
			Rhs: rhs,
		}
	}
	return
}

func (p *Parser) referenceEquality() (expr Expr) {
	expr = p.comparison()

	for p.match(lexer.Equals, lexer.NotEquals) {
		op := p.previous()
		rhs := p.comparison()

		expr = BinaryExpr{
			Lhs: expr,
			Op:  op.TokenType,
			Rhs: rhs,
		}
	}
	return
}

func (p *Parser) comparison() (expr Expr) {
	expr = p.addition()

	for p.match(lexer.GreaterEqual, lexer.RAngle, lexer.LesserEqual, lexer.LAngle) {
		op := p.previous()
		rhs := p.addition()

		expr = BinaryExpr{
			Lhs: expr,
			Op:  op.TokenType,
			Rhs: rhs,
		}
	}
	return
}

func (p *Parser) addition() (expr Expr) {
	expr = p.multiplication()

	for p.match(lexer.Add, lexer.Subtract) {
		op := p.previous()
		rhs := p.multiplication()
		expr = BinaryExpr{
			Lhs: expr,
			Op:  op.TokenType,
			Rhs: rhs,
		}
	}
	return
}

func (p *Parser) multiplication() (expr Expr) {
	expr = p.unary()

	for p.match(lexer.Multiply, lexer.Slash, lexer.Mod) {
		op := p.previous()
		rhs := p.unary()
		expr = BinaryExpr{
			Lhs: expr,
			Op:  op.TokenType,
			Rhs: rhs,
		}
	}
	return
}

func (p *Parser) unary() (expr Expr) {
	if p.match(lexer.Subtract, lexer.Not, lexer.Add) {
		op := p.previous()
		rhs := p.unary()
		expr = UnaryExpr{
			Op:  op.TokenType,
			Rhs: rhs,
		}
		return
	}
	expr = p.invoke()
	return
}

func (p *Parser) invoke() (expr Expr) {
	expr = p.funDef()

	for p.match(lexer.LParen, lexer.Dot) {
		switch p.previous().TokenType {
		case lexer.LParen:
			separator := lexer.Comma
			args := p.invocationParameters(&separator)

			expr = InvocationExpr{
				Invoker: expr,
				Args:    args,
			}
			break
		case lexer.Dot:
			id := p.consume(lexer.Identifier, "Expected identifier inside context getter/setter")

			expr = ContextExpr{
				Context:  expr,
				Variable: VariableExpr{Identifier: id.Text},
			}
			break
		}
	}
	return
}

func (p *Parser) funDef() (expr Expr) {
	if p.match(lexer.LParen) {
		isFunc := p.isFuncDef()
		if isFunc {
			args := p.functionArguments()

			var typ *Type

			if p.match(lexer.Identifier) {
				typ1 := p.typeContract()

				typ = &typ1
			}

			p.consume(lexer.Arrow, "Expected arrow at function definition")

			stmt := p.statement()

			expr = FuncDefExpr{
				Arguments:  args,
				ReturnType: typ,
				Statement:  stmt,
			}
			return
		} else {
			expr = p.expression()
			expr = GroupExpr{Group: expr}
			return
		}
	}
	return p.primary()
}

func (p *Parser) primary() (expr Expr) {
	var error error
	switch p.peek().TokenType {
	case lexer.String:
		str := p.consume(lexer.String, "Expected string")

		expr = StringLiteralExpr{Value: str.Text}
		break
	case lexer.Boolean:
		truth := p.consume(lexer.Boolean, "Expected boolean")

		boolVal, err := strconv.ParseBool(truth.Text)
		error = err
		expr = BooleanLiteralExpr{Value: boolVal}
		break
	case lexer.Int:
		integr := p.consume(lexer.Int, "Expected integer")
		intVal, err := strconv.ParseInt(integr.Text, 10, 64)
		error = err
		expr = IntegerLiteralExpr{Value: intVal}
		break
	case lexer.Float:
		flt := p.consume(lexer.Float, "Expected float")

		fltVal, err := strconv.ParseFloat(flt.Text, 64)
		error = err
		expr = FloatLiteralExpr{Value: fltVal}
		break
	case lexer.Identifier:
		str := p.consume(lexer.Identifier, "Expected identifier")
		expr = VariableExpr{Identifier: str.Text}
		break
	}
	if error != nil {
		panic(ParseError{
			token:   p.previous(),
			message: "Expected literal",
		})
	}

	if expr == nil {
		panic(ParseError{
			token:   p.peek(),
			message: "Invalid expression",
		})
	}
	return
}

func (p *Parser) syncError() {
	for !p.isAtEnd() && !p.check(lexer.NEWLINE) && !p.check(lexer.EOF) {
		p.advance()
	}
	for p.match(lexer.NEWLINE) {
	}
}
