let s = '';
for i in range(0, 4) {
    switch i {
    case 0, 1:
        s += 'a';

    case 2:
        s += 'b';
    }
}
assert(s == 'aab');
# Golem Tutorial

Golem is a general purpose, interpreted language, with first-class functions and a 
dynamic type system.  Golem aims to combine the clean semantics of Python, 
the concurrency of Go, the flexibility of Javascript, and the embeddability of Lua.

Golem doesn't yet have a [REPL](https://en.wikipedia.org/wiki/Read%E2%80%93eval%E2%80%93print_loop)
environment.  So, follow along with this tutorial by typing code into a 
text file ('tutorial.glm', for example), and then running the file
from the command line to look at the results.

First, build a golem executable via:

```
go build golem.go
```

Then, fire up your text editor of choice, and type the following program into 
a file named 'tutorial.glm':

```
println('Hello, world.');
```

and then run it like so:

```
./golem tutorial.glm
```

The function `println` is built in to Golem.  There are several of these builtin functions
in Golem.

## Basic Types

Golem has a simple, straightforward type system.  The basic primitive types 
include boolean, string, int and float.  There is also 'null', which 
represents the absence of a value.  Basic values are immutable.

Golem has the usual set of c-syntax-family operators that you would 
expect: `==`, `!=`, `||`, `&&`, `<`, `>`, `+`, `-`, and so forth.  

```
assert(1 + 2 == 3);
assert(42 / 7 == 8 - 2);
```

We will cover the operators in more detail later.  Note that we used another builtin 
function, `assert`, which will throw an exception if the value that is passed into 
it is not true.

Integer values in Golem are signed 64 bit integers.  Float values are 64-bit.  Ints 
are coerced to Floats during arithmetic and checks for equality:

```
assert(12 / 4.0 == 3.0);
assert(12 / 4.0 == 3);
```

Another builtin function, `str`, returns  the string representation of a value:

```
assert(str(3) == '3');
```

Strings can be delimited either with a single quote or a double quote:

```
assert('abc\n' == "abc\n");
```

During addition, if one of the values is a string, and the other is not, then
the other value is converted to a string, and the two strings are then 
concatenated together:

```
assert('a' + 1 == 'a1');
```

Unlike many other dynamic languages, Golem has no concept of 'truthiness'.  The only 
things that are true or false are boolean values:

```
assert(true);
assert(!false);
```

So, the empty string, zero, null, etc. are *not* boolean, and will throw a
TypeMismatch error if you attempt to evaluate them in a place where a boolean value is
expected.

```
assert('');
```

## Comments

Golem uses C-family comments:  `/* ... */` for a block comment, and `//` for 
a line comment.

## Variables

Values can be assigned to variables. Variables are declared via either the `let` 
or `const` keyword.  It is an error to refer to a variable before it has been
declared.

```
let a = 1;
const b = 2;
a = b + 3;
println(a);
```

`let` and `const` are statements -- they do not return a value.  Assignments, on the
other, *are* expressions:

```
let a = 1;
let b = (a = 2);
assert(a == b && b == 2);
```

## Collections

Golem has four collection data types: List, Dict, Set, and Tuple.

You can create a list by enclosing a comma-delimited sequence of values in
square brackets.  Once you've created a list, you can use square brackets to access
individual elements of a list (this is called the 'index operator').  

```
let a = [];
let b = [3,4,5];
assert(a.isEmpty());
assert(b[0] == 3);
```

If the index value is negative, values will be indexed from the beginning of the list.
This is a really handy way to get the last element of a list.
```
let a = [1,2,3];
assert(a[-1] == 3);
```

Use the 'slice operator' to create a new list from part of an existing list or string.
If you leave off the first or last value of the slice operation, the resulting slice
will start at the beginning or end.  Negative values work with slices in the same 
way that they do with lists.

```
let c = [4,5,6,7,8];
assert(c[1:3] == [5,6]);
assert(c[:3] == [4,5,6]);
assert(c[2:] == [6,7,8]);
assert(c[1:-1] == [5,6,7]);
```

Indexing and slicing works on strings too:

```
assert('abc'[1] == 'b');
assert('abc'[:-1] == 'ab');
```

Golem's `dict` type is similar to Python's 'dict', or 'HashMap' in java.  The
keys can be any value that supports hashing (currently str, int, float, bool, or tuple). 
A future version of Golem will probably allow for more types to act as a dict key.

```
let a = dict {'x': 1, 'y': 2};
assert(a['x'] == 1);
```

A `set` is a collection of distinct values.  Any value that can act as a key in a dict
can be a member of a set.

```
let a = set {'x', 'y'};
assert(a.contains('x'));
```

A `tuple` is similar to a Python tuple.  It is an immutable list-like data structure.
Tuples must have at least two values.

```
let a = (1, 2);
assert(a[0] == 1);
```

The builtin function `len` can be used to get the length of any of the collections.
`len` will also return the length of a string.

```
let a = [1, 2, 3];
let b = 'lmnop';
let c = dict {"x": 3};

assert([len(a), len(b), len(c)] == [3,5,1]);
```

## Intrinsic Functions

In Golem, an 'intrinsic function' is a function that is intrinsically present on a value.
Intrinsic functions are similar to what are called 'methods' in other languages, except 
that intrinsic functions are not user-definable -- they are baked into the language 
itself.

Here's an example of how to use some of the intrinsic functions that a list has:

```
let ls = [];
assert(ls.isEmpty());
ls.add('a');
ls.addAll(['b', 'c']);
assert(ls == ['a', 'b', 'c']);
assert(ls.contains('c'));
assert(ls.indexOf('b') == 1);
assert(ls.join(',') == 'a,b,c');

ls = [1, 2, 3, 4, 5];
let squares = ls.map(x => x * x);
let addedUp = ls.reduce(0, |acc, x| => acc + x);
let even = ls.filter(x => (x % 2 == 0));

assert(squares == [1, 4, 9, 16, 25]);
assert(addedUp == 15);
assert(even == [2, 4]);

```

Note that we jumped ahead a bit with the `map`, `reduce` and `filter` examples,
because we are using 'lambdas' there to define functions on the fly.  We will 
discuss functions later on in the tutorial.

See the [Golem Language Reference](https://github.com/mjarmy/golem/blob/master/docs/reference.md) 
for a complete description of all the intrinsic functions on the various types.

## Control Structures

Golem has a familiar set of control structures: `if`, `while`, `switch`, and `for`.

```
let a = 1;
while a < 12 {
    if a < 3 {
        a = a + 2;
    } else {
        a = 15;
    }
}
assert(a == 15);
```

Golem also has `break` and `continue`, which will break out of a `while` or `for` loop,
or continue at the top of the loop, as in other languages.

Golem has 'ternary-if' expressions as well:

```
const a = 10;
let b = a < 3 ? 4 : 5;
assert(b == 5);
```

`switch` works roughly the same way as it does in other languages, except that you 
can switch on any value, not just integers.  Also, there is no 'fall-thru' -- at most
only one case will be executed.  Therefore the `break` keyword is not applicable to 
switches.

```
let a = 'abc';
let b = 0;
switch a {
    case 0:
        b = 1;
    case 'abc':
        b = 2;
    default:
        b = 3;
}
assert(b == 2);
```

You can have multiple expressions in a case statement.  The body of the case
will be executed if at least on of the expressions matches:

```
let s = '';
let i = 0;
while i < 4 {
    switch i {
    case 0, 1:
        s += 'a';
    case 2:
        s += 'b';
    }
    i++;
}
assert(s == 'aab');
```

You can leave the expression out after the `switch` keyword. This lets you switch
on a sequence of boolean case statements, which is sometimes easier to read than
a cascade of 'if, else-if, else-if' statmements:

```
let b = 0;
switch {
    case 1 < 2:
        b = 1;
    default:
        b = 2;
}
assert(b == 1);
```

Golem's `for` statement is much like python's `for` statement (and therefore unlike
`for` in C-family langauges like Java, JS, C#, etc).

```
let a = [1, 2, 3];
let z = 0;
for e in a {
    z += e;
}
assert(z == 6);
```

`for` iterates over a sequence of values derived from an 'iterable' value.  Lists, dicts,
sets, and strings are iterable.

Use a tuple to capture the values iterated from a dict:

```
let d = dict { "x": 1, "y": 2, "z": 3 };
for (k, v) in d {
    println("key: ", k, ", value: ", v);
}
```

There is actually one more iterable type -- ranges.  Ranges are created via the `range`
builtin function.  A range is an immutable value that represents a sequence of integers.  

```
let list = ["frog", "cow", "rabbit"];
for i in range(0, len(list)) {
    if list[i] == "cow" {
        println("The cow is at element ", i);
        break;
    }
}
```

Note that ranges do not actually contain a list of all the specified integers.  
They simply represent a sequence that can be iterated over.

## Operators and Expressions

Golem has the following operators, with the following precedence (from low to high):

| Category       | Operators     |
| -------------  | ------------- |
| or             | <code>&#124;&#124;</code>  |
| and            | `&&`  |
| comparative    | `==`, `!=`, `>`, `>=`, `<`, `<=`, `<=>`, `has` |
| additive       | `+`, `-`, <code>&#124;</code>, `^` |
| multiplicative | `*`, `/`, `%`, `&`, `<<`, `>>` |
| unary          | `-`, `!`, `~` |
| postfix        | `++`, `--`  |

The 'spaceship' operator, `<=>`, returns -1, 0, or 1 if the left-hand operator is 
less than, equal to, or greater than the right-hand operator: 

```
assert((5 <=> 10) == -1);
```

We will discuss the `has` operator later on when we talk about structs.

Note that Golem does not have any prefix operators, so `++` and `--` are always postfix.

Golem also supports 'assignment operators`, which perform an operation and
do an assignment at the same time, e.g.:

```
a += b; // is the same as a = a + b;
```

Here are the assignment operators:

`=+`, `=-`, `=*`, `=/`, `=%`, `=^`, `=&`, `=|`, `=<<`, `=>>`

## Functions and Closures

Functions are first class values in Golem.  They are created with the 'fn' keyword, 
and they are invoked by adding parameters in parentheses to the end of an expression that 
evaluates to a function:

```
let a = fn(x) {
    return x * 7;
};
assert(a(6) == 42);
```

Functions do not have to have an explicit `return` statement. If there is no `return`,
they will return the last expression that was evaluated.  If no expression is 
evaluated, `null` is returned.

```
let a = fn() {};
let b = fn(x) { x * x; };
assert(a() == null);
assert(b(3) == 9);
```

Golem supports closures as well -- in fact closures are a fundamental mechanism
in Golem for managing state.  Here is an example of a closure that acts as a
[accumulator generator](http://www.paulgraham.com/accgen.html):

```
let foo = fn(n) {
    return fn(i) {
        return n += i;
    }; 
};
let f = foo(4);
assert([f(1), f(2), f(3)] == [5, 7, 10]);
```

You can declare the formal parameters of a function to be constant.  In the following
example, the formal parameter 'b' is constant, so it cannot be changed inside the 
function:

```
let a = 1;

fn foo(const b) {
    return a += b;
}

foo(2);
foo(3);
assert(a == 6);
```

Golem also supports 'lambda' syntax, via the `=>` operator.  Lambdas provide a 
lightweight way to define a function on the fly. The body of a lambda function is a 
single expression. A lambda that takes only one parameter can omit the surrounding pipes.

```
let a = || => 3;
let b = x => x * x;
let c = |x, y| => (x + y)*5;

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

```
fn a() {
    return b();
}
fn b() {
    return 42;
}
assert(a() == 42);
```

**TODO** arity()
**TODO** optional param values, variadic functions

## Passing In Arguments via main()

You can pass arguments into a Golem program by defining a `main()` function, that
accepts exactly one parameter.  The parameter will always be a list of the
command line arguments.

```
fn main(args) {
    for i in range(0, len(args)) {
        println('argument ', i, ' is "', args[i], '"');
    }
}
```

## Structs

Golem is not an object-oriented language.  It does not have classes, objects, 
inheritance, or constructors.  What it does have, however, are values 
which we call 'structs'.

Structs are created via the `struct` keyword.  
```
let s = struct { a: 1, b: 2 };
```

Structs are similar to dicts in some ways, but quite different in others.  The keys of
a struct can only be strings ('a' and 'b' in the above example), and furthermore they
must be valid identifiers -- they cannot have spaces or special characters.

The dot operator, `.`, is used on structs to get or set the value associated with a key.  You
can use the `has` operator to test whether a struct contains a given key.  

```
let s = struct { a: 1, b: 2 };
assert(s has 'a');
assert(s.a == 1);
```

There are three builtin functions that you can use to inspect and modify a struct: 
`fields()`, `getval()`, and `setval()`.

```
let s = struct { a: 1, b: 2 };
assert(fields(s) == set { 'a', 'b' });
assert(getval(s, 'a') == 1);
assert(setval(s, 'a', 3) == 3);
assert(getval(s, 'a') == 3);
```

Onnce a struct is created, it cannot have new keys added
to it, or existing keys removed.  The _values_ associated with the keys can be changed
though, as we saw in the previous example.

The `this` keyword is used in Golem to allow a struct to refer to itself. In Golem,
`this` is only valid inside a struct, and it is always lexically scoped to refer to 
the innermost enclosing struct. (Note: this is **very** different than the semantics 
of `this` in javascript.)

```
let s = struct { a: 1, b: 2, c: this.a + this.b };
println(s);
assert(s.c == 3);
```

The builtin-function `merge()` can be used to combine an arbitrary number of 
existing structs into a new struct.

```
let a = struct { x: 1, y: 2};
let b = struct { y: 3, z: 4};
let c = merge(a, b);
assert(c.x == 1);
assert(c.y == 2);
assert(c.z == 4);
a.x = 10;
assert(c.x == 10); // x is changed here too!
```

If there are any duplicated keys in the structs passed in to 'merge()', then the
value associated with the first such key is used.  

Also, note in the above example that if you change a value in one of the structs passed 
in to merge(), the value changes in the merged struct as well.  That is because the 
structs 'a', 'b', and 'c' actually share a common set of entries.  We will see in the 
next section that this behaviour has some quite useful ramifications.

## Combining Structs Together

The combination of closures, structs and merge() is very powerful.  With these tools, it
is possible to simulate various features from other languages, including 
inheritance, multiple-inheritance, prototype chains, and the like.

For instance, consider the following program:

```
fn Rectangle(w, h) {
    return struct {
        width: w,
        height: h,
        area: fn() { return this.width * this.height; }
    };
}

fn Box(rect, d) {
    return merge(
        rect, 
        struct {
            depth: d,
            volume: fn() { return rect.area() * this.depth; }
        });
}

let r = Rectangle(2, 3);

let b = Box(r, 4);
assert([b.width, b.height, b.depth, b.area(), b.volume()] == [2, 3, 4, 6, 24]);

r.width = 5;
assert([b.width, b.height, b.depth, b.area(), b.volume()] == [5, 3, 4, 15, 60]);

```

The functions 'Rectangle' and 'Box' are very much like what one might call 'constructors'
in another language.  The structs that they return have functions as entries 
(e.g. 'area()'), and these functions refer to the 'this' keyword, and to captured 
variables.  As such, the functions are an awful lot like what one might call a 
'method' in another language.

The use of the 'merge()' function to create a box out of a rectangle is similar to
how inheritance is used in other languages.  Does that mean that a Box is a subclass
of a Rectangle?  Not really, no.  There is no such thing as a 'class' in Golem.  However, 
due to the behaviour of merge(), they *are* inter-related in a way that is 
very much like inheritance.

By the way, note that the functions 'Rectangle' and 'Box' are capitalized.  It is 
considered idiomatic in Golem to capitalize 'constructor-like' functions that 
return complicated structs which have things like closures, 'this' references, 
merges from other complicated structs, and the like.

One of the primary goals of the Golem project is to explore the power provided by 
the simple building blocks of functions, closures, structs and merge().  It is hoped
that the simplicity and flexibility of these elements can be used to create a variety
of complex runtime structures that are easy to reason about and use.  

## Error Handling

Golem uses the familiar 'try-catch-finally` syntax that exists in many C-family 
languages.

```
try {
    let z = 4 / 0;
}
catch e {
    println(e);
    assert(e.kind == 'DivideByZero');
}
```

Exceptions are structs.  You can throw an exception by using a struct literal with the
`throw` keyword.

```
try {
    throw struct { msg: 'foo' };
    assert(false); // can't get here
}
catch e {
    println(e);
    assert(e.msg == 'foo');
}
```

## Concurrency

**TODO** explain this a bunch more.

Golem uses the Go Language's [concurrency system](https://tour.golang.org/concurrency/1).  This 
means that Golem has 'goroutines', channels and the ability to send and 
receive messages.  Go's concurrency capabilities are quite powerful, and Golem is capable
of using them all (**TODO** well, just not yet.  We need to figure out how to provide
access to mutexes and lots of other lower level things)

```
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

go sum(a[:n], c);
go sum(a[n:], c);
let x = c.recv();
let y = c.recv();
assert([x, y] == [-5, 17]);
```

## Immutabilty

Golem supports immutability via the `freeze()` builtin function, which makes a mutable
value become immutable.  You can check if a value is immutable via the `frozen()`
builtin function. `freeze()` always returns the value that you pass into it.

```
let s = freeze(struct { a: 1, b: 2 });
assert(frozen(s));

try {
    s.a = 0;       // This will throw an ImmutableValue error.
    assert(false); // We can't reach this statement.
} catch e {
    assert(e.kind == 'ImmutableValue');
}
```

`freeze()` only has an affect on Lists, Dicts, Sets and Structs.  All other values 
are already immutable, so calling `freeze()` on them has no effect

An import caveat regarding immutability is that even though closures, like all functions, 
are immutable, they can still have enclosed state that can be modified.  There
is no way in Golem to freeze a closure after the fact sot that it can no longer modify 
any of its captured variables.  It is up to you to manage state properly if you are 
using closures.  Here is the "accumulator generator" from a previous example.  
We freeze it this time, but it still has mutable state via the enclosed variable 'n':

```
let foo = fn(n) { return fn(i) { return n += i; }; };
freeze(foo);
let f = foo(4);
assert([f(1), f(2), f(3)] == [5, 7, 10]);
```

Using the `const` keyword, in conjuction with `freeze()`, is a good way to lock 
down your code so that inadvertent immutablity does not creep in.

Immutabilty and concurrency go hand in hand.  By using immutable values whenever 
possible, you can greatly reduce the likelyhood of bugs in your concurrency code, 
and make it much easier to reason about as well.
    
## Type Introspection

There is a builtin function called `type()` that will return a string 
describing the type of a value.  Here is a program that will print a list 
of every possible type:

```
let values = [
    true, "", 0, 0.0, fn(){}, 
    [], range(0,1), (0,1), dict{}, set{}, 
    struct{}, chan()];

let types = values.map(type);

println(types);
```

## Standard Library

**TODO** io, sys, regex, net, http, time, sql, json

## A Full-Fledged Program

As a final example, here is a  program that finds strings inside a text file or
files.  This program can be found in the 'examples' directory of the github repo.

```
import io;
import regex;
import sys;

/*
 * To run this program, do the following from the parent directory:
 * 
 *     go build .
 *     ./golem examples/searchFiles.glm let examples
 *
 * This will find every occurence of the word 'let' in all of the files 
 * in the 'examples' directory.
 * 
 * Note that this program doesn't yet understand file globbing, so for now you have to 
 * provide an explicit name for the file or directory that you want to search.
 */

fn traverse(pattern, file) {
    if file.isDir() {
        for child in file.items() {
            traverse(pattern, child);
        }
    } else {
        let lines = file.readLines();
        for i in range(0, len(lines)) {
            if pattern.match(lines[i]) {
                println([file.name, i, lines[i]].join(':'));
            } 
        }
    }
}

fn main(args) {

    if len(args) != 2 {
        println("Expected 2 arguments, got ", len(args));
        sys.exit(-1);
    }

    let pattern = regex.compile(args[0]);
    let file = io.File(args[1]);

    traverse(pattern, file);
}

```

