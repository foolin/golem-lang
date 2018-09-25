
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
#### `walk`

`walk` walks the file tree rooted at root, calling walkFn for each file or directory
in the tree, including root.  The walkFn must accept two parameters.  The first
parameter will be the path of the current file, and the second will be the
[fileInfo](TODO) for the file.

* signature: `Walk(root <Str>, walkFn <Func>) <Func>`
* example:

```
path.filepath.walk('.', fn(path, info) {
println([path, info.name(), info.isDir()])
})
```
