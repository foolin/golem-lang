
## golem

The golem module defines functions that perform introspection and transformation
on Golem values and code.

`golem` has the following fields:

* [fields](#fields)
* [getField](#getField)
* [hasField](#hasField)
* [setField](#setField)

### `fields`

`fields` returns a Set of the names of a value's fields.

* signature: `fields(value <Value>) <Set>`
* example:

```
import golem
println(golem.fields([]))
println(golem.fields(struct { a: 1, b: 2}))
```

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

### `hasField`

`hasField` returns whether a value has a field with a given name.

* signature: `getField(value <Value>, name <Str>) <Bool>`
* example:

```
import golem
let a = [1, 2]
println(golem.hasField(a, 'add'))
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

