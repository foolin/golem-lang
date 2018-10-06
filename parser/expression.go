// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package parser

import (
	"github.com/mjarmy/golem-lang/ast"
)

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
