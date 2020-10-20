package interpreter

import (
	"elara/interpreter/types"
	"elara/parser"
	"fmt"
)

type Value struct {
	Type  types.Type
	Value interface{}
}

var unitValue = Value{
	Type:  types.UnitType,
	Value: "Unit",
}

func UnitValue() *Value {
	return &unitValue
}

type Variable struct {
	Name    string
	Mutable bool
	Type    types.Type
	Value   Value
}

type Function struct {
	Signature types.FunctionType
	body      Command
}

func (f *Function) String() string {
	return fmt.Sprintf("Function %s => %s", f.Signature.Params, f.Signature.Output)
}

func (f *Function) exec(ctx *Context, parameters []Command) Value {

	for i, parameter := range parameters {
		paramValue := parameter.Exec(ctx)
		ctx.DefineParameter(f.Signature.Params[i].Name, &paramValue)
	}

	return f.body.Exec(ctx)
}

type Signature struct {
	parameters []parser.Type
	ReturnType parser.Type
}

func (v Variable) string() string {
	return fmt.Sprintf("Variable { name: %s, mutable: %T, type: %s, Value: %s", v.Name, v.Mutable, v.Type, v.Value)
}
