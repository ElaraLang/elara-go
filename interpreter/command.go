package interpreter

import (
	"elara/parser"
	"reflect"
)

type Command interface {
	exec(ctx *Context) Value
}

type DefineVarCommand struct {
	Name    string
	Mutable bool
	Type    parser.Type
	value   Value
}

func (c DefineVarCommand) exec(ctx *Context) Value {
	variable := Variable{
		Name:    c.Name,
		Mutable: c.Mutable,
		Type:    c.Type,
		Value:   c.value,
	}

	ctx.DefineVariable(c.Name, variable)
	return variable.Value
}

func ToCommand(statement parser.Stmt) Command {

	switch t := statement.(type) {
	case parser.VarDefStmt:
		Type := t.Type
		if Type == nil {
			Type = &parser.ElementaryTypeContract{Identifier: "Any"}
		}
		return DefineVarCommand{
			Name:    t.Identifier,
			Mutable: t.Mutable,
			Type:    Type,
			value:   Value{Type: Type, value: t.Value},
		}
	case parser.ExpressionStmt:
		return ExpressionToCommand(t.Expr)
	}

	panic("Could not handle " + reflect.TypeOf(statement).Name())
}

func ExpressionToCommand(expr parser.Expr) Command {

	switch t := expr.(type) {
	case parser.StringLiteralExpr:
		println(t.Value)
	}

	panic("Could not handle " + reflect.TypeOf(expr).Name())
}
