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
)

// Os is the "os" module in the standard library
var Os g.Struct

func init() {
	var err error
	Os, err = g.NewFrozenFieldStruct(
		map[string]g.Field{
			"exit":   g.NewField(exit),
			"open":   g.NewField(open),
			"create": g.NewField(create),
			"stat":   g.NewField(stat),
		})
	g.Assert(err == nil)
}

// exit exits the program
var exit g.Value = g.NewFixedNativeFunc(
	[]g.Type{g.IntType}, false,
	func(ev g.Eval, params []g.Value) (g.Value, g.Error) {
		n := params[0].(g.Int)

		os.Exit(int(n.IntVal()))

		// we will never actually get here
		return g.Null, nil
	})

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

// stat stats a file
var stat g.Value = g.NewFixedNativeFunc(
	[]g.Type{g.StrType}, false,
	func(ev g.Eval, params []g.Value) (g.Value, g.Error) {
		s := params[0].(g.Str)

		// TODO os.Lstat
		// TODO followSymLink == false, same as os.Lstat
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

// NewFileInfo creates a struct for 'os.FileInfo'
func NewFileInfo(info os.FileInfo) g.Struct {
	stc, err := g.NewMethodStruct(info, fileInfoMethods)
	g.Assert(err == nil)
	return stc
}

var fileInfoMethods = map[string]g.Method{

	"name": g.NewNullaryMethod(
		func(self interface{}, ev g.Eval) (g.Value, g.Error) {
			info := self.(os.FileInfo)

			// TODO: this means that the the file name
			// must be valid UTF-8.  Is that what we really want?
			return g.NewStr(info.Name())
		}),

	"size": g.NewWrapperMethod(
		func(self interface{}) g.Value {
			info := self.(os.FileInfo)
			return g.NewInt(info.Size())
		}),

	"mode": g.NewWrapperMethod(
		func(self interface{}) g.Value {
			info := self.(os.FileInfo)
			return g.NewInt(int64(info.Mode()))
		}),

	//		//g.NewField("modTime", true, ModTime() time.Time TODO

	"isDir": g.NewWrapperMethod(
		func(self interface{}) g.Value {
			info := self.(os.FileInfo)
			return g.NewBool(info.IsDir())
		}),
}

//-------------------------------------------------------------------------

func newFile(f *os.File) g.Struct {
	stc, err := g.NewMethodStruct(f, fileMethods)
	g.Assert(err == nil)
	return stc
}

var fileMethods = map[string]g.Method{

	"readLines": g.NewNullaryMethod(
		func(self interface{}, ev g.Eval) (g.Value, g.Error) {
			f := self.(*os.File)
			return readLines(f)
		}),

	"writeLines": g.NewFixedMethod(
		[]g.Type{g.ListType}, true,
		func(self interface{}, ev g.Eval, params []g.Value) (g.Value, g.Error) {
			f := self.(*os.File)
			lines := params[0].(g.List)
			return g.Null, writeLines(ev, f, lines)
		}),

	"close": g.NewNullaryMethod(
		func(self interface{}, ev g.Eval) (g.Value, g.Error) {
			f := self.(*os.File)
			err := closeFile(f)
			if err != nil {
				return nil, err
			}
			return g.Null, nil
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
