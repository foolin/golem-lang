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
		rgx := NewRegexpModule()
		libModules[name] = rgx
		return rgx, nil
	case "sys":
		sys := NewSysModule()
		libModules[name] = sys
		return sys, nil
	default:
		return nil, g.UndefinedModuleError(name)
	}
}
