## Tuple

A Tuple is an immutable sequence of two or more values.  Tuples are similar to Lists,
but they have special semantics in certain Golem expressions and statements.

Valid operators for Tuple are:

* The equality operators `==`, `!=`
* The [`index`](interfaces.html#indexable) operator `a[x]`

The index operator can return a value of any type.

Tuples are
[`lenable`](interfaces.html#lenable) and
[`hashable`](interfaces.html#hashable).

Tuples have no fields.

