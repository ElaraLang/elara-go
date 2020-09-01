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
