// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package ncore

import (
	"fmt"
	"reflect"
)

type (
	Field interface {
		Get(Context) (Value, Error)
		Invoke(Context, []Value) (Value, Error)
		Set(Context, Value) Error
	}

	FieldMap interface {
		Names() []string
		Has(string) bool
		Get(string, Context) (Value, Error)
		Invoke(string, Context, []Value) (Value, Error)
		Set(string, Context, Value) Error

		// InternalReplace is a 'secret' internal function that is used
		// by the Interpreter.  Please pretend its not here.
		InternalReplace(string, Field) Error
	}
)

//--------------------------------------------------------------
// Field
//--------------------------------------------------------------

type field struct {
	value Value
}

func NewField(val Value) Field {
	return &field{val}
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

func (f *field) Set(cx Context, val Value) Error {
	f.value = val
	return nil
}

//--------------------------------------------------------------
// Readonly Field
//--------------------------------------------------------------

type readonlyField struct {
	value Value
}

func NewReadonlyField(val Value) Field {
	return &readonlyField{val}
}

func (f *readonlyField) Get(cx Context) (Value, Error) {
	return f.value, nil
}

func (f *readonlyField) Invoke(cx Context, params []Value) (Value, Error) {
	fn, ok := f.value.(Func)
	if !ok {
		return nil, TypeMismatchError("Expected Func")
	}
	return fn.Invoke(cx, params)
}

func (f *readonlyField) Set(cx Context, val Value) Error {
	return fmt.Errorf("ReadonlyField")
}

//--------------------------------------------------------------
// Property
//--------------------------------------------------------------

type property struct {
	get Func
	set Func
}

func NewProperty(get Func, set Func) (Field, Error) {

	if !reflect.DeepEqual(Arity{FixedArity, 0, 0}, get.Arity()) {
		return nil, fmt.Errorf("InvalidGetterArity: %s", get.Arity().String())
	}

	if !reflect.DeepEqual(Arity{FixedArity, 1, 0}, set.Arity()) {
		return nil, fmt.Errorf("InvalidSetterArity: %s", set.Arity().String())
	}

	return &property{get, set}, nil
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

func (p *property) Set(cx Context, val Value) Error {
	_, err := p.set.Invoke(cx, []Value{val})
	return err
}

//--------------------------------------------------------------
// Readonly Property
//--------------------------------------------------------------

type readonlyProperty struct {
	get Func
}

func NewReadonlyProperty(get Func) (Field, Error) {

	if !reflect.DeepEqual(Arity{FixedArity, 0, 0}, get.Arity()) {
		return nil, fmt.Errorf("InvalidGetterArity: %s", get.Arity().String())
	}

	return &readonlyProperty{get}, nil
}

func (p *readonlyProperty) Get(cx Context) (Value, Error) {
	return p.get.Invoke(cx, []Value{})
}

func (p *readonlyProperty) Invoke(cx Context, params []Value) (Value, Error) {

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

func (p *readonlyProperty) Set(cx Context, val Value) Error {
	return fmt.Errorf("ReadonlyField")
}
