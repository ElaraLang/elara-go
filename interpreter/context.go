package interpreter

import (
	"elara/parser"
	"fmt"
)

type Context struct {
	variables map[string][]*Variable
}

func NewContext() *Context {
	return &Context{variables: map[string][]*Variable{}}
}

func (c Context) DefineVariable(name string, value Variable) {
	vars := c.variables[name]
	vars = append(vars, &value)
	c.variables[name] = vars
}

func (c Context) string() string {
	s := ""
	for key, values := range c.variables {
		s += fmt.Sprintf("%s = [\n", key)
		for _, val := range values {
			s += fmt.Sprintf("%s \n", val.string())
		}
		s += "]"
	}

	return s
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
