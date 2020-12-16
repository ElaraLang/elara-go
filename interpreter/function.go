package interpreter

import (
	"fmt"
	"math"
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

func (f *Function) Exec(ctx *Context, parameters []*Value) (val *Value) {
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

func (sig *Signature) Accepts(other Signature) bool {
	if len(sig.Parameters) != len(other.Parameters) {
		return false
	}
	for i, param := range sig.Parameters {
		if !param.Type.Accepts(other.Parameters[i].Type) {
			return false
		}
	}
	return sig.ReturnType.Accepts(other.ReturnType)
}

/**
Get the distance of a given signature to another one
This is similar to the Levenshtein Distance algorithm in that 0 means the signatures are the same,
and any additional value is the number of adjustments that must be made to *types only* to reach the other signature
*/
func (sig *Signature) Distance(other Signature) int {
	if sig == &other {
		return 0
	}
	changes := 0
	changes += int(math.Abs(float64(len(sig.Parameters) - len(other.Parameters)))) //distance between the param types themself

	if len(sig.Parameters) == len(other.Parameters) {
		for i, parameter := range sig.Parameters {
			if !parameter.Type.Accepts(other.Parameters[i].Type) {
				changes++
			}
		}
	}

	if !sig.ReturnType.Accepts(other.ReturnType) {
		changes++
	}
	return changes
}

type Parameter struct {
	Name string
	Type Type
}
