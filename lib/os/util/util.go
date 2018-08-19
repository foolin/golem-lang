// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package util

import (
	g "github.com/mjarmy/golem-lang/core"
	"os"
)

// NewInfo creates a struct for 'os.FileInfo'
func NewInfo(info os.FileInfo) g.Struct {

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
