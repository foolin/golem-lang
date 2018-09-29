## Int

Int is the set of all signed 64-bit integers (-9223372036854775808 to 9223372036854775807).

Integer literals can either be in decimal format, e.g. `123`, or hexidecimal format,
e.g. `0xabcd`.

Valid operators for an Int are:

* The equality operators `==`, `!=`
* The [`comparision`](interfaces.html#comparable) operators `>`, `>=`, `<`, `<=`, `<=>`
* The arithmetic operators `+`, `-`, `*`, `/`
* The integer arithmetic operators <code>&#124;</code>, `^`, `%`, `&`, `<<`, `>>`
* The unary integer complement operator `~`
* The postfix operators `++`, `--`

When applying an arithmetic operator `+`, `-`, `*`, `/` to an Int, if the other
operand is a Float, then the result will be a Float, otherwise the result will be an Int.

Ints are [`hashable`](interfaces.html#hashable)

An Int has no fields.

