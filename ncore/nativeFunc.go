// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package ncore

import (
	"fmt"
)

// NativeFunc is a Func that is implemented in Go rather than Golem
type NativeFunc interface {
	Func
	nativeFuncMarker()
}

//--------------------------------------------------------------
// nativeBaseFunc
//--------------------------------------------------------------

type nativeBaseFunc struct {
	arity         *Arity
	requiredTypes []Type
	allowNull     bool
	invoke        func(Context, []Value) (Value, Error)
}

func (f *nativeBaseFunc) nativeFuncMarker() {}

func (f *nativeBaseFunc) Type() Type { return FuncType }

func (f *nativeBaseFunc) Frozen(cx Context) (Bool, Error) {
	return True, nil
}

func (f *nativeBaseFunc) HashCode(cx Context) (Int, Error) {
	return nil, TypeMismatchError("Expected Hashable Type")
}

//func (f *nativeBaseFunc) GetField(cx Context, key Str) (Value, Error) {
//	return nil, NoSuchFieldError(key.String())
//}

func (f *nativeBaseFunc) Cmp(cx Context, v Value) (Int, Error) {
	return nil, TypeMismatchError("Expected Comparable Type")
}

func (f *nativeBaseFunc) ToStr(cx Context) Str {
	return NewStr(fmt.Sprintf("nativeFunc<%p>", f))
}

func (f *nativeBaseFunc) FieldNames() ([]string, Error) {
	return []string{}, nil
}

func (f *nativeBaseFunc) HasField(Context, Value) (Bool, Error) {
	return False, nil
}

func (f *nativeBaseFunc) Arity() *Arity { return f.arity }

func checkType(value Value, typ Type, allowNull bool) Error {

	// accept 'any' type
	if typ == AnyType {
		return nil
	}

	// skip over null values
	if value.Type() == NullType {
		if allowNull {
			return nil
		}
		return NullValueError()
	}

	// check type
	if value.Type() != typ {
		return TypeMismatchError(fmt.Sprintf("Expected %s", typ.String()))
	}

	// invoke
	return nil
}

//--------------------------------------------------------------
// nativeFixedFunc
//--------------------------------------------------------------

type nativeFixedFunc struct {
	*nativeBaseFunc
}

// NewFixedNativeFunc creates a new NativeFunc with fixed arity
func NewFixedNativeFunc(
	requiredTypes []Type,
	allowNull bool,
	invoke func(Context, []Value) (Value, Error)) NativeFunc {

	arity := &Arity{
		Kind:           FixedArity,
		RequiredParams: len(requiredTypes),
		OptionalParams: nil,
	}

	return &nativeFixedFunc{
		&nativeBaseFunc{arity, requiredTypes, allowNull, invoke},
	}
}

func (f *nativeFixedFunc) Freeze(cx Context) (Value, Error) {
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

	numValues := len(values)
	numReqs := len(f.requiredTypes)

	// arity mismatch
	if numValues != numReqs {
		return nil, ArityMismatchError(
			fmt.Sprintf("%d", numReqs),
			numValues)
	}

	// check types on required values
	for i := 0; i < numReqs; i++ {
		err := checkType(values[i], f.requiredTypes[i], f.allowNull)
		if err != nil {
			return nil, err
		}
	}

	return f.invoke(cx, values)
}

//--------------------------------------------------------------
// nativeVariadicFunc
//--------------------------------------------------------------

type nativeVariadicFunc struct {
	*nativeBaseFunc
	variadicType Type
}

// NewVariadicNativeFunc creates a new NativeFunc with variadic arity
func NewVariadicNativeFunc(
	requiredTypes []Type,
	variadicType Type,
	allowNull bool,
	invoke func(Context, []Value) (Value, Error)) NativeFunc {

	arity := &Arity{
		Kind:           VariadicArity,
		RequiredParams: len(requiredTypes),
		OptionalParams: nil,
	}

	return &nativeVariadicFunc{
		&nativeBaseFunc{arity, requiredTypes, allowNull, invoke},
		variadicType,
	}
}

func (f *nativeVariadicFunc) Freeze(cx Context) (Value, Error) {
	return f, nil
}

func (f *nativeVariadicFunc) Eq(cx Context, v Value) (Bool, Error) {
	switch t := v.(type) {
	case *nativeVariadicFunc:
		// equality is based on identity
		return NewBool(f == t), nil
	default:
		return False, nil
	}
}

func (f *nativeVariadicFunc) Invoke(cx Context, values []Value) (Value, Error) {

	numValues := len(values)
	numReqs := len(f.requiredTypes)

	// arity mismatch
	if numValues < numReqs {
		return nil, ArityMismatchError(
			fmt.Sprintf("at least %d", numReqs),
			numValues)
	}

	// check types on required values
	for i := 0; i < numReqs; i++ {
		err := checkType(values[i], f.requiredTypes[i], f.allowNull)
		if err != nil {
			return nil, err
		}
	}

	// check types on variadic values
	for i := numReqs; i < numValues; i++ {
		err := checkType(values[i], f.variadicType, f.allowNull)
		if err != nil {
			return nil, err
		}
	}

	// invoke
	return f.invoke(cx, values)
}

//--------------------------------------------------------------
// nativeMultipleFunc
//--------------------------------------------------------------

type nativeMultipleFunc struct {
	*nativeBaseFunc
	optionalTypes []Type
}

// NewMultipleNativeFunc creates a new NativeFunc with multiple arity
func NewMultipleNativeFunc(
	requiredTypes []Type,
	optionalTypes []Type,
	allowNull bool,
	invoke func(Context, []Value) (Value, Error)) NativeFunc {

	arity := &Arity{
		Kind:           MultipleArity,
		RequiredParams: len(requiredTypes),
		OptionalParams: nil,
	}

	return &nativeMultipleFunc{
		&nativeBaseFunc{arity, requiredTypes, allowNull, invoke},
		optionalTypes,
	}
}

func (f *nativeMultipleFunc) Freeze(cx Context) (Value, Error) {
	return f, nil
}

func (f *nativeMultipleFunc) Eq(cx Context, v Value) (Bool, Error) {
	switch t := v.(type) {
	case *nativeMultipleFunc:
		// equality is based on identity
		return NewBool(f == t), nil
	default:
		return False, nil
	}
}

func (f *nativeMultipleFunc) Invoke(cx Context, values []Value) (Value, Error) {

	numValues := len(values)
	numReqs := len(f.requiredTypes)
	numOpts := len(f.optionalTypes)

	// arity mismatch
	if numValues < numReqs {
		return nil, ArityMismatchError(
			fmt.Sprintf("at least %d", numReqs),
			numValues)
	}
	if numValues > (numReqs + numOpts) {
		return nil, ArityMismatchError(
			fmt.Sprintf("at most %d", numReqs+numOpts),
			numValues)
	}

	// check types on required values
	for i := 0; i < numReqs; i++ {
		err := checkType(values[i], f.requiredTypes[i], f.allowNull)
		if err != nil {
			return nil, err
		}
	}

	// check types on optional values
	for i := numReqs; i < numValues; i++ {
		err := checkType(values[i], f.optionalTypes[i-numReqs], f.allowNull)
		if err != nil {
			return nil, err
		}
	}

	// invoke
	return f.invoke(cx, values)
}

////--------------------------------------------------------------
//// virtualFunc
////--------------------------------------------------------------
//
//// A virtual function is a function that is an intrinsic
//// part of a given Type. To reduce the overhead required to make
//// new values, these functions are created on the fly.
//type virtualFunc struct {
//	owner Value
//	name  string
//	NativeFunc
//}
//
//func (f *virtualFunc) Eq(cx Context, v Value) (Bool, Error) {
//	switch t := v.(type) {
//	case *virtualFunc:
//		// equality for intrinsic functions is based on whether
//		// they have the same owner, and the same name
//		ownerEq, err := f.owner.Eq(cx, t.owner)
//		if err != nil {
//			return nil, err
//		}
//		return NewBool(ownerEq.BoolVal() && (f.name == t.name)), nil
//	default:
//		return False, nil
//	}
//}
