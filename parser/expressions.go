package parser

import (
	"elara/lexer"
	"strconv"
)

type Expr interface{ exprNode() }

type BinaryExpr struct {
	Lhs Expr
	Op  TokenType
	Rhs Expr
}

type UnaryExpr struct {
	Op  TokenType
	Rhs Expr
}

type GroupExpr struct {
	Group Expr
}

type VariableExpr struct {
	Identifier string
}

type AssignmentExpr struct {
	Context    *Expr
	Identifier string
	Value      Expr
}

type InvocationExpr struct {
	Invoker Expr
	Args    []Expr
}

type ContextExpr struct {
	Context  Expr
	Variable VariableExpr
}

type IfElseExpr struct {
	Condition  Expr
	MainBranch []Stmt
	MainResult Expr
	ElseBranch []Stmt
	ElseResult Expr
}

type FuncDefExpr struct {
	Arguments  []FunctionArgument
	ReturnType *Type
	Statement  Stmt
}

type StringLiteralExpr struct {
	Value string
}

type IntegerLiteralExpr struct {
	Value int64
}

type FloatLiteralExpr struct {
	Value float64
}

type BooleanLiteralExpr struct {
	Value bool
}

func (FuncDefExpr) exprNode()        {}
func (StringLiteralExpr) exprNode()  {}
func (IntegerLiteralExpr) exprNode() {}
func (FloatLiteralExpr) exprNode()   {}
func (BooleanLiteralExpr) exprNode() {}
func (UnaryExpr) exprNode()          {}
func (BinaryExpr) exprNode()         {}
func (GroupExpr) exprNode()          {}
func (ContextExpr) exprNode()        {}
func (IfElseExpr) exprNode()         {}
func (InvocationExpr) exprNode()     {}
func (AssignmentExpr) exprNode()     {}
func (VariableExpr) exprNode()       {}

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
	if p.check(lexer.LParen) {
		isFunc := p.isFuncDef()
		if isFunc {
			args := p.functionArguments()

			var typ *Type

			if p.check(lexer.Identifier) {
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
