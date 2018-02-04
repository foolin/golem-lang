// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package analyzer

import (
	"errors"
	"fmt"
	"github.com/mjarmy/golem-lang/ast"
	"sort"
)

// Analyzer analyzes an AST.
type Analyzer interface {
	ast.Visitor
	Module() *ast.FnExpr
	Analyze() []error
	scope() *scope
}

type analyzer struct {
	mod       *ast.FnExpr
	rootScope *scope
	curScope  *scope
	loops     []ast.Loop
	structs   []*ast.StructExpr
	errors    []error
}

// NewAnalyzer creates a new Analyzer
func NewAnalyzer(mod *ast.FnExpr) Analyzer {

	rootScope := newFuncScope(nil)

	return &analyzer{mod, rootScope, rootScope, []ast.Loop{}, []*ast.StructExpr{}, nil}
}

func (a *analyzer) scope() *scope {
	return a.rootScope
}

// Analyze analyzes an AST.
func (a *analyzer) Analyze() []error {

	// visit module block
	a.visitBlock(a.mod.Body)

	// save NumLocals
	fscope := a.curScope.funcScope
	a.mod.NumLocals = fscope.numLocals

	// sanity check for captures
	if len(fscope.captures) > 0 {
		panic("invalid module")
	}
	a.mod.NumCaptures = 0
	a.mod.ParentCaptures = nil

	// done
	return a.errors
}

func (a *analyzer) Module() *ast.FnExpr {
	return a.mod
}

func (a *analyzer) Visit(node ast.Node) {
	switch t := node.(type) {

	case *ast.BlockNode:
		a.visitBlock(t)

	case *ast.FnExpr:
		a.visitFunc(t)

	case *ast.ImportStmt:
		a.visitImport(t)

	case *ast.ConstStmt:
		a.visitDecls(t.Decls, true)

	case *ast.LetStmt:
		a.visitDecls(t.Decls, false)

	case *ast.AssignmentExpr:
		a.visitAssignment(t)

	case *ast.TryStmt:
		a.visitTry(t)

	case *ast.PostfixExpr:
		a.visitPostfixExpr(t)

	case *ast.IdentExpr:
		a.visitIdentExpr(t)

	case *ast.WhileStmt:
		a.loops = append(a.loops, t)
		t.Traverse(a)
		a.loops = a.loops[:len(a.loops)-1]

	case *ast.ForStmt:
		a.loops = append(a.loops, t)
		a.visitFor(t)
		a.loops = a.loops[:len(a.loops)-1]

	case *ast.BreakStmt:
		if len(a.loops) == 0 {
			a.errors = append(a.errors, errors.New("'break' outside of loop"))
		}

	case *ast.ContinueStmt:
		if len(a.loops) == 0 {
			a.errors = append(a.errors, errors.New("'continue' outside of loop"))
		}

	case *ast.StructExpr:
		a.visitStructExpr(t)

	case *ast.ThisExpr:
		a.visitThisExpr(t)

	default:
		t.Traverse(a)

	}
}

func (a *analyzer) visitDecls(decls []*ast.DeclNode, isConst bool) {

	for _, d := range decls {
		if d.Val != nil {
			a.Visit(d.Val)
		}
		a.defineIdent(d.Ident, isConst)
	}
}

func (a *analyzer) visitImport(imp *ast.ImportStmt) {
	a.defineIdent(imp.Ident, true)
}

func (a *analyzer) visitTry(t *ast.TryStmt) {

	a.Visit(t.TryBlock)
	if t.CatchToken != nil {
		a.curScope = newBlockScope(a.curScope)
		a.defineIdent(t.CatchIdent, true)
		a.Visit(t.CatchBlock)
		a.curScope = a.curScope.parent
	}
	if t.FinallyToken != nil {
		a.Visit(t.FinallyBlock)
	}
}

func (a *analyzer) defineIdent(ident *ast.IdentExpr, isConst bool) {
	sym := ident.Symbol.Text
	if _, ok := a.curScope.get(sym); ok {
		a.errors = append(a.errors,
			fmt.Errorf("Symbol '%s' is already defined", sym))
	} else {
		ident.Variable = a.curScope.put(sym, isConst)
	}
}

func (a *analyzer) visitBlock(blk *ast.BlockNode) {

	a.curScope = newBlockScope(a.curScope)

	// visit named funcs identifiers
	for _, n := range blk.Statements {
		if nf, ok := n.(*ast.NamedFnStmt); ok {
			a.defineIdent(nf.Ident, true)
		}
	}

	// visit everything, skipping named func identifiers
	for _, n := range blk.Statements {
		if nf, ok := n.(*ast.NamedFnStmt); ok {
			a.Visit(nf.Func)
		} else {
			a.Visit(n)
		}
	}

	a.curScope = a.curScope.parent
}

func (a *analyzer) visitFor(fr *ast.ForStmt) {

	// push block scope
	a.curScope = newBlockScope(a.curScope)

	// define identifiers
	for _, ident := range fr.Idents {
		a.defineIdent(ident, false)
	}

	// define the identifier for the iterable
	a.defineIdent(fr.IterableIdent, false)

	// visit the iterable and body
	a.Visit(fr.Iterable)
	a.visitBlock(fr.Body)

	// pop block scope
	a.curScope = a.curScope.parent
}

func (a *analyzer) visitFunc(fn *ast.FnExpr) {

	// push scope
	a.curScope = newFuncScope(a.curScope)

	// visit child nodes
	for _, f := range fn.FormalParams {
		f.Ident.Variable = a.curScope.put(f.Ident.Symbol.Text, f.IsConst)
	}
	a.visitBlock(fn.Body)

	// save function scope info
	fscope := a.curScope.funcScope
	fn.NumLocals = fscope.numLocals
	fn.NumCaptures = len(fscope.captures)
	fn.ParentCaptures = a.makeParentCaptures()

	// pop scope
	a.curScope = a.curScope.parent
}

func (a *analyzer) makeParentCaptures() []*ast.Variable {

	fscope := a.curScope.funcScope
	num := len(fscope.captures)

	if num != len(fscope.parentCaptures) {
		panic("capture length mismatch")
	}
	if num == 0 {
		return nil
	}

	// First, sort the captured Variables by index
	caps := make(byIndex, 0, num)
	for _, v := range fscope.captures {
		caps = append(caps, v)
	}
	sort.Sort(caps)

	// Then use the sorted list to create the proper ordering of parentCaptures
	parentCaps := make([]*ast.Variable, 0, num)
	for _, v := range caps {
		parentCaps = append(parentCaps, fscope.parentCaptures[v.Symbol])
	}
	return parentCaps
}

type byIndex []*ast.Variable

// Variables are sorted by Index
func (v byIndex) Len() int {
	return len(v)
}
func (v byIndex) Swap(i, j int) {
	v[i], v[j] = v[j], v[i]
}
func (v byIndex) Less(i, j int) bool {
	return v[i].Index < v[j].Index
}

func (a *analyzer) visitAssignment(asn *ast.AssignmentExpr) {

	switch t := asn.Assignee.(type) {

	case *ast.IdentExpr:
		a.Visit(asn.Val)
		a.doVisitAssignIdent(t)

	case *ast.FieldExpr:
		a.Visit(t.Operand)
		a.Visit(asn.Val)

	case *ast.IndexExpr:
		a.Visit(t.Operand)
		a.Visit(t.Index)
		a.Visit(asn.Val)

	default:
		panic("invalid assignee type")
	}
}

func (a *analyzer) visitPostfixExpr(ps *ast.PostfixExpr) {

	switch t := ps.Assignee.(type) {

	case *ast.IdentExpr:
		a.doVisitAssignIdent(t)

	case *ast.FieldExpr:
		a.Visit(t.Operand)

	case *ast.IndexExpr:
		a.Visit(t.Operand)
		a.Visit(t.Index)

	default:
		panic("invalid assignee type")
	}
}

// visit an Ident that is part of an assignment
func (a *analyzer) doVisitAssignIdent(ident *ast.IdentExpr) {
	sym := ident.Symbol.Text
	if v, ok := a.curScope.get(sym); ok {
		if v.IsConst {
			a.errors = append(a.errors,
				fmt.Errorf("Symbol '%s' is constant", sym))
		}
		ident.Variable = v
	} else {
		a.errors = append(a.errors,
			fmt.Errorf("Symbol '%s' is not defined", sym))
	}
}

func (a *analyzer) visitIdentExpr(ident *ast.IdentExpr) {

	sym := ident.Symbol.Text

	if v, ok := a.curScope.get(sym); ok {
		ident.Variable = v
	} else {
		a.errors = append(a.errors,
			fmt.Errorf("Symbol '%s' is not defined", sym))
	}
}

func (a *analyzer) visitStructExpr(stc *ast.StructExpr) {
	a.structs = append(a.structs, stc)

	a.curScope = newStructScope(a.curScope, stc)
	stc.Traverse(a)
	a.curScope = a.curScope.parent

	a.structs = a.structs[:len(a.structs)-1]
}

func (a *analyzer) visitThisExpr(this *ast.ThisExpr) {

	n := len(a.structs)
	if n == 0 {
		a.errors = append(a.errors, errors.New("'this' outside of loop"))
	} else {
		this.Variable = a.curScope.this()
	}
}
