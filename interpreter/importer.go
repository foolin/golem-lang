// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package interpreter

import (
	g "github.com/mjarmy/golem-lang/core"
)

// An Importer is responsible for managing imported Modules.
type Importer interface {
	GetModule(itp *Interpreter, name string) (g.Module, error)
}

type importer struct {
	moduleMap map[string]g.Module
}

// NewImporter creates an Importer for a pre-defined collection
// of modules.  No new modules will be discovered or added during run-time.
func NewImporter(modules []g.Module) Importer {

	var moduleMap = map[string]g.Module{}
	for _, m := range modules {
		moduleMap[m.Name()] = m
	}

	return &importer{moduleMap}
}

func (imp *importer) GetModule(itp *Interpreter, name string) (g.Module, error) {

	if m, ok := imp.moduleMap[name]; ok {
		return m, nil
	}
	return nil, g.UndefinedModule(name)
}
