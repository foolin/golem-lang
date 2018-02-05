// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package lib

import (
	"sync"

	g "github.com/mjarmy/golem-lang/core"
)

var libModules = make(map[string]g.Module)
var mutex = &sync.Mutex{}

// LookupModule resolves up the module with the given name.
func LookupModule(name string) (g.Module, g.Error) {

	mutex.Lock()
	defer mutex.Unlock()

	mod, ok := libModules[name]
	if ok {
		return mod, nil
	}

	switch name {
	case "regexp":
		m := NewRegexpModule()
		libModules[name] = m
		return m, nil
	case "os":
		m := NewOsModule()
		libModules[name] = m
		return m, nil
	case "path":
		m := NewPathModule()
		libModules[name] = m
		return m, nil
	default:
		return nil, g.UndefinedModuleError(name)
	}
}
