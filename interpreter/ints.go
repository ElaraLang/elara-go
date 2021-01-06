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
					Name:     "value",
					Type:     IntType,
					Position: 1,
				},
			},
			ReturnType: IntType,
		},
		Body: NewAbstractCommand(func(ctx *Context) ReturnedValue {
			this := ctx.FindParameter(0).Value.(int64)
			value := ctx.FindParameter(1).Value.(int64)

			returningValue := NonReturningValue(IntValue(this + value))
			return returningValue
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
					Name:     "value",
					Type:     IntType,
					Position: 1,
				},
			},
			ReturnType: IntType,
		},
		Body: NewAbstractCommand(func(ctx *Context) ReturnedValue {
			this := ctx.FindParameter(0).Value.(int64)
			value := ctx.FindParameter(1).Value.(int64)

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
					Name:     "value",
					Type:     IntType,
					Position: 1,
				},
			},
			ReturnType: IntType,
		},
		Body: NewAbstractCommand(func(ctx *Context) ReturnedValue {
			this := ctx.FindParameter(0).Value.(int64)
			value := ctx.FindParameter(1).Value.(int64)

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
					Name:     "value",
					Type:     IntType,
					Position: 1,
				},
			},
			ReturnType: IntType,
		},
		Body: NewAbstractCommand(func(ctx *Context) ReturnedValue {
			this := ctx.FindParameter(0).Value.(int64)
			value := ctx.FindParameter(1).Value.(int64)

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
					Name:     "value",
					Type:     IntType,
					Position: 1,
				},
			},
			ReturnType: IntType,
		},
		Body: NewAbstractCommand(func(ctx *Context) ReturnedValue {
			this := ctx.FindParameter(0).Value.(int64)
			value := ctx.FindParameter(1).Value.(int64)

			return NonReturningValue(IntValue(this % value))
		}),
	})
}

func define(ctx *Context, name string, function *Function) {
	function.name = &name
	funcType := NewFunctionType(function)
	ctx.DefineVariable(name, &Variable{
		Name:    name,
		Mutable: false,
		Type:    funcType,
		Value: &Value{
			Type:  funcType,
			Value: function,
		},
	})
}
