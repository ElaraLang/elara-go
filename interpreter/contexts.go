package interpreter

import (
	"fmt"
	"sync"
)

var contextPool = sync.Pool{
	New: func() interface{} {
		return &Context{
			variables:       map[string][]*Variable{},
			parameters:      map[string]*Value{},
			namespace:       "",
			name:            "",
			contextPath:     map[string][]*Context{},
			types:           map[string]Type{},
			parent:          nil,
			isFunctionScope: false,
		}
	},
}

func (c *Context) Clone() *Context {
	var parentClone *Context = nil
	if c.parent != nil {
		parentClone = c.parent.Clone()
	}
	fromPool := contextPool.Get().(*Context)
	fromPool.variables = c.variables
	fromPool.parameters = c.parameters
	fromPool.namespace = c.namespace
	fromPool.name = c.name
	fromPool.contextPath = c.contextPath
	fromPool.types = c.types
	fromPool.parent = parentClone
	fromPool.isFunctionScope = c.isFunctionScope
	return fromPool
}

func (c *Context) Cleanup() {
	c.isFunctionScope = false
	c.variables = map[string][]*Variable{}
	c.parameters = map[string]*Value{}
	c.namespace = ""
	c.name = ""
	c.contextPath = map[string][]*Context{}
	c.types = map[string]Type{}
	c.parent = nil
	contextPool.Put(c)
}

func NewContext(init bool) *Context {
	c := contextPool.Get().(*Context)
	if !init {
		return c
	}
	c.DefineVariable("stdout", &Variable{
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
		Body: NewAbstractCommand(func(ctx *Context) ReturnedValue {
			var input string
			_, err := fmt.Scanln(&input)
			if err != nil {
				panic(err)
			}

			return NonReturningValue(&Value{Value: input, Type: StringType})
		}),
	}

	inputContract := NewFunctionType(&inputFunction)

	c.DefineVariable(inputFunctionName, &Variable{
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
		Body: NewAbstractCommand(func(ctx *Context) ReturnedValue {
			return NonReturningValue(&Value{
				Type: NewCollectionTypeOf(AnyType),
				Value: &Collection{
					ElementType: AnyType,
					Elements:    []*Value{},
				},
			})
		}),
	}
	emptyContract := NewFunctionType(emptyFun)

	c.DefineVariable(emptyName, &Variable{
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
