package interpreter

import (
	"elara/parser"
	"fmt"
)

type Interpreter struct {
	lines []parser.Stmt
}

func NewInterpreter(code []parser.Stmt) *Interpreter {
	return &Interpreter{
		lines: code,
	}
}

func (s *Interpreter) Exec(scriptMode bool) []*Value {
	_ = NewStack()
	context := NewContext()

	values := make([]*Value, len(s.lines))

	for i := 0; i < len(s.lines); i++ {
		line := s.lines[i]
		command := ToCommand(line)

		res := command.Exec(context)
		values[i] = res
		if scriptMode {
			formatted := res.String()
			if formatted == nil {
				fmt.Println("<no value>")
			} else {
				fmt.Println(*formatted)
			}
		}
	}
	return values
}
