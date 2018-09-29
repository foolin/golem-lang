# A Tour of Golem 

Welcome to the tour of the Golem Programming Language.

## Hello, world

Let's get started with the proverbial hello world program.  In this tour, we will 
be running Golem code directly in the browser, so press the `Run` button below
to see the output of the program:

```
println('Hello, world.');
```

You may have noticed that there is a semicolon at the end of the [`println`](builtins.html#println) 
statement.  Semicolons are usually optional in Golem --  they are required only if 
you want to write multiple separate statements on a single line.  We will omit
unnecessary semicolons in the rest of the tour.

## Basic Types

Golem's basic primitive types include [bool](bool.html), [int](int.html), [float](float.html) 
and [string](str.html) There is also [null](null.html), which represents the absence 
of a value.  Basic values are immutable.


```
println(null)
println(false)
println(true)
println(1 + 2)
println(3.0 / 4.0)
println('abc' + "def")
```

Golem has the usual set of C-language-family [operators](syntax.html#operator-precedence) 
that you would expect: `==`, `!=`, `||`, `&&`, `<`, `>`, `+`, `-`, and so forth.  

Unlike many other dynamic languages, Golem has no concept of 'truthiness'.  The only 
things that are true or false are boolean values.  So, `''`, `0`, `null`, etc. 
are *not* boolean, and and error will be thrown if you attempt to evaluate them 
in a place where a boolean value is expected.

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
other, are [expressions](#TODO), and evaluating an expression returns a value:

```
let a = 1
let b = (a = 2)
assert(a == b && b == 2)
```

## Collections

Golem has four collection data types: [List](list.html), 
[Dict](dict.html), [Set](set.html), and [Tuple](tuple.html) 

### List

You can create a list by enclosing a comma-delimited sequence of values in
square brackets.  Once you've created a list, you can use square brackets to access
individual elements of a list (this is called the "index operator").  

```
let a = []
let b = [3,4,5]
assert(a.isEmpty())
assert(b[0] == 3)
println(b[-1]) // negative indexes start from the end
```

Use the "slice operator" to create a new list from part of an existing list.
If you leave off the first or last value of the slice operation, the resulting slice
will start at the beginning or end.  Negative values work with slices in the same 
way that they do with lists.

```
let c = ['x', 'y', 'z']
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

### Dict

Golem's `dict` type is an 
[associative array](https://en.wikipedia.org/wiki/Associative_array).  The keys of a 
dict can be any value that supports [hashing](interfaces.html#hashable). 

```
let a = dict {'x': 1, 'y': 2}
println(a['x'])
a['z'] = 3
println(a)
```

### Set

A `set` is a unordered collection of distinct values.  Any value that can act as a key 
in a dict can be a member of a set.

```
let a = set {'x', 'y'}
println(a)
println(a.contains('x'))
```

### Tuple

A `tuple` is an immutable list-like data structure.  Tuples must have at least two values.

```
let a = (1, 2)
println(a)
```

### `len`

The builtin function [`len`](builtins.html#len) can be used to get the length of any 
of the collections.  `len` will also return the length of a string.

```
let a = [1, 2, 3]
let b = 'lmnop'
let c = dict {"x": 3}
println([len(a), len(b), len(c)])
```

## Fields

A "field" in Golem is a named member of a value.  Each type has a collection 
of fields that are associated with a given value.  

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
`contains`, `indexOf`, and `join` -- are all [functions](#functions).  Most fields 
that are built in to the various Golem types are functions.

## Control Structures

Golem has a familiar set of control structures: `if`, `while`, `switch`, and `for`.

```
let a = 2
while a < 100 {
    a = a * a 
}
println(a)
```

Golem also has `break` and `continue`, which will break out of a `while` or `for` loop,
or continue at the top of the loop, as in other languages.

Golem has 'ternary-if' expressions as well:

```
const a = 10
println(a < 3 ? 4 : 5)
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
println(b)
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
println(s)
```

You can leave the expression out after the `switch` keyword. This lets you switch
on a sequence of boolean case statements, which is sometimes easier to read than
a cascade of 'if, else-if, else-if' statmements:

```
let a = 0
switch {
    case a < 1:
        println('foo')
    default:
        println('bar')
}
```

Golem's `for` statement iterates over a sequence of values derived from 
an [iterable](interfaces.html#iterable) value.

```
let a = [1, 2, 3]
let b = 0
for e in a {
    b += e
}
println(b)
```

By convention, iterating over a dict produces a sequence of tuples.  You 
can capture the values in the tuple directly in a `for` loop via
"tuple destructuring":

```
let d = dict { "x": 1, "y": 2, "z": 3 }
for e in d {
    println(e)
}
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

## Functions

A [Function](func.html) is a sequence of [`expressions`](#TODO) and [`statements`](#TODO) 
that can be invoked to perform a task. We have already encountered quite a few 
functions: builtin functions like `println`, and a few fields like the ones on a list.

### Function Syntax

Functions are created with the `fn` keyword:

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
let b = fn(x) { x * x; } // semicolon required here!
assert(a() == null)
assert(b(3) == 9)
```

A `return` statement without a value is syntactically invalid -- all return statements
must include a value to return.

You can declare the formal parameters of a function to be constant.  In the following
example, the formal parameter 'b' is constant, so it cannot be changed inside the 
function:

```
fn foo(const b) {
    return b + 42
}
println(foo(12))
```

### Lambdas

Golem also supports 'lambda' syntax, via the `=>` operator.  Lambdas provide a 
lightweight way to define a function on the fly. The body of a lambda function is a 
single expression.

```
let a = || => 3
let b = |x| => x * x
let c = |x, y| => (x + y) * 5
println(a())
println(b(2))
println(c(1, 2))
```

Here is an example which passes lambdas to some list functions that expect a function
as a parameter:

```
const ls = [1, 2, 3, 4, 5]
let squares = ls.map(|x| => x * x)
let addedUp = ls.reduce(0, |acc, x| => acc + x)
let even = ls.filter(|x| => (x % 2 == 0))
println(squares)
println(addedUp)
println(even)
```

### Named Functions

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
provides a feature called 'Named Functions' that offers this functionality.  Named functions 
are declared at the beginning of a given scope.  For example:

```
fn a() {
    return b()
}
fn b() {
    return 42
}
assert(a() == 42)
```

This example works because both `a` and `b` are declared "simultaneously" at the beginning
of the program, before any other declarations are processed by the compiler.  Note 
that the above program is identical in every way to the previous one, except for 
the forward references provided by the named function syntax.

Named functions are often times just easier to read as well.  It is considering idiomatic
in Golem to use named functions even when it is not strictly necessary.

### Closures

Golem supports [closures](https://en.wikipedia.org/wiki/Closure_(computer_programming)) as 
well.   

Here is an example of a closure that acts as a
[accumulator generator](http://www.paulgraham.com/accgen.html):

```
fn foo(n) {
    return fn(i) {
        return n += i
    }; 
}
let f = foo(4)
println([f(1), f(2), f(3)])
```

Closures are a fundamental mechanism in Golem for managing state.  We will have more 
to say about closures [later on](#merging-structs) in the tour.

### Optional Parameters

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

### Variadic Functions

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

### Arity

There is a builtin function called [arity](builtins.html#arity) that returns 
a [struct](#structs) that describes the [arity](https://en.wikipedia.org/wiki/Arity) 
of a function.  Here is a program that prints the arity of 3 of the builtin functions 
we have already used:

```
println(arity(len))
println(arity(range))
println(arity(println))
```

## Structs

Golem is not an object-oriented language.  It does not have classes, objects, 
inheritance, or constructors.  What it does have, however, are values 
which we call "structs".

### Struct Syntax

Structs are created via the `struct` keyword.  

```
let s = struct { a: 1, b: 2 }
println(s)
```

In the above example, we've created a struct that has two fields, `a` and `b`.  Remember that 
a "field" in Golem is a named member of a value.  Structs are values that have
an arbitrary collection of fields.

Structs are similar to dicts in some ways, but quite different in others.  The field 
names of a struct can only be strings, and furthermore they
must be valid identifiers -- they cannot have spaces or special characters.

The dot operator, `.`, is used on structs to get or set the fields of a struct:

```
let s = struct { a: 1, b: 2 }
println(s)
s.a = 3
println(s)
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
```

### Properties

Structs can have properties defined on them, so that a given field has a 'getter' 
function, and optional 'setter' function.  The getter function must take 0 parameters, 
and the setter function must take 1 parameter.  If the setter function is omitted,
the property is readonly.  Properties are useful for hiding the inner workings
of a struct behind a simpler facade.  Here is an example (which uses 
[try-catch](#errors)):

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

println([s.a, s.b, x, s.c()])
s.b = 3
println([s.a, s.b, x, s.c()])
x = 4
println([s.a, s.b, x, s.c()])
```

### Merging Structs

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
all three structs actually share a common set of fields.  We will see in the next section
that this behaviour can be quite useful.

### Using Structs to build complex values

By using structs, closures, and `merge()` together, it is possible to simulate various 
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

The functions `newRectangle` and `newBox` are very much like what one might call "constructors"
in another language.  The structs that they return have functions as entries 
(e.g. `area()`), and these functions refer to captured variables (`w`, `h`, and `d`) 
that are somewhat like member variables of a class.  As such, the functions are quite 
a bit like what one might call a "method" in another language.

The use of the `merge()` function to create a box out of a rectangle is similar to
how inheritance is used in other languages.  Does that mean that a Box is a subclass
of a Rectangle?  Not really, no.  There is no such thing as a "class" in Golem.  However, 
due to the behaviour of merge(), they *are* inter-related in a way that is 
very much like inheritance.

One of the primary goals of the Golem project is to explore the power provided by 
the simple building blocks of functions, closures, structs and merge().  It is hoped
that the simplicity and flexibility of these elements can be used to create a variety
of complex runtime structures that are easy to reason about and use.  

## Errors

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

Try statements can be nested, and errors that are not caught in a function are 
passed upwards in the call stack.

## Concurrency

Golem uses the Go Language's [concurrency system](https://tour.golang.org/concurrency/1).  This 
means that Golem has 'goroutines', channels and the ability to send and 
receive messages.  

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

Golem's concurrency is not finished yet.  In the near future it will be enhanced
with the `select` keyword, the ability to range over a channel, and various pieces of
functionality from Go's `sync` package.

## Immutability

Golem supports immutability via the `freeze()` builtin function, which makes a mutable
value become immutable.  You can check if a value is immutable via the `frozen()`
builtin function. `freeze()` always returns the value that you pass into it.

```
let s = freeze(struct { a: 1, b: 2 })
println(frozen(s))

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

## Introspection

## Embedding

## Command Line Interface

### Modules

private 

### The `main()` Function
### Standard Library


