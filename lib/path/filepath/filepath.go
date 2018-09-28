// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package path

import (
	"fmt"
	"os"
	"path/filepath"

	g "github.com/mjarmy/golem-lang/core"
	libOs "github.com/mjarmy/golem-lang/lib/os"
)

/*doc

## `path.filepath`

`path.filepath` implements utility routines for manipulating filename paths in a
way compatible with the target operating system-defined file paths.

*/

/*doc
`path.filepath` has the following fields:

* [ext](#ext)
* [walk](#walk)

*/

/*doc
### `ext`

`ext` returns the file name extension used by path. The extension is the suffix
beginning at the final dot in the final element of path; it is empty if there is no dot.

* signature: `ext(name <Str>) <Str>`
* example:

```
import path
println(path.filepath.ext('foo.txt'))
```

*/

// ext returns a file extension
var Ext g.Value = g.NewFixedNativeFunc(
	[]g.Type{g.StrType}, false,
	func(ev g.Eval, params []g.Value) (g.Value, g.Error) {
		s := params[0].(g.Str)

		// TODO: this means that the the path must be
		// valid UTF-8.  Is that what we really want?
		return g.NewStr(filepath.Ext(s.String()))
	})

/*doc
### `walk`

`walk` walks the file tree rooted at root, calling walkFn for each file or directory
in the tree, including root.  The walkFn must accept two parameters.  The first
parameter will be the path of the current file, and the second will be the
[fileInfo](TODO) for the file.

* signature: `Walk(root <Str>, walkFn <Func>) <Func>`
* example:

```
import path
path.filepath.walk('.', fn(path, info) {
	println([path, info.name(), info.isDir()])
})
```

*/

// walk walks a directory path
var Walk g.Value = g.NewFixedNativeFunc(
	[]g.Type{g.StrType, g.FuncType}, false,
	func(ev g.Eval, params []g.Value) (g.Value, g.Error) {

		dir := params[0].(g.Str)
		callback := params[1].(g.Func)

		arity := callback.Arity()
		if arity.Kind != g.FixedArity || arity.Required != 2 {
			return nil, g.ArityMismatch(2, int(arity.Required))
		}

		err := filepath.Walk(
			dir.String(),
			func(path string, info os.FileInfo, e error) error {
				if e != nil {
					return e
				}

				// TODO: this means that the the path must be
				// valid UTF-8.  Is that what we really want?
				s, err := g.NewStr(path)
				if err != nil {
					return err
				}

				_, err = callback.Invoke(ev,
					[]g.Value{s, libOs.NewFileInfo(info)})
				return err
			})

		if err != nil {
			if gerr, ok := err.(g.Error); ok {
				return nil, gerr
			}
			return nil, g.Error(fmt.Errorf("PathError: %s", err.Error()))
		}
		return g.Null, nil
	})
