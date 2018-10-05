## Float

Float is the set of all IEEE-754 64-bit floating-point numbers.

Valid operators for Float are:

* The equality operators `==`, `!=`
* The [`comparision`](interfaces.html#comparable) operators `>`, `>=`, `<`, `<=`, `<=>`
* The arithmetic operators `+`, `-`, `*`, `/`
* The postfix operators `++`, `--`

Applying an arithmetic operator to a Float always returns a Float.

Floats are [`hashable`](interfaces.html#hashable)

A Float has the following fields:

* [abs](#abs)
* [ceil](#ceil)
* [floor](#floor)
* [format](#format)
* [round](#round)

### `abs`

`abs` returns the absolute value of the float.

* signature: `abs() <Float>`
* example: `let n = -1.2; println(n.abs())`

### `ceil`

`ceil` returns the least integer value greater than or equal to the float.

* signature: `ceil() <Float>`
* example: `let n = -1.2; println(n.ceil())`

### `floor`

`floor` returns the greatest integer value less than or equal to the float.

* signature: `floor() <Float>`
* example: `let n = -1.2; println(n.floor())`

### `format`

`format`

* signature: `format(fmt <Str>, prec = -1 <Int>) <Str>`
* example: `let n = 1.23; println(n.format("f"))`

### `round`

`round` returns the nearest integer, rounding half away from zero.

* signature: `round() <Float>`
* example: `let n = -1.2; println(n.round())`

