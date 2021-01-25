package parser

import (
	"github.com/ElaraLang/elara/lexer"
	"github.com/ElaraLang/elara/util"
	"testing"
)

func TestBasicParsing(t *testing.T) {
	code := "let a = 5 * 3 + 3 * 8"
	codeChannel := util.ToChannel(code)
	outputChannel := make(chan lexer.Token)
	go lexer.Lex(codeChannel, outputChannel)
	parser := NewReplParser(outputChannel)
	parser.ParseStatement()
}
