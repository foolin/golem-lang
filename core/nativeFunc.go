// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package core

import (
	"fmt"
)

//--------------------------------------------------------------
// NativeFunc

// NativeFunc is a Func that is implemented in GoStmt rather than Golem
type NativeFunc interface {
	Func
}

type nativeFunc struct {
	minArity int
	maxArity int
	invoke   func(Context, []Value) (Value, Error)
}

// NewNativeFunc creates a new NativeFunc
func NewNativeFunc(minArity int, maxArity int, f func(Context, []Value) (Value, Error)) NativeFunc {
	return &nativeFunc{minArity, maxArity, f}
}

func (f *nativeFunc) funcMarker() {}

func (f *nativeFunc) Type() Type { return FuncType }

func (f *nativeFunc) Freeze() (Value, Error) {
	return f, nil
}

func (f *nativeFunc) Frozen() (Bool, Error) {
	return True, nil
}

func (f *nativeFunc) Eq(cx Context, v Value) (Bool, Error) {
	switch t := v.(type) {
	case NativeFunc:
		// equality is based on identity
		return NewBool(f == t), nil
	default:
		return False, nil
	}
}

func (f *nativeFunc) HashCode(cx Context) (Int, Error) {
	return nil, TypeMismatchError("Expected Hashable Type")
}

func (f *nativeFunc) GetField(cx Context, key Str) (Value, Error) {
	return nil, NoSuchFieldError(key.String())
}

func (f *nativeFunc) Cmp(cx Context, v Value) (Int, Error) {
	return nil, TypeMismatchError("Expected Comparable Type")
}

func (f *nativeFunc) ToStr(cx Context) Str {
	return NewStr(fmt.Sprintf("nativeFunc<%p>", f))
}

func (f *nativeFunc) MinArity() int { return f.minArity }
func (f *nativeFunc) MaxArity() int { return f.maxArity }

func (f *nativeFunc) Invoke(cx Context, values []Value) (Value, Error) {

	arity := len(values)
	min := f.MinArity()
	max := f.MaxArity()

	if min == max {
		if arity != min {
			return nil, ArityMismatchError(fmt.Sprintf("%d", min), arity)
		}
	} else {
		if arity < min {
			return nil, ArityMismatchError(fmt.Sprintf("at least %d", min), arity)
		} else if (max != -1) && (arity > max) {
			return nil, ArityMismatchError(fmt.Sprintf("at most %d", max), arity)
		}
	}

	return f.invoke(cx, values)
}

//---------------------------------------------------------------
// An intrinsic function is a function that is an intrinsic
// part of a given Type. These functions are created on the
// fly.

type intrinsicFunc struct {
	owner Value
	name  string
	*nativeFunc
}

func (f *intrinsicFunc) Eq(cx Context, v Value) (Bool, Error) {
	switch t := v.(type) {
	case *intrinsicFunc:
		// equality for intrinsic functions is based on whether
		// they have the same owner, and the same name
		ownerEq, err := f.owner.Eq(cx, t.owner)
		if err != nil {
			return nil, err
		}
		return NewBool(ownerEq.BoolVal() && (f.name == t.name)), nil
	default:
		return False, nil
	}
}
