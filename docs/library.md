# Standard Library

* [encoding](#encoding)
  * [json](#encodingjson)
* [golem](#golem)
* [os](#os)
* [path](#path)
  * [filepath](#pathfilepath)


## encoding

The encoding module defines functionality that converts
data to and from byte-level and textual representations.


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

## golem

The golem module defines functions that perform introspection and transformation
on Golem values and code.

#### `fields`

`fields` returns a Set of the names of a value's fields.

* signature: `fields(value <Value>) <Set>`
* example:

```
    import golem

    println(golem.fields([]))
    println(golem.fields(struct { a: 1, b: 2}))
```

#### `getField`

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

#### `hasField`

`hasField` returns whether a value has a field with a given name.

* signature: `getField(value <Value>, name <Str>) <Bool>`
* example:

```
    import golem

    let a = [1, 2]
    println(golem.hasField(a, 'add'))
```

#### `setField`

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


## os

Module os provides a platform-independent interface to operating system functionality.


### Functions

#### `create`

`create` creates the named file with mode 0666 (before umask), truncating it if
it already exists. If successful, methods on the returned [file](#file)
can be used for I/O.

* signature: `create(name <Str>) <Struct>`
* example:

```
    import os

    let f = os.create('foo.txt')
    try {
        // do something with the file
    } finally {
        f.close()
    }
```

#### `exit`

`exit` causes the current program to exit with the given status code. Conventionally,
code zero indicates success, non-zero an error. The program terminates immediately.

* signature: `exit(code <Int>) <Null>`
* example:

```
import os

os.exit(-1)
```

#### `open`

`open` opens the named file for reading. If successful, methods on the
returned [file](#file) can be used for reading.

* signature: `open(name <Str>) <Struct>`
* example:

```
    import os

    let f = os.open('foo.txt')
    try {
        // do something with the file
    } finally {
        f.close()
    }
```

#### `stat`

`stat` returns a [fileInfo](#fileinfo) describing the named file.

* signature: `stat(name <Str>) <Struct>`
* example:

```
    import os

    let s = os.stat('foo.txt')
    println([
        'name: ' + s.name(),
        'size: ' + s.size(),
        'mode: ' + s.mode(),
        'isDir: ' + s.isDir()
    ])
```


### Structs

#### `fileInfo`

A `fileInfo` is a struct that describes a file and is returned by stat.

##### `name`
`name` is the base name of the file
* signature: `name() <Str>`
##### `size`
`size` is the length in bytes for regular files; system-dependent for others
* signature: `size() <Int>`
##### `mode`
`mode` is the file mode bits
* signature: `mode() <Int>`
##### `isDir`
`isDir` is an abbreviation for Mode().IsDir()
* signature: `isDir() <Str>`
#### `file`

A `file` is a struct that represents an open file descriptor.

##### `readLines`
`readLines` returns a List of Strs, for each line of text in the file.
* signature: `readLines() <List>`
##### `writeLines`
`writeLines` writes a List of Strs to the file as a sequence of lines.
* signature: `writeLines(<List>) <Null>`
##### `close`
`close` closes the File, rendering it unusable for I/O. On files that support
SetDeadline, any pending I/O operations will be canceled and return immediately
with an error.
* signature: `close() <Null>`

## path

Module path implements utility routines for manipulating slash-separated paths.


### `path.filepath`

`path.filepath` implements utility routines for manipulating filename paths in a
way compatible with the target operating system-defined file paths.

#### `ext`

`ext` returns the file name extension used by path. The extension is the suffix
beginning at the final dot in the final element of path; it is empty if there is no dot.

* signature: `ext(name <Str>) <Str>`
* example:

```
    import path

    println(path.filepath.ext('foo.txt'))
```

_This document uses documentation from [go](https://github.com/golang/go), which
is licensed under the BSD-3-Clause license._
