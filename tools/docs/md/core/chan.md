## Chan

A Chan is a conduit through which you can send and receive values.
A new Chan is created by the [`chan()`](builtins.html#chan) builtin function.

Valid operators for Chan are:

* The equality operators `==`, `!=`

Chan has the following fields:

* [send](#send)
* [recv](#recv)

### `send`

`send` sends a value to the chan.

* signature: `send(val <Value>)`

### `recv`

`recv` receives a value from the chan.

* signature: `recv() <Value>`

