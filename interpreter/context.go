package interpreter

import (
	"elara/util"
	"fmt"
)

type Context struct {
	variables  map[string][]*Variable
	parameters map[string]*Value
	receiver   *Value
	namespace  string
	//A map from namespace -> context slice
	contextPath map[string][]*Context
}

var globalContext = &Context{
	namespace:   "__global__",
	variables:   map[string][]*Variable{},
	parameters:  map[string]*Value{},
	receiver:    nil,
	contextPath: map[string][]*Context{},
}

func (c *Context) Init(namespace string) {
	if c.namespace != "" {
		panic("Context has already been initialized!")
	}
	c.namespace = namespace
	globalContext.contextPath[c.namespace] = append(globalContext.contextPath[c.namespace], c)
}

func NewContext() *Context {
	c := &Context{
		variables:   map[string][]*Variable{},
		parameters:  map[string]*Value{},
		receiver:    nil,
		namespace:   "",
		contextPath: map[string][]*Context{},
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

	inputFunctionName := "input"
	inputFunction := Function{
		name: &inputFunctionName,
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
		for _, contexts := range c.contextPath {
			for _, context := range contexts {
				v := context.FindVariable(name)
				if v != nil {
					return v
				}
			}
		}
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

func (c *Context) Import(namespace string) {
	contexts := globalContext.contextPath[namespace]
	if contexts == nil {
		panic("Nothing found in namespace " + namespace)
	}
	ns := c.contextPath[namespace]
	if ns == nil {
		c.contextPath[namespace] = contexts
		return
	}
	ns = append(ns, contexts...)
	c.contextPath[namespace] = ns
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
