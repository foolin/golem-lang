// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package bytecode

import (
	"bytes"
	"fmt"

	g "github.com/mjarmy/golem-lang/core"
)

// Pool is a pool of the constants, function templates, and struct definitions
// used by a given Module.  Pools are created at compile time, and
// are immutable at run time.
type Pool struct {
	Constants  []g.Basic
	StructDefs [][]string
	Templates  []*FuncTemplate
}

func (p *Pool) String() string {

	var buf bytes.Buffer

	buf.WriteString("Constants:\n")
	for i, b := range p.Constants {
		s, err := b.ToStr(nil)
		g.Assert(err == nil)
		buf.WriteString(fmt.Sprintf("    %d: %v\n", i, s))
	}

	buf.WriteString("StructDefs:\n")
	for i, d := range p.StructDefs {
		buf.WriteString(fmt.Sprintf("    %d: %v\n", i, d))
	}

	buf.WriteString("Templates:\n")
	for i, t := range p.Templates {
		buf.WriteString(fmt.Sprintf("    %d: Template\n", i))
		buf.WriteString(fmt.Sprintf("        Arity: %s\n", t.Arity))
		buf.WriteString(fmt.Sprintf("        NumCaptures: %d\n", t.NumCaptures))
		buf.WriteString(fmt.Sprintf("        NumLocals: %d\n", t.NumLocals))

		//buf.WriteString(fmt.Sprintf("        Opcodes:\n"))
		//i := 0
		//for i < len(t.Opcodes) {
		//	buf.WriteString("            ")
		//	buf.WriteString(FmtOpcode(t.Opcodes, i))
		//	i += OpcodeSize(t.Opcodes[i])
		//}

		buf.WriteString(fmt.Sprintf("        LineNumberTable:\n"))
		for j, ln := range t.LineNumberTable {
			buf.WriteString(fmt.Sprintf("            %d: %v\n", j, ln))
		}

		buf.WriteString(fmt.Sprintf("        ExceptionHandlers:\n"))
		for j, eh := range t.ExceptionHandlers {
			buf.WriteString(fmt.Sprintf("            %d: %v\n", j, eh))
		}
	}

	return buf.String()
}
