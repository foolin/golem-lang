// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package ncore

import (
	"reflect"
)

type (
	Field interface {
		Get(Context) (Value, Error)
		Invoke(Context, []Value) (Value, Error)
		Set(Context, Value) (bool, Error)
	}

	Fields interface {
		Names() []string
		Has(string) bool
		Get(string, Context) (Value, Error)
		Invoke(string, Context, []Value) (Value, Error)
		Set(string, Context, Value) Error

		// InternalReplace is a 'secret' internal function.
		// Please pretend its not here.
		InternalReplace(string, Field) Error
	}
)

//--------------------------------------------------------------
// Value Field
//--------------------------------------------------------------

type field struct {
	value    Value
	readonly bool
}

func NewField(val Value) Field {
	return &field{val, false}
}

func NewReadonlyField(val Value) Field {
	return &field{val, true}
}

func (f *field) Get(cx Context) (Value, Error) {
	return f.value, nil
}

func (f *field) Invoke(cx Context, params []Value) (Value, Error) {
	fn, ok := f.value.(Func)
	if !ok {
		return nil, TypeMismatchError("Expected Func")
	}
	return fn.Invoke(cx, params)
}

func (f *field) Set(cx Context, val Value) (bool, Error) {
	if f.readonly {
		return false, nil
	}
	f.value = val
	return true, nil
}

//--------------------------------------------------------------
// Property
//--------------------------------------------------------------

type property struct {
	get      Func
	set      Func
	readonly bool
}

func NewProperty(get Func, set Func) (Field, Error) {

	if !reflect.DeepEqual(Arity{FixedArity, 0, 0}, get.Arity()) {
		return nil, ArityMismatchError(
			"FixedArity(0)",
			int(get.Arity().RequiredParams))
	}

	if !reflect.DeepEqual(Arity{FixedArity, 1, 0}, set.Arity()) {
		return nil, ArityMismatchError(
			"FixedArity(1)",
			int(get.Arity().RequiredParams))
	}

	return &property{get, set, false}, nil
}

func NewReadonlyProperty(get Func) (Field, Error) {

	if !reflect.DeepEqual(Arity{FixedArity, 0, 0}, get.Arity()) {
		return nil, ArityMismatchError(
			"FixedArity(0)",
			int(get.Arity().RequiredParams))
	}

	return &property{get, nil, false}, nil
}

func (p *property) Get(cx Context) (Value, Error) {
	return p.get.Invoke(cx, []Value{})
}

func (p *property) Invoke(cx Context, params []Value) (Value, Error) {

	val, err := p.get.Invoke(cx, []Value{})
	if err != nil {
		return nil, err
	}

	fn, ok := val.(Func)
	if !ok {
		return nil, TypeMismatchError("Expected Func")
	}
	return fn.Invoke(cx, params)
}

func (p *property) Set(cx Context, val Value) (bool, Error) {
	if p.readonly {
		return false, nil
	}
	_, err := p.set.Invoke(cx, []Value{val})
	return true, err
}
