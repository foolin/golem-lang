// Copyright 2017 The Golem Project Developers
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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

//-------------------------------------
// Token

type TokenKind int

const (
	UNEXPECTED_CHAR TokenKind = iota
	UNEXPECTED_EOF
	badKind

	EOF

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
	NULL
	TRUE
	FALSE
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

	case NULL:
		return "NULL"
	case TRUE:
		return "TRUE"
	case FALSE:
		return "FALSE"
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
