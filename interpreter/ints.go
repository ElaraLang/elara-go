package interpreter

var IntType = NewEmptyType("Int")

func InitInts(ctx *Context) {
	intAddName := "plus"
	intAdd := &Function{
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
		name: &intAddName,
		Body: NewAbstractCommand(func(ctx *Context) *Value {
			this := ctx.FindParameter("this").Value.(int64)
			value := ctx.FindParameter("value").Value.(int64)

			return IntValue(this + value)
		}),
	}
	intAddType := NewFunctionType(intAdd)

	ctx.DefineVariable(intAddName, Variable{
		Name:    intAddName,
		Mutable: false,
		Type:    intAddType,
		Value: &Value{
			Type:  intAddType,
			Value: intAdd,
		},
	})

	intMinusName := "minus"
	intMinus := &Function{
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
		name: &intAddName,
		Body: NewAbstractCommand(func(ctx *Context) *Value {
			this := ctx.FindParameter("this").Value.(int64)
			value := ctx.FindParameter("value").Value.(int64)

			return IntValue(this - value)
		}),
	}
	intMinusType := NewFunctionType(intAdd)

	ctx.DefineVariable(intMinusName, Variable{
		Name:    intMinusName,
		Mutable: false,
		Type:    intMinusType,
		Value: &Value{
			Type:  intMinusType,
			Value: intMinus,
		},
	})

	intTimesName := "times"
	intTimes := &Function{
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
		name: &intTimesName,
		Body: NewAbstractCommand(func(ctx *Context) *Value {
			this := ctx.FindParameter("this").Value.(int64)
			value := ctx.FindParameter("value").Value.(int64)

			return IntValue(this * value)
		}),
	}
	intTimesType := NewFunctionType(intTimes)

	ctx.DefineVariable(intTimesName, Variable{
		Name:    intTimesName,
		Mutable: false,
		Type:    intTimesType,
		Value: &Value{
			Type:  intTimesType,
			Value: intTimes,
		},
	})
}
