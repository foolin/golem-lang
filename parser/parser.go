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

//--------------------------------------------------------------
// Parser

// Parser parses Golem source code, and creates an Abstract Syntax Tree.
type Parser struct {
	scn       *scanner.Scanner
	isBuiltIn func(string) bool
	cur       tokenInfo
	next      tokenInfo
	synthetic int
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
func (p *Parser) ParseModule() (fn *ast.FnExpr, err error) {

	// In a recursive descent parser, errors can be generated deep
	// in the call stack.  We are going to use panic-recover to handle them.
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(runtime.Error); ok {
				panic(r)
			}
			fn = nil
			err = r.(error)
		}
	}()

	// read the first two tokens
	p.cur = p.advance()
	p.next = p.advance()

	// parse imports
	stmts := p.imports()

	// parse the module
	stmts = append(stmts, p.statements(ast.EOF)...)
	p.expect(ast.EOF)

	params := []*ast.FormalParam{}
	block := &ast.BlockNode{nil, stmts, nil}
	return &ast.FnExpr{nil, params, block, 0, 0, nil}, err
}

func (p *Parser) imports() []ast.Statement {

	stmts := []ast.Statement{}

	for {
		if p.cur.token.Kind != ast.Import {
			break
		}

		imp := &ast.ImportStmt{
			p.expect(ast.Import),
			&ast.IdentExpr{p.expect(ast.Ident), nil}}
		p.expectStatementDelimiter()
		stmts = append(stmts, imp)
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
		return &ast.ExprStmt{expr}

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
		return &ast.ExprStmt{expr}
	}
}

func (p *Parser) namedFn() *ast.NamedFnStmt {
	token := p.expect(ast.Fn)
	result := &ast.NamedFnStmt{
		token,
		&ast.IdentExpr{p.expect(ast.Ident), nil},
		p.fnExpr(token)}
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
			return &ast.ConstStmt{token, decls}
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
			return &ast.LetStmt{token, decls}
		default:
			panic(p.unexpected())
		}
	}
}

func (p *Parser) decl() *ast.DeclNode {

	ident := &ast.IdentExpr{p.expect(ast.Ident), nil}
	if p.accept(ast.Eq) {
		return &ast.DeclNode{ident, p.expression()}
	}
	return &ast.DeclNode{ident, nil}
}

func (p *Parser) ifStmt() *ast.IfStmt {

	token := p.expect(ast.If)
	cond := p.expression()
	then := p.block()

	if p.accept(ast.Else) {

		switch p.cur.token.Kind {

		case ast.Lbrace:
			result := &ast.IfStmt{token, cond, then, p.block()}
			p.expectStatementDelimiter()
			return result

		case ast.If:
			result := &ast.IfStmt{token, cond, then, p.ifStmt()}
			p.expectStatementDelimiter()
			return result

		default:
			panic(p.unexpected())
		}

	} else {
		p.expectStatementDelimiter()
		return &ast.IfStmt{token, cond, then, nil}
	}
}

func (p *Parser) whileStmt() *ast.WhileStmt {

	result := &ast.WhileStmt{p.expect(ast.While), p.expression(), p.block()}
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

	// make synthetic Identifier for iterable
	iblIdent := p.makeSyntheticIdent(tok.Position)

	// parse the rest
	iterable := p.expression()
	body := p.block()

	// done
	p.expectStatementDelimiter()
	return &ast.ForStmt{token, idents, iblIdent, iterable, body}
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
		panic(&parserError{InvalidFor, lparen})
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
	result := &ast.SwitchStmt{token, item, lbrace, cases, def, p.expect(ast.Rbrace)}
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
				panic(&parserError{InvalidSwitch, colon})
			}
			return &ast.CaseNode{token, matches, body}

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
		panic(&parserError{InvalidSwitch, colon})
	}

	return &ast.DefaultNode{token, body}
}

func (p *Parser) breakStmt() *ast.BreakStmt {
	result := &ast.BreakStmt{
		p.expect(ast.Break)}
	p.expectStatementDelimiter()
	return result
}

func (p *Parser) continueStmt() *ast.ContinueStmt {
	result := &ast.ContinueStmt{
		p.expect(ast.Continue)}
	p.expectStatementDelimiter()
	return result
}

func (p *Parser) returnStmt() *ast.ReturnStmt {

	token := p.expect(ast.Return)

	val := p.expression()
	p.expectStatementDelimiter()
	return &ast.ReturnStmt{token, val}
}

func (p *Parser) throwStmt() *ast.ThrowStmt {

	result := &ast.ThrowStmt{
		p.expect(ast.Throw),
		p.expression()}
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
		panic(&parserError{InvalidTry, tryToken})
	}

	// done
	p.expectStatementDelimiter()
	return &ast.TryStmt{
		tryToken, tryBlock,
		catchToken, catchIdent, catchBlock,
		finallyToken, finallyBlock}
}

func (p *Parser) goStmt() *ast.GoStmt {

	token := p.expect(ast.Go)

	prm := p.primary()
	if p.cur.token.Kind != ast.Lparen {
		panic(p.unexpected())
	}
	lparen, actual, rparen := p.actualParams()
	invocation := &ast.InvokeExpr{prm, lparen, actual, rparen}

	p.expectStatementDelimiter()
	return &ast.GoStmt{token, invocation}
}

// parse a sequence of stmts that are wrapped in curly braces
func (p *Parser) block() *ast.BlockNode {

	lbrace := p.expect(ast.Lbrace)
	stmts := p.statements(ast.Rbrace)
	rbrace := p.expect(ast.Rbrace)
	return &ast.BlockNode{lbrace, stmts, rbrace}
}

func (p *Parser) expression() ast.Expression {

	exp := p.ternaryExpr()

	if asn, ok := exp.(ast.Assignable); ok {

		if p.cur.token.Kind == ast.Eq {

			// assignment
			eq := p.expect(ast.Eq)
			exp = &ast.AssignmentExpr{asn, eq, p.expression()}

		} else if isAssignOp(p.cur.token.Kind) {

			// assignment operation
			op := p.consume().token
			exp = &ast.AssignmentExpr{
				asn,
				op,
				&ast.BinaryExpr{
					asn,
					fromAssignOp(op),
					p.expression()}}
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
		return &ast.TernaryExpr{lhs, then, _else}

	}
	return lhs
}

func (p *Parser) orExpr() ast.Expression {

	lhs := p.andExpr()
	for p.cur.token.Kind == ast.DblPipe {
		tok := p.cur.token
		p.consume()
		lhs = &ast.BinaryExpr{lhs, tok, p.andExpr()}
	}
	return lhs
}

func (p *Parser) andExpr() ast.Expression {

	lhs := p.comparativeExpr()
	for p.cur.token.Kind == ast.DblAmp {
		tok := p.cur.token
		p.consume()
		lhs = &ast.BinaryExpr{lhs, tok, p.comparativeExpr()}
	}
	return lhs
}

func (p *Parser) comparativeExpr() ast.Expression {

	lhs := p.additiveExpr()
	for isComparative(p.cur.token.Kind) {
		tok := p.cur.token
		p.consume()
		lhs = &ast.BinaryExpr{lhs, tok, p.additiveExpr()}
	}
	return lhs
}

func (p *Parser) additiveExpr() ast.Expression {

	lhs := p.multiplicativeExpr()
	for isAdditive(p.cur.token.Kind) {
		tok := p.cur.token
		p.consume()
		lhs = &ast.BinaryExpr{lhs, tok, p.multiplicativeExpr()}
	}
	return lhs
}

func (p *Parser) multiplicativeExpr() ast.Expression {

	lhs := p.unaryExpr()
	for isMultiplicative(p.cur.token.Kind) {
		tok := p.cur.token
		p.consume()
		lhs = &ast.BinaryExpr{lhs, tok, p.unaryExpr()}
	}
	return lhs
}

func (p *Parser) unaryExpr() ast.Expression {

	if isUnary(p.cur.token.Kind) {
		tok := p.cur.token
		p.consume()
		return &ast.UnaryExpr{tok, p.unaryExpr()}

	}
	return p.postfixExpr()
}

func (p *Parser) postfixExpr() ast.Expression {

	exp := p.primaryExpr()

	for isPostfix(p.cur.token.Kind) {

		if asn, ok := exp.(ast.Assignable); ok {
			tok := p.cur.token
			p.consume()
			exp = &ast.PostfixExpr{asn, tok}
		} else {
			panic(&parserError{InvalidPostfix, p.cur.token})
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
			prm = &ast.InvokeExpr{prm, lparen, actual, rparen}

		case ast.Lbracket:
			lbracket := p.consume().token

			switch p.cur.token.Kind {
			case ast.Colon:
				p.consume()
				prm = &ast.SliceToExpr{
					prm,
					lbracket,
					p.expression(),
					p.expect(ast.Rbracket)}

			default:
				exp := p.expression()

				switch p.cur.token.Kind {
				case ast.Rbracket:
					prm = &ast.IndexExpr{
						prm,
						lbracket,
						exp,
						p.expect(ast.Rbracket)}

				case ast.Colon:
					p.consume()

					switch p.cur.token.Kind {
					case ast.Rbracket:
						prm = &ast.SliceFromExpr{
							prm,
							lbracket,
							exp,
							p.expect(ast.Rbracket)}
					default:
						prm = &ast.SliceExpr{
							prm,
							lbracket,
							exp,
							p.expression(),
							p.expect(ast.Rbracket)}
					}

				default:
					panic(p.unexpected())
				}
			}

		case ast.Dot:
			p.expect(ast.Dot)
			prm = &ast.FieldExpr{prm, p.expect(ast.Ident)}

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
			return &ast.BuiltinExpr{p.consume().token}

		case p.next.token.Kind == ast.EqGt:
			return p.lambdaOne()
		default:
			return p.identExpr()
		}

	case p.cur.token.Kind == ast.This:
		return &ast.ThisExpr{p.consume().token, nil}

	case p.cur.token.Kind == ast.Fn:
		return p.fnExpr(p.consume().token)

	case p.cur.token.Kind == ast.Pipe:
		return p.lambda()

	case p.cur.token.Kind == ast.DblPipe:
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
	return &ast.IdentExpr{tok, nil}
}

func (p *Parser) fnExpr(token *ast.Token) *ast.FnExpr {

	p.expect(ast.Lparen)
	if p.accept(ast.Rparen) {
		return &ast.FnExpr{token, nil, p.block(), 0, 0, nil}
	}
	params := []*ast.FormalParam{}

	for {

		switch p.cur.token.Kind {
		case ast.Const:
			p.consume()
			params = append(params, &ast.FormalParam{p.identExpr(), true})
		case ast.Ident:
			params = append(params, &ast.FormalParam{p.identExpr(), false})
		default:
			panic(p.unexpected())
		}

		switch p.cur.token.Kind {
		case ast.Comma:
			p.consume()
		case ast.Rparen:
			p.consume()
			return &ast.FnExpr{token, params, p.block(), 0, 0, nil}
		default:
			panic(p.unexpected())
		}
	}

}

func (p *Parser) lambdaZero() *ast.FnExpr {

	token := p.expect(ast.DblPipe)

	p.expect(ast.EqGt)
	params := []*ast.FormalParam{}
	expr := &ast.ExprStmt{p.expression()}
	block := &ast.BlockNode{nil, []ast.Statement{expr}, nil}
	return &ast.FnExpr{token, params, block, 0, 0, nil}
}

func (p *Parser) lambdaOne() *ast.FnExpr {
	token := p.expect(ast.Ident)
	p.expect(ast.EqGt)
	params := []*ast.FormalParam{{&ast.IdentExpr{token, nil}, false}}
	expr := &ast.ExprStmt{p.expression()}
	block := &ast.BlockNode{nil, []ast.Statement{expr}, nil}
	return &ast.FnExpr{token, params, block, 0, 0, nil}
}

func (p *Parser) lambda() *ast.FnExpr {

	token := p.expect(ast.Pipe)

	params := []*ast.FormalParam{}
	switch p.cur.token.Kind {

	case ast.Ident:
		params = append(params, &ast.FormalParam{p.identExpr(), false})
	loop:
		for {
			switch p.cur.token.Kind {

			case ast.Comma:
				p.consume()
				params = append(params, &ast.FormalParam{p.identExpr(), false})

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

	expr := &ast.ExprStmt{p.expression()}
	block := &ast.BlockNode{nil, []ast.Statement{expr}, nil}
	return &ast.FnExpr{token, params, block, 0, 0, nil}
}

func (p *Parser) structExpr() ast.Expression {
	return p.structBody(p.expect(ast.Struct))
}

//func (p *Parser) propExpr() ast.Expression {
//	token := p.expect(ast.Prop)
//	return p.structBody(token)
//}

func (p *Parser) structBody(token *ast.Token) ast.Expression {

	// key-value pairs
	keys := []*ast.Token{}
	values := []ast.Expression{}
	var rbrace *ast.Token
	lbrace := p.expect(ast.Lbrace)

	switch p.cur.token.Kind {

	case ast.Ident:
		keys = append(keys, p.cur.token)
		p.consume()
		p.expect(ast.Colon)
		values = append(values, p.expression())
	loop:
		for {
			switch p.cur.token.Kind {

			case ast.Comma:
				p.consume()
				keys = append(keys, p.cur.token)
				p.consume()
				p.expect(ast.Colon)
				values = append(values, p.expression())

			case ast.Rbrace:
				rbrace = p.consume().token
				break loop

			default:
				panic(p.unexpected())
			}
		}

	case ast.Rbrace:
		rbrace = p.consume().token

	default:
		panic(p.unexpected())
	}

	// done
	return &ast.StructExpr{token, lbrace, keys, values, rbrace, -1}
}

func (p *Parser) dictExpr() ast.Expression {

	dictToken := p.expect(ast.Dict)

	entries := []*ast.DictEntryExpr{}
	var rbrace *ast.Token

	lbrace := p.expect(ast.Lbrace)

	switch p.cur.token.Kind {

	case ast.Rbrace:
		rbrace = p.consume().token

	default:
		key := p.expression()
		p.expect(ast.Colon)
		value := p.expression()
		entries = append(entries, &ast.DictEntryExpr{key, value})

	loop:
		for {
			switch p.cur.token.Kind {

			case ast.Comma:
				p.consume()

				key = p.expression()
				p.expect(ast.Colon)
				value = p.expression()
				entries = append(entries, &ast.DictEntryExpr{key, value})

			case ast.Rbrace:
				rbrace = p.consume().token
				break loop

			default:
				panic(p.unexpected())
			}
		}
	}

	return &ast.DictExpr{dictToken, lbrace, entries, rbrace}
}

func (p *Parser) setExpr() ast.Expression {

	setToken := p.expect(ast.Set)
	lbrace := p.expect(ast.Lbrace)

	if p.cur.token.Kind == ast.Rbrace {
		return &ast.SetExpr{setToken, lbrace, []ast.Expression{}, p.consume().token}
	}

	elems := []ast.Expression{p.expression()}
	for {
		switch p.cur.token.Kind {
		case ast.Rbrace:
			return &ast.SetExpr{setToken, lbrace, elems, p.consume().token}
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
		return &ast.ListExpr{lbracket, []ast.Expression{}, p.consume().token}
	}

	elems := []ast.Expression{p.expression()}
	for {
		switch p.cur.token.Kind {
		case ast.Rbracket:
			return &ast.ListExpr{lbracket, elems, p.consume().token}
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
			return &ast.TupleExpr{lparen, elems, p.consume().token}
		case ast.Comma:
			p.consume()
			elems = append(elems, p.expression())
		default:
			panic(p.unexpected())
		}
	}
}

func (p *Parser) basicExpr() ast.Expression {

	tok := p.cur.token

	switch {

	case tok.IsBasic():
		p.consume()
		return &ast.BasicExpr{tok}

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
			panic(&parserError{UnexpectedChar, token})

		case ast.UnexpectedEOF:
			panic(&parserError{UnexpectedEOF, token})

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
		return &parserError{UnexpectedEOF, p.cur.token}

	case ast.Reserved:
		return &parserError{UnexpectedReservedWork, p.cur.token}

	default:
		return &parserError{UnexpectedToken, p.cur.token}
	}
}

// make a synthetic identifier
func (p *Parser) makeSyntheticIdent(pos ast.Pos) *ast.IdentExpr {
	sym := fmt.Sprintf("#synthetic%d", p.synthetic)
	p.synthetic++
	return &ast.IdentExpr{
		&ast.Token{ast.Ident, sym, pos}, nil}
}

func isComparative(kind ast.TokenKind) bool {
	switch kind {
	case
		ast.DblEq,
		ast.NotEq,
		ast.Gt,
		ast.GtEq,
		ast.Lt,
		ast.LtEq,
		ast.Cmp,
		ast.Has:

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
		ast.DblLt,
		ast.DblGt:

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
		ast.DblPlus,
		ast.DblMinus:

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
		ast.DblLtEq,
		ast.DblGtEq:

		return true
	default:
		return false
	}
}

func fromAssignOp(t *ast.Token) *ast.Token {

	switch t.Kind {
	case ast.PlusEq:
		return &ast.Token{ast.Plus, "+", t.Position}
	case ast.MinusEq:
		return &ast.Token{ast.Minus, "-", t.Position}
	case ast.StarEq:
		return &ast.Token{ast.Star, "*", t.Position}
	case ast.SlashEq:
		return &ast.Token{ast.Slash, "/", t.Position}
	case ast.PercentEq:
		return &ast.Token{ast.Percent, "%", t.Position}
	case ast.CaretEq:
		return &ast.Token{ast.Caret, "^", t.Position}
	case ast.AmpEq:
		return &ast.Token{ast.Amp, "&", t.Position}
	case ast.PipeEq:
		return &ast.Token{ast.Pipe, "|", t.Position}
	case ast.DblLtEq:
		return &ast.Token{ast.DblLt, "<<", t.Position}
	case ast.DblGtEq:
		return &ast.Token{ast.DblGt, ">>", t.Position}

	default:
		panic("invalid op")
	}
}

//--------------------------------------------------------------
// parserError

type parserErrorKind int

// Parser Errors
const (
	UnexpectedChar parserErrorKind = iota
	UnexpectedToken
	UnexpectedReservedWork
	UnexpectedEOF
	InvalidPostfix
	InvalidFor
	InvalidSwitch
	InvalidTry
)

type parserError struct {
	kind  parserErrorKind
	token *ast.Token
}

func (e *parserError) Error() string {

	switch e.kind {

	case UnexpectedChar:
		return fmt.Sprintf("Unexpected Character '%v' at %v", e.token.Text, e.token.Position)

	case UnexpectedToken:
		return fmt.Sprintf("Unexpected Token '%v' at %v", e.token.Text, e.token.Position)

	case UnexpectedReservedWork:
		return fmt.Sprintf("Unexpected Reserved Word '%v' at %v", e.token.Text, e.token.Position)

	case UnexpectedEOF:
		return fmt.Sprintf("Unexpected EOF at %v", e.token.Position)

	case InvalidPostfix:
		return fmt.Sprintf("Invalid Postfix Expression at %v", e.token.Position)

	case InvalidFor:
		return fmt.Sprintf("Invalid ForStmt Expression at %v", e.token.Position)

	case InvalidSwitch:
		return fmt.Sprintf("Invalid SwitchStmt Expression at %v", e.token.Position)

	case InvalidTry:
		return fmt.Sprintf("Invalid Try Expression at %v", e.token.Position)

	default:
		panic("unreachable")
	}
}
