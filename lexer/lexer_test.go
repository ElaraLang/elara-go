package lexer

import (
	"reflect"
	"testing"
)

func TestIntAssignmentLexing(t *testing.T) {
	code := "let a = 30"
	tokens := *Lex(nil, code)

	expectedTokens := []Token{
		CreateToken(Let, "let", CreatePosition(nil, 0, 0)),
		CreateToken(Identifier, "a", CreatePosition(nil, 0, 4)),
		CreateToken(Equal, "=", CreatePosition(nil, 0, 6)),
		CreateToken(Int, "30", CreatePosition(nil, 0, 8)),
	}

	if !reflect.DeepEqual(tokens, expectedTokens) {
		t.Errorf("Incorrect lexing output, got %v but expected %v", tokens, expectedTokens)
	}
}

func TestFloatAssignmentLexing(t *testing.T) {
	code := "let a = 3.5"
	tokens := *Lex(nil, code)

	expectedTokens := []Token{
		CreateToken(Let, "let", CreatePosition(nil, 0, 0)),
		CreateToken(Identifier, "a", CreatePosition(nil, 0, 4)),
		CreateToken(Equal, "=", CreatePosition(nil, 0, 6)),
		CreateToken(Float, "3.5", CreatePosition(nil, 0, 8)),
	}

	if !reflect.DeepEqual(tokens, expectedTokens) {
		t.Errorf("Incorrect lexing output, got %v but expected %v", tokens, expectedTokens)
	}
}

func TestStringAssignmentLexing(t *testing.T) {
	code := `let a = "Hello"`
	tokens := *Lex(nil, code)

	expectedTokens := []Token{
		CreateToken(Let, "let", CreatePosition(nil, 0, 0)),
		CreateToken(Identifier, "a", CreatePosition(nil, 0, 4)),
		CreateToken(Equal, "=", CreatePosition(nil, 0, 6)),
		CreateToken(String, "Hello", CreatePosition(nil, 0, 8)),
	}

	if !reflect.DeepEqual(tokens, expectedTokens) {
		t.Errorf("Incorrect lexing output, got %v but expected %v", tokens, expectedTokens)
	}
}

func TestBooleanAssignmentLexing(t *testing.T) {
	code := `let a = true`
	tokens := *Lex(nil, code)

	expectedTokens := []Token{
		CreateToken(Let, "let", CreatePosition(nil, 0, 0)),
		CreateToken(Identifier, "a", CreatePosition(nil, 0, 4)),
		CreateToken(Equal, "=", CreatePosition(nil, 0, 6)),
		CreateToken(Boolean, "true", CreatePosition(nil, 0, 8)),
	}

	if !reflect.DeepEqual(tokens, expectedTokens) {
		t.Errorf("Incorrect lexing output, got %v but expected %v", tokens, expectedTokens)
	}
}

func TestSimpleFunctionLexing(t *testing.T) {
	code := `let a = () => {}`
	tokens := *Lex(nil, code)

	expectedTokens := []Token{
		CreateToken(Let, "let", CreatePosition(nil, 0, 0)),
		CreateToken(Identifier, "a", CreatePosition(nil, 0, 4)),
		CreateToken(Equal, "=", CreatePosition(nil, 0, 6)),
		CreateToken(LParen, "(", CreatePosition(nil, 0, 8)),
		CreateToken(RParen, ")", CreatePosition(nil, 0, 9)),
		CreateToken(Arrow, "=>", CreatePosition(nil, 0, 11)),
		CreateToken(LBrace, "{", CreatePosition(nil, 0, 14)),
		CreateToken(RBrace, "}", CreatePosition(nil, 0, 15)),
	}

	if !reflect.DeepEqual(tokens, expectedTokens) {
		t.Errorf("Incorrect lexing output, got %v but expected %v", tokens, expectedTokens)
	}
}

func TestHelloWorldLexing(t *testing.T) {
	code := `let hello-world => print "Hello World"
             hello-world()`
	tokens := *Lex(nil, code)

	expectedTokens := []Token{
		CreateToken(Let, "let", CreatePosition(nil, 0, 0)),
		CreateToken(Identifier, "hello-world", CreatePosition(nil, 0, 4)),
		CreateToken(Arrow, "=>", CreatePosition(nil, 0, 16)),
		CreateToken(Identifier, "print", CreatePosition(nil, 0, 19)),
		CreateToken(String, "Hello World", CreatePosition(nil, 0, 25)),
		CreateToken(NEWLINE, "\n", CreatePosition(nil, 0, 36)),
		//Note: These add 13 because there are 13 spaces in the raw string before the call
		CreateToken(Identifier, "hello-world", CreatePosition(nil, 1, 13+0)),
		CreateToken(LParen, "(", CreatePosition(nil, 1, 13+11)),
		CreateToken(RParen, ")", CreatePosition(nil, 1, 13+12)),
	}

	if !reflect.DeepEqual(tokens, expectedTokens) {
		t.Errorf("Incorrect lexing output, got %v but expected %v", tokens, expectedTokens)
	}
}

func TestBracketLexing(t *testing.T) {
	code := `()[]{}<>`
	tokens := *Lex(nil, code)

	expectedTokens := []Token{
		CreateToken(LParen, "(", CreatePosition(nil, 0, 0)),
		CreateToken(RParen, ")", CreatePosition(nil, 0, 1)),
		CreateToken(LSquare, "[", CreatePosition(nil, 0, 2)),
		CreateToken(RSquare, "]", CreatePosition(nil, 0, 3)),
		CreateToken(LBrace, "{", CreatePosition(nil, 0, 4)),
		CreateToken(RBrace, "}", CreatePosition(nil, 0, 5)),
		CreateToken(LAngle, "<", CreatePosition(nil, 0, 6)),
		CreateToken(RAngle, ">", CreatePosition(nil, 0, 7)),
	}

	if !reflect.DeepEqual(tokens, expectedTokens) {
		t.Errorf("Incorrect lexing output, got %v but expected %v", tokens, expectedTokens)
	}
}

func TestOperatorLexing(t *testing.T) {
	code := `+ - * / % && || ^ == != > >= < <= !`
	tokens := *Lex(nil, code)

	expectedTokens := []Token{
		CreateToken(Add, "+", CreatePosition(nil, 0, 0)),
		CreateToken(Subtract, "-", CreatePosition(nil, 0, 2)),
		CreateToken(Multiply, "*", CreatePosition(nil, 0, 4)),
		CreateToken(Slash, "/", CreatePosition(nil, 0, 6)),
		CreateToken(Mod, "%", CreatePosition(nil, 0, 8)),
		CreateToken(And, "&&", CreatePosition(nil, 0, 10)),
		CreateToken(Or, "||", CreatePosition(nil, 0, 13)),
		CreateToken(Xor, "^", CreatePosition(nil, 0, 16)),
		CreateToken(Equals, "==", CreatePosition(nil, 0, 18)),
		CreateToken(NotEquals, "!=", CreatePosition(nil, 0, 21)),
		CreateToken(RAngle, ">", CreatePosition(nil, 0, 24)),
		CreateToken(GreaterEqual, ">=", CreatePosition(nil, 0, 26)),
		CreateToken(LAngle, "<", CreatePosition(nil, 0, 29)),
		CreateToken(LesserEqual, "<=", CreatePosition(nil, 0, 31)),
		CreateToken(Not, "!", CreatePosition(nil, 0, 34)),
	}

	if !reflect.DeepEqual(tokens, expectedTokens) {
		t.Errorf("Incorrect lexing output %v but expected %v", tokens, expectedTokens)
	}
}

func TestUnderscoreLexing(t *testing.T) {
	code := `_`
	tokens := *Lex(nil, code)

	expectedTokens := []Token{
		CreateToken(Underscore, "_", CreatePosition(nil, 0, 0)),
	}

	if !reflect.DeepEqual(tokens, expectedTokens) {
		t.Errorf("Incorrect lexing output, got %v but expected %v", tokens, expectedTokens)
	}
}
