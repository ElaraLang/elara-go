package lexer

func Lex(file *string, code string) (tokens chan *Token) {
	chars := []rune(code)
	scanner := NewTokenReader(chars)

	tokens = make(chan *Token)

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
	}

	return tokens
}
