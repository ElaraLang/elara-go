package interpreter

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
func Init() {
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

	IntType.functions = map[string]Function{
		"add": {
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
				result := ctx.receiver.Value.(int) + parameter.Value.(int)
				return &Value{
					Type:  IntType,
					Value: result,
				}
			}),
		},
	}
}
