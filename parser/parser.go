// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package parser

import (
	"fmt"
	"runtime"

	"github.com/mjarmy/golem-lang/ast"
	"github.com/mjarmy/golem-lang/scanner"
)

//---------------------------------------------------------------
// The Golem Parser
//---------------------------------------------------------------

// Parser parses Golem source code, and creates an Abstract Syntax Tree.
type Parser struct {
	scn           *scanner.Scanner
	isBuiltIn     func(string) bool
	cur           tokenInfo
	next          tokenInfo
	iterIDCounter int
}

type tokenInfo struct {
	token  *ast.Token
	skipLF bool // whether or not we skipped and linefeeds while advancing to this token
}

// NewParser creates a new Parser
func NewParser(scn *scanner.Scanner, isBuiltIn func(string) bool) *Parser {
	return &Parser{scn, isBuiltIn, tokenInfo{}, tokenInfo{}, 0}
}

// ParseModule parses a Golem module
func (p *Parser) ParseModule() (mod *ast.Module, err error) {

	// In a recursive descent parser, errors can be generated deep
	// in the call stack.  We are going to use panic-recover to handle them.
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(runtime.Error); ok {
				panic(r)
			}
			mod = nil
			err = r.(error)
		}
	}()

	// read the first two tokens
	p.cur = p.advance()
	p.next = p.advance()

	// parse imports
	stmts := p.imports()

	// parse the rest of the statements
	stmts = append(stmts, p.statements(ast.EOF)...)
	p.expect(ast.EOF)

	// create initialization function
	initFunc := &ast.FnExpr{
		Token:    nil,
		Required: []*ast.Param{},
		Body: &ast.BlockNode{
			LBrace:     nil,
			Statements: stmts,
			RBrace:     nil,
			Scope:      ast.NewScope(),
		},
		Scope: ast.NewFuncScope(),
	}

	// done
	mod = &ast.Module{
		Name:     p.scn.Source.Name,
		Path:     p.scn.Source.Path,
		InitFunc: initFunc,
	}
	return
}

func (p *Parser) imports() []ast.Statement {

	stmts := []ast.Statement{}

	for {
		if p.cur.token.Kind != ast.Import {
			break
		}

		tok := p.expect(ast.Import)

		idents := []*ast.IdentExpr{}
		idents = append(idents, &ast.IdentExpr{
			Symbol:   p.expect(ast.Ident),
			Variable: nil,
		})

	loop:
		for {
			switch {
			case p.cur.token.Kind == ast.Comma:
				p.consume()
				idents = append(idents, &ast.IdentExpr{
					Symbol:   p.expect(ast.Ident),
					Variable: nil,
				})
			case p.atStatementDelimiter():
				break loop
			default:
				panic(p.unexpected())
			}
		}

		p.expectStatementDelimiter()
		stmts = append(stmts, &ast.ImportStmt{
			Token:  tok,
			Idents: idents,
		})
	}

	return stmts
}

// Parse a sequence of statements or expressions.
func (p *Parser) statements(endKind ast.TokenKind) []ast.Statement {

	stmts := []ast.Statement{}

	for {
		if p.cur.token.Kind == endKind {
			return stmts
		}

		stmts = append(stmts, p.statement())
	}

}

// Parse a sequence of statements or expressions, ending with any of the provided tokens.
func (p *Parser) statementsAny(endKinds ...ast.TokenKind) []ast.Statement {

	stmts := []ast.Statement{}

	for {
		for _, e := range endKinds {
			if p.cur.token.Kind == e {
				return stmts
			}
		}

		stmts = append(stmts, p.statement())
	}
}

// Parse a statement, or return nil if there is no statement
// waiting to be parsed.
func (p *Parser) statement() ast.Statement {

	switch p.cur.token.Kind {

	case ast.Const:
		return p.constStmt()

	case ast.Let:
		return p.letStmt()

	case ast.Fn:
		if p.next.token.Kind == ast.Ident {
			// named function
			return p.namedFn()
		}
		// anonymous function
		expr := p.fnExpr(p.consume().token)
		p.expectStatementDelimiter()
		return &ast.ExprStmt{Expr: expr}

	case ast.If:
		return p.ifStmt()

	case ast.While:
		return p.whileStmt()

	case ast.For:
		return p.forStmt()

	case ast.Switch:
		return p.switchStmt()

	case ast.Break:
		return p.breakStmt()

	case ast.Continue:
		return p.continueStmt()

	case ast.Return:
		return p.returnStmt()

	case ast.Throw:
		return p.throwStmt()

	case ast.Try:
		return p.tryStmt()

	case ast.Go:
		return p.goStmt()

	default:
		// we couldn't find a statement to parse, so parse an expression instead
		expr := p.expression()
		p.expectStatementDelimiter()
		return &ast.ExprStmt{Expr: expr}
	}
}

func (p *Parser) namedFn() *ast.NamedFnStmt {
	token := p.expect(ast.Fn)
	result := &ast.NamedFnStmt{
		Token: token,
		Ident: &ast.IdentExpr{
			Symbol:   p.expect(ast.Ident),
			Variable: nil,
		},
		Func: p.fnExpr(token),
	}
	p.expectStatementDelimiter()
	return result
}

func (p *Parser) constStmt() *ast.ConstStmt {

	token := p.expect(ast.Const)
	decls := []*ast.DeclNode{p.decl()}

	for {
		switch {
		case p.cur.token.Kind == ast.Comma:
			p.consume()
			decls = append(decls, p.decl())
		case p.atStatementDelimiter():
			p.expectStatementDelimiter()
			return &ast.ConstStmt{
				Token: token,
				Decls: decls,
			}
		default:
			panic(p.unexpected())
		}
	}
}

func (p *Parser) letStmt() *ast.LetStmt {

	token := p.expect(ast.Let)
	decls := []*ast.DeclNode{p.decl()}

	for {
		switch {
		case p.cur.token.Kind == ast.Comma:
			p.consume()
			decls = append(decls, p.decl())
		case p.atStatementDelimiter():
			p.expectStatementDelimiter()
			return &ast.LetStmt{
				Token: token,
				Decls: decls,
			}
		default:
			panic(p.unexpected())
		}
	}
}

func (p *Parser) decl() *ast.DeclNode {

	ident := &ast.IdentExpr{
		Symbol:   p.expect(ast.Ident),
		Variable: nil,
	}
	if p.accept(ast.Eq) {
		return &ast.DeclNode{
			Ident: ident,
			Val:   p.expression(),
		}
	}
	return &ast.DeclNode{
		Ident: ident,
		Val:   nil,
	}
}

func (p *Parser) ifStmt() *ast.IfStmt {

	token := p.expect(ast.If)
	cond := p.expression()
	then := p.block()

	if p.accept(ast.Else) {

		switch p.cur.token.Kind {

		case ast.Lbrace:
			result := &ast.IfStmt{
				Token: token,
				Cond:  cond,
				Then:  then,
				Else:  p.block(),
			}
			p.expectStatementDelimiter()
			return result

		case ast.If:
			result := &ast.IfStmt{
				Token: token,
				Cond:  cond,
				Then:  then,
				Else:  p.ifStmt(),
			}
			p.expectStatementDelimiter()
			return result

		default:
			panic(p.unexpected())
		}

	} else {
		p.expectStatementDelimiter()
		return &ast.IfStmt{
			Token: token,
			Cond:  cond,
			Then:  then,
			Else:  nil,
		}
	}
}

func (p *Parser) whileStmt() *ast.WhileStmt {

	result := &ast.WhileStmt{
		Token: p.expect(ast.While),
		Cond:  p.expression(),
		Body:  p.block(),
	}
	p.expectStatementDelimiter()
	return result
}

func (p *Parser) forStmt() *ast.ForStmt {

	token := p.expect(ast.For)

	// parse identifers -- either single ident, or 'tuple' of idents
	var idents []*ast.IdentExpr
	switch p.cur.token.Kind {

	case ast.Ident:
		idents = []*ast.IdentExpr{p.identExpr()}

	case ast.Lparen:
		idents = p.tupleIdents()

	default:
		panic(p.unexpected())
	}

	// parse 'in'
	tok := p.expect(ast.In)

	// make identifier for iterable
	iblIdent := p.makeIterIdent(tok.Position)

	// parse the rest
	iterable := p.expression()
	body := p.block()

	// done
	p.expectStatementDelimiter()
	return &ast.ForStmt{
		Token:         token,
		Idents:        idents,
		IterableIdent: iblIdent,
		Iterable:      iterable,
		Body:          body,
		Scope:         ast.NewScope(),
	}
}

// make an identifier for an iterable in a 'for' stmt
func (p *Parser) makeIterIdent(pos ast.Pos) *ast.IdentExpr {
	sym := fmt.Sprintf("#iter%d", p.iterIDCounter)
	p.iterIDCounter++
	return &ast.IdentExpr{
		Symbol:   &ast.Token{Kind: ast.Ident, Text: sym, Position: pos},
		Variable: nil,
	}
}

func (p *Parser) tupleIdents() []*ast.IdentExpr {

	lparen := p.expect(ast.Lparen)

	idents := []*ast.IdentExpr{}

	switch p.cur.token.Kind {

	case ast.Ident:
		idents = append(idents, p.identExpr())
	loop:
		for {
			switch p.cur.token.Kind {

			case ast.Comma:
				p.consume()
				idents = append(idents, p.identExpr())

			case ast.Rparen:
				p.consume()
				break loop

			default:
				panic(p.unexpected())
			}
		}

	case ast.Rparen:
		p.consume()

	default:
		panic(p.unexpected())
	}

	if len(idents) < 2 {
		panic(newParserError(p.scn.Source.Path, invalidFor, lparen))
	}

	return idents
}

func (p *Parser) switchStmt() *ast.SwitchStmt {

	token := p.expect(ast.Switch)

	var item ast.Expression
	if p.cur.token.Kind != ast.Lbrace {
		item = p.expression()
	}
	lbrace := p.expect(ast.Lbrace)

	// cases
	cases := []*ast.CaseNode{p.caseStmt()}
	for p.cur.token.Kind == ast.Case {
		cases = append(cases, p.caseStmt())
	}

	// default
	var def *ast.DefaultNode
	if p.cur.token.Kind == ast.Default {
		def = p.defaultStmt()
	}

	// done
	result := &ast.SwitchStmt{
		Token:       token,
		Item:        item,
		LBrace:      lbrace,
		Cases:       cases,
		DefaultNode: def,
		RBrace:      p.expect(ast.Rbrace),
	}
	p.expectStatementDelimiter()
	return result
}

func (p *Parser) caseStmt() *ast.CaseNode {

	token := p.expect(ast.Case)

	matches := []ast.Expression{p.expression()}
	for {
		switch p.cur.token.Kind {

		case ast.Comma:
			p.expect(ast.Comma)
			matches = append(matches, p.expression())

		case ast.Colon:
			colon := p.expect(ast.Colon)
			body := p.statementsAny(ast.Case, ast.Default, ast.Rbrace)
			if len(body) == 0 {
				panic(newParserError(p.scn.Source.Path, invalidSwitch, colon))
			}
			return &ast.CaseNode{
				Token:   token,
				Matches: matches,
				Body:    body,
			}

		default:
			panic(p.unexpected())
		}
	}
}

func (p *Parser) defaultStmt() *ast.DefaultNode {

	token := p.expect(ast.Default)
	colon := p.expect(ast.Colon)

	body := p.statements(ast.Rbrace)
	if len(body) == 0 {
		panic(newParserError(p.scn.Source.Path, invalidSwitch, colon))
	}

	return &ast.DefaultNode{
		Token: token,
		Body:  body,
	}
}

func (p *Parser) breakStmt() *ast.BreakStmt {
	result := &ast.BreakStmt{
		Token: p.expect(ast.Break),
	}
	p.expectStatementDelimiter()
	return result
}

func (p *Parser) continueStmt() *ast.ContinueStmt {
	result := &ast.ContinueStmt{
		Token: p.expect(ast.Continue),
	}
	p.expectStatementDelimiter()
	return result
}

func (p *Parser) returnStmt() *ast.ReturnStmt {

	token := p.expect(ast.Return)

	val := p.expression()
	p.expectStatementDelimiter()
	return &ast.ReturnStmt{
		Token: token,
		Val:   val,
	}
}

func (p *Parser) throwStmt() *ast.ThrowStmt {

	result := &ast.ThrowStmt{
		Token: p.expect(ast.Throw),
		Val:   p.expression(),
	}
	p.expectStatementDelimiter()
	return result
}

func (p *Parser) tryStmt() *ast.TryStmt {

	tryToken := p.expect(ast.Try)
	tryBlock := p.block()

	// catch
	var catchToken *ast.Token
	var catchIdent *ast.IdentExpr
	var catchBlock *ast.BlockNode

	if p.cur.token.Kind == ast.Catch {
		catchToken = p.expect(ast.Catch)
		catchIdent = p.identExpr()
		catchBlock = p.block()
	}

	// finally
	var finallyToken *ast.Token
	var finallyBlock *ast.BlockNode

	if p.cur.token.Kind == ast.Finally {
		finallyToken = p.expect(ast.Finally)
		finallyBlock = p.block()
	}

	// make sure we got at least one of try or catch
	if catchToken == nil && finallyToken == nil {
		panic(newParserError(p.scn.Source.Path, invalidTry, tryToken))
	}

	// done
	p.expectStatementDelimiter()
	return &ast.TryStmt{
		TryToken:     tryToken,
		TryBlock:     tryBlock,
		CatchToken:   catchToken,
		CatchIdent:   catchIdent,
		CatchBlock:   catchBlock,
		FinallyToken: finallyToken,
		FinallyBlock: finallyBlock,
		CatchScope:   ast.NewScope(),
	}
}

func (p *Parser) goStmt() *ast.GoStmt {

	token := p.expect(ast.Go)

	prm := p.primary()
	if p.cur.token.Kind != ast.Lparen {
		panic(p.unexpected())
	}
	lparen, actual, rparen := p.actualParams()
	invocation := &ast.InvokeExpr{
		Operand: prm,
		LParen:  lparen,
		Params:  actual,
		RParen:  rparen,
	}

	p.expectStatementDelimiter()
	return &ast.GoStmt{
		Token:      token,
		Invocation: invocation,
	}
}

// parse a sequence of stmts that are wrapped in curly braces
func (p *Parser) block() *ast.BlockNode {

	lbrace := p.expect(ast.Lbrace)
	stmts := p.statements(ast.Rbrace)
	rbrace := p.expect(ast.Rbrace)
	return &ast.BlockNode{
		LBrace:     lbrace,
		Statements: stmts,
		RBrace:     rbrace,
		Scope:      ast.NewScope(),
	}
}

func (p *Parser) expression() ast.Expression {

	exp := p.ternaryExpr()

	if asn, ok := exp.(ast.Assignable); ok {

		if p.cur.token.Kind == ast.Eq {

			// assignment
			eq := p.expect(ast.Eq)
			exp = &ast.AssignmentExpr{
				Assignee: asn,
				Eq:       eq,
				Val:      p.expression(),
			}

		} else if isAssignOp(p.cur.token.Kind) {

			// assignment operation
			op := p.consume().token
			exp = &ast.AssignmentExpr{
				Assignee: asn,
				Eq:       op,
				Val: &ast.BinaryExpr{
					LHS: asn,
					Op:  fromAssignOp(op),
					RHS: p.expression(),
				},
			}
		}
	}

	return exp
}

func (p *Parser) ternaryExpr() ast.Expression {

	lhs := p.orExpr()

	if p.cur.token.Kind == ast.Hook {

		p.consume()
		then := p.expression()
		p.expect(ast.Colon)
		_else := p.ternaryExpr()
		return &ast.TernaryExpr{
			Cond: lhs,
			Then: then,
			Else: _else,
		}

	}
	return lhs
}

func (p *Parser) orExpr() ast.Expression {

	lhs := p.andExpr()
	for p.cur.token.Kind == ast.DoublePipe {
		tok := p.cur.token
		p.consume()
		lhs = &ast.BinaryExpr{
			LHS: lhs,
			Op:  tok,
			RHS: p.andExpr(),
		}
	}
	return lhs
}

func (p *Parser) andExpr() ast.Expression {

	lhs := p.comparativeExpr()
	for p.cur.token.Kind == ast.DoubleAmp {
		tok := p.cur.token
		p.consume()
		lhs = &ast.BinaryExpr{
			LHS: lhs,
			Op:  tok,
			RHS: p.comparativeExpr(),
		}
	}
	return lhs
}

func (p *Parser) comparativeExpr() ast.Expression {

	lhs := p.additiveExpr()
	for isComparative(p.cur.token.Kind) {
		tok := p.cur.token
		p.consume()
		lhs = &ast.BinaryExpr{
			LHS: lhs,
			Op:  tok,
			RHS: p.additiveExpr(),
		}
	}
	return lhs
}

func (p *Parser) additiveExpr() ast.Expression {

	lhs := p.multiplicativeExpr()
	for isAdditive(p.cur.token.Kind) {
		tok := p.cur.token
		p.consume()
		lhs = &ast.BinaryExpr{
			LHS: lhs,
			Op:  tok,
			RHS: p.multiplicativeExpr(),
		}
	}
	return lhs
}

func (p *Parser) multiplicativeExpr() ast.Expression {

	lhs := p.unaryExpr()
	for isMultiplicative(p.cur.token.Kind) {
		tok := p.cur.token
		p.consume()
		lhs = &ast.BinaryExpr{
			LHS: lhs,
			Op:  tok,
			RHS: p.unaryExpr(),
		}
	}
	return lhs
}

func (p *Parser) unaryExpr() ast.Expression {

	if isUnary(p.cur.token.Kind) {
		tok := p.cur.token
		p.consume()
		return &ast.UnaryExpr{
			Op:      tok,
			Operand: p.unaryExpr(),
		}

	}
	return p.postfixExpr()
}

func (p *Parser) postfixExpr() ast.Expression {

	exp := p.primaryExpr()

	for isPostfix(p.cur.token.Kind) {

		if asn, ok := exp.(ast.Assignable); ok {
			tok := p.cur.token
			p.consume()
			exp = &ast.PostfixExpr{
				Assignee: asn,
				Op:       tok,
			}
		} else {
			panic(newParserError(p.scn.Source.Path, invalidPostfix, p.cur.token))
		}
	}

	return exp
}

func (p *Parser) primaryExpr() ast.Expression {
	prm := p.primary()

	for {
		// look for suffixes: Invoke, Select, Index, Slice
		switch p.cur.token.Kind {

		case ast.Lparen:
			lparen, actual, rparen := p.actualParams()
			prm = &ast.InvokeExpr{
				Operand: prm,
				LParen:  lparen,
				Params:  actual,
				RParen:  rparen,
			}

		case ast.Lbracket:
			lbracket := p.consume().token

			switch p.cur.token.Kind {
			case ast.Colon:
				p.consume()
				prm = &ast.SliceToExpr{
					Operand:  prm,
					LBracket: lbracket,
					To:       p.expression(),
					RBracket: p.expect(ast.Rbracket),
				}

			default:
				exp := p.expression()

				switch p.cur.token.Kind {
				case ast.Rbracket:
					prm = &ast.IndexExpr{
						Operand:  prm,
						LBracket: lbracket,
						Index:    exp,
						RBracket: p.expect(ast.Rbracket),
					}

				case ast.Colon:
					p.consume()

					switch p.cur.token.Kind {
					case ast.Rbracket:
						prm = &ast.SliceFromExpr{
							Operand:  prm,
							LBracket: lbracket,
							From:     exp,
							RBracket: p.expect(ast.Rbracket),
						}
					default:
						prm = &ast.SliceExpr{
							Operand:  prm,
							LBracket: lbracket,
							From:     exp,
							To:       p.expression(),
							RBracket: p.expect(ast.Rbracket),
						}
					}

				default:
					panic(p.unexpected())
				}
			}

		case ast.Dot:
			p.expect(ast.Dot)
			prm = &ast.FieldExpr{
				Operand: prm,
				Key:     p.expect(ast.Ident),
			}

		default:
			return prm
		}
	}
}

func (p *Parser) primary() ast.Expression {

	switch {

	case p.cur.token.Kind == ast.Lparen:
		lparen := p.consume().token
		expr := p.expression()

		switch p.cur.token.Kind {
		case ast.Rparen:
			p.expect(ast.Rparen)
			return expr

		case ast.Comma:
			p.expect(ast.Comma)
			return p.tupleExpr(lparen, expr)

		default:
			panic(p.unexpected())
		}

	case p.cur.token.Kind == ast.Ident:

		switch {
		case p.isBuiltIn(p.cur.token.Text):
			return &ast.BuiltinExpr{
				Fn: p.consume().token,
			}

		default:
			return p.identExpr()
		}

	case p.cur.token.Kind == ast.This:
		return &ast.ThisExpr{
			Token:    p.consume().token,
			Variable: nil,
		}

	case p.cur.token.Kind == ast.Fn:
		return p.fnExpr(p.consume().token)

	case p.cur.token.Kind == ast.Pipe:
		return p.lambda()

	case p.cur.token.Kind == ast.DoublePipe:
		return p.lambdaZero()

	case p.cur.token.Kind == ast.Struct:
		return p.structExpr()

	case p.cur.token.Kind == ast.Dict:
		return p.dictExpr()

	case p.cur.token.Kind == ast.Set:
		return p.setExpr()

	case p.cur.token.Kind == ast.Lbracket:
		return p.listExpr()

	default:
		return p.basicExpr()
	}
}

func (p *Parser) identExpr() *ast.IdentExpr {
	tok := p.cur.token
	p.expect(ast.Ident)
	return &ast.IdentExpr{
		Symbol:   tok,
		Variable: nil,
	}
}

func (p *Parser) fnExpr(token *ast.Token) *ast.FnExpr {

	p.expect(ast.Lparen)
	if p.accept(ast.Rparen) {
		return &ast.FnExpr{
			Token:    token,
			Required: nil,
			Body:     p.block(),
			Scope:    ast.NewFuncScope(),
		}
	}

	params := []*ast.Param{}
	optional := []*ast.OptionalParam{}

	for {

		switch p.cur.token.Kind {

		case ast.Const:
			p.consume()

			ident := p.identExpr()
			if p.accept(ast.Eq) {
				optional = append(optional, &ast.OptionalParam{
					Ident:   ident,
					IsConst: true,
					Value:   p.basicExpr(),
				})
			} else {
				if len(optional) > 0 {
					panic(p.unexpected())
				}

				params = append(params, &ast.Param{
					Ident:   ident,
					IsConst: true,
				})
			}

		case ast.Ident:

			ident := p.identExpr()
			if p.accept(ast.Eq) {
				optional = append(optional, &ast.OptionalParam{
					Ident:   ident,
					IsConst: false,
					Value:   p.basicExpr(),
				})
			} else {
				if len(optional) > 0 {
					panic(p.unexpected())
				}

				params = append(params, &ast.Param{
					Ident:   ident,
					IsConst: false,
				})
			}

		default:
			panic(p.unexpected())
		}

		switch p.cur.token.Kind {

		case ast.Comma:
			p.consume()

		// Variadic Arity
		case ast.TripleDot:

			if len(optional) > 0 {
				panic(p.unexpected())
			}

			p.consume()
			p.expect(ast.Rparen)

			n := len(params) - 1
			return &ast.FnExpr{
				Token:    token,
				Required: params[:n],
				Variadic: params[n],
				Body:     p.block(),
				Scope:    ast.NewFuncScope(),
			}

		case ast.Rparen:
			p.consume()

			// Fixed or Multiple Arity
			if len(optional) == 0 {
				return &ast.FnExpr{
					Token:    token,
					Required: params,
					Body:     p.block(),
					Scope:    ast.NewFuncScope(),
				}
			}

			// Multiple Arity
			return &ast.FnExpr{
				Token:    token,
				Required: params,
				Optional: optional,
				Body:     p.block(),
				Scope:    ast.NewFuncScope(),
			}

		default:
			panic(p.unexpected())
		}
	}

}

func (p *Parser) lambdaZero() *ast.FnExpr {

	token := p.expect(ast.DoublePipe)

	p.expect(ast.EqGt)
	params := []*ast.Param{}
	expr := &ast.ExprStmt{Expr: p.expression()}
	block := &ast.BlockNode{
		LBrace:     nil,
		Statements: []ast.Statement{expr},
		RBrace:     nil,
		Scope:      ast.NewScope(),
	}
	return &ast.FnExpr{
		Token:    token,
		Required: params,
		Body:     block,
		Scope:    ast.NewFuncScope(),
	}
}

func (p *Parser) lambda() *ast.FnExpr {

	token := p.expect(ast.Pipe)

	params := []*ast.Param{}
	switch p.cur.token.Kind {

	case ast.Ident:
		params = append(params, &ast.Param{
			Ident:   p.identExpr(),
			IsConst: false,
		})
	loop:
		for {
			switch p.cur.token.Kind {

			case ast.Comma:
				p.consume()
				params = append(params, &ast.Param{
					Ident:   p.identExpr(),
					IsConst: false,
				})

			case ast.Pipe:
				p.consume()
				break loop

			default:
				panic(p.unexpected())
			}
		}

	case ast.Pipe:
		p.consume()

	default:
		panic(p.unexpected())
	}

	p.expect(ast.EqGt)

	expr := &ast.ExprStmt{Expr: p.expression()}
	block := &ast.BlockNode{
		LBrace:     nil,
		Statements: []ast.Statement{expr},
		RBrace:     nil,
		Scope:      ast.NewScope(),
	}
	return &ast.FnExpr{
		Token:    token,
		Required: params,
		Body:     block,
		Scope:    ast.NewFuncScope(),
	}
}

func (p *Parser) structExpr() ast.Expression {

	structToken := p.expect(ast.Struct)
	lbrace := p.expect(ast.Lbrace)

	if p.cur.token.Kind == ast.Rbrace {
		return &ast.StructExpr{
			StructToken: structToken,
			LBrace:      lbrace,
			Entries:     []*ast.StructEntry{},
			RBrace:      p.consume().token,
			Scope:       ast.NewStructScope(),
		}
	}

	entry := p.structEntry()
	entries := []*ast.StructEntry{entry}

	names := make(map[string]bool)
	name := entry.Key.Text
	names[name] = true

	for {
		switch p.cur.token.Kind {
		case ast.Rbrace:
			return &ast.StructExpr{
				StructToken: structToken,
				LBrace:      lbrace,
				Entries:     entries,
				RBrace:      p.consume().token,
				Scope:       ast.NewStructScope(),
			}
		case ast.Comma:
			p.consume()
			entry := p.structEntry()

			name := entry.Key.Text
			if _, ok := names[name]; ok {
				panic(newParserError(p.scn.Source.Path, duplicateKey, entry.Key))
			}
			names[name] = true

			entries = append(entries, entry)
		default:
			panic(p.unexpected())
		}
	}
}

func (p *Parser) structEntry() *ast.StructEntry {

	key := p.expect(ast.Ident)
	p.expect(ast.Colon)

	if p.cur.token.Kind == ast.Prop {
		return &ast.StructEntry{
			Key:   key,
			Value: p.property(),
		}
	}
	return &ast.StructEntry{
		Key:   key,
		Value: p.expression(),
	}

}

func (p *Parser) property() *ast.PropNode {

	token := p.expect(ast.Prop)
	lbrace := p.expect(ast.Lbrace)

	get := p.propertyFunc()
	if len(get.Required) != 0 {
		panic(newParserError(p.scn.Source.Path, invalidPropertyGetter, get.Token))
	}

	var set *ast.FnExpr
	if p.accept(ast.Comma) {
		set = p.propertyFunc()
		if len(set.Required) != 1 {
			panic(newParserError(p.scn.Source.Path, invalidPropertySetter, set.Token))
		}
	}

	return &ast.PropNode{
		Token:  token,
		LBrace: lbrace,
		Get:    get,
		Set:    set,
		RBrace: p.expect(ast.Rbrace),
	}
}

func (p *Parser) propertyFunc() *ast.FnExpr {

	switch p.cur.token.Kind {

	case ast.Fn:
		return p.fnExpr(p.consume().token)

	case ast.DoublePipe:
		return p.lambdaZero()

	case ast.Pipe:
		return p.lambda()

	default:
		panic(p.unexpected())
	}
}

func (p *Parser) dictExpr() ast.Expression {

	dictToken := p.expect(ast.Dict)
	lbrace := p.expect(ast.Lbrace)

	if p.cur.token.Kind == ast.Rbrace {
		return &ast.DictExpr{
			DictToken: dictToken,
			LBrace:    lbrace,
			Entries:   []*ast.DictEntry{},
			RBrace:    p.consume().token,
		}
	}

	entries := []*ast.DictEntry{p.dictEntry()}
	for {
		switch p.cur.token.Kind {
		case ast.Rbrace:
			return &ast.DictExpr{
				DictToken: dictToken,
				LBrace:    lbrace,
				Entries:   entries,
				RBrace:    p.consume().token,
			}
		case ast.Comma:
			p.consume()
			entries = append(entries, p.dictEntry())
		default:
			panic(p.unexpected())
		}
	}
}

func (p *Parser) dictEntry() *ast.DictEntry {
	key := p.expression()
	p.expect(ast.Colon)
	value := p.expression()
	return &ast.DictEntry{
		Key:   key,
		Value: value,
	}
}

func (p *Parser) setExpr() ast.Expression {

	setToken := p.expect(ast.Set)
	lbrace := p.expect(ast.Lbrace)

	if p.cur.token.Kind == ast.Rbrace {
		return &ast.SetExpr{
			SetToken: setToken,
			LBrace:   lbrace,
			Elems:    []ast.Expression{},
			RBrace:   p.consume().token,
		}
	}

	elems := []ast.Expression{p.expression()}
	for {
		switch p.cur.token.Kind {
		case ast.Rbrace:
			return &ast.SetExpr{
				SetToken: setToken,
				LBrace:   lbrace,
				Elems:    elems,
				RBrace:   p.consume().token,
			}
		case ast.Comma:
			p.consume()
			elems = append(elems, p.expression())
		default:
			panic(p.unexpected())
		}
	}
}

func (p *Parser) listExpr() ast.Expression {

	lbracket := p.expect(ast.Lbracket)

	if p.cur.token.Kind == ast.Rbracket {
		return &ast.ListExpr{
			LBracket: lbracket,
			Elems:    []ast.Expression{},
			RBracket: p.consume().token,
		}
	}

	elems := []ast.Expression{p.expression()}
	for {
		switch p.cur.token.Kind {
		case ast.Rbracket:
			return &ast.ListExpr{
				LBracket: lbracket,
				Elems:    elems,
				RBracket: p.consume().token,
			}
		case ast.Comma:
			p.consume()
			elems = append(elems, p.expression())
		default:
			panic(p.unexpected())
		}
	}
}

func (p *Parser) tupleExpr(lparen *ast.Token, expr ast.Expression) ast.Expression {

	elems := []ast.Expression{expr, p.expression()}

	for {
		switch p.cur.token.Kind {
		case ast.Rparen:

			// tuples always have at least 2 elements
			if len(elems) < 2 {
				panic("unreachable")
			}

			return &ast.TupleExpr{
				LParen: lparen,
				Elems:  elems,
				RParen: p.consume().token,
			}
		case ast.Comma:
			p.consume()
			elems = append(elems, p.expression())
		default:
			panic(p.unexpected())
		}
	}
}

func (p *Parser) basicExpr() *ast.BasicExpr {

	tok := p.cur.token

	switch {

	case tok.IsBasic():
		p.consume()
		return &ast.BasicExpr{
			Token: tok,
		}

	default:
		panic(p.unexpected())
	}
}

func (p *Parser) actualParams() (*ast.Token, []ast.Expression, *ast.Token) {

	lparen := p.expect(ast.Lparen)

	params := []ast.Expression{}
	switch p.cur.token.Kind {

	case ast.Rparen:
		return lparen, params, p.consume().token

	default:
		params = append(params, p.expression())
		for {
			switch p.cur.token.Kind {

			case ast.Comma:
				p.consume()
				params = append(params, p.expression())

			case ast.Rparen:
				return lparen, params, p.consume().token

			default:
				panic(p.unexpected())
			}

		}
	}
}

// consume the current token if it has the given kind
func (p *Parser) accept(kind ast.TokenKind) bool {
	if p.cur.token.Kind == kind {
		p.consume()
		return true
	}
	return false
}

// consume the current token if it has the given kind, else panic
func (p *Parser) expect(kind ast.TokenKind) *ast.Token {
	if p.cur.token.Kind == kind {
		result := p.cur.token
		p.consume()
		return result
	}
	panic(p.unexpected())
}

func (p *Parser) expectStatementDelimiter() {
	switch {
	case
		p.cur.token.Kind == ast.Semicolon,
		p.cur.token.Kind == ast.EOF:
		p.consume()
	case p.cur.skipLF:
		// nothing to do
		return
	default:
		panic(p.unexpected())
	}
}

func (p *Parser) atStatementDelimiter() bool {

	switch {
	case
		p.cur.token.Kind == ast.Semicolon,
		p.cur.token.Kind == ast.EOF:
		return true
	case p.cur.skipLF:
		return true
	default:
		return false
	}
}

// consume the current token
func (p *Parser) consume() tokenInfo {
	result := p.cur
	p.cur, p.next = p.next, p.advance()
	return result
}

func (p *Parser) advance() tokenInfo {

	token := p.scn.Next()
	skipLF := false

	// skip over line_feed
	for token.Kind == ast.LineFeed {
		skipLF = true
		token = p.scn.Next()
	}

	// look for errors from the scanner
	if token.IsBad() {
		switch token.Kind {

		case ast.UnexpectedChar:
			panic(newParserError(p.scn.Source.Path, unexpectedChar, token))

		case ast.UnexpectedEOF:
			panic(newParserError(p.scn.Source.Path, unexpectedEOF, token))

		default:
			panic("unreachable")
		}
	}

	// done
	return tokenInfo{token, skipLF}
}

// create a error that we will panic with
func (p *Parser) unexpected() error {
	switch p.cur.token.Kind {
	case ast.EOF:
		return newParserError(p.scn.Source.Path, unexpectedEOF, p.cur.token)

	case ast.Reserved:
		return newParserError(p.scn.Source.Path, unexpectedReservedWord, p.cur.token)

	default:
		return newParserError(p.scn.Source.Path, unexpectedToken, p.cur.token)
	}
}

func isComparative(kind ast.TokenKind) bool {
	switch kind {
	case
		ast.DoubleEq,
		ast.NotEq,
		ast.Gt,
		ast.GtEq,
		ast.Lt,
		ast.LtEq,
		ast.Cmp:

		return true
	default:
		return false
	}
}

func isAdditive(kind ast.TokenKind) bool {
	switch kind {
	case
		ast.Plus,
		ast.Minus,
		ast.Pipe,
		ast.Caret:

		return true
	default:
		return false
	}
}

func isMultiplicative(kind ast.TokenKind) bool {
	switch kind {
	case
		ast.Star,
		ast.Slash,
		ast.Percent,
		ast.Amp,
		ast.DoubleLt,
		ast.DoubleGt:

		return true
	default:
		return false
	}
}

func isUnary(kind ast.TokenKind) bool {

	switch kind {
	case
		ast.Minus,
		ast.Not,
		ast.Tilde:

		return true
	default:
		return false
	}
}

func isPostfix(kind ast.TokenKind) bool {

	switch kind {
	case
		ast.DoublePlus,
		ast.DoubleMinus:

		return true
	default:
		return false
	}
}

func isAssignOp(kind ast.TokenKind) bool {
	switch kind {
	case
		ast.PlusEq,
		ast.MinusEq,
		ast.StarEq,
		ast.SlashEq,
		ast.PercentEq,
		ast.CaretEq,
		ast.AmpEq,
		ast.PipeEq,
		ast.DoubleLtEq,
		ast.DoubleGtEq:

		return true
	default:
		return false
	}
}

func fromAssignOp(t *ast.Token) *ast.Token {

	switch t.Kind {
	case ast.PlusEq:
		return &ast.Token{Kind: ast.Plus, Text: "+", Position: t.Position}
	case ast.MinusEq:
		return &ast.Token{Kind: ast.Minus, Text: "-", Position: t.Position}
	case ast.StarEq:
		return &ast.Token{Kind: ast.Star, Text: "*", Position: t.Position}
	case ast.SlashEq:
		return &ast.Token{Kind: ast.Slash, Text: "/", Position: t.Position}
	case ast.PercentEq:
		return &ast.Token{Kind: ast.Percent, Text: "%", Position: t.Position}
	case ast.CaretEq:
		return &ast.Token{Kind: ast.Caret, Text: "^", Position: t.Position}
	case ast.AmpEq:
		return &ast.Token{Kind: ast.Amp, Text: "&", Position: t.Position}
	case ast.PipeEq:
		return &ast.Token{Kind: ast.Pipe, Text: "|", Position: t.Position}
	case ast.DoubleLtEq:
		return &ast.Token{Kind: ast.DoubleLt, Text: "<<", Position: t.Position}
	case ast.DoubleGtEq:
		return &ast.Token{Kind: ast.DoubleGt, Text: ">>", Position: t.Position}

	default:
		panic("invalid op")
	}
}
