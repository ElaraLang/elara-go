package main

import (
	"elara/interpreter"
	"elara/lexer"
	"elara/parser"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"
)

func main() {
	start := time.Now()
	goPath := os.Getenv("GOPATH")
	filePath := path.Join(goPath, "elara.el")
	input, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	reader := strings.NewReader(string(input))
	scanner := lexer.NewScanner(reader)

	result := make([]lexer.Token, 0)
	for {
		tok, str := scanner.Read()
		result = append(result, CreateToken(tok, str))
		if tok == lexer.EOF {
			break
		}
	}

	psr := parser.NewParser(&result)
	parseRes, errs := psr.Parse()

	if len(errs) != 0 {
		println("Errors")
		fmt.Printf("%q\n", errs)
	}

	interpreter := interpreter.NewInterpreter(parseRes)

	interpreter.Exec()
	duration := time.Since(start)

	fmt.Printf("Executed in %s", duration)
}

func CreateToken(tokenType lexer.TokenType, text string) lexer.Token {
	return lexer.Token{
		TokenType: tokenType,
		Text:      text,
	}
}
