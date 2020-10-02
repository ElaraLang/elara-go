package interpreter

import (
	"elara/interpreter/types"
	"elara/util"
	"fmt"
)

type Context struct {
	variables map[string][]*Variable

	parameters []*Value
}

func NewContext() *Context {
	c := &Context{
		variables:  map[string][]*Variable{},
		parameters: []*Value{},
	}

	//Todo remove
	printContract := types.FunctionType{
		Params: []types.Type{types.AnyType},
		Output: types.UnitType,
	}
	c.DefineVariable("print", Variable{
		Name:    "print",
		Mutable: false,
		Type:    printContract,
		Value: Value{
			Type: printContract,
			Value: Function{
				Signature: printContract,
				body: []Command{
					NewAbstractCommand(func(ctx *Context) Value {
						value := ctx.FindParameter(0).Value
						fmt.Printf("%s\n", util.Stringify(value))

						return *UnitValue()
					}),
				},
			},
		},
	})
	return c
}

func (c Context) DefineVariable(name string, value Variable) {
	vars := c.variables[name]
	vars = append(vars, &value)
	c.variables[name] = vars
}

func (c Context) FindVariable(name string) *Variable {
	vars := c.variables[name]

	if vars == nil {
		return nil
	}
	return vars[0]
}

func (c *Context) DefineParameter(index int, value *Value) {
	if index >= len(c.parameters) {
		newParameters := make([]*Value, index+1)
		newParameters[index] = value
		c.parameters = newParameters
		return
	}
	c.parameters[index] = value
}

func (c Context) FindParameter(index int) *Value {
	return c.parameters[index]
}

func (c Context) string() string {
	s := ""
	for key, values := range c.variables {
		s += fmt.Sprintf("%s = [\n", key)
		for _, val := range values {
			s += fmt.Sprintf("%s \n", val.string())
		}
		s += "]\n"
	}

	return s
}
