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
	isConst    bool
	isProperty bool
	value      Value
}

// Name returns the name of a field
func (f *field) Name() string {
	return f.name
}

// NewField a name-value pair.
func NewField(name string, isConst bool, value Value) Field {
	return &field{name, isConst, false, value}
}

// NewProperty creates a Property using 'getter' and 'setter' functions.
// If the setter is nil, then the property is const.
func NewProperty(name string, getter Func, setter Func) Field {

	if getter.MinArity() != 0 || getter.MaxArity() != 0 {
		panic("Property getter does not have arity 0")
	}

	if setter != nil {
		if setter.MinArity() != 1 || setter.MaxArity() != 1 {
			panic("Property setter does not have arity 1")
		}
	}

	prop := NewTuple([]Value{getter, setter})
	return &field{name, setter == nil, true, prop}
}
