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

func MakeStr(s string) Str {
	return str(s)
}

func (s str) basicMarker() {}

func (s str) TypeOf() Type { return TSTR }

func (s str) ToStr() Str { return s }

func (s str) HashCode() (Int, Error) {
	h := strHash(string(s))
	return MakeInt(int64(h)), nil
}

func (s str) Eq(v Value) Bool {
	switch t := v.(type) {

	case str:
		return MakeBool(s == t)

	default:
		return FALSE
	}
}

func (s str) GetField(key Str) (Value, Error) {
	return nil, NoSuchFieldError(key.String())
}

func (s str) Cmp(v Value) (Int, Error) {
	switch t := v.(type) {

	case str:
		cmp := strings.Compare(string(s), string(t))
		return MakeInt(int64(cmp)), nil

	default:
		return nil, TypeMismatchError("Expected Comparable Type")
	}
}

func (s str) Plus(v Value) (Value, Error) {
	return strcat(s, v), nil
}

func (s str) Get(index Value) (Value, Error) {
	// TODO implement this more efficiently
	runes := []rune(string(s))

	idx, err := validateIndex(index, len(runes))
	if err != nil {
		return nil, err
	}

	return str(string(runes[idx.IntVal()])), nil
}

func (s str) Len() Int {
	n := utf8.RuneCountInString(string(s))
	return MakeInt(int64(n))
}

func (s str) Slice(from Value, to Value) (Value, Error) {
	runes := []rune(string(s))

	f, err := validateIndex(from, len(runes))
	if err != nil {
		return nil, err
	}

	t, err := validateIndex(to, len(runes)+1)
	if err != nil {
		return nil, err
	}

	if t.IntVal() < f.IntVal() {
		return nil, IndexOutOfBoundsError()
	}

	return str(string(runes[f.IntVal():t.IntVal()])), nil
}

func (s str) SliceFrom(from Value) (Value, Error) {
	runes := []rune(string(s))

	f, err := validateIndex(from, len(runes))
	if err != nil {
		return nil, err
	}

	return str(string(runes[f.IntVal():])), nil
}

func (s str) SliceTo(to Value) (Value, Error) {
	runes := []rune(string(s))

	t, err := validateIndex(to, len(runes)+1)
	if err != nil {
		return nil, err
	}

	return str(string(runes[:t.IntVal()])), nil
}

//---------------------------------------------------------------
// Iterator

type strIterator struct {
	Struct
	runes []rune
	n     int
}

func (s str) NewIterator() Iterator {

	stc, err := NewStruct([]*StructEntry{
		{"nextValue", true, false, NULL},
		{"getValue", true, false, NULL}})
	if err != nil {
		panic("invalid struct")
	}

	itr := &strIterator{stc, []rune(string(s)), -1}

	// TODO make the struct immutable once we have set the functions
	stc.InitField(MakeStr("nextValue"), &nativeFunc{
		func(values []Value) (Value, Error) {
			return itr.IterNext(), nil
		}})
	stc.InitField(MakeStr("getValue"), &nativeFunc{
		func(values []Value) (Value, Error) {
			return itr.IterGet()
		}})

	return itr
}

func (i *strIterator) IterNext() Bool {
	i.n++
	return MakeBool(i.n < len(i.runes))
}

func (i *strIterator) IterGet() (Value, Error) {

	if (i.n >= 0) && (i.n < len(i.runes)) {
		return str([]rune{i.runes[i.n]}), nil
	} else {
		return nil, NoSuchElementError()
	}
}
