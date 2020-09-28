package main

import (
	"elara/interpreter"
	"elara/lexer"
	"elara/parser"
	"fmt"
	"strings"
)

func main() {
	text :=
		"let a = 5\n" +
			"print(a)\n" +
			"print(\"Hello\")"
	reader := strings.NewReader(text)
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
	parseRes, err := psr.Parse()

	if len(err) != 0 {
		println("Errors")
		fmt.Printf("%q\n", err)
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
