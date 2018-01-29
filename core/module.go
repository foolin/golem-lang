// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package core

import (
	"bytes"
	"fmt"

	"github.com/mjarmy/golem-lang/core/opcodes"
)

//---------------------------------------------------------------
// Module

type Module interface {
	GetModuleName() string
	GetContents() Struct
}

//---------------------------------------------------------------
// BytecodeModule

type BytecodeModule struct {
	Pool       []Basic
	Refs       []*Ref
	StructDefs [][]Field
	Templates  []*Template
	Contents   Struct
}

func (m *BytecodeModule) GetModuleName() string {
	return ""
}

func (m *BytecodeModule) GetContents() Struct {
	return m.Contents
}

func (m *BytecodeModule) String() string {
	var buf bytes.Buffer
	buf.WriteString("----------------------------\n")
	buf.WriteString("BytecodeModule:\n")

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
// A Ref is a container for a value

type Ref struct {
	Val Value
}

func NewRef(val Value) *Ref {
	return &Ref{val}
}

func (r *Ref) String() string {
	return fmt.Sprintf("Ref(%v)", r.Val)
}
