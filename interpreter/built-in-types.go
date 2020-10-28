package interpreter

import "strconv"

var AnyType = EmptyType("Any")
var UnitType = EmptyType("Unit")

var IntType = EmptyType("Int")
var FloatType = EmptyType("Float")
var BooleanType = EmptyType("Boolean")
var StringType = EmptyType("String")

var types = map[string]*Type{
	"Any":  AnyType,
	"Unit": UnitType,

	"Int":     IntType,
	"Float":   FloatType,
	"Boolean": BooleanType,
	"String":  StringType,
}

func BuiltInTypeByName(name string) *Type {
	return types[name]
}

var isInitialized = false

func Init() {
	if isInitialized {
		return
	}
	isInitialized = true
	StringType.functions = map[string]Function{
		"plus": {
			Signature: Signature{
				Parameters: []Parameter{
					{
						Name: "value",
						Type: *StringType,
					},
				},
				ReturnType: *StringType,
			},
			Body: NewAbstractCommand(func(ctx *Context) *Value {
				parameter := ctx.FindParameter("value")
				concatenated := ctx.receiver.Value.(string) + parameter.Value.(string)
				return &Value{
					Type:  StringType,
					Value: concatenated,
				}
			}),
		},
		"to-int": {
			Signature: Signature{
				Parameters: []Parameter{},
				ReturnType: *IntType,
			},
			Body: NewAbstractCommand(func(ctx *Context) *Value {
				value, err := strconv.ParseInt(ctx.receiver.Value.(string), 10, 64)
				if err != nil {
					panic(err)
				}
				return &Value{
					Type:  IntType,
					Value: value,
				}
			}),
		},
		"equals": {
			Signature: Signature{
				Parameters: []Parameter{
					{
						Name: "value",
						Type: *StringType,
					},
				},
				ReturnType: *BooleanType,
			},
			Body: NewAbstractCommand(func(ctx *Context) *Value {
				parameter := ctx.FindParameter("value")
				eq := ctx.receiver.Value.(string) == parameter.Value.(string)
				return &Value{
					Type:  BooleanType,
					Value: eq,
				}
			}),
		},
	}

	BooleanType.functions = map[string]Function{
		"and": {
			Signature: Signature{
				Parameters: []Parameter{
					{
						Name: "value",
						Type: *BooleanType,
					},
				},
				ReturnType: *BooleanType,
			},
			Body: NewAbstractCommand(func(ctx *Context) *Value {
				parameter := ctx.FindParameter("value")
				and := ctx.receiver.Value.(bool) && parameter.Value.(bool)
				return &Value{
					Type:  BooleanType,
					Value: and,
				}
			}),
		},
		"not": {
			Signature: Signature{
				Parameters: []Parameter{},
				ReturnType: *BooleanType,
			},
			Body: NewAbstractCommand(func(ctx *Context) *Value {
				return &Value{
					Type:  BooleanType,
					Value: !ctx.receiver.Value.(bool),
				}
			}),
		},
	}

	intAdd := Function{
		Signature: Signature{
			Parameters: []Parameter{
				{
					Name: "value",
					Type: *IntType,
				},
			},
			ReturnType: *IntType,
		},
		Body: NewAbstractCommand(func(ctx *Context) *Value {
			parameter := ctx.FindParameter("value")
			result := ctx.receiver.Value.(int64) + parameter.Value.(int64)
			return &Value{
				Type:  IntType,
				Value: result,
			}
		}),
	}

	IntType.functions = map[string]Function{
		"plus": intAdd,
		"add":  intAdd,
		"minus": {
			Signature: Signature{
				Parameters: []Parameter{
					{
						Name: "value",
						Type: *IntType,
					},
				},
				ReturnType: *IntType,
			},
			Body: NewAbstractCommand(func(ctx *Context) *Value {
				parameter := ctx.FindParameter("value")
				result := ctx.receiver.Value.(int64) - parameter.Value.(int64)
				return &Value{
					Type:  IntType,
					Value: result,
				}
			}),
		},
		"times": {
			Signature: Signature{
				Parameters: []Parameter{
					{
						Name: "value",
						Type: *IntType,
					},
				},
				ReturnType: *IntType,
			},
			Body: NewAbstractCommand(func(ctx *Context) *Value {
				parameter := ctx.FindParameter("value")
				result := ctx.receiver.Value.(int64) * parameter.Value.(int64)
				return &Value{
					Type:  IntType,
					Value: result,
				}
			}),
		},
		"divide": {
			Signature: Signature{
				Parameters: []Parameter{
					{
						Name: "value",
						Type: *IntType,
					},
				},
				ReturnType: *IntType,
			},
			Body: NewAbstractCommand(func(ctx *Context) *Value {
				parameter := ctx.FindParameter("value")
				result := ctx.receiver.Value.(int64) / parameter.Value.(int64)
				return &Value{
					Type:  IntType,
					Value: result,
				}
			}),
		},
		"equals": {
			Signature: Signature{
				Parameters: []Parameter{
					{
						Name: "value",
						Type: *IntType,
					},
				},
				ReturnType: *BooleanType,
			},
			Body: NewAbstractCommand(func(ctx *Context) *Value {
				parameter := ctx.FindParameter("value")
				result := ctx.receiver.Value.(int64) == parameter.Value.(int64)
				return &Value{
					Type:  BooleanType,
					Value: result,
				}
			}),
		},
	}
}
