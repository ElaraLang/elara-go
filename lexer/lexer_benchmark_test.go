package lexer

import (
	"strings"
	"testing"
)

var result []Token

func BenchmarkHelloWorldAssignment(b *testing.B) {
	var lexed []Token
	code := strings.Repeat(`let hello-world = "hello world"\n`, 1_000)
	for n := 0; n < b.N; n++ {
		lexed = lex(code)
	}
	result = lexed
}
