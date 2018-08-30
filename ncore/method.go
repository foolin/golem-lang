// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package ncore

import (
//	"fmt"
//	"reflect"
)

type (
	// A Method knows how to invoke some native code on behalf
	// of a given 'self' parameter, without having to actually create
	// a NativeFunc to do the invocation.
	Method interface {
		Invoke(interface{}, Context, []Value) (Value, Error)
		ToFunc(interface{}) NativeFunc
	}

	method struct {
		arity  Arity
		invoke func(interface{}, Context, []Value) (Value, Error)
	}

	fixedMethod struct {
		*method
		requiredTypes []Type
		allowNull     bool
	}
)

// NewFixedMethod creates a new Method with fixed arity
func NewFixedMethod(
	requiredTypes []Type,
	allowNull bool,
	invoke func(interface{}, Context, []Value) (Value, Error)) Method {

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

func (m *fixedMethod) Invoke(self interface{}, cx Context, params []Value) (Value, Error) {

	err := vetFixedParams(params, m.requiredTypes, m.allowNull)
	if err != nil {
		return nil, err
	}
	return m.invoke(self, cx, params)
}

func (m *fixedMethod) ToFunc(self interface{}) NativeFunc {

	return NewFixedNativeFunc(
		m.requiredTypes,
		m.allowNull,
		func(cx Context, params []Value) (Value, Error) {
			return m.invoke(self, cx, params)
		})
}

//--------------------------------------------------------------
// methodFunc
//--------------------------------------------------------------

// A methodFunc is a function that is created only
// when we really need to have it. The 'same' methodFunc can end up
// being created more than once, so equality is based on whether
// the two funcs have the same owner, and the same name
type methodFunc struct {
	*nativeFunc

	owner Value
	name  string
}

func (f *methodFunc) Eq(cx Context, v Value) (Bool, Error) {
	switch t := v.(type) {
	case *methodFunc:
		ownerEq, err := f.owner.Eq(cx, t.owner)
		if err != nil {
			return nil, err
		}
		return NewBool(ownerEq.BoolVal() && (f.name == t.name)), nil
	default:
		return False, nil
	}
}
