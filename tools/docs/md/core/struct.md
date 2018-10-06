## Struct

A Struct is a collection of fields that are defined in golem code.

Valid operators for Struct are:

* The equality operators `==`, `!=`

Structs do not have any pre-defined fields.

Structs can have the following magic fields:

### `__eq__`

`__eq__` overrides the `==` operator

* signature: `__eq__(x <Value>) <Bool>`

### `__hashCode__`

`__hashCode__` causes a struct to be [hashable](interfaces.html#hashable), so it can be
used as a key in a dict, or an entry in a set.  Note that if you define `__hashCode__`,
you *must* also always define `__eq__`.

* signature: `__hashCode__() <Int>`

### `__str__`

`__str__` overrides the value returned by the builtin function [`str`](builtins.html#str)

* signature: `__str__() <Str>`

