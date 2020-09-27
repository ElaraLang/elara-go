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

type VariableCommand struct {
	Variable string
}

func (c VariableCommand) exec(ctx *Context) Value {
	variable := ctx.FindVariable(c.Variable)
	if variable == nil {
		panic("No such variable " + c.Variable)
	}
	return variable.Value
}

type InvocationCommand struct {
	Invoking Command
	args     []Command
}

func (c InvocationCommand) exec(ctx *Context) Value {
	panic("implement me")
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
	case parser.VariableExpr:
		return VariableCommand{Variable: t.Identifier}
	case parser.InvocationExpr:
		fun := ExpressionToCommand(t.Invoker)
		args := make([]Command, len(t.Args))
		for _, arg := range t.Args {
			args = append(args, ExpressionToCommand(arg))
		}

		return InvocationCommand{
			Invoking: fun,
			args:     args,
		}

	}

	panic("Could not handle " + reflect.TypeOf(expr).Name())
}
