// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package ioutil

import (
	"fmt"
	"io/ioutil"
	"os"

	g "github.com/mjarmy/golem-lang/core"
	gos "github.com/mjarmy/golem-lang/lib/os"
)

/*doc

## `io.ioutil`

`io.ioutil` implements some I/O utility functions.

*/

/*doc
`io.ioutil` has the following fields:

* [readDir](#readDir)
* [readFileString](#readFileString)
* [writeFileString](#writeFileString)

*/

/*doc
### `readDir`

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

*/

// ReadDir reads a directory
var ReadDir g.Value = g.NewFixedNativeFunc(
	[]g.Type{g.StrType}, false,
	func(ev g.Eval, params []g.Value) (g.Value, g.Error) {
		filename := params[0].(g.Str)

		infos, err := ioutil.ReadDir(filename.String())
		if err != nil {
			return nil, g.Error(fmt.Errorf("IoError: %s", err.Error()))
		}
		values := make([]g.Value, len(infos))
		for i, inf := range infos {
			values[i] = gos.NewFileInfo(inf)
		}
		return g.NewList(values), nil
	})

/*doc
### `readFileString`

`readFileString` reads an entire file as a stirng.

* signature: `readFile(filename <Str>) <Str>`
* example:

```
import io
println(io.ioutil.readFileString('testdata.txt'))
```

*/

// ReadFileString reads an entire file as a stirng.
var ReadFileString g.Value = g.NewFixedNativeFunc(
	[]g.Type{g.StrType}, false,
	func(ev g.Eval, params []g.Value) (g.Value, g.Error) {
		filename := params[0].(g.Str)

		content, err := ioutil.ReadFile(filename.String())
		if err != nil {
			return nil, g.Error(fmt.Errorf("IoError: %s", err.Error()))
		}
		return g.NewStr(string(content))
	})

/*doc
### `writeFileString`

`writeFileString` writes a string to a file

* signature: `writeFileString(filename <Str>, data <Str>) <Null>`
* example:

```
import io
io.ioutil.writeFileString('testdata.txt', 'abc')
```

*/

// WriteFileString writes a string to a file
var WriteFileString g.Value = g.NewFixedNativeFunc(
	[]g.Type{g.StrType, g.StrType}, false,
	func(ev g.Eval, params []g.Value) (g.Value, g.Error) {
		filename := params[0].(g.Str)
		data := params[1].(g.Str)

		// todo pass in FileMode
		fileMode := os.FileMode(0666)
		err := ioutil.WriteFile(filename.String(), []byte(data.String()), fileMode)
		if err != nil {
			return nil, g.Error(fmt.Errorf("IoError: %s", err.Error()))
		}
		return g.Null, nil
	})
