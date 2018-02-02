// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package ast

import (
	"bytes"
	"fmt"
	"reflect"
)

// Visitor visits Nodes in an AST
type Visitor interface {
	Visit(node Node)
}

// Traverse ImportStmt
func (imp *ImportStmt) Traverse(v Visitor) {
	v.Visit(imp.Ident)
}

// Traverse ConstStmt
func (cns *ConstStmt) Traverse(v Visitor) {
	for _, d := range cns.Decls {
		v.Visit(d.Ident)
		if d.Val != nil {
			v.Visit(d.Val)
		}
	}
}

// Traverse LetStmt
func (let *LetStmt) Traverse(v Visitor) {
	for _, d := range let.Decls {
		v.Visit(d.Ident)
		if d.Val != nil {
			v.Visit(d.Val)
		}
	}
}

// Traverse NamedFnStmt
func (nf *NamedFnStmt) Traverse(v Visitor) {
	v.Visit(nf.Ident)
	v.Visit(nf.Func)
}

// Traverse AssignmentExpr
func (asn *AssignmentExpr) Traverse(v Visitor) {
	v.Visit(asn.Assignee)
	v.Visit(asn.Val)
}

// Traverse IfStmt
func (ifn *IfStmt) Traverse(v Visitor) {
	v.Visit(ifn.Cond)
	v.Visit(ifn.Then)
	if ifn.Else != nil {
		v.Visit(ifn.Else)
	}
}

// Traverse WhileStmt
func (wh *WhileStmt) Traverse(v Visitor) {
	v.Visit(wh.Cond)
	v.Visit(wh.Body)
}

// Traverse ForStmt
func (fr *ForStmt) Traverse(v Visitor) {
	for _, n := range fr.Idents {
		v.Visit(n)
	}
	v.Visit(fr.IterableIdent)
	v.Visit(fr.Iterable)
	v.Visit(fr.Body)
}

// Traverse SwitchStmt
func (sw *SwitchStmt) Traverse(v Visitor) {
	if sw.Item != nil {
		v.Visit(sw.Item)
	}

	for _, cs := range sw.Cases {
		v.Visit(cs)
	}

	if sw.DefaultNode != nil {
		v.Visit(sw.DefaultNode)
	}
}

// Traverse CaseNode
func (cs *CaseNode) Traverse(v Visitor) {
	for _, n := range cs.Matches {
		v.Visit(n)
	}

	for _, n := range cs.Body {
		v.Visit(n)
	}
}

// Traverse DefaultNode
func (def *DefaultNode) Traverse(v Visitor) {
	for _, n := range def.Body {
		v.Visit(n)
	}
}

// Traverse BreakStmt
func (br *BreakStmt) Traverse(v Visitor) {
}

// Traverse ContinueStmt
func (cn *ContinueStmt) Traverse(v Visitor) {
}

// Traverse ReturnStmt
func (rt *ReturnStmt) Traverse(v Visitor) {
	if rt.Val != nil {
		v.Visit(rt.Val)
	}
}

// Traverse ThrowStmt
func (t *ThrowStmt) Traverse(v Visitor) {
	v.Visit(t.Val)
}

// Traverse TryStmt
func (t *TryStmt) Traverse(v Visitor) {
	v.Visit(t.TryBlock)
	if t.CatchToken != nil {
		v.Visit(t.CatchIdent)
		v.Visit(t.CatchBlock)
	}
	if t.FinallyToken != nil {
		v.Visit(t.FinallyBlock)
	}
}

// Traverse GoStmt
func (g *GoStmt) Traverse(v Visitor) {
	v.Visit(g.Invocation)
}

// Traverse ExprStmt
func (n *ExprStmt) Traverse(v Visitor) {
	v.Visit(n.Expr)
}

// Traverse BlockNode
func (blk *BlockNode) Traverse(v Visitor) {
	for _, n := range blk.Statements {
		v.Visit(n)
	}
}

// Traverse TernaryExpr
func (trn *TernaryExpr) Traverse(v Visitor) {
	v.Visit(trn.Cond)
	v.Visit(trn.Then)
	v.Visit(trn.Else)
}

// Traverse BinaryExpr
func (bin *BinaryExpr) Traverse(v Visitor) {
	v.Visit(bin.Lhs)
	v.Visit(bin.Rhs)
}

// Traverse UnaryExpr
func (un *UnaryExpr) Traverse(v Visitor) {
	v.Visit(un.Operand)
}

// Traverse PostfixExpr
func (pf *PostfixExpr) Traverse(v Visitor) {
	v.Visit(pf.Assignee)
}

// Traverse BasicExpr
func (basic *BasicExpr) Traverse(v Visitor) {
}

// Traverse IdentExpr
func (ident *IdentExpr) Traverse(v Visitor) {
}

// Traverse BuiltinExpr
func (ident *BuiltinExpr) Traverse(v Visitor) {
}

// Traverse FnExpr
func (fn *FnExpr) Traverse(v Visitor) {
	for _, n := range fn.FormalParams {
		v.Visit(n.Ident)
	}
	v.Visit(fn.Body)
}

// Traverse InvokeExpr
func (inv *InvokeExpr) Traverse(v Visitor) {
	v.Visit(inv.Operand)
	for _, n := range inv.Params {
		v.Visit(n)
	}
}

// Traverse ListExpr
func (ls *ListExpr) Traverse(v Visitor) {
	for _, val := range ls.Elems {
		v.Visit(val)
	}
}

// Traverse SetExpr
func (s *SetExpr) Traverse(v Visitor) {
	for _, val := range s.Elems {
		v.Visit(val)
	}
}

// Traverse TupleExpr
func (tp *TupleExpr) Traverse(v Visitor) {
	for _, val := range tp.Elems {
		v.Visit(val)
	}
}

// Traverse StructExpr
func (stc *StructExpr) Traverse(v Visitor) {
	for _, val := range stc.Values {
		v.Visit(val)
	}
}

// Traverse DictExpr
func (dict *DictExpr) Traverse(v Visitor) {
	for _, e := range dict.Entries {
		v.Visit(e)
	}
}

// Traverse DictEntryExpr
func (de *DictEntryExpr) Traverse(v Visitor) {
	v.Visit(de.Key)
	v.Visit(de.Value)
}

// Traverse ThisExpr
func (t *ThisExpr) Traverse(v Visitor) {
}

// Traverse FieldExpr
func (f *FieldExpr) Traverse(v Visitor) {
	v.Visit(f.Operand)
}

// Traverse IndexExpr
func (i *IndexExpr) Traverse(v Visitor) {
	v.Visit(i.Operand)
	v.Visit(i.Index)
}

// Traverse SliceExpr
func (i *SliceExpr) Traverse(v Visitor) {
	v.Visit(i.Operand)
	v.Visit(i.From)
	v.Visit(i.To)
}

// Traverse SliceFromExpr
func (i *SliceFromExpr) Traverse(v Visitor) {
	v.Visit(i.Operand)
	v.Visit(i.From)
}

// Traverse SliceToExpr
func (i *SliceToExpr) Traverse(v Visitor) {
	v.Visit(i.Operand)
	v.Visit(i.To)
}

//--------------------------------------------------------------
// ast debug

type dump struct {
	buf    bytes.Buffer
	indent int
}

// Dump creates a string representation of a Node and its
// descendant Nodes.
func Dump(node Node) string {
	p := &dump{}
	p.Visit(node)
	return p.buf.String()
}

func (p *dump) Visit(node Node) {

	for i := 0; i < p.indent; i++ {
		p.buf.WriteString(".   ")
	}

	switch t := node.(type) {

	case *BlockNode:
		p.buf.WriteString("BlockNode\n")

	case *ImportStmt:
		p.buf.WriteString("ImportStmt\n")
	case *ConstStmt:
		p.buf.WriteString("ConstStmt\n")
	case *LetStmt:
		p.buf.WriteString("LetStmt\n")
	case *NamedFnStmt:
		p.buf.WriteString("NamedFnStmt\n")

	case *IfStmt:
		p.buf.WriteString("IfStmt\n")
	case *WhileStmt:
		p.buf.WriteString("WhileStmt\n")
	case *ForStmt:
		p.buf.WriteString("ForStmt\n")
	case *BreakStmt:
		p.buf.WriteString("BreakStmt\n")
	case *ContinueStmt:
		p.buf.WriteString("ContinueStmt\n")
	case *ReturnStmt:
		p.buf.WriteString("ReturnStmt\n")
	case *ThrowStmt:
		p.buf.WriteString("ThrowStmt\n")
	case *TryStmt:
		p.buf.WriteString("TryStmt\n")
	case *GoStmt:
		p.buf.WriteString("GoStmt\n")

	case *ExprStmt:
		p.buf.WriteString("ExprStmt\n")

	case *AssignmentExpr:
		p.buf.WriteString("AssignmentExpr\n")
	case *BinaryExpr:
		p.buf.WriteString(fmt.Sprintf("BinaryExpr(%q)\n", t.Op.Text))
	case *UnaryExpr:
		p.buf.WriteString(fmt.Sprintf("UnaryExpr(%q)\n", t.Op.Text))
	case *PostfixExpr:
		p.buf.WriteString(fmt.Sprintf("PostfixExpr(%q)\n", t.Op.Text))
	case *BasicExpr:
		p.buf.WriteString(fmt.Sprintf("BasicExpr(%v,%q)\n", t.Token.Kind, t.Token.Text))
	case *IdentExpr:
		p.buf.WriteString(fmt.Sprintf("IdentExpr(%v,%v)\n", t.Symbol.Text, t.Variable))

	case *FnExpr:
		p.buf.WriteString(fmt.Sprintf("FnExpr(numLocals:%d", t.NumLocals))
		p.buf.WriteString(fmt.Sprintf(" numCaptures:%d", t.NumCaptures))
		p.buf.WriteString(" parentCaptures:")
		p.buf.WriteString(varsString(t.ParentCaptures))
		p.buf.WriteString(")\n")
	case *InvokeExpr:
		p.buf.WriteString("InvokeExpr\n")
	case *BuiltinExpr:
		p.buf.WriteString(fmt.Sprintf("BuiltinExpr(%q)\n", t.Fn.Text))

	case *StructExpr:
		p.buf.WriteString(fmt.Sprintf("StructExpr(%v,%d)\n", tokensString(t.Keys), t.LocalThisIndex))
	case *DictExpr:
		p.buf.WriteString("DictExpr\n")
	case *DictEntryExpr:
		p.buf.WriteString("DictEntryExpr\n")
	case *ThisExpr:
		p.buf.WriteString(fmt.Sprintf("ThisExpr(%v)\n", t.Variable))
	case *ListExpr:
		p.buf.WriteString("ListExpr\n")
	case *TupleExpr:
		p.buf.WriteString("TupleExpr\n")

	case *FieldExpr:
		p.buf.WriteString(fmt.Sprintf("FieldExpr(%v)\n", t.Key.Text))

	case *IndexExpr:
		p.buf.WriteString("IndexExpr\n")

	case *SliceExpr:
		p.buf.WriteString("SliceExpr\n")
	case *SliceFromExpr:
		p.buf.WriteString("SliceFromExpr\n")
	case *SliceToExpr:
		p.buf.WriteString("SliceToExpr\n")

	default:
		fmt.Println(reflect.TypeOf(node))
		panic("cannot visit")
	}

	p.indent++
	node.Traverse(p)
	p.indent--
}

func varsString(vars []*Variable) string {

	var buf bytes.Buffer
	buf.WriteString("[")
	n := 0
	for v := range vars {
		if n > 0 {
			buf.WriteString(", ")
		}
		n++
		buf.WriteString(fmt.Sprintf("%v", vars[v]))
	}
	buf.WriteString("]")
	return buf.String()
}

func tokensString(tokens []*Token) string {

	var buf bytes.Buffer
	buf.WriteString("[")
	n := 0
	for t := range tokens {
		if n > 0 {
			buf.WriteString(", ")
		}
		n++
		buf.WriteString(fmt.Sprintf("%v", tokens[t].Text))
	}
	buf.WriteString("]")
	return buf.String()
}
