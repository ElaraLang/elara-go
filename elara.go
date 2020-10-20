package main

import (
	"elara/interpreter"
	"elara/lexer"
	"elara/parser"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"time"
)

func main() {
	args := os.Args
	//This isn't really a repl, but it will be. for now, it's close enough in that it will print the output of every expression
	var repl = false
	for _, arg := range args {
		if arg == "--repl" {
			repl = true
			break
		}
	}

	fileName, input := loadElaraFile()

	start := time.Now()

	result := lexer.Lex(&fileName, string(input))

	psr := parser.NewParser(&result)
	parseRes, errs := psr.Parse()

	if len(errs) != 0 {
		_, _ = os.Stderr.WriteString("Parse Errors")
		fmt.Printf("%q\n", errs)
		return
	}

	evaluator := interpreter.NewInterpreter(parseRes)
	interpreter.Init()

	evaluator.Exec(repl)
	duration := time.Since(start)

	fmt.Printf("Executed in %s.\n", duration)
}

func loadElaraFile() (string, []byte) {
	goPath := os.Getenv("GOPATH")
	fileName := "elara.el"
	filePath := path.Join(goPath, fileName)

	input, err := ioutil.ReadFile(filePath)

	if err != nil {
		panic(err)
	}
	return fileName, input
}
