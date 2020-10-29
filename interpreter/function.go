package interpreter

import (
	"fmt"
)

type Function struct {
	Signature Signature
	Body      Command
	name      *string
}

func (f *Function) String() string {
	return fmt.Sprintf("Function %s => %s", f.Signature.Parameters, f.Signature.ReturnType)
}

func (f *Function) Exec(ctx *Context, receiver *Value, parameters []Command) (val *Value) {
	for i, parameter := range parameters {
		paramValue := parameter.Exec(ctx)
		expectedParameter := f.Signature.Parameters[i]

		if !expectedParameter.Type.Accepts(*paramValue.Type) {
			panic(fmt.Sprintf("Expected %s for parameter %s and got %s", expectedParameter.Type.Name, expectedParameter.Name, paramValue.Type.Name))
		}

		ctx.DefineParameter(expectedParameter.Name, paramValue)
	}

	ctx.receiver = receiver

	defer func() {
		s := recover()
		if s != nil {
			_, is := s.(*Value)
			if is {
				val = s.(*Value)
			} else {
				panic(s)
			}
		}
	}()

	value := f.Body.Exec(ctx)
	if !f.Signature.ReturnType.Accepts(*value.Type) {
		name := "<anonymous>"
		if f.name != nil {
			name = *f.name
		}
		panic(fmt.Sprintf("Function '%s' did not return value of type %s, instead was %s", name, f.Signature.ReturnType.Name, value.Type.Name))
	}
	return value
}

type Signature struct {
	Parameters []Parameter
	ReturnType Type
}

type Parameter struct {
	Name string
	Type Type
}
