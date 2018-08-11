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
		Name string
		Path string
		Code string

		reader    io.RuneReader
		cur       curRune
		pos       ast.Pos
		isDone    bool
		doneToken *ast.Token
	}
)

// NewScanner creates a new Scanner
func NewScanner(name, path, code string) *Scanner {

	s := &Scanner{
		Name: name,
		Path: path,
		Code: code,

		reader:    strings.NewReader(code),
		cur:       curRune{r: 0, size: 1, idx: -1},
		pos:       ast.Pos{Line: 1, Col: 0},
		isDone:    false,
		doneToken: nil,
	}

	s.consume()
	return s
}

// Next produces the next token in the stream.  By convention, if the stream is
// finished, the last token is produced over and over again.
// This makes it easier for the parser to implement lookahead.
func (s *Scanner) Next() *ast.Token {

	// If we are already finished, then return the last token
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
			return &ast.Token{Kind: ast.LineFeed, Text: "\n", Position: pos}

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
				return &ast.Token{Kind: ast.SlashEq, Text: "/=", Position: pos}

			default:
				return &ast.Token{Kind: ast.Slash, Text: "/", Position: pos}
			}

		case r == '+':
			s.consume()
			r = s.cur.r
			if r == '+' {
				s.consume()
				return &ast.Token{Kind: ast.DblPlus, Text: "++", Position: pos}
			} else if r == '=' {
				s.consume()
				return &ast.Token{Kind: ast.PlusEq, Text: "+=", Position: pos}
			} else {
				return &ast.Token{Kind: ast.Plus, Text: "+", Position: pos}
			}

		case r == '-':
			s.consume()
			r = s.cur.r
			if r == '-' {
				s.consume()
				return &ast.Token{Kind: ast.DblMinus, Text: "--", Position: pos}
			} else if r == '=' {
				s.consume()
				return &ast.Token{Kind: ast.MinusEq, Text: "-=", Position: pos}
			} else {
				return &ast.Token{Kind: ast.Minus, Text: "-", Position: pos}
			}

		case r == '*':
			s.consume()
			r = s.cur.r
			if r == '=' {
				s.consume()
				return &ast.Token{Kind: ast.StarEq, Text: "*=", Position: pos}
			}
			return &ast.Token{Kind: ast.Star, Text: "*", Position: pos}

		case r == '(':
			s.consume()
			return &ast.Token{Kind: ast.Lparen, Text: "(", Position: pos}
		case r == ')':
			s.consume()
			return &ast.Token{Kind: ast.Rparen, Text: ")", Position: pos}
		case r == '{':
			s.consume()
			return &ast.Token{Kind: ast.Lbrace, Text: "{", Position: pos}
		case r == '}':
			s.consume()
			return &ast.Token{Kind: ast.Rbrace, Text: "}", Position: pos}
		case r == '[':
			s.consume()
			return &ast.Token{Kind: ast.Lbracket, Text: "[", Position: pos}
		case r == ']':
			s.consume()
			return &ast.Token{Kind: ast.Rbracket, Text: "]", Position: pos}
		case r == ';':
			s.consume()
			return &ast.Token{Kind: ast.Semicolon, Text: ";", Position: pos}
		case r == ':':
			s.consume()
			return &ast.Token{Kind: ast.Colon, Text: ":", Position: pos}
		case r == ',':
			s.consume()
			return &ast.Token{Kind: ast.Comma, Text: ",", Position: pos}
		case r == '.':
			s.consume()
			return &ast.Token{Kind: ast.Dot, Text: ".", Position: pos}
		case r == '?':
			s.consume()
			return &ast.Token{Kind: ast.Hook, Text: "?", Position: pos}

		case r == '%':
			s.consume()
			r = s.cur.r
			if r == '=' {
				s.consume()
				return &ast.Token{Kind: ast.PercentEq, Text: "%=", Position: pos}
			}
			return &ast.Token{Kind: ast.Percent, Text: "%", Position: pos}

		case r == '^':
			s.consume()
			r = s.cur.r
			if r == '=' {
				s.consume()
				return &ast.Token{Kind: ast.CaretEq, Text: "^=", Position: pos}
			}
			return &ast.Token{Kind: ast.Caret, Text: "^", Position: pos}

		case r == '~':
			s.consume()
			return &ast.Token{Kind: ast.Tilde, Text: "~", Position: pos}

		case r == '=':
			s.consume()
			r = s.cur.r
			if r == '=' {
				s.consume()
				return &ast.Token{Kind: ast.DblEq, Text: "==", Position: pos}
			} else if r == '>' {
				s.consume()
				return &ast.Token{Kind: ast.EqGt, Text: "=>", Position: pos}
			} else {
				return &ast.Token{Kind: ast.Eq, Text: "=", Position: pos}
			}

		case r == '!':
			s.consume()
			r = s.cur.r
			if r == '=' {
				s.consume()
				return &ast.Token{Kind: ast.NotEq, Text: "!=", Position: pos}
			}
			return &ast.Token{Kind: ast.Not, Text: "!", Position: pos}
		case r == '>':
			s.consume()
			r = s.cur.r
			if r == '=' {
				s.consume()
				return &ast.Token{Kind: ast.GtEq, Text: ">=", Position: pos}
			} else if r == '>' {
				s.consume()
				r = s.cur.r
				if r == '=' {
					s.consume()
					return &ast.Token{Kind: ast.DblGtEq, Text: ">>=", Position: pos}
				}
				return &ast.Token{Kind: ast.DblGt, Text: ">>", Position: pos}
			} else {
				return &ast.Token{Kind: ast.Gt, Text: ">", Position: pos}
			}
		case r == '<':
			s.consume()
			r = s.cur.r
			if r == '=' {
				s.consume()
				r = s.cur.r
				if r == '>' {
					s.consume()
					return &ast.Token{Kind: ast.Cmp, Text: "<=>", Position: pos}
				}
				return &ast.Token{Kind: ast.LtEq, Text: "<=", Position: pos}
			} else if r == '<' {
				s.consume()
				r = s.cur.r
				if r == '=' {
					s.consume()
					return &ast.Token{Kind: ast.DblLtEq, Text: "<<=", Position: pos}
				}
				return &ast.Token{Kind: ast.DblLt, Text: "<<", Position: pos}
			} else {
				return &ast.Token{Kind: ast.Lt, Text: "<", Position: pos}
			}

		case r == '|':
			s.consume()
			r = s.cur.r
			if r == '|' {
				s.consume()
				return &ast.Token{Kind: ast.DblPipe, Text: "||", Position: pos}
			} else if r == '=' {
				s.consume()
				return &ast.Token{Kind: ast.PipeEq, Text: "|=", Position: pos}
			} else {
				return &ast.Token{Kind: ast.Pipe, Text: "|", Position: pos}
			}
		case r == '&':
			s.consume()
			r = s.cur.r
			if r == '&' {
				s.consume()
				return &ast.Token{Kind: ast.DblAmp, Text: "&&", Position: pos}
			} else if r == '=' {
				s.consume()
				return &ast.Token{Kind: ast.AmpEq, Text: "&=", Position: pos}
			} else {
				return &ast.Token{Kind: ast.Amp, Text: "&", Position: pos}
			}

		case r == '\'':
			return s.nextStr('\'')

		case r == '"':
			return s.nextStr('"')

		case r == '`':
			return s.nextRawStr()

		case isDigit(r):
			return s.nextNumber()

		case isIdentStart(r):
			return s.nextIdentOrKeyword()

		case r == eof:
			s.isDone = true
			s.doneToken = &ast.Token{Kind: ast.EOF, Text: "", Position: pos}
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

	text := s.Code[begin:s.cur.idx]
	switch text {

	case "_":
		return &ast.Token{Kind: ast.BlankIdent, Text: text, Position: pos}
	case "null":
		return &ast.Token{Kind: ast.Null, Text: text, Position: pos}
	case "true":
		return &ast.Token{Kind: ast.True, Text: text, Position: pos}
	case "false":
		return &ast.Token{Kind: ast.False, Text: text, Position: pos}
	case "if":
		return &ast.Token{Kind: ast.If, Text: text, Position: pos}
	case "else":
		return &ast.Token{Kind: ast.Else, Text: text, Position: pos}
	case "while":
		return &ast.Token{Kind: ast.While, Text: text, Position: pos}
	case "break":
		return &ast.Token{Kind: ast.Break, Text: text, Position: pos}
	case "continue":
		return &ast.Token{Kind: ast.Continue, Text: text, Position: pos}
	case "fn":
		return &ast.Token{Kind: ast.Fn, Text: text, Position: pos}
	case "return":
		return &ast.Token{Kind: ast.Return, Text: text, Position: pos}
	case "const":
		return &ast.Token{Kind: ast.Const, Text: text, Position: pos}
	case "let":
		return &ast.Token{Kind: ast.Let, Text: text, Position: pos}
	case "for":
		return &ast.Token{Kind: ast.For, Text: text, Position: pos}
	case "in":
		return &ast.Token{Kind: ast.In, Text: text, Position: pos}
	case "switch":
		return &ast.Token{Kind: ast.Switch, Text: text, Position: pos}
	case "case":
		return &ast.Token{Kind: ast.Case, Text: text, Position: pos}
	case "default":
		return &ast.Token{Kind: ast.Default, Text: text, Position: pos}
	case "prop":
		return &ast.Token{Kind: ast.Prop, Text: text, Position: pos}
	case "try":
		return &ast.Token{Kind: ast.Try, Text: text, Position: pos}
	case "catch":
		return &ast.Token{Kind: ast.Catch, Text: text, Position: pos}
	case "finally":
		return &ast.Token{Kind: ast.Finally, Text: text, Position: pos}
	case "throw":
		return &ast.Token{Kind: ast.Throw, Text: text, Position: pos}
	case "go":
		return &ast.Token{Kind: ast.Go, Text: text, Position: pos}
	case "module":
		return &ast.Token{Kind: ast.Module, Text: text, Position: pos}
	case "import":
		return &ast.Token{Kind: ast.Import, Text: text, Position: pos}
	case "struct":
		return &ast.Token{Kind: ast.Struct, Text: text, Position: pos}
	case "dict":
		return &ast.Token{Kind: ast.Dict, Text: text, Position: pos}
	case "set":
		return &ast.Token{Kind: ast.Set, Text: text, Position: pos}
	case "this":
		return &ast.Token{Kind: ast.This, Text: text, Position: pos}
	case "has":
		return &ast.Token{Kind: ast.Has, Text: text, Position: pos}

	case "byte", "defer", "goto", "like", "native", "package",
		"priv", "private", "prot", "protected", "pub", "public",
		"rune", "select", "static", "sync", "rsync", "with", "yield":

		// reserve a bunch of keywords just in case
		return &ast.Token{Kind: ast.Reserved, Text: text, Position: pos}

	default:
		return &ast.Token{Kind: ast.Ident, Text: text, Position: pos}
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
			return &ast.Token{Kind: ast.Str, Text: buf.String(), Position: pos}

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

	text := s.Code[begin:end]
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
	return rune(n), nil
}

func (s *Scanner) nextRawStr() *ast.Token {

	pos := s.pos
	s.consume()

	var buf bytes.Buffer

	for {
		r := s.cur.r

		switch {

		case r == '`':
			// end of string
			s.consume()
			return &ast.Token{Kind: ast.Str, Text: buf.String(), Position: pos}

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
		r = s.cur.r

		switch {

		case isDigit(r):
			return s.unexpectedChar(r, s.pos)

		case r == '.' || isExp(r):
			return s.nextFloat(begin, pos)

		case r == 'x':
			return s.nextHexInt(begin, pos)

		default:
			return &ast.Token{Kind: ast.Int, Text: "0", Position: pos}
		}

	} else {
		s.acceptWhile(isDigit)
		r = s.cur.r
		switch {
		case r == '.':
			return s.nextFloat(begin, pos)
		case isExp(r):
			s.consume()
			s.accept(func(r rune) bool { return (r == '+') || (r == '-') })
			t := s.expect(isDigit)
			if t != nil {
				return t
			}
			s.acceptWhile(isDigit)
			return &ast.Token{Kind: ast.Float, Text: s.Code[begin:s.cur.idx], Position: pos}
		}
		return &ast.Token{Kind: ast.Int, Text: s.Code[begin:s.cur.idx], Position: pos}
	}

}

func (s *Scanner) nextHexInt(begin int, pos ast.Pos) *ast.Token {

	s.consume()

	t := s.expect(isHexDigit)
	if t != nil {
		return t
	}
	s.acceptWhile(isHexDigit)

	return &ast.Token{Kind: ast.Int, Text: s.Code[begin:s.cur.idx], Position: pos}
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

	return &ast.Token{Kind: ast.Float, Text: s.Code[begin:s.cur.idx], Position: pos}
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
		s.doneToken = &ast.Token{Kind: ast.UnexpectedEOF, Text: "", Position: pos}
	} else {
		s.doneToken = &ast.Token{Kind: ast.UnexpectedChar, Text: string(r), Position: pos}
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
