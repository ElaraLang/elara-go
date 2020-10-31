namespace examples/calculator
print("Welcome to the official Elara calculator. Please input the first number:")

let a = input()
let left = a.to-int()

print("Please input the operation: +, -, *, /")
let op = input()

print("Please input the second number:")

let b = input()
let right = b.to-int()

let result = if op == "+" => left + right
             else if op == "-" => left - right
             else if op == "*" => left * right
             else if op == "/" => left / right
             else => "Unknown operation " + op

print("== Result == ")
print(result)