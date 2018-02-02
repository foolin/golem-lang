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

	// If we are already finished, then by convention we return
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
			return &ast.Token{ast.LINE_FEED, "\n", pos}

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
				return &ast.Token{ast.SLASH_EQ, "/=", pos}

			default:
				return &ast.Token{ast.SLASH, "/", pos}
			}

		case r == '+':
			s.consume()
			r = s.cur.r
			if r == '+' {
				s.consume()
				return &ast.Token{ast.DBL_PLUS, "++", pos}
			} else if r == '=' {
				s.consume()
				return &ast.Token{ast.PLUS_EQ, "+=", pos}
			} else {
				return &ast.Token{ast.PLUS, "+", pos}
			}

		case r == '-':
			s.consume()
			r = s.cur.r
			if r == '-' {
				s.consume()
				return &ast.Token{ast.DBL_MINUS, "--", pos}
			} else if r == '=' {
				s.consume()
				return &ast.Token{ast.MINUS_EQ, "-=", pos}
			} else {
				return &ast.Token{ast.MINUS, "-", pos}
			}

		case r == '*':
			s.consume()
			r = s.cur.r
			if r == '=' {
				s.consume()
				return &ast.Token{ast.STAR_EQ, "*=", pos}
			}
			return &ast.Token{ast.STAR, "*", pos}

		case r == '(':
			s.consume()
			return &ast.Token{ast.LPAREN, "(", pos}
		case r == ')':
			s.consume()
			return &ast.Token{ast.RPAREN, ")", pos}
		case r == '{':
			s.consume()
			return &ast.Token{ast.LBRACE, "{", pos}
		case r == '}':
			s.consume()
			return &ast.Token{ast.RBRACE, "}", pos}
		case r == '[':
			s.consume()
			return &ast.Token{ast.LBRACKET, "[", pos}
		case r == ']':
			s.consume()
			return &ast.Token{ast.RBRACKET, "]", pos}
		case r == ';':
			s.consume()
			return &ast.Token{ast.SEMICOLON, ";", pos}
		case r == ':':
			s.consume()
			return &ast.Token{ast.COLON, ":", pos}
		case r == ',':
			s.consume()
			return &ast.Token{ast.COMMA, ",", pos}
		case r == '.':
			s.consume()
			return &ast.Token{ast.DOT, ".", pos}
		case r == '?':
			s.consume()
			return &ast.Token{ast.HOOK, "?", pos}

		case r == '%':
			s.consume()
			r := s.cur.r
			if r == '=' {
				s.consume()
				return &ast.Token{ast.PERCENT_EQ, "%=", pos}
			}
			return &ast.Token{ast.PERCENT, "%", pos}

		case r == '^':
			s.consume()
			r := s.cur.r
			if r == '=' {
				s.consume()
				return &ast.Token{ast.CARET_EQ, "^=", pos}
			}
			return &ast.Token{ast.CARET, "^", pos}

		case r == '~':
			s.consume()
			return &ast.Token{ast.TILDE, "~", pos}

		case r == '=':
			s.consume()
			r := s.cur.r
			if r == '=' {
				s.consume()
				return &ast.Token{ast.DBL_EQ, "==", pos}
			} else if r == '>' {
				s.consume()
				return &ast.Token{ast.EQ_GT, "=>", pos}
			} else {
				return &ast.Token{ast.EQ, "=", pos}
			}

		case r == '!':
			s.consume()
			r := s.cur.r
			if r == '=' {
				s.consume()
				return &ast.Token{ast.NOT_EQ, "!=", pos}
			}
			return &ast.Token{ast.NOT, "!", pos}
		case r == '>':
			s.consume()
			r := s.cur.r
			if r == '=' {
				s.consume()
				return &ast.Token{ast.GT_EQ, ">=", pos}
			} else if r == '>' {
				s.consume()
				r := s.cur.r
				if r == '=' {
					s.consume()
					return &ast.Token{ast.DBL_GT_EQ, ">>=", pos}
				}
				return &ast.Token{ast.DBL_GT, ">>", pos}
			} else {
				return &ast.Token{ast.GT, ">", pos}
			}
		case r == '<':
			s.consume()
			r := s.cur.r
			if r == '=' {
				s.consume()
				r := s.cur.r
				if r == '>' {
					s.consume()
					return &ast.Token{ast.CMP, "<=>", pos}
				}
				return &ast.Token{ast.LT_EQ, "<=", pos}
			} else if r == '<' {
				s.consume()
				r := s.cur.r
				if r == '=' {
					s.consume()
					return &ast.Token{ast.DBL_LT_EQ, "<<=", pos}
				}
				return &ast.Token{ast.DBL_LT, "<<", pos}
			} else {
				return &ast.Token{ast.LT, "<", pos}
			}

		case r == '|':
			s.consume()
			r := s.cur.r
			if r == '|' {
				s.consume()
				return &ast.Token{ast.DBL_PIPE, "||", pos}
			} else if r == '=' {
				s.consume()
				return &ast.Token{ast.PIPE_EQ, "|=", pos}
			} else {
				return &ast.Token{ast.PIPE, "|", pos}
			}
		case r == '&':
			s.consume()
			r := s.cur.r
			if r == '&' {
				s.consume()
				return &ast.Token{ast.DBL_AMP, "&&", pos}
			} else if r == '=' {
				s.consume()
				return &ast.Token{ast.AMP_EQ, "&=", pos}
			} else {
				return &ast.Token{ast.AMP, "&", pos}
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
			s.doneToken = &ast.Token{ast.EOF, "", pos}
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
		return &ast.Token{ast.BLANK_IDENT, text, pos}
	case "null":
		return &ast.Token{ast.NullValue, text, pos}
	case "true":
		return &ast.Token{ast.True, text, pos}
	case "false":
		return &ast.Token{ast.False, text, pos}
	case "if":
		return &ast.Token{ast.IF, text, pos}
	case "else":
		return &ast.Token{ast.ELSE, text, pos}
	case "while":
		return &ast.Token{ast.WHILE, text, pos}
	case "break":
		return &ast.Token{ast.BREAK, text, pos}
	case "continue":
		return &ast.Token{ast.CONTINUE, text, pos}
	case "fn":
		return &ast.Token{ast.FN, text, pos}
	case "return":
		return &ast.Token{ast.RETURN, text, pos}
	case "const":
		return &ast.Token{ast.CONST, text, pos}
	case "let":
		return &ast.Token{ast.LET, text, pos}
	case "for":
		return &ast.Token{ast.FOR, text, pos}
	case "in":
		return &ast.Token{ast.IN, text, pos}
	case "switch":
		return &ast.Token{ast.SWITCH, text, pos}
	case "case":
		return &ast.Token{ast.CASE, text, pos}
	case "default":
		return &ast.Token{ast.DEFAULT, text, pos}
	case "prop":
		return &ast.Token{ast.PROP, text, pos}
	case "try":
		return &ast.Token{ast.TRY, text, pos}
	case "catch":
		return &ast.Token{ast.CATCH, text, pos}
	case "finally":
		return &ast.Token{ast.FINALLY, text, pos}
	case "throw":
		return &ast.Token{ast.THROW, text, pos}
	case "go":
		return &ast.Token{ast.GO, text, pos}
	case "module":
		return &ast.Token{ast.MODULE, text, pos}
	case "import":
		return &ast.Token{ast.IMPORT, text, pos}
	case "struct":
		return &ast.Token{ast.STRUCT, text, pos}
	case "dict":
		return &ast.Token{ast.DICT, text, pos}
	case "set":
		return &ast.Token{ast.SET, text, pos}
	case "this":
		return &ast.Token{ast.THIS, text, pos}
	case "has":
		return &ast.Token{ast.HAS, text, pos}

	case "byte", "defer", "goto", "like", "native", "package",
		"priv", "private", "prot", "protected", "pub", "public",
		"rune", "select", "static", "sync", "rsync", "with", "yield":

		// reserve a bunch of keywords just in case
		return &ast.Token{ast.RESERVED, text, pos}

	default:
		return &ast.Token{ast.IDENT, text, pos}
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
			return &ast.Token{ast.STR, buf.String(), pos}

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
			return &ast.Token{ast.STR, buf.String(), pos}

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
			return &ast.Token{ast.INT, "0", pos}
		}

	} else {
		s.acceptWhile(isDigit)
		r := s.cur.r
		if r == '.' || isExp(r) {
			return s.nextFloat(begin, pos)
		}
		return &ast.Token{ast.INT, s.source[begin:s.cur.idx], pos}
	}

}

func (s *Scanner) nextHexInt(begin int, pos ast.Pos) *ast.Token {

	s.consume()

	t := s.expect(isHexDigit)
	if t != nil {
		return t
	}
	s.acceptWhile(isHexDigit)

	return &ast.Token{ast.INT, s.source[begin:s.cur.idx], pos}
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

	return &ast.Token{ast.FLOAT, s.source[begin:s.cur.idx], pos}
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
		s.doneToken = &ast.Token{ast.UNEXPECTED_EOF, "", pos}
	} else {
		s.doneToken = &ast.Token{ast.UNEXPECTED_CHAR, string(r), pos}
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
