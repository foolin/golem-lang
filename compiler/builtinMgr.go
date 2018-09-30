// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package compiler

import (
	g "github.com/mjarmy/golem-lang/core"
	"sort"
)

type (
	// BuiltinManager manages the built-in functions (and other built-in values)
	// for a given instance of the Interpreter
	BuiltinManager interface {
		Builtins() []g.Value
		Contains(s string) bool
		IndexOf(s string) int
	}

	builtinManager struct {
		values []g.Value
		lookup map[string]int
	}
)

// NewBuiltinManager creates a new BuiltinManager
func NewBuiltinManager(builtins map[string]g.Value) BuiltinManager {

	// sort builtins by name
	type entry struct {
		name string
		val  g.Value
	}
	entries := make([]entry, len(builtins))
	i := 0
	for k, v := range builtins {
		entries[i] = entry{k, v}
		i++
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].name < entries[j].name
	})

	// create manager
	values := make([]g.Value, len(entries))
	lookup := make(map[string]int)
	for i, e := range entries {
		values[i] = e.val
		lookup[e.name] = i
	}
	return &builtinManager{values, lookup}
}

func (b *builtinManager) Builtins() []g.Value {
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
