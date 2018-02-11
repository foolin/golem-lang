// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package ast

import (
	"bytes"
	"fmt"
	"sort"
)

type (

	// Scope contains scoping information associated with a BlockNode, ForStmt, or TryStmt.
	Scope interface {
		GetVariable(string) (Variable, bool)
		PutVariable(string, Variable)
	}

	// StructScope contains scoping information associated with a StructExpr.  The
	// only symbol that will ever be defined in a StructScope is 'this'.
	StructScope interface {
		Scope
		structMarker()
	}

	// FuncScope contains scoping information associated with a FnExpr.
	FuncScope interface {
		Scope

		NumLocals() int
		IncrementNumLocals()

		NumCaptures() int
		GetCapture(string) (Variable, bool)
		PutCapture(Variable) Variable

		GetParentCaptures() []Variable
	}

	// A Variable contains information about how an Identifer is defined in a Node.
	// Variables are created directly via:
	//
	//      (1) A 'let' or 'const' stmt.
	//		(2) The 'for' clause of a for-loop.
	//		(3) The 'catch' clause of a try-catch block.
	//		(4) A 'this' expression
	// 		(5) The formal parameters of a function.
	//
	// Variables are created indirectly via the capture mechanism, in which
	// a Scope references a Variable that is defined in one of its ancestor Scopes
	// in the AST, and there are intervening FuncScopes that must capture the Variable.
	//
	Variable interface {
		Symbol() string
		Index() int
		IsConst() bool
		IsCapture() bool
	}
)

//-----------------------------------------------------------------------------
// Scope
//-----------------------------------------------------------------------------

type scope struct {
	defs map[string]Variable
}

// NewScope creates a new Scope
func NewScope() Scope {
	return &scope{make(map[string]Variable)}
}

func (s *scope) String() string {
	var buf bytes.Buffer

	buf.WriteString("Scope")
	buf.WriteString(" defs:")
	buf.WriteString(defString(s.defs))

	return buf.String()
}

func (s *scope) GetVariable(sym string) (Variable, bool) {
	v, ok := s.defs[sym]
	return v, ok
}

func (s *scope) PutVariable(sym string, v Variable) {
	s.defs[sym] = v
}

//-----------------------------------------------------------------------------
// StructScope
//-----------------------------------------------------------------------------

type structScope struct {
	scope
}

// NewStructScope creates a new StructScope
func NewStructScope() StructScope {
	return &structScope{
		scope{make(map[string]Variable)},
	}
}

func (s *structScope) String() string {
	var buf bytes.Buffer

	buf.WriteString("StructScope")
	buf.WriteString(" defs:")
	buf.WriteString(defString(s.defs))

	return buf.String()
}

func (s *structScope) structMarker() {}

//-----------------------------------------------------------------------------
// FuncScope
//-----------------------------------------------------------------------------

type funcScope struct {
	scope
	numLocals int
	captures  map[string]capture
}

type capture struct {
	parent Variable
	child  Variable
}

// NewFuncScope creates a new FuncScope
func NewFuncScope() FuncScope {
	return &funcScope{
		scope{make(map[string]Variable)},
		0,
		make(map[string]capture),
	}
}

func (s *funcScope) String() string {
	var buf bytes.Buffer

	buf.WriteString("FuncScope")
	buf.WriteString(" defs:")
	buf.WriteString(defString(s.defs))

	buf.WriteString(" captures:")
	buf.WriteString(capString(s.captures))
	buf.WriteString(fmt.Sprintf(" numLocals:%d", s.numLocals))

	return buf.String()
}

func (s *funcScope) NumLocals() int {
	return s.numLocals
}

func (s *funcScope) IncrementNumLocals() {
	s.numLocals++
}

func (s *funcScope) NumCaptures() int {
	return len(s.captures)
}

func (s *funcScope) GetCapture(sym string) (Variable, bool) {
	cap, ok := s.captures[sym]
	if ok {
		return cap.child, true
	}
	return nil, false
}

func (s *funcScope) PutCapture(parent Variable) Variable {

	sym := parent.Symbol()
	child := NewVariable(sym, len(s.captures), parent.IsConst(), true)
	s.captures[sym] = capture{parent, child}
	return child
}

func (s *funcScope) GetParentCaptures() []Variable {

	// First sort the captures by child index
	sorted := byChildIndex{}
	for _, v := range s.captures {
		sorted = append(sorted, v)
	}
	sort.Sort(sorted)

	// Then use the sorted list to create the proper ordering of parents
	parents := []Variable{}
	for _, c := range sorted {
		parents = append(parents, c.parent)
	}
	return parents
}

type byChildIndex []capture

// Variables are sorted by Index
func (c byChildIndex) Len() int {
	return len(c)
}
func (c byChildIndex) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}
func (c byChildIndex) Less(i, j int) bool {
	return c[i].child.Index() < c[j].child.Index()
}

//-----------------------------------------------------------------------------
// Variable
//-----------------------------------------------------------------------------

type variable struct {
	debugID int

	symbol    string
	index     int
	isConst   bool
	isCapture bool
}

var debugIDCounter int

// InternalResetDebugging is a 'secret' internal function.  Please pretend its not here.
func InternalResetDebugging() {
	debugIDCounter = 0
}

// NewVariable creates a new Variable
func NewVariable(symbol string, index int, isConst bool, isCapture bool) Variable {
	debugID := debugIDCounter
	debugIDCounter++
	return &variable{debugID, symbol, index, isConst, isCapture}
}

func (v *variable) String() string {
	return fmt.Sprintf("v(%d: %s,%d,%v,%v)", v.debugID, v.symbol, v.index, v.isConst, v.isCapture)
}

func (v *variable) Symbol() string {
	return v.symbol
}

func (v *variable) Index() int {
	return v.index
}

func (v *variable) IsConst() bool {
	return v.isConst
}

func (v *variable) IsCapture() bool {
	return v.isCapture
}

//-----------------------------------------------------------------------------
// util
//-----------------------------------------------------------------------------

func defString(defs map[string]Variable) string {

	// sort the keys alphabetically
	keys := make([]string, len(defs))
	i := 0
	for k := range defs {
		keys[i] = k
		i++
	}
	sort.Strings(keys)

	var buf bytes.Buffer
	buf.WriteString("{")
	n := 0
	for _, k := range keys {
		if n > 0 {
			buf.WriteString(", ")
		}
		n++
		buf.WriteString(fmt.Sprintf("%v: %v", k, defs[k]))
	}
	buf.WriteString("}")
	return buf.String()
}

func capString(caps map[string]capture) string {

	// sort the keys alphabetically
	keys := make([]string, len(caps))
	i := 0
	for k := range caps {
		keys[i] = k
		i++
	}
	sort.Strings(keys)

	var buf bytes.Buffer
	buf.WriteString("{")
	n := 0
	for _, k := range keys {
		if n > 0 {
			buf.WriteString(", ")
		}
		n++
		c := caps[k]
		buf.WriteString(fmt.Sprintf("%v: (parent: %v, child %v)", k, c.parent, c.child))
	}
	buf.WriteString("}")
	return buf.String()
}
