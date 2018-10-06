
## golem

The golem module defines functions that perform introspection and transformation
on Golem values and code, as well as providing a few useful utility functions.

`golem` has the following fields:

* [getField](#getfield)
* [setField](#setfield)
* [toDict](#todict)

### `getField`

`getField` returns the value associated with a field name.

* signature: `getField(value <Value>, name <Str>) <Value>`
* example:

```
import golem
let a = [1, 2]
let f = golem.getField(a, 'add')
f(3)
println(a)
```

### `makeHashCode`

`makeHashCode` generates a hashCode for a sequence of
[`hashable`](interfaces.html#hashable) values. `makeHashCode` uses the
[Jenkins hash function](https://en.wikipedia.org/wiki/Jenkins_hash_function)
algorithm.

* signature: `makeHashCode(values... <Value>) <Int>`
* example:

```
import golem
println(golem.makeHashCode(1, 2, 3))
```

### `setField`

`setField` sets the value associated with a field name. `setField` only works
on Structs -- you cannot set the fields of other types. `setField` returns `null`
if it was successful.

* signature: `setField(s <Struct>, name <Str>, value <Value>) <Null>`
* example:

```
import golem
let s = struct { a: 1, b: 2 }
golem.setField(s, 'a', 3)
println(s)
```

### `toDict`

`toDict` converts a Struct into a Dict.

* signature: `toDict(s <Struct>) <Dict>`
* example:

```
import golem
println(golem.toDict(struct { a: 1, b: 2 }))
```

