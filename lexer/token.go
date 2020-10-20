package lexer

import "fmt"

type Token struct {
	TokenType TokenType
	Text      string
	Position  Position
}

type Position struct {
	file   *string
	line   int
	column int
}

func CreateToken(tokenType TokenType, text string, position Position) Token {
	return Token{
		TokenType: tokenType,
		Text:      text,
		Position:  position,
	}
}

func CreatePosition(file *string, line int, column int) Position {
	return Position{
		file:   file,
		line:   line,
		column: column,
	}
}

func (t *Token) String() string {
	return fmt.Sprintf("%s '%s' at %s", t.TokenType.String(), t.Text, t.Position.String())
}

func (p *Position) String() string {
	if p.file != nil {
		return fmt.Sprintf("%s, %d:%d", *p.file, p.line, p.column)
	}
	return fmt.Sprintf("Unknown file, %d:%d", p.line, p.column)
}
