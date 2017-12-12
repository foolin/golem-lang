// Copyright 2017 The Golem Project Developers
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package core

import (
	//"fmt"
	"reflect"
	"testing"
)

func assert(t *testing.T, flag bool) {
	if !flag {
		t.Error("assertion failure")
	}
}

func ok(t *testing.T, val Value, err Error, expect Value) {

	if err != nil {
		panic("ok")
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
	assert(t, val.TypeOf() == expected)
}

func TestNull(t *testing.T) {
	okType(t, NULL, TNULL)

	var v Value
	var err Error

	v = NULL.ToStr()
	ok(t, v, nil, MakeStr("null"))

	v = NULL.Eq(NULL)
	ok(t, v, nil, TRUE)
	v = NULL.Eq(TRUE)
	ok(t, v, nil, FALSE)

	v, err = NULL.Cmp(TRUE)
	fail(t, v, err, "NullValue")

	v, err = NULL.Plus(MakeStr("a"))
	ok(t, v, err, MakeStr("nulla"))
}

func TestBool(t *testing.T) {

	s := TRUE.ToStr()
	ok(t, s, nil, MakeStr("true"))
	s = FALSE.ToStr()
	ok(t, s, nil, MakeStr("false"))

	okType(t, TRUE, TBOOL)
	okType(t, FALSE, TBOOL)

	assert(t, TRUE.BoolVal())
	assert(t, !FALSE.BoolVal())

	b := TRUE.Eq(TRUE)
	ok(t, b, nil, TRUE)
	b = FALSE.Eq(FALSE)
	ok(t, b, nil, TRUE)
	b = TRUE.Eq(FALSE)
	ok(t, b, nil, FALSE)
	b = FALSE.Eq(TRUE)
	ok(t, b, nil, FALSE)
	b = FALSE.Eq(MakeStr("a"))
	ok(t, b, nil, FALSE)

	i, err := TRUE.Cmp(FALSE)
	ok(t, i, err, ONE)
	i, err = FALSE.Cmp(TRUE)
	ok(t, i, err, NEG_ONE)
	i, err = TRUE.Cmp(TRUE)
	ok(t, i, err, ZERO)
	i, err = FALSE.Cmp(FALSE)
	ok(t, i, err, ZERO)
	i, err = TRUE.Cmp(MakeInt(1))
	fail(t, i, err, "TypeMismatch: Expected Comparable Type")

	val := TRUE.Not()
	ok(t, val, nil, FALSE)
	val = FALSE.Not()
	ok(t, val, nil, TRUE)

	v, err := TRUE.Plus(MakeInt(1))
	fail(t, v, err, "TypeMismatch: Expected Number Type")

	v, err = TRUE.Plus(MakeStr("a"))
	ok(t, v, err, MakeStr("truea"))
}

func TestStr(t *testing.T) {
	a := MakeStr("a")
	b := MakeStr("b")

	var v Value
	var err Error

	v = a.ToStr()
	ok(t, v, nil, MakeStr("a"))
	v = b.ToStr()
	ok(t, v, nil, MakeStr("b"))

	okType(t, a, TSTR)
	v = a.Eq(b)
	ok(t, v, nil, FALSE)
	v = b.Eq(a)
	ok(t, v, nil, FALSE)
	v = a.Eq(a)
	ok(t, v, nil, TRUE)
	v = a.Eq(MakeStr("a"))
	ok(t, v, nil, TRUE)

	v, err = a.Cmp(MakeInt(1))
	fail(t, v, err, "TypeMismatch: Expected Comparable Type")
	v, err = a.Cmp(a)
	ok(t, v, err, MakeInt(0))
	v, err = a.Cmp(b)
	ok(t, v, err, MakeInt(-1))
	v, err = b.Cmp(a)
	ok(t, v, err, MakeInt(1))

	v, err = a.Plus(MakeInt(1))
	ok(t, v, err, MakeStr("a1"))
	v, err = a.Plus(NULL)
	ok(t, v, err, MakeStr("anull"))

	ab := MakeStr("ab")
	v, err = ab.Get(MakeInt(0))
	ok(t, v, err, a)
	v, err = ab.Get(MakeInt(1))
	ok(t, v, err, b)

	v, err = ab.Get(MakeInt(-1))
	fail(t, v, err, "IndexOutOfBounds")

	v, err = ab.Get(MakeInt(2))
	fail(t, v, err, "IndexOutOfBounds")

	v = MakeStr("").Len()
	ok(t, v, nil, ZERO)

	v = MakeStr("a").Len()
	ok(t, v, nil, ONE)

	v = MakeStr("abcde").Len()
	ok(t, v, nil, MakeInt(5))

	//////////////////////////////
	// sliceable

	a = MakeStr("xyz")
	v, err = a.SliceFrom(ONE)
	ok(t, v, err, MakeStr("yz"))
	v, err = a.SliceTo(ONE)
	ok(t, v, err, MakeStr("x"))
	v, err = a.Slice(ZERO, ONE)
	ok(t, v, err, MakeStr("x"))
	v, err = a.Slice(ZERO, MakeInt(3))
	ok(t, v, err, MakeStr("xyz"))

	v, err = a.Slice(ZERO, ZERO)
	ok(t, v, err, MakeStr(""))

	v, err = a.Slice(MakeInt(2), ZERO)
	fail(t, nil, err, "IndexOutOfBounds")

	v, err = a.SliceFrom(MakeInt(7))
	fail(t, nil, err, "IndexOutOfBounds")
	v, err = a.SliceTo(MakeInt(7))
	fail(t, nil, err, "IndexOutOfBounds")
	v, err = a.Slice(MakeInt(7), MakeInt(7))
	fail(t, nil, err, "IndexOutOfBounds")

	//////////////////////////////
	// unicode

	a = MakeStr("日本語")
	v = a.Len()
	ok(t, v, nil, MakeInt(3))

	v, err = a.Get(MakeInt(2))
	ok(t, v, err, MakeStr("語"))
}

func TestInt(t *testing.T) {
	a := MakeInt(0)
	b := MakeInt(1)

	s := a.ToStr()
	ok(t, s, nil, MakeStr("0"))
	s = b.ToStr()
	ok(t, s, nil, MakeStr("1"))

	okType(t, a, TINT)

	z := a.Eq(b)
	ok(t, z, nil, FALSE)
	z = b.Eq(a)
	ok(t, z, nil, FALSE)
	z = a.Eq(a)
	ok(t, z, nil, TRUE)
	z = a.Eq(MakeInt(0))
	ok(t, z, nil, TRUE)
	z = a.Eq(MakeFloat(0.0))
	ok(t, z, nil, TRUE)

	n, err := a.Cmp(TRUE)
	fail(t, n, err, "TypeMismatch: Expected Comparable Type")
	n, err = a.Cmp(a)
	ok(t, n, err, MakeInt(0))
	n, err = a.Cmp(b)
	ok(t, n, err, MakeInt(-1))
	n, err = b.Cmp(a)
	ok(t, n, err, MakeInt(1))

	f := MakeFloat(0.0)
	g := MakeFloat(1.0)
	n, err = a.Cmp(f)
	ok(t, n, err, MakeInt(0))
	n, err = a.Cmp(g)
	ok(t, n, err, MakeInt(-1))
	n, err = g.Cmp(a)
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

	v1, err := MakeInt(3).Plus(MakeInt(2))
	ok(t, v1, err, MakeInt(5))
	v1, err = MakeInt(3).Plus(MakeFloat(2.0))
	ok(t, v1, err, MakeFloat(5.0))
	v1, err = MakeInt(3).Plus(MakeStr("a"))
	ok(t, v1, err, MakeStr("3a"))
	v2, err := MakeInt(3).Plus(FALSE)
	fail(t, v2, err, "TypeMismatch: Expected Number Type")
	v2, err = MakeInt(3).Plus(NULL)
	fail(t, v2, err, "TypeMismatch: Expected Number Type")

	v1, err = MakeInt(7).Rem(MakeInt(3))
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

	s := a.ToStr()
	ok(t, s, nil, MakeStr("0.1"))
	s = b.ToStr()
	ok(t, s, nil, MakeStr("1.2"))

	okType(t, a, TFLOAT)
	z := a.Eq(b)
	ok(t, z, nil, FALSE)
	z = b.Eq(a)
	ok(t, z, nil, FALSE)
	z = a.Eq(a)
	ok(t, z, nil, TRUE)
	z = a.Eq(MakeFloat(0.1))
	ok(t, z, nil, TRUE)

	f := MakeFloat(0.0)
	g := MakeFloat(1.0)
	i := MakeInt(0)
	j := MakeInt(1)
	n, err := f.Cmp(MakeStr("f"))
	fail(t, n, err, "TypeMismatch: Expected Comparable Type")
	n, err = f.Cmp(f)
	ok(t, n, err, MakeInt(0))
	n, err = f.Cmp(g)
	ok(t, n, err, MakeInt(-1))
	n, err = g.Cmp(f)
	ok(t, n, err, MakeInt(1))
	n, err = f.Cmp(i)
	ok(t, n, err, MakeInt(0))
	n, err = f.Cmp(j)
	ok(t, n, err, MakeInt(-1))
	n, err = j.Cmp(f)
	ok(t, n, err, MakeInt(1))

	z = MakeFloat(1.0).Eq(MakeInt(1))
	ok(t, z, nil, TRUE)

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

	v1, err := MakeFloat(3.3).Plus(MakeInt(2))
	ok(t, v1, err, MakeFloat(float64(3.3)+float64(int64(2))))
	v1, err = MakeFloat(3.3).Plus(MakeFloat(2.0))
	ok(t, v1, err, MakeFloat(float64(3.3)+float64(2.0)))
	v1, err = MakeFloat(3.3).Plus(MakeStr("a"))
	ok(t, v1, err, MakeStr("3.3a"))
	v1, err = MakeFloat(3.3).Plus(FALSE)
	fail(t, v1, err, "TypeMismatch: Expected Number Type")
	v1, err = MakeFloat(3.3).Plus(NULL)
	fail(t, v1, err, "TypeMismatch: Expected Number Type")
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
	h, err := NULL.HashCode()
	fail(t, h, err, "NullValue")

	h, err = TRUE.HashCode()
	ok(t, h, err, MakeInt(1009))

	h, err = FALSE.HashCode()
	ok(t, h, err, MakeInt(1013))

	h, err = MakeInt(123).HashCode()
	ok(t, h, err, MakeInt(123))

	h, err = MakeFloat(0).HashCode()
	ok(t, h, err, MakeInt(0))

	h, err = MakeFloat(1.0).HashCode()
	ok(t, h, err, MakeInt(4607182418800017408))

	h, err = MakeFloat(-1.23e45).HashCode()
	ok(t, h, err, MakeInt(-3941894481896550236))

	h, err = MakeStr("").HashCode()
	ok(t, h, err, MakeInt(0))

	h, err = MakeStr("abcdef").HashCode()
	ok(t, h, err, MakeInt(1928994870288439732))
}

func structFuncField(t *testing.T, stc Struct, name Str) NativeFunc {
	v, err := stc.GetField(name)
	assert(t, err == nil)
	f, ok := v.(NativeFunc)
	assert(t, ok)
	return f
}

func structInvokeFunc(t *testing.T, stc Struct, name Str) Value {
	f := structFuncField(t, stc, name)
	v, err := f.Invoke([]Value{})
	assert(t, err == nil)

	return v
}

func structInvokeBoolFunc(t *testing.T, stc Struct, name Str) Bool {
	v := structInvokeFunc(t, stc, name)
	b, ok := v.(Bool)
	assert(t, ok)
	return b
}

func TestStrIterator(t *testing.T) {

	var ibl Iterable = MakeStr("abc")

	var itr Iterator = ibl.NewIterator()
	v, err := itr.IterGet()
	fail(t, v, err, "NoSuchElement")
	s := MakeStr("")
	for itr.IterNext().BoolVal() {
		v, err = itr.IterGet()
		assert(t, err == nil)
		s = strcat(s, v)
	}
	ok(t, s, nil, MakeStr("abc"))
	v, err = itr.IterGet()
	fail(t, v, err, "NoSuchElement")

	itr = ibl.NewIterator()
	s = MakeStr("")
	for structInvokeBoolFunc(t, itr, MakeStr("nextValue")).BoolVal() {
		v := structInvokeFunc(t, itr, MakeStr("getValue"))
		s = strcat(s, v)
	}
	ok(t, s, nil, MakeStr("abc"))
}
