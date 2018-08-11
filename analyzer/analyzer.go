// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package analyzer

import (
	"fmt"
	"github.com/mjarmy/golem-lang/ast"
)

// Analyzer analyzes an AST.
type Analyzer interface {
	ast.Visitor
	Module() *ast.FnExpr
	Analyze() []error
}

type analyzer struct {
	modName string
	modPath string
	mod     *ast.FnExpr
	scopes  []ast.Scope
	loops   []ast.Loop
	structs []*ast.StructExpr
	errors  []error
}

// NewAnalyzer creates a new Analyzer
func NewAnalyzer(modName, modPath string, mod *ast.FnExpr) Analyzer {

	return &analyzer{modName, modPath, mod, []ast.Scope{mod.Scope}, []ast.Loop{}, []*ast.StructExpr{}, nil}
}

// Analyze analyzes an AST.
func (a *analyzer) Analyze() []error {

	// visit module block
	a.visitBlock(a.mod.Body)

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
			a.errors = append(a.errors,
				fmt.Errorf("'break' outside of loop, at %s:%v", a.modPath, t.Token.Position))
		}

	case *ast.ContinueStmt:
		if len(a.loops) == 0 {
			a.errors = append(a.errors,
				fmt.Errorf("'continue' outside of loop, at %s:%v", a.modPath, t.Token.Position))
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
	for _, ident := range imp.Idents {
		a.defineIdent(ident, true)
	}
}

func (a *analyzer) visitTry(t *ast.TryStmt) {

	a.Visit(t.TryBlock)
	if t.CatchToken != nil {
		a.pushScope(t.CatchScope)
		t.CatchIdent.Variable = a.putVariable(t.CatchIdent.Symbol.Text, true)

		a.Visit(t.CatchBlock)
		a.popScope()
	}
	if t.FinallyToken != nil {
		a.Visit(t.FinallyBlock)
	}
}

func (a *analyzer) defineIdent(ident *ast.IdentExpr, isConst bool) {
	sym := ident.Symbol.Text
	if _, ok := a.getVariable(sym); ok {
		a.errors = append(a.errors,
			fmt.Errorf("Symbol '%s' is already defined, at %s:%v", sym, a.modPath, ident.Symbol.Position))
	} else {
		ident.Variable = a.putVariable(sym, isConst)
	}
}

func (a *analyzer) visitBlock(blk *ast.BlockNode) {

	a.pushScope(blk.Scope)

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

	a.popScope()
}

func (a *analyzer) visitFor(fr *ast.ForStmt) {

	a.pushScope(fr.Scope)

	// define identifiers
	for _, ident := range fr.Idents {
		ident.Variable = a.putVariable(ident.Symbol.Text, false)
	}

	// define the identifier for the iterable
	a.defineIdent(fr.IterableIdent, false)

	// visit the iterable and body
	a.Visit(fr.Iterable)
	a.visitBlock(fr.Body)

	a.popScope()
}

func (a *analyzer) visitFunc(fn *ast.FnExpr) {

	a.pushScope(fn.Scope)

	// visit child nodes
	for _, f := range fn.FormalParams {
		f.Ident.Variable = a.putVariable(f.Ident.Symbol.Text, f.IsConst)
	}
	a.visitBlock(fn.Body)

	a.popScope()
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
	if v, ok := a.getVariable(sym); ok {
		if v.IsConst() {
			a.errors = append(a.errors,
				fmt.Errorf("Symbol '%s' is constant, at %s:%v", sym, a.modPath, ident.Symbol.Position))
		}
		ident.Variable = v
	} else {
		a.errors = append(a.errors,
			fmt.Errorf("Symbol '%s' is not defined, at %s:%v", sym, a.modPath, ident.Symbol.Position))
	}
}

func (a *analyzer) visitIdentExpr(ident *ast.IdentExpr) {

	sym := ident.Symbol.Text

	if v, ok := a.getVariable(sym); ok {
		ident.Variable = v
	} else {
		a.errors = append(a.errors,
			fmt.Errorf("Symbol '%s' is not defined, at %s:%v", sym, a.modPath, ident.Symbol.Position))
	}
}

func (a *analyzer) visitStructExpr(stc *ast.StructExpr) {
	a.structs = append(a.structs, stc)

	a.pushScope(stc.Scope)
	stc.Traverse(a)
	a.popScope()

	a.structs = a.structs[:len(a.structs)-1]
}

func (a *analyzer) visitThisExpr(this *ast.ThisExpr) {

	n := len(a.structs)
	if n == 0 {
		a.errors = append(a.errors,
			fmt.Errorf("'this' outside of struct, at %s:%v", a.modPath, this.Token.Position))
	} else {
		this.Variable = a.putThis()
	}
}

//-----------------------------------------------------------------------------
// Scope Management
//-----------------------------------------------------------------------------

func (a *analyzer) pushScope(scope ast.Scope) {
	a.scopes = append(a.scopes, scope)
}

func (a *analyzer) popScope() {
	a.scopes = a.scopes[:len(a.scopes)-1]
}

// Put a symbol entry into the current scope.
func (a *analyzer) putVariable(sym string, isConst bool) ast.Variable {

	s := a.scopes[len(a.scopes)-1]

	if _, ok := s.GetVariable(sym); ok {
		// its the caller's responsibility to ensure this never happens
		panic("symbol is already defined")
	}

	v := ast.NewVariable(sym, a.incrementNumLocals(len(a.scopes)-1), isConst, false)
	s.PutVariable(sym, v)
	return v
}

// Increment the number of local variables in the nearest parent FuncScope,
// staring from the given scope index.
func (a *analyzer) incrementNumLocals(from int) int {

	for i := from; i >= 0; i-- {
		s := a.scopes[i]

		if f, ok := s.(ast.FuncScope); ok {
			idx := f.NumLocals()
			f.IncrementNumLocals()
			if idx+1 >= (2 << 16) {
				panic("TODO wide index")
			}
			return idx
		}
	}

	panic("unreachable")
}

// Get a variable by walking up the scope stack, or return false if we can't find it.
func (a *analyzer) getVariable(sym string) (ast.Variable, bool) {

	// We must create a capture in each function that we skip over
	// while we were looking for the variable
	funcScopes := []ast.FuncScope{}

	for i := len(a.scopes) - 1; i >= 0; i-- {
		s := a.scopes[i]

		// we found the variable definition
		if v, ok := s.GetVariable(sym); ok {
			return a.applyCaptures(v, funcScopes), true
		}

		if f, ok := s.(ast.FuncScope); ok {
			// we found the variable in a function capture
			if cp, ok := f.GetCapture(sym); ok {
				return a.applyCaptures(cp.Child(), funcScopes), true
			}
			// Save the function so we can capture into it later.
			funcScopes = append(funcScopes, f)
		}
	}

	// The variable is not currently defined in any scope in the stack
	return nil, false
}

func (a *analyzer) applyCaptures(
	v ast.Variable,
	funcScopes []ast.FuncScope) ast.Variable {

	for i := len(funcScopes) - 1; i >= 0; i-- {
		v = funcScopes[i].PutCapture(v).Child()
	}
	return v
}

// This creates a variable for 'this', or returns an existing 'this' variable.
func (a *analyzer) putThis() ast.Variable {

	// walk up the stack, looking for a StructScope
	for i := len(a.scopes) - 1; i >= 0; i-- {
		s := a.scopes[i]

		if _, ok := s.(ast.StructScope); ok {

			// Define a 'this' variable on the structScope, if its not already defined.
			// NOTE: We increment the number of local variables
			// starting from the current scope, not the top of the stack.
			if _, ok := s.GetVariable("this"); !ok {
				s.PutVariable(
					"this",
					ast.NewVariable("this", a.incrementNumLocals(i), true, false))
			}

			// now call getVariable(), to trigger captures in any intervening functions.
			v, ok := a.getVariable("this")
			if !ok {
				panic("call to 'this' failed")
			}
			return v
		}
	}

	panic("unreachable")
}
