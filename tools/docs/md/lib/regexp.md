
# regexp

Module regexp implements regular expression search.

The syntax of the regular expressions accepted is the same general syntax used by
Perl, Python, and other languages.


`regexp` has the following fields:

* [compile](#compile)

`regexp` defines the following structs:

* [regex](#regex)


## Fields

### `compile`

`compile` parses a regular expression and returns, if successful, a
[regex](#regex) struct that can be used to match against text.

* signature: `compile(expr <Str>) <Struct>`


## Structs

### `regex`

`regex` is a Struct that is the representation of a compiled regular expression. A
Regexp is safe for concurrent use by multiple goroutines, except for configuration
methods, such as Longest.

A `regex` struct has the fields:

* [matchString](#matchstring)

#### `matchString`

`matchString` reports whether the string s contains any match of the regular expression.

* signature: `matchString(<Str>) <Bool>`

