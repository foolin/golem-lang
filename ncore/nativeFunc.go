// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package ncore

import (
	"fmt"
)

//--------------------------------------------------------------
// NativeFunc
//--------------------------------------------------------------

// NativeFunc is a Func that is implemented in Go
type (
	NativeFunc interface {
		Func
		nativeFuncMarker()
	}

	//Method struct {
	//	Arity  Arity
	//	Invoke func(Context, Value, []Value) (Value, Error)
	//}

	Invoke func(Context, []Value) (Value, Error)
)

func NewNativeFunc(arity Arity, invoke Invoke) NativeFunc {
	return &nativeFunc{arity, invoke}
}

// NewFixedNativeFunc is a convenience function for creating
// a NativeFunc with fixed arity
func NewFixedNativeFunc(
	requiredTypes []Type,
	allowNull bool,
	invoke Invoke) NativeFunc {

	arity, wrapper := NewFixedFuncDef(requiredTypes, allowNull, invoke)
	return &nativeFunc{arity, wrapper}
}

// NewFixedFuncDef creates the definition for a fixed NativeFunc
func NewFixedFuncDef(
	requiredTypes []Type,
	allowNull bool,
	invoker Invoke) (arity Arity, wrapper Invoke) {

	arity = Arity{
		Kind:           FixedArity,
		RequiredParams: uint16(len(requiredTypes)),
	}

	wrapper = func(cx Context, params []Value) (Value, Error) {
		numValues := len(params)
		numReqs := len(requiredTypes)

		// arity mismatch
		if numValues != numReqs {
			return nil, ArityMismatchError(
				fmt.Sprintf("%d", numReqs),
				numValues)
		}

		// check types on required params
		for i := 0; i < numReqs; i++ {
			err := vetFuncParam(requiredTypes[i], allowNull, params[i])
			if err != nil {
				return nil, err
			}
		}

		return invoker(cx, params)
	}

	return
}

func vetFuncParam(typ Type, allowNull bool, param Value) Error {

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

//--------------------------------------------------------------
// nativeFunc
//--------------------------------------------------------------

type nativeFunc struct {
	arity  Arity
	invoke Invoke
}

func (f *nativeFunc) nativeFuncMarker() {}

func (f *nativeFunc) Type() Type { return FuncType }

func (f *nativeFunc) Frozen(cx Context) (Bool, Error) {
	return True, nil
}

func (f *nativeFunc) Freeze(cx Context) (Value, Error) {
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

func (f *nativeFunc) HashCode(cx Context) (Int, Error) {
	return nil, TypeMismatchError("Expected Hashable Type")
}

func (f *nativeFunc) Cmp(cx Context, v Value) (Int, Error) {
	return nil, TypeMismatchError("Expected Comparable Type")
}

func (f *nativeFunc) ToStr(cx Context) Str {
	return NewStr(fmt.Sprintf("nativeFunc<%p>", f))
}

func (f *nativeFunc) Arity() Arity {
	return f.arity
}

func (f *nativeFunc) Invoke(cx Context, params []Value) (Value, Error) {
	return f.invoke(cx, params)
}

//--------------------------------
// fields

func (f *nativeFunc) FieldNames() ([]string, Error) {
	return []string{}, nil
}

func (f *nativeFunc) HasField(name string) (bool, Error) {
	return false, nil
}

func (f *nativeFunc) GetField(name string, cx Context) (Value, Error) {
	return nil, NoSuchFieldError(name)
}

func (f *nativeFunc) InvokeField(name string, cx Context, params []Value) (Value, Error) {
	return nil, NoSuchFieldError(name)
}

//--------------------------------------------------------------
// virtualFunc
//--------------------------------------------------------------

// virtualFunc is a function that is created only if
// we need to pass it around as a first class value.
// The 'same' virtual func can be instantiated more than once,
// so it has a special Eq() implementation.
type virtualFunc struct {
	owner Value
	name  string
	*nativeFunc
}

func newVirtualFunc(owner Value, name string, arity Arity, invoke Invoke) NativeFunc {
	return &virtualFunc{
		owner, name,
		&nativeFunc{arity, invoke},
	}
}

func (f *virtualFunc) Freeze(cx Context) (Value, Error) {
	return f, nil
}

func (f *virtualFunc) Eq(cx Context, v Value) (Bool, Error) {
	switch t := v.(type) {
	case *virtualFunc:
		// equality for virtual functions is based on whether
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

////--------------------------------------------------------------
////--------------------------------------------------------------
////--------------------------------------------------------------

////--------------------------------------------------------------
//// nativeBaseFunc
////--------------------------------------------------------------
//
//type nativeBaseFunc struct {
//	arity  Arity
//	invoke Invoke
//}
//
//func (f *nativeBaseFunc) nativeFuncMarker() {}
//
//func (f *nativeBaseFunc) Type() Type { return FuncType }
//
//func (f *nativeBaseFunc) Frozen(cx Context) (Bool, Error) {
//	return True, nil
//}
//
//func (f *nativeBaseFunc) HashCode(cx Context) (Int, Error) {
//	return nil, TypeMismatchError("Expected Hashable Type")
//}
//
////func (f *nativeBaseFunc) GetField(cx Context, key Str) (Value, Error) {
////	return nil, NoSuchFieldError(key.String())
////}
//
//func (f *nativeBaseFunc) Cmp(cx Context, v Value) (Int, Error) {
//	return nil, TypeMismatchError("Expected Comparable Type")
//}
//
//func (f *nativeBaseFunc) ToStr(cx Context) Str {
//	return NewStr(fmt.Sprintf("nativeFunc<%p>", f))
//}
//
////--------------------------------
//// fields
//
//func (f *nativeBaseFunc) FieldNames() ([]string, Error) {
//	return []string{}, nil
//}
//
//func (f *nativeBaseFunc) HasField(name string) (bool, Error) {
//	return false, nil
//}
//
//func (f *nativeBaseFunc) GetField(name string, cx Context) (Value, Error) {
//	return nil, NoSuchFieldError(name)
//}
//
//func (f *nativeBaseFunc) InvokeField(name string, cx Context, params []Value) (Value, Error) {
//	return nil, NoSuchFieldError(name)
//}
//
////--------------------------------
//// func
//
//func (f *nativeBaseFunc) Arity() Arity { return f.arity }
//
////--------------------------------------------------------------
//// nativeFunc
////--------------------------------------------------------------
//
//type nativeFunc struct {
//	*nativeBaseFunc
//}
//
//// NewNativeFunc creates a NativeFunc
//func NewNativeFunc(arity Arity, invoke Invoke) NativeFunc {
//	return &nativeFunc{
//		&nativeBaseFunc{arity, invoke},
//	}
//}
//
//func (f *nativeFunc) Freeze(cx Context) (Value, Error) {
//	return f, nil
//}
//
//func (f *nativeFunc) Eq(cx Context, v Value) (Bool, Error) {
//	switch t := v.(type) {
//	case *nativeFunc:
//		// equality is based on identity
//		return NewBool(f == t), nil
//	default:
//		return False, nil
//	}
//}
//
//func (f *nativeFunc) Invoke(cx Context, params []Value) (Value, Error) {
//	return f.invoke(cx, params)
//}
//
////--------------------------------------------------------------
//// nativeFixedFunc
////--------------------------------------------------------------
//
//type nativeFixedFunc struct {
//	*nativeBaseFunc
//	requiredTypes []Type
//	allowNull     bool
//}
//
//// NewFixedNativeFunc is a convenience function for creating
//// a NativeFunc with fixed arity
//func NewFixedNativeFunc(
//	requiredTypes []Type,
//	allowNull bool,
//	invoke Invoke) NativeFunc {
//
//	arity := Arity{
//		Kind:           FixedArity,
//		RequiredParams: uint16(len(requiredTypes)),
//	}
//
//	return &nativeFixedFunc{
//		&nativeBaseFunc{arity, invoke},
//		requiredTypes, allowNull,
//	}
//}
//
//func (f *nativeFixedFunc) Freeze(cx Context) (Value, Error) {
//	return f, nil
//}
//
//func (f *nativeFixedFunc) Eq(cx Context, v Value) (Bool, Error) {
//	switch t := v.(type) {
//	case *nativeFixedFunc:
//		// equality is based on identity
//		return NewBool(f == t), nil
//	default:
//		return False, nil
//	}
//}
//
//func (f *nativeFixedFunc) Invoke(cx Context, params []Value) (Value, Error) {
//
//	err := VetFixedFuncParams(f.requiredTypes, f.allowNull, params)
//	if err != nil {
//		return nil, err
//	}
//
//	return f.invoke(cx, params)
//}
//
////--------------------------------------------------------------
//// nativeVariadicFunc
////--------------------------------------------------------------
//
//type nativeVariadicFunc struct {
//	*nativeBaseFunc
//	requiredTypes []Type
//	allowNull     bool
//	variadicType  Type
//}
//
//// NewVariadicNativeFunc is a convenience function for creating
//// a NativeFunc with variadic arity
//func NewVariadicNativeFunc(
//	requiredTypes []Type,
//	variadicType Type,
//	allowNull bool,
//	invoke Invoke) NativeFunc {
//
//	arity := Arity{
//		Kind:           VariadicArity,
//		RequiredParams: uint16(len(requiredTypes)),
//	}
//
//	return &nativeVariadicFunc{
//		&nativeBaseFunc{arity, invoke},
//		requiredTypes, allowNull, variadicType,
//	}
//}
//
//func (f *nativeVariadicFunc) Freeze(cx Context) (Value, Error) {
//	return f, nil
//}
//
//func (f *nativeVariadicFunc) Eq(cx Context, v Value) (Bool, Error) {
//	switch t := v.(type) {
//	case *nativeVariadicFunc:
//		// equality is based on identity
//		return NewBool(f == t), nil
//	default:
//		return False, nil
//	}
//}
//
//func (f *nativeVariadicFunc) Invoke(cx Context, params []Value) (Value, Error) {
//
//	numValues := len(params)
//	numReqs := len(f.requiredTypes)
//
//	// arity mismatch
//	if numValues < numReqs {
//		return nil, ArityMismatchError(
//			fmt.Sprintf("at least %d", numReqs),
//			numValues)
//	}
//
//	// check types on required params
//	for i := 0; i < numReqs; i++ {
//		err := vetFuncParam(f.requiredTypes[i], f.allowNull, params[i])
//		if err != nil {
//			return nil, err
//		}
//	}
//
//	// check types on variadic params
//	for i := numReqs; i < numValues; i++ {
//		err := vetFuncParam(f.variadicType, f.allowNull, params[i])
//		if err != nil {
//			return nil, err
//		}
//	}
//
//	// invoke
//	return f.invoke(cx, params)
//}
//
////--------------------------------------------------------------
//// nativeMultipleFunc
////--------------------------------------------------------------
//
//type nativeMultipleFunc struct {
//	*nativeBaseFunc
//	requiredTypes []Type
//	allowNull     bool
//	optionalTypes []Type
//}
//
//// NewMultipleNativeFunc is a convenience function for creating
//// a NativeFunc with multiple arity
//func NewMultipleNativeFunc(
//	requiredTypes []Type,
//	optionalTypes []Type,
//	allowNull bool,
//	invoke Invoke) NativeFunc {
//
//	arity := Arity{
//		Kind:           MultipleArity,
//		RequiredParams: uint16(len(requiredTypes)),
//		OptionalParams: uint16(len(optionalTypes)),
//	}
//
//	return &nativeMultipleFunc{
//		&nativeBaseFunc{arity, invoke},
//		requiredTypes, allowNull, optionalTypes,
//	}
//}
//
//func (f *nativeMultipleFunc) Freeze(cx Context) (Value, Error) {
//	return f, nil
//}
//
//func (f *nativeMultipleFunc) Eq(cx Context, v Value) (Bool, Error) {
//	switch t := v.(type) {
//	case *nativeMultipleFunc:
//		// equality is based on identity
//		return NewBool(f == t), nil
//	default:
//		return False, nil
//	}
//}
//
//func (f *nativeMultipleFunc) Invoke(cx Context, params []Value) (Value, Error) {
//
//	numValues := len(params)
//	numReqs := len(f.requiredTypes)
//	numOpts := len(f.optionalTypes)
//
//	// arity mismatch
//	if numValues < numReqs {
//		return nil, ArityMismatchError(
//			fmt.Sprintf("at least %d", numReqs),
//			numValues)
//	}
//	if numValues > (numReqs + numOpts) {
//		return nil, ArityMismatchError(
//			fmt.Sprintf("at most %d", numReqs+numOpts),
//			numValues)
//	}
//
//	// check types on required params
//	for i := 0; i < numReqs; i++ {
//		err := vetFuncParam(f.requiredTypes[i], f.allowNull, params[i])
//		if err != nil {
//			return nil, err
//		}
//	}
//
//	// check types on optional params
//	for i := numReqs; i < numValues; i++ {
//		err := vetFuncParam(f.optionalTypes[i-numReqs], f.allowNull, params[i])
//		if err != nil {
//			return nil, err
//		}
//	}
//
//	// invoke
//	return f.invoke(cx, params)
//}
//
////--------------------------------------------------------------
//// virtualFunc
////--------------------------------------------------------------
//
//// A virtual function is a function that is created only if
//// we need to pass it around as a first class value.
//// The 'same' virtual func can be instantiated more than once,
//// so it has a special Eq() implementation.
//type virtualFunc struct {
//	owner Value
//	name  string
//	*nativeBaseFunc
//}
//
//func newVirtualFunc(owner Value, name string, arity Arity, invoke Invoke) NativeFunc {
//	return &virtualFunc{
//		owner, name,
//		&nativeBaseFunc{arity, invoke},
//	}
//}
//
//func (f *virtualFunc) Freeze(cx Context) (Value, Error) {
//	return f, nil
//}
//
//func (f *virtualFunc) Eq(cx Context, v Value) (Bool, Error) {
//	switch t := v.(type) {
//	case *virtualFunc:
//		// equality for virtual functions is based on whether
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
//
//func (f *virtualFunc) Invoke(cx Context, params []Value) (Value, Error) {
//	return f.invoke(cx, params)
//}
//
////--------------------------------------------------------------
//// Vet
////--------------------------------------------------------------
//
//func VetFixedFuncParams(requiredTypes []Type, allowNull bool, params []Value) Error {
//
//	numValues := len(params)
//	numReqs := len(requiredTypes)
//
//	// arity mismatch
//	if numValues != numReqs {
//		return ArityMismatchError(
//			fmt.Sprintf("%d", numReqs),
//			numValues)
//	}
//
//	// check types on required params
//	for i := 0; i < numReqs; i++ {
//		err := vetFuncParam(requiredTypes[i], allowNull, params[i])
//		if err != nil {
//			return err
//		}
//	}
//
//	return nil
//}
//
//func vetFuncParam(typ Type, allowNull bool, param Value) Error {
//
//	// accept 'any' type
//	if typ == AnyType {
//		return nil
//	}
//
//	// skip over null params
//	if param.Type() == NullType {
//		if allowNull {
//			return nil
//		}
//		return NullValueError()
//	}
//
//	// check type
//	if param.Type() != typ {
//		return TypeMismatchError(fmt.Sprintf("Expected %s", typ.String()))
//	}
//
//	// invoke
//	return nil
//}
//
//// NewFixedFuncDef
//func NewFixedFuncDef(
//	requiredTypes []Type,
//	allowNull bool,
//	invoke Invoke) (arity Arity, wrapper Invoke) {
//
//	arity = Arity{
//		Kind:           FixedArity,
//		RequiredParams: uint16(len(requiredTypes)),
//	}
//
//	wrapper = func(cx Context, params []Value) (Value, Error) {
//		err := VetFixedFuncParams(requiredTypes, allowNull, params)
//		if err != nil {
//			return nil, err
//		}
//		return invoke(cx, params)
//	}
//
//	return
//}
