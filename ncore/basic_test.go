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
		panic("ok")
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

	var val Value
	var err Error

	val = Null.ToStr(nil)
	ok(t, val, nil, NewStr("null"))

	val, err = Null.Eq(nil, Null)
	ok(t, val, err, True)
	val, err = Null.Eq(nil, True)
	ok(t, val, err, False)

	val, err = Null.Cmp(nil, True)
	fail(t, val, err, "NullValue")

	_, err = Null.FieldNames()
	fail(t, nil, err, "NullValue")

	_, err = Null.HasField("a")
	fail(t, nil, err, "NullValue")

	val, err = Null.GetField("a", nil)
	fail(t, val, err, "NullValue")

	val, err = Null.InvokeField("a", nil, []Value{})
	fail(t, val, err, "NullValue")
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

	var val Value
	var err Error

	val, err = True.Eq(nil, True)
	ok(t, val, err, True)
	val, err = False.Eq(nil, False)
	ok(t, val, err, True)
	val, err = True.Eq(nil, False)
	ok(t, val, err, False)
	val, err = False.Eq(nil, True)
	ok(t, val, err, False)
	val, err = False.Eq(nil, NewStr("a"))
	ok(t, val, err, False)

	val, err = True.Cmp(nil, False)
	ok(t, val, err, One)
	val, err = False.Cmp(nil, True)
	ok(t, val, err, NegOne)
	val, err = True.Cmp(nil, True)
	ok(t, val, err, Zero)
	val, err = False.Cmp(nil, False)
	ok(t, val, err, Zero)
	val, err = True.Cmp(nil, NewInt(1))
	fail(t, val, err, "TypeMismatch: Expected Comparable Type")

	val = True.Not()
	ok(t, val, nil, False)
	val = False.Not()
	ok(t, val, nil, True)

	names, err := True.FieldNames()
	okNames(t, names, err, []string{})

	val, err = True.GetField("a", nil)
	fail(t, val, err, "NoSuchField: Field 'a' not found")

	val, err = True.InvokeField("a", nil, []Value{})
	fail(t, val, err, "NoSuchField: Field 'a' not found")
}

func TestStr(t *testing.T) {
	a := NewStr("a")
	b := NewStr("b")

	var val Value
	var err Error

	val = a.ToStr(nil)
	ok(t, val, nil, NewStr("a"))
	val = b.ToStr(nil)
	ok(t, val, nil, NewStr("b"))

	okType(t, a, StrType)
	val, err = a.Eq(nil, b)
	ok(t, val, err, False)
	val, err = b.Eq(nil, a)
	ok(t, val, err, False)
	val, err = a.Eq(nil, a)
	ok(t, val, err, True)
	val, err = a.Eq(nil, NewStr("a"))
	ok(t, val, err, True)

	val, err = a.Cmp(nil, NewInt(1))
	fail(t, val, err, "TypeMismatch: Expected Comparable Type")
	val, err = a.Cmp(nil, a)
	ok(t, val, err, NewInt(0))
	val, err = a.Cmp(nil, b)
	ok(t, val, err, NewInt(-1))
	val, err = b.Cmp(nil, a)
	ok(t, val, err, NewInt(1))

	ab := NewStr("ab")
	val, err = ab.Get(nil, NewInt(0))
	ok(t, val, err, a)
	val, err = ab.Get(nil, NewInt(1))
	ok(t, val, err, b)

	val, err = ab.Get(nil, NewInt(-1))
	ok(t, val, err, b)

	val, err = ab.Get(nil, NewInt(2))
	fail(t, val, err, "IndexOutOfBounds: 2")

	val = NewStr("").Len(nil)
	ok(t, val, nil, Zero)

	val = NewStr("a").Len(nil)
	ok(t, val, nil, One)

	val = NewStr("abcde").Len(nil)
	ok(t, val, nil, NewInt(5))

	// unicode
	a = NewStr("日本語")
	val = a.Len(nil)
	ok(t, val, nil, NewInt(3))

	val, err = a.Get(nil, NewInt(2))
	ok(t, val, err, NewStr("語"))

	//////////////////////////////

	//names, err := a.FieldNames()
	//okNames(t, names, err, []string{
	//	"contains",
	//	//"index",
	//	//"lastIndex",
	//	//"startsWith",
	//	//"endsWith",
	//	//"replace",
	//	//"split",
	//})

	//var bv bool
	//bv, err = a.HasField("a")
	//ok(t, bv, err, false)

	//bv, err = a.HasField("contains")
	//ok(t, bv, err, true)

	////bv, err = a.HasField("index")
	////ok(t, bv, err, true)

	//val, err = a.GetField("a", nil)
	//fail(t, val, err, "NoSuchField: Field 'a' not found")

	//c1, err := a.GetField("contains", nil)
	//tassert(t, err == nil)
	//c2, err := a.GetField("contains", nil)
	//tassert(t, err == nil)
	//val, err = c1.Eq(nil, c2)
	//ok(t, val, err, True)

	////i1, err := a.GetField("index", nil)
	////tassert(t, err == nil)
	////i2, err := a.GetField("index", nil)
	////tassert(t, err == nil)
	////val, err = i1.Eq(nil, i2)
	////ok(t, val, err, True)

	////val, err = c1.Eq(nil, i1)
	////ok(t, val, err, False)
	////val, err = c1.Eq(nil, i2)
	////ok(t, val, err, False)
	////val, err = c2.Eq(nil, i1)
	////ok(t, val, err, False)
	////val, err = c2.Eq(nil, i2)
	////ok(t, val, err, False)
}

func TestInt(t *testing.T) {
	a := NewInt(0)
	b := NewInt(1)

	var bv bool
	var val Value
	var err Error

	val = a.ToStr(nil)
	ok(t, val, nil, NewStr("0"))
	val = b.ToStr(nil)
	ok(t, val, nil, NewStr("1"))

	okType(t, a, IntType)

	val, err = a.Eq(nil, b)
	ok(t, val, err, False)
	val, err = b.Eq(nil, a)
	ok(t, val, err, False)
	val, err = a.Eq(nil, a)
	ok(t, val, err, True)
	val, err = a.Eq(nil, NewInt(0))
	ok(t, val, err, True)
	val, err = a.Eq(nil, NewFloat(0.0))
	ok(t, val, err, True)

	val, err = a.Cmp(nil, True)
	fail(t, val, err, "TypeMismatch: Expected Comparable Type")
	val, err = a.Cmp(nil, a)
	ok(t, val, err, NewInt(0))
	val, err = a.Cmp(nil, b)
	ok(t, val, err, NewInt(-1))
	val, err = b.Cmp(nil, a)
	ok(t, val, err, NewInt(1))

	f := NewFloat(0.0)
	g := NewFloat(1.0)
	val, err = a.Cmp(nil, f)
	ok(t, val, err, NewInt(0))
	val, err = a.Cmp(nil, g)
	ok(t, val, err, NewInt(-1))
	val, err = g.Cmp(nil, a)
	ok(t, val, err, NewInt(1))

	val = a.Negate()
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

	bv, err = One.HasField("a")
	ok(t, bv, err, false)

	val, err = One.GetField("a", nil)
	fail(t, val, err, "NoSuchField: Field 'a' not found")

	val, err = One.InvokeField("a", nil, []Value{})
	fail(t, val, err, "NoSuchField: Field 'a' not found")
}

func TestFloat(t *testing.T) {
	a := NewFloat(0.1)
	b := NewFloat(1.2)

	okType(t, a, FloatType)

	var bv bool
	var val Value
	var err Error

	val = a.ToStr(nil)
	ok(t, val, nil, NewStr("0.1"))
	val = b.ToStr(nil)
	ok(t, val, nil, NewStr("1.2"))

	val, err = a.Eq(nil, b)
	ok(t, val, err, False)
	val, err = b.Eq(nil, a)
	ok(t, val, err, False)
	val, err = a.Eq(nil, a)
	ok(t, val, err, True)
	val, err = a.Eq(nil, NewFloat(0.1))
	ok(t, val, err, True)

	f := NewFloat(0.0)
	g := NewFloat(1.0)
	i := NewInt(0)
	j := NewInt(1)
	val, err = f.Cmp(nil, NewStr("f"))
	fail(t, val, err, "TypeMismatch: Expected Comparable Type")
	val, err = f.Cmp(nil, f)
	ok(t, val, err, NewInt(0))
	val, err = f.Cmp(nil, g)
	ok(t, val, err, NewInt(-1))
	val, err = g.Cmp(nil, f)
	ok(t, val, err, NewInt(1))
	val, err = f.Cmp(nil, i)
	ok(t, val, err, NewInt(0))
	val, err = f.Cmp(nil, j)
	ok(t, val, err, NewInt(-1))
	val, err = j.Cmp(nil, f)
	ok(t, val, err, NewInt(1))

	val, err = NewFloat(1.0).Eq(nil, NewInt(1))
	ok(t, val, err, True)

	val = a.Negate()
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

	bv, err = a.HasField("a")
	ok(t, bv, err, false)

	val, err = a.GetField("a", nil)
	fail(t, val, err, "NoSuchField: Field 'a' not found")

	val, err = a.InvokeField("a", nil, []Value{})
	fail(t, val, err, "NoSuchField: Field 'a' not found")
}

func TestBasic(t *testing.T) {
	// make sure all the Basic types can be used as hashmap key
	entries := make(map[Basic]Value)
	entries[Null] = True
	entries[Zero] = True
	entries[NewFloat(0.123)] = True
	entries[False] = True
}

func TestBasicHashCode(t *testing.T) {
	h, err := Null.HashCode(nil)
	fail(t, h, err, "NullValue")

	h, err = True.HashCode(nil)
	ok(t, h, err, NewInt(1009))

	h, err = False.HashCode(nil)
	ok(t, h, err, NewInt(1013))

	h, err = NewInt(123).HashCode(nil)
	ok(t, h, err, NewInt(123))

	h, err = NewFloat(0).HashCode(nil)
	ok(t, h, err, NewInt(0))

	h, err = NewFloat(1.0).HashCode(nil)
	ok(t, h, err, NewInt(4607182418800017408))

	h, err = NewFloat(-1.23e45).HashCode(nil)
	ok(t, h, err, NewInt(-3941894481896550236))

	h, err = NewStr("").HashCode(nil)
	ok(t, h, err, NewInt(0))

	h, err = NewStr("abcdef").HashCode(nil)
	ok(t, h, err, NewInt(1928994870288439732))
}
