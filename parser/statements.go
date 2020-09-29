package parser

import "elara/lexer"

type Stmt interface {
	stmtNode()
}

type ExpressionStmt struct {
	Expr Expr
}

type BlockStmt struct {
	Stmts []Stmt
}

type VarDefStmt struct {
	Mutable    bool
	Identifier string
	Type       Type
	Value      Expr
}

type StructDefStmt struct {
	Identifier   string
	StructFields []StructField
}

type IfElseStmt struct {
	Condition  Expr
	MainBranch Stmt
	ElseBranch *Stmt
}

type WhileStmt struct {
	Condition Expr
	Body      Stmt
}

type ExtendStmt struct {
	Condition Expr
	Body      Stmt
}
type TypeStmt struct {
	Identifier string
	Contract   Type
}
type GenerifiedStmt struct {
	Contracts []GenericContract
	Statement Stmt
}

func (ExpressionStmt) stmtNode() {}
func (BlockStmt) stmtNode()      {}
func (VarDefStmt) stmtNode()     {}
func (StructDefStmt) stmtNode()  {}
func (IfElseStmt) stmtNode()     {}
func (WhileStmt) stmtNode()      {}
func (ExtendStmt) stmtNode()     {}
func (GenerifiedStmt) stmtNode() {}
func (t TypeStmt) stmtNode()     {}

func (p *Parser) declaration() (stmt Stmt) {
	if p.check(lexer.Let) {
		return p.varDefStatement()
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

func (p *Parser) varDefStatement() (stmt Stmt) {
	p.consume(lexer.Let, "Expected variable declaration to start with let")
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
	p.cleanNewLines()
	for !p.check(lexer.RBrace) {
		decl := p.declaration()
		result = append(result, decl)
		p.consume(lexer.NEWLINE, "Expected newline after declaration in block")
		p.cleanNewLines()
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
