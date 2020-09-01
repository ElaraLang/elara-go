package main

import (
	"./lexer"
	"strings"
)

func main() {
	text := "let a = hello"
	reader := strings.NewReader(text)
	scanner := lexer.NewScanner(reader)
	for {
		tok, str := scanner.Read()
		if tok == lexer.EOF {
			break
		}

		println(lexer.TokenNames[tok] + ": '" + str + "'")
	}
}
