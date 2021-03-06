// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package core

import (
	"testing"
)

func TestFixedNativeFunc(t *testing.T) {

	fn := NewFixedNativeFunc(
		[]Type{},
		false,
		func(ev Eval, params []Value) (Value, Error) {
			return Zero, nil
		})

	ok(t, fn.Arity(), nil, Arity{FixedArity, 0, 0})

	val, err := fn.Invoke(nil, []Value{})
	ok(t, val, err, Zero)

	val, err = fn.Invoke(nil, []Value{One})
	fail(t, val, err, "ArityMismatch: Expected 0 parameters, got 1")

	//----------------------------------------------

	fn = NewFixedNativeFunc(
		[]Type{IntType},
		false,
		func(ev Eval, params []Value) (Value, Error) {
			n := params[0].(Int)
			return n.Add(One), nil
		})

	ok(t, fn.Arity(), nil, Arity{FixedArity, 1, 0})

	val, err = fn.Invoke(nil, []Value{NewInt(3)})
	ok(t, val, err, NewInt(4))

	val, err = fn.Invoke(nil, []Value{Null})
	fail(t, val, err, "NullValue")

	val, err = fn.Invoke(nil, []Value{})
	fail(t, val, err, "ArityMismatch: Expected 1 parameter, got 0")

	val, err = fn.Invoke(nil, []Value{MustStr("a")})
	fail(t, val, err, "TypeMismatch: Expected Int, not Str")

	//----------------------------------------------

	fn = NewFixedNativeFunc(
		[]Type{IntType},
		true,
		func(ev Eval, params []Value) (Value, Error) {
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
	fail(t, val, err, "ArityMismatch: Expected 1 parameter, got 0")

	val, err = fn.Invoke(nil, []Value{MustStr("a")})
	fail(t, val, err, "TypeMismatch: Expected Int, not Str")
}

func TestVariadicNativeFunc(t *testing.T) {

	fn := NewVariadicNativeFunc(
		[]Type{},
		AnyType,
		true,
		func(ev Eval, params []Value) (Value, Error) {
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
		func(ev Eval, params []Value) (Value, Error) {
			return NewInt(int64(len(params))), nil
		})
	ok(t, fn.Arity(), nil, Arity{VariadicArity, 1, 0})

	val, err = fn.Invoke(nil, []Value{})
	fail(t, val, err, "ArityMismatch: Expected at least 1 parameter, got 0")

	val, err = fn.Invoke(nil, []Value{MustStr("a")})
	fail(t, val, err, "TypeMismatch: Expected Int, not Str")

	val, err = fn.Invoke(nil, []Value{Null})
	fail(t, val, err, "NullValue")

	val, err = fn.Invoke(nil, []Value{Zero, MustStr("a"), False})
	fail(t, val, err, "TypeMismatch: Expected Bool, not Str")

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
		func(ev Eval, params []Value) (Value, Error) {

			if len(params) == 1 {
				params = append(params, MustStr("a"))
			}

			if len(params) == 2 {
				params = append(params, False)
			}

			s := MustStr("")
			for _, v := range params {
				val, err := v.ToStr(ev)
				if err != nil {
					return nil, err
				}

				s = s.Concat(val)
			}
			return s, nil
		})

	ok(t, fn.Arity(), nil, Arity{MultipleArity, 1, 2})

	val, err := fn.Invoke(nil, []Value{})
	fail(t, val, err, "ArityMismatch: Expected at least 1 parameter, got 0")

	val, err = fn.Invoke(nil, []Value{Zero, Zero, Zero, Zero})
	fail(t, val, err, "ArityMismatch: Expected at most 3 parameters, got 4")

	val, err = fn.Invoke(nil, []Value{Zero, Zero, MustStr("d")})
	fail(t, val, err, "TypeMismatch: Expected Str, not Int")

	val, err = fn.Invoke(nil, []Value{Zero, MustStr("c"), Null})
	fail(t, val, err, "NullValue")

	val, err = fn.Invoke(nil, []Value{Zero, MustStr("c"), Zero})
	fail(t, val, err, "TypeMismatch: Expected Bool, not Int")

	val, err = fn.Invoke(nil, []Value{Zero, MustStr("c"), True})
	ok(t, val, err, MustStr("0ctrue"))

	val, err = fn.Invoke(nil, []Value{Zero, MustStr("c")})
	ok(t, val, err, MustStr("0cfalse"))

	val, err = fn.Invoke(nil, []Value{Zero})
	ok(t, val, err, MustStr("0afalse"))
}
