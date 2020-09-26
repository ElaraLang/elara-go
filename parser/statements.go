package parser

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
	Type       *Type
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
