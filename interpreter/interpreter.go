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

func (s *Interpreter) Exec(scriptMode bool) {
	_ = NewStack()
	context := NewContext()

	for i := 0; i < len(s.lines); i++ {
		line := s.lines[i]
		command := ToCommand(line)

		res := command.Exec(context)
		if scriptMode {
			fmt.Println(res.String())
		}
	}
}
