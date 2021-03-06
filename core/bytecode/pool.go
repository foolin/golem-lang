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
	return p.DebugString(func(curLine int) string {
		return fmt.Sprintf("// line %d", curLine)
	})
}

// DebugString prints diagnositc information
func (p *Pool) DebugString(printLine func(int) string) string {

	var buf bytes.Buffer

	buf.WriteString("Constants:\n")
	for i, b := range p.Constants {
		s, err := b.ToStr(nil)
		g.Assert(err == nil)
		buf.WriteString(fmt.Sprintf("    %d: %s '%v'\n", i, b.Type(), s))
	}

	buf.WriteString("StructDefs:\n")
	for i, d := range p.StructDefs {
		buf.WriteString(fmt.Sprintf("    %d: %v\n", i, d))
	}

	buf.WriteString("Templates:\n")
	for i, t := range p.Templates {
		buf.WriteString(fmt.Sprintf("    %d: Template\n", i))
		buf.WriteString(fmt.Sprintf("        Arity: %s\n", t.Arity))
		buf.WriteString(fmt.Sprintf("        OptionalParams: %v\n", t.OptionalParams))
		buf.WriteString(fmt.Sprintf("        NumCaptures: %d\n", t.NumCaptures))
		buf.WriteString(fmt.Sprintf("        NumLocals: %d\n", t.NumLocals))

		buf.WriteString(fmt.Sprintf("        Bytecodes:\n"))
		btc := t.Bytecodes
		ip := 0
		curLine := -1
		for ip < len(btc) {
			buf.WriteString("            ")
			buf.WriteString(padRight(FmtBytecode(btc, ip), 28))

			//bc := btc[ip]
			//switch bc {
			//case GetField, ImportModule, LoadConst: //  InitField, SetField, IncField
			//	q := DecodeParam(btc, ip)
			//	c := p.Constants[q]
			//	buf.WriteString(fmt.Sprintf(" # %s '%v'", c.Type(), c))
			//case InvokeField:
			//	q, _ := DecodeWideParams(btc, ip)
			//	c := p.Constants[q]
			//	buf.WriteString(fmt.Sprintf(" # %s '%v'", c.Type(), c))
			//}

			if curLine != t.LineNumber(ip) {
				curLine = t.LineNumber(ip)
				buf.WriteString(printLine(curLine))
			}

			buf.WriteString("\n")
			ip += Size(btc[ip])
		}

		buf.WriteString(fmt.Sprintf("        LineNumberTable:\n"))
		for j, ln := range t.LineNumberTable {
			buf.WriteString(fmt.Sprintf("            %d: %v\n", j, ln))
		}

		buf.WriteString(fmt.Sprintf("        ErrorHandlers:\n"))
		for j, eh := range t.ErrorHandlers {
			buf.WriteString(fmt.Sprintf("            %d: %v\n", j, eh))
		}
	}

	return buf.String()
}
