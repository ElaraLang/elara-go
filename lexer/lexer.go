package lexer

var emptyToken = &Token{}

func Lex(file *string, code string) *[]Token {
	chars := []rune(code)
	scanner := NewTokenReader(chars)

	//Note: in our big benchmark, the token:chars ratio seems to be about 1:1.2 (5:6). Could be worth doing len(code) / 1.2 and rounding?
	tokens := make([]Token, len(code)/2)
	i := 0

	for {
		tok, runes, line, col := scanner.Read()
		if tok == EOF {
			break
		}
		token := Token{
			TokenType: tok,
			Text:      runes,
			Position:  CreatePosition(file, line, col),
		}

		if i <= len(tokens)-1 && tokens[i].Equals(emptyToken) {
			tokens[i] = token
			i++
		} else {
			tokens = append(tokens, token)
		}
	}

	for i := range tokens {
		if tokens[i].Equals(emptyToken) {
			tokens = tokens[:i]
			break
		}
	}
	return &tokens
}
