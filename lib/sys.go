// Copyright 2017 The Golem Project Developers
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package lib

import (
	g "github.com/mjarmy/golem-lang/core"
	"os"
)

type sysModule struct {
	contents g.Struct
}

func InitSysModule() g.Module {

	exit := g.NewNativeFunc(
		func(values []g.Value) (g.Value, g.Error) {
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
				return nil, g.ArityMismatchError("0 or 1", len(values))
			}

			// we will never actually get here
			return g.NULL, nil
		})

	contents, err := g.NewStruct([]*g.StructEntry{
		{"exit", true, false, exit}})
	g.Assert(err == nil, "InitSysModule")

	return &sysModule{contents}
}

func (m *sysModule) GetModuleName() string {
	return "sys"
}

func (m *sysModule) GetContents() g.Struct {
	return m.contents
}
