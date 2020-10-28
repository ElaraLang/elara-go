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
	printFunction := Function{
		Signature: Signature{
			Parameters: []Parameter{
				{
					Name: "value",
					Type: *AnyType,
				},
			},
			ReturnType: *UnitType,
		},
		Body: NewAbstractCommand(func(ctx *Context) *Value {
			value := ctx.FindParameter("value").Value
			fmt.Printf("%s\n", util.Stringify(value))

			return UnitValue()
		}),
	}
	printFunctionName := "print"
	printContract := FunctionType(&printFunctionName, printFunction)

	c.DefineVariable(printFunctionName, Variable{
		Name:    printFunctionName,
		Mutable: false,
		Type:    *printContract,
		Value: &Value{
			Type:  printContract,
			Value: printFunction,
		},
	})

	inputFunction := Function{
		Signature: Signature{
			Parameters: []Parameter{},
			ReturnType: *StringType,
		},
		Body: NewAbstractCommand(func(ctx *Context) *Value {
			var input string
			_, err := fmt.Scanln(&input)
			if err != nil {
				println(err)
			}

			return &Value{Value: input, Type: StringType}
		}),
	}
	inputFunctionName := "input"
	inputContract := FunctionType(&inputFunctionName, printFunction)

	c.DefineVariable(inputFunctionName, Variable{
		Name:    inputFunctionName,
		Mutable: false,
		Type:    *inputContract,
		Value: &Value{
			Type:  inputContract,
			Value: inputFunction,
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
