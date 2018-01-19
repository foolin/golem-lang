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
	"bufio"
	"fmt"
	g "github.com/mjarmy/golem-lang/core"
	"io/ioutil"
	"os"
)

type ioModule struct {
	contents g.Struct
}

func InitIoModule() g.Module {

	file := g.NewNativeFunc(
		1, 1,
		func(cx g.Context, values []g.Value) (g.Value, g.Error) {
			s, ok := values[0].(g.Str)
			if !ok {
				return nil, g.TypeMismatchError("Expected Str")
			}

			return makeFile(s), nil
		})

	contents, err := g.NewStruct([]g.Field{
		g.NewField("File", true, file)}, true)
	if err != nil {
		panic("InitIoModule")
	}

	return &ioModule{contents}
}

func makeFile(name g.Str) g.Struct {

	isDir := g.NewNativeFunc(
		0, 0,
		func(cx g.Context, values []g.Value) (g.Value, g.Error) {
			fi, err := os.Stat(name.String())
			if err != nil {
				return nil, g.MakeError("IoError", err.Error())
			}
			isDir := false
			if mode := fi.Mode(); mode.IsDir() {
				isDir = true
			}

			return g.MakeBool(isDir), nil
		})

	items := g.NewNativeFunc(

		0, 0,
		func(cx g.Context, values []g.Value) (g.Value, g.Error) {
			files, err := ioutil.ReadDir(name.String())
			if err != nil {
				return nil, g.MakeError("IoError", err.Error())
			}

			sep := fmt.Sprintf("%c", os.PathSeparator)
			list := []g.Value{}
			for _, f := range files {
				itemName := name.String() + sep + f.Name()
				list = append(list, makeFile(g.MakeStr(itemName)))
			}

			return g.NewList(list), nil
		})

	readLines := g.NewNativeFunc(
		0, 0,
		func(cx g.Context, values []g.Value) (g.Value, g.Error) {
			f, err := os.Open(name.String())
			if err != nil {
				return nil, g.MakeError("IoError", err.Error())
			}
			defer f.Close()

			scanner := bufio.NewScanner(f)
			list := []g.Value{}
			for scanner.Scan() {
				list = append(list, g.MakeStr(scanner.Text()))
			}

			if err := scanner.Err(); err != nil {
				return nil, g.MakeError("IoError", err.Error())
			}

			return g.NewList(list), nil
		})

	file, err := g.NewStruct([]g.Field{
		g.NewField("isDir", true, isDir),
		g.NewField("items", true, items),
		g.NewField("readLines", true, readLines),
		g.NewField("name", true, name)}, true)
	if err != nil {
		panic("InitIoModule")
	}

	return file
}

func (m *ioModule) GetModuleName() string {
	return "io"
}

func (m *ioModule) GetContents() g.Struct {
	return m.contents
}
