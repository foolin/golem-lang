// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

//import (
//	"os"
//	"path/filepath"
//
//	g "github.com/mjarmy/golem-lang/core"
//	osutil "github.com/mjarmy/golem-lang/lib/os/util"
//)
//
//type module struct{ contents g.Struct }
//
//func (m *module) GetModuleName() string { return "path" }
//func (m *module) GetContents() g.Struct { return m.contents }
//
//// LoadModule creates the 'path' module.
//func LoadModule() (g.Module, g.Error) {
//
//	contents, err := g.NewStruct([]g.Field{
//		g.NewField("filepath", true, newFilepath()),
//	}, true)
//	if err != nil {
//		return nil, err
//	}
//
//	return &module{contents}, nil
//}
//
//func newFilepath() g.Struct {
//
//	stc, err := g.NewStruct([]g.Field{
//		g.NewField("ext", true, ext()),
//		g.NewField("walk", true, walk()),
//	}, true)
//	if err != nil {
//		panic("unreachable")
//	}
//
//	return stc
//}
//
//func ext() g.NativeFunc {
//
//	return g.NewNativeFuncStr(
//		func(cx g.Context, name g.Str) (g.Value, g.Error) {
//			return g.NewStr(filepath.Ext(name.String())), nil
//		})
//}
//
//func walk() g.NativeFunc {
//
//	return g.NewNativeFunc(
//		2, 2,
//		func(cx g.Context, values []g.Value) (g.Value, g.Error) {
//
//			dir, ok := values[0].(g.Str)
//			if !ok {
//				return nil, g.TypeMismatchError("Expected Str")
//			}
//
//			callback, ok := values[1].(g.Func)
//			if !ok {
//				return nil, g.TypeMismatchError("Expected Func")
//			}
//			if callback.MinArity() != 2 || callback.MaxArity() != 2 {
//				return nil, g.ArityMismatchError("2", callback.MinArity())
//			}
//
//			err := filepath.Walk(
//				dir.String(),
//				func(path string, info os.FileInfo, err error) error {
//					if err != nil {
//						return err
//					}
//					_, gerr := callback.Invoke(cx,
//						[]g.Value{g.NewStr(path), osutil.NewInfo(info)})
//					return gerr
//				})
//
//			if err != nil {
//				if gerr, ok := err.(g.Error); ok {
//					return nil, gerr
//				}
//				return nil, g.NewError("PathError", err.Error())
//			}
//			return g.Null, nil
//		})
//}
