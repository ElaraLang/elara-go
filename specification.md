# Elara Programming Language Specification

## Chapter 0 - Grammar Syntax

This specification uses a modified form of [ANTLR4 Grammar Syntax](https://github.com/antlr/antlr4/blob/master/doc/index.md) to describe the grammar of Elara.

In many cases, natural language may be used instead of a strict ANTLR-compliant syntax for the sake of simplicity.

Literal patterns (for example a pattern matching the character sequence `/*`) should have their characters separated by spaces (the previous example becomes `/ *`).

## Chapter 1 - Lexical Structure

### 1.1 - Unicode

Elara programs are encoded in the [Unicode character set](https://unicode.org) and almost all Unicode characters are supported in source code.

Elara Source Code should either be encoded in a UTF-8 or UTF-16 format.

### 1.2 - Line Terminators

Elara recognises 3 different options for line terminator characters: `CR`, `LF`, or `CR` immediately followed by `LF`.

Each of these patterns are treated as 1 single line terminator.
These are used to determine when expressions start and end, and will determine line numbers in errors produced.

```antlr4
LineTerminator:
    | The ASCII CR Character
    | The ASCII LF Character
    | The ASCII CR Character immediately followed by the ASCII LF Character
```


### 1.3 - Whitespace

```antlr
WhiteSpace:
    | The ASCII SP Character, ` `
    | The ASCII HT Character, `\t`
    | LineTerminator

InputCharacter: Any Unicode Character except WhiteSpace
```

### 1.4 - Comments

Comment Properties:

- Comments cannot be nested
- 1 comment type's syntax has no special meaning in another comment
- Comments do not apply in string literals

Elara uses 2 main formats for comments:

#### Single Line Comments

Single line comments are denoted with the literal text `//`. Any source code following from this pattern until a **Line Terminator** character is encountered is ignored.

```antlr
SingleLineComment:
    / / InputCharacter+ 
```


#### Multi line comments

Multi line comments start with `/*` and continue until `*/` is encountered. 
If no closing comment is found (i.e `EOF` is reached before a `*/`), an error should be raised by the compiler

```antlr
StartMultiLineComment: / *

EndMultiLineComment: * /

MultiLineComment: StartMultiLineComment InputCharacter+ EndMultiLineComment
```

### 1.5 - Normal Identifiers

Identifiers are unlimited length sequences of any of the following characters:

- Any characters of the alphabet, upper or lower case
- Any denary digit

Additionally, a valid identifier must satisfy all of the following:

- Must start with a character that is **NOT** a digit (i.e 0-9)
- Must not directly match any reserved Elara keywords

```antlr
IdentifierCharacter: Any characters of the latin alphabet, upper or lowercase 

IdentifierCharacterOrDigit: IdentifierCharacter or 0-9

Identifier: IdentifierCharacter IdentifierCharacterOrDigit+ but not a Keyword
```

#### 1.5.1 - Type vs Binding Identifiers

If an identifier starts with an uppercase character, it is implied to be referencing a type's identifier.
On the other hand, starting with lowercase implies a "binding" identifier (that is, referring to some named expression)

This implies that all type names must start with Uppercase (excluding generics), and all bindings must start with lowercase.

### 1.6 Operator Identifiers

Elara supports custom operators and operator overloading. A separate type of identifier is defined for these.
These identifiers may only consist of the following symbols:

#### Valid Operator Symbols
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
- `\`

Additionally, they must not match any of the following "native operator" patterns:

- `=`
- `/*`
- `//`

When referenced in any context apart from infix application, the operator's identifier must also be surrounded in parentheses. For example, to define an operator `/=` we do 
```fs
let (/=) x y = [implementation]
```

```antlr
InvalidOperatorSymbol: 
    | =
    | / *
    | / /

OperatorSymbol: Any character described in "Valid Operator Symbols" except InvalidOperatorSymbol

QualifiedOperator: ( UnqualifiedOperator )

UnqualifiedOperator: OperatorSymbol+
```

### 1.7 - Number Literals
Number Literals are unlimited sequences of numeric characters. 

For clarity, any number literals may contain `_` which can be used in place of a comma or dot in real world numbers. These should be ignored by the lexer and do not affect the resultant number in any way. For example, the literal `21_539` is functionally identical to `21539`

```antlr
DecimalDigit: 0 | 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 | 9

BinaryDigit: 0 | 1

HexadecimalDigit: DecimalDigit | a | b | c | d | e | f | A | B | C | D | E | F

Separator: _

DecimalLiteral: (DecimalDigit | Separator)+

HexadecimalPrefix: 0x

HexadecimalLiteral: HexadecimalPrefix (HexadecimalDigit | Separator)+

BinaryPrefix: 0b

BinaryLiteral: BinaryPrefix (BinaryDigit | Separator)+


FloatingPointSeparator: .

FloatingPoint: DecimalLiteral FloatingPointSeparator DecimalLiteral
```

#### 1.7.1 - Integer Literals

##### 1.7.1.1 - Decimal Integer Literals
Decimal Integer Literals follow a base-10 notation and so any characters between `0` and `9` are accepted.

##### 1.7.1.2 - Hexadecimal Integer Literals
Hexadecimal Integer Literals represent numbers in a base-16 format. Any number literal directly preceeded by a `0x` sequence of characters should be treated as hexadecimal. Any characters between `0` and `9`, along with upper or lowercase letters between `a` and `f` are accepted.

##### 1.7.1.3 - Binary Integer Literals
These represent numbers in a base-2 format. Any number literal directly preceeded by a `0b` sequence of characters should be treated as binary. These only accept 2 characters: `0` and `1`

#### 1.7.2 - Floating Point Literals

Floating point literals represent numbers with a decimal point.

Only base-10 notation is supported in floating points. That is, Hexadecimal and Binary floats are illegal.
