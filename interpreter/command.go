package interpreter

import (
	"elara/parser"
	"reflect"
)

type Command interface {
	Exec(ctx *Context) Value
}

type DefineVarCommand struct {
	Name    string
	Mutable bool
	Type    parser.Type
	value   Command
}

func (c DefineVarCommand) Exec(ctx *Context) Value {
	variable := Variable{
		Name:    c.Name,
		Mutable: c.Mutable,
		Type:    c.Type,
		Value:   c.value.Exec(ctx),
	}

	ctx.DefineVariable(c.Name, variable)
	return variable.Value
}

type VariableCommand struct {
	Variable string
}

func (c VariableCommand) Exec(ctx *Context) Value {
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

func (c InvocationCommand) Exec(ctx *Context) Value {
	val := c.Invoking.Exec(ctx)
	fun, ok := val.value.(Function)
	if !ok {
		panic("Cannot invoke non-function")
	}

	return fun.exec(ctx, c.args)
}

type AbstractCommand struct {
	content func(ctx *Context) Value
}

func (c AbstractCommand) Exec(ctx *Context) Value {
	return c.content(ctx)
}

func NewAbstractCommand(content func(ctx *Context) Value) *AbstractCommand {
	return &AbstractCommand{
		content: content,
	}
}

type LiteralCommand struct {
	value Value
}

func (c LiteralCommand) Exec(ctx *Context) Value {
	return c.value
}
func ToCommand(statement parser.Stmt) Command {

	switch t := statement.(type) {
	case parser.VarDefStmt:
		Type := t.Type
		if Type == nil {
			Type = &parser.ElementaryTypeContract{Identifier: "Any"}
		}
		valueExpr := ExpressionToCommand(t.Value)
		return DefineVarCommand{
			Name:    t.Identifier,
			Mutable: t.Mutable,
			Type:    Type,
			value:   valueExpr,
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
		args := make([]Command, 0)
		for _, arg := range t.Args {
			command := ExpressionToCommand(arg)
			if command == nil {
				panic("Could not convert expression " + reflect.TypeOf(arg).Name() + " to command")
			}
			args = append(args, command)
		}

		return InvocationCommand{
			Invoking: fun,
			args:     args,
		}

	case parser.StringLiteralExpr:
		str := t.Value
		value := Value{
			Type:  parser.ElementaryTypeContract{Identifier: "String"},
			value: str,
		}
		return LiteralCommand{value: value}

	case parser.IntegerLiteralExpr:
		str := t.Value
		value := Value{
			Type:  parser.ElementaryTypeContract{Identifier: "Int"},
			value: str,
		}
		return LiteralCommand{value: value}

	}

	panic("Could not handle " + reflect.TypeOf(expr).Name())
}
