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

func (s str) Freeze(ev Eval) (Value, Error) {
	return s, nil
}

func (s str) Frozen(ev Eval) (Bool, Error) {
	return True, nil
}

func (s str) ToStr(ev Eval) (Str, Error) {
	return s, nil
}

func (s str) HashCode(ev Eval) (Int, Error) {
	h := strHash(string(s))
	return NewInt(int64(h)), nil
}

func (s str) Eq(ev Eval, v Value) (Bool, Error) {
	switch t := v.(type) {

	case str:
		return NewBool(s == t), nil

	default:
		return False, nil
	}
}

func (s str) Cmp(ev Eval, c Comparable) (Int, Error) {
	switch t := c.(type) {

	case str:
		cmp := strings.Compare(string(s), string(t))
		return NewInt(int64(cmp)), nil

	default:
		return nil, ComparableMismatchError(StrType, c.(Value).Type())
	}
}

func (s str) Get(ev Eval, index Value) (Value, Error) {
	// TODO implement this more efficiently
	runes := []rune(string(s))

	idx, err := boundedIndex(index, len(runes))
	if err != nil {
		return nil, err
	}

	return str(string(runes[idx])), nil
}

func (s str) Set(ev Eval, index Value, val Value) Error {
	return ImmutableValueError()
}

func (s str) Len(ev Eval) (Int, Error) {
	n := utf8.RuneCountInString(string(s))
	return NewInt(int64(n)), nil
}

func (s str) Slice(ev Eval, from Value, to Value) (Value, Error) {
	runes := []rune(string(s))

	f, t, err := sliceIndices(from, to, len(runes))
	if err != nil {
		return nil, err
	}

	return str(string(runes[f:t])), nil
}

func (s str) SliceFrom(ev Eval, from Value) (Value, Error) {
	runes := []rune(string(s))

	f, _, err := sliceIndices(from, NegOne, len(runes))
	if err != nil {
		return nil, err
	}

	return str(string(runes[f:])), nil
}

func (s str) SliceTo(ev Eval, to Value) (Value, Error) {
	runes := []rune(string(s))

	_, t, err := sliceIndices(Zero, to, len(runes))
	if err != nil {
		return nil, err
	}

	return str(string(runes[:t])), nil
}

//---------------------------------------------------------------
// Iterator

type strIterator struct {
	Struct
	runes []rune
	n     int
}

func (s str) NewIterator(ev Eval) (Iterator, Error) {

	itr := &strIterator{iteratorStruct(), []rune(string(s)), -1}

	next, get := iteratorFields(ev, itr)
	itr.Internal("next", next)
	itr.Internal("get", get)

	return itr, nil
}

func (i *strIterator) IterNext(ev Eval) (Bool, Error) {
	i.n++
	return NewBool(i.n < len(i.runes)), nil
}

func (i *strIterator) IterGet(ev Eval) (Value, Error) {

	if (i.n >= 0) && (i.n < len(i.runes)) {
		return str([]rune{i.runes[i.n]}), nil
	}
	return nil, NoSuchElementError()
}

//--------------------------------------------------------------

func (s str) Concat(that Str) Str {
	a := string(s)
	b := string(that.(str))
	return str(strcpy(a) + strcpy(b))
}

func strcpy(s string) string {
	c := make([]byte, len(s))
	copy(c, s)
	return string(c)
}

func (s str) Contains(substr Str) Bool {
	a := string(s)
	b := string(substr.(str))
	return NewBool(strings.Contains(a, b))
}

func (s str) Index(substr Str) Int {
	a := string(s)
	b := string(substr.(str))
	result := strings.Index(a, b)
	if result == -1 {
		return NegOne
	}
	result = utf8.RuneCountInString(a[:result])
	return NewInt(int64(result))
}

func (s str) LastIndex(substr Str) Int {
	a := string(s)
	b := string(substr.(str))
	result := strings.LastIndex(a, b)
	if result == -1 {
		return NegOne
	}
	result = utf8.RuneCountInString(a[:result])
	return NewInt(int64(result))
}

func (s str) HasPrefix(substr Str) Bool {
	a := string(s)
	b := string(substr.(str))
	return NewBool(strings.HasPrefix(a, b))
}

func (s str) HasSuffix(substr Str) Bool {
	a := string(s)
	b := string(substr.(str))
	return NewBool(strings.HasSuffix(a, b))
}

func (s str) Replace(old, new Str, n Int) Str {
	a := string(s)
	b := string(old.(str))
	c := string(new.(str))
	d := int(n.(_int))
	return str(strings.Replace(a, b, c, d))
}

func (s str) Split(sep Str) List {
	a := string(s)
	b := string(sep.(str))

	tokens := strings.Split(a, b)
	result := make([]Value, len(tokens))
	for i, t := range tokens {
		result[i] = str(t)
	}
	return NewList(result)
}

//--------------------------------------------------------------
// fields

var strMethods = map[string]Method{

	"contains": NewFixedMethod(
		[]Type{StrType}, false,
		func(self interface{}, ev Eval, params []Value) (Value, Error) {
			return self.(Str).Contains(params[0].(Str)), nil
		}),

	"index": NewFixedMethod(
		[]Type{StrType}, false,
		func(self interface{}, ev Eval, params []Value) (Value, Error) {
			return self.(Str).Index(params[0].(Str)), nil
		}),

	"lastIndex": NewFixedMethod(
		[]Type{StrType}, false,
		func(self interface{}, ev Eval, params []Value) (Value, Error) {
			return self.(Str).LastIndex(params[0].(Str)), nil
		}),

	"hasPrefix": NewFixedMethod(
		[]Type{StrType}, false,
		func(self interface{}, ev Eval, params []Value) (Value, Error) {
			return self.(Str).HasPrefix(params[0].(Str)), nil
		}),

	"hasSuffix": NewFixedMethod(
		[]Type{StrType}, false,
		func(self interface{}, ev Eval, params []Value) (Value, Error) {
			return self.(Str).HasSuffix(params[0].(Str)), nil
		}),

	"replace": NewMultipleMethod(
		[]Type{StrType, StrType},
		[]Type{IntType},
		false,
		func(self interface{}, ev Eval, params []Value) (Value, Error) {
			old := params[0].(Str)
			new := params[1].(Str)
			n := NegOne
			if len(params) == 3 {
				n = params[2].(Int)
			}
			return self.(Str).Replace(old, new, n), nil
		}),

	"split": NewFixedMethod(
		[]Type{StrType}, false,
		func(self interface{}, ev Eval, params []Value) (Value, Error) {
			return self.(Str).Split(params[0].(Str)), nil
		}),
}

func (s str) FieldNames() ([]string, Error) {
	names := make([]string, 0, len(strMethods))
	for name := range strMethods {
		names = append(names, name)
	}
	return names, nil
}

func (s str) HasField(name string) (bool, Error) {
	_, ok := strMethods[name]
	return ok, nil
}

func (s str) GetField(ev Eval, name string) (Value, Error) {
	if method, ok := strMethods[name]; ok {
		return method.ToFunc(s, name), nil
	}
	return nil, NoSuchFieldError(name)
}

func (s str) InvokeField(ev Eval, name string, params []Value) (Value, Error) {
	if method, ok := strMethods[name]; ok {
		return method.Invoke(s, ev, params)
	}
	return nil, NoSuchFieldError(name)
}
