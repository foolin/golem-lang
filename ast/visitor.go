// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package ast

import (
	"bytes"
	"fmt"
	"reflect"
)

//--------------------------------------------------------------
// Visitor

type Visitor interface {
	Visit(node Node)
}

func (n *Empty) Traverse(v Visitor) {
}

func (imp *Import) Traverse(v Visitor) {
	v.Visit(imp.Ident)
}

func (cns *Const) Traverse(v Visitor) {
	for _, d := range cns.Decls {
		v.Visit(d.Ident)
		if d.Val != nil {
			v.Visit(d.Val)
		}
	}
}

func (let *Let) Traverse(v Visitor) {
	for _, d := range let.Decls {
		v.Visit(d.Ident)
		if d.Val != nil {
			v.Visit(d.Val)
		}
	}
}

func (nf *NamedFn) Traverse(v Visitor) {
	v.Visit(nf.Ident)
	v.Visit(nf.Func)
}

func (asn *AssignmentExpr) Traverse(v Visitor) {
	v.Visit(asn.Assignee)
	v.Visit(asn.Val)
}

func (ifn *If) Traverse(v Visitor) {
	v.Visit(ifn.Cond)
	v.Visit(ifn.Then)
	if ifn.Else != nil {
		v.Visit(ifn.Else)
	}
}

func (wh *While) Traverse(v Visitor) {
	v.Visit(wh.Cond)
	v.Visit(wh.Body)
}

func (fr *For) Traverse(v Visitor) {
	for _, n := range fr.Idents {
		v.Visit(n)
	}
	v.Visit(fr.IterableIdent)
	v.Visit(fr.Iterable)
	v.Visit(fr.Body)
}

func (sw *Switch) Traverse(v Visitor) {
	if sw.Item != nil {
		v.Visit(sw.Item)
	}

	for _, cs := range sw.Cases {
		v.Visit(cs)
	}

	if sw.Default != nil {
		v.Visit(sw.Default)
	}
}

func (cs *Case) Traverse(v Visitor) {
	for _, n := range cs.Matches {
		v.Visit(n)
	}

	for _, n := range cs.Body {
		v.Visit(n)
	}
}

func (def *Default) Traverse(v Visitor) {
	for _, n := range def.Body {
		v.Visit(n)
	}
}

func (br *Break) Traverse(v Visitor) {
}

func (cn *Continue) Traverse(v Visitor) {
}

func (rt *Return) Traverse(v Visitor) {
	if rt.Val != nil {
		v.Visit(rt.Val)
	}
}

func (t *Throw) Traverse(v Visitor) {
	v.Visit(t.Val)
}

func (t *Try) Traverse(v Visitor) {
	v.Visit(t.TryBlock)
	if t.CatchToken != nil {
		v.Visit(t.CatchIdent)
		v.Visit(t.CatchBlock)
	}
	if t.FinallyToken != nil {
		v.Visit(t.FinallyBlock)
	}
}

func (sp *Go) Traverse(v Visitor) {
	v.Visit(sp.Invocation)
}

func (blk *Block) Traverse(v Visitor) {
	for _, n := range blk.Nodes {
		v.Visit(n)
	}
}

func (trn *TernaryExpr) Traverse(v Visitor) {
	v.Visit(trn.Cond)
	v.Visit(trn.Then)
	v.Visit(trn.Else)
}

func (bin *BinaryExpr) Traverse(v Visitor) {
	v.Visit(bin.Lhs)
	v.Visit(bin.Rhs)
}

func (un *UnaryExpr) Traverse(v Visitor) {
	v.Visit(un.Operand)
}

func (pf *PostfixExpr) Traverse(v Visitor) {
	v.Visit(pf.Assignee)
}

func (basic *BasicExpr) Traverse(v Visitor) {
}

func (ident *IdentExpr) Traverse(v Visitor) {
}

func (ident *BuiltinExpr) Traverse(v Visitor) {
}

func (fn *FnExpr) Traverse(v Visitor) {
	for _, n := range fn.FormalParams {
		v.Visit(n.Ident)
	}
	v.Visit(fn.Body)
}

func (inv *InvokeExpr) Traverse(v Visitor) {
	v.Visit(inv.Operand)
	for _, n := range inv.Params {
		v.Visit(n)
	}
}

func (ls *ListExpr) Traverse(v Visitor) {
	for _, val := range ls.Elems {
		v.Visit(val)
	}
}

func (s *SetExpr) Traverse(v Visitor) {
	for _, val := range s.Elems {
		v.Visit(val)
	}
}

func (tp *TupleExpr) Traverse(v Visitor) {
	for _, val := range tp.Elems {
		v.Visit(val)
	}
}

func (stc *StructExpr) Traverse(v Visitor) {
	for _, val := range stc.Values {
		v.Visit(val)
	}
}

func (dict *DictExpr) Traverse(v Visitor) {
	for _, e := range dict.Entries {
		v.Visit(e)
	}
}

func (de *DictEntryExpr) Traverse(v Visitor) {
	v.Visit(de.Key)
	v.Visit(de.Value)
}

func (this *ThisExpr) Traverse(v Visitor) {
}

func (f *FieldExpr) Traverse(v Visitor) {
	v.Visit(f.Operand)
}

func (i *IndexExpr) Traverse(v Visitor) {
	v.Visit(i.Operand)
	v.Visit(i.Index)
}

func (i *SliceExpr) Traverse(v Visitor) {
	v.Visit(i.Operand)
	v.Visit(i.From)
	v.Visit(i.To)
}

func (i *SliceFromExpr) Traverse(v Visitor) {
	v.Visit(i.Operand)
	v.Visit(i.From)
}

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

	case *Empty:
		p.buf.WriteString("Empty\n")

	case *Block:
		p.buf.WriteString("Block\n")

	case *Import:
		p.buf.WriteString("Import\n")
	case *Const:
		p.buf.WriteString("Const\n")
	case *Let:
		p.buf.WriteString("Let\n")
	case *NamedFn:
		p.buf.WriteString("NamedFn\n")

	case *If:
		p.buf.WriteString("If\n")
	case *While:
		p.buf.WriteString("While\n")
	case *For:
		p.buf.WriteString("For\n")
	case *Break:
		p.buf.WriteString("Break\n")
	case *Continue:
		p.buf.WriteString("Continue\n")
	case *Return:
		p.buf.WriteString("Return\n")
	case *Throw:
		p.buf.WriteString("Throw\n")
	case *Try:
		p.buf.WriteString("Try\n")
	case *Go:
		p.buf.WriteString("Go\n")

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
