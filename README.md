<p align="center">
  <img alt="elara logo" src="https://github.com/ElaraLang/Elara-Old/blob/master/4.jpg?raw=true">
</p>

# Elara
https://discord.gg/xu5gSTV

Elara is a multi-paradigm (although primarily functional) language with a clean syntax with little "noise"
 
Unlike most functional languages, Elara will not restrict the programmer for the sake of it,
and prioritises developer freedom wherever possible.

## Basic Syntax
*(all syntax is subject to change)*
### Variable Declaration
Variables are declared with the syntax 

`let [name] = [value]` with type inference

For example, `let age = 23`

Explicit type specification:

`let [name]: [Type] = [value]`

For example, `let name: String = "Bob"` 

By default, variables are reference immutable.
For mutability, the syntax `let mut [name] = [value]` should be used.
This is discouraged as immutable values should always be preferred


### Function Declaration

Functions are first class types, and are declared in a near identical way to variables:

Note that the **Arrow Syntax** (`=>`) is used with functions to distinguish from a function *call*
```
let printHello() => { 
  print("Hello World!")
}
```

Functions with a single expression can be declared with a more concise syntax:
```
let printHello => print("Hello World!")
```

Functions with parameters:
```
let printTwice(String message) => {
    print(message)
    print(message)
}
```

Functions with a clearly defined return type:
```
let add(Int a, Int b) => Int {
    a + b
}
```

### Function Calling

Elara supports a wide range of function calling syntax to try and make programming more natural and less restrictive:

#### Simple Calling
`printTwice("Hello")`

`printHello()`


#### OOP style calling on a receiver function
`"Hello".printTwice()`

This feature works with multiple parameters:
```
let addTo(Int a, Int b) => {
    a + b
}

3.addTo(4)
addTo(3, 4)
```

the 2 calls are identical

#### Infix Style calling
You can also omit the parentheses and commas with infix functions (functions with 2 parameters):
```
3 addTo 4
```

### Collections
Elara has collection literals for the 2 main types:

#### Lists
List literals are a comma separated list of elements, surrounded by square brackets

- Empty List: `[]`
- Single Element Lists: `[1]`
- Multi Element Lists: `[1, 2, 3, 4]`

Lists are immutable, and the recommended implementation is a persistent one to make copying more efficient

Lists should aim to be as homogeneous as possible - that is, 
Lists should try to form a union of all elements' types to form the List's type

List types are declared in the format `[ElementType]`
For example `[Any]`, `[Int]`, `[() => Unit]`

#### Maps 
Map literals are a comma separated list of **Entries**, surrounded by curly brackets

Entries are composed of a Key and a Value, separated by a colon.
An Entry's Key and Value must both be valid expressions

- Empty Map: `{}`
- Single Element Map: `{a: "b"}` (this assumes a variable named `a` is present in the current scope)
- Multi Element Map: 
```
{ 
    a: "b",
    c: "d"
}
```
(Again, this assumes the presence of `a` and `c`)

Maps are also immutable, and are typically implemented as a hash table.

Map types follow the format `{K : V}`
For example: `{Int : String}`, `{String : () => Unit}`, `{Person : Int}`

### Structs

Structs in Elara are **Data Only**
They are declared with the following syntax:
```
struct Person {
    String name
    mut Int age
    Int height = 110
}
``` 

And can be instantiated like so:
`let mark = Person("Mark", 32, 160)`


### OOP Structs
Structs can easily replicate objects with extension syntax, which is the most idiomatic way of adding functionality to structs:

```
struct Person {
    //blah
}
extend Person {
    let celebrateBirthday => {
        print("Happy Birthday " + name + "!")
        age += 1
    }
} 
```

from here we can do `somePerson.celebrateBirthday()` as if it was a method.

The `extend` syntax works with any type and can be done from any file

#### "Inheritance"
The `extend` syntax effectively adds inheritance too:

```
struct Person {
    //blah
}
extend Person {
    struct Student {
        Topic major
    }
}
```

This is not true "inheritance". Instead, the `Student` struct will copy all the properties of `Person`.
Because the type system is contract based (that is, type `B` can be assigned to type `A` if it has the same contract in its members),
this is effectively inheritance - we can use an instance of `Student` wherever we use a `Person`.
### Type System

Elara features a simple, linear type system. 
`Any` is the root of the type hierarchy, with subtypes such as `Int`, `String` and `Person`.

However, there are also a few quirks that aim to make the type system more flexible:

**Contract based type parameters**

Type parameters for generics support contract based boundary.
Take for example the simple generic function, ignoring the unnecessary generic (since T can be any type):
```
#T
let printAndReturn(T data) => T {
    print(data)
    return data
}
``` 
 
We cannot guarantee that every type will give a user-friendly value for `print`.

To work around this, we can add a boundary to `T`, that only accepts types that define a `toString` function:

```
<T { toString() => String } >
let printAndReturn(T data) => T {
    print(data.toString())
    return data
}
```

This gives programmers extra flexibility in that they can program to a specific contract, rather than a type


### Namespaces and Importing

The namespace system in Elara is simple

Declaring a namespace is usually done at the top of the file:
`namespace elara/core`

Namespaces follow the format `base/module`, similar to languages like Clojure

Importing a namespace is simple:
`import elara/core`

The files will now have access to the contents of all files in that namespace

### Functional Features
* Lambdas are defined identically to functions:
`let lambda = (Type name) => {}`
Parameter types can be omitted if possible to infer from context

* Functions are first class
```
let add1 = (Int a) => a + 1

let added1List = someList map add1
```

* Function chaining is trivial
```
someList.map(add1).filter(isEven).forEach(print)
```
### Conclusion

Elara is in very early stages, with the evaluator being nowhere near finished.

The eventual plan includes:
    * Static typing
    * Compiling to native code, but also supporting other backends such as JavaScript or JVM Bytecode
    * A proper standard library
    * Allowing type inference for function parameters

