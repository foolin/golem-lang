// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package core

import (
	"fmt"
	"reflect"
)

type (
	Field interface {
		Get(Eval) (Value, Error)
		Invoke(Eval, []Value) (Value, Error)
		IsReadonly() bool
		Set(Eval, Value) Error
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

func (f *field) Get(ev Eval) (Value, Error) {
	return f.value, nil
}

func (f *field) Invoke(ev Eval, params []Value) (Value, Error) {
	fn, ok := f.value.(Func)
	if !ok {
		return nil, TypeMismatchError(FuncType, f.value.Type())
	}
	return fn.Invoke(ev, params)
}

func (f *field) IsReadonly() bool {
	return false
}

func (f *field) Set(ev Eval, val Value) Error {
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

func (f *readonlyField) Get(ev Eval) (Value, Error) {
	return f.value, nil
}

func (f *readonlyField) Invoke(ev Eval, params []Value) (Value, Error) {
	fn, ok := f.value.(Func)
	if !ok {
		return nil, TypeMismatchError(FuncType, f.value.Type())
	}
	return fn.Invoke(ev, params)
}

func (f *readonlyField) IsReadonly() bool {
	return true
}

func (f *readonlyField) Set(ev Eval, val Value) Error {
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

func (p *property) Get(ev Eval) (Value, Error) {
	return p.get.Invoke(ev, []Value{})
}

func (p *property) Invoke(ev Eval, params []Value) (Value, Error) {

	val, err := p.get.Invoke(ev, []Value{})
	if err != nil {
		return nil, err
	}

	fn, ok := val.(Func)
	if !ok {
		return nil, TypeMismatchError(FuncType, val.Type())
	}
	return fn.Invoke(ev, params)
}

func (p *property) IsReadonly() bool {
	return false
}

func (p *property) Set(ev Eval, val Value) Error {
	_, err := p.set.Invoke(ev, []Value{val})
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

func (p *readonlyProperty) Get(ev Eval) (Value, Error) {
	return p.get.Invoke(ev, []Value{})
}

func (p *readonlyProperty) Invoke(ev Eval, params []Value) (Value, Error) {

	val, err := p.get.Invoke(ev, []Value{})
	if err != nil {
		return nil, err
	}

	fn, ok := val.(Func)
	if !ok {
		return nil, TypeMismatchError(FuncType, val.Type())
	}
	return fn.Invoke(ev, params)
}

func (p *readonlyProperty) IsReadonly() bool {
	return true
}

func (p *readonlyProperty) Set(ev Eval, val Value) Error {
	return fmt.Errorf("ReadonlyField")
}
