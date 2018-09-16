# Builtin Functions

* [Standard Builtins](#standard-builtins)
* [Unsandboxed Builtins](#unsandboxed-builtins)


## Standard Builtins

Golem has a collection of standard builtin functions that provide
various kinds of important functionality. All of the standard builtins
are "pure" functions that do not do any form of I/O. As such they
are suitable for use in sandboxed environments.

* [`arity()`](#arity)
* [`assert()`](#assert)
* [`chan()`](#chan)
* [`fields()`](#fields)
* [`freeze()`](#freeze)
* [`frozen()`](#frozen)
* [`getField()`](#getField)
* [`hasField()`](#hasField)
* [`iter()`](#iter)
* [`len()`](#len)
* [`merge()`](#merge)
* [`range()`](#range)
* [`setField()`](#setField)
* [`str()`](#str)
* [`type()`](#type)

#### `arity`

`arity` returns a Struct describing the [arity](https://en.wikipedia.org/wiki/Arity) of a Func.
A func's arity type is always either "Fixed", "Variadic", or "Multiple".

* signature: `arity(f <Func>) <Struct>`
* example:

```
    assert(arity(println) == struct { kind: "Variadic", required: 0 })
    assert(arity(len)     == struct { kind: "Fixed",    required: 1 })
    assert(arity(range)   == struct { kind: "Multiple", required: 2, optional: 1 })
```
#### `assert`

`assert` accepts a single boolean value, and throws an error
if the value is not equal to `true`.  `assert` returns `true`
if it does not throw an error.

* signature: `assert(b <Bool>) <Bool>`
* example: `assert(0 < 1)`

#### `chan`

`chan` creates a channel of values.  `chan` has a single optional size parameter that
defaults to 0.  If size is 0, an unbuffered channel will be created.
If the size is greater than 0, then a buffered channel of that size will be created.

* signature: `chan(size = 0 <Int>) <Chan>`
* example: `let ch = chan()`

#### `fields`

`fields` returns a Set of the names of a value's fields.

* signature: `fields(value <Value>) <Set>`
* example:

```
    println(fields([]))
    println(fields(struct { a: 1, b: 2}))
```

#### `freeze`

`freeze` freezes a value, if it is not already frozen.  Its OK to call `freeze`
on values that are already frozen.  The value is returned after it is frozen.

* signature: `freeze(value <Value>) <Freeze>`
* example: `freeze([1, 2])`

#### `frozen`

`frozen` returns whether or not a value is frozen.

* signature: `frozen(value <Value>) <Bool>`
* example:

```
    println(frozen('a'))
    println(frozen([3, 4]))
```

#### `getField`

`getField` returns the value associated with a field name.

* signature: `getField(value <Value>, name <Str>) <Value>`
* example:

```
    let a = [1, 2]
    let f = getField(a, 'add')
    f(3)
    println(a)
```

#### `hasField`

`hasField` returns whether a value has a field with a given name.

* signature: `getField(value <Value>, name <Str>) <Bool>`
* example:

```
    let a = [1, 2]
    println(hasField(a, 'add'))
```

#### `iter`

`iter` returns an iterator for an Iterable value.  Str, List, Range, Dict,
and Set are iterable.

An iterator is a Struct that has two fields:

* A `next()` function that
returns whether there are any more values in the iterator,
and advances the iterator forwards if there is another value.

* A `get()` function that returns the currently available value.

By convention, a new iterator has to have `next()` called on it to advance
to the first available value. Calling `get()` before the first call to `next()`
throws an error.

* signature: `iter(value <Iterable>) <Struct>`
* `next` signature: `next() <Bool>`
* `get` signature: `get() <Value>`
* example:

```
    let a = [1, 2, 3]
    let itr = iter(a)
    while itr.next() {
        println(itr.get())
    }
```

#### `len`

`len` returns the length of a value that has a length.  Str, List, Tuple, Range, Dict,
and Set have a length

* signature: `len(value <Lenable>) <Int>`
* example: `println(len('abc'))`

#### `merge`

`merge` merges structs together into a new struct.  Consult the [tour](#TODO)
for a detailed description of how `merge` works.

* signature: `merge(structs... <Struct>) <Struct>`

#### `range`

`range` creates a Range, starting at "from" (inclusive) and going until
"to" (exclusive).

The optional "step" parameter, which defaults to 1,
specifies the distance between succesive integers in the range.  You can
create a "backwards" range by specify a negative step value, and a "from"
that is less than "to".

* signature: `range(from <Int>, to <Int>, step = 1 <Int>) <Range>`
* example:

```
    for i in range(0, 5) {
        println(i)
    }
```
#### `setField`

`setField` sets the value associated with a field name. `setField` only works
on Structs -- you cannot set the fields of other types. `setField` returns `null`
if it was successful.

* signature: `setField(s <Struct>, name <Str>, value <Value>) <Null>`
* example:

```
    let s = struct { a: 1, b: 2 }
    setField(s, 'a', 3)
    println(s)
```

#### `str`

`str` returns a Str representation of a value.

* signature: `str(value <Value>) <Str>`
* example: `println(str([null, true, 1, 'abc']))`

#### `type`

`type` returns the type of a value.

* signature: `type(value <Value>) <Str>`
* example:

```
    println(type(1.23))
    let a = [null, true, 1, 'xyz']
    println(a.map(type))
```


## Unsandboxed Builtins

Golem also has a couple of "unsandboxed" builtins.  These functions
perform I/O, so they should not be included in sandboxed Golem
environments.

* [`print()`](#print)
* [`println()`](#println)

#### `print`

`print` prints a sequence of values to STDOUT.

* signature: `print(values... <Value>) <Null>`

#### `println`

`println` prints a sequence of values to STDOUT, followed by a "\n".

* signature: `println(values... <Value>) <Null>`


_This document uses documentation from [go](https://github.com/golang/go), which
is licensed under the BSD-3-Clause license._
