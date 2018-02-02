// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package ast

import (
	"fmt"
)

//-------------------------------------
// Pos

type Pos struct {
	Line int
	Col  int
}

func (p Pos) String() string {
	return fmt.Sprintf("(%d, %d)", p.Line, p.Col)
}

func (p Pos) Advance(len int) Pos {
	return Pos{p.Line, p.Col + len}
}

//-------------------------------------
// Token

type TokenKind int

const (
	UnexpectedChar TokenKind = iota
	UnexpectedEof
	badKind

	Eof
	LineFeed

	Plus
	DblPlus
	Minus
	DblMinus

	Star
	Slash
	Lparen
	Rparen
	Lbrace
	Rbrace
	Lbracket
	Rbracket
	Semicolon
	Colon
	Comma
	Dot
	Hook

	Eq
	DblEq
	EqGt
	Not
	NotEq
	Gt
	DblGt
	GtEq
	Lt
	DblLt
	LtEq
	Cmp

	Pipe
	DblPipe
	Amp
	DblAmp

	Percent
	Caret
	Tilde

	PlusEq
	MinusEq
	StarEq
	SlashEq
	PercentEq
	CaretEq
	AmpEq
	PipeEq
	DblLtEq
	DblGtEq

	basicBegin
	Null
	True
	False
	Str
	Int
	Float
	basicEnd

	Ident

	BlankDent
	If
	Else
	While
	Break
	Continue
	Fn
	Return
	Const
	Let
	For
	In
	Switch
	Case
	Default
	Prop

	Struct
	This
	Has
	Dict
	Set

	Try
	Catch
	Finally
	Throw

	Go

	Module
	Import

	Reserved
)

func (t TokenKind) String() string {
	switch t {
	case UnexpectedChar:
		return "UnexpectedChar"
	case UnexpectedEof:
		return "UnexpectedEof"

	case Eof:
		return "Eof"
	case LineFeed:
		return "LineFeed"

	case Plus:
		return "Plus"
	case DblPlus:
		return "DblPlus"
	case Minus:
		return "Minus"
	case DblMinus:
		return "DblMinus"
	case Star:
		return "Star"
	case Slash:
		return "Slash"
	case Lparen:
		return "Lparen"
	case Rparen:
		return "Rparen"
	case Lbrace:
		return "Lbrace"
	case Rbrace:
		return "Rbrace"
	case Lbracket:
		return "Lbracket"
	case Rbracket:
		return "Rbracket"

	case Eq:
		return "Eq"
	case DblEq:
		return "DblEq"
	case EqGt:
		return "EqGt"
	case Semicolon:
		return "Semicolon"
	case Colon:
		return "Colon"
	case Comma:
		return "Comma"
	case Dot:
		return "Dot"
	case Hook:
		return "Hook"

	case Percent:
		return "Percent"
	case Caret:
		return "Caret"
	case Tilde:
		return "Tilde"

	case PlusEq:
		return "PlusEq"
	case MinusEq:
		return "MinusEq"
	case StarEq:
		return "StarEq"
	case SlashEq:
		return "SlashEq"
	case PercentEq:
		return "PercentEq"
	case CaretEq:
		return "CaretEq"
	case AmpEq:
		return "AmpEq"
	case PipeEq:
		return "PipeEq"
	case DblLtEq:
		return "DblLtEq"
	case DblGtEq:
		return "DblGtEq"

	case Null:
		return "Null"
	case True:
		return "True"
	case False:
		return "False"
	case Str:
		return "Str"
	case Int:
		return "Int"
	case Float:
		return "Float"

	case Ident:
		return "Ident"

	case If:
		return "If"
	case Else:
		return "Else"
	case While:
		return "While"
	case Break:
		return "Break"
	case Continue:
		return "Continue"
	case Fn:
		return "Fn"
	case Return:
		return "Return"
	case Const:
		return "Const"
	case Let:
		return "Let"
	case For:
		return "For"
	case In:
		return "In"
	case Switch:
		return "Switch"
	case Case:
		return "Case"
	case Default:
		return "Default"
	case Prop:
		return "Prop"

	case Struct:
		return "Struct"
	case This:
		return "This"
	case Has:
		return "Has"
	case Dict:
		return "Dict"
	case Set:
		return "Set"

	case Try:
		return "Try"
	case Finally:
		return "Finally"
	case Catch:
		return "Catch"
	case Throw:
		return "Throw"

	case Go:
		return "Go"

	case Module:
		return "Module"
	case Import:
		return "Import"

	case Reserved:
		return "Reserved"

	default:
		panic("unreachable")
	}
}

type Token struct {
	Kind     TokenKind
	Text     string
	Position Pos
}

func (t *Token) String() string {
	return fmt.Sprintf("Token(%v, %q, %v)", t.Kind, t.Text, t.Position)
}

func (t *Token) IsBad() bool {
	return t.Kind < badKind
}

func (t *Token) IsBasic() bool {
	return t.Kind > basicBegin && t.Kind < basicEnd
}
