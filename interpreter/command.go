package interpreter

import (
	"elara/interpreter/types"
	"elara/lexer"
	"elara/parser"
	"reflect"
)

type Command interface {
	Exec(ctx *Context) Value
}

type DefineVarCommand struct {
	Name    string
	Mutable bool
	Type    types.Type
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
		param := ctx.FindParameter(c.Variable)
		if param == nil {
			panic("No such variable or parameter " + c.Variable)
		}
		return *param
	}
	return variable.Value
}

type InvocationCommand struct {
	Invoking Command
	args     []Command
}

func (c InvocationCommand) Exec(ctx *Context) Value {
	val := c.Invoking.Exec(ctx)
	fun, ok := val.Value.(Function)
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

func (c LiteralCommand) Exec(_ *Context) Value {
	return c.value
}

type BinaryOperatorCommand struct {
	lhs Command
	op  func(ctx *Context, lhs Value, rhs Value) Value
	rhs Command
}

func (c BinaryOperatorCommand) Exec(ctx *Context) Value {
	lhs := c.lhs.Exec(ctx)
	rhs := c.rhs.Exec(ctx)

	return c.op(ctx, lhs, rhs)
}

type BlockCommand struct {
	lines []Command
}

func (c *BlockCommand) Exec(ctx *Context) Value {
	var last Value
	for _, line := range c.lines {
		last = line.Exec(ctx)
	}
	return last
}

func ToCommand(statement parser.Stmt) Command {

	switch t := statement.(type) {
	case parser.VarDefStmt:
		Type := types.FromASTType(t.Type)
		if Type == nil {
			Type = types.AnyType
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

	case parser.BlockStmt:
		commands := make([]Command, len(t.Stmts))
		for i, stmt := range t.Stmts {
			commands[i] = ToCommand(stmt)
		}
		return &BlockCommand{lines: commands}
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
			Type:  types.StringType,
			Value: str,
		}
		return LiteralCommand{value: value}

	case parser.IntegerLiteralExpr:
		str := t.Value
		value := Value{
			Type:  types.IntType,
			Value: str,
		}
		return LiteralCommand{value: value}
	case parser.FloatLiteralExpr:
		str := t.Value
		value := Value{
			Type:  types.FloatType,
			Value: str,
		}
		return LiteralCommand{value: value}
	case parser.BinaryExpr:
		lhs := t.Lhs
		lhsCmd := ExpressionToCommand(lhs)
		op := t.Op
		rhs := t.Rhs
		rhsCmd := ExpressionToCommand(rhs)

		switch op {
		case lexer.Add:
			return BinaryOperatorCommand{
				lhs: lhsCmd,
				op: func(ctx *Context, lhs Value, rhs Value) Value {
					left := lhs.Value
					lhsInt, ok := left.(int64)
					if !ok {
						panic("LHS must be an int64")
					}
					rhsInt, ok := rhs.Value.(int64)
					if !ok {
						panic("RHS must be an int64")
					}

					return Value{
						Type:  types.IntType,
						Value: lhsInt + rhsInt,
					}
				},
				rhs: rhsCmd,
			}
		case lexer.Subtract:
			return BinaryOperatorCommand{
				lhs: lhsCmd,
				op: func(ctx *Context, lhs Value, rhs Value) Value {
					left := lhs.Value
					lhsInt, ok := left.(int64)
					if !ok {
						panic("LHS must be an int64")
					}
					rhsInt, ok := rhs.Value.(int64)
					if !ok {
						panic("RHS must be an int64")
					}

					return Value{
						Type:  types.IntType,
						Value: lhsInt - rhsInt,
					}
				},
				rhs: rhsCmd,
			}
		case lexer.Multiply:
			return BinaryOperatorCommand{
				lhs: lhsCmd,
				op: func(ctx *Context, lhs Value, rhs Value) Value {
					left := lhs.Value
					lhsInt, ok := left.(int64)
					if !ok {
						panic("LHS must be an int64")
					}
					rhsInt, ok := rhs.Value.(int64)
					if !ok {
						panic("RHS must be an int64")
					}

					return Value{
						Type:  types.IntType,
						Value: lhsInt * rhsInt,
					}
				},
				rhs: rhsCmd,
			}
		case lexer.Slash:
			return BinaryOperatorCommand{
				lhs: lhsCmd,
				op: func(ctx *Context, lhs Value, rhs Value) Value {
					left := lhs.Value
					lhsInt, ok := left.(int64)
					if !ok {
						panic("LHS must be an int64")
					}
					rhsInt, ok := rhs.Value.(int64)
					if !ok {
						panic("RHS must be an int64")
					}

					return Value{
						Type:  types.IntType,
						Value: lhsInt / rhsInt,
					}
				},
				rhs: rhsCmd,
			}
		}
	case parser.FuncDefExpr:
		params := make([]types.Parameter, len(t.Arguments))
		for i, parameter := range t.Arguments {
			paramType := types.FromASTType(parameter.Type)
			params[i] = types.Parameter{
				Type: paramType,
				Name: parameter.Name,
			}
		}

		functionType := types.FunctionType{
			Params: params,
			Output: types.FromASTType(t.ReturnType),
		}

		fun := Function{
			Signature: functionType,
			body:      ToCommand(t.Statement),
		}
		return &LiteralCommand{value: Value{
			Type:  functionType,
			Value: fun,
		}}
	}

	panic("Could not handle " + reflect.TypeOf(expr).Name())
}
