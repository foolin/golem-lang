// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package core

import (
	"bytes"
	"fmt"

	"github.com/mjarmy/golem-lang/core/opcodes"
)

// Module is a namespace containing compiled Golem code
type Module struct {
	Name       string
	Pool       []Basic
	Refs       []*Ref
	StructDefs [][]*FieldDef
	Templates  []*Template
	Contents   Struct
}

// GetModuleName for a Module returns the empty string.
func (m *Module) GetModuleName() string {
	return ""
}

// GetContents returns the contents of a Module.
func (m *Module) GetContents() Struct {
	return m.Contents
}

func (m *Module) String() string {
	var buf bytes.Buffer
	buf.WriteString("----------------------------\n")
	buf.WriteString("Module:\n")

	buf.WriteString("    Pool:\n")
	for i, val := range m.Pool {
		typeOf := val.Type()
		buf.WriteString("        ")
		buf.WriteString(fmt.Sprintf("%d: %v(%v)\n", i, typeOf, val))
	}

	buf.WriteString("    Refs:\n")
	for i, ref := range m.Refs {
		buf.WriteString("        ")
		buf.WriteString(fmt.Sprintf("%d: %v\n", i, ref))
	}

	buf.WriteString("    StructDefs:\n")
	for i, def := range m.StructDefs {
		buf.WriteString("        ")
		buf.WriteString(fmt.Sprintf("%d: %v\n", i, def))
	}

	for i, t := range m.Templates {

		buf.WriteString(fmt.Sprintf(
			"    Template(%d): Arity: %d, NumCaptures: %d, NumLocals: %d\n",
			i, t.Arity, t.NumCaptures, t.NumLocals))

		buf.WriteString("        OpCodes:\n")
		for i := 0; i < len(t.OpCodes); {
			text := opcodes.FmtOpcode(t.OpCodes, i)
			buf.WriteString("            ")
			buf.WriteString(text)
			i += opcodes.OpCodeSize(t.OpCodes[i])
		}

		buf.WriteString("        LineNumberTable:\n")
		for _, ln := range t.LineNumberTable {
			buf.WriteString("            ")
			buf.WriteString(fmt.Sprintf("%v\n", ln))
		}

		buf.WriteString("        ExceptionHandlers:\n")
		for _, eh := range t.ExceptionHandlers {
			buf.WriteString("            ")
			buf.WriteString(fmt.Sprintf("%v\n", eh))
		}
	}

	return buf.String()
}

//---------------------------------------------------------------

// Ref is a container for a Value.  Refs are used by the interpreter
// as a place to store the value of a variable.
type Ref struct {
	Val Value
}

// NewRef creates a new Ref
func NewRef(val Value) *Ref {
	return &Ref{val}
}

func (r *Ref) String() string {
	return fmt.Sprintf("Ref(%v)", r.Val)
}
