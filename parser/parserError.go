// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package parser

import (
	"fmt"

	"github.com/mjarmy/golem-lang/ast"
)

type parserErrorKind int

const (
	unexpectedChar parserErrorKind = iota
	unexpectedToken
	unexpectedReservedWork
	unexpectedEOF
	invalidPostfix
	invalidFor
	invalidSwitch
	invalidTry
	invalidPropertyGetter
	invalidPropertySetter
	duplicateKey
)

type parserError struct {
	path  string
	kind  parserErrorKind
	token *ast.Token
}

func newParserError(path string, kind parserErrorKind, token *ast.Token) *parserError {
	return &parserError{path: path, kind: kind, token: token}
}

func (e *parserError) Error() string {

	switch e.kind {

	case unexpectedChar:
		return fmt.Sprintf("Unexpected Character '%v' at %s:%v", e.token.Text, e.path, e.token.Position)

	case unexpectedToken:
		return fmt.Sprintf("Unexpected Token '%v' at %s:%v", e.token.Text, e.path, e.token.Position)

	case unexpectedReservedWork:
		return fmt.Sprintf("Unexpected Reserved Word '%v' at %s:%v", e.token.Text, e.path, e.token.Position)

	case unexpectedEOF:
		return fmt.Sprintf("Unexpected EOF at %s:%v", e.path, e.token.Position)

	case invalidPostfix:
		return fmt.Sprintf("Invalid Postfix Expression at %s:%v", e.path, e.token.Position)

	case invalidFor:
		return fmt.Sprintf("Invalid ForStmt Expression at %s:%v", e.path, e.token.Position)

	case invalidSwitch:
		return fmt.Sprintf("Invalid SwitchStmt Expression at %s:%v", e.path, e.token.Position)

	case invalidTry:
		return fmt.Sprintf("Invalid Try Expression at %s:%v", e.path, e.token.Position)

	case invalidPropertyGetter:
		return fmt.Sprintf("Invalid Property Getter at %s:%v", e.path, e.token.Position)

	case invalidPropertySetter:
		return fmt.Sprintf("Invalid Property Setter at %s:%v", e.path, e.token.Position)

	case duplicateKey:
		return fmt.Sprintf("Duplicate Key at %s:%v", e.path, e.token.Position)

	default:
		panic("unreachable")
	}
}
