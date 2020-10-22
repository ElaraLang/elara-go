package interpreter

import (
	"elara/parser"
	"fmt"
	"reflect"
	"strings"
)

type Type struct {
	Name      string
	functions map[string]Function
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
		functions: map[string]Function{},
	}
}
func SimpleType(name string, functions map[string]Function) *Type {
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
	functionCompound := map[string]Function{}
	for name, function := range a.functions {
		functionCompound[name] = function
	}
	for name, function := range b.functions {
		functionCompound[name] = function
	}

	return &Type{
		Name:      fmt.Sprintf("%sAnd%s", a.Name, b.Name),
		functions: functionCompound,
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
		newName := fmt.Sprintf("%sTo%sFunction", strings.Join(paramNames, ""), function.Signature.ReturnType.Name)
		return &Type{
			Name:      newName,
			functions: map[string]Function{newName: function},
		}
	}
	return &Type{
		Name:      *name,
		functions: map[string]Function{*name: function},
	}
}

func FromASTType(ast parser.Type) *Type {
	if ast == nil {
		return AnyType
	}
	switch t := ast.(type) {
	case parser.ElementaryTypeContract:
		identifier := t.Identifier
		builtIn := BuiltInTypeByName(identifier)
		if builtIn != nil {
			return builtIn
		}
		return SimpleType(identifier, map[string]Function{})
	}

	panic("Could not handle " + reflect.TypeOf(ast).Name())
}
