// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

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

// InitIoModule initializes the 'io' module.
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
				return nil, g.NewError("IoError", err.Error())
			}
			isDir := false
			if mode := fi.Mode(); mode.IsDir() {
				isDir = true
			}

			return g.NewBool(isDir), nil
		})

	items := g.NewNativeFunc(

		0, 0,
		func(cx g.Context, values []g.Value) (g.Value, g.Error) {
			files, err := ioutil.ReadDir(name.String())
			if err != nil {
				return nil, g.NewError("IoError", err.Error())
			}

			sep := fmt.Sprintf("%c", os.PathSeparator)
			list := []g.Value{}
			for _, f := range files {
				itemName := name.String() + sep + f.Name()
				list = append(list, makeFile(g.NewStr(itemName)))
			}

			return g.NewList(list), nil
		})

	readLines := g.NewNativeFunc(
		0, 0,
		func(cx g.Context, values []g.Value) (g.Value, g.Error) {
			f, err := os.Open(name.String())
			if err != nil {
				return nil, g.NewError("IoError", err.Error())
			}
			defer f.Close() // nolint: errcheck

			scanner := bufio.NewScanner(f)
			list := []g.Value{}
			for scanner.Scan() {
				list = append(list, g.NewStr(scanner.Text()))
			}

			if err := scanner.Err(); err != nil {
				return nil, g.NewError("IoError", err.Error())
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
