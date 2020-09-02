package lexer

import (
	"strings"
	"testing"
)

var result []Token

func BenchmarkEverySymbol(b *testing.B) {
	var lexed []Token
	code := strings.Repeat(`
		let a = 30
		let a = 3.5
		let a = "Hello"
		let a = true
		let a = () => {}
		let hello-world => print "Hello World"
		hello-world()
		()[]{}<>

		+ - * / % && || ^ == != > >= < <= !
		, : _
		let	mut	struct if else match while
		`, 1_000)
	for n := 0; n < b.N; n++ {
		lexed = lex(code)
	}
	result = lexed
}
