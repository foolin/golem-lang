// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package compiler

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/mjarmy/golem-lang/ast"
	g "github.com/mjarmy/golem-lang/core"
	bc "github.com/mjarmy/golem-lang/core/bytecode"
)

//---------------------------------------------------------------
// The Golem Compiler
//---------------------------------------------------------------

// Compiler compiles an AST into bc
type Compiler interface {
	ast.Visitor

	Compile() *bc.Module
}

type compiler struct {
	poolBuilder *poolBuilder
	builtInMgr  g.BuiltinManager
	mod         *bc.Module

	funcs   []*ast.FnExpr
	funcIdx int

	btc      []byte
	lnum     []bc.LineNumberEntry
	handlers []*bc.ErrorHandler
}

// NewCompiler creates a new Compiler
func NewCompiler(
	builtInMgr g.BuiltinManager,
	astMod *ast.Module) Compiler {

	mod := &bc.Module{
		Name: astMod.Name,
		Path: astMod.Path,
	}

	// the 'init' function is always the first function in the list
	funcs := []*ast.FnExpr{astMod.InitFunc}

	return &compiler{
		builtInMgr:  builtInMgr,
		poolBuilder: newPoolBuilder(),
		mod:         mod,
		funcs:       funcs,
		funcIdx:     0,
		btc:         nil,
		lnum:        nil,
		handlers:    nil,
	}
}

// Compile compiles a Module
func (c *compiler) Compile() *bc.Module {

	// compile all the funcs
	for c.funcIdx < len(c.funcs) {
		c.poolBuilder.addTemplate(c.compileFunc(c.funcs[c.funcIdx]))
		c.funcIdx++
	}

	// done
	c.mod.Pool = c.poolBuilder.build()
	c.mod.Contents = c.makeModuleContents()
	return c.mod
}

func (c *compiler) makeModuleContents() g.Struct {

	fields := make(map[string]g.Field)

	stmts := c.funcs[0].Body.Statements
	for _, st := range stmts {
		switch t := st.(type) {
		case *ast.LetStmt:
			for _, d := range t.Decls {
				name := d.Ident.Symbol.Text
				vbl := d.Ident.Variable
				fields[name] = c.makeModuleProperty(vbl.Index(), vbl.IsConst())
			}
		case *ast.ConstStmt:
			for _, d := range t.Decls {
				name := d.Ident.Symbol.Text
				vbl := d.Ident.Variable
				fields[name] = c.makeModuleProperty(vbl.Index(), vbl.IsConst())
			}
		case *ast.NamedFnStmt:
			name := t.Ident.Symbol.Text
			vbl := t.Ident.Variable
			fields[name] = c.makeModuleProperty(vbl.Index(), vbl.IsConst())
		}
	}

	stc, err := g.NewFieldStruct(fields)
	g.Assert(err == nil)
	return stc
}

func (c *compiler) makeModuleProperty(index int, isConst bool) g.Field {

	get := g.NewFixedNativeFunc(
		[]g.Type{}, false,
		func(ev g.Eval, values []g.Value) (g.Value, g.Error) {
			return c.mod.Refs[index].Val, nil
		})

	if isConst {
		prop, err := g.NewReadonlyProperty(get)
		if err != nil {
			panic("unreachable")
		}
		return prop
	}

	set := g.NewFixedNativeFunc(
		[]g.Type{g.AnyType}, false,
		func(ev g.Eval, values []g.Value) (g.Value, g.Error) {
			c.mod.Refs[index].Val = values[0]
			return g.Null, nil
		})
	prop, err := g.NewProperty(get, set)
	if err != nil {
		panic("unreachable")
	}
	return prop
}

func makeArity(fe *ast.FnExpr) (g.Arity, []g.Value) {

	if fe.Variadic != nil {
		return g.Arity{
			Kind:     g.VariadicArity,
			Required: uint16(len(fe.Required)),
			Optional: 0,
		}, nil
	}

	if fe.Optional != nil {

		opt := make([]g.Value, len(fe.Optional))
		for i, o := range fe.Optional {
			opt[i] = toBasicValue(o.Value)
		}

		return g.Arity{
			Kind:     g.MultipleArity,
			Required: uint16(len(fe.Required)),
			Optional: uint16(len(fe.Optional)),
		}, opt
	}

	return g.Arity{
		Kind:     g.FixedArity,
		Required: uint16(len(fe.Required)),
		Optional: 0,
	}, nil
}

func (c *compiler) compileFunc(fe *ast.FnExpr) *bc.FuncTemplate {

	arity, optional := makeArity(fe)

	tpl := &bc.FuncTemplate{
		Module:          c.mod,
		Arity:           arity,
		OptionalParams:  optional,
		NumCaptures:     fe.Scope.NumCaptures(),
		NumLocals:       fe.Scope.NumLocals(),
		Bytecodes:       nil,
		LineNumberTable: nil,
		ErrorHandlers:   nil,
	}

	// reset template info for current func
	c.btc = []byte{}
	c.lnum = []bc.LineNumberEntry{}
	c.handlers = []*bc.ErrorHandler{}

	// TODO LoadNull and ReturnStmt are workarounds for the fact that
	// we have not yet written a Control Flow Graph
	c.push(ast.Pos{}, bc.LoadNull)
	c.Visit(fe.Body)
	c.push(ast.Pos{}, bc.Return)

	// save template info
	tpl.Bytecodes = c.btc
	tpl.LineNumberTable = c.lnum
	tpl.ErrorHandlers = c.handlers

	return tpl
}

func (c *compiler) Visit(node ast.Node) {
	switch t := node.(type) {

	case *ast.BlockNode:
		c.visitBlock(t)

	case *ast.ImportStmt:
		c.visitImport(t)

	case *ast.ConstStmt:
		c.visitDecls(t.Decls)

	case *ast.LetStmt:
		c.visitDecls(t.Decls)

	case *ast.NamedFnStmt:
		c.visitNamedFn(t)

	case *ast.AssignmentExpr:
		c.visitAssignment(t)

	case *ast.IfStmt:
		c.visitIf(t)

	case *ast.WhileStmt:
		c.visitWhile(t)

	case *ast.ForStmt:
		c.visitFor(t)

	case *ast.SwitchStmt:
		c.visitSwitch(t)

	case *ast.BreakStmt:
		c.visitBreak(t)

	case *ast.ContinueStmt:
		c.visitContinue(t)

	case *ast.ReturnStmt:
		c.visitReturn(t)

	case *ast.TryStmt:
		c.visitTry(t)

	case *ast.ThrowStmt:
		c.visitThrow(t)

	case *ast.TernaryExpr:
		c.visitTernaryExpr(t)

	case *ast.BinaryExpr:
		c.visitBinaryExpr(t)

	case *ast.UnaryExpr:
		c.visitUnaryExpr(t)

	case *ast.PostfixExpr:
		c.visitPostfixExpr(t)

	case *ast.BasicExpr:
		c.visitBasicExpr(t)

	case *ast.IdentExpr:
		c.visitIdentExpr(t)

	case *ast.BuiltinExpr:
		c.visitBuiltinExpr(t)

	case *ast.FnExpr:
		c.visitFunc(t)

	case *ast.InvokeExpr:
		c.visitInvoke(t)

	case *ast.GoStmt:
		c.visitGo(t)

	case *ast.ExprStmt:
		c.visitExprStmt(t)

	case *ast.StructExpr:
		c.visitStructExpr(t)

	case *ast.ThisExpr:
		c.visitThisExpr(t)

	case *ast.FieldExpr:
		c.visitFieldExpr(t)

	case *ast.IndexExpr:
		c.visitIndexExpr(t)

	case *ast.SliceExpr:
		c.visitSliceExpr(t)

	case *ast.SliceFromExpr:
		c.visitSliceFromExpr(t)

	case *ast.SliceToExpr:
		c.visitSliceToExpr(t)

	case *ast.ListExpr:
		c.visitListExpr(t)

	case *ast.SetExpr:
		c.visitSetExpr(t)

	case *ast.TupleExpr:
		c.visitTupleExpr(t)

	case *ast.DictExpr:
		c.visitDictExpr(t)

	default:
		panic(fmt.Sprintf("cannot compile %v\n", node))
	}
}

func (c *compiler) visitBlock(blk *ast.BlockNode) {

	// TODO A 'standalone' expression is an expression that is evaluated
	// but whose result is never assigned.  The *last* of these type
	// of expressions that is evaluated at runtime should be left on the
	// stack, since it could end up being used as an implicit return value.
	// The rest of them must be popped once they've been evaluated, so we
	// don't fill up the stack with un-needed values
	//
	// However, at the moment we do not have a Control Flow Graph, and thus
	// have no way of knowing which expressions should be popped.
	// We need to write the Control Flow Graph to fix this problem.

	for _, stmt := range blk.Statements {
		c.Visit(stmt)

		// TODO
		//if (node is ast.Expression) && someControlFlowGraphCheck() {
		//	c.push(node.End(), g.Pop)
		//}
	}
}

func (c *compiler) visitDecls(decls []*ast.DeclNode) {

	for _, d := range decls {
		if d.Val == nil {
			c.push(d.Ident.Begin(), bc.LoadNull)
		} else {
			c.Visit(d.Val)
		}

		c.assignIdent(d.Ident)
	}
}

func (c *compiler) visitImport(imp *ast.ImportStmt) {

	for _, ident := range imp.Idents {

		// push the module onto the stack
		sym := ident.Symbol.Text
		c.pushBytecode(
			ident.Begin(),
			bc.ImportModule,
			c.poolBuilder.constIndex(mustStr(sym)))

		// store module in identifer
		v := ident.Variable
		c.pushBytecode(ident.Begin(), bc.StoreLocal, v.Index())
	}
}

func (c *compiler) assignIdent(ident *ast.IdentExpr) {

	v := ident.Variable
	if v.IsCapture() {
		c.pushBytecode(ident.Begin(), bc.StoreCapture, v.Index())
	} else {
		c.pushBytecode(ident.Begin(), bc.StoreLocal, v.Index())
	}
}

func (c *compiler) visitNamedFn(nf *ast.NamedFnStmt) {

	c.Visit(nf.Func)

	v := nf.Ident.Variable
	g.Assert(!v.IsCapture())
	c.pushBytecode(nf.Ident.Begin(), bc.StoreLocal, v.Index())
}

func (c *compiler) visitAssignment(asn *ast.AssignmentExpr) {

	switch t := asn.Assignee.(type) {

	case *ast.IdentExpr:

		c.Visit(asn.Val)
		c.push(asn.Eq.Position, bc.Dup)
		c.assignIdent(t)

	case *ast.FieldExpr:

		c.Visit(t.Operand)
		c.Visit(asn.Val)
		c.pushBytecode(
			t.Key.Position,
			bc.SetField,
			c.poolBuilder.constIndex(mustStr(t.Key.Text)))

	case *ast.IndexExpr:

		c.Visit(t.Operand)
		c.Visit(t.Index)
		c.Visit(asn.Val)
		c.push(t.Index.Begin(), bc.SetIndex)

	default:
		panic("invalid assignee type")
	}
}

func (c *compiler) visitPostfixExpr(pe *ast.PostfixExpr) {

	switch t := pe.Assignee.(type) {

	case *ast.IdentExpr:

		c.visitIdentExpr(t)
		c.push(t.Begin(), bc.Dup)

		switch pe.Op.Text {
		case "++":
			c.push(pe.Op.Position, bc.LoadOne)
		case "--":
			c.push(pe.Op.Position, bc.LoadNegOne)
		default:
			panic("invalid postfix operator")
		}

		c.push(pe.Op.Position, bc.Inc)
		c.assignIdent(t)

	case *ast.FieldExpr:

		c.Visit(t.Operand)

		switch pe.Op.Text {
		case "++":
			c.push(pe.Op.Position, bc.LoadOne)
		case "--":
			c.push(pe.Op.Position, bc.LoadNegOne)
		default:
			panic("invalid postfix operator")
		}

		c.pushBytecode(
			t.Key.Position,
			bc.IncField,
			c.poolBuilder.constIndex(mustStr(t.Key.Text)))

	case *ast.IndexExpr:

		c.Visit(t.Operand)
		c.Visit(t.Index)

		switch pe.Op.Text {
		case "++":
			c.push(pe.Op.Position, bc.LoadOne)
		case "--":
			c.push(pe.Op.Position, bc.LoadNegOne)
		default:
			panic("invalid postfix operator")
		}

		c.push(t.Index.Begin(), bc.IncIndex)

	default:
		panic("invalid assignee type")
	}
}

func (c *compiler) visitIf(f *ast.IfStmt) {

	c.Visit(f.Cond)

	j0 := c.push(f.Cond.End(), bc.JumpFalse, 0xFF, 0xFF)
	c.Visit(f.Then)

	if f.Else == nil {

		c.setJump(j0, c.btcLen())

	} else {

		j1 := c.push(f.Else.Begin(), bc.Jump, 0xFF, 0xFF)
		c.setJump(j0, c.btcLen())

		c.Visit(f.Else)
		c.setJump(j1, c.btcLen())
	}
}

func (c *compiler) visitTernaryExpr(f *ast.TernaryExpr) {

	c.Visit(f.Cond)
	j0 := c.push(f.Cond.End(), bc.JumpFalse, 0xFF, 0xFF)

	c.Visit(f.Then)
	j1 := c.push(f.Else.Begin(), bc.Jump, 0xFF, 0xFF)
	c.setJump(j0, c.btcLen())

	c.Visit(f.Else)
	c.setJump(j1, c.btcLen())
}

func (c *compiler) visitWhile(w *ast.WhileStmt) {

	begin := c.btcLen()
	c.Visit(w.Cond)
	j0 := c.push(w.Cond.End(), bc.JumpFalse, 0xFF, 0xFF)

	body := c.btcLen()
	c.Visit(w.Body)
	c.push(w.Body.End(), bc.Jump, begin.high, begin.low)

	end := c.btcLen()
	c.setJump(j0, end)

	c.fixBreakContinue(begin, body, end)
}

func (c *compiler) visitFor(f *ast.ForStmt) {

	tok := f.Iterable.Begin()
	idx := f.IterableIdent.Variable.Index()

	// put Iterable expression on stack
	c.Visit(f.Iterable)

	// call NewIterator()
	c.push(tok, bc.NewIter)

	// store iterator
	c.pushBytecode(tok, bc.StoreLocal, idx)

	// top of loop: load iterator and call IterNext()
	begin := c.btcLen()
	c.pushBytecode(tok, bc.LoadLocal, idx)
	c.push(tok, bc.IterNext)
	j0 := c.push(tok, bc.JumpFalse, 0xFF, 0xFF)

	// load iterator and call IterGet()
	c.pushBytecode(tok, bc.LoadLocal, idx)
	c.push(tok, bc.IterGet)

	if len(f.Idents) == 1 {
		// perform StoreLocal on the current item
		ident := f.Idents[0]
		c.pushBytecode(ident.Begin(), bc.StoreLocal, ident.Variable.Index())
	} else {
		// make sure the current item is really a tuple,
		// and is of the proper length
		c.pushBytecode(tok, bc.CheckTuple, len(f.Idents))

		// perform StoreLocal on each tuple element
		for i, ident := range f.Idents {
			c.push(tok, bc.Dup)
			c.pushInt(tok, int64(i))
			c.push(tok, bc.GetIndex)
			c.pushBytecode(ident.Begin(), bc.StoreLocal, ident.Variable.Index())
		}

		// pop the tuple
		c.push(tok, bc.Pop)
	}

	// compile the body
	body := c.btcLen()
	c.Visit(f.Body)
	c.push(f.Body.End(), bc.Jump, begin.high, begin.low)

	// jump to top of loop
	end := c.btcLen()
	c.setJump(j0, end)

	c.fixBreakContinue(begin, body, end)
}

func (c *compiler) fixBreakContinue(begin instPtr, body instPtr, end instPtr) {

	// replace BreakStmt and ContinueStmt with Jump
	for i := body.ip; i < end.ip; {
		switch c.btc[i] {
		case bc.Break:
			c.btc[i] = bc.Jump
			c.btc[i+1] = end.high
			c.btc[i+2] = end.low
		case bc.Continue:
			c.btc[i] = bc.Jump
			c.btc[i+1] = begin.high
			c.btc[i+2] = begin.low
		}
		i += bc.Size(c.btc[i])
	}
}

func (c *compiler) visitBreak(br *ast.BreakStmt) {
	c.push(br.Begin(), bc.Break, 0xFF, 0xFF)
}

func (c *compiler) visitContinue(cn *ast.ContinueStmt) {
	c.push(cn.Begin(), bc.Continue, 0xFF, 0xFF)
}

func (c *compiler) visitSwitch(sw *ast.SwitchStmt) {

	// visit the item, if there is one
	hasItem := false
	if sw.Item != nil {
		hasItem = true
		c.Visit(sw.Item)
	}

	// visit each case
	endJumps := []int{}
	for _, cs := range sw.Cases {
		endJumps = append(endJumps, c.visitCase(cs, hasItem))
	}

	// visit default
	if sw.DefaultNode != nil {
		for _, n := range sw.DefaultNode.Body {
			c.Visit(n)
		}
	}

	// if there is an item, pop it
	if hasItem {
		c.push(sw.End(), bc.Pop)
	}

	// set all the end jumps
	for _, j := range endJumps {
		c.setJump(j, c.btcLen())
	}
}

func (c *compiler) visitCase(cs *ast.CaseNode, hasItem bool) int {

	bodyJumps := []int{}

	// visit each match, and jump to body if true
	for _, m := range cs.Matches {

		if hasItem {
			// if there is an item, Dup it and do an Eq comparison against the match
			c.push(m.Begin(), bc.Dup)
			c.Visit(m)
			c.push(m.Begin(), bc.Eq)
		} else {
			// otherwise, evaluate the match and assume its a Bool
			c.Visit(m)
		}

		bodyJumps = append(bodyJumps, c.push(m.End(), bc.JumpTrue, 0xFF, 0xFF))
	}

	// no match -- jump to the end of the case
	caseEndJump := c.push(cs.End(), bc.Jump, 0xFF, 0xFF)

	// set all the body jumps
	for _, j := range bodyJumps {
		c.setJump(j, c.btcLen())
	}

	// visit body, and then push a jump to the very end of the switch
	for _, n := range cs.Body {
		c.Visit(n)
	}
	endJump := c.push(cs.End(), bc.Jump, 0xFF, 0xFF)

	// set the jump to the end of the case
	c.setJump(caseEndJump, c.btcLen())

	// return the jump to end of the switch
	return endJump
}

func (c *compiler) visitReturn(rt *ast.ReturnStmt) {
	c.Visit(rt.Val)
	c.push(rt.Begin(), bc.Return)
}

func (c *compiler) compileTryBlock(block *ast.BlockNode) {

	c.pushBytecode(block.Begin(), bc.PushTry, len(c.handlers))
	c.handlers = append(c.handlers, &bc.ErrorHandler{})

	c.Visit(block)
	c.push(block.End(), bc.PopTry)
}

func (c *compiler) compileCatchBlock(
	tryEnd ast.Pos,
	ident *ast.IdentExpr,
	block *ast.BlockNode) int {

	// push a jump, so that we'll skip the catch block during normal execution
	skipCatch := c.push(tryEnd, bc.Jump, 0xFF, 0xFF)

	begin := len(c.btc)

	// store the error (which the interpreter has put on the stack
	// for us as part of the error recovery process)
	v := ident.Variable
	g.Assert(!v.IsCapture())
	c.pushBytecode(ident.Begin(), bc.StoreLocal, v.Index())

	// compile the catch
	c.Visit(block)
	c.push(block.End(), bc.TryDone)

	// fix the jump
	c.setJump(skipCatch, c.btcLen())

	return begin
}

func (c *compiler) compileFinallyBlock(block *ast.BlockNode) int {

	begin := len(c.btc)
	c.Visit(block)
	c.push(block.End(), bc.TryDone)
	return begin
}

func (c *compiler) visitTry(t *ast.TryStmt) {

	// try
	c.compileTryBlock(t.TryBlock)

	begin := len(c.btc)

	// catch
	catch := -1
	if t.CatchBlock != nil {
		catch = c.compileCatchBlock(
			t.TryBlock.End(), t.CatchIdent, t.CatchBlock)
	}

	// finally
	finally := -1
	if t.FinallyBlock != nil {
		finally = c.compileFinallyBlock(t.FinallyBlock)
	}

	end := len(c.btc)

	// replace Return with TryReturn
	ip := begin
	for ip < end {
		if c.btc[ip] == bc.Return {
			c.btc[ip] = bc.TryReturn
		}

		ip += bc.Size(c.btc[ip])
	}

	// done
	g.Assert(!(catch == -1 && finally == -1)) // sanity check
	handler := c.handlers[len(c.handlers)-1]
	handler.Catch = catch
	handler.Finally = finally
	handler.End = end
}

func (c *compiler) visitThrow(t *ast.ThrowStmt) {
	c.Visit(t.Val)
	c.push(t.End(), bc.Throw)
}

func (c *compiler) visitBinaryExpr(b *ast.BinaryExpr) {

	switch b.Op.Kind {

	case ast.DoublePipe:
		c.visitOr(b.LHS, b.RHS)
	case ast.DoubleAmp:
		c.visitAnd(b.LHS, b.RHS)

	case ast.DoubleEq:
		b.Traverse(c)
		c.push(b.Op.Position, bc.Eq)
	case ast.NotEq:
		b.Traverse(c)
		c.push(b.Op.Position, bc.Ne)

	case ast.Gt:
		b.Traverse(c)
		c.push(b.Op.Position, bc.Gt)
	case ast.GtEq:
		b.Traverse(c)
		c.push(b.Op.Position, bc.Gte)
	case ast.Lt:
		b.Traverse(c)
		c.push(b.Op.Position, bc.Lt)
	case ast.LtEq:
		b.Traverse(c)
		c.push(b.Op.Position, bc.Lte)
	case ast.Cmp:
		b.Traverse(c)
		c.push(b.Op.Position, bc.Cmp)

	case ast.Plus:
		b.Traverse(c)
		c.push(b.Op.Position, bc.Plus)
	case ast.Minus:
		b.Traverse(c)
		c.push(b.Op.Position, bc.Sub)
	case ast.Star:
		b.Traverse(c)
		c.push(b.Op.Position, bc.Mul)
	case ast.Slash:
		b.Traverse(c)
		c.push(b.Op.Position, bc.Div)

	case ast.Percent:
		b.Traverse(c)
		c.push(b.Op.Position, bc.Rem)
	case ast.Amp:
		b.Traverse(c)
		c.push(b.Op.Position, bc.BitAnd)
	case ast.Pipe:
		b.Traverse(c)
		c.push(b.Op.Position, bc.BitOr)
	case ast.Caret:
		b.Traverse(c)
		c.push(b.Op.Position, bc.BitXor)
	case ast.DoubleLt:
		b.Traverse(c)
		c.push(b.Op.Position, bc.LeftShift)
	case ast.DoubleGt:
		b.Traverse(c)
		c.push(b.Op.Position, bc.RightShift)

	default:
		panic("unreachable")
	}
}

func (c *compiler) visitOr(lhs ast.Expression, rhs ast.Expression) {

	c.Visit(lhs)
	j0 := c.push(lhs.End(), bc.JumpTrue, 0xFF, 0xFF)

	c.Visit(rhs)
	j1 := c.push(rhs.End(), bc.JumpFalse, 0xFF, 0xFF)

	c.setJump(j0, c.btcLen())
	c.push(rhs.End(), bc.LoadTrue)
	j2 := c.push(rhs.End(), bc.Jump, 0xFF, 0xFF)

	c.setJump(j1, c.btcLen())
	c.push(rhs.End(), bc.LoadFalse)

	c.setJump(j2, c.btcLen())
}

func (c *compiler) visitAnd(lhs ast.Expression, rhs ast.Expression) {

	c.Visit(lhs)
	j0 := c.push(lhs.End(), bc.JumpFalse, 0xFF, 0xFF)

	c.Visit(rhs)
	j1 := c.push(rhs.End(), bc.JumpFalse, 0xFF, 0xFF)

	c.push(rhs.End(), bc.LoadTrue)
	j2 := c.push(rhs.End(), bc.Jump, 0xFF, 0xFF)

	c.setJump(j0, c.btcLen())
	c.setJump(j1, c.btcLen())
	c.push(rhs.End(), bc.LoadFalse)

	c.setJump(j2, c.btcLen())
}

func (c *compiler) visitUnaryExpr(u *ast.UnaryExpr) {

	switch u.Op.Kind {
	case ast.Minus:
		opn := u.Operand

		switch t := opn.(type) {
		case *ast.BasicExpr:
			switch t.Token.Kind {

			case ast.Int:
				i := parseInt(t.Token.Text)
				switch i {
				case 0:
					c.push(u.Op.Position, bc.LoadZero)
				case 1:
					c.push(u.Op.Position, bc.LoadNegOne)
				default:
					c.pushBytecode(
						u.Op.Position,
						bc.LoadConst,
						c.poolBuilder.constIndex(g.NewInt(-i)))
				}

			default:
				c.Visit(u.Operand)
				c.push(u.Op.Position, bc.Negate)
			}
		default:
			c.Visit(u.Operand)
			c.push(u.Op.Position, bc.Negate)
		}

	case ast.Not:
		c.Visit(u.Operand)
		c.push(u.Op.Position, bc.Not)

	case ast.Tilde:
		c.Visit(u.Operand)
		c.push(u.Op.Position, bc.Complement)

	default:
		panic("unreachable")
	}
}

func (c *compiler) visitBasicExpr(basic *ast.BasicExpr) {

	switch basic.Token.Kind {

	case ast.Null:
		c.push(basic.Token.Position, bc.LoadNull)

	case ast.True:
		c.push(basic.Token.Position, bc.LoadTrue)

	case ast.False:
		c.push(basic.Token.Position, bc.LoadFalse)

	case ast.Str:
		c.pushBytecode(
			basic.Token.Position,
			bc.LoadConst,
			c.poolBuilder.constIndex(mustStr(basic.Token.Text)))

	case ast.Int:
		c.pushInt(
			basic.Token.Position,
			parseInt(basic.Token.Text))

	case ast.Float:
		f := parseFloat(basic.Token.Text)
		c.pushBytecode(
			basic.Token.Position,
			bc.LoadConst,
			c.poolBuilder.constIndex(g.NewFloat(f)))

	default:
		panic("unreachable")
	}
}

func toBasicValue(basic *ast.BasicExpr) g.Value {

	switch basic.Token.Kind {

	case ast.Null:
		return g.Null

	case ast.True:
		return g.True

	case ast.False:
		return g.False

	case ast.Str:
		return mustStr(basic.Token.Text)

	case ast.Int:
		return g.NewInt(parseInt(basic.Token.Text))

	case ast.Float:
		return g.NewFloat(parseFloat(basic.Token.Text))

	default:
		panic("unreachable")
	}
}

func (c *compiler) visitIdentExpr(ident *ast.IdentExpr) {
	v := ident.Variable
	if v.IsCapture() {
		c.pushBytecode(ident.Begin(), bc.LoadCapture, v.Index())
	} else {
		c.pushBytecode(ident.Begin(), bc.LoadLocal, v.Index())
	}
}

func (c *compiler) visitBuiltinExpr(blt *ast.BuiltinExpr) {

	c.pushBytecode(blt.Fn.Position, bc.LoadBuiltin, c.builtInMgr.IndexOf(blt.Fn.Text))
}

func (c *compiler) visitFunc(fe *ast.FnExpr) {

	c.pushBytecode(fe.Begin(), bc.NewFunc, len(c.funcs))

	parents := getSortedCaptureParents(fe.Scope)
	for _, vbl := range parents {
		if vbl.IsCapture() {
			c.pushBytecode(fe.Begin(), bc.FuncCapture, vbl.Index())
		} else {
			c.pushBytecode(fe.Begin(), bc.FuncLocal, vbl.Index())
		}
	}

	c.funcs = append(c.funcs, fe)
}

func (c *compiler) visitInvoke(inv *ast.InvokeExpr) {

	// InvokeField
	if fe, ok := inv.Operand.(*ast.FieldExpr); ok {

		c.Visit(fe.Operand)
		for _, n := range inv.Params {
			c.Visit(n)
		}

		// push the field index, and number of params
		c.pushWideBytecode(
			fe.Key.Position,
			bc.InvokeField,
			c.poolBuilder.constIndex(mustStr(fe.Key.Text)),
			len(inv.Params))
		return
	}

	// Invoke
	c.Visit(inv.Operand)
	for _, n := range inv.Params {
		c.Visit(n)
	}
	// push the number of params
	c.pushBytecode(inv.Begin(), bc.Invoke, len(inv.Params))

}

func (c *compiler) visitGo(gw *ast.GoStmt) {

	inv := gw.Invocation
	c.Visit(inv.Operand)
	for _, n := range inv.Params {
		c.Visit(n)
	}
	c.pushBytecode(inv.Begin(), bc.Go, len(inv.Params))
}

func (c *compiler) visitExprStmt(es *ast.ExprStmt) {
	c.Visit(es.Expr)
}

func (c *compiler) visitStructExpr(stc *ast.StructExpr) {

	// add struct def to pool
	def := make([]string, len(stc.Keys))
	for i, k := range stc.Keys {
		def[i] = k.Text
	}
	defIdx := c.poolBuilder.structDefIndex(def)

	// create new struct
	c.pushBytecode(stc.Begin(), bc.NewStruct, defIdx)

	// if the struct is referenced by a 'this', then store local
	if this, ok := stc.Scope.GetVariable("this"); ok {
		c.push(stc.Begin(), bc.Dup)
		c.pushBytecode(stc.Begin(), bc.StoreLocal, this.Index())
	}

	// init each field
	for i, k := range stc.Keys {

		v := stc.Values[i]
		if p, ok := v.(*ast.PropNode); ok {

			if p.Set == nil {

				// InitReadonlyProperty
				c.Visit(p.Get)
				c.pushBytecode(
					v.Begin(),
					bc.InitReadonlyProperty,
					c.poolBuilder.constIndex(mustStr(k.Text)))

			} else {

				// InitProperty
				c.Visit(p.Get)
				c.Visit(p.Set)
				c.pushBytecode(
					v.Begin(),
					bc.InitProperty,
					c.poolBuilder.constIndex(mustStr(k.Text)))
			}
		} else {

			// InitField
			c.Visit(v)
			c.pushBytecode(
				v.Begin(),
				bc.InitField,
				c.poolBuilder.constIndex(mustStr(k.Text)))
		}
	}
}

func (c *compiler) visitThisExpr(this *ast.ThisExpr) {

	v := this.Variable
	if v.IsCapture() {
		c.pushBytecode(this.Begin(), bc.LoadCapture, v.Index())
	} else {
		c.pushBytecode(this.Begin(), bc.LoadLocal, v.Index())
	}
}

func (c *compiler) visitFieldExpr(fe *ast.FieldExpr) {

	c.Visit(fe.Operand)

	// push the field index
	c.pushBytecode(
		fe.Key.Position,
		bc.GetField,
		c.poolBuilder.constIndex(mustStr(fe.Key.Text)))
}

func (c *compiler) visitIndexExpr(ie *ast.IndexExpr) {
	c.Visit(ie.Operand)
	c.Visit(ie.Index)
	c.push(ie.Index.Begin(), bc.GetIndex)
}

func (c *compiler) visitSliceExpr(s *ast.SliceExpr) {
	c.Visit(s.Operand)
	c.Visit(s.From)
	c.Visit(s.To)
	c.push(s.From.Begin(), bc.Slice)
}

func (c *compiler) visitSliceFromExpr(s *ast.SliceFromExpr) {
	c.Visit(s.Operand)
	c.Visit(s.From)
	c.push(s.From.Begin(), bc.SliceFrom)
}

func (c *compiler) visitSliceToExpr(s *ast.SliceToExpr) {
	c.Visit(s.Operand)
	c.Visit(s.To)
	c.push(s.To.Begin(), bc.SliceTo)
}

func (c *compiler) visitListExpr(ls *ast.ListExpr) {

	for _, v := range ls.Elems {
		c.Visit(v)
	}
	c.pushBytecode(ls.Begin(), bc.NewList, len(ls.Elems))
}

func (c *compiler) visitSetExpr(s *ast.SetExpr) {

	for _, v := range s.Elems {
		c.Visit(v)
	}
	c.pushBytecode(s.Begin(), bc.NewSet, len(s.Elems))
}

func (c *compiler) visitTupleExpr(tp *ast.TupleExpr) {

	for _, v := range tp.Elems {
		c.Visit(v)
	}
	c.pushBytecode(tp.Begin(), bc.NewTuple, len(tp.Elems))
}

func (c *compiler) visitDictExpr(d *ast.DictExpr) {

	for _, de := range d.Entries {
		c.Visit(de.Key)
		c.Visit(de.Value)
	}

	c.pushBytecode(d.Begin(), bc.NewDict, len(d.Entries))
}

func (c *compiler) pushInt(pos ast.Pos, i int64) {
	switch i {
	case 0:
		c.push(pos, bc.LoadZero)
	case 1:
		c.push(pos, bc.LoadOne)
	default:
		c.pushBytecode(
			pos,
			bc.LoadConst,
			c.poolBuilder.constIndex(g.NewInt(i)))
	}
}

// returns the length of btc *before* the bytes are pushed
func (c *compiler) push(pos ast.Pos, bytes ...byte) int {

	n := len(c.btc)
	c.btc = append(c.btc, bytes...)

	ln := len(c.lnum)
	if (ln == 0) || (pos.Line != c.lnum[ln-1].LineNum) {
		c.lnum = append(c.lnum, bc.LineNumberEntry{
			Index:   n,
			LineNum: pos.Line,
		})
	}

	return n
}

// push a 3-byte bc
func (c *compiler) pushBytecode(pos ast.Pos, code byte, p int) int {
	high, low := bc.EncodeParam(p)
	return c.push(pos, code, high, low)
}

// push a 5-byte bc
func (c *compiler) pushWideBytecode(pos ast.Pos, code byte, p, q int) int {
	high1, low1, high2, low2 := bc.EncodeWideParams(p, q)
	return c.push(pos, code, high1, low1, high2, low2)
}

// replace a mocked-up jump value with the 'real' destination
func (c *compiler) setJump(jmp int, dest instPtr) {
	c.btc[jmp+1] = dest.high
	c.btc[jmp+2] = dest.low
}

func (c *compiler) btcLen() instPtr {
	high, low := bc.EncodeParam(len(c.btc))
	return instPtr{len(c.btc), high, low}
}

//--------------------------------------------------------------
// misc

type instPtr struct {
	ip   int
	high byte
	low  byte
}

func parseInt(text string) int64 {
	i, err := strconv.ParseInt(text, 10, 64)
	g.Assert(err == nil)
	g.Assert(i >= 0)
	return i
}

func parseFloat(text string) float64 {
	f, err := strconv.ParseFloat(text, 64)
	g.Assert(err == nil)
	g.Assert(f >= 0)
	return f
}

//--------------------------------------------------------------
// capture

func getSortedCaptureParents(f ast.FuncScope) []ast.Variable {

	// Sort the captures by child index
	caps := []ast.Capture{}
	caps = append(caps, f.GetCaptures()...)
	sort.Slice(caps, func(i, j int) bool {
		return caps[i].Child().Index() < caps[j].Child().Index()
	})

	// Use the sorted list to create the proper ordering of parents
	parents := []ast.Variable{}
	for _, c := range caps {
		parents = append(parents, c.Parent())
	}
	return parents
}

//--------------------------------------------------------------

// This is impossible, because every string in the AST is
// guaranteed by the scanner to be UTF-8.
func mustStr(s string) g.Str {
	sv, err := g.NewStr(s)
	if err != nil {
		panic("unreachable")
	}
	return sv
}
