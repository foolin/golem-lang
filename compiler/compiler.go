// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package compiler

import (
	"fmt"
	"github.com/mjarmy/golem-lang/analyzer"
	"github.com/mjarmy/golem-lang/ast"
	g "github.com/mjarmy/golem-lang/core"
	"sort"
	"strconv"
)

type Compiler interface {
	ast.Visitor
	Compile() *g.BytecodeModule
}

type compiler struct {
	builtInMgr g.BuiltinManager

	pool     *g.HashMap
	opc      []byte
	lnum     []g.LineNumberEntry
	handlers []g.ExceptionHandler

	funcs      []*ast.FnExpr
	templates  []*g.Template
	structDefs [][]g.Field
	idx        int
}

func NewCompiler(anl analyzer.Analyzer, builtInMgr g.BuiltinManager) Compiler {

	funcs := []*ast.FnExpr{anl.Module()}
	templates := []*g.Template{}
	structDefs := [][]g.Field{}

	return &compiler{builtInMgr, g.EmptyHashMap(), nil, nil, nil, funcs, templates, structDefs, 0}
}

func (c *compiler) Compile() *g.BytecodeModule {

	// compile all the funcs
	for c.idx < len(c.funcs) {
		c.templates = append(
			c.templates,
			c.compileFunc(c.funcs[c.idx]))
		c.idx += 1
	}

	// done
	mod := &g.BytecodeModule{makePoolSlice(c.pool), nil, c.structDefs, c.templates, nil}
	mod.Contents = c.makeModuleContents(mod)
	return mod
}

func (c *compiler) makeModuleContents(mod *g.BytecodeModule) g.Struct {

	entries := []g.Field{}
	nodes := c.funcs[0].Body.Nodes
	for _, n := range nodes {
		switch t := n.(type) {
		case *ast.Let:
			for _, d := range t.Decls {
				vbl := d.Ident.Variable
				entries = append(entries, c.makeModuleProperty(
					mod, d.Ident.Symbol.Text, vbl.Index, vbl.IsConst))
			}
		case *ast.Const:
			for _, d := range t.Decls {
				vbl := d.Ident.Variable
				entries = append(entries, c.makeModuleProperty(
					mod, d.Ident.Symbol.Text, vbl.Index, vbl.IsConst))
			}
		case *ast.NamedFn:
			vbl := t.Ident.Variable
			entries = append(entries, c.makeModuleProperty(
				mod, t.Ident.Symbol.Text, vbl.Index, vbl.IsConst))
		}
	}

	stc, err := g.NewStruct(entries, false)
	assert(err == nil)
	return stc
}

func (c *compiler) makeModuleProperty(
	mod *g.BytecodeModule,
	name string,
	refIndex int,
	isConst bool) g.Field {

	getter := g.NewNativeFunc(0, 0,
		func(cx g.Context, values []g.Value) (g.Value, g.Error) {
			return mod.Refs[refIndex].Val, nil
		})

	var setter g.Func = nil
	if !isConst {
		setter = g.NewNativeFunc(1, 1,
			func(cx g.Context, values []g.Value) (g.Value, g.Error) {
				mod.Refs[refIndex].Val = values[0]
				return nil, nil
			})
	}

	return g.NewProperty(name, getter, setter)
}

func (c *compiler) compileFunc(fe *ast.FnExpr) *g.Template {

	arity := len(fe.FormalParams)
	tpl := &g.Template{arity, fe.NumCaptures, fe.NumLocals, nil, nil, nil}

	c.opc = []byte{}
	c.lnum = []g.LineNumberEntry{}
	c.handlers = []g.ExceptionHandler{}

	// TODO LOAD_NULL and RETURN are workarounds for the fact that
	// we have not yet written a Control Flow Graph
	c.push(ast.Pos{}, g.LOAD_NULL)
	c.Visit(fe.Body)
	c.push(ast.Pos{}, g.RETURN)

	tpl.OpCodes = c.opc
	tpl.LineNumberTable = c.lnum
	tpl.ExceptionHandlers = c.handlers

	return tpl
}

func (c *compiler) Visit(node ast.Node) {
	switch t := node.(type) {

	case *ast.Empty:
		c.visitEmpty(t)

	case *ast.Block:
		c.visitBlock(t)

	case *ast.Import:
		c.visitImport(t)

	case *ast.Const:
		c.visitDecls(t.Decls)

	case *ast.Let:
		c.visitDecls(t.Decls)

	case *ast.NamedFn:
		c.visitNamedFn(t)

	case *ast.Assignment:
		c.visitAssignment(t)

	case *ast.If:
		c.visitIf(t)

	case *ast.While:
		c.visitWhile(t)

	case *ast.For:
		c.visitFor(t)

	case *ast.Switch:
		c.visitSwitch(t)

	case *ast.Break:
		c.visitBreak(t)

	case *ast.Continue:
		c.visitContinue(t)

	case *ast.Return:
		c.visitReturn(t)

	case *ast.Try:
		c.visitTry(t)

	case *ast.Throw:
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

	case *ast.Go:
		c.visitGo(t)

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

func (c *compiler) visitEmpty(e *ast.Empty) {
}

func (c *compiler) visitBlock(blk *ast.Block) {

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

	for _, node := range blk.Nodes {
		c.Visit(node)

		// TODO
		//if (node is ast.Expr) && someControlFlowGraphCheck() {
		//	c.push(node.End(), g.POP)
		//}
	}
}

func (c *compiler) visitDecls(decls []*ast.Decl) {

	for _, d := range decls {
		if d.Val == nil {
			c.push(d.Ident.Begin(), g.LOAD_NULL)
		} else {
			c.Visit(d.Val)
		}

		c.assignIdent(d.Ident)
	}
}

func (c *compiler) visitImport(imp *ast.Import) {

	ident := imp.Ident

	// push the module onto the stack
	sym := ident.Symbol.Text
	c.pushIndex(
		ident.Begin(),
		g.IMPORT_MODULE,
		poolIndex(c.pool, g.MakeStr(sym)))

	// store module in identifer
	v := ident.Variable
	c.pushIndex(ident.Begin(), g.STORE_LOCAL, v.Index)
}

func (c *compiler) assignIdent(ident *ast.IdentExpr) {

	v := ident.Variable
	if v.IsCapture {
		c.pushIndex(ident.Begin(), g.STORE_CAPTURE, v.Index)
	} else {
		c.pushIndex(ident.Begin(), g.STORE_LOCAL, v.Index)
	}
}

func (c *compiler) visitNamedFn(nf *ast.NamedFn) {

	c.Visit(nf.Func)

	v := nf.Ident.Variable
	assert(!v.IsCapture)
	c.pushIndex(nf.Ident.Begin(), g.STORE_LOCAL, v.Index)
}

func (c *compiler) visitAssignment(asn *ast.Assignment) {

	switch t := asn.Assignee.(type) {

	case *ast.IdentExpr:

		c.Visit(asn.Val)
		c.push(asn.Eq.Position, g.DUP)
		c.assignIdent(t)

	case *ast.FieldExpr:

		c.Visit(t.Operand)
		c.Visit(asn.Val)
		c.pushIndex(
			t.Key.Position,
			g.SET_FIELD,
			poolIndex(c.pool, g.MakeStr(t.Key.Text)))

	case *ast.IndexExpr:

		c.Visit(t.Operand)
		c.Visit(t.Index)
		c.Visit(asn.Val)
		c.push(t.Index.Begin(), g.SET_INDEX)

	default:
		panic("invalid assignee type")
	}
}

func (c *compiler) visitPostfixExpr(pe *ast.PostfixExpr) {

	switch t := pe.Assignee.(type) {

	case *ast.IdentExpr:

		c.visitIdentExpr(t)
		c.push(t.Begin(), g.DUP)

		switch pe.Op.Text {
		case "++":
			c.push(pe.Op.Position, g.LOAD_ONE)
		case "--":
			c.push(pe.Op.Position, g.LOAD_NEG_ONE)
		default:
			panic("invalid postfix operator")
		}

		c.push(pe.Op.Position, g.PLUS)
		c.assignIdent(t)

	case *ast.FieldExpr:

		c.Visit(t.Operand)

		switch pe.Op.Text {
		case "++":
			c.push(pe.Op.Position, g.LOAD_ONE)
		case "--":
			c.push(pe.Op.Position, g.LOAD_NEG_ONE)
		default:
			panic("invalid postfix operator")
		}

		c.pushIndex(
			t.Key.Position,
			g.INC_FIELD,
			poolIndex(c.pool, g.MakeStr(t.Key.Text)))

	case *ast.IndexExpr:

		c.Visit(t.Operand)
		c.Visit(t.Index)

		switch pe.Op.Text {
		case "++":
			c.push(pe.Op.Position, g.LOAD_ONE)
		case "--":
			c.push(pe.Op.Position, g.LOAD_NEG_ONE)
		default:
			panic("invalid postfix operator")
		}

		c.push(t.Index.Begin(), g.INC_INDEX)

	default:
		panic("invalid assignee type")
	}
}

func (c *compiler) visitIf(f *ast.If) {

	c.Visit(f.Cond)

	j0 := c.push(f.Cond.End(), g.JUMP_FALSE, 0xFF, 0xFF)
	c.Visit(f.Then)

	if f.Else == nil {

		c.setJump(j0, c.opcLen())

	} else {

		j1 := c.push(f.Else.Begin(), g.JUMP, 0xFF, 0xFF)
		c.setJump(j0, c.opcLen())

		c.Visit(f.Else)
		c.setJump(j1, c.opcLen())
	}
}

func (c *compiler) visitTernaryExpr(f *ast.TernaryExpr) {

	c.Visit(f.Cond)
	j0 := c.push(f.Cond.End(), g.JUMP_FALSE, 0xFF, 0xFF)

	c.Visit(f.Then)
	j1 := c.push(f.Else.Begin(), g.JUMP, 0xFF, 0xFF)
	c.setJump(j0, c.opcLen())

	c.Visit(f.Else)
	c.setJump(j1, c.opcLen())
}

func (c *compiler) visitWhile(w *ast.While) {

	begin := c.opcLen()
	c.Visit(w.Cond)
	j0 := c.push(w.Cond.End(), g.JUMP_FALSE, 0xFF, 0xFF)

	body := c.opcLen()
	c.Visit(w.Body)
	c.push(w.Body.End(), g.JUMP, begin.high, begin.low)

	end := c.opcLen()
	c.setJump(j0, end)

	c.fixBreakContinue(begin, body, end)
}

func (c *compiler) visitFor(f *ast.For) {

	tok := f.Iterable.Begin()
	idx := f.IterableIdent.Variable.Index

	// put Iterable expression on stack
	c.Visit(f.Iterable)

	// call NewIterator()
	c.push(tok, g.ITER)

	// store iterator
	c.pushIndex(tok, g.STORE_LOCAL, idx)

	// top of loop: load iterator and call IterNext()
	begin := c.opcLen()
	c.pushIndex(tok, g.LOAD_LOCAL, idx)
	c.push(tok, g.ITER_NEXT)
	j0 := c.push(tok, g.JUMP_FALSE, 0xFF, 0xFF)

	// load iterator and call IterGet()
	c.pushIndex(tok, g.LOAD_LOCAL, idx)
	c.push(tok, g.ITER_GET)

	if len(f.Idents) == 1 {
		// perform STORE_LOCAL on the current item
		ident := f.Idents[0]
		c.pushIndex(ident.Begin(), g.STORE_LOCAL, ident.Variable.Index)
	} else {
		// make sure the current item is really a tuple,
		// and is of the proper length
		c.pushIndex(tok, g.CHECK_TUPLE, len(f.Idents))

		// perform STORE_LOCAL on each tuple element
		for i, ident := range f.Idents {
			c.push(tok, g.DUP)
			c.loadInt(tok, int64(i))
			c.push(tok, g.GET_INDEX)
			c.pushIndex(ident.Begin(), g.STORE_LOCAL, ident.Variable.Index)
		}

		// pop the tuple
		c.push(tok, g.POP)
	}

	// compile the body
	body := c.opcLen()
	c.Visit(f.Body)
	c.push(f.Body.End(), g.JUMP, begin.high, begin.low)

	// jump to top of loop
	end := c.opcLen()
	c.setJump(j0, end)

	c.fixBreakContinue(begin, body, end)
}

func (c *compiler) fixBreakContinue(begin *instPtr, body *instPtr, end *instPtr) {

	// replace BREAK and CONTINUE with JUMP
	for i := body.ip; i < end.ip; {
		switch c.opc[i] {
		case g.BREAK:
			c.opc[i] = g.JUMP
			c.opc[i+1] = end.high
			c.opc[i+2] = end.low
		case g.CONTINUE:
			c.opc[i] = g.JUMP
			c.opc[i+1] = begin.high
			c.opc[i+2] = begin.low
		}
		i += g.OpCodeSize(c.opc[i])
	}
}

func (c *compiler) visitBreak(br *ast.Break) {
	c.push(br.Begin(), g.BREAK, 0xFF, 0xFF)
}

func (c *compiler) visitContinue(cn *ast.Continue) {
	c.push(cn.Begin(), g.CONTINUE, 0xFF, 0xFF)
}

func (c *compiler) visitSwitch(sw *ast.Switch) {

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
	if sw.Default != nil {
		for _, n := range sw.Default.Body {
			c.Visit(n)
		}
	}

	// if there is an item, pop it
	if hasItem {
		c.push(sw.End(), g.POP)
	}

	// set all the end jumps
	for _, j := range endJumps {
		c.setJump(j, c.opcLen())
	}
}

func (c *compiler) visitCase(cs *ast.Case, hasItem bool) int {

	bodyJumps := []int{}

	// visit each match, and jump to body if true
	for _, m := range cs.Matches {

		if hasItem {
			// if there is an item, DUP it and do an EQ comparison against the match
			c.push(m.Begin(), g.DUP)
			c.Visit(m)
			c.push(m.Begin(), g.EQ)
		} else {
			// otherwise, evaluate the match and assume its a Bool
			c.Visit(m)
		}

		bodyJumps = append(bodyJumps, c.push(m.End(), g.JUMP_TRUE, 0xFF, 0xFF))
	}

	// no match -- jump to the end of the case
	caseEndJump := c.push(cs.End(), g.JUMP, 0xFF, 0xFF)

	// set all the body jumps
	for _, j := range bodyJumps {
		c.setJump(j, c.opcLen())
	}

	// visit body, and then push a jump to the very end of the switch
	for _, n := range cs.Body {
		c.Visit(n)
	}
	endJump := c.push(cs.End(), g.JUMP, 0xFF, 0xFF)

	// set the jump to the end of the case
	c.setJump(caseEndJump, c.opcLen())

	// return the jump to end of the switch
	return endJump
}

func (c *compiler) visitReturn(rt *ast.Return) {
	if rt.Val != nil {
		c.Visit(rt.Val)
	}
	c.push(rt.Begin(), g.RETURN)
}

func (c *compiler) visitTry(t *ast.Try) {

	begin := len(c.opc)
	c.Visit(t.TryBlock)
	end := len(c.opc)

	//////////////////////////
	// catch

	catch := -1
	if t.CatchBlock != nil {

		// push a jump, so we'll skip the catch block during normal execution
		end := c.push(t.TryBlock.End(), g.JUMP, 0xFF, 0xFF)

		// save the beginning of the catch
		catch = len(c.opc)

		// store the exception that the interpreter has put on the stack for us
		v := t.CatchIdent.Variable
		assert(!v.IsCapture)
		c.pushIndex(t.CatchIdent.Begin(), g.STORE_LOCAL, v.Index)

		// compile the catch
		c.Visit(t.CatchBlock)

		// pop the exception
		c.push(t.CatchBlock.End(), g.POP)

		// add a DONE to mark the end of the catch block
		c.push(t.CatchBlock.End(), g.DONE)

		// fix the jump
		c.setJump(end, c.opcLen())
	}

	//////////////////////////
	// finally

	finally := -1
	if t.FinallyBlock != nil {

		// save the beginning of the finally
		finally = len(c.opc)

		// compile the finally
		c.Visit(t.FinallyBlock)

		// add a DONE to mark the end of the finally block
		c.push(t.FinallyBlock.End(), g.DONE)
	}

	//////////////////////////
	// done

	// sanity check
	assert(!(catch == -1 && finally == -1))
	c.handlers = append(c.handlers, g.ExceptionHandler{begin, end, catch, finally})
}

func (c *compiler) visitThrow(t *ast.Throw) {
	c.Visit(t.Val)
	c.push(t.End(), g.THROW)
}

func (c *compiler) visitBinaryExpr(b *ast.BinaryExpr) {

	switch b.Op.Kind {

	case ast.DBL_PIPE:
		c.visitOr(b.Lhs, b.Rhs)
	case ast.DBL_AMP:
		c.visitAnd(b.Lhs, b.Rhs)

	case ast.DBL_EQ:
		b.Traverse(c)
		c.push(b.Op.Position, g.EQ)
	case ast.NOT_EQ:
		b.Traverse(c)
		c.push(b.Op.Position, g.NE)

	case ast.GT:
		b.Traverse(c)
		c.push(b.Op.Position, g.GT)
	case ast.GT_EQ:
		b.Traverse(c)
		c.push(b.Op.Position, g.GTE)
	case ast.LT:
		b.Traverse(c)
		c.push(b.Op.Position, g.LT)
	case ast.LT_EQ:
		b.Traverse(c)
		c.push(b.Op.Position, g.LTE)
	case ast.CMP:
		b.Traverse(c)
		c.push(b.Op.Position, g.CMP)
	case ast.HAS:
		b.Traverse(c)
		c.push(b.Op.Position, g.HAS)

	case ast.PLUS:
		b.Traverse(c)
		c.push(b.Op.Position, g.PLUS)
	case ast.MINUS:
		b.Traverse(c)
		c.push(b.Op.Position, g.SUB)
	case ast.STAR:
		b.Traverse(c)
		c.push(b.Op.Position, g.MUL)
	case ast.SLASH:
		b.Traverse(c)
		c.push(b.Op.Position, g.DIV)

	case ast.PERCENT:
		b.Traverse(c)
		c.push(b.Op.Position, g.REM)
	case ast.AMP:
		b.Traverse(c)
		c.push(b.Op.Position, g.BIT_AND)
	case ast.PIPE:
		b.Traverse(c)
		c.push(b.Op.Position, g.BIT_OR)
	case ast.CARET:
		b.Traverse(c)
		c.push(b.Op.Position, g.BIT_XOR)
	case ast.DBL_LT:
		b.Traverse(c)
		c.push(b.Op.Position, g.LEFT_SHIFT)
	case ast.DBL_GT:
		b.Traverse(c)
		c.push(b.Op.Position, g.RIGHT_SHIFT)

	default:
		panic("unreachable")
	}
}

func (c *compiler) visitOr(lhs ast.Expr, rhs ast.Expr) {

	c.Visit(lhs)
	j0 := c.push(lhs.End(), g.JUMP_TRUE, 0xFF, 0xFF)

	c.Visit(rhs)
	j1 := c.push(rhs.End(), g.JUMP_FALSE, 0xFF, 0xFF)

	c.setJump(j0, c.opcLen())
	c.push(rhs.End(), g.LOAD_TRUE)
	j2 := c.push(rhs.End(), g.JUMP, 0xFF, 0xFF)

	c.setJump(j1, c.opcLen())
	c.push(rhs.End(), g.LOAD_FALSE)

	c.setJump(j2, c.opcLen())
}

func (c *compiler) visitAnd(lhs ast.Expr, rhs ast.Expr) {

	c.Visit(lhs)
	j0 := c.push(lhs.End(), g.JUMP_FALSE, 0xFF, 0xFF)

	c.Visit(rhs)
	j1 := c.push(rhs.End(), g.JUMP_FALSE, 0xFF, 0xFF)

	c.push(rhs.End(), g.LOAD_TRUE)
	j2 := c.push(rhs.End(), g.JUMP, 0xFF, 0xFF)

	c.setJump(j0, c.opcLen())
	c.setJump(j1, c.opcLen())
	c.push(rhs.End(), g.LOAD_FALSE)

	c.setJump(j2, c.opcLen())
}

func (c *compiler) visitUnaryExpr(u *ast.UnaryExpr) {

	switch u.Op.Kind {
	case ast.MINUS:
		opn := u.Operand

		switch t := opn.(type) {
		case *ast.BasicExpr:
			switch t.Token.Kind {

			case ast.INT:
				i := parseInt(t.Token.Text)
				switch i {
				case 0:
					c.push(u.Op.Position, g.LOAD_ZERO)
				case 1:
					c.push(u.Op.Position, g.LOAD_NEG_ONE)
				default:
					c.pushIndex(
						u.Op.Position,
						g.LOAD_CONST,
						poolIndex(c.pool, g.MakeInt(-i)))
				}

			default:
				c.Visit(u.Operand)
				c.push(u.Op.Position, g.NEGATE)
			}
		default:
			c.Visit(u.Operand)
			c.push(u.Op.Position, g.NEGATE)
		}

	case ast.NOT:
		c.Visit(u.Operand)
		c.push(u.Op.Position, g.NOT)

	case ast.TILDE:
		c.Visit(u.Operand)
		c.push(u.Op.Position, g.COMPLEMENT)

	default:
		panic("unreachable")
	}
}

func (c *compiler) visitBasicExpr(basic *ast.BasicExpr) {

	switch basic.Token.Kind {

	case ast.NULL:
		c.push(basic.Token.Position, g.LOAD_NULL)

	case ast.TRUE:
		c.push(basic.Token.Position, g.LOAD_TRUE)

	case ast.FALSE:
		c.push(basic.Token.Position, g.LOAD_FALSE)

	case ast.STR:
		c.pushIndex(
			basic.Token.Position,
			g.LOAD_CONST,
			poolIndex(c.pool, g.MakeStr(basic.Token.Text)))

	case ast.INT:
		c.loadInt(
			basic.Token.Position,
			parseInt(basic.Token.Text))

	case ast.FLOAT:
		f := parseFloat(basic.Token.Text)
		c.pushIndex(
			basic.Token.Position,
			g.LOAD_CONST,
			poolIndex(c.pool, g.MakeFloat(f)))

	default:
		panic("unreachable")
	}

}

func (c *compiler) visitIdentExpr(ident *ast.IdentExpr) {
	v := ident.Variable
	if v.IsCapture {
		c.pushIndex(ident.Begin(), g.LOAD_CAPTURE, v.Index)
	} else {
		c.pushIndex(ident.Begin(), g.LOAD_LOCAL, v.Index)
	}
}

func (c *compiler) visitBuiltinExpr(blt *ast.BuiltinExpr) {

	c.pushIndex(blt.Fn.Position, g.LOAD_BUILTIN, c.builtInMgr.IndexOf(blt.Fn.Text))
}

func (c *compiler) visitFunc(fe *ast.FnExpr) {

	c.pushIndex(fe.Begin(), g.NEW_FUNC, len(c.funcs))
	for _, pc := range fe.ParentCaptures {
		if pc.IsCapture {
			c.pushIndex(fe.Begin(), g.FUNC_CAPTURE, pc.Index)
		} else {
			c.pushIndex(fe.Begin(), g.FUNC_LOCAL, pc.Index)
		}
	}

	c.funcs = append(c.funcs, fe)
}

func (c *compiler) visitInvoke(inv *ast.InvokeExpr) {

	c.Visit(inv.Operand)
	for _, n := range inv.Params {
		c.Visit(n)
	}
	c.pushIndex(inv.Begin(), g.INVOKE, len(inv.Params))
}

func (c *compiler) visitGo(spawn *ast.Go) {

	inv := spawn.Invocation
	c.Visit(inv.Operand)
	for _, n := range inv.Params {
		c.Visit(n)
	}
	c.pushIndex(inv.Begin(), g.GO, len(inv.Params))
}

func (c *compiler) visitStructExpr(stc *ast.StructExpr) {

	// create def and entries
	def := []g.Field{}
	for _, k := range stc.Keys {
		def = append(def, g.NewField(k.Text, false, g.NULL))
	}
	defIdx := len(c.structDefs)
	c.structDefs = append(c.structDefs, def)

	// create new struct
	c.pushIndex(stc.Begin(), g.NEW_STRUCT, defIdx)

	// if the struct is referenced by a 'this', then store local
	if stc.LocalThisIndex != -1 {
		c.push(stc.Begin(), g.DUP)
		c.pushIndex(stc.Begin(), g.STORE_LOCAL, stc.LocalThisIndex)
	}

	// init each value
	for i, k := range stc.Keys {
		v := stc.Values[i]
		c.push(k.Position, g.DUP)
		c.Visit(v)
		c.pushIndex(
			v.Begin(),
			g.INIT_FIELD,
			poolIndex(c.pool, g.MakeStr(k.Text)))
		c.push(k.Position, g.POP)
	}
}

func (c *compiler) visitThisExpr(this *ast.ThisExpr) {
	v := this.Variable
	if v.IsCapture {
		c.pushIndex(this.Begin(), g.LOAD_CAPTURE, v.Index)
	} else {
		c.pushIndex(this.Begin(), g.LOAD_LOCAL, v.Index)
	}
}

func (c *compiler) visitFieldExpr(fe *ast.FieldExpr) {
	c.Visit(fe.Operand)
	c.pushIndex(
		fe.Key.Position,
		g.GET_FIELD,
		poolIndex(c.pool, g.MakeStr(fe.Key.Text)))
}

func (c *compiler) visitIndexExpr(ie *ast.IndexExpr) {
	c.Visit(ie.Operand)
	c.Visit(ie.Index)
	c.push(ie.Index.Begin(), g.GET_INDEX)
}

func (c *compiler) visitSliceExpr(s *ast.SliceExpr) {
	c.Visit(s.Operand)
	c.Visit(s.From)
	c.Visit(s.To)
	c.push(s.From.Begin(), g.SLICE)
}

func (c *compiler) visitSliceFromExpr(s *ast.SliceFromExpr) {
	c.Visit(s.Operand)
	c.Visit(s.From)
	c.push(s.From.Begin(), g.SLICE_FROM)
}

func (c *compiler) visitSliceToExpr(s *ast.SliceToExpr) {
	c.Visit(s.Operand)
	c.Visit(s.To)
	c.push(s.To.Begin(), g.SLICE_TO)
}

func (c *compiler) visitListExpr(ls *ast.ListExpr) {

	for _, v := range ls.Elems {
		c.Visit(v)
	}
	c.pushIndex(ls.Begin(), g.NEW_LIST, len(ls.Elems))
}

func (c *compiler) visitSetExpr(s *ast.SetExpr) {

	for _, v := range s.Elems {
		c.Visit(v)
	}
	c.pushIndex(s.Begin(), g.NEW_SET, len(s.Elems))
}

func (c *compiler) visitTupleExpr(tp *ast.TupleExpr) {

	for _, v := range tp.Elems {
		c.Visit(v)
	}
	c.pushIndex(tp.Begin(), g.NEW_TUPLE, len(tp.Elems))
}

func (c *compiler) visitDictExpr(d *ast.DictExpr) {

	for _, de := range d.Entries {
		c.Visit(de.Key)
		c.Visit(de.Value)
	}

	c.pushIndex(d.Begin(), g.NEW_DICT, len(d.Entries))
}

func (c *compiler) loadInt(pos ast.Pos, i int64) {
	switch i {
	case 0:
		c.push(pos, g.LOAD_ZERO)
	case 1:
		c.push(pos, g.LOAD_ONE)
	default:
		c.pushIndex(
			pos,
			g.LOAD_CONST,
			poolIndex(c.pool, g.MakeInt(i)))
	}
}

// returns the length of opc *before* the bytes are pushed
func (c *compiler) push(pos ast.Pos, bytes ...byte) int {
	n := len(c.opc)
	for _, b := range bytes {
		c.opc = append(c.opc, b)
	}

	ln := len(c.lnum)
	if (ln == 0) || (pos.Line != c.lnum[ln-1].LineNum) {
		c.lnum = append(c.lnum, g.LineNumberEntry{n, pos.Line})
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
	return int64(i)
}

func parseFloat(text string) float64 {
	f, err := strconv.ParseFloat(text, 64)
	assert(err == nil)
	assert(f >= 0)
	return float64(f)
}

//--------------------------------------------------------------
// pool

func poolIndex(pool *g.HashMap, key g.Basic) int {

	// Its OK for the Context to be nil here
	// The key is always Basic, so the Context will never be used.
	var cx g.Context = nil

	b, err := pool.ContainsKey(cx, key)
	assert(err == nil)

	if b.BoolVal() {
		v, err := pool.Get(cx, key)
		assert(err == nil)

		i, ok := v.(g.Int)
		assert(ok)
		return int(i.IntVal())
	} else {
		i := pool.Len()
		err := pool.Put(cx, key, i)
		assert(err == nil)
		return int(i.IntVal())
	}
}

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

	slice := make([]g.Basic, n, n)
	for i, e := range entries {
		b, ok := e.Key.(g.Basic)
		assert(ok)
		slice[i] = b
	}

	return slice
}

func assert(flag bool) {
	if !flag {
		panic("assertion failure")
	}
}
