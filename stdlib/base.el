namespace elara/std

let print = (Any msg) => {
    stdout.print(msg + "\n")
}
let print-raw = (Any msg) => {
    stdout.print(msg)
}

let test => print("Hi")
