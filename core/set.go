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
func NewSet(ev Evaluator, values []Value) (Set, Error) {

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

func (s *set) Freeze(ev Evaluator) (Value, Error) {
	s.frozen = true
	return s, nil
}

func (s *set) Frozen(ev Evaluator) (Bool, Error) {
	return NewBool(s.frozen), nil
}

func (s *set) ToStr(ev Evaluator) (Str, Error) {

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

func (s *set) HashCode(ev Evaluator) (Int, Error) {
	return nil, HashCodeMismatchError(SetType)
}

func (s *set) Eq(ev Evaluator, v Value) (Bool, Error) {
	switch t := v.(type) {
	case *set:
		return s.hashMap.Eq(ev, t.hashMap)
	default:
		return False, nil
	}
}

func (s *set) Len(ev Evaluator) (Int, Error) {
	return s.hashMap.Len(), nil
}

func (s *set) IsEmpty() Bool {
	return NewBool(s.hashMap.Len().IntVal() == 0)
}

func (s *set) Contains(ev Evaluator, key Value) (Bool, Error) {
	return s.hashMap.ContainsKey(ev, key)
}

//---------------------------------------------------------------
// Mutation

func (s *set) Add(ev Evaluator, val Value) Error {
	if s.frozen {
		return ImmutableValueError()
	}

	return s.hashMap.Put(ev, val, True)
}

func (s *set) AddAll(ev Evaluator, val Value) Error {
	if s.frozen {
		return ImmutableValueError()
	}

	panic("TODO")

	//	if ibl, ok := val.(Iterable); ok {
	//		itr := ibl.NewIterator(ev)
	//		for itr.IterNext(ev).BoolVal() {
	//			v, err := itr.IterGet()
	//			if err != nil {
	//				return err
	//			}
	//			err = s.hashMap.Put(ev, v, True)
	//			if err != nil {
	//				return err
	//			}
	//		}
	//		return nil
	//	}
	//
	//		return nil, TempMismatchError(
	//			fmt.Sprintf("Type %s has no iter()", values[0].Type()))
}

func (s *set) Clear() Error {
	if s.frozen {
		return ImmutableValueError()
	}

	s.hashMap = EmptyHashMap()
	return nil
}

func (s *set) Remove(ev Evaluator, key Value) (Bool, Error) {
	if s.frozen {
		return nil, ImmutableValueError()
	}

	return s.hashMap.Remove(ev, key)
}

//---------------------------------------------------------------
// Iterator

type setIterator struct {
	Struct
	s       *set
	itr     *HIterator
	hasNext bool
}

func (s *set) NewIterator(ev Evaluator) (Iterator, Error) {

	itr := &setIterator{iteratorStruct(), s, s.hashMap.Iterator(), false}

	next, get := iteratorFields(ev, itr)
	itr.Internal("next", next)
	itr.Internal("get", get)

	return itr, nil
}

func (i *setIterator) IterNext(ev Evaluator) (Bool, Error) {
	i.hasNext = i.itr.Next()
	return NewBool(i.hasNext), nil
}

func (i *setIterator) IterGet(ev Evaluator) (Value, Error) {

	if i.hasNext {
		entry := i.itr.Get()
		return entry.Key, nil
	}
	return nil, NoSuchElementError()
}

//--------------------------------------------------------------
// fields

var setMethods = map[string]Method{}

func (s *set) FieldNames() ([]string, Error) {
	names := make([]string, 0, len(setMethods))
	for name, _ := range setMethods {
		names = append(names, name)
	}
	return names, nil
}

func (s *set) HasField(name string) (bool, Error) {
	_, ok := setMethods[name]
	return ok, nil
}

func (s *set) GetField(name string, ev Evaluator) (Value, Error) {
	if method, ok := setMethods[name]; ok {
		return method.ToFunc(s, name), nil
	}
	return nil, NoSuchFieldError(name)
}

func (s *set) InvokeField(name string, ev Evaluator, params []Value) (Value, Error) {

	if method, ok := setMethods[name]; ok {
		return method.Invoke(s, ev, params)
	}
	return nil, NoSuchFieldError(name)
}

////--------------------------------------------------------------
//// intrinsic functions
//
//func (s *set) GetField(ev Evaluator, key Str) (Value, Error) {
//	switch sn := key.String(); sn {
//
//	case "add":
//		return &virtualFunc{s, sn, NewFixedNativeFunc(
//			[]Type{AnyType}, false,
//			func(ev Evaluator, values []Value) (Value, Error) {
//				err := s.Add(ev, values[0])
//				if err != nil {
//					return nil, err
//				}
//				return s, nil
//			})}, nil
//
//	case "addAll":
//		return &virtualFunc{s, sn, NewFixedNativeFunc(
//			[]Type{AnyType}, false,
//			func(ev Evaluator, values []Value) (Value, Error) {
//				err := s.AddAll(ev, values[0])
//				if err != nil {
//					return nil, err
//				}
//				return s, nil
//			})}, nil
//
//	case "clear":
//		return &virtualFunc{s, sn, NewFixedNativeFunc(
//			[]Type{}, false,
//			func(ev Evaluator, values []Value) (Value, Error) {
//				err := s.Clear()
//				if err != nil {
//					return nil, err
//				}
//				return s, nil
//			})}, nil
//
//	case "isEmpty":
//		return &virtualFunc{s, sn, NewFixedNativeFunc(
//			[]Type{}, false,
//			func(ev Evaluator, values []Value) (Value, Error) {
//				return s.IsEmpty(), nil
//			})}, nil
//
//	case "contains":
//		return &virtualFunc{s, sn, NewFixedNativeFunc(
//			[]Type{AnyType}, false,
//			func(ev Evaluator, values []Value) (Value, Error) {
//				return s.Contains(ev, values[0])
//			})}, nil
//
//	case "remove":
//		return &virtualFunc{s, sn, NewFixedNativeFunc(
//			[]Type{AnyType}, false,
//			func(ev Evaluator, values []Value) (Value, Error) {
//				return s.Remove(ev, values[0])
//			})}, nil
//
//	case "iterator":
//		return &virtualFunc{s, sn, NewFixedNativeFunc(
//			[]Type{}, false,
//			func(ev Evaluator, values []Value) (Value, Error) {
//				return s.NewIterator(ev), nil
//			})}, nil
//
//	default:
//		return nil, NoSuchFieldError(key.String())
//	}
//}
