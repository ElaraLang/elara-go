package tests

import (
	"github.com/ElaraLang/elara/base"
	"testing"
)

func BenchmarkElara(b *testing.B) {
	code := `namespace test/lol
	import elara/std
	
struct Person {
    String name
    Int age
}

let daveFactory = () => { Person("Dave", 50) }

let produceDaves(Int amount) => {
    let factories = [daveFactory] * amount
    factories.map(run)
}

let daves = produceDaves(2147483647)
	`
	base.LoadStdLib()
	for i := 0; i < b.N; i++ {
		//base.Execute("", code, false)
	}
}
