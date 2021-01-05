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

	scope := context.EnterScope(f.String())

	for i, paramValue := range parameters {
		expectedParameter := f.Signature.Parameters[i]

		if !expectedParameter.Type.Accepts(paramValue.Type) {
			panic(fmt.Sprintf("Expected %s for parameter %s and got %s", expectedParameter.Type.Name(), expectedParameter.Name, *paramValue.String()))
		}

		scope.DefineParameter(expectedParameter.Name, paramValue)
	}

	defer func() { //Catch returned values
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

	value := f.Body.Exec(scope)
	if value == nil {
		value = UnitValue()
	}
	if !f.Signature.ReturnType.Accepts(value.Type) {
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

func (s *Signature) Accepts(other *Signature, compareReturnTypes bool) bool {
	if len(s.Parameters) != len(other.Parameters) {
		return false
	}
	for i, parameter := range s.Parameters {
		otherParam := other.Parameters[i]
		if !parameter.Type.Accepts(otherParam.Type) {
			return false
		}
	}
	if compareReturnTypes {
		return s.ReturnType.Accepts(other.ReturnType)
	}
	return true
}

type Parameter struct {
	Name string
	Type Type
}
