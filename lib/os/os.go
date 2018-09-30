// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package os

import (
	"bufio"
	"fmt"
	"io"
	"os"

	g "github.com/mjarmy/golem-lang/core"
	"github.com/mjarmy/golem-lang/lib/os/exec"
)

/*doc

## os

Module os provides a platform-independent interface to operating system functionality.

*/

// Os is the "os" module in the standard library
var Os g.Module

func init() {

	exec, err := g.NewFrozenStruct(
		map[string]g.Field{
			"runCommand": g.NewField(exec.RunCommand),
		})
	g.Assert(err == nil)

	os, err := g.NewFrozenStruct(
		map[string]g.Field{
			"create": g.NewField(create),
			"exec":   g.NewField(exec),
			"exit":   g.NewField(exit),
			"open":   g.NewField(open),
			"stat":   g.NewField(stat),
		})
	g.Assert(err == nil)

	Os = g.NewNativeModule("os", os)
}

/*doc

`os` has the following fields:

* [create](#create)
* [exec](lib_osexec.html)
* [exit](#exit)
* [open](#open)
* [stat](#stat)

`os` defines the following structs:

* [fileInfo](#fileinfo)
* [file](#file)

*/

/*doc

## Fields

*/

/*doc
### `create`

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

*/

// create creates a file
var create g.Value = g.NewFixedNativeFunc(
	[]g.Type{g.StrType}, false,
	func(ev g.Eval, params []g.Value) (g.Value, g.Error) {
		s := params[0].(g.Str)

		f, err := os.Create(s.String())
		if err != nil {
			return nil, g.Error(fmt.Errorf("OsError: %s", err.Error()))
		}
		return newFile(f), nil
	})

/*doc
### `exit`

`exit` causes the current program to exit with the given status code. Conventionally,
code zero indicates success, non-zero an error. The program terminates immediately.

* signature: `exit(code <Int>) <Null>`
* example:

```
import os

os.exit(-1)
```

*/

// exit exits the program
var exit g.Value = g.NewFixedNativeFunc(
	[]g.Type{g.IntType}, false,
	func(ev g.Eval, params []g.Value) (g.Value, g.Error) {
		n := params[0].(g.Int)

		os.Exit(int(n.IntVal()))

		// we will never actually get here
		return g.Null, nil
	})

/*doc
### `open`

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

*/

// open opens a file
var open g.Value = g.NewFixedNativeFunc(
	[]g.Type{g.StrType}, false,
	func(ev g.Eval, params []g.Value) (g.Value, g.Error) {
		s := params[0].(g.Str)

		f, err := os.Open(s.String())
		if err != nil {
			return nil, g.Error(fmt.Errorf("OsError: %s", err.Error()))
		}
		return newFile(f), nil
	})

/*doc
### `stat`

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

*/

// stat stats a file
var stat g.Value = g.NewFixedNativeFunc(
	[]g.Type{g.StrType}, false,
	func(ev g.Eval, params []g.Value) (g.Value, g.Error) {
		s := params[0].(g.Str)

		// TODO os.Lstat
		// followSymLink == false, same as os.Lstat
		//
		//followSymLink := g.True
		//if len(params) == 2 {
		//	var ok bool
		//	followSymLink, ok = params[1].(g.Bool)
		//	if !ok {
		//		return nil, g.TypeMismatch("Expected Bool")
		//	}
		//}
		//var fn = os.Stat
		//if !followSymLink.BoolVal() {
		//	fn = os.Lstat
		//}

		info, err := os.Stat(s.String())
		if err != nil {
			return nil, g.Error(fmt.Errorf("OsError: %s", err.Error()))
		}
		return NewFileInfo(info), nil
	})

//-------------------------------------------------------------------------

/*doc

## Structs

*/

/*doc
### `fileInfo`

A `fileInfo` is a struct that describes a file and is returned by stat.

A `fileInfo` struct has the fields:

* [isDir](#isdir)
* [mode](#mode)
* [name](#name)
* [size](#size)

*/

// NewFileInfo creates a struct for 'os.FileInfo'
func NewFileInfo(info os.FileInfo) g.Struct {
	stc, err := g.NewMethodStruct(info, fileInfoMethods)
	g.Assert(err == nil)
	return stc
}

var fileInfoMethods = map[string]g.Method{

	/*doc
	#### `isDir`

	`isDir` is an abbreviation for Mode().IsDir()

	* signature: `isDir() <Str>`

	*/

	"isDir": g.NewWrapperMethod(
		func(self interface{}) g.Value {
			info := self.(os.FileInfo)
			return g.NewBool(info.IsDir())
		}),

	/*doc
	#### `mode`

	`mode` is the file mode bits

	* signature: `mode() <Int>`

	*/

	"mode": g.NewWrapperMethod(
		func(self interface{}) g.Value {
			info := self.(os.FileInfo)
			return g.NewInt(int64(info.Mode()))
		}),

	// TODO
	// g.NewField("modTime", true, ModTime() time.Time

	/*doc
	#### `name`

	`name` is the base name of the file

	* signature: `name() <Str>`

	*/

	"name": g.NewNullaryMethod(
		func(self interface{}, ev g.Eval) (g.Value, g.Error) {
			info := self.(os.FileInfo)

			// TODO: this means that the the file name
			// must be valid UTF-8.  Is that what we really want?
			return g.NewStr(info.Name())
		}),

	/*doc
	#### `size`

	`size` is the length in bytes for regular files; system-dependent for others

	* signature: `size() <Int>`

	*/

	"size": g.NewWrapperMethod(
		func(self interface{}) g.Value {
			info := self.(os.FileInfo)
			return g.NewInt(info.Size())
		}),
}

//-------------------------------------------------------------------------

/*doc
### `file`

A `file` is a struct that represents an open file descriptor.

A `file` struct has the fields:

* [close](#close)
* [readLines](#readlines)
* [writeLines](#writelines)

*/

func newFile(f *os.File) g.Struct {
	stc, err := g.NewMethodStruct(f, fileMethods)
	g.Assert(err == nil)
	return stc
}

var fileMethods = map[string]g.Method{

	/*doc
	#### `close`

	`close` closes the File, rendering it unusable for I/O. On files that support
	SetDeadline, any pending I/O operations will be canceled and return immediately
	with an error.

	* signature: `close() <Null>`

	*/
	"close": g.NewNullaryMethod(
		func(self interface{}, ev g.Eval) (g.Value, g.Error) {
			f := self.(*os.File)
			err := closeFile(f)
			if err != nil {
				return nil, err
			}
			return g.Null, nil
		}),

	/*doc
	#### `readLines`

	`readLines` returns a List of Strs, for each line of text in the file.

	* signature: `readLines() <List>`

	*/
	"readLines": g.NewNullaryMethod(
		func(self interface{}, ev g.Eval) (g.Value, g.Error) {
			f := self.(*os.File)
			return readLines(f)
		}),

	/*doc
	#### `writeLines`

	`writeLines` writes a List of Strs to the file as a sequence of lines.

	* signature: `writeLines(<List>) <Null>`

	*/
	"writeLines": g.NewFixedMethod(
		[]g.Type{g.ListType}, true,
		func(self interface{}, ev g.Eval, params []g.Value) (g.Value, g.Error) {
			f := self.(*os.File)
			lines := params[0].(g.List)
			return g.Null, writeLines(ev, f, lines)
		}),
}

func readLines(f io.Reader) (g.List, g.Error) {

	lines := []g.Value{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {

		// TODO: this means that the the file being read must consist entirely of
		// valid UTF-8.  Is that what we really want?
		s, e := g.NewStr(scanner.Text())
		if e != nil {
			return nil, e
		}

		lines = append(lines, s)
	}

	if e := scanner.Err(); e != nil {
		return nil, g.Error(fmt.Errorf("OsError: %s", e.Error()))
	}

	return g.NewList(lines), nil
}

func writeLines(ev g.Eval, f *os.File, lines g.List) g.Error {

	for _, v := range lines.Values() {
		s, err := v.ToStr(ev)
		if err != nil {
			return err
		}
		_, e := f.WriteString(s.String())
		if e != nil {
			return g.Error(fmt.Errorf("OsError: %s", e.Error()))
		}
		_, e = f.WriteString("\n")
		if e != nil {
			return g.Error(fmt.Errorf("OsError: %s", e.Error()))
		}
	}

	return nil
}

func closeFile(f io.Closer) g.Error {
	e := f.Close()
	if e != nil {
		return g.Error(fmt.Errorf("OsError: %s", e.Error()))
	}
	return nil
}
