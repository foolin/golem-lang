
## Sandbox Builtins

Golem has a collection of standard builtin functions that provide
various kinds of important functionality. All of the sandbox builtins
are "pure" functions that do not do any form of I/O. As such they
are suitable for use in sandboxed environments.

* [`arity()`](#arity)
* [`assert()`](#assert)
* [`chan()`](#chan)
* [`fields()`](#fields)
* [`freeze()`](#freeze)
* [`frozen()`](#frozen)
* [`has()`](#has)
* [`hashCode()`](#hashcode)
* [`iter()`](#iter)
* [`len()`](#len)
* [`merge()`](#merge)
* [`range()`](#range)
* [`stream()`](#stream)
* [`str()`](#str)
* [`type()`](#type)

### `arity`

`arity` returns a [Struct](struct.html) describing the [arity](https://en.wikipedia.org/wiki/Arity) of a Func.
A func's arity type is always either "Fixed", "Variadic", or "Multiple".

* signature: `arity(f <Func>) <Struct>`
* example:

```
println(arity(len))
println(arity(println))
println(arity(range))
```

### `assert`

`assert` accepts a single boolean value, and throws an error
if the value is not equal to `true`.  `assert` returns `true`
if it does not throw an error.

* signature: `assert(b <Bool>) <Bool>`
* example: `assert(0 < 1)`

### `chan`

`chan` creates a [channel](chan.html) of values.  `chan` has a single optional size parameter that
defaults to 0.  If size is 0, an unbuffered channel will be created.
If the size is greater than 0, then a buffered channel of that size will be created.

* signature: `chan(size = 0 <Int>) <Chan>`
* example: `let ch = chan()`

### `fields`

`fields` returns a [Set](set.html) of the names of a value's fields.

* signature: `fields(value <Value>) <Set>`
* example:

```
println(fields([]))
```

### `freeze`

`freeze` freezes a value, if it is not already frozen.  Its OK to call `freeze`
on values that are already frozen.  The value is returned after it is frozen.

* signature: `freeze(value <Value>) <Freeze>`
* example: `freeze([1, 2])`

### `frozen`

`frozen` returns whether or not a value is frozen.

* signature: `frozen(value <Value>) <Bool>`
* example:

```
println(frozen('a'))
println(frozen([3, 4]))
```

### `has`

`has` returns whether a value has a field with a given name.

* signature: `has(val <Value>, name <Str>) <Bool>`
* example:

```
let a = [1, 2]
println(has(a, 'add'))
```

### `hashCode`

`hashCode` returns the hash code of a [`hashable`](interfaces.html#hashable) value.

* signature: `hashCode(val <Value>) <Int>`
* example:

```
println(hashCode('abc'))
```

### `iter`

`iter` returns an "iterator" [Struct](struct.html) for an [iterable](interfaces.html#iterable)
value.

A new iterator must have `next()` called on it to advance
to the first available value. Calling `get()` before the first call to `next()`
throws an error.

* signature: `iter(itr <Iterable>) <Struct>`
* example:

```
let a = [1, 2, 3]
let itr = iter(a)
while itr.next() {
    println(itr.get())
}
```

An iterator struct has the following fields.

* `next()` returns whether there are any more values in the iterator,
and advances the iterator forwards if there is another value.

    * signature: `next() <Bool>`

* `get()` returns the currently available value.

    * signature: `get() <Value>`

### `len`

`len` returns the length of a value that has a [length](interfaces.html#lenable).

* signature: `len(value <Lenable>) <Int>`
* example: `println(len('abc'))`

### `merge`

`merge` combines an arbitrary number of existing [structs](Struct.html) into a new struct.  If
there are any duplicated keys in the structs passed in to 'merge()', then the
value associated with the *last* such key is used.

* signature: `merge(structs... <Struct>) <Struct>`
* example:

```
let a = struct { x: 1, y: 2 }
let b = struct { y: 3, z: 4 }
let c = merge(a, b)

println('a: ', a)
println('b: ', b)
println('c: ', c)

a.x = 10

println()
println('a: ', a)
println('b: ', b)
println('c: ', c) // x is changed here too!
```

### `range`

`range` creates a [Range](range.html), starting at "from" (inclusive) and going until
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

### `str`

`str` returns a Str representation of a value.

* signature: `str(value <Value>) <Str>`
* example: `println(str([null, true, 1, 'abc']))`

### `stream`

`stream` returns a "stream" [Struct](struct.html) for an [Iterable](interfaces.html#iterable)
value.

A stream performs a series of transforms on a sequence of iterated values, and then
collects the values into a final result.

Streams have two kinds of fields:

* Transformer functions, which perform some kind of transformation on the sequence of values in the stream.

* Collector functions, which turn the sequence of values into a final result.

Streams are lazy -- calling one of the transformer
functions doesn't do any processing, it simply registers a new transformation,
and then returns the modified stream.  Processing on the sequence of values in the stream
does not start until one of the collector functions is called.

Streams are "single use" values.  Once one of the collector functions has been called,
an error will be thrown if any of the stream's functions are called.

* signature: `stream(itr <Iterable>) <Struct>`
* example:

```
// print the sum of the even squares
let a = [1, 2, 3, 4, 5]
println(stream(a)
    .map(|e| => e*e)
    .filter(|e| => e % 2 == 0)
    .reduce(0, |acc, e| => acc + e))
```

A stream has the following fields (TODO add a bunch more):

#### transformer functions

* `filter()` adds a "filter" transformation to the stream, by removing elements
which do not match the provided predicate function.  The predicate function
must accept one value, and must return a boolean value.

    * signature: `filter(predicate <Func>) <Stream>`
    * predicate signature: `fn(val <Value>) <Bool>`

* `map()` adds a "map" transformation to the stream, by transforming elements
according to the provided mapping function.  The mapping function must accept
one value, and must return one value.

    * signature: `map(mapping <Func>) <Stream>`
    * mapping signature: `fn(val <Value>) <Value>`

#### collector functions

* `reduce()` reduces the stream to a single value, by applying a "reducer" function
to an accumulated value and each element in the stream.
Accumulation is done starting with the first element in the stream,
and ending with the last.  The reducer function must accept two values, and return one value.

    * signature: `reduce(initial <Value>, reducer <Func>) <List>`
    * reducer signature: `fn(accum <Value>, val <Value>) <Value>`

* `toList()` collects the stream's sequence of values into a [List](list.html).

    * signature: `toList() <List>`

### `type`

`type` returns the type of a value.

* signature: `type(value <Value>) <Str>`
* example:

```
println(type(1.23))
let a = [null, true, 1, 'xyz']
println(a.map(type))
```


## SideEffect Builtins

Golem also has "side-effect" builtins.  These functions
perform I/O, so they should not be included in sandboxed Golem
environments.

* [`print()`](#print)
* [`println()`](#println)

### `print`

`print` prints a sequence of values to STDOUT.

* signature: `print(values... <Value>) <Null>`

### `println`

`println` prints a sequence of values to STDOUT, followed by a linefeed.

* signature: `println(values... <Value>) <Null>`

