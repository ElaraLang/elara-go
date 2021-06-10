# Elara Programming Language Specification

## Chapter 1 - Lexical Structure

### 1.1 - Unicode

Elara programs are encoded in the [Unicode character set](https://unicode.org) and almost all Unicode characters are supported in source code.

### 1.2 - Line Terminators

Elara recognises 3 different options for line terminator characters: `CR`, `LF`, or `CR` immediately followed by `LF`.
Each of these patterns are treated as 1 single line terminator. These are used to determine when expressions start and end, and will detemine line numbers in errors produced.

### 1.3 - Whitespace

White space is defined as any of the following:

- The ASCII ` ` character, "space"
- The ASCII ` ` character, "tab"

### 1.4 - Comments

Comment Properties:

- Comments cannot be nested
- 1 comment type's syntax has no special meaning in another comment
- Comments do not apply in string literals

Elara uses 2 main formats for comments:

#### Single Line Comments

Single line comments are denoted with the literal text `//`. Any source code following from this pattern until a **Line Terminator** character is encountered is ignored.

#### Multi line comments

Multi line comments start with `/*` and continue until `*/` is encountered. If no closing comment is found (i.e `EOF` is reached), an error should be raised by the compiler

### 1.5 - Normal Identifiers

Identifiers are unlimited length sequences of any of the following characters:

- Any characters of the alphabet, upper or lower case
- Any denary digit

Additionally, a valid identifier must satisfy all of the following:

- Must start with a character that is **NOT** a digit (i.e 0-9)
- Must not directly match any reserved Elara keywords

#### 1.5.1 - Type vs Binding Identifiers

If an identifier starts with an uppercase character, it is implied to be referencing a type's identifier.
On the other hand, starting with lowercase implies a "binding" identifier (that is, referring to some named expression)

This implies that all type names must start with Uppercase (excluding generics), and all bindings must start with lowercase.

### 1.6 Operator Identifiers

Elara supports custom operators and operator overloading. A separate type of identifier is defined for these.
These identifiers may only consist of the following symbols:

- `.`
- `>`
- `<`
- `=`
- `-`
- `/`
- `+`
- `*`
- `&`
- `|`
- `!`
- `?`
- `%`
- `^`
- `:`
- `~`
- `#`
- `_`

Additionally, they must not match any of the following "native operator" tokens:

- `=`
- `/*`
- `//`

When referenced in any context apart from infix application, the operator's identifier must also be surrounded in parentheses. For example, to define an operator `/=` we do ```fs
let (/=) x y = [implementation]

```

```
