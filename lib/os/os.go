// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package os

import (
	"bufio"
	"io"
	"os"

	g "github.com/mjarmy/golem-lang/core"
)

// Os is the "os" module in the standard library
var Os g.Struct

func init() {
	var err error
	Os, err = g.NewStruct([]g.Field{
		g.NewField("exit", true, exit),
		g.NewField("open", true, open),
		g.NewField("stat", true, stat),
	}, true)
	if err != nil {
		panic("unreachable")
	}
}

// exit exits the program
var exit g.Value = g.NewNativeFuncInt(
	func(cx g.Context, n g.Int) (g.Value, g.Error) {

		os.Exit(int(n.IntVal()))

		// we will never actually get here
		return g.Null, nil
	})

// open opens a file
var open g.Value = g.NewNativeFuncStr(
	func(cx g.Context, s g.Str) (g.Value, g.Error) {

		f, err := os.Open(s.String())
		if err != nil {
			return nil, g.NewError("OsError", err.Error())
		}
		return newFile(f), nil
	})

// stat stats a file
var stat g.Value = g.NewNativeFuncStr(
	func(cx g.Context, s g.Str) (g.Value, g.Error) {

		// TODO os.Lstat
		// TODO followSymLink == false, same as os.Lstat
		//
		//followSymLink := g.True
		//if len(values) == 2 {
		//	var ok bool
		//	followSymLink, ok = values[1].(g.Bool)
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
			return nil, g.NewError("OsError", err.Error())
		}
		return NewFileInfo(info), nil
	})

//-------------------------------------------------------------------------

// NewFileInfo creates a struct for 'os.FileInfo'
func NewFileInfo(info os.FileInfo) g.Struct {

	stc, err := g.NewStruct([]g.Field{
		g.NewField("name", true, g.NewStr(info.Name())),
		g.NewField("size", true, g.NewInt(info.Size())),
		g.NewField("mode", true, g.NewInt(int64(info.Mode()))),
		//g.NewField("modTime", true, ModTime() time.Time TODO
		g.NewField("isDir", true, g.NewBool(info.IsDir())),
	}, true)
	if err != nil {
		panic("unreachable")
	}

	return stc
}

//-------------------------------------------------------------------------

func newFile(f *os.File) g.Struct {

	stc, err := g.NewStruct([]g.Field{
		g.NewField("readLines", true, readLines(f)),
		g.NewField("close", true, close(f)),
	}, true)
	if err != nil {
		panic("unreachable")
	}

	return stc
}

func readLines(f io.Reader) g.NativeFunc {
	return g.NewNativeFunc0(
		func(cx g.Context) (g.Value, g.Error) {

			lines := []g.Value{}
			scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				lines = append(lines, g.NewStr(scanner.Text()))
			}

			if err := scanner.Err(); err != nil {
				return nil, g.NewError("OsError", err.Error())
			}

			return g.NewList(lines), nil
		})
}

func close(f io.Closer) g.NativeFunc {
	return g.NewNativeFunc0(
		func(cx g.Context) (g.Value, g.Error) {
			err := f.Close()
			if err != nil {
				return nil, g.NewError("OsError", err.Error())
			}
			return g.Null, nil
		})
}
