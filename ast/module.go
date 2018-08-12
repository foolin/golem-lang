// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package ast

// Module is the basic unit of compilation in Golem
type Module struct {
	Name     string
	Path     string
	Imports  []*ImportStmt
	InitFunc *FnExpr
}
