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
// nativeSimpleFunc
//--------------------------------------------------------------

type nativeSimpleFunc struct {
	arity  Arity
	invoke Invoke
}

// NewNativeFunc creates a NativeFunc
func NewNativeFunc(arity Arity, invoke Invoke) NativeFunc {
	return &nativeSimpleFunc{arity, invoke}
}

func (f *nativeSimpleFunc) nativeFuncMarker() {}

func (f *nativeSimpleFunc) Type() Type { return FuncType }

func (f *nativeSimpleFunc) Frozen(cx Context) (Bool, Error) {
	return True, nil
}
func (f *nativeSimpleFunc) Freeze(cx Context) (Value, Error) {
	return f, nil
}

func (f *nativeSimpleFunc) Eq(cx Context, v Value) (Bool, Error) {
	switch t := v.(type) {
	case *nativeSimpleFunc:
		// equality is based on identity
		return NewBool(f == t), nil
	default:
		return False, nil
	}
}

func (f *nativeSimpleFunc) HashCode(cx Context) (Int, Error) {
	return nil, TypeMismatchError("Expected Hashable Type")
}

//func (f *nativeSimpleFunc) GetField(cx Context, key Str) (Value, Error) {
//	return nil, NoSuchFieldError(key.String())
//}

func (f *nativeSimpleFunc) Cmp(cx Context, v Value) (Int, Error) {
	return nil, TypeMismatchError("Expected Comparable Type")
}

func (f *nativeSimpleFunc) ToStr(cx Context) Str {
	return NewStr(fmt.Sprintf("nativeFunc<%p>", f))
}

//--------------------------------
// fields

func (f *nativeSimpleFunc) FieldNames() ([]string, Error) {
	return []string{}, nil
}

func (f *nativeSimpleFunc) HasField(name string) (bool, Error) {
	return false, nil
}

func (f nativeSimpleFunc) GetField(name string, cx Context) (Value, Error) {
	return nil, NoSuchFieldError(name)
}

func (f nativeSimpleFunc) InvokeField(name string, cx Context, params []Value) (Value, Error) {
	return nil, NoSuchFieldError(name)
}

//--------------------------------
// func

func (f *nativeSimpleFunc) Arity() Arity {
	return f.arity
}

func (f *nativeSimpleFunc) Invoke(cx Context, params []Value) (Value, Error) {

	return f.invoke(cx, params)
}

//--------------------------------------------------------------
// nativeBaseFunc
//--------------------------------------------------------------

type nativeBaseFunc struct {
	arity         Arity
	requiredTypes []Type
	allowNull     bool
	invoke        Invoke
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

//--------------------------------
// fields

func (f *nativeBaseFunc) FieldNames() ([]string, Error) {
	return []string{}, nil
}

func (f *nativeBaseFunc) HasField(name string) (bool, Error) {
	return false, nil
}

func (f *nativeBaseFunc) GetField(name string, cx Context) (Value, Error) {
	return nil, NoSuchFieldError(name)
}

func (f *nativeBaseFunc) InvokeField(name string, cx Context, params []Value) (Value, Error) {
	return nil, NoSuchFieldError(name)
}

//--------------------------------
// func

func (f *nativeBaseFunc) Arity() Arity { return f.arity }

//--------------------------------------------------------------
// nativeFixedFunc
//--------------------------------------------------------------

type nativeFixedFunc struct {
	*nativeBaseFunc
}

// NewFixedNativeFunc is a convenience function for creating
// a NativeFunc with fixed arity
func NewFixedNativeFunc(
	requiredTypes []Type,
	allowNull bool,
	invoke Invoke) NativeFunc {

	arity := Arity{
		Kind:           FixedArity,
		RequiredParams: uint16(len(requiredTypes)),
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

func (f *nativeFixedFunc) Invoke(cx Context, params []Value) (Value, Error) {

	err := VetFixedFuncParams(f.requiredTypes, f.allowNull, params)
	if err != nil {
		return nil, err
	}

	return f.invoke(cx, params)
}

//--------------------------------------------------------------
// nativeVariadicFunc
//--------------------------------------------------------------

type nativeVariadicFunc struct {
	*nativeBaseFunc
	variadicType Type
}

// NewVariadicNativeFunc is a convenience function for creating
// a NativeFunc with variadic arity
func NewVariadicNativeFunc(
	requiredTypes []Type,
	variadicType Type,
	allowNull bool,
	invoke Invoke) NativeFunc {

	arity := Arity{
		Kind:           VariadicArity,
		RequiredParams: uint16(len(requiredTypes)),
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

func (f *nativeVariadicFunc) Invoke(cx Context, params []Value) (Value, Error) {

	numValues := len(params)
	numReqs := len(f.requiredTypes)

	// arity mismatch
	if numValues < numReqs {
		return nil, ArityMismatchError(
			fmt.Sprintf("at least %d", numReqs),
			numValues)
	}

	// check types on required params
	for i := 0; i < numReqs; i++ {
		err := VetFuncParam(f.requiredTypes[i], f.allowNull, params[i])
		if err != nil {
			return nil, err
		}
	}

	// check types on variadic params
	for i := numReqs; i < numValues; i++ {
		err := VetFuncParam(f.variadicType, f.allowNull, params[i])
		if err != nil {
			return nil, err
		}
	}

	// invoke
	return f.invoke(cx, params)
}

//--------------------------------------------------------------
// nativeMultipleFunc
//--------------------------------------------------------------

type nativeMultipleFunc struct {
	*nativeBaseFunc
	optionalTypes []Type
}

// NewMultipleNativeFunc is a convenience function for creating
// a NativeFunc with multiple arity
func NewMultipleNativeFunc(
	requiredTypes []Type,
	optionalTypes []Type,
	allowNull bool,
	invoke Invoke) NativeFunc {

	arity := Arity{
		Kind:           MultipleArity,
		RequiredParams: uint16(len(requiredTypes)),
		OptionalParams: uint16(len(optionalTypes)),
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

func (f *nativeMultipleFunc) Invoke(cx Context, params []Value) (Value, Error) {

	numValues := len(params)
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

	// check types on required params
	for i := 0; i < numReqs; i++ {
		err := VetFuncParam(f.requiredTypes[i], f.allowNull, params[i])
		if err != nil {
			return nil, err
		}
	}

	// check types on optional params
	for i := numReqs; i < numValues; i++ {
		err := VetFuncParam(f.optionalTypes[i-numReqs], f.allowNull, params[i])
		if err != nil {
			return nil, err
		}
	}

	// invoke
	return f.invoke(cx, params)
}

//--------------------------------------------------------------

func VetFixedFuncParams(requiredTypes []Type, allowNull bool, params []Value) Error {

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
		err := VetFuncParam(requiredTypes[i], allowNull, params[i])
		if err != nil {
			return err
		}
	}

	return nil
}

func VetFuncParam(typ Type, allowNull bool, param Value) Error {

	// accept 'any' type
	if typ == AnyType {
		return nil
	}

	// skip over null params
	if param.Type() == NullType {
		if allowNull {
			return nil
		}
		return NullValueError()
	}

	// check type
	if param.Type() != typ {
		return TypeMismatchError(fmt.Sprintf("Expected %s", typ.String()))
	}

	// invoke
	return nil
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
