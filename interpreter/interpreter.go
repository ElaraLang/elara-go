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
		s := line.(parser.VarDefStmt)

		var Type parser.Type
		if s.Type == nil {
			Type = parser.ElementaryTypeContract{
				Identifier: "Any",
			}
		} else {
			Type = *s.Type
		}
		variable := Variable{
			Name:    s.Identifier,
			Mutable: s.Mutable,
			Type:    Type,
			Value:   Value{s.Value},
		}

		context.DefineVariable(s.Identifier, variable)
	}

	println(context.string())

}
