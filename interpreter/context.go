package interpreter

import (
	"fmt"
)

type Context struct {
	variables  map[string][]*Variable
	parameters map[string]*Value
	receiver   *Value
	namespace  string
	//A map from namespace -> context slice
	contextPath map[string][]*Context

	types map[string]Type
}

var globalContext = &Context{
	namespace:   "__global__",
	variables:   map[string][]*Variable{},
	parameters:  map[string]*Value{},
	receiver:    nil,
	contextPath: map[string][]*Context{},

	types: map[string]Type{},
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
		types:       map[string]Type{},
	}
	c.DefineVariable("stdout", Variable{
		Name:    "stdout",
		Mutable: false,
		Type:    *OutputType,
		Value: &Value{
			Type:  OutputType,
			Value: nil,
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

	inputContract := FunctionType(inputFunction)

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

func (c *Context) DefineVariable(name string, value Variable) {
	vars := c.variables[name]
	vars = append(vars, &value)
	c.variables[name] = vars
}

func (c *Context) FindVariable(name string) *Variable {
	vars := c.variables[name]
	if vars != nil {
		return vars[0]
	}

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

func (c *Context) DefineParameter(name string, value *Value) {
	c.parameters[name] = value
}

func (c *Context) FindParameterIndexed(index int) *Value {
	i := 0
	for _, value := range c.parameters {
		if i == index {
			return value
		}
		i++
	}
	return nil
}

func (c *Context) FindParameter(name string) *Value {
	return c.parameters[name]
}

func (c *Context) FindType(name string) *Type {
	t, ok := c.types[name]
	if ok {
		return &t
	}
	for _, contexts := range c.contextPath {
		for _, context := range contexts {
			t := context.FindType(name)
			if t != nil {
				return t
			}
		}
	}
	return nil
}

func (c *Context) FindConstructor(name string) *Value {

	t := c.FindType(name)
	if t == nil {
		return nil
	}
	constructorParams := make([]Parameter, 0)
	for _, v := range t.variables {
		if v.Value == nil {
			constructorParams = append(constructorParams, Parameter{
				Name: v.Name,
				Type: v.Type,
			})
		}
	}

	constructor := Function{
		Signature: Signature{
			Parameters: constructorParams,
			ReturnType: *t,
		},
		Body: NewAbstractCommand(func(ctx *Context) *Value {
			values := make(map[string]*Value, len(constructorParams))
			for _, param := range constructorParams {
				values[param.Name] = ctx.FindParameter(param.Name)
			}
			return &Value{
				Type: t,
				Value: Instance{
					Type:   t,
					Values: values,
				},
			}

		}),
		name: &name,
	}

	return &Value{
		Type:  FunctionType(constructor),
		Value: constructor,
	}
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

func (c *Context) string() string {
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
