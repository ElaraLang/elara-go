namespace examples/calculator
import elara/std

struct Test {
    String name
    Int i
    Int extra
}

let t = Test("Hello", 3, 4)
print(t.name)
print(t.i)
print(t.extra)