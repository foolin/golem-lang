// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package core

import (
	"sync"
)

type (
	// A Method knows how to invoke some native code on behalf
	// of a given 'self' parameter, without having to actually create
	// a NativeFunc to do the invocation.
	Method interface {

		// Invoke the Method
		Invoke(interface{}, Eval, []Value) (Value, Error)

		// Create a NativeFunc that can invoke the Method
		ToFunc(interface{}, string) NativeFunc
	}

	// MethodInvoke defines a func signature used for invoking a Method
	MethodInvoke func(interface{}, Eval, []Value) (Value, Error)

	// NullaryMethodInvoke defines a func signature used for invoking a nullary Method
	NullaryMethodInvoke func(interface{}, Eval) (Value, Error)

	// WrapperMethodInvoke defines a func signature used for invoking a wrapper Method
	WrapperMethodInvoke func(interface{}) Value
)

//--------------------------------------------------------------
// WrapperMethod
//--------------------------------------------------------------

type wrapperMethod struct {
	invoke WrapperMethodInvoke
	mx     sync.Mutex
	fn     NativeFunc
}

// NewWrapperMethod creates a new wrapper Method.
func NewWrapperMethod(wrapper WrapperMethodInvoke) Method {
	return &wrapperMethod{wrapper, sync.Mutex{}, nil}
}

func (m *wrapperMethod) Invoke(self interface{}, ev Eval, params []Value) (Value, Error) {

	Assert(len(params) == 0)
	return m.invoke(self), nil
}

func (m *wrapperMethod) ToFunc(self interface{}, methodName string) NativeFunc {

	m.mx.Lock()
	defer m.mx.Unlock()

	if m.fn == nil {
		m.fn = NewFixedNativeFunc(
			[]Type{}, false,
			func(ev Eval, params []Value) (Value, Error) {
				return m.invoke(self), nil
			})
	}

	return m.fn
}

//--------------------------------------------------------------
// NullaryMethod
//--------------------------------------------------------------

type nullaryMethod struct {
	invoke NullaryMethodInvoke
	mx     sync.Mutex
	fn     NativeFunc
}

// NewNullaryMethod creates a new nullary Method.
func NewNullaryMethod(nullary NullaryMethodInvoke) Method {
	return &nullaryMethod{nullary, sync.Mutex{}, nil}
}

func (m *nullaryMethod) Invoke(self interface{}, ev Eval, params []Value) (Value, Error) {
	Assert(len(params) == 0)
	return m.invoke(self, ev)
}

func (m *nullaryMethod) ToFunc(self interface{}, methodName string) NativeFunc {

	m.mx.Lock()
	defer m.mx.Unlock()

	if m.fn == nil {
		m.fn = NewFixedNativeFunc(
			[]Type{}, false,
			func(ev Eval, params []Value) (Value, Error) {
				return m.invoke(self, ev)
			})

	}

	return m.fn
}

//--------------------------------------------------------------
// embeddable struct for various method implementations
//--------------------------------------------------------------

type method struct {
	arity  Arity
	invoke MethodInvoke
	mx     sync.Mutex
	fn     NativeFunc
}

//--------------------------------------------------------------
// FixedMethod
//--------------------------------------------------------------

type fixedMethod struct {
	*method
	requiredTypes []Type
	allowNull     bool
}

// NewFixedMethod creates a new Method with fixed arity
func NewFixedMethod(
	requiredTypes []Type,
	allowNull bool,
	invoke MethodInvoke) Method {

	arity := Arity{
		Kind:     FixedArity,
		Required: uint16(len(requiredTypes)),
		Optional: 0,
	}

	return &fixedMethod{
		&method{arity, invoke, sync.Mutex{}, nil},
		requiredTypes, allowNull,
	}
}

func (m *fixedMethod) Invoke(self interface{}, ev Eval, params []Value) (Value, Error) {

	err := vetFixedParams(params, m.requiredTypes, m.allowNull)
	if err != nil {
		return nil, err
	}
	return m.invoke(self, ev, params)
}

func (m *fixedMethod) ToFunc(self interface{}, methodName string) NativeFunc {

	m.mx.Lock()
	defer m.mx.Unlock()

	if m.fn == nil {
		m.fn = NewFixedNativeFunc(
			m.requiredTypes,
			m.allowNull,
			func(ev Eval, params []Value) (Value, Error) {
				return m.invoke(self, ev, params)
			})
	}

	return m.fn
}

//--------------------------------------------------------------
// VariadicMethod
//--------------------------------------------------------------

type variadicMethod struct {
	*method
	requiredTypes []Type
	variadicType  Type
	allowNull     bool
}

// NewVariadicMethod creates a new Method with variadic arity
func NewVariadicMethod(
	requiredTypes []Type,
	variadicType Type,
	allowNull bool,
	invoke MethodInvoke) Method {

	arity := Arity{
		Kind:     VariadicArity,
		Required: uint16(len(requiredTypes)),
		Optional: 0,
	}

	return &variadicMethod{
		&method{arity, invoke, sync.Mutex{}, nil},
		requiredTypes, variadicType, allowNull,
	}
}

func (m *variadicMethod) Invoke(self interface{}, ev Eval, params []Value) (Value, Error) {

	err := vetVariadicParams(params, m.requiredTypes, m.variadicType, m.allowNull)
	if err != nil {
		return nil, err
	}
	return m.invoke(self, ev, params)
}

func (m *variadicMethod) ToFunc(self interface{}, methodName string) NativeFunc {

	m.mx.Lock()
	defer m.mx.Unlock()

	if m.fn == nil {
		m.fn = NewVariadicNativeFunc(
			m.requiredTypes,
			m.variadicType,
			m.allowNull,
			func(ev Eval, params []Value) (Value, Error) {
				return m.invoke(self, ev, params)
			})
	}

	return m.fn
}

//--------------------------------------------------------------
// MultipleMethod
//--------------------------------------------------------------

type multipleMethod struct {
	*method
	requiredTypes []Type
	optionalTypes []Type
	allowNull     bool
}

// NewMultipleMethod creates a new Method with multiple arity
func NewMultipleMethod(
	requiredTypes []Type,
	optionalTypes []Type,
	allowNull bool,
	invoke MethodInvoke) Method {

	arity := Arity{
		Kind:     MultipleArity,
		Required: uint16(len(requiredTypes)),
		Optional: uint16(len(optionalTypes)),
	}

	return &multipleMethod{
		&method{arity, invoke, sync.Mutex{}, nil},
		requiredTypes, optionalTypes, allowNull,
	}
}

func (m *multipleMethod) Invoke(self interface{}, ev Eval, params []Value) (Value, Error) {

	err := vetMultipleParams(params, m.requiredTypes, m.optionalTypes, m.allowNull)
	if err != nil {
		return nil, err
	}
	return m.invoke(self, ev, params)
}

func (m *multipleMethod) ToFunc(self interface{}, methodName string) NativeFunc {

	m.mx.Lock()
	defer m.mx.Unlock()

	if m.fn == nil {
		m.fn = NewMultipleNativeFunc(
			m.requiredTypes,
			m.optionalTypes,
			m.allowNull,
			func(ev Eval, params []Value) (Value, Error) {
				return m.invoke(self, ev, params)
			})
	}

	return m.fn
}
