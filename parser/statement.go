// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package parser

import (
	"fmt"

	"github.com/mjarmy/golem-lang/ast"
)

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
