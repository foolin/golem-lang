## Range

A Range is a representation of an immutable sequence of integers.
Note that a Range doesn't actually contain a list of all its
Ints -- it produces them one at a time on demand.
A new Range is created by the [`range()`](#TODO) builtin function.

Valid operators for Range are:

* The equality operators `==`, `!=`
* The index operator `a[x]`

The index operator always return an Int.

Ranges have a [`len()`](#TODO) and are [`iterable`](#TODO).

Range has the following fields:

  * [count](#count)
  * [from](#from)
  * [step](#step)
  * [to](#to)

### `count`

`count` is the total number of Ints in the range.

* signature: `count() <Int>`

### `from`

`from` is the first Int in the range, inclusive

* signature: `from() <Int>`

### `step`

`step` is the distance between succesive Ints in the range.

* signature: `step() <Int>`

### `to`

`to` is the last Int in the range, exclusive

* signature: `to() <Int>`

