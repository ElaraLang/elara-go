package interpreter

import (
	"elara/parser"
	"fmt"
	"reflect"
	"strings"
)

type Type struct {
	Name      string
	variables map[string]Variable
}

func (t *Type) Accepts(other Type) bool {
	if &other == t {
		return true
	}
	for i := range t.variables {
		fun1 := t.variables[i]
		fun2 := other.variables[i]
		if !reflect.DeepEqual(fun1, fun2) {
			return false
		}
	}
	return true
}

func EmptyType(name string) *Type {
	return &Type{
		Name:      name,
		variables: map[string]Variable{},
	}
}
func SimpleType(name string, functions map[string]Variable) *Type {
	return &Type{
		Name:      name,
		variables: functions,
	}
}

func AliasType(other Type) *Type {
	return &Type{
		Name:      other.Name,
		variables: other.variables,
	}
}

func CompoundType(a Type, b Type) *Type {
	functionCompound := map[string]Variable{}
	for name, function := range a.variables {
		functionCompound[name] = function
	}
	for name, function := range b.variables {
		functionCompound[name] = function
	}

	return &Type{
		Name:      fmt.Sprintf("%sAnd%s", a.Name, b.Name),
		variables: functionCompound,
	}
}

func UnionType(a Type, b Type) *Type {
	panic("TODO")
}

func FunctionType(function Function) *Type {
	if function.name == nil {
		paramNames := make([]string, len(function.Signature.Parameters))
		for i := range function.Signature.Parameters {
			paramNames[i] = function.Signature.Parameters[i].Name
		}
		newName := fmt.Sprintf("%sTo%sFunction", strings.Join(paramNames, ""), function.Signature.ReturnType.Name)

		//Dirty hack but I cba to define a new function, especially with no overloading...
		function.name = &newName
		t := *FunctionType(function)
		function.name = nil

		return &Type{
			Name: newName,
			variables: map[string]Variable{newName: {
				Name:    newName,
				Mutable: false,
				Type:    t,
				Value: &Value{
					Type:  &t,
					Value: function,
				},
			}},
		}
	}
	t := &Type{
		Name: *function.name,
	}
	t.variables = map[string]Variable{*function.name: {
		Name:    *function.name,
		Mutable: false,
		Type:    *t,
		Value: &Value{
			Type:  t,
			Value: function,
		},
	}}
	return t
}

func FromASTType(ast parser.Type, ctx *Context) *Type {
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
		defined, isDefined := ctx.types[identifier]
		if isDefined {
			return &defined
		}
		panic("No such type " + identifier)
	}

	panic("Could not handle " + reflect.TypeOf(ast).Name())
}
