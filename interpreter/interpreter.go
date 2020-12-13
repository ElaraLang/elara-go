package interpreter

import (
	"fmt"
	"github.com/ElaraLang/elara/parser"
)

type Interpreter struct {
	lines   []parser.Stmt
	context *Context
}

func NewInterpreter(code []parser.Stmt) *Interpreter {
	return &Interpreter{
		lines:   code,
		context: NewContext(),
	}
}
func NewEmptyInterpreter() *Interpreter {
	return &Interpreter{
		context: NewContext(),
	}
}

func (s *Interpreter) ResetLines(lines *[]parser.Stmt) {
	s.lines = *lines
}

func (s *Interpreter) Exec(scriptMode bool) []*Value {
	values := make([]*Value, len(s.lines))

	for i := 0; i < len(s.lines); i++ {
		line := s.lines[i]
		command := ToCommand(line)

		//Ignore any top level returns
		defer func() {
			r := recover()
			if r != nil {
				_, isValue := r.(*Value)
				if isValue {
					return
				}
				panic(r)
			}
		}()
		res := command.Exec(s.context)
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
