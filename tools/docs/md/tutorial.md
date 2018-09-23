# Golem Tutorial

Welcome to the tutorial for the Golem Programming Language.

* [Hello, world.](#hello-world.)
* [Basic Types](#basic-types)
* [Comments](#comments)
* [Variables](#variables)
* [Collections](#collections)
* [Fields](#fields)
* [Control Structures](#control-structures)
* [Operators and Expressions](#operators-and-expressions)
* [Functions and Closures](#functions-and-closures)
* [The `main()` function](#the-main-function)
* [Structs](#structs)
* [Properties](#properties)
* [Combining Structs Together](#combining-structs-together)
* [Error Handling](#error-handling)
* [Concurrency](#concurrency)
* [Immutabilty](#immutabilty)
* [Type Introspection](#type-introspection)
* [Modules](#modules)
* [Standard Library](#standard-library)

## Hello, world.

To get started, you must first compile a version of the Golem
Command Line Interface, or "CLI".  This requires that you have the Go language 
toolchain installed on your system, with at least version 1.9.

Clone the Golem [repository](https://github.com/mjarmy/golem-lang) into the proper
place in your go development environment, `cd` into the top level directory, 
and type `make`.  This will build Golem, and place the `golem` CLI executable 
in a sub-directory called `build`.

Golem's CLI  doesn't have a 
[REPL](https://en.wikipedia.org/wiki/Read%E2%80%93eval%E2%80%93print_loop), so 
follow along by entering Golem source code into a text file 
('tutorial.glm', for example), and then running the file from 
the command line to look at the results.

So, fire up your text editor of choice, and type the following 
program into a file named 'tutorial.glm':

```
println('Hello, world.');
```

and then run it like so:

```
./build/golem tutorial.glm
```

The function `println` is built in to the Golem CLI.  There are [several](#TODO)of these 
builtin functions in Golem.

If you are a vim user, there is a [vim plugin](https://github.com/mjarmy/golem-lang-vim) for 
Golem that you can install, that provides syntax highlighting.

You may have noticed that there is a semicolon at the end of the println 
statement.  Semicolons are optional in Golem --  you are not required to include a
semicolon after each statement in your code, although you can do so if you 
wish. 

However, semicolons are required if you want to write multiple separate statements 
on a single line:

```
print('Hello, '); println('world.')
```

We will omit any optional semicolons in the rest of this tutorial.

## Basic Types

Golem's basic primitive types include boolean, string, int and float.  There 
is also `null`, which represents the absence of a value.  Basic 
values are immutable.

The two boolean values are `true` and `false`.

Integer values in Golem are signed 64 bit integers, and Float values are 
IEEE-754 64-bit floating-point numbers.  Ints are coerced to floats during 
arithmetic and checks for equality:

```
assert(12 / 4.0 == 3.0)
assert(12 / 4.0 == 3)
```

Note that we used another builtin function, `assert`, which accepts a boolean value,
and will throw an exception if the value that is passed into it is not `true`.

Golem has the usual set of C-language-family operators that you would 
expect: `==`, `!=`, `||`, `&&`, `<`, `>`, `+`, `-`, and so forth.  

```
println(1 + 2)
println(42 / 7)
```

We will cover the operators in more detail later.  

Strings can be delimited either with a single quote or a double quote:

```
println('abc')
println("abc")
```

When adding two values together, if one of the values is a string, and the other 
is not, then the other value is converted to a string, and the two strings are then 
concatenated together:

```
println('a' + 1)
```

Unlike many other dynamic languages, Golem has no concept of 'truthiness'.  The only 
things that are true or false are boolean values:

```
assert(true)
assert(!false)
```

So, `''`, `0`, `null`, etc. are *not* boolean, and and error will be thrown 
if you attempt to evaluate them in a place where a boolean value is
expected.

```
assert('')
```

Another builtin function, `str`, returns  the string representation of a value:

```
assert(str(3) == '3')
```

## Comments

Golem uses C-language-family comments:  `/* ... */` for a block comment, and `//` for 
a line comment.

## Variables

Values can be assigned to variables. Variables are declared via either the `let` 
or `const` keyword.  It is an error to refer to a variable before it has been
declared.

```
let a = 1
const b = 2
a = b + 3
println(a)
```

As you might expect, the value of a const variable cannot be changed once it has
been assigned.

`let` and `const` are "statements" -- they do not return a value.  Assignments, on the
other, are "expressions", and evaluating an expression returns a value:

```
let a = 1
let b = (a = 2)
assert(a == b && b == 2)
```

## Collections

Golem has four collection data types: List, Dict, Set, and Tuple.

You can create a list by enclosing a comma-delimited sequence of values in
square brackets.  Once you've created a list, you can use square brackets to access
individual elements of a list (this is called the "index operator").  

```
let a = []
let b = [3,4,5]
assert(a.isEmpty())
assert(b[0] == 3)
```

If the index value is negative, values will be indexed from the end of the list.
This is a really handy way to get the last element of a list.
```
let a = [1,2,3]
println(a[-1])
```

Use the "slice operator" to create a new list from part of an existing list or string.
If you leave off the first or last value of the slice operation, the resulting slice
will start at the beginning or end.  Negative values work with slices in the same 
way that they do with lists.

```
let c = [4,5,6,7,8]
println(c[1:3])
println(c[:3])
println(c[2:])
println(c[1:-1])
```

Indexing and slicing works on strings too:

```
println('abc'[1])
println('abc'[:-1])
```

Golem's `dict` type is an 
[associative array](https://en.wikipedia.org/wiki/Associative_array).  The keys of a 
`dict` can be any value that supports hashing (str, int, float, bool, or tuple). 

```
let a = dict {'x': 1, 'y': 2}
assert(a['x'] == 1)
```

A `set` is a unordered collection of distinct values.  Any value that can act as a key 
in a dict can be a member of a set.

```
let a = set {'x', 'y'}
assert(a.contains('x'))
```

A `tuple` is an immutable list-like data structure.  Tuples must have at least two values.

```
let a = (1, 2)
assert(a[0] == 1)
```

The builtin function `len` can be used to get the length of any of the collections.
`len` will also return the length of a string.

```
let a = [1, 2, 3]
let b = 'lmnop'
let c = dict {"x": 3}

assert([len(a), len(b), len(c)] == [3,5,1])
```

## Fields

A "field" is a built-in named value that is associated with a parent 
value.  Fields are somewhat like what are called "methods" in other 
languages.  Each type has a collection of fields that are associated with a given value.  

As an example, here is  how to use some of the fields that are present on a list value:

```
let ls = []
assert(ls.isEmpty())
ls.add('a')
ls.addAll(['b', 'c'])
println(ls)
println(ls.contains('c'))
println(ls.indexOf('b'))
println(ls.join(','))
```

Note that the fields we see in the above example -- `isEmpty`, `add`, `addAll`,
`contains`, `indexOf`, and `join` -- are all [functions](#TODO).  Most fields 
that are built in to the various Golem types are functions.

See the [Type Reference](#TODO) for a complete description of all the fields 
that are defined on the various types.

## Control Structures

Golem has a familiar set of control structures: `if`, `while`, `switch`, and `for`.

```
let a = 1
while a < 12 {
    if a < 3 {
        a = a + 2
    } else {
        a = 15
    }
}
assert(a == 15)
```

Golem also has `break` and `continue`, which will break out of a `while` or `for` loop,
or continue at the top of the loop, as in other languages.

Golem has 'ternary-if' expressions as well:

```
const a = 10
let b = a < 3 ? 4 : 5
assert(b == 5)
```

`switch` works roughly the same way as it does in other languages, except that you 
can switch on any value, not just integers.  Also, there is no 'fall-through' -- at most
only one case will be executed.  Therefore the `break` keyword is not applicable to 
switches.

```
let a = 'abc'
let b = 0
switch a {
    case 0:
        b = 1
    case 'abc':
        b = 2
    default:
        b = 3
}
assert(b == 2)
```

You can have multiple expressions in a case statement.  The body of the case
will be executed if at least on of the expressions matches:

```
let s = ''
let i = 0
while i < 4 {
    switch i {
        case 0, 1:
            s += 'a'
        case 2:
            s += 'b'
    }
    i++
}
assert(s == 'aab')
```

You can leave the expression out after the `switch` keyword. This lets you switch
on a sequence of boolean case statements, which is sometimes easier to read than
a cascade of 'if, else-if, else-if' statmements:

```
let b = 0
switch {
    case 1 < 2:
        b = 1
    default:
        b = 2
}
assert(b == 1)
```

Golem's `for` statement iterates over 
a sequence of values derived from an 'iterable' value.  Lists, dicts, sets, and 
strings are iterable.

```
let a = [1, 2, 3]
let z = 0
for e in a {
    z += e
}
assert(z == 6)
```

By convention, iterating over a dict produces a sequence of tuples.  You 
can capture the values in the tuple directly in a `for` loop via
"tuple destructuring":

```
let d = dict { "x": 1, "y": 2, "z": 3 }
for (k, v) in d {
    println("key: ", k, ", value: ", v)
}
```

There is one more iterable type -- ranges.  Ranges are created via the `range`
builtin function.  A range is an immutable value that represents a sequence of integers.  

```
let list = ["frog", "cow", "rabbit"]
for i in range(0, len(list)) {
    if list[i] == "cow" {
        println("The cow is at element ", i)
        break
    }
}
```

Note that ranges do not actually contain a list of all the specified integers.  They 
simply represent a sequence that can be iterated over.

## Operators and Expressions

Golem has the following operators, with the following precedence (from low to high):

| Category       | Operators     |
| -------------  | ------------- |
| or             | <code>&#124;&#124;</code>  |
| and            | `&&`  |
| comparative    | `==`, `!=`, `>`, `>=`, `<`, `<=`, `<=>` |
| additive       | `+`, `-`, <code>&#124;</code>, `^` |
| multiplicative | `*`, `/`, `%`, `&`, `<<`, `>>` |
| unary          | `-`, `!`, `~` |
| postfix        | `++`, `--`  |

The 'spaceship' operator, `<=>`, returns -1, 0, or 1 if the left-hand operator is 
less than, equal to, or greater than the right-hand operator: 

```
assert((5 <=> 10) == -1)
```

Note that  `++` and `--` are postfix.  Golem does not have any prefix operators.

Golem also supports 'assignment operators`, which perform an operation and
do an assignment at the same time, e.g.:

```
let a = 1, b = 2
a += b // is the same as a = a + b
```

Here are the assignment operators:

`=+`, `=-`, `=*`, `=/`, `=%`, `=^`, `=&`, `=|`, `=<<`, `=>>`

## Functions and Closures

Functions are first class values in Golem.  They are created with the `fn` keyword. 

```
let a = fn(x) {
    return x * 7
}
assert(a(6) == 42)
```

Functions do not have to have an explicit `return` statement. If there is no `return`,
they will return the last expression that was evaluated.  If no expression is 
evaluated in the function, `null` is returned.

```
let a = fn() {}
let b = fn(x) { x * x; }
assert(a() == null)
assert(b(3) == 9)
```

A `return` statement without a value is syntactically invalid -- all return statements
must include a value to return.

Golem supports [closures](https://en.wikipedia.org/wiki/Closure_(computer_programming)) as 
well -- in fact closures are a fundamental mechanism
in Golem for managing state.  Here is an example of a closure that acts as a
[accumulator generator](http://www.paulgraham.com/accgen.html):

```
let foo = fn(n) {
    return fn(i) {
        return n += i
    }; 
}
let f = foo(4)
assert([f(1), f(2), f(3)] == [5, 7, 10])
```

You can declare the formal parameters of a function to be constant.  In the following
example, the formal parameter 'b' is constant, so it cannot be changed inside the 
function:

```
let a = 1

fn foo(const b) {
    return a += b
}

foo(2)
foo(3)
assert(a == 6)
```

Golem also supports 'lambda' syntax, via the `=>` operator.  Lambdas provide a 
lightweight way to define a function on the fly. The body of a lambda function is a 
single expression.

```
let a = || => 3
let b = |x| => x * x
let c = |x, y| => (x + y) * 5

assert(a() == 3)
assert(b(2) == 4)
assert(c(1, 2) == 15)
```

Consider the following program, in which function `a` calls function `b`:

```
const b = fn() {
    return 42
}
const a = fn() {
    return b()
}
println(a())
```

This program works because `b` is declared before `a`.  However, if we reverse the order
of declarations, we get a compilation error, because `b` has not yet been defined.

```
const a = fn() {
    return b()
}
const b = fn() {
    return 42
}
println(a())
```

It is often the case that we need to allow for forward references like the one above.  Golem 
provides a feature called 'Named Functions' that offers this functionality.  For example:

```
fn a() {
    return b()
}
fn b() {
    return 42
}
assert(a() == 42)
```

Note that the above program is identical in every way to the previous one, except for 
the forward references provided by the named function syntax.

Some of the fields on certain types, like `list`, accept functions as parameters.
Here is an example of how to use the `map`, `reduce`, and `filter` fields:

```
let ls = [1, 2, 3, 4, 5]
let squares = ls.map(|x| => x * x)
let addedUp = ls.reduce(0, |acc, x| => acc + x)
let even = ls.filter(|x| => (x % 2 == 0))
let strings = ls.map(str).reduce('', |acc, x| => acc + x)

assert(squares == [1, 4, 9, 16, 25])
assert(addedUp == 15)
assert(even == [2, 4])
assert(strings = '12345')
```

Functions can be declared with optional parameters as well.  In the following 
example, the `y` parameter defaults to 0 unless the function is invoked with 
two parameters:

```
fn a(x, y = 0) {
    return x + y
}

println(a(1))
println(a(1, 2))
```

You can declare more than one optional parameter, but they all must go at the 
end of the parameter declarations.

Functions can also be declared with 'variadic' parameters.  `println`
is actually a variadic function -- it will accept any number of parameters:

```
println('frog', 'cow', 'rabbit')
```

Use an ellipsis (three dots) to declare a variadic parameter:

```
fn a(x, ls...) {
    ls.map(|e| => e + x)
}

println(a(1))
println(a(1, 2))
println(a(1, 2, 3))
```

The "extra" parameters are gathered together into a list.  A variadic parameter must
always be the last formal parameter.  Also, you cannot mix optional parameters and
variadic parameters in a declaration.

There is a builtin function called `arity` that returns a [struct](#structs) that describes
the [arity](https://en.wikipedia.org/wiki/Arity) of a function.  Here is a snippet
that prints the arity of 3 of the builtin functions we have already used:

```
println(arity(len))
println(arity(range))
println(arity(println))
```

## The `main()` function

You can pass arguments into a Golem CLI program by defining a `main()` function, that
accepts exactly one parameter.  The parameter will always be a list of the
command line arguments.

```
fn main(args) {
    for i in range(0, len(args)) {
        println('argument ', i, ' is "', args[i], '"')
    }
}
```

## Structs

Golem is not an object-oriented language.  It does not have classes, objects, interfaces,
inheritance, or constructors.  What it does have, however, are values 
which we call "structs".

Structs are created via the `struct` keyword.  
```
let s = struct { a: 1, b: 2 }
```

In the above example, we've created a struct that has two fields, `a` and `b`.

Structs are similar to dicts in some ways, but quite different in others.  The field 
names of a struct can only be strings, and furthermore they
must be valid identifiers -- they cannot have spaces or special characters.

The dot operator, `.`, is used on structs to get or set the fields of a struct:

```
let s = struct { a: 1, b: 2 }
assert(s.a == 1)
s.a = 3
assert(s.a == 3)
```

Once a struct is created, it cannot have new fields added
to it, or existing fields removed.  The _values_ associated with the fields can be changed
though, as we saw in the previous example.

The `this` keyword is used in Golem to allow a struct to refer to itself. In Golem,
`this` is only valid inside a struct, and it is always lexically scoped to refer to 
the innermost enclosing struct. 

```
let s = struct { a: 1, b: 2, c: this.a + this.b }
println(s)
assert(s.c == 3)
```

The builtin-function `merge()` can be used to combine an arbitrary number of 
existing structs into a new struct.

```
let a = struct { x: 1, y: 2 }
let b = struct { y: 3, z: 4 }
let c = merge(a, b)

println(a)
println(b)
println(c)

a.x = 10

println(a)
println(b)
println(c) // x is changed here too!
```

If there are any duplicated keys in the structs passed in to 'merge()', then the
value associated with the first such key is used.  

Also, note in the above example that if you change a value in one of the structs passed 
in to merge(), the value changes in the merged struct as well.  That is because the 
all three structs actually share a common set of fields.  We will see later on 
that this behaviour can be quite useful.

## Properties

Structs can have properties defined on them, so that a given field has a 'getter' 
function, and optional 'setter' function.  The getter function must take 0 parameters, 
and the setter function must take 1 parameter.  If the setter function is omitted,
the property is readonly.  Properties are useful for hiding the inner workings
of a struct behind a simpler facade.  Here is an example (which uses 
[try-catch](#error-handling)):

```
let x = 2
let s = struct {

    // 'a' is a readonly property with a getter function.
    a: prop { || => 1 },

    // 'b' is a property with getter and setters functions.
    b: prop { || => x, |v| => x = v },

    // 'a' and 'b' act like normal fields here.
    c: || => this.a + this.b
}

try {
    s.a = 42
    assert(false) // we will never get here.
} catch e {
    println(e.error)
}

assert([s.a, s.b, x, s.c()] == [1, 2, 2, 3])

s.b = 3
assert([s.a, s.b, x, s.c()] == [1, 3, 3, 4])

x = 4
assert([s.a, s.b, x, s.c()] == [1, 4, 4, 5])
```

## Combining Structs Together

By using closures, structs and `merge()`, it is possible to simulate various 
features from other languages, including inheritance, multiple-inheritance, 
prototype chains, and the like.

For instance, consider the following program:

```
fn newRectangle(w, h) {
    return struct {
        width:  prop { || => w, |val| => w = val },
        height: prop { || => h, |val| => h = val },
        area:   || => w * h
    }
}

fn newBox(rect, d) {
    return merge(
        rect, 
        struct {
            depth:  prop { || => d, |val| => d = val },
            volume: || => rect.area() * d
        })
}

let r = newRectangle(2, 3)
let b = newBox(r, 4)

println([b.width, b.height, b.depth, b.area(), b.volume()])
r.width = 5
println([b.width, b.height, b.depth, b.area(), b.volume()])
```

The functions 'newRectangle' and 'newBox' are very much like what one might call 'constructors'
in another language.  The structs that they return have functions as entries 
(e.g. 'area()'), and these functions refer to the 'this' keyword, and to captured 
variables.  As such, the functions are quite a bit like what one might call a 
'method' in another language.

The use of the 'merge()' function to create a box out of a rectangle is similar to
how inheritance is used in other languages.  Does that mean that a Box is a subclass
of a Rectangle?  Not really, no.  There is no such thing as a 'class' in Golem.  However, 
due to the behaviour of merge(), they *are* inter-related in a way that is 
very much like inheritance.

One of the primary goals of the Golem project is to explore the power provided by 
the simple building blocks of functions, closures, structs and merge().  It is hoped
that the simplicity and flexibility of these elements can be used to create a variety
of complex runtime structures that are easy to reason about and use.  

## Error Handling

Golem uses the familiar 'try-catch-finally` syntax that exists in many C-family 
languages.

```
try {
    let z = 4 / 0
}
catch e {
    println(e.error) 
    println(e.stackTrace) 
}
```

The error value in a `catch` clause is always a struct with an `error` field and
a `stackTrace` field.

You can throw an exception using the `throw` keyword, followed by an expression that
evaluates to a string.

```
try {
    throw 'FooError: foo'
}
catch e {
    println(e.error) 
    println(e.stackTrace) 
}
```

There is also a `finally` clause, which is always executed no matter what happens
inside the try block or catch clause:

```
try {
    throw 'FooError: foo'
}
catch e {
    println(e.error) 
    println(e.stackTrace) 
}
finally {
    println('finally')
}
```

A try block must always have either a catch clause, a try clause, or both.

The try-catch mechanism works exactly the same way that it does in Java.  Try statements
can be nested, and errors that are not caught in a function are passed upwards in
the call stack.

## Concurrency

**TODO** explain this a bunch more.

Golem uses the Go Language's [concurrency system](https://tour.golang.org/concurrency/1).  This 
means that Golem has 'goroutines', channels and the ability to send and 
receive messages.  Go's concurrency capabilities are quite powerful, and Golem is capable
of using them all (**TODO** well, just not yet.  We need to figure out how to provide
access to `select`, mutexes and quite a few other things)

```
fn sum(s, c) {
    c.send(s.reduce(0, |acc, x| => acc+x))
}

let s = [7, 2, 8, -9, 4, 0]
let n = len(s)/2
let c = chan()

go sum(s[:n], c)
go sum(s[n:], c)

let result = [c.recv(), c.recv()]
println(result)
```

## Immutabilty

Golem supports immutability via the `freeze()` builtin function, which makes a mutable
value become immutable.  You can check if a value is immutable via the `frozen()`
builtin function. `freeze()` always returns the value that you pass into it.

```
let s = freeze(struct { a: 1, b: 2 })
assert(frozen(s))

try {
    s.a = 0;       // This will throw an error.
    assert(false); // We can't reach this statement.
} catch e {
    println(e.error)
}
```

`freeze()` only has an affect on Lists, Dicts, Sets and Structs.  All other values 
are already immutable, so calling `freeze()` on them has no effect

Immutabilty and concurrency go hand in hand.  By using immutable 
values whenever possible, you can reduce the likelyhood of bugs in 
your concurrency code, and make it much easier to reason about as well.

An important caveat regarding immutability is that although closures, like all functions, 
are immutable, they can still have enclosed state that can be modified.  There
is no way in Golem to freeze a closure after the fact so that it can no longer modify 
any of its captured variables.  It is up to you to manage state properly if you are 
using closures.  

Here is the "accumulator generator" from a previous example.  We 
freeze it this time, but it still has mutable state via the enclosed variable 'n':

```
let foo = freeze(fn(n) { 
    return |i| => n += i
})
let f = foo(4)
assert([f(1), f(2), f(3)] == [5, 7, 10])
```

## Type Introspection

There is a builtin function called `type()` that will return a string 
describing the type of a value.  Here is a program that will print a list 
of every possible type:

```
let values = [
    null, true, "", 0, 0.0, fn(){}, 
    [], range(0,1), (0,1), dict{}, set{}, 
    struct{}, chan()]

let types = values.map(type)

println(types)
```

## Modules

If you've been following along with the tutorial, you have been using a 
file called "tutorial.glm".  The Golem CLI actually compiles
this file into a `module` called "tutorial".  Modules are the fundamental
unit of compilation in Golem, and are also used for namespace 
management. All you need to do to create your own modules in the CLI is
create a file with the name you want.  As an example, lets create a module 
called foo, and reference in the tutorial module.

In a file called "foo.glm", place the following:

```
fn square(x) {
    return x*x
}
```

And then in your "tutorial.glm", you can reference the "foo" module 
via the `import` statement, like
this:

```
import foo

assert(foo.square(5) == 25)
```

## Standard Library

Golem has a [Standard Library](#TODO) that is implemented as a collection of modules. 

When embedding the Golem interpreter in a Go program, some or all of 
the standard library can be included in the sandboxed environment. The 
Golem CLI makes the entire standard library available.

To use one of the modules from the standard library, simply import it like
you would any module, e.g. `import os`.

