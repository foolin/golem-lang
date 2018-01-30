// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package core

import (
	//"fmt"
	"reflect"
	"testing"
)

var cx Context = nil

func tassert(t *testing.T, flag bool) {
	if !flag {
		t.Error("assertion failure")
	}
}

func ok(t *testing.T, val Value, err Error, expect Value) {

	if err != nil {
		t.Error(err, " != ", nil)
	}

	if !reflect.DeepEqual(val, expect) {
		t.Error(val, " != ", expect)
	}
}

func fail(t *testing.T, val Value, err Error, expect string) {

	if val != nil {
		t.Error(val, " != ", nil)
	}

	if err == nil || err.Error() != expect {
		t.Error(err.Error(), " != ", expect)
	}
}

func okType(t *testing.T, val Value, expected Type) {
	tassert(t, val.Type() == expected)
}

func TestNull(t *testing.T) {
	okType(t, NULL, TNULL)

	var v Value
	var err Error

	v = NULL.ToStr(cx)
	ok(t, v, nil, NewStr("null"))

	v, err = NULL.Eq(cx, NULL)
	ok(t, v, err, TRUE)
	v, err = NULL.Eq(cx, TRUE)
	ok(t, v, err, FALSE)

	v, err = NULL.Cmp(cx, TRUE)
	fail(t, v, err, "NullValue")
}

func TestBool(t *testing.T) {

	s := TRUE.ToStr(cx)
	ok(t, s, nil, NewStr("true"))
	s = FALSE.ToStr(cx)
	ok(t, s, nil, NewStr("false"))

	okType(t, TRUE, TBOOL)
	okType(t, FALSE, TBOOL)

	tassert(t, TRUE.BoolVal())
	tassert(t, !FALSE.BoolVal())

	b, err := TRUE.Eq(cx, TRUE)
	ok(t, b, err, TRUE)
	b, err = FALSE.Eq(cx, FALSE)
	ok(t, b, err, TRUE)
	b, err = TRUE.Eq(cx, FALSE)
	ok(t, b, err, FALSE)
	b, err = FALSE.Eq(cx, TRUE)
	ok(t, b, err, FALSE)
	b, err = FALSE.Eq(cx, NewStr("a"))
	ok(t, b, err, FALSE)

	i, err := TRUE.Cmp(cx, FALSE)
	ok(t, i, err, ONE)
	i, err = FALSE.Cmp(cx, TRUE)
	ok(t, i, err, NEG_ONE)
	i, err = TRUE.Cmp(cx, TRUE)
	ok(t, i, err, ZERO)
	i, err = FALSE.Cmp(cx, FALSE)
	ok(t, i, err, ZERO)
	i, err = TRUE.Cmp(cx, NewInt(1))
	fail(t, i, err, "TypeMismatch: Expected Comparable Type")

	val := TRUE.Not()
	ok(t, val, nil, FALSE)
	val = FALSE.Not()
	ok(t, val, nil, TRUE)
}

func TestStr(t *testing.T) {
	a := NewStr("a")
	b := NewStr("b")

	var v Value
	var err Error

	v = a.ToStr(cx)
	ok(t, v, nil, NewStr("a"))
	v = b.ToStr(cx)
	ok(t, v, nil, NewStr("b"))

	okType(t, a, TSTR)
	v, err = a.Eq(cx, b)
	ok(t, v, err, FALSE)
	v, err = b.Eq(cx, a)
	ok(t, v, err, FALSE)
	v, err = a.Eq(cx, a)
	ok(t, v, err, TRUE)
	v, err = a.Eq(cx, NewStr("a"))
	ok(t, v, err, TRUE)

	v, err = a.Cmp(cx, NewInt(1))
	fail(t, v, err, "TypeMismatch: Expected Comparable Type")
	v, err = a.Cmp(cx, a)
	ok(t, v, err, NewInt(0))
	v, err = a.Cmp(cx, b)
	ok(t, v, err, NewInt(-1))
	v, err = b.Cmp(cx, a)
	ok(t, v, err, NewInt(1))

	ab := NewStr("ab")
	v, err = ab.Get(cx, NewInt(0))
	ok(t, v, err, a)
	v, err = ab.Get(cx, NewInt(1))
	ok(t, v, err, b)

	v, err = ab.Get(cx, NewInt(-1))
	ok(t, v, err, b)

	v, err = ab.Get(cx, NewInt(2))
	fail(t, v, err, "IndexOutOfBounds: 2")

	v = NewStr("").Len()
	ok(t, v, nil, ZERO)

	v = NewStr("a").Len()
	ok(t, v, nil, ONE)

	v = NewStr("abcde").Len()
	ok(t, v, nil, NewInt(5))

	//////////////////////////////
	// unicode

	a = NewStr("日本語")
	v = a.Len()
	ok(t, v, nil, NewInt(3))

	v, err = a.Get(cx, NewInt(2))
	ok(t, v, err, NewStr("語"))
}

func TestInt(t *testing.T) {
	a := NewInt(0)
	b := NewInt(1)

	s := a.ToStr(cx)
	ok(t, s, nil, NewStr("0"))
	s = b.ToStr(cx)
	ok(t, s, nil, NewStr("1"))

	okType(t, a, TINT)

	z, err := a.Eq(cx, b)
	ok(t, z, err, FALSE)
	z, err = b.Eq(cx, a)
	ok(t, z, err, FALSE)
	z, err = a.Eq(cx, a)
	ok(t, z, err, TRUE)
	z, err = a.Eq(cx, NewInt(0))
	ok(t, z, err, TRUE)
	z, err = a.Eq(cx, NewFloat(0.0))
	ok(t, z, err, TRUE)

	n, err := a.Cmp(cx, TRUE)
	fail(t, n, err, "TypeMismatch: Expected Comparable Type")
	n, err = a.Cmp(cx, a)
	ok(t, n, err, NewInt(0))
	n, err = a.Cmp(cx, b)
	ok(t, n, err, NewInt(-1))
	n, err = b.Cmp(cx, a)
	ok(t, n, err, NewInt(1))

	f := NewFloat(0.0)
	g := NewFloat(1.0)
	n, err = a.Cmp(cx, f)
	ok(t, n, err, NewInt(0))
	n, err = a.Cmp(cx, g)
	ok(t, n, err, NewInt(-1))
	n, err = g.Cmp(cx, a)
	ok(t, n, err, NewInt(1))

	val := a.Negate()
	ok(t, val, nil, NewInt(0))

	val = b.Negate()
	ok(t, val, nil, NewInt(-1))

	val, err = NewInt(3).Sub(NewInt(2))
	ok(t, val, err, NewInt(1))
	val, err = NewInt(3).Sub(NewFloat(2.0))
	ok(t, val, err, NewFloat(1.0))
	val, err = NewInt(3).Sub(NewStr("a"))
	fail(t, val, err, "TypeMismatch: Expected Number Type")
	val, err = NewInt(3).Sub(FALSE)
	fail(t, val, err, "TypeMismatch: Expected Number Type")
	val, err = NewInt(3).Sub(NULL)
	fail(t, val, err, "TypeMismatch: Expected Number Type")

	val, err = NewInt(3).Mul(NewInt(2))
	ok(t, val, err, NewInt(6))
	val, err = NewInt(3).Mul(NewFloat(2.0))
	ok(t, val, err, NewFloat(6.0))
	val, err = NewInt(3).Mul(NewStr("a"))
	fail(t, val, err, "TypeMismatch: Expected Number Type")
	val, err = NewInt(3).Mul(FALSE)
	fail(t, val, err, "TypeMismatch: Expected Number Type")
	val, err = NewInt(3).Mul(NULL)
	fail(t, val, err, "TypeMismatch: Expected Number Type")

	val, err = NewInt(3).Div(NewInt(2))
	ok(t, val, err, NewInt(1))
	val, err = NewInt(3).Div(NewFloat(2.0))
	ok(t, val, err, NewFloat(1.5))
	val, err = NewInt(3).Div(NewStr("a"))
	fail(t, val, err, "TypeMismatch: Expected Number Type")
	val, err = NewInt(3).Div(FALSE)
	fail(t, val, err, "TypeMismatch: Expected Number Type")
	val, err = NewInt(3).Div(NULL)
	fail(t, val, err, "TypeMismatch: Expected Number Type")

	val, err = NewInt(3).Div(NewInt(0))
	fail(t, val, err, "DivideByZero")
	val, err = NewInt(3).Div(NewFloat(0.0))
	fail(t, val, err, "DivideByZero")

	v1, err := NewInt(7).Rem(NewInt(3))
	ok(t, v1, err, NewInt(1))
	v1, err = NewInt(8).BitAnd(NewInt(41))
	ok(t, v1, err, NewInt(8&41))
	v1, err = NewInt(8).BitOr(NewInt(41))
	ok(t, v1, err, NewInt(8|41))
	v1, err = NewInt(8).BitXOr(NewInt(41))
	ok(t, v1, err, NewInt(8^41))
	v1, err = NewInt(1).LeftShift(NewInt(3))
	ok(t, v1, err, NewInt(8))
	v1, err = NewInt(8).RightShift(NewInt(3))
	ok(t, v1, err, NewInt(1))

	v1, err = NewInt(8).RightShift(NewStr("a"))
	fail(t, v1, err, "TypeMismatch: Expected 'Int'")

	v1, err = NewInt(8).RightShift(NewInt(-1))
	fail(t, v1, err, "InvalidArgument: Shift count cannot be less than zero")
	v1, err = NewInt(8).LeftShift(NewInt(-1))
	fail(t, v1, err, "InvalidArgument: Shift count cannot be less than zero")

	v1 = NewInt(0).Complement()
	ok(t, v1, nil, NewInt(-1))
}

func TestFloat(t *testing.T) {
	a := NewFloat(0.1)
	b := NewFloat(1.2)

	s := a.ToStr(cx)
	ok(t, s, nil, NewStr("0.1"))
	s = b.ToStr(cx)
	ok(t, s, nil, NewStr("1.2"))

	okType(t, a, TFLOAT)
	z, err := a.Eq(cx, b)
	ok(t, z, err, FALSE)
	z, err = b.Eq(cx, a)
	ok(t, z, err, FALSE)
	z, err = a.Eq(cx, a)
	ok(t, z, err, TRUE)
	z, err = a.Eq(cx, NewFloat(0.1))
	ok(t, z, err, TRUE)

	f := NewFloat(0.0)
	g := NewFloat(1.0)
	i := NewInt(0)
	j := NewInt(1)
	n, err := f.Cmp(cx, NewStr("f"))
	fail(t, n, err, "TypeMismatch: Expected Comparable Type")
	n, err = f.Cmp(cx, f)
	ok(t, n, err, NewInt(0))
	n, err = f.Cmp(cx, g)
	ok(t, n, err, NewInt(-1))
	n, err = g.Cmp(cx, f)
	ok(t, n, err, NewInt(1))
	n, err = f.Cmp(cx, i)
	ok(t, n, err, NewInt(0))
	n, err = f.Cmp(cx, j)
	ok(t, n, err, NewInt(-1))
	n, err = j.Cmp(cx, f)
	ok(t, n, err, NewInt(1))

	z, err = NewFloat(1.0).Eq(cx, NewInt(1))
	ok(t, z, err, TRUE)

	val := a.Negate()
	ok(t, val, nil, NewFloat(-0.1))

	val, err = NewFloat(3.3).Sub(NewInt(2))
	ok(t, val, err, NewFloat(float64(3.3)-float64(int64(2))))
	val, err = NewFloat(3.3).Sub(NewFloat(2.0))
	ok(t, val, err, NewFloat(float64(3.3)-float64(2.0)))
	val, err = NewFloat(3.3).Sub(NewStr("a"))
	fail(t, val, err, "TypeMismatch: Expected Number Type")
	val, err = NewFloat(3.3).Sub(FALSE)
	fail(t, val, err, "TypeMismatch: Expected Number Type")
	val, err = NewFloat(3.3).Sub(NULL)
	fail(t, val, err, "TypeMismatch: Expected Number Type")

	val, err = NewFloat(3.3).Mul(NewInt(2))
	ok(t, val, err, NewFloat(float64(3.3)*float64(int64(2))))
	val, err = NewFloat(3.3).Mul(NewFloat(2.0))
	ok(t, val, err, NewFloat(float64(3.3)*float64(2.0)))
	val, err = NewFloat(3.3).Mul(NewStr("a"))
	fail(t, val, err, "TypeMismatch: Expected Number Type")
	val, err = NewFloat(3.3).Mul(FALSE)
	fail(t, val, err, "TypeMismatch: Expected Number Type")
	val, err = NewFloat(3.3).Mul(NULL)
	fail(t, val, err, "TypeMismatch: Expected Number Type")

	val, err = NewFloat(3.3).Div(NewInt(2))
	ok(t, val, err, NewFloat(float64(3.3)/float64(int64(2))))
	val, err = NewFloat(3.3).Div(NewFloat(2.0))
	ok(t, val, err, NewFloat(float64(3.3)/float64(2.0)))
	val, err = NewFloat(3.3).Div(NewStr("a"))
	fail(t, val, err, "TypeMismatch: Expected Number Type")
	val, err = NewFloat(3.3).Div(FALSE)
	fail(t, val, err, "TypeMismatch: Expected Number Type")
	val, err = NewFloat(3.3).Div(NULL)
	fail(t, val, err, "TypeMismatch: Expected Number Type")

	val, err = NewFloat(3.3).Div(NewInt(0))
	fail(t, val, err, "DivideByZero")
	val, err = NewFloat(3.3).Div(NewFloat(0.0))
	fail(t, val, err, "DivideByZero")
}

func TestBasic(t *testing.T) {
	// make sure all the Basic types can be used as hashmap key
	entries := make(map[Basic]Value)
	entries[NULL] = TRUE
	entries[ZERO] = TRUE
	entries[NewFloat(0.123)] = TRUE
	entries[FALSE] = TRUE
}

func TestBasicHashCode(t *testing.T) {
	h, err := NULL.HashCode(cx)
	fail(t, h, err, "NullValue")

	h, err = TRUE.HashCode(cx)
	ok(t, h, err, NewInt(1009))

	h, err = FALSE.HashCode(cx)
	ok(t, h, err, NewInt(1013))

	h, err = NewInt(123).HashCode(cx)
	ok(t, h, err, NewInt(123))

	h, err = NewFloat(0).HashCode(cx)
	ok(t, h, err, NewInt(0))

	h, err = NewFloat(1.0).HashCode(cx)
	ok(t, h, err, NewInt(4607182418800017408))

	h, err = NewFloat(-1.23e45).HashCode(cx)
	ok(t, h, err, NewInt(-3941894481896550236))

	h, err = NewStr("").HashCode(cx)
	ok(t, h, err, NewInt(0))

	h, err = NewStr("abcdef").HashCode(cx)
	ok(t, h, err, NewInt(1928994870288439732))
}

func structFuncField(t *testing.T, stc Struct, name Str) NativeFunc {
	v, err := stc.GetField(cx, name)
	tassert(t, err == nil)
	f, ok := v.(NativeFunc)
	tassert(t, ok)
	return f
}

func structInvokeFunc(t *testing.T, stc Struct, name Str) Value {
	f := structFuncField(t, stc, name)
	v, err := f.Invoke(nil, []Value{})
	tassert(t, err == nil)

	return v
}

func structInvokeBoolFunc(t *testing.T, stc Struct, name Str) Bool {
	v := structInvokeFunc(t, stc, name)
	b, ok := v.(Bool)
	tassert(t, ok)
	return b
}

func TestStrIterator(t *testing.T) {

	var ibl Iterable = NewStr("abc")

	var itr Iterator = ibl.NewIterator(cx)
	v, err := itr.IterGet()
	fail(t, v, err, "NoSuchElement")
	s := NewStr("")
	for itr.IterNext().BoolVal() {
		v, err = itr.IterGet()
		tassert(t, err == nil)
		s = s.Concat(v.ToStr(cx))
	}
	ok(t, s, nil, NewStr("abc"))
	v, err = itr.IterGet()
	fail(t, v, err, "NoSuchElement")

	itr = ibl.NewIterator(cx)
	s = NewStr("")
	for structInvokeBoolFunc(t, itr, NewStr("nextValue")).BoolVal() {
		v := structInvokeFunc(t, itr, NewStr("getValue"))
		s = s.Concat(v.ToStr(cx))
	}
	ok(t, s, nil, NewStr("abc"))
}
