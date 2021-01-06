package tests

import (
	"github.com/ElaraLang/elara/base"
	"testing"
)

func BenchmarkElara(b *testing.B) {
	code := `namespace test/lol
	import elara/std

	let c = [1, 2, 3, 4]
	let plus1(Int i) => { i + 1 }
	let times3(Int i) => { i * 3 }
	let isEven(Int i) => { i % 2 == 0 }
	map(c, plus1)
	c.map(times3)
	c.filter(isEven)
	`
	base.LoadStdLib()
	for i := 0; i < b.N; i++ {
		base.Execute(nil, code, false)
	}
}
