// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package bytecode

import (
	"bytes"
	"fmt"

	g "github.com/mjarmy/golem-lang/ncore"
)

// Module is a namespace containing compiled Golem code.  Modules are the fundamental
// unit of compilation in Golem.
type Module struct {
	// Name is the name of the module
	Name string
	// Path is the path to the source code of the module
	Path string
	// Pool is a pool of constants, function templates, and struct definitions
	Pool *Pool
	// Refs is a list of Refs for the values contained in the Contents
	Refs []*Ref

	// Contents is the top level exported values of the module
	Contents g.Struct
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

// Pool is a pool of the constants, function templates, and struct definitions
// used by a given Module.  Pools are created at compile time, and
// are immutable at run time.
type Pool struct {
	Constants  []g.Basic
	Templates  []*FuncTemplate
	StructDefs [][]string
}
