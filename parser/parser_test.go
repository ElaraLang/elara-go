package parser

import (
	"fmt"
	"github.com/ElaraLang/elara/ast"
	"github.com/ElaraLang/elara/lexer"
	"testing"
)

func TestBasicParsing(t *testing.T) {
	code := "let a = 5 * 3 + 3 * 8\n"
	tokens := lexer.Lex(code)
	inputChannel := make(chan lexer.Token)
	outChannel := make(chan ast.Statement)
	errChannel := make(chan ParseError)
	parser := NewParser(inputChannel, outChannel, errChannel)
	go postLexedTokens(inputChannel, tokens)
	go parser.Parse()
	stmt, err := collectParserResult(&parser)
	close(inputChannel)
	close(outChannel)
	close(errChannel)

	fmt.Println(stmt[0].ToString())
	fmt.Println(err)
}

func postLexedTokens(inChannel chan lexer.Token, tokens []lexer.Token) {
	for _, v := range tokens {
		inChannel <- v
	}
	inChannel <- lexer.Token{
		TokenType: lexer.EOF,
		Text:      nil,
		Position:  lexer.Position{},
	}
}
