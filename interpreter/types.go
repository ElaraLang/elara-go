package interpreter

import (
	"elara/parser"
	"fmt"
	"github.com/adam-hanna/arrayOperations"
	"reflect"
	"strings"
)

type Type struct {
	Name      string
	functions []Function
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

func EmptyType(name string) *Type {
	return &Type{
		Name:      name,
		functions: []Function{},
	}
}
func SimpleType(name string, functions []Function) *Type {
	return &Type{
		Name:      name,
		functions: functions,
	}
}

func AliasType(other Type) *Type {
	return &Type{
		Name:      other.Name,
		functions: other.functions,
	}
}

func CompoundType(a Type, b Type) *Type {
	functionCompound, ok := arrayOperations.Union(a.functions, b.functions)
	if !ok {
		panic("Cannot find the union of types " + a.Name + " and " + b.Name)
	}

	functions, ok := functionCompound.Interface().([]Function)
	if !ok {
		panic("Cannot convert to slice of functions")
	}

	return &Type{
		Name:      fmt.Sprintf("%sAnd%s", a.Name, b.Name),
		functions: functions,
	}
}

func UnionType(a Type, b Type) *Type {
	panic("TODO")
}

func FunctionType(name *string, function Function) *Type {
	if name == nil {
		paramNames := make([]string, len(function.Signature.Parameters))
		for i := range function.Signature.Parameters {
			paramNames[i] = function.Signature.Parameters[i].Name
		}
		return &Type{
			Name:      fmt.Sprintf("%sTo%sFunction", strings.Join(paramNames, ""), function.Signature.ReturnType.Name),
			functions: []Function{function},
		}
	}
	return &Type{
		Name:      *name,
		functions: []Function{function},
	}
}

func FromASTType(ast parser.Type) *Type {
	if ast == nil {
		return AnyType
	}
	switch t := ast.(type) {
	case parser.ElementaryTypeContract:
		return SimpleType(t.Identifier, []Function{})
	}

	panic("Could not handle " + reflect.TypeOf(ast).Name())
}
