# The Golem Programming Language Tour

Welcome to the tour of the Golem Programming Language.

## Hello, world

Let's get started with the proverbial hello world program.  In this tour, we will 
be running Golem code directly in the browser, so press the `Run` button below
to see the output of the program:

```
println('Hello, world.');
```

You may have noticed that there is a semicolon at the end of the `println` 
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

The builtin function `len` can be used to get the length of any of the collections.
`len` will also return the length of a string.

```
let a = [1, 2, 3]
let b = 'lmnop'
let c = dict {"x": 3}

println([len(a), len(b), len(c)])
```

