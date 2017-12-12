# Golem Tutorial

**Golem is a fledgling project right now.  Its not yet ready for production use.**  

Golem is a general purpose, interpreted language, with first-class functions and a 
dynamic type system.  Golem aims to combine the clean semantics of Python, 
the concurrency of Go, the flexibility of Javascript, and the embeddability of Lua.

Golem doesn't have a [REPL](https://en.wikipedia.org/wiki/Read%E2%80%93eval%E2%80%93print_loop)
environment yet (**TODO**).  So, follow along with this tutorial by typing things into a 
source code file ('tutorial.glm', for example), and then running `./golem tutorial.glm`
from the command line to look at the results.

First, you will need to build golem itself, via:
```
go build .
```
That will create a golem executable for you.

Next, add the following to your 'tutorial.glm' file, and run it:

```golem
println('Hello, world.');
```

The function `println` is built in to Golem.  There are several of these builtin functions
in Golem.

## Basic Types

Golem has a simple, straightforward type system.  The basic primitive types 
include boolean, string, int and float.  There is also 'null', which 
represents the absence of a value.  Basic values are immutable.

Golem has the usual set of c-syntax-family operators that you would 
expect: ==, !=, ||, &&, <, >, <=, >=, +, -, and so forth.  

```golem
assert(1 + 2 == 3);
assert(42 / 7 == 8 - 2);
```

We will cover the operators in more detail later.  Note that we used another builtin 
function, `assert`, which will throw an exception if the value that is passed into 
it is not true.

Integer values in Golem are signed 64 bit integers.  Float values are 64-bit.  Ints 
are coerced to Floats during arithmetic and checks for equality:

```golem
assert(12 / 4.0 == 3.0);
assert(12 / 4.0 == 3);
```

Another builtin function, `str`, returns  the string representation of a value:

```golem
assert(str(3) == '3');
```

Strings can be delimited either with a single quote or a double quote:

```golem
assert('abc\n' == "abc\n");
```

During addition, if one of the values is a string, and the other is not, then
the other value is converted to a string, and the two strings are then 
concatenated together:

```golem
assert('a' + 1 == 'a1');
```

Unlike many other dynamic languages, Golem has no concept of 'truthiness'.  The only 
things that are true or false are boolean values:

```golem
assert(true);
assert(!false);
```

So, the empty string, zero, null, etc. are *not* boolean, and will throw an
exception if you attempt to evaluate them in a place where a boolean value is
expected.

**TODO** intrinsic functions on basic types.


## Variables

Values can be assigned to variables. Variables are declared via either the `let` 
or `const` keyword.  It is an error to refer to a variable before it has been
declared.

```golem
let a = 1;
const b = 2;
a = b + 3;
println(a);
```

`let` and `const` are statements -- they do not return a value.  Assignments, on the
other, *are* expressions in golem:

```golem
let a = 1;
let b = (a = 2);
assert(a == b && b == 2);
```

## Collections

Golem has three collection data types: List, Dict, and Set.

Use square brackets to create a list:

```golem
let a = [];
let b = [3,4,5];
assert(a.isEmpty());
assert(b[0] == 3);
```

Use the 'slice' operator to create a new list from part of an existing list:

```golem
let c = [4,5,6,7,8];
assert(c[1:3] == [5,6]);
assert(c[:3] == [4,5,6]);
assert(c[2:] == [6,7,8]);
```

Golem's `dict` type is similar to Python's 'dict', or 'HashMap' in java.  The
keys can be any value that supports hashing (currently str, int, float, or bool). 
**TODO** A future version of Golem will allow for structs to act as a dict key.

```golem
let a = dict {'x': 1, 'y': 2};
assert(a['x'] == 1);
```

A `set` is a collection of distinct values.  Any value that can act as a key in a dict
can be a member of a set.

```golem
let a = set {'x', 'y'};
assert(a.contains('x'));
```

**TODO** assignments for list and dict

## Control Structures

**TODO** if, while, for, switch, ternary if

## Error Handling

**TODO** try, catch, finally, throw

## Operators and Expressions

**TODO** assigment, increment, decrement, ternary if

**TODO** precedence

## Functions and Closures

Functions are first class values in Golem.  They can be instantiated and passed 
around. They are created with the 'fn' keyword, and they are invoked the usual way,
by adding parameters in parentheses to the end of an expression that 
evaluates to a function:

```golem
let a = fn(x) {
    return x * 7;
};
assert(a(6) == 42);
```

Functions do not have to have an explicit `return` statement. If there is no `return`,
they will return the last expression that was evaluated.  If no expression is 
evaluated, `null` is returned.

```golem
let a = fn() {};
let b = fn(x) { x * x; };
assert(a() == null);
assert(b(3) == 9);
```

Golem supports closures as well -- in fact closures are a fundamental mechanism
in Golem for managing state.  If you are not familiar with closures, you should
definitely [read up on them](https://en.wikipedia.org/wiki/Closure_(computer_programming)). 
Here is an example of a closure that acts as a
[accumulator generator](http://www.paulgraham.com/accgen.html):

```golem
let foo = fn(n) {
    return fn(i) {
        return n += i;
    }; 
};
let f = foo(4);
assert([f(1), f(2), f(3)] == [5, 7, 10]);
```

Golem also supports 'lambda' syntax.  Lambdas provide a lightweight way
to define a function on the fly. The body of a lambda function is a single 
expression. A lambda that takes only one parameter can omit the surrounding pipes.
**TODO** provide easier to read examples

```golem
let z = 5;
let a = || => 3;
let b = x => x * x;
let c = |x, y| => (x + y)*z;
assert(a() == 3);
assert(b(2) == 4);
assert(c(1, 2) == 15);
```

'Named functions' in Golem are functions that are declared at the beginning of
a given scope, before any other declarations are processed by the compiler.  Using 
named function syntax allows for forward references -- you 
can refer to functions that have not been defined yet.

Note that named functions do not have a semicolon at the end of the closing 
curly brace.

```golem
fn a() {
    return b();
}
fn b() {
    return 42;
}
assert(a() == 42);
```

**TODO** optional param values, variadic functions

## Structs

Golem is *not* an object-oriented language.  It does not have classes, objects, or 
inheritance.  What it does have, though, are values which we call 'structs'.  A
struct is similar in spirit to what one might called a 
'[duck-typed](https://en.wikipedia.org/wiki/Duck_typing) object'.

`this`

merge()

## Putting it All Together

The combination of closures, structs and merge() is very powerful.  Show 
some examples.
    
## Reflection

**TODO** typeof(), meta()

## Immutability

**TODO** freeze()

## Concurrency

**TODO** explain spawn, chan(), ch.send(), ch.recv()

```golem
fn sum(a, c) {
    let total = 0;
    for v in a {
        total += v;
    }
    c.send(total);
}

let a = [7, 2, 8, -9, 4, 0];
let n = len(a) / 2;
let c = chan();

spawn sum(a[:n], c);
spawn sum(a[n:], c);
let x = c.recv();
let y = c.recv();
assert([x, y] == [-5, 17]);
```

## Standard Library

**TODO** io, net, http, time, sql, json

