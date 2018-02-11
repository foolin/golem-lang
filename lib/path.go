// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package lib

import (
	"os"
	"path/filepath"

	g "github.com/mjarmy/golem-lang/core"
)

type pathModule struct {
	contents g.Struct
}

func (m *pathModule) GetModuleName() string {
	return "path"
}

func (m *pathModule) GetContents() g.Struct {
	return m.contents
}

// NewPathModule creates the 'path' module.
func NewPathModule() g.Module {

	contents, err := g.NewStruct([]g.Field{
		g.NewField("filepath", true, newFilepath()),
	}, true)

	if err != nil {
		panic("unreachable")
	}

	return &pathModule{contents}
}

func newFilepath() g.Struct {

	stc, err := g.NewStruct([]g.Field{
		g.NewField("ext", true, ext()),
		g.NewField("walk", true, walk()),
	}, true)
	if err != nil {
		panic("unreachable")
	}

	return stc
}

func ext() g.NativeFunc {

	return g.NewNativeFunc(
		1, 1,
		func(cx g.Context, values []g.Value) (g.Value, g.Error) {

			name, ok := values[0].(g.Str)
			if !ok {
				return nil, g.TypeMismatchError("Expected Str")
			}

			return g.NewStr(filepath.Ext(name.String())), nil
		})
}

func walk() g.NativeFunc {

	return g.NewNativeFunc(
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
						[]g.Value{g.NewStr(path), newInfo(info)})
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
}
