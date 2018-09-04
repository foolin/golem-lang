// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package core

import (
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

func (s str) Freeze(ev Evaluator) (Value, Error) {
	return s, nil
}

func (s str) Frozen(ev Evaluator) (Bool, Error) {
	return True, nil
}

func (s str) ToStr(ev Evaluator) (Str, Error) {
	return s, nil
}

func (s str) HashCode(ev Evaluator) (Int, Error) {
	h := strHash(string(s))
	return NewInt(int64(h)), nil
}

func (s str) Eq(ev Evaluator, v Value) (Bool, Error) {
	switch t := v.(type) {

	case str:
		return NewBool(s == t), nil

	default:
		return False, nil
	}
}

func (s str) Cmp(ev Evaluator, c Comparable) (Int, Error) {
	switch t := c.(type) {

	case str:
		cmp := strings.Compare(string(s), string(t))
		return NewInt(int64(cmp)), nil

	default:
		return nil, ComparableMismatchError(StrType, c.(Value).Type())
	}
}

func (s str) Get(ev Evaluator, index Value) (Value, Error) {
	// TODO implement this more efficiently
	runes := []rune(string(s))

	idx, err := boundedIndex(index, len(runes))
	if err != nil {
		return nil, err
	}

	return str(string(runes[idx])), nil
}

func (s str) Set(ev Evaluator, index Value, val Value) Error {
	return ImmutableValueError()
}

func (s str) Len(ev Evaluator) (Int, Error) {
	n := utf8.RuneCountInString(string(s))
	return NewInt(int64(n)), nil
}

func (s str) Slice(ev Evaluator, from Value, to Value) (Value, Error) {
	runes := []rune(string(s))

	f, t, err := sliceIndices(from, to, len(runes))
	if err != nil {
		return nil, err
	}

	return str(string(runes[f:t])), nil
}

func (s str) SliceFrom(ev Evaluator, from Value) (Value, Error) {
	runes := []rune(string(s))

	f, _, err := sliceIndices(from, NegOne, len(runes))
	if err != nil {
		return nil, err
	}

	return str(string(runes[f:])), nil
}

func (s str) SliceTo(ev Evaluator, to Value) (Value, Error) {
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

//---------------------------------------------------------------
// Iterator

type strIterator struct {
	Struct
	runes []rune
	n     int
}

func (s str) NewIterator(ev Evaluator) (Iterator, Error) {

	itr := &strIterator{iteratorStruct(), []rune(string(s)), -1}

	next, get := iteratorFields(ev, itr)
	itr.Internal("next", next)
	itr.Internal("get", get)

	return itr, nil
}

func (i *strIterator) IterNext(ev Evaluator) (Bool, Error) {
	i.n++
	return NewBool(i.n < len(i.runes)), nil
}

func (i *strIterator) IterGet(ev Evaluator) (Value, Error) {

	if (i.n >= 0) && (i.n < len(i.runes)) {
		return str([]rune{i.runes[i.n]}), nil
	}
	return nil, NoSuchElementError()
}

//--------------------------------------------------------------
// fields

var strMethods = map[string]Method{

	"contains": NewFixedMethod(
		[]Type{StrType}, false,
		func(self interface{}, ev Evaluator, params []Value) (Value, Error) {

			s := self.(Str).String()
			substr := params[0].(Str).String()
			result := strings.Contains(s, substr)

			return NewBool(result), nil
		}),

	"index": NewFixedMethod(
		[]Type{StrType}, false,
		func(self interface{}, ev Evaluator, params []Value) (Value, Error) {
			s := self.(Str).String()
			substr := params[0].(Str).String()

			result := strings.Index(s, substr)
			if result == -1 {
				return NegOne, nil
			}
			result = utf8.RuneCountInString(s[:result])
			return NewInt(int64(result)), nil
		}),

	"lastIndex": NewFixedMethod(
		[]Type{StrType}, false,
		func(self interface{}, ev Evaluator, params []Value) (Value, Error) {
			s := self.(Str).String()
			substr := params[0].(Str).String()

			result := strings.LastIndex(s, substr)
			if result == -1 {
				return NegOne, nil
			}
			result = utf8.RuneCountInString(s[:result])
			return NewInt(int64(result)), nil
		}),

	"hasPrefix": NewFixedMethod(
		[]Type{StrType}, false,
		func(self interface{}, ev Evaluator, params []Value) (Value, Error) {
			s := self.(Str).String()
			prefix := params[0].(Str).String()
			return NewBool(strings.HasPrefix(s, prefix)), nil
		}),

	"hasSuffix": NewFixedMethod(
		[]Type{StrType}, false,
		func(self interface{}, ev Evaluator, params []Value) (Value, Error) {
			s := self.(Str).String()
			prefix := params[0].(Str).String()
			return NewBool(strings.HasPrefix(s, prefix)), nil
		}),

	"replace": NewMultipleMethod(
		[]Type{StrType, StrType},
		[]Type{IntType},
		false,
		func(self interface{}, ev Evaluator, params []Value) (Value, Error) {
			s := self.(Str).String()
			old := params[0].(Str).String()
			new := params[1].(Str).String()
			var n int = -1
			if len(params) == 3 {
				n = int(params[2].(Int).IntVal())
			}

			return NewStr(strings.Replace(
				s,
				old,
				new,
				n)), nil
		}),

	"split": NewFixedMethod(
		[]Type{StrType}, false,
		func(self interface{}, ev Evaluator, params []Value) (Value, Error) {
			s := self.(Str).String()
			sep := params[0].(Str).String()

			tokens := strings.Split(s, sep)
			result := make([]Value, len(tokens))
			for i, t := range tokens {
				result[i] = NewStr(t)
			}
			return NewList(result), nil
		}),
}

func (s str) FieldNames() ([]string, Error) {
	names := make([]string, 0, len(strMethods))
	for name, _ := range strMethods {
		names = append(names, name)
	}
	return names, nil
}

func (s str) HasField(name string) (bool, Error) {
	_, ok := strMethods[name]
	return ok, nil
}

func (s str) GetField(name string, ev Evaluator) (Value, Error) {
	if method, ok := strMethods[name]; ok {
		return method.ToFunc(s, name), nil
	}
	return nil, NoSuchFieldError(name)
}

func (s str) InvokeField(name string, ev Evaluator, params []Value) (Value, Error) {
	if method, ok := strMethods[name]; ok {
		return method.Invoke(s, ev, params)
	}
	return nil, NoSuchFieldError(name)
}
