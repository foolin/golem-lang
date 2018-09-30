// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package bytecode

import (
	"bytes"
	"fmt"

	g "github.com/mjarmy/golem-lang/core"
)

// Module is a namespace containing compiled Golem code.  Modules are the fundamental
// unit of compilation in Golem.
type Module struct {

	// Name is the name of the module
	name string

	// Contents is the top level exported values of the module.
	// The Contents are populated when the Interpreter instantiates the module.
	contents g.Struct

	// Path is the path to the source code of the module
	Path string

	// Pool is a pool of constants, function templates, and struct definitions
	Pool *Pool

	// Refs is a list of Refs for the values contained in the Contents.
	// The Refs are populated when the Interpreter instantiates the module.
	Refs []*Ref
}

func NewModule(name, path string) *Module {
	return &Module{
		name: name,
		Path: path,
	}
}

func (m *Module) String() string {
	var buf bytes.Buffer
	buf.WriteString("----------------------------\n")
	buf.WriteString("Module:\n")
	buf.WriteString(fmt.Sprintf("    Name: %s\n", m.name))
	buf.WriteString(fmt.Sprintf("    Path: %s\n", m.Path))

	buf.WriteString("    Refs:\n")
	for i, ref := range m.Refs {
		buf.WriteString("        ")
		buf.WriteString(fmt.Sprintf("%d: %v\n", i, ref))
	}

	return buf.String()
}

func (m *Module) Name() string {
	return m.name
}

func (m *Module) Contents() g.Struct {
	return m.contents
}

func (m *Module) SetContents(contents g.Struct) {
	m.contents = contents
}
