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
	Params []Type
	Output Type
}

func (t FunctionType) String() string {
	stringedParams := make([]string, len(t.Params))
	for i, param := range t.Params {
		stringedParams[i] = param.String()
	}
	return fmt.Sprintf("(%s) => %s", strings.Join(stringedParams, ", "), t.Output.String())
}
