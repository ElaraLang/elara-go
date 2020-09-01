package lexer

import (
	"reflect"
	"testing"
)

func TestIntAssignmentLexing(t *testing.T) {
	code := "let a = 30"
	tokens := lex(code)

	expectedTokens := []Token{
		CreateToken(Let, "let"),
		CreateToken(Identifier, "a"),
		CreateToken(Equal, "="),
		CreateToken(Int, "30"),
	}

	if !reflect.DeepEqual(tokens, expectedTokens) {
		t.Errorf("Incorrect lexing output, got %v but expected %v", tokens, expectedTokens)
	}
}

func TestFloatAssignmentLexing(t *testing.T) {
	code := "let a = 3.5"
	tokens := lex(code)

	expectedTokens := []Token{
		CreateToken(Let, "let"),
		CreateToken(Identifier, "a"),
		CreateToken(Equal, "="),
		CreateToken(Float, "3.5"),
	}

	if !reflect.DeepEqual(tokens, expectedTokens) {
		t.Errorf("Incorrect lexing output, got %v but expected %v", tokens, expectedTokens)
	}
}

func TestStringAssignmentLexing(t *testing.T) {
	code := `let a = "Hello"`
	tokens := lex(code)

	expectedTokens := []Token{
		CreateToken(Let, "let"),
		CreateToken(Identifier, "a"),
		CreateToken(Equal, "="),
		CreateToken(String, "Hello"),
	}

	if !reflect.DeepEqual(tokens, expectedTokens) {
		t.Errorf("Incorrect lexing output, got %v but expected %v", tokens, expectedTokens)
	}
}

func TestBooleanAssignmentLexing(t *testing.T) {
	code := `let a = true`
	tokens := lex(code)

	expectedTokens := []Token{
		CreateToken(Let, "let"),
		CreateToken(Identifier, "a"),
		CreateToken(Equal, "="),
		CreateToken(Boolean, "true"),
	}

	if !reflect.DeepEqual(tokens, expectedTokens) {
		t.Errorf("Incorrect lexing output, got %v but expected %v", tokens, expectedTokens)
	}
}

func TestSimpleFunctionLexing(t *testing.T) {
	code := `let a = () => {}`
	tokens := lex(code)

	expectedTokens := []Token{
		CreateToken(Let, "let"),
		CreateToken(Identifier, "a"),
		CreateToken(Equal, "="),
		CreateToken(LParen, "("),
		CreateToken(RParen, ")"),
		CreateToken(Arrow, "=>"),
		CreateToken(LBrace, "{"),
		CreateToken(RBrace, "}"),
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
		CreateToken(Let, "let"),
		CreateToken(Identifier, "hello-world"),
		CreateToken(Arrow, "=>"),
		CreateToken(Identifier, "print"),
		CreateToken(String, "Hello World"),
		CreateToken(Identifier, "hello-world"),
		CreateToken(LParen, "("),
		CreateToken(RParen, ")"),
	}

	if !reflect.DeepEqual(tokens, expectedTokens) {
		t.Errorf("Incorrect lexing output, got %v but expected %v", tokens, expectedTokens)
	}
}

func TestBracketLexing(t *testing.T) {
	code := `()[]{}<>`
	tokens := lex(code)

	expectedTokens := []Token{
		CreateToken(LParen, "("),
		CreateToken(RParen, ")"),
		CreateToken(LSquare, "["),
		CreateToken(RSquare, "]"),
		CreateToken(LBrace, "{"),
		CreateToken(RBrace, "}"),
		CreateToken(Lesser, "<"),
		CreateToken(Greater, ">"),
	}

	if !reflect.DeepEqual(tokens, expectedTokens) {
		t.Errorf("Incorrect lexing output, got %v but expected %v", tokens, expectedTokens)
	}
}
