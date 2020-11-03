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
	name := "Function"
	if f.name != nil {
		name = *f.name
	}
	return fmt.Sprintf("%s %s => %s", name, f.Signature.Parameters, f.Signature.ReturnType)
}

func (f *Function) Exec(ctx *Context, receiver *Value, parameters []*Value) (val *Value) {
	if len(parameters) != len(f.Signature.Parameters) {
		panic(fmt.Sprintf("Illegal number of arguments for function %s. Expected %d, received %d", *f.name, len(f.Signature.Parameters), len(parameters)))
	}

	ctx.EnterScope(f.String())

	for i, paramValue := range parameters {
		expectedParameter := f.Signature.Parameters[i]

		if !expectedParameter.Type.Accepts(*paramValue.Type) {
			panic(fmt.Sprintf("Expected %s for parameter %s and got %s", expectedParameter.Type.Name, expectedParameter.Name, *paramValue.String()))
		}

		ctx.DefineParameter(expectedParameter.Name, paramValue)
	}

	ctx.receiver = receiver

	defer func() {
		ctx.ExitScope()
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
