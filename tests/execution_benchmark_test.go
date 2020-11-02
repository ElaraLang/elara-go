package tests

import (
	"github.com/ElaraLang/elara/base"
	"testing"
)

var code = `let b = if true => "yes" else => "no"
	b`

func BenchmarkSimpleExecution(b *testing.B) {

	for n := 0; n < b.N; n++ {
		_, _, _, _ = base.Execute(nil, code, false)
	}
}
