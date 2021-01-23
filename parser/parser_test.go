package parser

import (
	"github.com/ElaraLang/elara/lexer"
	"testing"
)

func TestBasicParsing(t *testing.T) {
	code := "let a = 30"
	tokens := lexer.Lex(code)
	parser := NewReplParser(make(chan lexer.Token))
	go parser.parseStatement()
	for _, v := range tokens {
		parser.Tape.Channel <- v
	}
	//fmt.Println(res.ToString())
}
