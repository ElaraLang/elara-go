package interpreter

import (
	"fmt"
	"github.com/ElaraLang/elara/ast"
	"github.com/ElaraLang/elara/parser"
	"os"
	"reflect"
)

type Interpreter struct {
	lines       chan ast.Statement
	parseErrors chan parser.ParseError
	output      chan *Value
	context     *Context
}

func NewInterpreter(code chan ast.Statement, parseErrors chan parser.ParseError, output chan *Value) *Interpreter {
	return &Interpreter{
		lines:       code,
		parseErrors: parseErrors,
		output:      output,
		context:     NewContext(true),
	}
}
func NewEmptyInterpreter() *Interpreter {
	return NewInterpreter(nil, nil, nil)
}

func (s *Interpreter) ResetLines(lines chan ast.Statement) {
	s.lines = lines
}

func (s *Interpreter) Exec(scriptMode bool) {
	for {
		select {
		case line, ok := <-s.lines:
			if !ok {
				return
			}
			command := ToCommand(line)
			res := command.Exec(s.context).Unwrap()
			//s.output <- res
			if scriptMode {
				formatted := s.context.Stringify(res) + " " + reflect.TypeOf(res).String()
				fmt.Println(formatted)
			}
		case err := <-s.parseErrors:
			_, _ = os.Stderr.WriteString(fmt.Sprintf("%s\n", err))
			return
		default:
		}

	}
}
