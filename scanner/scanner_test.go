// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package scanner

import (
	//"fmt"
	"github.com/mjarmy/golem-lang/ast"
	"reflect"
	"testing"
)

func ok(t *testing.T, s *Scanner, tokenKind ast.TokenKind, text string, line int, col int) {

	token := &ast.Token{tokenKind, text, ast.Pos{line, col}}

	nextToken := s.Next()

	if !reflect.DeepEqual(*nextToken, *token) {
		t.Error(nextToken, " != ", token)
		panic("ok")
	}
}

func TestDelimiter(t *testing.T) {

	s := NewScanner("")
	ok(t, s, ast.Eof, "", 1, 1)
	ok(t, s, ast.Eof, "", 1, 1)
	ok(t, s, ast.Eof, "", 1, 1)
	ok(t, s, ast.Eof, "", 1, 1)

	s = NewScanner("#")
	ok(t, s, ast.UnexpectedChar, "#", 1, 1)
	ok(t, s, ast.UnexpectedChar, "#", 1, 1)
	ok(t, s, ast.UnexpectedChar, "#", 1, 1)

	s = NewScanner("+")
	ok(t, s, ast.Plus, "+", 1, 1)
	ok(t, s, ast.Eof, "", 1, 2)

	s = NewScanner("-\n/")
	ok(t, s, ast.Minus, "-", 1, 1)
	ok(t, s, ast.LineFeed, "\n", 2, 0)
	ok(t, s, ast.Slash, "/", 2, 1)
	ok(t, s, ast.Eof, "", 2, 2)

	s = NewScanner("+-*/)(")
	ok(t, s, ast.Plus, "+", 1, 1)
	ok(t, s, ast.Minus, "-", 1, 2)
	ok(t, s, ast.Star, "*", 1, 3)
	ok(t, s, ast.Slash, "/", 1, 4)
	ok(t, s, ast.Rparen, ")", 1, 5)
	ok(t, s, ast.Lparen, "(", 1, 6)
	ok(t, s, ast.Eof, "", 1, 7)

	s = NewScanner("}{==;=+ =,:.?[]=>")
	ok(t, s, ast.Rbrace, "}", 1, 1)
	ok(t, s, ast.Lbrace, "{", 1, 2)
	ok(t, s, ast.DblEq, "==", 1, 3)
	ok(t, s, ast.Semicolon, ";", 1, 5)
	ok(t, s, ast.Eq, "=", 1, 6)
	ok(t, s, ast.Plus, "+", 1, 7)
	ok(t, s, ast.Eq, "=", 1, 9)
	ok(t, s, ast.Comma, ",", 1, 10)
	ok(t, s, ast.Colon, ":", 1, 11)
	ok(t, s, ast.Dot, ".", 1, 12)
	ok(t, s, ast.Hook, "?", 1, 13)
	ok(t, s, ast.Lbracket, "[", 1, 14)
	ok(t, s, ast.Rbracket, "]", 1, 15)
	ok(t, s, ast.EqGt, "=>", 1, 16)
	ok(t, s, ast.Eof, "", 1, 18)

	s = NewScanner("! !=")
	ok(t, s, ast.Not, "!", 1, 1)
	ok(t, s, ast.NotEq, "!=", 1, 3)
	ok(t, s, ast.Eof, "", 1, 5)

	s = NewScanner("> >=")
	ok(t, s, ast.Gt, ">", 1, 1)
	ok(t, s, ast.GtEq, ">=", 1, 3)
	ok(t, s, ast.Eof, "", 1, 5)

	s = NewScanner("< <= <=>")
	ok(t, s, ast.Lt, "<", 1, 1)
	ok(t, s, ast.LtEq, "<=", 1, 3)
	ok(t, s, ast.Cmp, "<=>", 1, 6)
	ok(t, s, ast.Eof, "", 1, 9)

	s = NewScanner("& && | ||")
	ok(t, s, ast.Amp, "&", 1, 1)
	ok(t, s, ast.DblAmp, "&&", 1, 3)
	ok(t, s, ast.Pipe, "|", 1, 6)
	ok(t, s, ast.DblPipe, "||", 1, 8)
	ok(t, s, ast.Eof, "", 1, 10)

	s = NewScanner("%^~<<>>++--")
	ok(t, s, ast.Percent, "%", 1, 1)
	ok(t, s, ast.Caret, "^", 1, 2)
	ok(t, s, ast.Tilde, "~", 1, 3)
	ok(t, s, ast.DblLt, "<<", 1, 4)
	ok(t, s, ast.DblGt, ">>", 1, 6)
	ok(t, s, ast.DblPlus, "++", 1, 8)
	ok(t, s, ast.DblMinus, "--", 1, 10)
	ok(t, s, ast.Eof, "", 1, 12)

	s = NewScanner("+= -= *= /= %= ^= &= |= >>= <<= ")
	ok(t, s, ast.PlusEq, "+=", 1, 1)
	ok(t, s, ast.MinusEq, "-=", 1, 4)
	ok(t, s, ast.StarEq, "*=", 1, 7)
	ok(t, s, ast.SlashEq, "/=", 1, 10)
	ok(t, s, ast.PercentEq, "%=", 1, 13)
	ok(t, s, ast.CaretEq, "^=", 1, 16)
	ok(t, s, ast.AmpEq, "&=", 1, 19)
	ok(t, s, ast.PipeEq, "|=", 1, 22)
	ok(t, s, ast.DblGtEq, ">>=", 1, 25)
	ok(t, s, ast.DblLtEq, "<<=", 1, 29)
	ok(t, s, ast.Eof, "", 1, 33)
}

func TestInt(t *testing.T) {

	s := NewScanner("0")
	ok(t, s, ast.Int, "0", 1, 1)
	ok(t, s, ast.Eof, "", 1, 2)

	s = NewScanner("12+34 - 5 ")
	ok(t, s, ast.Int, "12", 1, 1)
	ok(t, s, ast.Plus, "+", 1, 3)
	ok(t, s, ast.Int, "34", 1, 4)
	ok(t, s, ast.Minus, "-", 1, 7)
	ok(t, s, ast.Int, "5", 1, 9)
	ok(t, s, ast.Eof, "", 1, 11)

	s = NewScanner("678")
	ok(t, s, ast.Int, "678", 1, 1)
	ok(t, s, ast.Eof, "", 1, 4)

	s = NewScanner("0 00")
	ok(t, s, ast.Int, "0", 1, 1)
	ok(t, s, ast.UnexpectedChar, "0", 1, 4)

	s = NewScanner("00 1")
	ok(t, s, ast.UnexpectedChar, "0", 1, 2)

	s = NewScanner("0xabcdef123456789")
	ok(t, s, ast.Int, "0xabcdef123456789", 1, 1)
	ok(t, s, ast.Eof, "", 1, 18)

	s = NewScanner("0xABCDEF")
	ok(t, s, ast.Int, "0xABCDEF", 1, 1)
	ok(t, s, ast.Eof, "", 1, 9)

	s = NewScanner("0x")
	ok(t, s, ast.UnexpectedEof, "", 1, 3)

	s = NewScanner("0xg")
	ok(t, s, ast.UnexpectedChar, "g", 1, 3)
}

func TestFloat(t *testing.T) {
	s := NewScanner("0.12 0.34")
	ok(t, s, ast.Float, "0.12", 1, 1)
	ok(t, s, ast.Float, "0.34", 1, 6)
	ok(t, s, ast.Eof, "", 1, 10)

	s = NewScanner("12.34 56.78")
	ok(t, s, ast.Float, "12.34", 1, 1)
	ok(t, s, ast.Float, "56.78", 1, 7)
	ok(t, s, ast.Eof, "", 1, 12)

	s = NewScanner("0.34E1")
	ok(t, s, ast.Float, "0.34E1", 1, 1)
	ok(t, s, ast.Eof, "", 1, 7)

	s = NewScanner("0.34E-1")
	ok(t, s, ast.Float, "0.34E-1", 1, 1)
	ok(t, s, ast.Eof, "", 1, 8)

	s = NewScanner("0.34E+1")
	ok(t, s, ast.Float, "0.34E+1", 1, 1)
	ok(t, s, ast.Eof, "", 1, 8)

	s = NewScanner("0.34e2")
	ok(t, s, ast.Float, "0.34e2", 1, 1)
	ok(t, s, ast.Eof, "", 1, 7)

	s = NewScanner("0e6")
	ok(t, s, ast.Float, "0e6", 1, 1)
	ok(t, s, ast.Eof, "", 1, 4)

	s = NewScanner("1e6")
	ok(t, s, ast.Float, "1e6", 1, 1)
	ok(t, s, ast.Eof, "", 1, 4)

	s = NewScanner("0.")
	ok(t, s, ast.UnexpectedEof, "", 1, 3)
	s = NewScanner("0. ")
	ok(t, s, ast.UnexpectedChar, " ", 1, 3)

	s = NewScanner("0.1e")
	ok(t, s, ast.UnexpectedEof, "", 1, 5)
	s = NewScanner("0.1e ")
	ok(t, s, ast.UnexpectedChar, " ", 1, 5)
}

func TestStr(t *testing.T) {
	s := NewScanner("''")
	ok(t, s, ast.Str, "", 1, 1)
	ok(t, s, ast.Eof, "", 1, 3)

	s = NewScanner("'a'")
	ok(t, s, ast.Str, "a", 1, 1)
	ok(t, s, ast.Eof, "", 1, 4)

	s = NewScanner("\"\"")
	ok(t, s, ast.Str, "", 1, 1)
	ok(t, s, ast.Eof, "", 1, 3)

	s = NewScanner("\"a\"")
	ok(t, s, ast.Str, "a", 1, 1)
	ok(t, s, ast.Eof, "", 1, 4)

	s = NewScanner("'ab' 'c'")
	ok(t, s, ast.Str, "ab", 1, 1)
	ok(t, s, ast.Str, "c", 1, 6)
	ok(t, s, ast.Eof, "", 1, 9)

	s = NewScanner("'ab")
	ok(t, s, ast.UnexpectedEof, "", 1, 4)

	s = NewScanner("'\n'")
	ok(t, s, ast.UnexpectedChar, "\n", 2, 0)

	s = NewScanner("'\\'\\n\\r\\t\\\\'")
	ok(t, s, ast.Str, "'\n\r\t\\", 1, 1)
	ok(t, s, ast.Eof, "", 1, 13)

	s = NewScanner("`a`")
	ok(t, s, ast.Str, "a", 1, 1)
	ok(t, s, ast.Eof, "", 1, 4)

	s = NewScanner("`a\nb`")
	ok(t, s, ast.Str, "a\nb", 1, 1)
	ok(t, s, ast.Eof, "", 2, 3)
}

func TestIdentOrKeyword(t *testing.T) {
	s := NewScanner("a bar")
	ok(t, s, ast.Ident, "a", 1, 1)
	ok(t, s, ast.Ident, "bar", 1, 3)
	ok(t, s, ast.Eof, "", 1, 6)

	s = NewScanner("_ zork")
	ok(t, s, ast.BlankDent, "_", 1, 1)
	ok(t, s, ast.Ident, "zork", 1, 3)
	ok(t, s, ast.Eof, "", 1, 7)

	s = NewScanner("null true false")
	ok(t, s, ast.Null, "null", 1, 1)
	ok(t, s, ast.True, "true", 1, 6)
	ok(t, s, ast.False, "false", 1, 11)
	ok(t, s, ast.Eof, "", 1, 16)

	s = NewScanner("if else")
	ok(t, s, ast.If, "if", 1, 1)
	ok(t, s, ast.Else, "else", 1, 4)
	ok(t, s, ast.Eof, "", 1, 8)

	s = NewScanner("while break continue")
	ok(t, s, ast.While, "while", 1, 1)
	ok(t, s, ast.Break, "break", 1, 7)
	ok(t, s, ast.Continue, "continue", 1, 13)
	ok(t, s, ast.Eof, "", 1, 21)

	s = NewScanner("fn return const let for in")
	ok(t, s, ast.Fn, "fn", 1, 1)
	ok(t, s, ast.Return, "return", 1, 4)
	ok(t, s, ast.Const, "const", 1, 11)
	ok(t, s, ast.Let, "let", 1, 17)
	ok(t, s, ast.For, "for", 1, 21)
	ok(t, s, ast.In, "in", 1, 25)
	ok(t, s, ast.Eof, "", 1, 27)

	s = NewScanner("switch case default prop")
	ok(t, s, ast.Switch, "switch", 1, 1)
	ok(t, s, ast.Case, "case", 1, 8)
	ok(t, s, ast.Default, "default", 1, 13)
	ok(t, s, ast.Prop, "prop", 1, 21)
	ok(t, s, ast.Eof, "", 1, 25)

	s = NewScanner("try catch finally throw")
	ok(t, s, ast.Try, "try", 1, 1)
	ok(t, s, ast.Catch, "catch", 1, 5)
	ok(t, s, ast.Finally, "finally", 1, 11)
	ok(t, s, ast.Throw, "throw", 1, 19)
	ok(t, s, ast.Eof, "", 1, 24)

	s = NewScanner("go module import")
	ok(t, s, ast.Go, "go", 1, 1)
	ok(t, s, ast.Module, "module", 1, 4)
	ok(t, s, ast.Import, "import", 1, 11)
	ok(t, s, ast.Eof, "", 1, 17)

	s = NewScanner("struct this has dict set")
	ok(t, s, ast.Struct, "struct", 1, 1)
	ok(t, s, ast.This, "this", 1, 8)
	ok(t, s, ast.Has, "has", 1, 13)
	ok(t, s, ast.Dict, "dict", 1, 17)
	ok(t, s, ast.Set, "set", 1, 22)
	ok(t, s, ast.Eof, "", 1, 25)
}

func TestComments(t *testing.T) {

	s := NewScanner("1 //foo\n2")
	ok(t, s, ast.Int, "1", 1, 1)
	ok(t, s, ast.LineFeed, "\n", 2, 0)
	ok(t, s, ast.Int, "2", 2, 1)
	ok(t, s, ast.Eof, "", 2, 2)

	s = NewScanner("1 2 //foo")
	ok(t, s, ast.Int, "1", 1, 1)
	ok(t, s, ast.Int, "2", 1, 3)
	ok(t, s, ast.Eof, "", 1, 10)

	s = NewScanner("1 /*foo*/2")
	ok(t, s, ast.Int, "1", 1, 1)
	ok(t, s, ast.Int, "2", 1, 10)
	ok(t, s, ast.Eof, "", 1, 11)

	s = NewScanner("1 2/**/")
	ok(t, s, ast.Int, "1", 1, 1)
	ok(t, s, ast.Int, "2", 1, 3)
	ok(t, s, ast.Eof, "", 1, 8)

	s = NewScanner("1 /*")
	ok(t, s, ast.Int, "1", 1, 1)
	ok(t, s, ast.UnexpectedEof, "", 1, 5)

	s = NewScanner("1 /* *")
	ok(t, s, ast.Int, "1", 1, 1)
	ok(t, s, ast.UnexpectedEof, "", 1, 7)

}

func TestUnicode(t *testing.T) {

	s := NewScanner("'\\u'")
	ok(t, s, ast.UnexpectedChar, "'", 1, 4)
	s = NewScanner("'\\u['")
	ok(t, s, ast.UnexpectedChar, "[", 1, 4)
	s = NewScanner("'\\u{z'")
	ok(t, s, ast.UnexpectedChar, "z", 1, 5)
	s = NewScanner("'\\u{a'")
	ok(t, s, ast.UnexpectedChar, "'", 1, 6)
	s = NewScanner("'\\u{a]'")
	ok(t, s, ast.UnexpectedChar, "]", 1, 6)
	s = NewScanner("'\\u{1234567}'")
	ok(t, s, ast.UnexpectedChar, "7", 1, 11)

	s = NewScanner("'\\u{24}'")
	ok(t, s, ast.Str, "$", 1, 1)
	s = NewScanner("'\\u{2665}'")
	ok(t, s, ast.Str, "♥", 1, 1)
	s = NewScanner("'\\u{1F496}'")
	ok(t, s, ast.Str, "💖", 1, 1)
	s = NewScanner("'\\u{1f496}\\u{2665}\\u{24}'")
	ok(t, s, ast.Str, "💖♥$", 1, 1)
}
