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
func NewSet(ev Eval, values []Value) (Set, Error) {

	hashMap := EmptyHashMap()
	for _, v := range values {
		err := hashMap.Put(ev, v, True)
		if err != nil {
			return nil, err
		}
	}

	return &set{hashMap, false}, nil
}

func (s *set) compositeMarker() {}

func (s *set) Type() Type { return SetType }

func (s *set) Freeze(ev Eval) (Value, Error) {
	s.frozen = true
	return s, nil
}

func (s *set) Frozen(ev Eval) (Bool, Error) {
	return NewBool(s.frozen), nil
}

func (s *set) ToStr(ev Eval) (Str, Error) {

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

		s, err := entry.Key.ToStr(ev)
		if err != nil {
			return nil, err
		}

		buf.WriteString(s.String())
	}

	buf.WriteString(" }")
	return NewStr(buf.String()), nil
}

func (s *set) HashCode(ev Eval) (Int, Error) {
	return nil, HashCodeMismatch(SetType)
}

func (s *set) Eq(ev Eval, v Value) (Bool, Error) {
	switch t := v.(type) {
	case *set:
		return s.hashMap.Eq(ev, t.hashMap)
	default:
		return False, nil
	}
}

func (s *set) Len(ev Eval) (Int, Error) {
	return s.hashMap.Len(), nil
}

//---------------------------------------------------------------

func (s *set) IsEmpty() Bool {
	return NewBool(s.hashMap.Len().IntVal() == 0)
}

func (s *set) Contains(ev Eval, key Value) (Bool, Error) {
	return s.hashMap.Contains(ev, key)
}

func (s *set) Add(ev Eval, val Value) (Set, Error) {
	if s.frozen {
		return nil, ImmutableValue()
	}

	err := s.hashMap.Put(ev, val, True)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (s *set) AddAll(ev Eval, val Value) (Set, Error) {
	if s.frozen {
		return nil, ImmutableValue()
	}

	ibl, ok := val.(Iterable)
	if !ok {
		return nil, IterableMismatch(val.Type())
	}

	itr, err := ibl.NewIterator(ev)
	if err != nil {
		return nil, err
	}

	b, err := itr.IterNext(ev)
	if err != nil {
		return nil, err
	}
	for b.BoolVal() {
		v, err := itr.IterGet(ev)
		if err != nil {
			return nil, err
		}

		err = s.hashMap.Put(ev, v, True)
		if err != nil {
			return nil, err
		}

		b, err = itr.IterNext(ev)
		if err != nil {
			return nil, err
		}
	}
	return s, nil
}

func (s *set) Clear() (Set, Error) {
	if s.frozen {
		return nil, ImmutableValue()
	}

	s.hashMap = EmptyHashMap()
	return s, nil
}

func (s *set) Remove(ev Eval, key Value) (Set, Error) {
	if s.frozen {
		return nil, ImmutableValue()
	}

	_, err := s.hashMap.Remove(ev, key)
	if err != nil {
		return nil, err
	}
	return s, nil
}

//---------------------------------------------------------------
// Iterator

type setIterator struct {
	Struct
	s       *set
	itr     *HIterator
	hasNext bool
}

func (s *set) NewIterator(ev Eval) (Iterator, Error) {

	itr := &setIterator{iteratorStruct(), s, s.hashMap.Iterator(), false}

	next, get := iteratorFields(ev, itr)
	itr.Internal("next", next)
	itr.Internal("get", get)

	return itr, nil
}

func (i *setIterator) IterNext(ev Eval) (Bool, Error) {
	i.hasNext = i.itr.Next()
	return NewBool(i.hasNext), nil
}

func (i *setIterator) IterGet(ev Eval) (Value, Error) {

	if i.hasNext {
		entry := i.itr.Get()
		return entry.Key, nil
	}
	return nil, NoSuchElement()
}

//--------------------------------------------------------------
// fields

var setMethods = map[string]Method{

	"isEmpty": NewNullaryMethod(
		func(self interface{}, ev Eval) (Value, Error) {
			s := self.(Set)
			return s.IsEmpty(), nil
		}),

	"contains": NewFixedMethod(
		[]Type{AnyType}, true,
		func(self interface{}, ev Eval, params []Value) (Value, Error) {
			s := self.(Set)
			return s.Contains(ev, params[0])
		}),

	"add": NewFixedMethod(
		[]Type{AnyType}, true,
		func(self interface{}, ev Eval, params []Value) (Value, Error) {
			s := self.(Set)
			return s.Add(ev, params[0])
		}),

	"addAll": NewFixedMethod(
		[]Type{AnyType}, false,
		func(self interface{}, ev Eval, params []Value) (Value, Error) {
			s := self.(Set)
			return s.AddAll(ev, params[0])
		}),

	"clear": NewNullaryMethod(
		func(self interface{}, ev Eval) (Value, Error) {
			s := self.(Set)
			return s.Clear()
		}),

	"remove": NewFixedMethod(
		[]Type{AnyType}, false,
		func(self interface{}, ev Eval, params []Value) (Value, Error) {
			s := self.(Set)
			return s.Remove(ev, params[0])
		}),
}

func (s *set) FieldNames() ([]string, Error) {
	names := make([]string, 0, len(setMethods))
	for name := range setMethods {
		names = append(names, name)
	}
	return names, nil
}

func (s *set) HasField(name string) (bool, Error) {
	_, ok := setMethods[name]
	return ok, nil
}

func (s *set) GetField(ev Eval, name string) (Value, Error) {
	if method, ok := setMethods[name]; ok {
		return method.ToFunc(s, name), nil
	}
	return nil, NoSuchField(name)
}

func (s *set) InvokeField(ev Eval, name string, params []Value) (Value, Error) {

	if method, ok := setMethods[name]; ok {
		return method.Invoke(s, ev, params)
	}
	return nil, NoSuchField(name)
}
