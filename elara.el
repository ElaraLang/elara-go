namespace examples/what
import elara/std

struct A {
    String name
}

extend A {
    let b = 3
    let print-name => print(name)
}

let a = A("Test")
a.print-name()
print(a.b)