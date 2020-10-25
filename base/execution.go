package base

import (
	"elara/interpreter"
	"elara/lexer"
	"elara/parser"
	"fmt"
	"os"
)

func Execute(fileName *string, code string, repl bool) []*interpreter.Value {

	result := lexer.Lex(fileName, code)

	psr := parser.NewParser(result)
	parseRes, errs := psr.Parse()

	if len(errs) != 0 {
		_, _ = os.Stderr.WriteString("Parse Errors: \n")
		_, _ = os.Stderr.WriteString(fmt.Sprintf("%q\n", errs))
		return []*interpreter.Value{}
	}

	evaluator := interpreter.NewInterpreter(parseRes)
	interpreter.Init()

	return evaluator.Exec(repl)
}
