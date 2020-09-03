package lexer

type Token struct {
	TokenType TokenType
	Text      string
}

func CreateToken(tokenType TokenType, text string) Token {
	return Token{
		TokenType: tokenType,
		Text:      text,
	}
}
