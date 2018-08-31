// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package ncore

import (
	"testing"
)

func testMethodEq(t *testing.T, m Method) {
	a1 := m.ToFunc(1, "foo")
	a2 := m.ToFunc(1, "foo")
	a3 := m.ToFunc(2, "foo")

	val, err := a1.Eq(nil, a2)
	ok(t, val, err, True)

	val, err = a1.Eq(nil, a3)
	ok(t, val, err, False)

	val, err = a2.Eq(nil, a3)
	ok(t, val, err, False)

	b1 := m.ToFunc(NewInt(1), "foo")
	b2 := m.ToFunc(NewInt(1), "foo")
	b3 := m.ToFunc(NewInt(2), "foo")

	val, err = b1.Eq(nil, b2)
	ok(t, val, err, True)

	val, err = b1.Eq(nil, b3)
	ok(t, val, err, False)

	val, err = b2.Eq(nil, b3)
	ok(t, val, err, False)

	val, err = a1.Eq(nil, b1)
	ok(t, val, err, False)
}

func TestFixedMethod(t *testing.T) {

	m := NewFixedMethod(
		[]Type{},
		false,
		func(self interface{}, ev Evaluator, params []Value) (Value, Error) {
			n := self.(int)
			return NewInt(int64(n * n)), nil
		})
	testMethodEq(t, m)

	val, err := m.Invoke(10, nil, []Value{})
	ok(t, val, err, NewInt(100))

	fn := m.ToFunc(10, "foo")

	ok(t, fn.Arity(), nil, Arity{FixedArity, 0, 0})

	val, err = fn.Invoke(nil, []Value{})
	ok(t, val, err, NewInt(100))

	val, err = fn.Invoke(nil, []Value{One})
	fail(t, val, err, "ArityMismatch: Expected 0 params, got 1")

	//----------------------------------------------

	m = NewFixedMethod(
		[]Type{IntType},
		false,
		func(self interface{}, ev Evaluator, params []Value) (Value, Error) {
			n := self.(Int)
			p := params[0].(Int)
			return n.Add(p)
		})

	val, err = m.Invoke(NewInt(2), nil, []Value{NewInt(3)})
	ok(t, val, err, NewInt(5))

	fn = m.ToFunc(NewInt(3), "foo")

	ok(t, fn.Arity(), nil, Arity{FixedArity, 1, 0})

	val, err = fn.Invoke(nil, []Value{NewInt(3)})
	ok(t, val, err, NewInt(6))

	val, err = fn.Invoke(nil, []Value{Null})
	fail(t, val, err, "NullValue")

	val, err = fn.Invoke(nil, []Value{})
	fail(t, val, err, "ArityMismatch: Expected 1 params, got 0")

	val, err = fn.Invoke(nil, []Value{NewStr("a")})
	fail(t, val, err, "TypeMismatch: Expected Int")
}

func TestVariadicMethod(t *testing.T) {

	m := NewVariadicMethod(
		[]Type{},
		IntType,
		false,
		func(self interface{}, ev Evaluator, params []Value) (Value, Error) {
			n := self.(int)
			for _, p := range params {
				n += int(p.(Int).IntVal())
			}
			return NewInt(int64(n)), nil
		})
	testMethodEq(t, m)

	val, err := m.Invoke(10, nil, []Value{NewInt(2), NewInt(3)})
	ok(t, val, err, NewInt(15))

	fn := m.ToFunc(10, "foo")

	ok(t, fn.Arity(), nil, Arity{VariadicArity, 0, 0})

	val, err = fn.Invoke(nil, []Value{NewInt(2), NewInt(3)})
	ok(t, val, err, NewInt(15))
}

func TestMultipleMethod(t *testing.T) {

	m := NewMultipleMethod(
		[]Type{IntType},
		[]Type{IntType},
		false,
		func(self interface{}, ev Evaluator, params []Value) (Value, Error) {
			n := self.(int)
			n += int(params[0].(Int).IntVal())
			if len(params) == 2 {
				n += int(params[1].(Int).IntVal())
			}

			return NewInt(int64(n)), nil
		})
	testMethodEq(t, m)

	val, err := m.Invoke(10, nil, []Value{NewInt(2)})
	ok(t, val, err, NewInt(12))

	val, err = m.Invoke(10, nil, []Value{NewInt(2), NewInt(3)})
	ok(t, val, err, NewInt(15))

	fn := m.ToFunc(10, "foo")

	ok(t, fn.Arity(), nil, Arity{MultipleArity, 1, 1})

	val, err = fn.Invoke(nil, []Value{NewInt(2)})
	ok(t, val, err, NewInt(12))

	val, err = fn.Invoke(nil, []Value{NewInt(2), NewInt(3)})
	ok(t, val, err, NewInt(15))
}

//--------------------------------------------------------------

func show(val Value) {
	//println(val.ToStr(nil).String())
}

const iterate = 2 * 1000 * 1000

//const iterate = 4

func TestBenchmarkDirectInvoke(t *testing.T) {

	var i1 Int = NewInt(1)
	var i2 Int = NewInt(2)

	for i := 0; i < iterate; i++ {
		val, _ := i1.Add(i2)
		show(val)
	}
}

func TestBenchmarkFuncInvoke(t *testing.T) {

	fn := NewFixedNativeFunc(
		[]Type{IntType, IntType},
		false,
		func(ev Evaluator, params []Value) (Value, Error) {
			a := params[0].(Int)
			b := params[1].(Int)
			return a.Add(b)
		})

	var i1 Int = NewInt(1)
	var i2 Int = NewInt(2)

	for i := 0; i < iterate; i++ {
		val, _ := fn.Invoke(nil, []Value{i1, i2})
		show(val)
	}
}

func TestBenchmarkMethodInvoke(t *testing.T) {

	m := NewFixedMethod(
		[]Type{IntType},
		false,
		func(self interface{}, ev Evaluator, params []Value) (Value, Error) {
			n := self.(Int)
			p := params[0].(Int)
			return n.Add(p)
		})

	var i1 Int = NewInt(1)
	var i2 Int = NewInt(2)

	for i := 0; i < iterate; i++ {
		val, _ := m.Invoke(i1, nil, []Value{i2})
		show(val)
	}
}

func TestBenchmarkVirtualFuncInvoke(t *testing.T) {

	m := NewFixedMethod(
		[]Type{IntType},
		false,
		func(self interface{}, ev Evaluator, params []Value) (Value, Error) {
			n := self.(Int)
			p := params[0].(Int)
			return n.Add(p)
		})

	var i1 Int = NewInt(1)
	var i2 Int = NewInt(2)

	for i := 0; i < iterate; i++ {
		fn := m.ToFunc(i1, "foo")
		val, _ := fn.Invoke(nil, []Value{i2})
		show(val)
	}
}
