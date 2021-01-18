package parser

type Program struct {
	Statements []Statement
}

type Statement interface {
	toString()
	statementNode()
}

type Expression interface {
	toString()
	expressionNode()
}
