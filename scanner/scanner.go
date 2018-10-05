// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package scanner

import (
	"bytes"
	"fmt"
	"github.com/mjarmy/golem-lang/ast"
	"io"
	"strconv"
	"strings"
	"unicode/utf8"
)

//---------------------------------------------------------------
// The Golem Scanner
//---------------------------------------------------------------

const eof rune = -1

type (
	curRune struct {
		r    rune
		size int
		idx  int
	}

	// Source contains the name, path, and source code for a Golem Module
	Source struct {
		Name string
		Path string
		Code string
	}

	// Scanner scans Golem source code and produces a stream of tokens.
	Scanner struct {
		Source *Source

		reader    io.RuneReader
		cur       curRune
		pos       ast.Pos
		isDone    bool
		doneToken *ast.Token
	}
)

// NewScanner creates a new Scanner
func NewScanner(source *Source) (*Scanner, error) {

	if !utf8.ValidString(source.Code) {
		return nil, fmt.Errorf("Source code is not a valid UTF-8-encoded string")
	}

	s := &Scanner{
		Source:    source,
		reader:    strings.NewReader(source.Code),
		cur:       curRune{r: 0, size: 1, idx: -1},
		pos:       ast.Pos{Line: 1, Col: 0},
		isDone:    false,
		doneToken: nil,
	}

	s.consume()
	return s, nil
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
				return &ast.Token{Kind: ast.DoublePlus, Text: "++", Position: pos}
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
				return &ast.Token{Kind: ast.DoubleMinus, Text: "--", Position: pos}
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
			r = s.cur.r
			if r == '.' {

				s.consume()
				r = s.cur.r
				if r == '.' {
					s.consume()
					return &ast.Token{Kind: ast.TripleDot, Text: "...", Position: pos}
					//} else {
					//	// "DoubleDot" isn't actually used for anything right now
					//	return &ast.Token{Kind: ast.DoubleDot, Text: "..", Position: pos}
				}

			} else {
				return &ast.Token{Kind: ast.Dot, Text: ".", Position: pos}
			}

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
				return &ast.Token{Kind: ast.DoubleEq, Text: "==", Position: pos}
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
					return &ast.Token{Kind: ast.DoubleGtEq, Text: ">>=", Position: pos}
				}
				return &ast.Token{Kind: ast.DoubleGt, Text: ">>", Position: pos}
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
					return &ast.Token{Kind: ast.DoubleLtEq, Text: "<<=", Position: pos}
				}
				return &ast.Token{Kind: ast.DoubleLt, Text: "<<", Position: pos}
			} else {
				return &ast.Token{Kind: ast.Lt, Text: "<", Position: pos}
			}

		case r == '|':
			s.consume()
			r = s.cur.r
			if r == '|' {
				s.consume()
				return &ast.Token{Kind: ast.DoublePipe, Text: "||", Position: pos}
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
				return &ast.Token{Kind: ast.DoubleAmp, Text: "&&", Position: pos}
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

		case IsIdentStart(r):
			return s.nextIdentOrKeyword()

		case r == '$':
			return s.nextMagicField()

		case IsIdentStart(r):

		case r == eof:
			s.isDone = true
			s.doneToken = &ast.Token{Kind: ast.EOF, Text: "", Position: pos}
			return s.doneToken

		default:
			return s.unexpectedChar(r, pos)
		}
	}
}

func (s *Scanner) nextMagicField() *ast.Token {

	pos := s.pos
	begin := s.cur.idx
	r := s.cur.r
	s.consume()

	s.acceptWhile(IsIdentContinue)

	text := s.Source.Code[begin:s.cur.idx]

	if !IsMagicField(text) {
		return s.unexpectedChar(r, pos)
	}
	return &ast.Token{Kind: ast.MagicField, Text: text, Position: pos}
}

func (s *Scanner) nextIdentOrKeyword() *ast.Token {

	pos := s.pos
	begin := s.cur.idx
	s.consume()

	s.acceptWhile(IsIdentContinue)

	text := s.Source.Code[begin:s.cur.idx]

	if kind, ok := keywords[text]; ok {
		return &ast.Token{Kind: kind, Text: text, Position: pos}
	}

	if _, ok := reservedWords[text]; ok {
		return &ast.Token{Kind: ast.Reserved, Text: text, Position: pos}
	}

	return &ast.Token{Kind: ast.Ident, Text: text, Position: pos}
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

	text := s.Source.Code[begin:end]
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
			return &ast.Token{Kind: ast.Float, Text: s.Source.Code[begin:s.cur.idx], Position: pos}
		}
		return &ast.Token{Kind: ast.Int, Text: s.Source.Code[begin:s.cur.idx], Position: pos}
	}

}

func (s *Scanner) nextHexInt(begin int, pos ast.Pos) *ast.Token {

	s.consume()

	t := s.expect(isHexDigit)
	if t != nil {
		return t
	}
	s.acceptWhile(isHexDigit)

	return &ast.Token{Kind: ast.Int, Text: s.Source.Code[begin:s.cur.idx], Position: pos}
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

	return &ast.Token{Kind: ast.Float, Text: s.Source.Code[begin:s.cur.idx], Position: pos}
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

func isExp(r rune) bool {
	return (r == 'e') || (r == 'E')
}
