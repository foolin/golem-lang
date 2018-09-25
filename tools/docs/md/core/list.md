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

  * [add](#add)
  * [addAll](#addall)
  * [clear](#clear)
  * [contains](#contains)
  * [filter](#filter)
  * [indexOf](#indexof)
  * [isEmpty](#isempty)
  * [join](#join)
  * [map](#map)
  * [reduce](#reduce)
  * [remove](#remove)
  * [sort](#sort)

### `add`

`add` adds a value to the end of the list, and returns the modified list.

* signature: `add(val <Value>) <List>`
* example:

```
let a = [1, 2, 3]
println(a.add(4))
```

### `addAll`

`addAll` adds all of the values in the given [Iterable](#TODO) to the end of the list,
and returns the modified list.

* signature: `addAll(itr <Iterable>) <List>`
* example:

```
let a = [1, 2]
println(a.addAll([3, 4]))
```

### `clear`

`clear` removes all of the values from the list, and returns the empty list.

* signature: `clear() <List>`
* example:

```
let a = [1, 2]
println(a.clear())
```

### `contains`

`contains` returns whether the given value is an element in the list.

* signature: `contains(val <Value>) <Bool>`
* example:

```
let a = [1, 2]
println(a.contains(2))
```

### `filter`

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

### `indexOf`

`indexOf` returns the index of the given value in the list, or -1 if the value
is not contained in the list.

* signature: `indexOf(val <Value>) <Int>`
* example:

```
let a = ['x', 'y', 'z']
println(a.indexOf('z'))
```

### `isEmpty`

`isEmpty` returns whether the list contains any values.

* signature: `isEmpty() <Bool>`
* example: `println([].isEmpty())`

### `join`

Join concatenates the [`str()`](#TODO) of the elements of the list to create a single string.
The separator string sep is placed between elements in the resulting string.
The sep parameter is optional, and defaults to the empty string `''`.

* signature: `join(sep = '' <Str>) <Str>`
* example: `println([1,2,3].join(', '))`

### `map`

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

### `reduce`

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

### `remove`

`remove` remove the value at the given index from the list, and returns the
modified list.

* signature: `remove(index <Int>) <List>`
* example: `println(['a','b','c'].remove(2))`


### `sort`

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

