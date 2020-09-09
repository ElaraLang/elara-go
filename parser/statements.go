package parser

type Stmt interface{ stmtNode() }

type ExpressionStmt struct {
	Expr Expr
}

type BlockStmt struct {
	Stmts []Stmt
}

type VarDefStmt struct {
	Mutable    bool
	Identifier string
	Type       *Type
	Value      Expr
}

type StructDefStmt struct {
	Identifier string
	Fields     []string
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

func (ExpressionStmt) stmtNode() {}
func (BlockStmt) stmtNode()      {}
func (VarDefStmt) stmtNode()     {}
func (StructDefStmt) stmtNode()  {}
func (IfElseStmt) stmtNode()     {}
func (WhileStmt) stmtNode()      {}
func (ExtendStmt) stmtNode()     {}
