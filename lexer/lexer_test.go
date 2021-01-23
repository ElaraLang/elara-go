package lexer

import (
	"reflect"
	"testing"
)

func TestIntAssignmentLexing(t *testing.T) {
	code := "let a = 30"
	tokens := Lex(code)

	expectedTokens := []Token{
		CreateToken(Let, "let", CreatePosition(0, 0)),
		CreateToken(Identifier, "a", CreatePosition(0, 4)),
		CreateToken(Equal, "=", CreatePosition(0, 6)),
		CreateToken(Int, "30", CreatePosition(0, 8)),
	}

	if !reflect.DeepEqual(tokens, expectedTokens) {
		t.Errorf("Incorrect lexing output, got %v but expected %v", tokens, expectedTokens)
	}
}

func TestFloatAssignmentLexing(t *testing.T) {
	code := "let a = 3.5"
	tokens := Lex(code)

	expectedTokens := []Token{
		CreateToken(Let, "let", CreatePosition(0, 0)),
		CreateToken(Identifier, "a", CreatePosition(0, 4)),
		CreateToken(Equal, "=", CreatePosition(0, 6)),
		CreateToken(Float, "3.5", CreatePosition(0, 8)),
	}

	if !reflect.DeepEqual(tokens, expectedTokens) {
		t.Errorf("Incorrect lexing output, got %v but expected %v", tokens, expectedTokens)
	}
}

func TestStringAssignmentLexing(t *testing.T) {
	code := `let a = "Hello"`
	tokens := Lex(code)

	expectedTokens := []Token{
		CreateToken(Let, "let", CreatePosition(0, 0)),
		CreateToken(Identifier, "a", CreatePosition(0, 4)),
		CreateToken(Equal, "=", CreatePosition(0, 6)),
		CreateToken(String, "Hello", CreatePosition(0, 8)),
	}

	if !reflect.DeepEqual(tokens, expectedTokens) {
		t.Errorf("Incorrect lexing output, got %v but expected %v", tokens, expectedTokens)
	}
}

func TestBooleanAssignmentLexing(t *testing.T) {
	code := `let a = true`
	tokens := Lex(code)

	expectedTokens := []Token{
		CreateToken(Let, "let", CreatePosition(0, 0)),
		CreateToken(Identifier, "a", CreatePosition(0, 4)),
		CreateToken(Equal, "=", CreatePosition(0, 6)),
		CreateToken(BooleanTrue, "true", CreatePosition(0, 8)),
	}

	if !reflect.DeepEqual(tokens, expectedTokens) {
		t.Errorf("Incorrect lexing output, got %v but expected %v", tokens, expectedTokens)
	}
}

func TestSimpleFunctionLexing(t *testing.T) {
	code := `let a = () => {}`
	tokens := Lex(code)

	expectedTokens := []Token{
		CreateToken(Let, "let", CreatePosition(0, 0)),
		CreateToken(Identifier, "a", CreatePosition(0, 4)),
		CreateToken(Equal, "=", CreatePosition(0, 6)),
		CreateToken(LParen, "(", CreatePosition(0, 8)),
		CreateToken(RParen, ")", CreatePosition(0, 9)),
		CreateToken(Arrow, "=>", CreatePosition(0, 11)),
		CreateToken(LBrace, "{", CreatePosition(0, 14)),
		CreateToken(RBrace, "}", CreatePosition(0, 15)),
	}

	if !reflect.DeepEqual(tokens, expectedTokens) {
		t.Errorf("Incorrect lexing output, got %v but expected %v", tokens, expectedTokens)
	}
}

func TestHelloWorldLexing(t *testing.T) {
	code := `let hello-world => print "Hello World"
             hello-world()`
	tokens := Lex(code)

	expectedTokens := []Token{
		CreateToken(Let, "let", CreatePosition(0, 0)),
		CreateToken(Identifier, "hello-world", CreatePosition(0, 4)),
		CreateToken(Arrow, "=>", CreatePosition(0, 16)),
		CreateToken(Identifier, "print", CreatePosition(0, 19)),
		CreateToken(String, "Hello World", CreatePosition(0, 25)),
		CreateToken(NEWLINE, "\n", CreatePosition(0, 36)),
		//Note: These add 13 because there are 13 spaces in the raw string before the call
		CreateToken(Identifier, "hello-world", CreatePosition(1, 13+0)),
		CreateToken(LParen, "(", CreatePosition(1, 13+11)),
		CreateToken(RParen, ")", CreatePosition(1, 13+12)),
	}

	if !reflect.DeepEqual(tokens, expectedTokens) {
		t.Errorf("Incorrect lexing output, got %v but expected %v", tokens, expectedTokens)
	}
}

func TestBracketLexing(t *testing.T) {
	code := `()[]{}<>`
	tokens := Lex(code)

	expectedTokens := []Token{
		CreateToken(LParen, "(", CreatePosition(0, 0)),
		CreateToken(RParen, ")", CreatePosition(0, 1)),
		CreateToken(LSquare, "[", CreatePosition(0, 2)),
		CreateToken(RSquare, "]", CreatePosition(0, 3)),
		CreateToken(LBrace, "{", CreatePosition(0, 4)),
		CreateToken(RBrace, "}", CreatePosition(0, 5)),
		CreateToken(LAngle, "<", CreatePosition(0, 6)),
		CreateToken(RAngle, ">", CreatePosition(0, 7)),
	}

	if !reflect.DeepEqual(tokens, expectedTokens) {
		t.Errorf("Incorrect lexing output, got %v but expected %v", tokens, expectedTokens)
	}
}

func TestOperatorLexing(t *testing.T) {
	code := `+ - * / % && || ^ == != > >= < <= !`
	tokens := Lex(code)

	expectedTokens := []Token{
		CreateToken(Add, "+", CreatePosition(0, 0)),
		CreateToken(Subtract, "-", CreatePosition(0, 2)),
		CreateToken(Multiply, "*", CreatePosition(0, 4)),
		CreateToken(Slash, "/", CreatePosition(0, 6)),
		CreateToken(Mod, "%", CreatePosition(0, 8)),
		CreateToken(And, "&&", CreatePosition(0, 10)),
		CreateToken(Or, "||", CreatePosition(0, 13)),
		CreateToken(Xor, "^", CreatePosition(0, 16)),
		CreateToken(Equals, "==", CreatePosition(0, 18)),
		CreateToken(NotEquals, "!=", CreatePosition(0, 21)),
		CreateToken(RAngle, ">", CreatePosition(0, 24)),
		CreateToken(GreaterEqual, ">=", CreatePosition(0, 26)),
		CreateToken(LAngle, "<", CreatePosition(0, 29)),
		CreateToken(LesserEqual, "<=", CreatePosition(0, 31)),
		CreateToken(Not, "!", CreatePosition(0, 34)),
	}

	if !reflect.DeepEqual(tokens, expectedTokens) {
		t.Errorf("Incorrect lexing output %v but expected %v", tokens, expectedTokens)
	}
}

func TestUnderscoreLexing(t *testing.T) {
	code := `_`
	tokens := Lex(code)

	expectedTokens := []Token{
		CreateToken(Underscore, "_", CreatePosition(0, 0)),
	}

	if !reflect.DeepEqual(tokens, expectedTokens) {
		t.Errorf("Incorrect lexing output, got %v but expected %v", tokens, expectedTokens)
	}
}

func TestHashLexing(t *testing.T) {
	code := "#"
	tokens := Lex(code)

	expectedTokens := []Token{
		CreateToken(Hash, "#", CreatePosition(0, 0)),
	}

	if !reflect.DeepEqual(tokens, expectedTokens) {
		t.Errorf("Incorrect lexing output, got %v but expected %v", tokens, expectedTokens)
	}
}
