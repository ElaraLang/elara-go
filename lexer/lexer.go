package lexer

import "fmt"

func Lex(file *string, code string, tokens chan *Token) {
	chars := []rune(code)
	scanner := NewTokenReader(chars)

	for {
		tok, runes, line, col := scanner.Read()
		if tok == EOF {
			tokens <- nil
			break
		}

		token := Token{
			TokenType: tok,
			Text:      runes,
			Position:  CreatePosition(file, line, col),
		}

		tokens <- &token
		fmt.Printf("Sent token to channel %s \n", tok.String())
	}
}
