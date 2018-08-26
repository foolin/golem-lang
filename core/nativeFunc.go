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

//--------------------------------------------------------------

type nativeBaseFunc struct {
	arity *Arity
}

func (f *nativeBaseFunc) funcMarker() {}

func (f *nativeBaseFunc) Type() Type { return FuncType }

func (f *nativeBaseFunc) Frozen() (Bool, Error) {
	return True, nil
}

func (f *nativeBaseFunc) HashCode(cx Context) (Int, Error) {
	return nil, TypeMismatchError("Expected Hashable Type")
}

func (f *nativeBaseFunc) GetField(cx Context, key Str) (Value, Error) {
	return nil, NoSuchFieldError(key.String())
}

func (f *nativeBaseFunc) Cmp(cx Context, v Value) (Int, Error) {
	return nil, TypeMismatchError("Expected Comparable Type")
}

func (f *nativeBaseFunc) ToStr(cx Context) Str {
	return NewStr(fmt.Sprintf("nativeFunc<%p>", f))
}

func (f *nativeBaseFunc) MinArity() int { panic("MinArity") }
func (f *nativeBaseFunc) MaxArity() int { panic("MaxArity") }
func (f *nativeBaseFunc) Arity() *Arity { return f.arity }

//--------------------------------------------------------------

type nativeFixedFunc struct {
	*nativeBaseFunc

	types  []Type
	invoke func(Context, []Value) (Value, Error)
}

// NewFixedNativeFunc creates a new NativeFunc with fixed arity
func NewFixedNativeFunc(
	types []Type,
	invoke func(Context, []Value) (Value, Error)) NativeFunc {

	arity := &Arity{
		Kind:           FixedArity,
		RequiredParams: len(types),
		OptionalParams: nil,
	}

	return &nativeFixedFunc{&nativeBaseFunc{arity}, types, invoke}
}

func (f *nativeFixedFunc) Freeze() (Value, Error) {
	return f, nil
}

func (f *nativeFixedFunc) Eq(cx Context, v Value) (Bool, Error) {
	switch t := v.(type) {
	case *nativeFixedFunc:
		// equality is based on identity
		return NewBool(f == t), nil
	default:
		return False, nil
	}
}

func (f *nativeFixedFunc) Invoke(cx Context, values []Value) (Value, Error) {

	// arity mismatch
	if len(values) != len(f.types) {
		return nil, ArityMismatchError(
			fmt.Sprintf("%d", len(f.types)),
			len(values))
	}

	// type mismatch
	for i, t := range f.types {
		// accept 'any' type
		if t == AnyType {
			continue
		}
		// skip over null values
		if values[i].Type() == NullType {
			continue
		}

		if values[i].Type() != t {
			return nil, TypeMismatchError(fmt.Sprintf("Expected %s", t.String()))
		}
	}

	return f.invoke(cx, values)
}
