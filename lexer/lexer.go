package lexer

import (
	"strings"
)

func lex(code string) []Token {
	reader := strings.NewReader(code)
	scanner := NewScanner(reader)

	tokens := make([]Token, 0)
	for {
		tok, str := scanner.Read()
		if tok == EOF {
			break
		}

		tokens = append(tokens, Token{
			TokenType: tok,
			Text:      str,
		})
	}
	return tokens
}
