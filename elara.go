package main

import (
	"elara/lexer"
	"strings"
)

func main() {
	text := "let a = hello"
	reader := strings.NewReader(text)
	scanner := lexer.NewScanner(reader)
	for {
		tok, str := scanner.Scan()
		if tok == lexer.EOF {
			break
		}

		println(lexer.TokenNames[tok] + ": '" + str + "'")
	}

}
