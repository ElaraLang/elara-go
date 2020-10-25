package lexer

import "fmt"

type Token struct {
	TokenType TokenType
	Text      []rune
	Position  Position
}

func (t *Token) Equals(other *Token) bool {
	if t.TokenType != other.TokenType {
		return false
	}
	if t.Position != other.Position {
		return false
	}
	if runeSliceEq(t.Text, other.Text) {
		return false
	}
	return true
}

type Position struct {
	file   *string
	line   int
	column int
}

func CreateToken(tokenType TokenType, text string, position Position) Token {
	return Token{
		TokenType: tokenType,
		Text:      []rune(text),
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
	return fmt.Sprintf("%s '%s' at %s", t.TokenType.String(), string(t.Text), t.Position.String())
}

func (p *Position) String() string {
	if p.file != nil {
		return fmt.Sprintf("%s, %d:%d", *p.file, p.line, p.column)
	}
	return fmt.Sprintf("Unknown file, %d:%d", p.line, p.column)
}
