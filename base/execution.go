package base

import (
	"fmt"
	"github.com/ElaraLang/elara/interpreter"
	"github.com/ElaraLang/elara/lexer"
	"github.com/ElaraLang/elara/parser"
	"os"
	"time"
)

func Execute(fileName *string, code string, scriptMode bool) (results []*interpreter.Value, lexTime, parseTime, execTime time.Duration) {

	start := time.Now()
	result := lexer.Lex(fileName, code)
	lexTime = time.Since(start)

	start = time.Now()
	psr := parser.NewParser(result)
	parseRes, errs := psr.Parse()
	parseTime = time.Since(start)

	if len(errs) != 0 {
		_, _ = os.Stderr.WriteString("Parse Errors: \n")
		_, _ = os.Stderr.WriteString(fmt.Sprintf("%q\n", errs))
		return []*interpreter.Value{}, lexTime, parseTime, time.Duration(-1)
	}

	start = time.Now()
	evaluator := interpreter.NewInterpreter(parseRes)
	interpreter.Init()

	results = evaluator.Exec(scriptMode)
	execTime = time.Since(start)
	return results, lexTime, parseTime, execTime
}
