// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package core

import (
	"fmt"
	"reflect"
)

type (
	// Field is a name-value pair in a Struct
	Field interface {
		Name() string
		fieldMarker()
	}

	// FieldDef defines a name-value pair in a Struct
	FieldDef struct {
		Name       string
		IsReadonly bool
		IsProperty bool
	}

	field struct {
		name       string
		isReadonly bool
		isProperty bool
		value      Value
	}
)

// Name returns the name of a field
func (f *field) Name() string {
	return f.name
}

func (f *field) fieldMarker() {}

// NewField creates a name-value pair.
func NewField(name string, isReadonly bool, value Value) Field {
	return &field{name, isReadonly, false, value}
}

var getterArity = &Arity{Kind: FixedArity, RequiredParams: 0, OptionalParams: nil}
var setterArity = &Arity{Kind: FixedArity, RequiredParams: 1, OptionalParams: nil}

// NewReadonlyNativeProperty creates a readonly Property using a 'getter' function.
// The 'getter' function must have an arity of 0.
func NewReadonlyNativeProperty(name string, getter NativeFunc) (Field, Error) {

	if !reflect.DeepEqual(getterArity, getter.Arity()) {
		return nil, ArityMismatchError("0", getter.Arity().RequiredParams)
	}

	return &field{name, true, true, NewTuple([]Value{getter, Null})}, nil
}

// NewNativeProperty creates a Property using 'getter' and 'setter' functions.
// The 'getter' function must have an arity of 0, and the 'setter' function
// must have an arity of 1.  By convention the setter function should
// return 'Null'; its return value will be ignored.
func NewNativeProperty(name string, getter NativeFunc, setter NativeFunc) (Field, Error) {

	if !reflect.DeepEqual(getterArity, getter.Arity()) {
		return nil, ArityMismatchError("0", getter.Arity().RequiredParams)
	}

	if !reflect.DeepEqual(setterArity, setter.Arity()) {
		return nil, ArityMismatchError("0", setter.Arity().RequiredParams)
	}

	return &field{name, false, true, NewTuple([]Value{getter, setter})}, nil
}

func (fd *FieldDef) String() string {
	return fmt.Sprintf("fieldDef(%s %v %v)", fd.Name, fd.IsReadonly, fd.IsProperty)
}
