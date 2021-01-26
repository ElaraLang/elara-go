package parser

import (
	"github.com/ElaraLang/elara/ast"
	"github.com/ElaraLang/elara/lexer"
	"testing"
)

func TestBasicParsing(t *testing.T) {
	code := "let a = 5 * 3 + 3 * 8"
	tokens := lexer.Lex(code)
	parser := NewParser(make(chan lexer.Token), make(chan ast.Statement))
	go parser.parseStatement()
	for _, v := range tokens {
		parser.Tape.Channel <- v
	}
}
