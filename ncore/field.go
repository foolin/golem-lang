// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package ncore

//--------------------------------------------------------------
// Field
//--------------------------------------------------------------

type (
	Getter func() (Value, Error)
	Setter func(Value) Error

	Field interface {
		Name() string
		Getter() Getter
		Setter() Setter
	}

	field struct {
		name   string
		getter Getter
		setter Setter
	}
)

func NewField(name string, getter Getter, setter Setter) Field {
	return &field{name, getter, setter}
}

func NewReadonlyField(name string, getter Getter) Field {
	return &field{
		name,
		getter,
		func(val Value) Error {
			return ReadonlyFieldError(name)
		},
	}
}

func (f *field) Name() string   { return f.name }
func (f *field) Getter() Getter { return f.getter }
func (f *field) Setter() Setter { return f.setter }

////--------------------------------------------------------------
//// Fields
////--------------------------------------------------------------
//
//type (
//	Fields interface {
//		Names() []string
//		Get(string) (Field, bool)
//	}
//)
