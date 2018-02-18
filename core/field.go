// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package core

import (
	"fmt"
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

	// PropertyGetter is the signature for a native 'getter' function.
	PropertyGetter func(Context) (Value, Error)

	// PropertySetter is the signature for a native 'setter' function.
	PropertySetter func(Context, Value) (Value, Error)
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

// NewNativeProperty creates a Property using 'getter' and 'setter' functions.
// If 'setter' is nil, the Property will be readonly.
func NewNativeProperty(name string, getter PropertyGetter, setter PropertySetter) Field {

	get := NewNativeFunc(0, 0,
		func(cx Context, values []Value) (Value, Error) {
			return getter(cx)
		})
	if setter == nil {
		return &field{name, true, true, NewTuple([]Value{get, nil})}
	}

	set := NewNativeFunc(1, 1,
		func(cx Context, values []Value) (Value, Error) {
			return setter(cx, values[0])
		})
	return &field{name, false, true, NewTuple([]Value{get, set})}
}

func (fd *FieldDef) String() string {
	return fmt.Sprintf("fieldDef(%s %v %v)", fd.Name, fd.IsReadonly, fd.IsProperty)
}
