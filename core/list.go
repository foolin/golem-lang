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
	"bytes"
	"strings"
)

//---------------------------------------------------------------
// list

type list struct {
	array []Value
}

func NewList(values []Value) List {
	return &list{values}
}

func (ls *list) compositeMarker() {}

func (ls *list) TypeOf() Type { return TLIST }

func (ls *list) ToStr() Str {

	var buf bytes.Buffer
	buf.WriteString("[")
	for idx, v := range ls.array {
		if idx > 0 {
			buf.WriteString(",")
		}
		buf.WriteString(" ")
		buf.WriteString(v.ToStr().String())
	}
	buf.WriteString(" ]")
	return MakeStr(buf.String())
}

func (ls *list) HashCode() (Int, Error) {
	return nil, TypeMismatchError("Expected Hashable Type")
}

func (ls *list) Eq(v Value) Bool {
	switch t := v.(type) {
	case *list:
		return valuesEq(ls.array, t.array)
	default:
		return FALSE
	}
}

func (ls *list) Cmp(v Value) (Int, Error) {
	return nil, TypeMismatchError("Expected Comparable Type")
}

func (ls *list) Plus(v Value) (Value, Error) {
	switch t := v.(type) {

	case Str:
		return strcat(ls, t), nil

	default:
		return nil, TypeMismatchError("Expected Number Type")
	}
}

func (ls *list) Get(index Value) (Value, Error) {
	idx, err := validateIndex(index, len(ls.array))
	if err != nil {
		return nil, err
	}
	return ls.array[idx.IntVal()], nil
}

func (ls *list) Set(index Value, val Value) Error {
	idx, err := validateIndex(index, len(ls.array))
	if err != nil {
		return err
	}

	ls.array[idx.IntVal()] = val
	return nil
}

func (ls *list) Add(val Value) Error {
	ls.array = append(ls.array, val)
	return nil
}

func (ls *list) AddAll(val Value) Error {
	if ibl, ok := val.(Iterable); ok {
		itr := ibl.NewIterator()
		for itr.IterNext().BoolVal() {
			v, err := itr.IterGet()
			if err != nil {
				return err
			}
			ls.array = append(ls.array, v)
		}
		return nil
	} else {
		return TypeMismatchError("Expected Iterable Type")
	}
}

func (ls *list) Contains(val Value) (Bool, Error) {
	return MakeBool(!ls.IndexOf(val).Eq(NEG_ONE).BoolVal()), nil
}

func (ls *list) IndexOf(val Value) Int {
	for i, v := range ls.array {
		if val.Eq(v).BoolVal() {
			return MakeInt(int64(i))
		}
	}
	return NEG_ONE
}

func (ls *list) Clear() {
	ls.array = []Value{}
}

func (ls *list) IsEmpty() Bool {
	return MakeBool(len(ls.array) == 0)
}

func (ls *list) Join(delim Str) Str {

	s := make([]string, len(ls.array), len(ls.array))
	for i, v := range ls.array {
		s[i] = v.ToStr().String()
	}

	return MakeStr(strings.Join(s, delim.ToStr().String()))
}

func (ls *list) Len() Int {
	return MakeInt(int64(len(ls.array)))
}

func (ls *list) Slice(from Value, to Value) (Value, Error) {

	f, err := validateIndex(from, len(ls.array))
	if err != nil {
		return nil, err
	}

	t, err := validateIndex(to, len(ls.array)+1)
	if err != nil {
		return nil, err
	}

	// TODO do we want a different error here?
	if t.IntVal() < f.IntVal() {
		return nil, IndexOutOfBoundsError()
	}

	a := ls.array[f.IntVal():t.IntVal()]
	b := make([]Value, len(a))
	copy(b, a)
	return NewList(b), nil
}

func (ls *list) SliceFrom(from Value) (Value, Error) {
	return ls.Slice(from, MakeInt(int64(len(ls.array))))
}

func (ls *list) SliceTo(to Value) (Value, Error) {
	return ls.Slice(ZERO, to)
}

func (ls *list) Values() []Value {
	return ls.array
}

//---------------------------------------------------------------
// Iterator

type listIterator struct {
	Struct
	ls *list
	n  int
}

func (ls *list) NewIterator() Iterator {

	stc, err := NewStruct([]*StructEntry{
		{"nextValue", true, false, NULL},
		{"getValue", true, false, NULL}})
	if err != nil {
		panic("invalid struct")
	}

	itr := &listIterator{stc, ls, -1}

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

func (i *listIterator) IterNext() Bool {
	i.n++
	return MakeBool(i.n < len(i.ls.array))
}

func (i *listIterator) IterGet() (Value, Error) {
	if (i.n >= 0) && (i.n < len(i.ls.array)) {
		return i.ls.array[i.n], nil
	} else {
		return nil, NoSuchElementError()
	}
}

//--------------------------------------------------------------
// intrinsic functions

func (ls *list) GetField(key Str) (Value, Error) {
	switch key.String() {

	case "add":
		return &intrinsicFunc{ls, "add", &nativeFunc{
			func(values []Value) (Value, Error) {
				if len(values) != 1 {
					return nil, ArityMismatchError("1", len(values))
				}
				err := ls.Add(values[0])
				if err != nil {
					return nil, err
				} else {
					return ls, nil
				}
			}}}, nil

	case "addAll":
		return &intrinsicFunc{ls, "addAll", &nativeFunc{
			func(values []Value) (Value, Error) {
				if len(values) != 1 {
					return nil, ArityMismatchError("1", len(values))
				}
				err := ls.AddAll(values[0])
				if err != nil {
					return nil, err
				} else {
					return ls, nil
				}
			}}}, nil

	case "clear":
		return &intrinsicFunc{ls, "clear", &nativeFunc{
			func(values []Value) (Value, Error) {
				if len(values) != 0 {
					return nil, ArityMismatchError("0", len(values))
				}
				ls.Clear()
				return ls, nil
			}}}, nil

	case "isEmpty":
		return &intrinsicFunc{ls, "isEmpty", &nativeFunc{
			func(values []Value) (Value, Error) {
				if len(values) != 0 {
					return nil, ArityMismatchError("0", len(values))
				}
				return ls.IsEmpty(), nil
			}}}, nil

	case "contains":
		return &intrinsicFunc{ls, "contains", &nativeFunc{
			func(values []Value) (Value, Error) {
				if len(values) != 1 {
					return nil, ArityMismatchError("1", len(values))
				}
				return ls.Contains(values[0])
			}}}, nil

	case "indexOf":
		return &intrinsicFunc{ls, "indexOf", &nativeFunc{
			func(values []Value) (Value, Error) {
				if len(values) != 1 {
					return nil, ArityMismatchError("1", len(values))
				}
				return ls.IndexOf(values[0]), nil
			}}}, nil

	case "join":
		return &intrinsicFunc{ls, "join", &nativeFunc{
			func(values []Value) (Value, Error) {
				var delim Str
				switch len(values) {
				case 0:
					delim = MakeStr("")
				case 1:
					if s, ok := values[0].(Str); ok {
						delim = s
					} else {
						return nil, TypeMismatchError("Expected Str")
					}
				default:
					return nil, ArityMismatchError("0 or 1", len(values))
				}

				return ls.Join(delim), nil
			}}}, nil

	default:
		return nil, NoSuchFieldError(key.String())
	}
}
