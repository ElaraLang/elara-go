package interpreter

import (
	"elara/parser"
	"fmt"
)

type Value struct {
	Type  parser.Type
	value interface{}
}

type Variable struct {
	Name    string
	Mutable bool
	Type    parser.Type
	Value   Value
}

type Function struct {
	Signature parser.InvocableTypeContract
	body      []Command
}

func (f *Function) String() string {
	return fmt.Sprintf("Function %s => %s", f.Signature.Args, f.Signature.ReturnType)
}

func (f *Function) exec(ctx *Context, parameters []Command) Value {
	var val Value
	for i, parameter := range parameters {
		paramValue := parameter.Exec(ctx)
		ctx.DefineParameter(i, &paramValue)
	}
	for _, line := range f.body {
		val = line.Exec(ctx)
	}

	return val
}

type Signature struct {
	parameters []parser.Type
	ReturnType parser.Type
}

func (v Variable) string() string {
	return fmt.Sprintf("Variable { name: %s, mutable: %T, type: %s, value: %s", v.Name, v.Mutable, v.Type, v.Value)
}
