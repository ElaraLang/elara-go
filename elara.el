namespace examples/main
import elara/std

let content = read-file("data.txt") split "\n" |> map to-int

content zip content |> filter (a, b) => a + b == 2020 |> for-each print