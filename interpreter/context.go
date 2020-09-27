package interpreter

import (
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

func (c Context) FindVariable(name string) *Variable {
	vars := c.variables[name]

	if vars == nil {
		return nil
	}
	return vars[0]
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
