// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package ast

import (
	"fmt"
)

//-------------------------------------
// Pos

// Pos represents a Line-and-Column location in Golem source code.
type Pos struct {
	Line int
	Col  int
}

func (p Pos) String() string {
	return fmt.Sprintf("%d:%d", p.Line, p.Col)
}

// Advance advances a Pos column forwards
func (p Pos) Advance(len int) Pos {
	return Pos{p.Line, p.Col + len}
}

// TokenKind defines all the various kinds of token
type TokenKind int

// The various kinds of token
const (
	UnexpectedChar TokenKind = iota
	UnexpectedEOF
	badKind

	EOF
	LineFeed

	Plus
	DoublePlus
	Minus
	DoubleMinus

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
	//DoubleDot
	TripleDot
	Hook

	Eq
	DoubleEq
	EqGt
	Not
	NotEq
	Gt
	DoubleGt
	GtEq
	Lt
	DoubleLt
	LtEq
	Cmp

	Pipe
	DoublePipe
	Amp
	DoubleAmp

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
	DoubleLtEq
	DoubleGtEq

	basicBegin
	Null
	True
	False
	Str
	Int
	Float
	basicEnd

	Ident
	MagicField

	BlankIdent
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
	Dict
	Set

	Try
	Catch
	Finally
	Throw

	Go

	Import

	Reserved
)

func (t TokenKind) String() string {
	switch t {
	case UnexpectedChar:
		return "UnexpectedChar"
	case UnexpectedEOF:
		return "UnexpectedEOF"

	case EOF:
		return "EOF"
	case LineFeed:
		return "LineFeed"

	case Plus:
		return "Plus"
	case DoublePlus:
		return "DoublePlus"
	case Minus:
		return "Minus"
	case DoubleMinus:
		return "DoubleMinus"
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
	case DoubleEq:
		return "DoubleEq"
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
		//	case DoubleDot:
		//		return "DoubleDot"
	case TripleDot:
		return "TripleDot"
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
	case DoubleLtEq:
		return "DoubleLtEq"
	case DoubleGtEq:
		return "DoubleGtEq"

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
	case MagicField:
		return "MagicField"

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

	case Import:
		return "Import"

	case Reserved:
		return "Reserved"

	default:
		panic("unreachable")
	}
}

// Token is produced by the Scanner.  The Parser uses Tokens to assemble
// an Abstract Syntax Tree.
type Token struct {
	Kind     TokenKind
	Text     string
	Position Pos
}

func (t *Token) String() string {
	return fmt.Sprintf("Token(%v, %q, %v)", t.Kind, t.Text, t.Position)
}

// IsBad returns whether or not a Token is considered to be invalid.
func (t *Token) IsBad() bool {
	return t.Kind < badKind
}

// IsBasic returns whether or not a token represents one of the basic types.
func (t *Token) IsBasic() bool {
	return t.Kind > basicBegin && t.Kind < basicEnd
}
