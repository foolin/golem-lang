
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
