
### `io.ioutil`

`io.ioutil` implements some I/O utility functions.

#### `readDir`

`readDir` reads the directory named by dirname and returns a
list of directory entries sorted by filename. The resulting
list will be [fileinfo](#TODO) Structs.

* signature: `readDir(filename <Str>) <List>`
* example:

```
import io

let files = io.ioutil.readDir('.')
    for f in files {
        println([f.name(), f.isDir()])
    }
```
#### `readFileString`

`readFileString` reads an entire file as a stirng.

* signature: `readFile(filename <Str>) <Str>`
* example:

```
import io

println(io.ioutil.readFileString('testdata.txt'))
```
#### `writeFileString`

`writeFileString` writes a string to a file

* signature: `writeFileString(filename <Str>, data <Str>) <Null>`
* example:

```
import io

io.ioutil.writeFileString('testdata.txt', 'abc')
```
