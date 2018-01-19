// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package lib

import (
	"errors"
	"fmt"
	g "github.com/mjarmy/golem-lang/core"
)

var libModules = make(map[string]g.Module)

func LookupModule(name string) (g.Module, error) {

	mod, ok := libModules[name]
	if ok {
		return mod, nil
	} else {
		switch name {
		case "io":
			io := InitIoModule()
			libModules[name] = io
			return io, nil
		case "regex":
			rgx := InitRegexModule()
			libModules[name] = rgx
			return rgx, nil
		case "sys":
			sys := InitSysModule()
			libModules[name] = sys
			return sys, nil
		default:
			return nil, errors.New(fmt.Sprintf("Module '%s' is not defined", name))
		}
	}
}
