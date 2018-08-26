// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package core

import (
	"fmt"
)

// NativeFunc is a Func that is implemented in Go rather than Golem
type NativeFunc interface {
	Func
}

//---------------------------------------------------------------
// An intrinsic function is a function that is an intrinsic
// part of a given Type. These functions are created on the
// fly.

type intrinsicFunc struct {
	owner Value
	name  string
	NativeFunc
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

//--------------------------------------------------------------

type baseFunc struct {
}

func (f *baseFunc) funcMarker() {}

func (f *baseFunc) Type() Type { return FuncType }

func (f *baseFunc) Frozen() (Bool, Error) {
	return True, nil
}

func (f *baseFunc) HashCode(cx Context) (Int, Error) {
	return nil, TypeMismatchError("Expected Hashable Type")
}

func (f *baseFunc) GetField(cx Context, key Str) (Value, Error) {
	return nil, NoSuchFieldError(key.String())
}

func (f *baseFunc) Cmp(cx Context, v Value) (Int, Error) {
	return nil, TypeMismatchError("Expected Comparable Type")
}

func (f *baseFunc) ToStr(cx Context) Str {
	return NewStr(fmt.Sprintf("nativeFunc<%p>", f))
}

//--------------------------------------------------------------

type nativeFunc struct {
	*baseFunc
	minArity int
	maxArity int
	invoke   func(Context, []Value) (Value, Error)
}

// NewNativeFunc creates a new NativeFunc
func NewNativeFunc(
	minArity int,
	maxArity int,
	f func(Context, []Value) (Value, Error)) NativeFunc {

	return &nativeFunc{&baseFunc{}, minArity, maxArity, f}
}

func (f *nativeFunc) Freeze() (Value, Error) {
	return f, nil
}

func (f *nativeFunc) Eq(cx Context, v Value) (Bool, Error) {
	switch t := v.(type) {
	case *nativeFunc:
		// equality is based on identity
		return NewBool(f == t), nil
	default:
		return False, nil
	}
}

func (f *nativeFunc) Arity() *Arity { panic("Arity") }
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

//--------------------------------------------------------------
// nativeFunc0

type nativeFunc0 struct {
	*baseFunc
	invoke func(Context) (Value, Error)
}

// NewNativeFunc0 creates a new NativeFunc that takes 0 parameters.
func NewNativeFunc0(f func(Context) (Value, Error)) NativeFunc {
	return &nativeFunc0{&baseFunc{}, f}
}

func (f *nativeFunc0) Freeze() (Value, Error) {
	return f, nil
}

func (f *nativeFunc0) Eq(cx Context, v Value) (Bool, Error) {
	switch t := v.(type) {
	case *nativeFunc0:
		// equality is based on identity
		return NewBool(f == t), nil
	default:
		return False, nil
	}
}

func (f *nativeFunc0) Arity() *Arity { panic("Arity") }
func (f *nativeFunc0) MinArity() int { return 0 }
func (f *nativeFunc0) MaxArity() int { return 0 }

func (f *nativeFunc0) Invoke(cx Context, values []Value) (Value, Error) {

	arity := len(values)
	if arity != 0 {
		return nil, ArityMismatchError("0", arity)
	}

	return f.invoke(cx)
}

//--------------------------------------------------------------
// nativeFuncValue

type nativeFuncValue struct {
	*baseFunc
	invoke func(Context, Value) (Value, Error)
}

// NewNativeFuncValue creates a new NativeFunc that takes 1 Value parameter.
func NewNativeFuncValue(f func(Context, Value) (Value, Error)) NativeFunc {
	return &nativeFuncValue{&baseFunc{}, f}
}

func (f *nativeFuncValue) Freeze() (Value, Error) {
	return f, nil
}

func (f *nativeFuncValue) Eq(cx Context, v Value) (Bool, Error) {
	switch t := v.(type) {
	case *nativeFuncValue:
		// equality is based on identity
		return NewBool(f == t), nil
	default:
		return False, nil
	}
}

func (f *nativeFuncValue) Arity() *Arity { panic("Arity") }
func (f *nativeFuncValue) MinArity() int { return 1 }
func (f *nativeFuncValue) MaxArity() int { return 1 }

func (f *nativeFuncValue) Invoke(cx Context, values []Value) (Value, Error) {

	arity := len(values)
	if arity != 1 {
		return nil, ArityMismatchError("1", arity)
	}

	return f.invoke(cx, values[0])
}

//--------------------------------------------------------------
// nativeFuncStr

type nativeFuncStr struct {
	*baseFunc
	invoke func(Context, Str) (Value, Error)
}

// NewNativeFuncStr creates a new NativeFunc that takes 1 Str parameter.
func NewNativeFuncStr(f func(Context, Str) (Value, Error)) NativeFunc {
	return &nativeFuncStr{&baseFunc{}, f}
}

func (f *nativeFuncStr) Freeze() (Value, Error) {
	return f, nil
}

func (f *nativeFuncStr) Eq(cx Context, v Value) (Bool, Error) {
	switch t := v.(type) {
	case *nativeFuncStr:
		// equality is based on identity
		return NewBool(f == t), nil
	default:
		return False, nil
	}
}

func (f *nativeFuncStr) Arity() *Arity { panic("Arity") }
func (f *nativeFuncStr) MinArity() int { return 1 }
func (f *nativeFuncStr) MaxArity() int { return 1 }

func (f *nativeFuncStr) Invoke(cx Context, values []Value) (Value, Error) {

	arity := len(values)
	if arity != 1 {
		return nil, ArityMismatchError("1", arity)
	}

	s, ok := values[0].(Str)
	if !ok {
		return nil, TypeMismatchError("Expected Str")
	}

	return f.invoke(cx, s)
}

//--------------------------------------------------------------
// nativeFuncInt

type nativeFuncInt struct {
	*baseFunc
	invoke func(Context, Int) (Value, Error)
}

// NewNativeFuncInt creates a new NativeFunc that takes 1 Int parameter.
func NewNativeFuncInt(f func(Context, Int) (Value, Error)) NativeFunc {
	return &nativeFuncInt{&baseFunc{}, f}
}

func (f *nativeFuncInt) Freeze() (Value, Error) {
	return f, nil
}

func (f *nativeFuncInt) Eq(cx Context, v Value) (Bool, Error) {
	switch t := v.(type) {
	case *nativeFuncInt:
		// equality is based on identity
		return NewBool(f == t), nil
	default:
		return False, nil
	}
}

func (f *nativeFuncInt) Arity() *Arity { panic("Arity") }
func (f *nativeFuncInt) MinArity() int { return 1 }
func (f *nativeFuncInt) MaxArity() int { return 1 }

func (f *nativeFuncInt) Invoke(cx Context, values []Value) (Value, Error) {

	arity := len(values)
	if arity != 1 {
		return nil, ArityMismatchError("1", arity)
	}

	i, ok := values[0].(Int)
	if !ok {
		return nil, TypeMismatchError("Expected Int")
	}

	return f.invoke(cx, i)
}

//--------------------------------------------------------------
// nativeFuncFunc

type nativeFuncFunc struct {
	*baseFunc
	invoke func(Context, Func) (Value, Error)
}

// NewNativeFuncFunc creates a new NativeFunc that takes 1 Func parameter.
func NewNativeFuncFunc(f func(Context, Func) (Value, Error)) NativeFunc {
	return &nativeFuncFunc{&baseFunc{}, f}
}

func (f *nativeFuncFunc) Freeze() (Value, Error) {
	return f, nil
}

func (f *nativeFuncFunc) Eq(cx Context, v Value) (Bool, Error) {
	switch t := v.(type) {
	case *nativeFuncFunc:
		// equality is based on identity
		return NewBool(f == t), nil
	default:
		return False, nil
	}
}

func (f *nativeFuncFunc) Arity() *Arity { panic("Arity") }
func (f *nativeFuncFunc) MinArity() int { return 1 }
func (f *nativeFuncFunc) MaxArity() int { return 1 }

func (f *nativeFuncFunc) Invoke(cx Context, values []Value) (Value, Error) {

	arity := len(values)
	if arity != 1 {
		return nil, ArityMismatchError("1", arity)
	}

	fn, ok := values[0].(Func)
	if !ok {
		return nil, TypeMismatchError("Expected Func")
	}

	return f.invoke(cx, fn)
}

//--------------------------------------------------------------
// nativeFuncBool

type nativeFuncBool struct {
	*baseFunc
	invoke func(Context, Bool) (Value, Error)
}

// NewNativeFuncBool creates a new NativeFunc that takes 1 Bool parameter.
func NewNativeFuncBool(f func(Context, Bool) (Value, Error)) NativeFunc {
	return &nativeFuncBool{&baseFunc{}, f}
}

func (f *nativeFuncBool) Freeze() (Value, Error) {
	return f, nil
}

func (f *nativeFuncBool) Eq(cx Context, v Value) (Bool, Error) {
	switch t := v.(type) {
	case *nativeFuncBool:
		// equality is based on identity
		return NewBool(f == t), nil
	default:
		return False, nil
	}
}

func (f *nativeFuncBool) Arity() *Arity { panic("Arity") }
func (f *nativeFuncBool) MinArity() int { return 1 }
func (f *nativeFuncBool) MaxArity() int { return 1 }

func (f *nativeFuncBool) Invoke(cx Context, values []Value) (Value, Error) {

	arity := len(values)
	if arity != 1 {
		return nil, ArityMismatchError("1", arity)
	}

	i, ok := values[0].(Bool)
	if !ok {
		return nil, TypeMismatchError("Expected Bool")
	}

	return f.invoke(cx, i)
}

//--------------------------------------------------------------
// nativeFuncStruct

type nativeFuncStruct struct {
	*baseFunc
	invoke func(Context, Struct) (Value, Error)
}

// NewNativeFuncStruct creates a new NativeFunc that takes 1 Struct parameter.
func NewNativeFuncStruct(f func(Context, Struct) (Value, Error)) NativeFunc {
	return &nativeFuncStruct{&baseFunc{}, f}
}

func (f *nativeFuncStruct) Freeze() (Value, Error) {
	return f, nil
}

func (f *nativeFuncStruct) Eq(cx Context, v Value) (Bool, Error) {
	switch t := v.(type) {
	case *nativeFuncStruct:
		// equality is based on identity
		return NewBool(f == t), nil
	default:
		return False, nil
	}
}

func (f *nativeFuncStruct) Arity() *Arity { panic("Arity") }
func (f *nativeFuncStruct) MinArity() int { return 1 }
func (f *nativeFuncStruct) MaxArity() int { return 1 }

func (f *nativeFuncStruct) Invoke(cx Context, values []Value) (Value, Error) {

	arity := len(values)
	if arity != 1 {
		return nil, ArityMismatchError("1", arity)
	}

	i, ok := values[0].(Struct)
	if !ok {
		return nil, TypeMismatchError("Expected Struct")
	}

	return f.invoke(cx, i)
}
