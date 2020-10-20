package typesystem

import (
	"elara/interpreter"
	"fmt"
	"github.com/adam-hanna/arrayOperations"
	"reflect"
)

type Type struct {
	name      string
	functions []interpreter.Function
}

func (t *Type) Accepts(other Type) bool {
	for i := range t.functions {
		fun1 := t.functions[i]
		for i2 := range other.functions {
			if reflect.DeepEqual(fun1, other.functions[i2]) {
				break
			}
		}
		return false
	}
	return true
}

func SimpleType(name string, functions []interpreter.Function) *Type {
	return &Type{
		name:      name,
		functions: functions,
	}
}

func AliasType(other Type) *Type {
	return &Type{
		name:      other.name,
		functions: other.functions,
	}
}

func CompoundType(a Type, b Type) *Type {
	functionCompound, ok := arrayOperations.Union(a.functions, b.functions)
	if !ok {
		panic("Cannot find the union of types " + a.name + " and " + b.name)
	}

	functions, ok := functionCompound.Interface().([]interpreter.Function)
	if !ok {
		panic("Cannot convert to slice of functions")
	}

	return &Type{
		name:      fmt.Sprintf("%sAnd%s", a.name, b.name),
		functions: functions,
	}
}

func UnionType(a Type, b Type) *Type {
	panic("TODO")
}
