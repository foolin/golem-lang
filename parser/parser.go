// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package parser

import (
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
