package lexer

func Lex(file string, code string) *[]Token {
	chars := []rune(code)
	scanner := NewTokenReader(chars)

	//Note: in our big benchmark, the token:chars ratio seems to be about 1:1.2 (5:6). Could be worth doing len(code) / 1.2 and rounding?
	estimateLength := len(code)
	if estimateLength > 10 {
		estimateLength /= 2
	}
	tokens := make([]Token, estimateLength)
	i := 0

	for {
		tok, runes, line, col := scanner.Read()
		if tok == EOF {
			tokens = tokens[:i]
			break
		}

		token := Token{
			TokenType: tok,
			Text:      runes,
			Position:  CreatePosition(&file, line, col),
		}

		if i <= len(tokens)-1 {
			tokens[i] = token
			i++
		} else {
			tokens = append(tokens, token)
			i = len(tokens)
		}
	}

	return &tokens
}
