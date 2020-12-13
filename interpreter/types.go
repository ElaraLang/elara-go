package interpreter

import (
	"fmt"
	"github.com/ElaraLang/elara/parser"
	"reflect"
	"strings"
)

type Type struct {
	Name      string
	variables *VariableMap
}

func (t *Type) Accepts(other Type) bool {
	if &other == t {
		return true
	}
	for i := range t.variables.m {
		fun1 := t.variables.m[i]
		fun2 := other.variables.m[i]
		if fun1 == fun2 {
			continue
		}
		if !fun1.Equals(*fun2) {
			return false
		}
	}
	return true
}

func EmptyType(name string) *Type {
	return &Type{
		Name:      name,
		variables: NewVariableMap(),
	}
}
func SimpleType(name string, functions VariableMap) *Type {
	return &Type{
		Name:      name,
		variables: &functions,
	}
}

func AliasType(other Type) *Type {
	return &Type{
		Name:      other.Name,
		variables: other.variables,
	}
}

func CompoundType(a Type, b Type) *Type {
	functionCompound := NewVariableMap()
	for name, function := range a.variables.m {
		functionCompound.Set(name, function)
	}
	for name, function := range b.variables.m {
		functionCompound.Set(name, function)
	}

	return &Type{
		Name:      fmt.Sprintf("%sAnd%s", a.Name, b.Name),
		variables: functionCompound,
	}
}

func UnionType(a Type, b Type) *Type {
	panic("TODO")
}

func FunctionType(function *Function) *Type {
	//Build the type name based on signature
	var parameters string
	if len(function.Signature.Parameters) == 0 {
		parameters = "()"
	} else {
		paramTypes := make([]string, len(function.Signature.Parameters))
		for i, param := range function.Signature.Parameters {
			paramTypes[i] = param.Type.Name
		}
		parameters = "(" + strings.Join(paramTypes, ", ") + ")"
	}

	functionName := parameters + " => " + function.Signature.ReturnType.Name

	t := &Type{
		Name: functionName,
	}

	t.variables = NewVariableMap()
	t.variables.Set(functionName, &Variable{
		Name:    "value",
		Mutable: false,
		Type:    *t,
		Value: &Value{
			Type:  t,
			Value: function,
		},
	})

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

type VariableMap struct {
	m    map[string]*Variable
	keys []string
}

func NewVariableMap() *VariableMap {
	return &VariableMap{
		m:    map[string]*Variable{},
		keys: []string{},
	}
}

func (n *VariableMap) Set(k string, v *Variable) {
	n.m[k] = v
	n.keys = append(n.keys, k)
}
