// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package scanner

import (
	"bytes"
	"github.com/mjarmy/golem-lang/ast"
	"io"
	"strconv"
	"strings"
	"unicode"
)

const eof rune = -1

type (
	curRune struct {
		r    rune
		size int
		idx  int
	}

	// Scanner scans Golem source code and produces a stream of tokens.
	Scanner struct {
		source    string
		reader    io.RuneReader
		cur       curRune
		pos       ast.Pos
		isDone    bool
		doneToken *ast.Token
	}
)

// NewScanner creates a new Scanner
func NewScanner(source string) *Scanner {
	reader := strings.NewReader(source)
	s := &Scanner{source, reader, curRune{0, 1, -1}, ast.Pos{1, 0}, false, nil}
	s.consume()
	return s
}

// Next produces the next token in the stream.  By convention, if the stream is
// finished, the last token is produced over and over again.
func (s *Scanner) Next() *ast.Token {

	// IfStmt we are already finished, then by convention we return
	// the last token again.  This makes it easier
	// for the parser to implement lookahead.
	if s.isDone {
		return s.doneToken
	}

	for {
		pos := s.pos
		r := s.cur.r

		switch {

		case isWhitespace(r):
			s.consume()

		case r == '\n':
			s.consume()
			return &ast.Token{ast.LineFeed, "\n", pos}

		case r == '/':
			s.consume()
			r = s.cur.r

			switch r {

			// line comment
			case '/':
				s.consume()
				r = s.cur.r
				for (r != '\n') && (r != eof) {
					s.consume()
					r = s.cur.r
				}

			// block comment
			case '*':
				s.consume()
				r = s.cur.r

			loop:
				for {
					switch r {
					case '*':
						s.consume()
						r = s.cur.r

						switch r {
						case '/':
							s.consume()
							break loop
						case eof:
							return s.unexpectedChar(r, s.pos)
						default:
							s.consume()
							r = s.cur.r
						}
					case eof:
						return s.unexpectedChar(r, s.pos)
					default:
						s.consume()
						r = s.cur.r
					}
				}

			case '=':
				s.consume()
				return &ast.Token{ast.SlashEq, "/=", pos}

			default:
				return &ast.Token{ast.Slash, "/", pos}
			}

		case r == '+':
			s.consume()
			r = s.cur.r
			if r == '+' {
				s.consume()
				return &ast.Token{ast.DblPlus, "++", pos}
			} else if r == '=' {
				s.consume()
				return &ast.Token{ast.PlusEq, "+=", pos}
			} else {
				return &ast.Token{ast.Plus, "+", pos}
			}

		case r == '-':
			s.consume()
			r = s.cur.r
			if r == '-' {
				s.consume()
				return &ast.Token{ast.DblMinus, "--", pos}
			} else if r == '=' {
				s.consume()
				return &ast.Token{ast.MinusEq, "-=", pos}
			} else {
				return &ast.Token{ast.Minus, "-", pos}
			}

		case r == '*':
			s.consume()
			r = s.cur.r
			if r == '=' {
				s.consume()
				return &ast.Token{ast.StarEq, "*=", pos}
			}
			return &ast.Token{ast.Star, "*", pos}

		case r == '(':
			s.consume()
			return &ast.Token{ast.Lparen, "(", pos}
		case r == ')':
			s.consume()
			return &ast.Token{ast.Rparen, ")", pos}
		case r == '{':
			s.consume()
			return &ast.Token{ast.Lbrace, "{", pos}
		case r == '}':
			s.consume()
			return &ast.Token{ast.Rbrace, "}", pos}
		case r == '[':
			s.consume()
			return &ast.Token{ast.Lbracket, "[", pos}
		case r == ']':
			s.consume()
			return &ast.Token{ast.Rbracket, "]", pos}
		case r == ';':
			s.consume()
			return &ast.Token{ast.Semicolon, ";", pos}
		case r == ':':
			s.consume()
			return &ast.Token{ast.Colon, ":", pos}
		case r == ',':
			s.consume()
			return &ast.Token{ast.Comma, ",", pos}
		case r == '.':
			s.consume()
			return &ast.Token{ast.Dot, ".", pos}
		case r == '?':
			s.consume()
			return &ast.Token{ast.Hook, "?", pos}

		case r == '%':
			s.consume()
			r := s.cur.r
			if r == '=' {
				s.consume()
				return &ast.Token{ast.PercentEq, "%=", pos}
			}
			return &ast.Token{ast.Percent, "%", pos}

		case r == '^':
			s.consume()
			r := s.cur.r
			if r == '=' {
				s.consume()
				return &ast.Token{ast.CaretEq, "^=", pos}
			}
			return &ast.Token{ast.Caret, "^", pos}

		case r == '~':
			s.consume()
			return &ast.Token{ast.Tilde, "~", pos}

		case r == '=':
			s.consume()
			r := s.cur.r
			if r == '=' {
				s.consume()
				return &ast.Token{ast.DblEq, "==", pos}
			} else if r == '>' {
				s.consume()
				return &ast.Token{ast.EqGt, "=>", pos}
			} else {
				return &ast.Token{ast.Eq, "=", pos}
			}

		case r == '!':
			s.consume()
			r := s.cur.r
			if r == '=' {
				s.consume()
				return &ast.Token{ast.NotEq, "!=", pos}
			}
			return &ast.Token{ast.Not, "!", pos}
		case r == '>':
			s.consume()
			r := s.cur.r
			if r == '=' {
				s.consume()
				return &ast.Token{ast.GtEq, ">=", pos}
			} else if r == '>' {
				s.consume()
				r := s.cur.r
				if r == '=' {
					s.consume()
					return &ast.Token{ast.DblGtEq, ">>=", pos}
				}
				return &ast.Token{ast.DblGt, ">>", pos}
			} else {
				return &ast.Token{ast.Gt, ">", pos}
			}
		case r == '<':
			s.consume()
			r := s.cur.r
			if r == '=' {
				s.consume()
				r := s.cur.r
				if r == '>' {
					s.consume()
					return &ast.Token{ast.Cmp, "<=>", pos}
				}
				return &ast.Token{ast.LtEq, "<=", pos}
			} else if r == '<' {
				s.consume()
				r := s.cur.r
				if r == '=' {
					s.consume()
					return &ast.Token{ast.DblLtEq, "<<=", pos}
				}
				return &ast.Token{ast.DblLt, "<<", pos}
			} else {
				return &ast.Token{ast.Lt, "<", pos}
			}

		case r == '|':
			s.consume()
			r := s.cur.r
			if r == '|' {
				s.consume()
				return &ast.Token{ast.DblPipe, "||", pos}
			} else if r == '=' {
				s.consume()
				return &ast.Token{ast.PipeEq, "|=", pos}
			} else {
				return &ast.Token{ast.Pipe, "|", pos}
			}
		case r == '&':
			s.consume()
			r := s.cur.r
			if r == '&' {
				s.consume()
				return &ast.Token{ast.DblAmp, "&&", pos}
			} else if r == '=' {
				s.consume()
				return &ast.Token{ast.AmpEq, "&=", pos}
			} else {
				return &ast.Token{ast.Amp, "&", pos}
			}

		case r == '\'':
			return s.nextStr('\'')

		case r == '"':
			return s.nextStr('"')

		case r == '`':
			return s.nextMultilineStr()

		case isDigit(r):
			return s.nextNumber()

		case isIdentStart(r):
			return s.nextIdentOrKeyword()

		case r == eof:
			s.isDone = true
			s.doneToken = &ast.Token{ast.Eof, "", pos}
			return s.doneToken

		default:
			return s.unexpectedChar(r, pos)
		}
	}
}

func (s *Scanner) nextIdentOrKeyword() *ast.Token {

	pos := s.pos
	begin := s.cur.idx
	s.consume()

	s.acceptWhile(isIdentContinue)

	text := s.source[begin:s.cur.idx]
	switch text {

	case "_":
		return &ast.Token{ast.BlankDent, text, pos}
	case "null":
		return &ast.Token{ast.Null, text, pos}
	case "true":
		return &ast.Token{ast.True, text, pos}
	case "false":
		return &ast.Token{ast.False, text, pos}
	case "if":
		return &ast.Token{ast.If, text, pos}
	case "else":
		return &ast.Token{ast.Else, text, pos}
	case "while":
		return &ast.Token{ast.While, text, pos}
	case "break":
		return &ast.Token{ast.Break, text, pos}
	case "continue":
		return &ast.Token{ast.Continue, text, pos}
	case "fn":
		return &ast.Token{ast.Fn, text, pos}
	case "return":
		return &ast.Token{ast.Return, text, pos}
	case "const":
		return &ast.Token{ast.Const, text, pos}
	case "let":
		return &ast.Token{ast.Let, text, pos}
	case "for":
		return &ast.Token{ast.For, text, pos}
	case "in":
		return &ast.Token{ast.In, text, pos}
	case "switch":
		return &ast.Token{ast.Switch, text, pos}
	case "case":
		return &ast.Token{ast.Case, text, pos}
	case "default":
		return &ast.Token{ast.Default, text, pos}
	case "prop":
		return &ast.Token{ast.Prop, text, pos}
	case "try":
		return &ast.Token{ast.Try, text, pos}
	case "catch":
		return &ast.Token{ast.Catch, text, pos}
	case "finally":
		return &ast.Token{ast.Finally, text, pos}
	case "throw":
		return &ast.Token{ast.Throw, text, pos}
	case "go":
		return &ast.Token{ast.Go, text, pos}
	case "module":
		return &ast.Token{ast.Module, text, pos}
	case "import":
		return &ast.Token{ast.Import, text, pos}
	case "struct":
		return &ast.Token{ast.Struct, text, pos}
	case "dict":
		return &ast.Token{ast.Dict, text, pos}
	case "set":
		return &ast.Token{ast.Set, text, pos}
	case "this":
		return &ast.Token{ast.This, text, pos}
	case "has":
		return &ast.Token{ast.Has, text, pos}

	case "byte", "defer", "goto", "like", "native", "package",
		"priv", "private", "prot", "protected", "pub", "public",
		"rune", "select", "static", "sync", "rsync", "with", "yield":

		// reserve a bunch of keywords just in case
		return &ast.Token{ast.Reserved, text, pos}

	default:
		return &ast.Token{ast.Ident, text, pos}
	}
}

func (s *Scanner) nextStr(delim rune) *ast.Token {

	pos := s.pos
	s.consume()

	var buf bytes.Buffer

	for {
		r := s.cur.r

		switch {

		case r == delim:
			// end of string
			s.consume()
			return &ast.Token{ast.Str, buf.String(), pos}

		case r == '\\':
			// escaped character
			s.consume()
			r = s.cur.r
			switch r {
			case '\\':
				s.consume()
				buf.WriteRune('\\')
			case 'n':
				s.consume()
				buf.WriteRune('\n')
			case 'r':
				s.consume()
				buf.WriteRune('\r')
			case 't':
				s.consume()
				buf.WriteRune('\t')
			case 'u':
				s.consume()
				u, err := s.unicodeRune()
				if err != nil {
					return err
				}
				buf.WriteRune(u)
			case delim:
				s.consume()
				buf.WriteRune(delim)
			default:
				return s.unexpectedChar(r, s.pos)
			}

		case r == eof:
			// unterminated string literal
			return s.unexpectedChar(r, s.pos)

		case r < ' ':
			// disallow embedded control characters
			return s.unexpectedChar(r, s.pos)

		default:
			buf.WriteRune(r)
			s.consume()
		}
	}
}

func (s *Scanner) unicodeRune() (rune, *ast.Token) {

	if s.cur.r != '{' {
		return 0, s.unexpectedChar(s.cur.r, s.pos)
	}
	s.consume()

	begin := s.cur.idx
	pos := s.pos
	t := s.expect(isHexDigit)
	if t != nil {
		return 0, t
	}
	s.acceptWhile(isHexDigit)
	end := s.cur.idx

	text := s.source[begin:end]
	if len(text) > 6 {
		// too long
		runes := []rune(text)
		return 0, s.unexpectedChar(runes[len(runes)-1], pos.Advance(6))
	}

	if s.cur.r != '}' {
		return 0, s.unexpectedChar(s.cur.r, s.pos)
	}
	s.consume()

	n, err := strconv.ParseInt(text, 16, 32)
	if err != nil {
		panic("unreachable")
	}
	return rune(int32(n)), nil
}

func (s *Scanner) nextMultilineStr() *ast.Token {

	pos := s.pos
	s.consume()

	var buf bytes.Buffer

	for {
		r := s.cur.r

		switch {

		case r == '`':
			// end of string
			s.consume()
			return &ast.Token{ast.Str, buf.String(), pos}

		case r == eof:
			// unterminated string literal
			return s.unexpectedChar(r, s.pos)

		case r < ' ' && !(r == '\n' || r == '\t' || r == '\r'):
			// disallow embedded control characters
			return s.unexpectedChar(r, s.pos)

		default:
			buf.WriteRune(r)
			s.consume()
		}
	}
}

func (s *Scanner) nextNumber() *ast.Token {

	pos := s.pos
	r := s.cur.r
	begin := s.cur.idx
	s.consume()

	if r == '0' {
		r := s.cur.r

		switch {

		case isDigit(r):
			return s.unexpectedChar(r, s.pos)

		case r == '.' || isExp(r):
			return s.nextFloat(begin, pos)

		case r == 'x':
			return s.nextHexInt(begin, pos)

		default:
			return &ast.Token{ast.Int, "0", pos}
		}

	} else {
		s.acceptWhile(isDigit)
		r := s.cur.r
		if r == '.' || isExp(r) {
			return s.nextFloat(begin, pos)
		}
		return &ast.Token{ast.Int, s.source[begin:s.cur.idx], pos}
	}

}

func (s *Scanner) nextHexInt(begin int, pos ast.Pos) *ast.Token {

	s.consume()

	t := s.expect(isHexDigit)
	if t != nil {
		return t
	}
	s.acceptWhile(isHexDigit)

	return &ast.Token{ast.Int, s.source[begin:s.cur.idx], pos}
}

func (s *Scanner) nextFloat(begin int, pos ast.Pos) *ast.Token {

	s.consume()

	t := s.expect(isDigit)
	if t != nil {
		return t
	}
	s.acceptWhile(isDigit)

	if s.accept(isExp) {
		s.accept(func(r rune) bool { return (r == '+') || (r == '-') })

		t := s.expect(isDigit)
		if t != nil {
			return t
		}
		s.acceptWhile(isDigit)
	}

	return &ast.Token{ast.Float, s.source[begin:s.cur.idx], pos}
}

// accept a rune that matches the given function
func (s *Scanner) accept(fn func(rune) bool) bool {

	r := s.cur.r
	if fn(r) {
		s.consume()
		return true
	}
	return false
}

// accept a sequence of runes that match the given function
func (s *Scanner) acceptWhile(fn func(rune) bool) {

	for {
		r := s.cur.r
		if fn(r) {
			s.consume()
		} else {
			return
		}
	}
}

// expect a rune that match the given function
func (s *Scanner) expect(fn func(rune) bool) *ast.Token {

	pos := s.pos
	r := s.cur.r

	if fn(r) {
		s.consume()
		return nil
	}
	return s.unexpectedChar(r, pos)
}

func (s *Scanner) unexpectedChar(r rune, pos ast.Pos) *ast.Token {
	s.isDone = true
	if r == eof {
		s.doneToken = &ast.Token{ast.UnexpectedEof, "", pos}
	} else {
		s.doneToken = &ast.Token{ast.UnexpectedChar, string(r), pos}
	}
	return s.doneToken
}

// consume the current rune
func (s *Scanner) consume() {

	lastSize := s.cur.size

	r, size, err := s.reader.ReadRune()
	s.cur.size = size
	s.cur.idx += lastSize

	// set eof if there was an error
	if err == nil {
		s.cur.r = r
	} else {
		s.cur.r = eof
	}

	// advance position
	if r == '\n' {
		s.pos.Line++
		s.pos.Col = 0
	} else {
		s.pos.Col += lastSize
	}

}

func isWhitespace(r rune) bool {
	return r == ' ' || r == '\t' || r == '\r'
}

func isDigit(r rune) bool {
	return (r >= '0') && (r <= '9')
}

func isHexDigit(r rune) bool {
	return (r >= '0') && (r <= '9') ||
		(r >= 'a') && (r <= 'f') ||
		(r >= 'A') && (r <= 'F')
}

func isIdentStart(r rune) bool {
	return unicode.IsLetter(r) || r == '_'
}

func isIdentContinue(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_'
}

func isExp(r rune) bool {
	return (r == 'e') || (r == 'E')
}
