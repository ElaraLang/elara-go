package parser

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
	Group Expr
}

type AssignmentExpr struct {
	Identifier string
	Value      Expr
}

type InvocationExpr struct {
	Invoker Expr
	Args    []Expr
}

type ContextExpr struct {
	Context Expr
	Expr    Expr
}

type IfElseExpr struct {
	Condition  Expr
	MainBranch []Stmt
	MainResult Expr
	ElseBranch []Stmt
	ElseResult Expr
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

func (StringLiteralExpr) exprNode()  {}
func (IntegerLiteralExpr) exprNode() {}
func (FloatLiteralExpr) exprNode()   {}
func (BooleanLiteralExpr) exprNode() {}
func (UnaryExpr) exprNode()          {}
func (BinaryExpr) exprNode()         {}
func (GroupExpr) exprNode()          {}
func (InvocationExpr) exprNode()     {}
func (AssignmentExpr) exprNode()     {}
func (VariableExpr) exprNode()       {}
