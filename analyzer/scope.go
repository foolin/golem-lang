// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package analyzer

import (
	"bytes"
	"fmt"
	"github.com/mjarmy/golem-lang/ast"
	"sort"
)

type scopeType int

const (
	funcType scopeType = iota
	blockType
	structType
)

func (s scopeType) String() string {
	if s == funcType {
		return "Func "
	} else if s == blockType {
		return "BlockNode"
	} else {
		return "Struct"
	}
}

type scope struct {
	parent *scope
	defs   map[string]*ast.Variable

	scopeType   scopeType
	funcScope   *funcScope
	structScope *structScope
}

type funcScope struct {
	numLocals      int
	captures       map[string]*ast.Variable
	parentCaptures map[string]*ast.Variable
}

type structScope struct {
	stc *ast.StructExpr
}

func newFuncScope(parent *scope) *scope {
	return newScope(
		parent, funcType,
		&funcScope{
			0,
			make(map[string]*ast.Variable),
			make(map[string]*ast.Variable)},
		nil)
}

func newBlockScope(parent *scope) *scope {
	return newScope(parent, blockType, nil, nil)
}

func newStructScope(parent *scope, stc *ast.StructExpr) *scope {
	return newScope(parent, structType, nil, &structScope{stc})
}

func newScope(
	parent *scope,
	scopeType scopeType,
	funcScope *funcScope,
	structScope *structScope) *scope {

	s := &scope{
		parent,
		make(map[string]*ast.Variable),
		scopeType,
		funcScope,
		structScope}

	return s
}

func (s *scope) String() string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("%v", s.scopeType))

	buf.WriteString(" defs:")
	buf.WriteString(mapString(s.defs))

	if s.scopeType == funcType {
		buf.WriteString(" captures:")
		buf.WriteString(mapString(s.funcScope.captures))
		buf.WriteString(" parentCaptures:")
		buf.WriteString(mapString(s.funcScope.parentCaptures))
		buf.WriteString(fmt.Sprintf(" numLocals:%d", s.funcScope.numLocals))
	}

	return buf.String()
}

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

// Define a Variable, either as a formal param for a Function,
// or via LetStmt or ConstStmt.
func (s *scope) put(sym string, isConst bool) *ast.Variable {

	_, ok := s.defs[sym]
	if ok {
		panic("symbol is already defined")
	}
	v := &ast.Variable{sym, incrementNumLocals(s), isConst, false}
	s.defs[sym] = v
	return v
}

// Create a variable for 'this', or return an existing 'this' variable.
func (s *scope) this() *ast.Variable {

	// find the nearest parent structScope
	os := s
	for os.scopeType != structType {
		os = os.parent
	}

	// define a 'this' variable on the structScope, if its not already defined
	_, ok := os.defs["this"]
	if !ok {
		idx := incrementNumLocals(os)
		os.defs["this"] = &ast.Variable{"this", idx, true, false}
		os.structScope.stc.LocalThisIndex = idx
	}

	// now call get(), from the original scope, to trigger captures in
	// any intervening functions.
	v, ok := s.get("this")
	if !ok {
		panic("call to 'this' failed")
	}

	// done
	return v
}

// Increment the number of local variables in the nearest parent funcScope.
func incrementNumLocals(s *scope) int {

	for {
		if s.scopeType == funcType {
			idx := s.funcScope.numLocals
			s.funcScope.numLocals++
			if s.funcScope.numLocals >= (2 << 16) {
				panic("TODO wide index")
			}
			return idx
		}
		s = s.parent
	}
}

// Get a Variable, by traversing up the scope stack.
func (s *scope) get(sym string) (*ast.Variable, bool) {

	// make a stack of every scope we encounter.
	stack := make([]*scope, 0, 4)
	z := s
	for {
		stack = append(stack, z)

		// Check to see if the symbol is already captured
		if z.scopeType == funcType {
			if v, ok := z.funcScope.captures[sym]; ok {
				return v, true
			}
		}

		// IfStmt the varable is defined in the current scope, then
		// capture anything that needs to be captured in the stack we have
		// created, and then return.
		if _, ok := z.defs[sym]; ok {
			return capture(sym, stack), true
		}

		// Still can't find it.  GoStmt to the parent, or return false.
		if z.parent == nil {
			return nil, false
		}
		z = z.parent
	}
}

// Create a succession of captures in the stack of scopes.
func capture(sym string, stack []*scope) *ast.Variable {

	n := len(stack)
	v := stack[n-1].defs[sym]

	if n == 1 {
		return v
	}

	// look for functions between the beginning and end, exclusive
	for a := n - 2; a >= 1; a-- {
		if stack[a].scopeType == funcType {
			s := stack[a]

			// capture a variable from the parent scope down into this parent scope
			idx := len(s.funcScope.parentCaptures)
			s.funcScope.parentCaptures[sym] = v

			// capture the variable in this scope
			v = &ast.Variable{sym, idx, v.IsConst, true}
			s.funcScope.captures[sym] = v
		}
	}

	return v
}
