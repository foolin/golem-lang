// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package ast

// Module is the basic unit of compilation in Golem
type Module struct {
	Name     string
	Path     string
	InitFunc *FnExpr
}

// Imports returns the names of the Modules that are imported by this Module.
func (m *Module) Imports() []string {

	imports := []string{}

	for _, stmt := range m.InitFunc.Body.Statements {
		imp, ok := stmt.(*ImportStmt)
		if !ok {
			break
		}
		for _, ident := range imp.Idents {
			imports = append(imports, ident.Symbol.Text)
		}
	}

	return imports
}
