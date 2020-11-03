package main

import (
	"github.com/ElaraLang/elara/base"
	"testing"
)

var code = `let fact = (Int n) => {
    if n == 1 => return 1
    return n * fact(n - 1)
}
fact(8)`

func BenchmarkElara(b *testing.B) {
	for i := 0; i < b.N; i++ {
		base.Execute(nil, code, false)
	}
}
