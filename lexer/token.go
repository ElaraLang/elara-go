package lexer

type Token struct {
	TokenType TokenType
	data      []rune
	line      int
	col       int
}
