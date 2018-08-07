// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package core

import (
	"bytes"
)

type set struct {
	hashMap *HashMap
	frozen  bool
}

// NewSet creates a new Set
func NewSet(cx Context, values []Value) (Set, Error) {

	hashMap := EmptyHashMap()
	for _, v := range values {
		err := hashMap.Put(cx, v, True)
		if err != nil {
			return nil, err
		}
	}

	return &set{hashMap, false}, nil
}

func (s *set) compositeMarker() {}

func (s *set) Type() Type { return SetType }

func (s *set) Freeze() (Value, Error) {
	s.frozen = true
	return s, nil
}

func (s *set) Frozen() (Bool, Error) {
	return NewBool(s.frozen), nil
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
	return NewStr(buf.String())
}

func (s *set) HashCode(cx Context) (Int, Error) {
	return nil, TypeMismatchError("Expected Hashable Type")
}

func (s *set) Eq(cx Context, v Value) (Bool, Error) {
	switch t := v.(type) {
	case *set:
		return s.hashMap.Eq(cx, t.hashMap)
	default:
		return False, nil
	}
}

func (s *set) Cmp(cx Context, v Value) (Int, Error) {
	return nil, TypeMismatchError("Expected Comparable Type")
}

func (s *set) Len() Int {
	return s.hashMap.Len()
}

func (s *set) IsEmpty() Bool {
	return NewBool(s.hashMap.Len().IntVal() == 0)
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

	return s.hashMap.Put(cx, val, True)
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
			err = s.hashMap.Put(cx, v, True)
			if err != nil {
				return err
			}
		}
		return nil
	}
	return TypeMismatchError("Expected Iterable Type")
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
	return NewBool(i.hasNext)
}

func (i *setIterator) IterGet() (Value, Error) {

	if i.hasNext {
		entry := i.itr.Get()
		return entry.Key, nil
	}
	return nil, NoSuchElementError()
}

//--------------------------------------------------------------
// intrinsic functions

func (s *set) GetField(cx Context, key Str) (Value, Error) {
	switch sn := key.String(); sn {

	case "add":
		return &intrinsicFunc{s, sn, NewNativeFunc(
			1, 1,
			func(cx Context, values []Value) (Value, Error) {
				err := s.Add(cx, values[0])
				if err != nil {
					return nil, err
				}
				return s, nil
			})}, nil

	case "addAll":
		return &intrinsicFunc{s, sn, NewNativeFunc(
			1, 1,
			func(cx Context, values []Value) (Value, Error) {
				err := s.AddAll(cx, values[0])
				if err != nil {
					return nil, err
				}
				return s, nil
			})}, nil

	case "clear":
		return &intrinsicFunc{s, sn, NewNativeFunc(
			0, 0,
			func(cx Context, values []Value) (Value, Error) {
				err := s.Clear()
				if err != nil {
					return nil, err
				}
				return s, nil
			})}, nil

	case "isEmpty":
		return &intrinsicFunc{s, sn, NewNativeFunc(
			0, 0,
			func(cx Context, values []Value) (Value, Error) {
				return s.IsEmpty(), nil
			})}, nil

	case "contains":
		return &intrinsicFunc{s, sn, NewNativeFunc(
			1, 1,
			func(cx Context, values []Value) (Value, Error) {
				return s.Contains(cx, values[0])
			})}, nil

	case "remove":
		return &intrinsicFunc{s, sn, NewNativeFunc(
			1, 1,
			func(cx Context, values []Value) (Value, Error) {
				return s.Remove(cx, values[0])
			})}, nil

	default:
		return nil, NoSuchFieldError(key.String())
	}
}
