package lexer

import (
	"strings"
	"testing"
)

var code = strings.Repeat(`
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
		`, 10_000)

func BenchmarkEverySymbol(b *testing.B) {

	for n := 0; n < b.N; n++ {
		_ = *Lex(nil, code)
	}
}

func BenchmarkScannerReading1(b *testing.B) {
	code := "314.159265359"
	for n := 0; n < b.N; n++ {
		reader := NewTokenReader([]rune(code))
		_, _ = reader.readNumber()
	}
}

func BenchmarkScannerReading2(b *testing.B) {
	code := "314.159265359"
	for n := 0; n < b.N; n++ {
		reader := NewTokenReader([]rune(code))
		_, _ = reader.readNumberNew()
	}
}
