package interpreter

import (
	"fmt"
	"github.com/ElaraLang/elara/util"
	"math"
)

type Context struct {
	variables  map[uint64][]*Variable
	parameters []*Value
	namespace  string
	name       string //The optional name of the context - may be empty

	//A map from namespace -> context slice
	contextPath map[string][]*Context

	extensions map[Type]map[string]*Extension
	types      map[string]Type
	parent     *Context
	function   *Function //Will only be nil if this is a Function scope
}

var globalContext = &Context{
	namespace:   "__global__",
	name:        "__global__",
	variables:   map[uint64][]*Variable{},
	extensions:  map[Type]map[string]*Extension{},
	parameters:  []*Value{},
	contextPath: map[string][]*Context{},

	types:    map[string]Type{},
	function: nil,
}

func (c *Context) Init(namespace string) {
	if c.namespace != "" {
		panic("Context has already been initialized!")
	}
	c.namespace = namespace
	globalContext.contextPath[c.namespace] = append(globalContext.contextPath[c.namespace], c)
}

func (c *Context) DefineVariableWithHash(hash uint64, value *Variable) {
	vars := c.variables[hash]
	vars = append(vars, value)
	c.variables[hash] = vars
}
func (c *Context) DefineVariable(value *Variable) {
	hash := util.Hash(value.Name)
	c.DefineVariableWithHash(hash, value)
}

func (c *Context) FindFunction(hash uint64, signature *Signature) *Function {
	vars := c.variables[hash]
	if vars != nil {
		matching := make([]*Variable, 0)
		for _, variable := range vars {
			asFunction, isFunction := variable.Value.Value.(*Function)
			if isFunction {
				if asFunction.Signature.Accepts(signature, c, false) {
					matching = append(matching, variable)
				}
			}
		}
		if len(matching) > 1 {
			_ = fmt.Errorf("multiple matching functions with name %s and signature %s", matching[0].Name, signature.String())
		}
		if len(matching) != 0 {
			return matching[0].Value.Value.(*Function)
		}
	}

	if c.parent != nil {
		parFound := c.parent.FindFunction(hash, signature)
		if parFound != nil {
			return parFound
		}
	}
	for _, contexts := range c.contextPath {
		for _, context := range contexts {
			v := context.FindFunction(hash, signature)
			if v != nil {
				return v
			}
		}
	}
	return nil
}
func (c *Context) FindVariable(hash uint64) *Variable {
	variable, _ := c.FindVariableMaxDepth(hash, -1)
	return variable
}

//TODO this needs optimising, it's a MASSIVE hotspot
func (c *Context) FindVariableMaxDepth(hash uint64, maxDepth int) (*Variable, int) {
	vars := c.variables[hash]
	if vars != nil {
		return vars[len(vars)-1], 0
	}

	i := 0
	for {
		if c.parent == nil {
			break
		}
		parentVar, depth := c.parent.FindVariableMaxDepth(hash, int(math.Max(float64(maxDepth-i), -1)))
		if parentVar != nil {
			if maxDepth < 0 || c.parent.function == nil && i+depth <= maxDepth {
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
			v, _ := context.FindVariableMaxDepth(hash, maxDepth-i)
			if v != nil {
				return v, i //0 for a variable from an import?
			}
		}
	}

	return nil, -1
}

func (c *Context) DefineParameter(pos uint, value *Value) {
	c.parameters[pos] = value
}

func (c *Context) FindParameter(pos uint) *Value {
	par := c.parameters[pos]
	if par != nil {
		return par
	}
	return nil
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

func (c *Context) EnterScope(name string, function *Function, paramLength uint) *Context {
	scope := NewContext(false)
	scope.parent = c
	scope.namespace = c.name
	scope.name = name
	scope.contextPath = c.contextPath
	scope.function = function
	scope.parameters = make([]*Value, paramLength)
	scope.extensions = c.extensions
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
	if asStruct.constructor != nil {
		return asStruct.constructor
	}

	constructorParams := make([]Parameter, 0)
	i := uint(0)
	for _, v := range asStruct.Properties {
		if v.DefaultValue == nil {
			constructorParams = append(constructorParams, Parameter{
				Position: i,
				Name:     v.Name,
				Type:     v.Type,
			})
		}
		i++
	}

	constructor := &Function{
		Signature: Signature{
			Parameters: constructorParams,
			ReturnType: t,
		},
		Body: NewAbstractCommand(func(ctx *Context) *ReturnedValue {
			values := make(map[string]*Value, len(constructorParams))
			for _, param := range constructorParams {
				values[param.Name] = ctx.FindParameter(param.Position)
			}
			return NonReturningValue(&Value{
				Type: t,
				Value: &Instance{
					Type:   asStruct,
					Values: values,
				},
			})
		}),
		name: &name,
	}

	constructorVal := &Value{
		Type:  NewFunctionType(constructor),
		Value: constructor,
	}
	asStruct.constructor = constructorVal
	return constructorVal
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
			s += fmt.Sprintf("%s \n", val.String())
		}
		s += "]\n"
	}

	return s
}

func (c *Context) Stringify(value *Value) string {
	if value == nil {
		return "<empty value>"
	}

	return util.Stringify(value.Value)

}

func (c *Context) DefineExtension(receiverType Type, name string, value *Extension) {
	extensions, present := c.extensions[receiverType]
	if !present {
		extensions = map[string]*Extension{}
	}
	_, exists := extensions[name]
	if exists {
		panic("Extension on " + receiverType.Name() + " with name " + name + " already exists.")
	}
	extensions[name] = value
	c.extensions[receiverType] = extensions
}

func (c *Context) FindExtension(receiverType Type, name string) *Extension {
	extensions, present := c.extensions[receiverType]
	if !present {
		extensions = map[string]*Extension{}
	}
	return extensions[name]
}
