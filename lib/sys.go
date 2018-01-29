// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package lib

import (
	g "github.com/mjarmy/golem-lang/core"
	"os"
)

type sysModule struct {
	contents g.Struct
}

// InitSysModule initializes the 'sys' module.
func InitSysModule() g.Module {

	exit := g.NewNativeFunc(
		0, 1,
		func(cx g.Context, values []g.Value) (g.Value, g.Error) {
			switch len(values) {
			case 0:
				os.Exit(0)
			case 1:
				if n, ok := values[0].(g.Int); ok {
					os.Exit(int(n.IntVal()))
				} else {
					return nil, g.TypeMismatchError("Expected Int")
				}
			default:
				panic("arity mismatch")
			}

			// we will never actually get here
			return g.NULL, nil
		})

	contents, err := g.NewStruct([]g.Field{g.NewField("exit", true, exit)}, true)
	if err != nil {
		panic("InitSysModule")
	}

	return &sysModule{contents}
}

func (m *sysModule) GetModuleName() string {
	return "sys"
}

func (m *sysModule) GetContents() g.Struct {
	return m.contents
}
