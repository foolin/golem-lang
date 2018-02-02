// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package ast

import (
	"bytes"
	"fmt"
	"strings"
)

//--------------------------------------------------------------
// Node

// interfaces
type (
	Node interface {
		fmt.Stringer
		Traverse(Visitor)
		Begin() Pos
		End() Pos
	}

	Statement interface {
		Node
		stmtMarker()
	}

	Loop interface {
		Statement
		loopMarker()
	}

	Expression interface {
		Node
		exprMarker()
	}

	Assignable interface {
		Expression
		assignableMarker()
	}
)

// structs
type (

	//---------------------
	// statement ndoes

	ImportStmt struct {
		Token *Token
		Ident *IdentExpr
	}

	ConstStmt struct {
		Token *Token
		Decls []*DeclNode
	}

	LetStmt struct {
		Token *Token
		Decls []*DeclNode
	}

	NamedFnStmt struct {
		Token *Token
		Ident *IdentExpr
		Func  *FnExpr
	}

	IfStmt struct {
		Token *Token
		Cond  Expression
		Then  *BlockNode
		Else  Node // either a BlockNode, or another IfStmt
	}

	WhileStmt struct {
		Token *Token
		Cond  Expression
		Body  *BlockNode
	}

	ForStmt struct {
		Token         *Token
		Idents        []*IdentExpr
		IterableIdent *IdentExpr
		Iterable      Expression
		Body          *BlockNode
	}

	SwitchStmt struct {
		Token   *Token
		Item    Expression
		LBrace  *Token
		Cases   []*CaseNode
		DefaultNode *DefaultNode
		RBrace  *Token
	}

	BreakStmt struct {
		Token *Token
	}

	ContinueStmt struct {
		Token *Token
	}

	ReturnStmt struct {
		Token *Token
		Val   Expression
	}

	ThrowStmt struct {
		Token *Token
		Val   Expression
	}

	TryStmt struct {
		TryToken     *Token
		TryBlock     *BlockNode
		CatchToken   *Token
		CatchIdent   *IdentExpr
		CatchBlock   *BlockNode
		FinallyToken *Token
		FinallyBlock *BlockNode
	}

	GoStmt struct {
		Token      *Token
		Invocation *InvokeExpr
	}

	// ExprStmt is a Statement that contains an Expression
	ExprStmt struct {
		Expr Expression
	}

	//--------------------------------------
	// nodes that are parts of a statement

	BlockNode struct {
		LBrace     *Token
		Statements []Statement
		RBrace     *Token
	}

	DeclNode struct {
		Ident *IdentExpr
		Val   Expression
	}

	CaseNode struct {
		Token   *Token
		Matches []Expression
		Body    []Statement
	}

	DefaultNode struct {
		Token *Token
		Body  []Statement
	}

	//---------------------
	// expression nodes

	AssignmentExpr struct {
		Assignee Assignable
		Eq       *Token
		Val      Expression
	}

	TernaryExpr struct {
		Cond Expression
		Then Expression
		Else Expression
	}

	BinaryExpr struct {
		Lhs Expression
		Op  *Token
		Rhs Expression
	}

	UnaryExpr struct {
		Op      *Token
		Operand Expression
	}

	PostfixExpr struct {
		Assignee Assignable
		Op       *Token
	}

	BasicExpr struct {
		Token *Token
	}

	IdentExpr struct {
		Symbol   *Token
		Variable *Variable
	}

	BuiltinExpr struct {
		Fn *Token
	}

	FnExpr struct {
		Token        *Token
		FormalParams []*FormalParam
		Body         *BlockNode

		// set by analyzer
		NumLocals      int
		NumCaptures    int
		ParentCaptures []*Variable
	}

	FormalParam struct {
		Ident   *IdentExpr
		IsConst bool
	}

	InvokeExpr struct {
		Operand Expression
		LParen  *Token
		Params  []Expression
		RParen  *Token
	}

	ListExpr struct {
		LBracket *Token
		Elems    []Expression
		RBracket *Token
	}

	SetExpr struct {
		SetToken *Token
		LBrace   *Token
		Elems    []Expression
		RBrace   *Token
	}

	TupleExpr struct {
		LParen *Token
		Elems  []Expression
		RParen *Token
	}

	StructExpr struct {
		StructToken *Token
		LBrace      *Token
		Keys        []*Token
		Values      []Expression
		RBrace      *Token

		// The index of the struct expression in the local variable array.
		// '-1' means that the struct is not referenced by a 'this', and thus
		// is not stored in the local variable array
		LocalThisIndex int
	}

	ThisExpr struct {
		Token    *Token
		Variable *Variable
	}

	FieldExpr struct {
		Operand Expression
		Key     *Token
	}

	DictExpr struct {
		DictToken *Token
		LBrace    *Token
		Entries   []*DictEntryExpr
		RBrace    *Token
	}

	DictEntryExpr struct {
		Key   Expression
		Value Expression
	}

	IndexExpr struct {
		Operand  Expression
		LBracket *Token
		Index    Expression
		RBracket *Token
	}

	SliceExpr struct {
		Operand  Expression
		LBracket *Token
		From     Expression
		To       Expression
		RBracket *Token
	}

	SliceFromExpr struct {
		Operand  Expression
		LBracket *Token
		From     Expression
		RBracket *Token
	}

	SliceToExpr struct {
		Operand  Expression
		LBracket *Token
		To       Expression
		RBracket *Token
	}
)

//--------------------------------------------------------------
// markers

func (*ImportStmt) stmtMarker()   {}
func (*ConstStmt) stmtMarker()    {}
func (*LetStmt) stmtMarker()      {}
func (*NamedFnStmt) stmtMarker()  {}
func (*IfStmt) stmtMarker()       {}
func (*WhileStmt) stmtMarker()    {}
func (*ForStmt) stmtMarker()      {}
func (*SwitchStmt) stmtMarker()   {}
func (*BreakStmt) stmtMarker()    {}
func (*ContinueStmt) stmtMarker() {}
func (*ReturnStmt) stmtMarker()   {}
func (*ThrowStmt) stmtMarker()    {}
func (*TryStmt) stmtMarker()      {}
func (*GoStmt) stmtMarker()       {}
func (*ExprStmt) stmtMarker() {}

func (*WhileStmt) loopMarker() {}
func (*ForStmt) loopMarker()   {}

func (*AssignmentExpr) exprMarker() {}
func (*TernaryExpr) exprMarker()    {}
func (*BinaryExpr) exprMarker()     {}
func (*UnaryExpr) exprMarker()      {}
func (*PostfixExpr) exprMarker()    {}
func (*BasicExpr) exprMarker()      {}
func (*IdentExpr) exprMarker()      {}
func (*BuiltinExpr) exprMarker()    {}
func (*FnExpr) exprMarker()         {}
func (*InvokeExpr) exprMarker()     {}
func (*ListExpr) exprMarker()       {}
func (*SetExpr) exprMarker()        {}
func (*TupleExpr) exprMarker()      {}
func (*StructExpr) exprMarker()     {}
func (*ThisExpr) exprMarker()       {}
func (*FieldExpr) exprMarker()      {}
func (*DictExpr) exprMarker()       {}
func (*DictEntryExpr) exprMarker()  {}
func (*IndexExpr) exprMarker()      {}
func (*SliceExpr) exprMarker()      {}
func (*SliceFromExpr) exprMarker()  {}
func (*SliceToExpr) exprMarker()    {}

func (*IdentExpr) assignableMarker()   {}
func (*BuiltinExpr) assignableMarker() {}
func (*FieldExpr) assignableMarker()   {}
func (*IndexExpr) assignableMarker()   {}

//--------------------------------------------------------------
// Begin, End

func (n *BlockNode) Begin() Pos { return n.LBrace.Position }
func (n *BlockNode) End() Pos {
	if n.RBrace == nil {
		return n.Statements[len(n.Statements)-1].End()
	} else {
		return n.RBrace.Position
	}
}

func (n *ImportStmt) Begin() Pos { return n.Token.Position }
func (n *ImportStmt) End() Pos   { return n.Ident.End() }

func (n *DeclNode) Begin() Pos { return n.Ident.Begin() }
func (n *DeclNode) End() Pos {
	if n.Val == nil {
		return n.Ident.End()
	} else {
		return n.Val.End()
	}
}

func (n *ConstStmt) Begin() Pos { return n.Token.Position }
func (n *ConstStmt) End() Pos   { return n.Decls[len(n.Decls)-1].End() }

func (n *LetStmt) Begin() Pos { return n.Token.Position }
func (n *LetStmt) End() Pos   { return n.Decls[len(n.Decls)-1].End() }

func (n *NamedFnStmt) Begin() Pos { return n.Token.Position }
func (n *NamedFnStmt) End() Pos   { return n.Func.End() }

func (n *IfStmt) Begin() Pos { return n.Token.Position }
func (n *IfStmt) End() Pos {
	if n.Else == nil {
		return n.Then.End()
	} else {
		return n.Else.End()
	}
}

func (n *WhileStmt) Begin() Pos { return n.Token.Position }
func (n *WhileStmt) End() Pos   { return n.Body.End() }

func (n *ForStmt) Begin() Pos { return n.Token.Position }
func (n *ForStmt) End() Pos   { return n.Body.End() }

func (n *SwitchStmt) Begin() Pos { return n.Token.Position }
func (n *SwitchStmt) End() Pos   { return n.RBrace.Position }

func (n *CaseNode) Begin() Pos { return n.Token.Position }
func (n *CaseNode) End() Pos   { return n.Body[len(n.Body)-1].End() }

func (n *DefaultNode) Begin() Pos { return n.Token.Position }
func (n *DefaultNode) End() Pos   { return n.Body[len(n.Body)-1].End() }

func (n *BreakStmt) Begin() Pos { return n.Token.Position }
func (n *BreakStmt) End() Pos   { return n.Token.Position.Advance(len("break") - 1) }

func (n *ContinueStmt) Begin() Pos { return n.Token.Position }
func (n *ContinueStmt) End() Pos   { return n.Token.Position.Advance(len("continue") - 1) }

func (n *ReturnStmt) Begin() Pos { return n.Token.Position }
func (n *ReturnStmt) End() Pos {
	if n.Val == nil {
		return n.Token.Position.Advance(len("return") - 1)
	} else {
		return n.Val.End()
	}
}

func (n *ThrowStmt) Begin() Pos { return n.Token.Position }
func (n *ThrowStmt) End() Pos   { return n.Val.End() }

func (n *TryStmt) Begin() Pos { return n.TryToken.Position }
func (n *TryStmt) End() Pos {
	if n.FinallyToken == nil {
		return n.CatchBlock.End()
	} else {
		return n.FinallyBlock.End()
	}
}

func (n *GoStmt) Begin() Pos { return n.Token.Position }
func (n *GoStmt) End() Pos   { return n.Invocation.End() }

func (n *ExprStmt) Begin() Pos { return n.Expr.Begin() }
func (n *ExprStmt) End() Pos   { return n.Expr.End() }

func (n *AssignmentExpr) Begin() Pos { return n.Assignee.Begin() }
func (n *AssignmentExpr) End() Pos   { return n.Val.End() }

func (n *TernaryExpr) Begin() Pos { return n.Cond.Begin() }
func (n *TernaryExpr) End() Pos   { return n.Else.End() }

func (n *BinaryExpr) Begin() Pos { return n.Lhs.Begin() }
func (n *BinaryExpr) End() Pos   { return n.Rhs.End() }

func (n *UnaryExpr) Begin() Pos { return n.Op.Position }
func (n *UnaryExpr) End() Pos   { return n.Operand.End() }

func (n *PostfixExpr) Begin() Pos { return n.Assignee.Begin() }
func (n *PostfixExpr) End() Pos   { return n.Op.Position }

func (n *BasicExpr) Begin() Pos { return n.Token.Position }
func (n *BasicExpr) End() Pos {
	return Pos{
		n.Token.Position.Line,
		n.Token.Position.Col + len(n.Token.Text) - 1}
}

func (n *IdentExpr) Begin() Pos { return n.Symbol.Position }
func (n *IdentExpr) End() Pos {
	return Pos{
		n.Symbol.Position.Line,
		n.Symbol.Position.Col + len(n.Symbol.Text) - 1}
}

func (n *BuiltinExpr) Begin() Pos { return n.Fn.Position }
func (n *BuiltinExpr) End() Pos {
	return Pos{
		n.Fn.Position.Line,
		n.Fn.Position.Col + len(n.Fn.Text) - 1}
}

func (n *FnExpr) Begin() Pos { return n.Token.Position }
func (n *FnExpr) End() Pos   { return n.Body.End() }

func (n *InvokeExpr) Begin() Pos { return n.Operand.Begin() }
func (n *InvokeExpr) End() Pos   { return n.RParen.Position }

func (n *ListExpr) Begin() Pos { return n.LBracket.Position }
func (n *ListExpr) End() Pos   { return n.RBracket.Position }

func (n *SetExpr) Begin() Pos { return n.SetToken.Position }
func (n *SetExpr) End() Pos   { return n.RBrace.Position }

func (n *TupleExpr) Begin() Pos { return n.LParen.Position }
func (n *TupleExpr) End() Pos   { return n.RParen.Position }

func (n *StructExpr) Begin() Pos { return n.StructToken.Position }
func (n *StructExpr) End() Pos   { return n.RBrace.Position }

func (n *ThisExpr) Begin() Pos { return n.Token.Position }
func (n *ThisExpr) End() Pos {
	return Pos{
		n.Token.Position.Line,
		n.Token.Position.Col + len("this") - 1}
}

func (n *FieldExpr) Begin() Pos { return n.Operand.Begin() }
func (n *FieldExpr) End() Pos   { return n.Key.Position }

func (n *DictExpr) Begin() Pos { return n.DictToken.Position }
func (n *DictExpr) End() Pos   { return n.RBrace.Position }

func (n *DictEntryExpr) Begin() Pos { return n.Key.Begin() }
func (n *DictEntryExpr) End() Pos   { return n.Value.End() }

func (n *IndexExpr) Begin() Pos { return n.Operand.Begin() }
func (n *IndexExpr) End() Pos   { return n.RBracket.Position }

func (n *SliceExpr) Begin() Pos     { return n.Operand.Begin() }
func (n *SliceExpr) End() Pos       { return n.RBracket.Position }
func (n *SliceFromExpr) Begin() Pos { return n.Operand.Begin() }
func (n *SliceFromExpr) End() Pos   { return n.RBracket.Position }
func (n *SliceToExpr) Begin() Pos   { return n.Operand.Begin() }
func (n *SliceToExpr) End() Pos     { return n.RBracket.Position }

//--------------------------------------------------------------
// string

func (blk *BlockNode) String() string {
	var buf bytes.Buffer
	buf.WriteString("{ ")
	writeStatements(blk.Statements, &buf)
	buf.WriteString(" }")
	return buf.String()
}

func (imp *ImportStmt) String() string {
	buf := new(bytes.Buffer)
	buf.WriteString("import ")
	buf.WriteString(imp.Ident.String())
	buf.WriteString(";")
	return buf.String()
}

func (cns *ConstStmt) String() string {
	buf := new(bytes.Buffer)
	buf.WriteString("const ")
	buf.WriteString(stringDecls(cns.Decls))
	buf.WriteString(";")
	return buf.String()
}

func (let *LetStmt) String() string {
	buf := new(bytes.Buffer)
	buf.WriteString("let ")
	buf.WriteString(stringDecls(let.Decls))
	buf.WriteString(";")
	return buf.String()
}

func (nf *NamedFnStmt) String() string {
	buf := new(bytes.Buffer)
	buf.WriteString("fn ")
	buf.WriteString(nf.Ident.String())
	buf.WriteString(stringFormalParams(nf.Func.FormalParams))
	buf.WriteString(" ")
	buf.WriteString(nf.Func.Body.String())
	buf.WriteString(";")
	return buf.String()
}

func stringDecls(decls []*DeclNode) string {
	buf := new(bytes.Buffer)
	for i, d := range decls {
		if i > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(fmt.Sprintf("%v", d.Ident))
		if d.Val != nil {
			buf.WriteString(fmt.Sprintf(" = %v", d.Val))
		}
	}
	return buf.String()
}

func (asn *AssignmentExpr) String() string {
	return fmt.Sprintf("(%v = %v)", asn.Assignee, asn.Val)
}

func (ifn *IfStmt) String() string {
	if ifn.Else == nil {
		return fmt.Sprintf("if %v %v;", ifn.Cond, ifn.Then)
	} else {
		return fmt.Sprintf("if %v %v else %v;", ifn.Cond, ifn.Then, ifn.Else)
	}
}

func (wh *WhileStmt) String() string {
	return fmt.Sprintf("while %v %v;", wh.Cond, wh.Body)
}

func (fr *ForStmt) String() string {
	if len(fr.Idents) == 1 {
		return fmt.Sprintf("for %v in %v %v;", fr.Idents[0], fr.Iterable, fr.Body)
	} else {
		return fmt.Sprintf("for %s in %v %v;", stringIdents(fr.Idents), fr.Iterable, fr.Body)
	}
}

func stringIdents(idents []*IdentExpr) string {
	var buf bytes.Buffer

	buf.WriteString("(")
	for idx, p := range idents {
		if idx > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(p.String())
	}
	buf.WriteString(")")

	return buf.String()
}

func (sw *SwitchStmt) String() string {
	var buf bytes.Buffer

	buf.WriteString("switch ")
	if sw.Item != nil {
		buf.WriteString(fmt.Sprintf("%v", sw.Item))
		buf.WriteString(" ")
	}

	buf.WriteString("{ ")
	for i, c := range sw.Cases {
		if i > 0 {
			buf.WriteString(" ")
		}
		buf.WriteString(fmt.Sprintf("%v", c))
	}
	if sw.DefaultNode != nil {
		buf.WriteString(fmt.Sprintf("%v", sw.DefaultNode))
	}
	buf.WriteString(" };")

	return buf.String()
}

func (cs *CaseNode) String() string {
	var buf bytes.Buffer

	buf.WriteString("case ")
	for i, m := range cs.Matches {
		if i > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(fmt.Sprintf("%v", m))
	}

	buf.WriteString(": ")
	writeStatements(cs.Body, &buf)

	return buf.String()
}

func (def *DefaultNode) String() string {
	var buf bytes.Buffer

	buf.WriteString(" default: ")
	writeStatements(def.Body, &buf)

	return buf.String()
}

func (br *BreakStmt) String() string {
	return "break;"
}

func (cn *ContinueStmt) String() string {
	return "continue;"
}

func (rt *ReturnStmt) String() string {
	if rt.Val == nil {
		return "return;"
	} else {
		return fmt.Sprintf("return %v;", rt.Val)
	}
}

func (t *ThrowStmt) String() string {
	return fmt.Sprintf("throw %v;", t.Val)
}

func (t *TryStmt) String() string {

	var buf bytes.Buffer

	buf.WriteString("try ")
	buf.WriteString(t.TryBlock.String())

	if t.CatchToken != nil {
		buf.WriteString(" catch ")
		buf.WriteString(t.CatchIdent.String())
		buf.WriteString(" ")
		buf.WriteString(t.CatchBlock.String())
	}

	if t.FinallyToken != nil {
		buf.WriteString(" finally ")
		buf.WriteString(t.FinallyBlock.String())
	}
	buf.WriteString(";")

	return buf.String()
}

func (g *GoStmt) String() string {
	return fmt.Sprintf("go %v;", g.Invocation)
}

func (es *ExprStmt) String() string {
	return fmt.Sprintf("%v;", es.Expr)
}

func (trn *TernaryExpr) String() string {
	return fmt.Sprintf("(%v ? %v : %v)", trn.Cond, trn.Then, trn.Else)
}

func (bin *BinaryExpr) String() string {
	return fmt.Sprintf("(%v %s %v)", bin.Lhs, bin.Op.Text, bin.Rhs)
}

func (unary *UnaryExpr) String() string {
	return fmt.Sprintf("%s%v", unary.Op.Text, unary.Operand)
}

func (pf *PostfixExpr) String() string {
	return fmt.Sprintf("%v%s", pf.Assignee, pf.Op.Text)
}

func (basic *BasicExpr) String() string {
	if basic.Token.Kind == Str {
		// TODO escape embedded delim, \n, \r, \t, \u
		return strings.Join([]string{"'", basic.Token.Text, "'"}, "")
	} else {
		return basic.Token.Text
	}
}

func (ident *IdentExpr) String() string {
	return ident.Symbol.Text
}

func (blt *BuiltinExpr) String() string {
	return blt.Fn.Text
}

func (fn *FnExpr) String() string {
	var buf bytes.Buffer

	buf.WriteString("fn")
	buf.WriteString(stringFormalParams(fn.FormalParams))
	buf.WriteString(" ")
	buf.WriteString(fn.Body.String())

	return buf.String()
}

func stringFormalParams(params []*FormalParam) string {
	var buf bytes.Buffer

	buf.WriteString("(")
	for idx, p := range params {
		if idx > 0 {
			buf.WriteString(", ")
		}
		if p.IsConst {
			buf.WriteString("const ")
		}
		buf.WriteString(p.Ident.String())
	}
	buf.WriteString(")")

	return buf.String()
}

func (inv *InvokeExpr) String() string {
	var buf bytes.Buffer
	buf.WriteString(inv.Operand.String())
	buf.WriteString("(")
	for idx, p := range inv.Params {
		if idx > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(p.String())
	}
	buf.WriteString(")")
	return buf.String()
}

func (ls *ListExpr) String() string {
	var buf bytes.Buffer
	buf.WriteString("[ ")
	for idx, v := range ls.Elems {
		if idx > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(v.String())
	}
	buf.WriteString(" ]")
	return buf.String()
}

func (s *SetExpr) String() string {
	var buf bytes.Buffer
	buf.WriteString("set { ")
	for idx, v := range s.Elems {
		if idx > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(v.String())
	}
	buf.WriteString(" }")
	return buf.String()
}

func (tp *TupleExpr) String() string {
	var buf bytes.Buffer
	buf.WriteString("(")
	for idx, v := range tp.Elems {
		if idx > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(v.String())
	}
	buf.WriteString(")")
	return buf.String()
}

func (stc *StructExpr) String() string {
	var buf bytes.Buffer
	buf.WriteString("struct")

	buf.WriteString(" { ")
	for idx, k := range stc.Keys {
		if idx > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(k.Text)
		buf.WriteString(": ")
		buf.WriteString(stc.Values[idx].String())
	}
	buf.WriteString(" }")
	return buf.String()
}

func (this *ThisExpr) String() string {
	return "this"
}

func (f *FieldExpr) String() string {
	var buf bytes.Buffer
	buf.WriteString(f.Operand.String())
	buf.WriteString(".")
	buf.WriteString(f.Key.Text)
	return buf.String()
}

func (dict *DictExpr) String() string {
	var buf bytes.Buffer
	buf.WriteString("dict { ")
	for idx, e := range dict.Entries {
		if idx > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(e.String())
	}
	buf.WriteString(" }")
	return buf.String()
}

func (de *DictEntryExpr) String() string {
	var buf bytes.Buffer
	buf.WriteString(de.Key.String())
	buf.WriteString(": ")
	buf.WriteString(de.Value.String())
	return buf.String()
}

func (i *IndexExpr) String() string {
	var buf bytes.Buffer
	buf.WriteString(i.Operand.String())
	buf.WriteString("[")
	buf.WriteString(i.Index.String())
	buf.WriteString("]")
	return buf.String()
}

func (s *SliceExpr) String() string {
	var buf bytes.Buffer
	buf.WriteString(s.Operand.String())
	buf.WriteString("[")
	buf.WriteString(s.From.String())
	buf.WriteString(":")
	buf.WriteString(s.To.String())
	buf.WriteString("]")
	return buf.String()
}

func (s *SliceFromExpr) String() string {
	var buf bytes.Buffer
	buf.WriteString(s.Operand.String())
	buf.WriteString("[")
	buf.WriteString(s.From.String())
	buf.WriteString(":]")
	return buf.String()
}

func (s *SliceToExpr) String() string {
	var buf bytes.Buffer
	buf.WriteString(s.Operand.String())
	buf.WriteString("[:")
	buf.WriteString(s.To.String())
	buf.WriteString("]")
	return buf.String()
}

func writeStatements(stmts []Statement, buf *bytes.Buffer) {
	for idx, n := range stmts {
		if idx > 0 {
			buf.WriteString(" ")
		}
		buf.WriteString(n.String())
		if _, ok := n.(Expression); ok {
			buf.WriteString(";")
		}
	}
}

//--------------------------------------------------------------
// A Variable points to a Ref.  Variables are defined either
// as formal params for a Function, or via LetStmt or ConstStmt, or via
// the capture mechanism.

type Variable struct {
	Symbol    string
	Index     int
	IsConst   bool
	IsCapture bool
}

func (v *Variable) String() string {
	return fmt.Sprintf("(%d,%v,%v)", v.Index, v.IsConst, v.IsCapture)
}

type VariableArray []*Variable

// Variables are sorted by Index
func (v VariableArray) Len() int {
	return len(v)
}
func (v VariableArray) Swap(i, j int) {
	v[i], v[j] = v[j], v[i]
}
func (v VariableArray) Less(i, j int) bool {
	return v[i].Index < v[j].Index
}
