// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package core

import (
	"fmt"
)

type Stream interface {

	// transformers
	//DropWhile(Predicate)
	//Filter(Predicate)
	//Flatten(Flattener)
	//Limit(Int)
	//Map(Mapper)
	//Peek(Consumer)
	//Skip(Int)
	//TakeWhile(Predicate)

	// terminators
	//AllMatch(Eval, Predicate) (Bool, Error)
	//AnyMatch(Eval, Predicate) (Bool, Error)
	//Count(Eval) (Int, Error)
	//ForEach(Eval, Consumer) Error
	//Max(Eval, Lesser) (Value, Error)
	//Min(Eval, Lesser) (Value, Error)
	//Reduce(Eval, Value, Reducer) (Value, Error)
	//ToDict(Eval) (Dict, Error)
	ToList(Eval) (List, Error)
	//ToSet(Eval) (Set, Error)
	//ToStruct(Eval) (Struct, Error)
	//ToTuple(Eval) (Tuple, Error)
}

type (
	stream struct {
		itr Iterator
	}
)

func NewStream(ev Eval, ibl Iterable) (Stream, Error) {

	itr, err := ibl.NewIterator(ev)
	if err != nil {
		return nil, err
	}

	return &stream{itr}, nil
}

//--------------------------------------------------------------
// transformers
//--------------------------------------------------------------

//--------------------------------------------------------------
// terminators
//--------------------------------------------------------------

func (s *stream) terminate(ev Eval, consume func(v Value) Error) Error {

	b, err := s.itr.IterNext(ev)
	if err != nil {
		return err
	}
	for b.BoolVal() {
		v, err := s.itr.IterGet(ev)
		if err != nil {
			return err
		}

		consume(v)

		b, err = s.itr.IterNext(ev)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *stream) ToList(ev Eval) (List, Error) {
	values := []Value{}
	err := s.terminate(ev, func(v Value) Error {
		values = append(values, v)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return NewList(values), nil
}

//--------------------------------------------------------------
// builtin
//--------------------------------------------------------------

var streamMethods = map[string]Method{

	"toList": NewFixedMethod(
		[]Type{}, false,
		func(self interface{}, ev Eval, params []Value) (Value, Error) {
			s := self.(Stream)
			return s.ToList(ev)
		}),
}

// BuiltinStream returns a stream Struct.
var BuiltinStream = NewFixedNativeFunc(
	[]Type{AnyType},
	false,
	func(ev Eval, params []Value) (Value, Error) {
		ibl, ok := params[0].(Iterable)
		if !ok {
			return nil, fmt.Errorf("TypeMismatch: stream() expected iterable value, got %s", params[0].Type())
		}

		s, err := NewStream(ev, ibl)
		if err != nil {
			return nil, err
		}

		return NewMethodStruct(s, streamMethods)
	})
