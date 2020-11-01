namespace examples/main
import elara/std

struct Test {
    String name
}

extend Test {
    let b = 3
    let print-name => print(name)
}

let test = Test("Test Struct")
test.print-name()
print(test.b)