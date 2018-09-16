# Types
Golem has a set of standard types.  Each type has a set of values,
together with operators and fields specific to those values.

* Basic Types:
  * [Null](#null)
  * [Bool](#bool)
  * [Int](#int)
  * [Float](#float)
  * [Str](#str)
* Composite Types:
  * [List](#list)
  * [Range](#range)
  * [Tuple](#tuple)
  * [Dict](#dict)
  * [Set](#set)
  * [Struct](#struct)
* Other Types:
  * [Func](#func)
  * [Chan](#chan)

# Basic Types

## Null

Null represents the absence of a value. The only instance of Null is `null`.

Null has no valid operators, and no fields.

## Bool

Bool represents boolean truth values.  The only instances of
Bool are `true` and `false`.

Valid operators for Bool are:
* The equality operators `==`, `!=`
* The boolean operators `||` `&&`
* The unary negation operator `!`

Bools are [`hashable`](#TODO)

Bool has no fields.

## Int

Int is the set of all signed 64-bit integers (-9223372036854775808 to 9223372036854775807).

Valid operators for Int are:
* The equality operators `==`, `!=`
* The comparison operators `>`, `>=`, `<`, `<=`, `<=>`
* The arithmetic operators `+`, `-`, `*`, `/`
* The integer arithmetic operators <code>&#124;</code>, `^`, `%`, `&`, `<<`, `>>`
* The unary integer complement operator `~`
* The postfix operators `++`, `--`

When applying an arithmetic operator `+`, `-`, `*`, `/`to an Int, if the other
operand is a Float, then the result will be a Float, otherwise the result will be an Int.

Ints are [`hashable`](#TODO)

Int has no fields.

## Float

Float is the set of all IEEE-754 64-bit floating-point numbers.

Valid operators for Float are:
* The equality operators `==`, `!=`
* The comparison operators `>`, `>=`, `<`, `<=`, `<=>`
* The arithmetic operators `+`, `-`, `*`, `/`
* The postfix operators `++`, `--`

Applying an arithmetic operator to a Float always returns a Float.

Floats are [`hashable`](#TODO)

Float has no fields.

## Str

Str is the set of all valid sequences of UTF-8-encoded "code points", otherwise
known as "runes".  Strs are immutable.

Valid operators for Str are:
* The equality operators `==`, `!=`,
* The comparison operators `>`, `>=`, `<`, `<=`, `<=>`
* The index operator `a[x]`
* The slice operators `a[x:y]`, `a[x:]`, `a[:y]`

The index operator always returns a single-rune Str.

The slice operators always return a Str.

Strs have a [`len()`](#TODO) and are [`iterable`](#TODO).

Strs are [`hashable`](#TODO)

Str has the following fields:

#### `contains`

`contains` reports whether a substring is within a string.

* signature: `contains(substr <Str>) <Bool>`
* example: `'abcdef'.contains('de')`

#### `hasPrefix`

`hasPrefix` tests whether a string begins with a prefix.

* signature: `hasPrefix(prefix <Str>) <Bool>`
* example: `'abcdef'.hasPrefix('ab')`

#### `hasSuffix`

`hasSuffix` tests whether a string ends with a suffix.

* signature: `hasSuffix(suffix <Str>) <Bool>`
* example: `'abcdef'.hasSuffix('ab')`

#### `index`

`index` returns the index of the first instance of a substring in a string.
or -1 if the substring is not present.

* signature: `index(substr <Str>) <Int>`
* example: `'abcab'.index('ab')`

#### `lastIndex`

`lastIndex` returns the index of the last instance of a substring in a string,
or -1 if the substring is not present.

* signature: `lastIndex(substr <Str>) <Int>`
* example: `'abcab'.lastIndex('ab')`

#### `map`

`map` returns a copy of the string with all its characters modified according to
the mapping function.

The mapping function must accept one Str parameter, and must return a Str.

* signature: `map(mapping <Func>) <Str>`
* mapping signature: `fn(s <Str>) <Str>`
* example:

```
let s = 'abc(def)[x,y,z]'
let t = s.map(fn(c) {
    return c >= 'a' && c <= 'z' ? c : ''
})
println(t)
```

#### `replace`

`replace` returns a copy of a string with the first n non-overlapping instances
of `old` replaced by `new`. If `old` is empty, it matches at the beginning of a string
and after each UTF-8 sequence, yielding up to k+1 replacements for a k-rune string.
If `n` < 0, there is no limit on the number of replacements.  The parameter `n` is
optional, and defaults to -1.

* signature: `replace(old <Str>, new <Str>, n = -1 <Int>) <Int>`
* example: `'abcab'.replace('a', 'x')`

#### `split`

`split` slices a string into all substrings separated by sep and returns a list
of the substrings between those separators.

If the string does not contain sep and sep is not empty, `split` returns a list
of length 1 whose only element is the string.

If sep is empty, `split` splits after each UTF-8 sequence. If both the string
and sep are empty, `split` returns an empty list.

* signature: `split(sep <Str>) <List>`
* example: `'a,b,c'.split(',')`

#### `toChars`

`toChars` splits a string into a list of single-rune Strs.

* signature: `toChars() <List>`
* example: `'xyz'.toChars()`

#### `trim`

`trim` returns a new string with all leading and trailing runes contained in cutset removed.

* signature: `trim(<Str>) <Str>`
* example: `'\t\tabc\n'.trim('\t\n')`

# Composite Types

## List

A List is an ordered array of values.

Valid operators for List are:
* The equality operators `==`, `!=`
* The index operator `a[x]`
* The slice operators `a[x:y]`, `a[x:]`, `a[:y]`

The index operator can return a value of any type.

The slice operators always return a List.

Lists have a [`len()`](#TODO) and are [`iterable`](#TODO).

List has the following fields:

#### `add`

`add` adds a value to the end of the list, and returns the modified list.

* signature: `add(val <Value>) <List>`
* example:

```
let a = [1, 2, 3]
println(a.add(4))
```

#### `addAll`

`addAll` adds all of the values in the given [Iterable](#TODO) to the end of the list,
and returns the modified list.

* signature: `addAll(itr <Iterable>) <List>`
* example:

```
let a = [1, 2]
println(a.addAll([3, 4]))
```

#### `clear`

`clear` removes all of the values from the list, and returns the empty list.

* signature: `clear() <List>`
* example:

```
let a = [1, 2]
println(a.clear())
```

#### `contains`

`contains` returns whether the given value is an element in the list.

* signature: `contains(val <Value>) <Bool>`
* example:

```
let a = [1, 2]
println(a.contains(2))
```

#### `filter`

`filter` returns a new list by passing each of the elements of the current list
into the given predicate.  If the predicate returns `true` for an element, that
element is added to the new list.  The original list is unmodified.

The predicate must accept one parameter of any type, and return a Bool.

* signature: `filter(predicate <Func>) <List>`
* predicate signature: `fn(val <Value>) <Bool>`
* example:

```
let a = [1, 2, 3, 4, 5]
println(a.filter(|e| => e % 2 == 0))
```

#### `indexOf`

`indexOf` returns the index of the given value in the list, or -1 if the value
is not contained in the list.

* signature: `indexOf(val <Value>) <Int>`
* example:

```
let a = ['x', 'y', 'z']
println(a.indexOf('z'))
```

#### `isEmpty`

`isEmpty` returns whether the list contains any values.

* signature: `isEmpty() <Bool>`
* example: `println([].isEmpty())`

#### `join`

Join concatenates the [`str()`](#TODO) of the elements of the list to create a single string.
The separator string sep is placed between elements in the resulting string.
The sep parameter is optional, and defaults to the empty string `''`.

* signature: `join(sep = '' <Str>) <Str>`
* example: `println([1,2,3].join(', '))`

#### `map`

`map` returns a copy of the list with all its elements modified according to
the mapping function.  The original list is unmodified.

The mapping function must accept one value, and must return one value.

* signature: `map(mapping <Func>) <List>`
* mapping signature: `fn(val <Value>) <Value>`
* example:

```
let a = [1,2,3]
let b = a.map(|e| => e * e)
println(b)
```

#### `reduce`

`reduce` reduces the list to a single value, by applying a "reducer" function
to an accumulated value and each element in the list.
Accumulation is done starting with the first element in the list,
and ending with the last.  The original list is unmodified.

The reducer function must accept two values, and return one value.

* signature: `reduce(start <Value>, reducer <Func>) <List>`
* reducer signature: `fn(accum <Value>, val <Value>) <Value>`
* example:

```
let a = [1,2,3]
let b = a.reduce(0, |acc, e| => acc + e)
println(b)
```

#### `remove`

`remove` remove the value at the given index from the list, and returns the
modified list.

* signature: `remove(index <Int>) <List>`
* example: `println(['a','b','c'].remove(2))`


#### `sort`

`sort` sorts the elements in the list and returns the modified list.  If the
optional "lesser" function is provided, it is used to compare values in the list.
If the lesser function is not provided, then the `<` operator is used.

* signature: `sort(lesser = null <Func>) <List>`
* lesser signature: `fn(val <Value>, val <Value>) <Bool>`
* example:

```
let a = [7, 4, 11, 13, 6, 2, 9, 1]
a.sort(|a, b| => b < a) // sort in reverse
println(a)
```

## Tuple

A Tuple is an immutable sequence of two or more values.  Tuples are similar to Lists,
but they have special semantics in certain Golem expressions and statements.

Valid operators for Tuple are:
* The equality operators `==`, `!=`
* The index operator `a[x]`

The index operator can return a value of any type.

Tuples have a [`len()`](#TODO).

Tuples are [`hashable`](#TODO)

Tuples have no fields.

## Range

A Range is a representation of an immutable sequence of integers.
Note that a Aange doesn't actually contain a list of all its
Ints -- it produces them one at a time on demand.
A new Range is created by the [`range()`](#TODO) builtin function.

Valid operators for Range are:
* The equality operators `==`, `!=`
* The index operator `a[x]`

The index operator always return an Int.

Ranges have a [`len()`](#TODO) and are [`iterable`](#TODO).

Range has the following fields:

#### `count`

`count` is the total number of Ints in the range.

* signature: `count() <Int>`

#### `from`

`from` is the first Int in the range, inclusive

* signature: `from() <Int>`

#### `step`

`step` is the distance between succesive Ints in the range.

* signature: `step() <Int>`

#### `to`

`to` is the last Int in the range, exclusive

* signature: `to() <Int>`

## Dict

A Dict is an [associative array](https://en.wikipedia.org/wiki/Associative_array),
in which the keys are all [`hashable`](#TODO).

Valid operators for Dict are:
* The equality operators `==`, `!=`
* The index operator `a[x]`

Dict have a [`len()`](#TODO) and are [`iterable`](#TODO).

Each iterated element in a Dict is a 2-Tuple containing a key-value pair.

Dict has the following fields:

#### `addAll`

`addAll` adds all of the values in the given [Iterable](#TODO) to the dict,
and returns the modified dict.
Each iterated element must be a 2-Tuple containing a key-value pair.

* signature: `addAll(itr <Iterable>) <Dict>`
* example:

```
let d = dict {'a': 1, 'b': 2}
println(d.addAll([('b', 2), ('c', 3)]))
```

#### `clear`

`clear` removes all of the entries from the dict, and returns the empty dict.

* signature: `clear() <Dict>`
* example:

```
let d = dict {'a': 1, 'b': 2}
println(d.clear())
```

#### `contains`

`contains` returns whether the given key is present in the dict.

* signature: `contains(key <Value>) <Bool>`
* example:

```
let d = dict {'a': 1, 'b': 2}
println(d.contains('b'))
```

#### `isEmpty`

`isEmpty` returns whether the dict contains any values.

* signature: `isEmpty() <Bool>`
* example: `println(dict {}.isEmpty())`

#### `remove`

`remove` remove the entry associated with the given key from the dict,
and returns modified dict.  If the key is not present in the dict, then
the dict is unmodified.

* signature: `remove(key <Value>) <Dict>`
* example:

```
let d = dict {'a': 1, 'b': 2}
println(d.remove('a'))
```

## Set

A Set is an un-ordered collection of unique, [`hashable`](#TODO) values.

Valid operators for Set are:
* The equality operators `==`, `!=`

Sets have a [`len()`](#TODO) and are [`iterable`](#TODO).

Set has the following fields:

#### `add`

`add` adds a value to the set, and returns the modified set.

* signature: `add(val <Value>) <Set>`
* example:

```
let a = set {1, 2, 3}
println(a.add(4))
```

#### `addAll`

`addAll` adds all of the values in the given [Iterable](#TODO) to the set,
and returns the modified set.

* signature: `addAll(itr <Iterable>) <Set>`
* example:

```
let a = set {1, 2}
println(a.addAll([3, 4]))
```

#### `clear`

`clear` removes all of the values from the set, and returns the empty set.

* signature: `clear() <Set>`
* example:

```
let a = set {1, 2}
println(a.clear())
```

#### `contains`

`contains` returns whether the given value is an element in the set.

* signature: `contains(val <Value>) <Bool>`
* example:

```
let a = set {1, 2}
println(a.contains(2))
```

#### `isEmpty`

`isEmpty` returns whether the set contains any values.

* signature: `isEmpty() <Bool>`
* example: `println(set{}.isEmpty())`

#### `remove`

`remove` remove the value from the set, and returns the
modified set.  If the value is not present in the set, then
the set is unmodified.

* signature: `remove(value <Value>) <Set>`
* example:

```
let a = set {1, 2, 3}
println(a.remove(2))
```

# Other Types

## Func

A Func is a sequence of [`expressions`](#TODO) and [`statements`](#TODO) that can
be invoked to perform a task.

Valid operators for Func are:
* The equality operators `==`, `!=`
* The invocation operator `a(x)`

Funcs have an [`arity()`](#TODO)

Funcs have no fields.

## Chan

A Chan is a conduit through which you can send and receive values.
A new Chan is created by the [`chan()`](#TODO) builtin function.

Valid operators for Chan are:
* The equality operators `==`, `!=`

Chan has the following fields:

#### `send`

`send` sends a value to the chan.

* signature: `send(val <Value>)`

#### `recv`

`recv` receives a value from the chan.

* signature: `recv() <Value>`


_This document uses documentation from [go](https://github.com/golang/go), which
is licensed under the BSD-3-Clause license._
