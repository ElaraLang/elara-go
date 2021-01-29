package interpreter

import (
	"fmt"
	"github.com/ElaraLang/elara/ast"
	"github.com/ElaraLang/elara/parser"
	"reflect"
)

type Interpreter struct {
	lines       chan ast.Statement
	parseErrors chan parser.ParseError
	output      chan *Value
	context     *Context
}

func NewInterpreter(code chan ast.Statement, parseErrors chan parser.ParseError, output chan *Value, io IO) *Interpreter {
	globalContext.io = io
	return &Interpreter{
		lines:       code,
		parseErrors: parseErrors,
		output:      output,
		context:     NewContext(true),
	}
}
func NewEmptyInterpreter() *Interpreter {
	globalContext.io = NewSysIO()
	return NewInterpreter(nil, nil, nil, globalContext.io)
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
				s.context.io.Println(formatted)
				//fmt.Println(formatted)
			}
		case err := <-s.parseErrors:
			s.context.io.Error(fmt.Sprintf("%s\n", err))
			// _, _ = os.Stderr.WriteString(fmt.Sprintf("%s\n", err))
			return
		default:
		}

	}
}
