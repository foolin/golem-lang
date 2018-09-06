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

// Path is the "path" module in the standard library
var Path g.Struct

func init() {

	filepath, err := g.NewFrozenFieldStruct(
		map[string]g.Field{
			"ext":  g.NewField(ext),
			"walk": g.NewField(walk),
		})
	g.Assert(err == nil)

	Path, err = g.NewFrozenFieldStruct(
		map[string]g.Field{
			"filepath": g.NewField(filepath),
		})
	g.Assert(err == nil)
}

// ext returns a file extension
var ext g.Value = g.NewFixedNativeFunc(
	[]g.Type{g.StrType}, false,
	func(ev g.Eval, params []g.Value) (g.Value, g.Error) {
		s := params[0].(g.Str)

		// TODO: this means that the the path must be
		// valid UTF-8.  Is that what we really want?
		return g.NewStr(filepath.Ext(s.String()))
	})

// walk walks a directory path
var walk g.Value = g.NewFixedNativeFunc(
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
