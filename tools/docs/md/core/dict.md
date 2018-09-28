## Dict

A Dict is an [associative array](https://en.wikipedia.org/wiki/Associative_array),
in which the keys are all [`hashable`](#TODO).

Valid operators for Dict are:

* The equality operators `==`, `!=`
* The index operator `a[x]`

Dict have a [`len()`](#TODO) and are [`iterable`](#TODO).

Each iterated element in a Dict is a 2-Tuple containing a key-value pair.

Dict has the following fields:

* [addAll](#addall)
* [clear](#clear)
* [contains](#contains)
* [isEmpty](#isempty)
* [remove](#remove)

### `addAll`

`addAll` adds all of the values in the given [Iterable](#TODO) to the dict,
and returns the modified dict.
Each iterated element must be a 2-Tuple containing a key-value pair.

* signature: `addAll(itr <Iterable>) <Dict>`
* example:

```
let d = dict {'a': 1, 'b': 2}
println(d.addAll([('b', 2), ('c', 3)]))
```

### `clear`

`clear` removes all of the entries from the dict, and returns the empty dict.

* signature: `clear() <Dict>`
* example:

```
let d = dict {'a': 1, 'b': 2}
println(d.clear())
```

### `contains`

`contains` returns whether the given key is present in the dict.

* signature: `contains(key <Value>) <Bool>`
* example:

```
let d = dict {'a': 1, 'b': 2}
println(d.contains('b'))
```

### `isEmpty`

`isEmpty` returns whether the dict contains any values.

* signature: `isEmpty() <Bool>`
* example: `println(dict {}.isEmpty())`

### `remove`

`remove` remove the entry associated with the given key from the dict,
and returns modified dict.  If the key is not present in the dict, then
the dict is unmodified.

* signature: `remove(key <Value>) <Dict>`
* example:

```
let d = dict {'a': 1, 'b': 2}
println(d.remove('a'))
```

