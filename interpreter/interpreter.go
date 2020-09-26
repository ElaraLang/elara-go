package interpreter

import (
	"elara/parser"
)

type Interpreter struct {
	lines []parser.Stmt
}

func NewInterpreter(code []parser.Stmt) *Interpreter {
	return &Interpreter{
		lines: code,
	}
}

func (s *Interpreter) Exec() {
	_ = NewStack()
	context := NewContext()

	for i := 0; i < len(s.lines); i++ {
		line := s.lines[i]
		command := ToCommand(line)

		command.exec(context)
	}

	println(context.string())

}
