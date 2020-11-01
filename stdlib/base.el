namespace elara/std

let print = (Any msg) => Unit {
    print-raw(msg + "\n")
}

let print-raw = (Any msg) => {
    stdout.write(msg)
}
