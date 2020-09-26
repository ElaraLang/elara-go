package interpreter

import (
	"elara/parser"
	"fmt"
)

type Value struct {
	Type  parser.Type
	value interface{}
}

type Variable struct {
	Name    string
	Mutable bool
	Type    parser.Type
	Value   Value
}

func (v Variable) string() string {
	return fmt.Sprintf("Variable { name: %s, mutable: %T, type: %s, value: %s", v.Name, v.Mutable, v.Type, v.Value)
}
