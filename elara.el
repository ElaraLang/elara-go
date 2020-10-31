namespace examples/calculator
import elara/std

struct Test {
    String name
    Int i
    Int extra
}

struct Test2 {
    String name
    Int i
}

let t: Test2 = Test("Hello", 3, 4)
print(t.name)
print(t.i)
print(t.extra)