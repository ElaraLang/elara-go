package interpreter

var AnyType = EmptyType("Any")
var UnitType = EmptyType("Unit")

var IntType = EmptyType("Int")
var FloatType = EmptyType("Float")
var BooleanType = EmptyType("Boolean")
var StringType = EmptyType("String")

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
	}
}
