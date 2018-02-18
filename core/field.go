// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package core

type (
	// Field is a name-value pair inside a Struct
	Field interface {
		Name() string
	}

	field struct {
		name       string
		isReadonly bool
		isProperty bool
		value      Value
	}

	PropertyGetter func(Context) (Value, Error)
	PropertySetter func(Context, Value) (Value, Error)
)

// Name returns the name of a field
func (f *field) Name() string {
	return f.name
}

// NewField a name-value pair.
func NewField(name string, isReadonly bool, value Value) Field {
	return &field{name, isReadonly, false, value}
}

//// NewProperty creates a Property using 'getter' and 'setter' functions.
//func NewProperty(name string, getter Func, setter Func) Field {
//
//	if getter.MinArity() != 0 || getter.MaxArity() != 0 {
//		panic("Property getter does not have arity 0")
//	}
//
//	if setter.MinArity() != 1 || setter.MaxArity() != 1 {
//		panic("Property setter does not have arity 1")
//	}
//
//	prop := NewTuple([]Value{getter, setter})
//	return &field{name, false, true, prop}
//}

// NewNativeProperty creates a Property using 'getter' and 'setter' functions.
// If 'setter' is nil, the Property will be readonly.
func NewNativeProperty(name string, getter PropertyGetter, setter PropertySetter) Field {

	get := NewNativeFunc(0, 0,
		func(cx Context, values []Value) (Value, Error) {
			return getter(cx)
		})
	if setter == nil {
		prop := NewTuple([]Value{get, nil})
		return &field{name, true, true, prop}

	}

	set := NewNativeFunc(1, 1,
		func(cx Context, values []Value) (Value, Error) {
			return setter(cx, values[0])
		})
	prop := NewTuple([]Value{get, set})
	return &field{name, false, true, prop}
}
