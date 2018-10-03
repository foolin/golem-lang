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

An Int has the following fields:

* [abs](#abs)
* [format](#format)
* [toChar](#tochar)
* [toFloat](#tofloat)

### `abs`

`abs` returns the absolute value of the int.

* signature: `abs() <Int>`
* example: `let n = -1; println(n.abs())`

### `format`

`format` returns the string representation of int in the given base,
for 2 <= base <= 36. The result uses the lower-case letters 'a' to 'z'
for digit values >= 10.  If the base is omitted, it defaults to 10.

* signature: `format(base = 10 <Int>) <Str>`
* example: `let n = 11259375; println(n.format(16))`

### `toChar`

`toChar` converts an int that is a valid rune into a string with a single
unicode character.

* signature: `toChar() <Str>`
* example: `let n = 19990; println(n.toChar())`

### `toFloat`

`toFloat` converts an int to a float

* signature: `toFloat() <Float>`
* example: `let n = 123; println(n.toFloat())`

