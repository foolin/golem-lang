// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"io"
	"os"

	g "github.com/mjarmy/golem-lang/core"
	osutil "github.com/mjarmy/golem-lang/lib/os/util"
)

// Exit exits the program
var Exit g.Value = g.NewNativeFuncInt(
	func(cx g.Context, n g.Int) (g.Value, g.Error) {

		os.Exit(int(n.IntVal()))

		// we will never actually get here
		return g.Null, nil
	})

// Open opens a file
var Open g.Value = g.NewNativeFuncStr(
	func(cx g.Context, s g.Str) (g.Value, g.Error) {

		f, err := os.Open(s.String())
		if err != nil {
			return nil, g.NewError("OsError", err.Error())
		}
		return newFile(f), nil
	})

// Stat stats a file
var Stat g.Value = g.NewNativeFuncStr(
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
		return osutil.NewInfo(info), nil
	})

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
