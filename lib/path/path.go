// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package path

import (
	"os"
	"path/filepath"

	g "github.com/mjarmy/golem-lang/core"
	libOs "github.com/mjarmy/golem-lang/lib/os"
)

// Path is the "path" module in the standard library
var Path g.Struct

func init() {
	var err error
	Path, err = g.NewStruct([]g.Field{
		g.NewField("ext", true, ext),
		g.NewField("walk", true, walk),
	}, true)
	if err != nil {
		panic("unreachable")
	}
}

// ext returns a file extension
var ext g.Value = g.NewNativeFuncStr(
	func(cx g.Context, name g.Str) (g.Value, g.Error) {
		return g.NewStr(filepath.Ext(name.String())), nil
	})

// wal walks a directory path
var walk g.Value = g.NewNativeFunc(
	2, 2,
	func(cx g.Context, values []g.Value) (g.Value, g.Error) {

		dir, ok := values[0].(g.Str)
		if !ok {
			return nil, g.TypeMismatchError("Expected Str")
		}

		callback, ok := values[1].(g.Func)
		if !ok {
			return nil, g.TypeMismatchError("Expected Func")
		}
		if callback.MinArity() != 2 || callback.MaxArity() != 2 {
			return nil, g.ArityMismatchError("2", callback.MinArity())
		}

		err := filepath.Walk(
			dir.String(),
			func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				_, gerr := callback.Invoke(cx,
					[]g.Value{g.NewStr(path), libOs.NewFileInfo(info)})
				return gerr
			})

		if err != nil {
			if gerr, ok := err.(g.Error); ok {
				return nil, gerr
			}
			return nil, g.NewError("PathError", err.Error())
		}
		return g.Null, nil
	})
