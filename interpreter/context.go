package interpreter

import (
	"fmt"
	"sort"
)

type Context struct {
	variables  map[string][]*Variable
	parameters map[string]*Value
	receiver   *Value
	namespace  string
	//A map from namespace -> context slice
	contextPath map[string][]*Context

	types  map[string]Type
	scopes []Scope
}

type Scope struct {
	name    string
	context Context
}

var globalContext = &Context{
	namespace:   "__global__",
	variables:   map[string][]*Variable{},
	parameters:  map[string]*Value{},
	receiver:    nil,
	contextPath: map[string][]*Context{},

	types:  map[string]Type{},
	scopes: []Scope{},
}

func (c *Context) Init(namespace string) {
	if c.namespace != "" {
		panic("Context has already been initialized!")
	}
	c.namespace = namespace
	globalContext.contextPath[c.namespace] = append(globalContext.contextPath[c.namespace], c)
	for _, t := range types {
		c.types[t.Name] = *t
	}
}

func NewContext() *Context {
	c := &Context{
		variables:   map[string][]*Variable{},
		parameters:  map[string]*Value{},
		receiver:    nil,
		namespace:   "",
		contextPath: map[string][]*Context{},
		types:       map[string]Type{},
		scopes:      []Scope{},
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
				panic(err)
			}

			return &Value{Value: input, Type: StringType}
		}),
	}

	inputContract := FunctionType(&inputFunction)

	c.DefineVariable(inputFunctionName, Variable{
		Name:    inputFunctionName,
		Mutable: false,
		Type:    *inputContract,
		Value: &Value{
			Type:  inputContract,
			Value: inputFunction,
		},
	})

	Init(c)
	return c
}

func (c *Context) DefineVariable(name string, value Variable) {

	if len(c.scopes) == 0 {
		vars := c.variables[name]
		vars = append(vars, &value)
		c.variables[name] = vars
		return
	}

	c.scopes[len(c.scopes)-1].context.DefineVariable(name, value)
}

func (c *Context) FindVariable(name string) *Variable {
	return c.FindVariableMaxDepth(name, 0)
}

func (c *Context) FindVariableMaxDepth(name string, maxDepth int) *Variable {
	if len(c.scopes) == 0 {
		vars := c.variables[name]
		if vars != nil {
			return vars[len(vars)-1]
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

	for i := len(c.scopes) - 1; i >= maxDepth; i-- {
		last := c.scopes[i]
		found := last.context.FindVariableMaxDepth(name, 0)
		if found != nil {
			return found
		}
	}
	return nil
}

func (c *Context) DefineParameter(name string, value *Value) {
	if len(c.scopes) == 0 {
		c.parameters[name] = value
		return
	}
	c.scopes[len(c.scopes)-1].context.DefineParameter(name, value)
}

func (c *Context) FindParameter(name string) *Value {
	if len(c.scopes) == 0 {
		return c.parameters[name]
	}

	for i := len(c.scopes) - 1; i >= 0; i-- {
		last := c.scopes[i]
		found := last.context.FindParameter(name)
		if found != nil {
			return found
		}
	}
	return nil
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

func (c *Context) EnterScope(name string) {
	scope := Scope{
		name:    name,
		context: *c,
	}
	c.scopes = append(c.scopes, scope)
}
func (c *Context) ExitScope() {
	//Drop the last scope
	c.scopes = c.scopes[:len(c.scopes)-1]
}

func (c *Context) FindConstructor(name string) *Value {

	t := c.FindType(name)
	if t == nil {
		return nil
	}
	constructorParams := make([]Parameter, 0)
	for _, key := range t.variables.keys {
		v := t.variables.m[key]
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
		Type:  FunctionType(&constructor),
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

func (c *Context) FindFunction(name string, signature Signature) *Function {
	vars := c.variables[name]
	functions := make([]Function, 0)
	distances := make([]int, 0)

	for _, v := range vars {
		fun, ok := v.Value.Value.(Function)
		if ok {
			if !fun.Signature.Accepts(signature) {
				continue
			}
			functions = append(functions, fun)

			distance := signature.Distance(fun.Signature)
			distances = append(distances, distance)
		}
	}

	sort.Slice(functions, func(i, j int) bool {
		return distances[i] > distances[j]
	})

	if len(functions) >= 1 {
		return &functions[0]
	}

	for i := range c.scopes {
		fun := c.scopes[i].context.FindFunction(name, signature)
		if fun != nil {
			return fun
		}
	}
	for i := range c.contextPath {
		contexts := c.contextPath[i]
		for j := range contexts {
			ctx := contexts[j]
			fun := ctx.FindFunction(name, signature)
			if fun != nil {
				return fun
			}
		}
	}

	return nil
}
