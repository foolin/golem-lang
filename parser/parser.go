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

type Parser struct {
	scn       *scanner.Scanner
	isBuiltIn func(string) bool
	cur       *ast.Token
	next      *ast.Token
	synthetic int
}

func NewParser(scn *scanner.Scanner, isBuiltIn func(string) bool) *Parser {
	return &Parser{scn, isBuiltIn, nil, nil, 0}
}

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
	block := &ast.Block{nil, stmts, nil}
	return &ast.FnExpr{nil, params, block, 0, 0, nil}, err
}

func (p *Parser) imports() []ast.Statement {

	stmts := []ast.Statement{}

	for {
		if p.cur.Kind != ast.IMPORT {
			break
		}

		imp := &ast.Import{
			p.expect(ast.IMPORT),
			&ast.IdentExpr{p.expect(ast.IDENT), nil}}
		p.expectStatementDelimiter()
		stmts = append(stmts, imp)
	}

	return stmts
}

// Parse a sequence of statements or expressions.
func (p *Parser) statements(endKind ast.TokenKind) []ast.Statement {

	stmts := []ast.Statement{}

	for {
		if p.cur.Kind == endKind {
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
			if p.cur.Kind == e {
				return stmts
			}
		}

		stmts = append(stmts, p.statement())
	}
}

// Parse a statement, or return nil if there is no statement
// waiting to be parsed.
func (p *Parser) statement() ast.Statement {

	switch p.cur.Kind {

	case ast.CONST:
		return p.constStmt()

	case ast.LET:
		return p.letStmt()

	case ast.FN:
		if p.next.Kind == ast.IDENT {
			// named function
			return p.namedFn()
		} else {
			// anonymous function
			expr := p.fnExpr(p.consume())
			p.expectStatementDelimiter()
			return &ast.ExprStmt{expr}
		}

	case ast.IF:
		return p.ifStmt()

	case ast.WHILE:
		return p.whileStmt()

	case ast.FOR:
		return p.forStmt()

	case ast.SWITCH:
		return p.switchStmt()

	case ast.BREAK:
		return p.breakStmt()

	case ast.CONTINUE:
		return p.continueStmt()

	case ast.RETURN:
		return p.returnStmt()

	case ast.THROW:
		return p.throwStmt()

	case ast.TRY:
		return p.tryStmt()

	case ast.GO:
		return p.goStmt()

	default:
		// we couldn't find a statement to parse, so parse an expression instead
		expr := p.expression()
		p.expectStatementDelimiter()
		return &ast.ExprStmt{expr}
	}
}

func (p *Parser) namedFn() *ast.NamedFn {
	token := p.expect(ast.FN)
	result := &ast.NamedFn{
		token,
		&ast.IdentExpr{p.expect(ast.IDENT), nil},
		p.fnExpr(token)}
	p.expectStatementDelimiter()
	return result
}

func (p *Parser) constStmt() *ast.Const {

	token := p.expect(ast.CONST)

	decls := []*ast.Decl{p.decl()}
	for {
		switch {
		case p.cur.Kind == ast.COMMA:
			p.consume()
			decls = append(decls, p.decl())
		case isStatementDelimiter(p.cur.Kind):
			p.expectStatementDelimiter()
			return &ast.Const{token, decls}
		default:
			panic(p.unexpected())
		}
	}
}

func (p *Parser) letStmt() *ast.Let {

	token := p.expect(ast.LET)

	decls := []*ast.Decl{p.decl()}
	for {
		switch {
		case p.cur.Kind == ast.COMMA:
			p.consume()
			decls = append(decls, p.decl())
		case isStatementDelimiter(p.cur.Kind):
			p.expectStatementDelimiter()
			return &ast.Let{token, decls}
		default:
			panic(p.unexpected())
		}
	}
}

func (p *Parser) decl() *ast.Decl {

	ident := &ast.IdentExpr{p.expect(ast.IDENT), nil}
	if p.accept(ast.EQ) {
		return &ast.Decl{ident, p.expression()}
	} else {
		return &ast.Decl{ident, nil}
	}
}

func (p *Parser) ifStmt() *ast.If {

	token := p.expect(ast.IF)
	cond := p.expression()
	then := p.block()

	if p.accept(ast.ELSE) {

		switch p.cur.Kind {

		case ast.LBRACE:
			result := &ast.If{token, cond, then, p.block()}
			p.expectStatementDelimiter()
			return result

		case ast.IF:
			result := &ast.If{token, cond, then, p.ifStmt()}
			p.expectStatementDelimiter()
			return result

		default:
			panic(p.unexpected())
		}

	} else {
		p.expectStatementDelimiter()
		return &ast.If{token, cond, then, nil}
	}
}

func (p *Parser) whileStmt() *ast.While {

	result := &ast.While{p.expect(ast.WHILE), p.expression(), p.block()}
	p.expectStatementDelimiter()
	return result
}

func (p *Parser) forStmt() *ast.For {

	token := p.expect(ast.FOR)

	// parse identifers -- either single ident, or 'tuple' of idents
	var idents []*ast.IdentExpr
	switch p.cur.Kind {

	case ast.IDENT:
		idents = []*ast.IdentExpr{p.identExpr()}

	case ast.LPAREN:
		idents = p.tupleIdents()

	default:
		panic(p.unexpected())
	}

	// parse 'in'
	tok := p.expect(ast.IN)

	// make synthetic Identifier for iterable
	iblIdent := p.makeSyntheticIdent(tok.Position)

	// parse the rest
	iterable := p.expression()
	body := p.block()

	// done
	p.expectStatementDelimiter()
	return &ast.For{token, idents, iblIdent, iterable, body}
}

func (p *Parser) tupleIdents() []*ast.IdentExpr {

	lparen := p.expect(ast.LPAREN)

	idents := []*ast.IdentExpr{}

	switch p.cur.Kind {

	case ast.IDENT:
		idents = append(idents, p.identExpr())
	loop:
		for {
			switch p.cur.Kind {

			case ast.COMMA:
				p.consume()
				idents = append(idents, p.identExpr())

			case ast.RPAREN:
				p.consume()
				break loop

			default:
				panic(p.unexpected())
			}
		}

	case ast.RPAREN:
		p.consume()

	default:
		panic(p.unexpected())
	}

	if len(idents) < 2 {
		panic(&parserError{INVALID_FOR, lparen})
	}

	return idents
}

func (p *Parser) switchStmt() *ast.Switch {

	token := p.expect(ast.SWITCH)

	var item ast.Expression = nil
	if p.cur.Kind != ast.LBRACE {
		item = p.expression()
	}
	lbrace := p.expect(ast.LBRACE)

	// cases
	cases := []*ast.Case{p.caseStmt()}
	for p.cur.Kind == ast.CASE {
		cases = append(cases, p.caseStmt())
	}

	// default
	var def *ast.Default = nil
	if p.cur.Kind == ast.DEFAULT {
		def = p.defaultStmt()
	}

	// done
	result := &ast.Switch{token, item, lbrace, cases, def, p.expect(ast.RBRACE)}
	p.expectStatementDelimiter()
	return result
}

func (p *Parser) caseStmt() *ast.Case {

	token := p.expect(ast.CASE)

	matches := []ast.Expression{p.expression()}
	for {
		switch p.cur.Kind {

		case ast.COMMA:
			p.expect(ast.COMMA)
			matches = append(matches, p.expression())

		case ast.COLON:
			colon := p.expect(ast.COLON)
			body := p.statementsAny(ast.CASE, ast.DEFAULT, ast.RBRACE)
			if len(body) == 0 {
				panic(&parserError{INVALID_SWITCH, colon})
			}
			return &ast.Case{token, matches, body}

		default:
			panic(p.unexpected())
		}
	}
}

func (p *Parser) defaultStmt() *ast.Default {

	token := p.expect(ast.DEFAULT)
	colon := p.expect(ast.COLON)

	body := p.statements(ast.RBRACE)
	if len(body) == 0 {
		panic(&parserError{INVALID_SWITCH, colon})
	}

	return &ast.Default{token, body}
}

func (p *Parser) breakStmt() *ast.Break {
	result := &ast.Break{
		p.expect(ast.BREAK)}
	p.expectStatementDelimiter()
	return result
}

func (p *Parser) continueStmt() *ast.Continue {
	result := &ast.Continue{
		p.expect(ast.CONTINUE)}
	p.expectStatementDelimiter()
	return result
}

func (p *Parser) returnStmt() *ast.Return {

	token := p.expect(ast.RETURN)

	if isStatementDelimiter(p.cur.Kind) {
		p.expectStatementDelimiter()
		return &ast.Return{token, nil}
	} else {
		val := p.expression()
		p.expectStatementDelimiter()
		return &ast.Return{token, val}
	}
}

func (p *Parser) throwStmt() *ast.Throw {

	result := &ast.Throw{
		p.expect(ast.THROW),
		p.expression()}
	p.expectStatementDelimiter()
	return result
}

func (p *Parser) tryStmt() *ast.Try {

	tryToken := p.expect(ast.TRY)
	tryBlock := p.block()

	// catch
	var catchToken *ast.Token = nil
	var catchIdent *ast.IdentExpr = nil
	var catchBlock *ast.Block = nil

	if p.cur.Kind == ast.CATCH {
		catchToken = p.expect(ast.CATCH)
		catchIdent = p.identExpr()
		catchBlock = p.block()
	}

	// finally
	var finallyToken *ast.Token = nil
	var finallyBlock *ast.Block = nil

	if p.cur.Kind == ast.FINALLY {
		finallyToken = p.expect(ast.FINALLY)
		finallyBlock = p.block()
	}

	// make sure we got at least one of try or catch
	if catchToken == nil && finallyToken == nil {
		panic(&parserError{INVALID_TRY, tryToken})
	}

	// done
	p.expectStatementDelimiter()
	return &ast.Try{
		tryToken, tryBlock,
		catchToken, catchIdent, catchBlock,
		finallyToken, finallyBlock}
}

func (p *Parser) goStmt() *ast.Go {

	token := p.expect(ast.GO)

	prm := p.primary()
	if p.cur.Kind != ast.LPAREN {
		panic(p.unexpected())
	}
	lparen, actual, rparen := p.actualParams()
	invocation := &ast.InvokeExpr{prm, lparen, actual, rparen}

	p.expectStatementDelimiter()
	return &ast.Go{token, invocation}
}

// parse a sequence of stmts that are wrapped in curly braces
func (p *Parser) block() *ast.Block {

	lbrace := p.expect(ast.LBRACE)
	stmts := p.statements(ast.RBRACE)
	rbrace := p.expect(ast.RBRACE)
	return &ast.Block{lbrace, stmts, rbrace}
}

func (p *Parser) expression() ast.Expression {

	exp := p.ternaryExpr()

	if asn, ok := exp.(ast.Assignable); ok {

		if p.cur.Kind == ast.EQ {

			// assignment
			eq := p.expect(ast.EQ)
			exp = &ast.AssignmentExpr{asn, eq, p.expression()}

		} else if isAssignOp(p.cur.Kind) {

			// assignment operation
			op := p.consume()
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

	if p.cur.Kind == ast.HOOK {

		p.consume()
		then := p.expression()
		p.expect(ast.COLON)
		_else := p.ternaryExpr()
		return &ast.TernaryExpr{lhs, then, _else}

	} else {
		return lhs
	}
}

func (p *Parser) orExpr() ast.Expression {

	lhs := p.andExpr()
	for p.cur.Kind == ast.DBL_PIPE {
		tok := p.cur
		p.consume()
		lhs = &ast.BinaryExpr{lhs, tok, p.andExpr()}
	}
	return lhs
}

func (p *Parser) andExpr() ast.Expression {

	lhs := p.comparativeExpr()
	for p.cur.Kind == ast.DBL_AMP {
		tok := p.cur
		p.consume()
		lhs = &ast.BinaryExpr{lhs, tok, p.comparativeExpr()}
	}
	return lhs
}

func (p *Parser) comparativeExpr() ast.Expression {

	lhs := p.additiveExpr()
	for isComparative(p.cur.Kind) {
		tok := p.cur
		p.consume()
		lhs = &ast.BinaryExpr{lhs, tok, p.additiveExpr()}
	}
	return lhs
}

func (p *Parser) additiveExpr() ast.Expression {

	lhs := p.multiplicativeExpr()
	for isAdditive(p.cur.Kind) {
		tok := p.cur
		p.consume()
		lhs = &ast.BinaryExpr{lhs, tok, p.multiplicativeExpr()}
	}
	return lhs
}

func (p *Parser) multiplicativeExpr() ast.Expression {

	lhs := p.unaryExpr()
	for isMultiplicative(p.cur.Kind) {
		tok := p.cur
		p.consume()
		lhs = &ast.BinaryExpr{lhs, tok, p.unaryExpr()}
	}
	return lhs
}

func (p *Parser) unaryExpr() ast.Expression {

	if isUnary(p.cur.Kind) {
		tok := p.cur
		p.consume()
		return &ast.UnaryExpr{tok, p.unaryExpr()}

	} else {
		return p.postfixExpr()
	}
}

func (p *Parser) postfixExpr() ast.Expression {

	exp := p.primaryExpr()

	for isPostfix(p.cur.Kind) {

		if asn, ok := exp.(ast.Assignable); ok {
			tok := p.cur
			p.consume()
			exp = &ast.PostfixExpr{asn, tok}
		} else {
			panic(&parserError{INVALID_POSTFIX, p.cur})
		}
	}

	return exp
}

func (p *Parser) primaryExpr() ast.Expression {
	prm := p.primary()

	for {
		// look for suffixes: Invoke, Select, Index, Slice
		switch p.cur.Kind {

		case ast.LPAREN:
			lparen, actual, rparen := p.actualParams()
			prm = &ast.InvokeExpr{prm, lparen, actual, rparen}

		case ast.LBRACKET:
			lbracket := p.consume()

			switch p.cur.Kind {
			case ast.COLON:
				p.consume()
				prm = &ast.SliceToExpr{
					prm,
					lbracket,
					p.expression(),
					p.expect(ast.RBRACKET)}

			default:
				exp := p.expression()

				switch p.cur.Kind {
				case ast.RBRACKET:
					prm = &ast.IndexExpr{
						prm,
						lbracket,
						exp,
						p.expect(ast.RBRACKET)}

				case ast.COLON:
					p.consume()

					switch p.cur.Kind {
					case ast.RBRACKET:
						prm = &ast.SliceFromExpr{
							prm,
							lbracket,
							exp,
							p.expect(ast.RBRACKET)}
					default:
						prm = &ast.SliceExpr{
							prm,
							lbracket,
							exp,
							p.expression(),
							p.expect(ast.RBRACKET)}
					}

				default:
					panic(p.unexpected())
				}
			}

		case ast.DOT:
			p.expect(ast.DOT)
			prm = &ast.FieldExpr{prm, p.expect(ast.IDENT)}

		default:
			return prm
		}
	}
}

func (p *Parser) primary() ast.Expression {

	switch {

	case p.cur.Kind == ast.LPAREN:
		lparen := p.consume()
		expr := p.expression()

		switch p.cur.Kind {
		case ast.RPAREN:
			p.expect(ast.RPAREN)
			return expr

		case ast.COMMA:
			p.expect(ast.COMMA)
			return p.tupleExpr(lparen, expr)

		default:
			panic(p.unexpected())
		}

	case p.cur.Kind == ast.IDENT:

		switch {
		case p.isBuiltIn(p.cur.Text):
			return &ast.BuiltinExpr{p.consume()}

		case p.next.Kind == ast.EQ_GT:
			return p.lambdaOne()
		default:
			return p.identExpr()
		}

	case p.cur.Kind == ast.THIS:
		return &ast.ThisExpr{p.consume(), nil}

	case p.cur.Kind == ast.FN:
		return p.fnExpr(p.consume())

	case p.cur.Kind == ast.PIPE:
		return p.lambda()

	case p.cur.Kind == ast.DBL_PIPE:
		return p.lambdaZero()

	case p.cur.Kind == ast.STRUCT:
		return p.structExpr()

	case p.cur.Kind == ast.DICT:
		return p.dictExpr()

	case p.cur.Kind == ast.SET:
		return p.setExpr()

	case p.cur.Kind == ast.LBRACKET:
		return p.listExpr()

	default:
		return p.basicExpr()
	}
}

func (p *Parser) identExpr() *ast.IdentExpr {
	tok := p.cur
	p.expect(ast.IDENT)
	return &ast.IdentExpr{tok, nil}
}

func (p *Parser) fnExpr(token *ast.Token) *ast.FnExpr {

	p.expect(ast.LPAREN)
	if p.accept(ast.RPAREN) {
		return &ast.FnExpr{token, nil, p.block(), 0, 0, nil}
	} else {
		params := []*ast.FormalParam{}

		for {

			switch p.cur.Kind {
			case ast.CONST:
				p.consume()
				params = append(params, &ast.FormalParam{p.identExpr(), true})
			case ast.IDENT:
				params = append(params, &ast.FormalParam{p.identExpr(), false})
			default:
				panic(p.unexpected())
			}

			switch p.cur.Kind {
			case ast.COMMA:
				p.consume()
			case ast.RPAREN:
				p.consume()
				return &ast.FnExpr{token, params, p.block(), 0, 0, nil}
			default:
				panic(p.unexpected())
			}
		}
	}
}

func (p *Parser) lambdaZero() *ast.FnExpr {

	token := p.expect(ast.DBL_PIPE)

	p.expect(ast.EQ_GT)
	params := []*ast.FormalParam{}
	expr := &ast.ExprStmt{p.expression()}
	block := &ast.Block{nil, []ast.Statement{expr}, nil}
	return &ast.FnExpr{token, params, block, 0, 0, nil}
}

func (p *Parser) lambdaOne() *ast.FnExpr {
	token := p.expect(ast.IDENT)
	p.expect(ast.EQ_GT)
	params := []*ast.FormalParam{{&ast.IdentExpr{token, nil}, false}}
	expr := &ast.ExprStmt{p.expression()}
	block := &ast.Block{nil, []ast.Statement{expr}, nil}
	return &ast.FnExpr{token, params, block, 0, 0, nil}
}

func (p *Parser) lambda() *ast.FnExpr {

	token := p.expect(ast.PIPE)

	params := []*ast.FormalParam{}
	switch p.cur.Kind {

	case ast.IDENT:
		params = append(params, &ast.FormalParam{p.identExpr(), false})
	loop:
		for {
			switch p.cur.Kind {

			case ast.COMMA:
				p.consume()
				params = append(params, &ast.FormalParam{p.identExpr(), false})

			case ast.PIPE:
				p.consume()
				break loop

			default:
				panic(p.unexpected())
			}
		}

	case ast.PIPE:
		p.consume()

	default:
		panic(p.unexpected())
	}

	p.expect(ast.EQ_GT)

	expr := &ast.ExprStmt{p.expression()}
	block := &ast.Block{nil, []ast.Statement{expr}, nil}
	return &ast.FnExpr{token, params, block, 0, 0, nil}
}

func (p *Parser) structExpr() ast.Expression {
	return p.structBody(p.expect(ast.STRUCT))
}

//func (p *Parser) propExpr() ast.Expression {
//	token := p.expect(ast.PROP)
//	return p.structBody(token)
//}

func (p *Parser) structBody(token *ast.Token) ast.Expression {

	// key-value pairs
	keys := []*ast.Token{}
	values := []ast.Expression{}
	var rbrace *ast.Token
	lbrace := p.expect(ast.LBRACE)

	switch p.cur.Kind {

	case ast.IDENT:
		keys = append(keys, p.cur)
		p.consume()
		p.expect(ast.COLON)
		values = append(values, p.expression())
	loop:
		for {
			switch p.cur.Kind {

			case ast.COMMA:
				p.consume()
				keys = append(keys, p.cur)
				p.consume()
				p.expect(ast.COLON)
				values = append(values, p.expression())

			case ast.RBRACE:
				rbrace = p.consume()
				break loop

			default:
				panic(p.unexpected())
			}
		}

	case ast.RBRACE:
		rbrace = p.consume()

	default:
		panic(p.unexpected())
	}

	// done
	return &ast.StructExpr{token, lbrace, keys, values, rbrace, -1}
}

func (p *Parser) dictExpr() ast.Expression {

	dictToken := p.expect(ast.DICT)

	entries := []*ast.DictEntryExpr{}
	var rbrace *ast.Token

	lbrace := p.expect(ast.LBRACE)

	switch p.cur.Kind {

	case ast.RBRACE:
		rbrace = p.consume()

	default:
		key := p.expression()
		p.expect(ast.COLON)
		value := p.expression()
		entries = append(entries, &ast.DictEntryExpr{key, value})

	loop:
		for {
			switch p.cur.Kind {

			case ast.COMMA:
				p.consume()

				key = p.expression()
				p.expect(ast.COLON)
				value = p.expression()
				entries = append(entries, &ast.DictEntryExpr{key, value})

			case ast.RBRACE:
				rbrace = p.consume()
				break loop

			default:
				panic(p.unexpected())
			}
		}
	}

	return &ast.DictExpr{dictToken, lbrace, entries, rbrace}
}

func (p *Parser) setExpr() ast.Expression {

	setToken := p.expect(ast.SET)
	lbrace := p.expect(ast.LBRACE)

	if p.cur.Kind == ast.RBRACE {
		return &ast.SetExpr{setToken, lbrace, []ast.Expression{}, p.consume()}
	} else {

		elems := []ast.Expression{p.expression()}
		for {
			switch p.cur.Kind {
			case ast.RBRACE:
				return &ast.SetExpr{setToken, lbrace, elems, p.consume()}
			case ast.COMMA:
				p.consume()
				elems = append(elems, p.expression())
			default:
				panic(p.unexpected())
			}
		}
	}
}

func (p *Parser) listExpr() ast.Expression {

	lbracket := p.expect(ast.LBRACKET)

	if p.cur.Kind == ast.RBRACKET {
		return &ast.ListExpr{lbracket, []ast.Expression{}, p.consume()}
	} else {

		elems := []ast.Expression{p.expression()}
		for {
			switch p.cur.Kind {
			case ast.RBRACKET:
				return &ast.ListExpr{lbracket, elems, p.consume()}
			case ast.COMMA:
				p.consume()
				elems = append(elems, p.expression())
			default:
				panic(p.unexpected())
			}
		}
	}
}

func (p *Parser) tupleExpr(lparen *ast.Token, expr ast.Expression) ast.Expression {

	elems := []ast.Expression{expr, p.expression()}

	for {
		switch p.cur.Kind {
		case ast.RPAREN:
			return &ast.TupleExpr{lparen, elems, p.consume()}
		case ast.COMMA:
			p.consume()
			elems = append(elems, p.expression())
		default:
			panic(p.unexpected())
		}
	}
}

func (p *Parser) basicExpr() ast.Expression {

	tok := p.cur

	switch {

	case tok.IsBasic():
		p.consume()
		return &ast.BasicExpr{tok}

	default:
		panic(p.unexpected())
	}
}

func (p *Parser) actualParams() (*ast.Token, []ast.Expression, *ast.Token) {

	lparen := p.expect(ast.LPAREN)

	params := []ast.Expression{}
	switch p.cur.Kind {

	case ast.RPAREN:
		return lparen, params, p.consume()

	default:
		params = append(params, p.expression())
		for {
			switch p.cur.Kind {

			case ast.COMMA:
				p.consume()
				params = append(params, p.expression())

			case ast.RPAREN:
				return lparen, params, p.consume()

			default:
				panic(p.unexpected())
			}

		}
	}
}

// consume the current token if it has the given kind
func (p *Parser) accept(kind ast.TokenKind) bool {
	if p.cur.Kind == kind {
		p.consume()
		return true
	} else {
		return false
	}
}

// consume the current token if it has the given kind, else panic
func (p *Parser) expect(kind ast.TokenKind) *ast.Token {
	if p.cur.Kind == kind {
		result := p.cur
		p.consume()
		return result
	} else {
		panic(p.unexpected())
	}
}

func (p *Parser) expectStatementDelimiter() {
	if isStatementDelimiter(p.cur.Kind) {
		p.consume()
	} else {
		panic(p.unexpected())
	}
}

// consume the current token
func (p *Parser) consume() *ast.Token {
	result := p.cur
	p.cur, p.next = p.next, p.advance()
	return result
}

func (p *Parser) advance() *ast.Token {

	tok := p.scn.Next()

	// skip over line_feed
	for tok.Kind == ast.LINE_FEED {
		tok = p.scn.Next()
	}

	// look for errors from the scanner
	if tok.IsBad() {
		switch tok.Kind {

		case ast.UNEXPECTED_CHAR:
			panic(&parserError{UNEXPECTED_CHAR, tok})

		case ast.UNEXPECTED_EOF:
			panic(&parserError{UNEXPECTED_EOF, tok})

		default:
			panic("unreachable")
		}
	}

	// done
	return tok
}

// create a error that we will panic with
func (p *Parser) unexpected() error {
	switch p.cur.Kind {
	case ast.EOF:
		return &parserError{UNEXPECTED_EOF, p.cur}

	case ast.RESERVED:
		return &parserError{UNEXPECTED_RESERVED_WORD, p.cur}

	default:
		return &parserError{UNEXPECTED_TOKEN, p.cur}
	}
}

// make a synthetic identifier
func (p *Parser) makeSyntheticIdent(pos ast.Pos) *ast.IdentExpr {
	sym := fmt.Sprintf("#synthetic%d", p.synthetic)
	p.synthetic++
	return &ast.IdentExpr{
		&ast.Token{ast.IDENT, sym, pos}, nil}
}

func isStatementDelimiter(kind ast.TokenKind) bool {

	switch kind {
	case
		ast.SEMICOLON:

		return true
	default:
		return false
	}
}

func isComparative(kind ast.TokenKind) bool {
	switch kind {
	case
		ast.DBL_EQ,
		ast.NOT_EQ,
		ast.GT,
		ast.GT_EQ,
		ast.LT,
		ast.LT_EQ,
		ast.CMP,
		ast.HAS:

		return true
	default:
		return false
	}
}

func isAdditive(kind ast.TokenKind) bool {
	switch kind {
	case
		ast.PLUS,
		ast.MINUS,
		ast.PIPE,
		ast.CARET:

		return true
	default:
		return false
	}
}

func isMultiplicative(kind ast.TokenKind) bool {
	switch kind {
	case
		ast.STAR,
		ast.SLASH,
		ast.PERCENT,
		ast.AMP,
		ast.DBL_LT,
		ast.DBL_GT:

		return true
	default:
		return false
	}
}

func isUnary(kind ast.TokenKind) bool {

	switch kind {
	case
		ast.MINUS,
		ast.NOT,
		ast.TILDE:

		return true
	default:
		return false
	}
}

func isPostfix(kind ast.TokenKind) bool {

	switch kind {
	case
		ast.DBL_PLUS,
		ast.DBL_MINUS:

		return true
	default:
		return false
	}
}

func isAssignOp(kind ast.TokenKind) bool {
	switch kind {
	case
		ast.PLUS_EQ,
		ast.MINUS_EQ,
		ast.STAR_EQ,
		ast.SLASH_EQ,
		ast.PERCENT_EQ,
		ast.CARET_EQ,
		ast.AMP_EQ,
		ast.PIPE_EQ,
		ast.DBL_LT_EQ,
		ast.DBL_GT_EQ:

		return true
	default:
		return false
	}
}

func fromAssignOp(t *ast.Token) *ast.Token {

	switch t.Kind {
	case ast.PLUS_EQ:
		return &ast.Token{ast.PLUS, "+", t.Position}
	case ast.MINUS_EQ:
		return &ast.Token{ast.MINUS, "-", t.Position}
	case ast.STAR_EQ:
		return &ast.Token{ast.STAR, "*", t.Position}
	case ast.SLASH_EQ:
		return &ast.Token{ast.SLASH, "/", t.Position}
	case ast.PERCENT_EQ:
		return &ast.Token{ast.PERCENT, "%", t.Position}
	case ast.CARET_EQ:
		return &ast.Token{ast.CARET, "^", t.Position}
	case ast.AMP_EQ:
		return &ast.Token{ast.AMP, "&", t.Position}
	case ast.PIPE_EQ:
		return &ast.Token{ast.PIPE, "|", t.Position}
	case ast.DBL_LT_EQ:
		return &ast.Token{ast.DBL_LT, "<<", t.Position}
	case ast.DBL_GT_EQ:
		return &ast.Token{ast.DBL_GT, ">>", t.Position}

	default:
		panic("invalid op")
	}
}

//--------------------------------------------------------------
// parserError

type parserErrorKind int

const (
	UNEXPECTED_CHAR parserErrorKind = iota
	UNEXPECTED_TOKEN
	UNEXPECTED_RESERVED_WORD
	UNEXPECTED_EOF
	INVALID_POSTFIX
	INVALID_FOR
	INVALID_SWITCH
	INVALID_TRY
)

type parserError struct {
	kind  parserErrorKind
	token *ast.Token
}

func (e *parserError) Error() string {

	switch e.kind {

	case UNEXPECTED_CHAR:
		return fmt.Sprintf("Unexpected Character '%v' at %v", e.token.Text, e.token.Position)

	case UNEXPECTED_TOKEN:
		return fmt.Sprintf("Unexpected Token '%v' at %v", e.token.Text, e.token.Position)

	case UNEXPECTED_RESERVED_WORD:
		return fmt.Sprintf("Unexpected Reserved Word '%v' at %v", e.token.Text, e.token.Position)

	case UNEXPECTED_EOF:
		return fmt.Sprintf("Unexpected EOF at %v", e.token.Position)

	case INVALID_POSTFIX:
		return fmt.Sprintf("Invalid Postfix Expression at %v", e.token.Position)

	case INVALID_FOR:
		return fmt.Sprintf("Invalid For Expression at %v", e.token.Position)

	case INVALID_SWITCH:
		return fmt.Sprintf("Invalid Switch Expression at %v", e.token.Position)

	case INVALID_TRY:
		return fmt.Sprintf("Invalid TRY Expression at %v", e.token.Position)

	default:
		panic("unreachable")
	}
}
