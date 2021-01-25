package base

import (
	"github.com/ElaraLang/elara/ast"
	"github.com/ElaraLang/elara/lexer"
	"github.com/ElaraLang/elara/parser"
	"time"
)

func Execute(fileName *string, code chan rune, scriptMode bool) (results ast.Statement, lexTime, parseTime, execTime time.Duration) {
	start := time.Now()
	output := make(chan lexer.Token)
	go lexer.Lex(code, output)

	start = time.Now()
	psr := parser.NewParser([]lexer.Token{}, output)
	parseRes := psr.ParseStatement()
	parseTime = time.Since(start)

	//if len(errs) != 0 {
	//	file := "Unknown File"
	//	if fileName != nil {
	//		file = *fileName
	//	}
	//	_, _ = os.Stderr.WriteString(fmt.Sprintf("Syntax Errors found in %s: \n", file))
	//	for _, err := range errs {
	//		_, _ = os.Stderr.WriteString(fmt.Sprintf("%s\n", err))
	//	}
	//	return []*interpreter.Value{}, lexTime, parseTime, time.Duration(-1)
	//}

	start = time.Now()
	//evaluator := interpreter.NewInterpreter(parseRes)

	//results = evaluator.Exec(scriptMode)
	execTime = time.Since(start)
	return parseRes, lexTime, parseTime, execTime
}
