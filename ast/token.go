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
	UNEXPECTED_CHAR TokenKind = iota
	UNEXPECTED_EOF
	badKind

	EOF
	LINE_FEED

	PLUS
	DBL_PLUS
	MINUS
	DBL_MINUS

	STAR
	SLASH
	LPAREN
	RPAREN
	LBRACE
	RBRACE
	LBRACKET
	RBRACKET
	SEMICOLON
	COLON
	COMMA
	DOT
	HOOK

	EQ
	DBL_EQ
	EQ_GT
	NOT
	NOT_EQ
	GT
	DBL_GT
	GT_EQ
	LT
	DBL_LT
	LT_EQ
	CMP

	PIPE
	DBL_PIPE
	AMP
	DBL_AMP

	PERCENT
	CARET
	TILDE

	PLUS_EQ
	MINUS_EQ
	STAR_EQ
	SLASH_EQ
	PERCENT_EQ
	CARET_EQ
	AMP_EQ
	PIPE_EQ
	DBL_LT_EQ
	DBL_GT_EQ

	basicBegin
	NullValue
	True
	False
	STR
	INT
	FLOAT
	basicEnd

	IDENT

	BLANK_IDENT
	IF
	ELSE
	WHILE
	BREAK
	CONTINUE
	FN
	RETURN
	CONST
	LET
	FOR
	IN
	SWITCH
	CASE
	DEFAULT
	PROP

	STRUCT
	THIS
	HAS
	DICT
	SET

	TRY
	CATCH
	FINALLY
	THROW

	GO

	MODULE
	IMPORT

	RESERVED
)

func (t TokenKind) String() string {
	switch t {
	case UNEXPECTED_CHAR:
		return "UNEXPECTED_CHAR"
	case UNEXPECTED_EOF:
		return "UNEXPECTED_EOF"

	case EOF:
		return "EOF"
	case LINE_FEED:
		return "LINE_FEED"

	case PLUS:
		return "PLUS"
	case DBL_PLUS:
		return "DBL_PLUS"
	case MINUS:
		return "MINUS"
	case DBL_MINUS:
		return "DBL_MINUS"
	case STAR:
		return "STAR"
	case SLASH:
		return "SLASH"
	case LPAREN:
		return "LPAREN"
	case RPAREN:
		return "RPAREN"
	case LBRACE:
		return "LBRACE"
	case RBRACE:
		return "RBRACE"
	case LBRACKET:
		return "LBRACKET"
	case RBRACKET:
		return "RBRACKET"

	case EQ:
		return "EQ"
	case DBL_EQ:
		return "DBL_EQ"
	case EQ_GT:
		return "EQ_GT"
	case SEMICOLON:
		return "SEMICOLON"
	case COLON:
		return "COLON"
	case COMMA:
		return "COMMA"
	case DOT:
		return "DOT"
	case HOOK:
		return "HOOK"

	case PERCENT:
		return "PERCENT"
	case CARET:
		return "CARET"
	case TILDE:
		return "TILDE"

	case PLUS_EQ:
		return "PLUS_EQ"
	case MINUS_EQ:
		return "MINUS_EQ"
	case STAR_EQ:
		return "STAR_EQ"
	case SLASH_EQ:
		return "SLASH_EQ"
	case PERCENT_EQ:
		return "PERCENT_EQ"
	case CARET_EQ:
		return "CARET_EQ"
	case AMP_EQ:
		return "AMP_EQ"
	case PIPE_EQ:
		return "PIPE_EQ"
	case DBL_LT_EQ:
		return "DBL_LT_EQ"
	case DBL_GT_EQ:
		return "DBL_GT_EQ"

	case NullValue:
		return "NullValue"
	case True:
		return "True"
	case False:
		return "False"
	case STR:
		return "STR"
	case INT:
		return "INT"
	case FLOAT:
		return "FLOAT"

	case IDENT:
		return "IDENT"

	case IF:
		return "IF"
	case ELSE:
		return "ELSE"
	case WHILE:
		return "WHILE"
	case BREAK:
		return "BREAK"
	case CONTINUE:
		return "CONTINUE"
	case FN:
		return "FN"
	case RETURN:
		return "RETURN"
	case CONST:
		return "CONST"
	case LET:
		return "LET"
	case FOR:
		return "FOR"
	case IN:
		return "IN"
	case SWITCH:
		return "SWITCH"
	case CASE:
		return "CASE"
	case DEFAULT:
		return "DEFAULT"
	case PROP:
		return "PROP"

	case STRUCT:
		return "STRUCT"
	case THIS:
		return "THIS"
	case HAS:
		return "HAS"
	case DICT:
		return "DICT"
	case SET:
		return "SET"

	case TRY:
		return "TRY"
	case FINALLY:
		return "FINALLY"
	case CATCH:
		return "CATCH"
	case THROW:
		return "THROW"

	case GO:
		return "GO"

	case MODULE:
		return "MODULE"
	case IMPORT:
		return "IMPORT"

	case RESERVED:
		return "RESERVED"

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
