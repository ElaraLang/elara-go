namespace examples/main
import elara/std

let fact = (Int n) => {
    if n == 1 => return 1
    return n * fact(n - 1)
}

print(fact(3) is Int)
