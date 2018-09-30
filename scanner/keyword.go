// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package scanner

import (
	"github.com/mjarmy/golem-lang/ast"
	"unicode"
)

var keywords = map[string]ast.TokenKind{
	"_":        ast.BlankIdent,
	"break":    ast.Break,
	"case":     ast.Case,
	"catch":    ast.Catch,
	"const":    ast.Const,
	"continue": ast.Continue,
	"default":  ast.Default,
	"dict":     ast.Dict,
	"else":     ast.Else,
	"false":    ast.False,
	"finally":  ast.Finally,
	"fn":       ast.Fn,
	"for":      ast.For,
	"go":       ast.Go,
	"if":       ast.If,
	"import":   ast.Import,
	"in":       ast.In,
	"let":      ast.Let,
	"null":     ast.Null,
	"prop":     ast.Prop,
	"return":   ast.Return,
	"set":      ast.Set,
	"struct":   ast.Struct,
	"switch":   ast.Switch,
	"this":     ast.This,
	"throw":    ast.Throw,
	"true":     ast.True,
	"try":      ast.Try,
	"while":    ast.While,
}

// reserve a bunch of keywords just in case
var reservedWords = map[string]bool{
	"as":        true,
	"byte":      true,
	"defer":     true,
	"goto":      true,
	"like":      true,
	"module":    true,
	"native":    true,
	"package":   true,
	"priv":      true,
	"private":   true,
	"prot":      true,
	"protected": true,
	"pub":       true,
	"public":    true,
	"pure":      true,
	"rsync":     true,
	"rune":      true,
	"select":    true,
	"static":    true,
	"sync":      true,
	"with":      true,
	"yield":     true,
}

// IsKeyword returns whether a string is a keyword
func IsKeyword(text string) bool {
	_, ok := keywords[text]
	if ok {
		return true
	}
	_, ok = reservedWords[text]
	return ok
}

// IsIdentifier returns whether a string is an identifier
func IsIdentifier(text string) bool {

	if IsKeyword(text) {
		return false
	}

	for i, r := range text {
		if i == 0 {
			if !IsIdentStart(r) {
				return false
			}
		} else {
			if !IsIdentContinue(r) {
				return false
			}
		}
	}

	return true
}

// IsIdentStart returns whether a rune can be the start of an identifier
func IsIdentStart(r rune) bool {
	return unicode.IsLetter(r) || r == '_'
}

// IsIdentContinue returns whether a rune can be in the middle of an identifier
func IsIdentContinue(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_'
}
