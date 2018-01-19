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
)

type set struct {
	hashMap *HashMap
	frozen  bool
}

func NewSet(cx Context, values []Value) Set {

	hashMap := EmptyHashMap()
	for _, v := range values {
		hashMap.Put(cx, v, TRUE)
	}

	return &set{hashMap, false}
}

func (s *set) compositeMarker() {}

func (s *set) Type() Type { return TSET }

func (s *set) Freeze() (Value, Error) {
	s.frozen = true
	return s, nil
}

func (s *set) Frozen() (Bool, Error) {
	return MakeBool(s.frozen), nil
}

func (s *set) ToStr(cx Context) Str {

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
		s := entry.Key.ToStr(cx)
		buf.WriteString(s.String())
	}

	buf.WriteString(" }")
	return MakeStr(buf.String())
}

func (s *set) HashCode(cx Context) (Int, Error) {
	return nil, TypeMismatchError("Expected Hashable Type")
}

func (s *set) Eq(cx Context, v Value) (Bool, Error) {
	switch t := v.(type) {
	case *set:
		return s.hashMap.Eq(cx, t.hashMap)
	default:
		return FALSE, nil
	}
}

func (s *set) Cmp(cx Context, v Value) (Int, Error) {
	return nil, TypeMismatchError("Expected Comparable Type")
}

func (s *set) Len() Int {
	return s.hashMap.Len()
}

func (s *set) IsEmpty() Bool {
	return MakeBool(s.hashMap.Len().IntVal() == 0)
}

func (s *set) Contains(cx Context, key Value) (Bool, Error) {
	return s.hashMap.ContainsKey(cx, key)
}

//---------------------------------------------------------------
// Mutation

func (s *set) Add(cx Context, val Value) Error {
	if s.frozen {
		return ImmutableValueError()
	}

	return s.hashMap.Put(cx, val, TRUE)
}

func (s *set) AddAll(cx Context, val Value) Error {
	if s.frozen {
		return ImmutableValueError()
	}

	if ibl, ok := val.(Iterable); ok {
		itr := ibl.NewIterator(cx)
		for itr.IterNext().BoolVal() {
			v, err := itr.IterGet()
			if err != nil {
				return err
			}
			s.hashMap.Put(cx, v, TRUE)
		}
		return nil
	} else {
		return TypeMismatchError("Expected Iterable Type")
	}
}

func (s *set) Clear() Error {
	if s.frozen {
		return ImmutableValueError()
	}

	s.hashMap = EmptyHashMap()
	return nil
}

func (s *set) Remove(cx Context, key Value) (Bool, Error) {
	if s.frozen {
		return nil, ImmutableValueError()
	}

	return s.hashMap.Remove(cx, key)
}

//---------------------------------------------------------------
// Iterator

type setIterator struct {
	Struct
	s       *set
	itr     *HIterator
	hasNext bool
}

func (s *set) NewIterator(cx Context) Iterator {
	return initIteratorStruct(cx,
		&setIterator{newIteratorStruct(), s, s.hashMap.Iterator(), false})
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

func (s *set) GetField(cx Context, key Str) (Value, Error) {
	switch sn := key.String(); sn {

	case "add":
		return &intrinsicFunc{s, sn, &nativeFunc{
			1, 1,
			func(cx Context, values []Value) (Value, Error) {
				err := s.Add(cx, values[0])
				if err != nil {
					return nil, err
				} else {
					return s, nil
				}
			}}}, nil

	case "addAll":
		return &intrinsicFunc{s, sn, &nativeFunc{
			1, 1,
			func(cx Context, values []Value) (Value, Error) {
				err := s.AddAll(cx, values[0])
				if err != nil {
					return nil, err
				} else {
					return s, nil
				}
			}}}, nil

	case "clear":
		return &intrinsicFunc{s, sn, &nativeFunc{
			0, 0,
			func(cx Context, values []Value) (Value, Error) {
				err := s.Clear()
				if err != nil {
					return nil, err
				} else {
					return s, nil
				}
			}}}, nil

	case "isEmpty":
		return &intrinsicFunc{s, sn, &nativeFunc{
			0, 0,
			func(cx Context, values []Value) (Value, Error) {
				return s.IsEmpty(), nil
			}}}, nil

	case "contains":
		return &intrinsicFunc{s, sn, &nativeFunc{
			1, 1,
			func(cx Context, values []Value) (Value, Error) {
				return s.Contains(cx, values[0])
			}}}, nil

	case "remove":
		return &intrinsicFunc{s, sn, &nativeFunc{
			1, 1,
			func(cx Context, values []Value) (Value, Error) {
				return s.Remove(cx, values[0])
			}}}, nil

	default:
		return nil, NoSuchFieldError(key.String())
	}
}
