package main

import (
	"elara/lexer"
	"elara/parser"
	"fmt"
	"strings"
)

func main() {
	text := "let a = 5\n" + "let b = a * 5"
	reader := strings.NewReader(text)
	scanner := lexer.NewScanner(reader)

	result := make([]lexer.Token, 0)
	for {
		tok, str := scanner.Read()
		result = append(result, CreateToken(tok, str))
		if tok == lexer.EOF {
			break
		}
		println(lexer.TokenNames[tok] + ": '" + str + "'")
	}

	println(fmt.Sprintf("%q\n", result))
	psr := parser.NewParser(&result)
	parseRes, err := psr.Parse()
	println("ParseResult")
	println(fmt.Sprintf("%q\n", parseRes))
	println("Errors")
	println(fmt.Sprintf("%q\n", err))
}
func CreateToken(tokenType lexer.TokenType, text string) lexer.Token {
	return lexer.Token{
		TokenType: tokenType,
		Text:      text,
	}
}
