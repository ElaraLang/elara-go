package interpreter

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

var contextPool = sync.Pool{
	New: func() interface{} {
		return &Context{
			variables:   map[uint64][]*Variable{},
			parameters:  []*Value{},
			namespace:   "",
			name:        "",
			contextPath: map[string][]*Context{},
			types:       map[string]Type{},
			parent:      nil,
			function:    nil,
			extensions:  map[Type]map[string]*Extension{},
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
	fromPool.function = c.function
	fromPool.extensions = c.extensions
	return fromPool
}

func (c *Context) Cleanup() {
	c.function = nil

	for s := range c.variables {
		delete(c.variables, s)
	}
	c.variables = map[uint64][]*Variable{}
	c.parameters = []*Value{}

	c.namespace = ""
	c.name = ""
	c.contextPath = map[string][]*Context{}
	c.types = map[string]Type{}
	c.extensions = map[Type]map[string]*Extension{}
	c.parent = nil
	contextPool.Put(c)
}

func NewContext(init bool) *Context {
	c := contextPool.Get().(*Context)
	if !init {
		return c
	}
	c.DefineVariable(&Variable{
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
		Body: NewAbstractCommand(func(ctx *Context) *ReturnedValue {
			var input string
			_, err := fmt.Scanln(&input)
			if err != nil {
				panic(err)
			}

			return NonReturningValue(&Value{Value: input, Type: StringType})
		}),
	}

	inputContract := NewFunctionType(&inputFunction)

	c.DefineVariable(&Variable{
		Name:    inputFunctionName,
		Mutable: false,
		Type:    inputContract,
		Value: &Value{
			Type:  inputContract,
			Value: inputFunction,
		},
	})

	// fetch(baseURL)

	fetchFunctionName := "fetch"
	fetchFunction := &Function{
		name: &fetchFunctionName,
		Signature: Signature{
			Parameters: []Parameter{
				{
					Name:     "this",
					Type:     StringType,
					Position: 0,
				},
			},
			ReturnType: StringType,
		},
		Body: NewAbstractCommand(func(ctx *Context) *ReturnedValue {

			var requestURL string = ctx.FindParameter(0).Value.(string)

			response, err := http.Get(requestURL)

			if err != nil {
				fmt.Print(err.Error())
				os.Exit(1)
			}

			responseData, err := ioutil.ReadAll(response.Body)
			if err != nil {
				log.Fatal(err)
			}

			return NonReturningValue(&Value{Value: string(responseData), Type: StringType})
			// return NonReturningValue(&Value{Value: fetch, Type: StringType})
		}),
	}

	fetchContract := NewFunctionType(fetchFunction)

	c.DefineVariable(&Variable{
		Name:    fetchFunctionName,
		Mutable: false,
		Type:    fetchContract,
		Value: &Value{
			Type:  fetchContract,
			Value: fetchFunction,
		},
	})

	// end of fetch()

	// setTimeout(function, ms)

	setTimeoutFunctionName := "setTimeout"
	setTimeoutFunction := &Function{
		name: &setTimeoutFunctionName,
		Signature: Signature{
			Parameters: []Parameter{
				{
					Name: "fx",
					Type: NewSignatureFunctionType ( Signature{
						Parameters: []Parameter{},
						ReturnType: UnitType,
					}),
					Position: 0,
				},
				{
					Name: "ms",
					Type: IntType,
					Position: 1,
				},
			},
			ReturnType: AnyType,
		},
		Body: NewAbstractCommand(func(ctx *Context) *ReturnedValue {

			var fxRun *Function = ctx.FindParameter(0).Value.(*Function)
			var msRun int64  = ctx.FindParameter(1).Value.(int64)

			time.Sleep(time.Duration(msRun * 1000000))

			fxRun.Exec(ctx, []*Value{})

			return NonReturningValue(UnitValue())

		}),
	}

	setTimeoutContract := NewFunctionType(setTimeoutFunction)

	c.DefineVariable(&Variable{
		Name:    setTimeoutFunctionName,
		Mutable: false,
		Type:    setTimeoutContract,
		Value: &Value{
			Type:  setTimeoutContract,
			Value: setTimeoutFunction,
		},
	})

	// end of setTimeout()

	emptyName := "empty"
	emptyFun := &Function{
		name: &emptyName,
		Signature: Signature{
			Parameters: []Parameter{},
			ReturnType: NewCollectionTypeOf(AnyType),
		},
		Body: NewAbstractCommand(func(ctx *Context) *ReturnedValue {
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

	c.DefineVariable(&Variable{
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
