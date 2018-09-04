// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package compiler

import (
	"sort"

	g "github.com/mjarmy/golem-lang/core"
	bc "github.com/mjarmy/golem-lang/core/bytecode"
)

// poolBuilder builds a Pool
type poolBuilder struct {
	constants  *g.HashMap
	templates  []*bc.FuncTemplate
	structDefs [][]string
}

func newPoolBuilder() *poolBuilder {
	return &poolBuilder{
		constants:  g.EmptyHashMap(),
		templates:  []*bc.FuncTemplate{},
		structDefs: [][]string{},
	}
}

func (p *poolBuilder) constIndex(key g.Basic) int {

	// Its OK for the Eval to be nil here.
	// The key is always g.Basic, so the Eval will never be used.
	var ev g.Eval

	b, err := p.constants.Contains(ev, key)
	g.Assert(err == nil)

	if b.BoolVal() {
		var v g.Value
		v, err = p.constants.Get(ev, key)
		g.Assert(err == nil)

		i, ok := v.(g.Int)
		g.Assert(ok)
		return int(i.IntVal())
	}
	i := p.constants.Len()
	err = p.constants.Put(ev, key, i)
	g.Assert(err == nil)
	return int(i.IntVal())
}

func (p *poolBuilder) addTemplate(tpl *bc.FuncTemplate) {
	p.templates = append(p.templates, tpl)
}

func (p *poolBuilder) structDefIndex(def []string) int {

	idx := len(p.structDefs)
	p.structDefs = append(p.structDefs, def)
	return idx
}

func (p *poolBuilder) build() *bc.Pool {
	return &bc.Pool{
		Constants:  p.makeConstants(),
		Templates:  p.templates,
		StructDefs: p.structDefs,
	}
}

//--------------------------------------------------------------

type constEntries []*g.HEntry

func (items constEntries) Len() int {
	return len(items)
}

func (items constEntries) Less(i, j int) bool {

	x, ok := items[i].Value.(g.Int)
	g.Assert(ok)

	y, ok := items[j].Value.(g.Int)
	g.Assert(ok)

	return x.IntVal() < y.IntVal()
}

func (items constEntries) Swap(i, j int) {
	items[i], items[j] = items[j], items[i]
}

func (p *poolBuilder) makeConstants() []g.Basic {

	n := int(p.constants.Len().IntVal())

	entries := make([]*g.HEntry, 0, n)
	itr := p.constants.Iterator()
	for itr.Next() {
		entries = append(entries, itr.Get())
	}

	sort.Sort(constEntries(entries))

	constants := make([]g.Basic, n)
	for i, e := range entries {
		b, ok := e.Key.(g.Basic)
		g.Assert(ok)
		constants[i] = b
	}

	return constants
}
