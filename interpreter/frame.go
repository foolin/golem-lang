// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package interpreter

import (
	"fmt"

	g "github.com/mjarmy/golem-lang/core"
	"github.com/mjarmy/golem-lang/core/bytecode"
)

//---------------------------------------------------------------
// An execution environment, a.k.a 'stack frame'.
//---------------------------------------------------------------

type frame struct {
	fn     bytecode.BytecodeFunc
	locals []*bytecode.Ref
	stack  []g.Value
	ip     int
}

func mustStr(val g.Value) g.Str {
	s, err := val.ToStr(nil)
	if err != nil {
		return g.NewStr("ToStr ERROR: " + err.Error())
	}
	return s
}

func dumpFrames(frames []*frame) {

	println("-----------------------------------------")

	f := frames[len(frames)-1]
	opc := f.fn.Template().Bytecodes

	fmt.Printf("%s\n", bytecode.FmtBytecode(opc, f.ip))

	for j, f := range frames {
		fmt.Printf("frame %d\n", j)
		dumpFrame(f)
	}
}

func dumpFrame(f *frame) {

	fmt.Printf("    locals:\n")
	for j, r := range f.locals {
		fmt.Printf("        %d: %s\n", j, mustStr(r.Val))
	}

	fmt.Printf("    stack:\n")
	for j, v := range f.stack {
		fmt.Printf("        %d: %s\n", j, mustStr(v))
	}

	fmt.Printf("    ip: %d\n", f.ip)
}
