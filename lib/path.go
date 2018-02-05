// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package lib

import (
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
