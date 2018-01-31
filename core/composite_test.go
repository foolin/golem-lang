// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package core

import (
	"reflect"
	"testing"
)

func TestChain(t *testing.T) {
	s0 := newStruct([]Field{
		NewField("a", true, ONE),
		NewField("b", true, NewInt(2))})
	s1 := newStruct([]Field{
		NewField("b", true, NewInt(3)),
		NewField("c", true, NewInt(4))})

	c := MergeStructs([]Struct{s0, s1})
	//println(c.ToStr(cx).String())
	tassert(t, len(c.FieldNames()) == 3)
	//ok(t, c.ToStr(cx), nil, NewStr("struct { b: 2, c: 4, a: 1 }"))

	v, err := c.GetField(cx, NewStr("a"))
	ok(t, v, err, ONE)
	v, err = c.GetField(cx, NewStr("b"))
	ok(t, v, err, NewInt(2))
	v, err = c.GetField(cx, NewStr("c"))
	ok(t, v, err, NewInt(4))
}

func newStruct(fields []Field) Struct {
	stc, err := NewStruct(fields, false)
	if err != nil {
		panic("invalid struct")
	}
	return stc
}

func TestStruct(t *testing.T) {
	stc := newStruct([]Field{})
	okType(t, stc, StructType)
	tassert(t, reflect.DeepEqual(stc.FieldNames(), []string{}))

	s := stc.ToStr(cx)
	ok(t, s, nil, NewStr("struct { }"))

	z, err := stc.Eq(cx, newStruct([]Field{}))
	ok(t, z, err, TRUE)
	z, err = stc.Eq(cx, newStruct([]Field{NewField("a", true, ONE)}))
	ok(t, z, err, FALSE)

	val, err := stc.GetField(cx, NewStr("a"))
	fail(t, val, err, "NoSuchField: Field 'a' not found")

	//////////////////

	stc = newStruct([]Field{NewField("a", false, ONE)})
	okType(t, stc, StructType)

	s = stc.ToStr(cx)
	ok(t, s, nil, NewStr("struct { a: 1 }"))

	z, err = stc.Eq(cx, newStruct([]Field{}))
	ok(t, z, err, FALSE)
	z, err = stc.Eq(cx, newStruct([]Field{NewField("a", true, ONE)}))
	ok(t, z, err, TRUE)

	val, err = stc.GetField(cx, NewStr("a"))
	ok(t, val, err, ONE)

	val, err = stc.GetField(cx, NewStr("b"))
	fail(t, val, err, "NoSuchField: Field 'b' not found")

	err = stc.SetField(cx, NewStr("a"), NewInt(123))
	if err != nil {
		panic("unexpected error")
	}

	val, err = stc.GetField(cx, NewStr("a"))
	ok(t, val, err, NewInt(123))

	err = stc.SetField(cx, NewStr("a"), NewInt(456))
	if err != nil {
		panic("unexpected error")
	}

	val, err = stc.GetField(cx, NewStr("a"))
	ok(t, val, err, NewInt(456))

	val, err = stc.Has(NewStr("a"))
	ok(t, val, err, TRUE)

	val, err = stc.Has(NewStr("abc"))
	ok(t, val, err, FALSE)

	val, err = stc.Has(ZERO)
	fail(t, val, err, "TypeMismatch: Expected 'Str'")

	stc, err = NewStruct([]Field{NewField("a", true, NULL)}, false)
	if err != nil {
		panic("oops")
	}
	val, err = stc.GetField(cx, NewStr("a"))
	ok(t, val, err, NULL)

	tassert(t, reflect.DeepEqual(stc.FieldNames(), []string{"a"}))

	stc, err = NewStruct([]Field{
		NewField("a", true, ONE),
		NewField("a", true, ZERO)}, false)
	fail(t, stc, err, "DuplicateField: Field 'a' is a duplicate")
}

func TestNativeProp(t *testing.T) {

	var propValue Value = ZERO

	getter := NewNativeFunc(0, 0,
		func(cx Context, values []Value) (Value, Error) {
			return propValue, nil
		})

	setter := NewNativeFunc(1, 1,
		func(cx Context, values []Value) (Value, Error) {
			propValue = values[0]
			return nil, nil
		})

	stc, err := NewStruct([]Field{NewProperty("a", getter, setter)}, false)
	tassert(t, err == nil)

	val, err := stc.GetField(cx, NewStr("a"))
	ok(t, val, err, ZERO)
	tassert(t, propValue == ZERO)

	err = stc.SetField(cx, NewStr("a"), ONE)
	tassert(t, err == nil)

	val, err = stc.GetField(cx, NewStr("a"))
	ok(t, val, err, ONE)
	tassert(t, propValue == ONE)
}

func TestList(t *testing.T) {
	ls := NewList([]Value{})
	okType(t, ls, ListType)

	var v Value
	var err Error

	v = ls.ToStr(cx)
	ok(t, v, nil, NewStr("[ ]"))

	v, err = ls.Eq(cx, NewList([]Value{}))
	ok(t, v, err, TRUE)

	v, err = ls.Eq(cx, NewList([]Value{NewStr("a")}))
	ok(t, v, err, FALSE)

	v, err = ls.Eq(cx, NULL)
	ok(t, v, err, FALSE)

	v = ls.Len()
	ok(t, v, nil, ZERO)

	ls.Add(cx, NewStr("a"))

	v, err = ls.Eq(cx, NewList([]Value{}))
	ok(t, v, err, FALSE)

	v, err = ls.Eq(cx, NewList([]Value{NewStr("a")}))
	ok(t, v, err, TRUE)

	v = ls.Len()
	ok(t, v, nil, ONE)

	v, err = ls.Get(cx, ZERO)
	ok(t, v, err, NewStr("a"))

	err = ls.Set(cx, ZERO, NewStr("b"))
	tassert(t, err == nil)

	v, err = ls.Get(cx, ZERO)
	ok(t, v, err, NewStr("b"))

	v, err = ls.Get(cx, NEG_ONE)
	ok(t, v, err, NewStr("b"))

	v, err = ls.Get(cx, ONE)
	fail(t, v, err, "IndexOutOfBounds: 1")

	err = ls.Set(cx, NEG_ONE, TRUE)
	tassert(t, err == nil)

	err = ls.Set(cx, ONE, TRUE)
	fail(t, nil, err, "IndexOutOfBounds: 1")

	v = ls.ToStr(cx)
	ok(t, v, nil, NewStr("[ true ]"))

	ls.Add(cx, NewStr("z"))

	v = ls.ToStr(cx)
	ok(t, v, nil, NewStr("[ true, z ]"))
}

func TestCompositeHashCode(t *testing.T) {
	h, err := NewDict(cx, []*HEntry{}).HashCode(cx)
	fail(t, h, err, "TypeMismatch: Expected Hashable Type")

	h, err = NewList([]Value{}).HashCode(cx)
	fail(t, h, err, "TypeMismatch: Expected Hashable Type")

	h, err = newStruct([]Field{}).HashCode(cx)
	fail(t, h, err, "TypeMismatch: Expected Hashable Type")
}

func TestDict(t *testing.T) {
	d := NewDict(cx, []*HEntry{})
	okType(t, d, DictType)

	var v Value
	var err Error

	v = d.ToStr(cx)
	ok(t, v, err, NewStr("dict { }"))

	v, err = d.Eq(cx, NewDict(cx, []*HEntry{}))
	ok(t, v, err, TRUE)

	v, err = d.Eq(cx, NULL)
	ok(t, v, err, FALSE)

	v = d.Len()
	ok(t, v, nil, ZERO)

	v, err = d.Get(cx, NewStr("a"))
	ok(t, v, err, NULL)

	err = d.Set(cx, NewStr("a"), ONE)
	tassert(t, err == nil)

	v, err = d.Get(cx, NewStr("a"))
	ok(t, v, err, ONE)

	v, err = d.Eq(cx, NewDict(cx, []*HEntry{}))
	ok(t, v, err, FALSE)

	v, err = d.Eq(cx, NewDict(cx, []*HEntry{{NewStr("a"), ONE}}))
	ok(t, v, err, TRUE)

	v = d.Len()
	ok(t, v, nil, ONE)

	v = d.ToStr(cx)
	ok(t, v, nil, NewStr("dict { a: 1 }"))

	err = d.Set(cx, NewStr("b"), NewInt(2))
	tassert(t, err == nil)

	v, err = d.Get(cx, NewStr("b"))
	ok(t, v, err, NewInt(2))

	v = d.ToStr(cx)
	ok(t, v, nil, NewStr("dict { b: 2, a: 1 }"))

	tp := NewTuple([]Value{ONE, ZERO})
	d = NewDict(cx, []*HEntry{{tp, TRUE}})

	v = d.ToStr(cx)
	ok(t, v, nil, NewStr("dict { (1, 0): true }"))

	v, err = d.Get(cx, tp)
	ok(t, v, err, TRUE)
}

func TestSet(t *testing.T) {
	s := NewSet(cx, []Value{})
	okType(t, s, SetType)

	var v Value
	var err Error

	v = s.ToStr(cx)
	ok(t, v, err, NewStr("set { }"))

	v, err = s.Eq(cx, NewSet(cx, []Value{}))
	ok(t, v, err, TRUE)

	v, err = s.Eq(cx, NewSet(cx, []Value{ONE}))
	ok(t, v, err, FALSE)

	v, err = s.Eq(cx, NULL)
	ok(t, v, err, FALSE)

	v = s.Len()
	ok(t, v, nil, ZERO)

	s = NewSet(cx, []Value{ONE})

	v = s.ToStr(cx)
	ok(t, v, err, NewStr("set { 1 }"))

	v, err = s.Eq(cx, NewSet(cx, []Value{}))
	ok(t, v, err, FALSE)

	v, err = s.Eq(cx, NewSet(cx, []Value{ONE, ONE, ONE}))
	ok(t, v, err, TRUE)

	v, err = s.Eq(cx, NULL)
	ok(t, v, err, FALSE)

	v = s.Len()
	ok(t, v, nil, ONE)

	s = NewSet(cx, []Value{ONE, ZERO, ZERO, ONE})

	v = s.ToStr(cx)
	ok(t, v, err, NewStr("set { 0, 1 }"))

	v = s.Len()
	ok(t, v, nil, NewInt(2))
}

func TestTuple(t *testing.T) {
	var v Value
	var err Error

	tp := NewTuple([]Value{ONE, ZERO})
	okType(t, tp, TupleType)

	v, err = tp.Eq(cx, NewTuple([]Value{ZERO, ZERO}))
	ok(t, v, err, FALSE)

	v, err = tp.Eq(cx, NewTuple([]Value{ONE, ZERO}))
	ok(t, v, err, TRUE)

	v, err = tp.Eq(cx, NULL)
	ok(t, v, err, FALSE)

	v, err = tp.Get(cx, ZERO)
	ok(t, v, err, ONE)

	v, err = tp.Get(cx, ONE)
	ok(t, v, err, ZERO)

	v, err = tp.Get(cx, NEG_ONE)
	println("adfasfda", v.ToStr(cx))
	ok(t, v, err, ZERO)

	v, err = tp.Get(cx, NewInt(2))
	fail(t, v, err, "IndexOutOfBounds: 2")

	v = tp.ToStr(cx)
	ok(t, v, nil, NewStr("(1, 0)"))

	v = tp.Len()
	ok(t, v, nil, NewInt(2))
}

func newRange(from int64, to int64, step int64) Range {
	r, err := NewRange(from, to, step)
	if err != nil {
		panic("invalid range")
	}
	return r
}

func TestRange(t *testing.T) {
	var v Value
	var err Error

	r := newRange(0, 5, 1)
	okType(t, r, RangeType)

	v, err = r.Eq(cx, newRange(0, 5, 2))
	ok(t, v, err, FALSE)

	v, err = r.Eq(cx, newRange(0, 5, 1))
	ok(t, v, err, TRUE)

	v, err = r.Eq(cx, NULL)
	ok(t, v, err, FALSE)

	v = r.Len()
	ok(t, v, nil, NewInt(5))

	v = newRange(0, 6, 3).Len()
	ok(t, v, nil, NewInt(2))
	v = newRange(0, 7, 3).Len()
	ok(t, v, nil, NewInt(3))
	v = newRange(0, 8, 3).Len()
	ok(t, v, nil, NewInt(3))
	v = newRange(0, 9, 3).Len()
	ok(t, v, nil, NewInt(3))

	v = newRange(0, 0, 3).Len()
	ok(t, v, nil, NewInt(0))
	v = newRange(1, 0, 1).Len()
	ok(t, v, nil, NewInt(0))

	v, err = NewRange(1, 0, 0)
	fail(t, v, err, "InvalidArgument: step cannot be 0")

	v = newRange(0, -5, -1).Len()
	ok(t, v, nil, NewInt(5))
	v = newRange(-1, -8, -3).Len()
	ok(t, v, nil, NewInt(3))

	r = newRange(0, 5, 1)
	v, err = r.Get(cx, ONE)
	ok(t, v, err, NewInt(1))

	r = newRange(3, 9, 2)
	v, err = r.Get(cx, NewInt(2))
	ok(t, v, err, NewInt(7))

	r = newRange(-9, -13, -1)
	v, err = r.Get(cx, ONE)
	ok(t, v, err, NewInt(-10))
}

func TestRangeIterator(t *testing.T) {

	var ibl Iterable = newRange(1, 5, 1)

	var itr Iterator = ibl.NewIterator(cx)
	v, err := itr.IterGet()
	fail(t, v, err, "NoSuchElement")
	var n int64 = 1
	for itr.IterNext().BoolVal() {
		v, err = itr.IterGet()
		tassert(t, err == nil)

		i, ok := v.(Int)
		tassert(t, ok)
		n *= i.IntVal()
	}
	tassert(t, n == 24)
	v, err = itr.IterGet()
	fail(t, v, err, "NoSuchElement")

	itr = ibl.NewIterator(cx)
	n = 1
	for structInvokeBoolFunc(t, itr, NewStr("nextValue")).BoolVal() {
		v := structInvokeFunc(t, itr, NewStr("getValue"))

		i, ok := v.(Int)
		tassert(t, ok)
		n *= i.IntVal()
	}
	tassert(t, n == 24)
}

func TestListIterator(t *testing.T) {

	var ibl Iterable = NewList(
		[]Value{NewInt(1), NewInt(2), NewInt(3), NewInt(4)})

	var itr Iterator = ibl.NewIterator(cx)
	v, err := itr.IterGet()
	fail(t, v, err, "NoSuchElement")
	var n int64 = 1
	for itr.IterNext().BoolVal() {
		v, err = itr.IterGet()
		tassert(t, err == nil)

		i, ok := v.(Int)
		tassert(t, ok)
		n *= i.IntVal()
	}
	tassert(t, n == 24)
	v, err = itr.IterGet()
	fail(t, v, err, "NoSuchElement")

	itr = ibl.NewIterator(cx)
	err = itr.SetField(cx, NewStr("nextValue"), NULL)
	fail(t, nil, err, "ImmutableValue")
	err = itr.SetField(cx, NewStr("getValue"), NULL)
	fail(t, nil, err, "ImmutableValue")

	n = 1
	for structInvokeBoolFunc(t, itr, NewStr("nextValue")).BoolVal() {
		v := structInvokeFunc(t, itr, NewStr("getValue"))

		i, ok := v.(Int)
		tassert(t, ok)
		n *= i.IntVal()
	}
	tassert(t, n == 24)
}

func TestDictIterator(t *testing.T) {

	var ibl Iterable = NewDict(cx,
		[]*HEntry{
			{NewStr("a"), ONE},
			{NewStr("b"), NewInt(2)},
			{NewStr("c"), NewInt(3)}})

	var itr Iterator = ibl.NewIterator(cx)
	v, err := itr.IterGet()
	fail(t, v, err, "NoSuchElement")
	s := NewStr("")
	for itr.IterNext().BoolVal() {
		v, err = itr.IterGet()
		tassert(t, err == nil)

		tp, ok := v.(Tuple)
		tassert(t, ok)
		s = s.Concat(tp.ToStr(cx))
	}
	ok(t, s, nil, NewStr("(b, 2)(a, 1)(c, 3)"))
	v, err = itr.IterGet()
	fail(t, v, err, "NoSuchElement")

	itr = ibl.NewIterator(cx)
	s = NewStr("")
	for structInvokeBoolFunc(t, itr, NewStr("nextValue")).BoolVal() {
		v := structInvokeFunc(t, itr, NewStr("getValue"))

		tp, ok := v.(Tuple)
		tassert(t, ok)
		s = s.Concat(tp.ToStr(cx))
	}
	ok(t, s, nil, NewStr("(b, 2)(a, 1)(c, 3)"))
}

func TestSetIterator(t *testing.T) {

	var ibl Iterable = NewSet(cx,
		[]Value{NewStr("a"), NewStr("b"), NewStr("c")})

	var itr Iterator = ibl.NewIterator(cx)
	v, err := itr.IterGet()
	fail(t, v, err, "NoSuchElement")
	s := NewStr("")
	for itr.IterNext().BoolVal() {
		v, err = itr.IterGet()
		tassert(t, err == nil)

		s = s.Concat(v.ToStr(cx))
	}
	ok(t, s, nil, NewStr("bac"))
	v, err = itr.IterGet()
	fail(t, v, err, "NoSuchElement")

	itr = ibl.NewIterator(cx)
	s = NewStr("")
	for structInvokeBoolFunc(t, itr, NewStr("nextValue")).BoolVal() {
		v := structInvokeFunc(t, itr, NewStr("getValue"))

		s = s.Concat(v.ToStr(cx))
	}
	ok(t, s, nil, NewStr("bac"))
}
