## Dict

A Dict is an [associative array](https://en.wikipedia.org/wiki/Associative_array),
in which the keys are all [`hashable`](interfaces.html#hashable).

Valid operators for Dict are:

* The equality operators `==`, `!=`
* The [`index`](interfaces.html#indexable) operator `a[x]`

Dicts are
[`lenable`](interfaces.html#lenable) and
[`iterable`](interfaces.html#iterable).

Each iterated element in a Dict is a 2-Tuple containing a key-value pair.

A Dict has the following fields:

* [addAll](#addall)
* [clear](#clear)
* [contains](#contains)
* [copy](#copy)
* [isEmpty](#isempty)
* [keys](#keys)
* [remove](#remove)
* [toStruct](#tostruct)

### `addAll`

`addAll` adds all of the values in the given [Iterable](interfaces.html#iterable)
to the dict, and returns the modified dict.
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

### `copy`

`copy` returns a shallow copy of the dict

* signature: `copy() <Dict>`
* example:

```
println(dict{'a':1,'b':2}.copy())
```

### `isEmpty`

`isEmpty` returns whether the dict contains any values.

* signature: `isEmpty() <Bool>`
* example: `println(dict {}.isEmpty())`

### `keys`

`keys` returns a Set of the dict's keys.

* signature: `keys() <Set>`
* example: `println(dict {'a': 1, 'b': 2}.keys())`

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

### `toStruct`

`toStruct` converts a Dict into a Struct.

* signature: `toStruct() <Struct>`
* example:

```
let d = dict {'a': 1, 'b': 2}; println(d.toStruct())
```

### `values`

`values` returns a Set of the dict's values.

* signature: `values() <Set>`
* example: `println(dict {'a': 1, 'b': 2}.values())`

