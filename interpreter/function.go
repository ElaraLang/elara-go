package interpreter

import (
	"fmt"
)

type Function struct {
	Signature Signature
	Body      Command
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
			}
		}
	}()
	return f.Body.Exec(ctx)
}

type Signature struct {
	Parameters []Parameter
	ReturnType Type
}

type Parameter struct {
	Name string
	Type Type
}
