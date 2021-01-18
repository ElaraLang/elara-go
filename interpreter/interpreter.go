package interpreter

import (
	"fmt"
	"github.com/ElaraLang/elara/parserlegacy"
	"reflect"
)

type Interpreter struct {
	lines   []parserlegacy.Stmt
	context *Context
}

func NewInterpreter(code []parserlegacy.Stmt) *Interpreter {
	return &Interpreter{
		lines:   code,
		context: NewContext(true),
	}
}
func NewEmptyInterpreter() *Interpreter {
	return NewInterpreter([]parserlegacy.Stmt{})
}

func (s *Interpreter) ResetLines(lines *[]parserlegacy.Stmt) {
	s.lines = *lines
}

func (s *Interpreter) Exec(scriptMode bool) []*Value {
	values := make([]*Value, len(s.lines))

	for i := 0; i < len(s.lines); i++ {
		line := s.lines[i]
		command := ToCommand(line)
		res := command.Exec(s.context).Unwrap()
		values[i] = res
		if scriptMode {
			formatted := s.context.Stringify(res) + " " + reflect.TypeOf(res).String()
			fmt.Println(formatted)
		}
	}
	return values
}
