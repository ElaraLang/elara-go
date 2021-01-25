package lexer

import (
	"fmt"
	"reflect"
	"testing"
)

const everySymbol = `
	(){}<>[]
	let extend return while mut lazy open struct namespace import type if else match as is try
	id
	+ - * / % && || ^ == != >= <= !
	| &
	= => .
	true
	false
	"Hello"
	'H'
	30
	0x12f
	0b10101
	3.465
	, :
	`

func BenchmarkLexer(t *testing.B) {
	for i := 0; i < t.N; i++ {
		input := everySymbol
		inputChan := make(chan rune)
		outputChan := make(chan Token, len(input))
		go Lex(inputChan, outputChan)
		go func() {
			for _, c := range input {
				inputChan <- c
			}
			inputChan <- eof //Signal to the channel that the input has terminated
		}()

		for {
			c := <-outputChan
			if c.TokenType == EOF {
				break
			}
		}
	}

}

func testLex(input string) []Token {
	inputChan := make(chan rune)
	outputChan := make(chan Token)
	go Lex(inputChan, outputChan)
	go func() {
		for _, c := range input {
			inputChan <- c
		}
		inputChan <- eof //Signal to the channel that the input has terminated
	}()
	output := make([]Token, 0)
	for {
		c := <-outputChan
		fmt.Printf("%+v\n", c)
		output = append(output, c)
		if c.TokenType == EOF {
			break
		}
	}
	return output
}

func TestHexNumberLexing(t *testing.T) {
	input := "0x3F"
	output := testLex(input)
	expectedOutput := []Token{
		{TokenType: HexadecimalInt, Data: []rune{'3', 'F'}, Line: 0, Col: 4},
		{TokenType: EOF, Data: nil, Line: 0, Col: 5},
	}
	if !reflect.DeepEqual(output, expectedOutput) {
		t.Fail()
	}
}

func TestDecNumberLexing(t *testing.T) {
	input := "04724"
	output := testLex(input)
	expectedOutput := []Token{
		{TokenType: DecimalInt, Data: []rune{'4', '7', '2', '4'}, Line: 0, Col: 5},
		{TokenType: EOF, Data: nil, Line: 0, Col: 6},
	}
	if !reflect.DeepEqual(output, expectedOutput) {
		t.Fail()
	}
}

func TestBinNumberLexing(t *testing.T) {
	input := "0b10110"
	output := testLex(input)
	expectedOutput := []Token{
		{TokenType: BinaryInt, Data: []rune{'1', '0', '1', '1', '0'}, Line: 0, Col: 7},
		{TokenType: EOF, Data: nil, Line: 0, Col: 8},
	}
	if !reflect.DeepEqual(output, expectedOutput) {
		t.Fail()
	}
}

func TestLetDeclarationLexing(t *testing.T) {
	input := "let a = 3"
	output := testLex(input)
	expectedOutput := []Token{
		{TokenType: Let, Data: nil, Line: 0, Col: 3},
		{TokenType: Identifier, Data: []rune{'a'}, Line: 0, Col: 5},
		{TokenType: Equal, Data: nil, Line: 0, Col: 7},
		{TokenType: DecimalInt, Data: []rune{'3'}, Line: 0, Col: 9},
		{TokenType: EOF, Data: nil, Line: 0, Col: 10},
	}
	if !reflect.DeepEqual(output, expectedOutput) {
		t.Fail()
	}
}

func TestLetMutDeclarationLexing(t *testing.T) {
	input := "let mut a = 3"
	output := testLex(input)
	expectedOutput := []Token{
		{TokenType: Let, Data: nil, Line: 0, Col: 3},
		{TokenType: Mut, Data: nil, Line: 0, Col: 7},
		{TokenType: Identifier, Data: []rune{'a'}, Line: 0, Col: 9},
		{TokenType: Equal, Data: nil, Line: 0, Col: 11},
		{TokenType: DecimalInt, Data: []rune{'3'}, Line: 0, Col: 13},
		{TokenType: EOF, Data: nil, Line: 0, Col: 14},
	}
	if !reflect.DeepEqual(output, expectedOutput) {
		t.Fail()
	}
}
func TestEveryTokenLexing(t *testing.T) {
	input := everySymbol
	output := testLex(input)
	expectedOutput := []Token{
		{TokenType: LParen, Data: nil}, {TokenType: RParen, Data: nil}, {TokenType: LBrace, Data: nil}, {TokenType: RBrace, Data: nil}, {TokenType: LAngle, Data: nil}, {TokenType: RAngle, Data: nil}, {TokenType: LSquare, Data: nil}, {TokenType: RSquare, Data: nil},
		{TokenType: Let, Data: nil}, {TokenType: Extend, Data: nil}, {TokenType: Return, Data: nil}, {TokenType: While, Data: nil}, {TokenType: Mut, Data: nil}, {TokenType: Lazy, Data: nil}, {TokenType: Open, Data: nil}, {TokenType: Struct, Data: nil}, {TokenType: Namespace, Data: nil}, {TokenType: Import, Data: nil}, {TokenType: Type, Data: nil}, {TokenType: If, Data: nil}, {TokenType: Else, Data: nil}, {TokenType: Match, Data: nil}, {TokenType: As, Data: nil}, {TokenType: Is, Data: nil}, {TokenType: Try, Data: nil},
		{TokenType: Identifier, Data: []rune{'i', 'd'}},
		{TokenType: Add, Data: nil}, {TokenType: Subtract, Data: nil}, {TokenType: Multiply, Data: nil}, {TokenType: Slash, Data: nil}, {TokenType: Mod, Data: nil}, {TokenType: And, Data: nil}, {TokenType: Or, Data: nil}, {TokenType: Xor, Data: nil}, {TokenType: Equals, Data: nil}, {TokenType: NotEquals, Data: nil}, {TokenType: GreaterEqual, Data: nil}, {TokenType: LesserEqual, Data: nil}, {TokenType: Not, Data: nil},
		{TokenType: TypeOr, Data: nil}, {TokenType: TypeAnd, Data: nil},
		{TokenType: Equal, Data: nil}, {TokenType: Arrow, Data: nil}, {TokenType: Dot, Data: nil},
		{TokenType: BooleanTrue, Data: nil}, {TokenType: BooleanFalse, Data: nil},
		{TokenType: String, Data: []rune{'H', 'e', 'l', 'l', 'o'}},
		{TokenType: Char, Data: []rune{'H'}},
		{TokenType: DecimalInt, Data: []rune{'3', '0'}},
		{TokenType: HexadecimalInt, Data: []rune{'1', '2', 'f'}},
		{TokenType: BinaryInt, Data: []rune{'1', '0', '1', '0', '1'}},
		{TokenType: Float, Data: []rune{'3', '.', '4', '6', '5'}},
		{TokenType: Comma, Data: nil},
		{TokenType: Colon, Data: nil},
		{TokenType: EOF, Data: nil},
	}
	if !eq(filterWhitespace(output), expectedOutput) {
		t.Fail()
	}
}
func filterWhitespace(tokens []Token) []Token {
	filtered := make([]Token, 0)
	for _, t := range tokens {
		if t.TokenType != NEWLINE {
			filtered = append(filtered, t)
		}
	}
	return filtered
}
func eq(output []Token, expected []Token) bool {
	if len(output) != len(expected) {
		return false
	}
	for i, token := range output {
		expect := expected[i]
		if expect.TokenType != token.TokenType || !reflect.DeepEqual(expect.Data, token.Data) {
			return false
		}
	}
	return true
}
