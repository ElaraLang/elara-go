package lexer

func Lex(file *string, code string) *[]Token {
	chars := []rune(code)
	scanner := NewTokenReader(chars)

	tokens := make([]Token, len(code)/3)
	i := 0
	emptyToken := Token{}
	for {
		tok, str, line, col := scanner.Read()
		if tok == EOF {
			break
		}
		token := Token{
			TokenType: tok,
			Text:      str,
			Position:  CreatePosition(file, line, col),
		}
		if i <= len(tokens)-1 && tokens[i] == emptyToken {
			tokens[i] = token
			i++
		} else {
			tokens = append(tokens, token)
		}
	}
	for i := range tokens {
		if tokens[i] == emptyToken {
			tokens = tokens[:i]
			break
		}
	}
	return &tokens
}
