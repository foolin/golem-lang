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

//// NewFixedMethod creates a new Method with fixed arity
//func NewFixedMethod(
//	requiredTypes []Type,
//	allowNull bool,
//	invoke func(interface{}, Context, []Value) (Value, Error)) Method {
//}

//--------------------------------------------------------------
// virtualFunc
//--------------------------------------------------------------

// A virtual function is a function that is created only
// when we really need to have it. The 'same' virtual function can end up
// being created more than once, so equality is based on whether
// the two funcs have the same owner, and the same name
type virtualFunc struct {
	*nativeFunc

	owner Value
	name  string
}

func (f *virtualFunc) Eq(cx Context, v Value) (Bool, Error) {
	switch t := v.(type) {
	case *virtualFunc:
		ownerEq, err := f.owner.Eq(cx, t.owner)
		if err != nil {
			return nil, err
		}
		return NewBool(ownerEq.BoolVal() && (f.name == t.name)), nil
	default:
		return False, nil
	}
}
