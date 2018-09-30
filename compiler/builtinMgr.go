// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package compiler

import (
	g "github.com/mjarmy/golem-lang/core"
)

type builtinManager struct {
	builtins []*g.Builtin
	lookup   map[string]int
}

// newBuiltinManager creates a new BuiltinManager
func newBuiltinManager(builtins []*g.Builtin) *builtinManager {
	lookup := make(map[string]int)
	for i, e := range builtins {
		lookup[e.Name] = i
	}
	return &builtinManager{builtins, lookup}
}

func (b *builtinManager) contains(s string) bool {
	_, ok := b.lookup[s]
	return ok
}

func (b *builtinManager) indexOf(s string) int {
	index, ok := b.lookup[s]
	if !ok {
		panic("unknown builtin")
	}
	return index
}
