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
	UnexpectedChar parserErrorKind = iota
	UnexpectedToken
	UnexpectedReservedWork
	UnexpectedEOF
	InvalidPostfix
	InvalidFor
	InvalidSwitch
	InvalidTry
	InvalidPropertyGetter
	InvalidPropertySetter
	DuplicateKey
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

	case UnexpectedChar:
		return fmt.Sprintf("Unexpected Character '%v' at %s:%v", e.token.Text, e.path, e.token.Position)

	case UnexpectedToken:
		return fmt.Sprintf("Unexpected Token '%v' at %s:%v", e.token.Text, e.path, e.token.Position)

	case UnexpectedReservedWork:
		return fmt.Sprintf("Unexpected Reserved Word '%v' at %s:%v", e.token.Text, e.path, e.token.Position)

	case UnexpectedEOF:
		return fmt.Sprintf("Unexpected EOF at %s:%v", e.path, e.token.Position)

	case InvalidPostfix:
		return fmt.Sprintf("Invalid Postfix Expression at %s:%v", e.path, e.token.Position)

	case InvalidFor:
		return fmt.Sprintf("Invalid ForStmt Expression at %s:%v", e.path, e.token.Position)

	case InvalidSwitch:
		return fmt.Sprintf("Invalid SwitchStmt Expression at %s:%v", e.path, e.token.Position)

	case InvalidTry:
		return fmt.Sprintf("Invalid Try Expression at %s:%v", e.path, e.token.Position)

	case InvalidPropertyGetter:
		return fmt.Sprintf("Invalid Property Getter at %s:%v", e.path, e.token.Position)

	case InvalidPropertySetter:
		return fmt.Sprintf("Invalid Property Setter at %s:%v", e.path, e.token.Position)

	case DuplicateKey:
		return fmt.Sprintf("Duplicate Key at %s:%v", e.path, e.token.Position)

	default:
		panic("unreachable")
	}
}
