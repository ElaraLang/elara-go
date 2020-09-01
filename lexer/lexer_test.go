package lexer

import (
	"reflect"
	"testing"
)

func TestIntAssignmentLexing(t *testing.T) {
	code := "let a = 30"
	tokens := lex(code)

	expectedTokens := []Token{
		CreateToken(LET, "let"),
		CreateToken(IDENTIFIER, "a"),
		CreateToken(EQUAL, "="),
		CreateToken(INT, "30"),
	}

	if !reflect.DeepEqual(tokens, expectedTokens) {
		t.Errorf("Incorrect lexing output, got %v but expected %v", tokens, expectedTokens)
	}
}

func TestFloatAssignmentLexing(t *testing.T) {
	code := "let a = 3.5"
	tokens := lex(code)

	expectedTokens := []Token{
		CreateToken(LET, "let"),
		CreateToken(IDENTIFIER, "a"),
		CreateToken(EQUAL, "="),
		CreateToken(FLOAT, "3.5"),
	}

	if !reflect.DeepEqual(tokens, expectedTokens) {
		t.Errorf("Incorrect lexing output, got %v but expected %v", tokens, expectedTokens)
	}
}

func TestStringAssignmentLexing(t *testing.T) {
	code := `let a = "Hello"`
	tokens := lex(code)

	expectedTokens := []Token{
		CreateToken(LET, "let"),
		CreateToken(IDENTIFIER, "a"),
		CreateToken(EQUAL, "="),
		CreateToken(STRING, "Hello"),
	}

	if !reflect.DeepEqual(tokens, expectedTokens) {
		t.Errorf("Incorrect lexing output, got %v but expected %v", tokens, expectedTokens)
	}
}

func TestSimpleFunctionLexing(t *testing.T) {
	code := `let a = () => {}`
	tokens := lex(code)

	expectedTokens := []Token{
		CreateToken(LET, "let"),
		CreateToken(IDENTIFIER, "a"),
		CreateToken(EQUAL, "="),
		CreateToken(LPAREN, "("),
		CreateToken(RPAREN, ")"),
		CreateToken(ARROW, "=>"),
		CreateToken(LBRACE, "{"),
		CreateToken(RBRACE, "}"),
	}

	if !reflect.DeepEqual(tokens, expectedTokens) {
		t.Errorf("Incorrect lexing output, got %v but expected %v", tokens, expectedTokens)
	}
}

func TestHelloWorldLexing(t *testing.T) {
	code := `let hello-world => print "Hello World" 
             hello-world()`
	tokens := lex(code)

	expectedTokens := []Token{
		CreateToken(LET, "let"),
		CreateToken(IDENTIFIER, "hello-world"),
		CreateToken(ARROW, "=>"),
		CreateToken(IDENTIFIER, "print"),
		CreateToken(STRING, "Hello World"),
		CreateToken(IDENTIFIER, "hello-world"),
		CreateToken(LPAREN, "("),
		CreateToken(RPAREN, ")"),
	}

	if !reflect.DeepEqual(tokens, expectedTokens) {
		t.Errorf("Incorrect lexing output, got %v but expected %v", tokens, expectedTokens)
	}
}
