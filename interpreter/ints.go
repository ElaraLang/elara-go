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
}
