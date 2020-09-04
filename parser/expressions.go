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
