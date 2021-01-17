package interpreter

import (
	"fmt"
	"github.com/ElaraLang/elara/util"
	"strings"
)

type Function struct {
	Signature Signature
	Body      Command
	name      *string
	context   *Context
}

func (f *Function) String() string {
	name := "Function"
	if f.name != nil {
		name = *f.name
	}
	return name + f.Signature.String()
}

func (f *Function) Exec(ctx *Context, parameters []*Value) (val *Value) {
	context := ctx
	if f.context != nil {
		//The cached context has highest priority for things like variables, but we set the parent to ensure that we can correctly inherit things like imports
		context = f.context.Clone()
		context.parent = ctx
	}
	if len(parameters) != len(f.Signature.Parameters) {
		panic(fmt.Sprintf("Illegal number of arguments for function %s. Expected %d, received %d", util.NillableStringify(f.name, "<anonymous>"), len(f.Signature.Parameters), len(parameters)))
	}

	var name string
	if f.name == nil {
		name = f.String()
	} else {
		name = *f.name
	}
	scope := context.EnterScope(name, f, uint(len(f.Signature.Parameters)))

	for i, paramValue := range parameters {
		expectedParameter := f.Signature.Parameters[i]

		if !expectedParameter.Type.Accepts(paramValue.Type, ctx) {
			panic(fmt.Sprintf("Expected %s for parameter %s and got %s (%s)", expectedParameter.Type.Name(), expectedParameter.Name, paramValue.String(), paramValue.Type.Name()))
		}
		//
		//if paramValue.Value == nil {
		//	panic("Value must not be nil")
		//}
		scope.DefineParameter(expectedParameter.Position, paramValue.Copy()) //Passing by value
	}

	value := f.Body.Exec(scope).Value //Can't unwrap because it might have returned from the function
	scope.Cleanup()                   //Exit out of the scope
	if value == nil {
		value = UnitValue()
	}
	if !f.Signature.ReturnType.Accepts(value.Type, ctx) {
		name := "<anonymous>"
		if f.name != nil {
			name = *f.name
		}
		panic(fmt.Sprintf("Function '%s' did not return value of type %s, instead was %s", name, f.Signature.ReturnType.Name(), value.Type.Name()))
	}
	return value
}

type Signature struct {
	Parameters []Parameter
	ReturnType Type
}

func (s *Signature) String() string {
	paramNames := make([]string, len(s.Parameters))
	for i := range s.Parameters {
		paramNames[i] = s.Parameters[i].Type.Name()
	}
	return fmt.Sprintf("(%s) => %s", strings.Join(paramNames, ", "), s.ReturnType.Name())
}

func (s *Signature) Accepts(other *Signature, ctx *Context, compareReturnTypes bool) bool {
	if len(s.Parameters) != len(other.Parameters) {
		return false
	}
	for i, parameter := range s.Parameters {
		otherParam := other.Parameters[i]
		if !parameter.Type.Accepts(otherParam.Type, ctx) {
			return false
		}
	}
	if compareReturnTypes {
		return s.ReturnType.Accepts(other.ReturnType, ctx)
	}
	return true
}

type Parameter struct {
	Name     string
	Position uint
	Type     Type
}
