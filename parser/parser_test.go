package parser

import (
	"fmt"
	"github.com/ElaraLang/elara/ast"
	"github.com/ElaraLang/elara/lexer"
	"testing"
)

func TestBasicParsing(t *testing.T) {
	code := "let n = 3\n" + "print(n)\n"

	runeChannel := make(chan rune)
	inputChannel := make(chan lexer.Token)
	outChannel := make(chan ast.Statement)
	errChannel := make(chan ParseError)

	parser := NewParser(inputChannel, outChannel, errChannel)
	go postRunes(runeChannel, code)
	go lexer.Lex(runeChannel, inputChannel)
	go parser.Parse()
	stmt, err := collectParserResult(&parser)
	close(inputChannel)
	close(outChannel)
	close(errChannel)
	close(runeChannel)
	printStmt(&stmt)
	printError(&err)
}

func postRunes(inChannel chan rune, code string) {
	inp := []rune(code)
	for _, v := range inp {
		inChannel <- v
	}
	inChannel <- rune(-1)
}

func printStmt(s *[]ast.Statement) {
	for _, v := range *s {
		fmt.Println(v.ToString())
	}
}

func printError(s *[]ParseError) {
	for _, v := range *s {
		fmt.Println(v.ErrorToken.String() + " " + v.Message)
	}
}
