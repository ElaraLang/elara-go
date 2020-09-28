package main

import (
	"elara/interpreter"
	"elara/lexer"
	"elara/parser"
	"fmt"
	"os"
	"strings"
)

func main() {
	elaraFile, err := os.Open("elara.el")
	if err != nil {
		panic(err)
	}
	var bytes []byte
	_, _ = elaraFile.Read(bytes)
	reader := strings.NewReader(string(bytes))
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
}

func CreateToken(tokenType lexer.TokenType, text string) lexer.Token {
	return lexer.Token{
		TokenType: tokenType,
		Text:      text,
	}
}
