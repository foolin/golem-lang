// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package interpreter

import (
	//"fmt"

	g "github.com/mjarmy/golem-lang/core"
	bc "github.com/mjarmy/golem-lang/core/bytecode"
)

//---------------------------------------------------------------
// An execution environment, a.k.a 'stack frame'.
//---------------------------------------------------------------

type frame struct {
	fn     bc.Func
	locals []*bc.Ref
	stack  []g.Value
	ip     int
}

//func dumpStr(val g.Value) string {
//	s, err := val.ToStr(nil)
//	if err != nil {
//		return "ToStr ERROR: " + err.Error()
//	}
//	return s.String()
//}
//
//func dumpFrames(frames []*frame) {
//
//	fmt.Printf("-----------------------------------------\n")
//
//	f := frames[len(frames)-1]
//	opc := f.fn.Template().Bytecodes
//
//	fmt.Printf("%s\n", bc.FmtBytecode(opc, f.ip))
//
//	for j, f := range frames {
//		fmt.Printf("frame %d\n", j)
//		dumpFrame(f)
//	}
//}
//
//func dumpFrame(f *frame) {
//
//	fmt.Printf("    locals:\n")
//	for j, r := range f.locals {
//		fmt.Printf("        %d: %s\n", j, dumpStr(r.Val))
//	}
//
//	fmt.Printf("    stack:\n")
//	for j, v := range f.stack {
//		fmt.Printf("        %d: %s\n", j, dumpStr(v))
//	}
//
//	fmt.Printf("    ip: %d\n", f.ip)
//}
