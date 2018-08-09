// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

//import (
//	"bufio"
//	"io"
//	"os"
//
//	g "github.com/mjarmy/golem-lang/core"
//	osutil "github.com/mjarmy/golem-lang/lib/os/util"
//)
//
//type module struct{ contents g.Struct }
//
//func (m *module) GetModuleName() string { return "os" }
//func (m *module) GetContents() g.Struct { return m.contents }
//
//// LoadModule creates the 'os' module.
//func LoadModule() (g.Module, g.Error) {
//
//	contents, err := g.NewStruct([]g.Field{
//		g.NewField("exit", true, exit()),
//		g.NewField("open", true, open()),
//		g.NewField("stat", true, stat()),
//	}, true)
//	if err != nil {
//		return nil, err
//	}
//
//	return &module{contents}, nil
//}
//
//func exit() g.NativeFunc {
//
//	return g.NewNativeFunc(
//		0, 1,
//		func(cx g.Context, values []g.Value) (g.Value, g.Error) {
//			switch len(values) {
//			case 0:
//				os.Exit(0)
//			case 1:
//				if n, ok := values[0].(g.Int); ok {
//					os.Exit(int(n.IntVal()))
//				} else {
//					return nil, g.TypeMismatchError("Expected Int")
//				}
//			default:
//				return nil, g.ArityMismatchError("0 or 1", len(values))
//			}
//
//			// we will never actually get here
//			return g.Null, nil
//		})
//}
//
//func open() g.NativeFunc {
//
//	return g.NewNativeFuncStr(
//		func(cx g.Context, s g.Str) (g.Value, g.Error) {
//
//			f, err := os.Open(s.String())
//			if err != nil {
//				return nil, g.NewError("OsError", err.Error())
//			}
//			return newFile(f), nil
//		})
//}
//
//func stat() g.NativeFunc {
//
//	return g.NewNativeFunc(
//		//1, 2,
//		1, 1,
//		func(cx g.Context, values []g.Value) (g.Value, g.Error) {
//
//			name, ok := values[0].(g.Str)
//			if !ok {
//				return nil, g.TypeMismatchError("Expected Str")
//			}
//
//			// TODO os.Lstat
//			// TODO followSymLink == false, same as os.Lstat
//			//
//			//followSymLink := g.True
//			//if len(values) == 2 {
//			//	var ok bool
//			//	followSymLink, ok = values[1].(g.Bool)
//			//	if !ok {
//			//		return nil, g.TypeMismatchError("Expected Bool")
//			//	}
//			//}
//			//var fn = os.Stat
//			//if !followSymLink.BoolVal() {
//			//	fn = os.Lstat
//			//}
//
//			info, err := os.Stat(name.String())
//			if err != nil {
//				return nil, g.NewError("OsError", err.Error())
//			}
//			return osutil.NewInfo(info), nil
//		})
//}
//
////-------------------------------------------------------------------------
//
//func newFile(f *os.File) g.Struct {
//
//	stc, err := g.NewStruct([]g.Field{
//		g.NewField("readLines", true, readLines(f)),
//		g.NewField("close", true, close(f)),
//	}, true)
//	if err != nil {
//		panic("unreachable")
//	}
//
//	return stc
//}
//
//func readLines(f io.Reader) g.NativeFunc {
//	return g.NewNativeFunc0(
//		func(cx g.Context) (g.Value, g.Error) {
//
//			lines := []g.Value{}
//			scanner := bufio.NewScanner(f)
//			for scanner.Scan() {
//				lines = append(lines, g.NewStr(scanner.Text()))
//			}
//
//			if err := scanner.Err(); err != nil {
//				return nil, g.NewError("OsError", err.Error())
//			}
//
//			return g.NewList(lines), nil
//		})
//}
//
//func close(f io.Closer) g.NativeFunc {
//	return g.NewNativeFunc0(
//		func(cx g.Context) (g.Value, g.Error) {
//			err := f.Close()
//			if err != nil {
//				return nil, g.NewError("OsError", err.Error())
//			}
//			return g.Null, nil
//		})
//}
