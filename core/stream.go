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
		adv advancer
	}

	advancer interface {
		advance(Eval) (Value, Error)
	}
)

func NewStream(ev Eval, ibl Iterable) (Stream, Error) {

	itr, err := ibl.NewIterator(ev)
	if err != nil {
		return nil, err
	}

	return &stream{&iteratorAdvancer{itr}}, nil
}

//--------------------------------------------------------------
// transformers
//--------------------------------------------------------------

type iteratorAdvancer struct {
	itr Iterator
}

func (i *iteratorAdvancer) advance(ev Eval) (Value, Error) {

	b, err := i.itr.IterNext(ev)
	if err != nil {
		return nil, err
	}

	if !b.BoolVal() {
		return nil, nil
	}

	val, err := i.itr.IterGet(ev)
	if err != nil {
		return nil, err
	}

	return val, nil
}

//--------------------------------------------------------------
// terminators
//--------------------------------------------------------------

func (s *stream) ToList(ev Eval) (List, Error) {

	values := []Value{}

	//--------------
	val, err := s.adv.advance(ev)
	if err != nil {
		return nil, err
	}
	for val != nil {
		//--------------

		values = append(values, val)

		//--------------
		val, err = s.adv.advance(ev)
		if err != nil {
			return nil, err
		}
	}
	//--------------

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
