package types

import (
	"fmt"
	"strings"
)

type Type interface {
	String() string
}

type SimpleType struct {
	name string
}

func (receiver *SimpleType) String() string {
	return receiver.name
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
