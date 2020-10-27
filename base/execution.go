package base

import (
	"elara/interpreter"
	"elara/lexer"
	"elara/parser"
	"fmt"
	"os"
	"time"
)

func Execute(fileName *string, code string, repl bool) (results []*interpreter.Value, lexTime time.Duration, parseTime time.Duration, execTime time.Duration) {
	tokens := make(chan *lexer.Token)
	go lexer.Lex(fileName, code, tokens)

	psr := parser.NewParser(tokens)
	parseRes, errs := psr.Parse()

	if len(errs) != 0 {
		_, _ = os.Stderr.WriteString("Parse Errors: \n")
		_, _ = os.Stderr.WriteString(fmt.Sprintf("%q\n", errs))
		return []*interpreter.Value{}, lexTime, parseTime, time.Duration(-1)
	}

	evaluator := interpreter.NewInterpreter(parseRes)
	interpreter.Init()

	results = evaluator.Exec(repl)
	return results, lexTime, parseTime, execTime
}
