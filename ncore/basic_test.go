// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package ncore

import (
	//"fmt"
	"reflect"
	"sort"
	"testing"
)

func tassert(t *testing.T, flag bool) {
	if !flag {
		t.Error("assertion failure")
		panic("tassert")
	}
}

func ok(t *testing.T, val interface{}, err Error, expect interface{}) {

	if err != nil {
		t.Error(err, " != ", nil)
	}

	if !reflect.DeepEqual(val, expect) {
		t.Error(val, " != ", expect)
		panic("ok")
	}
}

func okNames(t *testing.T, val []string, err Error, expect []string) {

	sort.Slice(val, func(i, j int) bool {
		return val[i] < val[j]
	})

	sort.Slice(expect, func(i, j int) bool {
		return expect[i] < expect[j]
	})

	ok(t, val, err, expect)
}

func fail(t *testing.T, val interface{}, err Error, expect string) {

	if val != nil {
		t.Error(val, " != ", nil)
		panic("fail")
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

	v = Null.ToStr(nil)
	ok(t, v, nil, NewStr("null"))

	v, err = Null.Eq(nil, Null)
	ok(t, v, err, True)
	v, err = Null.Eq(nil, True)
	ok(t, v, err, False)

	v, err = Null.Cmp(nil, True)
	fail(t, v, err, "NullValue")

	_, err = Null.FieldNames()
	fail(t, nil, err, "NullValue")

	v, err = Null.HasField(nil, NewStr("a"))
	fail(t, v, err, "NullValue")
}

func TestBool(t *testing.T) {

	s := True.ToStr(nil)
	ok(t, s, nil, NewStr("true"))
	s = False.ToStr(nil)
	ok(t, s, nil, NewStr("false"))

	okType(t, True, BoolType)
	okType(t, False, BoolType)

	tassert(t, True.BoolVal())
	tassert(t, !False.BoolVal())

	b, err := True.Eq(nil, True)
	ok(t, b, err, True)
	b, err = False.Eq(nil, False)
	ok(t, b, err, True)
	b, err = True.Eq(nil, False)
	ok(t, b, err, False)
	b, err = False.Eq(nil, True)
	ok(t, b, err, False)
	b, err = False.Eq(nil, NewStr("a"))
	ok(t, b, err, False)

	i, err := True.Cmp(nil, False)
	ok(t, i, err, One)
	i, err = False.Cmp(nil, True)
	ok(t, i, err, NegOne)
	i, err = True.Cmp(nil, True)
	ok(t, i, err, Zero)
	i, err = False.Cmp(nil, False)
	ok(t, i, err, Zero)
	i, err = True.Cmp(nil, NewInt(1))
	fail(t, i, err, "TypeMismatch: Expected Comparable Type")

	val := True.Not()
	ok(t, val, nil, False)
	val = False.Not()
	ok(t, val, nil, True)

	names, err := True.FieldNames()
	okNames(t, names, err, []string{})

	val, err = True.HasField(nil, NewStr("a"))
	ok(t, val, err, False)
}

func TestStr(t *testing.T) {
	a := NewStr("a")
	b := NewStr("b")

	var v Value
	var err Error

	v = a.ToStr(nil)
	ok(t, v, nil, NewStr("a"))
	v = b.ToStr(nil)
	ok(t, v, nil, NewStr("b"))

	okType(t, a, StrType)
	v, err = a.Eq(nil, b)
	ok(t, v, err, False)
	v, err = b.Eq(nil, a)
	ok(t, v, err, False)
	v, err = a.Eq(nil, a)
	ok(t, v, err, True)
	v, err = a.Eq(nil, NewStr("a"))
	ok(t, v, err, True)

	v, err = a.Cmp(nil, NewInt(1))
	fail(t, v, err, "TypeMismatch: Expected Comparable Type")
	v, err = a.Cmp(nil, a)
	ok(t, v, err, NewInt(0))
	v, err = a.Cmp(nil, b)
	ok(t, v, err, NewInt(-1))
	v, err = b.Cmp(nil, a)
	ok(t, v, err, NewInt(1))

	ab := NewStr("ab")
	v, err = ab.Get(nil, NewInt(0))
	ok(t, v, err, a)
	v, err = ab.Get(nil, NewInt(1))
	ok(t, v, err, b)

	v, err = ab.Get(nil, NewInt(-1))
	ok(t, v, err, b)

	v, err = ab.Get(nil, NewInt(2))
	fail(t, v, err, "IndexOutOfBounds: 2")

	v = NewStr("").Len(nil)
	ok(t, v, nil, Zero)

	v = NewStr("a").Len(nil)
	ok(t, v, nil, One)

	v = NewStr("abcde").Len(nil)
	ok(t, v, nil, NewInt(5))

	// unicode
	a = NewStr("日本語")
	v = a.Len(nil)
	ok(t, v, nil, NewInt(3))

	v, err = a.Get(nil, NewInt(2))
	ok(t, v, err, NewStr("語"))

	//////////////////////////////

	names, err := a.FieldNames()
	okNames(t, names, err, []string{
		"contains",
		"index",
		"lastIndex",
		"startsWith",
		"endsWith",
		"replace",
		"split",
	})

	v, err = a.HasField(nil, NewStr("a"))
	ok(t, v, err, False)

	v, err = a.HasField(nil, NewStr("contains"))
	ok(t, v, err, True)
}

func TestInt(t *testing.T) {
	a := NewInt(0)
	b := NewInt(1)

	s := a.ToStr(nil)
	ok(t, s, nil, NewStr("0"))
	s = b.ToStr(nil)
	ok(t, s, nil, NewStr("1"))

	okType(t, a, IntType)

	z, err := a.Eq(nil, b)
	ok(t, z, err, False)
	z, err = b.Eq(nil, a)
	ok(t, z, err, False)
	z, err = a.Eq(nil, a)
	ok(t, z, err, True)
	z, err = a.Eq(nil, NewInt(0))
	ok(t, z, err, True)
	z, err = a.Eq(nil, NewFloat(0.0))
	ok(t, z, err, True)

	n, err := a.Cmp(nil, True)
	fail(t, n, err, "TypeMismatch: Expected Comparable Type")
	n, err = a.Cmp(nil, a)
	ok(t, n, err, NewInt(0))
	n, err = a.Cmp(nil, b)
	ok(t, n, err, NewInt(-1))
	n, err = b.Cmp(nil, a)
	ok(t, n, err, NewInt(1))

	f := NewFloat(0.0)
	g := NewFloat(1.0)
	n, err = a.Cmp(nil, f)
	ok(t, n, err, NewInt(0))
	n, err = a.Cmp(nil, g)
	ok(t, n, err, NewInt(-1))
	n, err = g.Cmp(nil, a)
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

	names, err := One.FieldNames()
	okNames(t, names, err, []string{})

	has, err := One.HasField(nil, NewStr("a"))
	ok(t, has, err, False)
}

func TestFloat(t *testing.T) {
	a := NewFloat(0.1)
	b := NewFloat(1.2)

	s := a.ToStr(nil)
	ok(t, s, nil, NewStr("0.1"))
	s = b.ToStr(nil)
	ok(t, s, nil, NewStr("1.2"))

	okType(t, a, FloatType)
	z, err := a.Eq(nil, b)
	ok(t, z, err, False)
	z, err = b.Eq(nil, a)
	ok(t, z, err, False)
	z, err = a.Eq(nil, a)
	ok(t, z, err, True)
	z, err = a.Eq(nil, NewFloat(0.1))
	ok(t, z, err, True)

	f := NewFloat(0.0)
	g := NewFloat(1.0)
	i := NewInt(0)
	j := NewInt(1)
	n, err := f.Cmp(nil, NewStr("f"))
	fail(t, n, err, "TypeMismatch: Expected Comparable Type")
	n, err = f.Cmp(nil, f)
	ok(t, n, err, NewInt(0))
	n, err = f.Cmp(nil, g)
	ok(t, n, err, NewInt(-1))
	n, err = g.Cmp(nil, f)
	ok(t, n, err, NewInt(1))
	n, err = f.Cmp(nil, i)
	ok(t, n, err, NewInt(0))
	n, err = f.Cmp(nil, j)
	ok(t, n, err, NewInt(-1))
	n, err = j.Cmp(nil, f)
	ok(t, n, err, NewInt(1))

	z, err = NewFloat(1.0).Eq(nil, NewInt(1))
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

	names, err := a.FieldNames()
	okNames(t, names, err, []string{})

	has, err := a.HasField(nil, NewStr("a"))
	ok(t, has, err, False)
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
//	h, err := Null.HashCode(nil)
//	fail(t, h, err, "NullValue")
//
//	h, err = True.HashCode(nil)
//	ok(t, h, err, NewInt(1009))
//
//	h, err = False.HashCode(nil)
//	ok(t, h, err, NewInt(1013))
//
//	h, err = NewInt(123).HashCode(nil)
//	ok(t, h, err, NewInt(123))
//
//	h, err = NewFloat(0).HashCode(nil)
//	ok(t, h, err, NewInt(0))
//
//	h, err = NewFloat(1.0).HashCode(nil)
//	ok(t, h, err, NewInt(4607182418800017408))
//
//	h, err = NewFloat(-1.23e45).HashCode(nil)
//	ok(t, h, err, NewInt(-3941894481896550236))
//
//	h, err = NewStr("").HashCode(nil)
//	ok(t, h, err, NewInt(0))
//
//	h, err = NewStr("abcdef").HashCode(nil)
//	ok(t, h, err, NewInt(1928994870288439732))
//}
