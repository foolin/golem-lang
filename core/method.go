// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package core

import (
	"fmt"
)

type (
	// A Method knows how to invoke some native code on behalf
	// of a given 'self' parameter, without having to actually create
	// a NativeFunc to do the invocation.
	Method interface {
		Invoke(interface{}, Eval, []Value) (Value, Error)
		ToFunc(interface{}, string) NativeFunc
	}

	MethodInvoke        func(interface{}, Eval, []Value) (Value, Error)
	NullaryMethodInvoke func(interface{}, Eval) (Value, Error)
	WrapperMethodInvoke func(interface{}) Value
)

//--------------------------------------------------------------
// WrapperMethod
//--------------------------------------------------------------

type wrapperMethod struct {
	invoke WrapperMethodInvoke
}

// NewWrapperMethod creates a new wrapper Method.
func NewWrapperMethod(wrapper WrapperMethodInvoke) Method {
	return &wrapperMethod{wrapper}
}

func (m *wrapperMethod) Invoke(self interface{}, ev Eval, params []Value) (Value, Error) {

	fmt.Printf("Wrapper Method Invoke\n")

	Assert(len(params) == 0)
	return m.invoke(self), nil
}

func (m *wrapperMethod) ToFunc(self interface{}, methodName string) NativeFunc {

	fmt.Printf("Wrapper Method ToFunc\n")

	nf := NewFixedNativeFunc(
		[]Type{}, false,
		func(ev Eval, params []Value) (Value, Error) {
			return m.invoke(self), nil
		})

	return &fixedMethodFunc{
		self,
		methodName,
		nf.(*nativeFixedFunc)}
}

//--------------------------------------------------------------
// NullaryMethod
//--------------------------------------------------------------

type nullaryMethod struct {
	invoke NullaryMethodInvoke
}

// NewNullaryMethod creates a new nullary Method.
func NewNullaryMethod(nullary NullaryMethodInvoke) Method {
	return &nullaryMethod{nullary}
}

func (m *nullaryMethod) Invoke(self interface{}, ev Eval, params []Value) (Value, Error) {

	fmt.Printf("Nullary Method Invoke\n")

	Assert(len(params) == 0)
	return m.invoke(self, ev)
}

func (m *nullaryMethod) ToFunc(self interface{}, methodName string) NativeFunc {

	fmt.Printf("Nullary Method ToFunc\n")

	nf := NewFixedNativeFunc(
		[]Type{}, false,
		func(ev Eval, params []Value) (Value, Error) {
			return m.invoke(self, ev)
		})

	return &fixedMethodFunc{
		self,
		methodName,
		nf.(*nativeFixedFunc)}
}

//--------------------------------------------------------------
// embeddable struct for various method implementations
//--------------------------------------------------------------

type method struct {
	arity  Arity
	invoke MethodInvoke
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
		Kind:           FixedArity,
		RequiredParams: uint16(len(requiredTypes)),
		OptionalParams: 0,
	}

	return &fixedMethod{
		&method{arity, invoke},
		requiredTypes, allowNull,
	}
}

func (m *fixedMethod) Invoke(self interface{}, ev Eval, params []Value) (Value, Error) {

	fmt.Printf("Method Invoke\n")

	err := vetFixedParams(params, m.requiredTypes, m.allowNull)
	if err != nil {
		return nil, err
	}
	return m.invoke(self, ev, params)
}

func (m *fixedMethod) ToFunc(self interface{}, methodName string) NativeFunc {

	fmt.Printf("Method ToFunc\n")

	nf := NewFixedNativeFunc(
		m.requiredTypes,
		m.allowNull,
		func(ev Eval, params []Value) (Value, Error) {
			return m.invoke(self, ev, params)
		})

	return &fixedMethodFunc{
		self,
		methodName,
		nf.(*nativeFixedFunc)}
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
		Kind:           VariadicArity,
		RequiredParams: uint16(len(requiredTypes)),
		OptionalParams: 0,
	}

	return &variadicMethod{
		&method{arity, invoke},
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

	nf := NewVariadicNativeFunc(
		m.requiredTypes,
		m.variadicType,
		m.allowNull,
		func(ev Eval, params []Value) (Value, Error) {
			return m.invoke(self, ev, params)
		})

	return &variadicMethodFunc{
		self,
		methodName,
		nf.(*nativeVariadicFunc)}
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
		Kind:           MultipleArity,
		RequiredParams: uint16(len(requiredTypes)),
		OptionalParams: uint16(len(optionalTypes)),
	}

	return &multipleMethod{
		&method{arity, invoke},
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

	nf := NewMultipleNativeFunc(
		m.requiredTypes,
		m.optionalTypes,
		m.allowNull,
		func(ev Eval, params []Value) (Value, Error) {
			return m.invoke(self, ev, params)
		})

	return &multipleMethodFunc{
		self,
		methodName,
		nf.(*nativeMultipleFunc)}
}

//--------------------------------------------------------------
// methodFunc
//--------------------------------------------------------------

// A methodFunc is a function that is created only
// when we really need to have it. The 'same' methodFunc can end up
// being created more than once, so equality has special semantics.

// For self parameters that are Values, equality is based on Eq(),
// otherwise its based on '=='
func selfEq(ev Eval, this, that interface{}) (Bool, Error) {

	valThis, okThis := this.(Value)
	valThat, okThat := that.(Value)

	switch {
	case okThis && okThat:
		return valThis.Eq(ev, valThat)

	case !okThis && !okThat:
		return NewBool(this == that), nil

	default:
		return False, nil
	}
}

//---------------------------------------------
// fixed

type fixedMethodFunc struct {
	self interface{}
	name string
	*nativeFixedFunc
}

func (f *fixedMethodFunc) Eq(ev Eval, v Value) (Bool, Error) {
	switch t := v.(type) {
	case *fixedMethodFunc:
		// Equality is based on whether the two funcs have the same self,
		// and the same name.
		eq, err := selfEq(ev, f.self, t.self)
		if err != nil {
			return nil, err
		}
		return NewBool(eq.BoolVal() && (f.name == t.name)), nil
	default:
		return False, nil
	}
}

//---------------------------------------------
// variadic

type variadicMethodFunc struct {
	self interface{}
	name string
	*nativeVariadicFunc
}

func (f *variadicMethodFunc) Eq(ev Eval, v Value) (Bool, Error) {
	switch t := v.(type) {
	case *variadicMethodFunc:
		// Equality is based on whether the two funcs have the same self,
		// and the same name.
		eq, err := selfEq(ev, f.self, t.self)
		if err != nil {
			return nil, err
		}
		return NewBool(eq.BoolVal() && (f.name == t.name)), nil
	default:
		return False, nil
	}
}

//---------------------------------------------
// multiple

type multipleMethodFunc struct {
	self interface{}
	name string
	*nativeMultipleFunc
}

func (f *multipleMethodFunc) Eq(ev Eval, v Value) (Bool, Error) {
	switch t := v.(type) {
	case *multipleMethodFunc:
		// Equality is based on whether the two funcs have the same self,
		// and the same name.
		eq, err := selfEq(ev, f.self, t.self)
		if err != nil {
			return nil, err
		}
		return NewBool(eq.BoolVal() && (f.name == t.name)), nil
	default:
		return False, nil
	}
}
