// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package core

import (
	"testing"
)

func TestField(t *testing.T) {

	field := NewField(Zero)

	val, err := field.Get(nil)
	ok(t, val, err, Zero)

	val, err = field.Invoke(nil, []Value{})
	fail(t, val, err, "TypeMismatch: Expected Func, not Int")

	Tassert(t, !field.IsReadonly())
	err = field.Set(nil, One)
	Tassert(t, err == nil)

	val, err = field.Get(nil)
	ok(t, val, err, One)

	val, err = field.Invoke(nil, []Value{})
	fail(t, val, err, "TypeMismatch: Expected Func, not Int")

	//----------------------------

	fn := NewFixedNativeFunc(
		[]Type{IntType}, false,
		func(ev Evaluator, values []Value) (Value, Error) {
			n := values[0].(Int)
			return n.Add(One), nil
		})
	field = NewField(fn)

	val, err = field.Get(nil)
	Tassert(t, err == nil && val == fn)

	val, err = field.Invoke(nil, []Value{Zero})
	ok(t, val, err, One)

	Tassert(t, !field.IsReadonly())
	err = field.Set(nil, One)
	Tassert(t, err == nil)

	val, err = field.Get(nil)
	ok(t, val, err, One)

	val, err = field.Invoke(nil, []Value{})
	fail(t, val, err, "TypeMismatch: Expected Func, not Int")
}

func TestReadonlyField(t *testing.T) {

	field := NewReadonlyField(Zero)

	val, err := field.Get(nil)
	ok(t, val, err, Zero)

	val, err = field.Invoke(nil, []Value{})
	fail(t, val, err, "TypeMismatch: Expected Func, not Int")

	Tassert(t, field.IsReadonly())
	err = field.Set(nil, One)
	fail(t, nil, err, "ReadonlyField")

	//----------------------------

	fn := NewFixedNativeFunc(
		[]Type{IntType}, false,
		func(ev Evaluator, values []Value) (Value, Error) {
			n := values[0].(Int)
			return n.Add(One), nil
		})
	field = NewReadonlyField(fn)

	val, err = field.Get(nil)
	Tassert(t, err == nil && val == fn)

	val, err = field.Invoke(nil, []Value{Zero})
	ok(t, val, err, One)

	Tassert(t, field.IsReadonly())
	err = field.Set(nil, One)
	fail(t, nil, err, "ReadonlyField")
}

func TestPropertyField(t *testing.T) {

	var propVal Value = Zero

	get := NewFixedNativeFunc(
		[]Type{}, false,
		func(ev Evaluator, values []Value) (Value, Error) {
			return propVal, nil
		})

	set := NewFixedNativeFunc(
		[]Type{IntType}, false,
		func(ev Evaluator, values []Value) (Value, Error) {
			propVal = values[0]
			return Null, nil
		})

	bogus := NewFixedNativeFunc(
		[]Type{AnyType, AnyType}, false,
		func(ev Evaluator, values []Value) (Value, Error) {
			panic("unreachable")
		})

	_, err := NewProperty(bogus, set)
	fail(t, nil, err, "InvalidGetterArity: Arity(Fixed,2,0)")

	_, err = NewProperty(get, bogus)
	fail(t, nil, err, "InvalidSetterArity: Arity(Fixed,2,0)")

	field, err := NewProperty(get, set)
	Tassert(t, err == nil)

	val, err := field.Get(nil)
	ok(t, val, err, Zero)

	val, err = field.Invoke(nil, []Value{})
	fail(t, val, err, "TypeMismatch: Expected Func, not Int")

	Tassert(t, !field.IsReadonly())
	err = field.Set(nil, One)
	Tassert(t, err == nil)

	val, err = field.Get(nil)
	ok(t, val, err, One)

	val, err = field.Invoke(nil, []Value{})
	fail(t, val, err, "TypeMismatch: Expected Func, not Int")

	//----------------------------

	propVal = NewFixedNativeFunc(
		[]Type{IntType}, false,
		func(ev Evaluator, values []Value) (Value, Error) {
			n := values[0].(Int)
			return n.Add(One), nil
		})

	val, err = field.Get(nil)
	Tassert(t, err == nil && val == propVal)

	val, err = field.Invoke(nil, []Value{Zero})
	ok(t, val, err, One)

	Tassert(t, !field.IsReadonly())
	err = field.Set(nil, One)
	Tassert(t, err == nil)

	val, err = field.Get(nil)
	ok(t, val, err, One)

	val, err = field.Invoke(nil, []Value{})
	fail(t, val, err, "TypeMismatch: Expected Func, not Int")
}

func TestReadonlyPropertyField(t *testing.T) {

	var propVal Value = Zero

	get := NewFixedNativeFunc(
		[]Type{}, false,
		func(ev Evaluator, values []Value) (Value, Error) {
			return propVal, nil
		})

	bogus := NewFixedNativeFunc(
		[]Type{AnyType, AnyType}, false,
		func(ev Evaluator, values []Value) (Value, Error) {
			panic("unreachable")
		})

	_, err := NewReadonlyProperty(bogus)
	fail(t, nil, err, "InvalidGetterArity: Arity(Fixed,2,0)")

	field, err := NewReadonlyProperty(get)
	Tassert(t, err == nil)

	val, err := field.Get(nil)
	ok(t, val, err, Zero)

	val, err = field.Invoke(nil, []Value{})
	fail(t, val, err, "TypeMismatch: Expected Func, not Int")

	Tassert(t, field.IsReadonly())
	err = field.Set(nil, One)
	fail(t, nil, err, "ReadonlyField")

	//----------------------------

	propVal = NewFixedNativeFunc(
		[]Type{IntType}, false,
		func(ev Evaluator, values []Value) (Value, Error) {
			n := values[0].(Int)
			return n.Add(One), nil
		})

	val, err = field.Get(nil)
	Tassert(t, err == nil && val == propVal)

	val, err = field.Invoke(nil, []Value{Zero})
	ok(t, val, err, One)

	Tassert(t, field.IsReadonly())
	err = field.Set(nil, One)
	fail(t, nil, err, "ReadonlyField")
}
