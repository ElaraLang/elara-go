package parser

import (
	"github.com/ElaraLang/elara/lexer"
	"testing"
)

func TestBasicParsing(t *testing.T) {
	code := "let a = 5 + 3\n"
	tokens := lexer.Lex(code)
	parser := NewReplParser(make(chan lexer.Token))
	go parser.parseLetStatement()
	for _, v := range tokens {
		parser.Tape.Channel <- v
	}
}
