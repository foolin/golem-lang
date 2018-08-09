// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package compiler

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/mjarmy/golem-lang/analyzer"
	"github.com/mjarmy/golem-lang/ast"
	g "github.com/mjarmy/golem-lang/core"
	o "github.com/mjarmy/golem-lang/core/opcodes"
)

// Compiler compiles an AST into bytecode
type Compiler interface {
	ast.Visitor
	Compile() *g.Module
}

type compiler struct {
	builtInMgr g.BuiltinManager

	pool     *g.HashMap
	opc      []byte
	lnum     []g.LineNumberEntry
	handlers []g.ExceptionHandler

	funcs      []*ast.FnExpr
	templates  []*g.Template
	structDefs [][]*g.FieldDef
	idx        int
}

// NewCompiler creates a new Compiler
func NewCompiler(anl analyzer.Analyzer, builtInMgr g.BuiltinManager) Compiler {

	funcs := []*ast.FnExpr{anl.Module()}
	templates := []*g.Template{}
	structDefs := [][]*g.FieldDef{}

	return &compiler{
		builtInMgr, g.EmptyHashMap(), nil, nil, nil,
		funcs, templates, structDefs, 0}
}

func (c *compiler) Compile() *g.Module {

	// compile all the funcs
	for c.idx < len(c.funcs) {
		c.templates = append(
			c.templates,
			c.compileFunc(c.funcs[c.idx]))
		c.idx++
	}

	// done
	mod := &g.Module{
		Pool:       makePoolSlice(c.pool),
		Refs:       nil,
		StructDefs: c.structDefs,
		Templates:  c.templates,
		Contents:   nil,
	}
	mod.Contents = c.makeModuleContents(mod)
	return mod
}

func (c *compiler) makeModuleContents(mod *g.Module) g.Struct {

	entries := []g.Field{}
	stmts := c.funcs[0].Body.Statements
	for _, st := range stmts {
		switch t := st.(type) {
		case *ast.LetStmt:
			for _, d := range t.Decls {
				vbl := d.Ident.Variable
				entries = append(entries, c.makeModuleProperty(
					mod, d.Ident.Symbol.Text, vbl.Index(), vbl.IsConst()))
			}
		case *ast.ConstStmt:
			for _, d := range t.Decls {
				vbl := d.Ident.Variable
				entries = append(entries, c.makeModuleProperty(
					mod, d.Ident.Symbol.Text, vbl.Index(), vbl.IsConst()))
			}
		case *ast.NamedFnStmt:
			vbl := t.Ident.Variable
			entries = append(entries, c.makeModuleProperty(
				mod, t.Ident.Symbol.Text, vbl.Index(), vbl.IsConst()))
		}
	}

	stc, err := g.NewStruct(entries, false)
	assert(err == nil)
	return stc
}

func (c *compiler) makeModuleProperty(
	mod *g.Module,
	name string,
	refIndex int,
	isConst bool) g.Field {

	getter := g.NewNativeFunc0(
		func(cx g.Context) (g.Value, g.Error) {
			return mod.Refs[refIndex].Val, nil
		})
	if isConst {
		prop, err := g.NewReadonlyNativeProperty(name, getter)
		if err != nil {
			panic("unreachable")
		}
		return prop
	}

	setter := g.NewNativeFuncValue(
		func(cx g.Context, val g.Value) (g.Value, g.Error) {
			mod.Refs[refIndex].Val = val
			return g.Null, nil
		})
	prop, err := g.NewNativeProperty(name, getter, setter)
	if err != nil {
		panic("unreachable")
	}
	return prop
}

func (c *compiler) compileFunc(fe *ast.FnExpr) *g.Template {

	arity := len(fe.FormalParams)
	tpl := &g.Template{
		Arity:             arity,
		NumCaptures:       fe.Scope.NumCaptures(),
		NumLocals:         fe.Scope.NumLocals(),
		OpCodes:           nil,
		LineNumberTable:   nil,
		ExceptionHandlers: nil,
	}

	c.opc = []byte{}
	c.lnum = []g.LineNumberEntry{}
	c.handlers = []g.ExceptionHandler{}

	// TODO LoadNull and ReturnStmt are workarounds for the fact that
	// we have not yet written a Control Flow Graph
	c.push(ast.Pos{}, o.LoadNull)
	c.Visit(fe.Body)
	c.push(ast.Pos{}, o.ReturnStmt)

	tpl.OpCodes = c.opc
	tpl.LineNumberTable = c.lnum
	tpl.ExceptionHandlers = c.handlers

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

	case *ast.PropNode:
		c.visitPropNode(t)

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
	// So we need to write the Control Flow Graph to fix this problem.

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
			c.push(d.Ident.Begin(), o.LoadNull)
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
		c.pushIndex(
			ident.Begin(),
			o.ImportModule,
			poolIndex(c.pool, g.NewStr(sym)))

		// store module in identifer
		v := ident.Variable
		c.pushIndex(ident.Begin(), o.StoreLocal, v.Index())
	}
}

func (c *compiler) assignIdent(ident *ast.IdentExpr) {

	v := ident.Variable
	if v.IsCapture() {
		c.pushIndex(ident.Begin(), o.StoreCapture, v.Index())
	} else {
		c.pushIndex(ident.Begin(), o.StoreLocal, v.Index())
	}
}

func (c *compiler) visitNamedFn(nf *ast.NamedFnStmt) {

	c.Visit(nf.Func)

	v := nf.Ident.Variable
	assert(!v.IsCapture())
	c.pushIndex(nf.Ident.Begin(), o.StoreLocal, v.Index())
}

func (c *compiler) visitAssignment(asn *ast.AssignmentExpr) {

	switch t := asn.Assignee.(type) {

	case *ast.IdentExpr:

		c.Visit(asn.Val)
		c.push(asn.Eq.Position, o.Dup)
		c.assignIdent(t)

	case *ast.FieldExpr:

		c.Visit(t.Operand)
		c.Visit(asn.Val)
		c.pushIndex(
			t.Key.Position,
			o.SetField,
			poolIndex(c.pool, g.NewStr(t.Key.Text)))

	case *ast.IndexExpr:

		c.Visit(t.Operand)
		c.Visit(t.Index)
		c.Visit(asn.Val)
		c.push(t.Index.Begin(), o.SetIndex)

	default:
		panic("invalid assignee type")
	}
}

func (c *compiler) visitPostfixExpr(pe *ast.PostfixExpr) {

	switch t := pe.Assignee.(type) {

	case *ast.IdentExpr:

		c.visitIdentExpr(t)
		c.push(t.Begin(), o.Dup)

		switch pe.Op.Text {
		case "++":
			c.push(pe.Op.Position, o.LoadOne)
		case "--":
			c.push(pe.Op.Position, o.LoadNegOne)
		default:
			panic("invalid postfix operator")
		}

		c.push(pe.Op.Position, o.Plus)
		c.assignIdent(t)

	case *ast.FieldExpr:

		c.Visit(t.Operand)

		switch pe.Op.Text {
		case "++":
			c.push(pe.Op.Position, o.LoadOne)
		case "--":
			c.push(pe.Op.Position, o.LoadNegOne)
		default:
			panic("invalid postfix operator")
		}

		c.pushIndex(
			t.Key.Position,
			o.IncField,
			poolIndex(c.pool, g.NewStr(t.Key.Text)))

	case *ast.IndexExpr:

		c.Visit(t.Operand)
		c.Visit(t.Index)

		switch pe.Op.Text {
		case "++":
			c.push(pe.Op.Position, o.LoadOne)
		case "--":
			c.push(pe.Op.Position, o.LoadNegOne)
		default:
			panic("invalid postfix operator")
		}

		c.push(t.Index.Begin(), o.IncIndex)

	default:
		panic("invalid assignee type")
	}
}

func (c *compiler) visitIf(f *ast.IfStmt) {

	c.Visit(f.Cond)

	j0 := c.push(f.Cond.End(), o.JumpFalse, 0xFF, 0xFF)
	c.Visit(f.Then)

	if f.Else == nil {

		c.setJump(j0, c.opcLen())

	} else {

		j1 := c.push(f.Else.Begin(), o.Jump, 0xFF, 0xFF)
		c.setJump(j0, c.opcLen())

		c.Visit(f.Else)
		c.setJump(j1, c.opcLen())
	}
}

func (c *compiler) visitTernaryExpr(f *ast.TernaryExpr) {

	c.Visit(f.Cond)
	j0 := c.push(f.Cond.End(), o.JumpFalse, 0xFF, 0xFF)

	c.Visit(f.Then)
	j1 := c.push(f.Else.Begin(), o.Jump, 0xFF, 0xFF)
	c.setJump(j0, c.opcLen())

	c.Visit(f.Else)
	c.setJump(j1, c.opcLen())
}

func (c *compiler) visitWhile(w *ast.WhileStmt) {

	begin := c.opcLen()
	c.Visit(w.Cond)
	j0 := c.push(w.Cond.End(), o.JumpFalse, 0xFF, 0xFF)

	body := c.opcLen()
	c.Visit(w.Body)
	c.push(w.Body.End(), o.Jump, begin.high, begin.low)

	end := c.opcLen()
	c.setJump(j0, end)

	c.fixBreakContinue(begin, body, end)
}

func (c *compiler) visitFor(f *ast.ForStmt) {

	tok := f.Iterable.Begin()
	idx := f.IterableIdent.Variable.Index()

	// put Iterable expression on stack
	c.Visit(f.Iterable)

	// call NewIterator()
	c.push(tok, o.Iter)

	// store iterator
	c.pushIndex(tok, o.StoreLocal, idx)

	// top of loop: load iterator and call IterNext()
	begin := c.opcLen()
	c.pushIndex(tok, o.LoadLocal, idx)
	c.push(tok, o.IterNext)
	j0 := c.push(tok, o.JumpFalse, 0xFF, 0xFF)

	// load iterator and call IterGet()
	c.pushIndex(tok, o.LoadLocal, idx)
	c.push(tok, o.IterGet)

	if len(f.Idents) == 1 {
		// perform StoreLocal on the current item
		ident := f.Idents[0]
		c.pushIndex(ident.Begin(), o.StoreLocal, ident.Variable.Index())
	} else {
		// make sure the current item is really a tuple,
		// and is of the proper length
		c.pushIndex(tok, o.CheckTuple, len(f.Idents))

		// perform StoreLocal on each tuple element
		for i, ident := range f.Idents {
			c.push(tok, o.Dup)
			c.loadInt(tok, int64(i))
			c.push(tok, o.GetIndex)
			c.pushIndex(ident.Begin(), o.StoreLocal, ident.Variable.Index())
		}

		// pop the tuple
		c.push(tok, o.Pop)
	}

	// compile the body
	body := c.opcLen()
	c.Visit(f.Body)
	c.push(f.Body.End(), o.Jump, begin.high, begin.low)

	// jump to top of loop
	end := c.opcLen()
	c.setJump(j0, end)

	c.fixBreakContinue(begin, body, end)
}

func (c *compiler) fixBreakContinue(begin *instPtr, body *instPtr, end *instPtr) {

	// replace BreakStmt and ContinueStmt with Jump
	for i := body.ip; i < end.ip; {
		switch c.opc[i] {
		case o.BreakStmt:
			c.opc[i] = o.Jump
			c.opc[i+1] = end.high
			c.opc[i+2] = end.low
		case o.ContinueStmt:
			c.opc[i] = o.Jump
			c.opc[i+1] = begin.high
			c.opc[i+2] = begin.low
		}
		i += o.OpCodeSize(c.opc[i])
	}
}

func (c *compiler) visitBreak(br *ast.BreakStmt) {
	c.push(br.Begin(), o.BreakStmt, 0xFF, 0xFF)
}

func (c *compiler) visitContinue(cn *ast.ContinueStmt) {
	c.push(cn.Begin(), o.ContinueStmt, 0xFF, 0xFF)
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
		c.push(sw.End(), o.Pop)
	}

	// set all the end jumps
	for _, j := range endJumps {
		c.setJump(j, c.opcLen())
	}
}

func (c *compiler) visitCase(cs *ast.CaseNode, hasItem bool) int {

	bodyJumps := []int{}

	// visit each match, and jump to body if true
	for _, m := range cs.Matches {

		if hasItem {
			// if there is an item, Dup it and do an Eq comparison against the match
			c.push(m.Begin(), o.Dup)
			c.Visit(m)
			c.push(m.Begin(), o.Eq)
		} else {
			// otherwise, evaluate the match and assume its a Bool
			c.Visit(m)
		}

		bodyJumps = append(bodyJumps, c.push(m.End(), o.JumpTrue, 0xFF, 0xFF))
	}

	// no match -- jump to the end of the case
	caseEndJump := c.push(cs.End(), o.Jump, 0xFF, 0xFF)

	// set all the body jumps
	for _, j := range bodyJumps {
		c.setJump(j, c.opcLen())
	}

	// visit body, and then push a jump to the very end of the switch
	for _, n := range cs.Body {
		c.Visit(n)
	}
	endJump := c.push(cs.End(), o.Jump, 0xFF, 0xFF)

	// set the jump to the end of the case
	c.setJump(caseEndJump, c.opcLen())

	// return the jump to end of the switch
	return endJump
}

func (c *compiler) visitReturn(rt *ast.ReturnStmt) {
	c.Visit(rt.Val)
	c.push(rt.Begin(), o.ReturnStmt)
}

func (c *compiler) visitTry(t *ast.TryStmt) {

	begin := len(c.opc)
	c.Visit(t.TryBlock)
	end := len(c.opc)

	//////////////////////////
	// catch

	catch := -1
	if t.CatchBlock != nil {

		// push a jump, so we'll skip the catch block during normal execution
		catchEnd := c.push(t.TryBlock.End(), o.Jump, 0xFF, 0xFF)

		// save the beginning of the catch
		catch = len(c.opc)

		// store the exception that the interpreter has put on the stack for us
		v := t.CatchIdent.Variable
		assert(!v.IsCapture())
		c.pushIndex(t.CatchIdent.Begin(), o.StoreLocal, v.Index())

		// compile the catch
		c.Visit(t.CatchBlock)

		// pop the exception
		c.push(t.CatchBlock.End(), o.Pop)

		// add a Done to mark the end of the catch block
		c.push(t.CatchBlock.End(), o.Done)

		// fix the jump
		c.setJump(catchEnd, c.opcLen())
	}

	//////////////////////////
	// finally

	finally := -1
	if t.FinallyBlock != nil {

		// save the beginning of the finally
		finally = len(c.opc)

		// compile the finally
		c.Visit(t.FinallyBlock)

		// add a Done to mark the end of the finally block
		c.push(t.FinallyBlock.End(), o.Done)
	}

	//////////////////////////
	// done

	// sanity check
	assert(!(catch == -1 && finally == -1))
	c.handlers = append(c.handlers, g.ExceptionHandler{
		Begin:   begin,
		End:     end,
		Catch:   catch,
		Finally: finally,
	})
}

func (c *compiler) visitThrow(t *ast.ThrowStmt) {
	c.Visit(t.Val)
	c.push(t.End(), o.ThrowStmt)
}

func (c *compiler) visitBinaryExpr(b *ast.BinaryExpr) {

	switch b.Op.Kind {

	case ast.DblPipe:
		c.visitOr(b.LHS, b.RHS)
	case ast.DblAmp:
		c.visitAnd(b.LHS, b.RHS)

	case ast.DblEq:
		b.Traverse(c)
		c.push(b.Op.Position, o.Eq)
	case ast.NotEq:
		b.Traverse(c)
		c.push(b.Op.Position, o.Ne)

	case ast.Gt:
		b.Traverse(c)
		c.push(b.Op.Position, o.Gt)
	case ast.GtEq:
		b.Traverse(c)
		c.push(b.Op.Position, o.Gte)
	case ast.Lt:
		b.Traverse(c)
		c.push(b.Op.Position, o.Lt)
	case ast.LtEq:
		b.Traverse(c)
		c.push(b.Op.Position, o.Lte)
	case ast.Cmp:
		b.Traverse(c)
		c.push(b.Op.Position, o.Cmp)
	case ast.Has:
		b.Traverse(c)
		c.push(b.Op.Position, o.Has)

	case ast.Plus:
		b.Traverse(c)
		c.push(b.Op.Position, o.Plus)
	case ast.Minus:
		b.Traverse(c)
		c.push(b.Op.Position, o.Sub)
	case ast.Star:
		b.Traverse(c)
		c.push(b.Op.Position, o.Mul)
	case ast.Slash:
		b.Traverse(c)
		c.push(b.Op.Position, o.Div)

	case ast.Percent:
		b.Traverse(c)
		c.push(b.Op.Position, o.Rem)
	case ast.Amp:
		b.Traverse(c)
		c.push(b.Op.Position, o.BitAnd)
	case ast.Pipe:
		b.Traverse(c)
		c.push(b.Op.Position, o.BitOr)
	case ast.Caret:
		b.Traverse(c)
		c.push(b.Op.Position, o.BitXor)
	case ast.DblLt:
		b.Traverse(c)
		c.push(b.Op.Position, o.LeftShift)
	case ast.DblGt:
		b.Traverse(c)
		c.push(b.Op.Position, o.RightShift)

	default:
		panic("unreachable")
	}
}

func (c *compiler) visitOr(lhs ast.Expression, rhs ast.Expression) {

	c.Visit(lhs)
	j0 := c.push(lhs.End(), o.JumpTrue, 0xFF, 0xFF)

	c.Visit(rhs)
	j1 := c.push(rhs.End(), o.JumpFalse, 0xFF, 0xFF)

	c.setJump(j0, c.opcLen())
	c.push(rhs.End(), o.LoadTrue)
	j2 := c.push(rhs.End(), o.Jump, 0xFF, 0xFF)

	c.setJump(j1, c.opcLen())
	c.push(rhs.End(), o.LoadFalse)

	c.setJump(j2, c.opcLen())
}

func (c *compiler) visitAnd(lhs ast.Expression, rhs ast.Expression) {

	c.Visit(lhs)
	j0 := c.push(lhs.End(), o.JumpFalse, 0xFF, 0xFF)

	c.Visit(rhs)
	j1 := c.push(rhs.End(), o.JumpFalse, 0xFF, 0xFF)

	c.push(rhs.End(), o.LoadTrue)
	j2 := c.push(rhs.End(), o.Jump, 0xFF, 0xFF)

	c.setJump(j0, c.opcLen())
	c.setJump(j1, c.opcLen())
	c.push(rhs.End(), o.LoadFalse)

	c.setJump(j2, c.opcLen())
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
					c.push(u.Op.Position, o.LoadZero)
				case 1:
					c.push(u.Op.Position, o.LoadNegOne)
				default:
					c.pushIndex(
						u.Op.Position,
						o.LoadConst,
						poolIndex(c.pool, g.NewInt(-i)))
				}

			default:
				c.Visit(u.Operand)
				c.push(u.Op.Position, o.Negate)
			}
		default:
			c.Visit(u.Operand)
			c.push(u.Op.Position, o.Negate)
		}

	case ast.Not:
		c.Visit(u.Operand)
		c.push(u.Op.Position, o.Not)

	case ast.Tilde:
		c.Visit(u.Operand)
		c.push(u.Op.Position, o.Complement)

	default:
		panic("unreachable")
	}
}

func (c *compiler) visitBasicExpr(basic *ast.BasicExpr) {

	switch basic.Token.Kind {

	case ast.Null:
		c.push(basic.Token.Position, o.LoadNull)

	case ast.True:
		c.push(basic.Token.Position, o.LoadTrue)

	case ast.False:
		c.push(basic.Token.Position, o.LoadFalse)

	case ast.Str:
		c.pushIndex(
			basic.Token.Position,
			o.LoadConst,
			poolIndex(c.pool, g.NewStr(basic.Token.Text)))

	case ast.Int:
		c.loadInt(
			basic.Token.Position,
			parseInt(basic.Token.Text))

	case ast.Float:
		f := parseFloat(basic.Token.Text)
		c.pushIndex(
			basic.Token.Position,
			o.LoadConst,
			poolIndex(c.pool, g.NewFloat(f)))

	default:
		panic("unreachable")
	}

}

func (c *compiler) visitIdentExpr(ident *ast.IdentExpr) {
	v := ident.Variable
	if v.IsCapture() {
		c.pushIndex(ident.Begin(), o.LoadCapture, v.Index())
	} else {
		c.pushIndex(ident.Begin(), o.LoadLocal, v.Index())
	}
}

func (c *compiler) visitBuiltinExpr(blt *ast.BuiltinExpr) {

	c.pushIndex(blt.Fn.Position, o.LoadBuiltin, c.builtInMgr.IndexOf(blt.Fn.Text))
}

func (c *compiler) visitFunc(fe *ast.FnExpr) {

	c.pushIndex(fe.Begin(), o.NewFunc, len(c.funcs))

	parents := getSortedCaptureParents(fe.Scope)
	for _, vbl := range parents {
		if vbl.IsCapture() {
			c.pushIndex(fe.Begin(), o.FuncCapture, vbl.Index())
		} else {
			c.pushIndex(fe.Begin(), o.FuncLocal, vbl.Index())
		}
	}

	c.funcs = append(c.funcs, fe)
}

func (c *compiler) visitInvoke(inv *ast.InvokeExpr) {

	c.Visit(inv.Operand)
	for _, n := range inv.Params {
		c.Visit(n)
	}
	c.pushIndex(inv.Begin(), o.Invoke, len(inv.Params))
}

func (c *compiler) visitGo(gw *ast.GoStmt) {

	inv := gw.Invocation
	c.Visit(inv.Operand)
	for _, n := range inv.Params {
		c.Visit(n)
	}
	c.pushIndex(inv.Begin(), o.GoStmt, len(inv.Params))
}

func (c *compiler) visitExprStmt(es *ast.ExprStmt) {
	c.Visit(es.Expr)
}

func (c *compiler) visitStructExpr(stc *ast.StructExpr) {

	// create def and entries
	def := []*g.FieldDef{}
	for i, k := range stc.Keys {
		v := stc.Values[i]

		if p, ok := v.(*ast.PropNode); ok {
			def = append(def, &g.FieldDef{
				Name:       k.Text,
				IsReadonly: p.Setter == nil,
				IsProperty: true,
			})
		} else {
			def = append(def, &g.FieldDef{
				Name:       k.Text,
				IsReadonly: false,
				IsProperty: false,
			})
		}
	}
	defIdx := len(c.structDefs)
	c.structDefs = append(c.structDefs, def)

	// create new struct
	c.pushIndex(stc.Begin(), o.DefineStruct, defIdx)

	// if the struct is referenced by a 'this', then store local
	if this, ok := stc.Scope.GetVariable("this"); ok {
		c.push(stc.Begin(), o.Dup)
		c.pushIndex(stc.Begin(), o.StoreLocal, this.Index())
	}

	// init each value
	for i, k := range stc.Keys {
		v := stc.Values[i]
		c.push(k.Position, o.Dup)
		c.Visit(v)
		c.pushIndex(
			v.Begin(),
			o.InitField,
			poolIndex(c.pool, g.NewStr(k.Text)))
		c.push(k.Position, o.Pop)
	}
}

func (c *compiler) visitPropNode(pn *ast.PropNode) {

	c.Visit(pn.Getter)

	if pn.Setter == nil {
		c.push(pn.Begin(), o.LoadNull)
	} else {
		c.Visit(pn.Setter)
	}

	c.pushIndex(pn.Begin(), o.NewTuple, 2)
}

func (c *compiler) visitThisExpr(this *ast.ThisExpr) {
	v := this.Variable
	if v.IsCapture() {
		c.pushIndex(this.Begin(), o.LoadCapture, v.Index())
	} else {
		c.pushIndex(this.Begin(), o.LoadLocal, v.Index())
	}
}

func (c *compiler) visitFieldExpr(fe *ast.FieldExpr) {
	c.Visit(fe.Operand)
	c.pushIndex(
		fe.Key.Position,
		o.GetField,
		poolIndex(c.pool, g.NewStr(fe.Key.Text)))
}

func (c *compiler) visitIndexExpr(ie *ast.IndexExpr) {
	c.Visit(ie.Operand)
	c.Visit(ie.Index)
	c.push(ie.Index.Begin(), o.GetIndex)
}

func (c *compiler) visitSliceExpr(s *ast.SliceExpr) {
	c.Visit(s.Operand)
	c.Visit(s.From)
	c.Visit(s.To)
	c.push(s.From.Begin(), o.Slice)
}

func (c *compiler) visitSliceFromExpr(s *ast.SliceFromExpr) {
	c.Visit(s.Operand)
	c.Visit(s.From)
	c.push(s.From.Begin(), o.SliceFrom)
}

func (c *compiler) visitSliceToExpr(s *ast.SliceToExpr) {
	c.Visit(s.Operand)
	c.Visit(s.To)
	c.push(s.To.Begin(), o.SliceTo)
}

func (c *compiler) visitListExpr(ls *ast.ListExpr) {

	for _, v := range ls.Elems {
		c.Visit(v)
	}
	c.pushIndex(ls.Begin(), o.NewList, len(ls.Elems))
}

func (c *compiler) visitSetExpr(s *ast.SetExpr) {

	for _, v := range s.Elems {
		c.Visit(v)
	}
	c.pushIndex(s.Begin(), o.NewSet, len(s.Elems))
}

func (c *compiler) visitTupleExpr(tp *ast.TupleExpr) {

	for _, v := range tp.Elems {
		c.Visit(v)
	}
	c.pushIndex(tp.Begin(), o.NewTuple, len(tp.Elems))
}

func (c *compiler) visitDictExpr(d *ast.DictExpr) {

	for _, de := range d.Entries {
		c.Visit(de.Key)
		c.Visit(de.Value)
	}

	c.pushIndex(d.Begin(), o.NewDict, len(d.Entries))
}

func (c *compiler) loadInt(pos ast.Pos, i int64) {
	switch i {
	case 0:
		c.push(pos, o.LoadZero)
	case 1:
		c.push(pos, o.LoadOne)
	default:
		c.pushIndex(
			pos,
			o.LoadConst,
			poolIndex(c.pool, g.NewInt(i)))
	}
}

// returns the length of opc *before* the bytes are pushed
func (c *compiler) push(pos ast.Pos, bytes ...byte) int {
	n := len(c.opc)
	c.opc = append(c.opc, bytes...)

	ln := len(c.lnum)
	if (ln == 0) || (pos.Line != c.lnum[ln-1].LineNum) {
		c.lnum = append(c.lnum, g.LineNumberEntry{
			Index:   n,
			LineNum: pos.Line,
		})
	}

	return n
}

// push a 3-byte, indexed opcode
func (c *compiler) pushIndex(pos ast.Pos, opcode byte, idx int) int {
	high, low := index(idx)
	return c.push(pos, opcode, high, low)
}

// replace a mocked-up jump value with the 'real' destination
func (c *compiler) setJump(jmp int, dest *instPtr) {
	c.opc[jmp+1] = dest.high
	c.opc[jmp+2] = dest.low
}

func (c *compiler) opcLen() *instPtr {
	high, low := index(len(c.opc))
	return &instPtr{len(c.opc), high, low}
}

//--------------------------------------------------------------
// misc

type instPtr struct {
	ip   int
	high byte
	low  byte
}

func index(n int) (byte, byte) {
	assert(n < (2 << 16))
	return byte((n >> 8) & 0xFF), byte(n & 0xFF)
}

func parseInt(text string) int64 {
	i, err := strconv.ParseInt(text, 10, 64)
	assert(err == nil)
	assert(i >= 0)
	return i
}

func parseFloat(text string) float64 {
	f, err := strconv.ParseFloat(text, 64)
	assert(err == nil)
	assert(f >= 0)
	return f
}

func assert(flag bool) {
	if !flag {
		panic("assertion failure")
	}
}

//--------------------------------------------------------------
// pool

func poolIndex(pool *g.HashMap, key g.Basic) int {

	// Its OK for the Context to be nil here
	// The key is always Basic, so the Context will never be used.
	var cx g.Context

	b, err := pool.ContainsKey(cx, key)
	assert(err == nil)

	if b.BoolVal() {
		var v g.Value
		v, err = pool.Get(cx, key)
		assert(err == nil)

		i, ok := v.(g.Int)
		assert(ok)
		return int(i.IntVal())
	}
	i := pool.Len()
	err = pool.Put(cx, key, i)
	assert(err == nil)
	return int(i.IntVal())
}

// PoolItems is the contstant pool created by the compiler
type PoolItems []*g.HEntry

func (items PoolItems) Len() int {
	return len(items)
}

func (items PoolItems) Less(i, j int) bool {

	x, ok := items[i].Value.(g.Int)
	assert(ok)

	y, ok := items[j].Value.(g.Int)
	assert(ok)

	return x.IntVal() < y.IntVal()
}

func (items PoolItems) Swap(i, j int) {
	items[i], items[j] = items[j], items[i]
}

func makePoolSlice(pool *g.HashMap) []g.Basic {

	n := int(pool.Len().IntVal())

	entries := make([]*g.HEntry, 0, n)
	itr := pool.Iterator()
	for itr.Next() {
		entries = append(entries, itr.Get())
	}

	sort.Sort(PoolItems(entries))

	slice := make([]g.Basic, n)
	for i, e := range entries {
		b, ok := e.Key.(g.Basic)
		assert(ok)
		slice[i] = b
	}

	return slice
}

//--------------------------------------------------------------
// capture

func getSortedCaptureParents(f ast.FuncScope) []ast.Variable {

	// First sort the captures by child index
	sorted := byChildIndex{}
	for _, v := range f.GetCaptures() {
		sorted = append(sorted, v)
	}
	sort.Sort(sorted)

	// Then use the sorted list to create the proper ordering of parents
	parents := []ast.Variable{}
	for _, c := range sorted {
		parents = append(parents, c.Parent())
	}
	return parents
}

type byChildIndex []ast.Capture

// Variables are sorted by Index
func (c byChildIndex) Len() int {
	return len(c)
}
func (c byChildIndex) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}
func (c byChildIndex) Less(i, j int) bool {
	return c[i].Child().Index() < c[j].Child().Index()
}
