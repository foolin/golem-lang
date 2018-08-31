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
}

//--------------------------------------------------------------
// nativeFunc
//--------------------------------------------------------------

type nativeFunc struct {
	arity  Arity
	invoke func(Evaluator, []Value) (Value, Error)
}

func (f *nativeFunc) funcMarker() {}

func (f *nativeFunc) Type() Type { return FuncType }

func (f *nativeFunc) Eq(ev Evaluator, v Value) (Bool, Error) {
	switch t := v.(type) {
	case *nativeFunc:
		// equality is based on identity
		return NewBool(f == t), nil
	default:
		return False, nil
	}
}

func (f *nativeFunc) Freeze(ev Evaluator) (Value, Error) {
	return f, nil
}

func (f *nativeFunc) Frozen(ev Evaluator) (Bool, Error) {
	return True, nil
}

func (f *nativeFunc) HashCode(ev Evaluator) (Int, Error) {
	return nil, TypeMismatchError("Expected Hashable Type")
}

func (f *nativeFunc) Cmp(ev Evaluator, v Value) (Int, Error) {
	return nil, TypeMismatchError("Expected Comparable Type")
}

func (f *nativeFunc) ToStr(ev Evaluator) Str {
	return NewStr(fmt.Sprintf("nativeFunc<%p>", f))
}

func (f *nativeFunc) Arity() Arity { return f.arity }

//--------------------------------
// fields

func (f *nativeFunc) FieldNames() ([]string, Error) {
	return []string{}, nil
}

func (f *nativeFunc) HasField(name string) (bool, Error) {
	return false, nil
}

func (f *nativeFunc) GetField(name string, ev Evaluator) (Value, Error) {
	return nil, NoSuchFieldError(name)
}

func (f *nativeFunc) InvokeField(name string, ev Evaluator, params []Value) (Value, Error) {
	return nil, NoSuchFieldError(name)
}

//--------------------------------------------------------------
// nativeFixedFunc
//--------------------------------------------------------------

type nativeFixedFunc struct {
	*nativeFunc
	requiredTypes []Type
	allowNull     bool
}

// NewFixedNativeFunc creates a new NativeFunc with fixed arity
func NewFixedNativeFunc(
	requiredTypes []Type,
	allowNull bool,
	invoke func(Evaluator, []Value) (Value, Error)) NativeFunc {

	arity := Arity{
		Kind:           FixedArity,
		RequiredParams: uint16(len(requiredTypes)),
		OptionalParams: 0,
	}

	return &nativeFixedFunc{
		&nativeFunc{arity, invoke},
		requiredTypes, allowNull,
	}
}

func (f *nativeFixedFunc) Freeze(ev Evaluator) (Value, Error) {
	return f, nil
}

func (f *nativeFixedFunc) Eq(ev Evaluator, v Value) (Bool, Error) {
	switch t := v.(type) {
	case *nativeFixedFunc:
		// equality is based on identity
		return NewBool(f == t), nil
	default:
		return False, nil
	}
}

func (f *nativeFixedFunc) Invoke(ev Evaluator, params []Value) (Value, Error) {

	err := vetFixedParams(params, f.requiredTypes, f.allowNull)
	if err != nil {
		return nil, err
	}

	return f.invoke(ev, params)
}

//--------------------------------------------------------------
// nativeVariadicFunc
//--------------------------------------------------------------

type nativeVariadicFunc struct {
	*nativeFunc
	requiredTypes []Type
	variadicType  Type
	allowNull     bool
}

// NewVariadicNativeFunc creates a new NativeFunc with variadic arity
func NewVariadicNativeFunc(
	requiredTypes []Type,
	variadicType Type,
	allowNull bool,
	invoke func(Evaluator, []Value) (Value, Error)) NativeFunc {

	arity := Arity{
		Kind:           VariadicArity,
		RequiredParams: uint16(len(requiredTypes)),
		OptionalParams: 0,
	}

	return &nativeVariadicFunc{
		&nativeFunc{arity, invoke},
		requiredTypes, variadicType, allowNull,
	}
}

func (f *nativeVariadicFunc) Freeze(ev Evaluator) (Value, Error) {
	return f, nil
}

func (f *nativeVariadicFunc) Eq(ev Evaluator, v Value) (Bool, Error) {
	switch t := v.(type) {
	case *nativeVariadicFunc:
		// equality is based on identity
		return NewBool(f == t), nil
	default:
		return False, nil
	}
}

func (f *nativeVariadicFunc) Invoke(ev Evaluator, params []Value) (Value, Error) {

	err := vetVariadicParams(params, f.requiredTypes, f.variadicType, f.allowNull)
	if err != nil {
		return nil, err
	}

	return f.invoke(ev, params)
}

//--------------------------------------------------------------
// nativeMultipleFunc
//--------------------------------------------------------------

type nativeMultipleFunc struct {
	*nativeFunc
	requiredTypes []Type
	optionalTypes []Type
	allowNull     bool
}

// NewMultipleNativeFunc creates a new NativeFunc with multiple arity
func NewMultipleNativeFunc(
	requiredTypes []Type,
	optionalTypes []Type,
	allowNull bool,
	invoke func(Evaluator, []Value) (Value, Error)) NativeFunc {

	arity := Arity{
		Kind:           MultipleArity,
		RequiredParams: uint16(len(requiredTypes)),
		OptionalParams: uint16(len(optionalTypes)),
	}

	return &nativeMultipleFunc{
		&nativeFunc{arity, invoke},
		requiredTypes, optionalTypes, allowNull,
	}
}

func (f *nativeMultipleFunc) Freeze(ev Evaluator) (Value, Error) {
	return f, nil
}

func (f *nativeMultipleFunc) Eq(ev Evaluator, v Value) (Bool, Error) {
	switch t := v.(type) {
	case *nativeMultipleFunc:
		// equality is based on identity
		return NewBool(f == t), nil
	default:
		return False, nil
	}
}

func (f *nativeMultipleFunc) Invoke(ev Evaluator, params []Value) (Value, Error) {

	err := vetMultipleParams(params, f.requiredTypes, f.optionalTypes, f.allowNull)
	if err != nil {
		return nil, err
	}

	return f.invoke(ev, params)
}

//--------------------------------------------------------------
// vet params
//--------------------------------------------------------------

func vetFixedParams(
	params []Value,
	requiredTypes []Type,
	allowNull bool) Error {

	numValues := len(params)
	numReqs := len(requiredTypes)

	// arity mismatch
	if numValues != numReqs {
		return ArityMismatchError(
			fmt.Sprintf("%d", numReqs),
			numValues)
	}

	// check types on required params
	for i := 0; i < numReqs; i++ {
		err := vetParam(params[i], requiredTypes[i], allowNull)
		if err != nil {
			return err
		}
	}

	return nil
}

func vetVariadicParams(
	params []Value,
	requiredTypes []Type,
	variadicType Type,
	allowNull bool) Error {

	numValues := len(params)
	numReqs := len(requiredTypes)

	// arity mismatch
	if numValues < numReqs {
		return ArityMismatchError(
			fmt.Sprintf("at least %d", numReqs),
			numValues)
	}

	// check types on required params
	for i := 0; i < numReqs; i++ {
		err := vetParam(params[i], requiredTypes[i], allowNull)
		if err != nil {
			return err
		}
	}

	// check types on variadic params
	for i := numReqs; i < numValues; i++ {
		err := vetParam(params[i], variadicType, allowNull)
		if err != nil {
			return err
		}
	}

	return nil
}

func vetMultipleParams(
	params []Value,
	requiredTypes []Type,
	optionalTypes []Type,
	allowNull bool) Error {

	numValues := len(params)
	numReqs := len(requiredTypes)
	numOpts := len(optionalTypes)

	// arity mismatch
	if numValues < numReqs {
		return ArityMismatchError(
			fmt.Sprintf("at least %d", numReqs),
			numValues)
	}
	if numValues > (numReqs + numOpts) {
		return ArityMismatchError(
			fmt.Sprintf("at most %d", numReqs+numOpts),
			numValues)
	}

	// check types on required params
	for i := 0; i < numReqs; i++ {
		err := vetParam(params[i], requiredTypes[i], allowNull)
		if err != nil {
			return err
		}
	}

	// check types on optional params
	for i := numReqs; i < numValues; i++ {
		err := vetParam(params[i], optionalTypes[i-numReqs], allowNull)
		if err != nil {
			return err
		}
	}

	return nil
}

func vetParam(value Value, typ Type, allowNull bool) Error {

	// accept 'any' type
	if typ == AnyType {
		return nil
	}

	// skip over null params
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

	return nil
}
