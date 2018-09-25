
### `encoding.json`

`encoding.json` implements encoding and decoding of JSON as defined in RFC 7159

#### `marshal`

`marshal` returns the JSON encoding of a value.  Null, Bool, Float Int, Str, and List
are marshalled as their corresponding JSON elements.  Structs and Dicts are marshalled
as JSON objects. Other golem types cannot be marshalled.

* signature: `marshal(value <Value>) <Str>`
* example:

```
import encoding

let s = struct { a: [1, 2, 3]}
println(encoding.json.marshal(s))
```
#### `marshalIndent`

`marshalIndent` is like `marshal` but applies indent to format the output.
Each JSON element in the output will begin on a new line beginning with prefix
followed by one or more copies of indent according to the indentation nesting.

* signature: `marshalIndent(value <Value>, prefix <Str>, indent <Str>) <Str>`
* example:

```
import encoding

let s = struct { a: [1, 2, 3]}
println(encoding.json.marshalIndent(s, '', '  '))
```
#### `unmarshal`

`unmarshal` parses JSON-encoded data.
The optional useStructs parameter, which defaults to false, specifies whether
the data should be marshalled into structs rather than dicts.  If this parameter
is set to true, then the keys of the JSON objects in the data must all be valid
Golem identifiers.

* signature: `unmarshal(text <Str>, useStructs = false <Bool>) <Value>`
* example:

```
import encoding

let text = `{
  "a": [ 1, 2, 3 ]
}`

println(encoding.json.unmarshal(text))
println(encoding.json.unmarshal(text, true))
```
