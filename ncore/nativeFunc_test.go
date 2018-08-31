// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package ncore

import (
	"testing"
)

func TestFixedNativeFunc(t *testing.T) {

	fn := NewFixedNativeFunc(
		[]Type{},
		false,
		func(ev Evaluator, params []Value) (Value, Error) {
			return Zero, nil
		})

	ok(t, fn.Arity(), nil, Arity{FixedArity, 0, 0})

	val, err := fn.Invoke(nil, []Value{})
	ok(t, val, err, Zero)

	val, err = fn.Invoke(nil, []Value{One})
	fail(t, val, err, "ArityMismatch: Expected 0 params, got 1")

	//----------------------------------------------

	fn = NewFixedNativeFunc(
		[]Type{IntType},
		false,
		func(ev Evaluator, params []Value) (Value, Error) {
			n := params[0].(Int)
			return n.Add(One)
		})

	ok(t, fn.Arity(), nil, Arity{FixedArity, 1, 0})

	val, err = fn.Invoke(nil, []Value{NewInt(3)})
	ok(t, val, err, NewInt(4))

	val, err = fn.Invoke(nil, []Value{Null})
	fail(t, val, err, "NullValue")

	val, err = fn.Invoke(nil, []Value{})
	fail(t, val, err, "ArityMismatch: Expected 1 params, got 0")

	val, err = fn.Invoke(nil, []Value{NewStr("a")})
	fail(t, val, err, "TypeMismatch: Expected Int")

	//----------------------------------------------

	fn = NewFixedNativeFunc(
		[]Type{IntType},
		true,
		func(ev Evaluator, params []Value) (Value, Error) {
			if params[0] == Null {
				return True, nil
			}
			return False, nil
		})
	ok(t, fn.Arity(), nil, Arity{FixedArity, 1, 0})

	val, err = fn.Invoke(nil, []Value{Zero})
	ok(t, val, err, False)

	val, err = fn.Invoke(nil, []Value{Null})
	ok(t, val, err, True)

	val, err = fn.Invoke(nil, []Value{})
	fail(t, val, err, "ArityMismatch: Expected 1 params, got 0")

	val, err = fn.Invoke(nil, []Value{NewStr("a")})
	fail(t, val, err, "TypeMismatch: Expected Int")
}

func TestVariadicNativeFunc(t *testing.T) {

	fn := NewVariadicNativeFunc(
		[]Type{},
		AnyType,
		true,
		func(ev Evaluator, params []Value) (Value, Error) {
			return NewInt(int64(len(params))), nil
		})

	ok(t, fn.Arity(), nil, Arity{VariadicArity, 0, 0})

	val, err := fn.Invoke(nil, []Value{})
	ok(t, val, err, Zero)

	val, err = fn.Invoke(nil, []Value{True, Null})
	ok(t, val, err, NewInt(2))

	//----------------------------------------------

	fn = NewVariadicNativeFunc(
		[]Type{IntType},
		BoolType,
		false,
		func(ev Evaluator, params []Value) (Value, Error) {
			return NewInt(int64(len(params))), nil
		})
	ok(t, fn.Arity(), nil, Arity{VariadicArity, 1, 0})

	val, err = fn.Invoke(nil, []Value{})
	fail(t, val, err, "ArityMismatch: Expected at least 1 params, got 0")

	val, err = fn.Invoke(nil, []Value{NewStr("a")})
	fail(t, val, err, "TypeMismatch: Expected Int")

	val, err = fn.Invoke(nil, []Value{Null})
	fail(t, val, err, "NullValue")

	val, err = fn.Invoke(nil, []Value{Zero, NewStr("a"), False})
	fail(t, val, err, "TypeMismatch: Expected Bool")

	val, err = fn.Invoke(nil, []Value{Zero, Null, False})
	fail(t, val, err, "NullValue")

	val, err = fn.Invoke(nil, []Value{Zero, True, False})
	ok(t, val, err, NewInt(3))
}

func TestMultipleNativeFunc(t *testing.T) {

	fn := NewMultipleNativeFunc(
		[]Type{IntType},
		[]Type{StrType, BoolType},
		false,
		func(ev Evaluator, params []Value) (Value, Error) {

			if len(params) == 1 {
				params = append(params, NewStr("a"))
			}

			if len(params) == 2 {
				params = append(params, False)
			}

			str := NewStr("")
			for _, v := range params {
				val, err := v.ToStr(ev)
				if err != nil {
					return nil, err
				}

				str = str.Concat(val)
			}
			return str, nil
		})

	ok(t, fn.Arity(), nil, Arity{MultipleArity, 1, 2})

	val, err := fn.Invoke(nil, []Value{})
	fail(t, val, err, "ArityMismatch: Expected at least 1 params, got 0")

	val, err = fn.Invoke(nil, []Value{Zero, Zero, Zero, Zero})
	fail(t, val, err, "ArityMismatch: Expected at most 3 params, got 4")

	val, err = fn.Invoke(nil, []Value{Zero, Zero, NewStr("d")})
	fail(t, val, err, "TypeMismatch: Expected Str")

	val, err = fn.Invoke(nil, []Value{Zero, NewStr("c"), Null})
	fail(t, val, err, "NullValue")

	val, err = fn.Invoke(nil, []Value{Zero, NewStr("c"), Zero})
	fail(t, val, err, "TypeMismatch: Expected Bool")

	val, err = fn.Invoke(nil, []Value{Zero, NewStr("c"), True})
	ok(t, val, err, NewStr("0ctrue"))

	val, err = fn.Invoke(nil, []Value{Zero, NewStr("c")})
	ok(t, val, err, NewStr("0cfalse"))

	val, err = fn.Invoke(nil, []Value{Zero})
	ok(t, val, err, NewStr("0afalse"))
}
