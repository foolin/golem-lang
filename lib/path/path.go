// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"os"
	"path/filepath"

	g "github.com/mjarmy/golem-lang/core"
	osutil "github.com/mjarmy/golem-lang/lib/os/util"
)

// Ext returns a file extension
var Ext g.Value = g.NewNativeFuncStr(
	func(cx g.Context, name g.Str) (g.Value, g.Error) {
		return g.NewStr(filepath.Ext(name.String())), nil
	})

// Walk walks a directory path
var Walk g.Value = g.NewNativeFunc(
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
					[]g.Value{g.NewStr(path), osutil.NewFileInfo(info)})
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
