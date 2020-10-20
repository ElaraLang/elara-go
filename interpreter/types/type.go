package types

import (
	"fmt"
	"strings"
)

type Type interface {
	String() string
	Accepts(other Type) bool
}

type SimpleType struct {
	name string
}

func (t *SimpleType) String() string {
	return t.name
}

func (t *SimpleType) Accepts(other Type) bool {
	simple, ok := other.(*SimpleType)
	if !ok {
		return false
	}
	return simple.name == t.name
}

type FunctionType struct {
	Params []Parameter
	Output Type
}

type Parameter struct {
	Type Type
	Name string
}

func (t *Parameter) String() string {
	return fmt.Sprintf("Parameter %s of type %s", t.Name, t.Type.String())
}

func (t FunctionType) String() string {
	stringedParams := make([]string, len(t.Params))
	for i, param := range t.Params {
		stringedParams[i] = param.String()
	}
	return fmt.Sprintf("(%s) => %s", strings.Join(stringedParams, ", "), t.Output.String())
}

func (t FunctionType) Accepts(other Type) bool {
	function, ok := other.(*FunctionType)
	if !ok {
		return false
	}
	if len(function.Params) != len(t.Params) {
		return false
	}
	for i, param := range t.Params {
		if !param.Type.Accepts(function.Params[i].Type) {
			return false
		}
	}
	return t.Output.Accepts(function.Output)
}

type SimpleAnonymousType struct {
	name    string
	accepts func(other Type) bool
}

func (t SimpleAnonymousType) String() string {
	return t.name
}

func (t SimpleAnonymousType) Accepts(other Type) bool {
	return t.accepts(other)
}
