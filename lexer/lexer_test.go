package lexer

import (
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
