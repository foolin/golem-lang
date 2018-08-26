// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package core

import (
	"fmt"
)

// ObsoleteFunc is a Func that is implemented in Go rather than Golem
type ObsoleteFunc interface {
	Func
}

//--------------------------------------------------------------

type obsoleteBaseFunc struct {
}

func (f *obsoleteBaseFunc) funcMarker() {}

func (f *obsoleteBaseFunc) Type() Type { return FuncType }

func (f *obsoleteBaseFunc) Frozen() (Bool, Error) {
	return True, nil
}

func (f *obsoleteBaseFunc) HashCode(cx Context) (Int, Error) {
	return nil, TypeMismatchError("Expected Hashable Type")
}

func (f *obsoleteBaseFunc) GetField(cx Context, key Str) (Value, Error) {
	return nil, NoSuchFieldError(key.String())
}

func (f *obsoleteBaseFunc) Cmp(cx Context, v Value) (Int, Error) {
	return nil, TypeMismatchError("Expected Comparable Type")
}

func (f *obsoleteBaseFunc) ToStr(cx Context) Str {
	return NewStr(fmt.Sprintf("obsoleteFunc<%p>", f))
}

//--------------------------------------------------------------

type obsoleteFunc struct {
	*obsoleteBaseFunc
	minArity int
	maxArity int
	invoke   func(Context, []Value) (Value, Error)
}

// NewObsoleteFunc creates a new ObsoleteFunc
func NewObsoleteFunc(
	minArity int,
	maxArity int,
	f func(Context, []Value) (Value, Error)) ObsoleteFunc {

	return &obsoleteFunc{&obsoleteBaseFunc{}, minArity, maxArity, f}
}

func (f *obsoleteFunc) Freeze() (Value, Error) {
	return f, nil
}

func (f *obsoleteFunc) Eq(cx Context, v Value) (Bool, Error) {
	switch t := v.(type) {
	case *obsoleteFunc:
		// equality is based on identity
		return NewBool(f == t), nil
	default:
		return False, nil
	}
}

func (f *obsoleteFunc) Arity() *Arity { panic("Arity") }
func (f *obsoleteFunc) MinArity() int { return f.minArity }
func (f *obsoleteFunc) MaxArity() int { return f.maxArity }

func (f *obsoleteFunc) Invoke(cx Context, values []Value) (Value, Error) {

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
// obsoleteFunc0

type obsoleteFunc0 struct {
	*obsoleteBaseFunc
	invoke func(Context) (Value, Error)
}

// NewObsoleteFunc0 creates a new ObsoleteFunc that takes 0 parameters.
func NewObsoleteFunc0(f func(Context) (Value, Error)) ObsoleteFunc {
	return &obsoleteFunc0{&obsoleteBaseFunc{}, f}
}

func (f *obsoleteFunc0) Freeze() (Value, Error) {
	return f, nil
}

func (f *obsoleteFunc0) Eq(cx Context, v Value) (Bool, Error) {
	switch t := v.(type) {
	case *obsoleteFunc0:
		// equality is based on identity
		return NewBool(f == t), nil
	default:
		return False, nil
	}
}

func (f *obsoleteFunc0) Arity() *Arity { panic("Arity") }
func (f *obsoleteFunc0) MinArity() int { return 0 }
func (f *obsoleteFunc0) MaxArity() int { return 0 }

func (f *obsoleteFunc0) Invoke(cx Context, values []Value) (Value, Error) {

	arity := len(values)
	if arity != 0 {
		return nil, ArityMismatchError("0", arity)
	}

	return f.invoke(cx)
}

//--------------------------------------------------------------
// obsoleteFuncValue

type obsoleteFuncValue struct {
	*obsoleteBaseFunc
	invoke func(Context, Value) (Value, Error)
}

// NewObsoleteFuncValue creates a new ObsoleteFunc that takes 1 Value parameter.
func NewObsoleteFuncValue(f func(Context, Value) (Value, Error)) ObsoleteFunc {
	return &obsoleteFuncValue{&obsoleteBaseFunc{}, f}
}

func (f *obsoleteFuncValue) Freeze() (Value, Error) {
	return f, nil
}

func (f *obsoleteFuncValue) Eq(cx Context, v Value) (Bool, Error) {
	switch t := v.(type) {
	case *obsoleteFuncValue:
		// equality is based on identity
		return NewBool(f == t), nil
	default:
		return False, nil
	}
}

func (f *obsoleteFuncValue) Arity() *Arity { panic("Arity") }
func (f *obsoleteFuncValue) MinArity() int { return 1 }
func (f *obsoleteFuncValue) MaxArity() int { return 1 }

func (f *obsoleteFuncValue) Invoke(cx Context, values []Value) (Value, Error) {

	arity := len(values)
	if arity != 1 {
		return nil, ArityMismatchError("1", arity)
	}

	return f.invoke(cx, values[0])
}

//--------------------------------------------------------------
// obsoleteFuncStr

type obsoleteFuncStr struct {
	*obsoleteBaseFunc
	invoke func(Context, Str) (Value, Error)
}

// NewObsoleteFuncStr creates a new ObsoleteFunc that takes 1 Str parameter.
func NewObsoleteFuncStr(f func(Context, Str) (Value, Error)) ObsoleteFunc {
	return &obsoleteFuncStr{&obsoleteBaseFunc{}, f}
}

func (f *obsoleteFuncStr) Freeze() (Value, Error) {
	return f, nil
}

func (f *obsoleteFuncStr) Eq(cx Context, v Value) (Bool, Error) {
	switch t := v.(type) {
	case *obsoleteFuncStr:
		// equality is based on identity
		return NewBool(f == t), nil
	default:
		return False, nil
	}
}

func (f *obsoleteFuncStr) Arity() *Arity { panic("Arity") }
func (f *obsoleteFuncStr) MinArity() int { return 1 }
func (f *obsoleteFuncStr) MaxArity() int { return 1 }

func (f *obsoleteFuncStr) Invoke(cx Context, values []Value) (Value, Error) {

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
// obsoleteFuncInt

type obsoleteFuncInt struct {
	*obsoleteBaseFunc
	invoke func(Context, Int) (Value, Error)
}

// NewObsoleteFuncInt creates a new ObsoleteFunc that takes 1 Int parameter.
func NewObsoleteFuncInt(f func(Context, Int) (Value, Error)) ObsoleteFunc {
	return &obsoleteFuncInt{&obsoleteBaseFunc{}, f}
}

func (f *obsoleteFuncInt) Freeze() (Value, Error) {
	return f, nil
}

func (f *obsoleteFuncInt) Eq(cx Context, v Value) (Bool, Error) {
	switch t := v.(type) {
	case *obsoleteFuncInt:
		// equality is based on identity
		return NewBool(f == t), nil
	default:
		return False, nil
	}
}

func (f *obsoleteFuncInt) Arity() *Arity { panic("Arity") }
func (f *obsoleteFuncInt) MinArity() int { return 1 }
func (f *obsoleteFuncInt) MaxArity() int { return 1 }

func (f *obsoleteFuncInt) Invoke(cx Context, values []Value) (Value, Error) {

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
// obsoleteFuncFunc

type obsoleteFuncFunc struct {
	*obsoleteBaseFunc
	invoke func(Context, Func) (Value, Error)
}

// NewObsoleteFuncFunc creates a new ObsoleteFunc that takes 1 Func parameter.
func NewObsoleteFuncFunc(f func(Context, Func) (Value, Error)) ObsoleteFunc {
	return &obsoleteFuncFunc{&obsoleteBaseFunc{}, f}
}

func (f *obsoleteFuncFunc) Freeze() (Value, Error) {
	return f, nil
}

func (f *obsoleteFuncFunc) Eq(cx Context, v Value) (Bool, Error) {
	switch t := v.(type) {
	case *obsoleteFuncFunc:
		// equality is based on identity
		return NewBool(f == t), nil
	default:
		return False, nil
	}
}

func (f *obsoleteFuncFunc) Arity() *Arity { panic("Arity") }
func (f *obsoleteFuncFunc) MinArity() int { return 1 }
func (f *obsoleteFuncFunc) MaxArity() int { return 1 }

func (f *obsoleteFuncFunc) Invoke(cx Context, values []Value) (Value, Error) {

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
// obsoleteFuncBool

type obsoleteFuncBool struct {
	*obsoleteBaseFunc
	invoke func(Context, Bool) (Value, Error)
}

// NewObsoleteFuncBool creates a new ObsoleteFunc that takes 1 Bool parameter.
func NewObsoleteFuncBool(f func(Context, Bool) (Value, Error)) ObsoleteFunc {
	return &obsoleteFuncBool{&obsoleteBaseFunc{}, f}
}

func (f *obsoleteFuncBool) Freeze() (Value, Error) {
	return f, nil
}

func (f *obsoleteFuncBool) Eq(cx Context, v Value) (Bool, Error) {
	switch t := v.(type) {
	case *obsoleteFuncBool:
		// equality is based on identity
		return NewBool(f == t), nil
	default:
		return False, nil
	}
}

func (f *obsoleteFuncBool) Arity() *Arity { panic("Arity") }
func (f *obsoleteFuncBool) MinArity() int { return 1 }
func (f *obsoleteFuncBool) MaxArity() int { return 1 }

func (f *obsoleteFuncBool) Invoke(cx Context, values []Value) (Value, Error) {

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
// obsoleteFuncStruct

type obsoleteFuncStruct struct {
	*obsoleteBaseFunc
	invoke func(Context, Struct) (Value, Error)
}

// NewObsoleteFuncStruct creates a new ObsoleteFunc that takes 1 Struct parameter.
func NewObsoleteFuncStruct(f func(Context, Struct) (Value, Error)) ObsoleteFunc {
	return &obsoleteFuncStruct{&obsoleteBaseFunc{}, f}
}

func (f *obsoleteFuncStruct) Freeze() (Value, Error) {
	return f, nil
}

func (f *obsoleteFuncStruct) Eq(cx Context, v Value) (Bool, Error) {
	switch t := v.(type) {
	case *obsoleteFuncStruct:
		// equality is based on identity
		return NewBool(f == t), nil
	default:
		return False, nil
	}
}

func (f *obsoleteFuncStruct) Arity() *Arity { panic("Arity") }
func (f *obsoleteFuncStruct) MinArity() int { return 1 }
func (f *obsoleteFuncStruct) MaxArity() int { return 1 }

func (f *obsoleteFuncStruct) Invoke(cx Context, values []Value) (Value, Error) {

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
