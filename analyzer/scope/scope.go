// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package scope

import (
	"bytes"
	"fmt"
	"github.com/mjarmy/golem-lang/ast"
	"sort"
)

type (
	// Scope is scoping information associated with an AST Node
	Scope struct {
		Parent *Scope
		defs   map[string]*ast.Variable

		scopeType   scopeType
		FuncScope   *FuncScope
		structScope *structScope
	}

	// FuncScope is scoping information associated with an AST Function Node
	FuncScope struct {
		NumLocals      int
		Captures       map[string]*ast.Variable
		ParentCaptures map[string]*ast.Variable
	}

	structScope struct {
		stc *ast.StructExpr
	}
)

// NewBlockScope creates scope for a block Node
func NewBlockScope(parent *Scope) *Scope {
	return newScope(parent, blockType, nil, nil)
}

// NewFuncScope creates scope for a function Node
func NewFuncScope(parent *Scope) *Scope {
	return newScope(
		parent, funcType,
		&FuncScope{
			0,
			make(map[string]*ast.Variable),
			make(map[string]*ast.Variable)},
		nil)
}

// NewStructScope creates scope for a struct Node
func NewStructScope(parent *Scope, stc *ast.StructExpr) *Scope {
	return newScope(parent, structType, nil, &structScope{stc})
}

func newScope(
	parent *Scope,
	scopeType scopeType,
	FuncScope *FuncScope,
	structScope *structScope) *Scope {

	s := &Scope{
		parent,
		make(map[string]*ast.Variable),
		scopeType,
		FuncScope,
		structScope}

	return s
}

func (s *Scope) String() string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("%v", s.scopeType))

	buf.WriteString(" defs:")
	buf.WriteString(mapString(s.defs))

	if s.scopeType == funcType {
		buf.WriteString(" captures:")
		buf.WriteString(mapString(s.FuncScope.Captures))
		buf.WriteString(" parentCaptures:")
		buf.WriteString(mapString(s.FuncScope.ParentCaptures))
		buf.WriteString(fmt.Sprintf(" numLocals:%d", s.FuncScope.NumLocals))
	}

	return buf.String()
}

// Put defines a Variable, either as a formal param for a Function,
// or via Let or Const.
func (s *Scope) Put(sym string, isConst bool) *ast.Variable {

	_, ok := s.defs[sym]
	if ok {
		panic("symbol is already defined")
	}
	v := &ast.Variable{sym, incrementNumLocals(s), isConst, false}
	s.defs[sym] = v
	return v
}

// Get a Variable, by traversing up the Scope stack.
func (s *Scope) Get(sym string) (*ast.Variable, bool) {

	// make a stack of every Scope we encounter.
	stack := make([]*Scope, 0, 4)
	z := s
	for {
		stack = append(stack, z)

		// Check to see if the symbol is already captured
		if z.scopeType == funcType {
			if v, ok := z.FuncScope.Captures[sym]; ok {
				return v, true
			}
		}

		// If the varable is defined in the current Scope, then
		// capture anything that needs to be captured in the stack we have
		// created, and then return.
		if _, ok := z.defs[sym]; ok {
			return capture(sym, stack), true
		}

		// Still can't find it.  Go to the parent, or return false.
		if z.Parent == nil {
			return nil, false
		}
		z = z.Parent
	}
}

// This creates a variable for 'this', or return an existing 'this' variable.
func (s *Scope) This() *ast.Variable {

	// find the nearest parent structScope
	os := s
	for os.scopeType != structType {
		os = os.Parent
	}

	// define a 'this' variable on the structScope, if its not already defined
	_, ok := os.defs["this"]
	if !ok {
		idx := incrementNumLocals(os)
		os.defs["this"] = &ast.Variable{"this", idx, true, false}
		os.structScope.stc.LocalThisIndex = idx
	}

	// now call get(), from the original Scope, to trigger Captures in
	// any intervening functions.
	v, ok := s.Get("this")
	if !ok {
		panic("call to 'this' failed")
	}

	// done
	return v
}

// Increment the number of local variables in the nearest parent FuncScope.
func incrementNumLocals(s *Scope) int {

	for {
		if s.scopeType == funcType {
			idx := s.FuncScope.NumLocals
			s.FuncScope.NumLocals++
			if s.FuncScope.NumLocals >= (2 << 16) {
				panic("TODO wide index")
			}
			return idx
		}
		s = s.Parent
	}
}

// Create a succession of Captures in the stack of scopes.
func capture(sym string, stack []*Scope) *ast.Variable {

	n := len(stack)
	v := stack[n-1].defs[sym]

	if n == 1 {
		return v
	}

	// look for functions between the beginning and end, exclusive
	for a := n - 2; a >= 1; a-- {
		if stack[a].scopeType == funcType {
			s := stack[a]

			// capture a variable from the parent Scope down into this parent Scope
			idx := len(s.FuncScope.ParentCaptures)
			s.FuncScope.ParentCaptures[sym] = v

			// capture the variable in this Scope
			v = &ast.Variable{sym, idx, v.IsConst, true}
			s.FuncScope.Captures[sym] = v
		}
	}

	return v
}

//-----------------------------------------------------------------------------

type scopeType int

const (
	blockType scopeType = iota
	funcType
	structType
)

func (s scopeType) String() string {
	if s == funcType {
		return "Func"
	} else if s == structType {
		return "Struct"
	} else {
		return "Block"
	}
}

//-----------------------------------------------------------------------------

func mapString(m map[string]*ast.Variable) string {

	keys := make([]string, len(m))
	i := 0
	for k := range m {
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
		buf.WriteString(fmt.Sprintf("%v: %v", k, m[k]))
	}
	buf.WriteString("}")
	return buf.String()
}
