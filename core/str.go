// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package core

import (
	//"fmt"
	"strings"
	"unicode/utf8"
)

type str string

func (s str) String() string {
	return string(s)
}

// NewStr creates a new String
func NewStr(s string) Str {
	return str(s)
}

func (s str) basicMarker() {}

func (s str) Type() Type { return StrType }

func (s str) Freeze() (Value, Error) {
	return s, nil
}

func (s str) Frozen() (Bool, Error) {
	return True, nil
}

func (s str) ToStr(cx Context) Str { return s }

func (s str) HashCode(cx Context) (Int, Error) {
	h := strHash(string(s))
	return NewInt(int64(h)), nil
}

func (s str) Eq(cx Context, v Value) (Bool, Error) {
	switch t := v.(type) {

	case str:
		return NewBool(s == t), nil

	default:
		return False, nil
	}
}

func (s str) Cmp(cx Context, v Value) (Int, Error) {
	switch t := v.(type) {

	case str:
		cmp := strings.Compare(string(s), string(t))
		return NewInt(int64(cmp)), nil

	default:
		return nil, TypeMismatchError("Expected Comparable Type")
	}
}

func (s str) Get(cx Context, index Value) (Value, Error) {
	// TODO implement this more efficiently
	runes := []rune(string(s))

	idx, err := boundedIndex(index, len(runes))
	if err != nil {
		return nil, err
	}

	return str(string(runes[idx])), nil
}

func (s str) Len() Int {
	n := utf8.RuneCountInString(string(s))
	return NewInt(int64(n))
}

func (s str) Slice(cx Context, from Value, to Value) (Value, Error) {
	runes := []rune(string(s))

	f, t, err := sliceIndices(from, to, len(runes))
	if err != nil {
		return nil, err
	}

	return str(string(runes[f:t])), nil
}

func (s str) SliceFrom(cx Context, from Value) (Value, Error) {
	runes := []rune(string(s))

	f, _, err := sliceIndices(from, NegOne, len(runes))
	if err != nil {
		return nil, err
	}

	return str(string(runes[f:])), nil
}

func (s str) SliceTo(cx Context, to Value) (Value, Error) {
	runes := []rune(string(s))

	_, t, err := sliceIndices(Zero, to, len(runes))
	if err != nil {
		return nil, err
	}

	return str(string(runes[:t])), nil
}

func (s str) Concat(that Str) Str {
	a := string(s)
	b := string(that.(str))
	return str(strcpy(a) + strcpy(b))
}

// copy to avoid memory leaks
func strcpy(s string) string {
	c := make([]byte, len(s))
	copy(c, s)
	return string(c)
}

//---------------------------------------------------------------
// Iterator

type strIterator struct {
	Struct
	runes []rune
	n     int
}

func (s str) NewIterator(cx Context) Iterator {
	return initIteratorStruct(cx,
		&strIterator{newIteratorStruct(), []rune(string(s)), -1})
}

func (i *strIterator) IterNext() Bool {
	i.n++
	return NewBool(i.n < len(i.runes))
}

func (i *strIterator) IterGet() (Value, Error) {

	if (i.n >= 0) && (i.n < len(i.runes)) {
		return str([]rune{i.runes[i.n]}), nil
	}
	return nil, NoSuchElementError()
}

//--------------------------------------------------------------
// intrinsic functions

func (s str) GetField(cx Context, key Str) (Value, Error) {
	switch sn := key.String(); sn {

	case "contains":
		return &intrinsicFunc{s, sn, NewNativeFunc(
			1, 1,
			func(cx Context, values []Value) (Value, Error) {
				z, ok := values[0].(str)
				if !ok {
					return nil, TypeMismatchError("Expected Str")
				}
				return NewBool(strings.Contains(string(s), string(z))), nil
			})}, nil

	case "index":
		return &intrinsicFunc{s, sn, NewNativeFunc(
			1, 1,
			func(cx Context, values []Value) (Value, Error) {
				z, ok := values[0].(str)
				if !ok {
					return nil, TypeMismatchError("Expected Str")
				}
				return NewInt(int64(strings.Index(string(s), string(z)))), nil
			})}, nil

	case "lastIndex":
		return &intrinsicFunc{s, sn, NewNativeFunc(
			1, 1,
			func(cx Context, values []Value) (Value, Error) {
				z, ok := values[0].(str)
				if !ok {
					return nil, TypeMismatchError("Expected Str")
				}
				return NewInt(int64(strings.LastIndex(string(s), string(z)))), nil
			})}, nil

	case "startsWith":
		return &intrinsicFunc{s, sn, NewNativeFunc(
			1, 1,
			func(cx Context, values []Value) (Value, Error) {
				z, ok := values[0].(str)
				if !ok {
					return nil, TypeMismatchError("Expected Str")
				}
				return NewBool(strings.HasPrefix(string(s), string(z))), nil
			})}, nil

	case "endsWith":
		return &intrinsicFunc{s, sn, NewNativeFunc(
			1, 1,
			func(cx Context, values []Value) (Value, Error) {
				z, ok := values[0].(str)
				if !ok {
					return nil, TypeMismatchError("Expected Str")
				}
				return NewBool(strings.HasSuffix(string(s), string(z))), nil
			})}, nil

	case "replace":
		return &intrinsicFunc{s, sn, NewNativeFunc(
			2, 3,
			func(cx Context, values []Value) (Value, Error) {
				a, ok := values[0].(str)
				if !ok {
					return nil, TypeMismatchError("Expected Str")
				}
				b, ok := values[1].(str)
				if !ok {
					return nil, TypeMismatchError("Expected Str")
				}
				n := -1
				if len(values) == 3 {
					z, ok := values[2].(_int)
					if !ok {
						return nil, TypeMismatchError("Expected Int")
					}
					n = int(z)
				}
				return NewStr(strings.Replace(string(s), string(a), string(b), n)), nil
			})}, nil

	case "split":
		return &intrinsicFunc{s, sn, NewNativeFunc(
			1, 1,
			func(cx Context, values []Value) (Value, Error) {
				z, ok := values[0].(str)
				if !ok {
					return nil, TypeMismatchError("Expected Str")
				}
				tokens := strings.Split(string(s), string(z))
				result := make([]Value, len(tokens))
				for i, t := range tokens {
					result[i] = NewStr(t)
				}
				return NewList(result), nil

			})}, nil

		//func Split(s, sep string) []string
	default:
		return nil, NoSuchFieldError(key.String())
	}
}
