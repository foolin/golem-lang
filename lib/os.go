// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package lib

import (
	"bufio"
	g "github.com/mjarmy/golem-lang/core"
	"io"
	"os"
)

type osModule struct {
	contents g.Struct
}

func (m *osModule) GetModuleName() string {
	return "os"
}

func (m *osModule) GetContents() g.Struct {
	return m.contents
}

// NewOsModule creates the 'os' module.
func NewOsModule() g.Module {

	contents, err := g.NewStruct([]g.Field{
		g.NewField("open", true, open()),
		g.NewField("exit", true, exit())},
		true)

	if err != nil {
		panic("NewOsModule")
	}

	return &osModule{contents}
}

func open() g.NativeFunc {

	return g.NewNativeFunc(
		1, 1,
		func(cx g.Context, values []g.Value) (g.Value, g.Error) {
			s, ok := values[0].(g.Str)
			if !ok {
				return nil, g.TypeMismatchError("Expected Str")
			}

			f, err := os.Open(s.String())
			if err != nil {
				return nil, g.NewError("OsError", err.Error())
			}
			return makeFile(f), nil
		})

}

func exit() g.NativeFunc {

	return g.NewNativeFunc(
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
				return nil, g.ArityMismatchError("0 or 1", len(values))
			}

			// we will never actually get here
			return g.NullValue, nil
		})

}

//-------------------------------------------------------------------------

func makeFile(f *os.File) g.Struct {

	file, err := g.NewStruct([]g.Field{
		g.NewField("readLines", true, readLines(f)),
		g.NewField("close", true, close(f))},
		true)
	if err != nil {
		panic("NewOsModule")
	}

	return file
}

func readLines(f io.Reader) g.NativeFunc {
	return g.NewNativeFunc(
		0, 0,
		func(cx g.Context, values []g.Value) (g.Value, g.Error) {

			lines := []g.Value{}
			scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				lines = append(lines, g.NewStr(scanner.Text()))
			}

			if err := scanner.Err(); err != nil {
				return nil, g.NewError("OsError", err.Error())
			}

			return g.NewList(lines), nil
		})
}

func close(f io.Closer) g.NativeFunc {
	return g.NewNativeFunc(
		0, 0,
		func(cx g.Context, values []g.Value) (g.Value, g.Error) {
			err := f.Close()
			if err != nil {
				return nil, g.NewError("OsError", err.Error())
			}
			return g.NullValue, nil
		})
}
