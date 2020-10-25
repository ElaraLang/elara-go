package lexer

import (
	"strings"
)

func Lex(file *string, code string) []Token {
	reader := strings.NewReader(code)
	scanner := NewScanner(reader)

	tokens := make([]Token, len(code)/4)
	for {
		tok, str, line, col := scanner.Read()
		if tok == EOF {
			break
		}

		tokens = append(tokens, Token{
			TokenType: tok,
			Text:      str,
			Position:  CreatePosition(file, line, col),
		})
	}
	return tokens
}
