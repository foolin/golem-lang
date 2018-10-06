## Set

A Set is an un-ordered collection of unique, [`hashable`](interfaces.html#hashable) values.

Valid operators for Set are:

* The equality operators `==`, `!=`

Sets are
[`lenable`](interfaces.html#lenable) and
[`iterable`](interfaces.html#iterable).

A Set has the following fields:

* [add](#add)
* [addAll](#addall)
* [clear](#clear)
* [contains](#contains)
* [containsAll](#containsall)
* [containsAny](#containsany)
* [copy](#copy)
* [isEmpty](#isempty)
* [remove](#remove)

### `add`

`add` adds a value to the set, and returns the modified set.

* signature: `add(val <Value>) <Set>`
* example:

```
let a = set {1, 2, 3}
println(a.add(4))
```

### `addAll`

`addAll` adds all of the values in the given [Iterable](interfaces.html#iterable) to
the set, and returns the modified set.

* signature: `addAll(itr <Iterable>) <Set>`
* example:

```
let a = set {1, 2}
println(a.addAll([3, 4]))
```

### `clear`

`clear` removes all of the values from the set, and returns the empty set.

* signature: `clear() <Set>`
* example:

```
let a = set {1, 2}
println(a.clear())
```

### `contains`

`contains` returns whether the given value is an element in the set.

* signature: `contains(val <Value>) <Bool>`
* example:

```
let a = set {1, 2}
println(a.contains(2))
```

### `containsAll`

`containsAll` returns whether the set contains all of the values
in the given [Iterable](interfaces.html#iterable).

* signature: `containsAll(itr <Iterable>) <Bool>`
* example:

```
let a = set {1, 2}
println(a.containsAll([1, 2]))
```

### `containsAny`

`containsAny` returns whether the set contains all of the values
in the given [Iterable](interfaces.html#iterable).

* signature: `containsAny(itr <Iterable>) <Bool>`
* example:

```
let a = set {1, 2}
println(a.containsAny([1, 2]))
```

### `copy`

`copy` returns a shallow copy of the set

* signature: `copy() <Set>`
* example:

```
println(set {1,2}.copy())
```

### `isEmpty`

`isEmpty` returns whether the set contains any values.

* signature: `isEmpty() <Bool>`
* example: `println(set{}.isEmpty())`

### `remove`

`remove` remove the value from the set, and returns the
modified set.  If the value is not present in the set, then
the set is unmodified.

* signature: `remove(value <Value>) <Set>`
* example:

```
let a = set {1, 2, 3}
println(a.remove(2))
```

