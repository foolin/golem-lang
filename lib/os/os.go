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
			"exit": g.NewField(exit),
			"open": g.NewField(open),
			"stat": g.NewField(stat),
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
		//		return nil, g.TypeMismatchError("Expected Bool")
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

	"name": g.NewWrapperMethod(
		func(self interface{}) g.Value {
			info := self.(os.FileInfo)
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

	"close": g.NewNullaryMethod(
		func(self interface{}, ev g.Eval) (g.Value, g.Error) {
			f := self.(*os.File)
			err := close(f)
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
		lines = append(lines, g.NewStr(scanner.Text()))
	}

	if err := scanner.Err(); err != nil {
		return nil, g.Error(fmt.Errorf("OsError: %s", err.Error()))
	}

	return g.NewList(lines), nil
}

func close(f io.Closer) g.Error {
	err := f.Close()
	if err != nil {
		return g.Error(fmt.Errorf("OsError: %s", err.Error()))
	}
	return nil
}
