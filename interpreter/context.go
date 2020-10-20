package interpreter

import (
	"elara/util"
	"fmt"
)

type Context struct {
	variables map[string][]*Variable

	parameters map[string]*Value

	receiver *Value
}

func NewContext() *Context {
	c := &Context{
		variables:  map[string][]*Variable{},
		parameters: map[string]*Value{},
		receiver:   nil,
	}

	//Todo remove
	function := Function{
		Signature: Signature{
			Parameters: []Parameter{
				{
					Name: "value",
					Type: *AnyType,
				},
			},
			ReturnType: *UnitType,
		},
		Body: NewAbstractCommand(func(ctx *Context) Value {
			value := ctx.FindParameter("value").Value
			fmt.Printf("%s\n", util.Stringify(value))

			return *UnitValue()
		}),
	}
	funName := "print"
	printContract := FunctionType(&funName, function)

	c.DefineVariable(funName, Variable{
		Name:    funName,
		Mutable: false,
		Type:    *printContract,
		Value: Value{
			Type:  printContract,
			Value: function,
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

func (c *Context) DefineParameter(name string, value *Value) {
	c.parameters[name] = value
}

func (c Context) FindParameterIndexed(index int) *Value {
	i := 0
	for _, value := range c.parameters {
		if i == index {
			return value
		}
		i++
	}
	return nil
}

func (c Context) FindParameter(name string) *Value {
	return c.parameters[name]
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
