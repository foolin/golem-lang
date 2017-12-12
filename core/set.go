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
	"reflect"
)

type set struct {
	hashMap *HashMap
}

func NewSet(values []Value) Set {

	hashMap := EmptyHashMap()
	for _, v := range values {
		hashMap.Put(v, TRUE)
	}

	return &set{hashMap}
}

func (s *set) compositeMarker() {}

func (s *set) TypeOf() Type { return TDICT }

func (s *set) ToStr() Str {

	var buf bytes.Buffer
	buf.WriteString("set {")
	idx := 0
	itr := s.hashMap.Iterator()

	for itr.Next() {
		entry := itr.Get()
		if idx > 0 {
			buf.WriteString(",")
		}
		idx++

		buf.WriteString(" ")
		s := entry.Key.ToStr()
		buf.WriteString(s.String())
	}

	buf.WriteString(" }")
	return MakeStr(buf.String())
}

func (s *set) HashCode() (Int, Error) {
	return nil, TypeMismatchError("Expected Hashable Type")
}

func (s *set) Eq(v Value) Bool {
	switch t := v.(type) {
	case *set:
		return MakeBool(reflect.DeepEqual(s.hashMap, t.hashMap))
	default:
		return FALSE
	}
}

func (s *set) Cmp(v Value) (Int, Error) {
	return nil, TypeMismatchError("Expected Comparable Type")
}

func (s *set) Plus(v Value) (Value, Error) {
	switch t := v.(type) {

	case Str:
		return strcat(s, t), nil

	default:
		return nil, TypeMismatchError("Expected Number Type")
	}
}

func (s *set) Len() Int {
	return s.hashMap.Len()
}

func (s *set) Add(val Value) Error {
	return s.hashMap.Put(val, TRUE)
}

func (s *set) AddAll(val Value) Error {
	if ibl, ok := val.(Iterable); ok {
		itr := ibl.NewIterator()
		for itr.IterNext().BoolVal() {
			v, err := itr.IterGet()
			if err != nil {
				return err
			}
			s.hashMap.Put(v, TRUE)
		}
		return nil
	} else {
		return TypeMismatchError("Expected Iterable Type")
	}
}

func (s *set) Clear() {
	s.hashMap = EmptyHashMap()
}

func (s *set) IsEmpty() Bool {
	return MakeBool(s.hashMap.Len().IntVal() == 0)
}

func (s *set) Contains(key Value) (Bool, Error) {
	return s.hashMap.ContainsKey(key)
}

//---------------------------------------------------------------
// Iterator

type setIterator struct {
	Struct
	s       *set
	itr     *HIterator
	hasNext bool
}

func (s *set) NewIterator() Iterator {

	stc, err := NewStruct([]*StructEntry{
		{"nextValue", true, false, NULL},
		{"getValue", true, false, NULL}})
	if err != nil {
		panic("invalid struct")
	}

	itr := &setIterator{stc, s, s.hashMap.Iterator(), false}

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

func (i *setIterator) IterNext() Bool {
	i.hasNext = i.itr.Next()
	return MakeBool(i.hasNext)
}

func (i *setIterator) IterGet() (Value, Error) {

	if i.hasNext {
		entry := i.itr.Get()
		return entry.Key, nil
	} else {
		return nil, NoSuchElementError()
	}
}

//--------------------------------------------------------------
// intrinsic functions

func (s *set) GetField(key Str) (Value, Error) {
	switch key.String() {

	case "add":
		return &intrinsicFunc{s, "add", &nativeFunc{
			func(values []Value) (Value, Error) {
				if len(values) != 1 {
					return nil, ArityMismatchError("1", len(values))
				}
				err := s.Add(values[0])
				if err != nil {
					return nil, err
				} else {
					return s, nil
				}
			}}}, nil

	case "addAll":
		return &intrinsicFunc{s, "addAll", &nativeFunc{
			func(values []Value) (Value, Error) {
				if len(values) != 1 {
					return nil, ArityMismatchError("1", len(values))
				}
				err := s.AddAll(values[0])
				if err != nil {
					return nil, err
				} else {
					return s, nil
				}
			}}}, nil

	case "clear":
		return &intrinsicFunc{s, "clear", &nativeFunc{
			func(values []Value) (Value, Error) {
				if len(values) != 0 {
					return nil, ArityMismatchError("0", len(values))
				}
				s.Clear()
				return s, nil
			}}}, nil

	case "isEmpty":
		return &intrinsicFunc{s, "isEmpty", &nativeFunc{
			func(values []Value) (Value, Error) {
				if len(values) != 0 {
					return nil, ArityMismatchError("0", len(values))
				}
				return s.IsEmpty(), nil
			}}}, nil

	case "contains":
		return &intrinsicFunc{s, "contains", &nativeFunc{
			func(values []Value) (Value, Error) {
				if len(values) != 1 {
					return nil, ArityMismatchError("1", len(values))
				}
				return s.Contains(values[0])
			}}}, nil

	default:
		return nil, NoSuchFieldError(key.String())
	}
}
