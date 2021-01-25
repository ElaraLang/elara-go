package lexer

import (
	"reflect"
	"testing"
)

func BenchmarkLexer(t *testing.B) {
	for i := 0; i < t.N; i++ {
		input := "'h' + '2'"
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
