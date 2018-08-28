// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package ncore

import (
	"testing"
)

func TestField(t *testing.T) {

	var num Value = NewInt(3)
	getter := func() (Value, Error) {
		return num, nil
	}
	setter := func(val Value) Error {
		num = val
		return nil
	}

	field := NewField("foo", getter, setter)
	tassert(t, field.Name() == "foo")

	val, err := field.Getter()()
	ok(t, val, err, NewInt(3))

	err = field.Setter()(NewInt(4))
	tassert(t, err == nil)

	val, err = field.Getter()()
	ok(t, val, err, NewInt(4))
}

func TestReadonlyField(t *testing.T) {

	var num Value = NewInt(3)
	getter := func() (Value, Error) {
		return num, nil
	}

	field := NewReadonlyField("foo", getter)
	tassert(t, field.Name() == "foo")

	val, err := field.Getter()()
	ok(t, val, err, NewInt(3))

	err = field.Setter()(NewInt(4))
	fail(t, nil, err, "ReadonlyField: Field 'foo' is readonly")
}
