// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package ncore

import (
	//"fmt"
	"reflect"
	"testing"
)

var cx Context

func tassert(t *testing.T, flag bool) {
	if !flag {
		t.Error("assertion failure")
		panic("tassert")
	}
}

func ok(t *testing.T, val Value, err Error, expect Value) {

	if err != nil {
		t.Error(err, " != ", nil)
	}

	if !reflect.DeepEqual(val, expect) {
		t.Error(val, " != ", expect)
		panic("ok")
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
	okType(t, Null, NullType)

	var v Value
	var err Error

	v = Null.ToStr(cx)
	ok(t, v, nil, NewStr("null"))

	v, err = Null.Eq(cx, Null)
	ok(t, v, err, True)
	v, err = Null.Eq(cx, True)
	ok(t, v, err, False)

	v, err = Null.Cmp(cx, True)
	fail(t, v, err, "NullValue")
}

func TestBool(t *testing.T) {

	s := True.ToStr(cx)
	ok(t, s, nil, NewStr("true"))
	s = False.ToStr(cx)
	ok(t, s, nil, NewStr("false"))

	okType(t, True, BoolType)
	okType(t, False, BoolType)

	tassert(t, True.BoolVal())
	tassert(t, !False.BoolVal())

	b, err := True.Eq(cx, True)
	ok(t, b, err, True)
	b, err = False.Eq(cx, False)
	ok(t, b, err, True)
	b, err = True.Eq(cx, False)
	ok(t, b, err, False)
	b, err = False.Eq(cx, True)
	ok(t, b, err, False)
	b, err = False.Eq(cx, NewStr("a"))
	ok(t, b, err, False)

	i, err := True.Cmp(cx, False)
	ok(t, i, err, One)
	i, err = False.Cmp(cx, True)
	ok(t, i, err, NegOne)
	i, err = True.Cmp(cx, True)
	ok(t, i, err, Zero)
	i, err = False.Cmp(cx, False)
	ok(t, i, err, Zero)
	i, err = True.Cmp(cx, NewInt(1))
	fail(t, i, err, "TypeMismatch: Expected Comparable Type")

	val := True.Not()
	ok(t, val, nil, False)
	val = False.Not()
	ok(t, val, nil, True)
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

	okType(t, a, StrType)
	v, err = a.Eq(cx, b)
	ok(t, v, err, False)
	v, err = b.Eq(cx, a)
	ok(t, v, err, False)
	v, err = a.Eq(cx, a)
	ok(t, v, err, True)
	v, err = a.Eq(cx, NewStr("a"))
	ok(t, v, err, True)

	v, err = a.Cmp(cx, NewInt(1))
	fail(t, v, err, "TypeMismatch: Expected Comparable Type")
	v, err = a.Cmp(cx, a)
	ok(t, v, err, NewInt(0))
	v, err = a.Cmp(cx, b)
	ok(t, v, err, NewInt(-1))
	v, err = b.Cmp(cx, a)
	ok(t, v, err, NewInt(1))

	//	ab := NewStr("ab")
	//	v, err = ab.Get(cx, NewInt(0))
	//	ok(t, v, err, a)
	//	v, err = ab.Get(cx, NewInt(1))
	//	ok(t, v, err, b)
	//
	//	v, err = ab.Get(cx, NewInt(-1))
	//	ok(t, v, err, b)
	//
	//	v, err = ab.Get(cx, NewInt(2))
	//	fail(t, v, err, "IndexOutOfBounds: 2")
	//
	//	v = NewStr("").Len(nil)
	//	ok(t, v, nil, Zero)
	//
	//	v = NewStr("a").Len(nil)
	//	ok(t, v, nil, One)
	//
	//	v = NewStr("abcde").Len(nil)
	//	ok(t, v, nil, NewInt(5))
	//
	//	//////////////////////////////
	//	// unicode
	//
	//	a = NewStr("日本語")
	//	v = a.Len(nil)
	//	ok(t, v, nil, NewInt(3))
	//
	//	v, err = a.Get(cx, NewInt(2))
	//	ok(t, v, err, NewStr("語"))
}

func TestInt(t *testing.T) {
	a := NewInt(0)
	b := NewInt(1)

	s := a.ToStr(cx)
	ok(t, s, nil, NewStr("0"))
	s = b.ToStr(cx)
	ok(t, s, nil, NewStr("1"))

	okType(t, a, IntType)

	z, err := a.Eq(cx, b)
	ok(t, z, err, False)
	z, err = b.Eq(cx, a)
	ok(t, z, err, False)
	z, err = a.Eq(cx, a)
	ok(t, z, err, True)
	z, err = a.Eq(cx, NewInt(0))
	ok(t, z, err, True)
	z, err = a.Eq(cx, NewFloat(0.0))
	ok(t, z, err, True)

	n, err := a.Cmp(cx, True)
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
	val, err = NewInt(3).Sub(False)
	fail(t, val, err, "TypeMismatch: Expected Number Type")
	val, err = NewInt(3).Sub(Null)
	fail(t, val, err, "TypeMismatch: Expected Number Type")

	val, err = NewInt(3).Mul(NewInt(2))
	ok(t, val, err, NewInt(6))
	val, err = NewInt(3).Mul(NewFloat(2.0))
	ok(t, val, err, NewFloat(6.0))
	val, err = NewInt(3).Mul(NewStr("a"))
	fail(t, val, err, "TypeMismatch: Expected Number Type")
	val, err = NewInt(3).Mul(False)
	fail(t, val, err, "TypeMismatch: Expected Number Type")
	val, err = NewInt(3).Mul(Null)
	fail(t, val, err, "TypeMismatch: Expected Number Type")

	val, err = NewInt(3).Div(NewInt(2))
	ok(t, val, err, NewInt(1))
	val, err = NewInt(3).Div(NewFloat(2.0))
	ok(t, val, err, NewFloat(1.5))
	val, err = NewInt(3).Div(NewStr("a"))
	fail(t, val, err, "TypeMismatch: Expected Number Type")
	val, err = NewInt(3).Div(False)
	fail(t, val, err, "TypeMismatch: Expected Number Type")
	val, err = NewInt(3).Div(Null)
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
	fail(t, v1, err, "TypeMismatch: Expected Int")

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

	okType(t, a, FloatType)
	z, err := a.Eq(cx, b)
	ok(t, z, err, False)
	z, err = b.Eq(cx, a)
	ok(t, z, err, False)
	z, err = a.Eq(cx, a)
	ok(t, z, err, True)
	z, err = a.Eq(cx, NewFloat(0.1))
	ok(t, z, err, True)

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
	ok(t, z, err, True)

	val := a.Negate()
	ok(t, val, nil, NewFloat(-0.1))

	val, err = NewFloat(3.3).Sub(NewInt(2))
	ok(t, val, err, NewFloat(float64(3.3)-float64(int64(2))))
	val, err = NewFloat(3.3).Sub(NewFloat(2.0))
	ok(t, val, err, NewFloat(float64(3.3)-float64(2.0)))
	val, err = NewFloat(3.3).Sub(NewStr("a"))
	fail(t, val, err, "TypeMismatch: Expected Number Type")
	val, err = NewFloat(3.3).Sub(False)
	fail(t, val, err, "TypeMismatch: Expected Number Type")
	val, err = NewFloat(3.3).Sub(Null)
	fail(t, val, err, "TypeMismatch: Expected Number Type")

	val, err = NewFloat(3.3).Mul(NewInt(2))
	ok(t, val, err, NewFloat(float64(3.3)*float64(int64(2))))
	val, err = NewFloat(3.3).Mul(NewFloat(2.0))
	ok(t, val, err, NewFloat(float64(3.3)*float64(2.0)))
	val, err = NewFloat(3.3).Mul(NewStr("a"))
	fail(t, val, err, "TypeMismatch: Expected Number Type")
	val, err = NewFloat(3.3).Mul(False)
	fail(t, val, err, "TypeMismatch: Expected Number Type")
	val, err = NewFloat(3.3).Mul(Null)
	fail(t, val, err, "TypeMismatch: Expected Number Type")

	val, err = NewFloat(3.3).Div(NewInt(2))
	ok(t, val, err, NewFloat(float64(3.3)/float64(int64(2))))
	val, err = NewFloat(3.3).Div(NewFloat(2.0))
	ok(t, val, err, NewFloat(float64(3.3)/float64(2.0)))
	val, err = NewFloat(3.3).Div(NewStr("a"))
	fail(t, val, err, "TypeMismatch: Expected Number Type")
	val, err = NewFloat(3.3).Div(False)
	fail(t, val, err, "TypeMismatch: Expected Number Type")
	val, err = NewFloat(3.3).Div(Null)
	fail(t, val, err, "TypeMismatch: Expected Number Type")

	val, err = NewFloat(3.3).Div(NewInt(0))
	fail(t, val, err, "DivideByZero")
	val, err = NewFloat(3.3).Div(NewFloat(0.0))
	fail(t, val, err, "DivideByZero")
}

//func TestBasic(t *testing.T) {
//	// make sure all the Basic types can be used as hashmap key
//	entries := make(map[Basic]Value)
//	entries[Null] = True
//	entries[Zero] = True
//	entries[NewFloat(0.123)] = True
//	entries[False] = True
//}
//
//func TestBasicHashCode(t *testing.T) {
//	h, err := Null.HashCode(cx)
//	fail(t, h, err, "NullValue")
//
//	h, err = True.HashCode(cx)
//	ok(t, h, err, NewInt(1009))
//
//	h, err = False.HashCode(cx)
//	ok(t, h, err, NewInt(1013))
//
//	h, err = NewInt(123).HashCode(cx)
//	ok(t, h, err, NewInt(123))
//
//	h, err = NewFloat(0).HashCode(cx)
//	ok(t, h, err, NewInt(0))
//
//	h, err = NewFloat(1.0).HashCode(cx)
//	ok(t, h, err, NewInt(4607182418800017408))
//
//	h, err = NewFloat(-1.23e45).HashCode(cx)
//	ok(t, h, err, NewInt(-3941894481896550236))
//
//	h, err = NewStr("").HashCode(cx)
//	ok(t, h, err, NewInt(0))
//
//	h, err = NewStr("abcdef").HashCode(cx)
//	ok(t, h, err, NewInt(1928994870288439732))
//}
