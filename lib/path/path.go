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
	Path, err = g.NewFieldStruct(
		map[string]g.Field{
			"ext":  g.NewField(ext),
			"walk": g.NewField(walk),
		}, true)
	g.Assert(err == nil)
}

// ext returns a file extension
var ext g.Value = g.NewFixedNativeFunc(
	[]g.Type{g.StrType}, false,
	func(ev g.Eval, params []g.Value) (g.Value, g.Error) {
		s := params[0].(g.Str)
		return g.NewStr(filepath.Ext(s.String())), nil
	})

// walk walks a directory path
var walk g.Value = g.NewFixedNativeFunc(
	[]g.Type{g.StrType, g.FuncType}, false,
	func(ev g.Eval, params []g.Value) (g.Value, g.Error) {

		dir := params[0].(g.Str)
		callback := params[1].(g.Func)

		arity := callback.Arity()
		if arity.Kind != g.FixedArity || arity.RequiredParams != 2 {
			return nil, g.ArityMismatchError("2", int(arity.RequiredParams))
		}

		err := filepath.Walk(
			dir.String(),
			func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				_, gerr := callback.Invoke(ev,
					[]g.Value{g.NewStr(path), libOs.NewFileInfo(info)})
				return gerr
			})

		if err != nil {
			if gerr, ok := err.(g.Error); ok {
				return nil, gerr
			}
			return nil, g.NewError("PathError: " + err.Error())
		}
		return g.Null, nil
	})
