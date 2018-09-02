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

		buf.WriteString(fmt.Sprintf("        Bytecodes:\n"))
		btc := t.Bytecodes
		ip := 0
		curLine := -1
		for ip < len(btc) {
			buf.WriteString("            ")
			buf.WriteString(FmtBytecode(btc, ip))

			bc := btc[ip]
			switch bc {
			case GetField, ImportModule, LoadConst: //  InitField, SetField, IncField
				prm := DecodeParam(btc, ip)
				buf.WriteString(fmt.Sprintf(" # '%v'", p.Constants[prm]))
			case InvokeField:
				prm, _ := DecodeWideParams(btc, ip)
				buf.WriteString(fmt.Sprintf(" # '%v'", p.Constants[prm]))
			}

			if curLine != t.LineNumber(ip) {
				curLine = t.LineNumber(ip)
				buf.WriteString(fmt.Sprintf("  // line %d", curLine))
			}

			buf.WriteString("\n")
			ip += BytecodeSize(btc[ip])
		}

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
