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
	ok(t, v, nil, MakeStr("null"))

	v, err = NULL.Eq(cx, NULL)
	ok(t, v, err, TRUE)
	v, err = NULL.Eq(cx, TRUE)
	ok(t, v, err, FALSE)

	v, err = NULL.Cmp(cx, TRUE)
	fail(t, v, err, "NullValue")
}

func TestBool(t *testing.T) {

	s := TRUE.ToStr(cx)
	ok(t, s, nil, MakeStr("true"))
	s = FALSE.ToStr(cx)
	ok(t, s, nil, MakeStr("false"))

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
	b, err = FALSE.Eq(cx, MakeStr("a"))
	ok(t, b, err, FALSE)

	i, err := TRUE.Cmp(cx, FALSE)
	ok(t, i, err, ONE)
	i, err = FALSE.Cmp(cx, TRUE)
	ok(t, i, err, NEG_ONE)
	i, err = TRUE.Cmp(cx, TRUE)
	ok(t, i, err, ZERO)
	i, err = FALSE.Cmp(cx, FALSE)
	ok(t, i, err, ZERO)
	i, err = TRUE.Cmp(cx, MakeInt(1))
	fail(t, i, err, "TypeMismatch: Expected Comparable Type")

	val := TRUE.Not()
	ok(t, val, nil, FALSE)
	val = FALSE.Not()
	ok(t, val, nil, TRUE)
}

func TestStr(t *testing.T) {
	a := MakeStr("a")
	b := MakeStr("b")

	var v Value
	var err Error

	v = a.ToStr(cx)
	ok(t, v, nil, MakeStr("a"))
	v = b.ToStr(cx)
	ok(t, v, nil, MakeStr("b"))

	okType(t, a, TSTR)
	v, err = a.Eq(cx, b)
	ok(t, v, err, FALSE)
	v, err = b.Eq(cx, a)
	ok(t, v, err, FALSE)
	v, err = a.Eq(cx, a)
	ok(t, v, err, TRUE)
	v, err = a.Eq(cx, MakeStr("a"))
	ok(t, v, err, TRUE)

	v, err = a.Cmp(cx, MakeInt(1))
	fail(t, v, err, "TypeMismatch: Expected Comparable Type")
	v, err = a.Cmp(cx, a)
	ok(t, v, err, MakeInt(0))
	v, err = a.Cmp(cx, b)
	ok(t, v, err, MakeInt(-1))
	v, err = b.Cmp(cx, a)
	ok(t, v, err, MakeInt(1))

	ab := MakeStr("ab")
	v, err = ab.Get(cx, MakeInt(0))
	ok(t, v, err, a)
	v, err = ab.Get(cx, MakeInt(1))
	ok(t, v, err, b)

	v, err = ab.Get(cx, MakeInt(-1))
	ok(t, v, err, b)

	v, err = ab.Get(cx, MakeInt(2))
	fail(t, v, err, "IndexOutOfBounds: 2")

	v = MakeStr("").Len()
	ok(t, v, nil, ZERO)

	v = MakeStr("a").Len()
	ok(t, v, nil, ONE)

	v = MakeStr("abcde").Len()
	ok(t, v, nil, MakeInt(5))

	//////////////////////////////
	// unicode

	a = MakeStr("日本語")
	v = a.Len()
	ok(t, v, nil, MakeInt(3))

	v, err = a.Get(cx, MakeInt(2))
	ok(t, v, err, MakeStr("語"))
}

func TestInt(t *testing.T) {
	a := MakeInt(0)
	b := MakeInt(1)

	s := a.ToStr(cx)
	ok(t, s, nil, MakeStr("0"))
	s = b.ToStr(cx)
	ok(t, s, nil, MakeStr("1"))

	okType(t, a, TINT)

	z, err := a.Eq(cx, b)
	ok(t, z, err, FALSE)
	z, err = b.Eq(cx, a)
	ok(t, z, err, FALSE)
	z, err = a.Eq(cx, a)
	ok(t, z, err, TRUE)
	z, err = a.Eq(cx, MakeInt(0))
	ok(t, z, err, TRUE)
	z, err = a.Eq(cx, MakeFloat(0.0))
	ok(t, z, err, TRUE)

	n, err := a.Cmp(cx, TRUE)
	fail(t, n, err, "TypeMismatch: Expected Comparable Type")
	n, err = a.Cmp(cx, a)
	ok(t, n, err, MakeInt(0))
	n, err = a.Cmp(cx, b)
	ok(t, n, err, MakeInt(-1))
	n, err = b.Cmp(cx, a)
	ok(t, n, err, MakeInt(1))

	f := MakeFloat(0.0)
	g := MakeFloat(1.0)
	n, err = a.Cmp(cx, f)
	ok(t, n, err, MakeInt(0))
	n, err = a.Cmp(cx, g)
	ok(t, n, err, MakeInt(-1))
	n, err = g.Cmp(cx, a)
	ok(t, n, err, MakeInt(1))

	val := a.Negate()
	ok(t, val, nil, MakeInt(0))

	val = b.Negate()
	ok(t, val, nil, MakeInt(-1))

	val, err = MakeInt(3).Sub(MakeInt(2))
	ok(t, val, err, MakeInt(1))
	val, err = MakeInt(3).Sub(MakeFloat(2.0))
	ok(t, val, err, MakeFloat(1.0))
	val, err = MakeInt(3).Sub(MakeStr("a"))
	fail(t, val, err, "TypeMismatch: Expected Number Type")
	val, err = MakeInt(3).Sub(FALSE)
	fail(t, val, err, "TypeMismatch: Expected Number Type")
	val, err = MakeInt(3).Sub(NULL)
	fail(t, val, err, "TypeMismatch: Expected Number Type")

	val, err = MakeInt(3).Mul(MakeInt(2))
	ok(t, val, err, MakeInt(6))
	val, err = MakeInt(3).Mul(MakeFloat(2.0))
	ok(t, val, err, MakeFloat(6.0))
	val, err = MakeInt(3).Mul(MakeStr("a"))
	fail(t, val, err, "TypeMismatch: Expected Number Type")
	val, err = MakeInt(3).Mul(FALSE)
	fail(t, val, err, "TypeMismatch: Expected Number Type")
	val, err = MakeInt(3).Mul(NULL)
	fail(t, val, err, "TypeMismatch: Expected Number Type")

	val, err = MakeInt(3).Div(MakeInt(2))
	ok(t, val, err, MakeInt(1))
	val, err = MakeInt(3).Div(MakeFloat(2.0))
	ok(t, val, err, MakeFloat(1.5))
	val, err = MakeInt(3).Div(MakeStr("a"))
	fail(t, val, err, "TypeMismatch: Expected Number Type")
	val, err = MakeInt(3).Div(FALSE)
	fail(t, val, err, "TypeMismatch: Expected Number Type")
	val, err = MakeInt(3).Div(NULL)
	fail(t, val, err, "TypeMismatch: Expected Number Type")

	val, err = MakeInt(3).Div(MakeInt(0))
	fail(t, val, err, "DivideByZero")
	val, err = MakeInt(3).Div(MakeFloat(0.0))
	fail(t, val, err, "DivideByZero")

	v1, err := MakeInt(7).Rem(MakeInt(3))
	ok(t, v1, err, MakeInt(1))
	v1, err = MakeInt(8).BitAnd(MakeInt(41))
	ok(t, v1, err, MakeInt(8&41))
	v1, err = MakeInt(8).BitOr(MakeInt(41))
	ok(t, v1, err, MakeInt(8|41))
	v1, err = MakeInt(8).BitXOr(MakeInt(41))
	ok(t, v1, err, MakeInt(8^41))
	v1, err = MakeInt(1).LeftShift(MakeInt(3))
	ok(t, v1, err, MakeInt(8))
	v1, err = MakeInt(8).RightShift(MakeInt(3))
	ok(t, v1, err, MakeInt(1))

	v1, err = MakeInt(8).RightShift(MakeStr("a"))
	fail(t, v1, err, "TypeMismatch: Expected 'Int'")

	v1, err = MakeInt(8).RightShift(MakeInt(-1))
	fail(t, v1, err, "InvalidArgument: Shift count cannot be less than zero")
	v1, err = MakeInt(8).LeftShift(MakeInt(-1))
	fail(t, v1, err, "InvalidArgument: Shift count cannot be less than zero")

	v1 = MakeInt(0).Complement()
	ok(t, v1, nil, MakeInt(-1))
}

func TestFloat(t *testing.T) {
	a := MakeFloat(0.1)
	b := MakeFloat(1.2)

	s := a.ToStr(cx)
	ok(t, s, nil, MakeStr("0.1"))
	s = b.ToStr(cx)
	ok(t, s, nil, MakeStr("1.2"))

	okType(t, a, TFLOAT)
	z, err := a.Eq(cx, b)
	ok(t, z, err, FALSE)
	z, err = b.Eq(cx, a)
	ok(t, z, err, FALSE)
	z, err = a.Eq(cx, a)
	ok(t, z, err, TRUE)
	z, err = a.Eq(cx, MakeFloat(0.1))
	ok(t, z, err, TRUE)

	f := MakeFloat(0.0)
	g := MakeFloat(1.0)
	i := MakeInt(0)
	j := MakeInt(1)
	n, err := f.Cmp(cx, MakeStr("f"))
	fail(t, n, err, "TypeMismatch: Expected Comparable Type")
	n, err = f.Cmp(cx, f)
	ok(t, n, err, MakeInt(0))
	n, err = f.Cmp(cx, g)
	ok(t, n, err, MakeInt(-1))
	n, err = g.Cmp(cx, f)
	ok(t, n, err, MakeInt(1))
	n, err = f.Cmp(cx, i)
	ok(t, n, err, MakeInt(0))
	n, err = f.Cmp(cx, j)
	ok(t, n, err, MakeInt(-1))
	n, err = j.Cmp(cx, f)
	ok(t, n, err, MakeInt(1))

	z, err = MakeFloat(1.0).Eq(cx, MakeInt(1))
	ok(t, z, err, TRUE)

	val := a.Negate()
	ok(t, val, nil, MakeFloat(-0.1))

	val, err = MakeFloat(3.3).Sub(MakeInt(2))
	ok(t, val, err, MakeFloat(float64(3.3)-float64(int64(2))))
	val, err = MakeFloat(3.3).Sub(MakeFloat(2.0))
	ok(t, val, err, MakeFloat(float64(3.3)-float64(2.0)))
	val, err = MakeFloat(3.3).Sub(MakeStr("a"))
	fail(t, val, err, "TypeMismatch: Expected Number Type")
	val, err = MakeFloat(3.3).Sub(FALSE)
	fail(t, val, err, "TypeMismatch: Expected Number Type")
	val, err = MakeFloat(3.3).Sub(NULL)
	fail(t, val, err, "TypeMismatch: Expected Number Type")

	val, err = MakeFloat(3.3).Mul(MakeInt(2))
	ok(t, val, err, MakeFloat(float64(3.3)*float64(int64(2))))
	val, err = MakeFloat(3.3).Mul(MakeFloat(2.0))
	ok(t, val, err, MakeFloat(float64(3.3)*float64(2.0)))
	val, err = MakeFloat(3.3).Mul(MakeStr("a"))
	fail(t, val, err, "TypeMismatch: Expected Number Type")
	val, err = MakeFloat(3.3).Mul(FALSE)
	fail(t, val, err, "TypeMismatch: Expected Number Type")
	val, err = MakeFloat(3.3).Mul(NULL)
	fail(t, val, err, "TypeMismatch: Expected Number Type")

	val, err = MakeFloat(3.3).Div(MakeInt(2))
	ok(t, val, err, MakeFloat(float64(3.3)/float64(int64(2))))
	val, err = MakeFloat(3.3).Div(MakeFloat(2.0))
	ok(t, val, err, MakeFloat(float64(3.3)/float64(2.0)))
	val, err = MakeFloat(3.3).Div(MakeStr("a"))
	fail(t, val, err, "TypeMismatch: Expected Number Type")
	val, err = MakeFloat(3.3).Div(FALSE)
	fail(t, val, err, "TypeMismatch: Expected Number Type")
	val, err = MakeFloat(3.3).Div(NULL)
	fail(t, val, err, "TypeMismatch: Expected Number Type")

	val, err = MakeFloat(3.3).Div(MakeInt(0))
	fail(t, val, err, "DivideByZero")
	val, err = MakeFloat(3.3).Div(MakeFloat(0.0))
	fail(t, val, err, "DivideByZero")
}

func TestBasic(t *testing.T) {
	// make sure all the Basic types can be used as hashmap key
	entries := make(map[Basic]Value)
	entries[NULL] = TRUE
	entries[ZERO] = TRUE
	entries[MakeFloat(0.123)] = TRUE
	entries[FALSE] = TRUE
}

func TestBasicHashCode(t *testing.T) {
	h, err := NULL.HashCode(cx)
	fail(t, h, err, "NullValue")

	h, err = TRUE.HashCode(cx)
	ok(t, h, err, MakeInt(1009))

	h, err = FALSE.HashCode(cx)
	ok(t, h, err, MakeInt(1013))

	h, err = MakeInt(123).HashCode(cx)
	ok(t, h, err, MakeInt(123))

	h, err = MakeFloat(0).HashCode(cx)
	ok(t, h, err, MakeInt(0))

	h, err = MakeFloat(1.0).HashCode(cx)
	ok(t, h, err, MakeInt(4607182418800017408))

	h, err = MakeFloat(-1.23e45).HashCode(cx)
	ok(t, h, err, MakeInt(-3941894481896550236))

	h, err = MakeStr("").HashCode(cx)
	ok(t, h, err, MakeInt(0))

	h, err = MakeStr("abcdef").HashCode(cx)
	ok(t, h, err, MakeInt(1928994870288439732))
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

	var ibl Iterable = MakeStr("abc")

	var itr Iterator = ibl.NewIterator(cx)
	v, err := itr.IterGet()
	fail(t, v, err, "NoSuchElement")
	s := MakeStr("")
	for itr.IterNext().BoolVal() {
		v, err = itr.IterGet()
		tassert(t, err == nil)
		s = s.Concat(v.ToStr(cx))
	}
	ok(t, s, nil, MakeStr("abc"))
	v, err = itr.IterGet()
	fail(t, v, err, "NoSuchElement")

	itr = ibl.NewIterator(cx)
	s = MakeStr("")
	for structInvokeBoolFunc(t, itr, MakeStr("nextValue")).BoolVal() {
		v := structInvokeFunc(t, itr, MakeStr("getValue"))
		s = s.Concat(v.ToStr(cx))
	}
	ok(t, s, nil, MakeStr("abc"))
}
