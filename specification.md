# Elara Programming Language Specification

## Chapter 0 - Grammar Syntax

This specification uses a modified form of
[ANTLR4 Grammar Syntax](https://github.com/antlr/antlr4/blob/master/doc/index.md)
to describe the grammar of Elara.

In many cases, natural language may be used instead of a
strictly ANTLR-compliant syntax for the sake of simplicity.

Literal patterns (for example a pattern matching the character sequence `/*`)
should have their characters separated by spaces
(the previous example becomes `/ *`).

## Chapter 1 - Lexical Structure

### 1.1 - Unicode

Elara programs are encoded in the [Unicode character set](https://unicode.org)
and almost all Unicode characters are supported in source code.

Elara Source Code should either be encoded in a UTF-8 or UTF-16 format.

### 1.2 - Line Terminators

Elara recognises 3 different options for line terminator characters:
`CR`, `LF`, or `CR` immediately followed by `LF`.

Each of these patterns are treated as 1 single line terminator.
These are used to determine when expressions start and end,
and will determine line numbers in errors produced.

```antlr
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

Single line comments are denoted with the literal text `//`.
Any source code following from this pattern until a **Line Terminator**
character is encountered is ignored.

```antlr
SingleLineComment:
    / / InputCharacter+
```

#### Multi line comments

Multi line comments start with `/*` and continue until `*/` is encountered.
If no closing comment is found (i.e `EOF` is reached before a `*/`),
an error should be raised by the compiler.

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

If an identifier starts with an uppercase character,
it is implied to be referencing a type's identifier.

On the other hand, starting with lowercase implies a "binding" identifier
(that is, referring to some named expression)

This means that all type names must start with uppercase
(excluding generics), and all bindings must start with lowercase.

### 1.6 Operator Identifiers

Elara supports custom operators and operator overloading.
A separate type of identifier is defined for these.
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

Additionally, they must not match any of the following patterns:

- `=`
- `/*` or `*/`
- `//`

When referenced in any context apart from infix application,
the operator's identifier must also be surrounded in parentheses.
For example, to define an operator `/=` we do

```fsharp
let (/=) x y = [implementation]
```

```antlr
InvalidOperatorSymbol:
    | =
    | / *
    | * /
    | / /

OperatorSymbol: Any character described in "Valid Operator Symbols" except InvalidOperatorSymbol

QualifiedOperator: ( UnqualifiedOperator )

UnqualifiedOperator: OperatorSymbol+
```

### 1.7 - Number Literals

Number Literals are unlimited sequences of numeric characters.

For clarity, any number literals may contain `_` which can be used in place of
a comma or dot in real world numbers. These should be ignored by the lexer and
do not affect the resultant number in any way.
For example, the literal `21_539` is functionally identical to `21539`

Number Literals are translated to values of some type implementing
the `Num` type class. That is, all number literals are polymorphic by default.

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

Decimal Integer Literals follow a base-10 notation and
so any characters between `0` and `9` are accepted.

##### 1.7.1.2 - Hexadecimal Integer Literals

Hexadecimal Integer Literals represent numbers in a base-16 format.
Any number literal directly preceeded by a `0x` sequence of characters
should be treated as hexadecimal.
Any characters between `0` and `9`, along with upper or lowercase letters
between `a` and `f` are accepted.

##### 1.7.1.3 - Binary Integer Literals

These represent numbers in a base-2 format. Any number literal directly preceeded
by a `0b` sequence of characters should be treated as binary.
These only accept 2 characters: `0` and `1`

#### 1.7.2 - Floating Point Literals

Floating point literals represent numbers with a decimal point.

Only base-10 notation is supported in floating points.
That is, Hexadecimal and Binary floats are illegal and should
raise a compiler error if a `.` is encountered.

### 1.8 - Character Literals

```antlr
EscapeSequence:
    | \ b (backspace character, Unicode \u0008)
    | \ s (space character, Unicode \u0020)
    | \ h (horizontal tab, Unicode \u0009)
    | \ n (newline LF, Unicode \u000a)
    | \ f (form feed FF, Unicode \u000c)
    | \ r (carriage return CR, Unicode \u000d)
    | \ " (double quote ")
    | \ ' (single quote ')
    | \ \ (backslash \)
    | UnicodeEscape

UnicodeEscape: \ u {HexadecimalDigit, 4}

CharacterLiteralValue: EscapeSequence | Any Raw Unicode character

CharacterLiteral: ' CharacterLiteralValue '
```

It is a compile time error for a character literal to contain
more than 1 `CharacterLiteralValue` in between the single quotes,
or for a `LineTerminator` character to appear between the quotes.

`UnicodeEscape` values may only represent UTF-16 code units.
They are limited to `\u0000` to `\uffff`.

Character literals are converted to values of the `Char` type.

### 1.9 - String Literals

A String Literal consists of zero or more characters
surrounded by double quotes (`"`).
Special characters (such as newlines) must be represented
as escapes, rather than their literal values.

```antlr
StringCharacter:
    | Any Unicode Character except [", \]
    | EscapeSequence

StringLiteral: " StringCharacter+ "
```

If any Line Terminator Character appears between the opening and closing `"`
, an compiler error will be thrown.

String literals are translated to values of the `[Char]` type (a list of `Char`s).

### 1.10 - Text Blocks

A Text Block ("Multiline String") consists of zero or more characters
surrounded by 3 double quotes (`"""`). Similar to String Literals,
special characters must be represented with escape sequences.

```antlr
TextBlockCharacter:
    | Any Unicode Character except \
    | EscapeSequence

TextBlock: " " " TextBlockCharacter+ " " "
```

The output of a Text Block should trim any consistent indentation
at compile time, and replace raw `LineTerminator` characters with
their corresponding escape sequences. For example, the following code:

```fsharp
let message = """
    hello
    world
    """
```

should be translated to a normal String Literal equivalent to the following: `"\nhello\nworld\n"`

#### 1.10.1 - Raw Text Blocks

If the indentation trimming functionality of a standard Text Block
is not desired, we can use a Raw Text Block.
These are defined as zero or more characters directly preceeded by
`!"""` and directly terminated by `"""!`.

They behave identically to normal text blocks,
except not attempting to trim any indentation.

Adjusting the previous example to use a raw text block gives:

```fsharp
let message = !"""
    hello
    world
    """!
```

Which translates to the string literal `"\n    hello\n    world\n"`

### 1.11 - Keywords

Keywords are special sequences of characters that are  allowed as a
standard identifier. They should also generally be treated as separate tokens.

The following sequences of characters are reserved as keywords
and are not permitted for use as identifiers:

- `let`
- `def`
- `type`
- `in`
- `where`
- `class`
- `instance`
- `if`
- `else`
- `then`

## Chapter 2 - Parser and Grammar

### 2.1 - Bindings

Bindings can be declared anywhere in a source file with the following syntax:

```fsharp
    let [identifier] = [value]
```

Bindings declared in this form are accessible from anywhere in the code
that is after the binding declaration, and in the same, or a child scope.

Bindings are **NOT** expressions, and so they cannot be recursively defined.

#### 2.1.1 - Verbose Binding Syntax

An alternative binding form is provided that **is** an expression:

```fsharp
    let [identifier] = [value] in [expression]
```

In this form, the binding is only present in the scope of `expression`.
The entire construct evaluates to the result of `expression`.

#### 2.1.2 - Binding Type Declarations

Bindings can optionally declare an expected type of their value.

This is done with a colon (`:`) after the `[identifier]` section
followed by the expected type.

For example:

```fsharp
    let [identifier]: [Type] = [value]
```

If `value` (or `expression`) does not evaluate to a value of type `Type`,
a compiler error should be raised.

### 2.2 - Functions

The syntax for defining functions is very similar to bindings,
however function definitions also define the parameters of the function.

```fsharp
    let [identifier] [param1] [param2] [...] = [function_body]
```

where `param1`, `param2`, ... `paramN` are valid [Identifiers](#15---normal-identifiers)
and `function_body` is any expression.

For example:

```fsharp
   let f param1 param2 = param1 + param2
```

#### 2.2.1 - Multiple Expression Function Bodies

Functions can also have multiple expressions in their body.
This is denoted by a new line, and then indentation, so that every line of
the body lines up with the function name.

For example:

```fsharp
    let f param1 param2 =
        print param1
        print param2
        let double = param1 * 2
        double + param2
```

Every line in the body is evaluated, and the last line is returned.
Therefore, a function body must end with an expression.

##### 2.2.1.1 - Semicolons in Multiple Expression Function Bodies

Semicolons (`;`) can also be used as an alternative or supplement to line terminators:

```fsharp
    let f param1 param2 =
        print param1 ; print param2
        let double = param1 * 2
        double + param2
```

#### 2.2.2 - Def forms

It is often useful to explicitly specify the type of a function,
for drafting out implementations or for readability in future.

Whilst the let binding colon syntax can be used for functions, it is often messy.
For example:

```fsharp
    let add a b : Int -> Int -> Int = a + b
```

To avoid this issue, the type of the function can be declared on the line above
its implementation.

```fsharp
    def [identifier] : [type]
```

For example:

```fsharp
    def add : Int -> Int -> Int
    let add a b = a + b
```

This form can also be used for standard bindings (i.e not functions):

```fsharp
    def pi : Double
    let pi = 3.141529
```

The def line must be directly above the let line

#### 2.2.3 - Function Types

In Elara, a function is a mapping from a single input value to a single output value. 

There are 2 types of functions, which are categorised based on their **purity** (for a function to be pure, it must have no _visible_ side effects)

The purity of a function can always be inferred by the compiler based on its context. If a function _only_ calls other pure functions, it is itself pure. Otherwise, it is impure.

##### 2.2.3.1 - Pure Functions

A pure function mapping a value of type `a` to a value of type `b` has the type `a -> b`. 
A pure function has no side effects, meaning a compiler is permitted to replace the implementation with an optimised version. For example the function may be automatically memoised or inlined for performance benefits. This functionality can be prevented with the annotation `@NoInline`.

A pure function may only call other pure functions. Calling an impure function makes this function itself impure.

##### 2.2.3.2 - Impure Functions

An impure function behaves similarly to a pure function, but is expected to perform side effects. An impure function mapping a value of type `a` to a value of type `b` has the type `a => b`.
Because side effects are expected, compilers should be conservative when optimising impure functions to avoid subtle changes in behaviour.

An impure function may call any other function, pure or impure.

### 2.2.4 - Function Application

Calling functions is done with the form `function arg`, where `arg` is an expression and `function` is some pure or impure function

This applies the function to the given argument and evaluates to some value whose type is the function's return type.

Function application is left associative, so `f a b` is the same as `(f a) b`.


#### 2.2.5 - The main function

In order for an Elara program to compile to an executable, it must define a main function that serves as an entry point for the program.

This function must have the type `() => ()` and be named `main`.

If the main function is present in a project, the compiler must emit an executable. Similarly, if this function is not present, the compiler is not required to emit an executable (although it is not prohibited).

If multiple files in a project have the main function, the ambiguity must be resolved by the user. The process of resolving this is undefined as it is left up to the compiler implementation.


### 2.4 - Lists

Lists are written `[elem1, elem2, ..., elemN]`, where `N` is any integer `N >= 0`.

Lists are homogeneous, meaning all elements in a list must share a single type `a`.
The type of this list is written as `[a]`

The **empty list** is written as `[]` and has type `[a]` where `a` is a generic type. 
This means that the empty list can be used anywhere irrespective of the type of list.

The **cons operator** is written as `:` and is reserved for list construction.
This operator takes an element `a` and a list `[a]` and prepends the element to the head of the list.
For example, `3 : [4]` gives `[3, 4]`.

The cons operator is right associative, so `3 : 4 : [5]` is the same as `3 : (4 : [5])`

List literals are merely a syntax sugar for repeated application of the cons operator. 
For example, `[1, 2, 3]` is de-sugared to `1 : 2 : 3 : []`.

### 2.5 - Tuples

Tuples are similar to lists, except heterogeneous and fixed-length. 

A tuple is written as `(elem1, elem2, ..., elemN)` where `N` is any integer `N >= 2`.

The type of a tuple of length `N` is written as `(t1, t2, ..., tN)`.
For example, the tuple `(3, "Hello")` has type `(Int, String)`

The type constructor of a tuple of length `N` can be de-sugared to a parenthesised series of commas,
where there are `N - 1` commas. 
For example, a tuple of length 2 has the type constructor `(,)`. Thus, `(Int, String)` is identical to `(,) Int String`.

### 2.6 - The Unit Value

The unit value is written as `()` and has type `()`.
It is the only value of type `()`.

The unit value can be used to represent computations with no useful input or output, i.e side effects.

A pure function `f` with the signature `a -> ()` must have by definition `let f _ = ()`,
as no other pure implementations exist.

*Note that this could be thought of as an empty tuple (length 0) if empty tuples were permitted* 
(see section 2.5).

### 2.7 - Generic Types

Generic types are a type that can be any type. 
They are denoted with a lowercase identifier to distinguish them from normal types.
All generic types are assumed to be universally quantified. 
For example, the type expression `a -> a` means `∀a. a -> a`, i.e it must hold for every possible type.

Once a generic type has been "realised", it must remain the same all throughout a type expression.
For example, in the expression `a -> b -> (a, b)`, if `a` is inferred to `Int` then it must be `Int` everywhere in the expression.

The realised type of a generic type can usually be inferred from usage by the compiler.
If not, then they should be kept as generics.
If the realised type can never be inferred, then a compiler error should be raised to prompt
the user to resolve the ambiguity.


### 2.8 - Type Class Constraints
Type classes are Elara's mechanism for polymorphism. 
Their semantics and functionality is described later in this specification (TODO link).

A type expression can have place a **constraint** on some generic types to create a smaller set of possible types.
This is written as `(Class1 i1, Class2 i2, ..., ClassN iN) := expr` 
where `N >= 1`, `iN` is a valid Generic Type identifier, `ClassN` is the identifier of a type class,
and `expr` is some Type Expression.

For example, to constrain a function `print : a => ()` to any type that is an instance of the `Show` class, 
we could write `(Show a) := a => ()`.

### 2.9 - Standard Types

Standard types are the simplest form of type. 
These are simply some valid type identifier (see Section 1.5.1) 
that matches the name of any one of the following:

- Alias Type
- Algebraic Data Type
- Record Type

### 2.10 - Type Expressions

A type expression is some combination of types and constraints that forms an expression that describes
to a single type.

Valid elements in a type expression include:
- Lists (Section 2.4)
- Tuples (Section 2.5)
- The Unit Type (Section 2.6)
- Generic Types (Section 2.7)
- Type Class Constraints (Section 2.8)
- Standard Types (Section 2.9)
- Function Types (Section 2.2.3)

A type expression must have at least 1 element in it, but can be arbitrarily sized otherwise.

Examples of valid type expressions include:

- `(Show a) := [a] => ()`
- `(Int, String) -> (String, Int)`
- `String`
- `a -> a`
- `a -> b -> Int -> [(a, b)]`
- `a`

