package main

import (
	interpreter2 "elara/interpreter"
	"elara/lexer"
	"elara/parser"
	"fmt"
	"strings"
)

func main() {
	text := "let a = 5"
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

	println("Errors")
	println(fmt.Sprintf("%q\n", err))

	interpreter := interpreter2.NewInterpreter(parseRes)

	interpreter.Exec()
}

func CreateToken(tokenType lexer.TokenType, text string) lexer.Token {
	return lexer.Token{
		TokenType: tokenType,
		Text:      text,
	}
}
