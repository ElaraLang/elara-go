package parser

import (
	"fmt"
	"github.com/ElaraLang/elara/lexer"
	"testing"
)

func TestBasicParsing(t *testing.T) {
	code := "let a = 30"
	tokens := lexer.Lex(code)
	parser := NewParser(tokens, make(chan lexer.Token))
	res := parser.parseStatement()
	fmt.Println(res.ToString())
}
