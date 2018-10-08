// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package core

import (
	"fmt"
)

type Stream interface {

	// TODO implement all of these

	// transformers
	//DropWhile(Predicate) Error
	Filter(Predicate) Error
	//Flatten(Flattener) Error
	//Limit(Int) Error
	Map(Mapper) Error
	//Peek(Consumer) Error
	//Skip(Int) Error
	//TakeWhile(Predicate) Error

	// collectors
	//AllMatch(Eval, Predicate) (Bool, Error)
	//AnyMatch(Eval, Predicate) (Bool, Error)
	//Count(Eval) (Int, Error)
	//ForEach(Eval, Consumer) Error
	//Max(Eval, Lesser) (Value, Error)
	//Min(Eval, Lesser) (Value, Error)
	Reduce(Eval, Value, Reducer) (Value, Error)
	//ToDict(Eval) (Dict, Error)
	ToList(Eval) (List, Error)
	//ToSet(Eval) (Set, Error)
	//ToStruct(Eval) (Struct, Error)
	//ToTuple(Eval) (Tuple, Error)
}

type stream struct {
	adv       advancer
	collected bool
	this      Struct
}

func NewStream(itr Iterator) (Stream, Error) {
	return &stream{&iteratorAdvancer{itr}, false, nil}, nil
}

//--------------------------------------------------------------
// advancers
//--------------------------------------------------------------

type (
	advancer interface {
		advance(Eval) (Value, Error)
	}

	iteratorAdvancer struct {
		itr Iterator
	}

	filterAdvancer struct {
		base advancer
		pred Predicate
	}

	mapAdvancer struct {
		base   advancer
		mapper Mapper
	}
)

func (a *iteratorAdvancer) advance(ev Eval) (Value, Error) {

	b, err := a.itr.IterNext(ev)
	if err != nil {
		return nil, err
	}

	if !b.BoolVal() {
		return nil, nil
	}

	val, err := a.itr.IterGet(ev)
	if err != nil {
		return nil, err
	}

	return val, nil
}

func (a *filterAdvancer) advance(ev Eval) (Value, Error) {

	val, err := a.base.advance(ev)
	if val == nil || err != nil {
		return val, err
	}
	b, err := a.pred(ev, val)
	if err != nil {
		return nil, err
	}

	for !b.BoolVal() {
		val, err = a.base.advance(ev)
		if val == nil || err != nil {
			return val, err
		}
		b, err = a.pred(ev, val)
		if err != nil {
			return nil, err
		}
	}

	return val, nil
}

func (a *mapAdvancer) advance(ev Eval) (Value, Error) {

	v1, err := a.base.advance(ev)
	if v1 == nil || err != nil {
		return v1, err
	}
	v2, err := a.mapper(ev, v1)
	if err != nil {
		return nil, err
	}

	return v2, nil
}

//--------------------------------------------------------------
// transformers
//--------------------------------------------------------------

func (s *stream) Filter(pred Predicate) Error {

	if s.collected {
		return fmt.Errorf("stream has already been collected")
	}

	s.adv = &filterAdvancer{s.adv, pred}
	return nil
}

func (s *stream) Map(mapper Mapper) Error {

	if s.collected {
		return fmt.Errorf("stream has already been collected")
	}

	s.adv = &mapAdvancer{s.adv, mapper}
	return nil
}

//--------------------------------------------------------------
// collectors
//--------------------------------------------------------------

func (s *stream) Reduce(ev Eval, initial Value, reducer Reducer) (Value, Error) {

	if s.collected {
		return nil, fmt.Errorf("stream has already been collected")
	}
	s.collected = true

	acc := initial

	//--------------
	val, err := s.adv.advance(ev)
	if err != nil {
		return nil, err
	}
	for val != nil {

		//--------------
		acc, err = reducer(ev, acc, val)
		if err != nil {
			return nil, err
		}
		//--------------

		val, err = s.adv.advance(ev)
		if err != nil {
			return nil, err
		}
	}
	//--------------

	return acc, nil
}

func (s *stream) ToList(ev Eval) (List, Error) {

	if s.collected {
		return nil, fmt.Errorf("stream has already been collected")
	}
	s.collected = true

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
