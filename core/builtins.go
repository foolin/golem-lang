// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package core

import (
//	"fmt"
)

type (
	// BuiltinManager manages the built-in functions (and other built-in values)
	// for a given instance of the Interpreter
	BuiltinManager interface {
		Builtins() []Value
		Contains(s string) bool
		IndexOf(s string) int
	}

	// BuiltinEntry is an entry in a BuiltinManager
	BuiltinEntry struct {
		Name  string
		Value Value
	}

	builtinManager struct {
		values []Value
		lookup map[string]int
	}
)

// NewBuiltinManager creates a new BuiltinManager
func NewBuiltinManager(entries []*BuiltinEntry) BuiltinManager {
	values := make([]Value, len(entries))
	lookup := make(map[string]int)
	for i, e := range entries {
		values[i] = e.Value
		lookup[e.Name] = i
	}
	return &builtinManager{values, lookup}
}

func (b *builtinManager) Builtins() []Value {
	return b.values
}

func (b *builtinManager) Contains(s string) bool {
	_, ok := b.lookup[s]
	return ok

}

func (b *builtinManager) IndexOf(s string) int {
	index, ok := b.lookup[s]
	if !ok {
		panic("unknown builtin")
	}
	return index
}
