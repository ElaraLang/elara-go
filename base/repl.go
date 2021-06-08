package base

import (
	"fmt"
	"github.com/ElaraLang/elara/interpreter"
	"github.com/ElaraLang/elara/lexer"
	"github.com/ElaraLang/elara/parser"
)

var replFile = "Repl"

type ReplSession struct {
	Parser    parser.Parser
	Evaluator interpreter.Interpreter
}

func NewReplSession() ReplSession {
	return ReplSession{
		Parser:    *parser.NewEmptyParser(),
		Evaluator: *interpreter.NewEmptyInterpreter(),
	}
}

func (repl *ReplSession) Execute(input string) interface{} {
	tokens := lexer.Lex(input)
	repl.Parser.Reset(tokens)
	result, err := repl.Parser.Parse()
	if len(err) > 0 {
		fmt.Println("Errors found: ", err)
		return nil
	}
	repl.Evaluator.ResetLines(&result)
	evalRes := repl.Evaluator.Exec(true)
	return evalRes
}
