package tests

import (
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

let daves = produceDaves(50)
print(daves)
	`
	executeTest(code)
}
