package lexer

type Token struct {
	tokenType TokenType
	text      string
}

func CreateToken(tokenType TokenType, text string) Token {
	return Token{
		tokenType: tokenType,
		text:      text,
	}
}
