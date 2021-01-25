package lexer

import "fmt"

type Token struct {
	TokenType TokenType
	Data      []rune
	Line      int
	Col       int
}

func (t Token) String() string {
	data := ""
	if t.Data != nil {
		data = fmt.Sprintf("(%s)", string(t.Data))
	}
	return fmt.Sprintf("%s %s at %d:%d", t.TokenType.String(), data, t.Line, t.Col)
}
