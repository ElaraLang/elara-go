package main

import (
	"elara/base"
	"testing"
)

var code = `let some-func = () => {
    if true => {
        return 4
    } else => {
        return 5
    }
}

some-func() * some-func()`

func BenchmarkElara(b *testing.B) {
	for i := 0; i < b.N; i++ {
		base.Execute(nil, code, false)
	}
}
