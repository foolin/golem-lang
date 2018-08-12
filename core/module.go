// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package core

import (
	"bytes"
	"fmt"
	//"github.com/mjarmy/golem-lang/core/opcodes"
)

//--------------------------------------------------------------
// Module
//--------------------------------------------------------------

// Module is a namespace containing compiled Golem code.  Modules are the fundamental
// unit of compilation in Golem.
type Module struct {
	// Name is the name of the module
	Name string
	// Path is the path to the source code of the module
	Path string
	// Pool is a pool of constants, function templates, and struct definitions
	Pool *Pool
	// Contents is the top level exported values of the module
	Contents Struct
	// Refs is a list of Refs for the values contained in the Contents
	Refs []*Ref
}

func (m *Module) String() string {
	var buf bytes.Buffer
	buf.WriteString("----------------------------\n")
	buf.WriteString("Module:\n")
	buf.WriteString(fmt.Sprintf("    Name: %s\n", m.Name))
	buf.WriteString(fmt.Sprintf("    Path: %s\n", m.Path))

	buf.WriteString("    Refs:\n")
	for i, ref := range m.Refs {
		buf.WriteString("        ")
		buf.WriteString(fmt.Sprintf("%d: %v\n", i, ref))
	}

	return buf.String()
}

//--------------------------------------------------------------
// Ref
//--------------------------------------------------------------

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

//--------------------------------------------------------------
// Pool
//--------------------------------------------------------------

// Pool is a pool of constants, function templates, and struct definitions
// used by a given Module.  Pools are created at compile time, and
// are immutable at run time.
type Pool struct {
	Constants  []Basic
	Templates  []*FuncTemplate
	StructDefs [][]*FieldDef
}
