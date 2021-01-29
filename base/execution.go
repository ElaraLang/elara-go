package base

import (
	"github.com/ElaraLang/elara/ast"
	"github.com/ElaraLang/elara/interpreter"
	"github.com/ElaraLang/elara/lexer"
	"github.com/ElaraLang/elara/parser"
)

func Execute(fileName string, code chan rune, scriptMode bool, io interpreter.IO) {
	lexerOutput := make(chan lexer.Token)
	parserOutput := make(chan ast.Statement)
	parseErrors := make(chan parser.ParseError)
	go lexer.Lex(code, lexerOutput)

	psr := parser.NewParser(lexerOutput, parserOutput, parseErrors)
	go psr.Parse(fileName)

	interpreterOutput := make(chan *interpreter.Value)
	evaluator := interpreter.NewInterpreter(parserOutput, parseErrors, interpreterOutput, io)

	evaluator.Exec(scriptMode)
	return
}
