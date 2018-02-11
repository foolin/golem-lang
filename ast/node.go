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

type (
	// Node is a node in an Abstract Syntax Tree
	Node interface {
		fmt.Stringer
		Traverse(Visitor)
		Begin() Pos
		End() Pos
	}

	// Statement is a Node that is a statement
	Statement interface {
		Node
		stmtMarker()
	}

	// Loop is a Statement that is a loop
	Loop interface {
		Statement
		loopMarker() // nolint: megacheck
	}

	// Expression is a Node that is an expression
	Expression interface {
		Node
		exprMarker()
	}

	// Assignable is an  Expression that is assignable
	Assignable interface {
		Expression
		assignableMarker()
	}
)

type (

	// ImportStmt is an 'import' statement
	ImportStmt struct {
		Token  *Token
		Idents []*IdentExpr
	}

	// ConstStmt is a 'const' statement
	ConstStmt struct {
		Token *Token
		Decls []*DeclNode
	}

	// LetStmt is a 'let' statement
	LetStmt struct {
		Token *Token
		Decls []*DeclNode
	}

	// NamedFnStmt is a named function statement
	NamedFnStmt struct {
		Token *Token
		Ident *IdentExpr
		Func  *FnExpr
	}

	// IfStmt is a 'if' statement
	IfStmt struct {
		Token *Token
		Cond  Expression
		Then  *BlockNode
		Else  Node // either a BlockNode, or another IfStmt
	}

	// WhileStmt is a 'while' statement
	WhileStmt struct {
		Token *Token
		Cond  Expression
		Body  *BlockNode
	}

	// ForStmt is a 'for' statement
	ForStmt struct {
		Token         *Token
		Idents        []*IdentExpr
		IterableIdent *IdentExpr
		Iterable      Expression
		Body          *BlockNode

		// Scope defines the scope for the Idents
		Scope Scope
	}

	// SwitchStmt is a 'switch' statement
	SwitchStmt struct {
		Token       *Token
		Item        Expression
		LBrace      *Token
		Cases       []*CaseNode
		DefaultNode *DefaultNode
		RBrace      *Token
	}

	// BreakStmt is a 'break' statement
	BreakStmt struct {
		Token *Token
	}

	// ContinueStmt is a 'continue' statement
	ContinueStmt struct {
		Token *Token
	}

	// ReturnStmt is a 'return' statement
	ReturnStmt struct {
		Token *Token
		Val   Expression
	}

	// ThrowStmt is a 'throw' statement
	ThrowStmt struct {
		Token *Token
		Val   Expression
	}

	// TryStmt is a 'try' statement
	TryStmt struct {
		TryToken     *Token
		TryBlock     *BlockNode
		CatchToken   *Token
		CatchIdent   *IdentExpr
		CatchBlock   *BlockNode
		FinallyToken *Token
		FinallyBlock *BlockNode

		// CatchScope defines the scope for the CatchIdent
		CatchScope Scope
	}

	// GoStmt is a 'go' statement
	GoStmt struct {
		Token      *Token
		Invocation *InvokeExpr
	}

	// ExprStmt is a Statement that contains an Expression
	ExprStmt struct {
		Expr Expression
	}

	// BlockNode is a sequence of Statements
	BlockNode struct {
		LBrace     *Token
		Statements []Statement
		RBrace     *Token

		Scope Scope
	}

	// DeclNode is a declaration
	DeclNode struct {
		Ident *IdentExpr
		Val   Expression
	}

	// CaseNode is a 'case' clause in a 'switch' statement.
	CaseNode struct {
		Token   *Token
		Matches []Expression
		Body    []Statement
	}

	// DefaultNode is a 'default' clause in a 'switch' statement.
	DefaultNode struct {
		Token *Token
		Body  []Statement
	}

	// AssignmentExpr is an assigment expressions
	AssignmentExpr struct {
		Assignee Assignable
		Eq       *Token
		Val      Expression
	}

	// TernaryExpr is a ternary expression
	TernaryExpr struct {
		Cond Expression
		Then Expression
		Else Expression
	}

	// BinaryExpr is a binary expression
	BinaryExpr struct {
		LHS Expression
		Op  *Token
		RHS Expression
	}

	// UnaryExpr is a unary expression
	UnaryExpr struct {
		Op      *Token
		Operand Expression
	}

	// PostfixExpr is a postfix expression
	PostfixExpr struct {
		Assignee Assignable
		Op       *Token
	}

	// BasicExpr is a basic expression
	BasicExpr struct {
		Token *Token
	}

	// IdentExpr is an identifier expression
	IdentExpr struct {
		Symbol   *Token
		Variable Variable
	}

	// BuiltinExpr is a builtin-value expression
	BuiltinExpr struct {
		Fn *Token
	}

	// FnExpr is a function expression
	FnExpr struct {
		Token        *Token
		FormalParams []*FormalParam
		Body         *BlockNode

		Scope FuncScope
	}

	// FormalParam is a formal parameter in a function expression
	FormalParam struct {
		Ident   *IdentExpr
		IsConst bool
	}

	// InvokeExpr is an invocation expression
	InvokeExpr struct {
		Operand Expression
		LParen  *Token
		Params  []Expression
		RParen  *Token
	}

	// ListExpr is a list expression
	ListExpr struct {
		LBracket *Token
		Elems    []Expression
		RBracket *Token
	}

	// SetExpr is a 'set' expression
	SetExpr struct {
		SetToken *Token
		LBrace   *Token
		Elems    []Expression
		RBrace   *Token
	}

	// TupleExpr is a tuple expression
	TupleExpr struct {
		LParen *Token
		Elems  []Expression
		RParen *Token
	}

	// StructExpr is a struct expression
	StructExpr struct {
		StructToken *Token
		LBrace      *Token
		Keys        []*Token
		Values      []Expression
		RBrace      *Token

		// ThisScope will always either be empty, or contain
		// a single 'this' Variable.
		Scope StructScope
	}

	// ThisExpr is a 'this' expression
	ThisExpr struct {
		Token    *Token
		Variable Variable
	}

	// FieldExpr is a field expression
	FieldExpr struct {
		Operand Expression
		Key     *Token
	}

	// DictExpr is a 'dict' expression
	DictExpr struct {
		DictToken *Token
		LBrace    *Token
		Entries   []*DictEntryExpr
		RBrace    *Token
	}

	// DictEntryExpr is an entry in a DictExpr
	DictEntryExpr struct {
		Key   Expression
		Value Expression
	}

	// IndexExpr is an index expression
	IndexExpr struct {
		Operand  Expression
		LBracket *Token
		Index    Expression
		RBracket *Token
	}

	// SliceExpr is a slice expression
	SliceExpr struct {
		Operand  Expression
		LBracket *Token
		From     Expression
		To       Expression
		RBracket *Token
	}

	// SliceFromExpr is a slice expression
	SliceFromExpr struct {
		Operand  Expression
		LBracket *Token
		From     Expression
		RBracket *Token
	}

	// SliceToExpr is a slice expression
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
func (*ExprStmt) stmtMarker()     {}

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

// Begin BlockNode
func (n *BlockNode) Begin() Pos { return n.LBrace.Position }

// End BlockNode
func (n *BlockNode) End() Pos {
	if n.RBrace == nil {
		return n.Statements[len(n.Statements)-1].End()
	}
	return n.RBrace.Position
}

// Begin ImportStmt
func (n *ImportStmt) Begin() Pos { return n.Token.Position }

// End ImportStmt
func (n *ImportStmt) End() Pos { return n.Idents[len(n.Idents)-1].End() }

// Begin DeclNode
func (n *DeclNode) Begin() Pos { return n.Ident.Begin() }

// End DeclNode
func (n *DeclNode) End() Pos {
	if n.Val == nil {
		return n.Ident.End()
	}
	return n.Val.End()
}

// Begin ConstStmt
func (n *ConstStmt) Begin() Pos { return n.Token.Position }

// End ConstStmt
func (n *ConstStmt) End() Pos { return n.Decls[len(n.Decls)-1].End() }

// Begin LetStmt
func (n *LetStmt) Begin() Pos { return n.Token.Position }

// End LetStmt
func (n *LetStmt) End() Pos { return n.Decls[len(n.Decls)-1].End() }

// Begin NamedFnStmt
func (n *NamedFnStmt) Begin() Pos { return n.Token.Position }

// End NamedFnStmt
func (n *NamedFnStmt) End() Pos { return n.Func.End() }

// Begin IfStmt
func (n *IfStmt) Begin() Pos { return n.Token.Position }

// End IfStmt
func (n *IfStmt) End() Pos {
	if n.Else == nil {
		return n.Then.End()
	}
	return n.Else.End()
}

// Begin WhileStmt
func (n *WhileStmt) Begin() Pos { return n.Token.Position }

// End WhileStmt
func (n *WhileStmt) End() Pos { return n.Body.End() }

// Begin ForStmt
func (n *ForStmt) Begin() Pos { return n.Token.Position }

// End ForStmt
func (n *ForStmt) End() Pos { return n.Body.End() }

// Begin SwitchStmt
func (n *SwitchStmt) Begin() Pos { return n.Token.Position }

// End SwitchStmt
func (n *SwitchStmt) End() Pos { return n.RBrace.Position }

// Begin CaseNode
func (n *CaseNode) Begin() Pos { return n.Token.Position }

// End CaseNode
func (n *CaseNode) End() Pos { return n.Body[len(n.Body)-1].End() }

// Begin DefaultNode
func (n *DefaultNode) Begin() Pos { return n.Token.Position }

// End DefaultNode
func (n *DefaultNode) End() Pos { return n.Body[len(n.Body)-1].End() }

// Begin BreakStmt
func (n *BreakStmt) Begin() Pos { return n.Token.Position }

// End BreakStmt
func (n *BreakStmt) End() Pos { return n.Token.Position.Advance(len("break") - 1) }

// Begin ContinueStmt
func (n *ContinueStmt) Begin() Pos { return n.Token.Position }

// End ContinueStmt
func (n *ContinueStmt) End() Pos { return n.Token.Position.Advance(len("continue") - 1) }

// Begin ReturnStmt
func (n *ReturnStmt) Begin() Pos { return n.Token.Position }

// End ReturnStmt
func (n *ReturnStmt) End() Pos {
	if n.Val == nil {
		return n.Token.Position.Advance(len("return") - 1)
	}
	return n.Val.End()
}

// Begin ThrowStmt
func (n *ThrowStmt) Begin() Pos { return n.Token.Position }

// End ThrowStmt
func (n *ThrowStmt) End() Pos { return n.Val.End() }

// Begin TryStmt
func (n *TryStmt) Begin() Pos { return n.TryToken.Position }

// End TryStmt
func (n *TryStmt) End() Pos {
	if n.FinallyToken == nil {
		return n.CatchBlock.End()
	}
	return n.FinallyBlock.End()
}

// Begin GoStmt
func (n *GoStmt) Begin() Pos { return n.Token.Position }

// End GoStmt
func (n *GoStmt) End() Pos { return n.Invocation.End() }

// Begin ExprStmt
func (n *ExprStmt) Begin() Pos { return n.Expr.Begin() }

// End ExprStmt
func (n *ExprStmt) End() Pos { return n.Expr.End() }

// Begin AssignmentExpr
func (n *AssignmentExpr) Begin() Pos { return n.Assignee.Begin() }

// End AssignmentExpr
func (n *AssignmentExpr) End() Pos { return n.Val.End() }

// Begin TernaryExpr
func (n *TernaryExpr) Begin() Pos { return n.Cond.Begin() }

// End TernaryExpr
func (n *TernaryExpr) End() Pos { return n.Else.End() }

// Begin BinaryExpr
func (n *BinaryExpr) Begin() Pos { return n.LHS.Begin() }

// End BinaryExpr
func (n *BinaryExpr) End() Pos { return n.RHS.End() }

// Begin UnaryExpr
func (n *UnaryExpr) Begin() Pos { return n.Op.Position }

// End UnaryExpr
func (n *UnaryExpr) End() Pos { return n.Operand.End() }

// Begin PostfixExpr
func (n *PostfixExpr) Begin() Pos { return n.Assignee.Begin() }

// End PostfixExpr
func (n *PostfixExpr) End() Pos { return n.Op.Position }

// Begin BasicExpr
func (n *BasicExpr) Begin() Pos { return n.Token.Position }

// End BasicExpr
func (n *BasicExpr) End() Pos {
	return Pos{
		n.Token.Position.Line,
		n.Token.Position.Col + len(n.Token.Text) - 1}
}

// Begin IdentExpr
func (n *IdentExpr) Begin() Pos { return n.Symbol.Position }

// End IdentExpr
func (n *IdentExpr) End() Pos {
	return Pos{
		n.Symbol.Position.Line,
		n.Symbol.Position.Col + len(n.Symbol.Text) - 1}
}

// Begin BuiltinExpr
func (n *BuiltinExpr) Begin() Pos { return n.Fn.Position }

// End BuiltinExpr
func (n *BuiltinExpr) End() Pos {
	return Pos{
		n.Fn.Position.Line,
		n.Fn.Position.Col + len(n.Fn.Text) - 1}
}

// Begin FnExpr
func (n *FnExpr) Begin() Pos { return n.Token.Position }

// End FnExpr
func (n *FnExpr) End() Pos { return n.Body.End() }

// Begin InvokeExpr
func (n *InvokeExpr) Begin() Pos { return n.Operand.Begin() }

// End InvokeExpr
func (n *InvokeExpr) End() Pos { return n.RParen.Position }

// Begin ListExpr
func (n *ListExpr) Begin() Pos { return n.LBracket.Position }

// End ListExpr
func (n *ListExpr) End() Pos { return n.RBracket.Position }

// Begin SetExpr
func (n *SetExpr) Begin() Pos { return n.SetToken.Position }

// End SetExpr
func (n *SetExpr) End() Pos { return n.RBrace.Position }

// Begin TupleExpr
func (n *TupleExpr) Begin() Pos { return n.LParen.Position }

// End TupleExpr
func (n *TupleExpr) End() Pos { return n.RParen.Position }

// Begin StructExpr
func (n *StructExpr) Begin() Pos { return n.StructToken.Position }

// End StructExpr
func (n *StructExpr) End() Pos { return n.RBrace.Position }

// Begin ThisExpr
func (n *ThisExpr) Begin() Pos { return n.Token.Position }

// End ThisExpr
func (n *ThisExpr) End() Pos {
	return Pos{
		n.Token.Position.Line,
		n.Token.Position.Col + len("this") - 1}
}

// Begin FieldExpr
func (n *FieldExpr) Begin() Pos { return n.Operand.Begin() }

// End FieldExpr
func (n *FieldExpr) End() Pos { return n.Key.Position }

// Begin DictExpr
func (n *DictExpr) Begin() Pos { return n.DictToken.Position }

// End DictExpr
func (n *DictExpr) End() Pos { return n.RBrace.Position }

// Begin DictEntryExpr
func (n *DictEntryExpr) Begin() Pos { return n.Key.Begin() }

// End DictEntryExpr
func (n *DictEntryExpr) End() Pos { return n.Value.End() }

// Begin IndexExpr
func (n *IndexExpr) Begin() Pos { return n.Operand.Begin() }

// End IndexExpr
func (n *IndexExpr) End() Pos { return n.RBracket.Position }

// Begin SliceExpr
func (n *SliceExpr) Begin() Pos { return n.Operand.Begin() }

// End SliceExpr
func (n *SliceExpr) End() Pos { return n.RBracket.Position }

// Begin SliceFromExpr
func (n *SliceFromExpr) Begin() Pos { return n.Operand.Begin() }

// End SliceFromExpr
func (n *SliceFromExpr) End() Pos { return n.RBracket.Position }

// Begin SliceToExpr
func (n *SliceToExpr) Begin() Pos { return n.Operand.Begin() }

// End SliceToExpr
func (n *SliceToExpr) End() Pos { return n.RBracket.Position }

//--------------------------------------------------------------
// string

func (n *BlockNode) String() string {
	var buf bytes.Buffer
	buf.WriteString("{ ")
	writeStatements(n.Statements, &buf)
	buf.WriteString(" }")
	return buf.String()
}

func (n *ImportStmt) String() string {
	buf := new(bytes.Buffer)
	buf.WriteString("import ")
	for i, ident := range n.Idents {
		if i > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(ident.String())
	}
	buf.WriteString(";")
	return buf.String()
}

func (n *ConstStmt) String() string {
	buf := new(bytes.Buffer)
	buf.WriteString("const ")
	buf.WriteString(stringDecls(n.Decls))
	buf.WriteString(";")
	return buf.String()
}

func (n *LetStmt) String() string {
	buf := new(bytes.Buffer)
	buf.WriteString("let ")
	buf.WriteString(stringDecls(n.Decls))
	buf.WriteString(";")
	return buf.String()
}

func (n *NamedFnStmt) String() string {
	buf := new(bytes.Buffer)
	buf.WriteString("fn ")
	buf.WriteString(n.Ident.String())
	buf.WriteString(stringFormalParams(n.Func.FormalParams))
	buf.WriteString(" ")
	buf.WriteString(n.Func.Body.String())
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

func (n *AssignmentExpr) String() string {
	return fmt.Sprintf("(%v = %v)", n.Assignee, n.Val)
}

func (n *IfStmt) String() string {
	if n.Else == nil {
		return fmt.Sprintf("if %v %v;", n.Cond, n.Then)
	}
	return fmt.Sprintf("if %v %v else %v;", n.Cond, n.Then, n.Else)
}

func (n *WhileStmt) String() string {
	return fmt.Sprintf("while %v %v;", n.Cond, n.Body)
}

func (n *ForStmt) String() string {
	if len(n.Idents) == 1 {
		return fmt.Sprintf("for %v in %v %v;", n.Idents[0], n.Iterable, n.Body)
	}
	return fmt.Sprintf("for %s in %v %v;", stringIdents(n.Idents), n.Iterable, n.Body)
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

func (n *SwitchStmt) String() string {
	var buf bytes.Buffer

	buf.WriteString("switch ")
	if n.Item != nil {
		buf.WriteString(fmt.Sprintf("%v", n.Item))
		buf.WriteString(" ")
	}

	buf.WriteString("{ ")
	for i, c := range n.Cases {
		if i > 0 {
			buf.WriteString(" ")
		}
		buf.WriteString(fmt.Sprintf("%v", c))
	}
	if n.DefaultNode != nil {
		buf.WriteString(fmt.Sprintf("%v", n.DefaultNode))
	}
	buf.WriteString(" };")

	return buf.String()
}

func (n *CaseNode) String() string {
	var buf bytes.Buffer

	buf.WriteString("case ")
	for i, m := range n.Matches {
		if i > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(fmt.Sprintf("%v", m))
	}

	buf.WriteString(": ")
	writeStatements(n.Body, &buf)

	return buf.String()
}

func (n *DefaultNode) String() string {
	var buf bytes.Buffer

	buf.WriteString(" default: ")
	writeStatements(n.Body, &buf)

	return buf.String()
}

func (n *BreakStmt) String() string {
	return "break;"
}

func (n *ContinueStmt) String() string {
	return "continue;"
}

func (n *ReturnStmt) String() string {
	if n.Val == nil {
		return "return;"
	}
	return fmt.Sprintf("return %v;", n.Val)
}

func (n *ThrowStmt) String() string {
	return fmt.Sprintf("throw %v;", n.Val)
}

func (n *TryStmt) String() string {

	var buf bytes.Buffer

	buf.WriteString("try ")
	buf.WriteString(n.TryBlock.String())

	if n.CatchToken != nil {
		buf.WriteString(" catch ")
		buf.WriteString(n.CatchIdent.String())
		buf.WriteString(" ")
		buf.WriteString(n.CatchBlock.String())
	}

	if n.FinallyToken != nil {
		buf.WriteString(" finally ")
		buf.WriteString(n.FinallyBlock.String())
	}
	buf.WriteString(";")

	return buf.String()
}

func (n *GoStmt) String() string {
	return fmt.Sprintf("go %v;", n.Invocation)
}

func (n *ExprStmt) String() string {
	return fmt.Sprintf("%v;", n.Expr)
}

func (n *TernaryExpr) String() string {
	return fmt.Sprintf("(%v ? %v : %v)", n.Cond, n.Then, n.Else)
}

func (n *BinaryExpr) String() string {
	return fmt.Sprintf("(%v %s %v)", n.LHS, n.Op.Text, n.RHS)
}

func (n *UnaryExpr) String() string {
	return fmt.Sprintf("%s%v", n.Op.Text, n.Operand)
}

func (n *PostfixExpr) String() string {
	return fmt.Sprintf("%v%s", n.Assignee, n.Op.Text)
}

func (n *BasicExpr) String() string {
	if n.Token.Kind == Str {
		// TODO escape embedded delim, \n, \r, \t, \u
		return strings.Join([]string{"'", n.Token.Text, "'"}, "")
	}
	return n.Token.Text
}

func (n *IdentExpr) String() string {
	return n.Symbol.Text
}

func (n *BuiltinExpr) String() string {
	return n.Fn.Text
}

func (n *FnExpr) String() string {
	var buf bytes.Buffer

	buf.WriteString("fn")
	buf.WriteString(stringFormalParams(n.FormalParams))
	buf.WriteString(" ")
	buf.WriteString(n.Body.String())

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

func (n *InvokeExpr) String() string {
	var buf bytes.Buffer
	buf.WriteString(n.Operand.String())
	buf.WriteString("(")
	for idx, p := range n.Params {
		if idx > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(p.String())
	}
	buf.WriteString(")")
	return buf.String()
}

func (n *ListExpr) String() string {
	var buf bytes.Buffer
	buf.WriteString("[ ")
	for idx, v := range n.Elems {
		if idx > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(v.String())
	}
	buf.WriteString(" ]")
	return buf.String()
}

func (n *SetExpr) String() string {
	var buf bytes.Buffer
	buf.WriteString("set { ")
	for idx, v := range n.Elems {
		if idx > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(v.String())
	}
	buf.WriteString(" }")
	return buf.String()
}

func (n *TupleExpr) String() string {
	var buf bytes.Buffer
	buf.WriteString("(")
	for idx, v := range n.Elems {
		if idx > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(v.String())
	}
	buf.WriteString(")")
	return buf.String()
}

func (n *StructExpr) String() string {
	var buf bytes.Buffer
	buf.WriteString("struct")

	buf.WriteString(" { ")
	for idx, k := range n.Keys {
		if idx > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(k.Text)
		buf.WriteString(": ")
		buf.WriteString(n.Values[idx].String())
	}
	buf.WriteString(" }")
	return buf.String()
}

func (n *ThisExpr) String() string {
	return "this"
}

func (n *FieldExpr) String() string {
	var buf bytes.Buffer
	buf.WriteString(n.Operand.String())
	buf.WriteString(".")
	buf.WriteString(n.Key.Text)
	return buf.String()
}

func (n *DictExpr) String() string {
	var buf bytes.Buffer
	buf.WriteString("dict { ")
	for idx, e := range n.Entries {
		if idx > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(e.String())
	}
	buf.WriteString(" }")
	return buf.String()
}

func (n *DictEntryExpr) String() string {
	var buf bytes.Buffer
	buf.WriteString(n.Key.String())
	buf.WriteString(": ")
	buf.WriteString(n.Value.String())
	return buf.String()
}

func (n *IndexExpr) String() string {
	var buf bytes.Buffer
	buf.WriteString(n.Operand.String())
	buf.WriteString("[")
	buf.WriteString(n.Index.String())
	buf.WriteString("]")
	return buf.String()
}

func (n *SliceExpr) String() string {
	var buf bytes.Buffer
	buf.WriteString(n.Operand.String())
	buf.WriteString("[")
	buf.WriteString(n.From.String())
	buf.WriteString(":")
	buf.WriteString(n.To.String())
	buf.WriteString("]")
	return buf.String()
}

func (n *SliceFromExpr) String() string {
	var buf bytes.Buffer
	buf.WriteString(n.Operand.String())
	buf.WriteString("[")
	buf.WriteString(n.From.String())
	buf.WriteString(":]")
	return buf.String()
}

func (n *SliceToExpr) String() string {
	var buf bytes.Buffer
	buf.WriteString(n.Operand.String())
	buf.WriteString("[:")
	buf.WriteString(n.To.String())
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
