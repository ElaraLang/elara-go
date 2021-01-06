package interpreter

import (
	"fmt"
	"github.com/ElaraLang/elara/util"
)

type Context struct {
	variables  map[string][]*Variable
	parameters map[string]*Value
	namespace  string
	name       string //The optional name of the context - may be empty

	//A map from namespace -> context slice
	contextPath map[string][]*Context

	types           map[string]Type
	parent          *Context
	isFunctionScope bool
}

var globalContext = &Context{
	namespace:   "__global__",
	name:        "__global__",
	variables:   map[string][]*Variable{},
	parameters:  map[string]*Value{},
	contextPath: map[string][]*Context{},

	types:           map[string]Type{},
	isFunctionScope: false,
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
		variables:       map[string][]*Variable{},
		parameters:      map[string]*Value{},
		namespace:       "",
		name:            "",
		contextPath:     map[string][]*Context{},
		types:           map[string]Type{},
		isFunctionScope: false,
	}
	c.DefineVariable("stdout", Variable{
		Name:    "stdout",
		Mutable: false,
		Type:    OutputType,
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
			ReturnType: StringType,
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

	inputContract := NewFunctionType(&inputFunction)

	c.DefineVariable(inputFunctionName, Variable{
		Name:    inputFunctionName,
		Mutable: false,
		Type:    inputContract,
		Value: &Value{
			Type:  inputContract,
			Value: inputFunction,
		},
	})

	emptyName := "empty"
	emptyFun := &Function{
		name: &emptyName,
		Signature: Signature{
			Parameters: []Parameter{},
			ReturnType: NewCollectionTypeOf(AnyType),
		},
		Body: NewAbstractCommand(func(ctx *Context) *Value {
			return &Value{
				Type: NewCollectionTypeOf(AnyType),
				Value: &Collection{
					ElementType: AnyType,
					Elements:    []*Value{},
				},
			}
		}),
	}
	emptyContract := NewFunctionType(emptyFun)

	c.DefineVariable(emptyName, Variable{
		Name:    emptyName,
		Mutable: false,
		Type:    emptyContract,
		Value: &Value{
			Type:  emptyContract,
			Value: emptyFun,
		},
	})

	Init(c)
	return c
}

func (c *Context) DefineVariable(name string, value Variable) {

	vars := c.variables[name]
	vars = append(vars, &value)
	c.variables[name] = vars
}

func (c *Context) FindFunction(name string, signature *Signature) *Function {
	vars := c.variables[name]
	if vars != nil {
		matching := make([]*Variable, 0)
		for _, variable := range vars {
			asFunction, isFunction := variable.Value.Value.(*Function)
			if isFunction {
				if asFunction.Signature.Accepts(signature, false) {
					matching = append(matching, variable)
				}
			}
		}
		if len(matching) > 1 {
			_ = fmt.Errorf("multiple matching functions with name %s and signature %s", name, signature.String())
		}
		if len(matching) != 0 {
			return matching[0].Value.Value.(*Function)
		}
	}

	if c.parent != nil {
		parFound := c.parent.FindFunction(name, signature)
		if parFound != nil {
			return parFound
		}
	}
	for _, contexts := range c.contextPath {
		for _, context := range contexts {
			v := context.FindFunction(name, signature)
			if v != nil {
				return v
			}
		}
	}
	return nil
}
func (c *Context) FindVariable(name string) *Variable {
	variable, _ := c.FindVariableMaxDepth(name, -1)
	return variable
}

func (c *Context) FindVariableMaxDepth(name string, maxDepth int) (*Variable, int) {
	vars := c.variables[name]
	if vars != nil {
		return vars[len(vars)-1], 0
	}

	i := 0
	for {
		if c.parent == nil {
			break
		}
		parentVar, depth := c.parent.FindVariableMaxDepth(name, maxDepth-i)
		if parentVar != nil {
			if maxDepth == -1 || !c.parent.isFunctionScope && i+depth <= maxDepth {
				return parentVar, i + depth
			} else {
				return nil, i + depth
			}
		}
		i++
		if i >= maxDepth {
			break
		}
	}

	for _, contexts := range c.contextPath {
		for _, context := range contexts {
			v, _ := context.FindVariableMaxDepth(name, maxDepth-i)
			if v != nil {
				return v, 0 //0 for a variable from an import?
			}
		}
	}

	return nil, -1
}

func (c *Context) DefineParameter(name string, value *Value) {
	c.parameters[name] = value
}

func (c *Context) FindParameter(name string) *Value {
	par := c.parameters[name]
	if par != nil {
		return par
	}
	if c.parent == nil {
		return nil
	}
	return c.parent.FindParameter(name)
}

func (c *Context) FindType(name string) Type {
	t, ok := c.types[name]
	if ok {
		return t
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

func (c *Context) EnterScope(name string) *Context {
	scope := NewContext()
	scope.parent = c
	scope.namespace = c.name
	scope.name = name
	scope.contextPath = c.contextPath
	scope.isFunctionScope = true
	return scope
}

func (c *Context) FindConstructor(name string) *Value {

	t := c.FindType(name)
	if t == nil {
		return nil
	}
	asStruct, isStruct := t.(*StructType)
	if !isStruct {
		panic("Cannot construct non struct type")
	}
	constructorParams := make([]Parameter, 0)
	for _, v := range asStruct.Properties {
		if v.DefaultValue == nil {
			constructorParams = append(constructorParams, Parameter{
				Name: v.Name,
				Type: v.Type,
			})
		}
	}

	constructor := Function{
		Signature: Signature{
			Parameters: constructorParams,
			ReturnType: t,
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
		Type:  NewFunctionType(&constructor),
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

func (c *Context) Clone() *Context {
	var parentClone *Context = nil
	if c.parent != nil {
		parentClone = c.parent.Clone()
	}
	return &Context{
		variables:       c.variables,
		parameters:      c.parameters,
		namespace:       c.namespace,
		name:            c.name,
		contextPath:     c.contextPath,
		types:           c.types,
		parent:          parentClone,
		isFunctionScope: c.isFunctionScope,
	}
}

func (c *Context) Stringify(value *Value) string {
	if value == nil {
		return "<empty value>"
	}
	toString := c.FindFunction("toString", &Signature{
		Parameters: []Parameter{
			{Name: "this",
				Type: value.Type,
			},
		},
		ReturnType: StringType,
	})
	if toString != nil {
		asString := toString.Exec(c, []*Value{value})
		return asString.Value.(string)
	} else {
		return util.Stringify(value.Value)
	}
}
