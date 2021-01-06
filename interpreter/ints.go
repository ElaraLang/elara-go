package interpreter

var IntType = NewEmptyType("Int")

func InitInts(ctx *Context) {
	define(ctx, "plus", &Function{
		Signature: Signature{
			Parameters: []Parameter{
				{
					Name: "this",
					Type: IntType,
				},
				{
					Name: "value",
					Type: IntType,
				},
			},
			ReturnType: IntType,
		},
		Body: NewAbstractCommand(func(ctx *Context) ReturnedValue {
			this := ctx.FindParameter("this").Value.(int64)
			value := ctx.FindParameter("value").Value.(int64)

			return NonReturningValue(IntValue(this + value))
		}),
	})

	define(ctx, "minus", &Function{
		Signature: Signature{
			Parameters: []Parameter{
				{
					Name: "this",
					Type: IntType,
				},
				{
					Name: "value",
					Type: IntType,
				},
			},
			ReturnType: IntType,
		},
		Body: NewAbstractCommand(func(ctx *Context) ReturnedValue {
			this := ctx.FindParameter("this").Value.(int64)
			value := ctx.FindParameter("value").Value.(int64)

			return NonReturningValue(IntValue(this - value))
		}),
	})

	define(ctx, "times", &Function{
		Signature: Signature{
			Parameters: []Parameter{
				{
					Name: "this",
					Type: IntType,
				},
				{
					Name: "value",
					Type: IntType,
				},
			},
			ReturnType: IntType,
		},
		Body: NewAbstractCommand(func(ctx *Context) ReturnedValue {
			this := ctx.FindParameter("this").Value.(int64)
			value := ctx.FindParameter("value").Value.(int64)

			return NonReturningValue(IntValue(this * value))
		}),
	})

	define(ctx, "divide", &Function{
		Signature: Signature{
			Parameters: []Parameter{
				{
					Name: "this",
					Type: IntType,
				},
				{
					Name: "value",
					Type: IntType,
				},
			},
			ReturnType: IntType,
		},
		Body: NewAbstractCommand(func(ctx *Context) ReturnedValue {
			this := ctx.FindParameter("this").Value.(int64)
			value := ctx.FindParameter("value").Value.(int64)

			return NonReturningValue(IntValue(this / value))
		}),
	})

	define(ctx, "mod", &Function{
		Signature: Signature{
			Parameters: []Parameter{
				{
					Name: "this",
					Type: IntType,
				},
				{
					Name: "value",
					Type: IntType,
				},
			},
			ReturnType: IntType,
		},
		Body: NewAbstractCommand(func(ctx *Context) ReturnedValue {
			this := ctx.FindParameter("this").Value.(int64)
			value := ctx.FindParameter("value").Value.(int64)

			return NonReturningValue(IntValue(this % value))
		}),
	})
}

func define(ctx *Context, name string, function *Function) {
	function.name = &name
	funcType := NewFunctionType(function)
	ctx.DefineVariable(name, Variable{
		Name:    name,
		Mutable: false,
		Type:    funcType,
		Value: &Value{
			Type:  funcType,
			Value: function,
		},
	})
}
