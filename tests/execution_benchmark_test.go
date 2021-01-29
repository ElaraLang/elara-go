package tests

import (
	"testing"
)

var code = `let b = if true => "yes" else => "no"
	b`

func BenchmarkSimpleExecution(b *testing.B) {
	/*
		for n := 0; n < b.N; n++ {
			res, _, _, _ := base.Execute(nil, code, false)
			if res[1].Value != "yes" {
				b.Fail()
			}
		}*/
}
