// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package core

// Field is a name-value pair inside a Struct
type Field interface {
	Name() string
}

type field struct {
	name       string
	isReadonly bool
	isProperty bool
	value      Value
}

// Name returns the name of a field
func (f *field) Name() string {
	return f.name
}

// NewField a name-value pair.
func NewField(name string, isReadonly bool, value Value) Field {
	return &field{name, isReadonly, false, value}
}

// NewReadonlyProperty creates a readonly Property using a 'getter' function.
func NewReadonlyProperty(name string, getter Func) Field {

	if getter.MinArity() != 0 || getter.MaxArity() != 0 {
		panic("Property getter does not have arity 0")
	}

	prop := NewTuple([]Value{getter, nil})
	return &field{name, true, true, prop}
}

// NewProperty creates a Property using 'getter' and 'setter' functions.
func NewProperty(name string, getter Func, setter Func) Field {

	if getter.MinArity() != 0 || getter.MaxArity() != 0 {
		panic("Property getter does not have arity 0")
	}

	if setter.MinArity() != 1 || setter.MaxArity() != 1 {
		panic("Property setter does not have arity 1")
	}

	prop := NewTuple([]Value{getter, setter})
	return &field{name, false, true, prop}
}
