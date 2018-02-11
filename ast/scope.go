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

		GetCapture(string) (Capture, bool)
		PutCapture(Variable) Capture

		NumCaptures() int
		GetCaptures() []Capture
	}

	// Capture defines a parent/child relationship between a 'child' captured Variable,
	// and a 'parent' Variable that is found higher up in the AST.
	Capture interface {
		Symbol() string
		Parent() Variable
		Child() Variable
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
	captures  map[string]Capture
}

// NewFuncScope creates a new FuncScope
func NewFuncScope() FuncScope {
	return &funcScope{
		scope{make(map[string]Variable)},
		0,
		make(map[string]Capture),
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

func (s *funcScope) GetCapture(sym string) (Capture, bool) {
	cp, ok := s.captures[sym]
	return cp, ok
}

func (s *funcScope) PutCapture(parent Variable) Capture {
	sym := parent.Symbol()
	child := NewVariable(sym, len(s.captures), parent.IsConst(), true)
	cp := &capture{sym, parent, child}
	s.captures[sym] = cp
	return cp
}

func (s *funcScope) NumCaptures() int {
	return len(s.captures)
}

func (s *funcScope) GetCaptures() []Capture {
	cp := make([]Capture, len(s.captures))
	n := 0
	for _, v := range s.captures {
		cp[n] = v
		n++
	}
	return cp
}

//-----------------------------------------------------------------------------
// Capture
//-----------------------------------------------------------------------------

type capture struct {
	symbol string
	parent Variable
	child  Variable
}

func (c *capture) Symbol() string   { return c.symbol }
func (c *capture) Parent() Variable { return c.parent }
func (c *capture) Child() Variable  { return c.child }

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

func capString(caps map[string]Capture) string {

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
		buf.WriteString(fmt.Sprintf("%v: (parent: %v, child %v)", k, c.Parent(), c.Child()))
	}
	buf.WriteString("}")
	return buf.String()
}
